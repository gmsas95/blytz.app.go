# ==============================
# DEPLOYMENT INSTRUCTIONS
# ==============================

# Prerequisites:
# - Ubuntu 20.04+ or CentOS 8+ VPS
# - Docker & Docker Compose installed
# - Domain pointing to VPS IP

# ==============================
# 1. QUICK DEPLOYMENT
# ==============================

# Clone repository
git clone https://github.com/gmsas95/blytz.live.remake.git
cd blytz.live.remake

# Configure environment
cp .env.example .env
nano .env  # Edit with your actual values

# Deploy everything
./deploy.sh

# ==============================
# 2. MANUAL DEPLOYMENT (if script fails)
# ==============================

# Build and start services
docker-compose up --build -d

# Check logs
docker-compose logs -f

# ==============================
# 3. DOMAIN CONFIGURATION
# ==============================

# DNS Records (point to your VPS IP):
# A     blytz.live         -> YOUR_VPS_IP
# A     www.blytz.live      -> YOUR_VPS_IP
# A     api.blytz.live      -> YOUR_VPS_IP

# ==============================
# 4. SSL CERTIFICATES (Let's Encrypt)
# ==============================

# Install Certbot
sudo apt update
sudo apt install certbot python3-certbot-nginx

# Get SSL certificate
sudo certbot --nginx -d blytz.live -d www.blytz.live

# Auto-renewal
sudo crontab -e
# Add: 0 12 * * * /usr/bin/certbot renew --quiet

# ==============================
# 5. NGINX REVERSE PROXY (Optional)
# ==============================

# Copy production nginx config
sudo cp nginx/prod.conf /etc/nginx/sites-available/blytz.live

# Enable site
sudo ln -s /etc/nginx/sites-available/blytz.live /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx

# ==============================
# 6. FIREWALL CONFIGURATION
# ==============================

# Open required ports
sudo ufw allow 22/tcp      # SSH
sudo ufw allow 80/tcp      # HTTP
sudo ufw allow 443/tcp     # HTTPS
sudo ufw allow 8080/tcp    # Backend (optional)
sudo ufw enable

# ==============================
# 7. MONITORING SETUP
# ==============================

# Check container status
docker-compose ps

# Check logs
docker-compose logs backend
docker-compose logs frontend

# Scale services
docker-compose up --scale backend=2 --scale frontend=2

# ==============================
# 8. BACKUP CONFIGURATION
# ==============================

# Database backup
docker-compose exec postgres pg_dump -U blytz_user blytz_marketplace > backup_$(date +%Y%m%d).sql

# Volume backup
sudo tar -czf blytz_backup_$(date +%Y%m%d).tar.gz /var/lib/docker/volumes/

# ==============================
# 9. TROUBLESHOOTING
# ==============================

# Restart services
docker-compose restart [service-name]

# Rebuild service
docker-compose up --build -d [service-name]

# Clear Docker cache
docker system prune -a

# Check port conflicts
sudo netstat -tlnp | grep :80

# Check disk space
df -h

# Check memory usage
free -h

# ==============================
# 10. PERFORMANCE TUNING
# ==============================

# Optimize PostgreSQL
docker-compose exec postgres psql -U blytz_user -d blytz_marketplace -c "ALTER SYSTEM SET shared_buffers = '256MB';"
docker-compose exec postgres psql -U blytz_user -d blytz_marketplace -c "ALTER SYSTEM SET effective_cache_size = '1GB';"

# Optimize Redis
docker-compose exec redis redis-cli CONFIG SET maxmemory 512mb
docker-compose exec redis redis-cli CONFIG SET maxmemory-policy allkeys-lru

# Container resource limits
# Edit docker-compose.yml to add:
# deploy:
#   resources:
#     limits:
#       cpus: '1.0'
#       memory: 512M
#     reservations:
#       cpus: '0.5'
#       memory: 256M

# ==============================
# 11. SECURITY HARDENING
# ==============================

# Update system
sudo apt update && sudo apt upgrade -y

# Fail2ban
sudo apt install fail2ban
sudo systemctl enable fail2ban
sudo cp /etc/fail2ban/jail.conf /etc/fail2ban/jail.local

# Security headers (already in nginx/prod.conf)
# - HSTS
# - CSP
# - X-Frame-Options
# - X-Content-Type-Options

# Docker security
# - Use non-root users (configured)
# - Read-only filesystems where possible
# - Resource limits

# ==============================
# 12. MONITORING DASHBOARD
# ==============================

# Install Portainer (optional)
docker volume create portainer_data
docker run -d -p 8000:8000 --name portainer --restart always -v /var/run/docker.sock:/var/run/docker.sock -v portainer_data:/data portainer/portainer-ce:latest

# Install Node Exporter for Prometheus (optional)
docker run -d --net="host" --pid="host" -v "/:/host:ro,rslave" quay.io/prometheus/node-exporter

# ==============================
# REQUIRED ENVIRONMENT VARIABLES
# ==============================

# MINIMUM REQUIRED FOR PRODUCTION:
# - POSTGRES_PASSWORD (strong password)
# - JWT_SECRET (32+ character secret)
# - MINIO_ROOT_PASSWORD (strong password)
# - REDIS_PASSWORD (strong password)
# - SMTP_USER & SMTP_PASSWORD (email credentials)
# - CORS_ORIGINS (your domain)
# - VITE_API_URL (your domain API URL)

# ==============================
# SERVICE URLs AFTER DEPLOYMENT
# ==============================

# Frontend: https://blytz.live
# Backend API: https://blytz.live/api/v1
# MinIO Console: https://blytz.live:9001
# Health Check: https://blytz.live/health

# ==============================
# DEPLOYMENT CHECKLIST
# ==============================

□ VPS specifications (4GB+ RAM, 50GB+ SSD)
□ Domain DNS configured
□ Firewall ports open (22, 80, 443)
□ SSL certificates installed
□ Environment variables configured
□ Database connection working
□ Redis connection working
□ MinIO storage working
□ Frontend loading
□ Backend API responding
□ Health checks passing
□ Monitoring setup
□ Backup strategy configured
□ Security hardening complete

# ==============================
# CONTACT & SUPPORT
# ==============================

# Issues: https://github.com/gmsas95/blytz.live.remake/issues
# Documentation: Check repository README files
# Logs: docker-compose logs [service-name]