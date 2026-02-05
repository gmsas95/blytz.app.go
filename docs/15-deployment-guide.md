# Deployment Guide

## Overview

Production deployment procedures for the Blytz platform.

## Infrastructure

### Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                           Cloudflare                                │
│  - DNS (blytz.app)                                                 │
│  - CDN (Static assets)                                             │
│  - DDoS Protection                                                 │
│  - SSL/TLS                                                         │
└─────────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                         Kubernetes Cluster                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │
│  │  Frontend    │  │   Backend    │  │   Mobile     │              │
│  │  (Next.js)   │  │   (Go/Bun)   │  │   API        │              │
│  │  Pods: 3     │  │   Pods: 3    │  │   (Optional) │              │
│  └──────────────┘  └──────────────┘  └──────────────┘              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │
│  │  LiveKit     │  │   Socket.io  │  │   Ingress    │              │
│  │  Server      │  │   Gateway    │  │   Controller │              │
│  └──────────────┘  └──────────────┘  └──────────────┘              │
└─────────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                           Data Layer                                │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │
│  │  PostgreSQL  │  │    Redis     │  │   Cloudflare │              │
│  │  (Primary +  │  │   Cluster    │  │      R2      │              │
│  │   Replica)   │  │              │  │              │              │
│  └──────────────┘  └──────────────┘  └──────────────┘              │
└─────────────────────────────────────────────────────────────────────┘
```

## Prerequisites

### Tools
- `kubectl` - Kubernetes CLI
- `helm` - Kubernetes package manager
- `docker` - Container runtime
- `terraform` - Infrastructure as code
- `doctl` - DigitalOcean CLI (or AWS/GCP equivalent)

### Accounts
- Cloudflare account
- DigitalOcean / AWS / GCP account
- Docker Hub / GitHub Container Registry

## Environment Setup

### 1. Create Kubernetes Cluster

**DigitalOcean:**
```bash
# Create cluster
doctl kubernetes cluster create blytz-prod \
  --region sgp1 \
  --version 1.29 \
  --node-pool "name=worker;size=s-4vcpu-8gb;count=3;auto-scale=true;min-nodes=3;max-nodes=10"

# Get kubeconfig
doctl kubernetes cluster kubeconfig save blytz-prod
```

**AWS EKS:**
```bash
# Using eksctl
eksctl create cluster \
  --name blytz-prod \
  --region ap-southeast-1 \
  --node-type t3.medium \
  --nodes 3 \
  --nodes-min 3 \
  --nodes-max 10
```

### 2. Install Ingress Controller

```bash
# NGINX Ingress Controller
helm upgrade --install ingress-nginx ingress-nginx \
  --repo https://kubernetes.github.io/ingress-nginx \
  --namespace ingress-nginx \
  --create-namespace

# Get load balancer IP
kubectl get service ingress-nginx-controller -n ingress-nginx
```

### 3. Install Cert Manager (SSL)

```bash
helm upgrade --install cert-manager cert-manager \
  --repo https://charts.jetstack.io \
  --namespace cert-manager \
  --create-namespace \
  --set installCRDs=true

# Create ClusterIssuer for Let's Encrypt
kubectl apply -f - <<EOF
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: admin@blytz.app
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - http01:
        ingress:
          class: nginx
EOF
```

### 4. Install PostgreSQL

```bash
# Using Helm with Bitnami chart
helm upgrade --install postgres bitnami/postgresql \
  --namespace database \
  --create-namespace \
  --set global.postgresql.auth.password=YOUR_SECURE_PASSWORD \
  --set global.postgresql.auth.database=blytz \
  --set architecture=replication \
  --set auth.replicationPassword=REPLICATION_PASSWORD \
  --set primary.persistence.size=100Gi \
  --set readReplicas.replicaCount=1

