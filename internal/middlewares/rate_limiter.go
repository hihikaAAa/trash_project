package middlewares

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type clientState struct {
	limiter  *rate.Limiter
	lastSeen time.Time
	blocked  bool
}

type RateLimiter struct {
	nextCleanup  time.Time
	clients      map[string]*clientState
	cleanupEvery time.Duration
	mu           sync.Mutex
}

func NewRateLimiter() *RateLimiter {
	interval := 3 * time.Minute
	rl := &RateLimiter{
		clients:      make(map[string]*clientState),
		cleanupEvery: interval,
		nextCleanup:  time.Now().Add(interval),
	}

	go rl.cleanupLoop()

	return rl
}

func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.cleanupEvery)
	defer ticker.Stop()

	for now := range ticker.C {
		rl.mu.Lock()
		for key, cl := range rl.clients {
			if cl.blocked {
				cl.blocked = false
				cl.limiter = rate.NewLimiter(5, 10)
				cl.lastSeen = now
				continue
			}
			if now.Sub(cl.lastSeen) > rl.cleanupEvery {
				delete(rl.clients, key)
			}
		}
		rl.nextCleanup = now.Add(rl.cleanupEvery)
		rl.mu.Unlock()
	}
}

func makeClientKey(c *gin.Context) string {
	ip := c.ClientIP()
	ua := strings.TrimSpace(c.GetHeader("User-Agent"))
	if ua == "" {
		ua = "unknown"
	}
	return ip + "|" + ua
}

func (rl *RateLimiter) getClient(key string) *clientState {
	cl, exists := rl.clients[key]
	if !exists {
		cl = &clientState{
			limiter:  rate.NewLimiter(5, 10),
			lastSeen: time.Now(),
		}
		rl.clients[key] = cl
		return cl
	}
	if !cl.blocked {
		cl.lastSeen = time.Now()
	}
	return cl
}

func (rl *RateLimiter) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := makeClientKey(c)

		rl.mu.Lock()

		cl := rl.getClient(key)
		retryAfter := max(int(time.Until(rl.nextCleanup).Seconds()), 1)

		if cl.blocked {
			rl.mu.Unlock()
			c.Header("Retry-After", strconv.Itoa(retryAfter))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "temporarily blocked",
				"retry-after": retryAfter,
			})
			c.Abort()
			return
		}

		if !cl.limiter.Allow() {
			cl.blocked = true
			rl.mu.Unlock()

			c.Header("Retry-After", strconv.Itoa(retryAfter))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "rate limit exceeded",
				"retry-after": retryAfter,
			})
			c.Abort()
			return
		}
		rl.mu.Unlock()
		c.Next()
	}
}
