package raft

import (
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

var ServerLogger *log.Logger

func randMs(minDuration time.Duration, maxDuration time.Duration) time.Duration {
	randRange := maxDuration.Milliseconds() - minDuration.Milliseconds()
	randDuration := rand.Int63n(randRange)
	randDuration += minDuration.Milliseconds()
	return time.Duration(randDuration) * time.Millisecond
}

func redirectStdout() (closeBuf func() (out []byte)) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	closeBuf = func() (out []byte) {
		w.Close()
		out, _ = io.ReadAll(r)
		os.Stdout = rescueStdout
		return out
	}

	return closeBuf
}

func SetLogger(logger *log.Logger) {
	if logger != nil {
		ServerLogger = logger
	}
}
