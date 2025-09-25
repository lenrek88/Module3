package main

import (
	"lenrek88/config"
	"lenrek88/handlers"
	"lenrek88/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.LoadConfig("config.json"); err != nil {
		panic(err)
	}
	r := gin.Default()
	limiter := middleware.NewClientLimiter(5, time.Minute)
	cache := middleware.NewCacheMiddleware(30 * time.Second)
	r.Use(limiter.Middleware(), cache.Middleware())

	r.GET("/today", handlers.TodayHandler)
	r.GET("/weekly", handlers.WeeklyHandler)
	if err := r.Run(config.AppConfig.Port); err != nil {
		panic(err)
	}

}
