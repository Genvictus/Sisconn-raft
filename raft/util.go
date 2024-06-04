package raft

import (
	"math/rand"
	"time"
)

func randMs(minMs int, maxMs int) time.Duration {
	randRange := maxMs - minMs
	randMs := rand.Intn(randRange)
	randMs += minMs
	return time.Duration(randRange) * time.Millisecond
}
