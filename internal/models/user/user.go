package user

import "time"

// User represents a user entity
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Password  string    `json:"-"` // Never include password in JSON response
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserRequest represents request to create user
type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// UpdateUserRequest represents request to update user
type UpdateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password,omitempty"` // Optional field for password update
}

// Envelope structures for secure data transmission (simplified)
type LoginEnvelope struct {
	Data string `json:"data"` // Base64 encoded encrypted login data
}

// LoginRequest represents authentication request (internal use)
type LoginRequest struct {
	Identifier string `json:"identifier"` // Can be email or phone
	Password   string `json:"password"`
}

// LoginResponse represents authentication response
type LoginResponse struct {
	Token string    `json:"token"`
	User  *AuthUser `json:"user"`
}

// SecureLoginRequest represents the secure envelope for login
type SecureLoginRequest struct {
	Envelope LoginEnvelope `json:"envelope"`
}

// AuthUser represents authenticated user info for JWT claims
type AuthUser struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}
