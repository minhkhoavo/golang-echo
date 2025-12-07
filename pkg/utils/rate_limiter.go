package utils

import (
	"context"
	"sync"
	"time"
)

// RateLimiter defines the interface for rate limiting implementations
type RateLimiter interface {
	Allow(key string) bool
	AllowContext(ctx context.Context, key string) bool
	Reset()
	Close() error
}

// TokenBucketLimiter implements token bucket algorithm rate limiter
// This is suitable for single-instance deployment
type TokenBucketLimiter struct {
	buckets map[string]*tokenBucket
	mu      sync.RWMutex
	rate    int           // tokens per refill period
	window  time.Duration // refill period
	ticker  *time.Ticker
	done    chan struct{}
}

type tokenBucket struct {
	tokens    int
	lastRefill time.Time
}

// NewTokenBucketLimiter creates a new token bucket rate limiter
// rate: number of requests allowed per window period
// window: the time period for refill (e.g., 1 minute)
func NewTokenBucketLimiter(rate int, window time.Duration) *TokenBucketLimiter {
	limiter := &TokenBucketLimiter{
		buckets: make(map[string]*tokenBucket),
		rate:    rate,
		window:  window,
		ticker:  time.NewTicker(window),
		done:    make(chan struct{}),
	}

	// Goroutine to refill tokens periodically
	go limiter.refillLoop()

	return limiter
}

// Allow checks if a request from the given key is allowed
func (tbl *TokenBucketLimiter) Allow(key string) bool {
	tbl.mu.Lock()
	defer tbl.mu.Unlock()

	bucket, exists := tbl.buckets[key]
	if !exists {
		// New key, initialize bucket
		tbl.buckets[key] = &tokenBucket{
			tokens:     tbl.rate - 1,
			lastRefill: time.Now(),
		}
		return true
	}

	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}

	return false
}

// AllowContext is the context-aware version of Allow
func (tbl *TokenBucketLimiter) AllowContext(ctx context.Context, key string) bool {
	select {
	case <-ctx.Done():
		return false
	default:
		return tbl.Allow(key)
	}
}

// refillLoop refills tokens at regular intervals
func (tbl *TokenBucketLimiter) refillLoop() {
	for {
		select {
		case <-tbl.ticker.C:
			tbl.mu.Lock()
			for _, bucket := range tbl.buckets {
				bucket.tokens = tbl.rate
			}
			tbl.mu.Unlock()

		case <-tbl.done:
			tbl.ticker.Stop()
			return
		}
	}
}

// Reset clears all rate limit data
func (tbl *TokenBucketLimiter) Reset() {
	tbl.mu.Lock()
	defer tbl.mu.Unlock()
	tbl.buckets = make(map[string]*tokenBucket)
}

// Close stops the rate limiter
func (tbl *TokenBucketLimiter) Close() error {
	close(tbl.done)
	return nil
}

// ==============================================================
// SlidingWindowLimiter - Alternative implementation using sliding window
// Better for accurate rate limiting
// ==============================================================

type SlidingWindowLimiter struct {
	windows map[string][]time.Time
	mu      sync.RWMutex
	limit   int           // max requests
	window  time.Duration // time window
	cleanup *time.Ticker
	done    chan struct{}
}

// NewSlidingWindowLimiter creates a new sliding window rate limiter
func NewSlidingWindowLimiter(limit int, window time.Duration) *SlidingWindowLimiter {
	limiter := &SlidingWindowLimiter{
		windows: make(map[string][]time.Time),
		limit:   limit,
		window:  window,
		cleanup: time.NewTicker(window),
		done:    make(chan struct{}),
	}

	// Cleanup old entries periodically
	go limiter.cleanupLoop()

	return limiter
}

// Allow checks if a request from the given key is allowed
func (swl *SlidingWindowLimiter) Allow(key string) bool {
	swl.mu.Lock()
	defer swl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-swl.window)

	// Get or create window for this key
	timestamps, exists := swl.windows[key]
	if !exists {
		timestamps = []time.Time{}
	}

	// Remove timestamps outside the window
	validTimestamps := []time.Time{}
	for _, ts := range timestamps {
		if ts.After(windowStart) {
			validTimestamps = append(validTimestamps, ts)
		}
	}

	// Check if limit exceeded
	if len(validTimestamps) >= swl.limit {
		swl.windows[key] = validTimestamps
		return false
	}

	// Add current request
	validTimestamps = append(validTimestamps, now)
	swl.windows[key] = validTimestamps

	return true
}

// AllowContext is the context-aware version of Allow
func (swl *SlidingWindowLimiter) AllowContext(ctx context.Context, key string) bool {
	select {
	case <-ctx.Done():
		return false
	default:
		return swl.Allow(key)
	}
}

// Reset clears all rate limit data
func (swl *SlidingWindowLimiter) Reset() {
	swl.mu.Lock()
	defer swl.mu.Unlock()
	swl.windows = make(map[string][]time.Time)
}

// cleanupLoop removes old entries to prevent memory leak
func (swl *SlidingWindowLimiter) cleanupLoop() {
	for {
		select {
		case <-swl.cleanup.C:
			swl.mu.Lock()
			now := time.Now()
			windowStart := now.Add(-swl.window)

			for key, timestamps := range swl.windows {
				validTimestamps := []time.Time{}
				for _, ts := range timestamps {
					if ts.After(windowStart) {
						validTimestamps = append(validTimestamps, ts)
					}
				}

				if len(validTimestamps) == 0 {
					delete(swl.windows, key)
				} else {
					swl.windows[key] = validTimestamps
				}
			}
			swl.mu.Unlock()

		case <-swl.done:
			swl.cleanup.Stop()
			return
		}
	}
}

// Close stops the rate limiter
func (swl *SlidingWindowLimiter) Close() error {
	close(swl.done)
	return nil
}
