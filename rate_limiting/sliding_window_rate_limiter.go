package rate_limiting

import (
	"sync"
	"time"
)

type SlidingWindowLimiter struct {
	limit    int
	interval time.Duration
	mutex    sync.Mutex

	currentTime   time.Time
	previousCount int
	currentCount  int
}

func NewSlidingWindowLimiter(limit int, interval time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		limit:       limit,
		interval:    interval,
		currentTime: time.Now(),
	}
}

func (l *SlidingWindowLimiter) Allow() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	newPeriodTime := l.currentTime.Add(l.interval)
	if time.Now().After(newPeriodTime) {
		l.currentTime = time.Now()
		l.previousCount = l.currentCount
		l.currentCount = 0
	}

	interval := float64(l.interval)
	currentCount := float64(l.currentCount)
	previousCount := float64(l.previousCount)
	elapsed := time.Now().Sub(l.currentTime).Seconds()
	count := (previousCount * (interval - elapsed) / interval) + currentCount
	if int(count) >= l.limit {
		return false
	}

	l.currentCount++
	return true
}
