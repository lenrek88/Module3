package main

import (
	"context"
	"encoding/json"
	"fmt"
	"lenrek88/exchanger"
	"net/http"
)

func rateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	rate, err := exchanger.FetchRate(ctx, from, to)
	if err != nil {
		context.WithCancel(ctx)
		fmt.Println(err)
		return
	}

	data, _ := json.Marshal(rate)
	w.Write(data)

}

func main() {

	http.HandleFunc("/rate", rateHandler)
	http.ListenAndServe(":8080", nil)
}
