package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"server_kufatech/internal/handlers"
)

func SetupHealthRoutes(r chi.Router, healthHandler *handlers.HealthHandler, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/health", func(r chi.Router) {
		r.Get("/", healthHandler.HealthCheck)
	})
}
