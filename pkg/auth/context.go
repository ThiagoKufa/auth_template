package auth

import (
	"context"
)

type contextKey string

const userEmailKey contextKey = "userEmail"

// WithUserEmail adiciona o email do usuário ao contexto
func WithUserEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, userEmailKey, email)
}

// GetUserEmail obtém o email do usuário do contexto
func GetUserEmail(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(userEmailKey).(string)
	return email, ok
}
