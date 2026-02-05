# Environment Configuration

## Overview

Environment variable configuration for all environments: development, staging, and production.

## Environment Files

### File Structure
```
├── .env.development           # Local development
├── .env.staging               # Staging environment
├── .env.production            # Production environment
└── .env.example               # Template (committed to repo)
```

## Core Configuration

### Backend (.env)

```bash
# ==========================================
# Server Configuration
# ==========================================
PORT=8080
ENV=development                    # development | staging | production
LOG_LEVEL=debug                    # debug | info | warn | error

# ==========================================
# Database (PostgreSQL)
# ==========================================
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=blytz
DATABASE_USER=postgres
DATABASE_PASSWORD=your_secure_password
DATABASE_SSL_MODE=disable          # disable | require | verify-ca | verify-full
DATABASE_URL=postgresql://${DATABASE_USER}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}

# Connection Pool
DATABASE_MAX_OPEN_CONNS=100
DATABASE_MAX_IDLE_CONNS=50
DATABASE_CONN_MAX_LIFETIME=1h

# ==========================================
# Cache (Redis)
# ==========================================
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_URL=redis://${REDIS_HOST}:${REDIS_PORT}/${REDIS_DB}

# Redis Key Prefix (for multi-tenant)
REDIS_KEY_PREFIX=blytz:dev

# ==========================================
# Authentication (JWT)
# ==========================================
JWT_SECRET=your-super-secret-jwt-key-min-32-chars
JWT_ACCESS_EXPIRY=15m              # Access token expiry
JWT_REFRESH_EXPIRY=7d              # Refresh token expiry
JWT_ISSUER=blytz.app
JWT_AUDIENCE=blytz.app

# ==========================================
# Storage (Cloudflare R2)
# ==========================================
R2_ACCOUNT_ID=your-account-id
R2_ACCESS_KEY_ID=your-access-key
R2_SECRET_ACCESS_KEY=your-secret-key
R2_BUCKET_NAME=blytz-storage
R2_ENDPOINT=https://${R2_ACCOUNT_ID}.r2.cloudflarestorage.com
R2_PUBLIC_URL=https://cdn.blytz.app

# ==========================================
# Payments (Stripe)
# ==========================================
STRIPE_SECRET_KEY=sk_test_...
STRIPE_WEBHOOK_SECRET=whsec_...
STRIPE_PUBLISHABLE_KEY=pk_test_...

# Stripe Connect (for seller payouts)
STRIPE_CONNECT_CLIENT_ID=ca_...

# Payment Methods
STRIPE_ENABLE_FPX=true
STRIPE_ENABLE_GRABPAY=true
STRIPE_ENABLE_TOUCHNGO=true

# ==========================================
# Live Streaming (LiveKit)
# ==========================================
LIVEKIT_API_KEY=your-api-key
LIVEKIT_API_SECRET=your-api-secret
LIVEKIT_WS_URL=wss://livekit.blytz.app
LIVEKIT_HTTP_URL=https://livekit-api.blytz.app

# Recording
LIVEKIT_RECORDING_BUCKET=${R2_BUCKET_NAME}
LIVEKIT_RECORDING_ENDPOINT=${R2_ENDPOINT}

# ==========================================
# Shipping (NinjaVan)
# ==========================================
NINJAVAN_BASE_URL=https://api.ninjavan.co
NINJAVAN_API_KEY=your-api-key
NINJAVAN_API_SECRET=your-api-secret
NINJAVAN_WEBHOOK_SECRET=njv_webhook_...

# Sandbox for testing
NINJAVAN_SANDBOX=true

# ==========================================
# Real-time (Socket.io)
# ==========================================
SOCKET_CORS_ORIGIN=http://localhost:3000
SOCKET_REDIS_ADAPTER=true

# ==========================================
# Email (SendGrid)
# ==========================================
SENDGRID_API_KEY=SG.xxx
EMAIL_FROM=noreply@blytz.app
EMAIL_FROM_NAME=Blytz

# Templates
SENDGRID_TEMPLATE_WELCOME=d-xxx
SENDGRID_TEMPLATE_ORDER_CONFIRMATION=d-xxx
SENDGRID_TEMPLATE_SHIPPING=d-xxx

# ==========================================
# Push Notifications (Firebase)
# ==========================================
FIREBASE_PROJECT_ID=blytz-app
FIREBASE_SERVICE_ACCOUNT_KEY_PATH=/path/to/service-account.json

# ==========================================
# Monitoring & Logging
# ==========================================
SENTRY_DSN=https://xxx@xxx.ingest.sentry.io/xxx
SENTRY_ENVIRONMENT=${ENV}

# Datadog (optional)
DATADOG_API_KEY=xxx
DATADOG_APP_KEY=xxx

# ==========================================
# Feature Flags
# ==========================================
ENABLE_LIVE_STREAMING=true
ENABLE_AUCTIONS=true
ENABLE_FPX_PAYMENTS=true
ENABLE_AUTO_SHIPPING=true
```

