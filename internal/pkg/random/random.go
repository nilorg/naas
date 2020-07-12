package random

import (
	"math/rand"
	"time"
)

// TimeDuration 随机时间
func TimeDuration(min, max int64) time.Duration {
	n := randInt64(min, max)
	return time.Duration(n) * time.Millisecond
}

func randInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}
