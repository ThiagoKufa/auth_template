package auth

import (
	"sync"
	"time"
)

type BruteForceProtector struct {
	attempts     sync.Map
	maxAttempts  int
	blockTime    time.Duration
	cleanupTimer *time.Ticker
}

type attemptInfo struct {
	count     int
	lastTry   time.Time
	blockedAt time.Time
}

func NewBruteForceProtector(maxAttempts int, blockTime time.Duration) *BruteForceProtector {
	protector := &BruteForceProtector{
		maxAttempts:  maxAttempts,
		blockTime:    blockTime,
		cleanupTimer: time.NewTicker(time.Hour),
	}
	go protector.cleanup()
	return protector
}

func (p *BruteForceProtector) RecordAttempt(identifier string) bool {
	now := time.Now()

	// Carregar ou criar informações de tentativa
	value, _ := p.attempts.LoadOrStore(identifier, &attemptInfo{
		lastTry: now,
	})
	info := value.(*attemptInfo)

	// Verificar se está bloqueado
	if !info.blockedAt.IsZero() && now.Sub(info.blockedAt) < p.blockTime {
		return false
	}

	// Resetar contagem se passou muito tempo desde a última tentativa
	if now.Sub(info.lastTry) > time.Hour {
		info.count = 0
		info.blockedAt = time.Time{}
	}

	// Incrementar contagem
	info.count++
	info.lastTry = now

	// Bloquear se excedeu tentativas
	if info.count >= p.maxAttempts {
		info.blockedAt = now
		return false
	}

	return true
}

func (p *BruteForceProtector) Reset(identifier string) {
	p.attempts.Delete(identifier)
}

func (p *BruteForceProtector) cleanup() {
	for range p.cleanupTimer.C {
		now := time.Now()
		p.attempts.Range(func(key, value interface{}) bool {
			if info, ok := value.(*attemptInfo); ok {
				// Limpar tentativas antigas
				if now.Sub(info.lastTry) > time.Hour*24 {
					p.attempts.Delete(key)
				}
			}
			return true
		})
	}
}

func (p *BruteForceProtector) Close() {
	p.cleanupTimer.Stop()
}
