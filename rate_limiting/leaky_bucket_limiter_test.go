package rate_limiting

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLeakyBucketLimiter(t *testing.T) {
	t.Parallel()

	rateLimiter := NewLeakyBucketLimiter(context.Background(), 3, time.Second)

	assert.True(t, rateLimiter.Allow())
	assert.True(t, rateLimiter.Allow())
	assert.True(t, rateLimiter.Allow())
	assert.False(t, rateLimiter.Allow())

	time.Sleep(350 * time.Millisecond)

	assert.True(t, rateLimiter.Allow())
	assert.False(t, rateLimiter.Allow())
}

func TestLeakyBucketLimiterWithGoroutines(t *testing.T) {
	t.Parallel()

	goroutinesNumber := 10
	rateLimiter := NewLeakyBucketLimiter(context.Background(), goroutinesNumber, time.Second)

	wg := sync.WaitGroup{}
	wg.Add(goroutinesNumber)
	for i := 0; i < goroutinesNumber; i++ {
		go func() {
			defer wg.Done()
			assert.True(t, rateLimiter.Allow())
		}()
	}

	wg.Wait()
	assert.False(t, rateLimiter.Allow())
}

func TestLeakyBucketLimiterWithCancel(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	rateLimiter := NewLeakyBucketLimiter(ctx, 3, 500*time.Millisecond)

	assert.True(t, rateLimiter.Allow())
	assert.True(t, rateLimiter.Allow())
	assert.True(t, rateLimiter.Allow())
	assert.False(t, rateLimiter.Allow())

	cancel()

	time.Sleep(550 * time.Millisecond)
	assert.False(t, rateLimiter.Allow())
}
