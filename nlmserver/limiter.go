package main

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type itemLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	limiters map[string]*itemLimiter
	header   string
	rate     rate.Limit
	burst    int
	next     http.Handler
	done     chan bool
}

func NewRateLimiter(header string, r rate.Limit, n int, next http.Handler) *RateLimiter {
	rl := &RateLimiter{
		limiters: make(map[string]*itemLimiter),
		header:   header,
		rate:     r,
		burst:    n,
		next:     next,
		done:     make(chan bool),
	}

	go rl.cleanupLimiters(rl.done)

	return rl
}

func (rl *RateLimiter) cleanupLimiters(done <-chan bool) {
	for {
		select {
		case <-time.After(time.Minute):
			rl.mu.Lock()
			for k, v := range rl.limiters {
				if time.Since(v.lastSeen) > 2*time.Minute {
					delete(rl.limiters, k)
				}
			}

			rl.mu.Unlock()
		case <-done:
			return
		}
	}
}

func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.limiters[key]
	if !exists {
		limiter := rate.NewLimiter(rl.rate, rl.burst)

		rl.limiters[key] = &itemLimiter{limiter, time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func (rl *RateLimiter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	limiterKey := r.Header.Get(rl.header)

	limiter := rl.getLimiter(limiterKey)

	if limiter.Allow() == false {
		log.Printf("Too many requests for key: %s", limiterKey)
		w.Header().Add("HX-Location", "/speed")
		w.Header().Set("Retry-After", strconv.FormatFloat(1.0/float64(rl.rate), 'f', 0, 64))
		http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		return
	}

	rl.next.ServeHTTP(w, r)
}

func (rl *RateLimiter) Close() {
	close(rl.done)
}
