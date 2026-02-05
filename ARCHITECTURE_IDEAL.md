# Blytz.live - Ideal Architecture

## High-Level System Architecture

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                                    CLIENT LAYER                                          │
├─────────────────────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────────┐  ┌─────────────────────┐  ┌─────────────────────┐             │
│  │   Web App (Next.js) │  │  Mobile App (RN)    │  │  Seller Dashboard   │             │
│  │   - Live streaming  │  │  - Push notifications│  │  - Stream management │            │
│  │   - Real-time bids  │  │  - Mobile bidding   │  │  - Analytics        │             │
│  │   - Chat            │  │  - Camera streaming │  │  - Inventory        │             │
│  └──────────┬──────────┘  └──────────┬──────────┘  └──────────┬──────────┘             │
│             │                        │                        │                        │
│             └────────────────────────┼────────────────────────┘                        │
│                                      │                                                   │
└──────────────────────────────────────┼───────────────────────────────────────────────────┘
                                       │
                              ┌────────┴────────┐
                              │   CDN (Edge)    │
                              │  - Static assets│
                              │  - HLS streams  │
                              │  - API caching  │
                              └────────┬────────┘
                                       │
┌──────────────────────────────────────┼───────────────────────────────────────────────────┐
│                              API GATEWAY LAYER                                           │
├──────────────────────────────────────┼───────────────────────────────────────────────────┤
│                                      ▼                                                   │
│  ┌─────────────────────────────────────────────────────────────────────────────────────┐ │
│  │                        Kong/AWS API Gateway / Traefik                               │ │
│  │  - Rate limiting (global)  - Authentication  - Load balancing  - SSL termination    │ │
│  └─────────────────────────────────────────────────────────────────────────────────────┘ │
│                                      │                                                   │
└──────────────────────────────────────┼───────────────────────────────────────────────────┘
                                       │
