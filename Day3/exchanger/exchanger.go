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

func request(from, to string, ctx context.Context) Rate {

	url := "https://latest.currency-api.pages.dev/v1/currencies/" + from + ".json"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil) // запрос с контекстом

	if err != nil {
		fmt.Println(err)
	}
	delay := time.Second * 2
	time.Sleep(delay)
	client := &http.Client{}
	resp, err2 := client.Do(req)

	if err2 != nil {
		fmt.Println("ERROR TO REQUEST ERR2:", err2)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.Status)
	}
	body, _ := io.ReadAll(resp.Body)

	var res apiResponse
	json.Unmarshal(body, &res)
	rates, ok := res[from].(map[string]any)

	if !ok {
		fmt.Println("error")
	}
	rate, ok := rates[to].(float64)

	return Rate(rate)
}

func FetchRate(ctx context.Context, from, to string) (Rate, error) {

	resultChan := make(chan Rate)

	go func() {
		rate := request(from, to, ctx)
		resultChan <- rate
	}()

	select {
	case rate := <-resultChan:
		return rate, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}

}
