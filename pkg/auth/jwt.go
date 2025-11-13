package auth

import (
	"fmt"
	"os"
	"time"

	userModels "github.com/Christyan39/test-eDot/internal/models/user"
	"github.com/golang-jwt/jwt"
)

// JWTClaims represents the JWT claims
type JWTClaims struct {
	User *userModels.AuthUser `json:"user"`
	jwt.StandardClaims
}

// GenerateToken generates a JWT token for the user
func GenerateToken(userModel *userModels.User) (string, error) {
	jwtSecret := getJWTSecret()

	authUser := &userModels.AuthUser{
		ID:    userModel.ID,
		Name:  userModel.Name,
		Email: userModel.Email,
		Phone: userModel.Phone,
	}

	claims := JWTClaims{
		User: authUser,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), // 24 hours
			IssuedAt:  time.Now().Unix(),
			Subject:   fmt.Sprintf("%d", userModel.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// ValidateToken validates a JWT token and returns the user claims
func ValidateToken(tokenString string) (*userModels.AuthUser, error) {
	jwtSecret := getJWTSecret()

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims.User, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// getJWTSecret returns the JWT secret from environment or default
func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-super-secret-jwt-key-change-in-production"
	}
	return secret
}
