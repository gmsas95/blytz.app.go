# API Specifications

## Base URL

```
Development: http://localhost:8080/api/v1
Staging:     https://api.staging.blytz.app/api/v1
Production:  https://api.blytz.app/api/v1
```

## Authentication

All protected endpoints require a Bearer token in the Authorization header:

```
Authorization: Bearer <jwt_token>
```

### Token Response Format
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 3600,
  "token_type": "Bearer"
}
```

## Response Format

### Success Response
```json
{
  "success": true,
  "data": { ... },
  "meta": {
    "timestamp": "2025-02-05T08:00:00Z"
  }
}
```

### Error Response
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": [
      { "field": "email", "message": "Email is required" }
    ]
  },
  "meta": {
    "timestamp": "2025-02-05T08:00:00Z"
  }
}
```

### Pagination
```json
{
  "success": true,
  "data": [ ... ],
  "meta": {
    "pagination": {
      "page": 1,
      "per_page": 20,
      "total": 100,
      "total_pages": 5,
      "has_next": true,
      "has_prev": false
    }
  }
}
```

---

## Authentication Endpoints

### POST /auth/register
Register a new user account.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "first_name": "Ahmad",
  "last_name": "Bin Abdullah",
  "phone": "+60123456789",
  "role": "buyer"
}
```

**Response (201):**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "first_name": "Ahmad",
      "last_name": "Bin Abdullah",
      "role": "buyer",
      "created_at": "2025-02-05T08:00:00Z"
    },
    "tokens": {
      "access_token": "eyJhbGci...",
      "refresh_token": "eyJhbGci...",
      "expires_in": 3600
    }
  }
}
```

**Errors:**
- `409` - Email already registered
- `400` - Invalid input (password too weak, etc.)

---

### POST /auth/login
Authenticate and receive tokens.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**Response (200):**
```json
{
  "success": true,
  "data": {
    "user": { ... },
    "tokens": { ... }
  }
}
```

**Errors:**
- `401` - Invalid credentials
- `403` - Account suspended

---

### POST /auth/refresh
Refresh access token using refresh token.

**Request:**
```json
{
  "refresh_token": "eyJhbGci..."
}
```

**Response (200):**
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGci...",
    "expires_in": 3600
  }
}
```

---

### POST /auth/logout
Invalidate tokens.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200):**
```json
{
  "success": true,
  "data": {
    "message": "Logged out successfully"
  }
}
```

---

### GET /auth/me
Get current user profile.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200):**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "first_name": "Ahmad",
    "last_name": "Bin Abdullah",
    "phone": "+60123456789",
    "avatar_url": "https://r2.blytz.app/avatars/uuid.jpg",
    "role": "seller",
    "seller_profile": {
      "store_name": "Ahmad's Collectibles",
      "store_slug": "ahmad-collectibles",
      "rating": 4.8,
      "verified": true
    },
    "created_at": "2025-02-05T08:00:00Z"
  }
}
```

---

## User Endpoints

### PUT /users/me
Update user profile.

**Request:**
```json
{
  "first_name": "Ahmad",
  "last_name": "Bin Abdullah",
  "phone": "+60123456789",
  "avatar_url": "https://..."
}
```

**Response (200):**
```json
{
  "success": true,
  "data": {
    "user": { ... }
  }
}
```

---

### GET /users/me/addresses
Get user's saved addresses.

**Response (200):**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "label": "Home",
      "recipient_name": "Ahmad Bin Abdullah",
      "phone": "+60123456789",
      "address_line1": "123 Jalan Bukit Bintang",
      "address_line2": "Apt 5B",
      "city": "Kuala Lumpur",
      "state": "Wilayah Persekutuan",
      "postal_code": "50200",
      "is_default": true
    }
  ]
}
```

---

### POST /users/me/addresses
Add new address.

**Request:**
```json
{
  "label": "Office",
  "recipient_name": "Ahmad Bin Abdullah",
  "phone": "+60123456789",
  "address_line1": "456 Menara KL",
  "city": "Kuala Lumpur",
  "state": "Wilayah Persekutuan",
  "postal_code": "50450",
  "is_default": false
}
```

**Response (201):**
```json
{
  "success": true,
  "data": {
    "address": { ... }
  }
}
```

---

## Seller Endpoints

### POST /sellers/register
Register as a seller (requires authentication).

**Request:**
```json
{
  "store_name": "Ahmad's Collectibles",
  "description": "Rare collectibles and antiques",
  "logo_url": "https://...",
  "banner_url": "https://..."
}
```

**Response (201):**
```json
{
  "success": true,
  "data": {
    "seller": {
      "id": "uuid",
      "store_name": "Ahmad's Collectibles",
      "store_slug": "ahmad-collectibles",
      "status": "pending_verification"
    }
  }
}
```

---

### GET /sellers/:slug
Get seller public profile.

**Response (200):**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "store_name": "Ahmad's Collectibles",
    "store_slug": "ahmad-collectibles",
    "description": "...",
    "logo_url": "...",
    "banner_url": "...",
    "rating": 4.8,
    "review_count": 150,
    "total_sales": 500,
    "verified": true,
    "followers_count": 1200,
    "is_following": false
  }
}
```

