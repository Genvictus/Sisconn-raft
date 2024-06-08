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

	// current node's state (follower, candidate, leader)
	currentState atomic.Uint32

	raftLock      sync.RWMutex
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
				// TODO write log
				conn.Close()
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
		log.Println("Current state:", r.currentState.Load())
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
			r.requestVotes()

		case _Follower:
			timer := time.NewTimer(randMs(ELECTION_TIMEOUT_MIN, ELECTION_TIMEOUT_MAX))
			select {
			case <-timer.C:
				// Reached random timeout, begin election
				r.currentState.Store(_Candidate)
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
			r.requestVotes()

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
func (r *RaftNode) countNodes(isQuorum bool) int {
	r.membership.logLock.RLock()
	n := len(r.membership.replicatedState)
	r.membership.logLock.RUnlock()
	if isQuorum {
		n = n/2 + 1
	}
	return n
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

	totalNodes := r.countNodes(false)
	// if single server, just commit
	if totalNodes == 1 {
		r.log.commitEntries(lastIndex)
		committedCh <- true
		return
	}
	// this waits at least for majority of RPC to complete replication
	// somewhat a partial barrier (for majority of goroutines)
	majority := totalNodes/2 + 1
	var signaled = false
	var successCount = 0
	// check for each append entries if successful
	for i := 0; i < totalNodes-1; i++ {
		replicated := <-completedCh
		// if append entries successful (logs replicated to other nodes)
		if replicated {
			successCount++
			log.Printf("%d successful replication\n", successCount)
			if successCount == (majority - 1) {
				// if majority has replicated, notify replication's committed
				log.Println("Successfully committed")
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
		prevLogIndex, prevLogTerm := r.getFollowerIndex(address)
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
			log.Println(err.Error())
			break
		}

		// follower term is newer than this term, immediately step down
		r.compareTerm(result.Term)

		// check if append is successful
		r.indexLock.Lock()
		if !result.Success {
			// decrement the index for next appendEntries()
			// TODO probably faulty
			r.nextIndex[address]--
		} else {
			appendSuccessful = true
			// update the index
			r.nextIndex[address] = lastIndex + 1
		}
		r.indexLock.Unlock()

		// if call is for heartbeat
		if isHeartbeat {
			break
		}
	}
	return appendSuccessful
}

func (r *RaftNode) requestVotes() {
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
	voteCh := make(chan bool, r.countNodes(false)-1)
	retryCh := make(chan string, r.countNodes(false)-1)
	// send vote request to all nodes
	for address := range r.raftClient {
		if address != r.address {
			go func(addr string) {
				voteCh <- r.singleRequestVote(addr, lastIndex, lastTerm, retryCh)
			}(address)
		}
	}

	electionDuration := time.NewTimer(randMs(ELECTION_TIMEOUT_MIN, ELECTION_TIMEOUT_MAX))

	totalNodes := r.countNodes(true)
	for {
		select {
		case voteGranted := <-voteCh:
			if voteGranted {
				voteCount++
			}

			if voteCount >= totalNodes {
				log.Println("Election won by ", r.address)
				r.initiateLeader()
				return
			}
		case address := <-retryCh:
			// retry if addr is active
			log.Println("Vote request failed to ", address)
			// go func(addr string) {
			// 	voteCh <- r.singleRequestVote(addr, lastIndex, lastTerm, retryCh)
			// }(address)

		case <-electionDuration.C:
			// run the election again
			if voteCount >= totalNodes {
				log.Println("Election won by ", r.address)
				r.initiateLeader()
			} else {
				log.Println("Election timeout")
				r.currentState.Store(_Follower)
			}
			return
		}

	}
}

func (r *RaftNode) singleRequestVote(address string, lastLogIndex uint64, lastLogTerm uint64, retryChannel chan<- string) bool {
	targetClient := r.raftClient[address]
	// map?
	lastLogIndexMap := map[uint32]uint64{0: lastLogIndex}
	args := pb.RequestVoteArg{
		Term:         r.currentTerm,
		CandidateId:  r.address,
		LastLogIndex: lastLogIndexMap,
		LastLogTerm:  lastLogTerm,
	}

	log.Println("Requesting vote from ", address)

	// RPC
	ctx, cancel := context.WithTimeout(context.Background(), SERVER_RPC_TIMEOUT)
	result, err := targetClient.RequestVote(ctx, &args)
	defer cancel()

	if err != nil {
		log.Println("Error requesting vote from ", address)
		log.Println(err.Error())
		retryChannel <- address
		return false
	}

	log.Println("Vote result ", result)
	// if follower term is newer immediately step down
	if result.Term > r.currentTerm {
		log.Println("Step down")
		r.stateChange <- _StepDown
		return false
	}

	// if vote is granted, return true
	if result.VoteGranted {
		return true
	}
	return false
}

func (r *RaftNode) initiateLeader() {
	// set current state as leader
	r.currentState.Store(_Leader)

	// get next log index for initialization
	r.log.indexLock.RLock()
	nextIndex := r.log.lastIndex + 1
	r.log.indexLock.RUnlock()

	// initialize indexes to track followers
	r.indexLock.Lock()
	r.nextIndex = map[string]uint64{}
	r.matchIndex = map[string]uint64{}
	// initiate indexes after each election
	r.membership.stateLock.RLock()
	for address := range r.membership.replicatedState {
		r.nextIndex[address] = nextIndex
		r.matchIndex[address] = 0
	}
	r.membership.stateLock.RUnlock()
	r.indexLock.Unlock()

	// send initial heartbeat
	go r.appendEntries(true, nil)
}

// handles RPC request or response's term
func (r *RaftNode) compareTerm(receivedTerm uint64) {
	currentTerm := r.currentTerm
	if receivedTerm > currentTerm {
		// handle new term
		r.raftLock.Lock()
		r.currentTerm = receivedTerm
		r.raftLock.Unlock()
		switch r.currentState.Load() {
		case _Leader:
			// if found a newer term, current leader is stale
			r.stateChange <- _StepDown
		}
	}
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
