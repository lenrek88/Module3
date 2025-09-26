package middleware

import (
	"bytes"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type CacheItem struct {
	Data      []byte
	Timestamp time.Time
}

type CacheMiddleware struct {
	Cache map[string]*CacheItem
	Mu    sync.Mutex
	Ttl   time.Duration
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
	return &CacheMiddleware{Cache: make(map[string]*CacheItem), Ttl: ttl}
}

func (cm *CacheMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Request.URL.RequestURI()
		cm.Mu.Lock()
		defer cm.Mu.Unlock()
		item, exists := cm.Cache[key]

		if exists && time.Since(item.Timestamp) < cm.Ttl {
			c.Data(http.StatusOK, "application/json", item.Data)
			c.Abort()
			return
		}
		writer := &bodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = writer
		c.Next()
		cm.Cache[key] = &CacheItem{
			Data:      writer.body.Bytes(),
			Timestamp: time.Now(),
		}

	}
}
