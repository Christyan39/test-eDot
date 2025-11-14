package order

import (
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	orderModel "github.com/Christyan39/test-eDot/internal/models/order"
	orderUsecase "github.com/Christyan39/test-eDot/internal/usecases/order"
)

// OrderHandler defines the order HTTP handler interface
type OrderHandler interface {
	CreateOrder(c echo.Context) error
}

// orderHandler implements OrderHandler
type orderHandler struct {
	orderUsecase orderUsecase.OrderUsecase
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(orderUsecase orderUsecase.OrderUsecase) OrderHandler {
	return &orderHandler{
		orderUsecase: orderUsecase,
	}
}

// CreateOrder creates a new order
// @Summary Create a new order
// @Description Create a new order for a product
// @Tags orders
// @Accept json
// @Produce json
// @Param order body orderModel.CreateOrderRequest true "Order creation data"
// @Success 201 {object} orderModel.CreateOrderResponse "Order created successfully"
// @Failure 400 {object} map[string]string "Bad request - invalid input or validation failed"
// @Failure 404 {object} map[string]string "Product not found"
// @Failure 409 {object} map[string]string "Insufficient stock"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /orders [post]
// @Security BearerAuth
func (h *orderHandler) CreateOrder(c echo.Context) error {
	// Bind request data
	var req orderModel.CreateOrderRequest
	if err := c.Bind(&req); err != nil {
		log.Printf("[CreateOrder] Failed to bind request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Validate request
	if req.UserID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid user ID",
		})
	}
	if req.ShopID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid shop ID",
		})
	}
	if len(req.Items) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "at least one item is required",
		})
	}

	log.Printf("[CreateOrder] Order data: UserID=%d, ShopID=%d, Items=%d",
		req.UserID, req.ShopID, len(req.Items)) // Call usecase to create order
	response, err := h.orderUsecase.CreateOrder(c.Request().Context(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Printf("[CreateOrder] Product not found: %v", err)
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}
		if strings.Contains(err.Error(), "insufficient stock") {
			log.Printf("[CreateOrder] Insufficient stock: %v", err)
			return c.JSON(http.StatusConflict, map[string]string{
				"error": err.Error(),
			})
		}
		if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "must be") {
			log.Printf("[CreateOrder] Validation error: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}
		log.Printf("[CreateOrder] Usecase error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create order",
		})
	}

	log.Printf("[CreateOrder] Successfully created orders: IDs=%v, TotalItems=%d, TotalPrice=%.2f", response.OrderIDs, response.TotalItems, response.TotalPrice)
	return c.JSON(http.StatusCreated, response)
}
