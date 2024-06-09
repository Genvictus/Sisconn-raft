package raft

import (
	pb "Sisconn-raft/raft/raftpc"
	"Sisconn-raft/raft/transport"
	"context"
	"log"
	"net"
	"reflect"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc"
)

func TestNewNode(t *testing.T) {
	serverAddress := transport.NewAddress("localhost", 2001)
	node := NewNode(serverAddress.String())
	if node == nil {
		t.Fatal("New node is nil")
	}

	if node.address != "localhost:2001" {
		t.Errorf("Node address is incorrect Expected: %s, but got: %s", serverAddress.String(), node.address)
	}

	if node.currentState.Load() != 0 {
		t.Errorf("Node current state is incorrect Expected: %d, but got: %d", 0, node.currentState.Load())
	}
}

func TestRaftNode_AddConnections(t *testing.T) {
	// Test 3 node connection
	serverAddress1 := transport.NewAddress("localhost", 2001)
	serverAddress2 := transport.NewAddress("localhost", 2002)
	serverAddress3 := transport.NewAddress("localhost", 2003)

	node := NewNode(serverAddress1.String())
	node.AddConnections([]string{
		serverAddress1.String(),
		serverAddress2.String(),
		serverAddress3.String(),
	})

	if len(node.raftClient) != 2 {
		t.Errorf("Node raft client count is incorrect Expected: %d, but got: %d", 2, len(node.raftClient))
	}

	if len(node.conn) != 2 {
		t.Errorf("Node connections length is incorrect Expected: %d, but got: %d", 2, len(node.conn))
	}

	if len(node.membership.logEntries) != 4 {
		t.Errorf("Node membership length is incorrect Expected: %d, but got: %d", 4, len(node.membership.logEntries))
	}

	// Test new connection fail
	closeBuf := redirectStdout()

	node.AddConnections([]string{
		"localhost%:2001",
	})

	out := closeBuf()

	string_out := string(out)

	if !strings.Contains(string_out, "fails") {
		t.Errorf("Expected error message to contain: %s, but got: %s", "fails", string_out)
	}
}

func TestRaftNode_RemoveConnections(t *testing.T) {
	serverAddress1 := transport.NewAddress("localhost", 2001)
	serverAddress2 := transport.NewAddress("localhost", 2002)
	serverAddress3 := transport.NewAddress("localhost", 2003)
	serverAddress4 := transport.NewAddress("localhost", 2004)
	serverAddress5 := transport.NewAddress("localhost", 2005)

	node := NewNode(serverAddress1.String())
	node.AddConnections([]string{
		serverAddress1.String(),
		serverAddress2.String(),
		serverAddress3.String(),
		serverAddress4.String(),
		serverAddress5.String(),
	})

	node.RemoveConnections([]string{
		serverAddress3.String(),
		serverAddress4.String(),
	})

	if len(node.raftClient) != 2 {
		t.Errorf("Node raft client count is incorrect Expected: %d, but got: %d", 2, len(node.raftClient))
	}

	if len(node.conn) != 2 {
		t.Errorf("Node connections length is incorrect Expected: %d, but got: %d", 2, len(node.conn))
	}

	if node.conn[serverAddress3.String()] != nil {
		t.Errorf("Node connections for server %s should be nil", serverAddress3.String())
	}

	if node.conn[serverAddress4.String()] != nil {
		t.Errorf("Node connections for server %s should be nil", serverAddress4.String())
	}
}

func TestRaftNode_Run(t *testing.T) {
	serverAddress := transport.NewAddress("localhost", 2024)
	node := NewNode(serverAddress.String())
	ListenServer(node, grpc.NewServer())

	node.AddConnections([]string{
		serverAddress.String(),
	})

	go node.Run()

	// Prepare config
	PREVIOUS_HEARTBEAT_INTERVAL := HEARTBEAT_INTERVAL
	PREVIOUS_ELECTION_TIMEOUT_MIN := ELECTION_TIMEOUT_MIN
	PREVIOUS_ELECTION_TIMEOUT_MAX := ELECTION_TIMEOUT_MAX

	SetRaftIntervals(1*time.Millisecond, 2*time.Millisecond, 3*time.Millisecond)

	// Test stepdown leader
	node.currentState.Store(_Leader)
	node.stateChange <- _StepDown

	state := node.currentState.Load()
	if state != _Follower {
		t.Errorf("Expected leader state to be _Follower, but got: %d", state)
	}

	// Test stepdown candidate
	// refresh follower
	node.currentState.Store(_Candidate)
	node.stateChange <- _RefreshFollower

	node.stateChange <- _StepDown

	state = node.currentState.Load()
	if state != _Follower {
		t.Errorf("Expected candidate state to be _Follower, but got: %d", state)
	}

	// TODO test timer timeout if can

	// Restore config
	SetRaftIntervals(PREVIOUS_HEARTBEAT_INTERVAL, PREVIOUS_ELECTION_TIMEOUT_MIN, PREVIOUS_ELECTION_TIMEOUT_MAX)
}

