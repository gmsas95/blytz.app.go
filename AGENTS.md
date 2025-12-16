# Blytz.live.remake - Agents Guide

## Project Overview

Blytz.live.remake is a modern live marketplace platform designed for real-time auctions, bidding, and live streaming capabilities. The platform connects buyers and sellers through interactive live sessions with real-time product demonstrations and instant bidding functionality.

### Business Objectives
- **Primary**: Create engaging live commerce experiences with real-time auctions
- **Secondary**: Build a scalable marketplace supporting high concurrent users
- **Tertiary**: Enable mobile-first experiences for on-the-go bidding

## Current Development Status

### Phase 1: Backend Foundation ✅ COMPLETED (2025-12-15)
- Project structure with clean architecture
- Docker configuration for development
- PostgreSQL + Redis database setup with SQLite fallback
- GORM ORM integration with UUID primary keys
- Core data models (User, Category, Product)
- API foundation with Gin framework
- CORS and middleware support
- Health check endpoint
- Environment configuration management
- Logging infrastructure with logrus

### Phase 2: Authentication System ✅ COMPLETED (2025-12-15)
- JWT token management (access + refresh tokens)
- User registration and login with email validation
- Secure password hashing with bcrypt
- Protected routes with middleware
- Role-based access control (buyer/seller/admin)
- Rate limiting (5 req/min for auth, 100 req/min general)
- Input validation and structured error handling
- User profile management
- Password change functionality
- SQLite in-memory database for demo

**API Endpoints Implemented:**
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User authentication
- `POST /api/v1/auth/refresh` - Token refresh
- `GET /api/v1/auth/profile` - User profile (protected)
- `POST /api/v1/auth/change-password` - Password change (protected)
- `POST /api/v1/auth/logout` - User logout (protected)

### Phase 3: Product Management ✅ COMPLETED (2025-12-16)
Features implemented:
- Product CRUD operations (Create, Read, Update, Delete)
- JSON field support for images, specifications, and shipping info
- Category associations with validation
- Advanced search and filtering with pagination
- Product status management (draft, active, sold, cancelled)
- Price management (starting price, reserve price, buy now price)
- View count tracking
- Ownership validation and security controls

**API Endpoints Implemented:**
- `GET /api/v1/products` - List products with filtering and pagination
- `POST /api/v1/products` - Create product (protected)
- `GET /api/v1/products/:id` - Get product details
- `PUT /api/v1/products/:id` - Update product (protected)
- `DELETE /api/v1/products/:id` - Delete product (protected)
- `GET /api/v1/products/my-products` - List seller's products (protected)

### Phase 4: Auction System 🔄 IN PROGRESS
Target Start: TBD
Features to implement:
- Live auction sessions
- Real-time bidding system
- WebSocket integration
- Bid history tracking
- Reserve price enforcement
- Winner determination
- Auto-bid functionality

### Phase 5: Basic E-commerce System 📋 PLANNED
Target Start: After Phase 4
Timeline: ~14-19 days
Features to implement:
- Shopping cart system (Cart, CartItem models)
- Order management system (Order, OrderItem models)
- Basic payment processing (Payment model + Stripe integration)
- Address management (UserAddress model)
- Inventory management (Stock, StockReservation models)
- Tax and shipping calculations
- Order status tracking (pending, processing, shipped, delivered, cancelled)

**API Endpoints to Implement:**
- Cart: GET /api/v1/cart, POST /api/v1/cart/items, PUT /api/v1/cart/items/:id, DELETE /api/v1/cart/items/:id
- Orders: POST /api/v1/orders, GET /api/v1/orders, GET /api/v1/orders/:id
- Payments: POST /api/v1/payments/intent, POST /api/v1/payments/confirm, GET /api/v1/payments/:id
- Addresses: GET /api/v1/addresses, POST /api/v1/addresses, PUT /api/v1/addresses/:id, DELETE /api/v1/addresses/:id

### Phase 6: Live Streaming 📋 PLANNED
Target Start: After Phase 5
- LiveKit integration
- Video streaming capabilities
- Live chat during auctions

