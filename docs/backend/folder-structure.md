# Backend Folder Structure Reference

This document provides a complete reference of the backend folder structure following Clean Architecture principles.

## Root Structure

```
backend/
├── cmd/                          # Application entry points
│   └── server/
│       └── main.go              # Main server entry point
├── internal/                     # Private application code
│   ├── app/                     # Application composition root
│   ├── domain/                  # Domain layer (business logic)
│   ├── application/             # Application layer (use cases)
│   ├── infrastructure/          # Infrastructure layer (implementations)
│   └── interfaces/              # Interface adapters (HTTP, CLI)
├── pkg/                         # Public packages (importable by others)
├── go.mod                       # Go module definition
├── go.sum                       # Go dependencies
└── README.md                    # Backend README
```

## Layer-by-Layer Breakdown

### 1. CMD Layer (`cmd/`)

**Purpose:** Application entry points

```
cmd/
└── server/
    └── main.go                  # Server entry point
        # - Loads configuration
        # - Creates application
        # - Starts server
        # - Handles graceful shutdown
```

**Key Rule:** Keep main.go thin - just wire up dependencies and start the app.

### 2. APP Layer (`internal/app/`)

**Purpose:** Dependency injection and application lifecycle

```
app/
├── app.go                       # Application container
│   # - Config struct
│   # - New() - creates app with all dependencies
│   # - Run() - starts HTTP server + WebSocket hub
│   # - Shutdown() - graceful shutdown
└── stubs.go                     # Temporary stubs
    # - Placeholder implementations for missing pieces
```

**Key Rule:** This is the only place where all layers meet.

### 3. DOMAIN Layer (`internal/domain/`)

**Purpose:** Business entities and rules (no external dependencies)

```
domain/
├── auction/                     # Auction domain
│   └── auction.go
│       # - Auction entity
│       # - Bid entity
│       # - AuctionStatus enum
│       # - Repository interface
│       # - Domain methods (PlaceBid, etc.)
└── user/                        # User domain
    └── user.go
        # - User entity
        # - Role enum
        # - TokenManager interface
        # - Domain methods
```

**Key Rule:** Domain has ZERO external dependencies. No imports from other layers.

**Example Domain Entity:**
```go
// domain/auction/auction.go
package auction

type Auction struct {
    ID          uuid.UUID
    SellerID    uuid.UUID
    Title       string
    Status      AuctionStatus
    Bids        []Bid
}

type Repository interface {
    GetByID(ctx context.Context, id uuid.UUID) (*Auction, error)
    Save(ctx context.Context, auction *Auction) error
}

// Domain method with business logic
func (a *Auction) PlaceBid(bid Bid) error {
    if a.Status != StatusActive {
        return errors.New("auction not active")
    }
    // ... more business rules
    a.Bids = append(a.Bids, bid)
    return nil
}
```

### 4. APPLICATION Layer (`internal/application/`)

**Purpose:** Use cases and application services

```
application/
├── auction/                     # Auction use cases
│   └── service.go
│       # - CreateAuction()
│       # - GetAuction()
│       # - PlaceBid()
│       # - EndAuction()
└── auth/                        # Auth use cases
    └── service.go
        # - Register()
        # - Login()
        # - RefreshToken()
        # - Logout()
```

**Key Rule:** Orchestrates domain objects. No business logic here.

**Example Application Service:**
```go
// application/auction/service.go
package auction

type Service struct {
    repo    domain.AuctionRepository
    cache   Cache
    eventBus EventBus
}

func (s *Service) PlaceBid(ctx context.Context, dto PlaceBidDTO) error {
    // 1. Get auction from repository
    auction, err := s.repo.GetByID(ctx, dto.AuctionID)
    if err != nil {
        return err
    }
    
    // 2. Create bid
    bid := domain.Bid{
        BidderID: dto.BidderID,
        Amount:   dto.Amount,
    }
    
    // 3. Execute domain logic
    if err := auction.PlaceBid(bid); err != nil {
        return err
    }
    
    // 4. Save result
    if err := s.repo.Save(ctx, auction); err != nil {
        return err
    }
    
    // 5. Publish event
    s.eventBus.Publish("bid.placed", bid)
    
    return nil
}
```

