package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	apperrors "server_kufatech/internal/errors"
	"server_kufatech/internal/interfaces/service"
	"server_kufatech/pkg/logger"
)

type AuthHandler struct {
	authService service.AuthService
	log         *logger.Logger
}

func NewAuthHandler(authService service.AuthService, log *logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		log:         log,
	}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Request: %s %s", r.Method, r.URL.Path)

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Erro ao decodificar requisição: %v", err)
		h.writeError(w, apperrors.NewValidationError("requisição inválida"))
		return
	}

	err := h.authService.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		h.log.Error("Erro no registro: %v", err)
		h.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Request: %s %s", r.Method, r.URL.Path)

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Erro ao decodificar requisição: %v", err)
		h.writeError(w, apperrors.NewValidationError("requisição inválida"))
		return
	}

	tokens, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.log.Error("Erro no login: %v", err)
		h.writeError(w, err)
		return
	}

	resp := tokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	h.writeJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Request: %s %s", r.Method, r.URL.Path)

	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Erro ao decodificar requisição: %v", err)
		h.writeError(w, apperrors.NewValidationError("requisição inválida"))
		return
	}

	tokens, err := h.authService.RefreshTokens(r.Context(), req.RefreshToken)
	if err != nil {
		h.log.Error("Erro no refresh: %v", err)
		h.writeError(w, err)
		return
	}

	resp := tokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	h.writeJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Request: %s %s", r.Method, r.URL.Path)

	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Erro ao decodificar requisição: %v", err)
		h.writeError(w, apperrors.NewValidationError("requisição inválida"))
		return
	}

	if err := h.authService.Logout(r.Context(), req.RefreshToken); err != nil {
		h.log.Error("Erro no logout: %v", err)
		h.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.log.Error("Erro ao codificar resposta: %v", err)
	}
}

func (h *AuthHandler) writeError(w http.ResponseWriter, err error) {
	var status int
	var message string

	switch e := err.(type) {
	case *apperrors.AppError:
		status = e.StatusCode()
		message = e.Error()
	default:
		status = http.StatusInternalServerError
		message = "erro interno do servidor"
	}

	h.writeJSON(w, status, map[string]interface{}{
		"error": message,
		"code":  status,
	})
}

func (h *AuthHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			http.Error(w, "token não fornecido", http.StatusUnauthorized)
			return
		}

		if err := h.authService.ValidateAccessToken(r.Context(), token); err != nil {
			http.Error(w, "token inválido", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
