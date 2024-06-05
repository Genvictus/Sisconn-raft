package raft

import "time"

// TIMEOUTS
var CLIENT_TIMEOUT time.Duration = 500 * time.Millisecond

var SERVER_RPC_TIMEOUT time.Duration = 500 * time.Millisecond

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
	_NodeActive   = "active"
	_NodeInactive = "inactive"
)

const (
	_RefreshFollower = iota
	_NewAppendEntry
	_StepDown
)
