package response

import (
	"net/http"
)

type AppError struct {
	Code     int
	Key      string
	Message  string
	Err      error
	FieldErr map[string]string // ← Changed to map format (field name → error message)
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(code int, key string, msg string, err error) *AppError {
	return &AppError{
		Code:    code,
		Key:     key,
		Message: msg,
		Err:     err,
	}
}

// NewAppErrorWithFieldErrors creates an AppError with field-level validation errors (map format)
func NewAppErrorWithFieldErrors(code int, key string, msg string, fieldErrs map[string]string) *AppError {
	return &AppError{
		Code:     code,
		Key:      key,
		Message:  msg,
		FieldErr: fieldErrs,
	}
}

// BadRequest returns a 400 Bad Request error
func BadRequest(key string, message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Key:     key,
		Message: message,
		Err:     err,
	}
}

// BadRequestWithFields returns a 400 Bad Request error with field errors (map format)
func BadRequestWithFields(key string, message string, fieldErrs map[string]string) *AppError {
	return &AppError{
		Code:     http.StatusBadRequest,
		Key:      key,
		Message:  message,
		FieldErr: fieldErrs,
	}
}

// Unauthorized returns a 401 Unauthorized error
func Unauthorized(key string, message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Key:     key,
		Message: message,
		Err:     err,
	}
}

// Forbidden returns a 403 Forbidden error
func Forbidden(key string, message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusForbidden,
		Key:     key,
		Message: message,
		Err:     err,
	}
}

// NotFound returns a 404 Not Found error
func NotFound(key string, message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Key:     key,
		Message: message,
		Err:     err,
	}
}

// Conflict returns a 409 Conflict error
func Conflict(key string, message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusConflict,
		Key:     key,
		Message: message,
		Err:     err,
	}
}

// Internal returns a 500 Internal Server Error
func Internal(err error) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Key:     "INTERNAL_SERVER_ERROR",
		Message: "Internal server error",
		Err:     err,
	}
}

// TooManyRequests returns a 429 Too Many Requests error
func TooManyRequests(key string, message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusTooManyRequests,
		Key:     key,
		Message: message,
		Err:     err,
	}
}
