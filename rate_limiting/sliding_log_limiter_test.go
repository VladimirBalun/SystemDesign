package rate_limiting

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSlidingLogLimiter(t *testing.T) {
	t.Parallel()

	rateLimiter := NewSlidingLogLimiter(3, 100*time.Millisecond)

	assert.True(t, rateLimiter.Allow())
	assert.True(t, rateLimiter.Allow())
	assert.True(t, rateLimiter.Allow())
	assert.False(t, rateLimiter.Allow())

	time.Sleep(150 * time.Millisecond)

	assert.True(t, rateLimiter.Allow())
	assert.True(t, rateLimiter.Allow())
	assert.True(t, rateLimiter.Allow())
	assert.False(t, rateLimiter.Allow())
}

func TestSlidingLogLimiterWithGoroutines(t *testing.T) {
	t.Parallel()

	goroutinesNumber := 10
	rateLimiter := NewSlidingLogLimiter(goroutinesNumber, time.Second)

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
