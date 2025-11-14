# Test eDot - Clean Architecture Go API

A Golang API application with clean 3-layer architecture (Handler, Usecase, Repository) using interfaces, Echo framework, and comprehensive Swagger documentation.

## Architecture

### 3-Layer Clean Architecture:

1. **Handler Layer** (`internal/handlers/user/`)
   - Handles HTTP requests and responses
   - Input validation from HTTP requests
   - Communicates with Usecase layer through interfaces

2. **Usecase Layer** (`internal/usecases/user/`)
   - Business logic and business validation
   - Orchestration between various repositories
   - Communicates with Repository layer through interfaces

3. **Repository Layer** (`internal/repositories/user/`)
   - Data access layer
   - Database operations (CRUD)
   - Interface implementation for database abstraction

### Project Structure:
```
.
├── cmd/server/user/main.go           # Application entry point
├── internal/
│   ├── handlers/user/                # HTTP handlers
│   │   ├── user_handler.go
│   │   └── auth_handler.go
│   ├── usecases/user/                # Business logic
│   │   ├── user_usecase.go
│   │   └── auth_usecase.go
│   ├── repositories/user/            # Data access
│   │   └── user_repository.go
│   └── models/user/                  # Data models
│       └── user.go
├── pkg/                              # Shared utilities
│   ├── auth/                         # JWT authentication
│   ├── database/                     # Database utilities
│   └── logger/                       # Logging utilities
├── docs/                             # API documentation
│   ├── swagger.json                  # Generated Swagger spec
│   ├── swagger.yaml                  # Generated Swagger spec
│   ├── docs.go                       # Generated Swagger docs
│   └── API_DOCUMENTATION.md          # API guide
├── configs/                          # Configuration files
└── .env.example                      # Environment variables template
```

## API Documentation

### 📚 Interactive Swagger UI
**Access the complete API documentation at:** [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

The Swagger UI provides:
- Interactive API testing
- Complete endpoint documentation
- Request/response examples
- Model definitions
- Authentication testing

### API Endpoints

#### System
- `GET /health` - Server health check
- `GET /swagger/*` - Interactive Swagger UI

#### Authentication
- `POST /api/v1/auth/login` - User login (returns JWT token)
- `GET /api/v1/auth/profile` - Get user profile (protected)

#### User Management
- `GET /api/v1/users` - Get all users
- `POST /api/v1/users` - Create new user
- `GET /api/v1/users/{id}` - Get user by ID
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user

## Quick Start

### 1. Clone and Setup
```bash
git clone <repository-url>
cd test-eDot
cp .env.example .env
# Edit .env with your database configuration
```

### 2. Install Dependencies
```bash
make deps
# or
go mod tidy
```

### 3. Setup Database (MySQL)
```sql
CREATE DATABASE IF NOT EXISTS testdb CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE testdb;

CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20),
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### 4. Run the Application
```bash
make run
# or
go run ./cmd/server/user
```

Server will start at `http://localhost:8080`

## Development Commands

```bash
# Build the application
make build

# Run the application  
make run

# Generate Swagger documentation
make swagger

# Run tests
make test

# Clean build artifacts
make clean

# Install development tools
make tools

# Show all available commands
make help
```

## Authentication

The API uses JWT (JSON Web Tokens) for authentication:

1. **Login** via `POST /api/v1/auth/login` to get a JWT token
2. **Include token** in requests: `Authorization: Bearer <your-jwt-token>`
3. **Access protected endpoints** like `/api/v1/auth/profile`

### Example Login:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "user@example.com",
    "password": "your-password"
  }'
```

## Example Requests

### Create User
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+1234567890",
    "password": "secure-password"
  }'
```

### Get All Users
```bash
curl -X GET http://localhost:8080/api/v1/users
```

### Update User (with Authentication)
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "John Updated",
    "email": "john.updated@example.com",
    "phone": "+1234567891"
  }'
```

## Technology Stack

- **Go 1.21** - Programming language
- **Echo v4** - Web framework
- **MySQL** - Database
- **JWT** - Authentication
- **Swagger** - API documentation
- **Clean Architecture** - Project structure

## Dependencies

```go
require (
    github.com/go-sql-driver/mysql v1.7.1
    github.com/golang-jwt/jwt v3.2.2+incompatible
    github.com/labstack/echo/v4 v4.11.3
    github.com/swaggo/echo-swagger v1.4.1
    github.com/swaggo/swag v1.16.6
    golang.org/x/crypto v0.32.0
)
```

## Key Features

✅ **Clean Architecture** - 3-layer separation with interfaces  
✅ **Comprehensive API Documentation** - Interactive Swagger UI  
✅ **JWT Authentication** - Secure token-based auth  
✅ **MySQL Integration** - Full database connectivity  
✅ **Input Validation** - Request validation at handler level  
✅ **Error Handling** - Proper error propagation and responses  
✅ **Middleware Support** - CORS, logging, authentication  
✅ **Development Tools** - Makefile, auto-documentation generation  

## Contributing

1. Follow the existing project structure
2. Add Swagger annotations for new endpoints
3. Regenerate documentation: `make swagger`
4. Update tests for new features
5. Follow clean architecture principles

## License

This project is licensed under the MIT License.