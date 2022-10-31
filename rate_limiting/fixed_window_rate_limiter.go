package rate_limiting

import (
	"context"
	"sync/atomic"
	"time"
)

type FixedWindowLimiter struct {
	count int32
	limit int
}

func NewFixedWindowLimiter(ctx context.Context, limit int, interval time.Duration) *FixedWindowLimiter {
	limiter := &FixedWindowLimiter{
		count: 0,
		limit: limit,
	}

	go limiter.startPeriodicCountRefresh(ctx, interval)
	return limiter
}

func (l *FixedWindowLimiter) startPeriodicCountRefresh(ctx context.Context, interval time.Duration) {
	timer := time.NewTicker(interval)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			atomic.StoreInt32(&l.count, 0)
		}
	}
}

func (l *FixedWindowLimiter) Allow() bool {
	count := atomic.LoadInt32(&l.count)
	if count > int32(l.limit) {
		return false
	}

	for !atomic.CompareAndSwapInt32(&l.count, count, count+1) {
		count = atomic.LoadInt32(&l.count)
	}

	return count+1 <= int32(l.limit)
}
