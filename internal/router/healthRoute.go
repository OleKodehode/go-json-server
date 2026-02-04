package router

import (
	"encoding/json"
	"net/http"
)

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	responseItem := map[string]string{"status": "ok"}
	json.NewEncoder(w).Encode(responseItem)
}