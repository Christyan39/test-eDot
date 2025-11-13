package user

import (
	"database/sql"
	"fmt"
	"time"

	models "github.com/Christyan39/test-eDot/internal/models/user"
)

// UserRepositoryInterface defines user repository contract
type UserRepositoryInterface interface {
	GetByEmailOrPhone(identifier string) (*models.User, error)
	Create(req *models.CreateUserRequest) error
}

// UserRepository implements UserRepositoryInterface
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates new user repository
func NewUserRepository(db *sql.DB) UserRepositoryInterface {
	return &UserRepository{
		db: db,
	}
}

// GetByEmailOrPhone retrieves user by email or phone (for authentication)
func (r *UserRepository) GetByEmailOrPhone(identifier string) (*models.User, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database connection is not available")
	}

	query := `SELECT id, name, email, phone, password, created_at, updated_at FROM users WHERE email = ? OR phone = ?`

	user := &models.User{}
	err := r.db.QueryRow(query, identifier, identifier).Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}

// Create creates new user
func (r *UserRepository) Create(req *models.CreateUserRequest) error {
	if r.db == nil {
		return fmt.Errorf("database connection is not available")
	}

	query := `INSERT INTO users (name, email, phone, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`

	now := time.Now()
	_, err := r.db.Exec(query, req.Name, req.Email, req.Phone, req.Password, now, now)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	// Get the created user
	return nil
}
