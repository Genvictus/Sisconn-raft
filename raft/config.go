package raft

import "time"

// TIMEOUTS
var CLIENT_TIMEOUT time.Duration = 500 * time.Millisecond

var (
	RPC_TIMEOUT          time.Duration = 500 * time.Millisecond
	HEARTBEAT_INTERVAL   time.Duration = time.Second
	ELECTION_TIMEOUT_MIN time.Duration = 2 * time.Second
	ELECTION_TIMEOUT_MAX time.Duration = 3 * time.Second
)

func SetClientTimeout(t time.Duration) {
	CLIENT_TIMEOUT = t
}

func SetServerRPCTimeout(t time.Duration) {
	RPC_TIMEOUT = t
}

func SetRaftIntervals(heartbeat time.Duration, election_min time.Duration, election_max time.Duration) {
	HEARTBEAT_INTERVAL = heartbeat
	ELECTION_TIMEOUT_MIN = election_min
	ELECTION_TIMEOUT_MAX = election_max
}

// Raft Constants
const (
	MembershipLog = iota
	DataLog
)

const (
	Leader = iota
	Candidate
	Follower
)

const (
	NodeActive   = "active"
	NodeInactive = "inactive"
)
