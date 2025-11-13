# Shared Middleware Package

This package provides reusable middleware components for service-to-service authentication across the microservices architecture.

## Service Authentication Middleware

### Overview
The `ServiceAuthMiddleware` validates API key authentication for internal service endpoints using the `X-API-Key` header.

### Usage

#### Basic Usage (Environment-based)
```go
import (
    pkgMiddleware "github.com/Christyan39/test-eDot/pkg/middleware"
)

// Apply to specific routes
products.PATCH("/:id/hold-stock", productHandler.UpdateOnHoldStock, pkgMiddleware.ServiceAuthMiddleware())

// Or apply to route groups
internalAPI := api.Group("/internal", pkgMiddleware.ServiceAuthMiddleware())
internalAPI.POST("/sync", syncHandler.SyncData)
```

#### Custom API Key Usage
```go
import (
    pkgMiddleware "github.com/Christyan39/test-eDot/pkg/middleware"
)

// Use a specific API key instead of environment variable
customKey := "my-custom-service-key"
products.PATCH("/:id/special", handler.SpecialEndpoint, pkgMiddleware.ServiceAuthMiddlewareWithKey(customKey))
```

### Environment Configuration
Set the `API_KEY` environment variable in your service configuration:

```env
# Service API key for internal authentication
API_KEY=your-internal-service-api-key
```

### HTTP Request Format
Include the API key in the request header:

```bash
curl -X PATCH "http://localhost:8081/api/v1/products/1/hold-stock" \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-internal-service-api-key" \
  -d '{"on_hold_stock": 5}'
```

### Response Codes
- **200**: Success - Valid API key, request processed
- **401**: Unauthorized - Missing `X-API-Key` header
- **403**: Forbidden - Invalid API key value

### Error Response Format
```json
{
  "error": "API key required for internal endpoints"
}
```

```json
{
  "error": "Invalid API key"
}
```

## Integration Examples

### Product Service
```go
// Internal inventory management endpoint
products.PATCH("/:id/hold-stock", productHandler.UpdateOnHoldStock, pkgMiddleware.ServiceAuthMiddleware())
```

### User Service (Example)
```go
import pkgMiddleware "github.com/Christyan39/test-eDot/pkg/middleware"

// Internal user sync endpoint
users.POST("/sync", userHandler.SyncUsers, pkgMiddleware.ServiceAuthMiddleware())
```

### Order Service (Example)
```go
import pkgMiddleware "github.com/Christyan39/test-eDot/pkg/middleware"

// Internal order processing endpoints
internal := api.Group("/internal", pkgMiddleware.ServiceAuthMiddleware())
internal.POST("/orders/process", orderHandler.ProcessOrder)
internal.PATCH("/orders/:id/status", orderHandler.UpdateStatus)
```

## Security Best Practices

1. **Use strong API keys**: Generate cryptographically secure random strings
2. **Rotate keys regularly**: Update API keys periodically
3. **Environment isolation**: Use different keys per environment (dev, staging, prod)
4. **Secure transmission**: Always use HTTPS in production
5. **Logging**: Log authentication failures for security monitoring

## Testing

Use the provided test script to verify middleware functionality:

```bash
# Test the middleware authentication
./scripts/test-hold-stock.sh
```

This will test:
- Missing API key (401 response)
- Invalid API key (403 response)  
- Valid API key (200 response)