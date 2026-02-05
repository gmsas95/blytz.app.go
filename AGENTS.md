# Blytz.app.go - Agents Guide

## Project Overview

Blytz.app.go is a modern live marketplace platform designed for real-time auctions, bidding, and live streaming capabilities. The platform connects buyers and sellers through interactive live sessions with real-time product demonstrations and instant bidding functionality.

### Business Objectives
- **Primary**: Create engaging live commerce experiences with real-time auctions
- **Secondary**: Build a scalable marketplace supporting high concurrent users
- **Tertiary**: Enable mobile-first experiences for on-the-go bidding

## Architecture Overview

The backend follows **Clean Architecture** principles with a well-structured monolith design:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  Interface Layer    → HTTP Handlers, Middleware, WebSocket                │
│       ↑                                                                     │
│  Application Layer  → Use Cases, Services, DTOs                           │
│       ↑                                                                     │
│  Domain Layer       → Entities, Value Objects, Repository Interfaces      │
│       ↑                                                                     │
│  Infrastructure     → PostgreSQL, Redis, WebSocket, External APIs         │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Project Structure
```
backend/
├── cmd/server/main.go                 # Application entry point
├── internal/
│   ├── app/                          # Dependency injection & app container
│   ├── domain/                       # Business entities (no deps)
│   │   ├── auction/                  # Auction entity
│   │   └── user/                     # User entity
│   ├── application/                  # Use cases
│   │   ├── auction/                  # Auction service
│   │   └── auth/                     # Auth service
│   ├── infrastructure/               # Implementations
│   │   ├── cache/redis/              # Redis cache
│   │   ├── http/                     # JWT, HTTP server
│   │   ├── messaging/redis/          # Event bus
│   │   ├── persistence/postgres/     # Repositories
│   │   └── websocket/                # WebSocket hub
│   └── interfaces/                   # Adapters
│       ├── http/handlers/            # HTTP handlers
│       └── middleware/               # Auth, rate limiting
└── pkg/errors/                       # Custom errors
```

## Current Development Status

### ✅ COMPLETED: Clean Architecture Foundation (2025-02-04)
Migrated from module-based to Clean Architecture:
- Domain layer with entities (User, Auction)
- Application layer with services
- Infrastructure layer with PostgreSQL, Redis, WebSocket
- Interface layer with HTTP handlers
- Dependency injection in app.go

### ✅ COMPLETED: Authentication System
- JWT token management (access + refresh tokens)
- User registration and login
- Password hashing with bcrypt
- Protected routes middleware
- Role-based access control (buyer/seller/admin)

### 🔄 IN PROGRESS: Product & Auction System
**Current State:**
- ✅ Auction domain entity
- ✅ Auction application service
- ✅ Auction repository (Postgres)
- ✅ WebSocket hub infrastructure
- ✅ Event bus (Redis Pub/Sub)

**TODO:**
- ❌ Product domain (needs migration from old codebase)
- ❌ WebSocket bidding endpoints
- ❌ Auction lifecycle management (start/pause/end)
- ❌ Real-time bid broadcasting
- ❌ Winner determination

### 📋 PLANNED: E-commerce System
After Auction system completion:
- Shopping cart (Cart, CartItem)
- Order management (Order, OrderItem)
- Payment processing (Stripe integration)
- Address management
- Inventory management

### 📋 PLANNED: Live Streaming
- LiveKit integration
- Video streaming
- Live chat during auctions

## Domain Entities

### User
```go
// domain/user/user.go
type User struct {
    ID           uuid.UUID
    Email        string
    PasswordHash string
    Role         Role // buyer, seller, admin
    FirstName    string
    LastName     string
    // ...
}
```

### Auction
```go
// domain/auction/auction.go
type Auction struct {
    ID           uuid.UUID
    SellerID     uuid.UUID
    Title        string
    Description  string
    StartingBid  decimal.Decimal
    Status       AuctionStatus // draft, active, ended, cancelled
    StartTime    time.Time
    EndTime      time.Time
    Bids         []Bid
}

type Bid struct {
    ID        uuid.UUID
    BidderID  uuid.UUID
    Amount    decimal.Decimal
    CreatedAt time.Time
}
```

## API Endpoints

### Authentication
```
POST /api/v1/auth/register
POST /api/v1/auth/login
POST /api/v1/auth/refresh
GET  /api/v1/auth/profile
POST /api/v1/auth/logout
```

### Auctions
```
GET    /api/v1/auctions          # List auctions
POST   /api/v1/auctions          # Create auction
GET    /api/v1/auctions/:id      # Get auction details
PUT    /api/v1/auctions/:id      # Update auction
DELETE /api/v1/auctions/:id      # Delete auction
POST   /api/v1/auctions/:id/bid  # Place bid
GET    /api/v1/auctions/:id/bids # Get bid history
```

