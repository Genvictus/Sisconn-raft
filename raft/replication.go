package raft

import "sync"

type keyValuelogPlayer interface {
	// modify accordingly

	// Applies log from starting index to the last index specified
	replayLog(startIndex uint64, lastIndex uint64)

	// Append a new entry to the log, and then play the appended log entry
	appendLog(term uint64, key string, value string)

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

	indexLock   sync.RWMutex
	commitIndex uint64
	lastApplied uint64

	// example of mutex lock in golang
	logLock    sync.RWMutex
	logEntries []keyValueReplicationEntry

	replicatedState map[string]string
}

func newKeyValueReplication() keyValueReplication {
	return keyValueReplication{
		commitIndex: 0,
		lastApplied: 0,

		replicatedState: map[string]string{},
	}
}

func (k *keyValueReplication) replayLog(startIndex uint64, lastIndex uint64) {
	for i := startIndex; i <= lastIndex; i++ {
		k.replicatedState[k.logEntries[i].key] = k.logEntries[i].value
	}
}

func (k *keyValueReplication) appendLog(term uint64, key string, value string) {
	k.logEntries = append(k.logEntries, keyValueReplicationEntry{term: term, key: key, value: value})
	if k.replicatedState == nil {
		k.replicatedState = map[string]string{}
	}
	k.replicatedState[key] = value
}

func (k *keyValueReplication) replaceLog(startIndex uint64, logEntries []keyValueReplicationEntry) {
	k.logEntries = append(k.logEntries[:startIndex], logEntries...)
	k.replicatedState = map[string]string{}
	k.replayLog(0, uint64(len(k.logEntries))-1)
}

func (k *keyValueReplication) get(key string) string {
	return k.replicatedState[key]
}
