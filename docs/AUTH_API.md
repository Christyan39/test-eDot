# Authentication API Documentation

This API provides secure authentication using JWT tokens with request envelope protection against MITM attacks.

## Security Features

- **Password Hashing**: Uses bcrypt for secure password storage
- **JWT Tokens**: Secure token-based authentication with 24-hour expiration
- **Request Envelope**: All sensitive requests require proper formatting to prevent MITM attacks
- **Protected Routes**: Most API endpoints require authentication
- **Indonesian Phone Validation**: Supports formats: `081234567890` or `+6281234567890`

## Authentication Endpoints

### 1. Login (POST /auth/login)
Authenticate with email/phone and password.

**Request:**
```json
{
  "identifier": "john@example.com", // Can be email or phone
  "password": "password123"
}
```

**Response:**
```json
{
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com", 
      "phone": "+6281234567890",
      "created_at": "2025-11-13T10:00:00Z",
      "updated_at": "2025-11-13T10:00:00Z"
    }
  },
  "message": "Login successful"
}
```

### 2. Get Profile (GET /auth/profile)
Get current user profile (requires authentication).

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+6281234567890"
  },
  "message": "Profile retrieved successfully"
}
```

## User Management Endpoints

### 1. Create User (POST /api/v1/users) - Public
```json
{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "phone": "+6281234567891",
  "password": "securepassword123"
}
```

### 2. Get All Users (GET /api/v1/users) - Protected
**Headers:** `Authorization: Bearer <jwt_token>`

### 3. Get User by ID (GET /api/v1/users/:id) - Protected
**Headers:** `Authorization: Bearer <jwt_token>`

### 4. Update User (PUT /api/v1/users/:id) - Protected
**Headers:** `Authorization: Bearer <jwt_token>`
```json
{
  "name": "Updated Name",
  "email": "updated@example.com", 
  "phone": "081234567999",
  "password": "newpassword123" // Optional
}
```

### 5. Delete User (DELETE /api/v1/users/:id) - Protected
**Headers:** `Authorization: Bearer <jwt_token>`

## Sample Credentials

Use these for testing:
- **Email**: `john@example.com` or **Phone**: `+6281234567890`
- **Password**: `password123`

## Phone Number Validation

Supported formats:
- `081234567890` (local format)
- `+6281234567890` (international format)

## Environment Variables

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=testdb

# JWT Configuration  
JWT_SECRET=your-super-secret-jwt-key-change-in-production

# Server Configuration
PORT=8080
```

## Security Notes

1. **JWT Secret**: Change the default JWT secret in production
2. **Password Requirements**: Implement password strength requirements as needed
3. **Rate Limiting**: Consider adding rate limiting for authentication endpoints
4. **HTTPS**: Always use HTTPS in production
5. **Request Validation**: All requests are validated for proper envelope structure
6. **Error Messages**: Authentication errors return generic messages to prevent user enumeration

## Error Responses

```json
{
  "error": "Invalid credentials"           // 401 Unauthorized
}
{
  "error": "Authorization header required" // 401 Unauthorized  
}
{
  "error": "Invalid or expired token"      // 401 Unauthorized
}
{
  "error": "Invalid request body"          // 400 Bad Request
}
```