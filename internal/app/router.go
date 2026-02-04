package app

import (
	"encoding/json"
	"net/http"
)

// Helper function
func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
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
