package order

import (
	"time"
)

// Order represents an order in the system
type Order struct {
	ID         int                    `json:"id" db:"id"`
	UserID     int                    `json:"user_id" db:"user_id"`
	ShopID     int                    `json:"shop_id" db:"shop_id"`
	ProductID  int                    `json:"product_id" db:"product_id"`
	Quantity   int                    `json:"quantity" db:"quantity"`
	TotalPrice float64                `json:"total_price" db:"total_price"`
	Status     string                 `json:"status" db:"status"`
	OrderData  map[string]interface{} `json:"order_data" db:"order_data"`
	CreatedAt  time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at" db:"updated_at"`
}

// OrderGroup represents a group of related orders from the same request
type OrderGroup struct {
	Orders     []Order   `json:"orders"`
	TotalPrice float64   `json:"total_price"`
	TotalItems int       `json:"total_items"`
	UserID     int       `json:"user_id"`
	ShopID     int       `json:"shop_id"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

// OrderItem represents a single item in an order
type OrderItem struct {
	ProductID int     `json:"product_id" validate:"required,min=1"`
	Quantity  int     `json:"quantity" validate:"required,min=1"`
	Price     float64 `json:"price" validate:"required,min=0"`
}

// CreateOrderRequest represents the request to create a new order with multiple products
type CreateOrderRequest struct {
	UserID     int                    `json:"user_id" validate:"required,min=1"`
	ShopID     int                    `json:"shop_id" validate:"required,min=1"`
	TotalPrice float64                `json:"total_price"`
	Items      []OrderItem            `json:"items" validate:"required,min=1,dive"`
	OrderData  map[string]interface{} `json:"order_data,omitempty"`
}

// CreateOrderResponse represents the response after creating an order
type CreateOrderResponse struct {
	OrderIDs   []int   `json:"order_ids"`
	TotalPrice float64 `json:"total_price"`
	TotalItems int     `json:"total_items"`
	Status     string  `json:"status"`
	Message    string  `json:"message"`
}

// OrderStatus constants
const (
	OrderStatusPending   = "pending"
	OrderStatusConfirmed = "confirmed"
	OrderStatusShipped   = "shipped"
	OrderStatusDelivered = "delivered"
	OrderStatusCancelled = "cancelled"
)
