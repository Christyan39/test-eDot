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
	handlers "github.com/Christyan39/test-eDot/internal/handlers/product"
	repositories "github.com/Christyan39/test-eDot/internal/repositories/product"
	usecases "github.com/Christyan39/test-eDot/internal/usecases/product"
	pkgMiddleware "github.com/Christyan39/test-eDot/pkg/middleware"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Product Service API
// @version 1.0
// @description API documentation for Product Service - manages products and categories
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Configure logging to output to stdout with timestamps
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Load environment variables from .env file
	loadEnvFile()

	// Connect to MySQL database
	db, err := connectToMySQL()
	if err != nil {
		log.Fatalf("Warning: Failed to connect to MySQL: %v", err)
	}

	fmt.Println("Starting Product Service...")
	log.Println("=== PRODUCT SERVICE STARTING ===")
	log.Println("[STARTUP] Product Service initialization started")
	log.Printf("[STARTUP] Process ID: %d", os.Getpid())
	if wd, err := os.Getwd(); err == nil {
		log.Printf("[STARTUP] Working directory: %s", wd)
	}

	// Initialize layers
	log.Println("[STARTUP] Initializing repository layer...")
	productRepo := repositories.NewProductRepository(db)
	log.Println("[STARTUP] Initializing usecase layer...")
	productUsecase := usecases.NewProductUsecase(productRepo)
	log.Println("[STARTUP] Initializing handler layer...")
	productHandler := handlers.NewProductHandler(productUsecase)

	// Initialize Echo
	log.Println("[STARTUP] Initializing Echo web framework...")
	e := echo.New()
	e.Debug = true

	// Configure Echo logger
	e.Logger.SetOutput(os.Stdout)

	// Middleware
	log.Println("[STARTUP] Setting up middleware...")
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} [${method}] ${uri} -> ${status} (${latency_human}) from ${remote_ip}\n",
		Output: os.Stdout,
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	log.Println("[STARTUP] Middleware configured successfully")

	// Swagger documentation route
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Health check route
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "ok",
			"service": "product-service",
			"version": getEnv("SERVICE_VERSION", "1.0.0"),
			"message": "Product service is running",
		})
	})

	// API routes
	log.Println("[STARTUP] Setting up API routes...")

	// Product routes
	products := e.Group("/products")
	products.POST("", productHandler.CreateProduct)
	products.GET("", productHandler.ListProducts)

	// Internal service endpoint with service authentication
	products.PATCH("/:id/hold-stock", productHandler.UpdateOnHoldStock, pkgMiddleware.ServiceAuthMiddleware())
	log.Println("[STARTUP] Routes configured successfully")

	// Start server
	port := getEnv("PORT", "8081")
	log.Printf("[STARTUP] ========================")
	log.Printf("[STARTUP] PRODUCT SERVICE READY!")
	log.Printf("[STARTUP] ========================")
	fmt.Printf("Product Service starting on port %s\n", port)
	log.Printf("[INFO] Health check: http://localhost:%s/health", port)
	log.Printf("[INFO] API Base URL: http://localhost:%s/api/v1", port)
	log.Printf("[INFO] Swagger UI: http://localhost:%s/swagger/index.html", port)
	log.Println("[STARTUP] Server starting... (Press Ctrl+C to stop)")

	if err := e.Start(":" + port); err != nil {
		log.Fatal("Product Service failed to start:", err)
	}
}

// loadEnvFile loads environment variables from .env files
func loadEnvFile() {
	// Try to load from configs/product/.env first
	envPath := filepath.Join("configs", "product", ".env")
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
	dbName := getEnv("DB_NAME", "edot_product")

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

// Helper function to safely get working directory
func getCurrentWorkingDir() string {
	if wd, err := os.Getwd(); err == nil {
		return wd
	}
	return "unknown"
}
