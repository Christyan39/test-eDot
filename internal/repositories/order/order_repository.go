package order

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	orderModel "github.com/Christyan39/test-eDot/internal/models/order"
)

// OrderRepository defines the order repository interface
type OrderRepository interface {
	CreateOrder(tx *sql.Tx, req *orderModel.CreateOrderRequest, itemPrices map[int]float64) (*orderModel.OrderGroup, error)
	GetByID(ctx context.Context, id int) (*orderModel.Order, error)
	BeginTx(ctx context.Context) (*sql.Tx, error)
}

// orderRepository implements OrderRepository
type orderRepository struct {
	db *sql.DB
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

// CreateOrder creates multiple order records for a multi-product order
func (r *orderRepository) CreateOrder(tx *sql.Tx, req *orderModel.CreateOrderRequest, itemPrices map[int]float64) (*orderModel.OrderGroup, error) {
	orderDataJSON, err := json.Marshal(req.OrderData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal order data: %w", err)
	}

	// Begin transaction for atomic order creation
	query := `
		INSERT INTO orders (user_id, shop_id, product_id, quantity, total_price, status, order_data, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	var orders []orderModel.Order
	var totalPrice float64
	var totalItems int

	// Create an order record for each product
	for _, item := range req.Items {
		itemTotal := itemPrices[item.ProductID] * float64(item.Quantity)
		totalPrice += itemTotal
		totalItems += item.Quantity

		result, err := tx.Exec(query,
			req.UserID,
			req.ShopID,
			item.ProductID,
			item.Quantity,
			itemTotal,
			orderModel.OrderStatusPending,
			orderDataJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create order for product %d: %w", item.ProductID, err)
		}

		// Get the inserted ID
		orderID, err := result.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("failed to get order ID for product %d: %w", item.ProductID, err)
		}

		// Create order object
		order := orderModel.Order{
			ID:         int(orderID),
			UserID:     req.UserID,
			ShopID:     req.ShopID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			TotalPrice: itemTotal,
			Status:     orderModel.OrderStatusPending,
			OrderData:  req.OrderData,
		}
		orders = append(orders, order)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return order group
	return &orderModel.OrderGroup{
		Orders:     orders,
		TotalPrice: totalPrice,
		TotalItems: totalItems,
		UserID:     req.UserID,
		ShopID:     req.ShopID,
		Status:     orderModel.OrderStatusPending,
	}, nil
}

// GetByID retrieves an order by ID
func (r *orderRepository) GetByID(ctx context.Context, id int) (*orderModel.Order, error) {
	query := `
		SELECT id, user_id, shop_id, product_id, quantity, total_price, status, order_data, created_at, updated_at
		FROM orders
		WHERE id = ?
	`

	var order orderModel.Order
	var orderDataJSON []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID,
		&order.UserID,
		&order.ShopID,
		&order.ProductID,
		&order.Quantity,
		&order.TotalPrice,
		&order.Status,
		&orderDataJSON,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Unmarshal order data
	if err := json.Unmarshal(orderDataJSON, &order.OrderData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order data: %w", err)
	}

	return &order, nil
}

func (r *orderRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return tx, nil
}
