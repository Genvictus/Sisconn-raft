package raft

import (
	"os"
	"strconv"
	"testing"
	"time"
)

func TestSetClientTimeout(t *testing.T) {
    tests := []struct {
        expected time.Duration
    }{
        {5 * time.Second},
        {10 * time.Second},
        {1000 * time.Millisecond},
        {500 * time.Millisecond},
        {1500 * time.Millisecond},
    }

    for _, tt := range tests {
        t.Run("TestSetClientTimeout", func(t *testing.T) {
            SetClientTimeout(tt.expected)
            if CLIENT_TIMEOUT != tt.expected {
                t.Errorf("SetClientTimeout failed, expected %v, got %v", tt.expected, CLIENT_TIMEOUT)
            }
        })
    }
}

func TestSetServerRPCTimeout(t *testing.T) {
    tests := []struct {
        expected time.Duration
    }{
        {5 * time.Second},
        {10 * time.Second},
        {1000 * time.Millisecond},
        {500 * time.Millisecond},
        {1500 * time.Millisecond},
    }

    for _, tt := range tests {
        t.Run("TestSetServerRPCTimeout", func(t *testing.T) {
            SetServerRPCTimeout(tt.expected)
            if SERVER_RPC_TIMEOUT != tt.expected {
                t.Errorf("SetServerRPCTimeout failed, expected %v, got %v", tt.expected, SERVER_RPC_TIMEOUT)
            }
        })
    }
}

func TestSetRaftIntervals(t *testing.T) {
    tests := []struct {
        expectedHeartbeat time.Duration
        expectedElectionMin time.Duration
        expectedElectionMax time.Duration
    }{
        {5 * time.Second, 3 * time.Second, 5 * time.Second, },
        {10 * time.Second, 8 * time.Second, 9 * time.Second, },
        {1000 * time.Millisecond, 100 * time.Millisecond, 200 * time.Millisecond, },
        {500 * time.Millisecond, 400 * time.Millisecond, 600 * time.Millisecond, },
        {1500 * time.Millisecond, 1400 * time.Millisecond, 1600 * time.Millisecond, },
    }

    for _, tt := range tests {
        t.Run("TestSetRaftIntervals", func(t *testing.T) {
            SetRaftIntervals(tt.expectedHeartbeat, tt.expectedElectionMin, tt.expectedElectionMax)

            if HEARTBEAT_INTERVAL != tt.expectedHeartbeat {
                t.Errorf("SetRaftIntervals failed for HEARTBEAT_INTERVAL, expected %v, got %v", tt.expectedHeartbeat, HEARTBEAT_INTERVAL)
            }

            if ELECTION_TIMEOUT_MIN != tt.expectedElectionMin {
                t.Errorf("SetRaftIntervals failed for ELECTION_TIMEOUT_MIN, expected %v, got %v", tt.expectedElectionMin, ELECTION_TIMEOUT_MIN)
            }

            if ELECTION_TIMEOUT_MAX != tt.expectedElectionMax {
                t.Errorf("SetRaftIntervals failed for ELECTION_TIMEOUT_MAX, expected %v, got %v", tt.expectedElectionMax, ELECTION_TIMEOUT_MAX)
            }
        })
    }
}

