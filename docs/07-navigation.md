# Navigation Structure

## App Routes (Next.js App Router)

### Route Groups

```
(app)/                    # Root group
├── (auth)/               # Auth group - minimal layout
│   ├── layout.tsx        # Clean layout (no navbar)
│   ├── login/
│   │   └── page.tsx
│   ├── register/
│   │   └── page.tsx
│   ├── forgot-password/
│   │   └── page.tsx
│   └── reset-password/
│       └── page.tsx
│
├── (main)/               # Main app group - full layout
│   ├── layout.tsx        # With navbar, sidebar
│   ├── page.tsx          # Home/Feed
│   ├──
│   ├── auctions/
│   │   ├── page.tsx      # Auction listing
│   │   └── [id]/
│   │       └── page.tsx  # Auction detail
│   ├──
│   ├── streams/
│   │   ├── page.tsx      # Live streams
│   │   └── [id]/
│   │       └── page.tsx  # Stream viewer
│   ├──
│   ├── products/
│   │   ├── page.tsx      # Product catalog
│   │   └── [slug]/
│   │       └── page.tsx  # Product detail
│   ├──
│   ├── sellers/
│   │   └── [slug]/
│   │       └── page.tsx  # Seller profile
│   ├──
│   ├── categories/
│   │   └── [slug]/
│   │       └── page.tsx  # Category page
│   ├──
│   ├── search/
│   │   └── page.tsx      # Search results
│   ├──
│   ├── cart/
│   │   └── page.tsx      # Shopping cart
│   ├──
│   ├── checkout/
│   │   └── page.tsx      # Checkout flow
│   ├──
│   ├── orders/
│   │   ├── page.tsx      # Order history
│   │   └── [id]/
│   │       └── page.tsx  # Order detail
│   ├──
│   ├── profile/
│   │   ├── page.tsx      # User profile
│   │   ├── addresses/
│   │   │   └── page.tsx  # Manage addresses
│   │   ├── payments/
│   │   │   └── page.tsx  # Payment methods
│   │   └── settings/
│   │       └── page.tsx  # Account settings
│   └──
│
├── (seller)/             # Seller dashboard group
│   ├── layout.tsx        # Seller sidebar layout
│   ├── dashboard/
│   │   └── page.tsx      # Seller analytics
│   ├── products/
│   │   ├── page.tsx      # Manage products
│   │   └── new/
│   │       └── page.tsx  # Add product
│   ├── auctions/
│   │   ├── page.tsx      # Manage auctions
│   │   └── new/
│   │       └── page.tsx  # Create auction
│   ├── streams/
│   │   ├── page.tsx      # Stream management
│   │   └── go-live/
│   │       └── page.tsx  # Start streaming
│   ├── orders/
│   │   ├── page.tsx      # Seller orders
│   │   └── [id]/
│   │       └── page.tsx  # Order detail
│   └── settings/
│       └── page.tsx      # Store settings
│
└── (admin)/              # Admin panel group
    ├── layout.tsx        # Admin layout
    ├── dashboard/
    │   └── page.tsx      # Admin overview
    ├── users/
    │   └── page.tsx      # User management
    ├── sellers/
    │   └── page.tsx      # Seller verification
    ├── products/
    │   └── page.tsx      # Product moderation
    ├── orders/
    │   └── page.tsx      # Order oversight
    └── settings/
        └── page.tsx      # Platform settings
```

## Navigation Menus

### Main Navigation (Buyer/Seller)

```typescript
const mainNav = [
  {
    label: 'Home',
    href: '/',
    icon: 'Home',
  },
  {
    label: 'Live Streams',
    href: '/streams',
    icon: 'Video',
    badge: 'LIVE', // when live streams active
  },
  {
    label: 'Auctions',
    href: '/auctions',
    icon: 'Gavel',
  },
  {
    label: 'Products',
    href: '/products',
    icon: 'Package',
  },
  {
    label: 'Categories',
    href: '/categories',
    icon: 'Grid',
    children: [
      { label: 'Fashion', href: '/categories/fashion' },
      { label: 'Electronics', href: '/categories/electronics' },
      { label: 'Collectibles', href: '/categories/collectibles' },
      { label: 'Home & Living', href: '/categories/home-living' },
    ],
  },
];

const userNav = [
  {
    label: 'My Profile',
    href: '/profile',
    icon: 'User',
  },
  {
    label: 'My Orders',
    href: '/orders',
    icon: 'ShoppingBag',
  },
  {
    label: 'Saved Items',
    href: '/profile/saved',
    icon: 'Heart',
  },
  {
    label: 'Following',
    href: '/profile/following',
    icon: 'Users',
  },
  {
    label: 'Notifications',
    href: '/profile/notifications',
    icon: 'Bell',
    badge: 'count', // unread count
  },
  {
    divider: true,
  },
  {
    label: 'Seller Dashboard',
    href: '/seller/dashboard',
    icon: 'Store',
    showIf: 'isSeller',
  },
  {
    label: 'Settings',
    href: '/profile/settings',
    icon: 'Settings',
  },
  {
    label: 'Logout',
    action: 'logout',
    icon: 'LogOut',
  },
];
```

