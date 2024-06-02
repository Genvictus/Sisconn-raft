package raft

import (
	pb "Sisconn-raft/raft/raftpc"
	"Sisconn-raft/raft/transport"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type raftState struct {
	// also used as server ID
	address string

	// current node's state (follower, candidate, leader)
	currentState  uint
	leaderAddress string

	// persistent states, at least in the paper
	currentTerm uint64
	votedFor    string
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

	connLock   sync.RWMutex
	conn       map[string]*grpc.ClientConn
	raftClient map[string]pb.RaftClient
}

func NewNode(address string) *RaftNode {
	return &RaftNode{
		raftState: raftState{
			address:      address,
			currentState: Follower,

			// currentTerm: 0,
			// votedFor: "",
			log: newKeyValueReplication(),

			membership: newKeyValueReplication(),
		},
	}
}

func (r *RaftNode) AddConnections(targets []*transport.Address) {
	r.membership.logLock.RLock()
	defer r.membership.logLock.RUnlock()

	r.connLock.Lock()
	defer r.connLock.Unlock()
	for _, address := range targets {
		addressStr := address.String()

		r.membership.appendLog(r.currentTerm, addressStr, NodeActive)

		// create connection
		var conn, err = grpc.NewClient(
			addressStr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		if err != nil {
			// TODO write log
			conn.Close()
			continue
		}

		// if previous connection exists, close it first
		if r.conn[addressStr] != nil {
			r.conn[addressStr].Close()
		}
		r.conn[addressStr] = conn

		// create grpc client
		r.raftClient[addressStr] = pb.NewRaftClient(conn)
	}
}

func (r *RaftNode) minVote() int {
	r.membership.logLock.RLock()
	n := len(r.membership.replicatedState)
	r.membership.logLock.RUnlock()

	return n
}

func (r *RaftNode) appendEntries() {
	// TODO
}

func (r *RaftNode) singleAppendEntries(address string) {
	// TODO
}

func (r *RaftNode) heartBeat() {
	// TODO
}

func (r *RaftNode) requestVote() {
	// TODO
}

// util
func (r *RaftNode) GetVotedFor() string {
	return r.votedFor
}

func (r *RaftNode) SetCurrentTerm(term uint64) {
	r.currentTerm = term
}
