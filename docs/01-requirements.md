# Requirements Documentation

## Overview

Blytz is a livestream ecommerce platform combining real-time video streaming with auction and buy-now capabilities, specifically designed for the Malaysian market with NinjaVan logistics integration.

## User Personas

### 1. Buyer (Pembeli)
**Demographics:** 18-45 years old, mobile-first users
**Goals:**
- Discover unique products through live demonstrations
- Get deals through auctions
- Interact with sellers in real-time
- Track orders easily

**Pain Points:**
- Can't verify product quality online
- Miss interactive shopping experience
- Uncertain about shipping reliability

### 2. Seller (Penjual)
**Demographics:** Small business owners, entrepreneurs, collectors
**Goals:**
- Showcase products through live video
- Engage directly with customers
- Run auctions to maximize prices
- Manage inventory and orders

**Pain Points:**
- High commission fees on other platforms
- Limited interaction with buyers
- Complex logistics management

### 3. Admin (Pentadbir)
**Goals:**
- Platform moderation
- User management
- Analytics and reporting
- Ensure compliance

## Feature Requirements

### FR-001: User Authentication
**Priority:** High
**Description:** Users can register, login, and manage their accounts

**Acceptance Criteria:**
- [ ] Email/password registration with verification
- [ ] Social login (Google, Facebook)
- [ ] JWT token-based authentication
- [ ] Password reset functionality
- [ ] Profile management (name, avatar, address)
- [ ] Role selection: Buyer or Seller

**User Stories:**
- As a buyer, I want to create an account so I can participate in auctions
- As a seller, I want to verify my identity so buyers trust me

---

### FR-002: Live Streaming
**Priority:** High
**Description:** Sellers can broadcast live video streams

**Acceptance Criteria:**
- [ ] Start/stop live stream
- [ ] Multi-host support (co-streaming)
- [ ] Screen sharing capability
- [ ] Stream quality adaptation (720p, 1080p)
- [ ] Stream recording for replay
- [ ] Viewer count display
- [ ] Stream scheduling (go live at specific time)

**User Stories:**
- As a seller, I want to go live so buyers can see my products in real-time
- As a buyer, I want to watch live streams so I can evaluate products

---

### FR-003: Real-Time Chat
**Priority:** High
**Description:** Live chat during streams

**Acceptance Criteria:**
- [ ] Send/receive messages in real-time
- [ ] Emojis and reactions
- [ ] Moderator controls (delete messages, timeout users)
- [ ] Mention users (@username)
- [ ] Chat history persistence
- [ ] Pin important messages

**User Stories:**
- As a buyer, I want to ask questions in chat so I can learn about products
- As a seller, I want to see chat messages so I can respond to questions

---

### FR-004: Product Catalog
**Priority:** High
**Description:** Sellers can manage their products

**Acceptance Criteria:**
- [ ] Create product listings with images/videos
- [ ] Product categories and tags
- [ ] Inventory management (stock quantity)
- [ ] Product variants (size, color)
- [ ] Product description with rich text
- [ ] Product condition (new, used, refurbished)

**User Stories:**
- As a seller, I want to list products so buyers can find them
- As a buyer, I want to browse products so I can find what I need

---

### FR-005: Auction System
**Priority:** High
**Description:** Real-time bidding on products

**Acceptance Criteria:**
- [ ] Create auction with starting price, reserve price, buy-now price
- [ ] Set auction duration (countdown timer)
- [ ] Real-time bid placement
- [ ] Bid history and current highest bid display
- [ ] Auto-bid functionality (proxy bidding)
- [ ] Reserve price enforcement
- [ ] Automatic winner determination
- [ ] Bid notifications (outbid alerts)

**User Stories:**
- As a buyer, I want to place bids so I can win auctions
- As a seller, I want to set reserve prices so I don't sell below minimum

---

### FR-006: Shopping Cart
**Priority:** Medium
**Description:** Buyers can add items to cart

**Acceptance Criteria:**
- [ ] Add products to cart
- [ ] Update quantities
- [ ] Remove items
- [ ] Save for later
- [ ] Cart persistence across sessions
- [ ] Show total price with shipping estimate

**User Stories:**
- As a buyer, I want to add items to cart so I can purchase multiple items

---

### FR-007: Checkout & Payments
**Priority:** High
**Description:** Secure payment processing

