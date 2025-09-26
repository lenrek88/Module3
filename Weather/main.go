package main

import (
	"flag"
	"fmt"
	"lenrek88/cmd"
	"lenrek88/config"
	"lenrek88/docs"
	"lenrek88/handlers"
	"lenrek88/middleware"
	"lenrek88/weather"
	"time"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// общие анотации для всего сервера
// @title Документация – заголовок всей доки
// @version 1.0 – версия АПИ
// @description API для погоды – текстовое описание АПИ
// @host localhost:8080 – адрес для запросов
// @BasePath / – базовый адрес (может быть /api или /api/v1)
func main() {
	if err := config.LoadConfig("config.json"); err != nil {
		panic(err)
	}

	cliMode := flag.Bool("cli", false, "Run as CLI")
	flag.Parse()

	if *cliMode {
		fmt.Println("Погодный информер")
		cache := middleware.NewCacheMiddleware(30 * time.Second)

		w := weather.NewInformer(cache)
		cli := cmd.NewCmd(w)
		cli.Run()
	} else {
		if err := config.LoadConfig("config.json"); err != nil {
			panic(err)
		}
		r := gin.Default()
		limiter := middleware.NewClientLimiter(100, time.Minute)
		cache := middleware.NewCacheMiddleware(30 * time.Second)
		r.Use(limiter.Middleware(), cache.Middleware(), gin.Logger(), middleware.Recovery())

		docs.SwaggerInfo.BasePath = "/"
		r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

		r.GET("/today", handlers.TodayHandler)
		r.GET("/weekly", handlers.WeeklyHandler)

		if err := r.Run(config.AppConfig.Port); err != nil {
			panic(err)
		}
	}

}
