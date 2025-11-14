#!/bin/bash

# Test script for Order Service API endpoints
echo "🧪 Testing Order Service API Endpoints"
echo "====================================="

# Configuration
BASE_URL="http://localhost:8082"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}Base URL: $BASE_URL${NC}"
echo -e "${BLUE}Testing Order Service functionality...${NC}"
echo ""

# Test 1: Health Check
echo -e "${YELLOW}Test 1: Health Check${NC}"
echo "GET $BASE_URL/health"
curl -X GET "$BASE_URL/health" \
  -w "\n\nResponse time: %{time_total}s\nHTTP Status: %{http_code}\n" -s
echo -e "\n${GREEN}✓ Health check test completed${NC}\n"

# Test 2: Create Order (valid request)
echo -e "${YELLOW}Test 2: Create Order (valid request)${NC}"
echo "POST $BASE_URL/orders"
curl -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{
    "user_id": 1,
    "shop_id": 1,
    "product_id": 1,
    "quantity": 2,
    "order_data": {
      "shipping_address": "123 Main St, City, Country",
      "payment_method": "credit_card",
      "notes": "Please handle with care"
    }
  }' \
  -w "\n\nResponse time: %{time_total}s\nHTTP Status: %{http_code}\n" -s
echo -e "\n${GREEN}✓ Valid order creation test completed${NC}\n"

# Test 3: Create Order with invalid data (missing required fields)
echo -e "${YELLOW}Test 3: Create Order with invalid data${NC}"
echo "POST $BASE_URL/orders"
curl -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{
    "user_id": 1,
    "quantity": 2
  }' \
  -w "\n\nResponse time: %{time_total}s\nHTTP Status: %{http_code}\n" -s
echo -e "\n${RED}✓ Invalid data test completed (should be 400)${NC}\n"

# Test 4: Create Order with zero quantity
echo -e "${YELLOW}Test 4: Create Order with zero quantity${NC}"
echo "POST $BASE_URL/orders"
curl -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{
    "user_id": 1,
    "shop_id": 1,
    "product_id": 1,
    "quantity": 0
  }' \
  -w "\n\nResponse time: %{time_total}s\nHTTP Status: %{http_code}\n" -s
echo -e "\n${RED}✓ Zero quantity test completed (should be 400)${NC}\n"

# Test 5: Create Order for non-existent product
echo -e "${YELLOW}Test 5: Create Order for non-existent product${NC}"
echo "POST $BASE_URL/orders"
curl -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{
    "user_id": 1,
    "shop_id": 1,
    "product_id": 999,
    "quantity": 1
  }' \
  -w "\n\nResponse time: %{time_total}s\nHTTP Status: %{http_code}\n" -s
echo -e "\n${RED}✓ Non-existent product test completed (should be 404)${NC}\n"

# Test 6: Create Order with excessive quantity (stock check)
echo -e "${YELLOW}Test 6: Create Order with excessive quantity${NC}"
echo "POST $BASE_URL/orders"
curl -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{
    "user_id": 1,
    "shop_id": 1,
    "product_id": 1,
    "quantity": 1000
  }' \
  -w "\n\nResponse time: %{time_total}s\nHTTP Status: %{http_code}\n" -s
echo -e "\n${RED}✓ Excessive quantity test completed (should be 409)${NC}\n"

echo -e "${BLUE}====================================="
echo -e "🏁 Order Service API Tests Completed!"
echo -e "=====================================${NC}"
echo ""
echo -e "${GREEN}Expected Results:${NC}"
echo -e "- Test 1: ${GREEN}200 OK${NC} (health check)"
echo -e "- Test 2: ${GREEN}201 Created${NC} (successful order)"
echo -e "- Test 3, 4: ${RED}400 Bad Request${NC} (validation errors)"
echo -e "- Test 5: ${RED}404 Not Found${NC} (product not found)"
echo -e "- Test 6: ${RED}409 Conflict${NC} (insufficient stock)"
echo ""
echo -e "${BLUE}💡 Check terminal logs for detailed handler output${NC}"