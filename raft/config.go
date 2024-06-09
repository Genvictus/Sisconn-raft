package raft

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// TIMEOUTS
var CLIENT_TIMEOUT time.Duration = 500 * time.Millisecond

var SERVER_RPC_TIMEOUT time.Duration = 500 * time.Millisecond
var REPLICATION_TIMEOUT time.Duration = 500 * time.Millisecond

var (
	HEARTBEAT_INTERVAL   time.Duration = 3000 * time.Millisecond
	ELECTION_TIMEOUT_MIN time.Duration = 6000 * time.Millisecond
	ELECTION_TIMEOUT_MAX time.Duration = 9000 * time.Millisecond
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
	timeoutClientInt, err := strconv.Atoi(os.Getenv("CLIENT_TIMEOUT"))
	timeoutClient := time.Duration(timeoutClientInt) * time.Millisecond
	if err != nil {
		timeoutClient = CLIENT_TIMEOUT
	}

	timeoutServerInt, err := strconv.Atoi(os.Getenv("SERVER_RPC_TIMEOUT"))
	timeoutServer := time.Duration(timeoutServerInt) * time.Millisecond
	if err != nil {
		timeoutServer = SERVER_RPC_TIMEOUT
	}

	heartbeatInt, err := strconv.Atoi(os.Getenv("HEARTBEAT_INTERVAL"))
	heartbeat := time.Duration(heartbeatInt) * time.Millisecond
	if err != nil {
		heartbeat = HEARTBEAT_INTERVAL
	}

	electionMinInt, err := strconv.Atoi(os.Getenv("ELECTION_TIMEOUT_MIN"))
	electionMin := time.Duration(electionMinInt) * time.Millisecond
	if err != nil {
		electionMin = ELECTION_TIMEOUT_MIN
	}

	electionMaxInt, err := strconv.Atoi(os.Getenv("ELECTION_TIMEOUT_MAX"))
	electionMax := time.Duration(electionMaxInt) * time.Millisecond
	if err != nil {
		electionMax = ELECTION_TIMEOUT_MAX
	}

	SetClientTimeout(timeoutClient)
	SetServerRPCTimeout(timeoutServer)
	SetRaftIntervals(heartbeat, electionMin, electionMax)
}

func PrintConfig() {
	fmt.Println("Client Timeout:", CLIENT_TIMEOUT)
	fmt.Println("Server RPC Timeout:", SERVER_RPC_TIMEOUT)
	fmt.Println("Heartbeat Interval:", HEARTBEAT_INTERVAL)
	fmt.Println("Election Timeout Min:", ELECTION_TIMEOUT_MIN)
	fmt.Println("Election Timeout Max:", ELECTION_TIMEOUT_MAX)
}

func LoadClientConfig() {
	// from os.env
	timeoutClientInt, err := strconv.Atoi(os.Getenv("CLIENT_TIMEOUT"))
	timeoutClient := time.Duration(timeoutClientInt) * time.Millisecond
	if err != nil {
		timeoutClient = CLIENT_TIMEOUT
	}

	SetClientTimeout(timeoutClient)
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
