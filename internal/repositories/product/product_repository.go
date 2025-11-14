package product

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	productModel "github.com/Christyan39/test-eDot/internal/models/product"
)

// ProductRepository defines the product repository interface
type ProductRepository interface {
	Create(product *productModel.CreateProductRequest) error
	GetByIDForUpdateTx(tx *sql.Tx, id int) (*productModel.Product, error)
	List(req *productModel.ProductListRequest) (*productModel.ProductListResponse, error)
	UpdateTx(tx *sql.Tx, id int, req *productModel.UpdateProductRequest) error
	TxBegin(ctx context.Context) (*sql.Tx, error)
}

// productRepository implements ProductRepository
type productRepository struct {
	db *sql.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

// Create creates a new product
func (r *productRepository) Create(req *productModel.CreateProductRequest) error {
	shopMetadataJSON, err := json.Marshal(req.ShopMetadata)
	if err != nil {
		return fmt.Errorf("failed to marshal shop metadata: %w", err)
	}

	query := `
		INSERT INTO products (name, description, price, stock, on_hold_stock, shop_id, shop_metadata, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, 'active', NOW(), NOW())
	`

	_, err = r.db.Exec(query,
		req.Name,
		req.Description,
		req.Price,
		req.Stock,
		req.OnHoldStock,
		req.ShopID,
		shopMetadataJSON,
	)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

// GetByIDForUpdateTx retrieves a product by ID within a transaction with row lock
func (r *productRepository) GetByIDForUpdateTx(tx *sql.Tx, id int) (*productModel.Product, error) {
	query := `
		SELECT id, name, description, price, stock, on_hold_stock, shop_id, shop_metadata, status, created_at, updated_at
		FROM products
		WHERE id = ? FOR UPDATE
	`

	var product productModel.Product
	var shopMetadataJSON []byte
	err := tx.QueryRow(query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.OnHoldStock,
		&product.ShopID,
		&shopMetadataJSON,
		&product.Status,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Unmarshal shop metadata
	if err := json.Unmarshal(shopMetadataJSON, &product.ShopMetadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal shop metadata: %w", err)
	}

	return &product, nil
}

// List retrieves products with filtering and pagination
func (r *productRepository) List(req *productModel.ProductListRequest) (*productModel.ProductListResponse, error) {
	countQuery := "SELECT COUNT(*) FROM products WHERE 1=1"
	query := `
		SELECT id, name, description, price, stock, on_hold_stock, shop_id, shop_metadata, status, created_at, updated_at
		FROM products
		WHERE 1=1
	`

	args := []interface{}{}
	conditions := []string{}

	// Add filters
	if req.ShopID > 0 {
		conditions = append(conditions, "shop_id = ?")
		args = append(args, req.ShopID)
	}
	if req.MinPrice > 0 {
		conditions = append(conditions, "price >= ?")
		args = append(args, req.MinPrice)
	}
	if req.MaxPrice > 0 {
		conditions = append(conditions, "price <= ?")
		args = append(args, req.MaxPrice)
	}
	if req.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, req.Status)
	}
	if req.Search != "" {
		conditions = append(conditions, "(name LIKE ? OR description LIKE ?)")
		searchParam := "%" + req.Search + "%"
		args = append(args, searchParam, searchParam)
	}

	// Apply conditions
	if len(conditions) > 0 {
		conditionStr := " AND " + strings.Join(conditions, " AND ")
		countQuery += conditionStr
		query += conditionStr
	}

	// Get total count
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// Add pagination
	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	offset := (req.Page - 1) * req.Limit
	args = append(args, req.Limit, offset)

	// Execute query
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	defer rows.Close()

	var products []productModel.Product
	for rows.Next() {
		var product productModel.Product
		var shopMetadataJSON []byte
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.OnHoldStock,
			&product.ShopID,
			&shopMetadataJSON,
			&product.Status,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}

		// Unmarshal shop metadata
		if err := json.Unmarshal(shopMetadataJSON, &product.ShopMetadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal shop metadata: %w", err)
		}

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	// Calculate pages
	pages := (total + req.Limit - 1) / req.Limit

	return &productModel.ProductListResponse{
		Products: products,
		Total:    total,
		Page:     req.Page,
		Limit:    req.Limit,
		Pages:    pages,
	}, nil
}

// UpdateTx updates a product by ID with partial data within a transaction
func (r *productRepository) UpdateTx(tx *sql.Tx, id int, req *productModel.UpdateProductRequest) error {
	// Build dynamic update query based on provided fields
	setClauses := []string{}
	args := []interface{}{}

	if req.Name != "" {
		setClauses = append(setClauses, "name = ?")
		args = append(args, req.Name)
	}

	if req.Description != "" {
		setClauses = append(setClauses, "description = ?")
		args = append(args, req.Description)
	}

	if req.Price > 0 {
		setClauses = append(setClauses, "price = ?")
		args = append(args, req.Price)
	}

	setClauses = append(setClauses, "stock = ?")
	args = append(args, req.Stock)

	setClauses = append(setClauses, "on_hold_stock = ?")
	args = append(args, req.OnHoldStock)

	if req.ShopMetadata != nil {
		shopMetadataJSON, err := json.Marshal(req.ShopMetadata)
		if err != nil {
			return fmt.Errorf("failed to marshal shop metadata: %w", err)
		}
		setClauses = append(setClauses, "shop_metadata = ?")
		args = append(args, string(shopMetadataJSON))
	}

	if req.Status != "" {
		setClauses = append(setClauses, "status = ?")
		args = append(args, req.Status)
	}

	// Always update updated_at timestamp
	setClauses = append(setClauses, "updated_at = NOW()")

	if len(setClauses) == 1 { // Only updated_at was added
		return fmt.Errorf("no fields provided for update")
	}

	// Add product ID as the last parameter for WHERE clause
	args = append(args, id)

	// Build the complete query
	query := fmt.Sprintf(`
		UPDATE products 
		SET %s 
		WHERE id = ?
	`, strings.Join(setClauses, ", "))

	// Execute the update within transaction
	result, err := tx.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

func (r *productRepository) TxBegin(ctx context.Context) (*sql.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return tx, nil
}
