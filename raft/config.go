package raft

import (
	"fmt"
	"os"
	"time"
)

// TIMEOUTS
var CLIENT_TIMEOUT time.Duration = 500 * time.Millisecond

var SERVER_RPC_TIMEOUT time.Duration = 500 * time.Millisecond
var REPLICATION_TIMEOUT time.Duration = 500 * time.Millisecond

var (
	HEARTBEAT_INTERVAL   time.Duration = 1000 * time.Millisecond
	ELECTION_TIMEOUT_MIN time.Duration = 2000 * time.Millisecond
	ELECTION_TIMEOUT_MAX time.Duration = 3000 * time.Millisecond
)

func SetClientTimeout(t time.Duration) {
	CLIENT_TIMEOUT = t
}

func SetServerRPCTimeout(t time.Duration) {
	SERVER_RPC_TIMEOUT = t
}

func SetRaftIntervals(heartbeat time.Duration, election_min time.Duration, election_max time.Duration) {
	HEARTBEAT_INTERVAL = heartbeat
	ELECTION_TIMEOUT_MIN = election_min
	ELECTION_TIMEOUT_MAX = election_max
}

func LoadRaftConfig() {
	// from os.env
	timeoutClientStr := os.Getenv("CLIENT_TIMEOUT") + "ms"
	timeoutClient, err := time.ParseDuration(timeoutClientStr)
	if err != nil {
		// Use default value
		timeoutClient = CLIENT_TIMEOUT
	}

	timeoutServerStr := os.Getenv("SERVER_RPC_TIMEOUT") + "ms"
	timeoutServer, err := time.ParseDuration(timeoutServerStr)
	if err != nil {
		// Use default value
		// log.Println(err)
		timeoutServer = SERVER_RPC_TIMEOUT
	}

	heartbeatStr := os.Getenv("HEARTBEAT_INTERVAL") + "ms"
	heartbeat, err := time.ParseDuration(heartbeatStr)
	if err != nil {
		// Use default value
		heartbeat = HEARTBEAT_INTERVAL
	}

	electionMinStr := os.Getenv("ELECTION_TIMEOUT_MIN") + "ms"

	electionMin, err := time.ParseDuration(electionMinStr)
	if err != nil {
		// Use default value
		electionMin = ELECTION_TIMEOUT_MIN
	}

	electionMaxStr := os.Getenv("ELECTION_TIMEOUT_MAX") + "ms"
	electionMax, err := time.ParseDuration(electionMaxStr)
	if err != nil {
		// Use default value
		electionMax = ELECTION_TIMEOUT_MAX
	}

	SetClientTimeout(timeoutClient)
	SetServerRPCTimeout(timeoutServer)
	SetRaftIntervals(heartbeat, electionMin, electionMax)

	// PrintConfig()
}

func PrintConfig() {
	fmt.Println("Client Timeout:", CLIENT_TIMEOUT)
	fmt.Println("Server RPC Timeout:", SERVER_RPC_TIMEOUT)
	fmt.Println("Heartbeat Interval:", HEARTBEAT_INTERVAL)
	fmt.Println("Election Timeout Min:", ELECTION_TIMEOUT_MIN)
	fmt.Println("Election Timeout Max:", ELECTION_TIMEOUT_MAX)
}

// Raft Constants
const (
	_MembershipLog = iota
	_DataLog
)

const (
	_Follower = iota
	_Candidate
	_Leader
)

const (
	OkResponse        = "OK"
	PendingResponse   = "PENDING"
	NotLeaderResponse = "PLEASE CONTACT LEADER ADDRESS "
)

const (
	_NodeActive   = "active"
	_NodeInactive = "inactive"
)

const (
	_RefreshFollower = iota
	_NewAppendEntry
	_StepDown
)
