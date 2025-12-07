package middleware

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
)

// RequestLoggerMiddleware logs HTTP requests and responses
func RequestLoggerMiddleware(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			duration := time.Since(start)

			logger.InfoContext(
				c.Request().Context(),
				"http_request",
				slog.String("method", c.Request().Method),
				slog.String("path", c.Request().URL.Path),
				slog.Int("status", c.Response().Status),
				slog.String("remote_ip", c.RealIP()),
				slog.String("duration_ms", duration.String()),
				slog.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)),
			)

			return err
		}
	}
}
