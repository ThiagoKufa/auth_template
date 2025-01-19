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

type TokenManager struct {
	accessSecret  string
	refreshSecret string
}

func NewTokenManager(accessSecret, refreshSecret string) *TokenManager {
	return &TokenManager{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
	}
}

// GenerateToken gera um token JWT com o tipo e duração especificados
func (m *TokenManager) GenerateToken(userID string, tokenType TokenType, duration time.Duration) (string, error) {
	claims := Claims{
		UserID: userID,
		Type:   tokenType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	// Escolher a chave secreta baseado no tipo do token
	secret := m.accessSecret
	if tokenType == TokenTypeRefresh {
		secret = m.refreshSecret
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateToken valida um token JWT e retorna suas claims
func (m *TokenManager) ValidateToken(tokenString string, expectedType TokenType) (*Claims, error) {
	// Escolher a chave secreta baseado no tipo esperado
	secret := m.accessSecret
	if expectedType == TokenTypeRefresh {
		secret = m.refreshSecret
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
		return nil, fmt.Errorf("claims inválidas")
	}

	if claims.Type != expectedType {
		return nil, fmt.Errorf("tipo de token inválido")
	}

	return claims, nil
}
