package handlers

import (
	"lenrek88/exchanger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
		return
	}

	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil || amountFloat <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount parameter"})
		//http.Error(w, "Invalid amount parameter", http.StatusBadRequest)
		return
	}

	if amountFloat <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be positive"})
		//err := fmt.Errorf("handler: amount must be positive, got %f", amountFloat)
		//logger.Error("ExchangeHandler error", err)
		return
	}
	rate, err := h.exchangeService.FetchRate(c.Request.Context(), from, to)
	if err != nil {
		//err = fmt.Errorf("handler: failed to fetch rate for %s to %s: %w", from, to, err)
		//logger.Error("ExchangeHandler error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get exchange rate"})
		return
	}
	result := h.exchangeService.ConvertAmount(rate["dev"], amountFloat)

	c.JSON(http.StatusOK, gin.H{"exchange": result})

	//logger.Info("exchange : from " + from + ", to " + to + ", amount " + amount + ", exchange " + string(data))

}
