package app

import (
	"encoding/json"
	"fmt"
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

// validateBody checks whether a request has a valid JSON body
func validateBody (r *http.Request) (map[string]any, error) {
	item := map[string]any{}

	// Check to see if there is any body
	if r.Body == nil {
		return nil, fmt.Errorf("Request body is required")
	}

	// Try to decode - Return early if there is any error
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		return nil, fmt.Errorf("Invalid JSON format: %v", err)
	}

	// Check the decode result to make sure it's not empty
	if len(item) == 0 {
		return nil, fmt.Errorf("Request body can't be empty")
	}

	return item, nil
}