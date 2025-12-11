package middleware

import (
	"strings"

	"golang-echo/pkg/constants"
	"golang-echo/pkg/response"
	"golang-echo/pkg/utils"

	"github.com/labstack/echo/v4"
)

func JWTMiddleware(jwtManager *utils.JWTManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return response.Unauthorized("MISSING_TOKEN", "Authorization header is missing", nil)
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return response.Unauthorized("INVALID_TOKEN_FORMAT", "Invalid token format. Use 'Bearer <token>'", nil)
			}
			tokenString := parts[1]

			claims, err := jwtManager.VerifyToken(tokenString)
			if err != nil {
				return response.Unauthorized("INVALID_TOKEN", "Token is invalid or expired", err)
			}

			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)
			c.Set("name", claims.Name)
			c.Set("role", claims.Role)
			return next(c)
		}
	}
}

func AdminMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, ok := c.Get("role").(string)
			if !ok || !constants.IsAdmin(role) {
				return response.Forbidden("FORBIDDEN", "You do not have permission to access this resource", nil)
			}
			return next(c)
		}
	}
}

func GetUserIDFromContext(c echo.Context) int {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return 0
	}
	return userID
}

func GetEmailFromContext(c echo.Context) string {
	email, ok := c.Get("email").(string)
	if !ok {
		return ""
	}
	return email
}