# Get connection details
kubectl get secret --namespace database postgres-postgresql -o jsonpath="{.data.password}" | base64 -d
```

### 5. Install Redis

```bash
helm upgrade --install redis bitnami/redis \
  --namespace database \
  --set global.redis.password=YOUR_REDIS_PASSWORD \
  --set architecture=standalone \
  --set master.persistence.size=20Gi
```

## Application Deployment

### 1. Create Secrets

```bash
# Create namespace
kubectl create namespace blytz

# Create secrets from env file
kubectl create secret generic blytz-secrets \
  --namespace blytz \
  --from-env-file=.env.production

# Or create individually
kubectl create secret generic db-credentials \
  --namespace blytz \
  --from-literal=url=postgresql://...

kubectl create secret generic jwt-secret \
  --namespace blytz \
  --from-literal=secret=YOUR_JWT_SECRET

kubectl create secret generic stripe-credentials \
  --namespace blytz \
  --from-literal=secret_key=sk_live_... \
  --from-literal=webhook_secret=whsec_...
```

### 2. Backend Deployment

```yaml
# k8s/backend-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: blytz-backend
  namespace: blytz
spec:
  replicas: 3
  selector:
    matchLabels:
      app: blytz-backend
  template:
    metadata:
      labels:
        app: blytz-backend
    spec:
      containers:
      - name: backend
        image: ghcr.io/gmsas95/blytz-backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
        - name: ENV
          value: "production"
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: blytz-secrets
              key: DATABASE_URL
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: blytz-secrets
              key: JWT_SECRET
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: blytz-backend
  namespace: blytz
spec:
  selector:
    app: blytz-backend
  ports:
  - port: 80
    targetPort: 8080
```

### 3. Frontend Deployment

```yaml
# k8s/frontend-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: blytz-frontend
  namespace: blytz
spec:
  replicas: 3
  selector:
    matchLabels:
      app: blytz-frontend
  template:
    metadata:
      labels:
        app: blytz-frontend
    spec:
      containers:
      - name: frontend
        image: ghcr.io/gmsas95/blytz-frontend:latest
        ports:
        - containerPort: 3000
        env:
        - name: NEXT_PUBLIC_API_URL
          value: "https://api.blytz.app/api/v1"
        - name: NEXT_PUBLIC_SOCKET_URL
          value: "https://api.blytz.app"
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
---
apiVersion: v1
kind: Service
metadata:
  name: blytz-frontend
  namespace: blytz
spec:
  selector:
    app: blytz-frontend
  ports:
  - port: 80
    targetPort: 3000
```

### 4. Ingress Configuration

```yaml
# k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: blytz-ingress
  namespace: blytz
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
spec:
  tls:
  - hosts:
    - blytz.app
    - api.blytz.app
    secretName: blytz-tls
  rules:
  - host: blytz.app
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: blytz-frontend
            port:
              number: 80
  - host: api.blytz.app
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: blytz-backend
            port:
              number: 80
```

## CI/CD Pipeline

### GitHub Actions Workflow

```yaml
# .github/workflows/deploy.yml
name: Deploy to Production

