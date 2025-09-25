package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Config struct {
	APIBaseURL map[string]string `json:"api_base_url"`
	Appid      string            `json:"appid"`
	Timeout    time.Duration     `json:"timeout_s"`
	Port       string            `json:"port"`
}

var AppConfig Config

func LoadConfig(fileName string) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	err = json.Unmarshal(data, &AppConfig)
	if err != nil {
		return fmt.Errorf("error parsing config: %w", err)
	}

	return nil
}
