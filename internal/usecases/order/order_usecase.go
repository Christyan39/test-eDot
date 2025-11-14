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
	CreateOrder(ctx context.Context, req *orderModel.CreateOrderRequest) error
}

// orderUsecase implements OrderUsecase
type orderUsecase struct {
	orderRepo     orderRepo.OrderRepository
	productClient clients.ProductServiceClientInterface
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
func (u *orderUsecase) CreateOrder(ctx context.Context, req *orderModel.CreateOrderRequest) error {
	// Validate each item and collect product information via HTTP calls
	totalPrice := 0.0

	itemIDs := make([]int64, 0, len(req.Items))
	for _, item := range req.Items {
		itemIDs = append(itemIDs, item.ProductID)
	}

	// Fetch product details in batch
	products, err := u.productClient.GetProductByIDs(itemIDs)
	if err != nil {
		log.Printf("Failed to fetch products from Product Service: %v", err)
		return fmt.Errorf("failed to fetch product details")
	}

	productMap := make(map[int64]*productModels.Product)
	for i, product := range products {
		if product.ShopID != req.ShopID {
			return fmt.Errorf("product %d does not belong to shop %d", product.ID, req.ShopID)
		}
		productMap[product.ID] = &products[i]
	}

	for _, item := range req.Items {
		product := productMap[item.ProductID]

		if item.ProductID <= 0 {
			return fmt.Errorf("invalid product ID: %d", item.ProductID)
		}
		if item.Quantity <= 0 {
			return fmt.Errorf("quantity must be greater than 0 for product %d", item.ProductID)
		}
		// Check stock availability
		if item.Quantity > product.Stock {
			return fmt.Errorf("insufficient stock for product %d: requested %d, available %d",
				item.ProductID, item.Quantity, product.Stock)
		}

		if item.Price != product.Price {
			return fmt.Errorf("price mismatch for product %d: expected %.2f, got %.2f",
				item.ProductID, product.Price, item.Price)
		}

		totalPrice += product.Price * float64(item.Quantity)
	}

	if totalPrice != req.TotalPrice {
		return fmt.Errorf("total price mismatch: expected %.2f, got %.2f", totalPrice, req.TotalPrice)
	}

	// Create orders for all products
	log.Printf("[CreateOrder] Creating order records in database")
	tx, err := u.orderRepo.BeginTx(ctx)
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Failed to rollback transaction: %v", rollbackErr)
			}
		}
	}()

	orderID, err := u.orderRepo.CreateOrder(tx, req)
	if err != nil {
		log.Printf("Failed to create orders: %v", err)
		return fmt.Errorf("failed to create orders: %w", err)
	}
	productOnHoldStockUpdates := make([]productModels.UpdateProductRequest, 0, len(req.Items))

	for i, _ := range req.Items {
		req.Items[i].OrderID = orderID
		productOnHoldStockUpdates = append(productOnHoldStockUpdates, productModels.UpdateProductRequest{})
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
