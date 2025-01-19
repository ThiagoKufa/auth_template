package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"server_kufatech/internal/config"
	"server_kufatech/internal/di"
	"server_kufatech/internal/middleware"
	"server_kufatech/internal/routes"
)

func main() {
	// Carregar configuração
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// Inicializar container de dependências
	container, err := di.InitializeContainer(cfg)
	if err != nil {
		panic(err)
	}

	// Criar o router Chi
	r := chi.NewRouter()

	// Configurar middlewares globais
	rateLimiter := middleware.NewRateLimiter(10, time.Minute) // 10 requisições por minuto para teste
	r.Use(
		middleware.SecurityHeaders,                       // headers de segurança primeiro
		middleware.CORS(&container.Config.Security.CORS), // depois CORS
		middleware.Compress,                              // depois compressão
		middleware.Timeout(30*time.Second),               // depois timeout
		rateLimiter.RateLimit,                            // rate limiting por último
	)

	// Setup das rotas
	routes.SetupRoutes(r, container.Logger, container.AuthHandler, container.HealthHandler)

	// Iniciar o servidor
	container.Logger.Info("Servidor iniciado na porta %s", container.Config.Server.Port)
	if err := http.ListenAndServe(container.Config.Server.Port, r); err != nil {
		container.Logger.Error("Erro ao iniciar o servidor: %v", err)
	}
}
