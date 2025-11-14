package product

import "time"

// ShopMetadata represents shop information stored as JSON
type ShopMetadata struct {
	ShopName string `json:"shop_name"`
	ShopID   int64  `json:"shop_id"`
	Status   string `json:"status"`
}

// Product represents a product entity
type Product struct {
	ID           int64        `json:"id" db:"id"`
	Name         string       `json:"name" db:"name"`
	Description  string       `json:"description" db:"description"`
	Price        float64      `json:"price" db:"price"`
	Stock        int          `json:"stock" db:"stock"`
	OnHoldStock  int          `json:"on_hold_stock" db:"on_hold_stock"`
	ShopID       int          `json:"shop_id" db:"shop_id"`
	ShopMetadata ShopMetadata `json:"shop_metadata" db:"shop_metadata"`
	Status       string       `json:"status" db:"status"` // active, inactive, discontinued
	CreatedAt    time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at" db:"updated_at"`
}

// CreateProductRequest represents request to create product
type CreateProductRequest struct {
	Name         string       `json:"name" validate:"required,min=2,max=100"`
	Description  string       `json:"description" validate:"required,min=10,max=1000"`
	Price        float64      `json:"price" validate:"required,min=0"`
	Stock        int          `json:"stock" validate:"required,min=0"`
	OnHoldStock  int          `json:"on_hold_stock" validate:"min=0"`
	ShopID       int          `json:"shop_id" validate:"required,min=1"`
	ShopMetadata ShopMetadata `json:"shop_metadata" validate:"required"`
}

// UpdateProductRequest represents request to update product
type UpdateProductRequest struct {
	ID           int64         `json:"id" validate:"required,min=1"`
	Name         string        `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description  string        `json:"description,omitempty" validate:"omitempty,min=10,max=1000"`
	Price        float64       `json:"price,omitempty" validate:"omitempty,min=0"`
	Stock        int           `json:"stock,omitempty" validate:"omitempty,min=0"`
	OnHoldStock  int           `json:"on_hold_stock,omitempty" validate:"omitempty,min=0"`
	ShopID       int           `json:"shop_id,omitempty" validate:"omitempty,min=1"`
	ShopMetadata *ShopMetadata `json:"shop_metadata,omitempty"`
	Status       string        `json:"status,omitempty" validate:"omitempty,oneof=active inactive discontinued"`
}

// ProductListRequest represents request for product listing with filters
type ProductListRequest struct {
	Page     int     `json:"page" query:"page" validate:"min=1"`
	Limit    int     `json:"limit" query:"limit" validate:"min=1,max=100"`
	ShopID   int     `json:"shop_id" query:"shop_id" validate:"omitempty,min=1"`
	MinPrice float64 `json:"min_price" query:"min_price" validate:"omitempty,min=0"`
	MaxPrice float64 `json:"max_price" query:"max_price" validate:"omitempty,min=0"`
	Status   string  `json:"status" query:"status" validate:"omitempty,oneof=active inactive discontinued"`
	Search   string  `json:"search" query:"search" validate:"omitempty,max=100"`
	IDs      []int   `json:"ids" query:"ids" validate:"omitempty,dive,min=1"`
}

// ProductListResponse represents paginated product list response
type ProductListResponse struct {
	Products []Product `json:"products"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	Limit    int       `json:"limit"`
	Pages    int       `json:"pages"`
}

type HoldStockRequest struct {
	OrderID  int64     `json:"order_id" validate:"required,min=1"`
	Products []Product `json:"products" validate:"required,dive"`
}

type HoldStockAudit struct {
	ID        int64     `json:"id" db:"id"`
	OrderID   int64     `json:"order_id" db:"order_id"`
	ProductID int64     `json:"product_id" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	Status    string    `json:"status" db:"status"` // held, success, cancelled
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
