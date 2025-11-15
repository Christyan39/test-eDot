package order

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Christyan39/test-eDot/internal/clients"
	orderModel "github.com/Christyan39/test-eDot/internal/models/order"
	productModels "github.com/Christyan39/test-eDot/internal/models/product"
	orderRepo "github.com/Christyan39/test-eDot/internal/repositories/order"
	"github.com/Christyan39/test-eDot/pkg/config"
	"github.com/Christyan39/test-eDot/pkg/nsq"
	nsqio "github.com/nsqio/go-nsq"
)

// OrderUsecase defines the order business logic interface
type OrderUsecase interface {
	CreateOrder(ctx context.Context, req *orderModel.CreateOrderRequest) error
	ProcessOrderMessage(msg *nsqio.Message) error
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
	expDuration, _ := strconv.Atoi(config.GetEnv("ORDER_EXPIRATION_DURATION_SECONDS", "1"))
	expDurationMs := expDuration * 1000 // convert to milliseconds

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

	req.ExpiresAt = time.Now().Add(time.Duration(expDuration) * time.Millisecond)
	orderID, err := u.orderRepo.CreateOrder(tx, req)
	if err != nil {
		log.Printf("Failed to create orders: %v", err)
		return fmt.Errorf("failed to create orders: %w", err)
	}

	holdStockRequest := productModels.HoldStockRequest{
		OrderID:  orderID,
		Products: []productModels.Product{},
	}

	for i, item := range req.Items {
		req.Items[i].OrderID = orderID
		holdStockRequest.Products = append(holdStockRequest.Products, productModels.Product{
			ID:          item.ProductID,
			OnHoldStock: item.Quantity,
		})
	}

	err = u.orderRepo.CreateOrderItem(tx, req.Items)
	if err != nil {
		log.Printf("Failed to create order items: %v", err)
		return fmt.Errorf("failed to create order items: %w", err)
	}

	err = u.productClient.HoldStockInBulk(ctx, &holdStockRequest)
	if err != nil {
		log.Printf("Failed to hold stock in Product Service: %v", err)
		return fmt.Errorf("failed to hold stock in Product Service: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Publish order creation event to NSQ for further processing
	nsqAddress := config.GetEnv("NSQD_HOST", "http://localhost:4151")
	topic := config.GetEnv("NSQ_TOPIC_ORDER", "")

	req.OrderID = orderID
	nsqErr := nsq.PublishHTTP(nsqAddress, topic, req, expDurationMs)
	if nsqErr != nil {
		log.Printf("[NSQERROR] Failed to publish order %d to NSQ: %v", orderID, nsqErr)
	}

	return nil
}

// ProcessOrderMessage processes an order message from NSQ
func (u *orderUsecase) ProcessOrderMessage(msg *nsqio.Message) error {
	ctx := context.Background()

	// Example: Unmarshal and log the message. Replace with real logic.
	var req orderModel.CreateOrderRequest
	if err := json.Unmarshal(msg.Body, &req); err != nil {
		log.Printf("[NSQ] Failed to unmarshal order message: %v", err)
		return err
	}
	log.Printf("[NSQ] Processing order message: %+v", req)
	// TODO: Add your business logic her

	tx, err := u.orderRepo.BeginTx(ctx)
	if err != nil {
		log.Printf("[NSQ] Failed to begin transaction: %v", err)
		return err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("[NSQ] Failed to rollback transaction: %v", rollbackErr)
			}
		}
	}()

	// GetOrderByIDForUpdate
	checkedOrder, err := u.orderRepo.GetByIDForUpdateTx(ctx, tx, int(req.OrderID))
	if err != nil {
		log.Printf("[NSQ] Failed to get order %d for update: %v", req.OrderID, err)
		return err
	}

	if checkedOrder.Status != "pending" {
		log.Printf("[NSQ] Order %d already processed with status %s", req.OrderID, checkedOrder.Status)
		tx.Commit()
		return nil
	}

	if time.Now().Before(checkedOrder.ExpiresAt) {
		log.Printf("[NSQ] Order %d not expired yet, skipping", req.OrderID)
		tx.Commit()

		// requeue the message for later processing
		nsqAddress := config.GetEnv("NSQD_HOST", "http://localhost:4151")
		topic := config.GetEnv("NSQ_TOPIC_ORDER", "")
		duration := checkedOrder.ExpiresAt.Sub(time.Now()).Milliseconds()

		nsqErr := nsq.PublishHTTP(nsqAddress, topic, req, int(duration))
		if nsqErr != nil {
			log.Printf("[NSQERROR] Failed to publish order %d to NSQ: %v", req.OrderID, nsqErr)
		}

		return nil
	}

	checkedOrder.Status = "expired"
	err = u.orderRepo.UpdateOrderStatusTx(ctx, tx, checkedOrder.ID, checkedOrder.Status)
	if err != nil {
		log.Printf("[NSQ] Failed to update order %d status to EXPIRED: %v", req.OrderID, err)
		return err
	}

	err = u.productClient.ReleaseHeldStockInBulk(ctx, &productModels.ReleaseHeldStockRequest{
		OrderID: req.OrderID,
	})
	if err != nil {
		log.Printf("[NSQ] Failed to release held stock for order %d: %v", req.OrderID, err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("[NSQ] Failed to commit transaction for order %d: %v", req.OrderID, err)
		return err
	}

	log.Printf("[NSQ] Order %d marked as EXPIRED due to non-payment", req.OrderID)
	return nil
}
