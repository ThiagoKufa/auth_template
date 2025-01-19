package apperrors

import (
	"fmt"
)

type AppError struct {
	Message string
	Code    int
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) StatusCode() int {
	return e.Code
}

func NewValidationError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    400,
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    401,
	}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    403,
	}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    404,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    409,
	}
}

func NewInternalError(err error) *AppError {
	return &AppError{
		Message: "erro interno do servidor",
		Code:    500,
		Err:     err,
	}
}

func NewRateLimitError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    429,
	}
}
