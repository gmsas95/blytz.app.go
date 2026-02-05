# Backend Architecture Documentation

## Overview

The Blytz.live backend follows **Clean Architecture** principles with a well-structured monolith design. This architecture provides:
- Clear separation of concerns
- Domain-driven design
- Testability at all layers
- Technology independence
- Future microservices extraction capability

## Clean Architecture Layers

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           CLEAN ARCHITECTURE                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                        INTERFACE LAYER                                │  │
│  │   (HTTP Handlers, Middleware, WebSocket, CLI)                        │  │
│  │                        ↑ Depends on Application                       │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                    ↑                                         │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                      APPLICATION LAYER                                │  │
│  │   (Services, Use Cases, DTOs, Event Handlers)                        │  │
│  │                        ↑ Depends on Domain                            │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                    ↑                                         │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                         DOMAIN LAYER                                  │  │
│  │   (Entities, Value Objects, Domain Services, Repository Interfaces)  │  │
│  │                        ↑ No external dependencies                     │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                    ↑                                         │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                      INFRASTRUCTURE LAYER                             │  │
│  │   (Database, Cache, Message Bus, External APIs)                      │  │
│  │                        ↓ Implements Domain Interfaces                 │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go                    # Application entry point
├── internal/
│   ├── app/                          # Application composition root
│   │   ├── app.go                    # Application container & dependency injection
│   │   └── stubs.go                  # Temporary stubs for missing implementations
│   ├── domain/                       # Domain layer (business entities)
│   │   ├── auction/                  # Auction domain
│   │   │   └── auction.go            # Auction entity, value objects
│   │   └── user/                     # User domain
│   │       └── user.go               # User entity, token manager interface
│   ├── application/                  # Application layer (use cases)
│   │   ├── auction/                  # Auction application services
│   │   │   └── service.go            # Auction use cases
│   │   └── auth/                     # Auth application services
│   │       └── service.go            # Authentication use cases
│   ├── infrastructure/               # Infrastructure layer (implementations)
│   │   ├── cache/redis/              # Redis cache implementation
│   │   │   ├── cache.go              # Redis client wrapper
│   │   │   └── auction.go            # Auction cache implementation
│   │   ├── http/                     # HTTP infrastructure
│   │   │   ├── jwt.go                # JWT token manager implementation
│   │   │   └── server.go             # HTTP server setup
│   │   ├── messaging/redis/          # Event bus implementation
│   │   │   └── event_bus.go          # Redis-based event bus
│   │   ├── persistence/postgres/     # Database repository implementations
│   │   │   ├── connection.go         # Database connection
│   │   │   ├── models.go             # GORM models
│   │   │   ├── auction_repository.go # Auction repository implementation
│   │   │   └── user_repository.go    # User repository implementation
│   │   └── websocket/                # WebSocket infrastructure
│   │       └── hub.go                # WebSocket hub for real-time communication
│   └── interfaces/                   # Interface adapters
│       ├── http/handlers/            # HTTP handlers
│       │   ├── auth.go               # Auth HTTP handlers
│       │   └── auction.go            # Auction HTTP handlers
│       └── middleware/               # HTTP middleware
│           ├── auth.go               # Authentication middleware
│           └── ratelimit.go          # Rate limiting middleware
├── pkg/                              # Public packages
│   └── errors/                       # Custom error types
├── go.mod                            # Go module definition
├── go.sum                            # Go dependencies
└── README.md                         # Backend README
```

## Layer Details

### 1. Domain Layer (`internal/domain/`)
**The heart of the application - no external dependencies.**

- **Entities**: Core business objects (User, Auction, Product)
- **Value Objects**: Immutable objects with no identity (Money, Address)
- **Domain Services**: Complex business logic that doesn't fit in entities
- **Repository Interfaces**: Define data access contracts
- **Domain Events**: Business events that occur within the domain

**Key Files:**
- `domain/auction/auction.go` - Auction entity with business rules
- `domain/user/user.go` - User entity with token manager interface

**Rules:**
- No imports from other layers
- Pure business logic
- Framework-agnostic

### 2. Application Layer (`internal/application/`)
**Orchestrates use cases and coordinates domain objects.**

- **Application Services**: Implement use cases
- **DTOs**: Data transfer objects for input/output
- **Event Handlers**: Handle domain events
- **Transaction Management**: Coordinate transactions across multiple aggregates

**Key Files:**
- `application/auction/service.go` - Auction use cases (create, bid, end)
- `application/auth/service.go` - Authentication use cases (login, register)

**Rules:**
- Depends only on Domain layer
- Contains no business rules (orchestrates domain)
- Framework-agnostic

### 3. Infrastructure Layer (`internal/infrastructure/`)
**Implements interfaces defined in Domain layer.**

- **Persistence**: Database repositories (PostgreSQL/GORM)
- **Cache**: Redis cache implementations
- **Messaging**: Event bus implementations (Redis Pub/Sub)
- **External APIs**: Third-party service integrations
- **WebSocket**: Real-time communication infrastructure

**Key Files:**
- `infrastructure/persistence/postgres/*` - Repository implementations
- `infrastructure/cache/redis/*` - Cache implementations
- `infrastructure/http/jwt.go` - JWT token manager implementation
- `infrastructure/websocket/hub.go` - WebSocket hub

**Rules:**
- Implements domain interfaces
- Contains no business logic
- Technology-specific code

### 4. Interface Layer (`internal/interfaces/`)
**Adapts external interfaces to application layer.**

- **HTTP Handlers**: REST API endpoints
- **Middleware**: Cross-cutting concerns (auth, rate limiting)
- **WebSocket Handlers**: Real-time connection handlers
- **CLI Commands**: Command-line interfaces

**Key Files:**
- `interfaces/http/handlers/auth.go` - Auth HTTP handlers
- `interfaces/http/handlers/auction.go` - Auction HTTP handlers
- `interfaces/middleware/auth.go` - Auth middleware

**Rules:**
- Depends on Application layer
- Handles HTTP-specific concerns only
- No business logic

### 5. Application Composition (`internal/app/`)
**Wires all dependencies together.**

- **Dependency Injection**: Creates and connects all components
- **Configuration**: Application configuration
- **Lifecycle Management**: Startup and shutdown

**Key Files:**
- `app/app.go` - Application container with DI

## Dependency Rule

**Dependencies point inward (Domain ← Application ← Infrastructure/Interface):**

```go
// Domain layer - no external dependencies
package domain

type AuctionRepository interface {
    GetByID(ctx context.Context, id uuid.UUID) (*Auction, error)
    Save(ctx context.Context, auction *Auction) error
}

// Application layer - depends on Domain
package application

type AuctionService struct {
    repo domain.AuctionRepository  // Interface from domain
}

// Infrastructure layer - implements Domain interface
package persistence

type AuctionRepository struct {
    db *gorm.DB
}

func (r *AuctionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Auction, error) {
    // Implementation using GORM
}
```

## Data Flow

### HTTP Request Flow
1. **HTTP Request** → Interface Layer (Handler)
2. **Handler** → Validates input → Calls Application Service
3. **Application Service** → Orchestrates → Calls Domain methods
4. **Domain** → Executes business logic
5. **Repository Interface** → Implemented by Infrastructure
6. **Infrastructure** → Persists to PostgreSQL / Redis
7. **Response** flows back through layers

### WebSocket Flow (Real-time Bidding)
1. **Client** → WebSocket connection
2. **WebSocket Hub** → (Infrastructure layer)
3. **Event Bus** → Publishes bid events
4. **Application Service** → Processes bid
5. **Domain** → Validates bid rules
6. **Broadcast** → All connected clients receive update

## Technology Stack

### Backend
- **Language**: Go 1.23+
- **Framework**: Gin (HTTP routing)
- **ORM**: GORM (database abstraction)
- **Database**: PostgreSQL 17.7
- **Cache**: Redis 8+
- **WebSocket**: Gorilla WebSocket
- **Events**: Redis Pub/Sub

### Domain Entities

```go
// domain/auction/auction.go
type Auction struct {
    ID          uuid.UUID
    SellerID    uuid.UUID
    ProductID   uuid.UUID
    Title       string
    Description string
    StartingBid decimal.Decimal
    ReservePrice decimal.Decimal
    Status      AuctionStatus
    StartTime   time.Time
    EndTime     time.Time
    CreatedAt   time.Time
    UpdatedAt   time.Time
    Bids        []Bid
}

type Bid struct {
    ID        uuid.UUID
    AuctionID uuid.UUID
    BidderID  uuid.UUID
    Amount    decimal.Decimal
    CreatedAt time.Time
}
```

## Design Principles

### 1. Clean Architecture
- **Dependency Inversion**: Dependencies point inward
- **Separation of Concerns**: Each layer has single responsibility
- **Testability**: Each layer can be tested in isolation
- **Independence**: Frameworks/UI can be swapped

### 2. Domain-Driven Design
- **Ubiquitous Language**: Code matches business terms
- **Bounded Contexts**: Clear domain boundaries
- **Rich Domain Models**: Entities contain behavior
- **Domain Events**: Loose coupling between aggregates

### 3. SOLID Principles
- **Single Responsibility**: One reason to change per class
- **Open/Closed**: Open for extension, closed for modification
- **Liskov Substitution**: Subtypes fully substitutable
- **Interface Segregation**: Small, focused interfaces
- **Dependency Inversion**: Depend on abstractions

## Security Architecture

### Authentication
- **JWT Tokens**: Stateless with access/refresh tokens
- **Token Manager Interface**: Defined in Domain layer
- **JWT Implementation**: In Infrastructure layer
- **Role-Based Access**: Buyer, seller, admin roles

### Input Validation
- **Handler Layer**: HTTP request validation
- **Domain Layer**: Business rule validation
- **Sanitization**: Prevent SQL injection via GORM
- **Rate Limiting**: Middleware in Interface layer

### Data Protection
- **Sensitive Fields**: Passwords never exposed
- **HTTPS**: Enforced in production
- **Environment Variables**: Secrets externalized

## Testing Strategy

### Unit Tests
```
domain/         # Test business logic without dependencies
application/    # Test use cases with mocked repositories
```

### Integration Tests
```
infrastructure/ # Test with real database/cache
interfaces/     # Test HTTP handlers with test server
```

### Test Isolation
```go
// Domain tests - no dependencies
func TestAuction_PlaceBid(t *testing.T) {
    auction := domain.NewAuction(...)
    err := auction.PlaceBid(bid)
    // Assert
}

// Application tests - mocked repositories
func TestAuctionService_PlaceBid(t *testing.T) {
    mockRepo := &MockAuctionRepository{}
    service := application.NewAuctionService(mockRepo, ...)
    // Test
}
```

## Migration from Old Structure

The codebase was migrated from a **module-based structure** to **Clean Architecture**:

| Old Structure | New Structure | Purpose |
|--------------|---------------|---------|
| `internal/auth/` | `domain/user/` + `application/auth/` | Separated domain from use cases |
| `internal/models/` | `internal/domain/*/` | Domain entities with behavior |
| `internal/database/` | `infrastructure/persistence/` | Repository implementations |
| `internal/cache/` | `infrastructure/cache/` | Cache implementations |
| `internal/middleware/` | `interfaces/middleware/` | HTTP-specific middleware |
| `internal/handlers/` | `interfaces/http/handlers/` | HTTP handlers |

## Benefits of Clean Architecture

1. **Testability**: Domain logic tested without database/HTTP
2. **Flexibility**: Swap PostgreSQL for MongoDB without touching domain
3. **Clarity**: Clear boundaries between business and technical concerns
4. **Maintainability**: Changes isolated to specific layers
5. **Future-Proof**: Easy to extract microservices later

## Current Implementation Status

✅ **Completed:**
- Domain layer (User, Auction entities)
- Application layer (Auth, Auction services)
- Infrastructure layer (Postgres, Redis, WebSocket)
- Interface layer (HTTP handlers, middleware)
- Dependency injection in app.go

📋 **TODO:**
- Product domain (was in old structure, needs migration)
- Cart domain
- Order domain
- Payment domain
- Address domain
- Complete auction WebSocket integration

## Development Guidelines

### Adding a New Feature

1. **Domain First**: Define entities and repository interfaces in `domain/`
2. **Application Second**: Implement use cases in `application/`
3. **Infrastructure Third**: Create repository implementations
4. **Interface Last**: Add HTTP handlers

### Example: Adding Product Catalog

```
1. domain/product/product.go          # Product entity, repository interface
2. application/product/service.go     # Product use cases
3. infrastructure/persistence/postgres/product_repository.go
4. interfaces/http/handlers/product.go # HTTP handlers
5. internal/app/app.go                # Wire dependencies
```

---

**This architecture provides a solid foundation for Blytz.live while maintaining flexibility for future evolution.**
