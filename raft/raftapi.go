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
	// maybe add a dedicated logger
	// TODO: get sender's IP to be outputted to log
	log.Println("ping received")
	return &pb.MessageResponse{Response: "OK"}, nil
}

func (s *ServiceServer) Get(ctx context.Context, in *pb.KeyedRequest) (*pb.ValueResponse, error) {
	log.Println("get key:", in.Key)
	log.Println(s.Server.log.replicatedState)
	return &pb.ValueResponse{Value: s.Server.log.get(in.Key)}, nil
}

func (s *ServiceServer) Set(ctx context.Context, in *pb.KeyValuedRequest) (*pb.MessageResponse, error) {
	log.Println("set key:", in.Key, "with value:", in.Value)
	s.Server.log.appendLog(s.Server.currentTerm, in.Key, in.Value)
	s.Server.replicateEntry()
	return &pb.MessageResponse{Response: "OK"}, nil
}

func (s *ServiceServer) Strln(ctx context.Context, in *pb.KeyedRequest) (*pb.ValueResponse, error) {
	log.Println("get strln for key:", in.Key)
	return &pb.ValueResponse{Value: strconv.Itoa(len(s.Server.log.get(in.Key)))}, nil
}

func (s *ServiceServer) Del(ctx context.Context, in *pb.KeyedRequest) (*pb.ValueResponse, error) {
	log.Println("del key:", in.Key)
	val := s.Server.log.get(in.Key)
	s.Server.log.appendLog(s.Server.currentTerm, _DELETE_KEY, in.Key)
	s.Server.replicateEntry()
	return &pb.ValueResponse{Value: val}, nil
}

func (s *ServiceServer) Append(ctx context.Context, in *pb.KeyValuedRequest) (*pb.MessageResponse, error) {
	log.Println("append key:", in.Key, "with", in.Value)
	s.Server.log.appendLog(s.Server.currentTerm, in.Key, s.Server.log.get(in.Key)+in.Value)
	s.Server.replicateEntry()
	return &pb.MessageResponse{Response: "OK"}, nil
}

func (s *ServiceServer) ReqLog(ctx context.Context, in *pb.LogRequest) (*pb.LogResponse, error) {
	log.Println("request log")
	var logEntries []*pb.LogEntry
	for _, log := range s.Server.log.logEntries {
		logEntries = append(logEntries, &pb.LogEntry{Term: log.term, Key: log.key, Value: log.value})
	}
	return &pb.LogResponse{LogEntries: logEntries}, nil
}

func (s *ServiceServer) Commit(ctx context.Context, in *pb.CommitRequest) (*pb.MessageResponse, error) {
	log.Println("commit")
	var entries []TransactionEntry
	for _, entry := range in.CommitEntries {
		log.Println(entry.Type, entry.Key, entry.Value)
		entries = append(entries, TransactionEntry{command: entry.Type, key: entry.Key, value: entry.Value})
	}
	s.Server.log.appendTransaction(s.Server.currentTerm, entries)
	s.Server.replicateEntry()
	return &pb.MessageResponse{Response: "OK (" + strconv.Itoa(len(entries)) + " commands executed)"}, nil
}

func (s *ServiceServer) AddNode(ctx context.Context, in *pb.KeyValuedRequest) (*pb.MessageResponse, error) {
	log.Println("add node")
	return &pb.MessageResponse{Response: "Not Implemented"}, nil
}

func (s *ServiceServer) RemoveNode(ctx context.Context, in *pb.KeyedRequest) (*pb.MessageResponse, error) {
	log.Println("add node")
	return &pb.MessageResponse{Response: "Not Implemented"}, nil
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

	// TODO proper response from follower
	callerNode := s.Server

	// fmt.Printf("RequestVote: %v\n", in.Term)
	// fmt.Printf("Caller: %v\n", callerNode.address)

	if in.Term < callerNode.currentTerm {
		return &pb.VoteResult{VoteGranted: false}, nil
	}

	if callerNode.votedFor == "" || callerNode.votedFor == in.CandidateId {

		// TODO proper log check
		followerLog := &callerNode.log

		if followerLog.lastIndex == 0 {
			callerNode.votedFor = in.CandidateId
			return &pb.VoteResult{VoteGranted: true, Term: in.Term}, nil
		}

		checkTerm := followerLog.logEntries[followerLog.lastIndex].term <= in.LastLogTerm
		checkIdx := len(followerLog.logEntries)-1 <= int(in.LastLogIndex[0])

		if checkTerm && checkIdx {
			callerNode.votedFor = in.CandidateId
			return &pb.VoteResult{VoteGranted: true, Term: in.Term}, nil
		}
	}

	return &pb.VoteResult{VoteGranted: false}, nil
}