## E-commerce Implementation Plan

### E-commerce Architecture Overview

The e-commerce system will be built as an extension to the existing backend architecture, maintaining the same clean architecture principles and modular design.

### Module Structure for E-commerce

```
internal/
├── cart/                           # Shopping Cart Module
│   ├── models.go                    # Cart, CartItem DTOs
│   ├── service.go                   # Cart business logic
│   └── handlers.go                  # Cart HTTP handlers
├── orders/                         # Order Management Module
│   ├── models.go                    # Order, OrderItem, OrderStatus
│   ├── service.go                   # Order business logic
│   └── handlers.go                  # Order HTTP handlers
├── payments/                       # Payment Processing Module
│   ├── models.go                    # Payment, PaymentMethod, Transaction
│   ├── service.go                   # Payment business logic
│   └── handlers.go                  # Payment HTTP handlers
├── addresses/                      # Address Management Module
│   ├── models.go                    # UserAddress, AddressType
│   ├── service.go                   # Address business logic
│   └── handlers.go                  # Address HTTP handlers
└── inventory/                      # Inventory Management Module
    ├── models.go                    # Stock, StockReservation
    ├── service.go                   # Inventory business logic
    └── handlers.go                  # Inventory HTTP handlers
```

### Database Schema for E-commerce

#### Order Management Tables
```go
// Order represents a customer order
type Order struct {
    ID              uuid.UUID      `gorm:"primaryKey;type:uuid" json:"id"`
    UserID          uuid.UUID      `gorm:"not null;references:ID" json:"user_id"`
    Status          string         `gorm:"not null;default:'pending'" json:"status"` // pending, processing, shipped, delivered, cancelled
    TotalAmount     float64        `gorm:"not null" json:"total_amount"`
    Subtotal        float64        `gorm:"not null" json:"subtotal"`
    TaxAmount       float64        `gorm:"default:0" json:"tax_amount"`
    ShippingCost    float64        `gorm:"default:0" json:"shipping_cost"`
    DiscountAmount  float64        `gorm:"default:0" json:"discount_amount"`
    ShippingAddress *Address       `gorm:"embedded;embeddedPrefix:shipping_" json:"shipping_address"`
    BillingAddress  *Address       `gorm:"embedded;embeddedPrefix:billing_" json:"billing_address"`
    PaymentID       *uuid.UUID     `gorm:"references:ID" json:"payment_id"`
    TrackingNumber  *string        `json:"tracking_number"`
    Notes           *string        `json:"notes"`
    common.BaseModel                 // ID, CreatedAt, UpdatedAt, DeletedAt
}

// OrderItem represents items in an order
type OrderItem struct {
    ID          uuid.UUID  `gorm:"primaryKey;type:uuid" json:"id"`
    OrderID     uuid.UUID  `gorm:"not null;references:ID" json:"order_id"`
    ProductID   uuid.UUID  `gorm:"not null;references:ID" json:"product_id"`
    Quantity    int        `gorm:"not null" json:"quantity"`
    UnitPrice   float64    `gorm:"not null" json:"unit_price"`
    Total       float64    `gorm:"not null" json:"total"`
    Product     Product    `gorm:"foreignKey:ProductID" json:"product,omitempty"`
    common.BaseModel
}
```

#### Shopping Cart Tables
```go
// Cart represents a shopping cart
type Cart struct {
    ID     uuid.UUID  `gorm:"primaryKey;type:uuid" json:"id"`
    UserID *uuid.UUID `gorm:"references:ID" json:"user_id,omitempty"` // nullable for guest carts
    Token  string     `gorm:"uniqueIndex;not null" json:"token"`    // for guest carts
    ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
    Items  []CartItem `gorm:"foreignKey:CartID" json:"items"`
    common.BaseModel
}

// CartItem represents items in a cart
type CartItem struct {
    ID        uuid.UUID  `gorm:"primaryKey;type:uuid" json:"id"`
    CartID     uuid.UUID  `gorm:"not null;references:ID" json:"cart_id"`
    ProductID  uuid.UUID  `gorm:"not null;references:ID" json:"product_id"`
    Quantity   int        `gorm:"not null" json:"quantity"`
    AddedAt    time.Time  `gorm:"autoCreateTime" json:"added_at"`
    Product    Product    `gorm:"foreignKey:ProductID" json:"product,omitempty"`
    common.BaseModel
}
```

