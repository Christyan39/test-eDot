package user

import (
	"net/http"
	"os"
	"strings"

	models "github.com/Christyan39/test-eDot/internal/models/user"
	usecases "github.com/Christyan39/test-eDot/internal/usecases/user"
	"github.com/Christyan39/test-eDot/pkg/envelope"
	"github.com/labstack/echo/v4"
)

// AuthHandler handles HTTP requests for authentication
type AuthHandler struct {
	authUsecase     usecases.AuthUsecaseInterface
	envelopeService *envelope.EnvelopeService
}

type AuthHandlerInterface interface {
	HandleEnvelopeLogin(c echo.Context) error
	HandleDirectLogin(c echo.Context) error
	CreateEnvelope(c echo.Context) error
}

// NewAuthHandler creates new authentication handler
func NewAuthHandler(authUsecase usecases.AuthUsecaseInterface) AuthHandlerInterface {
	// Get envelope secret from environment
	envelopeSecret := os.Getenv("ENVELOPE_SECRET")
	if envelopeSecret == "" {
		envelopeSecret = "default-envelope-secret-change-in-production"
	}

	return &AuthHandler{
		authUsecase:     authUsecase,
		envelopeService: envelope.NewEnvelopeService(envelopeSecret),
	}
}

// HandleEnvelopeLogin processes secure envelope login
func (h *AuthHandler) HandleEnvelopeLogin(c echo.Context) error {
	var req models.LoginRequest
	var secureReq models.SecureLoginRequest
	err := c.Bind(&secureReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Decrypt envelope (simplified)
	err = h.envelopeService.DecryptData(secureReq.Envelope.Data, &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid envelope",
		})
	}

	// Validate credentials
	if req.Identifier == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Both identifier and password are required",
		})
	}

	return h.processLogin(c, &req)
}

// handleDirectLogin processes direct login (for backward compatibility)
func (h *AuthHandler) HandleDirectLogin(c echo.Context) error {
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Validate credentials
	if req.Identifier == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Both identifier and password are required",
		})
	}

	return h.processLogin(c, &req)
}

// processLogin handles the actual login logic
func (h *AuthHandler) processLogin(c echo.Context, req *models.LoginRequest) error {
	response, err := h.authUsecase.Login(req)
	if err != nil {
		// Don't expose detailed error messages for security
		if strings.Contains(err.Error(), "invalid credentials") {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid credentials",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Authentication failed",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":    response,
		"message": "Login successful",
	})
}

// CreateEnvelope handles POST /auth/create-envelope
// @Summary Create secure envelope
// @Description Create a secure envelope for sensitive data transmission
// @Tags auth
// @Accept json
// @Produce json
// @Param data body interface{} true "Data to encrypt"
// @Success 200 {object} map[string]interface{} "Envelope created successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Envelope creation failed"
// @Router /auth/create-envelope [post]
func (h *AuthHandler) CreateEnvelope(c echo.Context) error {
	var data interface{}
	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Create envelope (simplified)
	encryptedData, err := h.envelopeService.EncryptData(data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create envelope",
		})
	}

	envelope := map[string]interface{}{
		"data": encryptedData,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"envelope": envelope,
		},
		"message": "Envelope created successfully",
	})
}
