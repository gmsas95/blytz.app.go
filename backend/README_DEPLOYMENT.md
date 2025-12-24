# 🚀 Blytz.live Backend Deployment

## Docker Compose Setup

The project is now ready for deployment with Docker Compose. The following services are configured:

### Services
- **postgres**: PostgreSQL database (port 5432)
- **redis**: Redis cache (port 6379)  
- **backend**: Go API server (port 8080)
- **frontend**: React/Vite frontend (port 3000)

## Quick Start

1. **Clone and navigate to backend directory:**
   ```bash
   cd /home/sas/blytz.live.remake/backend
   ```

2. **Configure environment variables:**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Build and run all services:**
   ```bash
   docker compose up -d --build --remove-orphans
   ```

4. **Access the applications:**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Health Check: http://localhost:8080/health
   - API Documentation: http://localhost:8080/api/v1

## Environment Variables

### Backend (.env)
```bash
# Application
ENV=development
PORT=8080
LOG_LEVEL=debug

# Database
DATABASE_URL=postgres://postgres:postgres@postgres:5432/blytz_dev

# Redis
REDIS_URL=redis://redis:6379

# JWT
JWT_SECRET=your-super-secret-jwt-key

# LiveKit (optional)
LIVEKIT_API_KEY=
LIVEKIT_API_SECRET=
LIVEKIT_HOST=

# AWS S3 (optional)
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
AWS_REGION=
AWS_S3_BUCKET=
```

### Frontend
The frontend service is configured to use:
- API URL: http://localhost:8080/api/v1
- Development server: Vite dev server
- Port: 3000

## Database Setup

### Initial Setup
The database will be automatically migrated on first run. A test category will be created if no categories exist.

### Manual Database Access
```bash
# Connect to PostgreSQL container
docker exec -it blytz-postgres psql -U postgres -d blytz_dev

# View tables
\dt

# Connect to Redis
docker exec -it blytz-redis redis-cli
```

## Health Checks

### Backend Health
```bash
# Check backend health
curl http://localhost:8080/health

# Expected response
{
  "status": "ok",
  "database": "connected",
  "redis": "connected",
  "env": "development"
}
```

### Service Status
```bash
# Check all service logs
docker compose logs

# Check specific service logs
docker compose logs backend
docker compose logs frontend
docker compose logs postgres
docker compose logs redis

# Check service status
docker compose ps
```

## API Testing

### Authentication
```bash
# Register user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","username":"testuser"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

### Products
```bash
# List products
curl http://localhost:8080/api/v1/products

# Create product (requires auth)
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"name":"Test Product","description":"A test product","price":99.99}'
```

### Auctions
```bash
# List auctions
curl http://localhost:8080/api/v1/auctions

# Create auction (requires auth)
curl -X POST http://localhost:8080/api/v1/auctions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"product_id":"product-uuid","starting_bid":10.0,"reserve_price":50.0,"end_time":"2024-12-31T23:59:59Z"}'
```

## Development Workflow

### Backend Development
```bash
# Run with hot reload (development)
docker compose up --build backend

# View backend logs
docker compose logs -f backend
```

### Frontend Development
```bash
# Run frontend with hot reload
docker compose up --build frontend

# View frontend logs
docker compose logs -f frontend
```

### Database Migrations
```bash
# Database migrations are handled automatically by GORM AutoMigrate
# For custom migrations, see internal/database/migration.go
```

## Troubleshooting

### Common Issues

1. **Port Conflicts**
   ```bash
   # Check if ports are in use
   netstat -tulpn | grep :8080
   netstat -tulpn | grep :3000
   
   # Kill processes using ports
   sudo kill -9 PID
   ```

2. **Database Connection Issues**
   ```bash
   # Check PostgreSQL container status
   docker compose ps postgres
   
   # Restart PostgreSQL
   docker compose restart postgres
   
   # Check PostgreSQL logs
   docker compose logs postgres
   ```

3. **Redis Connection Issues**
   ```bash
   # Check Redis container status
   docker compose ps redis
   
   # Restart Redis
   docker compose restart redis
   
   # Check Redis logs
   docker compose logs redis
   ```

4. **Build Errors**
   ```bash
   # Clean build cache
   docker compose down --volumes --rmi all
   
   # Rebuild from scratch
   docker compose build --no-cache
   ```

### Logs and Debugging

```bash
# View all logs
docker compose logs

# Follow logs in real-time
docker compose logs -f

# View logs for specific service
docker compose logs -f backend

# View recent logs (last 50 lines)
docker compose logs --tail=50 backend

# Get container shell for debugging
docker exec -it blytz-backend sh
docker exec -it blytz-frontend sh
```

## Production Deployment

### Security Considerations
1. Change default passwords and secrets
2. Use HTTPS in production
3. Set up proper CORS origins
4. Enable rate limiting
5. Configure firewall rules
6. Use environment-specific configurations

### Performance Optimization
1. Enable Redis for caching
2. Configure database connection pooling
3. Set up reverse proxy (nginx/traefik)
4. Enable gzip compression
5. Monitor resource usage

### Monitoring
1. Set up application monitoring
2. Configure log aggregation
3. Set up health check alerts
4. Monitor database performance
5. Track API response times

## Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend     │    │     Backend     │    │   Database      │
│   (React)      │◄──►│     (Go)       │◄──►│  (PostgreSQL)   │
│   Port: 3000   │    │   Port: 8080   │    │   Port: 5432    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │      Redis      │
                       │     Cache       │
                       │   Port: 6379   │
                       └─────────────────┘
```

## Next Steps

1. Configure production environment variables
2. Set up SSL certificates
3. Configure reverse proxy
4. Set up monitoring and logging
5. Deploy to staging environment
6. Perform load testing
7. Deploy to production

## Support

For issues and questions:
1. Check container logs: `docker compose logs`
2. Verify environment variables
3. Check network connectivity between containers
4. Review health endpoints
5. Consult implementation documentation: `IMPLEMENTATION_SUMMARY.md`