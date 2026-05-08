package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type ipRateLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

const (
	limiterTTL      = 10 * time.Minute
	cleanupInterval = 1 * time.Minute
)

func RateLimitByIP(limit rate.Limit, burst int) gin.HandlerFunc {
	var (
		mu          sync.Mutex
		limiters    = make(map[string]*ipRateLimiter)
		lastCleanup time.Time
	)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		mu.Lock()
		if now.Sub(lastCleanup) >= cleanupInterval {
			for key, entry := range limiters {
				if now.Sub(entry.lastSeen) > limiterTTL {
					delete(limiters, key)
				}
			}
			lastCleanup = now
		}

		entry, exists := limiters[ip]
		if !exists {
			entry = &ipRateLimiter{
				limiter:  rate.NewLimiter(limit, burst),
				lastSeen: now,
			}
			limiters[ip] = entry
		} else {
			entry.lastSeen = now
		}
		limiter := entry.limiter
		mu.Unlock()

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			return
		}

		c.Next()
	}
}
