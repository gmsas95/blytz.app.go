#!/bin/bash

# ==============================
# GEMINI3-MOCK + GO BACKEND DEPLOYMENT
# ==============================

echo "🚀 BLYTZ.LIVE - Gemini3 Integration Deployment"
echo "================================================"
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}✅${NC} $1"
}

print_info() {
    echo -e "${BLUE}ℹ️${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠️${NC} $1"
}

# 1. Stop existing services
echo "🛑 Stopping existing services..."
docker-compose down

# 2. Build enhanced backend
echo ""
print_info "Building enhanced Go backend with flash sales & reviews..."
cd backend
go mod tidy
docker build -t blytz-backend-gemini .
cd ..
print_status "Backend built successfully"

# 3. Build Gemini3-mock frontend
echo ""
print_info "Building Gemini3-mock frontend with backend integration..."
cd ../blytz-gemini-frontend
docker build -t blytz-frontend-gemini .
cd ../blytz.live.remake
print_status "Frontend built successfully"

# 4. Update docker-compose for Gemini3 integration
echo ""
print_info "Updating docker-compose.yml for Gemini3 integration..."
cat > docker-compose.gemini.yml << 'EOF'
version: '3.8'

services:
  frontend:
    image: blytz-frontend-gemini:latest
    container_name: blytz-frontend
    restart: unless-stopped
    environment:
      VITE_API_URL: ${VITE_API_URL:-http://backend:8080/api/v1}
      VITE_GEMINI_API_KEY: ${GEMINI_API_KEY:-}
      VITE_APP_NAME: ${VITE_APP_NAME:-Blytz.app Marketplace}
      VITE_APP_VERSION: ${VITE_APP_VERSION:-1.0.0}
    networks:
      - blytz-network
    ports:
      - "3000:80"

  backend:
    image: blytz-backend-gemini:latest
    container_name: blytz-backend
    restart: unless-stopped
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: ${DB_USER:-blytz_user}
      DB_PASSWORD: ${DB_PASSWORD}
      DB_NAME: ${DB_NAME:-blytz_marketplace}
      DB_SSL_MODE: disable
      REDIS_HOST: redis
      REDIS_PORT: 6379
      REDIS_PASSWORD: ${REDIS_PASSWORD}
      JWT_SECRET: ${JWT_SECRET}
      JWT_EXPIRES_IN: ${JWT_EXPIRES_IN:-168h}
      PORT: 8080
      CORS_ORIGINS: ${CORS_ORIGINS:-http://frontend:3000,http://localhost:3000}
      ENVIRONMENT: production
      GEMINI_API_KEY: ${GEMINI_API_KEY:-}
    networks:
      - blytz-network
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    ports:
      - "8080:8080"

  postgres:
    image: postgres:15-alpine
    container_name: blytz-postgres
    restart: unless-stopped
    environment:
      POSTGRES_DB: ${DB_NAME:-blytz_marketplace}
      POSTGRES_USER: ${DB_USER:-blytz_user}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./scripts/init.sql:/docker-entrypoint-initdb.d/init.sql:ro
    networks:
      - blytz-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER:-blytz_user} -d ${DB_NAME:-blytz_marketplace}"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    container_name: blytz-redis
    restart: unless-stopped
    command: redis-server --requirepass ${REDIS_PASSWORD}
    volumes:
      - redis_data:/data
    networks:
      - blytz-network
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 3s
      retries: 5

volumes:
  postgres_data:
  redis_data:

networks:
  blytz-network:
    driver: bridge
EOF

print_status "docker-compose.gemini.yml created"

# 5. Start Gemini3 integration
echo ""
print_info "Starting Gemini3-mock + Go backend integration..."
docker-compose -f docker-compose.gemini.yml up -d

# 6. Wait for services to be healthy
echo ""
print_info "Waiting for services to be healthy..."
sleep 10

# 7. Check service health
echo ""
print_info "Checking service health..."

# Backend health
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    print_status "Backend is healthy"
else
    print_warning "Backend health check failed"
fi

# Frontend health
if curl -f http://localhost:3000 > /dev/null 2>&1; then
    print_status "Frontend is running"
else
    print_warning "Frontend health check failed"
fi

# 8. Show service URLs
echo ""
echo "🌐 Service URLs:"
echo "=================="
echo ""
print_info "Frontend: http://localhost:3000"
print_info "Backend API: http://localhost:8080/api/v1"
print_info "Health Check: http://localhost:8080/health"
echo ""

# 9. Test API endpoints
echo "🧪 Testing API endpoints..."
echo ""

echo "Testing products endpoint..."
curl -s http://localhost:8080/api/v1/products | head -c 100

echo ""
echo "Testing flash products endpoint..."
curl -s http://localhost:8080/api/v1/products/flash | head -c 100

echo ""
echo "Testing hot products endpoint..."
curl -s http://localhost:8080/api/v1/products/hot | head -c 100

echo ""
print_status "🎉 Gemini3-mock + Go Backend Integration deployed successfully!"
echo ""
echo "📝 Next Steps:"
echo "1. Visit http://localhost:3000 to test the full integration"
echo "2. Test AI chat, flash sales, hot products"
echo "3. Verify authentication and cart functionality"
echo "4. Update environment variables in .env.local"
echo "5. Run docker-compose -f docker-compose.gemini.yml logs to debug"
echo ""
echo "🚀 Your Blytz.live marketplace is now running with:"
echo "   • Advanced Gemini3-mock frontend"
echo "   • Enhanced Go backend with flash sales"
echo "   • AI-powered chat assistant"
echo "   • Complete e-commerce functionality"
echo ""