package main

import (
	"net/http"
	"os"

	"github.com/Scalingo/sclng-backend-test-v1/internal/handlers"
)

func main() {
	http.HandleFunc("/health", handlers.HealthHandler)
	http.HandleFunc("/repos", handlers.NewReposHandler())

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		os.Exit(1)
	}
}
