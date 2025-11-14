package product

import (
	"context"
	"fmt"
	"log"

	productModel "github.com/Christyan39/test-eDot/internal/models/product"
	productRepo "github.com/Christyan39/test-eDot/internal/repositories/product"
)

// ProductUsecase defines the product business logic interface
type ProductUsecase interface {
	CreateProduct(ctx context.Context, req *productModel.CreateProductRequest) error
	ListProducts(ctx context.Context, req *productModel.ProductListRequest) (*productModel.ProductListResponse, error)
	UpdateOnHoldStock(ctx context.Context, id, newOnHoldStock int) error
}

// productUsecase implements ProductUsecase
type productUsecase struct {
	productRepo productRepo.ProductRepository
}

// NewProductUsecase creates a new product usecase
func NewProductUsecase(productRepo productRepo.ProductRepository) ProductUsecase {
	return &productUsecase{
		productRepo: productRepo,
	}
}

// CreateProduct creates a new product
func (u *productUsecase) CreateProduct(ctx context.Context, req *productModel.CreateProductRequest) error {
	// Validate shop metadata
	if req.ShopMetadata.ShopID <= 0 {
		return fmt.Errorf("invalid shop ID in metadata")
	}
	if req.ShopMetadata.ShopName == "" {
		return fmt.Errorf("shop name is required in metadata")
	}

	// Create the product
	err := u.productRepo.Create(req)
	if err != nil {
		log.Printf("Failed to create product: %v", err)
		return fmt.Errorf("failed to create product: %w", err)
	}

	log.Printf("Product created successfully for shop: %s", req.ShopMetadata.ShopName)
	return nil
}

// ListProducts retrieves products with filtering and pagination
func (u *productUsecase) ListProducts(ctx context.Context, req *productModel.ProductListRequest) (*productModel.ProductListResponse, error) {
	// Set default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// Validate price range
	if req.MinPrice > 0 && req.MaxPrice > 0 && req.MinPrice > req.MaxPrice {
		return nil, fmt.Errorf("minimum price cannot be greater than maximum price")
	}

	response, err := u.productRepo.List(req)
	if err != nil {
		log.Printf("Failed to list products: %v", err)
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	log.Printf("Listed %d products (page %d, limit %d)", len(response.Products), req.Page, req.Limit)
	return response, nil
}

// UpdateOnHoldStock updates product on-hold stock with transaction safety
func (u *productUsecase) UpdateOnHoldStock(ctx context.Context, id, newOnHoldStock int) error {
	if id <= 0 {
		return fmt.Errorf("invalid product ID")
	}
	if newOnHoldStock < 0 {
		return fmt.Errorf("on-hold stock cannot be negative")
	}

	// Begin transaction for atomic operation
	tx, err := u.productRepo.TxBegin(ctx)
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

	// Get product with row lock within transaction
	product, err := u.productRepo.GetByIDForUpdateTx(tx, id)
	if err != nil {
		log.Printf("Failed to get product with ID %d: %v", id, err)
		return fmt.Errorf("product not found")
	}

	// Validate that on-hold stock doesn't exceed available stock
	if newOnHoldStock > product.Stock {
		return fmt.Errorf("on-hold stock (%d) cannot exceed available stock (%d)", newOnHoldStock, product.Stock)
	}

	// Update on-hold stock within transaction
	err = u.productRepo.UpdateTx(tx, id, &productModel.UpdateProductRequest{
		OnHoldStock: product.OnHoldStock + newOnHoldStock,
		Stock:       product.Stock - newOnHoldStock,
	})
	if err != nil {
		log.Printf("Failed to update on-hold stock for product ID %d: %v", id, err)
		return fmt.Errorf("failed to update on-hold stock: %w", err)
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Successfully updated on-hold stock for product ID %d: %d -> %d", id, product.OnHoldStock, product.OnHoldStock+newOnHoldStock)
	return nil
}
