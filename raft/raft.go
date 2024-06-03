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

	connLock sync.RWMutex
	// Connection to other raft nodes, be sure to close the connection
	// when handling membership changes
	conn map[string]*grpc.ClientConn
	// Client for connection to other raft nodes
	raftClient map[string]pb.RaftClient
}

func NewNode(address string) *RaftNode {
	return &RaftNode{
		raftState: raftState{
			address:      address,
			currentState: _Follower,

			// currentTerm: 0,
			// votedFor: "",

			log: newKeyValueReplication(),

			membership: newKeyValueReplication(),
		},
	}
}

func (r *RaftNode) AddConnections(targets []string) {
	r.membership.logLock.Lock()
	defer r.membership.logLock.Unlock()

	r.connLock.Lock()
	defer r.connLock.Unlock()
	for _, address := range targets {

		r.membership.appendLog(r.currentTerm, address, _NodeActive)

		// create connection
		var conn, err = grpc.NewClient(
			address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		// if any fails, stop the process
		if err != nil {
			// TODO write log
			conn.Close()
			return
		}

		// if successful, add connections
		r.membership.appendLog(r.currentTerm, address, _NodeActive)
		r.conn[address] = conn
		// create grpc client
		r.raftClient[address] = pb.NewRaftClient(conn)
	}
}

func (r *RaftNode) RemoveConnections(targets []*transport.Address) {
	//TODO
}

func (r *RaftNode) Run() {
	// TODO
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
