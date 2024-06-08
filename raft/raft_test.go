package raft

import (
	"Sisconn-raft/raft/transport"
	"testing"
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

	if node.currentState != 0 {
		t.Errorf("Node current state is incorrect Expected: %d, but got: %d", 0, node.currentState)
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

	// TODO test connection fail
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
	// TODO: Add test cases.
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

	nodeCount = node.countNodes(false)
	if nodeCount != 5 {
		t.Errorf("Node count is incorrect Expected: %d, but got: %d", 5, nodeCount)
	}

	quorum = node.countNodes(true)
	if quorum != 3 {
		t.Errorf("Quorum count is incorrect Expected: %d, but got: %d", 3, quorum)
	}

	node.AddConnections([]string{
		serverAddress6.String(),
		serverAddress7.String(),
		serverAddress8.String(),
	})

	nodeCount = node.countNodes(false)
	if nodeCount != 8 {
		t.Errorf("Node count is incorrect Expected: %d, but got: %d", 8, nodeCount)
	}

	quorum = node.countNodes(true)
	if quorum != 5 {
		t.Errorf("Quorum count is incorrect Expected: %d, but got: %d", 5, quorum)
	}

	node.RemoveConnections([]string{
		serverAddress3.String(),
		serverAddress4.String(),
		serverAddress6.String(),
	})

	nodeCount = node.countNodes(false)
	if nodeCount != 5 {
		t.Errorf("Node count is incorrect Expected: %d, but got: %d", 5, nodeCount)
	}

	quorum = node.countNodes(true)
	if quorum != 3 {
		t.Errorf("Quorum count is incorrect Expected: %d, but got: %d", 3, quorum)
	}
}

func TestRaftNode_replicateEntry(t *testing.T) {
	// TODO: Add test cases.
}

func TestRaftNode_appendEntries(t *testing.T) {
	// TODO: Add test cases.
}

func TestRaftNode_singleAppendEntries(t *testing.T) {
	// TODO: Add test cases.
}

func TestRaftNode_requestVotes(t *testing.T) {
	// TODO: Add test cases.
}

func TestRaftNode_singleRequestVote(t *testing.T) {
	// TODO: Add test cases.
}

func TestRaftNode_getFollowerIndex(t *testing.T) {
	// TODO: Add test cases.
}

func TestRaftNode_createLogEntryArgs(t *testing.T) {
	// TODO: Add test cases.
}
