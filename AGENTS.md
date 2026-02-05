# Blytz.app.go - Agents Guide

## Project Overview

Blytz is a **livestream ecommerce and auction platform** that enables:
- 🎥 **Live streaming** - Sellers broadcast product demonstrations
- 💰 **Real-time auctions** - Buyers place bids during live streams
- 🛒 **Instant purchases** - Buy-now functionality alongside auctions
- 💬 **Live chat** - Real-time interaction between sellers and buyers
- 🚚 **Integrated logistics** - Seamless shipping with NinjaVan (Malaysia)

### Business Objectives
- **Primary**: Create engaging live commerce experiences with real-time auctions
- **Secondary**: Build a scalable marketplace supporting high concurrent users
- **Tertiary**: Enable mobile-first experiences for on-the-go bidding

## Tech Stack

| Layer | Technology | Purpose |
|-------|------------|---------|
| **Frontend** | Next.js 15+ | React web application |
| **Backend** | Go 1.23+ / Bun | API server, business logic |
| **Mobile** | Flutter 3+ | iOS/Android apps |
| **Database** | PostgreSQL 17+ | Primary data storage |
| **Cache** | Redis 8+ | Sessions, real-time data |
| **Storage** | Cloudflare R2 | Image/video storage |
| **Payments** | Stripe | Payment processing (Cards, FPX, E-wallets) |
| **Streaming** | LiveKit | Live video streaming |
| **Real-time** | Socket.io | Chat, bid notifications |
| **Logistics** | NinjaVan | Malaysia shipping |

## Documentation Index

> 📚 **All documentation is in `/docs/` folder. Start with [00-getting-started.md](docs/00-getting-started.md)**

### 📖 Quick Navigation

| Learning Path | Start Here |
|---------------|------------|
| **New Developer** | [00-getting-started.md](docs/00-getting-started.md) → [02-architecture.md](docs/02-architecture.md) |
| **Frontend Dev** | [05-components.md](docs/05-components.md) → [04-api-specifications.md](docs/04-api-specifications.md) |
| **Backend Dev** | [backend/architecture.md](docs/backend/architecture.md) → [03-database-schema.md](docs/03-database-schema.md) |
| **Mobile Dev** | [mobile/architecture.md](docs/mobile/architecture.md) → [04-api-specifications.md](docs/04-api-specifications.md) |
| **DevOps** | [11-environment-config.md](docs/11-environment-config.md) → [15-deployment-guide.md](docs/15-deployment-guide.md) |

### 📋 Tier 1: Getting Started

| Document | Description |
|----------|-------------|
| [README.md](docs/README.md) | Documentation hub index |
| [00-getting-started.md](docs/00-getting-started.md) | Main entry point, learning paths, quick reference |

### 📋 Tier 2: Core Documentation

| # | Document | Description | Status |
|---|----------|-------------|--------|
| 01 | [01-requirements.md](docs/01-requirements.md) | User requirements, acceptance criteria, user stories | ✅ Complete |
| 02 | [02-architecture.md](docs/02-architecture.md) | System architecture, tech stack, data flow diagrams | ✅ Complete |
| 03 | [03-database-schema.md](docs/03-database-schema.md) | PostgreSQL schema, entities, relationships | ✅ Complete |
| 04 | [04-api-specifications.md](docs/04-api-specifications.md) | RESTful APIs, WebSocket events, error codes | ✅ Complete |
| 05 | [05-components.md](docs/05-components.md) | Next.js components, pages, state management | ✅ Complete |
| 06 | [06-permissions.md](docs/06-permissions.md) | RBAC system (buyer, seller, admin) | ✅ Complete |
| 07 | [07-navigation.md](docs/07-navigation.md) | App routes, menu structure, deep linking | ✅ Complete |
| 08 | [08-implementation-phases.md](docs/08-implementation-phases.md) | Phase-by-phase development roadmap | ✅ Complete |
| 09 | [09-file-changes.md](docs/09-file-changes.md) | Project file organization reference | ✅ Complete |
| 10 | [10-testing.md](docs/10-testing.md) | Unit, integration, E2E testing strategy | ✅ Complete |

### 📋 Tier 3: Operational Documentation

