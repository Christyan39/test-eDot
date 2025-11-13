#!/bin/bash

echo "🧪 Quick Logging Test"
echo "Starting service for 5 seconds to test logs..."

# Build first to ensure we have the binary
go build -o bin/product ./cmd/server/product

# Start service in background  
./bin/product &
PID=$!

# Wait for startup
sleep 2

echo -e "\n📞 Making test API call..."

# Make a simple API call to trigger handler logs
curl -s "http://localhost:8081/api/v1/products" > /dev/null

sleep 1

# Stop the service
kill $PID 2>/dev/null
wait $PID 2>/dev/null

echo -e "\n✅ Test complete! You should see:"
echo "   ✓ [STARTUP] logs during initialization" 
echo "   ✓ Echo middleware logs for the API call"
echo "   ✓ [ListProducts] handler logs"