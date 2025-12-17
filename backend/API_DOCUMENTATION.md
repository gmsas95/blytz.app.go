# Blytz.live Marketplace API Documentation

## Overview

This document provides comprehensive API documentation for the Blytz.live Marketplace platform. The API is built with Go using the Gin framework and follows RESTful principles with OpenAPI 3.0 specification.

**Current Version:** 3.0.0  
**Base URL:** `http://localhost:8080/api/v1` (Development)  
**Base URL:** `https://api.blytz.live/api/v1` (Production)

## Features

- ✅ **Authentication & User Management** - JWT-based authentication with role-based access control
- ✅ **Product Catalog** - Full CRUD operations with advanced filtering and search
- ✅ **Shopping Cart System** - Guest and user cart support with 7-day expiration
- ✅ **Order Management** - Professional order lifecycle with admin dashboard
- 🔄 **Live Auction System** - In development (Phase 4)
- 🔄 **Real-time Bidding** - In development (Phase 4)
- 🔄 **Live Streaming** - In development (Phase 5)

## Quick Start

### Authentication Flow

1. **Register User**
   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/register \
     -H "Content-Type: application/json" \
     -d '{
       "email": "user@example.com",
       "password": "password123",
       "first_name": "John",
       "last_name": "Doe"
     }'
   ```

2. **Login**
   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{
       "email": "user@example.com",
       "password": "password123"
     }'
   ```

3. **Use Token**
   ```bash
   curl -X GET http://localhost:8080/api/v1/auth/profile \
     -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
   ```

### Shopping Cart Flow

1. **Create Cart** (Optional - auto-created)
   ```bash
   curl -X POST http://localhost:8080/api/v1/cart
   ```

2. **Add Item to Cart**
   ```bash
   curl -X POST http://localhost:8080/api/v1/cart/items \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "product_id": "550e8400-e29b-41d4-a716-446655440000",
       "quantity": 2
     }'
   ```

3. **Get Cart Details**
   ```bash
   curl -X GET http://localhost:8080/api/v1/cart
   ```

### Order Creation Flow

1. **Create Order from Cart**
   ```bash
   curl -X POST http://localhost:8080/api/v1/orders \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "cart_id": "550e8400-e29b-41d4-a716-446655440000",
       "shipping_address": {
         "first_name": "John",
         "last_name": "Doe",
         "address_line1": "123 Main St",
         "city": "New York",
         "state": "NY",
         "postal_code": "10001",
         "country": "USA"
       },
       "payment_method": "credit_card"
     }'
   ```

## API Endpoints

### System Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/health` | System health check | No |

### Authentication Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/auth/register` | User registration | No |
| POST | `/auth/login` | User login | No |
| POST | `/auth/refresh` | Refresh JWT token | No |
| GET | `/auth/profile` | Get user profile | Yes |
| POST | `/auth/change-password` | Change password | Yes |
| POST | `/auth/logout` | User logout | Yes |

### Product Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/products` | List products with filtering | No |
| POST | `/products` | Create product | Yes (Seller/Admin) |
| GET | `/products/{id}` | Get product by ID | No |
| PUT | `/products/{id}` | Update product | Yes (Owner/Admin) |
| DELETE | `/products/{id}` | Delete product | Yes (Owner/Admin) |
| GET | `/products/my-products` | List seller's products | Yes (Seller) |

### Shopping Cart Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/cart` | Get cart details | No |
| POST | `/cart` | Create cart | No |
| POST | `/cart/items` | Add item to cart | Yes |
| PUT | `/cart/items/{id}` | Update item quantity | Yes |
| DELETE | `/cart/items/{id}` | Remove item from cart | Yes |
| DELETE | `/cart` | Clear cart | Yes |
| POST | `/cart/merge` | Merge guest cart to user cart | Yes |

### Order Management Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/orders` | Create order from cart | Yes |
| GET | `/orders` | List user orders | Yes |
| GET | `/orders/{id}` | Get order by ID | Yes |
| PUT | `/orders/{id}/status` | Update order status | Yes (Seller/Admin) |
| DELETE | `/orders/{id}` | Cancel order | Yes |

### Admin Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/admin/orders/statistics` | Get order statistics | Yes (Admin) |

