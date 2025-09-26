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

// TodayHandler - имя хэндлера
// TodayHandler
// @Summary today handle
// @Param q query string true "город \n Example: Kazan" moscow
// @param lang query string true "локализация \n Example: ru" ru
// @param unit query string true "шкала градусов \n Example: Standard,Imperial,Metric"
// @Success 200 {string} rate
// @Router /today [get]
func TodayHandler(c *gin.Context) {
	unit := c.DefaultQuery("unit", "metric")
	lang := c.DefaultQuery("lang", "ru")
	appid := c.DefaultQuery("appid", config.AppConfig.Appid)
	q := c.DefaultQuery("q", "Moscow")
	url := config.AppConfig.APIBaseURL["today"]
	fullURL := fmt.Sprintf("%s?q=%s&units=%s&lang=%s&appid=%s",
		url, q, unit, lang, appid)
	body, err := api.Fetch(c, fullURL)
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
