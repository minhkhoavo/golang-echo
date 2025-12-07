package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Response represents a standardized response wrapper
// All responses should follow this structure
type Response[T any] struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Data      T      `json:"data"`
	RequestID string `json:"request_id,omitempty"`
}

// ListResponse represents a list response with pagination metadata
type ListResponse[T any] struct {
	Code       string          `json:"code"`
	Message    string          `json:"message"`
	Data       T               `json:"data"`
	Pagination *PaginationMeta `json:"pagination,omitempty"`
	RequestID  string          `json:"request_id,omitempty"`
}

// PaginationMeta contains pagination metadata
type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// ErrorResponse represents a standardized error response
// Errors field is now a dict (map) for better client-side consumption
// Example: { "email": "Email is required", "name": "Name is required" }
type ErrorResponse struct {
	Code      string            `json:"code"`
	Message   string            `json:"message"`
	Errors    map[string]string `json:"errors,omitempty"` // ← Changed from []FieldError to map
	RequestID string            `json:"request_id,omitempty"`
}

// ToErrorResponse converts AppError to ErrorResponse with request context
func (e *AppError) ToErrorResponse(c echo.Context) ErrorResponse {
	return ErrorResponse{
		Code:      e.Key,
		Message:   e.Message,
		RequestID: c.Response().Header().Get(echo.HeaderXRequestID),
	}
}

// Success returns a 200 OK success response with data
func Success[T any](c echo.Context, code string, message string, data T) error {
	return c.JSON(http.StatusOK, Response[T]{
		Code:      code,
		Message:   message,
		Data:      data,
		RequestID: c.Response().Header().Get(echo.HeaderXRequestID),
	})
}

// Created returns a 201 Created success response with data
func Created[T any](c echo.Context, code string, message string, data T) error {
	return c.JSON(http.StatusCreated, Response[T]{
		Code:      code,
		Message:   message,
		Data:      data,
		RequestID: c.Response().Header().Get(echo.HeaderXRequestID),
	})
}

// List returns a 200 OK list response without pagination
func List[T any](c echo.Context, code string, message string, data T) error {
	return c.JSON(http.StatusOK, ListResponse[T]{
		Code:      code,
		Message:   message,
		Data:      data,
		RequestID: c.Response().Header().Get(echo.HeaderXRequestID),
	})
}

// ListWithPagination returns a 200 OK list response with pagination
func ListWithPagination[T any](c echo.Context, code string, message string, data T, pagination *PaginationMeta) error {
	return c.JSON(http.StatusOK, ListResponse[T]{
		Code:       code,
		Message:    message,
		Data:       data,
		Pagination: pagination,
		RequestID:  c.Response().Header().Get(echo.HeaderXRequestID),
	})
}

// NoContent returns a 204 No Content response
func NoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
