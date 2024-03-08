package main

import (
	"log"
	"net/http"

	"github.com/Scalingo/sclng-backend-test-v1/internal/handlers"
)

func main() {
	l := log.Default()

	http.HandleFunc("/health", handlers.HealthHandler)
	http.HandleFunc("/repos", handlers.NewReposHandler(l))

	l.Println("Starting server")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
