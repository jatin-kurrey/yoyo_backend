package middleware

import (
	"net/http"
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
