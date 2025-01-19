package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"server_kufatech/internal/handlers"
	"server_kufatech/internal/middleware"
	"server_kufatech/pkg/logger"
)

func SetupRoutes(
	r chi.Router,
	log *logger.Logger,
	authHandler *handlers.AuthHandler,
	healthHandler *handlers.HealthHandler,
) {
	// Middleware básicos
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.NewErrorHandler(log).Handle)

	// Middleware de log personalizado
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info("Request: %s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})

	// Setup das rotas
	SetupAuthRoutes(r, authHandler)
	SetupHealthRoutes(r, healthHandler, authHandler.AuthMiddleware)
}
