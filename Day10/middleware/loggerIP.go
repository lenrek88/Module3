package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type LoggerIP struct {
	filename string
}

func NewLoggerIP(filename string) *LoggerIP {
	return &LoggerIP{filename: filename}
}

func (l *LoggerIP) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := os.OpenFile(l.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "couldn't open the file"})
		}
		IPLogger := log.New(file, "IP_ADDR: ", log.Ldate|log.Ltime|log.Lshortfile)
		msg := "IP address : " + c.ClientIP() + ", URL: " + c.Request.URL.String()
		IPLogger.Output(2, msg)
		c.Next()
		return
	}
}
