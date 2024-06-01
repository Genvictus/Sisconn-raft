package raft

// TIMEOUTS
const CLIENT_TIMEOUT = 500

const (
	MembershipLog = iota
	DataLog
)

const (
	NodeActive   = "active"
	NodeInactive = "inactive"
)
