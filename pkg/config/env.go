package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// LoadEnvFile loads environment variables from .env files
// It tries to load from the specified service config directory first,
// then falls back to root .env file
func LoadEnvFile(serviceName string) {
	// Try to load from configs/{serviceName}/.env first
	// Try to load from /app/{service}.env (for Docker volume mount)
	dockerEnvPath := "/app/" + serviceName + ".env"
	if err := godotenv.Load(dockerEnvPath); err != nil {
		// Try to load from configs/{serviceName}/.env (for local dev)
		log.Printf("Failed to load Docker env file: %s, trying local config...", dockerEnvPath)
		envPath := filepath.Join("configs", serviceName, ".env")
		if err := godotenv.Load(envPath); err != nil {
			// Try to load from root .env as fallback
			if err := godotenv.Load(); err != nil {
				log.Printf("Warning: No .env file found: %v", err)
			} else {
				log.Println("[STARTUP] Loaded environment from root .env")
			}
		} else {
			log.Printf("[STARTUP] Loaded environment from %s", envPath)
		}
	} else {
		log.Printf("[STARTUP] Loaded environment from %s (Docker volume)", dockerEnvPath)
	}
}

// GetEnv gets environment variable with default value
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
