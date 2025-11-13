package middleware

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

// getEnv gets environment variable with default fallback
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// ServiceAuthMiddleware validates service-to-service authentication using X-API-Key header
func ServiceAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check for API key authentication header
			apiKey := c.Request().Header.Get("X-API-Key")
			expectedKey := getEnv("API_KEY", "")

			if apiKey == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "API key required for internal endpoints",
				})
			}

			if apiKey != expectedKey {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "Invalid API key",
				})
			}

			return next(c)
		}
	}
}

// ServiceAuthMiddlewareWithKey validates service-to-service authentication with custom API key
func ServiceAuthMiddlewareWithKey(expectedKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check for API key authentication header
			apiKey := c.Request().Header.Get("X-API-Key")

			if apiKey == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "API key required for internal endpoints",
				})
			}

			if apiKey != expectedKey {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "Invalid API key",
				})
			}

			return next(c)
		}
	}
}
