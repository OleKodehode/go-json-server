package app

import (
	"net/http"

	"github.com/OleKodehode/go-json-server/internal/service"
)

func NewRouter(s *service.Service) http.Handler {
	mux := http.NewServeMux()
	h := NewHandler(s)
	
	// health check
	mux.HandleFunc("GET /health", HandleHealth)

	// Serve the same index.html file that the original used. No need for a handler
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	// GET collections or entries
	mux.HandleFunc("GET /{name}", h.GetAll)
	mux.HandleFunc("GET /{name}/{id}", h.GetByID)

	// Create new collections
	mux.HandleFunc("POST /{name}", h.Create)

	// Update entries
	mux.HandleFunc("PUT /{name}/{id}", h.Replace)
	mux.HandleFunc("PATCH /{name}/{id}", h.Update)

	// Delete entries
	mux.HandleFunc("DELETE /{name}/{id}", h.Delete)

	// Alternatively, wrap cors outside to omit OPTIONS requests logging
	return LoggingMiddleWare(CORSMiddleware(mux))

}

