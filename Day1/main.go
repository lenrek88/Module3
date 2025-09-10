package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Домашняя страница!"))

	})

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Привет, мир!"))
	})

	http.HandleFunc("/greet", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		w.Write([]byte("Привет, " + name + "!"))
	})

	e := http.ListenAndServe(":8080", nil)
	if e != nil {
		panic(e)
	}
}
