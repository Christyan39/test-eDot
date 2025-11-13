package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "github.com/Christyan39/test-eDot/docs"
	handlers "github.com/Christyan39/test-eDot/internal/handlers/user"
	repositories "github.com/Christyan39/test-eDot/internal/repositories/user"
	usecases "github.com/Christyan39/test-eDot/internal/usecases/user"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Test-eDot API
// @version 1.0
// @description API documentation for Test-eDot user management system
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load environment variables from .env file
	loadEnvFile()

	// Connect to MySQL database
	db, err := connectToMySQL()
	if err != nil {
		log.Fatalf("Warning: Failed to connect to MySQL: %v", err)
	}

	fmt.Println("Starting server with database configuration...")

	// Initialize layers
	userRepo := repositories.NewUserRepository(db)
	userUsecase := usecases.NewUserUsecase(userRepo)
	userHandler := handlers.NewUserHandler(userUsecase)

	authUsecase := usecases.NewAuthUsecase(userRepo)
	authHandler := handlers.NewAuthHandler(authUsecase)

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Swagger documentation route
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Health check route
	// @Summary Health check
	// @Description Check if the server is running
	// @Tags system
	// @Accept json
	// @Produce json
	// @Success 200 {object} map[string]string "Server is running"
	// @Router /health [get]
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// API routes
	api := e.Group("/api/v1")

	// Auth routes
	auth := api.Group("/auth")
	auth.POST("/login", authHandler.HandleDirectLogin)
	auth.POST("/secure-login", authHandler.HandleEnvelopeLogin)
	auth.POST("/create-envelope", authHandler.CreateEnvelope)

	// User routes
	api.POST("/users", userHandler.CreateUser)

	// Start server
	port := getEnv("PORT", "8080")
	fmt.Printf("Server starting on port %s\n", port)
	fmt.Printf("Health check: http://localhost:%s/health\n", port)
	fmt.Printf("API Base URL: http://localhost:%s/api/v1\n", port)

	if err := e.Start(":" + port); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

// loadEnvFile loads environment variables from .env files
func loadEnvFile() {
	// Try to load from configs/user/.env first
	envPath := filepath.Join("configs", "user", ".env")
	if err := godotenv.Load(envPath); err != nil {
		// Try to load from root .env as fallback
		if err := godotenv.Load(".env"); err != nil {
			log.Printf("No .env file found in %s or root directory, using system environment variables", envPath)
		} else {
			log.Printf("Loaded environment variables from root .env file")
		}
	} else {
		log.Printf("Loaded environment variables from %s", envPath)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// connectToMySQL establishes connection to MySQL database
func connectToMySQL() (*sql.DB, error) {
	// MySQL connection parameters from environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "edot_user")

	// MySQL DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Open database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Printf("Successfully connected to MySQL database: %s@%s:%s/%s\n",
		dbUser, dbHost, dbPort, dbName)

	return db, nil
}
