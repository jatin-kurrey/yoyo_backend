package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type bucket struct {
	count int
	reset time.Time
}

func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	if limit <= 0 {
		limit = 60
	}

	var mu sync.Mutex
	buckets := map[string]*bucket{}

	// Periodic cleanup to prevent memory leak
	go func() {
		for {
			time.Sleep(10 * time.Minute)
			mu.Lock()
			now := time.Now()
			for ip, b := range buckets {
				if now.After(b.reset) {
					delete(buckets, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		now := time.Now()
		key := c.ClientIP()

		mu.Lock()
		item, exists := buckets[key]
		if !exists || now.After(item.reset) {
			item = &bucket{count: 0, reset: now.Add(window)}
			buckets[key] = item
		}
		item.count++
		remaining := limit - item.count
		mu.Unlock()

		c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(max(0, remaining)))
		c.Header("X-RateLimit-Reset", strconv.Itoa(int(time.Until(item.reset).Seconds())))

		if remaining < 0 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "Too many requests. Please try again later.",
			})
			return
		}

		c.Next()
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
