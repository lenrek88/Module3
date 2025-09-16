package middleware

import (
	"bytes"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type CacheItem struct {
	data      []byte
	timestamp time.Time
}

type CacheMiddleware struct {
	cache map[string]*CacheItem
	mu    sync.Mutex
	ttl   time.Duration
}

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func NewCacheMiddleware(ttl time.Duration) *CacheMiddleware {
	return &CacheMiddleware{cache: make(map[string]*CacheItem), ttl: ttl}
}

func (cm *CacheMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Request.URL.RequestURI()
		cm.mu.Lock()
		defer cm.mu.Unlock()
		item, exists := cm.cache[key]

		if exists && time.Since(item.timestamp) < cm.ttl {
			c.Data(http.StatusOK, "application/json", item.data)
			c.Abort()
			return
		}
		writer := &bodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = writer
		c.Next()
		cm.cache[key] = &CacheItem{
			data:      writer.body.Bytes(),
			timestamp: time.Now(),
		}

	}
}