┌──────────────────────────────────────┼───────────────────────────────────────────────────┐
│                         APPLICATION LAYER (Kubernetes/Docker Swarm)                      │
├──────────────────────────────────────┼───────────────────────────────────────────────────┤
│                                      │                                                   │
│  ┌─────────────────────────────────────────────────────────────────────────────────────┐ │
│  │                           EVENT BUS (Apache Kafka / NATS JetStream)                 │ │
│  │                                                                                     │ │
│  │   Topics:                                                                           │ │
│  │   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐               │ │
│  │   │  bid.placed │  │ auction.end │  │  user.join  │  │  payment.ok │               │ │
│  │   └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘               │ │
│  │   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐               │ │
│  │   │   chat.msg  │  │ stream.start│  │ order.create│  │notification │               │ │
│  │   └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘               │ │
│  └─────────────────────────────────────────────────────────────────────────────────────┘ │
│                                      │                                                   │
│  ┌───────────────────────────────────┼─────────────────────────────────────────────────┐ │
│  │                                   │                                                 │ │
│  │  ┌───────────────────────────────────────────────────────────────────────────────┐  │ │
│  │  │                     API SERVICES (Stateless, Horizontally Scalable)           │  │ │
│  │  │                                                                               │  │ │
│  │  │  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐          │  │ │
│  │  │  │  Auth        │ │  Catalog     │ │  Auction     │ │  Order       │          │  │ │
│  │  │  │  Service     │ │  Service     │ │  Service     │ │  Service     │          │  │ │
│  │  │  │              │ │              │ │              │ │              │          │  │ │
│  │  │  │ - JWT auth   │ │ - Products   │ │ - Bidding    │ │ - Checkout   │          │  │ │
│  │  │  │ - OAuth      │ │ - Categories │ │ - Auto-bid   │ │ - Payment    │          │  │ │
│  │  │  │ - Sessions   │ │ - Search     │ │ - Scheduling │ │ - Fulfillment│          │  │ │
│  │  │  └──────────────┘ └──────────────┘ └──────────────┘ └──────────────┘          │  │ │
│  │  │                                                                               │  │ │
│  │  │  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐          │  │ │
│  │  │  │  Payment     │ │  User        │ │  Analytics   │ │  Notification│          │  │ │
│  │  │  │  Service     │ │  Service     │ │  Service     │ │  Service     │          │  │ │
│  │  │  │              │ │              │ │              │ │              │          │  │ │
│  │  │  │ - Stripe     │ │ - Profiles   │ │ - Events     │ │ - Email      │          │  │ │
│  │  │  │ - Webhooks   │ │ - Addresses  │ │ - Reports    │ │ - Push       │          │  │ │
│  │  │  │ - Refunds    │ │ - Preferences│ │ - ML models  │ │ - SMS        │          │  │ │
│  │  │  └──────────────┘ └──────────────┘ └──────────────┘ └──────────────┘          │  │ │
│  │  └───────────────────────────────────────────────────────────────────────────────┘  │ │
│  │                                   │                                                 │ │
│  │  ┌───────────────────────────────────────────────────────────────────────────────┐  │ │
│  │  │                    REAL-TIME SERVICES (WebSocket/SSE)                         │  │ │
│  │  │                                                                               │  │ │
│  │  │  ┌─────────────────────────────────────────────────────────────────────────┐  │  │ │
│  │  │  │                    WebSocket Gateway Service (3+ replicas)              │  │  │ │
│  │  │  │                                                                           │  │  │ │
│  │  │  │  ┌─────────────┐     ┌─────────────┐     ┌─────────────┐               │  │  │ │
│  │  │  │  │  Replica 1  │◄───►│  Replica 2  │◄───►│  Replica 3  │               │  │  │ │
│  │  │  │  │             │     │             │     │             │               │  │  │ │
│  │  │  │  │ Local       │     │ Local       │     │ Local       │               │  │  │ │
│  │  │  │  │ Connections │     │ Connections │     │ Connections │               │  │  │ │
│  │  │  │  └──────┬──────┘     └──────┬──────┘     └──────┬──────┘               │  │  │ │
│  │  │  │         │                   │                   │                       │  │  │ │
│  │  │  │         └───────────────────┼───────────────────┘                       │  │  │ │
│  │  │  │                             │                                           │  │  │ │
│  │  │  │                    ┌────────┴────────┐                                   │  │  │ │
│  │  │  │                    │  Redis Pub/Sub  │◄────── Cross-instance sync       │  │  │ │
│  │  │  │                    └─────────────────┘                                   │  │  │ │
│  │  │  └─────────────────────────────────────────────────────────────────────────┘  │  │ │
│  │  └───────────────────────────────────────────────────────────────────────────────┘  │ │
│  │                                                                                     │ │
│  │  ┌───────────────────────────────────────────────────────────────────────────────┐  │ │
│  │  │                    BACKGROUND WORKERS (Async Processing)                      │  │ │
│  │  │                                                                               │  │ │
│  │  │  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐          │  │ │
│  │  │  │ Email        │ │ Search       │ │ Image        │ │ Report       │          │  │ │
│  │  │  │ Worker       │ │ Indexer      │ │ Processor    │ │ Generator    │          │  │ │
│  │  │  └──────────────┘ └──────────────┘ └──────────────┘ └──────────────┘          │  │ │
│  │  └───────────────────────────────────────────────────────────────────────────────┘  │ │
│  │                                                                                     │ │
│  │  ┌───────────────────────────────────────────────────────────────────────────────┐  │ │
│  │  │                    STREAMING INFRASTRUCTURE (LiveKit Cluster)                 │  │ │
│  │  │                                                                               │  │ │
│  │  │  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐                          │  │ │
│  │  │  │ LiveKit      │ │ LiveKit      │ │ LiveKit      │                          │  │ │
│  │  │  │ Server 1     │ │ Server 2     │ │ Server 3     │                          │  │ │
│  │  │  │ (us-east)    │ │ (eu-west)    │ │ (ap-south)   │                          │  │ │
│  │  │  └──────────────┘ └──────────────┘ └──────────────┘                          │  │ │
│  │  └───────────────────────────────────────────────────────────────────────────────┘  │ │
│  └─────────────────────────────────────────────────────────────────────────────────────┘ │
│                                                                                          │
└──────────────────────────────────────────────────────────────────────────────────────────┘
                                       │
