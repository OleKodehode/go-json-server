package app

import (
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	RespondJSON(w, http.StatusOK, HealthResponse{Status: "ok"})
}