| # | Document | Description | Status |
|---|----------|-------------|--------|
| 11 | [11-environment-config.md](docs/11-environment-config.md) | Environment variables, secrets management | ✅ Complete |
| 12 | ~~12-data-seeding.md~~ | Sample data, initial setup scripts | ⏳ *Create when implementing seeding* |
| 13 | ~~13-ui-design-system.md~~ | Colors, typography, Tailwind config | ⏳ *Create when building UI* |
| 14 | ~~14-error-handling.md~~ | Error codes, boundaries, logging | ⏳ *Create when implementing error handling* |
| 15 | [15-deployment-guide.md](docs/15-deployment-guide.md) | CI/CD, Kubernetes, deployment procedures | ✅ Complete |
| 16 | [16-glossary.md](docs/16-glossary.md) | Business terms, Malay translations | ✅ Complete |
| 17 | [17-integration-guide.md](docs/17-integration-guide.md) | Stripe, LiveKit, Socket.io, NinjaVan | ✅ Complete |
| 18 | ~~18-hooks-utilities.md~~ | Custom React hooks, Go utilities | ⏳ *Create when building features* |
| 19 | ~~19-security-guidelines.md~~ | Security best practices, PCI compliance | ⏳ *Create before security audit* |
| 20 | ~~20-accessibility.md~~ | WCAG compliance, keyboard navigation | ⏳ *Create during a11y pass* |

### 📁 Platform-Specific Docs

| Platform | Path |
|----------|------|
| **Backend** | [backend/architecture.md](docs/backend/architecture.md) - Clean Architecture details |
| **Backend** | [backend/folder-structure.md](docs/backend/folder-structure.md) - Code organization |
| **Frontend** | [frontend/architecture.md](docs/frontend/architecture.md) - Next.js structure |
| **Mobile** | [mobile/architecture.md](docs/mobile/architecture.md) - Flutter structure |

### 🔌 Integration Docs

| Integration | Path |
|-------------|------|
| **Stripe** | [17-integration-guide.md](docs/17-integration-guide.md) (section) |
| **LiveKit** | [17-integration-guide.md](docs/17-integration-guide.md) (section) |
| **Socket.io** | [17-integration-guide.md](docs/17-integration-guide.md) (section) |
| **NinjaVan** | [17-integration-guide.md](docs/17-integration-guide.md) (section) |
| **Cloudflare R2** | [17-integration-guide.md](docs/17-integration-guide.md) (section) |

## Current Development Status

### ✅ COMPLETED: Clean Architecture Foundation
Migrated from module-based to Clean Architecture:
- Domain layer with entities (User, Auction)
- Application layer with services
- Infrastructure layer with PostgreSQL, Redis, WebSocket
- Interface layer with HTTP handlers
- Dependency injection in app.go

### ✅ COMPLETED: Documentation Structure
Complete documentation framework:
- Tier 1: Getting Started
- Tier 2: Core Documentation (10 docs)
- Tier 3: Operational Documentation (in progress)

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

## Architecture Overview

The backend follows **Clean Architecture** principles:

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
GET  /api/v1/auth/me
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

## Running the Application

```bash
# Backend
cd backend
go mod tidy
go run cmd/server/main.go

# Server: http://localhost:8080
# Health: http://localhost:8080/health
```

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

## Integration Configuration

### Stripe
```bash
STRIPE_SECRET_KEY=sk_test_...
STRIPE_PUBLISHABLE_KEY=pk_test_...
STRIPE_WEBHOOK_SECRET=whsec_...
```

### LiveKit
```bash
LIVEKIT_API_KEY=...
LIVEKIT_API_SECRET=...
LIVEKIT_WS_URL=wss://...
```

### NinjaVan
```bash
NINJAVAN_API_KEY=...
NINJAVAN_API_SECRET=...
```

### Cloudflare R2
```bash
R2_ACCESS_KEY_ID=...
R2_SECRET_ACCESS_KEY=...
R2_BUCKET_NAME=...
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

## Documentation Quick Links

- **[Getting Started](docs/00-getting-started.md)** - Start here!
- **[Architecture](docs/02-architecture.md)** - System overview
- **[Database Schema](docs/03-database-schema.md)** - Entity relationships
- **[API Specs](docs/04-api-specifications.md)** - RESTful endpoints
- **[Deployment Guide](docs/15-deployment-guide.md)** - Production setup
- **[Glossary](docs/16-glossary.md)** - Business terms & Malay translations

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

---

**Full documentation available in `/docs/` folder**
