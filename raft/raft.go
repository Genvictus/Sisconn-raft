package raft

import (
	pb "Sisconn-raft/raft/raftpc"
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type raftState struct {
	// also used as server ID
	address string

	raftLock sync.RWMutex
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

	stateChange chan uint
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

		conn:       map[string]*grpc.ClientConn{},
		raftClient: map[string]pb.RaftClient{},

		stateChange: make(chan uint),
	}
}

func (r *RaftNode) AddConnections(targets []string) {
	r.membership.logLock.Lock()
	defer r.membership.logLock.Unlock()

	r.connLock.Lock()
	defer r.connLock.Unlock()
	for _, address := range targets {
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

func (r *RaftNode) RemoveConnections(targets []string) {
	r.membership.logLock.Lock()
	defer r.membership.logLock.Unlock()

	r.connLock.Lock()
	defer r.connLock.Unlock()
	for _, address := range targets {
		// delete the grpc client
		delete(r.raftClient, address)
		// close the connection, then update the log
		r.conn[address].Close()
		r.membership.appendLog(r.currentTerm, _DELETE_KEY, address)
	}
}

func (r *RaftNode) Run() {
	for {
		switch r.currentState {
		case _Leader:
			timer := time.NewTimer(time.Duration(time.Duration(HEARTBEAT_INTERVAL) * time.Millisecond))
			select {
			case <-timer.C:
				// send heartbeat, synchronize log content at that
				r.appendEntries(true)
			case change := <-r.stateChange:
				if !timer.Stop() {
					<-timer.C
				}
				// received notification, check which to do
				switch change {
				case _StepDown:
					r.raftLock.Lock()
					r.currentState = _Follower
					r.raftLock.Unlock()
				case _NewAppendEntry:
					// server received new command from client, append new entry
					r.appendEntries(false)
				}
			}

		case _Candidate:
			r.requestVote()
		case _Follower:
			timer := time.NewTimer(randMs(ELECTION_TIMEOUT_MIN, ELECTION_TIMEOUT_MAX))
			select {
			case <-timer.C:
				// Reached random timeout, begin election
				r.currentState = _Candidate
			case change := <-r.stateChange:
				if !timer.Stop() {
					<-timer.C
				}
				switch change {
				case _StepDown:
					r.raftLock.Lock()
					r.currentState = _Follower
					r.raftLock.Unlock()
				}
			}
		}
	}
}

func (r *RaftNode) minVote() int {
	r.membership.logLock.RLock()
	n := len(r.membership.replicatedState)
	r.membership.logLock.RUnlock()

	return n
}

func (r *RaftNode) appendEntries(isHeartbeat bool) {
	r.membership.stateLock.RLock()
	defer r.membership.stateLock.RUnlock()

	var successCount = atomic.Int32{}
	successCount.Store(0)
	for address, _ := range r.membership.replicatedState {
		if address != r.address {
			go r.singleAppendEntries(address, isHeartbeat, &successCount)
		}
	}
}

func (r *RaftNode) singleAppendEntries(address string, untilMatch bool, successCount *atomic.Int32) {
	targetClient := r.raftClient[address]
	var appendSuccessful = false

	for !appendSuccessful {
		// get the current status of node
		// if node checks that it's no longer leader, stop
		if r.currentState != _Follower {
			break
		}

		// get follower's log (supposedly) last index and term
		prevLogIndex, prevLogTerm := r.getFollowerIndex(address)
		// prepare log entries
		logEntries := r.createLogEntryArgs(prevLogIndex+1, r.log.lastIndex)

		// get leader's last commit index
		r.log.indexLock.RLock()
		leaderCommit := r.log.commitIndex
		r.log.indexLock.RUnlock()

		// prepare arguments
		args := pb.AppendEntriesArg{
			Term:         r.currentTerm,
			LeaderId:     r.address,
			PrevLogIndex: prevLogIndex,
			PrevLogTerm:  prevLogTerm,
			LogEntries:   logEntries,
			LeaderCommit: leaderCommit,
			LogType:      _DataLog,
		}

		// create RPC
		ctx, _ := context.WithTimeout(context.Background(), SERVER_RPC_TIMEOUT)
		result, err := targetClient.AppendEntries(ctx, &args)

		// if error, stop
		if err != nil {
			log.Println(err.Error())
			break
		}

		// follower term is newer than this term, immediately step down
		if result.Term > r.currentTerm {
			r.stateChange <- _StepDown
		}
		// check if append is successful
		if !result.Success {
			// decrement the index for next appendEntries()
			r.indexLock.Lock()
			r.nextIndex[address]--
			r.indexLock.Unlock()
		} else {
			// add to successful append entries
			successCount.Add(1)

			// increment the index
			r.indexLock.Lock()
			r.nextIndex[address]++
			r.indexLock.Unlock()
		}

		// if call is not set to retry until log matches
		if !untilMatch {
			break
		}
	}
}

func (r *RaftNode) requestVote() {
	// TODO
}

// additional funcs to help with implementation
func (r *RaftNode) getFollowerIndex(address string) (uint64, uint64) {
	r.indexLock.RLock()
	var prevLogIndex = r.nextIndex[address] - 1
	r.indexLock.RUnlock()
	var prevLogTerm = r.log.getEntries(prevLogIndex, prevLogIndex)[0].term

	return prevLogIndex, prevLogTerm
}

func (r *RaftNode) createLogEntryArgs(startIndex uint64, lastIndex uint64) []*pb.LogEntry {
	logEntries := r.log.getEntries(startIndex, lastIndex)
	var argEntries = make([]*pb.LogEntry, len(logEntries))

	for index, entry := range logEntries {
		argEntries[index] = &pb.LogEntry{
			Term:  entry.term,
			Key:   entry.key,
			Value: entry.value,
		}
	}
	return argEntries
}
