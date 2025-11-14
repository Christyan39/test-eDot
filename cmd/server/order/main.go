package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/Christyan39/test-eDot/docs"
	orderHandlers "github.com/Christyan39/test-eDot/internal/handlers/order"
	orderRepositories "github.com/Christyan39/test-eDot/internal/repositories/order"
	orderUsecases "github.com/Christyan39/test-eDot/internal/usecases/order"
	"github.com/Christyan39/test-eDot/pkg/config"
	"github.com/Christyan39/test-eDot/pkg/database"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Order Service API
// @version 1.0
// @description API documentation for Order Service - manages orders and order processing
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8082
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	log.Println("=== ORDER SERVICE STARTING ===")
	log.Println("[STARTUP] Order Service initialization started")
	log.Printf("[STARTUP] Process ID: %d", os.Getpid())
	if wd, err := os.Getwd(); err == nil {
		log.Printf("[STARTUP] Working directory: %s", wd)
	}

	// Load environment variables
	log.Println("[STARTUP] Loading environment variables...")
	config.LoadEnvFile("order")

	// Initialize database connection
	log.Println("[STARTUP] Connecting to database...")
	db, err := database.InitMySQL("order")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("[STARTUP] Database connection established successfully")

	// Initialize layers
	log.Println("[STARTUP] Initializing repository layer...")
	orderRepo := orderRepositories.NewOrderRepository(db)
	log.Println("[STARTUP] Initializing usecase layer (with HTTP Product Service client)...")
	orderUsecase := orderUsecases.NewOrderUsecase(orderRepo)
	log.Println("[STARTUP] Initializing handler layer...")
	orderHandler := orderHandlers.NewOrderHandler(orderUsecase)

	// Initialize Echo
	log.Println("[STARTUP] Initializing Echo web framework...")
	e := echo.New()
	e.Debug = true

	// Configure Echo logger
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
			"service": "order-service",
			"version": config.GetEnv("SERVICE_VERSION", "1.0.0"),
			"message": "Order service is running",
		})
	})

	// API routes
	log.Println("[STARTUP] Setting up API routes...")

	// Order routes
	orders := e.Group("/orders")
	orders.POST("", orderHandler.CreateOrder)

	log.Println("[STARTUP] Routes configured successfully")

	// Start server
	port := config.GetEnv("PORT", "8082")
	log.Printf("[STARTUP] ========================")
	log.Printf("[STARTUP] ORDER SERVICE READY!")
	log.Printf("[STARTUP] ========================")
	fmt.Printf("Order Service starting on port %s\n", port)
	log.Printf("[INFO] Health check: http://localhost:%s/health", port)
	log.Printf("[INFO] API Base URL: http://localhost:%s", port)
	log.Printf("[INFO] Swagger UI: http://localhost:%s/swagger/index.html", port)
	log.Println("[STARTUP] Server starting... (Press Ctrl+C to stop)")

	if err := e.Start(":" + port); err != nil {
		log.Fatal("Order Service failed to start:", err)
	}
}
