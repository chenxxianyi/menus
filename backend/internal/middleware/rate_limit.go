package middleware

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"menu-recommend/internal/pkg/response"
)

type fixedWindowEntry struct {
	startedAt time.Time
	count     int
}

// FixedWindowLimiter is deliberately process-local: it provides a safe default
// for a single instance. Multi-instance deployments should replace it with a
// shared Redis-backed limiter before scaling out.
type FixedWindowLimiter struct {
	mu          sync.Mutex
	entries     map[string]fixedWindowEntry
	window      time.Duration
	lastCleanup time.Time
}

func NewFixedWindowLimiter(window time.Duration) *FixedWindowLimiter {
	if window <= 0 {
		window = time.Minute
	}
	return &FixedWindowLimiter{entries: make(map[string]fixedWindowEntry), window: window}
}

func (l *FixedWindowLimiter) Allow(key string, limit int) (bool, time.Duration) {
	if l == nil || limit <= 0 {
		return true, 0
	}
	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()
	if now.Sub(l.lastCleanup) >= l.window {
		for entryKey, entry := range l.entries {
			if now.Sub(entry.startedAt) >= l.window {
				delete(l.entries, entryKey)
			}
		}
		l.lastCleanup = now
	}
	entry := l.entries[key]
	if entry.startedAt.IsZero() || now.Sub(entry.startedAt) >= l.window {
		l.entries[key] = fixedWindowEntry{startedAt: now, count: 1}
		return true, 0
	}
	if entry.count >= limit {
		return false, l.window - now.Sub(entry.startedAt)
	}
	entry.count++
	l.entries[key] = entry
	return true, 0
}

func RateLimitByIP(limiter *FixedWindowLimiter, scope string, limit int) gin.HandlerFunc {
	return rateLimit(limiter, scope, limit, func(c *gin.Context) string {
		return c.ClientIP()
	})
}

func RateLimitByUserOrIP(limiter *FixedWindowLimiter, scope string, limit int) gin.HandlerFunc {
	return rateLimit(limiter, scope, limit, func(c *gin.Context) string {
		if userID := GetUserID(c); userID > 0 {
			return "user:" + strconv.FormatUint(uint64(userID), 10)
		}
		return "ip:" + c.ClientIP()
	})
}

func rateLimit(limiter *FixedWindowLimiter, scope string, limit int, keyFn func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		allowed, retryAfter := limiter.Allow(scope+":"+keyFn(c), limit)
		if allowed {
			c.Next()
			return
		}
		seconds := int(retryAfter.Round(time.Second).Seconds())
		if seconds < 1 {
			seconds = 1
		}
		c.Header("Retry-After", strconv.Itoa(seconds))
		response.TooManyRequests(c, fmt.Sprintf("操作过于频繁，请在 %d 秒后重试", seconds))
		c.Abort()
	}
}

type ConcurrencyLimiter struct {
	semaphore chan struct{}
}

func NewConcurrencyLimiter(maxConcurrent int) *ConcurrencyLimiter {
	if maxConcurrent <= 0 {
		return nil
	}
	return &ConcurrencyLimiter{semaphore: make(chan struct{}, maxConcurrent)}
}

func LimitConcurrency(limiter *ConcurrencyLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if limiter == nil {
			c.Next()
			return
		}
		select {
		case limiter.semaphore <- struct{}{}:
			defer func() { <-limiter.semaphore }()
			c.Next()
		default:
			response.TooManyRequests(c, "AI 服务繁忙，请稍后重试")
			c.Abort()
		}
	}
}
