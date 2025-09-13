package main

import (
	"context"
	"encoding/json"
	"fmt"
	"lenrek88/exchanger"
	"lenrek88/logger"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Config struct {
	APIBaseURL string        `json:"api_base_url"`
	Timeout    time.Duration `json:"timeout"`
}

var appConfig Config

func loadConfig() error {
	data, err := os.ReadFile("config.json")
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	err = json.Unmarshal(data, &appConfig)
	if err != nil {
		return fmt.Errorf("error parsing config: %w", err)
	}

	return nil
}

func rateHandler(w http.ResponseWriter, r *http.Request) {
	err := loadConfig()
	if err != nil {
		fmt.Println(err)
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*appConfig.Timeout)
	defer cancel()

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	rate, err := exchanger.FetchRate(appConfig.APIBaseURL, ctx, from, to)
	if err != nil {
		context.WithCancel(ctx)
		fmt.Println("GET : ", err)
		return
	}

	data, _ := json.Marshal(rate)
	w.Write(data)

}

func exchangeHandler(w http.ResponseWriter, r *http.Request) {

	err := loadConfig()
	if err != nil {
		fmt.Println(err)
	}
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*appConfig.Timeout)
	defer cancel()

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	amount := r.URL.Query().Get("amount")
	rate, err := exchanger.FetchRate(appConfig.APIBaseURL, ctx, from, to)
	if err != nil {
		context.WithCancel(ctx)
		return
	}
	fmt.Println(rate)
	rateFloat := float64(rate)
	amountFloat, _ := strconv.ParseFloat(amount, 64)

	data, _ := json.Marshal(rateFloat * amountFloat)
	w.Write(data)

}

func main() {
	e := logger.Init()
	if e != nil {
		fmt.Println(e)
	}
	http.HandleFunc("/rate", rateHandler)
	http.HandleFunc("/exchange", exchangeHandler)
	http.ListenAndServe(":8080", nil)
}
