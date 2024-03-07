package handlers

import "net/http"

// ReposHandler returns the latest public github repos.
func ReposHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}
}