┌──────────────────────────────────────┼───────────────────────────────────────────────────┐
│                              DATA LAYER                                                  │
├──────────────────────────────────────┼───────────────────────────────────────────────────┤
│                                      │                                                   │
│  ┌─────────────────────────────────────────────────────────────────────────────────────┐ │
│  │                         PRIMARY DATABASE CLUSTER (PostgreSQL)                       │ │
│  │                                                                                     │ │
│  │   ┌─────────────┐         ┌─────────────┐         ┌─────────────┐                  │ │
│  │   │   Primary   │◄───────►│  Replica 1  │◄───────►│  Replica 2  │                  │ │
│  │   │  (Writes)   │  Sync   │  (Reads)    │  Async  │  (Reads)    │                  │ │
│  │   └─────────────┘         └─────────────┘         └─────────────┘                  │ │
│  │                                                                                     │ │
│  │   Connection Pooling: PgBouncer                                                   │ │
│  │   Sharding: By auction_id (future)                                                │ │
│  └─────────────────────────────────────────────────────────────────────────────────────┘ │
│                                      │                                                   │
│  ┌───────────────────────────────────┼─────────────────────────────────────────────────┐ │
│  │                                   │                                                 │ │
│  │  ┌───────────────────────────────────────────────────────────────────────────────┐  │ │
│  │  │                         CACHE LAYER (Redis Cluster)                           │  │ │
│  │  │                                                                               │  │ │
│  │  │   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │  │ │
│  │  │   │  Master 1   │  │  Master 2   │  │  Master 3   │  │  Sentinel   │        │  │ │
│  │  │   │  (Hash 0)   │  │  (Hash 1)   │  │  (Hash 2)   │  │  (HA)       │        │  │ │
│  │  │   └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘        │  │ │
│  │  │                                                                               │  │ │
│  │  │   Use Cases:                                                                  │  │ │
│  │  │   - Session store          - Auction state cache                              │  │ │
│  │  │   - Rate limiting          - Real-time leaderboards                           │  │ │
│  │  │   - Pub/Sub                - Query result caching                             │  │ │
│  │  └───────────────────────────────────────────────────────────────────────────────┘  │ │
│  │                                                                                     │ │
│  │  ┌───────────────────────────────────────────────────────────────────────────────┐  │ │
│  │  │                      SEARCH ENGINE (Elasticsearch/OpenSearch)                 │  │ │
│  │  │                                                                               │  │ │
│  │  │   - Product search with filters                                               │  │ │
│  │  │   - Full-text search                                                          │  │ │
│  │  │   - Auto-complete suggestions                                                 │  │ │
│  │  │   - Analytics aggregations                                                    │  │ │
│  │  └───────────────────────────────────────────────────────────────────────────────┘  │ │
│  │                                                                                     │ │
│  │  ┌───────────────────────────────────────────────────────────────────────────────┐  │ │
│  │  │                      OBJECT STORAGE (S3 / MinIO)                              │  │ │
│  │  │                                                                               │  │ │
│  │  │   - Product images                                                            │  │ │
│  │  │   - Stream recordings                                                         │  │ │
│  │  │   - User avatars                                                              │  │ │
│  │  │   - Export files                                                              │  │ │
│  │  └───────────────────────────────────────────────────────────────────────────────┘  │ │
│  └─────────────────────────────────────────────────────────────────────────────────────┘ │
│                                                                                          │
└──────────────────────────────────────────────────────────────────────────────────────────┘
                                       │
