package handlers

import (
	"encoding/json"
	"fmt"
	"lenrek88/api"
	"lenrek88/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// WeeklyHandler - имя хэндлера
// WeeklyHandler
// @Summary Weekly handle
// @Param q query string true "город \n Example: Kazan" moscow
// @param lang query string true "локализация \n Example: ru" ru
// @param unit query string true "шкала градусов \n Example: Standard,Imperial,Metric"
// @Success 200 {string} rate
// @Router /weekly [get]
func WeeklyHandler(c *gin.Context) {

	unit := c.DefaultQuery("unit", "metric")
	lang := c.DefaultQuery("lang", "ru")
	appid := c.DefaultQuery("appid", config.AppConfig.Appid)
	q := c.DefaultQuery("q", "Moscow")
	url := config.AppConfig.APIBaseURL["forecast"]
	fullURL := fmt.Sprintf("%s?q=%s&units=%s&lang=%s&appid=%s",
		url, q, unit, lang, appid)
	body, err := api.Fetch(c, fullURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var weatherResp WeeklyResponse
	if err := json.Unmarshal(body, &weatherResp); err != nil {
		fmt.Println(err)
		panic(err)
	}

	simplifiedResponse := SimplifiedResponse{
		City: weatherResp.City.Name,
	}

	for _, forecast := range weatherResp.List {
		simplifiedForecast := SimplifiedForecast{
			DateTime:     forecast.DtTxt,
			Temperature:  forecast.Main.Temp,
			FeelsLike:    forecast.Main.FeelsLike,
			Humidity:     forecast.Main.Humidity,
			WindSpeed:    forecast.Wind.Speed,
			ChanceOfRain: forecast.Pop * 100,
		}
		if len(forecast.Weather) > 0 {
			simplifiedForecast.Description = forecast.Weather[0].Description
		}

		simplifiedResponse.Forecasts = append(simplifiedResponse.Forecasts, simplifiedForecast)
	}

	c.JSON(http.StatusOK, simplifiedResponse)

}
