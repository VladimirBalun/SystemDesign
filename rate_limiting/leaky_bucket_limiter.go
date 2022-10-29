package rate_limiting

import (
	"context"
	"time"
)

type LeakyBucketLimiter struct {
	leakyBucketCh chan struct{}
}

func NewLeakyBucketLimiter(ctx context.Context, limit int, period time.Duration) *LeakyBucketLimiter {
	limiter := &LeakyBucketLimiter{
		leakyBucketCh: make(chan struct{}, limit),
	}

	leakInterval := period.Nanoseconds() / int64(limit)
	go limiter.startPeriodicLeak(ctx, time.Duration(leakInterval))
	return limiter
}

func (l *LeakyBucketLimiter) startPeriodicLeak(ctx context.Context, interval time.Duration) {
	timer := time.NewTicker(interval)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			select {
			case <-l.leakyBucketCh:
			default:
			}
		}
	}
}

func (l *LeakyBucketLimiter) Allow() bool {
	select {
	case l.leakyBucketCh <- struct{}{}:
		return true
	default:
		return false
	}
}
