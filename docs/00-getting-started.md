# Getting Started with Blytz Documentation

Welcome to the Blytz livestream ecommerce platform documentation! This guide will help you navigate our documentation and find what you need quickly.

## What is Blytz?

Blytz is a **livestream ecommerce and auction platform** that enables:
- 🎥 **Live streaming** - Sellers broadcast product demonstrations
- 💰 **Real-time auctions** - Buyers place bids during live streams
- 🛒 **Instant purchases** - Buy-now functionality alongside auctions
- 💬 **Live chat** - Real-time interaction between sellers and buyers
- 🚚 **Integrated logistics** - Seamless shipping with NinjaVan (Malaysia)

## Tech Stack Overview

```
┌─────────────────────────────────────────────────────────────┐
│  CLIENT LAYER                                                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │  Next.js    │  │  Flutter    │  │  LiveKit Client     │  │
│  │  (Web)      │  │  (Mobile)   │  │  (Streaming)        │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└──────────────────────────────┬──────────────────────────────┘
                               │
┌──────────────────────────────┼──────────────────────────────┐
│  API GATEWAY (Go/Bun)        │                              │
│  - REST APIs                 │                              │
│  - WebSocket (Socket.io)     │                              │
│  - Authentication (JWT)      │                              │
└──────────────────────────────┼──────────────────────────────┘
                               │
┌──────────────────────────────┼──────────────────────────────┐
│  SERVICES LAYER              │                              │
│  ┌─────────────┐  ┌─────────┴──────────┐  ┌──────────────┐  │
│  │  Auction    │  │  LiveKit Server    │  │  Chat        │  │
│  │  Service    │  │  (Video/Audio)     │  │  Service     │  │
│  └─────────────┘  └────────────────────┘  └──────────────┘  │
│  ┌─────────────┐  ┌────────────────────┐  ┌──────────────┐  │
│  │  Payment    │  │  Product Catalog   │  │  Order       │  │
│  │  (Stripe)   │  │  Service           │  │  Service     │  │
│  └─────────────┘  └────────────────────┘  └──────────────┘  │
└─────────────────────────────────────────────────────────────┘
                               │
┌──────────────────────────────┼──────────────────────────────┐
│  DATA LAYER                  │                              │
│  ┌─────────────┐  ┌─────────┴──────────┐  ┌──────────────┐  │
│  │  PostgreSQL │  │  Redis             │  │  Cloudflare  │  │
│  │  (Primary)  │  │  (Cache/Sessions)  │  │  R2 (Files)  │  │
│  └─────────────┘  └────────────────────┘  └──────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## Documentation Learning Paths

### Path 1: New Developer (Full Stack)
**Goal:** Understand the entire system

1. Start with [02-architecture.md](02-architecture.md) - System overview
2. Read [03-database-schema.md](03-database-schema.md) - Data model
3. Review [04-api-specifications.md](04-api-specifications.md) - API contracts
4. Check [09-file-changes.md](09-file-changes.md) - Project structure
5. Follow [08-implementation-phases.md](08-implementation-phases.md) - Development roadmap

### Path 2: Frontend Developer
**Goal:** Build the Next.js web app

1. [05-components.md](05-components.md) - Component specifications
2. [04-api-specifications.md](04-api-specifications.md) - API integration
3. [18-hooks-utilities.md](18-hooks-utilities.md) - Custom hooks
4. [13-ui-design-system.md](13-ui-design-system.md) - Styling guide
5. [frontend/architecture.md](frontend/architecture.md) - Frontend structure

### Path 3: Backend Developer
**Goal:** Build the Go/Bun API

1. [backend/architecture.md](backend/architecture.md) - Clean Architecture
2. [03-database-schema.md](03-database-schema.md) - Database design
3. [04-api-specifications.md](04-api-specifications.md) - API endpoints
4. [17-integration-guide.md](17-integration-guide.md) - Third-party APIs
5. [19-security-guidelines.md](19-security-guidelines.md) - Security practices

### Path 4: Mobile Developer
**Goal:** Build the Flutter app

1. [mobile/architecture.md](mobile/architecture.md) - Flutter structure
2. [04-api-specifications.md](04-api-specifications.md) - API integration
3. [07-navigation.md](07-navigation.md) - Screen flow
4. [integrations/livekit.md](integrations/livekit.md) - Streaming integration

### Path 5: DevOps/Deployment
**Goal:** Deploy and maintain the platform

1. [11-environment-config.md](11-environment-config.md) - Environment setup
2. [15-deployment-guide.md](15-deployment-guide.md) - Deployment procedures
3. [19-security-guidelines.md](19-security-guidelines.md) - Security checklist
4. [14-error-handling.md](14-error-handling.md) - Monitoring and alerts

## Documentation Tiers Explained

### Tier 1: Getting Started
- **00-getting-started.md** (this file)
- **Purpose:** Navigation guide, learning paths
- **Audience:** All team members

### Tier 2: Core Documentation
Documents that define WHAT we're building:

| Document | What It Defines |
|----------|-----------------|
| 01-requirements.md | Features, user stories, acceptance criteria |
| 02-architecture.md | System design, tech choices, data flow |
| 03-database-schema.md | Database entities, relationships, constraints |
| 04-api-specifications.md | API contracts, request/response formats |
| 05-components.md | UI components, props, state management |
| 06-permissions.md | User roles, permissions matrix |
| 07-navigation.md | Routes, navigation flow, deep linking |
| 08-implementation-phases.md | Development roadmap, milestones |
| 09-file-changes.md | File organization, naming conventions |
| 10-testing.md | Test strategies, coverage requirements |

### Tier 3: Operational Documentation
Documents that define HOW we build and operate:

| Document | What It Defines |
|----------|-----------------|
| 11-environment-config.md | Environment variables, secrets |
| 12-data-seeding.md | Initial data, sample users/products |
| 13-ui-design-system.md | Design tokens, Tailwind config |
| 14-error-handling.md | Error codes, logging strategy |
| 15-deployment-guide.md | CI/CD pipelines, rollback |
| 16-glossary.md | Business terms, translations |
| 17-integration-guide.md | Third-party API details |
| 18-hooks-utilities.md | Reusable code patterns |
| 19-security-guidelines.md | Security practices, PCI compliance |
| 20-accessibility.md | WCAG compliance, a11y testing |

## Key Concepts

### Livestream + Auction Flow
```
1. Seller starts live stream (LiveKit)
2. Buyers join stream and chat (Socket.io)
3. Seller showcases products
4. Buyers place bids in real-time
5. Auction timer counts down
6. Winner is determined
7. Payment processed (Stripe)
8. Order created, shipped (NinjaVan)
```

### User Roles
| Role | Description |
|------|-------------|
| **Buyer** | Browse, bid, purchase, watch streams |
| **Seller** | Create auctions, go live, manage inventory |
| **Admin** | Manage users, moderate content, analytics |

### Core Entities
- **User** - Buyers, sellers, admins
- **Stream** - Live broadcast session
- **Product** - Items for sale/auction
- **Auction** - Bidding session for a product
- **Bid** - Buyer's price offer
- **Order** - Completed purchase
- **Payment** - Stripe transaction
- **Shipment** - NinjaVan delivery

## Quick Reference

### Environment URLs
```
Development:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- LiveKit: ws://localhost:7880

