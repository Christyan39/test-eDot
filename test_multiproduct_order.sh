#!/bin/bash

# Test script for Multi-Product Order API endpoint
echo "🛒 Testing Multi-Product Order API Endpoint"
echo "============================================="

# Configuration
BASE_URL="http://localhost:8082"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}Base URL: $BASE_URL${NC}"
echo -e "${BLUE}Testing Multi-Product Order functionality...${NC}"
echo ""

# Test 1: Create order with multiple products
echo -e "${YELLOW}Test 1: Create order with multiple products${NC}"
echo "POST $BASE_URL/orders"
curl -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{
    "user_id": 1,
    "shop_id": 1,
    "items": [
      {
        "product_id": 1,
        "quantity": 2
      },
      {
        "product_id": 2,
        "quantity": 1
      }
    ],
    "order_data": {
      "shipping_address": "123 Main St, City, Country",
      "payment_method": "credit_card",
      "notes": "Multi-product order test"
    }
  }' \
  -w "\n\nResponse time: %{time_total}s\nHTTP Status: %{http_code}\n" -s
echo -e "\n${GREEN}✓ Multi-product order test completed${NC}\n"

# Test 2: Create order with single product (backward compatibility)
echo -e "${YELLOW}Test 2: Create order with single product${NC}"
echo "POST $BASE_URL/orders"
curl -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{
    "user_id": 1,
    "shop_id": 1,
    "items": [
      {
        "product_id": 1,
        "quantity": 3
      }
    ],
    "order_data": {
      "shipping_address": "456 Oak Ave, City, Country",
      "payment_method": "debit_card"
    }
  }' \
  -w "\n\nResponse time: %{time_total}s\nHTTP Status: %{http_code}\n" -s
echo -e "\n${GREEN}✓ Single product order test completed${NC}\n"

# Test 3: Create order with large quantity (multiple of same product)
echo -e "${YELLOW}Test 3: Create order with large quantities${NC}"
echo "POST $BASE_URL/orders"
curl -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{
    "user_id": 2,
    "shop_id": 1,
    "items": [
      {
        "product_id": 1,
        "quantity": 5
      },
      {
        "product_id": 2,
        "quantity": 3
      },
      {
        "product_id": 3,
        "quantity": 2
      }
    ],
    "order_data": {
      "shipping_address": "789 Pine St, City, Country",
      "payment_method": "paypal",
      "priority": "express"
    }
  }' \
  -w "\n\nResponse time: %{time_total}s\nHTTP Status: %{http_code}\n" -s
echo -e "\n${GREEN}✓ Large quantity order test completed${NC}\n"

# Test 4: Invalid request - empty items array
echo -e "${YELLOW}Test 4: Invalid request - empty items array${NC}"
echo "POST $BASE_URL/orders"
curl -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{
    "user_id": 1,
    "shop_id": 1,
    "items": [],
    "order_data": {
      "shipping_address": "Test Address"
    }
  }' \
  -w "\n\nResponse time: %{time_total}s\nHTTP Status: %{http_code}\n" -s
echo -e "\n${RED}✓ Empty items test completed (should be 400)${NC}\n"

# Test 5: Invalid request - negative quantity
echo -e "${YELLOW}Test 5: Invalid request - negative quantity${NC}"
echo "POST $BASE_URL/orders"
curl -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{
    "user_id": 1,
    "shop_id": 1,
    "items": [
      {
        "product_id": 1,
        "quantity": -1
      }
    ]
  }' \
  -w "\n\nResponse time: %{time_total}s\nHTTP Status: %{http_code}\n" -s
echo -e "\n${RED}✓ Negative quantity test completed (should be 400)${NC}\n"

# Test 6: Non-existent product
echo -e "${YELLOW}Test 6: Non-existent product${NC}"
echo "POST $BASE_URL/orders"
curl -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{
    "user_id": 1,
    "shop_id": 1,
    "items": [
      {
        "product_id": 999,
        "quantity": 1
      }
    ]
  }' \
  -w "\n\nResponse time: %{time_total}s\nHTTP Status: %{http_code}\n" -s
echo -e "\n${RED}✓ Non-existent product test completed (should be 404)${NC}\n"

echo -e "${BLUE}=============================================="
echo -e "🏁 Multi-Product Order API Tests Completed!"
echo -e "=============================================${NC}"
echo ""
echo -e "${GREEN}Expected Results:${NC}"
echo -e "- Test 1, 2 & 3: ${GREEN}201 Created${NC} (successful multi-product orders)"
echo -e "- Test 4 & 5: ${RED}400 Bad Request${NC} (validation errors)"
echo -e "- Test 6: ${RED}404 Not Found${NC} (non-existent product)"
echo ""
echo -e "${BLUE}💡 New Multi-Product Features:${NC}"
echo -e "✅ Multiple products in single order"
echo -e "✅ Individual product validation"
echo -e "✅ Combined price calculation"
echo -e "✅ Multiple order IDs returned"
echo -e "✅ Total items counting"
echo ""
echo -e "${BLUE}📊 Response Format:${NC}"
echo -e '{"order_ids": [1,2,3], "total_price": 299.99, "total_items": 6, "status": "pending", "message": "Multi-product order created successfully"}'