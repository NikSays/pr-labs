package main

import (
	"math/rand"
	"time"
)

const (
	electionTimeoutMin = 5000 * time.Millisecond
	electionTimeoutMax = 6000 * time.Millisecond
)

func randomElectionTimeout() time.Duration {
	return time.Duration(rand.Int63n(int64(electionTimeoutMax-electionTimeoutMin))) + electionTimeoutMin
}
