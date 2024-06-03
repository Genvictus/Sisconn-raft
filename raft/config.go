package raft

import "time"

// TIMEOUTS
var CLIENT_TIMEOUT time.Duration = 500 * time.Millisecond

var SERVER_RPC_TIMEOUT time.Duration = 500 * time.Millisecond

var (
	// The raft intervals, in milliseconds
	HEARTBEAT_INTERVAL   int = 1000
	ELECTION_TIMEOUT_MIN int = 2000
	ELECTION_TIMEOUT_MAX int = 3000
)

func SetClientTimeout(t time.Duration) {
	CLIENT_TIMEOUT = t
}

func SetServerRPCTimeout(t time.Duration) {
	SERVER_RPC_TIMEOUT = t
}

func SetRaftIntervals(heartbeat int, election_min int, election_max int) {
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