┌──────────────────────────────────────┼───────────────────────────────────────────────────┐
│                         EXTERNAL SERVICES                                                │
├──────────────────────────────────────┼───────────────────────────────────────────────────┤
│                                      ▼                                                   │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐                  │
│  │   Stripe     │  │   SendGrid   │  │   Firebase   │  │   Datadog    │                  │
│  │  (Payments)  │  │  (Email)     │  │  (Push Notif)│  │ (Monitoring) │                  │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘                  │
│                                                                                          │
└──────────────────────────────────────────────────────────────────────────────────────────┘
```

---

## Component Deep Dive

### 1. API Services Architecture

Each service follows **Clean Architecture**:

```
auction-service/
├── cmd/
│   └── server/
│       └── main.go                    # Entry point (30 lines)
├── internal/
│   ├── domain/                        # Business logic (no deps)
│   │   ├── auction.go                 # Domain entities
│   │   ├── bid.go                     # Bid entity & rules
│   │   ├── repository.go              # Repository interfaces
│   │   └── service.go                 # Business logic interface
│   │
│   ├── application/                   # Use cases
│   │   ├── place_bid.go               # PlaceBid command handler
│   │   ├── start_auction.go           # StartAuction command
│   │   ├── end_auction.go             # EndAuction command
│   │   └── queries.go                 # Read models
│   │
│   └── infrastructure/                # External implementations
│       ├── persistence/
│       │   ├── postgres/
│       │   │   ├── auction_repo.go    # GORM implementation
│       │   │   └── models.go          # GORM models
│       │   └── cache/
│       │       └── auction_cache.go   # Redis cache
│       │
│       ├── messaging/
│       │   └── kafka/
│       │       └── event_publisher.go # Event publishing
│       │
│       └── http/
│           ├── handlers.go            # HTTP handlers
│           ├── routes.go              # Route definitions
│           └── dto.go                 # Request/response DTOs
│
├── pkg/
│   └── auction/
│       └── client.go                  # Public API client
│
└── Dockerfile
```

### 2. Real-Time Data Flow (Bid Placement)

```
User places bid
       │
       ▼
┌──────────────┐
│   WebSocket  │  1. Bid received via WebSocket
│   Gateway    │     (sticky session routes to specific replica)
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   Auction    │  2. Bid validation
│   Service    │     - Check auction status
└──────┬───────┘     - Validate bid amount
       │             - Check user balance
       ▼
┌──────────────┐
│  PostgreSQL  │  3. Atomic transaction (pessimistic lock)
│  (Primary)   │     INSERT bid
└──────┬───────┘     UPDATE auction.current_bid
       │             UPDATE user.balance (reserve)
       ▼
┌──────────────┐
│    Kafka     │  4. Publish event
│  (bid.placed)│     Topic: auction.{id}.events
└──────┬───────┘
       │
       ├──────────────────┬──────────────────┐
       ▼                  ▼                  ▼
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│   WebSocket  │  │   Analytics  │  │ Notification │
│   Consumers  │  │   Consumer   │  │  Consumer    │
│              │  │              │  │              │
│ Broadcast to │  │ Record event │  │ Send push    │
│ all replicas │  │ Update stats │  │ to outbid user│
└──────┬───────┘  └──────────────┘  └──────────────┘
       │
       ▼
┌──────────────┐
│  Redis Pub/  │  5. Cross-replica sync
│     Sub      │     All WebSocket replicas receive event
└──────┬───────┘
       │
       ├──────────┬──────────┐
       ▼          ▼          ▼
┌──────────┐ ┌──────────┐ ┌──────────┐
│ Replica 1│ │ Replica 2│ │ Replica 3│  6. Broadcast to
│  Users   │ │  Users   │ │  Users   │     connected clients
│  A,B,C   │ │  D,E,F   │ │  G,H,I   │
└──────────┘ └──────────┘ └──────────┘
```

### 3. Streaming Architecture (LiveKit)

```
Seller starts stream
       │
       ▼
