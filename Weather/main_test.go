package main

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(t *testing.T) {
	var context *gin.Context
	url := "https://api.openweathermap.org/data/2.5/forecast?q=moscow&appid=85b3329d911fdb03a5450fbb580353ab&units=metric"
	resp, err := apiFetch(context, url)
	if err != nil {
		panic(err)
	}
}
