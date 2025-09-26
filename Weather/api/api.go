package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func Fetch(c context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(c, http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(err)
		panic(err)

	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		panic(err)

	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		panic(err)

	}

	return body, nil
}
