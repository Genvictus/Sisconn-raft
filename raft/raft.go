package raft

import (
	pb "Sisconn-raft/raft/raftpc"
	"Sisconn-raft/raft/transport"
	"sync"

	"google.golang.org/grpc"
)

type raftState struct {
	// also used as server ID
	address transport.Address

	// persistent states, at least in the paper
	currentTerm uint64
	votedFor    transport.Address
	log         keyValueReplication

	// Volatile states, already in keyValueReplication
	// commitIndex uint64
	// lastApplied uint64

	// Leader States
	nextIndex  []uint64
	matchIndex []uint64

	// Cluster Membership
	membership keyValueReplication
}

type raftLog struct {
	log keyValueReplication

	// commitIndex uint64
	// lastApplied uint64

	// nextIndex  []uint64
	// matchIndex []uint64
}

type commited struct {
	sync.Mutex
	commited chan struct{}
}

type raftNode struct {
	conn       *grpc.ClientConn
	raftClient *pb.RaftClient
	raftState
}
