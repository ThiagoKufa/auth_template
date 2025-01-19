package di

import (
	"server_kufatech/internal/config"
	"server_kufatech/internal/handlers"
	"server_kufatech/internal/interfaces/repository"
	"server_kufatech/internal/interfaces/service"
	"server_kufatech/internal/services"
	"server_kufatech/pkg/auth"
	"server_kufatech/pkg/logger"

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
