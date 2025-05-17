package main

import (
	"time"
)

type TokenBucketRateLimiter struct {
	tokenBucketCh chan struct{}
	interval      time.Duration
}

func NewTokenBucketLimiter(size int, interval time.Duration) *TokenBucketRateLimiter {
	l := &TokenBucketRateLimiter{
		tokenBucketCh: make(chan struct{}, size),
		interval:      interval,
	}

	for i := 0; i < cap(l.tokenBucketCh); i++ {
		l.tokenBucketCh <- struct{}{}
	}

	go l.replenish()

	return l
}

func (l *TokenBucketRateLimiter) replenish() {
	t := time.NewTicker(l.interval)
	for range t.C {
		select {
		case l.tokenBucketCh <- struct{}{}:
		default:
		}
	}
}

func (l *TokenBucketRateLimiter) Allow() bool {
	select {
	case <-l.tokenBucketCh:
		return true
	default:
		return false
	}
}

type LeakyBucketLimiter struct {
	leakyBucketCh chan struct{}
	interval      time.Duration
}

// technically not leacky bucket.
func NewLeakyBucketLimiter(size int, interval time.Duration) *LeakyBucketLimiter {
	l := &LeakyBucketLimiter{
		leakyBucketCh: make(chan struct{}, size),
		interval:      interval,
	}

	go l.replenish()

	return l
}

func (l *LeakyBucketLimiter) replenish() {
	t := time.NewTicker(l.interval)
	// this version is busywaiting, plus it is not the shortest
	// for {
	// 	select {
	// 	case <-t.C:
	// 		select {
	// 		case l.leakyBucketCh <- struct{}{}:
	// 		default:
	// 		}
	// 	default:
	// 		continue
	// 	}
	// }
	for range t.C {
		select {
		case l.leakyBucketCh <- struct{}{}:
		default:
		}
	}
}

func (l *LeakyBucketLimiter) Allow() bool {
	select {
	case <-l.leakyBucketCh:
		return true
	default:
		return false
	}
}

func main() {
	// limiter := NewLeakyBucketLimiter(4, time.Millisecond*100)
	limiter := NewTokenBucketLimiter(4, time.Millisecond*100)

	for i := 0; i < 20; i++ {
		if limiter.Allow() {
			println("Request allowed", i)
		} else {
			println("Rate limited", i)
		}
		time.Sleep(50 * time.Millisecond)
	}
}
