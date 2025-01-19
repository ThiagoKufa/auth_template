package services

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"server_kufatech/internal/config"
	"server_kufatech/internal/entity"
	apperrors "server_kufatech/internal/errors"
	"server_kufatech/internal/interfaces/repository"
	"server_kufatech/internal/interfaces/service"
	"server_kufatech/internal/validation"
	"server_kufatech/pkg/auth"
)

type AuthService struct {
	userRepo       repository.UserRepository
	tokenManager   *auth.TokenManager
	tokenBlacklist *TokenBlacklist
	config         *config.Config
}

func NewAuthService(userRepo repository.UserRepository, tokenManager *auth.TokenManager, tokenBlacklist *TokenBlacklist, config *config.Config) service.AuthService {
	return &AuthService{
		userRepo:       userRepo,
		tokenManager:   tokenManager,
		tokenBlacklist: tokenBlacklist,
		config:         config,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password string) error {
	// Validar email
	sanitizedEmail, err := validation.ValidateEmail(email)
	if err != nil {
		return apperrors.NewValidationError("email inválido")
	}

	// Validar senha
	if err := validation.ValidatePassword(password, validation.DefaultPasswordPolicy); err != nil {
		return apperrors.NewValidationError(err.Error())
	}

	// Verificar se email já existe
	exists, err := s.userRepo.ExistsByEmail(ctx, sanitizedEmail)
	if err != nil {
		return fmt.Errorf("erro ao verificar email: %w", err)
	}
	if exists {
		return apperrors.NewConflictError("email já cadastrado")
	}

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("erro ao gerar hash da senha: %w", err)
	}

	// Criar usuário
	user := &entity.User{
		Email:    sanitizedEmail,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return fmt.Errorf("erro ao criar usuário: %w", err)
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*service.TokenPair, error) {
	// Buscar usuário pelo email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, apperrors.NewUnauthorizedError("credenciais inválidas")
	}

	// Verificar senha
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, apperrors.NewUnauthorizedError("credenciais inválidas")
	}

	// Gerar tokens
	userID := fmt.Sprintf("%d", user.ID)
	accessToken, err := s.tokenManager.GenerateToken(userID, auth.TokenTypeAccess)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar access token: %w", err)
	}

	refreshToken, err := s.tokenManager.GenerateToken(userID, auth.TokenTypeRefresh)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar refresh token: %w", err)
	}

	return &service.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (*service.TokenPair, error) {
	// Validar refresh token
	claims, err := s.tokenManager.ValidateToken(refreshToken, auth.TokenTypeRefresh)
	if err != nil {
		return nil, apperrors.NewUnauthorizedError("refresh token inválido")
	}

	// Verificar se o token está na blacklist
	blacklisted, err := s.tokenBlacklist.IsBlacklisted(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar token: %w", err)
	}
	if blacklisted {
		return nil, apperrors.NewUnauthorizedError("refresh token inválido")
	}

	// Adicionar o token atual à blacklist
	if err := s.tokenBlacklist.Add(ctx, refreshToken, s.config.Auth.RefreshTokenTTL); err != nil {
		return nil, fmt.Errorf("erro ao invalidar token: %w", err)
	}

	// Gerar novos tokens
	accessToken, err := s.tokenManager.GenerateToken(claims.UserID, auth.TokenTypeAccess)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar access token: %w", err)
	}

	newRefreshToken, err := s.tokenManager.GenerateToken(claims.UserID, auth.TokenTypeRefresh)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar refresh token: %w", err)
	}

	return &service.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) ValidateAccessToken(ctx context.Context, token string) error {
	_, err := s.tokenManager.ValidateToken(token, auth.TokenTypeAccess)
	if err != nil {
		return apperrors.NewUnauthorizedError("token inválido")
	}
	return nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	// Validar refresh token
	_, err := s.tokenManager.ValidateToken(refreshToken, auth.TokenTypeRefresh)
	if err != nil {
		return apperrors.NewUnauthorizedError("refresh token inválido")
	}

	// Adicionar o token à blacklist
	if err := s.tokenBlacklist.Add(ctx, refreshToken, s.config.Auth.RefreshTokenTTL); err != nil {
		return fmt.Errorf("erro ao invalidar token: %w", err)
	}

	return nil
}
