package rate_limiting

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLeakyBucketLimiter(t *testing.T) {
	rateLimiter := NewLeakyBucketLimiter(context.Background(), 3, time.Second)

	assert.True(t, rateLimiter.Allow())
	assert.True(t, rateLimiter.Allow())
	assert.True(t, rateLimiter.Allow())
	assert.False(t, rateLimiter.Allow())

	time.Sleep(350 * time.Millisecond)

	assert.True(t, rateLimiter.Allow())
	assert.False(t, rateLimiter.Allow())
}

func TestLeakyBucketLimiterWithCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	rateLimiter := NewLeakyBucketLimiter(ctx, 3, 500*time.Millisecond)

	assert.True(t, rateLimiter.Allow())
	assert.True(t, rateLimiter.Allow())
	assert.True(t, rateLimiter.Allow())
	assert.False(t, rateLimiter.Allow())

	cancel()

	time.Sleep(500 * time.Millisecond)
	assert.False(t, rateLimiter.Allow())
}