### WebSocket
```
WS /ws/auctions/:id              # Real-time auction updates
```

## Development Guidelines

### Adding a New Domain

**Step 1: Create Domain Entity**
```go
// internal/domain/product/product.go
package product

type Product struct {
    ID          uuid.UUID
    SellerID    uuid.UUID
    Title       string
    Description string
    Price       decimal.Decimal
    // ...
}

type Repository interface {
    GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
    Save(ctx context.Context, product *Product) error
    // ...
}
```

**Step 2: Create Application Service**
```go
// internal/application/product/service.go
package product

type Service struct {
    repo domain.ProductRepository
}

func (s *Service) CreateProduct(ctx context.Context, dto CreateProductDTO) (*domain.Product, error) {
    // Business logic
}
```

**Step 3: Create Infrastructure Repository**
```go
// internal/infrastructure/persistence/postgres/product_repository.go
package postgres

type ProductRepository struct {
    db *gorm.DB
}

func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
    // GORM implementation
}
```

**Step 4: Create HTTP Handler**
```go
// internal/interfaces/http/handlers/product.go
package handlers

type ProductHandler struct {
    service *product.Service
}

func (h *ProductHandler) Create(c *gin.Context) {
    // HTTP handling
}
```

**Step 5: Wire in app.go**
```go
// internal/app/app.go
func (a *Application) initServices() error {
    // ...
    productRepo := postgres.NewProductRepository(a.db)
    a.productService = product.NewService(productRepo)
    // ...
}
```

### Dependency Rule
**Dependencies always point inward:**
- Domain ← Application ← Infrastructure
- Domain ← Application ← Interface

**Never:**
- Domain imports Application/Infrastructure
- Application imports Infrastructure

## Testing

### Domain Tests (No Dependencies)
```go
func TestAuction_PlaceBid(t *testing.T) {
    auction := domain.NewAuction(...)
    bid := domain.Bid{Amount: decimal.NewFromFloat(100)}
    
    err := auction.PlaceBid(bid)
    
    assert.NoError(t, err)
    assert.Equal(t, 1, len(auction.Bids))
}
```

### Application Tests (Mocked Repositories)
```go
func TestAuctionService_PlaceBid(t *testing.T) {
    mockRepo := &MockAuctionRepository{}
    service := application.NewAuctionService(mockRepo, ...)
    
    // Test use case
}
```

## Technology Stack

- **Backend**: Go 1.23+, Gin framework
- **Database**: PostgreSQL 17.7 with GORM
- **Cache**: Redis 8+
- **WebSocket**: Gorilla WebSocket
- **Events**: Redis Pub/Sub
- **Auth**: JWT tokens

## Environment Variables

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=blytz

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# Auth
JWT_SECRET=your-secret-key

# Server
PORT=8080
ENV=development
```

## Running the Application

```bash
cd backend
go mod tidy
go run cmd/server/main.go

# Server: http://localhost:8080
# Health: http://localhost:8080/health
```

## Key Files Reference

| Purpose | Path |
|---------|------|
| Entry Point | `cmd/server/main.go` |
| DI Container | `internal/app/app.go` |
| User Entity | `internal/domain/user/user.go` |
| Auction Entity | `internal/domain/auction/auction.go` |
| Auth Service | `internal/application/auth/service.go` |
| Auction Service | `internal/application/auction/service.go` |
| User Repository | `internal/infrastructure/persistence/postgres/user_repository.go` |
| Auction Repository | `internal/infrastructure/persistence/postgres/auction_repository.go` |
| Auth Handler | `internal/interfaces/http/handlers/auth.go` |
| Auction Handler | `internal/interfaces/http/handlers/auction.go` |
| Auth Middleware | `internal/interfaces/middleware/auth.go` |
| WebSocket Hub | `internal/infrastructure/websocket/hub.go` |

## Documentation

- **[Architecture PRD](BLYTZ_LIVE_ARCHITECTURE_PRD.md)** - Complete system architecture
- **[Backend Architecture](docs/backend/architecture.md)** - Clean Architecture details
- **[Documentation Index](docs/DOCUMENTATION_INDEX.md)** - All documentation links

## Important Notes

1. **Clean Architecture**: Always follow the dependency rule - domain has no external dependencies
2. **Migrations Lost**: Product domain was in old codebase, needs to be re-implemented in Clean Architecture
3. **WebSocket Ready**: Infrastructure is set up, need to implement bidding logic
4. **Test Coverage**: Focus on domain layer tests first (no dependencies)
5. **Database**: Uses GORM with PostgreSQL, migrations auto-run on startup

## Next Priorities

1. **Product Domain** - Essential for auctions and e-commerce
2. **Auction WebSocket** - Real-time bidding functionality
3. **Cart/Orders** - E-commerce foundation
4. **LiveKit Integration** - Video streaming
