package main

import (
	"fmt"
	Vu "gophercises/cyoa/views"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "up and running...")
	})
	mux.HandleFunc("/view", Vu.Views())
	http.ListenAndServe(":3000", mux)
	return
}
