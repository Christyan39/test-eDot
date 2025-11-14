package order

import (
	"fmt"
	"log"

	"github.com/Christyan39/test-eDot/internal/clients"
	orderModel "github.com/Christyan39/test-eDot/internal/models/order"
	orderRepo "github.com/Christyan39/test-eDot/internal/repositories/order"
	"github.com/Christyan39/test-eDot/pkg/config"
)

// OrderUsecase defines the order business logic interface
type OrderUsecase interface {
	CreateOrder(req *orderModel.CreateOrderRequest) (*orderModel.CreateOrderResponse, error)
}

// orderUsecase implements OrderUsecase
type orderUsecase struct {
	orderRepo     orderRepo.OrderRepository
	productClient *clients.ProductServiceClient
}

// NewOrderUsecase creates a new order usecase
func NewOrderUsecase(orderRepo orderRepo.OrderRepository) OrderUsecase {
	productServiceURL := config.GetEnv("PRODUCT_SERVICE_URL", "http://localhost:8081")
	apiKey := config.GetEnv("PRODUCT_SERVICE_API_KEY", "")

	return &orderUsecase{
		orderRepo:     orderRepo,
		productClient: clients.NewProductServiceClient(productServiceURL, apiKey),
	}
}

// CreateOrder creates a new order with multiple products
func (u *orderUsecase) CreateOrder(req *orderModel.CreateOrderRequest) (*orderModel.CreateOrderResponse, error) {
	// Validate request
	if req.UserID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}
	if req.ShopID <= 0 {
		return nil, fmt.Errorf("invalid shop ID")
	}
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("at least one item is required")
	}

	// Validate each item and collect product information via HTTP calls
	itemPrices := make(map[int]float64)
	var orderIDs []int
	var stockUpdates []clients.UpdateProductRequest
	totalPrice := 0.0

	log.Printf("[CreateOrder] Making HTTP calls to Product Service at %s", u.productClient.BaseURL)

	for _, item := range req.Items {
		if item.ProductID <= 0 {
			return nil, fmt.Errorf("invalid product ID: %d", item.ProductID)
		}
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("quantity must be greater than 0 for product %d", item.ProductID)
		}

		// Get product via HTTP call to Product Service
		log.Printf("[CreateOrder] Fetching product %d from Product Service", item.ProductID)
		product, err := u.productClient.GetProductByID(item.ProductID)
		if err != nil {
			log.Printf("Failed to get product with ID %d from Product Service: %v", item.ProductID, err)
			return nil, fmt.Errorf("product %d not found", item.ProductID)
		}

		// Check if product belongs to the specified shop
		if product.ShopID != req.ShopID {
			return nil, fmt.Errorf("product %d does not belong to shop %d", item.ProductID, req.ShopID)
		}

		// Check stock availability
		availableStock := product.Stock - product.OnHoldStock
		if item.Quantity > availableStock {
			return nil, fmt.Errorf("insufficient stock for product %d: requested %d, available %d",
				item.ProductID, item.Quantity, availableStock)
		}

		// Calculate item price and add to total
		itemPrice := product.Price * float64(item.Quantity)
		itemPrices[item.ProductID] = product.Price
		totalPrice += itemPrice

		// Prepare stock update for Product Service
		stockUpdates = append(stockUpdates, clients.UpdateProductRequest{
			OnHoldStock: product.OnHoldStock + item.Quantity,
			Stock:       product.Stock,
		})

		log.Printf("[CreateOrder] Product %d validated: Price=%.2f, Available=%d, Requested=%d",
			item.ProductID, product.Price, availableStock, item.Quantity)
	}

	// Create orders for all products
	log.Printf("[CreateOrder] Creating order records in database")
	orderGroup, err := u.orderRepo.CreateMultiple(req, itemPrices)
	if err != nil {
		log.Printf("Failed to create orders: %v", err)
		return nil, fmt.Errorf("failed to create orders: %w", err)
	}

	// Update stock for all products via HTTP calls to Product Service
	log.Printf("[CreateOrder] Updating stock for %d products via HTTP calls", len(req.Items))
	for i, item := range req.Items {
		err = u.productClient.UpdateProductStock(item.ProductID, &stockUpdates[i])
		if err != nil {
			log.Printf("Failed to update product stock for product ID %d via HTTP: %v", item.ProductID, err)
			// Note: In a production system, this should be handled with proper transaction rollback
			return nil, fmt.Errorf("failed to reserve product stock for product %d", item.ProductID)
		}
		orderIDs = append(orderIDs, orderGroup.Orders[i].ID)
		log.Printf("[CreateOrder] Updated stock for product %d: OnHold=%d", item.ProductID, stockUpdates[i].OnHoldStock)
	}

	log.Printf("Multi-product order created successfully: OrderIDs=%v, UserID=%d, ShopID=%d, TotalItems=%d, TotalPrice=%.2f",
		orderIDs, req.UserID, req.ShopID, orderGroup.TotalItems, totalPrice)

	return &orderModel.CreateOrderResponse{
		OrderIDs:   orderIDs,
		TotalPrice: totalPrice,
		TotalItems: orderGroup.TotalItems,
		Status:     orderGroup.Status,
		Message:    "Multi-product order created successfully",
	}, nil
}
