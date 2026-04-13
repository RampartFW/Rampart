package security

import (
	"golang.org/x/crypto/bcrypt"
	"sync"
	"time"
)

// HashKey returns a bcrypt hash of the key
func HashKey(key string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(key), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyKey compares a key with a bcrypt hash
func VerifyKey(key, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(key)) == nil
}

// RateLimiter manages rate limits for API keys
type RateLimiter struct {
	mu     sync.Mutex
	limits map[string]*tokenBucket
	rate   float64
	burst  int
	stop   chan struct{}
}

type tokenBucket struct {
	tokens float64
	last   time.Time
}

func NewRateLimiter(rate float64, burst int) *RateLimiter {
	l := &RateLimiter{
		limits: make(map[string]*tokenBucket),
		rate:   rate,
		burst:  burst,
		stop:   make(chan struct{}),
	}
	go l.cleanupLoop()
	return l
}

func (l *RateLimiter) Stop() {
	close(l.stop)
}

func (l *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			l.cleanup()
		case <-l.stop:
			return
		}
	}
}

func (l *RateLimiter) cleanup() {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	for k, tb := range l.limits {
		// If the bucket hasn't been used in 10 minutes, remove it
		if now.Sub(tb.last) > 10*time.Minute {
			delete(l.limits, k)
		}
	}
}

func (l *RateLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	tb, ok := l.limits[key]
	if !ok {
		tb = &tokenBucket{
			tokens: float64(l.burst),
			last:   time.Now(),
		}
		l.limits[key] = tb
	}

	now := time.Now()
	elapsed := now.Sub(tb.last).Seconds()
	tb.tokens += elapsed * l.rate
	if tb.tokens > float64(l.burst) {
		tb.tokens = float64(l.burst)
	}
	tb.last = now

	if tb.tokens >= 1 {
		tb.tokens--
		return true
	}
	return false
}

// Initialize performs security-related initialization tasks,
// such as dropping unnecessary process capabilities.
func Initialize() error {
	return DropCapabilities()
}
