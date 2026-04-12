package security

import (
	"testing"
	"time"
)

func TestHashAndVerify(t *testing.T) {
	key := "rmp_test_key_123"
	hash, err := HashKey(key)
	if err != nil {
		t.Fatalf("failed to hash key: %v", err)
	}

	if !VerifyKey(key, hash) {
		t.Errorf("failed to verify correct key")
	}

	if VerifyKey("wrong_key", hash) {
		t.Errorf("verified incorrect key")
	}
}

func TestRateLimiter(t *testing.T) {
	// 10 requests per second, burst of 2
	limiter := NewRateLimiter(10, 2)
	key := "test_api_key"

	// First two should pass (burst)
	if !limiter.Allow(key) {
		t.Errorf("first request should be allowed")
	}
	if !limiter.Allow(key) {
		t.Errorf("second request should be allowed")
	}

	// Third should be rejected (burst exhausted)
	if limiter.Allow(key) {
		t.Errorf("third request should be rejected")
	}

	// Wait for refill
	time.Sleep(150 * time.Millisecond)

	// Should allow one now
	if !limiter.Allow(key) {
		t.Errorf("request after refill should be allowed")
	}
}
