package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/OleKodehode/go-json-server/internal/app"
	"github.com/OleKodehode/go-json-server/internal/db"
	"github.com/OleKodehode/go-json-server/internal/model"
	"github.com/OleKodehode/go-json-server/internal/service"
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

	db, err := db.Load[model.Data]("db")
	if err != nil {
		logger.Error("Failure to load DB - ", "Database Error: ", err)
		os.Exit(1)
	}

	serviceLayer := service.New(db)

	router := app.NewRouter(serviceLayer)

	logger.Info("Server starting", "port", port)
	fmt.Printf("http://%s:%s/", host, port) // convenience log
	if err := http.ListenAndServe(":"+port, router); err != nil {
		logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}