package raft

import (
	"log"
	"sync"
)

const (
	_DELETE_KEY  = ""
	_START_INDEX = 1
)

var (
	_DUMMY_ENTRY = keyValueReplicationEntry{
		term:  0,
		key:   "",
		value: "",
	}
)

type keyValuelogPlayer interface {
	// Append a new entry to the log
	appendLog(term uint64, key string, value string)
	// Replace the log's entry starting from the startIndex
	replaceLog(startIndex uint64, logEntries []keyValueReplicationEntry)
	// commit the entire log entries up to the specified index, and applying its logs as well
	commitEntries(lastIndex uint64)
	// Applies log from starting index to the last index specified
	replayLog(startIndex uint64, lastIndex uint64)
	// Get the log entries specified from startIndex to lastIndex (inclusive)
	getEntries(startIndex uint64, lastIndex uint64) []keyValueReplicationEntry

	appendTransaction(term uint64, entries []TransactionEntry)
}

type keyValueReplicationEntry struct {
	// probably need more than this

	term uint64
	// if key is an empty string, deletes the key in the replicated state
	// the deleted key in replicated state is in the value field instead
	key string
	// maybe store old value for rollback?
	// oldValue string
	value string
}

type keyValueReplication struct {
	indexLock sync.RWMutex
	// Indexes start at 0. However, logEntries use 1 indexing
	lastIndex   uint64
	commitIndex uint64
	lastApplied uint64

	// example of mutex lock in golang
	logLock    sync.RWMutex
	logEntries []keyValueReplicationEntry

	stateLock       sync.RWMutex
	replicatedState map[string]string

	// TO AVOID DEADLOCK, LOCK ACQUIRING ORDER: indexLock -> logLock -> stateLock
	// follow the lock from the top of the list to the bottom
}

type TransactionEntry struct {
	command string
	key     string
	value   string
}

func newKeyValueReplication() keyValueReplication {
	return keyValueReplication{
		lastIndex:   0,
		commitIndex: 0,
		lastApplied: 0,

		logEntries: []keyValueReplicationEntry{_DUMMY_ENTRY},

		replicatedState: map[string]string{},
	}
}

func (k *keyValueReplication) replayLog(startIndex uint64, lastIndex uint64) {
	k.indexLock.Lock()
	k.logLock.RLock()
	k.stateLock.Lock()

	for i := startIndex; i <= lastIndex; i++ {
		currentKey := k.logEntries[i].key
		// if the log entry is deletion
		if currentKey == _DELETE_KEY {
			delete(k.replicatedState, k.logEntries[i].value)
		} else {
			// else append the new value
			k.replicatedState[currentKey] = k.logEntries[i].value
		}
	}
	// update appliy index
	k.lastApplied = lastIndex

	k.stateLock.Unlock()
	k.logLock.RUnlock()
	k.indexLock.Unlock()
}

func (k *keyValueReplication) appendLog(term uint64, key string, value string) {
	k.indexLock.Lock()
	k.logLock.Lock()

	// append and update index
	k.logEntries = append(k.logEntries, keyValueReplicationEntry{term: term, key: key, value: value})
	k.lastIndex++

	k.logLock.Unlock()
	k.indexLock.Unlock()
}

func (k *keyValueReplication) appendTransaction(term uint64, entries []TransactionEntry) {
	k.indexLock.Lock()
	k.logLock.Lock()

	for _, entry := range entries {
		newval := k.get(entry.key)
		if entry.command == "append" {
			newval = newval + entry.value
		} else {
			newval = entry.value
		}
		// append and update index
		k.logEntries = append(k.logEntries, keyValueReplicationEntry{term: term, key: entry.key, value: newval})
		k.lastIndex++
	}

	k.logLock.Unlock()
	k.indexLock.Unlock()
}

func (k *keyValueReplication) replaceLog(startIndex uint64, logEntries []keyValueReplicationEntry) {
	k.indexLock.Lock()
	k.logLock.Lock()

	k.logEntries = append(k.logEntries[:startIndex], logEntries...)
	k.lastIndex = startIndex - 1 + uint64(len(logEntries))

	k.logLock.Unlock()
	k.indexLock.Unlock()
}

func (k *keyValueReplication) commitEntries(lastIndex uint64) {
	if lastIndex <= k.commitIndex {
		return
	}
	k.replayLog(k.lastApplied+1, lastIndex)
	k.indexLock.Lock()
	k.commitIndex = lastIndex
	log.Printf("Log Committed at index: %d\n", k.commitIndex)
	k.indexLock.Unlock()
}

func (k *keyValueReplication) getEntries(startIndex uint64, lastIndex uint64) []keyValueReplicationEntry {
	k.logLock.RLock()
	entries := k.logEntries[startIndex : lastIndex+1]
	k.logLock.RUnlock()
	return entries
}

func (k *keyValueReplication) get(key string) string {
	k.stateLock.RLock()
	val := k.replicatedState[key]
	k.stateLock.RUnlock()
	return val
}