### Seller Dashboard Navigation

```typescript
const sellerNav = [
  {
    section: 'Overview',
    items: [
      { label: 'Dashboard', href: '/seller/dashboard', icon: 'LayoutDashboard' },
      { label: 'Analytics', href: '/seller/analytics', icon: 'BarChart3' },
    ],
  },
  {
    section: 'Products',
    items: [
      { label: 'All Products', href: '/seller/products', icon: 'Package' },
      { label: 'Add Product', href: '/seller/products/new', icon: 'Plus' },
      { label: 'Categories', href: '/seller/categories', icon: 'Tags' },
    ],
  },
  {
    section: 'Auctions',
    items: [
      { label: 'All Auctions', href: '/seller/auctions', icon: 'Gavel' },
      { label: 'Create Auction', href: '/seller/auctions/new', icon: 'Plus' },
      { label: 'Bid History', href: '/seller/bids', icon: 'History' },
    ],
  },
  {
    section: 'Live Streams',
    items: [
      { label: 'Go Live', href: '/seller/streams/go-live', icon: 'Video' },
      { label: 'Stream History', href: '/seller/streams', icon: 'Clock' },
      { label: 'Schedule Stream', href: '/seller/streams/schedule', icon: 'Calendar' },
    ],
  },
  {
    section: 'Orders',
    items: [
      { label: 'All Orders', href: '/seller/orders', icon: 'ShoppingBag' },
      { label: 'Pending', href: '/seller/orders?status=pending', icon: 'Clock' },
      { label: 'Shipped', href: '/seller/orders?status=shipped', icon: 'Truck' },
      { label: 'Returns', href: '/seller/orders?status=returned', icon: 'RotateCcw' },
    ],
  },
  {
    section: 'Settings',
    items: [
      { label: 'Store Settings', href: '/seller/settings', icon: 'Store' },
      { label: 'Shipping', href: '/seller/shipping', icon: 'Truck' },
      { label: 'Payouts', href: '/seller/payouts', icon: 'Wallet' },
    ],
  },
];
```

### Admin Navigation

```typescript
const adminNav = [
  {
    section: 'Overview',
    items: [
      { label: 'Dashboard', href: '/admin/dashboard', icon: 'LayoutDashboard' },
      { label: 'Analytics', href: '/admin/analytics', icon: 'BarChart3' },
    ],
  },
  {
    section: 'Users',
    items: [
      { label: 'All Users', href: '/admin/users', icon: 'Users' },
      { label: 'Sellers', href: '/admin/sellers', icon: 'Store' },
      { label: 'Pending Verification', href: '/admin/sellers/pending', icon: 'UserCheck', badge: 'count' },
    ],
  },
  {
    section: 'Content',
    items: [
      { label: 'Products', href: '/admin/products', icon: 'Package' },
      { label: 'Streams', href: '/admin/streams', icon: 'Video' },
      { label: 'Categories', href: '/admin/categories', icon: 'Tags' },
    ],
  },
  {
    section: 'Transactions',
    items: [
      { label: 'Orders', href: '/admin/orders', icon: 'ShoppingBag' },
      { label: 'Payments', href: '/admin/payments', icon: 'CreditCard' },
      { label: 'Disputes', href: '/admin/disputes', icon: 'AlertTriangle' },
    ],
  },
  {
    section: 'Settings',
    items: [
      { label: 'Platform Settings', href: '/admin/settings', icon: 'Settings' },
      { label: 'Integrations', href: '/admin/integrations', icon: 'Plug' },
    ],
  },
];
```

## URL Parameters

### Dynamic Routes

| Route | Parameter | Example |
|-------|-----------|---------|
| `/products/[slug]` | `slug` - Product URL slug | `/products/vintage-rolex-submariner` |
| `/auctions/[id]` | `id` - Auction UUID | `/auctions/550e8400-e29b-41d4-a716-446655440000` |
| `/streams/[id]` | `id` - Stream UUID | `/streams/550e8400-e29b-41d4-a716-446655440000` |
| `/sellers/[slug]` | `slug` - Seller store slug | `/sellers/ahmad-collectibles` |
| `/categories/[slug]` | `slug` - Category slug | `/categories/watches` |
| `/orders/[id]` | `id` - Order UUID | `/orders/550e8400-e29b-41d4-a716-446655440000` |

### Query Parameters

#### Product Listing (`/products`)
```
?category=watches          # Filter by category
&min_price=1000           # Minimum price
&max_price=5000           # Maximum price
&condition=new            # Product condition
&q=rolex                  # Search query
&sort=price_asc           # Sort: price_asc, price_desc, newest, popular
&page=1                   # Page number
&per_page=20              # Items per page
```

#### Auction Listing (`/auctions`)
```
?status=active            # Filter: active, ending_soon, ended
&category=watches         # Category filter
&seller=ahmad-collectibles # Seller filter
&sort=ending_soon         # Sort: ending_soon, newest, bid_count
&page=1
```

