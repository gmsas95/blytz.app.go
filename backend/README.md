# Blytz Backend (Clean Architecture)

## 🎯 Overview

This is the clean architecture rewrite of the Blytz.live backend.

**Key improvements:**
- ✅ Clean 40-line main.go (was 448 lines)
- ✅ Domain-driven design (no external dependencies in domain)
- ✅ Repository pattern (testable, swappable)
- ✅ Redis Pub/Sub for WebSocket scaling
- ✅ Structured error handling
- ✅ 100 DB connections (was 25)

## 📁 Structure

```
backend/
├── cmd/
│   └── server/main.go              # Entry point (40 lines)
├── internal/
│   ├── app/app.go                  # Application container (DI)
│   ├── domain/                      # Business logic (pure)
│   │   ├── auction/                # Auction, Bid, AutoBid entities
│   │   ├── user/                   # User, auth interfaces
│   │   ├── product/                # Product, Category entities
│   │   ├── order/                  # Order, Cart entities
│   │   └── payment/                # Payment, Gateway interfaces
│   ├── application/                 # Use cases (TODO)
│   │   ├── auth/                   # Auth service
│   │   ├── auction/                # Auction service
│   │   ├── catalog/                # Product service
│   │   └── order/                  # Order service
│   ├── infrastructure/              # External implementations
│   │   ├── persistence/postgres/   # GORM repositories
│   │   ├── cache/redis/            # Redis cache
│   │   ├── messaging/redis/        # Event bus (Pub/Sub)
│   │   ├── http/                   # HTTP server
│   │   └── websocket/              # WebSocket hub (TODO)
│   └── interfaces/                  # HTTP handlers, middleware
│       ├── http/handlers/          # Route handlers
│       └── middleware/             # Auth, rate limiting (TODO)
├── pkg/
│   └── errors/                     # Structured errors
└── deployments/                     # Docker, Swarm configs (TODO)
```

## 🚀 Quick Start

```bash
# 1. Install dependencies
cd backend
go mod init github.com/blytz/live/backend
go mod tidy

# 2. Set up environment
cp .env.example .env
# Edit .env with your DB and Redis credentials

# 3. Run
go run cmd/server/main.go
```

## 🏗️ Architecture Principles

1. **Domain Layer**: Pure business logic, no external dependencies
2. **Application Layer**: Use cases, orchestrates domain objects
3. **Infrastructure Layer**: External concerns (DB, HTTP, etc.)
4. **Interfaces Layer**: HTTP handlers, middleware

## 📊 Progress

| Component | Status |
|-----------|--------|
| Domain layer | ✅ Done |
| Infrastructure (DB, Redis) | ✅ Done |
| Application services | 🚧 In Progress |
| WebSocket (Redis Pub/Sub) | 🚧 In Progress |
| HTTP handlers | 🚧 In Progress |
| Rate limiting (Redis) | 🚧 In Progress |
| Docker Swarm config | 🚧 In Progress |
| Tests | 🚧 In Progress |

## 📝 Notes

- Old backend backed up to `../backend-old-backup.tar.gz`
- This is a work in progress - not production ready yet
- See `internal/domain/` for business rules
- See `internal/infrastructure/` for implementations