func TestLoadRaftConfig(t *testing.T) {
	// Test null environment variables
	// Remove environment variables after testing
	os.Unsetenv("CLIENT_TIMEOUT")
	os.Unsetenv("SERVER_RPC_TIMEOUT")
	os.Unsetenv("HEARTBEAT_INTERVAL")
	os.Unsetenv("ELECTION_TIMEOUT_MIN")
	os.Unsetenv("ELECTION_TIMEOUT_MAX")

	expectedClientTimeout := CLIENT_TIMEOUT
	expectedServerRPCTimeout := SERVER_RPC_TIMEOUT
	expectedHeartbeatInterval := HEARTBEAT_INTERVAL
	expectedElectionTimeoutMin := ELECTION_TIMEOUT_MIN
	expectedElectionTimeoutMax := ELECTION_TIMEOUT_MAX

	LoadRaftConfig()

	if CLIENT_TIMEOUT != expectedClientTimeout {
		t.Errorf("LoadRaftConfig failed for CLIENT_TIMEOUT, expected %v, got %v", expectedClientTimeout, CLIENT_TIMEOUT)
	}

	if SERVER_RPC_TIMEOUT != expectedServerRPCTimeout {
		t.Errorf("LoadRaftConfig failed for SERVER_RPC_TIMEOUT, expected %v, got %v", expectedServerRPCTimeout, SERVER_RPC_TIMEOUT)
	}

	if HEARTBEAT_INTERVAL != expectedHeartbeatInterval {
		t.Errorf("LoadRaftConfig failed for HEARTBEAT_INTERVAL, expected %v, got %v", expectedHeartbeatInterval, HEARTBEAT_INTERVAL)
	}

	if ELECTION_TIMEOUT_MIN != expectedElectionTimeoutMin {
		t.Errorf("LoadRaftConfig failed for ELECTION_TIMEOUT_MIN, expected %v, got %v", expectedElectionTimeoutMin, ELECTION_TIMEOUT_MIN)
	}

	if ELECTION_TIMEOUT_MAX != expectedElectionTimeoutMax {
		t.Errorf("LoadRaftConfig failed for ELECTION_TIMEOUT_MAX, expected %v, got %v", expectedElectionTimeoutMax, ELECTION_TIMEOUT_MAX)
	}

	// Test set environment variables
	tests := []struct {
		clientTimeout      int
		serverRpcTimeout   int
		heartbeatInterval  int
		electionTimeoutMin int
		electionTimeoutMax int
	}{
		{100, 200, 300, 400, 500},
		{1000, 2000, 3000, 4000, 5000},
		{20, 35, 100, 30, 40},
	}
	for _, tt := range tests {
		t.Run("TestLoadRaftConfig", func(t *testing.T) {
			// Set environment variables for testing
			os.Setenv("CLIENT_TIMEOUT", strconv.Itoa(tt.clientTimeout))
			os.Setenv("SERVER_RPC_TIMEOUT", strconv.Itoa(tt.serverRpcTimeout))
			os.Setenv("HEARTBEAT_INTERVAL", strconv.Itoa(tt.heartbeatInterval))
			os.Setenv("ELECTION_TIMEOUT_MIN", strconv.Itoa(tt.electionTimeoutMin))
			os.Setenv("ELECTION_TIMEOUT_MAX", strconv.Itoa(tt.electionTimeoutMax))

			LoadRaftConfig()

			expectedClientTimeout := time.Duration(tt.clientTimeout) * time.Millisecond
			expectedServerRPCTimeout := time.Duration(tt.serverRpcTimeout) * time.Millisecond
			expectedHeartbeatInterval := time.Duration(tt.heartbeatInterval) * time.Millisecond
			expectedElectionTimeoutMin := time.Duration(tt.electionTimeoutMin) * time.Millisecond
			expectedElectionTimeoutMax := time.Duration(tt.electionTimeoutMax) * time.Millisecond

			if CLIENT_TIMEOUT != expectedClientTimeout {
				t.Errorf("LoadRaftConfig failed for CLIENT_TIMEOUT, expected %v, got %v", expectedClientTimeout, CLIENT_TIMEOUT)
			}

			if SERVER_RPC_TIMEOUT != expectedServerRPCTimeout {
				t.Errorf("LoadRaftConfig failed for SERVER_RPC_TIMEOUT, expected %v, got %v", expectedServerRPCTimeout, SERVER_RPC_TIMEOUT)
			}

			if HEARTBEAT_INTERVAL != expectedHeartbeatInterval {
				t.Errorf("LoadRaftConfig failed for HEARTBEAT_INTERVAL, expected %v, got %v", expectedHeartbeatInterval, HEARTBEAT_INTERVAL)
			}

			if ELECTION_TIMEOUT_MIN != expectedElectionTimeoutMin {
				t.Errorf("LoadRaftConfig failed for ELECTION_TIMEOUT_MIN, expected %v, got %v", expectedElectionTimeoutMin, ELECTION_TIMEOUT_MIN)
			}

			if ELECTION_TIMEOUT_MAX != expectedElectionTimeoutMax {
				t.Errorf("LoadRaftConfig failed for ELECTION_TIMEOUT_MAX, expected %v, got %v", expectedElectionTimeoutMax, ELECTION_TIMEOUT_MAX)
			}
		})
	}

	// Clean up environment variables after testing
	os.Unsetenv("CLIENT_TIMEOUT")
	os.Unsetenv("SERVER_RPC_TIMEOUT")
	os.Unsetenv("HEARTBEAT_INTERVAL")
	os.Unsetenv("ELECTION_TIMEOUT_MIN")
	os.Unsetenv("ELECTION_TIMEOUT_MAX")
}

func TestPrintConfig(t *testing.T) {
    // Capture the output
	closeBuf := redirectStdout()

    SetClientTimeout(100 * time.Millisecond)
    SetServerRPCTimeout(200 * time.Millisecond)
    SetRaftIntervals(300 * time.Millisecond, 400 * time.Millisecond, 500 * time.Millisecond)

	PrintConfig()
    
	out := closeBuf()

    string_out := string(out)

	expectedOutput := `Client Timeout: 100ms
Server RPC Timeout: 200ms
Heartbeat Interval: 300ms
Election Timeout Min: 400ms
Election Timeout Max: 500ms
`

	if string_out != expectedOutput {
		t.Errorf("expected %q but got %s", expectedOutput, string_out)
	}
}

func TestLoadClientConfig(t *testing.T) {
	os.Unsetenv("CLIENT_TIMEOUT")

	expectedClientTimeout := CLIENT_TIMEOUT

	LoadClientConfig()

	if CLIENT_TIMEOUT != expectedClientTimeout {
		t.Errorf("LoadClientConfig failed for CLIENT_TIMEOUT, expected %v, got %v", expectedClientTimeout, CLIENT_TIMEOUT)
	}

	tests := []struct {
		clientTimeout int
	}{
		{5},
        {10},
        {1000},
        {500},
        {1500},
	}
	
	for _, tt := range tests {
		t.Run("TestLoadClientConfig", func(t *testing.T) {
			os.Setenv("CLIENT_TIMEOUT", strconv.Itoa(tt.clientTimeout))

			LoadClientConfig()

			expectedClientTimeout := time.Duration(tt.clientTimeout) * time.Millisecond

			if CLIENT_TIMEOUT != expectedClientTimeout {
				t.Errorf("LoadClientConfig failed for CLIENT_TIMEOUT, expected %v, got %v", expectedClientTimeout, CLIENT_TIMEOUT)
			}
		})
	}

	os.Unsetenv("CLIENT_TIMEOUT")
}
