package service

import (
	"context"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthService interface {
	Register(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) (*TokenPair, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*TokenPair, error)
	ValidateAccessToken(ctx context.Context, token string) error
	Logout(ctx context.Context, refreshToken string) error
}
