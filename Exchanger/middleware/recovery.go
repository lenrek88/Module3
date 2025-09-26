package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("painc recovered: %v\nstack=%s method=%s path=%s", r, debug.Stack(), c.Request.Method, c.FullPath())
				c.AbortWithStatusJSON(
					http.StatusInternalServerError,
					gin.H{"error": "internal server error"},
				)
			}
		}()
		c.Next()
	}
}
