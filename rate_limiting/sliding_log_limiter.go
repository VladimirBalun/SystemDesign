package rate_limiting

import (
	"context"
	"sync"
	"time"
)

type SlidingLogLimiter struct {
	limit    int
	interval time.Duration
	logs     []time.Time
	mutex    sync.Mutex
}

func NewSlidingLogLimiter(ctx context.Context, limit int, interval time.Duration) *SlidingLogLimiter {
	return &SlidingLogLimiter{
		limit:    limit,
		interval: interval,
		logs:     make([]time.Time, 0),
	}
}

func (l *SlidingLogLimiter) Allow() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	lastPeriod := time.Now().Add(-l.interval)
	for len(l.logs) != 0 && l.logs[0].Before(lastPeriod) {
		l.logs = l.logs[1:]
	}

	l.logs = append(l.logs, time.Now())
	return len(l.logs) <= l.limit
}
