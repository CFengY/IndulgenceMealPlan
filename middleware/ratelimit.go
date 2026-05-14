package middleware

import (
	"IndulgenceMealPlan/global"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type rateLimiterStore struct {
	ips map[string]*ipLimiter
	mu  sync.RWMutex
}

var rlStore *rateLimiterStore

func (s *rateLimiterStore) getLimiter(ip string) *rate.Limiter {
	s.mu.RLock()
	entry, exists := s.ips[ip]
	s.mu.RUnlock()

	if exists {
		entry.lastSeen = time.Now()
		return entry.limiter
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	entry, exists = s.ips[ip]
	if !exists {
		limiter := rate.NewLimiter(
			rate.Limit(global.Config.RateLimit.Rate),
			global.Config.RateLimit.Burst,
		)
		s.ips[ip] = &ipLimiter{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}

	entry.lastSeen = time.Now()
	return entry.limiter
}

func (s *rateLimiterStore) cleanup(intervalSec int) {
	ticker := time.NewTicker(time.Duration(intervalSec) * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		s.mu.Lock()
		for ip, entry := range s.ips {
			if time.Since(entry.lastSeen) > time.Duration(intervalSec)*time.Second {
				delete(s.ips, ip)
			}
		}
		s.mu.Unlock()
	}
}

func RateLimitMiddleware() gin.HandlerFunc {
	if !global.Config.RateLimit.Enabled {
		return func(c *gin.Context) { c.Next() }
	}

	store := &rateLimiterStore{
		ips: make(map[string]*ipLimiter),
	}
	go store.cleanup(global.Config.RateLimit.CleanupInterval)
	rlStore = store

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := rlStore.getLimiter(ip)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
