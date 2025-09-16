package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"lenrek88/exchanger"
	"lenrek88/logger"
	"net/http"
	"strconv"
)

type Exchange struct {
	exchangeService *exchanger.ExchangeService
}

func NewExchangeHandler(service *exchanger.ExchangeService) *Exchange {
	return &Exchange{exchangeService: service}
}

func (h *Exchange) ExchangeHandler(c *gin.Context) {

	from := c.Query("from")
	to := c.Query("to")
	amount := c.Query("amount")

	if from == "" || to == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "both from and to params are required"})
		err := fmt.Errorf("handler: both from and to params are required")
		logger.Error("ExchangeHandler error", err)
		return
	}

	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil || amountFloat <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount parameter"})
		err := fmt.Errorf("handler: invalid amount parameter, got %f", amountFloat)
		logger.Error("ExchangeHandler error", err)
		return
	}

	if amountFloat <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be positive"})
		err := fmt.Errorf("handler: amount must be positive, got %f", amountFloat)
		logger.Error("ExchangeHandler error", err)
		return
	}
	rate, err := h.exchangeService.FetchRate(c.Request.Context(), from, to)
	if err != nil {
		err = fmt.Errorf("handler: failed to fetch rate for %s to %s: %w", from, to, err)
		logger.Error("ExchangeHandler error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get exchange rate"})
		return
	}
	result := h.exchangeService.ConvertAmount(rate["dev"], amountFloat)

	data := strconv.FormatFloat(result, 'f', -1, 64)
	logger.Info("exchange : from " + from + ", to " + to + ", amount " + amount + ", exchange " + data)

	c.JSON(http.StatusOK, gin.H{"exchange": result})

}