#### Payment Processing Tables
```go
// Payment represents a payment transaction
type Payment struct {
    ID             uuid.UUID  `gorm:"primaryKey;type:uuid" json:"id"`
    OrderID        uuid.UUID  `gorm:"not null;references:ID" json:"order_id"`
    PaymentMethod  string     `gorm:"not null" json:"payment_method"` // stripe, paypal, credit_card
    Amount         float64    `gorm:"not null" json:"amount"`
    Currency       string     `gorm:"not null;default:'USD'" json:"currency"`
    Status         string     `gorm:"default:'pending'" json:"status"` // pending, processing, completed, failed, refunded
    TransactionID  string     `gorm:"uniqueIndex" json:"transaction_id"`
    GatewayRef     string     `json:"gateway_ref"` // stripe_payment_id, paypal_id, etc.
    FailureReason  *string    `json:"failure_reason"`
    RefundedAmount float64    `gorm:"default:0" json:"refunded_amount"`
    common.BaseModel
}

// PaymentMethod represents a user's saved payment method
type PaymentMethod struct {
    ID         uuid.UUID  `gorm:"primaryKey;type:uuid" json:"id"`
    UserID     uuid.UUID  `gorm:"not null;references:ID" json:"user_id"`
    Type       string     `gorm:"not null" json:"type"` // credit_card, paypal, bank_account
    Provider   string     `gorm:"not null" json:"provider"` // stripe, paypal, etc.
    MethodRef  string     `gorm:"not null" json:"method_ref"` // tokenized reference
    IsDefault  bool       `gorm:"default:false" json:"is_default"`
    Last4      *string    `json:"last4"`
    ExpiryDate *time.Time `json:"expiry_date"`
    common.BaseModel
}
```

#### Address Management Tables
```go
// Address represents a user's shipping/billing address
type Address struct {
    ID          uuid.UUID  `gorm:"primaryKey;type:uuid" json:"id"`
    UserID      uuid.UUID  `gorm:"not null;references:ID" json:"user_id"`
    Type        string     `gorm:"not null" json:"type"` // shipping, billing
    Label       string     `gorm:"not null" json:"label"` // Home, Work, etc.
    FirstName   string     `gorm:"not null" json:"first_name"`
    LastName    string     `gorm:"not null" json:"last_name"`
    Company     *string    `json:"company"`
    AddressLine1 string     `gorm:"not null" json:"address_line1"`
    AddressLine2 *string   `json:"address_line2"`
    City        string     `gorm:"not null" json:"city"`
    State       string     `gorm:"not null" json:"state"`
    PostalCode  string     `gorm:"not null" json:"postal_code"`
    Country     string     `gorm:"not null" json:"country"`
    Phone       *string    `json:"phone"`
    IsDefault   bool       `gorm:"default:false" json:"is_default"`
    common.BaseModel
}
```

#### Inventory Management Tables
```go
// Stock represents product inventory
type Stock struct {
    ID            uuid.UUID  `gorm:"primaryKey;type:uuid" json:"id"`
    ProductID     uuid.UUID  `gorm:"not null;uniqueIndex;references:ID" json:"product_id"`
    Quantity      int        `gorm:"not null;default:0" json:"quantity"`
    Reserved      int        `gorm:"not null;default:0" json:"reserved"` // reserved in carts
    Available     int        `gorm:"not null;default:0" json:"available"` // calculated field
    LowStockAlert int        `gorm:"default:10" json:"low_stock_alert"`
    LastUpdated   time.Time   `gorm:"autoUpdateTime" json:"last_updated"`
    common.BaseModel
}

// StockReservation represents stock reserved in carts
type StockReservation struct {
    ID        uuid.UUID  `gorm:"primaryKey;type:uuid" json:"id"`
    StockID   uuid.UUID  `gorm:"not null;references:ID" json:"stock_id"`
    CartID    *uuid.UUID `gorm:"references:ID" json:"cart_id"`
    OrderID   *uuid.UUID `gorm:"references:ID" json:"order_id"`
    Quantity  int        `gorm:"not null" json:"quantity"`
    ExpiresAt time.Time   `gorm:"not null" json:"expires_at"`
    common.BaseModel
}
```

