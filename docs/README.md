# Blytz Documentation Hub

Welcome to the Blytz livestream ecommerce platform documentation.

## Quick Navigation

| Tier | Section | Purpose |
|------|---------|---------|
| **Tier 1** | [Getting Started](00-getting-started.md) | Main entry point, learning path |
| **Tier 2** | [Core Documentation](#tier-2-core-documentation) | Architecture, APIs, components |
| **Tier 3** | [Operations](#tier-3-operational-documentation) | Deployment, config, security |

## Documentation Index

### Tier 1: Getting Started
- [00-getting-started.md](00-getting-started.md) - How to navigate docs, learning paths

### Tier 2: Core Documentation

| # | Document | Description |
|---|----------|-------------|
| 01 | [Requirements](01-requirements.md) | User requirements, acceptance criteria, user stories |
| 02 | [System Architecture](02-architecture.md) | Overall system design, tech stack overview |
| 03 | [Database Schema](03-database-schema.md) | PostgreSQL schema, entities, relationships |
| 04 | [API Specifications](04-api-specifications.md) | RESTful APIs, WebSocket events |
| 05 | [Frontend Components](05-components.md) | Next.js components, pages, state management |
| 06 | [Permissions & RBAC](06-permissions.md) | Role-based access control (buyer, seller, admin) |
| 07 | [Navigation Structure](07-navigation.md) | App routes, menu structure, navigation flow |
| 08 | [Implementation Phases](08-implementation-phases.md) | Phase-by-phase development roadmap |
| 09 | [File Structure](09-file-changes.md) | Project file organization reference |
| 10 | [Testing Strategy](10-testing.md) | Unit, integration, E2E testing |

### Tier 3: Operational Documentation

| # | Document | Description |
|---|----------|-------------|
| 11 | [Environment Configuration](11-environment-config.md) | Environment variables, secrets management |
| 12 | [Data Seeding](12-data-seeding.md) | Sample data, initial setup scripts |
| 13 | [UI Design System](13-ui-design-system.md) | Colors, typography, Tailwind config |
| 14 | [Error Handling](14-error-handling.md) | Error codes, boundaries, logging |
| 15 | [Deployment Guide](15-deployment-guide.md) | CI/CD, deployment, rollback procedures |
| 16 | [Glossary](16-glossary.md) | Business terms, Malay translations |
| 17 | [Integration Guide](17-integration-guide.md) | Third-party integrations (Stripe, LiveKit, etc.) |
| 18 | [Hooks & Utilities](18-hooks-utilities.md) | Custom React hooks, Go utilities |
| 19 | [Security Guidelines](19-security-guidelines.md) | Security best practices, PCI compliance |
| 20 | [Accessibility](20-accessibility.md) | WCAG compliance, keyboard navigation |

## Platform-Specific Docs

### Backend (Go/Bun)
- [Backend Architecture](backend/architecture.md) - Clean Architecture details
- [Backend Folder Structure](backend/folder-structure.md) - Code organization

### Frontend (Next.js)
- [Frontend Architecture](frontend/architecture.md) - Next.js App Router structure
- [Component Library](frontend/components.md) - Reusable UI components

### Mobile (Flutter)
- [Mobile Architecture](mobile/architecture.md) - Flutter project structure
- [Mobile Navigation](mobile/navigation.md) - Routing and navigation

### Integrations
- [Stripe Payments](integrations/stripe.md) - Payment processing
- [LiveKit Streaming](integrations/livekit.md) - Livestream engine
- [Socket.io Chat](integrations/socketio.md) - Real-time chat
- [NinjaVan Logistics](integrations/ninjavan.md) - Malaysia shipping

## Tech Stack Summary

| Layer | Technology | Purpose |
|-------|------------|---------|
| Frontend | Next.js 15+ | React web application |
| Backend | Go 1.23+ / Bun | API server, business logic |
| Mobile | Flutter 3+ | iOS/Android apps |
| Database | PostgreSQL 17+ | Primary data storage |
| Cache | Redis 8+ | Sessions, real-time data |
| Storage | Cloudflare R2 | Image/video storage |
| Payments | Stripe | Payment processing |
| Streaming | LiveKit | Live video streaming |
| Real-time | Socket.io | Chat, notifications |
| Logistics | NinjaVan | Malaysia shipping |

## Development Quick Start

```bash
# Backend
cd backend
bun install
bun run dev

# Frontend
cd frontend
npm install
npm run dev

# Mobile
cd mobile
flutter pub get
flutter run
```

---

*Last updated: 2025-02-05*
