package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"lenrek88/exchanger"
	"lenrek88/logger"
	"net/http"
)

type Rate struct {
	exchangeService *exchanger.ExchangeService
}

func NewRateHandler(service *exchanger.ExchangeService) *Rate {
	return &Rate{exchangeService: service}
}

func (h *Rate) RateHandler(c *gin.Context) {

	from := c.Query("from")
	to := c.Query("to")

	if from == "" || to == "" {
		err := fmt.Errorf("handler: missing required parameters 'from' or 'to'")
		logger.Error("RateHandler error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "both from and to params are required"})
		return
	}

	rate, err := h.exchangeService.FetchRate(c.Request.Context(), from, to)
	if len(rate) >= 1 {
		c.JSON(http.StatusOK, gin.H{"rate": rate})
	}
	if err != nil && len(rate) == 0 {
		err = fmt.Errorf("handler: failed to fetch rate for %s to %s: %w", from, to, err)
		logger.Error("RateHandler error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "data not available"})
		return
	}

	data, err := json.Marshal(rate)
	logger.Info("rate : from " + from + ", to " + to + ", rate " + string(data))

}
