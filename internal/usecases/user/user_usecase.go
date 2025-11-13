package user

import (
	"fmt"
	"regexp"

	models "github.com/Christyan39/test-eDot/internal/models/user"
	repositories "github.com/Christyan39/test-eDot/internal/repositories/user"
	"github.com/Christyan39/test-eDot/pkg/auth"
)

// UserUsecaseInterface defines user usecase contract
type UserUsecaseInterface interface {
	CreateUser(req *models.CreateUserRequest) error
}

// UserUsecase implements UserUsecaseInterface
type UserUsecase struct {
	userRepo repositories.UserRepositoryInterface
}

// NewUserUsecase creates new user usecase
func NewUserUsecase(userRepo repositories.UserRepositoryInterface) UserUsecaseInterface {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

// validatePhone validates phone number format (optional validation)
func (u *UserUsecase) validatePhone(phone string) error {
	if phone == "" {
		return nil // Phone is optional
	}

	// Indonesian phone validation - accepts formats: 081234567890 or +6281234567890
	phoneRegex := regexp.MustCompile(`^(\+628|08)[0-9]{8,11}$`)
	if !phoneRegex.MatchString(phone) {
		return fmt.Errorf("invalid phone format, use format: 081234567890 or +6281234567890")
	}

	return nil
}

// CreateUser creates new user
func (u *UserUsecase) CreateUser(req *models.CreateUserRequest) error {
	// Validation
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	if req.Password == "" {
		return fmt.Errorf("password is required")
	}
	if err := u.validatePhone(req.Phone); err != nil {
		return err
	}

	// Hash password before storing
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}
	req.Password = hashedPassword

	err = u.userRepo.Create(req)
	if err != nil {
		return fmt.Errorf("usecase error: %v", err)
	}

	return nil
}
