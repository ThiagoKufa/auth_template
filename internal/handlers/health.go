package handlers

import (
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Verifica se há um parâmetro de delay para teste de timeout
	if delayStr := r.URL.Query().Get("delay"); delayStr != "" {
		if delay, err := strconv.Atoi(delayStr); err == nil {
			time.Sleep(time.Duration(delay) * time.Second)
		}
	}

	// Verifica a conexão com o banco de dados
	sqlDB, err := h.db.DB()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status": "error", "message": "database connection error"}`))
		return
	}

	err = sqlDB.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status": "error", "message": "database ping failed"}`))
		return
	}

	// Gera uma resposta grande para testar a compressão
	response := `{
		"status": "healthy",
		"database": "connected",
		"timestamp": "` + time.Now().Format(time.RFC3339) + `",
		"details": {
			"database_stats": {
				"open_connections": 0,
				"in_use": 0,
				"idle": 0,
				"wait_count": 0,
				"wait_duration": 0,
				"max_idle_closed": 0,
				"max_lifetime_closed": 0
			},
			"system_info": {
				"go_version": "1.22",
				"os": "linux",
				"arch": "amd64",
				"cpu_count": 8,
				"memory_stats": {
					"alloc": 0,
					"total_alloc": 0,
					"sys": 0,
					"num_gc": 0
				}
			}
		}
	}`

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(response))
}
