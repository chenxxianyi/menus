package middleware

import (
	"testing"
	"time"
)

func TestFixedWindowLimiterLimitsAndResets(t *testing.T) {
	limiter := NewFixedWindowLimiter(20 * time.Millisecond)
	if allowed, _ := limiter.Allow("user:1", 2); !allowed {
		t.Fatal("first request should be allowed")
	}
	if allowed, _ := limiter.Allow("user:1", 2); !allowed {
		t.Fatal("second request should be allowed")
	}
	if allowed, retry := limiter.Allow("user:1", 2); allowed || retry <= 0 {
		t.Fatalf("third request should be limited, allowed=%v retry=%v", allowed, retry)
	}
	time.Sleep(25 * time.Millisecond)
	if allowed, _ := limiter.Allow("user:1", 2); !allowed {
		t.Fatal("request should be allowed after the window resets")
	}
}
