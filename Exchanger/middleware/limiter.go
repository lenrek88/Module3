package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type ClientLimiter struct {
	clients  map[string]*ClientData
	mu       sync.Mutex
	limit    int
	interval time.Duration
}

type ClientData struct {
	count     int
	lastReset time.Time
}

func NewClientLimiter(limit int, time time.Duration) *ClientLimiter {
	return &ClientLimiter{clients: make(map[string]*ClientData), limit: limit, interval: time}
}

func (l *ClientLimiter) reset(clientIP string) {
	l.clients[clientIP] = &ClientData{count: 0, lastReset: time.Now()}
}

func (l *ClientLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		l.mu.Lock()
		defer l.mu.Unlock()
		now := time.Now()
		data, exists := l.clients[clientIP]
		if !exists {
			l.reset(clientIP)
			c.Next()
			return
		}
		if now.Sub(data.lastReset) > l.interval {
			l.reset(clientIP)
			c.Next()
			return
		}
		if data.count < l.limit {
			l.clients[clientIP].count++
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests, try again later"})
	}
}
