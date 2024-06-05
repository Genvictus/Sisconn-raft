package test

import (
	"Sisconn-raft/raft"
	"testing"
	"time"
)

func TestSetClientTimeout(t *testing.T) {
    expected := 5 * time.Second
    raft.SetClientTimeout(expected)
    
    if raft.CLIENT_TIMEOUT != expected {
        t.Errorf("SetClientTimeout failed, expected %v, got %v", expected, raft.CLIENT_TIMEOUT)
    }

    expected = 10 * time.Second
    raft.SetClientTimeout(expected)
    
    if raft.CLIENT_TIMEOUT != expected {
        t.Errorf("SetClientTimeout failed, expected %v, got %v", expected, raft.CLIENT_TIMEOUT)
    }
}

func TestSetServerRPCTimeout(t *testing.T) {
    expected := 5 * time.Second
    raft.SetServerRPCTimeout(expected)
    
    if raft.SERVER_RPC_TIMEOUT != expected {
        t.Errorf("SetServerRPCTimeout failed, expected %v, got %v", expected, raft.SERVER_RPC_TIMEOUT)
    }

    expected = 10 * time.Second
    raft.SetServerRPCTimeout(expected)
    
    if raft.SERVER_RPC_TIMEOUT != expected {
        t.Errorf("SetServerRPCTimeout failed, expected %v, got %v", expected, raft.SERVER_RPC_TIMEOUT)
    }
}

func TestSetRaftIntervals(t *testing.T) {
    expectedHeartbeat := 1 * time.Second
    expectedElectionMin := 2 * time.Second
    expectedElectionMax := 3 * time.Second
    
    raft.SetRaftIntervals(expectedHeartbeat, expectedElectionMin, expectedElectionMax)
    
    if raft.HEARTBEAT_INTERVAL != expectedHeartbeat {
        t.Errorf("SetRaftIntervals failed for HEARTBEAT_INTERVAL, expected %v, got %v", expectedHeartbeat, raft.HEARTBEAT_INTERVAL)
    }
    if raft.ELECTION_TIMEOUT_MIN != expectedElectionMin {
        t.Errorf("SetRaftIntervals failed for ELECTION_TIMEOUT_MIN, expected %v, got %v", expectedElectionMin, raft.ELECTION_TIMEOUT_MIN)
    }
    if raft.ELECTION_TIMEOUT_MAX != expectedElectionMax {
        t.Errorf("SetRaftIntervals failed for ELECTION_TIMEOUT_MAX, expected %v, got %v", expectedElectionMax, raft.ELECTION_TIMEOUT_MAX)
    }

    expectedHeartbeat = 10 * time.Second
    expectedElectionMin = 20 * time.Second
    expectedElectionMax = 30 * time.Second
    
    raft.SetRaftIntervals(expectedHeartbeat, expectedElectionMin, expectedElectionMax)
    
    if raft.HEARTBEAT_INTERVAL != expectedHeartbeat {
        t.Errorf("SetRaftIntervals failed for HEARTBEAT_INTERVAL, expected %v, got %v", expectedHeartbeat, raft.HEARTBEAT_INTERVAL)
    }
    if raft.ELECTION_TIMEOUT_MIN != expectedElectionMin {
        t.Errorf("SetRaftIntervals failed for ELECTION_TIMEOUT_MIN, expected %v, got %v", expectedElectionMin, raft.ELECTION_TIMEOUT_MIN)
    }
    if raft.ELECTION_TIMEOUT_MAX != expectedElectionMax {
        t.Errorf("SetRaftIntervals failed for ELECTION_TIMEOUT_MAX, expected %v, got %v", expectedElectionMax, raft.ELECTION_TIMEOUT_MAX)
    }
}
