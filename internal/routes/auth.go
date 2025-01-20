package routes

import (
	"time"

	"github.com/go-chi/chi/v5"

	"auth-template/internal/handlers"
	"auth-template/internal/middleware"
)

func SetupAuthRoutes(r chi.Router, authHandler *handlers.AuthHandler) {
	// Rate limiter específico para autenticação
	authLimiter := middleware.NewAuthRateLimiter(100, time.Hour) // 100 requisições por hora

	r.Route("/auth", func(r chi.Router) {
		// Aplicar rate limiting em todas as rotas de auth
		r.Use(authLimiter.LimitAuthEndpoints)

		r.Post("/login", authHandler.Login)
		r.Post("/register", authHandler.Register)
		r.Post("/refresh", authHandler.Refresh)
		r.Post("/logout", authHandler.Logout)

		// Rotas protegidas
		r.Group(func(r chi.Router) {
			r.Use(authHandler.AuthMiddleware)
			r.Get("/me", authHandler.Me)
		})
	})
}
