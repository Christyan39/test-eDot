#!/bin/bash

echo "🧪 Testing Handler Logs - this will show visible logs!"
echo "Starting service..."

# Start the service in the background
./bin/product &
SERVICE_PID=$!

# Wait for service to start up
sleep 3

echo -e "\n🔍 Making API calls to trigger handler logs...\n"

# Test 1: ListProducts - should show 📋 LIST PRODUCTS REQUEST
echo "1. Calling GET /products (should show LIST PRODUCTS logs with 📋 emoji)"
curl -s "http://localhost:8081/api/v1/products" > /dev/null
echo "   ↳ Check above for: 📋 LIST PRODUCTS REQUEST"

sleep 1

# Test 2: CreateProduct - should show 🚀 CREATE PRODUCT REQUEST  
echo -e "\n2. Calling POST /products (should show CREATE PRODUCT logs with 🚀 emoji)"
curl -s -X POST "http://localhost:8081/api/v1/products" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Visible Log Test Product",
    "description": "This product is created to test visible logging in handlers",
    "price": 123.45,
    "stock": 10,
    "shop_id": 999,
    "shop_metadata": {
      "shop_name": "Test Logging Shop",
      "shop_id": 999,
      "status": "active"
    }
  }' > /dev/null
echo "   ↳ Check above for: 🚀 CREATE PRODUCT REQUEST RECEIVED"

sleep 1

# Test 3: UpdateOnHoldStock - should show 🔄 UPDATE HOLD STOCK REQUEST
echo -e "\n3. Calling PATCH /products/1/hold-stock (should show UPDATE HOLD STOCK logs with 🔄 emoji)"  
curl -s -X PATCH "http://localhost:8081/api/v1/products/1/hold-stock" \
  -H "Content-Type: application/json" \
  -H "X-API-Key: internal-api-key-change-in-production" \
  -d '{"on_hold_stock": 7}' > /dev/null
echo "   ↳ Check above for: 🔄 UPDATE HOLD STOCK REQUEST"

sleep 2

echo -e "\n🛑 Stopping service..."
kill $SERVICE_PID 2>/dev/null
wait $SERVICE_PID 2>/dev/null

echo -e "\n✅ Test complete! You should have seen:"
echo "   📋 List Products logs with timestamps and IP"
echo "   🚀 Create Product logs with colored boxes"  
echo "   🔄 Update Hold Stock logs with product ID"
echo ""
echo "💡 If you didn't see the emoji logs, the handlers might not be getting called."
echo "   Try running './bin/product' manually and make API calls in another terminal."