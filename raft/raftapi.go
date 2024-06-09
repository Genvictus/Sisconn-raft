package raft

import (
	pb "Sisconn-raft/raft/raftpc"
	"context"
	"log"
	"strconv"

	"google.golang.org/grpc/peer"
)

/*
	Raft client RPC handling implementation
*/

type ServiceServer struct {
	pb.UnimplementedRaftServiceServer
	Server *RaftNode
}

func (s *ServiceServer) Ping(ctx context.Context, in *pb.PingRequest) (*pb.MessageResponse, error) {
	log.Println(s.Server.address)
	log.Println(s.Server.leaderAddress)
	if s.Server.currentState.Load() != _Leader {
		return &pb.MessageResponse{
			Response:      NotLeaderResponse + s.Server.leaderAddress,
			LeaderAddress: s.Server.leaderAddress,
		}, nil
	}
	ServerLogger.Println("ping received")

	p, ok := peer.FromContext(ctx)
	if ok {
		ServerLogger.Println("Sender's IP:", p.Addr.String())
	}
	return &pb.MessageResponse{
		Response:      OkResponse,
		LeaderAddress: "",
	}, nil
}

func (s *ServiceServer) Get(ctx context.Context, in *pb.KeyedRequest) (*pb.ValueResponse, error) {
	if s.Server.currentState.Load() != _Leader {
		return &pb.ValueResponse{
			Value:         NotLeaderResponse + s.Server.leaderAddress,
			LeaderAddress: s.Server.leaderAddress,
		}, nil
	}
	log.Println("get key:", in.Key)
	log.Println(s.Server.log.replicatedState)
	return &pb.ValueResponse{
		Value:         s.Server.log.get(in.Key),
		LeaderAddress: "",
	}, nil
}

func (s *ServiceServer) Set(ctx context.Context, in *pb.KeyValuedRequest) (*pb.MessageResponse, error) {
	if s.Server.currentState.Load() != _Leader {
		return &pb.MessageResponse{
			Response:      NotLeaderResponse + s.Server.leaderAddress,
			LeaderAddress: s.Server.leaderAddress,
		}, nil
	}
	log.Println("set key:", in.Key, "with value:", in.Value)
	s.Server.log.appendLog(s.Server.currentTerm, in.Key, in.Value)
	// Replicate Entries to commit
	commitCtx, cancel := context.WithTimeout(context.Background(), REPLICATION_TIMEOUT)
	defer cancel()
	var response string
	if s.Server.replicateEntry(commitCtx) {
		response = OkResponse
	} else {
		response = PendingResponse
	}
	return &pb.MessageResponse{
		Response:      response,
		LeaderAddress: "",
	}, nil
}

func (s *ServiceServer) Strln(ctx context.Context, in *pb.KeyedRequest) (*pb.ValueResponse, error) {
	if s.Server.currentState.Load() != _Leader {
		return &pb.ValueResponse{
			Value:         NotLeaderResponse + s.Server.leaderAddress,
			LeaderAddress: s.Server.leaderAddress,
		}, nil
	}
	log.Println("get strln for key:", in.Key)
	return &pb.ValueResponse{
		Value:         strconv.Itoa(len(s.Server.log.get(in.Key))),
		LeaderAddress: "",
	}, nil
}

func (s *ServiceServer) Del(ctx context.Context, in *pb.KeyedRequest) (*pb.ValueResponse, error) {
	if s.Server.currentState.Load() != _Leader {
		return &pb.ValueResponse{
			Value:         NotLeaderResponse + s.Server.leaderAddress,
			LeaderAddress: s.Server.leaderAddress,
		}, nil
	}
	log.Println("del key:", in.Key)
	val := s.Server.log.get(in.Key)
	s.Server.log.appendLog(s.Server.currentTerm, in.Key, _DELETE_KEY)
	// Replicate Entries to commit
	commitCtx, cancel := context.WithTimeout(context.Background(), REPLICATION_TIMEOUT)
	defer cancel()
	s.Server.replicateEntry(commitCtx)
	return &pb.ValueResponse{
		Value:         val,
		LeaderAddress: "",
	}, nil
}

func (s *ServiceServer) Append(ctx context.Context, in *pb.KeyValuedRequest) (*pb.MessageResponse, error) {
	if s.Server.currentState.Load() != _Leader {
		return &pb.MessageResponse{
			Response:      NotLeaderResponse + s.Server.leaderAddress,
			LeaderAddress: s.Server.leaderAddress,
		}, nil
	}
	log.Println("append key:", in.Key, "with", in.Value)
	s.Server.log.appendLog(s.Server.currentTerm, in.Key, s.Server.log.tempReplicatedState[in.Key]+in.Value)
	// Replicate Entries to commit
	commitCtx, cancel := context.WithTimeout(context.Background(), REPLICATION_TIMEOUT)
	defer cancel()
	var response string
	if s.Server.replicateEntry(commitCtx) {
		response = OkResponse
	} else {
		response = PendingResponse
	}
	return &pb.MessageResponse{
		Response:      response,
		LeaderAddress: "",
	}, nil
}