### Frontend (.env.local)

```bash
# ==========================================
# API Configuration
# ==========================================
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_SOCKET_URL=http://localhost:8080

# ==========================================
# Live Streaming
# ==========================================
NEXT_PUBLIC_LIVEKIT_URL=wss://livekit.blytz.app

# ==========================================
# Payments
# ==========================================
NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY=pk_test_...

# ==========================================
# App Configuration
# ==========================================
NEXT_PUBLIC_APP_NAME=Blytz
NEXT_PUBLIC_APP_URL=http://localhost:3000
NEXT_PUBLIC_CDN_URL=https://cdn.blytz.app

# ==========================================
# Feature Flags
# ==========================================
NEXT_PUBLIC_ENABLE_ANALYTICS=false
NEXT_PUBLIC_MAINTENANCE_MODE=false
```

### Mobile (.env)

```bash
# ==========================================
# API Configuration
# ==========================================
API_URL=http://localhost:8080/api/v1
SOCKET_URL=http://localhost:8080

# ==========================================
# Services
# ==========================================
LIVEKIT_URL=wss://livekit.blytz.app
STRIPE_PUBLISHABLE_KEY=pk_test_...

# ==========================================
# App
# ==========================================
APP_NAME=Blytz
APP_VERSION=1.0.0
```

## Environment-Specific Values

### Development
```bash
# .env.development
ENV=development
LOG_LEVEL=debug
DATABASE_HOST=localhost
REDIS_HOST=localhost
STRIPE_SECRET_KEY=sk_test_...
SENDGRID_API_KEY=SG.test_...
```

### Staging
```bash
# .env.staging
ENV=staging
LOG_LEVEL=info
DATABASE_HOST=staging-db.blytz.internal
REDIS_HOST=staging-redis.blytz.internal
STRIPE_SECRET_KEY=sk_test_...
SENDGRID_API_KEY=SG.live_...
```

### Production
```bash
# .env.production
ENV=production
LOG_LEVEL=warn
DATABASE_HOST=prod-db.blytz.internal
REDIS_HOST=prod-redis.blytz.internal
STRIPE_SECRET_KEY=sk_live_...
SENDGRID_API_KEY=SG.live_...
```

## Secrets Management

### Local Development
Use `.env.local` (not committed to git):
```bash
# Copy from template
cp .env.example .env.local

# Edit with your values
nano .env.local
```

### Staging/Production
Use secret management service:

**Option 1: GitHub Secrets**
```yaml
# .github/workflows/deploy.yml
- name: Deploy
  env:
    DATABASE_URL: ${{ secrets.STAGING_DATABASE_URL }}
    JWT_SECRET: ${{ secrets.JWT_SECRET }}
```

**Option 2: AWS Secrets Manager**
```bash
# Fetch secrets at runtime
aws secretsmanager get-secret-value --secret-id blytz/production
```

