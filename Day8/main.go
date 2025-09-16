package main

import (
	"lenrek88/config"
	"lenrek88/exchanger"
	"lenrek88/handlers"
	"lenrek88/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := logger.Init(); err != nil {
		panic(err)
	}
	if err := config.LoadConfig("config.json"); err != nil {
		logger.Error("Failed to load config", err)
		panic(err)
	}
	apiClient := exchanger.NewAPIClient(config.AppConfig.APIBaseURL)
	exchangeService := exchanger.NewExchangeService(apiClient)

	rateHandler := handlers.NewRateHandler(exchangeService).RateHandler
	exchangeHandler := handlers.NewExchangeHandler(exchangeService).ExchangeHandler

	r := gin.Default()
	r.GET("/rate", rateHandler)
	r.GET("/exchange", exchangeHandler)
	r.GET("/stats", handlers.StatsHandler)
	if err := r.Run(":" + config.AppConfig.Port); err != nil {
		logger.Error("Failed to start server", err)
		panic(err)
	}

}
