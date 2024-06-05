package raft

import (
	"math/rand"
	"time"
)

func randMs(minMs time.Duration, maxMs time.Duration) time.Duration {
	randRange := maxMs.Milliseconds() - minMs.Milliseconds()
	randMs := rand.Int63n(randRange)
	randMs += minMs.Milliseconds()
	return time.Duration(randRange)
}
