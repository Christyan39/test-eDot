package main

import (
	"fmt"
	"log"
	"net/http"

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
	// Load environment variables from .env file
	config.LoadEnvFile("user")

	// Connect to MySQL database
	db, err := database.InitMySQL("user")
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
	port := config.GetEnv("PORT", "8080")
	fmt.Printf("Server starting on port %s\n", port)
	fmt.Printf("Health check: http://localhost:%s/health\n", port)
	fmt.Printf("API Base URL: http://localhost:%s/api/v1\n", port)

	if err := e.Start(":" + port); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
