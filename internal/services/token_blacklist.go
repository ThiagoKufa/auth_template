package services

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	blacklistKeyPrefix = "blacklist:token:"
)

type TokenBlacklist struct {
	redis *redis.Client
}

func NewTokenBlacklist(redis *redis.Client) *TokenBlacklist {
	return &TokenBlacklist{
		redis: redis,
	}
}

// Add adiciona um token à blacklist com um TTL específico
func (b *TokenBlacklist) Add(ctx context.Context, token string, ttl time.Duration) error {
	key := fmt.Sprintf("%s%s", blacklistKeyPrefix, token)
	return b.redis.Set(ctx, key, "revoked", ttl).Err()
}

// IsBlacklisted verifica se um token está na blacklist
func (b *TokenBlacklist) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("%s%s", blacklistKeyPrefix, token)
	exists, err := b.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("erro ao verificar token na blacklist: %w", err)
	}
	return exists > 0, nil
}

// Remove remove um token da blacklist (útil para testes)
func (b *TokenBlacklist) Remove(ctx context.Context, token string) error {
	key := fmt.Sprintf("%s%s", blacklistKeyPrefix, token)
	return b.redis.Del(ctx, key).Err()
}