#### Stream Listing (`/streams`)
```
?status=live              # Filter: live, scheduled, ended
&category=fashion         # Category
&sort=popular             # Sort: popular, newest
&page=1
```

#### Order History (`/orders`)
```
?status=pending           # Filter status
&page=1
```

#### Search (`/search`)
```
?q=vintage+watch          # Search query
&type=product             # Type: product, seller, stream
&category=all             # Category filter
&sort=relevance           # Sort order
```

## Deep Linking (Mobile)

### URL Schemes

```
blytz://                    # App root
blytz://product/[slug]      # Product detail
blytz://auction/[id]        # Auction detail
blytz://stream/[id]         # Stream viewer
blytz://seller/[slug]       # Seller profile
blytz://order/[id]          # Order detail
blytz://cart                # Shopping cart
```

### Universal Links (iOS) / App Links (Android)

```
https://blytz.app/p/[slug]     # Product
https://blytz.app/a/[id]       # Auction
https://blytz.app/s/[id]       # Stream
https://blytz.app/u/[slug]     # Seller
```

### Push Notification Routing

```typescript
// Notification data structure
interface NotificationData {
  type: 'auction_ending' | 'bid_won' | 'order_shipped' | 'stream_starting';
  action_url: string;  // Deep link URL
  params: {
    id?: string;
    slug?: string;
  };
}

// Routing logic
function handleNotification(data: NotificationData) {
  switch (data.type) {
    case 'auction_ending':
      navigate(`/auctions/${data.params.id}`);
      break;
    case 'bid_won':
      navigate(`/orders/${data.params.id}`);
      break;
    case 'order_shipped':
      navigate(`/orders/${data.params.id}?tab=tracking`);
      break;
    case 'stream_starting':
      navigate(`/streams/${data.params.id}`);
      break;
  }
}
```

## Breadcrumb Navigation

### Breadcrumb Structure

```typescript
// components/breadcrumb.tsx
interface BreadcrumbItem {
  label: string;
  href?: string;  // undefined for current page
}

// Example breadcrumbs by page:

// Product detail
const productBreadcrumbs = [
  { label: 'Home', href: '/' },
  { label: 'Products', href: '/products' },
  { label: 'Watches', href: '/categories/watches' },
  { label: 'Vintage Rolex Submariner' },  // Current (no href)
];

// Auction detail
const auctionBreadcrumbs = [
  { label: 'Home', href: '/' },
  { label: 'Auctions', href: '/auctions' },
  { label: 'Vintage Rolex Auction' },
];

// Seller profile
const sellerBreadcrumbs = [
  { label: 'Home', href: '/' },
  { label: 'Sellers', href: '/sellers' },
  { label: "Ahmad's Collectibles" },
];

// Seller dashboard
const sellerDashboardBreadcrumbs = [
  { label: 'Home', href: '/' },
  { label: 'Seller Dashboard', href: '/seller/dashboard' },
  { label: 'Products', href: '/seller/products' },
  { label: 'Edit Product' },
];
```

## Navigation Guards

### Route Protection

```typescript
// middleware.ts (Next.js)
import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const token = request.cookies.get('token')?.value;
  
  // Public routes
  const publicRoutes = ['/', '/products', '/auctions', '/streams', '/sellers'];
  if (publicRoutes.some(route => pathname.startsWith(route)) && !pathname.includes('/new')) {
    return NextResponse.next();
  }
  
  // Auth required routes
  const authRoutes = ['/profile', '/orders', '/cart', '/checkout'];
  if (authRoutes.some(route => pathname.startsWith(route))) {
    if (!token) {
      return NextResponse.redirect(new URL('/login', request.url));
    }
  }
  
  // Seller only routes
  if (pathname.startsWith('/seller')) {
    if (!token) {
      return NextResponse.redirect(new URL('/login', request.url));
    }
    // Additional role check in layout
  }
  
  // Admin only routes
  if (pathname.startsWith('/admin')) {
    if (!token) {
      return NextResponse.redirect(new URL('/login', request.url));
    }
    // Additional admin check in layout
  }
  
  return NextResponse.next();
}

export const config = {
  matcher: ['/((?!api|_next/static|_next/image|favicon.ico).*)'],
};
```

### Client-Side Guards

```typescript
// hooks/use-route-guard.ts
export function useRouteGuard(requiredRole?: 'buyer' | 'seller' | 'admin') {
  const { user, isLoading } = useAuth();
  const router = useRouter();
  
  useEffect(() => {
    if (isLoading) return;
    
    if (!user) {
      router.push('/login');
      return;
    }
    
    if (requiredRole && user.role !== requiredRole && user.role !== 'admin') {
      router.push('/');
      return;
    }
  }, [user, isLoading, requiredRole]);
  
  return { isLoading, user };
}

// Usage in protected pages
function SellerDashboardPage() {
  const { isLoading, user } = useRouteGuard('seller');
  
  if (isLoading) return <Loading />;
  if (!user) return null;  // Redirecting
  
  return <Dashboard />;
}
```

---

*Last updated: 2025-02-05*
