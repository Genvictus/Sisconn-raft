package raft

import (
	"math/rand"
	"time"
)

func randMs(minDuration time.Duration, maxDuration time.Duration) time.Duration {
	randRange := maxDuration.Milliseconds() - minDuration.Milliseconds()
	randDuration := rand.Int63n(randRange)
	randDuration += minDuration.Milliseconds()
	return time.Duration(randDuration) * time.Millisecond
}