### E-commerce API Design

#### Cart Management API
```
GET    /api/v1/cart              # Get user's cart or guest cart by token
POST   /api/v1/cart              # Create new cart (guest)
POST   /api/v1/cart/items        # Add item to cart
PUT    /api/v1/cart/items/:id     # Update item quantity
DELETE /api/v1/cart/items/:id     # Remove item from cart
DELETE /api/v1/cart              # Clear entire cart
POST   /api/v1/cart/merge        # Merge guest cart to user cart
```

#### Order Management API
```
POST   /api/v1/orders            # Create order from cart
GET    /api/v1/orders            # List user orders with pagination
GET    /api/v1/orders/:id        # Get order details
PUT    /api/v1/orders/:id/cancel # Cancel order (user)
PUT    /api/v1/orders/:id/status # Update order status (admin/seller)
GET    /api/v1/orders/:id/tracking # Get tracking info
```

#### Payment Processing API
```
POST   /api/v1/payments/intent   # Create payment intent (Stripe)
POST   /api/v1/payments/confirm  # Confirm payment processing
GET    /api/v1/payments/:id       # Get payment status
POST   /api/v1/payments/refund    # Process refund (admin)
GET    /api/v1/payments/methods   # Get user's saved payment methods
POST   /api/v1/payments/methods   # Save new payment method
DELETE /api/v1/payments/methods/:id # Remove payment method
```

#### Address Management API
```
GET    /api/v1/addresses         # Get user addresses
POST   /api/v1/addresses         # Add new address
GET    /api/v1/addresses/:id     # Get specific address
PUT    /api/v1/addresses/:id     # Update address
DELETE /api/v1/addresses/:id     # Delete address
PUT    /api/v1/addresses/:id/default # Set as default address
```

#### Inventory Management API
```
GET    /api/v1/inventory/products/:id # Get product stock info
PUT    /api/v1/inventory/products/:id # Update stock levels (admin)
GET    /api/v1/inventory/low-stock  # Get low stock products (admin)
GET    /api/v1/inventory/reservations # Get stock reservations (admin)
```

### Implementation Priority and Timeline

#### Priority 1: Shopping Cart (3-4 days)
- Cart creation and management
- Item operations (add, update, remove)
- Guest cart support with tokens
- Cart persistence and expiration
- Guest cart merging on login

#### Priority 2: Order Management (4-5 days)
- Order creation from cart
- Order status tracking
- Order history and details
- Order cancellation
- Admin order management

#### Priority 3: Payment Processing (3-4 days)
- Stripe integration for payments
- Payment intent creation
- Payment confirmation and webhooks
- Basic refund processing
- Payment method saving

#### Priority 4: Address Management (2-3 days)
- Address CRUD operations
- Default address selection
- Address validation
- Guest checkout addresses

#### Priority 5: Inventory Management (2-3 days)
- Stock tracking and reservation
- Stock updates on purchase
- Low stock alerts
- Admin stock management

**Total Estimated Timeline: 14-19 days**

### Payment Gateway Integration

#### Stripe Integration
```go
// Stripe service interface
type StripeService interface {
    CreatePaymentIntent(amount int64, currency string, metadata map[string]string) (*stripe.PaymentIntent, error)
    ConfirmPaymentIntent(paymentIntentID string) (*stripe.PaymentIntent, error)
    ProcessRefund(chargeID string, amount int64) (*stripe.Refund, error)
    CreateCustomer(email string, name string) (*stripe.Customer, error)
    AttachPaymentMethod(customerID string, paymentMethodID string) (*stripe.PaymentMethod, error)
}
```

