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

func request(urlConf string, from, to string, ctx context.Context) Rate {

	url := urlConf + from + ".json"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil) // запрос с контекстом

	if err != nil {
		fmt.Println("ERROR TO REQUEST ERR:", err)
	}
	delay := time.Second * 4
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

func FetchRate(urlConf string, ctx context.Context, from, to string) (Rate, error) {

	resultChan := make(chan Rate)

	go func() {
		rate := request(urlConf, from, to, ctx)
		resultChan <- rate
	}()

	select {
	case rate := <-resultChan:
		return rate, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}

}
