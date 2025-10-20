package security

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"go-backend-api/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
	requests map[string]*TokenBucket
	mutex    sync.RWMutex
	rate     int           // requests per minute
	capacity int           // burst capacity
	cleanup  time.Duration // cleanup interval
}

// TokenBucket represents a token bucket for rate limiting
type TokenBucket struct {
	tokens     int
	lastRefill time.Time
	rate       int
	capacity   int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate, capacity int) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]*TokenBucket),
		rate:     rate,
		capacity: capacity,
		cleanup:  time.Minute * 5,
	}

	// Start cleanup goroutine
	go rl.cleanupExpiredBuckets()

	return rl
}

// Allow checks if a request is allowed for the given key
func (rl *RateLimiter) Allow(key string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	bucket, exists := rl.requests[key]

	if !exists {
		bucket = &TokenBucket{
			tokens:     rl.capacity - 1,
			lastRefill: now,
			rate:       rl.rate,
			capacity:   rl.capacity,
		}
		rl.requests[key] = bucket
		return true
	}

	// Refill tokens based on time elapsed
	elapsed := now.Sub(bucket.lastRefill)
	tokensToAdd := int(elapsed.Minutes()) * bucket.rate

	if tokensToAdd > 0 {
		bucket.tokens = min(bucket.capacity, bucket.tokens+tokensToAdd)
		bucket.lastRefill = now
	}

	// Check if tokens available
	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}

	return false
}

// cleanupExpiredBuckets removes expired token buckets
func (rl *RateLimiter) cleanupExpiredBuckets() {
	ticker := time.NewTicker(rl.cleanup)
	defer ticker.Stop()

	for range ticker.C {
		rl.mutex.Lock()
		now := time.Now()
		for key, bucket := range rl.requests {
			// Remove buckets that haven't been used for 10 minutes
			if now.Sub(bucket.lastRefill) > time.Minute*10 {
				delete(rl.requests, key)
			}
		}
		rl.mutex.Unlock()
	}
}

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware(rate, capacity int) gin.HandlerFunc {
	limiter := NewRateLimiter(rate, capacity)

	return func(c *gin.Context) {
		// Get client IP
		clientIP := c.ClientIP()
		key := fmt.Sprintf("%s:%s", clientIP, c.Request.URL.Path)

		if !limiter.Allow(key) {
			c.JSON(http.StatusTooManyRequests, response.Response{
				Success: false,
				Error: &response.ErrorInfo{
					Code:    http.StatusTooManyRequests,
					Message: "Rate limit exceeded. Please try again later.",
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AuthRateLimitMiddleware creates a rate limiting middleware for auth endpoints
func AuthRateLimitMiddleware() gin.HandlerFunc {
	// Stricter rate limiting for auth endpoints
	return RateLimitMiddleware(5, 10) // 5 requests per minute, burst of 10
}

// APIRateLimitMiddleware creates a rate limiting middleware for API endpoints
func APIRateLimitMiddleware() gin.HandlerFunc {
	// More lenient rate limiting for API endpoints
	return RateLimitMiddleware(100, 200) // 100 requests per minute, burst of 200
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