#### Webhook Endpoints
```go
// Webhook handlers for payment confirmations
POST /api/v1/webhooks/stripe    # Stripe webhook endpoint
// Events: payment_intent.succeeded, payment_intent.payment_failed, charge.dispute.created
```

### Security Considerations for E-commerce

1. **Payment Security**
   - PCI compliance considerations
   - Secure handling of payment data
   - Use of payment gateway tokens instead of raw card data

2. **Order Security**
   - Order ownership validation
   - Secure order status transitions
   - Protection against order manipulation

3. **Cart Security**
   - Cart ownership validation
   - Prevent cart manipulation
   - Secure guest cart tokens

4. **Address Security**
   - Address ownership validation
   - PII protection
   - Secure address storage

### Performance Considerations

1. **Database Optimization**
   - Proper indexing on frequently queried fields
   - Optimized queries for order history
   - Cart item count optimization

2. **Caching Strategy**
   - Redis for cart sessions
   - Product stock caching
   - User preference caching

3. **Scalability**
   - Async payment processing
   - Inventory reservation timeouts
   - Order processing queue

## Architecture Vision

### Target Architecture: Well-Structured Monolith with Future Extraction Points

The project avoids the distributed monolith anti-pattern by using a unified backend with proper module boundaries. This maintains simplicity while allowing future microservices extraction if business needs justify it.

### Backend (Go)
- **Framework**: Gin (confirmed from PRD)
- **Architecture**: Clean architecture with modular domains
- **Database**: PostgreSQL 17.7 with GORM
- **Cache**: Redis 8+ for sessions and caching
- **Real-time**: LiveKit for streaming, Gorilla WebSocket for messaging
- **Queue**: Redis Streams for async processing

### Frontend (Next.js)
- **Version**: Next.js 16+ with App Router (updated from PRD)
- **Language**: TypeScript 5+
- **Styling**: Tailwind CSS
- **Components**: Radix UI + custom components
- **State**: React Query + Context API
- **Forms**: React Hook Form + Zod validation
- **Real-time**: Socket.IO client

### Mobile (React Native)
- **Framework**: React Native with Expo
- **State**: Redux Toolkit + RTK Query
- **Navigation**: React Navigation 6
- **UI Components**: React Native Elements or NativeBase
- **Push Notifications**: Firebase Cloud Messaging
- **Authentication**: Firebase Auth or custom JWT

## Expected Project Structure

```
blytz-live-latest/
├── backend/                          # Unified Go Backend
│   ├── cmd/
│   │   └── server/
│   │       └── main.go              # Single entry point
│   ├── internal/                    # Private packages
│   │   ├── auth/                   # Auth module
│   │   │   ├── handlers.go
│   │   │   ├── models.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   ├── products/                # Product module
│   │   │   ├── handlers.go
│   │   │   ├── models.go
│   │   │   └── service.go
│   │   ├── cart/                    # Shopping cart module
│   │   │   ├── handlers.go
│   │   │   ├── models.go
│   │   │   └── service.go
│   │   ├── orders/                  # Order management module
│   │   │   ├── handlers.go
│   │   │   ├── models.go
│   │   │   └── service.go
│   │   ├── payments/                # Payment processing module
│   │   │   ├── handlers.go
│   │   │   ├── models.go
│   │   │   └── service.go
│   │   ├── addresses/               # Address management module
│   │   │   ├── handlers.go
│   │   │   ├── models.go
│   │   │   └── service.go
│   │   ├── inventory/              # Inventory management module
│   │   │   ├── handlers.go
│   │   │   ├── models.go
│   │   │   └── service.go
│   │   ├── auctions/                # Auction system module
│   │   ├── chat/                   # Chat module
│   │   ├── livekit/                # LiveKit integration module
│   │   └── logistics/              # Logistics module
│   │   ├── config/                  # Configuration
│   │   ├── middleware/              # HTTP middleware
│   │   ├── database/               # Database setup
│   │   └── common/                 # Shared utilities
│   ├── pkg/                       # Public packages
│   │   ├── http/                   # HTTP utilities
│   │   ├── logging/                # Logging setup
│   │   └── validation/            # Input validation
│   ├── migrations/                 # Database migrations
│   ├── tests/                     # Test files
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   └── README.md
├── frontend/                      # Next.js Web Application
│   ├── src/
│   │   ├── app/                   # App Router pages
│   │   ├── components/             # React components
│   │   ├── hooks/                  # Custom hooks
│   │   ├── lib/                   # Utilities
│   │   └── types/                 # TypeScript types
│   ├── public/
│   ├── package.json
│   ├── next.config.js
│   ├── tailwind.config.js
│   └── Dockerfile
├── mobile/                        # React Native Application
│   ├── src/
│   │   ├── components/
│   │   ├── screens/
│   │   ├── navigation/
│   │   ├── services/
│   │   └── utils/
│   ├── android/
│   ├── ios/
│   ├── package.json
│   └── metro.config.js
├── docker-compose.yml
├── docker-compose.dev.yml
├── docker-compose.prod.yml
├── .env.example
├── README.md
└── ARCHITECTURE.md
```