┌─────────────────────────────────────────────────────┐
│              LiveKit Cluster                         │
│                                                      │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │
│  │ SFU Server  │  │ SFU Server  │  │ SFU Server  │  │
│  │ (us-east)   │  │ (eu-west)   │  │ (ap-south)  │  │
│  │             │  │             │  │             │  │
│  │ ┌─────────┐ │  │ ┌─────────┐ │  │ ┌─────────┐ │  │
│  │ │ Room    │ │  │ │ Room    │ │  │ │ Room    │ │  │
│  │ │Router   │ │  │ │Router   │ │  │ │Router   │ │  │
│  │ └────┬────┘ │  │ └────┬────┘ │  │ └────┬────┘ │  │
│  │      │      │  │      │      │  │      │      │  │
│  │ ┌────┴────┐ │  │ ┌────┴────┐ │  │ ┌────┴────┐ │  │
│  │ │Publisher│ │  │ │Publisher│ │  │ │Publisher│ │  │
│  │ │ (Seller)│ │  │ │ (Seller)│ │  │ │ (Seller)│ │  │
│  │ └────┬────┘ │  │ └────┬────┘ │  │ └────┬────┘ │  │
│  │      │      │  │      │      │  │      │      │  │
│  │ ┌────┴────┐ │  │ ┌────┴────┐ │  │ ┌────┴────┐ │  │
│  │ │ Viewers │ │  │ │ Viewers │ │  │ │ Viewers │ │  │
│  │ │ A,B,C   │ │  │ │ D,E,F   │ │  │ │ G,H,I   │ │  │
│  │ └─────────┘ │  │ └─────────┘ │  │ └─────────┘ │  │
│  └─────────────┘  └─────────────┘  └─────────────┘  │
│                                                      │
│  Recording → S3  │  RTMP → YouTube (optional)       │
└─────────────────────────────────────────────────────┘
```

### 4. Database Schema (Sharding Strategy)

```
┌─────────────────────────────────────────────────────────────┐
│                    Database Cluster                          │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────────────┐  ┌──────────────────────┐        │
│  │   Shard 1            │  │   Shard 2            │        │
│  │   (Auctions A-M)     │  │   (Auctions N-Z)     │        │
│  │                      │  │                      │        │
│  │   auctions           │  │   auctions           │        │
│  │   bids               │  │   bids               │        │
│  │   auction_watchers   │  │   auction_watchers   │        │
│  └──────────────────────┘  └──────────────────────┘        │
│                                                              │
│  ┌─────────────────────────────────────────────────────────┐│
│  │   Global Tables (All Shards)                            ││
│  │   - users                                               ││
│  │   - products                                            ││
│  │   - categories                                          ││
│  │   - orders (with auction_id for sharding reference)     ││
│  │   - payments                                            ││
│  └─────────────────────────────────────────────────────────┘│
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## Technology Stack

### Core Services
| Component | Technology | Reason |
|-----------|-----------|---------|
| Language | Go 1.21+ | Performance, concurrency |
| Web Framework | Gin or Echo | Fast, mature |
| Database | PostgreSQL 15+ | ACID, JSON support |
| Cache | Redis Cluster | Pub/Sub, fast access |
| Message Queue | Apache Kafka or NATS | Durability, replay |
| Search | Elasticsearch | Full-text, aggregations |
| WebSocket | Custom + Redis Pub/Sub | Horizontal scaling |
| Streaming | LiveKit | WebRTC, scalable |
| Storage | S3/MinIO | Object storage |

### Infrastructure
| Component | Technology | Reason |
|-----------|-----------|---------|
| Container Orchestration | Kubernetes | Industry standard |
| Service Mesh | Istio or Linkerd | mTLS, traffic management |
| API Gateway | Kong or Traefik | Rate limiting, auth |
| Monitoring | Prometheus + Grafana | Metrics, visualization |
| Logging | ELK Stack or Loki | Centralized logging |
| Tracing | Jaeger or Tempo | Distributed tracing |
| Secrets | HashiCorp Vault | Secure secret management |

