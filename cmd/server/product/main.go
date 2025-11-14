package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/Christyan39/test-eDot/docs"
	handlers "github.com/Christyan39/test-eDot/internal/handlers/product"
	repositories "github.com/Christyan39/test-eDot/internal/repositories/product"
	usecases "github.com/Christyan39/test-eDot/internal/usecases/product"
	"github.com/Christyan39/test-eDot/pkg/config"
	"github.com/Christyan39/test-eDot/pkg/database"
	pkgMiddleware "github.com/Christyan39/test-eDot/pkg/middleware"

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
	config.LoadEnvFile("product")

	// Connect to MySQL database
	db, err := database.InitMySQL("product")
	if err != nil {
		log.Fatalf("Warning: Failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	// Initialize layers
	productRepo := repositories.NewProductRepository(db)
	productUsecase := usecases.NewProductUsecase(productRepo)
	productHandler := handlers.NewProductHandler(productUsecase)

	// Initialize Echo
	e := echo.New()
	e.Debug = true

	// Configure Echo logger
	e.Logger.SetOutput(os.Stdout)

	// Middleware
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
			"version": config.GetEnv("SERVICE_VERSION", "1.0.0"),
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
	port := config.GetEnv("PORT", "8081")
	log.Printf("[STARTUP] ========================")
	log.Printf("[STARTUP] PRODUCT SERVICE READY!")
	log.Printf("[STARTUP] ========================")
	log.Printf("Product Service starting on port %s\n", port)
	log.Printf("[INFO] Health check: http://localhost:%s/health", port)
	log.Printf("[INFO] API Base URL: http://localhost:%s/api/v1", port)
	log.Printf("[INFO] Swagger UI: http://localhost:%s/swagger/index.html", port)
	log.Println("[STARTUP] Server starting... (Press Ctrl+C to stop)")

	if err := e.Start(":" + port); err != nil {
		log.Fatal("Product Service failed to start:", err)
	}
}
