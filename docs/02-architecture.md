# System Architecture

## Overview

Blytz uses a modern microservices-oriented architecture with clean separation between frontend, backend, and infrastructure layers. The system is designed for scalability, real-time communication, and high availability.

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                                    CLIENT LAYER                                          │
├─────────────────────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────────┐  ┌─────────────────────┐  ┌─────────────────────┐               │
│  │   Next.js Web App   │  │  Flutter Mobile App │  │  LiveKit Client SDK │               │
│  │   (Port 3000)       │  │  (iOS/Android)      │  │  (Video/Audio)      │               │
│  │                     │  │                     │  │                     │               │
│  │ - Product browsing  │  │ - Native camera     │  │ - Stream publishing │               │
│  │ - Auction UI        │  │ - Push notifs       │  │ - Stream playback   │               │
│  │ - Chat interface    │  │ - Deep linking      │  │ - Screen share      │               │
│  └──────────┬──────────┘  └──────────┬──────────┘  └──────────┬──────────┘               │
│             │                        │                        │                          │
│             └────────────────────────┼────────────────────────┘                          │
│                                      │                                                    │
└──────────────────────────────────────┼────────────────────────────────────────────────────┘
                                       │
                              ┌────────┴────────┐
                              │   Cloudflare    │
                              │   (CDN + DNS)   │
                              │                 │
                              │ - Static assets │
                              │ - DDoS protection│
                              │ - Edge caching  │
                              └────────┬────────┘
                                       │
