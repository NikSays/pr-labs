package timeout

import (
	"math/rand"
	"time"
)

const (
	electionTimeoutMin = 5000 * time.Millisecond
	electionTimeoutMax = 6000 * time.Millisecond
)

func RandomElectionTimeout() time.Duration {
	return time.Duration(rand.Int63n(int64(electionTimeoutMax-electionTimeoutMin))) + electionTimeoutMin
}
