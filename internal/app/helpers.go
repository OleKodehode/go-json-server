package app

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Helper function
func RespondJSON(w http.ResponseWriter, status int, data any) {
	// If there is no body or no content
	// Allows me to reuse this for responses that might not send any data.
	if status == http.StatusNoContent || data == nil {
		w.WriteHeader(status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

type ErrorMessage struct {
	Error string `json:"error"`
}

func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, ErrorMessage{Error:message})
}

func totalHeader(w http.ResponseWriter, total int) {
	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
}