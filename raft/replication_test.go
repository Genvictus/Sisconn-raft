package raft

import (
	"testing"
)

var (
	dummyTransactions = []TransactionEntry{
		{command: "set", key: "key1", value: "value1"},
		{command: "set", key: "key2", value: "value2"},
		{command: "set", key: "key3", value: "value3"},
		{command: "append", key: "key1", value: " append1"}, // key1: "value1 append1"
		{command: "set", key: "key2", value: "replace2"},
		{command: "delete", key: "key3", value: ""},
		{command: "set", key: "key1", value: "newvalue1"},
		{command: "append", key: "key2", value: " append2"}, // key2: "replace2 append2"
		{command: "set", key: "key3", value: "value3"},
	}

	dummyTransactions2 = []TransactionEntry{
		{command: "set", key: "1yek", value: "1eulav"},
		{command: "set", key: "2yek", value: "2eulav"},
		{command: "set", key: "3yek", value: "3eulav"},
		{command: "append", key: "1yek", value: " 1dneppa"}, // 1yek: "1eulav 1dneppa"
		{command: "set", key: "2yek", value: "2esalper"},
		{command: "delete", key: "3yek", value: ""},
		{command: "set", key: "1yek", value: "1evitadwen"},
		{command: "append", key: "2yek", value: " 2dneppa"}, // 2yek: "2esalper 2dneppa"
		{command: "set", key: "3yek", value: "3eulav"},
	}

	dummyTransactions3 = []TransactionEntry{
		{command: "set", key: "1key", value: "1value"},
		{command: "set", key: "2key", value: "2value"},
		{command: "set", key: "3key", value: "3value"},
		{command: "append", key: "1key", value: " 1append"}, // 1key: "1value 1append"
		{command: "set", key: "2key", value: "2replace"},
		{command: "delete", key: "3key", value: ""},
		{command: "set", key: "1key", value: "1newvalue"},
		{command: "append", key: "2key", value: " 2append"}, // 2key: "2replace 2append"
		{command: "set", key: "3key", value: "3value"},
	}
)

func dummyReplicationHelper() *keyValueReplication {
	replication := newKeyValueReplication()
	replication.appendTransaction(1, dummyTransactions)
	return &replication
}

func dummyReplicationHelperInt(num int) *keyValueReplication {
	replication := newKeyValueReplication()
	switch num {
	case 1:
		replication.appendTransaction(1, dummyTransactions)
	case 2:
		replication.appendTransaction(1, dummyTransactions2)
	case 3:
		replication.appendTransaction(1, dummyTransactions3)
	}
	return &replication
}

func TestNewKeyValueReplication(t *testing.T) {
	replication := newKeyValueReplication()

	if replication.lastIndex != 0 {
		t.Errorf("lastIndex not initialized properly, Expected: %d, but got: %d", 0, replication.lastIndex)
	}

	if replication.commitIndex != 0 {
		t.Errorf("commitIndex not initialized properly, Expected: %d, but got: %d", 0, replication.commitIndex)
	}

	if replication.lastApplied != 0 {
		t.Errorf("lastApplied not initialized properly, Expected: %d, but got: %d", 0, replication.lastApplied)
	}

	if len(replication.logEntries) != 1 {
		t.Errorf("logEntries not initialized properly, Expected: %d, but got: %d", 1, len(replication.logEntries))
	}
}

func TestKeyValueReplication_appendLog(t *testing.T) {
	replication := newKeyValueReplication()
	replication.appendLog(1, "key1", "value1")

	if replication.lastIndex != 1 {
		t.Errorf("lastIndex not updated properly, Expected: %d, but got: %d", 1, replication.lastIndex)
	}

	if replication.commitIndex != 0 {
		t.Errorf("commitIndex not updated properly, Expected: %d, but got: %d", 0, replication.commitIndex)
	}

	if replication.lastApplied != 0 {
		t.Errorf("lastApplied not updated properly, Expected: %d, but got: %d", 0, replication.lastApplied)
	}

	if len(replication.logEntries) != 2 {
		t.Errorf("logEntries not updated properly, Expected: %d, but got: %d", 2, len(replication.logEntries))
	}

	if replication.logEntries[1].term != 1 {
		t.Errorf("term not append properly, Expected: %d, but got: %d", 1, replication.logEntries[1].term)
	}

	if replication.logEntries[1].key != "key1" {
		t.Errorf("key not append properly, Expected: %s, but got: %s", "key1", replication.logEntries[1].key)
	}

	if replication.logEntries[1].value != "value1" {
		t.Errorf("value not append properly, Expected: %s, but got: %s", "value1", replication.logEntries[1].value)
	}
}

