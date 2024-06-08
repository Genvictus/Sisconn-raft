package raft

import (
	"testing"
	"time"
)

func TestSetClientTimeout(t *testing.T) {
    expected := 5 * time.Second
    SetClientTimeout(expected)
    
    if CLIENT_TIMEOUT != expected {
        t.Errorf("SetClientTimeout failed, expected %v, got %v", expected, CLIENT_TIMEOUT)
    }

    expected = 10 * time.Second
    SetClientTimeout(expected)
    
    if CLIENT_TIMEOUT != expected {
        t.Errorf("SetClientTimeout failed, expected %v, got %v", expected, CLIENT_TIMEOUT)
    }
}

func TestSetServerRPCTimeout(t *testing.T) {
    expected := 5 * time.Second
    SetServerRPCTimeout(expected)
    
    if SERVER_RPC_TIMEOUT != expected {
        t.Errorf("SetServerRPCTimeout failed, expected %v, got %v", expected, SERVER_RPC_TIMEOUT)
    }

    expected = 10 * time.Second
    SetServerRPCTimeout(expected)
    
    if SERVER_RPC_TIMEOUT != expected {
        t.Errorf("SetServerRPCTimeout failed, expected %v, got %v", expected, SERVER_RPC_TIMEOUT)
    }
}

func TestSetRaftIntervals(t *testing.T) {
    expectedHeartbeat := 1 * time.Second
    expectedElectionMin := 2 * time.Second
    expectedElectionMax := 3 * time.Second
    
    SetRaftIntervals(expectedHeartbeat, expectedElectionMin, expectedElectionMax)
    
    if HEARTBEAT_INTERVAL != expectedHeartbeat {
        t.Errorf("SetRaftIntervals failed for HEARTBEAT_INTERVAL, expected %v, got %v", expectedHeartbeat, HEARTBEAT_INTERVAL)
    }
    if ELECTION_TIMEOUT_MIN != expectedElectionMin {
        t.Errorf("SetRaftIntervals failed for ELECTION_TIMEOUT_MIN, expected %v, got %v", expectedElectionMin, ELECTION_TIMEOUT_MIN)
    }
    if ELECTION_TIMEOUT_MAX != expectedElectionMax {
        t.Errorf("SetRaftIntervals failed for ELECTION_TIMEOUT_MAX, expected %v, got %v", expectedElectionMax, ELECTION_TIMEOUT_MAX)
    }

    expectedHeartbeat = 10 * time.Second
    expectedElectionMin = 20 * time.Second
    expectedElectionMax = 30 * time.Second
    
    SetRaftIntervals(expectedHeartbeat, expectedElectionMin, expectedElectionMax)
    
    if HEARTBEAT_INTERVAL != expectedHeartbeat {
        t.Errorf("SetRaftIntervals failed for HEARTBEAT_INTERVAL, expected %v, got %v", expectedHeartbeat, HEARTBEAT_INTERVAL)
    }
    if ELECTION_TIMEOUT_MIN != expectedElectionMin {
        t.Errorf("SetRaftIntervals failed for ELECTION_TIMEOUT_MIN, expected %v, got %v", expectedElectionMin, ELECTION_TIMEOUT_MIN)
    }
    if ELECTION_TIMEOUT_MAX != expectedElectionMax {
        t.Errorf("SetRaftIntervals failed for ELECTION_TIMEOUT_MAX, expected %v, got %v", expectedElectionMax, ELECTION_TIMEOUT_MAX)
    }
}
