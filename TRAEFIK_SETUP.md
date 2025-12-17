# 🚀 Traefik + Dokploy Setup Guide

## 🎯 **Perfect Minimal Setup**

You've got the right approach! **Traefik + Dokploy** is the cleanest way to deploy. This setup gives you:

✅ **4 Containers Only:**
- Frontend (React/Vite) 
- Backend (Go API)
- PostgreSQL (Database)
- Redis (Cache)

✅ **Automatic SSL** via Traefik
✅ **Automatic Routing** via Traefik
✅ **Dokploy Integration** perfect

## 📦 **Files Created:**

1. **`docker-compose.traefik.yml`** - Minimal 4-service setup
2. **`traefik-deploy.sh`** - Automated deployment script
3. **`.env.traefik`** - Environment template for Traefik

## 🚀 **Quick Deployment:**

```bash
# 1. Clone fresh copy
git clone https://github.com/gmsas95/blytz.live.remake.git
cd blytz.live.remake

# 2. Configure environment
cp .env.traefik .env
nano .env  # Update DOMAIN, passwords, etc.

# 3. Deploy to Dokploy
# In Dokploy, use: docker-compose.traefik.yml

# OR manual deploy:
./traefik-deploy.sh
```

## ⚙️ **Environment Setup (.env):**

```bash
# CRITICAL - Update these:
DOMAIN=blytz.live
POSTGRES_PASSWORD=your_secure_db_password
JWT_SECRET=your_super_secure_jwt_secret_minimum_32_characters
REDIS_PASSWORD=your_secure_redis_password

# Optional - Update if needed:
GEMINI_API_KEY=your_gemini_ai_api_key
SMTP_USER=your_email@gmail.com
SMTP_PASSWORD=your_email_app_password
```

## 🌐 **Traefik URL Mapping:**

| Container | Traefik URL | Internal Port |
|-----------|--------------|---------------|
| Frontend | `https://blytz.live` | 3000 |
| Backend | `https://api.blytz.live` | 8080 |
| PostgreSQL | Internal only | 5432 |
| Redis | Internal only | 6379 |

## 🎯 **What Traefik Handles:**

✅ **Automatic SSL** - Let's Encrypt certificates
✅ **HTTPS Routing** - Main domain → Frontend
✅ **API Routing** - Subdomain → Backend  
✅ **Load Balancing** - Multiple requests
✅ **Health Checks** - Container monitoring
✅ **SSL Termination** - HTTPS to HTTP inside

## 📝 **Dokploy Configuration:**

In your Dokploy dashboard:

1. **Compose File:** `docker-compose.traefik.yml`
2. **Environment:** Use `.env` file content
3. **Traefik:** Pre-configured with labels
4. **Networks:** All on `blytz-network`

## 🔧 **Traefik Labels Explained:**

```yaml
labels:
  - "traefik.enable=true"  # Enable Traefik
  - "traefik.http.routers.frontend.rule=Host(`blytz.live`)"  # Route domain
  - "traefik.http.routers.frontend.entrypoints=websecure"  # HTTPS only
  - "traefik.http.routers.frontend.tls=true"  # SSL enabled
  - "traefik.http.services.frontend.loadbalancer.server.port=3000"  # Container port
```

## 🚀 **Deployment Process:**

### **Step 1: Configure Environment**
```bash
cp .env.traefik .env
nano .env  # Set your DOMAIN and passwords
```

### **Step 2: Deploy via Dokploy**
1. Add repository in Dokploy
2. Select `docker-compose.traefik.yml`
3. Add environment variables
4. Deploy!

### **Step 3: Verify Services**
```bash
# Check container status
docker-compose -f docker-compose.traefik.yml ps

# Check logs
docker-compose -f docker-compose.traefik.yml logs -f
```

## 🎉 **What You Get:**

✅ **Production-ready** Blytz.live marketplace
✅ **Automatic HTTPS** with Let's Encrypt
✅ **Domain routing** - blytz.live + api.blytz.live
✅ **Database persistence** with PostgreSQL
✅ **Performance caching** with Redis
✅ **Cyberpunk design** preserved
✅ **Full e-commerce** functionality
✅ **AI integration** ready
✅ **Email notifications** ready

## 📱 **Final URLs:**

- **Frontend:** `https://blytz.live` 
- **Backend API:** `https://api.blytz.live/api/v1`
- **Health Check:** `https://api.blytz.live/health`

## 🌟 **Perfect Setup Benefits:**

✅ **Minimal containers** - Less overhead
✅ **Traefik managed** - Zero SSL config
✅ **Dokploy integrated** - One-click deployment
✅ **Auto-scaling** - Easy to add replicas
✅ **Health monitored** - Automatic restarts
✅ **Secure by default** - HTTPS only
✅ **Domain ready** - Just set DNS

## 🚀 **Ready for Production!**

This is exactly how modern containerized apps should be deployed:
- **Simple architecture**
- **Managed routing**
- **Automatic security**
- **Easy maintenance**

**🎯 Perfect setup for your Blytz.live marketplace!** 🎊