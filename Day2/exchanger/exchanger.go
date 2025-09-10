package exchanger

import (
	"context"
	"time"
)

func fakeRequest(from, to string) float64 {
	delay := 1 * time.Second
	time.Sleep(delay)

	if from == "USD" && to == "EUR" {
		return 0.85
	}
	if from == "EUR" && to == "USD" {
		return 1.17
	}
	return 1.0
}

func FetchRate(ctx context.Context, from, to string) (float64, error) {
	resultChan := make(chan float64)

	go func() {
		rate := fakeRequest(from, to)
		resultChan <- rate
	}()

	select {
	case rate := <-resultChan:
		return rate, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}

}