func TestKeyValueReplication_appendTransaction(t *testing.T) {
	replication := newKeyValueReplication()
	replication.appendTransaction(1, dummyTransactions)

	if replication.lastIndex != 9 {
		t.Errorf("lastIndex not updated properly, Expected: %d, but got: %d", 9, replication.lastIndex)
	}

	if replication.commitIndex != 0 {
		t.Errorf("commitIndex not updated properly, Expected: %d, but got: %d", 0, replication.commitIndex)
	}

	if replication.lastApplied != 0 {
		t.Errorf("lastApplied not updated properly, Expected: %d, but got: %d", 0, replication.lastApplied)
	}

	if len(replication.logEntries) != 10 {
		t.Errorf("logEntries not updated properly, Expected: %d, but got: %d", 10, len(replication.logEntries))
	}
}

func TestKeyValueReplication_replayLog(t *testing.T) {
	replication := dummyReplicationHelper()

	replication.replayLog(3, 6)

	if replication.lastApplied != 6 {
		t.Errorf("lastApplied not updated properly, Expected: %d, but got: %d", 6, replication.lastApplied)
	}

	expectedValue1 := "value1 append1"
	expectedValue2 := "replace2"
	expectedValue3 := ""

	if replication.replicatedState["key1"] != expectedValue1 {
		t.Errorf("key1 not replayed properly, Expected: %s, but got: %s", expectedValue1, replication.replicatedState["key1"])
	}

	if replication.replicatedState["key2"] != expectedValue2 {
		t.Errorf("key2 not replayed properly, Expected: %s, but got: %s", expectedValue2, replication.replicatedState["key2"])
	}

	if replication.replicatedState["key3"] != expectedValue3 {
		t.Errorf("key3 not replayed properly, Expected: %s, but got: %s", expectedValue3, replication.replicatedState["key3"])
	}
}

func TestKeyValueReplication_replaceLog(t *testing.T) {
	newTransactions := []TransactionEntry{
		{command: "set", key: "key1", value: "newvalue1"},
		{command: "append", key: "key2", value: " newappend2"},
		{command: "append", key: "key3", value: " newappend3"},
	}
	replication := newKeyValueReplication()
	replication.appendTransaction(2, newTransactions)

	sliceLogEntries := replication.logEntries[1:]
	replication.replaceLog(4, sliceLogEntries)

	if replication.lastIndex != 6 {
		t.Errorf("lastIndex not updated properly, Expected: %d, but got: %d", 6, replication.lastIndex)
	}

	if replication.commitIndex != 0 {
		t.Errorf("commitIndex not updated properly, Expected: %d, but got: %d", 0, replication.commitIndex)
	}

	if replication.lastApplied != 0 {
		t.Errorf("lastApplied not updated properly, Expected: %d, but got: %d", 0, replication.lastApplied)
	}

	if len(replication.logEntries) != 7 {
		t.Errorf("logEntries not updated properly, Expected: %d, but got: %d", 7, len(replication.logEntries))
	}

	// Test replaying newly replaced log
	replication.replayLog(1, 6)

	expectedValue1 := "newvalue1"
	expectedValue2 := " newappend2"
	expectedValue3 := " newappend3"

	if replication.replicatedState["key1"] != expectedValue1 {
		t.Errorf("key1 not replayed properly, Expected: %s, but got: %s", expectedValue1, replication.replicatedState["key1"])
	}

	if replication.replicatedState["key2"] != expectedValue2 {
		t.Errorf("key2 not replayed properly, Expected: %s, but got: %s", expectedValue2, replication.replicatedState["key2"])
	}

	if replication.replicatedState["key3"] != expectedValue3 {
		t.Errorf("key3 not replayed properly, Expected: %s, but got: %s", expectedValue3, replication.replicatedState["key3"])
	}

	// Test replace log after replay
	replication.replaceLog(7, sliceLogEntries)

	t.Log("replicated state: ", len(replication.replicatedState))
	
	if replication.lastIndex != 9 {
		t.Errorf("lastIndex not updated properly, Expected: %d, but got: %d", 9, replication.lastIndex)
	}

	if replication.commitIndex != 0 {
		t.Errorf("commitIndex not updated properly, Expected: %d, but got: %d", 0, replication.commitIndex)
	}

	if replication.lastApplied != 6 {
		t.Errorf("lastApplied not updated properly, Expected: %d, but got: %d", 6, replication.lastApplied)
	}

	if len(replication.logEntries) != 10 {
		t.Errorf("logEntries not updated properly, Expected: %d, but got: %d", 10, len(replication.logEntries))
	}

	// Test replace commited log
	replication.commitEntries(6)

	// handle panics
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Replacing commited log did not panic")
		}
	}()

	replication.replaceLog(3, sliceLogEntries)
}

