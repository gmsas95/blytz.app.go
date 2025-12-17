# 🚀 Traefik + Dokploy Deployment
# Frontend + Backend + Database + Cache only

#!/bin/bash

set -e

echo "🚀 Blytz.live Traefik + Dokploy Deployment"
echo "============================================"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print() {
    case $1 in
        "error") echo -e "${RED}❌ $2${NC}" ;;
        "success") echo -e "${GREEN}✅ $2${NC}" ;;
        "warning") echo -e "${YELLOW}⚠️  $2${NC}" ;;
        "info") echo -e "${BLUE}ℹ️  $2${NC}" ;;
        *) echo -e "$2" ;;
    esac
}

# Check if docker-compose.traefik.yml exists
if [ ! -f "docker-compose.traefik.yml" ]; then
    print "error" "docker-compose.traefik.yml not found!"
    exit 1
fi

# Environment setup
if [ ! -f ".env" ]; then
    print "info" "Creating environment file..."
    cp .env.minimum .env
    
    print "warning" "⚠️  CRITICAL: Update these values in .env:"
    print "warning" "   - DOMAIN (your domain: blytz.live)"
    print "warning" "   - POSTGRES_PASSWORD (database password)"
    print "warning" "   - JWT_SECRET (32+ character secret)"
    print "warning" "   - REDIS_PASSWORD (cache password)"
    print "warning" "   - GEMINI_API_KEY (AI features)"
    echo ""
    print "info" "Edit .env file now with your values:"
    print "info" "nano .env"
    echo ""
    read -p "Press Enter to continue (make sure .env is configured)..."
fi

# Deploy with Traefik
print "info" "Deploying with Traefik configuration..."
docker-compose -f docker-compose.traefik.yml down --remove-orphans 2>/dev/null || true

# Build and deploy
print "info" "Building and starting services..."
docker-compose -f docker-compose.traefik.yml up --build -d

# Wait for services
print "info" "Waiting for services to initialize..."
sleep 45

# Check health
print "info" "Checking service health..."
for i in {1..5}; do
    if docker-compose -f docker-compose.traefik.yml ps | grep -q "unhealthy\|Up.*Exit"; then
        print "warning" "Some services are still starting... ($i/5)"
        sleep 15
    else
        break
    fi
done

# Run migrations
print "info" "Running database migrations..."
docker-compose -f docker-compose.traefik.yml exec -T backend /app/blytz-server migrate up || true

# Check status
print "success" "Traefik deployment completed!"
echo ""
print "info" "Service Status:"
docker-compose -f docker-compose.traefik.yml ps

echo ""
echo "🌐 Service URLs (via Traefik):"
echo "================================"
print "success" "Frontend: https://${DOMAIN:-blytz.live}"
print "success" "Backend API: https://api.${DOMAIN:-blytz.live}"

echo ""
print "success" "🎉 Traefik deployment successful!"
print "info" "Traefik will automatically:"
print "info" "  - Handle SSL certificates"
print "info" "  - Route traffic to containers"
print "info" "  - Load balance requests"
print "info" "  - Provide HTTPS termination"

echo ""
print "info" "Commands:"
print "info" "  View logs: docker-compose -f docker-compose.traefik.yml logs -f [service]"
print "info" "  Stop services: docker-compose -f docker-compose.traefik.yml down"
print "info" "  Restart: docker-compose -f docker-compose.traefik.yml restart"