on:
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Run backend tests
        working-directory: backend
        run: |
          go test ./... -cover
      
      - name: Run frontend tests
        working-directory: frontend
        run: |
          npm ci
          npm test

  build:
    needs: test
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3
      
      - name: Login to Container Registry
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      
      - name: Build and push backend
        uses: docker/build-push-action@v5
        with:
          context: ./backend
          push: true
          tags: |
            ghcr.io/gmsas95/blytz-backend:${{ github.sha }}
            ghcr.io/gmsas95/blytz-backend:latest
          cache-from: type=gha
          cache-to: type=gha,mode=max
      
      - name: Build and push frontend
        uses: docker/build-push-action@v5
        with:
          context: ./frontend
          push: true
          tags: |
            ghcr.io/gmsas95/blytz-frontend:${{ github.sha }}
            ghcr.io/gmsas95/blytz-frontend:latest
          cache-from: type=gha
          cache-to: type=gha,mode=max

  deploy:
    needs: build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Configure kubectl
        uses: azure/setup-kubectl@v3
      
      - name: Set up Kustomize
        run: |
          curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh" | bash
          sudo mv kustomize /usr/local/bin/
      
      - name: Configure kubeconfig
        run: |
          echo "${{ secrets.KUBECONFIG }}" | base64 -d > kubeconfig
          export KUBECONFIG=kubeconfig
      
      - name: Deploy to Kubernetes
        run: |
          export KUBECONFIG=kubeconfig
          
          # Update image tags
          cd k8s
          kustomize edit set image \
            backend=ghcr.io/gmsas95/blytz-backend:${{ github.sha }} \
            frontend=ghcr.io/gmsas95/blytz-frontend:${{ github.sha }}
          
          # Apply
          kustomize build . | kubectl apply -f -
          
          # Wait for rollout
          kubectl rollout status deployment/blytz-backend -n blytz
          kubectl rollout status deployment/blytz-frontend -n blytz
      
      - name: Verify deployment
        run: |
          export KUBECONFIG=kubeconfig
          kubectl get pods -n blytz
          kubectl get ingress -n blytz
```

## Database Migrations

### Running Migrations

```bash
# Using golang-migrate
migrate -path ./migrations \
  -database "${DATABASE_URL}" \
  up

# Rollback
migrate -path ./migrations \
  -database "${DATABASE_URL}" \
  down 1
```

### Migration Job

```yaml
# k8s/migration-job.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: db-migration
  namespace: blytz
spec:
  template:
    spec:
      containers:
      - name: migrate
        image: migrate/migrate:latest
        command:
        - migrate
        - -path=/migrations
        - -database=$(DATABASE_URL)
        - up
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: blytz-secrets
              key: DATABASE_URL
        volumeMounts:
        - name: migrations
          mountPath: /migrations
      volumes:
      - name: migrations
        configMap:
          name: db-migrations
      restartPolicy: Never
  backoffLimit: 3
```

## Monitoring

### Install Prometheus & Grafana

```bash
helm upgrade --install prometheus prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  --create-namespace \
  --set grafana.enabled=true \
  --set grafana.adminPassword=ADMIN_PASSWORD

# Port forward to access Grafana
kubectl port-forward svc/prometheus-grafana 3000:80 -n monitoring
```

### Application Metrics

```go
// Add to backend
import "github.com/prometheus/client_golang/prometheus"

var (
  httpRequestsTotal = prometheus.NewCounterVec(
    prometheus.CounterOpts{
      Name: "http_requests_total",
      Help: "Total HTTP requests",
    },
    []string{"method", "endpoint", "status"},
  )
  
  bidCounter = prometheus.NewCounter(
    prometheus.CounterOpts{
      Name: "auctions_bids_total",
      Help: "Total bids placed",
    },
  )
)

func init() {
  prometheus.MustRegister(httpRequestsTotal)
  prometheus.MustRegister(bidCounter)
}
```

## Rollback Procedure

### Quick Rollback
```bash
# Rollback to previous version
kubectl rollout undo deployment/blytz-backend -n blytz
kubectl rollout undo deployment/blytz-frontend -n blytz

# Check status
kubectl rollout status deployment/blytz-backend -n blytz
```

### Database Rollback
```bash
# Rollback migrations
migrate -path ./migrations -database "$DATABASE_URL" down 1
```

## Disaster Recovery

### Backup Strategy
```bash
# Database backup
cronjob: 0 2 * * * pg_dump $DATABASE_URL | gzip > backup_$(date +%Y%m%d).sql.gz

# Upload to S3/R2
aws s3 cp backup_*.sql.gz s3://blytz-backups/database/
```

### Recovery Procedure
```bash
# Restore database
gunzip < backup_20250205.sql.gz | psql $DATABASE_URL

# Redeploy
kubectl apply -f k8s/
```

---

*Last updated: 2025-02-05*
