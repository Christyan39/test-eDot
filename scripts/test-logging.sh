#!/bin/bash

echo "🚀 Testing Product Service Logging..."
echo "Starting the service in background for 10 seconds to capture logs..."

# Start the service in background
./bin/product &
SERVICE_PID=$!

# Wait for service to start
sleep 2

echo -e "\n📝 Making API calls to generate logs...\n"

# Test 1: Health check (should show in Echo middleware logs)
echo "1. Testing health check:"
curl -s "http://localhost:8081/health" > /dev/null

sleep 1

# Test 2: List products (should show handler logs)
echo "2. Testing list products:"
curl -s "http://localhost:8081/api/v1/products" > /dev/null

sleep 1

# Test 3: Create product (should show detailed handler logs)
echo "3. Testing create product:"
curl -s -X POST "http://localhost:8081/api/v1/products" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Product",
    "description": "This is a test product for logging verification",
    "price": 99.99,
    "stock": 10,
    "shop_id": 1,
    "shop_metadata": {
      "shop_name": "Test Shop",
      "shop_id": 1,
      "status": "active"
    }
  }' > /dev/null

sleep 1

# Test 4: Update hold stock (should show middleware auth logs)
echo "4. Testing hold stock update:"
curl -s -X PATCH "http://localhost:8081/api/v1/products/1/hold-stock" \
  -H "Content-Type: application/json" \
  -H "X-API-Key: internal-api-key-change-in-production" \
  -d '{"on_hold_stock": 5}' > /dev/null

sleep 2

echo -e "\n🔚 Stopping service..."
kill $SERVICE_PID
wait $SERVICE_PID 2>/dev/null

echo -e "\n✅ Test completed! Check the above output for detailed logs from:"
echo "   - [STARTUP] messages during service initialization"
echo "   - [CreateProduct], [ListProducts], [UpdateOnHoldStock] messages from handlers"
echo "   - Echo middleware request/response logs"
echo ""
echo "💡 To run the service manually and see live logs:"
echo "   ./bin/product"