## Development Commands (Planned)

### Backend
```bash
cd backend
go mod tidy                    # Install dependencies
go run cmd/server/main.go      # Run development server
go test ./...                  # Run all tests
go build ./cmd/server          # Build binary
```

### Frontend
```bash
cd frontend
npm install                    # Install dependencies
npm run dev                    # Development server
npm run build                  # Production build
npm test                       # Run tests
npm run lint                   # Lint code
```

### Mobile
```bash
cd mobile
npm install                    # Install dependencies
npx expo start                 # Start development server
npx expo run android           # Run on Android
npx expo run ios               # Run on iOS
```

## Key Implementation Details

### Database Schema
The project uses a unified PostgreSQL database with the following core entities:
- Users and Authentication (with role-based access)
- Categories (hierarchical structure)
- Products (with condition, pricing, and media)
- Auction Sessions (with LiveKit integration)
- Bids (with auto-bid support)
- Orders and Payments (with multiple gateway support)
- Chat Messages (for auction interaction)
- User Sessions (for caching)

### API Design
RESTful API with `/api/v1` prefix and consistent patterns:
- Authentication endpoints: `/api/v1/auth/*`
- Product endpoints: `/api/v1/products/*`
- Auction endpoints: `/api/v1/auctions/*`
- All endpoints follow REST conventions with proper HTTP verbs

### Authentication
- JWT tokens with refresh token rotation
- Role-based access control (buyer, seller, admin)
- Multi-factor authentication for sensitive operations
- Session management with device tracking

### Real-time Features
- LiveKit for video streaming (room creation, token generation)
- Gorilla WebSocket for chat and bidding
- Socket.IO client for frontend connections
- Redis Streams for async processing

### Backend Module Structure
Each module follows clean architecture principles:
- Models (data structures with GORM tags)
- Services (business logic interfaces)
- Handlers (HTTP request handlers)
- Repositories (data access interfaces)

## Security Requirements
- HTTPS everywhere (Cloudflare SSL)
- Rate limiting per user/IP
- Input validation with Zod schemas
- SQL injection prevention with parameterized queries
- XSS protection with proper escaping
- Password hashing with bcrypt
- Database encryption at rest and in transit
- Redis authentication and network isolation
- Container security with non-root users

## Performance Targets

### Response Time Targets
- API endpoints: <100ms (95th percentile)
- Database queries: <50ms average
- Cache hits: <10ms
- WebSocket messages: <20ms

### Concurrent User Targets
- Phase 1: 1,000 concurrent users
- Phase 2: 5,000 concurrent users  
- Phase 3: 10,000+ concurrent users

### Availability Targets
- Uptime: 99.9% (all services)
- Database: 99.95% uptime
- Redis: 99.99% uptime
- Live streaming: 99.5% uptime

## Implementation Phases

