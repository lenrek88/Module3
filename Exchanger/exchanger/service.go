package exchanger

import (
	"context"
	"fmt"
	"sync"
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

	type Exchanger struct {
		code string
		fn   func(context.Context, string, string) (Rate, error)
	}

	exchanges := []Exchanger{
		{code: "cbr",
			fn: s.apiClient.fetchCbrRates},
		{code: "dev",
			fn: s.apiClient.fetchDevRates},
	}

	resultChan := make(chan result)
	var res result
	res.rates = make(map[string]Rate)

	var wg sync.WaitGroup
	wg.Add(len(exchanges))

	for _, ex := range exchanges {
		go func(ex Exchanger) {

			rate, err := ex.fn(ctx, from, to)
			if err != nil {
				res.err = fmt.Errorf("fetch %s rate: %w", ex.code, err)
			} else {
				res.rates[ex.code] = rate
			}
			wg.Done()
		}(ex)
	}
	go func() {
		wg.Wait()
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
