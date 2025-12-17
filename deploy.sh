#!/bin/bash

# Blytz.live Deployment Script
# Production deployment for VPS

set -e  # Exit on any error

echo "🚀 Blytz.live Marketplace Deployment Started"
echo "================================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print() {
    case $1 in
        "error") echo -e "${RED}❌ $2${NC}" ;;
        "success") echo -e "${GREEN}✅ $2${NC}" ;;
        "warning") echo -e "${YELLOW}⚠️  $2${NC}" ;;
        "info") echo -e "${BLUE}ℹ️  $2${NC}" ;;
        *) echo -e "$2" ;;
    esac
}

# Check if .env file exists
if [ ! -f ".env" ]; then
    print "error" ".env file not found!"
    print "info" "Please copy .env.example to .env and fill in your values"
    exit 1
fi

print "info" "Loading environment variables..."
source .env

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    print "error" "Docker is not installed!"
    print "info" "Please install Docker first: https://docs.docker.com/get-docker/"
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    print "error" "Docker Compose is not installed!"
    print "info" "Please install Docker Compose first: https://docs.docker.com/compose/install/"
    exit 1
fi

# Stop existing containers
print "info" "Stopping existing containers..."
docker-compose down 2>/dev/null || true

# Pull latest images
print "info" "Pulling latest Docker images..."
docker-compose pull

# Build and start containers
print "info" "Building and starting containers..."
docker-compose up --build -d

# Wait for containers to be healthy
print "info" "Waiting for services to be healthy..."
sleep 30

# Check container health
print "info" "Checking service health..."

if docker-compose ps | grep -q "unhealthy\|Up.*Exit"; then
    print "error" "Some services are not healthy!"
    print "info" "Checking logs..."
    docker-compose logs --tail=50
    exit 1
fi

# Run database migrations
print "info" "Running database migrations..."
docker-compose exec -T backend /app/blytz-server migrate up || true

# Create MinIO buckets
print "info" "Creating S3 buckets..."
docker-compose exec -T backend /app/blytz-server bucket create || true

# Show running containers
print "success" "Deployment completed successfully!"
echo ""
print "info" "Running containers:"
docker-compose ps

echo ""
echo "🌐 Service URLs:"
echo "=================="
print "success" "Frontend: http://localhost:${FRONTEND_PORT:-80}"
print "success" "Backend API: http://localhost:${BACKEND_PORT:-8080}"
print "success" "MinIO Console: http://localhost:${MINIO_CONSOLE_PORT:-9001}"
print "success" "Redis: localhost:${REDIS_PORT:-6379}"
print "success" "PostgreSQL: localhost:${POSTGRES_PORT:-5432}"

echo ""
print "info" "Default credentials:"
print "info" "MinIO: ${MINIO_ROOT_USER:-minioadmin} / ${MINIO_ROOT_PASSWORD}"
print "info" "PostgreSQL: ${POSTGRES_USER:-blytz_user} / [your password]"
print "info" "Redis: Password: ${REDIS_PASSWORD}"

echo ""
print "success" "🎉 Blytz.live Marketplace is now running!"
print "info" "Check logs with: docker-compose logs -f [service-name]"
print "info" "Stop services with: docker-compose down"
print "info" "Restart services with: docker-compose restart"

echo ""
print "warning" "Don't forget to:"
print "warning" "1. Configure your domain DNS"
print "warning" "2. Set up SSL certificates (Let's Encrypt recommended)"
print "warning" "3. Configure environment variables for production"
print "warning" "4. Set up monitoring and backups"

echo ""
print "info" "Deployment script completed! 🚀"