func TestKeyValueReplication_commitEntries(t *testing.T) {
	replication := dummyReplicationHelper()
	replication.commitEntries(6)

	if replication.commitIndex != 6 {
		t.Errorf("commitIndex not updated properly, Expected: %d, but got: %d", 6, replication.commitIndex)
	}

	if replication.lastApplied != 6 {
		t.Errorf("lastApplied not updated properly, Expected: %d, but got: %d", 6, replication.lastApplied)
	}

	expectedValue1 := "value1 append1"
	expectedValue2 := "replace2"
	expectedValue3 := ""

	if replication.replicatedState["key1"] != expectedValue1 {
		t.Errorf("key1 not replayed properly, Expected: %s, but got: %s", expectedValue1, replication.replicatedState["key1"])
	}

	if replication.replicatedState["key2"] != expectedValue2 {
		t.Errorf("key2 not replayed properly, Expected: %s, but got: %s", expectedValue2, replication.replicatedState["key2"])
	}

	if replication.replicatedState["key3"] != expectedValue3 {
		t.Errorf("key3 not replayed properly, Expected: %s, but got: %s", expectedValue3, replication.replicatedState["key3"])
	}

	// Commit again with lower commit index
	replication.commitEntries(3)

	if replication.commitIndex != 6 {
		t.Errorf("commitIndex got updated wrongly, Expected: %d, but got: %d", 6, replication.commitIndex)
	}

	if replication.lastApplied != 6 {
		t.Errorf("lastApplied got updated wrongly, Expected: %d, but got: %d", 6, replication.lastApplied)
	}
}

func TestKeyValueReplication_getEntries(t *testing.T) {
	replication := dummyReplicationHelper()
	entries := replication.getEntries(4, 6)

	if len(entries) != 3 {
		t.Errorf("getEntries not fetched properly, Expected: %d, but got: %d", 3, len(entries))
	}

	expectedEntries := []keyValueReplicationEntry{
		{term: 1, key: "key1", value: "value1 append1"},
		{term: 1, key: "key2", value: "replace2"},
		{term: 1, key: _DELETE_KEY, value: "key3"},
	}

	for i, entry := range entries {
		if entry != expectedEntries[i] {
			t.Errorf("entry %d not fetched properly, Expected: %v, but got: %v", i+4, expectedEntries[i], entry)
		}
	}
}

func TestKeyValueReplication_get(t *testing.T) {
	replication := dummyReplicationHelper()
	replication.replayLog(3, 6)

	value1 := replication.get("key1")
	value2 := replication.get("key2")
	value3 := replication.get("key3")

	expectedValue1 := "value1 append1"
	expectedValue2 := "replace2"
	expectedValue3 := ""

	if value1 != expectedValue1 {
		t.Errorf("key1 not fetched properly, Expected: %s, but got: %s", expectedValue1, value1)
	}

	if value2 != expectedValue2 {
		t.Errorf("key2 not fetched properly, Expected: %s, but got: %s", expectedValue2, value2)
	}

	if value3 != expectedValue3 {
		t.Errorf("key3 not fetched properly, Expected: %s, but got: %s", expectedValue3, value3)
	}
}
