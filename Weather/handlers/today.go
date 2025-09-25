package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"lenrek88/api"
	"lenrek88/config"
	"net/http"
)

func TodayHandler(c *gin.Context) {
	url := config.AppConfig.APIBaseURL["today"]

	body, err := api.Fetch(c, url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var weatherResp TodayResponse
	if err := json.Unmarshal(body, &weatherResp); err != nil {
		fmt.Println(err)
		panic(err)
	}

	date := time.Unix(int64(weatherResp.Dt), 0).Format("2006-01-02 15:04:05")

	simplifiedResponse := SimplifiedResponse{
		City: weatherResp.Name,
	}

	simplifiedForecast := SimplifiedForecast{
		DateTime:    date,
		Temperature: weatherResp.Main.Temp,
		FeelsLike:   weatherResp.Main.FeelsLike,
		Humidity:    weatherResp.Main.Humidity,
		WindSpeed:   weatherResp.Wind.Speed,
	}
	if len(weatherResp.Weather) > 0 {
		simplifiedForecast.Description = weatherResp.Weather[0].Description
	}

	simplifiedResponse.Forecasts = append(simplifiedResponse.Forecasts, simplifiedForecast)

	c.JSON(http.StatusOK, simplifiedResponse)

}
