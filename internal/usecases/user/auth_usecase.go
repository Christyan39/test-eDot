package user

import (
	"fmt"

	models "github.com/Christyan39/test-eDot/internal/models/user"
	repositories "github.com/Christyan39/test-eDot/internal/repositories/user"
	"github.com/Christyan39/test-eDot/pkg/auth"
)

// AuthUsecaseInterface defines authentication usecase contract
type AuthUsecaseInterface interface {
	Login(req *models.LoginRequest) (*models.LoginResponse, error)
}

// AuthUsecase implements AuthUsecaseInterface
type AuthUsecase struct {
	userRepo repositories.UserRepositoryInterface
}

// NewAuthUsecase creates new authentication usecase
func NewAuthUsecase(userRepo repositories.UserRepositoryInterface) AuthUsecaseInterface {
	return &AuthUsecase{
		userRepo: userRepo,
	}
}

// Login authenticates user with email/phone and password
func (a *AuthUsecase) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// Validation
	if req.Identifier == "" {
		return nil, fmt.Errorf("email or phone is required")
	}
	if req.Password == "" {
		return nil, fmt.Errorf("password is required")
	}

	// Find user by email or phone
	user, err := a.userRepo.GetByEmailOrPhone(req.Identifier)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %v", err)
	}
	if user == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password
	if !auth.CheckPassword(req.Password, user.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	// Create AuthUser from User (exclude password)
	authUser := &models.AuthUser{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
	}

	// Return response with token and user info (password excluded)
	return &models.LoginResponse{
		Token: token,
		User:  authUser,
	}, nil
}
