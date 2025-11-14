package product

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	productModel "github.com/Christyan39/test-eDot/internal/models/product"
	productUsecase "github.com/Christyan39/test-eDot/internal/usecases/product"
)

// ProductHandler defines the product HTTP handler interface
type ProductHandler interface {
	CreateProduct(c echo.Context) error
	ListProducts(c echo.Context) error
	UpdateOnHoldStock(c echo.Context) error
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

// UpdateOnHoldStock updates product on-hold stock
// @Summary Update product on-hold stock (Internal Service Only)
// @Description Update the on-hold stock quantity for a specific product (for inventory reservation). This endpoint is restricted to internal service access only.
// @Tags products
// @Accept json
// @Produce json
// @Param X-API-Key header string true "API key for internal service access"
// @Param id path int true "Product ID" minimum(1)
// @Param on_hold_stock body map[string]int true "On-hold stock data" example({"on_hold_stock": 5})
// @Success 200 {object} map[string]string "On-hold stock updated successfully"
// @Failure 400 {object} map[string]string "Bad request - invalid input or on-hold stock exceeds available stock"
// @Failure 401 {object} map[string]string "Unauthorized - API key required"
// @Failure 403 {object} map[string]string "Forbidden - invalid API key"
// @Failure 404 {object} map[string]string "Product not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /products/{id}/hold-stock [patch]
func (h *productHandler) UpdateOnHoldStock(c echo.Context) error {
	// Force log output to stdout
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Printf("[UpdateOnHoldStock] Invalid product ID: %s, error: %v", idParam, err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid product ID",
		})
	}

	var req productModel.Product
	if err := c.Bind(&req); err != nil {
		log.Printf("[UpdateOnHoldStock] Failed to bind request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	if req.OnHoldStock <= 0 {
		log.Printf("[UpdateOnHoldStock] Invalid on_hold_stock value: %d (cannot be negative or zero)", req.OnHoldStock)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "On-hold stock cannot be negative or zero",
		})
	}

	err = h.productUsecase.UpdateOnHoldStock(c.Request().Context(), id, req.OnHoldStock)
	if err != nil {
		if err.Error() == "product not found" {
			log.Printf("[UpdateOnHoldStock] Product not found: ID=%d", id)
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Product not found",
			})
		}
		if strings.Contains(err.Error(), "cannot exceed available stock") {
			log.Printf("[UpdateOnHoldStock] Stock validation failed: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}
		log.Printf("[UpdateOnHoldStock] Usecase error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update on-hold stock",
		})
	}

	log.Printf("[UpdateOnHoldStock] Successfully updated on-hold stock for product %d to %d", id, req.OnHoldStock)
	return c.JSON(http.StatusOK, map[string]string{
		"message": "On-hold stock updated successfully",
	})
}
