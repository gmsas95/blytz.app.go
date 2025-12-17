#!/bin/bash

# Development Deployment Script
# For local development with Docker

set -e

echo "🛠️  Blytz.live Development Environment Started"
echo "================================================"

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

# Check .env file
if [ ! -f ".env" ]; then
    print "error" ".env file not found!"
    print "info" "Copy .env.example to .env and fill in your development values"
    exit 1
fi

print "info" "Loading environment variables..."
source .env

# Stop existing development containers
print "info" "Stopping existing development containers..."
docker-compose -f docker-compose.dev.yml down 2>/dev/null || true

# Build and start development containers
print "info" "Building and starting development containers..."
docker-compose -f docker-compose.dev.yml up --build -d

# Wait for containers
print "info" "Waiting for services to start..."
sleep 20

# Check container health
print "info" "Checking service health..."

if docker-compose -f docker-compose.dev.yml ps | grep -q "unhealthy\|Up.*Exit"; then
    print "error" "Some services are not healthy!"
    print "info" "Checking logs..."
    docker-compose -f docker-compose.dev.yml logs --tail=30
    exit 1
fi

# Run migrations
print "info" "Running database migrations..."
docker-compose -f docker-compose.dev.yml exec -T backend /app/blytz-server migrate up || true

# Show development URLs
print "success" "Development environment started successfully!"
echo ""
print "info" "Development URLs:"
echo "======================"
print "success" "Frontend: http://localhost:3000"
print "success" "Backend API: http://localhost:8080"
print "success" "API Docs: http://localhost:8080/docs"

echo ""
print "info" "Development commands:"
print "info" "View logs: docker-compose -f docker-compose.dev.yml logs -f [service]"
print "info" "Stop services: docker-compose -f docker-compose.dev.yml down"
print "info" "Rebuild service: docker-compose -f docker-compose.dev.yml up --build -d [service]"

echo ""
print "success" "🚀 Development environment ready!"
print "warning" "Note: This is for development only. Use docker-compose.yml for production."