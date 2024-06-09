package raft

import (
	pb "Sisconn-raft/raft/raftpc"
	"context"
	"os"
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

	replications []keyValueReplication

	// Volatile states, already in keyValueReplication
	// commitIndex uint64
	// lastApplied uint64

	// Leader States
	followerIndex [](map[string]*followerIndex)
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

			replications: []keyValueReplication{
				// _ConfigReplicaiton
				newKeyValueReplication(),
				// _StateReplication
				newKeyValueReplication(),
			},

			followerIndex: make([]map[string]*followerIndex, _ReplicationEnd),
		},

		conn:       map[string]*grpc.ClientConn{},
		raftClient: map[string]pb.RaftClient{},

		stateChange: make(chan uint),
	}
	new.currentState.Store(_Follower)
	return &new
}

func (r *RaftNode) SetConnections(targets []string) {
	for _, address := range targets {
		// create connection
		r.replications[_ConfigReplication].appendLog(r.currentTerm, address, _NodeActive)
	}
	r.replications[_ConfigReplication].commitEntries(r.replications[_ConfigReplication].lastIndex)
	r.makeConnections()
}

func (r *RaftNode) makeConnections() {
	membership := &r.replications[_ConfigReplication]

	newConfig := membership.getUnionState()
	r.connLock.Lock()
	defer r.connLock.Unlock()
	// remove old connections
	for address := range r.raftClient {
		// check if connection is still in new config
		_, ok := newConfig[address]
		if ok {
			continue
		}
		// connection no longer in new config
		if address == r.address {
			// if this node is removed, quit
			os.Exit(69)
		}
		r.conn[address].Close()
		delete(r.conn, address)
		delete(r.raftClient, address)
	}
	// add new connections
	for address := range newConfig {
		if address != r.address {
			// if connection already exists, skip
			_, ok := r.conn[address]
			if ok {
				continue
			}

			conn, err := grpc.NewClient(
				address,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)

			// if any fails, probably bad address
			if err != nil {
				ServerLogger.Println("Add connection %s fails, error: %v\n", address, err)
				continue
			}
			// if successful, add connections
			r.conn[address] = conn
			// create grpc client
			r.raftClient[address] = pb.NewRaftClient(conn)
		}
	}
}

