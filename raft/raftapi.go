package raft

import (
	pb "Sisconn-raft/raft/raftpc"
	"context"
	"log"
)

/*
	Raft client RPC handling implementation
*/

type ServiceServer struct {
	pb.UnimplementedRaftServiceServer
}

func (s *ServiceServer) Ping(ctx context.Context, in *pb.PingRequest) (*pb.MessageResponse, error) {
	// maybe add a dedicated logger
	// TODO: get sender's IP to be outputted to log
	log.Println("ping received")
	return &pb.MessageResponse{Response: "OK"}, nil
}

func (s *ServiceServer) Get(ctx context.Context, in *pb.KeyedRequest) (*pb.ValueResponse, error) {
	log.Println("get key:", in.Key)
	return &pb.ValueResponse{Value: "1"}, nil
}

func (s *ServiceServer) Set(ctx context.Context, in *pb.KeyValuedRequest) (*pb.MessageResponse, error) {
	log.Println("set key:", in.Key, "with value:", in.Value)
	return &pb.MessageResponse{Response: "OK"}, nil
}

func (s *ServiceServer) Strln(ctx context.Context, in *pb.KeyedRequest) (*pb.ValueResponse, error) {
	log.Println("get strln for key:", in.Key)
	return &pb.ValueResponse{Value: "1"}, nil
}

func (s *ServiceServer) Del(ctx context.Context, in *pb.KeyedRequest) (*pb.ValueResponse, error) {
	log.Println("del key:", in.Key)
	return &pb.ValueResponse{Value: "1"}, nil
}

func (s *ServiceServer) Append(ctx context.Context, in *pb.KeyValuedRequest) (*pb.MessageResponse, error) {
	log.Println("append key:", in.Key, "with", in.Value)
	return &pb.MessageResponse{Response: "OK"}, nil
}

func (s *ServiceServer) ReqLog(ctx context.Context, in *pb.LogRequest) (*pb.LogResponse, error) {
	log.Println("request log")
	var logEntries []*pb.LogEntry
	return &pb.LogResponse{LogEntries: logEntries}, nil
}

func (s *ServiceServer) AddNode(ctx context.Context, in *pb.KeyValuedRequest) (*pb.MessageResponse, error) {
	log.Println("add node")
	return &pb.MessageResponse{Response: "Not Implemented"}, nil
}

func (s *ServiceServer) RemoveNode(ctx context.Context, in *pb.KeyedRequest) (*pb.MessageResponse, error) {
	log.Println("add node")
	return &pb.MessageResponse{Response: "Not Implemented"}, nil
}

// TODO: implement other RPCs

/*
	Raft internudes RPC implementation
*/

type RaftServer struct {
	pb.UnimplementedRaftServer
}

// TODO: implement raft protocol RPCs