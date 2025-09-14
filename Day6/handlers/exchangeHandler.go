package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"lenrek88/config"
	"lenrek88/exchanger"
	"lenrek88/logger"
	"net/http"
	"strconv"
	"time"
)

type Exchange struct {
	exchangeService *exchanger.ExchangeService
}

func NewExchangeHandler(service *exchanger.ExchangeService) *Exchange {
	return &Exchange{exchangeService: service}
}

func (h *Exchange) ExchangeHandler(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*config.AppConfig.Timeout)
	defer cancel()

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	amount := r.URL.Query().Get("amount")

	if from == "" || to == "" {
		err := fmt.Errorf("handler: missing required parameters 'from' or 'to'")
		logger.Error("RateHandler error", err)
		http.Error(w, "Missing parameters: from and to required", http.StatusBadRequest)
		return
	}

	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil || amountFloat <= 0 {
		err = fmt.Errorf("handler: invalid amount parameter '%s': %w", amount, err)
		logger.Error("ExchangeHandler error", err)
		http.Error(w, "Invalid amount parameter", http.StatusBadRequest)
		return
	}

	if amountFloat <= 0 {
		err := fmt.Errorf("handler: amount must be positive, got %f", amountFloat)
		logger.Error("ExchangeHandler error", err)
		return
	}
	rate, err := h.exchangeService.FetchRate(ctx, from, to)
	if err != nil {
		err := fmt.Errorf("handler: conversion failed for %f %s to %s: %w", amountFloat, from, to)
		logger.Error("ExchangeHandler error", err)
		return
	}
	result := h.exchangeService.ConvertAmount(rate, amountFloat)

	data, err := json.Marshal(result)
	if err != nil {
		err = fmt.Errorf("handler: failed to marshal response: %w", err)
		logger.Error("ExchangeHandler error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Info("exchange : from " + from + ", to " + to + ", amount " + amount + ", exchange " + string(data))

	w.Write(data)
}
