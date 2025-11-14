package order

import (
	"context"
	"fmt"
	"log"

	"github.com/Christyan39/test-eDot/internal/clients"
	orderModel "github.com/Christyan39/test-eDot/internal/models/order"
	productModels "github.com/Christyan39/test-eDot/internal/models/product"
	orderRepo "github.com/Christyan39/test-eDot/internal/repositories/order"
	"github.com/Christyan39/test-eDot/pkg/config"
)

// OrderUsecase defines the order business logic interface
type OrderUsecase interface {
	CreateOrder(ctx context.Context, req *orderModel.CreateOrderRequest) (*orderModel.CreateOrderResponse, error)
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
func (u *orderUsecase) CreateOrder(ctx context.Context, req *orderModel.CreateOrderRequest) (*orderModel.CreateOrderResponse, error) {
	// Validate each item and collect product information via HTTP calls
	itemPrices := make(map[int]float64)
	var orderIDs []int
	var stockUpdates []productModels.UpdateProductRequest
	totalPrice := 0.0

	log.Printf("[CreateOrder] Making HTTP calls to Product Service at %s", u.productClient.BaseURL)

	itemIDs := make([]int, 0, len(req.Items))
	for _, item := range req.Items {
		itemIDs = append(itemIDs, item.ProductID)
	}

	// Fetch product details in batch
	products, err := u.productClient.GetProductByIDs(itemIDs)
	if err != nil {
		log.Printf("Failed to fetch products from Product Service: %v", err)
		return nil, fmt.Errorf("failed to fetch product details")
	}

	productMap := make(map[int]*productModels.Product)
	for i, product := range products {
		if product.ShopID != req.ShopID {
			return nil, fmt.Errorf("product %d does not belong to shop %d", product.ID, req.ShopID)
		}
		productMap[product.ID] = &products[i]
	}

	for _, item := range req.Items {
		product := productMap[item.ProductID]

		if item.ProductID <= 0 {
			return nil, fmt.Errorf("invalid product ID: %d", item.ProductID)
		}
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("quantity must be greater than 0 for product %d", item.ProductID)
		}
		// Check stock availability
		if item.Quantity > product.Stock {
			return nil, fmt.Errorf("insufficient stock for product %d: requested %d, available %d",
				item.ProductID, item.Quantity, product.Stock)
		}

		if item.Price != product.Price {
			return nil, fmt.Errorf("price mismatch for product %d: expected %.2f, got %.2f",
				item.ProductID, product.Price, item.Price)
		}

		itemPrices[item.ProductID] = product.Price
		totalPrice += product.Price * float64(item.Quantity)

	}

	// Create orders for all products
	log.Printf("[CreateOrder] Creating order records in database")
	tx, err := u.orderRepo.BeginTx(ctx)
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Failed to rollback transaction: %v", rollbackErr)
			}
		}
	}()

	orderGroup, err := u.orderRepo.CreateOrder(tx, req, itemPrices)
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

	err = tx.Commit()
	if err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
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
