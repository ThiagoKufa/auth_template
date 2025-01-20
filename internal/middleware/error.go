package middleware

import (
	"encoding/json"
	"net/http"

	apperrors "auth-template/internal/errors"
	"auth-template/pkg/logger"
)

type ErrorHandler struct {
	log *logger.Logger
}

func NewErrorHandler(log *logger.Logger) *ErrorHandler {
	return &ErrorHandler{
		log: log,
	}
}

type errorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func (h *ErrorHandler) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				var appErr *apperrors.AppError
				if e, ok := err.(*apperrors.AppError); ok {
					appErr = e
				} else if e, ok := err.(error); ok {
					appErr = apperrors.NewInternalError(e)
				} else {
					appErr = apperrors.NewInternalError(nil)
				}

				h.log.Error("Erro na requisição: %v", appErr)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(appErr.StatusCode())
				json.NewEncoder(w).Encode(errorResponse{
					Error: appErr.Message,
					Code:  appErr.Code,
				})
			}
		}()

		next.ServeHTTP(w, r)
	})
}
