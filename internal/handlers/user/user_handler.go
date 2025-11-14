package user

import (
	"net/http"

	models "github.com/Christyan39/test-eDot/internal/models/user"
	usecases "github.com/Christyan39/test-eDot/internal/usecases/user"
	"github.com/Christyan39/test-eDot/pkg/envelope"
	"github.com/labstack/echo/v4"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	userUsecase     usecases.UserUsecaseInterface
	envelopeService *envelope.EnvelopeService
}

// NewUserHandler creates new user handler
func NewUserHandler(userUsecase usecases.UserUsecaseInterface) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

type UserHandlerInterface interface {
	CreateUser(c echo.Context) error
	HandleEnvelopeLogin(c echo.Context) error
	HandleDirectLogin(c echo.Context) error
	CreateEnvelope(c echo.Context) error
}

// CreateUser handles POST /users
// @Summary Create a new user
// @Description Create a new user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param user body user.CreateUserRequest true "User creation data"
// @Success 201 {object} map[string]interface{} "User created successfully"
// @Failure 400 {object} map[string]string "Invalid request body or validation error"
// @Router /users [post]
func (h *UserHandler) CreateUser(c echo.Context) error {
	var req models.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	err := h.userUsecase.CreateUser(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "User created successfully",
	})
}
