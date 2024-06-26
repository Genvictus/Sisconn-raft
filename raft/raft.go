package raft

import (
	pb "Sisconn-raft/raft/raftpc"
	"context"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type followerIndex struct {
	indexLock  sync.Mutex
	nextIndex  uint64
	matchIndex uint64
}

type raftState struct {
	// also used as server ID
	address string

	// current node's state (follower, candidate, leader)
	currentState atomic.Uint32

	raftLock      sync.RWMutex
	leaderAddress string
	// persistent states, at least in the paper
	currentTerm uint64
	votedFor    string

	log keyValueReplication

	// Volatile states, already in keyValueReplication
	// commitIndex uint64
	// lastApplied uint64

	// Cluster Membership
	membership keyValueReplication

	// Leader States
	followerIndex map[string]*followerIndex
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
	new := RaftNode{
		raftState: raftState{
			address:      address,
			currentState: atomic.Uint32{},

			// currentTerm: 0,
			// votedFor: "",

			log: newKeyValueReplication(),

			membership: newKeyValueReplication(),
		},

		conn:       map[string]*grpc.ClientConn{},
		raftClient: map[string]pb.RaftClient{},

		stateChange: make(chan uint),
	}
	new.currentState.Store(_Follower)
	return &new
}

func (r *RaftNode) AddConnections(targets []string) {
	r.connLock.Lock()
	defer r.connLock.Unlock()

	for _, address := range targets {
		// create connection
		if address != r.address {
			var conn, err = grpc.NewClient(
				address,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)

			// if any fails, stop the process
			if err != nil {
				fmt.Printf("Add connection %s fails, error: %v\n", address, err)
				return
			}

			// if successful, add connections
			r.conn[address] = conn
			// create grpc client
			r.raftClient[address] = pb.NewRaftClient(conn)
		}

		r.membership.appendLog(r.currentTerm, address, _NodeActive)
	}
	r.membership.commitEntries(r.membership.lastIndex)
}

func (r *RaftNode) RemoveConnections(targets []string) {
	r.connLock.Lock()
	defer r.connLock.Unlock()

	for _, address := range targets {
		// delete the grpc client
		delete(r.raftClient, address)
		// close the connection, then update the log
		r.conn[address].Close()
		delete(r.conn, address)
		r.membership.appendLog(r.currentTerm, _DELETE_KEY, address)
	}
	r.membership.commitEntries(r.membership.lastIndex)
}

func (r *RaftNode) Run() {
	for {
		ServerLogger.Println("Current state:", r.currentState.Load())
		ServerLogger.Println("Current voted for:", r.votedFor)
		switch r.currentState.Load() {
		case _Leader:
			timer := time.NewTimer(HEARTBEAT_INTERVAL)
			select {
			case <-timer.C:
				// send heartbeat, synchronize log content at that
				go r.appendEntries(true, nil)
			case change := <-r.stateChange:
				if !timer.Stop() {
					<-timer.C
				}
				// received notification, check which to do
				switch change {
				case _StepDown:
					r.currentState.Store(_Follower)
				case _NewAppendEntry:
					// server received new command from client, append entry sent
					// do nothing, let the timer refresh
				}
			}

		case _Candidate:
			change := <-r.stateChange
			switch change {
			case _CancelElection:
				r.currentState.Store(_Follower)
			}

		case _Follower:
			timer := time.NewTimer(randMs(ELECTION_TIMEOUT_MIN, ELECTION_TIMEOUT_MAX))
			select {
			case <-timer.C:
				// Reached random timeout, begin election
				go r.requestVotes()
			case change := <-r.stateChange:
				switch change {
				case _RefreshFollower:
					// received new AppendEntries, leader still alive,
					// refresh timer
					if !timer.Stop() {
						<-timer.C
					}
				}
			}
		}
	}
}

func (r *RaftNode) runTest() {
	for {
		switch r.currentState.Load() {
		case _Leader:
			change := <-r.stateChange
			// received notification, check which to do
			switch change {
			case _StepDown:
				r.currentState.Store(_Follower)
			case _NewAppendEntry:
				// server received new command from client, append entry sent
				// do nothing, let the timer refresh
			}

		case _Candidate:
			change := <-r.stateChange
			switch change {
			case _CancelElection:
				r.currentState.Store(_Follower)
			}

		case _Follower:
			change := <-r.stateChange
			switch change {
			case _RefreshFollower:
				// received new AppendEntries, leader still alive,
				// do nothing
			}
		}
	}
}

// count number of nodes in the cluster
// isQuorum returns the quorum for the nodes majority
func (r *RaftNode) countNodes() (int, int) {
	r.membership.logLock.RLock()
	total := len(r.membership.replicatedState)
	r.membership.logLock.RUnlock()
	quorum := total/2 + 1
	return total, quorum
}

func (r *RaftNode) replicateEntry(ctx context.Context) bool {
	// r.stateChange <- _NewAppendEntry
	wait := func() chan bool {
		committedCh := make(chan bool)
		go r.appendEntries(false, committedCh)
		return committedCh
	}

	// wait until committed
	select {
	case committed := <-wait():
		return committed
	case <-ctx.Done():
		return false
	}
}

func (r *RaftNode) appendEntries(isHeartbeat bool, committedCh chan<- bool) {
	// refresh heartbeat timer, appendEntries can only done by leaders
	r.stateChange <- _NewAppendEntry

	r.log.indexLock.RLock()
	lastIndex := r.log.lastIndex
	r.log.indexLock.RUnlock()

	r.membership.stateLock.RLock()
	var completedCh = make(chan bool)
	for address := range r.membership.replicatedState {
		if address != r.address {
			go func() {
				completedCh <- r.singleAppendEntries(address, isHeartbeat)
			}()
		}
	}
	r.membership.stateLock.RUnlock()

	totalNodes, majority := r.countNodes()
	// if single server, just commit
	if totalNodes == 1 {
		r.log.commitEntries(lastIndex)
		committedCh <- true
		return
	}
	// this waits at least for majority of RPC to complete replication
	// somewhat a partial barrier (for majority of goroutines)
	var signaled = false
	var successCount = 0
	// check for each append entries if successful
	for i := 0; i < totalNodes-1; i++ {
		replicated := <-completedCh
		// if append entries successful (logs replicated to other nodes)
		if replicated {
			successCount++
			ServerLogger.Printf("%d successful replication\n", successCount)
			if successCount == (majority - 1) {
				// if majority has replicated, notify replication's committed
				ServerLogger.Println("Successfully committed")
				r.log.commitEntries(lastIndex)
				committedCh <- true
				signaled = true
			}
		}
		// continue on until all nodes has completed replication
	}
	if !signaled {
		committedCh <- false
	}
}

func (r *RaftNode) singleAppendEntries(address string, isHeartbeat bool) bool {
	// check for ongoing appendEntries for target node
	// if it is for heartbeat, try appendentries
	if isHeartbeat {
		// if failed to acquire lock, there exists an appendEntries RPCs
		// ongoing for this follower, heartbeat is OK
		if !r.followerIndex[address].indexLock.TryLock() {
			return true
		}
	} else {
		r.followerIndex[address].indexLock.Lock()
	}
	defer r.followerIndex[address].indexLock.Unlock()

	targetClient := r.raftClient[address]
	var appendSuccessful = false

	// get server's last log index
	r.log.indexLock.RLock()
	lastIndex := r.log.lastIndex
	r.log.indexLock.RUnlock()

	// maybe add counter to limit retries
	for !appendSuccessful {
		// get the current status of node
		// if node checks that it's no longer leader, stop
		if r.currentState.Load() != _Leader {
			break
		}

		// get follower's log (supposedly) last index and term
		var prevLogIndex = r.followerIndex[address].nextIndex - 1
		var prevLogTerm = r.log.getEntries(prevLogIndex, prevLogIndex)[0].term
		// prepare log entries
		var logEntries []*pb.LogEntry
		if isHeartbeat {
			logEntries = make([]*pb.LogEntry, 0)
		} else {
			logEntries = r.createLogEntryArgs(prevLogIndex+1, lastIndex)
		}

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
		ctx, cancel := context.WithTimeout(context.Background(), SERVER_RPC_TIMEOUT)
		defer cancel()
		result, err := targetClient.AppendEntries(ctx, &args)

		// if error, stop
		if err != nil {
			ServerLogger.Println(err.Error())
			break
		}

		// follower term is newer than this term, immediately step down
		r.compareTerm(result.Term)

		// check if append is successful
		if !result.Success {
			// decrement the index for next appendEntries()
			r.followerIndex[address].nextIndex--
		} else {
			appendSuccessful = true
			// if it is heartbeat, append is actually not successful
			// DO NOT update the index because it resets the progress
			if isHeartbeat {
				break
			}
			// update the index
			r.followerIndex[address].nextIndex = lastIndex + 1
		}

		// if call is for heartbeat
		if isHeartbeat {
			break
		}
	}
	return appendSuccessful
}

func (r *RaftNode) requestVotes() {
	r.currentState.Store(_Candidate)
	// increment term
	r.raftLock.Lock()
	r.currentTerm++
	r.votedFor = r.address
	r.raftLock.Unlock()

	// get last log index and term
	r.log.indexLock.RLock()
	lastIndex := r.log.lastIndex
	lastTerm := r.log.getEntries(lastIndex, lastIndex)[0].term
	r.log.indexLock.RUnlock()

	// count votes
	voteCount := 1
	voteCh := make(chan bool)
	// send vote request to all nodes
	for address := range r.raftClient {
		if address != r.address {
			go func(addr string) {
				voteCh <- r.singleRequestVote(addr, lastIndex, lastTerm)
			}(address)
		}
	}

	totalNodes, majority := r.countNodes()
	for i := 1; i < totalNodes; i++ {
		voteGranted := <-voteCh
		if voteGranted {
			voteCount++
		}

		if voteCount >= majority {
			ServerLogger.Println("Election won by ", r.address)
			r.initiateLeader()
			return
		}
	}

	// case: 1 node
	if voteCount >= majority {
		ServerLogger.Println("Election won by ", r.address)
		r.initiateLeader()
		return
	}
	// majority not reached, revert to follower
	ServerLogger.Println("I lose the election with " + strconv.Itoa(voteCount) + " votes")
	r.currentState.Store(_Follower)
}

func (r *RaftNode) singleRequestVote(address string, lastLogIndex uint64, lastLogTerm uint64) bool {
	targetClient := r.raftClient[address]
	// map?
	lastLogIndexMap := map[uint32]uint64{0: lastLogIndex}
	args := pb.RequestVoteArg{
		Term:         r.currentTerm,
		CandidateId:  r.address,
		LastLogIndex: lastLogIndexMap,
		LastLogTerm:  lastLogTerm,
	}

	ServerLogger.Println("Requesting vote from ", address)

	// RPC
	ctx, cancel := context.WithTimeout(context.Background(), SERVER_RPC_TIMEOUT)
	result, err := targetClient.RequestVote(ctx, &args)
	defer cancel()

	if err != nil {
		ServerLogger.Println("Error requesting vote from ", address)
		ServerLogger.Println(err.Error())
		return false
	}

	ServerLogger.Println("Vote result ", result)
	// if follower term is newer immediately step down
	r.compareTerm(result.Term)

	// if vote is granted, return true
	return result.VoteGranted
}

func (r *RaftNode) initiateLeader() {
	// set current state as leader
	r.currentState.Store(_Leader)

	// get next log index for initialization
	r.log.indexLock.RLock()
	nextIndex := r.log.lastIndex + 1
	r.log.indexLock.RUnlock()

	// initialize indexes to track followers
	r.followerIndex = map[string]*followerIndex{}
	// initiate indexes after each election
	r.membership.stateLock.RLock()
	for address := range r.membership.replicatedState {
		r.followerIndex[address] = &followerIndex{
			nextIndex:  nextIndex,
			matchIndex: 0,
		}
	}
	r.membership.stateLock.RUnlock()

	// send initial heartbeat
	go r.appendEntries(true, nil)
}

// handles RPC request or response's term
func (r *RaftNode) compareTerm(receivedTerm uint64) {
	r.raftLock.Lock()
	currentTerm := r.currentTerm
	if receivedTerm > currentTerm {
		// handle new term
		r.currentTerm = receivedTerm
		r.votedFor = ""
		switch r.currentState.Load() {
		case _Leader:
			// if found a newer term, current leader is stale
			r.stateChange <- _StepDown
		case _Candidate:
			r.stateChange <- _CancelElection
		}
	}
	r.raftLock.Unlock()
}

// additional funcs to help with implementation
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
