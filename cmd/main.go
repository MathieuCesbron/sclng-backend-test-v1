package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Scalingo/sclng-backend-test-v1/internal/handlers"
)

func main() {
	http.HandleFunc("/health", handlers.HealthHandler)

	port := 8080
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		os.Exit(1)
	}
}