func (s *ServiceServer) ReqLog(ctx context.Context, in *pb.LogRequest) (*pb.LogResponse, error) {
	var leaderAddress string
	if s.Server.currentState.Load() != _Leader {
		leaderAddress = s.Server.leaderAddress
	} else {
		leaderAddress = ""
	}
	log.Println("request log")
	var logEntries []*pb.LogEntry
	for _, log := range s.Server.log.logEntries {
		logEntries = append(logEntries, &pb.LogEntry{Term: log.term, Key: log.key, Value: log.value})
	}
	return &pb.LogResponse{
		LogEntries:    logEntries,
		LeaderAddress: leaderAddress,
	}, nil
}

func (s *ServiceServer) Commit(ctx context.Context, in *pb.CommitRequest) (*pb.MessageResponse, error) {
	if s.Server.currentState.Load() != _Leader {
		return &pb.MessageResponse{
			Response:      NotLeaderResponse + s.Server.leaderAddress,
			LeaderAddress: s.Server.leaderAddress,
		}, nil
	}
	log.Println("commit")
	var entries []TransactionEntry
	for _, entry := range in.CommitEntries {
		log.Println(entry.Type, entry.Key, entry.Value)
		entries = append(entries, TransactionEntry{command: entry.Type, key: entry.Key, value: entry.Value})
	}
	s.Server.log.appendTransaction(s.Server.currentTerm, entries)
	// Replicate Entries to commit
	commitCtx, cancel := context.WithTimeout(context.Background(), REPLICATION_TIMEOUT)
	defer cancel()
	var response string
	if s.Server.replicateEntry(commitCtx) {
		response = OkResponse
	} else {
		response = PendingResponse
	}
	return &pb.MessageResponse{
		Response:      response + " (" + strconv.Itoa(len(entries)) + " commands execution)",
		LeaderAddress: "",
	}, nil
}

func (s *ServiceServer) ReqClusterInfo(ctx context.Context, in *pb.ClusterInfoRequest) (*pb.ClusterInfoResponse, error) {
	if s.Server.currentState.Load() != _Leader {
		return &pb.ClusterInfoResponse{
			NodeInfo:      nil,
			LeaderAddress: s.Server.leaderAddress,
		}, nil
	}

	log.Println("Addresses: ")
	clusterInfoResponses := &pb.ClusterInfoResponse{}
	for key, _ := range s.Server.membership.replicatedState {
		log.Println("key" + key)
		clusterInfoResponses.NodeInfo = append(clusterInfoResponses.NodeInfo, &pb.NodeInfo{
			Address: key,
		})
	}

	clusterInfoResponses.LeaderAddress = ""
	return clusterInfoResponses, nil
}

func (s *ServiceServer) AddNode(ctx context.Context, in *pb.KeyValuedRequest) (*pb.MessageResponse, error) {
	if s.Server.currentState.Load() != _Leader {
		return &pb.MessageResponse{
			Response:      NotLeaderResponse + s.Server.leaderAddress,
			LeaderAddress: s.Server.leaderAddress,
		}, nil
	}
	log.Println("Adding node " + in.GetKey())
	s.Server.AddConnections([]string{
		in.GetKey(),
	})
	// s.Server.membership.appendLog(s.Server.currentTerm, in.GetKey(), _NodeInactive)

	response := "Internal Server Error"
	return &pb.MessageResponse{
		Response:      response,
		LeaderAddress: "",
	}, nil
}

func (s *ServiceServer) RemoveNode(ctx context.Context, in *pb.KeyedRequest) (*pb.MessageResponse, error) {
	if s.Server.currentState.Load() != _Leader {
		return &pb.MessageResponse{
			Response:      NotLeaderResponse + s.Server.leaderAddress,
			LeaderAddress: s.Server.leaderAddress,
		}, nil
	}
	log.Println("Removing node " + in.Key)
	s.Server.RemoveConnections([]string{
		in.GetKey(),
	})
	// s.Server.raftState.membership.appendLog(s.Server.currentTerm, _DELETE_KEY, in.Key)

	return &pb.MessageResponse{
		Response:      "OK",
		LeaderAddress: "",
	}, nil
}