---

### POST /sellers/:slug/follow
Follow a seller.

**Response (200):**
```json
{
  "success": true,
  "data": {
    "message": "Now following Ahmad's Collectibles"
  }
}
```

---

## Product Endpoints

### GET /products
List products with filters.

**Query Parameters:**
```
?page=1&per_page=20&category=fashion&min_price=10&max_price=100&sort=newest&q=watch
```

**Response (200):**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "name": "Vintage Rolex Submariner",
      "slug": "vintage-rolex-submariner",
      "description": "...",
      "condition": "used",
      "base_price": 15000.00,
      "stock_quantity": 1,
      "seller": {
        "id": "uuid",
        "store_name": "Ahmad's Collectibles",
        "store_slug": "ahmad-collectibles"
      },
      "category": {
        "id": "uuid",
        "name": "Watches",
        "slug": "watches"
      },
      "images": [
        {
          "url": "https://r2.blytz.app/products/uuid-1.jpg",
          "thumbnail_url": "https://r2.blytz.app/products/uuid-1-thumb.jpg",
          "is_primary": true
        }
      ],
      "created_at": "2025-02-05T08:00:00Z"
    }
  ],
  "meta": {
    "pagination": { ... }
  }
}
```

---

### GET /products/:slug
Get product details.

**Response (200):**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "Vintage Rolex Submariner",
    "slug": "vintage-rolex-submariner",
    "description": "...",
    "condition": "used",
    "base_price": 15000.00,
    "compare_at_price": 18000.00,
    "stock_quantity": 1,
    "attributes": {
      "brand": "Rolex",
      "model": "Submariner",
      "year": "1985"
    },
    "seller": { ... },
    "category": { ... },
    "images": [ ... ],
    "active_auction": {
      "id": "uuid",
      "start_price": 10000.00,
      "current_bid": 12500.00,
      "bid_count": 15,
      "end_time": "2025-02-05T12:00:00Z"
    }
  }
}
```

---

### POST /products
Create product (seller only).

**Request:**
```json
{
  "name": "Vintage Rolex Submariner",
  "description": "Original 1985 model...",
  "category_id": "uuid",
  "condition": "used",
  "base_price": 15000.00,
  "stock_quantity": 1,
  "attributes": {
    "brand": "Rolex",
    "model": "Submariner",
    "year": "1985"
  },
  "images": [
    { "url": "https://r2.blytz.app/products/uuid-1.jpg", "is_primary": true },
    { "url": "https://r2.blytz.app/products/uuid-2.jpg", "is_primary": false }
  ]
}
```

**Response (201):**
```json
{
  "success": true,
  "data": {
    "product": { ... }
  }
}
```

---

## Auction Endpoints

### GET /auctions
List active auctions.

**Query Parameters:**
```
?status=active&category=watches&sort=ending_soon&page=1
```