**Acceptance Criteria:**
- [ ] Multiple payment methods (credit card, FPX, e-wallets)
- [ ] Stripe integration for card payments
- [ ] Shipping address selection
- [ ] Order summary before payment
- [ ] Payment confirmation and receipt
- [ ] Failed payment handling
- [ ] Refund processing

**User Stories:**
- As a buyer, I want to pay securely so I can complete my purchase
- As a seller, I want to receive payments so I can fulfill orders

---

### FR-008: Order Management
**Priority:** High
**Description:** Track and manage orders

**Acceptance Criteria:**
- [ ] Order creation after payment
- [ ] Order status tracking (pending, processing, shipped, delivered)
- [ ] Order history for buyers
- [ ] Order management dashboard for sellers
- [ ] Invoice generation
- [ ] Order cancellation and refunds

**User Stories:**
- As a buyer, I want to track my order so I know when it will arrive
- As a seller, I want to manage orders so I can fulfill them efficiently

---

### FR-009: Shipping & Logistics (NinjaVan)
**Priority:** High
**Description:** Integrated shipping with NinjaVan

**Acceptance Criteria:**
- [ ] Automatic shipment creation
- [ ] Shipping label generation
- [ ] Real-time tracking updates
- [ ] Shipping cost calculation
- [ ] Delivery confirmation
- [ ] Support for Malaysian addresses
- [ ] Multiple shipping options (standard, express)

**User Stories:**
- As a buyer, I want to track my shipment so I know its location
- As a seller, I want to print shipping labels so I can send orders

---

### FR-010: Notifications
**Priority:** Medium
**Description:** Push and in-app notifications

**Acceptance Criteria:**
- [ ] Outbid notifications
- [ ] Auction ending soon alerts
- [ ] Order status updates
- [ ] New follower notifications
- [ ] Stream start alerts (scheduled streams)
- [ ] Push notifications (mobile)
- [ ] Email notifications

**User Stories:**
- As a buyer, I want to be notified when I'm outbid so I can bid again
- As a seller, I want to know when someone follows me

---

### FR-011: Search & Discovery
**Priority:** Medium
**Description:** Find products and streams

**Acceptance Criteria:**
- [ ] Full-text product search
- [ ] Filter by category, price, condition
- [ ] Sort by relevance, price, newest
- [ ] Live stream discovery
- [ ] Recommended products
- [ ] Trending auctions
- [ ] Seller search

**User Stories:**
- As a buyer, I want to search products so I can find what I need
- As a buyer, I want to see live streams so I can join ongoing auctions

---

### FR-012: Reviews & Ratings
**Priority:** Medium
**Description:** Buyer feedback system

**Acceptance Criteria:**
- [ ] Rate sellers (1-5 stars)
- [ ] Write text reviews
- [ ] Review with photos
- [ ] Seller response to reviews
- [ ] Review helpfulness voting
- [ ] Average rating display

**User Stories:**
- As a buyer, I want to read reviews so I can trust sellers
- As a seller, I want to receive ratings so I can build reputation

---

## Non-Functional Requirements

### Performance
- Page load time < 2 seconds
- Video stream latency < 3 seconds
- API response time < 200ms (95th percentile)
- Support 1000+ concurrent viewers per stream

### Security
- PCI DSS compliance for payments
- JWT token expiration and refresh
- Input validation and sanitization
- HTTPS everywhere
- Rate limiting on APIs

### Scalability
- Horizontal scaling for API servers
- Database read replicas
- CDN for static assets
- Caching layer (Redis)

### Availability
- 99.9% uptime target
- Graceful degradation
- Automated backups
- Disaster recovery plan

## Localization Requirements

### Language Support
- **Primary:** Malay (Bahasa Malaysia)
- **Secondary:** English
- **Future:** Mandarin, Tamil

### Regional Requirements
- Malaysia timezone (GMT+8)
- Malaysian Ringgit (MYR) currency
- Malaysian address format
- NinjaVan shipping integration
- FPX payment support

## Compliance Requirements

- **PDPA:** Personal Data Protection Act (Malaysia)
- **PCI DSS:** Payment Card Industry standards
- **Content Moderation:** Prohibited items policy
- **Tax:** SST (Sales and Service Tax) compliance

---

*Last updated: 2025-02-05*
