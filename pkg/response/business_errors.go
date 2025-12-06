package response

import (
	"net/http"
)

// Generic error helpers by HTTP status code
// This approach is more scalable and reduces boilerplate

// BadRequest returns a 400 Bad Request error
func BadRequest(key string, message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Key:     key,
		Message: message,
		Err:     err,
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
