package handlers

import (
	"net/http"
	"os"
)

// HealthHandler handles the reporting of the current health of the service.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		os.Exit(1)
	}
}
