package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	apperrors "server_kufatech/internal/errors"
)

type AuthRateLimiter struct {
	attempts sync.Map
	window   time.Duration
	limit    int
}

type authAttempt struct {
	count     int
	startTime time.Time
}

func NewAuthRateLimiter(limit int, window time.Duration) *AuthRateLimiter {
	limiter := &AuthRateLimiter{
		window: window,
		limit:  limit,
	}
	go limiter.cleanup()
	return limiter
}

func (l *AuthRateLimiter) LimitAuthEndpoints(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Aplicar apenas em endpoints de autenticação
		if !strings.HasPrefix(r.URL.Path, "/auth/") {
			next.ServeHTTP(w, r)
			return
		}

		ip := getClientIP(r)
		now := time.Now()

		// Carregar ou criar janela de tempo
		value, _ := l.attempts.LoadOrStore(ip, &authAttempt{
			startTime: now,
		})
		attempt := value.(*authAttempt)

		// Resetar contador se passou a janela de tempo
		if now.Sub(attempt.startTime) > l.window {
			attempt.count = 0
			attempt.startTime = now
		}

		// Verificar limite
		if attempt.count >= l.limit {
			panic(apperrors.NewRateLimitError("muitas tentativas de autenticação"))
		}

		// Incrementar contador
		attempt.count++
		next.ServeHTTP(w, r)
	})
}

func (l *AuthRateLimiter) cleanup() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		l.attempts.Range(func(key, value interface{}) bool {
			if attempt, ok := value.(*authAttempt); ok {
				if now.Sub(attempt.startTime) > l.window {
					l.attempts.Delete(key)
				}
			}
			return true
		})
	}
}

func getClientIP(r *http.Request) string {
	// Tentar X-Real-IP
	ip := r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// Tentar X-Forwarded-For
	ip = r.Header.Get("X-Forwarded-For")
	if ip != "" {
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0])
	}

	// Usar RemoteAddr como fallback
	return strings.Split(r.RemoteAddr, ":")[0]
}