func (s *ServiceServer) NodeState(ctx context.Context, in *pb.StateRequest) (*pb.StateResponse, error) {
	var state string
	switch s.Server.currentState.Load() {
	case _Leader:
		state = "LEADER"
	case _Follower:
		state = "FOLLOWER"
	case _Candidate:
		state = "CANDIDATE"
	}

	return &pb.StateResponse{
		State: state,
	}, nil
}

/*
	Raft internudes RPC implementation
*/

type RaftServer struct {
	pb.UnimplementedRaftServer
	Server *RaftNode
}

func (s *RaftServer) RequestVote(ctx context.Context, in *pb.RequestVoteArg) (*pb.VoteResult, error) {
	log.Println("RequestVote")
	log.Println("Current voted for:", s.Server.votedFor)
	defer func() {
		log.Println("Final voted for:", s.Server.votedFor)
	}()

	// fmt.Printf("RequestVote: %v\n", in.Term)
	// fmt.Printf("Caller: %v\n", callerNode.address)

	result := pb.VoteResult{
		Term: s.Server.currentTerm,
	}

	if in.Term < s.Server.currentTerm {
		result.VoteGranted = false
		return &result, nil
	} else {
		s.Server.compareTerm(in.Term)
	}

	if s.Server.votedFor == "" || s.Server.votedFor == in.CandidateId {

		followerLog := &s.Server.log

		if followerLog.lastIndex == 0 {
			s.Server.votedFor = in.CandidateId
			s.Server.stateChange <- _RefreshFollower
			result.VoteGranted = true
			return &result, nil
		}

		checkTerm := followerLog.logEntries[followerLog.lastIndex].term <= in.LastLogTerm
		checkIdx := followerLog.lastIndex <= in.LastLogIndex[0]

		if checkTerm && checkIdx {
			s.Server.votedFor = in.CandidateId
			s.Server.stateChange <- _RefreshFollower
			result.VoteGranted = true
			return &result, nil
		}
	}

	result.VoteGranted = false
	return &result, nil
}

func (s *RaftServer) AppendEntries(ctx context.Context, in *pb.AppendEntriesArg) (*pb.AppendResult, error) {
	currentTerm := s.Server.currentTerm
	res := pb.AppendResult{
		Term: currentTerm,
	}
	// 1. Reply false if term < currentTerm (§5.1)
	// implies requesting node is outdated
	if in.Term < currentTerm {
		log.Println("Request failed: outdated term")
		res.Success = false
		return &res, nil
	} else {
		s.Server.compareTerm(in.Term)
	}
	// if node is follower, refresh election timer
	if s.Server.currentState.Load() == _Follower {
		// request is valid, refresh follower
		s.Server.stateChange <- _RefreshFollower
		s.Server.raftLock.Lock()
		s.Server.leaderAddress = in.GetLeaderId()
		s.Server.raftLock.Unlock()
	}

	// 2. Reply false if log doesn’t contain an entry at prevLogIndex
	// whose term matches prevLogTerm (§5.3)
	s.Server.log.indexLock.RLock()
	lastIndex := s.Server.log.lastIndex
	commitIndex := s.Server.log.commitIndex
	s.Server.log.indexLock.RUnlock()
	if lastIndex < in.PrevLogIndex {
		// if index is out of bound (newer entries)
		ServerLogger.Println("Request failed: index is out of bound")
		res.Success = false
		return &res, nil
	}
	entryTerm := s.Server.log.getEntries(in.PrevLogIndex, in.PrevLogIndex)[0].term
	if entryTerm != in.PrevLogTerm {
		res.Success = false
		return &res, nil
	}

	// 3. If an existing entry conflicts with a new one (same index
	// but different terms), delete the existing entry and all that
	// follow it (§5.3)
	// 4. Append any new entries not already in the log
	logEntries := make([]keyValueReplicationEntry, len(in.LogEntries))
	for index, entry := range in.LogEntries {
		logEntries[index] = keyValueReplicationEntry{
			term:  entry.GetTerm(),
			key:   entry.GetKey(),
			value: entry.GetValue(),
		}
	}
	// start with goroutine? is this long running?
	s.Server.log.replaceLog(in.PrevLogIndex+1, logEntries)

	// 5. If leaderCommit > commitIndex, set commitIndex =
	// min(leaderCommit, index of last new entry)
	if in.LeaderCommit > commitIndex {
		// should be long running, especially if this deals with persistence
		// i.e. writing to stable storage
		go s.Server.log.commitEntries(min(in.LeaderCommit, lastIndex))
	}

	res.Success = true
	return &res, nil
}