Staging:
- Frontend: https://staging.blytz.app
- Backend API: https://api.staging.blytz.app

Production:
- Frontend: https://blytz.app
- Backend API: https://api.blytz.app
```

### Important Commands
```bash
# Backend
cd backend && bun run dev          # Start dev server
cd backend && bun test             # Run tests
cd backend && bun run migrate      # Run database migrations

# Frontend
cd frontend && npm run dev         # Start dev server
cd frontend && npm run build       # Production build
cd frontend && npm test            # Run tests

# Mobile
cd mobile && flutter run           # Run on device
cd mobile && flutter build apk     # Build Android
cd mobile && flutter build ios     # Build iOS
```

### Integration Keys (Required)
```bash
# Stripe
STRIPE_PUBLIC_KEY=pk_test_...
STRIPE_SECRET_KEY=sk_test_...
STRIPE_WEBHOOK_SECRET=whsec_...

# LiveKit
LIVEKIT_API_KEY=...
LIVEKIT_API_SECRET=...
LIVEKIT_WS_URL=wss://...

# NinjaVan (Malaysia)
NINJAVAN_API_KEY=...
NINJAVAN_API_SECRET=...

# Cloudflare R2
R2_ACCESS_KEY_ID=...
R2_SECRET_ACCESS_KEY=...
R2_BUCKET_NAME=...
R2_ENDPOINT=...
```

## Need Help?

- **Can't find something?** Check the [Documentation Index](README.md)
- **Architecture questions?** See [02-architecture.md](02-architecture.md)
- **API questions?** See [04-api-specifications.md](04-api-specifications.md)
- **Integration issues?** See [17-integration-guide.md](17-integration-guide.md)

---

**Ready to start?** Pick your learning path above or jump to [02-architecture.md](02-architecture.md) for the system overview.