---

## Scaling Patterns

### Horizontal Scaling by Service

```
┌─────────────────────────────────────────────────────────┐
│                 Load Balancer                            │
└────────────────────┬────────────────────────────────────┘
                     │
        ┌────────────┼────────────┐
        ▼            ▼            ▼
┌──────────────┐ ┌──────────┐ ┌──────────────┐
│ Auth Service │ │   ...    │ │ Order Service │
│  (3 replicas)│ │          │ │  (2 replicas) │
└──────────────┘ └──────────┘ └──────────────┘
        │                           │
        ▼                           ▼
┌──────────────┐            ┌──────────────┐
│ Redis        │            │ PostgreSQL   │
│ (Sessions)   │            │ (Read: 3,    │
│              │            │  Write: 1)   │
└──────────────┘            └──────────────┘
```

### Caching Strategy (Multi-Layer)

```
User Request
     │
     ▼
┌──────────────┐  L1: CDN (Cloudflare)
│     CDN      │  - Static assets
│   (Global)   │  - Product images
└──────┬───────┘  - API responses (short TTL)
       │
       ▼ (miss)
┌──────────────┐  L2: Edge Cache (Varnish/Nginx)
│  Edge Cache  │  - Auction state (1s TTL)
│  (Regional)  │  - Hot products (60s TTL)
└──────┬───────┘
       │
       ▼ (miss)
┌──────────────┐  L3: Application Cache (Redis)
│  App Cache   │  - Session data
│   (Redis)    │  - Rate limit counters
└──────┬───────┘  - Query results
       │
       ▼ (miss)
┌──────────────┐  L4: Database
│  PostgreSQL  │  - Persistent storage
│   (Primary   │  - Read replicas for queries
│   + Replicas)│
└──────────────┘
```

---

## Data Flow Examples

### 1. User Places Bid (Happy Path)

```
┌─────────┐     ┌─────────────┐     ┌──────────────┐     ┌──────────┐
│  User   │────►│ WebSocket   │────►│   Auction    │────►│   DB     │
│         │     │  Gateway    │     │   Service    │     │  (Lock)  │
└─────────┘     └─────────────┘     └──────────────┘     └────┬─────┘
     ▲                                                         │
     │                                                         │
     │    ┌─────────────┐     ┌─────────────┐                 │
     └────┤   Kafka     │◄────┤  WebSocket  │◄────────────────┘
          │  (Event)    │     │  Broadcast  │
          └─────────────┘     └─────────────┘
```

### 2. Auction Ends & Order Created

```
Cron Job / Timer
      │
      ▼
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   Auction    │────►│   Auction    │────►│    Kafka     │
│   Service    │     │   Ends       │     │  (auction.   │
│   (Scheduler)│     │   (Winner)   │     │   ended)     │
└──────────────┘     └──────────────┘     └──────┬───────┘
                                                  │
                    ┌─────────────────────────────┼─────────────┐
                    │                             │             │
                    ▼                             ▼             ▼
            ┌──────────────┐           ┌──────────────┐ ┌──────────────┐
            │   Order      │           │Notification  │ │  Analytics   │
            │   Service    │           │  Service     │ │   Service    │
            │ (Create      │           │ (Notify      │ │  (Record)    │
            │   order)     │           │  winner)     │ │              │
            └──────────────┘           └──────────────┘ └──────────────┘
```

### 3. Stream Recording & Playback

```
Seller Stream
      │
      ▼
┌──────────────┐
│   LiveKit    │
│   Server     │
└──────┬───────┘
       │
       ├──────────┬──────────┐
       │          │          │
       ▼          ▼          ▼
┌──────────┐ ┌──────────┐ ┌──────────┐
│  Live    │ │ Record   │ │  HLS     │
│  Viewers │ │  to S3   │ │  Segment │
└──────────┘ └────┬─────┘ └────┬─────┘
                  │            │
                  ▼            ▼
           ┌──────────┐ ┌──────────┐
           │   S3     │ │  CDN     │
           │Storage   │ │Playback  │
           └──────────┘ └──────────┘
```