## Authentication

The API uses JWT (JSON Web Token) authentication. Include the access token in the Authorization header:

```
Authorization: Bearer YOUR_ACCESS_TOKEN
```

### Token Management

- **Access Token**: Short-lived (1 hour) token for API access
- **Refresh Token**: Long-lived token for getting new access tokens
- **Token Refresh**: Use `/auth/refresh` endpoint with refresh token

## Rate Limiting

- **Authentication endpoints**: 5 requests per minute
- **General endpoints**: 100 requests per minute
- **Per-IP rate limiting**: Applied to all endpoints

## Response Formats

### Success Response
```json
{
  "data": { ... },
  "message": "Operation completed successfully",
  "success": true
}
```

### Error Response
```json
{
  "error": "Invalid input",
  "details": "Email field is required",
  "code": "VALIDATION_ERROR"
}
```

### Pagination Response
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

## Data Models

### User Roles
- **buyer**: Can browse, bid, and purchase products
- **seller**: Can create and manage products
- **admin**: Full system access and management

### Product Status
- **draft**: Product not visible to buyers
- **active**: Product available for purchase
- **sold**: Product has been sold
- **cancelled**: Product listing cancelled

### Order Status
- **pending**: Order created, awaiting payment
- **processing**: Payment received, order being processed
- **shipped**: Order shipped to customer
- **delivered**: Order delivered to customer
- **cancelled**: Order cancelled

### Cart Features
- **Guest Carts**: 7-day expiration with token-based access
- **User Carts**: Persistent until converted to order
- **Cart Merging**: Transfer guest cart items to user cart
- **Stock Validation**: Real-time stock checking
- **Quantity Limits**: 1-10 items per product

## Error Codes

| Code | Description |
|------|-------------|
| 400 | Bad Request - Invalid input data |
| 401 | Unauthorized - Invalid or expired token |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource not found |
| 409 | Conflict - Resource already exists |
| 422 | Unprocessable Entity - Validation error |
| 429 | Too Many Requests - Rate limit exceeded |
| 500 | Internal Server Error - Server error |

## Testing

### Using Swagger UI
Access the interactive API documentation at:
```
http://localhost:8080/swagger/index.html
```

### Using cURL Examples
Detailed cURL examples for each endpoint are provided in the OpenAPI specification files.

### Using Postman
Import the OpenAPI JSON file into Postman for testing:
1. Open Postman
2. Click "Import"
3. Select the `openapi.json` file
4. Start testing the endpoints

## Development

### Local Development Setup
```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

### API Documentation Files
- `openapi.yaml` - Complete OpenAPI specification (YAML format)
- `openapi.json` - Complete OpenAPI specification (JSON format)
- `API_DOCUMENTATION.md` - This documentation file

### Updating Documentation
When adding new endpoints or modifying existing ones:
1. Update the OpenAPI specification files
2. Regenerate documentation if needed
3. Test all endpoints with the updated spec

## Performance

- **Response Time**: Sub-100ms average response time
- **Database Queries**: Optimized with proper indexing
- **Caching**: Redis caching for frequently accessed data
- **Pagination**: All list endpoints support pagination

## Security

- **HTTPS**: All production endpoints use HTTPS
- **Input Validation**: All inputs are validated and sanitized
- **SQL Injection**: Protected with parameterized queries
- **Rate Limiting**: Per-user and per-IP rate limiting
- **JWT Security**: Secure token generation and validation

## Support

For API support and questions:
- **Email**: support@blytzventures.com
- **GitHub Issues**: https://github.com/gmsas95/blytz.live.remake/issues
- **Documentation**: Check the OpenAPI specification files

## Versioning

The API uses semantic versioning (MAJOR.MINOR.PATCH):
- **Current Version**: 3.0.0
- **Base Path**: `/api/v1`
- **Backward Compatibility**: Maintained within major versions

## Roadmap

### Phase 4: Auction System (In Development)
- Live auction sessions
- Real-time bidding system
- WebSocket integration

### Phase 5: Live Streaming (Planned)
- LiveKit integration
- Video streaming capabilities
- Live chat during auctions

### Phase 6: Payment Processing (Planned)
- Multiple payment gateways
- Address management
- Payment confirmation system