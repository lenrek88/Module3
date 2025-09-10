package main

import (
	"context"
	"encoding/json"
	"fmt"
	"lenrek88/exchanger"
	"net/http"
	"strconv"
	"time"
)

func rateHandler(w http.ResponseWriter, r *http.Request) {

	delay := 3 * time.Second
	ctx, cancel := context.WithTimeout(r.Context(), delay)
	defer cancel()

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	rate, err := exchanger.FetchRate(ctx, from, to)
	if err != nil {
		context.WithCancel(ctx)
		fmt.Println("GET : ", err)
		return
	}

	data, _ := json.Marshal(rate)
	w.Write(data)

}

func exchangeHandler(w http.ResponseWriter, r *http.Request) {

	delay := 3 * time.Second
	ctx, cancel := context.WithTimeout(r.Context(), delay)
	defer cancel()

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	amount := r.URL.Query().Get("amount")
	rate, err := exchanger.FetchRate(ctx, from, to)
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

	http.HandleFunc("/rate", rateHandler)
	http.HandleFunc("/exchange", exchangeHandler)
	http.ListenAndServe(":8080", nil)
}
