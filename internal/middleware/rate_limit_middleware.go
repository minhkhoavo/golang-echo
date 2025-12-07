package middleware

import (
	"log/slog"
	"golang-echo/pkg/response"
	"golang-echo/pkg/utils"

	"github.com/labstack/echo/v4"
)

var (
	rateLimiter utils.RateLimiter
	logger      *slog.Logger
)

// InitRateLimiter initializes the rate limiter
func InitRateLimiter(requestsPerMin int, loggerInstance *slog.Logger) {
	rateLimiter = utils.NewSlidingWindowLimiter(requestsPerMin, 60*1)
	logger = loggerInstance
}

// RateLimitMiddleware returns a rate limiting middleware
func RateLimitMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if rateLimiter == nil {
				return next(c)
			}

			key := c.RealIP()

			if !rateLimiter.AllowContext(c.Request().Context(), key) {
				if logger != nil {
					logger.WarnContext(c.Request().Context(), "rate_limit_exceeded",
						slog.String("ip", key),
					)
				}
				return response.TooManyRequests("RATE_LIMIT_EXCEEDED", "Too many requests. Please try again later.", nil)
			}

			return next(c)
		}
	}
}
