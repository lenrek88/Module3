package handlers

import (
	"fmt"
	"lenrek88/logger"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func StatsHandler(c *gin.Context) {

	data, err := os.ReadFile("./app.log")
	if err != nil {
		err := fmt.Errorf("error reading app.log: %s", err.Error())
		logger.Error("StatsHandler error", err)
		return
	}
	dataLines := strings.Split(string(data), "\n")

	var amountRateContains int64
	var amountExchangeContains int64
	start := c.Query("start")
	end := c.Query("end")
	if start != "" || end != "" {

		timeLayout := "2006-01-02"
		startDateParse, err := time.Parse(timeLayout, start)
		if err != nil {
			err := fmt.Errorf("StatsHandler: time parse error: %s", err.Error())
			logger.Error("StatsHandler error", err)
			return
		}

		endDateParse, err := time.Parse(timeLayout, end)
		if err != nil {
			err := fmt.Errorf("StatsHandler: time parse error: %s", err.Error())
			logger.Error("StatsHandler error", err)
			return
		}
		for _, line := range dataLines {
			parts := strings.Split(line, " ")
			if len(parts) >= 2 {
				dateString := parts[1]
				parsedTime, err := time.Parse("2006/01/02", dateString)
				if err != nil {
					err := fmt.Errorf("StatsHandler: date parse error: %s", err.Error())
					logger.Error("StatsHandler error", err)
					return
				}
				newFormat := parsedTime.Format(timeLayout)
				timeParseNewFormat, _ := time.Parse(timeLayout, newFormat)
				fmt.Println(timeParseNewFormat, startDateParse, endDateParse)
				if (timeParseNewFormat.After(startDateParse) || timeParseNewFormat.Equal(startDateParse)) && (timeParseNewFormat.Before(endDateParse) || (timeParseNewFormat.Equal(endDateParse))) {
					if strings.Contains(line, "rate") {
						amountRateContains = amountRateContains + 1
					}
					if strings.Contains(line, "exchange") {
						amountExchangeContains = amountExchangeContains + 1
					}
				}
			}

		}
	} else {
		for _, line := range dataLines {
			if strings.Contains(line, "rate") {
				amountRateContains = amountRateContains + 1
			}
			if strings.Contains(line, "exchange") {
				amountExchangeContains = amountExchangeContains + 1
			}
		}
	}

	strRate := fmt.Sprintf("%d", amountRateContains)
	strExchange := fmt.Sprintf("%d", amountExchangeContains)
	if start == "" || end == "" {
		start = "2025-01-02"
		end = "текущее время"
	}
	t := "на период c " + start + " по " + end + " были выполнены запросы : " + "rate : " + strRate + " exchange : " + strExchange

	c.JSON(http.StatusOK, gin.H{"rate": t})

}
