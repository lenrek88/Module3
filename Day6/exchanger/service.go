package exchanger

import (
	"context"
	"fmt"
)

type ExchangeService struct {
	apiClient *APIClient
}

func NewExchangeService(apiClient *APIClient) *ExchangeService {
	return &ExchangeService{apiClient: apiClient}
}

func (s *ExchangeService) FetchRate(ctx context.Context, from, to string) (Rate, error) {

	if from == "" || to == "" {
		return 0, fmt.Errorf("service: currency codes cannot be empty")
	}

	resultChan := make(chan Rate)
	errorChan := make(chan error)

	go func() {
		rate, err := s.apiClient.MakeRequest(ctx, from, to)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- rate
	}()

	select {
	case rate := <-resultChan:
		return rate, nil
	case err := <-errorChan:
		return 0, fmt.Errorf("service: API request failed: %w", err)
	case <-ctx.Done():
		return 0, fmt.Errorf("service: request timeout for %s to %s: %w", from, to, ctx.Err())
	}
}

func (s *ExchangeService) ConvertAmount(rate Rate, amount float64) float64 {
	return float64(rate) * amount
}
