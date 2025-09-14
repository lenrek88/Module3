package exchanger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Rate float64
type apiResponse map[string]any
type APIClient struct {
	baseURL string
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{baseURL: baseURL}
}

func (c *APIClient) MakeRequest(ctx context.Context, from, to string) (Rate, error) {
	url := c.baseURL + from + ".json"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil) // запрос с контекстом

	if err != nil {
		return 0, fmt.Errorf("api client: failed to create request for %s: %w", from, err)
	}

	delay := time.Second * 1
	time.Sleep(delay)

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
