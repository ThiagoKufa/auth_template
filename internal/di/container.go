package di

import (
	"auth-template/internal/config"
	"auth-template/internal/handlers"
	"auth-template/internal/interfaces/repository"
	"auth-template/internal/interfaces/service"
	"auth-template/internal/services"
	"auth-template/pkg/auth"
	"auth-template/pkg/logger"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Container struct {
	Config         *config.Config
	Logger         *logger.Logger
	DB             *gorm.DB
	Redis          *redis.Client
	UserRepo       repository.UserRepository
	TokenManager   *auth.TokenManager
	TokenBlacklist *services.TokenBlacklist
	AuthService    service.AuthService
	AuthHandler    *handlers.AuthHandler
	HealthHandler  *handlers.HealthHandler
}
