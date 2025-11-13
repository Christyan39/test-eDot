# API Documentation

This document describes the REST API endpoints for the Test-eDot user management system.

## Swagger UI

The interactive Swagger UI is available at: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

The API uses JWT (JSON Web Tokens) for authentication. To access protected endpoints:

1. Login via `/auth/login` to get a JWT token
2. Include the token in the `Authorization` header: `Bearer <your-jwt-token>`

## Available Endpoints

### Authentication
- `POST /auth/login` - User login
- `GET /auth/profile` - Get user profile (requires authentication)

### Users
- `GET /users` - Get all users
- `POST /users` - Create a new user  
- `GET /users/{id}` - Get user by ID
- `PUT /users/{id}` - Update user
- `DELETE /users/{id}` - Delete user

### Health Check
- `GET /health` - Server health status

## Request/Response Examples

### Login Request
```json
{
  "identifier": "user@example.com",
  "password": "your-password"
}
```

### Create User Request  
```json
{
  "name": "John Doe",
  "email": "john@example.com", 
  "phone": "+1234567890",
  "password": "secure-password"
}
```

### Update User Request
```json
{
  "name": "John Updated",
  "email": "john.updated@example.com",
  "phone": "+1234567891",
  "password": "new-password"
}
```

## Response Format

All API responses follow this format:

### Success Response
```json
{
  "data": { /* response data */ },
  "message": "Operation successful"
}
```

### Error Response  
```json
{
  "error": "Error message description"
}
```

## Status Codes

- `200` - Success
- `201` - Created  
- `400` - Bad Request
- `401` - Unauthorized
- `404` - Not Found
- `500` - Internal Server Error

## Development

To regenerate Swagger documentation after making changes to API comments:

```bash
swag init -g cmd/server/user/main.go -o docs
```

Make sure to add Swagger annotations to your handler functions following the format:

```go
// @Summary Brief description
// @Description Detailed description  
// @Tags tag-name
// @Accept json
// @Produce json
// @Param name type dataType required "Description"
// @Success 200 {object} ResponseType "Success message"
// @Failure 400 {object} ErrorType "Error message"  
// @Router /endpoint [method]
func YourHandler(c echo.Context) error {
    // implementation
}
```