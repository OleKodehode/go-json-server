package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/OleKodehode/go-json-server/internal/app"
)

func main() {
	// setup of logger using slog
	logger := slog.Default()

	port := os.Getenv("PORT")
	if port == "" { // dev env
		port = "8080"
	}

	host := os.Getenv("HOST")
	if host == "" { // dev env
		host = "localhost"
	}

	router := app.NewRouter()

	logger.Info("Server starting", "port", port)
	fmt.Printf("http://%s:%s/health", host, port) // convenience log
	if err := http.ListenAndServe(":"+port, router); err != nil {
		logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}