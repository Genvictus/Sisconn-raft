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
	indexLock  sync.RWMutex
	nextIndex  map[string]uint64
	matchIndex map[string]uint64

	// Cluster Membership
	membership keyValueReplication
}

type commited struct {
	sync.Mutex
	commited chan struct{}
}

type RaftNode struct {
	raftState
	conn       grpc.ClientConn
	raftClient pb.RaftClient
}

func NewNode(address *transport.Address) (*RaftNode, error) {
	return &RaftNode{
		raftState: raftState{
			address: *address,
		},
	}, nil
}

func (r *RaftNode) AddConnections(targets []*transport.Address) {
	r.membership.logLock.RLock()
	defer r.membership.logLock.RUnlock()

	for _, address := range targets {
		r.membership.appendLog(r.currentTerm, address.String(), NodeActive)
	}
}

func (r *RaftNode) appendEntries() error {
	// TODO
}

func (r *RaftNode) heartBeat() error {
	// TODO
}
