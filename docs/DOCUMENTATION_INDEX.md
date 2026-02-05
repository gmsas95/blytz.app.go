# Documentation Index

Welcome to the Blytz.app.go documentation hub. This index provides quick access to all project documentation.

## 📋 Quick Links

### Getting Started
- **[Project README](../README.md)** - Project overview and quick start guide
- **[AGENTS Guide](../AGENTS.md)** - AI/agent development guide with Clean Architecture
- **[Architecture PRD](../BLYTZ_LIVE_ARCHITECTURE_PRD.md)** - Complete architecture specification

### Backend Development
- **[Backend Architecture](backend/architecture.md)** - Clean Architecture design and patterns
- **[Development Guide](backend/development-guide.md)** - Setup, coding standards, and workflows
- **[API Reference](api/backend-api.md)** - RESTful API documentation

### Project Documentation
- **[E-commerce Guide](ECOMMERCE_IMPLEMENTATION_GUIDE.md)** - E-commerce implementation plan

## 📊 Current Development Status

### ✅ Completed
1. **Clean Architecture Foundation** - Domain-driven design with proper layering
2. **Authentication System** - JWT auth, user management, security features
3. **Infrastructure Layer** - PostgreSQL, Redis, WebSocket setup

### 🔄 In Progress
4. **Product & Auction System** - Domain entities, WebSocket bidding

### 📋 Planned
5. **E-commerce System** - Shopping cart, orders, payments, addresses
6. **Live Streaming** - LiveKit video streaming integration

## 🛠️ Development Resources

### Backend Quick Start
```bash
cd backend
go mod tidy
go run cmd/server/main.go

# Server: http://localhost:8080
# Health: http://localhost:8080/health
```

### Authentication Test
```bash
# Register user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","first_name":"Test"}'

# Login user
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

## 📚 Documentation Categories

### 1. Architecture & Design
- [Architecture PRD](../BLYTZ_LIVE_ARCHITECTURE_PRD.md) - Complete system architecture
- [Backend Architecture](backend/architecture.md) - Clean Architecture details
- [AGENTS Guide](../AGENTS.md) - Development guide for AI agents

### 2. Development Guides
- [Backend Development Guide](backend/development-guide.md) - Backend coding standards
- [E-commerce Implementation Guide](ECOMMERCE_IMPLEMENTATION_GUIDE.md) - E-commerce development plan

### 3. API Documentation
- [Backend API Reference](api/backend-api.md) - RESTful API endpoints

## 🔍 Search Documentation

### Setting up development environment?
→ [Backend Development Guide](backend/development-guide.md)

### Understanding the system architecture?
→ [Architecture PRD](../BLYTZ_LIVE_ARCHITECTURE_PRD.md)
→ [Backend Architecture](backend/architecture.md)

### API endpoint documentation?
→ [Backend API Reference](api/backend-api.md)

### Adding a new domain/feature?
→ [AGENTS Guide](../AGENTS.md) - See "Adding a New Domain" section

### Current development progress?
→ [AGENTS Guide](../AGENTS.md) - See "Current Development Status"

### E-commerce implementation plan?
→ [E-commerce Implementation Guide](ECOMMERCE_IMPLEMENTATION_GUIDE.md)

## 🗂️ File Structure

```
docs/
├── README.md                    # This documentation index
├── DOCUMENTATION_INDEX.md       # Quick navigation (this file)
├── ECOMMERCE_IMPLEMENTATION_GUIDE.md  # E-commerce plan
├── backend/                     # Backend documentation
│   ├── architecture.md          # Clean Architecture design
│   └── development-guide.md    # Development setup & standards
├── api/                        # API documentation
│   └── backend-api.md          # RESTful API reference
├── frontend/                   # Frontend documentation (planned)
│   └── architecture.md
└── mobile/                     # Mobile documentation (planned)
    └── architecture.md

Root Level:
├── README.md                    # Project overview & quick start
├── AGENTS.md                   # AI agent development guide
├── BLYTZ_LIVE_ARCHITECTURE_PRD.md  # Complete architecture spec
├── ARCHITECTURE_IDEAL.md       # Ideal future architecture
└── backend/                     # Backend source code
    └── internal/               # Clean Architecture layers
        ├── app/                # Dependency injection
        ├── domain/             # Business entities
        ├── application/        # Use cases
        ├── infrastructure/     # Implementations
        └── interfaces/         # HTTP handlers
```

## 🏗️ Clean Architecture Overview

The backend follows Clean Architecture with these layers:

```
┌─────────────────────────────────────────┐
│  Interface Layer                        │
│  - HTTP Handlers                        │
│  - Middleware                           │
│  - WebSocket                            │
├─────────────────────────────────────────┤
│  Application Layer                      │
│  - Services                             │
│  - Use Cases                            │
│  - DTOs                                 │
├─────────────────────────────────────────┤
│  Domain Layer                           │
│  - Entities                             │
│  - Value Objects                        │
│  - Repository Interfaces                │
├─────────────────────────────────────────┤
│  Infrastructure Layer                   │
│  - PostgreSQL Repositories              │
│  - Redis Cache                          │
│  - WebSocket Hub                        │
│  - External APIs                        │
└─────────────────────────────────────────┘
```

Dependencies point **inward** - Domain has no external dependencies.

## 📈 Documentation Roadmap

### Backend
- [x] Clean Architecture documentation
- [x] Development setup guide
- [ ] Performance optimization guide
- [ ] Security best practices
- [ ] Testing strategies
- [ ] Database schema documentation
- [ ] Deployment guide

### Frontend
- [ ] Component library documentation
- [ ] State management patterns
- [ ] API integration guide

### Mobile
- [ ] Platform-specific features
- [ ] API integration guide

### API
- [ ] Interactive API documentation (Swagger)
- [ ] Postman collections
- [ ] WebSocket documentation

## 🔗 External Resources

### Technology Documentation
- **Go**: https://golang.org/doc/
- **Gin Framework**: https://gin-gonic.com/docs/
- **GORM**: https://gorm.io/docs/
- **PostgreSQL**: https://www.postgresql.org/docs/
- **Redis**: https://redis.io/documentation
- **Clean Architecture**: https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html

### Development Tools
- **Docker**: https://docs.docker.com/
- **GitHub CLI**: https://cli.github.com/manual/

## 📞 Support & Feedback

### Documentation Issues
- Found outdated information? Create an issue
- Missing information? Request additional documentation
- Confusing explanations? Suggest improvements

### Contribution Guidelines
1. Follow Clean Architecture principles
2. Update related documentation when making changes
3. Include code examples where helpful
4. Keep domain logic framework-agnostic

---

**This documentation is actively maintained. Last updated: 2025-02-05**
