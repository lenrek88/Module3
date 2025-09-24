package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/weather", weatherHandler)
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}

func weatherHandler(c *gin.Context) {
	url := "https://api.openweathermap.org/data/2.5/forecast?q=moscow&appid=85b3329d911fdb03a5450fbb580353ab&units=metric"
	req, err := http.NewRequestWithContext(c, http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var res map[string]any
	if err := json.Unmarshal(body, &res); err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, res)
}
