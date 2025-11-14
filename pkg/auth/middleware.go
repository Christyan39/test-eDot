package auth

import (
	"net/http"

	"github.com/Christyan39/test-eDot/pkg/config"
	"github.com/labstack/echo/v4"
)

// JWTAuthMiddleware checks JWT, extracts user data, and saves to context
func JWTAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing Authorization header")
		}
		var tokenString string
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid Authorization header format")
		}

		user, err := ValidateToken(tokenString)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
		}

		// Save user info to context
		c.Set("user", user)
		return next(c)
	}
}

// ServiceAuthMiddleware validates service-to-service authentication using X-API-Key header
func ServiceAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Check for API key authentication header
		apiKey := c.Request().Header.Get("X-API-Key")
		expectedKey := config.GetEnv("API_KEY", "")

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
