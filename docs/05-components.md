# Frontend Components

## Overview

Next.js 15+ component architecture using TypeScript, Tailwind CSS, and shadcn/ui.

## Component Hierarchy

```
app/
├── layout.tsx              # Root layout (providers, fonts)
├── page.tsx                # Home page
├── globals.css             # Global styles
│
├── (auth)/                 # Auth route group
│   ├── layout.tsx          # Auth layout (clean, minimal)
│   ├── login/
│   │   └── page.tsx
│   └── register/
│       └── page.tsx
│
├── (main)/                 # Main app route group
│   ├── layout.tsx          # Main layout (navbar, sidebar)
│   ├── page.tsx            # Feed/dashboard
│   ├── auctions/
│   │   ├── page.tsx        # Auction listing
│   │   └── [id]/
│   │       └── page.tsx    # Auction detail
│   ├── streams/
│   │   ├── page.tsx        # Stream discovery
│   │   └── [id]/
│   │       └── page.tsx    # Stream viewer
│   ├── products/
│   │   ├── page.tsx        # Product catalog
│   │   └── [slug]/
│   │       └── page.tsx    # Product detail
│   ├── sellers/
│   │   └── [slug]/
│   │       └── page.tsx    # Seller profile
│   ├── orders/
│   │   └── page.tsx        # Order history
│   ├── cart/
│   │   └── page.tsx        # Shopping cart
│   └── checkout/
│       └── page.tsx        # Checkout flow
│
└── api/                    # API routes
    └── webhooks/
        └── stripe/route.ts

components/
├── ui/                     # shadcn/ui components
│   ├── button.tsx
│   ├── card.tsx
│   ├── dialog.tsx
│   ├── input.tsx
│   └── ...
│
├── layout/                 # Layout components
│   ├── navbar.tsx
│   ├── sidebar.tsx
│   ├── footer.tsx
│   └── mobile-nav.tsx
│
├── auction/                # Auction-specific
│   ├── auction-card.tsx
│   ├── auction-grid.tsx
│   ├── bid-form.tsx
│   ├── bid-history.tsx
│   ├── countdown-timer.tsx
│   └── current-bid.tsx
│
├── stream/                 # Streaming
│   ├── stream-player.tsx
│   ├── stream-grid.tsx
│   ├── stream-preview.tsx
│   ├── go-live-button.tsx
│   └── stream-info.tsx
│
├── chat/                   # Chat
│   ├── chat-panel.tsx
│   ├── chat-message.tsx
│   ├── chat-input.tsx
│   └── emoji-picker.tsx
│
├── product/                # Product
│   ├── product-card.tsx
│   ├── product-grid.tsx
│   ├── product-gallery.tsx
│   ├── product-info.tsx
│   └── add-to-cart.tsx
│
├── seller/                 # Seller
│   ├── seller-card.tsx
│   ├── seller-header.tsx
│   ├── follow-button.tsx
│   └── seller-stats.tsx
│
├── order/                  # Order
│   ├── order-card.tsx
│   ├── order-list.tsx
│   ├── order-status.tsx
│   └── tracking-info.tsx
│
├── payment/                # Payment
│   ├── payment-form.tsx
│   ├── card-element.tsx
│   └── checkout-summary.tsx
│
└── shared/                 # Shared/Reusable
    ├── loading.tsx
    ├── error-boundary.tsx
    ├── empty-state.tsx
    ├── search-bar.tsx
    ├── category-filter.tsx
    ├── price-range.tsx
    ├── sort-dropdown.tsx
    └── avatar.tsx

hooks/
├── use-auth.ts
├── use-auction.ts
├── use-bid.ts
├── use-stream.ts
├── use-chat.ts
├── use-product.ts
├── use-cart.ts
├── use-order.ts
├── use-notification.ts
├── use-scroll.ts
└── use-media-query.ts

lib/
├── api.ts                  # API client
├── socket.ts               # Socket.io client
├── livekit.ts              # LiveKit client
├── utils.ts                # Utilities
└── constants.ts            # Constants

stores/
├── auth-store.ts
├── cart-store.ts
├── ui-store.ts
└── notification-store.ts

types/
├── index.ts
├── user.ts
├── product.ts
├── auction.ts
├── stream.ts
├── order.ts
└── api.ts
```

