package apperr

import (
	"net/http"
	"warehouse/pkg/utils"
)

type AppError struct {
	StatusCode int    `json:"-" swaggerignore:"true"`
	Code       string `json:"code" example:"NOT_FOUND"`
	Message    string `json:"message" example:"Resource not found"`
}

func (e *AppError) Error() string {
	return e.Message
}

// Predefined common errors
var (
	ErrNotFound     = &AppError{StatusCode: http.StatusNotFound, Code: "NOT_FOUND", Message: "Resource not found"}
	ErrBadRequest   = &AppError{StatusCode: http.StatusBadRequest, Code: "BAD_REQUEST", Message: "Bad request"}
	ErrConflict     = &AppError{StatusCode: http.StatusConflict, Code: "CONFLICT", Message: "Resource already exists"}
	ErrInternal     = &AppError{StatusCode: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Internal server error"}
	ErrUnauthorized = &AppError{StatusCode: http.StatusUnauthorized, Code: "UNAUTHORIZED", Message: "Unauthorized access"}
	ErrForbidden    = &AppError{StatusCode: http.StatusForbidden, Code: "FORBIDDEN", Message: "Access forbidden"}
)

// Helper functions
func New(statusCode int, code, message string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
	}
}

func NotFoundError(resource string) *AppError {
	return &AppError{
		StatusCode: http.StatusNotFound,
		Code:       "NOT_FOUND",
		Message:    resource + " not found",
	}
}

func BadRequestError(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		Code:       "BAD_REQUEST",
		Message:    utils.DefaultIfEmpty(message, "Bad request"),
	}
}
func ConflictError(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusConflict,
		Code:       "CONFLICT",
		Message:    utils.DefaultIfEmpty(message, "Resource already exists"),
	}
}

func InternalServerError(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusInternalServerError,
		Code:       "INTERNAL_ERROR",
		Message:    utils.DefaultIfEmpty(message, "Internal server error"),
	}
}
func UnauthorizedError(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusUnauthorized,
		Code:       "UNAUTHORIZED",
		Message:    utils.DefaultIfEmpty(message, "Unauthorized access"),
	}
}
func ForbiddenError(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusForbidden,
		Code:       "FORBIDDEN",
		Message:    utils.DefaultIfEmpty(message, "Access forbidden"),
	}
}
