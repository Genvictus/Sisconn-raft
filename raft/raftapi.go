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

// TODO: implement other RPCs

/*
	Raft internudes RPC implementation
*/

type RaftServer struct {
	pb.UnimplementedRaftServer
}

// TODO: implement raft protocol RPCs
