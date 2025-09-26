package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"lenrek88/api"
	"lenrek88/config"
	"lenrek88/handlers"
	"lenrek88/middleware"
	"time"
)

type Informer struct {
	informerMap handlers.SimplifiedResponse
	cache       *middleware.CacheMiddleware
}

func NewInformer(c *middleware.CacheMiddleware) *Informer {
	return &Informer{informerMap: handlers.SimplifiedResponse{}, cache: c}
}

func (in *Informer) WeeklyHandler(city, unit, lang string) (handlers.SimplifiedResponse, error) {
	appid := config.AppConfig.Appid
	url := config.AppConfig.APIBaseURL["forecast"]
	if url == "" {
		return handlers.SimplifiedResponse{}, fmt.Errorf("API base URL is empty")
	}
	if appid == "" {
		return handlers.SimplifiedResponse{}, fmt.Errorf("API key (appid) is missing")
	}

	fullURL := fmt.Sprintf("%s?q=%s&units=%s&lang=%s&appid=%s",
		url, city, unit, lang, appid)

	keyCache := fullURL
	in.cache.Mu.Lock()
	defer in.cache.Mu.Unlock()

	item, exists := in.cache.Cache[keyCache]
	var body []byte
	if exists && time.Since(item.Timestamp) < in.cache.Ttl {
		body = item.Data
	} else {
		ctx := context.Background()
		body, err := api.Fetch(ctx, fullURL)
		if err != nil {
			return handlers.SimplifiedResponse{}, err
		}
		in.cache.Cache[keyCache] = &middleware.CacheItem{
			Data:      body,
			Timestamp: time.Now(),
		}
	}

	var weatherResp handlers.WeeklyResponse
	if err := json.Unmarshal(body, &weatherResp); err != nil {
		return handlers.SimplifiedResponse{}, err
	}

	simplifiedResponse := handlers.SimplifiedResponse{
		City: weatherResp.City.Name,
	}

	for _, forecast := range weatherResp.List {
		simplifiedForecast := handlers.SimplifiedForecast{
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
	return simplifiedResponse, nil
}

func (in *Informer) TodayHandler(city, unit, lang string) (handlers.SimplifiedResponse, error) {
	appid := config.AppConfig.Appid
	url := config.AppConfig.APIBaseURL["today"]
	fullURL := fmt.Sprintf("%s?q=%s&units=%s&lang=%s&appid=%s",
		url, city, unit, lang, appid)

	keyCache := fullURL
	in.cache.Mu.Lock()
	defer in.cache.Mu.Unlock()

	item, exists := in.cache.Cache[keyCache]
	var body []byte
	var err error
	if exists && time.Since(item.Timestamp) < in.cache.Ttl {
		body = item.Data
	} else {
		ctx := context.Background()
		body, err = api.Fetch(ctx, fullURL)
		if err != nil {
			return handlers.SimplifiedResponse{}, err
		}
		in.cache.Cache[keyCache] = &middleware.CacheItem{
			Data:      body,
			Timestamp: time.Now(),
		}
	}

	var weatherResp handlers.TodayResponse
	if err := json.Unmarshal(body, &weatherResp); err != nil {
		return handlers.SimplifiedResponse{}, err
	}

	date := time.Unix(int64(weatherResp.Dt), 0).Format("2006-01-02 15:04:05")

	simplifiedResponse := handlers.SimplifiedResponse{
		City: weatherResp.Name,
	}

	simplifiedForecast := handlers.SimplifiedForecast{
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

	return simplifiedResponse, nil
}