func TestRaftNode_RunTest(t *testing.T) {
	serverAddress := transport.NewAddress("localhost", 2021)
	node := NewNode(serverAddress.String())
	ListenServer(node, grpc.NewServer())

	node.AddConnections([]string{
		serverAddress.String(),
	})

	go node.runTest()

	// Test stepdown leader
	node.currentState.Store(_Leader)
	node.stateChange <- _StepDown

	state := node.currentState.Load()
	if state != _Follower {
		t.Errorf("Expected leader state to be _Follower, but got: %d", state)
	}

	// Test stepdown candidate
	// refresh follower
	node.currentState.Store(_Candidate)
	node.stateChange <- _RefreshFollower

	node.stateChange <- _StepDown

	state = node.currentState.Load()
	if state != _Follower {
		t.Errorf("Expected candidate state to be _Follower, but got: %d", state)
	}
}

func TestRaftNode_countNodes(t *testing.T) {
	serverAddress1 := transport.NewAddress("localhost", 2001)
	serverAddress2 := transport.NewAddress("localhost", 2002)
	serverAddress3 := transport.NewAddress("localhost", 2003)
	serverAddress4 := transport.NewAddress("localhost", 2004)
	serverAddress5 := transport.NewAddress("localhost", 2005)
	serverAddress6 := transport.NewAddress("localhost", 2006)
	serverAddress7 := transport.NewAddress("localhost", 2007)
	serverAddress8 := transport.NewAddress("localhost", 2008)

	node := NewNode(serverAddress1.String())
	node.AddConnections([]string{
		serverAddress1.String(),
		serverAddress2.String(),
		serverAddress3.String(),
		serverAddress4.String(),
		serverAddress5.String(),
	})

	var nodeCount int
	var quorum int

	nodeCount, quorum = node.countNodes()
	if nodeCount != 5 {
		t.Errorf("Node count is incorrect Expected: %d, but got: %d", 5, nodeCount)
	}
	if quorum != 3 {
		t.Errorf("Quorum count is incorrect Expected: %d, but got: %d", 3, quorum)
	}

	node.AddConnections([]string{
		serverAddress6.String(),
		serverAddress7.String(),
		serverAddress8.String(),
	})

	nodeCount, quorum = node.countNodes()
	if nodeCount != 8 {
		t.Errorf("Node count is incorrect Expected: %d, but got: %d", 8, nodeCount)
	}
	if quorum != 5 {
		t.Errorf("Quorum count is incorrect Expected: %d, but got: %d", 5, quorum)
	}

	node.RemoveConnections([]string{
		serverAddress3.String(),
		serverAddress4.String(),
		serverAddress6.String(),
	})

	nodeCount, quorum = node.countNodes()
	if nodeCount != 5 {
		t.Errorf("Node count is incorrect Expected: %d, but got: %d", 5, nodeCount)
	}
	if quorum != 3 {
		t.Errorf("Quorum count is incorrect Expected: %d, but got: %d", 3, quorum)
	}
}

func TestRaftNode_replicateEntry(t *testing.T) {
	serverAddress1 := transport.NewAddress("localhost", 2341)
	serverAddress2 := transport.NewAddress("localhost", 5435)

	node := NewNode(serverAddress1.String())
	node2 := NewNode(serverAddress2.String())
	ListenServer(node, grpc.NewServer())
	ListenServer(node2, grpc.NewServer())

	node.AddConnections([]string{
		serverAddress1.String(),
		serverAddress2.String(),
	})

	node.log = *dummyReplicationHelperInt(2)
	node.initiateLeader()
	go node.runTest()
	go node2.runTest()

	ctx, cancel := context.WithCancel(context.Background())

	// Replicate Entry
	val := node.replicateEntry(ctx)

	if val != true {
		t.Errorf("Expected success, but got: %t", val)
	}

	// Testing ctx timeout
	cancel()
	val = node.replicateEntry(ctx)

	if val != false {
		t.Errorf("Expected failed, but got: %t", val)
	}
}

