package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	healthHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	}

	http.HandleFunc("/health", healthHandler)

	port := 8080
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		os.Exit(1)
	}
}
