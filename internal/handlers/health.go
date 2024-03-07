package handlers

import (
	"net/http"
)

// HealthHandler handles the reporting of the current health of the service.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	_, ok := w.Write([]byte("OK"))
	if ok != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
