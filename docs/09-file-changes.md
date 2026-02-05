# File Structure Reference

## Project Root

```
blytz.app/
├── README.md                    # Project overview
├── AGENTS.md                    # AI agent guide
├── docs/                        # Documentation
├── docker-compose.yml           # Local development
├── .env.example                 # Environment template
├── .gitignore
│
├── frontend/                    # Next.js application
├── backend/                     # Go/Bun API
├── mobile/                      # React Native app
└── infrastructure/              # Terraform/k8s configs
```

## Frontend Structure (Next.js)

```
frontend/
├── app/                         # Next.js App Router
│   ├── (auth)/                 # Auth route group
│   │   ├── layout.tsx
│   │   ├── login/
│   │   ├── register/
│   │   └── forgot-password/
│   ├── (main)/                 # Main app group
│   │   ├── layout.tsx
│   │   ├── page.tsx            # Home
│   │   ├── auctions/
│   │   ├── streams/
│   │   ├── products/
│   │   ├── sellers/
│   │   ├── cart/
│   │   ├── checkout/
│   │   ├── orders/
│   │   └── profile/
│   ├── (seller)/               # Seller dashboard
│   │   ├── layout.tsx
│   │   ├── dashboard/
│   │   ├── products/
│   │   ├── auctions/
│   │   ├── streams/
│   │   └── orders/
│   ├── (admin)/                # Admin panel
│   ├── api/                    # API routes
│   │   └── webhooks/
│   ├── globals.css
│   └── layout.tsx
│
├── components/                  # React components
│   ├── ui/                     # shadcn/ui components
│   ├── layout/                 # Layout components
│   ├── auction/                # Auction components
│   ├── stream/                 # Stream components
│   ├── chat/                   # Chat components
│   ├── product/                # Product components
│   ├── seller/                 # Seller components
│   ├── order/                  # Order components
│   ├── payment/                # Payment components
│   └── shared/                 # Shared components
│
├── hooks/                       # Custom React hooks
│   ├── use-auth.ts
│   ├── use-auction.ts
│   ├── use-bid.ts
│   ├── use-stream.ts
│   ├── use-chat.ts
│   ├── use-product.ts
│   ├── use-cart.ts
│   └── use-notification.ts
│
├── lib/                         # Utilities
│   ├── api.ts                  # API client
│   ├── socket.ts               # Socket.io
│   ├── livekit.ts              # LiveKit client
│   ├── utils.ts                # Helpers
│   └── constants.ts
│
├── stores/                      # Zustand stores
│   ├── auth-store.ts
│   ├── cart-store.ts
│   ├── ui-store.ts
│   └── notification-store.ts
│
├── types/                       # TypeScript types
│   ├── index.ts
│   ├── user.ts
│   ├── product.ts
│   ├── auction.ts
│   ├── stream.ts
│   └── order.ts
│
├── public/                      # Static assets
│   ├── images/
│   └── fonts/
│
├── next.config.js
├── tailwind.config.ts
├── tsconfig.json
└── package.json
```

## Backend Structure (Go/Bun)

```
backend/
├── cmd/
│   └── server/
│       └── main.go             # Entry point
│
├── src/
│   ├── domain/                 # Domain layer
│   │   ├── user/
│   │   │   ├── entity.go
│   │   │   └── repository.go   # Interface
│   │   ├── product/
│   │   ├── auction/
│   │   ├── order/
│   │   ├── stream/
│   │   └── chat/
│   │
│   ├── application/            # Application layer
│   │   ├── auth/
│   │   │   └── service.go
│   │   ├── product/
│   │   ├── auction/
│   │   ├── order/
│   │   ├── payment/
│   │   ├── stream/
│   │   └── notification/
│   │
│   ├── infrastructure/         # Infrastructure layer
│   │   ├── persistence/
│   │   │   └── postgres/
│   │   │       ├── connection.go
│   │   │       ├── user_repository.go
│   │   │       ├── product_repository.go
│   │   │       └── ...
│   │   ├── cache/
│   │   │   └── redis/
│   │   │       ├── client.go
│   │   │       └── auction_cache.go
│   │   ├── http/
│   │   │   ├── server.go
│   │   │   ├── jwt.go
│   │   │   └── middleware/
│   │   ├── messaging/
│   │   │   └── redis/
│   │   │       └── event_bus.go
│   │   ├── websocket/
│   │   │   └── hub.go
│   │   └── external/
│   │       ├── stripe/
│   │       ├── livekit/
│   │       └── ninjavan/
│   │
│   ├── interfaces/             # Interface layer
│   │   ├── http/
│   │   │   ├── handlers/
│   │   │   │   ├── auth.go
│   │   │   │   ├── product.go
│   │   │   │   ├── auction.go
│   │   │   │   └── ...
│   │   │   └── middleware/
│   │   │       ├── auth.go
│   │   │       ├── cors.go
│   │   │       └── ratelimit.go
│   │   └── websocket/
│   │       └── handlers.go
│   │
│   ├── config/
│   │   └── config.go           # Environment config
│   │
│   └── app/
│       └── app.go              # DI container
│
├── pkg/                        # Public packages
│   └── errors/
│       └── errors.go
│
├── tests/                      # Test files
│   ├── unit/
│   ├── integration/
│   └── e2e/
│
├── drizzle/                    # Database migrations
│   ├── schema.ts
│   └── migrations/
│
├── go.mod
├── go.sum
└── Dockerfile
```

## Mobile Structure (React Native with Expo)