**Option 3: HashiCorp Vault**
```bash
# Inject via sidecar
vault kv get -format=json secret/blytz/production
```

**Option 4: Docker Secrets (Swarm/K8s)**
```yaml
# docker-compose.yml
secrets:
  db_password:
    external: true

services:
  api:
    secrets:
      - db_password
```

## Configuration Validation

### Backend (Go)
```go
// config/config.go
package config

type Config struct {
  Port        string `env:"PORT" envDefault:"8080"`
  Environment string `env:"ENV" envDefault:"development"`
  
  Database struct {
    Host     string `env:"DATABASE_HOST,required"`
    Port     int    `env:"DATABASE_PORT" envDefault:"5432"`
    Name     string `env:"DATABASE_NAME,required"`
    User     string `env:"DATABASE_USER,required"`
    Password string `env:"DATABASE_PASSWORD,required"`
  }
  
  JWT struct {
    Secret        string        `env:"JWT_SECRET,required"`
    AccessExpiry  time.Duration `env:"JWT_ACCESS_EXPIRY" envDefault:"15m"`
    RefreshExpiry time.Duration `env:"JWT_REFRESH_EXPIRY" envDefault:"7d"`
  }
  
  Stripe struct {
    SecretKey      string `env:"STRIPE_SECRET_KEY,required"`
    WebhookSecret  string `env:"STRIPE_WEBHOOK_SECRET,required"`
  }
}

func Load() (*Config, error) {
  var cfg Config
  
  if err := env.Parse(&cfg); err != nil {
    return nil, fmt.Errorf("failed to parse config: %w", err)
  }
  
  // Validate required configs
  if err := validate(&cfg); err != nil {
    return nil, err
  }
  
  return &cfg, nil
}

func validate(cfg *Config) error {
  if len(cfg.JWT.Secret) < 32 {
    return errors.New("JWT_SECRET must be at least 32 characters")
  }
  
  if cfg.Environment == "production" {
    if !strings.HasPrefix(cfg.Stripe.SecretKey, "sk_live_") {
      return errors.New("production requires live Stripe keys")
    }
  }
  
  return nil
}
```

### Health Check Endpoint
```go
// Include config validation in health check
func HealthCheck(c *gin.Context) {
  checks := map[string]bool{
    "database": checkDatabase(),
    "redis": checkRedis(),
    "stripe": checkStripe(),
  }
  
  allHealthy := true
  for _, healthy := range checks {
    if !healthy {
      allHealthy = false
      break
    }
  }
  
  status := http.StatusOK
  if !allHealthy {
    status = http.StatusServiceUnavailable
  }
  
  c.JSON(status, gin.H{
    "status": map[string]bool{
      "healthy": allHealthy,
    },
    "checks": checks,
  })
}
```

## Docker Configuration

### docker-compose.yml
```yaml
version: '3.8'

services:
  api:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - ENV=${ENV:-development}
      - DATABASE_URL=${DATABASE_URL}
      - REDIS_URL=${REDIS_URL}
      - JWT_SECRET=${JWT_SECRET}
      - STRIPE_SECRET_KEY=${STRIPE_SECRET_KEY}
    env_file:
      - .env.local
    depends_on:
      - postgres
      - redis

  postgres:
    image: postgres:17
    environment:
      - POSTGRES_DB=blytz
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=${DATABASE_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  redis:
    image: redis:8-alpine
    ports:
      - "6379:6379"

volumes:
  postgres_data:
```

## Security Checklist

- [ ] `.env*.local` in `.gitignore`
- [ ] Production secrets never in code
- [ ] Secrets rotated regularly
- [ ] Different secrets per environment
- [ ] Strong JWT secrets (32+ chars)
- [ ] Database passwords complex
- [ ] API keys scoped to minimum permissions
- [ ] Webhook secrets validated

---

*Last updated: 2025-02-05*