**Response (200):**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "product": {
        "id": "uuid",
        "name": "Vintage Rolex Submariner",
        "image_url": "https://r2.blytz.app/products/uuid-1-thumb.jpg"
      },
      "seller": {
        "id": "uuid",
        "store_name": "Ahmad's Collectibles"
      },
      "start_price": 10000.00,
      "current_bid": 12500.00,
      "bid_count": 15,
      "buy_now_price": 20000.00,
      "start_time": "2025-02-05T08:00:00Z",
      "end_time": "2025-02-05T12:00:00Z",
      "time_remaining": 7200,
      "status": "active",
      "stream": {
        "id": "uuid",
        "status": "live",
        "viewer_count": 45
      }
    }
  ]
}
```

---

### GET /auctions/:id
Get auction details.

**Response (200):**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "product": { ... },
    "seller": { ... },
    "start_price": 10000.00,
    "reserve_price": 12000.00,
    "buy_now_price": 20000.00,
    "current_bid": 12500.00,
    "highest_bidder": {
      "id": "uuid",
      "first_name": "Mohammed",
      "last_name": "H."
    },
    "bid_count": 15,
    "bid_history": [
      {
        "id": "uuid",
        "bidder": { ... },
        "amount": 12500.00,
        "created_at": "2025-02-05T09:30:00Z"
      }
    ],
    "start_time": "2025-02-05T08:00:00Z",
    "end_time": "2025-02-05T12:00:00Z",
    "time_remaining": 7200,
    "status": "active",
    "stream": { ... }
  }
}
```

---

### POST /auctions
Create auction (seller only).

**Request:**
```json
{
  "product_id": "uuid",
  "start_price": 10000.00,
  "reserve_price": 12000.00,
  "buy_now_price": 20000.00,
  "min_bid_increment": 500.00,
  "start_time": "2025-02-05T08:00:00Z",
  "duration_minutes": 240
}
```

**Response (201):**
```json
{
  "success": true,
  "data": {
    "auction": { ... }
  }
}
```

---

### POST /auctions/:id/bids
Place a bid (requires authentication).

**Request:**
```json
{
  "amount": 13000.00,
  "is_auto_bid": false
}
```

**Response (201):**
```json
{
  "success": true,
  "data": {
    "bid": {
      "id": "uuid",
      "amount": 13000.00,
      "is_winning": true,
      "created_at": "2025-02-05T09:35:00Z"
    },
    "auction": {
      "current_bid": 13000.00,
      "bid_count": 16,
      "highest_bidder": { ... }
    }
  }
}
```

**Errors:**
- `400` - Bid too low (must be higher than current + increment)
- `400` - Auction not active
- `409` - Outbid (someone bid higher simultaneously)

---

## Stream Endpoints

### GET /streams
List active and upcoming streams.

**Query Parameters:**
```
?status=live&category=fashion&sort=popular
```

**Response (200):**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "title": "Rare Watch Collection Auction!",
      "description": "Join me for an exciting auction...",
      "seller": {
        "id": "uuid",
        "store_name": "Ahmad's Collectibles",
        "store_slug": "ahmad-collectibles"
      },
      "status": "live",
      "thumbnail_url": "https://r2.blytz.app/streams/uuid-thumb.jpg",
      "viewer_count": 45,
      "max_viewers": 120,
      "started_at": "2025-02-05T08:00:00Z",
      "scheduled_at": null,
      "products_count": 5,
      "featured_products": [ ... ]
    }
  ]
}
```

---

### GET /streams/:id
Get stream details with LiveKit token.

**Response (200):**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "title": "Rare Watch Collection Auction!",
    "seller": { ... },
    "status": "live",
    "viewer_count": 45,
    "livekit": {
      "room_name": "stream_uuid",
      "token": "eyJhbGciOiJIUzI1NiIs...",
      "ws_url": "wss://livekit.blytz.app"
    },
    "products": [ ... ],
    "active_auctions": [ ... ],
    "is_following_seller": false
  }
}
```

---

### POST /streams
Create stream (seller only).

**Request:**
```json
{
  "title": "Rare Watch Collection Auction!",
  "description": "Join me for...",
  "scheduled_at": "2025-02-05T10:00:00Z",
  "product_ids": ["uuid1", "uuid2"]
}
```

**Response (201):**
```json
{
  "success": true,
  "data": {
    "stream": { ... },
    "livekit": {
      "room_name": "stream_uuid",
      "token": "eyJhbGci...",
      "ws_url": "wss://livekit.blytz.app"
    }
  }
}
```

---

## Order Endpoints

### GET /orders
Get user's orders.

**Query Parameters:**
```
?status=pending&page=1
```

