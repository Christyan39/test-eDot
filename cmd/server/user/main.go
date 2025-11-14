package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/Christyan39/test-eDot/docs"
	handlers "github.com/Christyan39/test-eDot/internal/handlers/user"
	repositories "github.com/Christyan39/test-eDot/internal/repositories/user"
	usecases "github.com/Christyan39/test-eDot/internal/usecases/user"
	"github.com/Christyan39/test-eDot/pkg/config"
	"github.com/Christyan39/test-eDot/pkg/database"
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
	// Configure logging to output to stdout with timestamps
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Load environment variables from .env file
	config.LoadEnvFile("user")

	// Connect to MySQL database
	db, err := database.InitMySQL("user")
	if err != nil {
		log.Fatalf("Warning: Failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	// Initialize layers
	userRepo := repositories.NewUserRepository(db)
	userUsecase := usecases.NewUserUsecase(userRepo)
	userHandler := handlers.NewUserHandler(userUsecase)

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

	// Swagger documentation route
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Health check route
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
	auth.POST("/login", userHandler.HandleDirectLogin)
	auth.POST("/secure-login", userHandler.HandleEnvelopeLogin)
	auth.POST("/create-envelope", userHandler.CreateEnvelope)

	// User routes
	api.POST("/users", userHandler.CreateUser)

	// Start server
	port := config.GetEnv("PORT", "8080")
	log.Printf("[STARTUP] ========================")
	log.Printf("[STARTUP] USER SERVICE READY!")
	log.Printf("[STARTUP] ========================")
	log.Printf("User Service starting on port %s\n", port)
	log.Printf("[INFO] Health check: http://localhost:%s/health", port)
	log.Printf("[INFO] API Base URL: http://localhost:%s/api/v1", port)
	log.Printf("[INFO] Swagger UI: http://localhost:%s/swagger/index.html", port)
	log.Println("[STARTUP] Server starting... (Press Ctrl+C to stop)")

	if err := e.Start(":" + port); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
