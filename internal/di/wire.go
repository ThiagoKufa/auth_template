//go:build wireinject
// +build wireinject

package di

import (
	"fmt"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"server_kufatech/internal/config"
	"server_kufatech/internal/handlers"
	"server_kufatech/internal/interfaces/repository"
	"server_kufatech/internal/interfaces/service"
	repo "server_kufatech/internal/repository"
	"server_kufatech/internal/services"
	"server_kufatech/pkg/auth"
	"server_kufatech/pkg/database"
	"server_kufatech/pkg/logger"
)

var containerSet = wire.NewSet(
	logger.NewLogger,
	database.NewDB,
	provideRedis,
	provideUserRepository,
	provideTokenManager,
	provideTokenBlacklist,
	provideAuthService,
	handlers.NewAuthHandler,
	handlers.NewHealthHandler,
	wire.Struct(new(Container), "*"),
)

func provideRedis(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})
}

func provideUserRepository(db *gorm.DB) repository.UserRepository {
	return repo.NewUserRepository(db)
}

func provideTokenManager(cfg *config.Config) *auth.TokenManager {
	return auth.NewTokenManager(
		cfg.Auth.AccessTokenSecret,
		cfg.Auth.RefreshTokenSecret,
		cfg.Auth.AccessTokenTTL,
		cfg.Auth.RefreshTokenTTL,
	)
}

func provideTokenBlacklist(redis *redis.Client) *services.TokenBlacklist {
	return services.NewTokenBlacklist(redis)
}

func provideAuthService(
	userRepo repository.UserRepository,
	tokenManager *auth.TokenManager,
	tokenBlacklist *services.TokenBlacklist,
	cfg *config.Config,
) service.AuthService {
	return services.NewAuthService(userRepo, tokenManager, tokenBlacklist, cfg)
}

// InitializeContainer inicializa o container de dependências
func InitializeContainer(cfg *config.Config) (*Container, error) {
	wire.Build(containerSet)
	return nil, nil
}
