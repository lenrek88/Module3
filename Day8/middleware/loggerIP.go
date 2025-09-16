package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type loggerIP struct {
	filename string
}

func NewLoggerIP(filename string) *loggerIP {
	return &loggerIP{filename: filename}
}

func (l *loggerIP) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := os.OpenFile(l.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "couldn't open the file"})
		}
		IPLogger := log.New(file, "IP_ADDR: ", log.Ldate|log.Ltime|log.Lshortfile)
		msg := c.ClientIP() + c.Request.RemoteAddr
		IPLogger.Output(2, msg)
		c.Next()
		return
	}
}
