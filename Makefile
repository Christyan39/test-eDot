# Makefile for Test-eDot Microservices

.PHONY: build build-user build-product run run-user run-product swagger clean test help

# Build all services
build: build-user build-product build-order
	@echo "All services built successfully!"

# Build user service
build-user:
	@echo "Building user service..."
	go build -o bin/user ./cmd/server/user

# Build product service  
build-product:
	@echo "Building product service..."
	go build -o bin/product ./cmd/server/product

# Build order service
build-order:
	@echo "Building order service..."
	go build -o bin/order ./cmd/server/order

# Run user service
run-user: build-user
	@echo "Starting user service..."
	./bin/user

# Run product service
run-product: 
	go run ./cmd/server/product/main.go

# Run order service
run-order:
	go run ./cmd/server/order/main.go

# Run all services (legacy support)
run: run-user


# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/server/user/main.go -o docs

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f user product
	rm -f docs/docs.go docs/swagger.json docs/swagger.yaml

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Development server with auto-reload (requires air)
dev:
	@echo "Starting development server..."
	air

# Install development tools
tools:
	@echo "Installing development tools..."
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/cosmtrek/air@latest

# Docker build for user service
docker-build-user:
	@echo "Building user service Docker image..."
	docker build -f cmd/server/user/Dockerfile -t user-service:latest .

# Docker build for product service  
docker-build-product:
	@echo "Building product service Docker image..."
	docker build -f cmd/server/product/Dockerfile -t product-service:latest .

# Docker build for order service
docker-build-order:
	@echo "Building order service Docker image..."
	docker build -f cmd/server/order/Dockerfile -t order-service:latest .

# Build all Docker images
docker-build-all: docker-build-user docker-build-product docker-build-order
	@echo "All Docker images built successfully!"

# Push all Docker images (edit tags as needed for your registry)
docker-push-all:
	docker tag user-service:latest yourrepo/user-service:latest
	docker tag product-service:latest yourrepo/product-service:latest
	docker tag order-service:latest yourrepo/order-service:latest
	docker push yourrepo/user-service:latest
	docker push yourrepo/product-service:latest
	docker push yourrepo/order-service:latest

# Docker compose up
docker-compose-up: 
	@echo "Starting all services with docker-compose..."
	docker-compose up --build -d

docker-up:
	@echo "Starting all services with docker-compose..."
	make docker-build-all
	make docker-compose-up
	@echo "All services are up with docker-compose!"

# Docker compose down
docker-down:
	@echo "Stopping all services..."
	docker-compose down

# Database setup
db-create:
	@echo "Creating database..."
	mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS edot_user CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

db-migrate:
	@echo "Running database migrations..."
	mysql -u root -p edot_user < migration/001_create_users_table.sql

db-setup: db-create db-migrate
	@echo "Database setup complete!"

# Show help
help:
	@echo "Available commands:"
	@echo ""
	@echo "Building:"
	@echo "  build        - Build all services"
	@echo "  build-user   - Build user service only"
	@echo "  build-product - Build product service only"
	@echo ""
	@echo "Running:"
	@echo "  run-user     - Build and run user service (port 8080)"
	@echo "  run-product  - Build and run product service (port 8081)"
	@echo ""
	@echo "Development:"
	@echo "  swagger      - Generate Swagger documentation"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  deps         - Install/update dependencies"
	@echo "  dev          - Start development server with auto-reload"
	@echo "  tools        - Install development tools"
	@echo ""
	@echo "Database:"
	@echo "  db-create    - Create MySQL databases"
	@echo "  db-migrate   - Run database migrations"
	@echo "  db-setup     - Complete database setup (create + migrate)"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build-user - Build user service Docker image"
	@echo "  docker-build-product - Build product service Docker image"
	@echo "  docker-build-order - Build order service Docker image"
	@echo "  docker-build-all - Build all Docker images"
	@echo "  docker-push-all - Push all Docker images"
	@echo "  docker-up    - Start all services with docker-compose"
	@echo "  docker-down  - Stop all services"
	@echo ""
	@echo "  help         - Show this help message"

# Default target
all: deps swagger build