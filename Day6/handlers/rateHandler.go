package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"lenrek88/config"
	"lenrek88/exchanger"
	"lenrek88/logger"
	"net/http"
	"time"
)

type Rate struct {
	exchangeService *exchanger.ExchangeService
}

func NewRateHandler(service *exchanger.ExchangeService) *Rate {
	return &Rate{exchangeService: service}
}

func (h *Rate) RateHandler(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*config.AppConfig.Timeout)
	defer cancel()

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	if from == "" || to == "" {
		err := fmt.Errorf("handler: missing required parameters 'from' or 'to'")
		logger.Error("RateHandler error", err)
		http.Error(w, "Missing parameters: from and to required", http.StatusBadRequest)
		return
	}

	rate, err := h.exchangeService.FetchRate(ctx, from, to)
	if err != nil {
		err = fmt.Errorf("handler: failed to fetch rate for %s to %s: %w", from, to, err)
		logger.Error("RateHandler error", err)
		http.Error(w, "Failed to get exchange rate", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(rate)
	if err != nil {
		err = fmt.Errorf("handler: failed to marshal response: %w", err)
		logger.Error("RateHandler error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	logger.Info("rate : from " + from + ", to " + to + ", rate " + string(data))
	w.Write(data)

}
