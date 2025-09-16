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

func (s *ExchangeService) FetchRate(ctx context.Context, from, to string) (map[string]Rate, error) {

	if from == "" || to == "" {
		return nil, fmt.Errorf("service: currency codes cannot be empty")
	}

	type result struct {
		rates map[string]Rate
		err   error
	}

	resultChan := make(chan result)

	go func() {
		var res result
		res.rates = make(map[string]Rate)
		devRate, err := s.apiClient.fetchDevRates(ctx, from, to)
		if err != nil {
			res.err = fmt.Errorf("fetch dev rate: %w", err)
		} else {
			res.rates["dev"] = devRate
		}
		cbrRate, err := s.apiClient.fetchCbrRates(ctx, from, to)
		if err != nil {
			if res.err != nil {
				res.err = fmt.Errorf("dev failed, cbr also falied: %w", err)
			} else {
				res.err = fmt.Errorf("cbr fetch failed: %w", err)
			}
		} else {
			res.rates["cbr"] = cbrRate
		}

		if res.rates["dev"] != 0 || res.rates["cbr"] != 0 {
			res.err = nil
		}
		resultChan <- res
	}()

	select {
	case res := <-resultChan:
		return res.rates, res.err
	case <-ctx.Done():
		return nil, fmt.Errorf("FetchRate: timeout during first fetch: %w", ctx.Err())
	}

}

func (s *ExchangeService) ConvertAmount(rate Rate, amount float64) float64 {
	return float64(rate) * amount
}
