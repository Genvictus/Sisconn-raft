package raft

import (
	pb "Sisconn-raft/raft/raftpc"
	"context"
	"log"
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
	log.Println(s.Server.address)
	log.Println(s.Server.leaderAddress)
	if s.Server.currentState != _Leader {
		return &pb.MessageResponse{
			Response:      NotLeaderResponse + s.Server.leaderAddress,
			LeaderAddress: s.Server.leaderAddress,
		}, nil
	}
	// maybe add a dedicated logger
	// TODO: get sender's IP to be outputted to log
	log.Println("ping received")
	return &pb.MessageResponse{
		Response:      OkResponse,
		LeaderAddress: "",
	}, nil
}

func (s *ServiceServer) Get(ctx context.Context, in *pb.KeyedRequest) (*pb.ValueResponse, error) {
	if s.Server.currentState != _Leader {
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
	if s.Server.currentState != _Leader {
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
	if s.Server.currentState != _Leader {
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
	if s.Server.currentState != _Leader {
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
	if s.Server.currentState != _Leader {
		return &pb.MessageResponse{
			Response:      NotLeaderResponse + s.Server.leaderAddress,
			LeaderAddress: s.Server.leaderAddress,
		}, nil
	}
	log.Println("append key:", in.Key, "with", in.Value)
	s.Server.log.appendLog(s.Server.currentTerm, in.Key, s.Server.log.getLatest(in.Key)+in.Value)
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
	if s.Server.currentState != _Leader {
		return &pb.LogResponse{
			LogEntries:    nil,
			LeaderAddress: s.Server.leaderAddress,
		}, nil
	}
	log.Println("request log")
	var logEntries []*pb.LogEntry
	for _, log := range s.Server.log.logEntries {
		logEntries = append(logEntries, &pb.LogEntry{Term: log.term, Key: log.key, Value: log.value})
	}
	return &pb.LogResponse{
		LogEntries:    logEntries,
		LeaderAddress: "",
	}, nil
}

func (s *ServiceServer) Commit(ctx context.Context, in *pb.CommitRequest) (*pb.MessageResponse, error) {
	if s.Server.currentState != _Leader {
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

func (s *ServiceServer) AddNode(ctx context.Context, in *pb.KeyValuedRequest) (*pb.MessageResponse, error) {
	if s.Server.address != s.Server.leaderAddress {
		return &pb.MessageResponse{
			Response:      NotLeaderResponse + s.Server.leaderAddress,
			LeaderAddress: s.Server.leaderAddress,
		}, nil
	}
	log.Println("add node")
	return &pb.MessageResponse{
		Response:      "Not Implemented",
		LeaderAddress: "",
	}, nil
}

func (s *ServiceServer) RemoveNode(ctx context.Context, in *pb.KeyedRequest) (*pb.MessageResponse, error) {
	if s.Server.currentState != _Leader {
		return &pb.MessageResponse{
			Response:      NotLeaderResponse + s.Server.leaderAddress,
			LeaderAddress: s.Server.leaderAddress,
		}, nil
	}
	log.Println("add node")
	return &pb.MessageResponse{
		Response:      "Not Implemented",
		LeaderAddress: "",
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

	followerNode := s.Server

	// fmt.Printf("RequestVote: %v\n", in.Term)
	// fmt.Printf("Caller: %v\n", callerNode.address)

	if in.Term < followerNode.currentTerm {
		return &pb.VoteResult{VoteGranted: false}, nil
	}

	if followerNode.votedFor == "" || followerNode.votedFor == in.CandidateId {

		// TODO proper log check
		followerLog := &followerNode.log

		if followerLog.lastIndex == 0 {
			followerNode.votedFor = in.CandidateId
			return &pb.VoteResult{VoteGranted: true, Term: in.Term}, nil
		}

		checkTerm := followerLog.logEntries[followerLog.lastIndex].term <= in.LastLogTerm
		checkIdx := len(followerLog.logEntries)-1 <= int(in.LastLogIndex[0])

		if checkTerm && checkIdx {
			followerNode.votedFor = in.CandidateId
			return &pb.VoteResult{VoteGranted: true, Term: in.Term}, nil
		}
	}

	return &pb.VoteResult{VoteGranted: false}, nil
}

func (s *RaftServer) AppendEntries(ctx context.Context, in *pb.AppendEntriesArg) (*pb.AppendResult, error) {
	currentTerm := s.Server.currentTerm
	res := pb.AppendResult{
		Term: currentTerm,
	}
	// 1. Reply false if term < currentTerm (§5.1)
	if in.Term < currentTerm {
		res.Success = false
		return &res, nil
	}
	// 2. Reply false if log doesn’t contain an entry at prevLogIndex
	// whose term matches prevLogTerm (§5.3)
	if s.Server.log.lastIndex < in.PrevLogIndex {
		// if index is out of bound (newer entries)
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
	s.Server.log.indexLock.RLock()
	commitIndex := s.Server.log.commitIndex
	lastIndex := s.Server.log.lastIndex
	s.Server.log.indexLock.RUnlock()
	if in.LeaderCommit > commitIndex {
		go s.Server.log.commitEntries(min(in.LeaderCommit, lastIndex))
	}

	return &res, nil
}
