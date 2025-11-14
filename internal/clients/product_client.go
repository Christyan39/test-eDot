package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	productModels "github.com/Christyan39/test-eDot/internal/models/product"
)

type ProductServiceClient struct {
	BaseURL    string
	HTTPClient *http.Client
	APIKey     string
}

// NewProductServiceClient creates a new product service HTTP client
func NewProductServiceClient(baseURL, apiKey string) ProductServiceClientInterface {
	return &ProductServiceClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		APIKey: apiKey,
	}
}

type ProductServiceClientInterface interface {
	GetProductByIDs(productIDs []int64) ([]productModels.Product, error)
	UpdateProductStock(productID int64, req *productModels.UpdateProductRequest) error
}

// GetProductByID makes HTTP call to product service to get product details
func (p *ProductServiceClient) GetProductByIDs(productIDs []int64) ([]productModels.Product, error) {
	productIDsStr := make([]string, len(productIDs))
	for i, id := range productIDs {
		productIDsStr[i] = fmt.Sprintf("%d", id)
	}

	url := fmt.Sprintf("%s/products?ids=%s", p.BaseURL, strings.Join(productIDsStr, ","))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", p.APIKey)

	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("product service returned status %d: %s", resp.StatusCode, string(body))
	}

	var product productModels.ProductListResponse
	if err := json.Unmarshal(body, &product); err != nil {
		return nil, fmt.Errorf("failed to unmarshal product: %w", err)
	}

	return product.Products, nil
}

// UpdateProductStock makes HTTP call to product service to update stock
func (p *ProductServiceClient) UpdateProductStock(productID int64, req *productModels.UpdateProductRequest) error {
	url := fmt.Sprintf("%s/products/%d/hold-stock", p.BaseURL, productID)

	reqBody := map[string]int{
		"on_hold_stock": req.OnHoldStock,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-Key", p.APIKey)

	resp, err := p.HTTPClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("product service returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
