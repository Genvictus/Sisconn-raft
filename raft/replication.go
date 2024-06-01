package raft

import (
	"sync"
)

type keyValuelogPlayer interface {
	// modify accordingly

	// Applies log from starting index to the last index specified
	replayLog(startIndex uint64, lastIndex uint64)

	// Append a new entry to the log, and then play the appended log entry
	appendLog(term uint64, key string, value string)
	appendMultipleLog(term uint64, entries []TransactionEntry)

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

type TransactionEntry struct {
	command string
	key     string
	value   string
}

func newKeyValueReplication() keyValueReplication {
	return keyValueReplication{
		commitIndex: 0,
		lastApplied: 0,

		replicatedState: map[string]string{},
	}
}

func (k *keyValueReplication) replayLog(startIndex uint64, lastIndex uint64) {
	k.logLock.RLock()
	for i := startIndex; i <= lastIndex; i++ {
		k.replicatedState[k.logEntries[i].key] = k.logEntries[i].value
	}

	k.logLock.RUnlock()
}

func (k *keyValueReplication) appendLog(term uint64, key string, value string) {
	k.indexLock.Lock()
	k.logLock.Lock()

	k.logEntries = append(k.logEntries, keyValueReplicationEntry{term: term, key: key, value: value})
	if k.replicatedState == nil {
		k.replicatedState = map[string]string{}
	}
	k.replicatedState[key] = value

	k.indexLock.Unlock()
	k.logLock.Unlock()
}

func (k *keyValueReplication) appendMultipleLog(term uint64, entries []TransactionEntry) {
	k.indexLock.Lock()
	k.logLock.Lock()

	if k.replicatedState == nil {
		k.replicatedState = map[string]string{}
	}
	for _, entry := range entries {
		newval := k.get(entry.key)
		if entry.command == "append" {
			newval = newval + entry.value
		} else {
			newval = entry.value
		}
		k.logEntries = append(k.logEntries, keyValueReplicationEntry{term: term, key: entry.key, value: newval})
		k.replicatedState[entry.key] = newval
	}

	k.indexLock.Unlock()
	k.logLock.Unlock()
}

func (k *keyValueReplication) replaceLog(startIndex uint64, logEntries []keyValueReplicationEntry) {
	k.indexLock.Lock()
	k.logLock.Lock()

	k.logEntries = append(k.logEntries[:startIndex], logEntries...)

	k.indexLock.Unlock()
	k.logLock.Unlock()

	k.replicatedState = map[string]string{}
	k.replayLog(0, uint64(len(k.logEntries))-1)
}

func (k *keyValueReplication) get(key string) string {
	return k.replicatedState[key]
}
