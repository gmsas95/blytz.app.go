# Permissions & RBAC

## Overview

Role-Based Access Control (RBAC) system for Blytz platform.

## User Roles

| Role | Description | Typical Users |
|------|-------------|---------------|
| **buyer** | Can browse, bid, purchase | General public |
| **seller** | Can list products, host streams | Business owners |
| **admin** | Full platform access | Platform staff |

## Permission Matrix

### Feature Access by Role

| Feature | Buyer | Seller | Admin |
|---------|-------|--------|-------|
| **Authentication** ||||
| Register/Login | ✅ | ✅ | ✅ |
| Update Profile | Own | Own | Any |
| Delete Account | Own | Own | Any |
| **Product Catalog** ||||
| View Products | ✅ | ✅ | ✅ |
| Create Product | ❌ | ✅ | ✅ |
| Edit Product | ❌ | Own | Any |
| Delete Product | ❌ | Own | Any |
| **Auctions** ||||
| View Auctions | ✅ | ✅ | ✅ |
| Create Auction | ❌ | ✅ | ✅ |
| Edit Auction | ❌ | Own | Any |
| Cancel Auction | ❌ | Own (if no bids) | Any |
| Place Bid | ✅ | ❌ | ❌ |
| **Streams** ||||
| View Streams | ✅ | ✅ | ✅ |
| Start Stream | ❌ | ✅ | ✅ |
| End Stream | ❌ | Own | Any |
| Moderate Chat | ❌ | Own stream | Any |
| **Orders** ||||
| View Orders | Own | Own | Any |
| Create Order | ✅ | ❌ | ❌ |
| Update Order Status | ❌ | Own (ship) | Any |
| Cancel Order | Own (if pending) | Own (if pending) | Any |
| **Payments** ||||
| Make Payment | ✅ | ❌ | ❌ |
| View Payment History | Own | Own | Any |
| Process Refund | ❌ | ❌ | ✅ |
| **Shipping** ||||
| Create Shipment | ❌ | Own orders | Any |
| Print Label | ❌ | Own orders | Any |
| Track Shipment | ✅ (own) | ✅ (own) | Any |
| **Reviews** ||||
| Write Review | ✅ (after purchase) | ❌ | ❌ |
| Respond to Review | ❌ | Own reviews | Any |
| Delete Review | Own | ❌ | ✅ |
| **Admin** ||||
| Verify Sellers | ❌ | ❌ | ✅ |
| Suspend Users | ❌ | ❌ | ✅ |
| View Analytics | ❌ | Own | Platform |
| Manage Categories | ❌ | ❌ | ✅ |
| Moderate Content | ❌ | ❌ | ✅ |

## Permission Definitions

### Permission Format
```typescript
interface Permission {
  resource: string;      // 'product', 'auction', 'order', etc.
  action: string;        // 'create', 'read', 'update', 'delete'
  scope: 'own' | 'any';  // 'own' = only own resources, 'any' = all resources
}
```

### Role Permissions

#### Buyer Role
```typescript
const buyerPermissions = [
  { resource: 'user', action: 'read', scope: 'own' },
  { resource: 'user', action: 'update', scope: 'own' },
  { resource: 'user', action: 'delete', scope: 'own' },
  { resource: 'product', action: 'read', scope: 'any' },
  { resource: 'auction', action: 'read', scope: 'any' },
  { resource: 'auction', action: 'bid', scope: 'any' },
  { resource: 'stream', action: 'read', scope: 'any' },
  { resource: 'chat', action: 'create', scope: 'any' },
  { resource: 'order', action: 'create', scope: 'own' },
  { resource: 'order', action: 'read', scope: 'own' },
  { resource: 'order', action: 'update', scope: 'own' }, // cancel own order
  { resource: 'payment', action: 'create', scope: 'own' },
  { resource: 'payment', action: 'read', scope: 'own' },
  { resource: 'shipment', action: 'read', scope: 'own' },
  { resource: 'review', action: 'create', scope: 'own' },
  { resource: 'review', action: 'update', scope: 'own' },
  { resource: 'seller', action: 'read', scope: 'any' },
  { resource: 'seller', action: 'follow', scope: 'any' },
  { resource: 'notification', action: 'read', scope: 'own' },
  { resource: 'notification', action: 'update', scope: 'own' },
];
```

