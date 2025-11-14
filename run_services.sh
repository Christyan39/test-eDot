#!/bin/sh
# Run Go microservices in separate Terminal windows (macOS, no Docker)


# Start Product Service
echo "Starting product-service..."
osascript -e 'tell application "Terminal" to do script "cd '$(pwd)' && go run ./cmd/server/product/main.go"'

# Start Order Service
echo "Starting order-service..."
osascript -e 'tell application "Terminal" to do script "cd '$(pwd)' && go run ./cmd/server/order/main.go"'

# Start User Service
echo "Starting user-service..."
osascript -e 'tell application "Terminal" to do script "cd '$(pwd)' && go run ./cmd/server/user/main.go"'

echo "All Go services started in new Terminal windows."
