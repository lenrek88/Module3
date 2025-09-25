package exchanger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Rate float64
type apiResponse map[string]any
type APIClient struct {
	baseURL map[string]string
}

func NewAPIClient(baseURL map[string]string) *APIClient {
	return &APIClient{baseURL}
}

func (c *APIClient) fetchDevRates(ctx context.Context, from, to string) (Rate, error) {

	value, ok := c.baseURL["dev"]
	var url string
	if ok {
		url = value + strings.ToLower(from) + ".json"
	} else {
		return 0, fmt.Errorf("api client: dev the key is missing  %s: %w", from, "dev")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return 0, fmt.Errorf("api client: failed to create request for %s: %w", from, err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return 0, fmt.Errorf("api client: failed to execute request for %s: %w", from, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("api client: API returned non-200 status for %s: %s", from, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("api client: failed to read response for %s: %w", from, err)
	}

	var res map[string]any
	if err := json.Unmarshal(body, &res); err != nil {
		return 0, fmt.Errorf("api client: failed to parse JSON for %s: %w", from, err)
	}

	rates, ok := res[from].(map[string]any)
	if !ok {
		return 0, fmt.Errorf("api client: invalid response format for currency %s", from)
	}

	rate, ok := rates[to].(float64)
	if !ok {
		return 0, fmt.Errorf("api client: rate not found for currency pair %s to %s", from, to)
	}

	return Rate(rate), nil

}

func (c *APIClient) fetchCbrRates(ctx context.Context, from, to string) (Rate, error) {

	from = strings.ToUpper(from)
	to = strings.ToUpper(to)

	url, ok := c.baseURL["cbr"]
	if !ok {
		return 0, fmt.Errorf("fetchCbrRates: dev the key is missing  %s: %w", from, "dev")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return 0, fmt.Errorf("fetchCbrRates: failed to create request for %s: %w", from, err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return 0, fmt.Errorf("fetchCbrRates: failed to execute request for %s: %w", from, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("fetchCbrRates: API returned non-200 status for %s: %s", from, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("fetchCbrRates: failed to read response for %s: %w", from, err)
	}

	var res map[string]any
	if err := json.Unmarshal(body, &res); err != nil {
		return 0, fmt.Errorf("fetchCbrRates: failed to parse JSON for %s: %w", from, err)
	}
	rates, ok := res["Valute"].(map[string]any)

	if from == "RUB" && to == "RUB" {
		return Rate(1), nil
	}
	if to == "RUB" {
		rateTo, ok := rates[from].(map[string]any)
		rateToValue, ok := rateTo["Value"].(float64)
		if !ok {
			return 0, fmt.Errorf("fetchCbrRates: invalid response format for %s to %s", from, to)
		}
		return Rate(rateToValue), nil
	}
	if from == "RUB" {
		rateFrom, ok := rates[to].(map[string]any)
		rateFromValue, ok := rateFrom["Value"].(float64)
		if !ok {
			return 0, fmt.Errorf("fetchCbrRates: invalid response format for %s to %s", from, to)
		}
		return Rate(1 / rateFromValue), nil
	}
	rateFrom, ok := rates[from].(map[string]any)
	rateTo, ok := rates[to].(map[string]any)

	rateFromValue, ok := rateFrom["Value"].(float64)
	if !ok {
		return 0, fmt.Errorf("fetchCbrRates: invalid response format for currency %s", from)
	}
	rateToValue, ok := rateTo["Value"].(float64)
	if !ok {
		return 0, fmt.Errorf("fetchCbrRates: invalid response format for currency %s", from)
	}

	rate := rateFromValue / rateToValue
	return Rate(rate), nil
}
