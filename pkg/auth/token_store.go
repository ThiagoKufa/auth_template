package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const tokenBlacklistPrefix = "blacklist:"

type TokenStore struct {
	redis *redis.Client
}

func NewTokenStore(redis *redis.Client) *TokenStore {
	return &TokenStore{
		redis: redis,
	}
}

// Add adiciona um token à blacklist com um TTL específico
func (s *TokenStore) Add(ctx context.Context, token string, ttl time.Duration) error {
	key := s.formatKey(token)
	if err := s.redis.Set(ctx, key, true, ttl).Err(); err != nil {
		return fmt.Errorf("erro ao adicionar token à blacklist: %w", err)
	}
	return nil
}

// IsBlacklisted verifica se um token está na blacklist
func (s *TokenStore) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	key := s.formatKey(token)
	exists, err := s.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("erro ao verificar token na blacklist: %w", err)
	}
	return exists > 0, nil
}

// Remove remove um token da blacklist
func (s *TokenStore) Remove(ctx context.Context, token string) error {
	key := s.formatKey(token)
	if err := s.redis.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("erro ao remover token da blacklist: %w", err)
	}
	return nil
}

// formatKey formata a chave para o Redis
func (s *TokenStore) formatKey(token string) string {
	return tokenBlacklistPrefix + token
}
