package raft

import "sync"

type keyValuelogPlayer interface {
	// modify accordingly

	// Applies log from starting index to the last index specified
	replayLog(startIndex uint64, lastIndex uint64)

	// Append a new entry to the log, and then play the appended log entry
	appendLog(key string, value string)

	// Replace the log's entry starting from the startIndex
	// (rolling back the state as well), then append the new entries
	replaceLog(startIndex uint64, logEntries []keyValueReplicationEntry)
}

type keyValueReplicationEntry struct {
	// probably need more than this

	term uint64
	key  string
	// maybe store old value for rollback?
	// oldValue string
	value string
}

type keyValueReplication struct {
	// modify accordingly

	indexLock sync.RWMutex
	lastIndex uint64
	// maybe?
	lastComittedIndex uint64

	// example of mutex lock in golang
	logLock    sync.RWMutex
	logEntries []keyValueReplicationEntry

	replicatedState map[string]string
}