## Core Components

### Layout Components

#### Navbar
```typescript
// components/layout/navbar.tsx
interface NavbarProps {
  transparent?: boolean;
}

// Features:
// - Logo/brand
// - Search bar
// - Navigation links
// - Cart icon with count
// - Notification bell
// - User menu (avatar dropdown)
// - Mobile hamburger menu
```

#### Sidebar
```typescript
// components/layout/sidebar.tsx
interface SidebarProps {
  className?: string;
}

// Features:
// - Categories menu
// - Trending auctions
// - Following sellers
// - Mobile: slide-out drawer
```

### Auction Components

#### AuctionCard
```typescript
// components/auction/auction-card.tsx
interface AuctionCardProps {
  auction: {
    id: string;
    product: Product;
    seller: Seller;
    currentBid: number;
    bidCount: number;
    endTime: Date;
    thumbnailUrl: string;
    isLive: boolean;
    viewerCount?: number;
  };
  variant?: 'default' | 'compact' | 'featured';
}

// Usage:
<AuctionCard 
  auction={auction} 
  variant="featured" 
/>
```

#### BidForm
```typescript
// components/auction/bid-form.tsx
interface BidFormProps {
  auctionId: string;
  currentBid: number;
  minIncrement: number;
  onBidPlaced: (amount: number) => void;
  isActive: boolean;
}

// Features:
// - Quick bid buttons (+RM50, +RM100)
// - Custom amount input
// - Auto-bid toggle
// - Max auto-bid input
// - Submit button with validation
```

#### CountdownTimer
```typescript
// components/auction/countdown-timer.tsx
interface CountdownTimerProps {
  endTime: Date;
  onExpire?: () => void;
  size?: 'sm' | 'md' | 'lg';
}

// Features:
// - Days, hours, minutes, seconds
// - Warning color when < 1 hour
// - Critical color when < 5 minutes
```

### Stream Components

#### StreamPlayer
```typescript
// components/stream/stream-player.tsx
interface StreamPlayerProps {
  streamId: string;
  livekitToken: string;
  livekitUrl: string;
  isHost?: boolean;
}

// Features:
// - LiveKit video/audio
// - Screen sharing
// - Video controls (mute, camera)
// - Connection quality indicator
// - Viewer count
```

#### StreamPreview
```typescript
// components/stream/stream-preview.tsx
interface StreamPreviewProps {
  stream: {
    id: string;
    title: string;
    seller: Seller;
    thumbnailUrl: string;
    viewerCount: number;
    isLive: boolean;
  };
}

// Features:
// - Thumbnail with LIVE badge
// - Seller info
// - Viewer count
// - Hover play preview
```

### Chat Components

#### ChatPanel
```typescript
// components/chat/chat-panel.tsx
interface ChatPanelProps {
  streamId: string;
  isModerator?: boolean;
}

// Features:
// - Message list (auto-scroll)
// - Message input
// - Emoji picker
// - Moderation tools (delete, timeout)
// - User mentions
```

#### ChatMessage
```typescript
// components/chat/chat-message.tsx
interface ChatMessageProps {
  message: {
    id: string;
    userId: string;
    userName: string;
    userAvatar?: string;
    content: string;
    type: 'text' | 'system';
    createdAt: Date;
    isDeleted?: boolean;
  };
  isOwn?: boolean;
}

// Features:
// - User avatar
// - Name with role badge
// - Message content
// - Timestamp
// - Delete button (for moderators)
```

### Product Components

#### ProductCard
```typescript
// components/product/product-card.tsx
interface ProductCardProps {
  product: Product;
  variant?: 'default' | 'horizontal' | 'mini';
  showAddToCart?: boolean;
}

// Features:
// - Image with hover zoom
// - Product name
// - Price
// - Condition badge
// - Seller info
// - Add to cart button
```

#### ProductGallery
```typescript
// components/product/product-gallery.tsx
interface ProductGalleryProps {
  images: {
    url: string;
    thumbnailUrl: string;
    alt: string;
  }[];
}

// Features:
// - Main image display
// - Thumbnail strip
// - Lightbox on click
// - Swipe on mobile
// - Zoom on hover (desktop)
```

