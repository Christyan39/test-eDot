# Test eDot - Clean Architecture Go API

Aplikasi API Golang dengan arsitektur 3 layer (Handler, Usecase, Repository) menggunakan interface dan Echo framework.

## Arsitektur

### 3 Layer Architecture:

1. **Handler Layer** (`internal/handlers/`)
   - Menangani HTTP requests dan responses
   - Validasi input dari HTTP request
   - Komunikasi dengan Usecase layer melalui interface

2. **Usecase Layer** (`internal/usecases/`)
   - Business logic dan validasi bisnis
   - Orchestration antara berbagai repository
   - Komunikasi dengan Repository layer melalui interface

3. **Repository Layer** (`internal/repositories/`)
   - Data access layer
   - Database operations (CRUD)
   - Implementasi interface untuk abstraksi database

### Struktur Folder:
```
.
├── main.go                           # Entry point aplikasi
├── internal/
│   ├── handlers/                     # HTTP handlers
│   │   └── user_handler.go
│   ├── usecases/                     # Business logic
│   │   └── user_usecase.go
│   ├── repositories/                 # Data access
│   │   └── user_repository.go
│   └── models/                       # Data models
│       └── user.go
├── configs/
│   └── database.sql                  # Database schema
└── .env.example                      # Environment variables
```

## Interface Communication

### Repository Interface:
```go
type UserRepositoryInterface interface {
    GetAll() ([]*models.User, error)
    GetByID(id int) (*models.User, error)
    Create(req *models.CreateUserRequest) (*models.User, error)
    Update(id int, req *models.UpdateUserRequest) (*models.User, error)
    Delete(id int) error
}
```

### Usecase Interface:
```go
type UserUsecaseInterface interface {
    GetAllUsers() ([]*models.User, error)
    GetUserByID(id int) (*models.User, error)
    CreateUser(req *models.CreateUserRequest) (*models.User, error)
    UpdateUser(id int, req *models.UpdateUserRequest) (*models.User, error)
    DeleteUser(id int) error
}
```

## API Endpoints

### Health Check
- `GET /health` - Status server

### Users API
- `GET /api/v1/users` - Get all users
- `POST /api/v1/users` - Create new user
- `GET /api/v1/users/{id}` - Get user by ID
- `PUT /api/v1/users/{id}` - Update user by ID
- `DELETE /api/v1/users/{id}` - Delete user by ID

## Setup dan Jalankan

### 1. Setup Database (MySQL)
```sql
-- Jalankan script di configs/database.sql
CREATE DATABASE IF NOT EXISTS testdb CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE testdb;

CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### 2. Environment Variables
```bash
cp .env.example .env
# Edit .env sesuai konfigurasi database Anda
```

### 3. Jalankan Aplikasi
```bash
go run main.go
```

Server akan berjalan di `http://localhost:8080`

## Contoh Request

### Create User
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com"
  }'
```

### Get All Users
```bash
curl -X GET http://localhost:8080/api/v1/users
```

### Get User by ID
```bash
curl -X GET http://localhost:8080/api/v1/users/1
```

### Update User
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Smith",
    "email": "johnsmith@example.com"
  }'
```

### Delete User
```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

## Dependencies

- `database/sql` - Standard Go SQL driver interface
- `github.com/go-sql-driver/mysql` - MySQL driver for Go
- `github.com/labstack/echo/v4` - Echo web framework

## Catatan Implementation

1. **Dependency Injection**: Setiap layer menerima dependency melalui constructor
2. **Interface Abstraction**: Setiap layer berkomunikasi melalui interface
3. **Error Handling**: Proper error propagation dari repository ke handler
4. **Validation**: Business validation di usecase layer, input validation di handler
5. **Echo Framework**: Menggunakan Echo web framework untuk HTTP handling