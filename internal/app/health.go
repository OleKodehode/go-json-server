package app

import (
	"net/http"
)

// simple endpoint to check if the server is running - No need to add service to this
type HealthResponse struct {
	Status string `json:"status"`
}

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	RespondJSON(w, http.StatusOK, HealthResponse{Status: "ok"})
}