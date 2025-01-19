package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

type visitor struct {
	lastSeen time.Time
	count    int
}

type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     int
	per      time.Duration
}

func NewRateLimiter(rate int, per time.Duration) *RateLimiter {
	limiter := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		per:      per,
	}

	// Inicia limpeza em background
	go limiter.cleanup()
	return limiter
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.per {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) getIP(r *http.Request) string {
	// Tenta X-Real-IP primeiro
	ip := r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// Tenta X-Forwarded-For
	forwardedFor := r.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		ips := strings.Split(forwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}

	// Usa RemoteAddr como fallback
	return strings.Split(r.RemoteAddr, ":")[0]
}

func (rl *RateLimiter) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := rl.getIP(r)

		rl.mu.Lock()
		v, exists := rl.visitors[ip]
		if !exists {
			rl.visitors[ip] = &visitor{
				lastSeen: time.Now(),
				count:    1,
			}
			rl.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		// Reset contador se passou o período
		if time.Since(v.lastSeen) > rl.per {
			v.count = 1
			v.lastSeen = time.Now()
			rl.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		// Incrementa contador
		v.count++
		v.lastSeen = time.Now()

		// Verifica limite
		if v.count > rl.rate {
			rl.mu.Unlock()
			w.Header().Set("Retry-After", rl.per.String())
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		rl.mu.Unlock()
		next.ServeHTTP(w, r)
	})
}