func (r *RaftNode) Run() {
	for {
		ServerLogger.Println("Current state:", r.currentState.Load())
		ServerLogger.Println("Current voted for:", r.votedFor)
		switch r.currentState.Load() {
		case _Leader:
			timer := time.NewTimer(HEARTBEAT_INTERVAL)
			configTimer := time.NewTimer(CONFIG_UPDATE_INTERVAL)
			select {
			case <-timer.C:
				// send heartbeat, synchronize log content at that
				go r.appendEntries(true, _StateReplication, nil)
			case <-configTimer.C:
				if !timer.Stop() {
					<-timer.C
				}
				// periodically update configuration
				go r.replicateEntry(context.Background(), _ConfigReplication)
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
			case _StepDown:
				r.currentState.Store(_Follower)
			}

		case _Follower:
			timer := time.NewTimer(randMs(ELECTION_TIMEOUT_MIN, ELECTION_TIMEOUT_MAX))
			select {
			case <-timer.C:
				// Reached random timeout, begin election
				r.requestVotes()
			case change := <-r.stateChange:
				if !timer.Stop() {
					<-timer.C
				}
				switch change {
				case _RefreshFollower:
					// received new AppendEntries, leader still alive,
					// do nothing
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
			case _StepDown:
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
	membership := &r.replications[_ConfigReplication]
	membership.logLock.RLock()
	total := len(membership.replicatedState)
	membership.logLock.RUnlock()
	quorum := total/2 + 1
	return total, quorum
}

func (r *RaftNode) replicateEntry(ctx context.Context, replicationType uint64) bool {
	// r.stateChange <- _NewAppendEntry
	wait := func() chan bool {
		committedCh := make(chan bool)
		go r.appendEntries(false, replicationType, committedCh)
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

func (r *RaftNode) appendEntries(isHeartbeat bool, replicationType uint64, committedCh chan<- bool) {
	r.connLock.RLock()
	defer r.connLock.RUnlock()
	// refresh heartbeat timer, appendEntries can only done by leaders
	r.stateChange <- _NewAppendEntry

	log := &r.replications[replicationType]

	log.indexLock.RLock()
	lastIndex := log.lastIndex
	log.indexLock.RUnlock()

	membership := r.replications[_ConfigReplication].getUnionState()
	var completedCh = make(chan bool)
	for address := range membership {
		if address != r.address {
			go func() {
				completedCh <- r.singleAppendEntries(address, replicationType, isHeartbeat)
			}()
		}
	}

	totalNodes, majority := r.countNodes()
	// if single server, just commit
	if totalNodes == 1 {
		log.commitEntries(lastIndex)
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
				log.commitEntries(lastIndex)
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

func (r *RaftNode) singleAppendEntries(address string, replicationType uint64, isHeartbeat bool) bool {
	// check for ongoing appendEntries for target node
	// if it is for heartbeat, try appendentries
	followerIndexx := r.followerIndex[replicationType]
	log := &r.replications[replicationType]

	_, ok := followerIndexx[address]
	if !ok {
		log.indexLock.RLock()
		nextIndex := log.lastIndex + 1
		log.indexLock.RUnlock()
		followerIndexx[address] = &followerIndex{
			nextIndex:  nextIndex,
			matchIndex: 0,
		}
	}
	if isHeartbeat {
		// if failed to acquire lock, there exists an appendEntries RPCs
		// ongoing for this follower, heartbeat is OK
		if !followerIndexx[address].indexLock.TryLock() {
			return true
		}
	} else {
		followerIndexx[address].indexLock.Lock()
	}
	defer followerIndexx[address].indexLock.Unlock()

	targetClient := r.raftClient[address]
	var appendSuccessful = false

	// get server's last log index
	log.indexLock.RLock()
	lastIndex := log.lastIndex
	log.indexLock.RUnlock()

	// maybe add counter to limit retries
	for !appendSuccessful {
		// get the current status of node
		// if node checks that it's no longer leader, stop
		if r.currentState.Load() != _Leader {
			break
		}

		// get follower's log (supposedly) last index and term
		var prevLogIndex = followerIndexx[address].nextIndex - 1
		var prevLogTerm = log.getEntries(prevLogIndex, prevLogIndex)[0].term
		// prepare log entries
		var logEntries []*pb.LogEntry
		if isHeartbeat {
			logEntries = make([]*pb.LogEntry, 0)
		} else {
			logEntries = r.createLogEntryArgs(prevLogIndex+1, lastIndex, replicationType)
		}

		// get leader's last commit index
		log.indexLock.RLock()
		leaderCommit := log.commitIndex
		log.indexLock.RUnlock()

		// prepare arguments
		args := pb.AppendEntriesArg{
			Term:            r.currentTerm,
			LeaderId:        r.address,
			PrevLogIndex:    prevLogIndex,
			PrevLogTerm:     prevLogTerm,
			LogEntries:      logEntries,
			LeaderCommit:    leaderCommit,
			ReplicationType: replicationType,
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
			followerIndexx[address].nextIndex--
		} else {
			appendSuccessful = true
			// update the index
			followerIndexx[address].nextIndex = lastIndex + 1
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

	lastIndex := make([]uint64, _ReplicationEnd)
	lastTerm := make([]uint64, _ReplicationEnd)
	// get last log index and term
	for i := 0; i < _ReplicationEnd; i++ {
		log := &r.replications[i]

		log.indexLock.RLock()
		lastIndex[i] = log.lastIndex
		lastTerm[i] = log.getEntries(lastIndex[i], lastIndex[i])[0].term
		log.indexLock.RUnlock()
	}

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

func (r *RaftNode) singleRequestVote(address string, lastLogIndex []uint64, lastLogTerm []uint64) bool {
	targetClient := r.raftClient[address]
	args := pb.RequestVoteArg{
		Term:         r.currentTerm,
		CandidateId:  r.address,
		LastLogIndex: lastLogIndex,
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

	membership := r.replications[_ConfigReplication].getUnionState()
	for i := 0; i < _ReplicationEnd; i++ {
		log := &r.replications[i]
		// get next log index for initialization
		log.indexLock.RLock()
		nextIndex := log.lastIndex + 1
		log.indexLock.RUnlock()

		// initialize indexes to track followers
		r.followerIndex[i] = map[string]*followerIndex{}
		// initiate indexes after each election
		for address := range membership {
			(r.followerIndex[i])[address] = &followerIndex{
				nextIndex:  nextIndex,
				matchIndex: 0,
			}
		}
	}

	// send initial heartbeat
	go r.appendEntries(true, _StateReplication, nil)
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
			r.stateChange <- _StepDown
		}
	}
	r.raftLock.Unlock()
}

// additional funcs to help with implementation
func (r *RaftNode) createLogEntryArgs(startIndex uint64, lastIndex uint64, replicationType uint64) []*pb.LogEntry {
	log := &r.replications[replicationType]

	logEntries := log.getEntries(startIndex, lastIndex)
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