**Response (200):**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "order_number": "BLTZ-20250205-001",
      "status": "shipped",
      "total_amount": 13500.00,
      "seller": {
        "store_name": "Ahmad's Collectibles"
      },
      "items": [
        {
          "product_name": "Vintage Rolex Submariner",
          "product_image": "...",
          "quantity": 1,
          "unit_price": 13000.00
        }
      ],
      "shipment": {
        "tracking_number": "NVMY123456789",
        "carrier": "NinjaVan",
        "status": "in_transit"
      },
      "created_at": "2025-02-05T12:05:00Z"
    }
  ]
}
```

---

### GET /orders/:id
Get order details.

**Response (200):**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "order_number": "BLTZ-20250205-001",
    "status": "shipped",
    "subtotal": 13000.00,
    "shipping_cost": 500.00,
    "total_amount": 13500.00,
    "shipping_address": { ... },
    "items": [ ... ],
    "payment": {
      "method": "card",
      "status": "completed",
      "paid_at": "2025-02-05T12:05:30Z"
    },
    "shipment": {
      "tracking_number": "NVMY123456789",
      "carrier": "NinjaVan",
      "label_url": "https://...",
      "status": "in_transit",
      "tracking_history": [
        {
          "status": "picked_up",
          "location": "Kuala Lumpur",
          "timestamp": "2025-02-05T14:00:00Z"
        }
      ]
    },
    "timeline": [
      { "status": "created", "timestamp": "..." },
      { "status": "paid", "timestamp": "..." },
      { "status": "shipped", "timestamp": "..." }
    ]
  }
}
```

---

## Payment Endpoints

### POST /payments/intent
Create payment intent (Stripe).

**Request:**
```json
{
  "order_id": "uuid",
  "payment_method": "card"
}
```

**Response (200):**
```json
{
  "success": true,
  "data": {
    "client_secret": "pi_123_secret_456",
    "publishable_key": "pk_test_..."
  }
}
```

---

### POST /payments/confirm
Confirm payment after Stripe processing.

**Request:**
```json
{
  "order_id": "uuid",
  "payment_intent_id": "pi_123"
}
```

**Response (200):**
```json
{
  "success": true,
  "data": {
    "status": "completed",
    "order": { ... }
  }
}
```

---

## Notification Endpoints

### GET /notifications
Get user notifications.

**Query Parameters:**
```
?unread_only=true&page=1
```

**Response (200):**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "type": "bid_outbid",
      "title": "You've been outbid!",
      "message": "Someone placed a higher bid on Vintage Rolex Submariner",
      "data": {
        "auction_id": "uuid",
        "product_name": "Vintage Rolex Submariner"
      },
      "action_url": "/auctions/uuid",
      "is_read": false,
      "created_at": "2025-02-05T09:35:00Z"
    }
  ],
  "meta": {
    "unread_count": 3
  }
}
```

---

### PUT /notifications/:id/read
Mark notification as read.

**Response (200):**
```json
{
  "success": true,
  "data": {
    "message": "Notification marked as read"
  }
}
```

---

## WebSocket Events (Socket.io)

### Connection
```javascript
const socket = io('wss://api.blytz.app', {
  auth: { token: 'jwt_token' }
});
```

### Join Stream
```javascript
socket.emit('stream:join', { stream_id: 'uuid' });
```

### Chat Message
```javascript
// Send
socket.emit('chat:message', {
  stream_id: 'uuid',
  content: 'Great product!'
});

// Receive
socket.on('chat:message', (data) => {
  console.log(data);
  // { id, user_id, user_name, content, created_at }
});
```

### Bid Updates
```javascript
// Receive bid updates
socket.on('auction:bid', (data) => {
  console.log(data);
  // { auction_id, amount, bidder_id, bid_count, time_remaining }
});

// Auction ended
socket.on('auction:ended', (data) => {
  console.log(data);
  // { auction_id, winner_id, winning_amount }
});
```

### Stream Events
```javascript
socket.on('stream:viewer_count', (data) => {
  console.log(data);
  // { stream_id, count }
});

socket.on('stream:ended', (data) => {
  console.log(data);
  // { stream_id }
});
```

---

## Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `UNAUTHORIZED` | 401 | Invalid or missing authentication |
| `FORBIDDEN` | 403 | Insufficient permissions |
| `NOT_FOUND` | 404 | Resource not found |
| `VALIDATION_ERROR` | 400 | Invalid input data |
| `RATE_LIMITED` | 429 | Too many requests |
| `AUCTION_ENDED` | 400 | Auction already ended |
| `BID_TOO_LOW` | 400 | Bid below minimum increment |
| `INSUFFICIENT_FUNDS` | 400 | Payment declined |
| `SERVER_ERROR` | 500 | Internal server error |

---

*Last updated: 2025-02-05*