func TestRaftNode_appendEntries(t *testing.T) {
	serverAddress1 := transport.NewAddress("localhost", 5321)
	serverAddress2 := transport.NewAddress("localhost", 6465)

	node := NewNode(serverAddress1.String())
	node2 := NewNode(serverAddress2.String())

	ListenServer(node, grpc.NewServer())
	ListenServer(node2, grpc.NewServer())

	node.AddConnections([]string{
		serverAddress1.String(),
		serverAddress2.String(),
	})

	node.log = *dummyReplicationHelperInt(2)
	node.initiateLeader()

	go node.runTest()
	go node2.runTest()

	// Append Entries
	commitedChan := make(chan bool, 1)
	node.appendEntries(false, commitedChan)

	val := <-commitedChan
	if val != true {
		t.Errorf("Expected success, but got: %t", val)
	}
}

func TestRaftNode_singleAppendEntries(t *testing.T) {
	serverAddress1 := transport.NewAddress("localhost", 8011)
	serverAddress2 := transport.NewAddress("localhost", 8012)

	node := NewNode(serverAddress1.String())
	node2 := NewNode(serverAddress2.String())

	ListenServer(node, grpc.NewServer())
	ListenServer(node2, grpc.NewServer())

	node.AddConnections([]string{
		serverAddress1.String(),
		serverAddress2.String(),
	})

	node.log = *dummyReplicationHelperInt(1)
	node.initiateLeader()

	go node.runTest()
	go node2.runTest()

	// Single Append Entries
	val := node.singleAppendEntries(node2.address, false)
	if val != true {
		t.Errorf("Expected success, but got: %t", val)
	}

	// Heartbeat
	_ = node.singleAppendEntries(node2.address, true)

	if !reflect.DeepEqual(&node.log, &node2.log) {
		log.Println("Node log ", &node.log)
		log.Println("Node 2 log ", &node2.log)
		t.Errorf("Expected log to be equal, but got: %v", &node2.log)
	}

	// Test not leader
	node.currentState.Store(_Follower)

	val = node.singleAppendEntries(node2.address, false)
	if val != false {
		t.Errorf("Expected failed, but got: %t", val)
	}
}

func TestRaftNode_requestVotes(t *testing.T) {
	serverAddress1 := transport.NewAddress("localhost", 2010)
	serverAddress2 := transport.NewAddress("localhost", 2011)
	serverAddress3 := transport.NewAddress("localhost", 2012)
	serverAddress := transport.NewAddress("10.255.255.1", 80)
	serverAddress4 := transport.NewAddress("localhost", 2006)
	serverAddress5 := transport.NewAddress("localhost", 2007)
	serverAddress6 := transport.NewAddress("localhost", 2008)

	node := NewNode(serverAddress1.String())
	ListenServer(node, grpc.NewServer())
	node2 := NewNode(serverAddress2.String())
	ListenServer(node2, grpc.NewServer())
	node3 := NewNode(serverAddress3.String())
	ListenServer(node3, grpc.NewServer())
	node4 := NewNode(serverAddress4.String())
	ListenServer(node4, grpc.NewServer())
	node5 := NewNode(serverAddress5.String())
	ListenServer(node5, grpc.NewServer())
	node6 := NewNode(serverAddress6.String())
	ListenServer(node6, grpc.NewServer())

	node.AddConnections([]string{
		serverAddress1.String(),
		serverAddress2.String(),
		serverAddress3.String(),
		serverAddress4.String(),
		serverAddress.String(),
		serverAddress5.String(),
		serverAddress6.String(),
	})

	// start node
	go node.runTest()
	go node2.runTest()
	go node3.runTest()
	go node4.runTest()
	go node5.runTest()
	go node6.runTest()

	// Timeout to Follower
	t.Log("Testing Timeout to Follower")
	node.currentState.Store(_Candidate)
	node2.currentTerm = 69
	node3.currentTerm = 69
	node4.currentTerm = 69
	node5.currentTerm = 69
	node.requestVotes()

	state := node.currentState.Load()
	if state != _Follower {
		t.Errorf("Expected state to be _Follower, but got: %d", state)
	}

	// Leader
	t.Log("Testing Leader")
	node.currentState.Store(_Candidate)
	node.currentTerm = 70
	node.requestVotes()

	state = node.currentState.Load()
	if state != _Leader {
		t.Errorf("Expected state to be _Leader, but got: %d", state)
	}

	// Test single node case
	serverAddress7 := transport.NewAddress("localhost", 2017)
	node7 := NewNode(serverAddress1.String())
	node7.AddConnections([]string{
		serverAddress7.String(),
	})
	go node7.runTest()

	node7.currentState.Store(_Candidate)

	node7.requestVotes()

	state = node7.currentState.Load()
	if state != _Leader {
		t.Errorf("Expected state to be _Leader, but got: %d", state)
	}
}