```
mobile/
├── app/                        # Expo Router app directory
│   ├── (auth)/                 # Auth group
│   │   ├── login.tsx
│   │   ├── register.tsx
│   │   └── _layout.tsx
│   ├── (tabs)/                 # Main tab navigation
│   │   ├── index.tsx           # Home/Feed
│   │   ├── auctions.tsx
│   │   ├── products.tsx
│   │   ├── streams.tsx
│   │   ├── profile.tsx
│   │   └── _layout.tsx
│   ├── auctions/
│   │   ├── [id].tsx            # Auction detail
│   │   └── _layout.tsx
│   ├── products/
│   │   ├── [slug].tsx          # Product detail
│   │   └── _layout.tsx
│   ├── streams/
│   │   ├── [id].tsx            # Stream viewer
│   │   └── _layout.tsx
│   ├── _layout.tsx             # Root layout
│   └── +html.tsx               # HTML wrapper
│
├── src/
│   ├── components/             # React components
│   │   ├── ui/                # Base UI components
│   │   ├── auction/           # Auction components
│   │   ├── product/           # Product components
│   │   ├── stream/            # Stream components
│   │   └── common/            # Shared components
│   │
│   ├── hooks/                  # Custom hooks
│   │   ├── use-auth.ts
│   │   ├── use-auction.ts
│   │   ├── use-product.ts
│   │   ├── use-stream.ts
│   │   └── use-api.ts
│   │
│   ├── stores/                 # Zustand stores
│   │   ├── auth-store.ts
│   │   ├── auction-store.ts
│   │   └── ui-store.ts
│   │
│   ├── services/               # API services
│   │   ├── api.ts             # Axios/fetch client
│   │   ├── auth.ts
│   │   ├── auction.ts
│   │   ├── product.ts
│   │   ├── socket.ts          # Socket.io
│   │   └── livekit.ts         # LiveKit SDK
│   │
│   ├── utils/                  # Utilities
│   │   ├── constants.ts
│   │   ├── helpers.ts
│   │   └── validation.ts
│   │
│   └── types/                  # TypeScript types
│       ├── index.ts
│       ├── user.ts
│       ├── product.ts
│       └── auction.ts
│
├── assets/                     # Static assets
│   ├── images/
│   ├── fonts/
│   └── icons/
│
├── components/                 # External components (if any)
├── constants/                  # App constants
├── hooks/                      # Global hooks
│
├── app.json                    # Expo config
├── babel.config.js
├── metro.config.js
├── tailwind.config.js          # NativeWind
├── tsconfig.json
└── package.json
```

## File Naming Conventions

### Frontend
```
Components:      PascalCase.tsx      (AuctionCard.tsx)
Hooks:           use-camelCase.ts    (use-auction.ts)
Utils:           camelCase.ts        (api.ts)
Types:           PascalCase.ts       (Auction.ts)
Styles:          camelCase.module.css (auction.module.css)
Pages:           page.tsx            (Next.js convention)
Layouts:         layout.tsx          (Next.js convention)
```

### Backend (Go)
```
Files:           snake_case.go       (user_repository.go)
Packages:        lowercase           (user, auction)
Interfaces:      PascalCase          (Repository, Service)
Structs:         PascalCase          (User, Auction)
Methods:         PascalCase (public) (GetByID)
                 camelCase (private) (validateEmail)
Constants:       PascalCase          (DefaultTimeout)
```

### Mobile (React Native/Dart)
```
Files:           snake_case.dart     (auction_card.dart)
Classes:         PascalCase          (AuctionCard)
Functions:       camelCase           (placeBid)
Variables:       camelCase           (currentBid)
Constants:       camelCase           (defaultTimeout)
Private:         _prefix             (_validateBid)
```

## Environment Files

### Frontend (.env.local)
```
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_SOCKET_URL=ws://localhost:8080
NEXT_PUBLIC_LIVEKIT_URL=wss://livekit.blytz.app
NEXT_PUBLIC_STRIPE_PUBLIC_KEY=pk_test_...
```

### Backend (.env)
```
# Server
PORT=8080
ENV=development

# Database
DATABASE_URL=postgresql://user:pass@localhost:5432/blytz

# Redis
REDIS_URL=redis://localhost:6379

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRY=1h

# Stripe
STRIPE_SECRET_KEY=sk_test_...
STRIPE_WEBHOOK_SECRET=whsec_...

# LiveKit
LIVEKIT_API_KEY=...
LIVEKIT_API_SECRET=...

# NinjaVan
NINJAVAN_API_KEY=...
NINJAVAN_API_SECRET=...

# R2
R2_ACCESS_KEY_ID=...
R2_SECRET_ACCESS_KEY=...
R2_BUCKET_NAME=...
```

### Mobile (.env)
```
API_URL=http://localhost:8080/api/v1
SOCKET_URL=ws://localhost:8080
LIVEKIT_URL=wss://livekit.blytz.app
STRIPE_PUBLIC_KEY=pk_test_...
```

## Key Files Reference

| Purpose | Frontend | Backend | Mobile |
|---------|----------|---------|--------|
| Entry | `app/layout.tsx` | `cmd/server/main.go` | `lib/main.dart` |
| Config | `next.config.js` | `src/config/config.go` | `pubspec.yaml` |
| Routes | `app/**/page.tsx` | `src/interfaces/http/handlers/*.go` | `lib/presentation/pages/**/*.dart` |
| API Client | `lib/api.ts` | N/A | `lib/services/api_service.dart` |
| Auth | `hooks/use-auth.ts` | `src/application/auth/service.go` | `lib/presentation/bloc/auth/` |
| State | `stores/*.ts` | N/A | `lib/presentation/providers/` |

---

*Last updated: 2025-02-05*
