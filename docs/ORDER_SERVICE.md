# Order Service Documentation

## Overview
The Order Service is a microservice responsible for managing customer orders, order processing, and inventory management integration with the Product Service.

## Architecture
- **Port**: 8082
- **Database**: MySQL
- **Framework**: Echo v4
- **Documentation**: Swagger/OpenAPI

## Features
- ✅ **Create Order**: Place new orders with stock validation
- ✅ **Stock Integration**: Automatic product stock reservation
- ✅ **Price Calculation**: Dynamic pricing based on product data
- ✅ **Validation**: Comprehensive input and business rule validation
- ✅ **Error Handling**: Detailed error responses and logging

## API Endpoints

### Health Check
```http
GET /health
```
Returns service status and health information.

### Create Order
```http
POST /orders
Content-Type: application/json
Authorization: Bearer {token}
```

#### Request Body
```json
{
  "user_id": 1,
  "shop_id": 1,
  "product_id": 1,
  "quantity": 2,
  "order_data": {
    "shipping_address": "123 Main St, City, Country",
    "payment_method": "credit_card",
    "notes": "Optional order notes"
  }
}
```

#### Success Response (201 Created)
```json
{
  "id": 1,
  "total_price": 299.98,
  "status": "pending",
  "message": "Order created successfully"
}
```

#### Error Responses
- **400 Bad Request**: Invalid input or validation failed
- **404 Not Found**: Product not found
- **409 Conflict**: Insufficient stock
- **500 Internal Server Error**: Server error

## Business Logic

### Order Creation Flow
1. **Validation**: Validate user, shop, product IDs and quantity
2. **Product Check**: Verify product exists and belongs to specified shop
3. **Stock Check**: Ensure sufficient available stock (stock - on_hold_stock)
4. **Price Calculation**: Calculate total price (product_price × quantity)
5. **Order Creation**: Insert order record with "pending" status
6. **Stock Reservation**: Update product on_hold_stock to reserve inventory
7. **Response**: Return order details with total price

### Stock Management
- **Available Stock**: `stock - on_hold_stock`
- **Stock Reservation**: Increases `on_hold_stock` when order is created
- **Validation**: Prevents ordering more than available stock

## Database Schema

### Orders Table
```sql
CREATE TABLE orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    shop_id INT NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    total_price DECIMAL(10,2) NOT NULL,
    status ENUM('pending', 'confirmed', 'shipped', 'delivered', 'cancelled') DEFAULT 'pending',
    order_data JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

## Configuration

### Environment Variables
```bash
# Server Configuration
PORT=8082
SERVICE_VERSION=1.0.0

# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=ecommerce

# Logging Configuration
LOG_LEVEL=info
```

## Running the Service

### Development
```bash
# Build and run
make build-order
./bin/order

# Or run directly
make run-order
```

### Testing
```bash
# Run comprehensive API tests
./test_order_service.sh

# Manual testing with curl
curl -X POST http://localhost:8082/orders \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"shop_id":1,"product_id":1,"quantity":2}'
```

### Database Migration
```bash
# Execute the migration script
mysql -u root -p ecommerce < migration/order/001_create_orders_table.sql
```

## Integration with Other Services

### Product Service
- **Dependency**: Retrieves product details and manages stock
- **Stock Updates**: Modifies product on_hold_stock for inventory management
- **Validation**: Ensures product belongs to specified shop

### User Service (Future)
- **Validation**: Verify user existence and permissions
- **Authentication**: JWT token validation

## Logging and Monitoring

### Log Levels
- **INFO**: Order creation, successful operations
- **ERROR**: Validation failures, database errors
- **DEBUG**: Detailed request/response data

### Key Metrics
- Order creation rate
- Stock validation failures
- Average order value
- Error rates by type

## Error Handling

### Validation Errors (400)
- Missing required fields
- Invalid user/shop/product IDs
- Zero or negative quantity

### Business Logic Errors (409)
- Insufficient stock
- Product doesn't belong to shop

### System Errors (500)
- Database connection issues
- Stock update failures

## Security Considerations
- Input validation and sanitization
- SQL injection prevention with parameterized queries
- JSON data validation
- Error message sanitization

## Future Enhancements
- Order status updates (confirm, ship, deliver, cancel)
- Order history and listing
- Bulk order creation
- Order cancellation with stock release
- Payment integration
- Inventory rollback on order failure
- Real-time stock availability checks