func TestRaftNode_singleRequestVote(t *testing.T) {
	serverAddress1 := transport.NewAddress("localhost", 2000)
	serverAddress2 := transport.NewAddress("localhost", 2004)
	serverAddress3 := transport.NewAddress("localhost", 2005)
	serverAddress4 := transport.NewAddress("10.255.255.1", 80)

	node := NewNode(serverAddress1.String())
	ListenServer(node, grpc.NewServer())
	node2 := NewNode(serverAddress2.String())
	ListenServer(node2, grpc.NewServer())
	node3 := NewNode(serverAddress3.String())
	ListenServer(node3, grpc.NewServer())

	node.AddConnections([]string{
		serverAddress1.String(),
		serverAddress2.String(),
		serverAddress3.String(),
		serverAddress4.String(),
	})

	// start node
	go node.runTest()
	go node2.runTest()
	go node3.runTest()

	// Test vote request
	lastIndex := node.log.lastIndex
	lastTerm := node.log.getEntries(lastIndex, lastIndex)[0].term

	// vote granted
	node.currentTerm = 2
	result := node.singleRequestVote(node2.address, lastIndex, lastTerm)

	if result != true {
		t.Errorf("Vote request failed")
	}
	t.Log("Vote request success : ", result)

	// step down
	node3.currentTerm = 69
	result = node.singleRequestVote(node3.address, lastIndex, lastTerm)

	if result != false {
		t.Errorf("Vote request failed")
	}
	t.Log("Vote request failed : ", result)

	// state := <-node.stateChange
	// if state != _StepDown {
	// 	t.Errorf("Expected state change to _StepDown, but got: %d", state)
	// }

	// error vote
	result = node.singleRequestVote(serverAddress4.String(), lastIndex, lastTerm)

	if result != false {
		t.Errorf("Vote request failed")
	}
}

func TestRaftNode_initiateLeader(t *testing.T) {
	serverAddress1 := transport.NewAddress("localhost", 2014)
	serverAddress2 := transport.NewAddress("localhost", 2015)
	serverAddress3 := transport.NewAddress("localhost", 2016)

	node1 := NewNode(serverAddress1.String())
	ListenServer(node1, grpc.NewServer())
	node2 := NewNode(serverAddress2.String())
	ListenServer(node2, grpc.NewServer())
	node3 := NewNode(serverAddress3.String())
	ListenServer(node3, grpc.NewServer())

	node1.AddConnections([]string{
		serverAddress1.String(),
		serverAddress2.String(),
		serverAddress3.String(),
	})

	// start node
	go node1.runTest()
	go node2.runTest()
	go node3.runTest()

	node1.initiateLeader()

	state := node1.currentState.Load()
	if state != _Leader {
		t.Errorf("Expected state to be _Leader, but got: %d", state)
	}
}

func TestRaftNode_CompareTerm(t *testing.T) {
	serverAddress1 := transport.NewAddress("localhost", 2343)

	node1 := NewNode(serverAddress1.String())
	ListenServer(node1, grpc.NewServer())
	go node1.runTest()
	node1.currentTerm = 1
	node1.currentState.Store(_Leader)

	go node1.compareTerm(2)
	state := <-node1.stateChange

	if state != _StepDown {
		t.Errorf("Expected state change to _StepDown, but got: %d", state)
	}
}

func TestRaftNode_createLogEntryArgs(t *testing.T) {
	serverAddress := transport.NewAddress("localhost", 2344)

	node := NewNode(serverAddress.String())

	node.log.logEntries = []keyValueReplicationEntry{
		{term: 1, key: "key1", value: "value1"},
		{term: 1, key: "key2", value: "value2"},
		{term: 1, key: "key3", value: "value3"},
		{term: 2, key: "key1", value: " append1"},
		{term: 2, key: "key2", value: "replace2"},
		{term: 2, key: "key3", value: ""},
		{term: 3, key: "key1", value: "newvalue1"},
		{term: 3, key: "key2", value: " append2"},
		{term: 3, key: "key3", value: "value3"},
	}

	args := node.createLogEntryArgs(3, 5)

	for i, entry := range args {
		if entry.Term != node.log.logEntries[i+3].term {
			t.Errorf("Expected term to be %d, but got: %d", node.log.logEntries[i+3].term, entry.Term)
		}
	}
}

func ListenServer(node *RaftNode, server *grpc.Server) {
	serverAddress := node.address
	lis, err := net.Listen("tcp", serverAddress)

	if err != nil {
		log.Fatalf("failed to listen on %s", serverAddress)
	}

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	pb.RegisterRaftServiceServer(server, &ServiceServer{Server: node})
	pb.RegisterRaftServer(server, &RaftServer{Server: node})
}