#### Seller Role
```typescript
const sellerPermissions = [
  // Inherits buyer permissions except bidding
  { resource: 'user', action: 'read', scope: 'own' },
  { resource: 'user', action: 'update', scope: 'own' },
  { resource: 'product', action: 'create', scope: 'own' },
  { resource: 'product', action: 'read', scope: 'any' },
  { resource: 'product', action: 'update', scope: 'own' },
  { resource: 'product', action: 'delete', scope: 'own' },
  { resource: 'auction', action: 'create', scope: 'own' },
  { resource: 'auction', action: 'read', scope: 'any' },
  { resource: 'auction', action: 'update', scope: 'own' },
  { resource: 'auction', action: 'delete', scope: 'own' }, // if no bids
  { resource: 'stream', action: 'create', scope: 'own' },
  { resource: 'stream', action: 'read', scope: 'any' },
  { resource: 'stream', action: 'update', scope: 'own' },
  { resource: 'stream', action: 'delete', scope: 'own' },
  { resource: 'chat', action: 'moderate', scope: 'own' }, // own stream
  { resource: 'order', action: 'read', scope: 'own' },
  { resource: 'order', action: 'update', scope: 'own' }, // update status
  { resource: 'payment', action: 'read', scope: 'own' },
  { resource: 'shipment', action: 'create', scope: 'own' },
  { resource: 'shipment', action: 'read', scope: 'own' },
  { resource: 'shipment', action: 'print_label', scope: 'own' },
  { resource: 'review', action: 'read', scope: 'own' },
  { resource: 'review', action: 'respond', scope: 'own' },
  { resource: 'seller', action: 'read', scope: 'any' },
  { resource: 'analytics', action: 'read', scope: 'own' },
];
```

#### Admin Role
```typescript
const adminPermissions = [
  { resource: '*', action: '*', scope: 'any' }, // Full access
];
```

## Middleware Implementation

### Backend (Go/Bun)
```go
// middleware/authorization.go

func RequireAuth() gin.HandlerFunc {
  return func(c *gin.Context) {
    token := c.GetHeader("Authorization")
    if token == "" {
      c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
      return
    }
    
    // Validate JWT and set user context
    user, err := validateToken(token)
    if err != nil {
      c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
      return
    }
    
    c.Set("user_id", user.ID)
    c.Set("user_role", user.Role)
    c.Next()
  }
}

func RequireRole(roles ...string) gin.HandlerFunc {
  return func(c *gin.Context) {
    userRole := c.GetString("user_role")
    
    for _, role := range roles {
      if userRole == role {
        c.Next()
        return
      }
    }
    
    c.AbortWithStatusJSON(403, gin.H{"error": "Forbidden"})
  }
}

func RequirePermission(resource, action string) gin.HandlerFunc {
  return func(c *gin.Context) {
    userID := c.GetString("user_id")
    userRole := c.GetString("user_role")
    
    // Check if user has permission
    hasPermission := checkPermission(userRole, resource, action, userID, c)
    if !hasPermission {
      c.AbortWithStatusJSON(403, gin.H{"error": "Insufficient permissions"})
      return
    }
    
    c.Next()
  }
}
```

### Route Protection
```go
// routes setup
func SetupRoutes(r *gin.Engine) {
  api := r.Group("/api/v1")
  
  // Public routes
  api.POST("/auth/login", authHandler.Login)
  api.POST("/auth/register", authHandler.Register)
  api.GET("/products", productHandler.List)
  api.GET("/auctions", auctionHandler.List)
  api.GET("/streams", streamHandler.List)
  
  // Protected routes
  auth := api.Group("/")
  auth.Use(middleware.RequireAuth())
  {
    // Buyer + Seller routes
    auth.GET("/auth/me", userHandler.Me)
    auth.PUT("/users/me", userHandler.Update)
    
    // Seller only
    seller := auth.Group("/")
    seller.Use(middleware.RequireRole("seller", "admin"))
    {
      seller.POST("/products", productHandler.Create)
      seller.PUT("/products/:id", productHandler.Update)
      seller.DELETE("/products/:id", productHandler.Delete)
      seller.POST("/auctions", auctionHandler.Create)
      seller.POST("/streams", streamHandler.Create)
    }
    
    // Buyer only (can't bid as seller)
    buyer := auth.Group("/")
    buyer.Use(middleware.RequireRole("buyer"))
    {
      buyer.POST("/auctions/:id/bids", auctionHandler.PlaceBid)
      buyer.POST("/orders", orderHandler.Create)
    }
    
    // Admin only
    admin := auth.Group("/admin")
    admin.Use(middleware.RequireRole("admin"))
    {
      admin.GET("/users", adminHandler.ListUsers)
      admin.POST("/sellers/:id/verify", adminHandler.VerifySeller)
      admin.POST("/users/:id/suspend", adminHandler.SuspendUser)
    }
  }
}
```

