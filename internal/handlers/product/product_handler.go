package product

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	productModel "github.com/Christyan39/test-eDot/internal/models/product"
	productUsecase "github.com/Christyan39/test-eDot/internal/usecases/product"
)

// ProductHandler defines the product HTTP handler interface
type ProductHandler interface {
	CreateProduct(c echo.Context) error
	ListProducts(c echo.Context) error
	HoldStockInBulk(c echo.Context) error
}

// productHandler implements ProductHandler
type productHandler struct {
	productUsecase productUsecase.ProductUsecase
}

// NewProductHandler creates a new product handler
func NewProductHandler(productUsecase productUsecase.ProductUsecase) ProductHandler {
	return &productHandler{
		productUsecase: productUsecase,
	}
}

// CreateProduct creates a new product
// @Summary Create a new product
// @Description Create a new product in the marketplace with shop information
// @Tags products
// @Accept json
// @Produce json
// @Param product body productModel.CreateProductRequest true "Product creation data"
// @Success 201 {object} map[string]string "Product created successfully"
// @Failure 400 {object} map[string]string "Bad request - invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /products [post]
// @Security BearerAuth
func (h *productHandler) CreateProduct(c echo.Context) error {
	var req productModel.CreateProductRequest
	if err := c.Bind(&req); err != nil {
		log.Printf("[CreateProduct] Failed to bind request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	err := h.productUsecase.CreateProduct(c.Request().Context(), &req)
	if err != nil {
		log.Printf("[CreateProduct] Usecase error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create product",
		})
	}

	log.Printf("[CreateProduct] Product created successfully: %s", req.Name)
	return c.JSON(http.StatusCreated, map[string]string{
		"message": "Product created successfully",
	})
}

// ListProducts retrieves products with filtering and pagination
// @Summary List products with filters and pagination
// @Description Get a paginated list of products with optional filtering by shop, price, status, and search
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)" minimum(1)
// @Param limit query int false "Items per page (default: 10, max: 100)" minimum(1) maximum(100)
// @Param shop_id query int false "Filter by shop ID" minimum(1)
// @Param min_price query number false "Minimum price filter" minimum(0)
// @Param max_price query number false "Maximum price filter" minimum(0)
// @Param status query string false "Filter by status" Enums(active,inactive,discontinued)
// @Param search query string false "Search in product name and description" maxlength(100)
// @Success 200 {object} productModel.ProductListResponse "Successfully retrieved products"
// @Failure 400 {object} map[string]string "Bad request - invalid parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /products [get]
// @Security BearerAuth
func (h *productHandler) ListProducts(c echo.Context) error {
	// IMMEDIATE VISIBLE OUTPUT
	var req productModel.ProductListRequest
	if err := c.Bind(&req); err != nil {
		log.Printf("[ListProducts] Failed to bind parameters: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request parameters",
		})
	}

	response, err := h.productUsecase.ListProducts(c.Request().Context(), &req)
	if err != nil {
		log.Printf("[ListProducts] Usecase error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to list products",
		})
	}

	log.Printf("[ListProducts] Successfully retrieved %d products (total: %d)", len(response.Products), response.Total)
	return c.JSON(http.StatusOK, response)
}

// HoldStockInBulk updates on-hold stock for multiple products in bulk
// @Summary Update on-hold stock for multiple products in bulk
// @Description Update the on-hold stock for multiple products within a single request
// @Tags products
// @Accept json
// @Produce json
// @Param products body []productModel.UpdateProductRequest true "List of products to update on-hold stock"
// @Success 200 {object} map[string]string "Successfully updated on-hold stock in bulk"
// @Failure 400 {object} map[string]string "Bad request - invalid parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /products/hold-stock [patch]
// @Security BearerAuth
func (h *productHandler) HoldStockInBulk(c echo.Context) error {
	// Force log output to stdout
	var req productModel.HoldStockRequest
	if err := c.Bind(&req); err != nil {
		log.Printf("[HoldStockInBulk] Failed to bind request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	err := h.productUsecase.HoldStockInBulk(c.Request().Context(), &req)
	if err != nil {
		log.Printf("[UpdateOnHoldStockInBulk] Usecase error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update on-hold stock in bulk",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "On-hold stock updated successfully in bulk",
	})
}
