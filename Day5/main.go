package main

import (
	"lenrek88/config"
	"lenrek88/exchanger"
	"lenrek88/handlers"
	"lenrek88/logger"
	"net/http"
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
	http.HandleFunc("/rate", rateHandler)
	http.HandleFunc("/exchange", exchangeHandler)
	http.HandleFunc("/stats", handlers.StatsHandler)
	logger.Info("Server starting on port" + " : " + config.AppConfig.Port)
	if err := http.ListenAndServe(":"+config.AppConfig.Port, nil); err != nil {
		logger.Error("Failed to start server", err)
		panic(err)
	}
}
