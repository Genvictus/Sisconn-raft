package raft

import (
	pb "Sisconn-raft/raft/raftpc"
	"context"
	"strconv"
)

/*
	Raft client RPC handling implementation
*/

type ServiceServer struct {
	pb.UnimplementedRaftServiceServer
	Server *RaftNode
}

func (s *ServiceServer) Ping(ctx context.Context, in *pb.PingRequest) (*pb.MessageResponse, error) {
	ServerLogger.Println(s.Server.address)
	ServerLogger.Println(s.Server.leaderAddress)
	if s.Server.currentState.Load() != _Leader {
		return &pb.MessageResponse{
			Response:      NotLeaderResponse + s.Server.leaderAddress,
			LeaderAddress: s.Server.leaderAddress,
		}, nil
	}
	// maybe add a dedicated logger
	// TODO: get sender's IP to be outputted to log
	ServerLogger.Println("ping received")
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
	ServerLogger.Println("get key:", in.Key)
	log := s.Server.replications[_StateReplication]
	ServerLogger.Println(log.replicatedState)
	return &pb.ValueResponse{
		Value:         log.get(in.Key),
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
	ServerLogger.Println("set key:", in.Key, "with value:", in.Value)
	s.Server.replications[_StateReplication].appendLog(s.Server.currentTerm, in.Key, in.Value)
	// Replicate Entries to commit
	commitCtx, cancel := context.WithTimeout(context.Background(), REPLICATION_TIMEOUT)
	defer cancel()
	var response string
	if s.Server.replicateEntry(commitCtx, _StateReplication) {
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
	ServerLogger.Println("get strln for key:", in.Key)
	return &pb.ValueResponse{
		Value:         strconv.Itoa(len(s.Server.replications[_StateReplication].get(in.Key))),
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
	ServerLogger.Println("del key:", in.Key)
	log := &s.Server.replications[_StateReplication]
	val := log.get(in.Key)
	log.appendLog(s.Server.currentTerm, in.Key, _DELETE_KEY)
	// Replicate Entries to commit
	commitCtx, cancel := context.WithTimeout(context.Background(), REPLICATION_TIMEOUT)
	defer cancel()
	s.Server.replicateEntry(commitCtx, _StateReplication)
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
	ServerLogger.Println("append key:", in.Key, "with", in.Value)
	log := &s.Server.replications[_StateReplication]
	log.appendLog(s.Server.currentTerm, in.Key, log.tempReplicatedState[in.Key]+in.Value)
	// Replicate Entries to commit
	commitCtx, cancel := context.WithTimeout(context.Background(), REPLICATION_TIMEOUT)
	defer cancel()
	var response string
	if s.Server.replicateEntry(commitCtx, _StateReplication) {
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
	if s.Server.currentState.Load() != _Leader {
		return &pb.LogResponse{
			LogEntries:    nil,
			LeaderAddress: s.Server.leaderAddress,
		}, nil
	}
	ServerLogger.Println("request log")
	var logEntries []*pb.LogEntry
	for _, log := range s.Server.replications[_StateReplication].logEntries {
		logEntries = append(logEntries, &pb.LogEntry{Term: log.term, Key: log.key, Value: log.value})
	}
	return &pb.LogResponse{
		LogEntries:    logEntries,
		LeaderAddress: "",
	}, nil
}

func (s *ServiceServer) Commit(ctx context.Context, in *pb.CommitRequest) (*pb.MessageResponse, error) {
	if s.Server.currentState.Load() != _Leader {
		return &pb.MessageResponse{
			Response:      NotLeaderResponse + s.Server.leaderAddress,
			LeaderAddress: s.Server.leaderAddress,
		}, nil
	}
	ServerLogger.Println("commit")
	var entries []TransactionEntry
	for _, entry := range in.CommitEntries {
		ServerLogger.Println(entry.Type, entry.Key, entry.Value)
		entries = append(entries, TransactionEntry{command: entry.Type, key: entry.Key, value: entry.Value})
	}
	s.Server.replications[_StateReplication].appendTransaction(s.Server.currentTerm, entries)
	// Replicate Entries to commit
	commitCtx, cancel := context.WithTimeout(context.Background(), REPLICATION_TIMEOUT)
	defer cancel()
	var response string
	if s.Server.replicateEntry(commitCtx, _StateReplication) {
		response = OkResponse
	} else {
		response = PendingResponse
	}
	return &pb.MessageResponse{
		Response:      response + " (" + strconv.Itoa(len(entries)) + " commands execution)",
		LeaderAddress: "",
	}, nil
}

func (s *ServiceServer) AddNode(ctx context.Context, in *pb.KeyValuedRequest) (*pb.MessageResponse, error) {
	return nil, nil
}

func (s *ServiceServer) RemoveNode(ctx context.Context, in *pb.KeyedRequest) (*pb.MessageResponse, error) {
	return nil, nil
}

/*
	Raft internudes RPC implementation
*/

type RaftServer struct {
	pb.UnimplementedRaftServer
	Server *RaftNode
}

func (s *RaftServer) RequestVote(ctx context.Context, in *pb.RequestVoteArg) (*pb.VoteResult, error) {
	ServerLogger.Println("RequestVote")
	ServerLogger.Println("Current voted for:", s.Server.votedFor)
	defer func() {
		ServerLogger.Println("Final voted for:", s.Server.votedFor)
	}()

	result := pb.VoteResult{
		Term: s.Server.currentTerm,
	}

	if in.Term < s.Server.currentTerm {
		result.VoteGranted = false
		return &result, nil
	} else {
		s.Server.compareTerm(in.Term)
	}

	s.Server.stateChange <- _RefreshFollower
	if s.Server.votedFor == "" || s.Server.votedFor == in.CandidateId {

		grantVote := true
		for i := 0; i < _ReplicationEnd; i++ {
			log := &s.Server.replications[i]

			log.indexLock.RLock()
			lastIndex := log.lastIndex
			log.indexLock.RUnlock()

			checkIdx := lastIndex <= in.LastLogIndex[i]
			checkTerm := log.getEntries(lastIndex, lastIndex)[0].term <= in.LastLogTerm[i]

			if !(checkIdx && checkTerm) {
				grantVote = false
				break
			}
		}

		if grantVote {
			s.Server.votedFor = in.CandidateId
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
		ServerLogger.Println("Request failed: outdated term")
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

	// get log type
	log := &s.Server.replications[in.GetReplicationType()]

	// 2. Reply false if log doesn’t contain an entry at prevLogIndex
	// whose term matches prevLogTerm (§5.3)
	log.indexLock.RLock()
	lastIndex := log.lastIndex
	commitIndex := log.commitIndex
	log.indexLock.RUnlock()
	if lastIndex < in.PrevLogIndex {
		// if index is out of bound (newer entries)
		ServerLogger.Println("Request failed: ")
		res.Success = false
		return &res, nil
	}
	entryTerm := log.getEntries(in.PrevLogIndex, in.PrevLogIndex)[0].term
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
	log.replaceLog(in.PrevLogIndex+1, logEntries)

	// 5. If leaderCommit > commitIndex, set commitIndex =
	// min(leaderCommit, index of last new entry)
	if in.LeaderCommit > commitIndex {
		// should be long running, especially if this deals with persistence
		// i.e. writing to stable storage
		go log.commitEntries(min(in.LeaderCommit, lastIndex))
	}

	res.Success = true
	return &res, nil
}
