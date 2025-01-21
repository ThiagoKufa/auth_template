package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/redis/go-redis/v9"
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

// Claims representa os dados armazenados no token JWT
type Claims struct {
	UserID string    `json:"user_id"`
	Type   TokenType `json:"type"`
	jwt.StandardClaims
}

func (c *Claims) Valid() error {
	return c.StandardClaims.Valid()
}

// TokenManager gerencia a geração e validação de tokens JWT
type TokenManager struct {
	accessSecret    string
	refreshSecret   string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewTokenManager(accessSecret, refreshSecret string, accessTokenTTL, refreshTokenTTL time.Duration) *TokenManager {
	return &TokenManager{
		accessSecret:    accessSecret,
		refreshSecret:   refreshSecret,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

// GenerateToken gera um novo token JWT do tipo especificado
func (m *TokenManager) GenerateToken(userID string, tokenType TokenType) (string, error) {
	var duration time.Duration
	var secret string

	switch tokenType {
	case TokenTypeAccess:
		duration = m.accessTokenTTL
		secret = m.accessSecret
	case TokenTypeRefresh:
		duration = m.refreshTokenTTL
		secret = m.refreshSecret
	default:
		return "", fmt.Errorf("tipo de token inválido: %s", tokenType)
	}

	claims := &Claims{
		UserID: userID,
		Type:   tokenType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateToken valida um token JWT e retorna suas claims
func (m *TokenManager) ValidateToken(tokenString string, expectedType TokenType) (*Claims, error) {
	var secret string
	switch expectedType {
	case TokenTypeAccess:
		secret = m.accessSecret
	case TokenTypeRefresh:
		secret = m.refreshSecret
	default:
		return nil, fmt.Errorf("tipo de token inválido: %s", expectedType)
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("erro ao validar token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("token inválido")
	}

	if claims.Type != expectedType {
		return nil, fmt.Errorf("tipo de token inválido: esperado %s, recebido %s", expectedType, claims.Type)
	}

	return claims, nil
}

// TokenStore gerencia o armazenamento de tokens na blacklist
type TokenStore struct {
	redis *redis.Client
}

func NewTokenStore(redis *redis.Client) *TokenStore {
	return &TokenStore{
		redis: redis,
	}
}

// Add adiciona um token à blacklist
func (s *TokenStore) Add(ctx context.Context, token string, ttl time.Duration) error {
	return s.redis.Set(ctx, fmt.Sprintf("blacklist:%s", token), true, ttl).Err()
}

// IsBlacklisted verifica se um token está na blacklist
func (s *TokenStore) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	exists, err := s.redis.Exists(ctx, fmt.Sprintf("blacklist:%s", token)).Result()
	if err != nil {
		return false, fmt.Errorf("erro ao verificar blacklist: %w", err)
	}
	return exists > 0, nil
}