┌──────────────────────────────────────┼────────────────────────────────────────────────────┐
│                              API GATEWAY LAYER                                           │
├──────────────────────────────────────────────────────────────────────────────────────────┤
│                                      ▼                                                   │
│  ┌─────────────────────────────────────────────────────────────────────────────────────┐ │
│  │                        Go/Bun API Server (Port 8080)                                │ │
│  │                                                                                     │ │
│  │  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐               │ │
│  │  │  Auth Module │ │ Auction API  │ │ Product API  │ │ Order API    │               │ │
│  │  │  /api/v1/auth│ │  /api/v1/auctions│ │ /api/v1/products│ │ /api/v1/orders│          │ │
│  │  └──────────────┘ └──────────────┘ └──────────────┘ └──────────────┘               │ │
│  │  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐               │ │
│  │  │ Payment API  │ │ User API     │ │ Chat API     │ │ Webhook API  │               │ │
│  │  │ /api/v1/pay  │ │ /api/v1/users│ │ /api/v1/chat │ │ /webhooks/*  │               │ │
│  │  └──────────────┘ └──────────────┘ └──────────────┘ └──────────────┘               │ │
│  │                                                                                     │ │
│  │  Middleware: JWT Auth, Rate Limiting, CORS, Request Logging                        │ │
│  └─────────────────────────────────────────────────────────────────────────────────────┘ │
│                                      │                                                   │
│  ┌─────────────────────────────────────────────────────────────────────────────────────┐ │
│  │                    Socket.io Gateway (Real-time Communication)                      │ │
│  │                                                                                     │ │
│  │  Events:                                                                            │ │
│  │  - chat:message      - Bid placed notifications                                   │ │
│  │  - auction:bid       - User join/leave notifications                              │ │
│  │  - auction:update    - Stream status changes                                      │ │
│  │  - user:joined       - Order status updates                                       │ │
│  └─────────────────────────────────────────────────────────────────────────────────────┘ │
│                                      │                                                   │
└──────────────────────────────────────┼────────────────────────────────────────────────────┘
                                       │
┌──────────────────────────────────────┼────────────────────────────────────────────────────┐
│                         SERVICE LAYER (Internal Communication)                          │
├──────────────────────────────────────────────────────────────────────────────────────────┤
│                                      │                                                   │
│  ┌───────────────────────────────────┼─────────────────────────────────────────────────┐ │
│  │                                   │                                                 │ │
│  │  ┌──────────────┐ ┌──────────────┐┴┌──────────────┐ ┌──────────────┐               │ │
│  │  │  Auction     │ │  Product     │ │  Payment     │ │  Order       │               │ │
│  │  │  Service     │ │  Service     │ │  Service     │ │  Service     │               │ │
│  │  │              │ │              │ │              │ │              │               │ │
│  │  │ - Bid logic  │ │ - Catalog    │ │ - Stripe     │ │ - Checkout   │               │ │
│  │  │ - Timer mgmt │ │ - Inventory  │ │ - FPX        │ │ - NinjaVan   │               │ │
│  │  │ - Winner det │ │ - Search     │ │ - Refunds    │ │ - Tracking   │               │ │
│  │  └──────────────┘ └──────────────┘ └──────────────┘ └──────────────┘               │ │
│  │                                                                                     │ │
│  │  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐               │ │
│  │  │  User        │ │  Notification│ │  Stream      │ │  Chat        │               │ │
│  │  │  Service     │ │  Service     │ │  Service     │ │  Service     │               │ │
│  │  │              │ │              │ │              │ │              │               │ │
│  │  │ - Profiles   │ │ - Push       │ │ - LiveKit    │ │ - History    │               │ │
│  │  │ - Followers  │ │ - Email      │ │ - Room mgmt  │ │ - Moderation │               │ │
│  │  │ - Auth       │ │ - SMS        │ │ - Recording  │ │ - Socket.io  │               │ │
│  │  └──────────────┘ └──────────────┘ └──────────────┘ └──────────────┘               │ │
│  │                                                                                     │ │
│  └─────────────────────────────────────────────────────────────────────────────────────┘ │
│                                      │                                                   │
└──────────────────────────────────────┼────────────────────────────────────────────────────┘
                                       │
┌──────────────────────────────────────┼────────────────────────────────────────────────────┐
│                              DATA LAYER                                                  │
├──────────────────────────────────────────────────────────────────────────────────────────┤
│                                      │                                                   │
│  ┌───────────────────────────────────┼─────────────────────────────────────────────────┐ │
│  │                                   │                                                 │ │
│  │  ┌───────────────────────────────────────────────────────────────────────────────┐  │ │
│  │  │                         PostgreSQL (Primary Database)                         │  │ │
│  │  │                                                                               │  │ │
│  │  │   Tables:                                                                     │  │ │
│  │  │   - users, sellers, buyers                                                    │  │ │
│  │  │   - products, categories, inventory                                           │  │ │
│  │  │   - auctions, bids                                                            │  │ │
│  │  │   - orders, order_items, payments                                             │  │ │
│  │  │   - streams, chat_messages                                                    │  │ │
│  │  │   - shipments, tracking                                                       │  │ │
│  │  │                                                                               │  │ │
│  │  │   Features:                                                                   │  │ │
│  │  │   - Connection pooling                                                        │  │ │
│  │  │   - Read replicas (future)                                                    │  │ │
│  │  │   - Automated backups                                                         │  │ │
│  │  └───────────────────────────────────────────────────────────────────────────────┘  │ │
│  │                                   │                                                 │ │
│  │  ┌───────────────────────────────────────────────────────────────────────────────┐  │ │
│  │  │                         Redis (Cache & Sessions)                              │  │ │
│  │  │                                                                               │  │ │
│  │  │   Use Cases:                                                                  │  │ │
│  │  │   - Session store (JWT blacklist)                                             │  │ │
│  │  │   - Auction state cache (current bid, timer)                                  │  │ │
│  │  │   - Rate limiting counters                                                    │  │ │
│  │  │   - Real-time leaderboards                                                    │  │ │
│  │  │   - Pub/Sub for events                                                        │  │ │
│  │  │   - Query result caching                                                      │  │ │
│  │  └───────────────────────────────────────────────────────────────────────────────┘  │ │
│  │                                   │                                                 │ │
│  │  ┌───────────────────────────────────────────────────────────────────────────────┐  │ │
│  │  │                      Cloudflare R2 (Object Storage)                           │  │ │
│  │  │                                                                               │  │ │
│  │  │   Storage:                                                                    │  │ │
│  │  │   - Product images (multiple sizes)                                           │  │ │
│  │  │   - Stream recordings                                                         │  │ │
│  │  │   - User avatars                                                              │  │ │
│  │  │   - Invoice PDFs                                                              │  │ │
│  │  └───────────────────────────────────────────────────────────────────────────────┘  │ │
│  │                                                                                     │ │
│  └─────────────────────────────────────────────────────────────────────────────────────┘ │
│                                      │                                                   │
└──────────────────────────────────────┼────────────────────────────────────────────────────┘
                                       │
┌──────────────────────────────────────┼────────────────────────────────────────────────────┐
│                         EXTERNAL INTEGRATIONS                                            │
├──────────────────────────────────────┼────────────────────────────────────────────────────┤
│                                      ▼                                                   │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐                  │
│  │   Stripe     │  │   LiveKit    │  │   NinjaVan   │  │   Socket.io  │                  │
│  │  (Payments)  │  │ (Streaming)  │  │  (Shipping)  │  │  (Real-time) │                  │
│  │              │  │              │  │              │  │              │                  │
│  │ - Cards      │  │ - Video      │  │ - Labels     │  │ - Chat       │                  │
│  │ - FPX        │  │ - Audio      │  │ - Tracking   │  │ - Events     │                  │
│  │ - Wallets    │  │ - Recording  │  │ - Malaysia   │  │ - Presence   │                  │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘                  │
│                                                                                          │
└──────────────────────────────────────────────────────────────────────────────────────────┘
```

## Component Details

### Frontend (Next.js 15+)

**Framework:** Next.js with App Router
**Language:** TypeScript
**Styling:** Tailwind CSS + shadcn/ui
**State Management:** Zustand + React Query

**Key Features:**
- Server-side rendering for SEO
- Real-time updates via Socket.io client
- LiveKit client SDK for streaming
- Responsive mobile-first design

**Directory Structure:**
```
frontend/
├── app/                    # Next.js App Router
│   ├── (auth)/            # Auth group (login, register)
│   ├── (main)/            # Main app group
│   │   ├── auctions/      # Auction pages
│   │   ├── products/      # Product catalog
│   │   ├── streams/       # Live streams
│   │   └── profile/       # User profiles
│   ├── api/               # API routes (webhooks)
│   └── layout.tsx         # Root layout
├── components/            # React components
│   ├── ui/               # shadcn components
│   ├── auction/          # Auction components
│   ├── stream/           # Streaming components
│   └── chat/             # Chat components
├── lib/                   # Utilities
├── hooks/                 # Custom hooks
└── stores/                # Zustand stores
```

### Backend (Go/Bun)

**Runtime:** Bun (Node.js compatible)
**Framework:** Elysia.js or Express
**Database:** PostgreSQL (via Drizzle ORM)
**Architecture:** Clean Architecture / Hexagonal

**Key Modules:**
1. **Auth Service** - JWT, refresh tokens, social login
2. **Auction Service** - Bid processing, timer management
3. **Product Service** - Catalog, inventory, search
4. **Order Service** - Checkout, order lifecycle
5. **Payment Service** - Stripe integration, webhooks
6. **Stream Service** - LiveKit room management
7. **Chat Service** - Message persistence, moderation
8. **Notification Service** - Push, email, SMS

**Directory Structure:**
```
backend/
├── src/
│   ├── domain/           # Entities, repository interfaces
│   ├── application/      # Use cases, services
│   ├── infrastructure/   # DB, cache, external APIs
│   ├── interfaces/       # HTTP handlers, middleware
│   └── config/           # Environment config
├── tests/                # Test suites
└── drizzle/              # Database migrations
```

### Mobile (Flutter)

**Framework:** Flutter 3+
**Language:** Dart
**State Management:** Riverpod
**Architecture:** Clean Architecture

**Key Features:**
- Native camera integration for streaming
- Push notifications (Firebase)
- Deep linking
- Offline support

**Directory Structure:**
```
mobile/
├── lib/
│   ├── core/             # Utils, constants
│   ├── data/             # Repositories, models
│   ├── domain/           # Entities, use cases
│   ├── presentation/     # Screens, widgets
│   │   ├── screens/
│   │   └── widgets/
│   └── main.dart
├── android/              # Android config
├── ios/                  # iOS config
└── test/                 # Tests
```

### Database (PostgreSQL)

**Version:** PostgreSQL 17+
**ORM:** Drizzle (TypeScript) or GORM (Go)
**Migrations:** Version controlled, automated

**Key Tables:**
- `users` - User accounts
- `sellers` - Seller profiles
- `products` - Product catalog
- `auctions` - Auction sessions
- `bids` - Bid history
- `orders` - Customer orders
- `payments` - Payment records
- `streams` - Live stream sessions
- `shipments` - Shipping records

### Cache (Redis)

**Version:** Redis 8+
**Use Cases:**
- Session management
- Auction state (current bid, timer)
- Rate limiting
- Real-time leaderboards
- Pub/Sub for events

### Storage (Cloudflare R2)

**Type:** S3-compatible object storage
**Content:**
- Product images (original + thumbnails)
- Stream recordings
- User avatars
- Invoice PDFs

**Optimization:**
- Image resizing on upload
- CDN delivery via Cloudflare
- Signed URLs for private content

## Integration Details

### Stripe (Payments)
**Integration Type:** REST API + Webhooks
**Features:**
- Credit/debit card payments
- FPX (Malaysia bank transfer)
- E-wallets (GrabPay, Touch 'n Go)
- Subscription billing (for sellers)
- Refunds and disputes

**Webhooks:**
- `payment_intent.succeeded`
- `payment_intent.payment_failed`
- `charge.dispute.created`

### LiveKit (Streaming)
**Integration Type:** WebSocket + REST API
**Features:**
- WebRTC video/audio streaming
- Screen sharing
- Recording
- Multi-host support

**Components:**
- LiveKit Server (self-hosted or cloud)
- LiveKit Client SDK (JS, Flutter)

### NinjaVan (Logistics)
**Integration Type:** REST API
**Features:**
- Shipment creation
- Label generation
- Real-time tracking
- Malaysia-only coverage

**API Endpoints:**
- Create order
- Get tracking status
- Print shipping label
- Cancel shipment

### Socket.io (Real-time)
**Integration Type:** WebSocket
**Features:**
- Real-time chat
- Bid notifications
- Presence indicators
- Room-based communication

**Events:**
- `join:stream` - Join stream room
- `chat:message` - Send chat message
- `auction:bid` - Place bid
- `user:joined` - User joined notification

## Data Flow Examples

### Auction Bid Flow
```
1. Buyer clicks "Place Bid" in Next.js app
2. Frontend sends POST /api/v1/auctions/:id/bids
3. Backend validates JWT, checks bid amount
4. Auction Service processes bid logic
5. Bid saved to PostgreSQL
6. Auction state updated in Redis
7. Socket.io broadcasts bid to all viewers
8. Frontend updates UI with new bid
```

### Live Stream Flow
```
1. Seller clicks "Go Live" in Flutter app
2. LiveKit client creates room
3. Stream record created in PostgreSQL
4. Buyers receive push notification
5. Buyers join via Next.js or Flutter
6. LiveKit establishes WebRTC connections
7. Chat messages via Socket.io
8. Stream recording saved to R2
```

### Order Fulfillment Flow
```
1. Auction ends, winner determined
2. Order automatically created
3. Payment captured via Stripe
4. NinjaVan shipment created
5. Shipping label generated
6. Seller prints label, ships item
7. Tracking updates via NinjaVan webhook
8. Buyer receives delivery confirmation
```

## Scalability Considerations

### Horizontal Scaling
- Stateless API servers (Bun/Go)
- Database read replicas
- Redis Cluster for cache
- CDN for static assets

### Performance Optimization
- Database query optimization
- Redis caching layer
- Image optimization (WebP, responsive)
- Lazy loading for streams

### Future Architecture
- Event-driven with message queue (Kafka/NATS)
- Microservices extraction (when needed)
- Edge computing for low-latency streams
- Multi-region deployment

---

*Last updated: 2025-02-05*