### Payment Components

#### PaymentForm
```typescript
// components/payment/payment-form.tsx
interface PaymentFormProps {
  orderId: string;
  amount: number;
  onSuccess: () => void;
  onError: (error: Error) => void;
}

// Features:
// - Stripe Elements integration
// - Card input
// - FPX bank selection
// - E-wallet options
// - Submit button with loading
// - Error display
```

## Custom Hooks

### useAuction
```typescript
// hooks/use-auction.ts
function useAuction(auctionId: string) {
  return {
    auction: Auction | null;
    isLoading: boolean;
    error: Error | null;
    placeBid: (amount: number) => Promise<void>;
    refresh: () => void;
  };
}

// Usage:
const { auction, placeBid, isLoading } = useAuction(id);
```

### useBid
```typescript
// hooks/use-bid.ts
function useBid(auctionId: string) {
  return {
    currentBid: number;
    bidCount: number;
    highestBidder: User | null;
    timeRemaining: number;
    placeBid: (amount: number) => Promise<void>;
    isPlacing: boolean;
  };
}
```

### useStream
```typescript
// hooks/use-stream.ts
function useStream(streamId: string) {
  return {
    stream: Stream | null;
    isLive: boolean;
    viewerCount: number;
    joinStream: () => void;
    leaveStream: () => void;
  };
}
```

### useChat
```typescript
// hooks/use-chat.ts
function useChat(streamId: string) {
  return {
    messages: Message[];
    sendMessage: (content: string) => void;
    deleteMessage: (messageId: string) => void;
    isConnected: boolean;
  };
}
```

## State Management (Zustand)

### Auth Store
```typescript
// stores/auth-store.ts
interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
  register: (data: RegisterData) => Promise<void>;
}

const useAuthStore = create<AuthState>((set) => ({
  // ... implementation
}));
```

### Cart Store
```typescript
// stores/cart-store.ts
interface CartState {
  items: CartItem[];
  total: number;
  itemCount: number;
  addItem: (product: Product, quantity: number) => void;
  removeItem: (productId: string) => void;
  updateQuantity: (productId: string, quantity: number) => void;
  clearCart: () => void;
}
```

## Styling Guidelines

### Tailwind Classes
```typescript
// Button variants
const buttonVariants = cva(
  "inline-flex items-center justify-center rounded-md text-sm font-medium",
  {
    variants: {
      variant: {
        default: "bg-primary text-primary-foreground hover:bg-primary/90",
        destructive: "bg-destructive text-destructive-foreground",
        outline: "border border-input bg-background hover:bg-accent",
        ghost: "hover:bg-accent hover:text-accent-foreground",
      },
      size: {
        default: "h-10 px-4 py-2",
        sm: "h-9 rounded-md px-3",
        lg: "h-11 rounded-md px-8",
        icon: "h-10 w-10",
      },
    },
  }
);
```

### Responsive Breakpoints
```
sm: 640px   - Mobile landscape
md: 768px   - Tablet
lg: 1024px  - Desktop
xl: 1280px  - Large desktop
2xl: 1536px - Extra large
```

### Spacing Scale
```
1: 0.25rem  (4px)
2: 0.5rem   (8px)
3: 0.75rem  (12px)
4: 1rem     (16px)
5: 1.25rem  (20px)
6: 1.5rem   (24px)
8: 2rem     (32px)
10: 2.5rem  (40px)
12: 3rem    (48px)
16: 4rem    (64px)
```

## Component Checklist

### New Component Checklist
- [ ] TypeScript interface defined
- [ ] Props documented with JSDoc
- [ ] Default props set
- [ ] Responsive design tested
- [ ] Loading state handled
- [ ] Error state handled
- [ ] Accessibility (ARIA labels, keyboard nav)
- [ ] Storybook story created (optional)
- [ ] Unit tests written

### Performance Checklist
- [ ] Memoized with React.memo if needed
- [ ] useMemo for expensive calculations
- [ ] useCallback for event handlers
- [ ] Images optimized (next/image)
- [ ] Lazy loading for below-fold content

---

*Last updated: 2025-02-05*