### Frontend (Next.js)
```typescript
// hooks/use-permissions.ts
export function usePermissions() {
  const { user } = useAuth();
  
  return {
    can: (resource: string, action: string, ownerId?: string) => {
      if (!user) return false;
      
      // Admin can do anything
      if (user.role === 'admin') return true;
      
      // Check specific permissions
      const permission = getPermission(user.role, resource, action);
      if (!permission) return false;
      
      // Check scope
      if (permission.scope === 'own') {
        return ownerId === user.id;
      }
      
      return true;
    },
    
    isRole: (...roles: string[]) => roles.includes(user?.role || ''),
  };
}

// components/guard/require-auth.tsx
export function RequireAuth({ 
  children, 
  fallback = <LoginPrompt /> 
}: RequireAuthProps) {
  const { isAuthenticated } = useAuth();
  
  if (!isAuthenticated) return fallback;
  return children;
}

// components/guard/require-role.tsx
export function RequireRole({ 
  roles, 
  children, 
  fallback = <Forbidden /> 
}: RequireRoleProps) {
  const { user } = useAuth();
  
  if (!roles.includes(user?.role || '')) return fallback;
  return children;
}

// components/guard/require-permission.tsx
export function RequirePermission({
  resource,
  action,
  ownerId,
  children,
  fallback = <Forbidden />,
}: RequirePermissionProps) {
  const { can } = usePermissions();
  
  if (!can(resource, action, ownerId)) return fallback;
  return children;
}

// Usage in components
function ProductActions({ product }: { product: Product }) {
  const { user } = useAuth();
  const { can } = usePermissions();
  
  return (
    <div>
      {can('product', 'update', product.seller_id) && (
        <Button>Edit Product</Button>
      )}
      {can('product', 'delete', product.seller_id) && (
        <Button variant="destructive">Delete</Button>
      )}
    </div>
  );
}
```

## Resource Ownership

### Ownership Rules

| Resource | Owner Field | Ownership Logic |
|----------|-------------|-----------------|
| User | `id` | User owns their own profile |
| Seller | `user_id` | User owns their seller profile |
| Product | `seller_id` | Seller owns their products |
| Auction | `seller_id` | Seller owns their auctions |
| Stream | `seller_id` | Seller owns their streams |
| Order | `buyer_id` + `seller_id` | Both parties have access |
| Bid | `bidder_id` | Buyer owns their bids |
| Review | `reviewer_id` | Buyer owns their reviews |
| Address | `user_id` | User owns their addresses |

### Ownership Check Helper
```go
// utils/ownership.go

func IsOwner(c *gin.Context, resourceType string, resourceID string) bool {
  userID := c.GetString("user_id")
  userRole := c.GetString("user_role")
  
  // Admin bypass
  if userRole == "admin" {
    return true
  }
  
  // Fetch resource and check ownership
  switch resourceType {
  case "product":
    product := productRepo.GetByID(resourceID)
    return product.SellerID == userID
  case "auction":
    auction := auctionRepo.GetByID(resourceID)
    return auction.SellerID == userID
  case "order":
    order := orderRepo.GetByID(resourceID)
    return order.BuyerID == userID || order.SellerID == userID
  // ... other resources
  }
  
  return false
}
```

## Common Authorization Patterns

### 1. Public Read, Protected Write
```go
// Anyone can view, only owner can edit
r.GET("/products/:id", productHandler.Get)     // Public
r.PUT("/products/:id", 
  RequireAuth(),
  RequirePermission("product", "update"),
  productHandler.Update)                        // Protected
```

### 2. Role-Based Actions
```go
// Only sellers can create
r.POST("/auctions",
  RequireAuth(),
  RequireRole("seller", "admin"),
  auctionHandler.Create)

// Only buyers can bid
r.POST("/auctions/:id/bids",
  RequireAuth(),
  RequireRole("buyer"),
  auctionHandler.PlaceBid)
```

### 3. Ownership with Admin Override
```go
// Owner or admin can delete
r.DELETE("/products/:id",
  RequireAuth(),
  func(c *gin.Context) {
    if !IsOwner(c, "product", c.Param("id")) && 
       c.GetString("user_role") != "admin" {
      c.AbortWithStatus(403)
      return
    }
    c.Next()
  },
  productHandler.Delete)
```

---

*Last updated: 2025-02-05*
