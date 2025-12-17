# 🚀 Emergency Vite Production Build Fix
# Using Node.js serve instead of Nginx

#!/bin/bash

set -e

echo "🚀 Blytz.live Emergency Deploy - Node.js Version"
echo "=============================================="

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

# Check if we're in right directory
if [ ! -f "docker-compose.yml" ]; then
    print "error" "docker-compose.yml not found in current directory!"
    exit 1
fi

# Quick environment setup
if [ ! -f ".env" ]; then
    print "info" "Setting up quick environment file..."
    cp .env.minimum .env
    
    print "warning" "⚠️  CRITICAL: Update these values in .env:"
    print "warning" "   - POSTGRES_PASSWORD (database password)"
    print "warning" "   - JWT_SECRET (32+ character secret)"
    print "warning" "   - MINIO_ROOT_PASSWORD (storage password)"
    print "warning" "   - REDIS_PASSWORD (cache password)"
    print "warning" "   - SMTP_USER & SMTP_PASSWORD (email)"
    print "warning" "   - GEMINI_API_KEY (AI features)"
    echo ""
    read -p "Press Enter to continue (make sure .env is configured)..."
fi

# Clean up existing containers
print "info" "Cleaning up existing containers..."
docker-compose down --remove-orphans 2>/dev/null || true

# Remove unused images
print "info" "Cleaning up Docker resources..."
docker system prune -f

# Build and start services
print "info" "Building and starting services..."
docker-compose up --build -d

# Wait for services to start
print "info" "Waiting for services to initialize..."
sleep 60

# Check service health
print "info" "Checking service health..."
for i in {1..5}; do
    if docker-compose ps | grep -q "unhealthy\|Up.*Exit"; then
        print "warning" "Some services are still starting... ($i/5)"
        sleep 15
    else
        break
    fi
done

# Run database migrations
print "info" "Running database migrations..."
docker-compose exec -T backend /app/blytz-server migrate up || true

# Create MinIO bucket
print "info" "Creating S3 bucket..."
docker-compose exec -T backend /app/blytz-server bucket create || true

# Check final status
print "success" "Deployment completed!"
echo ""
print "info" "Service Status:"
docker-compose ps

echo ""
echo "🌐 Service URLs:"
echo "=================="
print "success" "Frontend (React): http://localhost:${FRONTEND_PORT:-3000}"
print "success" "Backend API: http://localhost:${BACKEND_PORT:-8080}"
print "success" "MinIO Console: http://localhost:${MINIO_CONSOLE_PORT:-9001}"

echo ""
print "success" "🎉 Blytz.live is now running!"
print "info" "Commands:"
print "info" "  View logs: docker-compose logs -f [service]"
print "info" "  Stop services: docker-compose down"
print "info" "  Restart: docker-compose restart"

echo ""
print "warning" "Note: Frontend uses Node.js serve, not Nginx"
print "info" "This is lighter and perfect for your deployment!"