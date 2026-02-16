package app

import (
	"log/slog"
	"net/http"
)

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer handlePanic(w,r)

		next.ServeHTTP(w,r)
	})
}

func handlePanic(w http.ResponseWriter, r *http.Request) {
				if err := recover(); err != nil {
				slog.Error("Panic recovered",
					"error", err,
					"path", r.URL.Path,
				)

				msg := map[string]string{"error": "Internal Server Error"}

				RespondJSON(w, http.StatusInternalServerError, msg)
			}
}