### Phase 1: Monolith Consolidation (0-3 months)
1. Merge existing services into unified backend
2. Implement single database with proper relations
3. Add comprehensive testing (unit + integration)
4. Deploy as single container for simplicity
5. Add monitoring and logging

### Phase 2: Performance Optimization (3-6 months)
1. Implement Redis caching layer
2. Add database indexing and query optimization
3. Implement CDN for static assets
4. Add background job processing for async tasks
5. Optimize mobile app performance

### Phase 3: Microservices Extraction (6-12 months, if needed)
Only split when business needs justify it:
- Team size > 15 engineers
- Different scaling requirements per domain
- Technology stack divergence needs
- Independent deployment requirements

#### Extraction Order:
1. Authentication Service (stateless)
2. Notification Service (I/O heavy)
3. File Upload Service (different infrastructure)
4. Analytics Service (ML/recommendations)

## Testing Strategy

### Backend
- Unit tests with Go's testing package
- Integration tests for API endpoints
- Database tests with test containers
- Target: 85%+ test coverage

### Frontend
- Jest + React Testing Library
- E2E tests with Cypress or Playwright
- Component testing with Storybook

### Mobile
- Jest for unit tests
- Detox for E2E testing

### Documentation Requirements
- GoDoc comments for all public functions
- TypeScript JSDoc for component props
- API documentation with OpenAPI/Swagger
- Database schema documentation
- README files for each major component

## Important Gotchas

1. **Single Database**: Avoid distributed monolith anti-pattern - use one unified database initially
2. **Future Extraction Points**: Design modules with clear boundaries for potential microservices extraction
3. **Real-time Complexity**: LiveKit integration requires careful handling of connections and room management
4. **Auction State Management**: Implement proper state transitions for auction lifecycle
5. **Bid Validation**: Prevent race conditions and ensure bid consistency
6. **WebSocket Scale**: Implement proper connection pooling for high-concurrency scenarios
7. **File Storage**: Use AWS S3 or equivalent for media files - don't store in database

## Environment Variables

Key environment variables to configure:
- `DATABASE_URL`: PostgreSQL connection string
- `REDIS_URL`: Redis connection string
- `JWT_SECRET`: JWT signing secret
- `LIVEKIT_API_KEY`: LiveKit API key
- `LIVEKIT_API_SECRET`: LiveKit API secret
- `NODE_ENV`: Development/production flag

## Infrastructure Stack

### Containerization
- Docker + Docker Compose
- Multi-stage builds for optimization
- Non-root user in containers

### Services
- PostgreSQL 17.7 (unified database)
- Redis 8+ (caching and sessions)
- Traefik (reverse proxy)
- Cloudflare (SSL/CDN)

### Monitoring & Logging
- Prometheus + Grafana (monitoring)
- ELK stack or Papertrail (logging)
- Health check endpoints for all services

## Code Conventions

### Go
- Follow standard Go formatting
- Use clean architecture patterns
- Interfaces in service layer
- Repository pattern for data access
- GORM tags for database models
- Proper error handling with wrapped errors

### TypeScript
- Strict TypeScript configuration
- Explicit return types
- Component props as interfaces
- Custom hooks for shared logic
- Zod schemas for validation

### Success Metrics

### Technical Metrics
- Code quality: 85%+ test coverage
- Performance: <100ms average response time
- Reliability: 99.9% uptime
- Security: Zero critical vulnerabilities

### Business Metrics
- User engagement: 10+ minutes average session
- Conversion rate: 5%+ auction participation / 3%+ purchase conversion
- Mobile adoption: 40%+ traffic from mobile
- Customer satisfaction: 4.5+ star rating
- Cart abandonment rate: <60% (industry average)
- Average order value: $75+ target
- Payment success rate: 95%+
- Order fulfillment time: <48 hours average

## Before You Start

1. Read the full architecture document (`BLYTZ_LIVE_ARCHITECTURE_PRD.md`)
2. Set up PostgreSQL 17.7 and Redis 8+ locally
3. Install Go 1.21+, Node.js 18+, and Docker
4. Consider creating a minimal proof-of-concept for the core auction functionality first
5. Plan for mobile-first responsive design