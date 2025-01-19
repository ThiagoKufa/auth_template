package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

type Claims struct {
	UserID string    `json:"user_id"`
	Type   TokenType `json:"type"`
	jwt.StandardClaims
}

func (c *Claims) Valid() error {
	return c.StandardClaims.Valid()
}

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