### 5. INFRASTRUCTURE Layer (`internal/infrastructure/`)

**Purpose:** Implement interfaces defined in Domain layer

```
infrastructure/
├── cache/
│   └── redis/                   # Redis cache implementations
│       ├── cache.go            # Redis client wrapper
│       └── auction.go          # Auction cache implementation
├── http/
│   ├── jwt.go                  # JWT token manager implementation
│   └── server.go               # HTTP server setup
├── messaging/
│   └── redis/                  # Event bus implementation
│       └── event_bus.go
├── persistence/
│   └── postgres/               # Repository implementations
│       ├── connection.go       # Database connection
│       ├── models.go           # GORM models
│       ├── auction_repository.go
│       └── user_repository.go
└── websocket/
    └── hub.go                  # WebSocket hub
```

**Key Rule:** Implements domain interfaces. Technology-specific code.

**Example Repository Implementation:**
```go
// infrastructure/persistence/postgres/auction_repository.go
package postgres

type AuctionRepository struct {
    db *gorm.DB
}

func NewAuctionRepository(db *gorm.DB) *AuctionRepository {
    return &AuctionRepository{db: db}
}

func (r *AuctionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Auction, error) {
    var model AuctionModel
    if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return model.ToDomain(), nil
}

func (r *AuctionRepository) Save(ctx context.Context, auction *domain.Auction) error {
    model := FromDomain(auction)
    return r.db.WithContext(ctx).Save(model).Error
}
```

### 6. INTERFACE Layer (`internal/interfaces/`)

**Purpose:** Adapt external interfaces to application layer

```
interfaces/
├── http/
│   └── handlers/               # HTTP handlers
│       ├── auth.go            # Auth endpoints
│       └── auction.go         # Auction endpoints
└── middleware/                # HTTP middleware
    ├── auth.go                # JWT validation
    └── ratelimit.go           # Rate limiting
```

**Key Rule:** Handles HTTP concerns only. No business logic.

**Example HTTP Handler:**
```go
// interfaces/http/handlers/auction.go
package handlers

type AuctionHandler struct {
    service *auction.Service
}

func NewAuctionHandler(service *auction.Service) *AuctionHandler {
    return &AuctionHandler{service: service}
}

func (h *AuctionHandler) PlaceBid(c *gin.Context) {
    // 1. Parse request
    var req PlaceBidRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, ErrorResponse{Error: err.Error()})
        return
    }
    
    // 2. Get user from context (set by auth middleware)
    userID := c.GetString("user_id")
    
    // 3. Call application service
    dto := auction.PlaceBidDTO{
        AuctionID: uuid.MustParse(c.Param("id")),
        BidderID:  uuid.MustParse(userID),
        Amount:    req.Amount,
    }
    
    if err := h.service.PlaceBid(c.Request.Context(), dto); err != nil {
        c.JSON(500, ErrorResponse{Error: err.Error()})
        return
    }
    
    // 4. Return response
    c.JSON(200, SuccessResponse{Message: "Bid placed"})
}
```

### 7. PKG Layer (`pkg/`)

**Purpose:** Public packages that can be imported by other projects

```
pkg/
└── errors/                     # Custom error types
    └── errors.go
```

**Key Rule:** No internal dependencies. Self-contained utilities.

## Dependency Flow

```
┌─────────────────────────────────────────────────────────────┐
│  EXTERNAL WORLD                                              │
│  (HTTP, Database, Cache, External APIs)                     │
└─────────────┬───────────────────────────────────────────────┘
              │
┌─────────────▼───────────────────────────────────────────────┐
│  INTERFACE LAYER (internal/interfaces/)                     │
│  - HTTP Handlers                                            │
│  - Middleware                                               │
│  - WebSocket Handlers                                       │
└─────────────┬───────────────────────────────────────────────┘
              │ calls
┌─────────────▼───────────────────────────────────────────────┐
│  APPLICATION LAYER (internal/application/)                  │
│  - Use Cases                                                │
│  - Services                                                 │
│  - DTOs                                                     │
└─────────────┬───────────────────────────────────────────────┘
              │ orchestrates
┌─────────────▼───────────────────────────────────────────────┐
│  DOMAIN LAYER (internal/domain/)                            │
│  - Entities                                                 │
│  - Value Objects                                            │
│  - Repository Interfaces                                    │
└─────────────┬───────────────────────────────────────────────┘
              │ implements
┌─────────────▼───────────────────────────────────────────────┐
│  INFRASTRUCTURE LAYER (internal/infrastructure/)            │
│  - PostgreSQL Repositories                                  │
│  - Redis Cache                                              │
│  - WebSocket Hub                                            │
└─────────────────────────────────────────────────────────────┘
```

