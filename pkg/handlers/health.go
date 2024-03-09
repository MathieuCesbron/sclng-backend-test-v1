package handlers

import (
	"fmt"
	"net/http"
)

// HealthHandler handles the reporting of the current health of the service.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	_, err := w.Write([]byte("OK"))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed writing to the conection: %s", err.Error()), http.StatusInternalServerError)
	}
}