---

## Deployment Architecture (Kubernetes)

```
┌─────────────────────────────────────────────────────────────┐
│                     Kubernetes Cluster                       │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌─────────────────────────────────────────────────────────┐│
│  │ Namespace: blytz-production                             ││
│  │                                                          ││
│  │  Deployment: api-gateway (3 replicas)                   ││
│  │  ├── Pod: kong-xxx (Container: kong)                    ││
│  │  ├── Pod: kong-yyy                                      ││
│  │  └── Pod: kong-zzz                                      ││
│  │                                                          ││
│  │  Deployment: auction-service (5 replicas)               ││
│  │  ├── Pod: auction-xxx (Container: auction-svc)          ││
│  │  ├── Pod: auction-yyy                                   ││
│  │  └── ...                                                ││
│  │                                                          ││
│  │  Deployment: websocket-gateway (3 replicas)             ││
│  │  └── ...                                                ││
│  │                                                          ││
│  │  StatefulSet: kafka (3 brokers)                         ││
│  │  └── ...                                                ││
│  │                                                          ││
│  │  CronJob: auction-end-processor                         ││
│  │  └── Runs every minute                                  ││
│  │                                                          ││
│  │  Job: db-migration (run on deploy)                      ││
│  │  └── ...                                                ││
│  │                                                          ││
│  └─────────────────────────────────────────────────────────┘│
│                                                              │
│  Services:                                                   │
│  ├── api-gateway (LoadBalancer → External IP)               │
│  ├── auction-service (ClusterIP → Internal)                 │
│  ├── websocket-gateway (LoadBalancer → External IP)         │
│  └── kafka (ClusterIP → Internal)                           │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## Monitoring & Observability

```
┌─────────────────────────────────────────────────────────────┐
│                     Observability Stack                      │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │ Prometheus  │  │   Grafana   │  │   Alertmanager│         │
│  │  (Metrics)  │  │ (Dashboards)│  │  (Alerts)   │          │
│  └──────┬──────┘  └─────────────┘  └─────────────┘          │
│         │                                                    │
│         │ Pull metrics from:                                 │
│         ├── Application /metrics endpoint                    │
│         ├── Node Exporter (system metrics)                   │
│         ├── PostgreSQL Exporter                              │
│         ├── Redis Exporter                                   │
│         └── Kafka Exporter                                   │
│                                                              │
│  ┌─────────────┐  ┌─────────────┐                           │
│  │    Loki     │  │   Tempo     │                           │
│  │   (Logs)    │  │  (Traces)   │                           │
│  └─────────────┘  └─────────────┘                           │
│                                                              │
│  Key Metrics:                                                │
│  - Request latency (p50, p95, p99)                          │
│  - Error rate by endpoint                                    │
│  - Active WebSocket connections                              │
│  - Bid throughput per auction                                │
│  - Database connection pool usage                            │
│  - Cache hit/miss rates                                      │
│  - Stream latency and bitrate                                │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## Summary

This architecture supports:

| Metric | Capacity |
|--------|----------|
| Concurrent Users | 1,000,000+ |
| Concurrent Auctions | 10,000+ |
| Bids/Second | 50,000+ |
| Stream Viewers | 100,000 per auction |
| Availability | 99.99% |
| Latency (p95) | <100ms API, <500ms WebSocket |

**Key Principles:**
1. **Stateless services** - Easy horizontal scaling
2. **Event-driven** - Decoupled, async processing
3. **CQRS** - Separate read/write paths
4. **Multi-layer caching** - Performance at scale
5. **Observability** - Metrics, logs, traces everywhere

Want me to implement any specific part of this architecture?