## Adding a New Feature: Checklist

When adding a new domain (e.g., Product):

- [ ] **Domain Layer**
  - [ ] Create `domain/product/product.go` with entity and repository interface
  - [ ] Define value objects (if any)
  - [ ] Write domain method tests

- [ ] **Application Layer**
  - [ ] Create `application/product/service.go`
  - [ ] Define DTOs for input/output
  - [ ] Implement use cases
  - [ ] Write service tests with mocked repository

- [ ] **Infrastructure Layer**
  - [ ] Create `infrastructure/persistence/postgres/product_repository.go`
  - [ ] Implement GORM model
  - [ ] Implement repository interface methods
  - [ ] Add cache implementation (if needed)

- [ ] **Interface Layer**
  - [ ] Create `interfaces/http/handlers/product.go`
  - [ ] Define request/response structs
  - [ ] Implement HTTP handlers
  - [ ] Add routes in server setup

- [ ] **App Layer**
  - [ ] Add product service initialization in `app/app.go`
  - [ ] Add product repository to DI container
  - [ ] Wire up HTTP handlers

## Common Patterns

### Repository Pattern
```go
// Domain: Define interface
type Repository interface {
    GetByID(ctx context.Context, id uuid.UUID) (*Entity, error)
    Save(ctx context.Context, entity *Entity) error
    Delete(ctx context.Context, id uuid.UUID) error
}

// Infrastructure: Implement interface
type Repository struct {
    db *gorm.DB
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Entity, error) {
    // Implementation
}
```

### Service Pattern
```go
// Application: Use repository interface
type Service struct {
    repo domain.Repository
}

func (s *Service) DoSomething(ctx context.Context, dto DTO) error {
    entity, err := s.repo.GetByID(ctx, dto.ID)
    if err != nil {
        return err
    }
    // Business logic
    return s.repo.Save(ctx, entity)
}
```

### Handler Pattern
```go
// Interface: Handle HTTP
type Handler struct {
    service *application.Service
}

func (h *Handler) Endpoint(c *gin.Context) {
    var req Request
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, Error(err))
        return
    }
    
    result, err := h.service.DoSomething(c.Request.Context(), req.ToDTO())
    if err != nil {
        c.JSON(500, Error(err))
        return
    }
    
    c.JSON(200, Success(result))
}
```

## File Naming Conventions

| Layer | File Name Pattern | Example |
|-------|------------------|---------|
| Domain | `entity.go` | `auction.go` |
| Application | `service.go` | `service.go` |
| Infrastructure | `entity_repository.go` | `auction_repository.go` |
| Interface | `entity.go` | `auction.go` |

## Import Rules

**Allowed:**
```go
// Application can import Domain
import "github.com/blytz/live/backend/internal/domain/auction"

// Infrastructure can import Domain
import "github.com/blytz/live/backend/internal/domain/auction"

// Interface can import Application and Domain
import "github.com/blytz/live/backend/internal/application/auction"
import "github.com/blytz/live/backend/internal/domain/auction"
```

**Forbidden:**
```go
// Domain CANNOT import other layers
import "github.com/blytz/live/backend/internal/application/..." // ❌
import "github.com/blytz/live/backend/internal/infrastructure/..." // ❌

// Application CANNOT import Infrastructure
import "github.com/blytz/live/backend/internal/infrastructure/..." // ❌
```

---

**Remember: Dependencies always point inward toward the Domain!**
