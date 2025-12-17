# Blytz.live Frontend Development Guide

## 🚀 Quick Start

This comprehensive guide covers everything needed to build the Blytz.live marketplace frontend.

### Prerequisites
- Node.js 18+
- npm or yarn
- Git
- VS Code (recommended)

## 📁 Project Structure

```
frontend/
├── README.md
├── package.json
├── next.config.js
├── tailwind.config.js
├── tsconfig.json
├── .env.local
├── public/
│   ├── manifest.json
│   ├── sw.ts
│   ├── icons/
│   └── images/
├── src/
│   ├── components/
│   │   ├── ui/
│   │   ├── layout/
│   │   ├── product/
│   │   ├── forms/
│   │   └── index.ts
│   ├── pages/
│   ├── hooks/
│   ├── store/
│   ├── services/
│   ├── utils/
│   ├── types/
│   ├── styles/
│   └── app.tsx
└── docs/
    ├── FRONTEND_IMPLEMENTATION.md
    ├── PWA_IMPLEMENTATION.md
    └── COMPONENT_LIBRARY.md
```

## 🛠 Development Setup

### 1. Initialize Next.js Project
```bash
npx create-next-app@latest blytz-frontend --typescript --tailwind --eslint --app
cd blytz-frontend
```

### 2. Install Dependencies
```bash
# Core Dependencies
npm install @next/bundle-analyzer @types/node @types/react @types/react-dom
npm install react react-dom next typescript

# UI & Styling
npm install tailwindcss postcss autoprefixer
npm install @headlessui/react @heroicons/react
npm install class-variance-authority tailwind-merge clsx

# State Management
npm install zustand @tanstack/react-query

# Forms & Validation
npm install react-hook-form @hookform/resolvers zod

# PWA & Performance
npm install next-pwa workbox-window

# Testing
npm install @testing-library/react @testing-library/jest-dom vitest jsdom

# Development Tools
npm install @types/uuid uuid
npm install axios date-fns
```

### 3. Configure Tailwind CSS
```javascript
// tailwind.config.js
/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#eff6ff',
          100: '#dbeafe',
          200: '#bfdbfe',
          300: '#93c5fd',
          400: '#60a5fa',
          500: '#3b82f6',
          600: '#2563eb',
          700: '#1d4ed8',
          800: '#1e40af',
          900: '#1e3a8a',
        },
        success: {
          50: '#ecfdf5',
          100: '#d1fae5',
          200: '#a7f3d0',
          300: '#6ee7b7',
          400: '#34d399',
          500: '#10b981',
          600: '#059669',
          700: '#047857',
          800: '#065f46',
          900: '#064e3b',
        }
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
      },
      animation: {
        'fade-in': 'fadeIn 0.5s ease-in-out',
        'slide-up': 'slideUp 0.3s ease-out',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideUp: {
          '0%': { transform: 'translateY(10px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
      },
    },
  },
  plugins: [],
}
```

### 4. Next.js Configuration
```javascript
// next.config.js
const withPWA = require('next-pwa')({
  dest: 'public',
  disable: process.env.NODE_ENV === 'development',
  register: true,
  skipWaiting: true,
  runtimeCaching: [
    {
      urlPattern: /^https?.\/\/fonts\.googleapis\.com\/.*/i,
      handler: 'CacheFirst',
      options: {
        cacheName: 'google-fonts',
        expiration: {
          maxEntries: 4,
          maxAgeSeconds: 365 * 24 * 60 * 60, // 365 days
        },
      },
    },
    {
      urlPattern: /\.(?:eot|otf|ttc|ttf|woff|woff2|css)$/i,
      handler: 'StaleWhileRevalidate',
      options: {
        cacheName: 'static-resources',
        expiration: {
          maxEntries: 60,
          maxAgeSeconds: 365 * 24 * 60 * 60, // 365 days
        },
      },
    },
  ],
})

/** @type {import('next').NextConfig} */
const nextConfig = {
  experimental: {
    appDir: true,
  },
  images: {
    domains: ['localhost', 'api.blytz.live'],
    formats: ['image/webp', 'image/avif'],
  },
  env: {
    NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL,
    NEXT_PUBLIC_WS_URL: process.env.NEXT_PUBLIC_WS_URL,
  },
}

module.exports = withPWA(nextConfig)
```

### 5. TypeScript Configuration
```json
{
  "compilerOptions": {
    "target": "es5",
    "lib": ["dom", "dom.iterable", "esnext"],
    "allowJs": true,
    "skipLibCheck": true,
    "strict": true,
    "forceConsistentCasingInFileNames": true,
    "noEmit": true,
    "esModuleInterop": true,
    "module": "esnext",
    "moduleResolution": "bundler",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "jsx": "preserve",
    "incremental": true,
    "plugins": [
      {
        "name": "next"
      }
    ],
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"],
      "@/components/*": ["./src/components/*"],
      "@/hooks/*": ["./src/hooks/*"],
      "@/utils/*": ["./src/utils/*"],
      "@/types/*": ["./src/types/*"]
    }
  },
  "include": ["next-env.d.ts", "**/*.ts", "**/*.tsx", ".next/types/**/*.ts"],
  "exclude": ["node_modules"]
}
```

### 6. Environment Variables
```bash
# .env.local
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_WS_URL=ws://localhost:8080
NEXT_PUBLIC_APP_NAME=Blytz.live Marketplace
NEXT_PUBLIC_APP_VERSION=1.0.0
```

## 🧩 Core Components Implementation

### Basic Component Structure
```typescript
// src/components/ui/Button.tsx
import React from 'react'
import { twMerge } from 'tailwind-merge'
import { cva, type VariantProps } from 'class-variance-authority'

const buttonVariants = cva(
  'inline-flex items-center justify-center rounded-md font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:pointer-events-none disabled:opacity-50',
  {
    variants: {
      variant: {
        default: 'bg-primary-600 text-white hover:bg-primary-700 focus:ring-primary-500',
        secondary: 'bg-gray-100 text-gray-900 hover:bg-gray-200 focus:ring-gray-500',
        outline: 'border border-gray-300 bg-transparent text-gray-700 hover:bg-gray-50 focus:ring-primary-500',
      },
      size: {
        sm: 'h-8 px-3 text-sm',
        md: 'h-10 px-4 py-2',
        lg: 'h-12 px-8 text-lg',
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'md',
    },
  }
)

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  loading?: boolean
  fullWidth?: boolean
}

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, loading, fullWidth, children, disabled, ...props }, ref) => {
    return (
      <button
        className={twMerge(
          buttonVariants({ variant, size }),
          fullWidth && 'w-full',
          className
        )}
        ref={ref}
        disabled={disabled || loading}
        {...props}
      >
        {loading ? (
          <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
          </svg>
        ) : children}
      </button>
    )
  }
)

Button.displayName = 'Button'
```

### Advanced Component Example - Product Card
```typescript
// src/components/product/ProductCard.tsx
import React from 'react'
import Image from 'next/image'
import Link from 'next/link'
import { formatCurrency } from '@/utils/format'
import { Button } from '@/components/ui/Button'
import { Badge } from '@/components/ui/Badge'
import { Product } from '@/types/product'

interface ProductCardProps {
  product: Product
  onAddToCart?: (productId: string) => void
  className?: string
}

export const ProductCard: React.FC<ProductCardProps> = ({ product, onAddToCart, className }) => {
  const [imageLoading, setImageLoading] = React.useState(true)
  const [isLoading, setIsLoading] = React.useState(false)

  const handleAddToCart = async () => {
    if (onAddToCart && !isLoading) {
      setIsLoading(true)
      try {
        await onAddToCart(product.id)
      } finally {
        setIsLoading(false)
      }
    }
  }

  return (
    <div className={`bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden hover:shadow-lg transition-shadow ${className}`}>
      {/* Product Image */}
      <div className="relative aspect-square">
        <Link href={`/products/${product.id}`}>
          <Image
            src={product.images?.[0] || '/placeholder-product.jpg'}
            alt={product.title}
            fill
            className={`object-cover ${imageLoading ? 'blur-sm' : ''}`}
            onLoadingComplete={() => setImageLoading(false)}
            sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 25vw"
          />
        </Link>

        {/* Condition Badge */}
        <Badge variant="secondary" className="absolute top-2 left-2">
          {product.condition}
        </Badge>

        {/* Featured Badge */}
        {product.featured && (
          <Badge variant="primary" className="absolute top-2 right-2">
            Featured
          </Badge>
        )}
      </div>

      {/* Product Info */}
      <div className="p-4">
        {/* Title */}
        <Link href={`/products/${product.id}`}>
          <h3 className="font-medium text-gray-900 line-clamp-2 hover:text-primary-600 transition-colors">
            {product.title}
          </h3>
        </Link>

        {/* Price */}
        <div className="mt-2">
          {product.buyNowPrice ? (
            <div className="flex items-center gap-2">
              <span className="text-lg font-bold text-success-600">
                {formatCurrency(product.buyNowPrice)}
              </span>
              <span className="text-sm text-gray-500 line-through">
                {formatCurrency(product.startingPrice)}
              </span>
            </div>
          ) : (
            <span className="text-lg font-bold text-gray-900">
              {formatCurrency(product.startingPrice)}
            </span>
          )}
        </div>

        {/* Seller Info */}
        <div className="mt-2 flex items-center text-sm text-gray-500">
          <div className="w-5 h-5 bg-gray-200 rounded-full mr-2" />
          <span>{product.seller?.firstName} {product.seller?.lastName}</span>
        </div>

        {/* Location */}
        {product.location && (
          <div className="mt-1 flex items-center text-sm text-gray-500">
            <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            {product.location}
          </div>
        )}

        {/* Add to Cart Button */}
        <Button
          className="mt-3 w-full"
          onClick={handleAddToCart}
          loading={isLoading}
          disabled={!product.hasVariants && !product.buyNowPrice}
        >
          {product.hasVariants ? 'Select Options' : 
           product.buyNowPrice ? 'Add to Cart' : 'Place Bid'}
        </Button>
      </div>
    </div>
  )
}
```

## 🔌 API Integration

### API Service Setup
```typescript
// src/services/api.ts
import axios from 'axios'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'

export const api = axios.create({
  baseURL: API_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor for adding auth token
api.interceptors.request.use((config) => {
  if (typeof window !== 'undefined') {
    const token = localStorage.getItem('access_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
  }
  return config
})

// Response interceptor for handling errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Handle unauthorized - redirect to login
      localStorage.removeItem('access_token')
      window.location.href = '/auth/login'
    }
    return Promise.reject(error)
  }
)
```

### Product Service
```typescript
// src/services/products.ts
import { api } from './api'
import { Product, ProductFilter, PaginatedResponse } from '@/types/product'

export const productService = {
  async getProducts(params?: ProductFilter): Promise<PaginatedResponse<Product>> {
    const response = await api.get('/products', { params })
    return response.data
  },

  async getProduct(id: string): Promise<Product> {
    const response = await api.get(`/products/${id}`)
    return response.data.data || response.data.product
  },

  async searchProducts(query: string, params?: ProductFilter): Promise<PaginatedResponse<Product>> {
    const response = await api.get('/catalog/search/products', {
      params: { q: query, ...params }
    })
    return response.data
  },

  async getFeaturedProducts(limit = 10): Promise<Product[]> {
    const response = await api.get('/catalog/search/products/featured', {
      params: { limit }
    })
    return response.data.data
  },

  async getRelatedProducts(id: string, limit = 6): Promise<Product[]> {
    const response = await api.get(`/catalog/search/products/${id}/related`, {
      params: { limit }
    })
    return response.data.data
  }
}
```

## 🗄 State Management

### Zustand Store Setup
```typescript
// src/store/authStore.ts
import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { User } from '@/types/user'

interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  login: (user: User, token: string) => void
  logout: () => void
  updateUser: (user: Partial<User>) => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      isAuthenticated: false,

      login: (user, token) => {
        localStorage.setItem('access_token', token)
        set({ user, token, isAuthenticated: true })
      },

      logout: () => {
        localStorage.removeItem('access_token')
        set({ user: null, token: null, isAuthenticated: false })
      },

      updateUser: (userData) => {
        const currentUser = get().user
        if (currentUser) {
          set({ user: { ...currentUser, ...userData } })
        }
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({ 
        user: state.user, 
        token: state.token, 
        isAuthenticated: state.isAuthenticated 
      }),
    }
  )
)
```

### Cart Store
```typescript
// src/store/cartStore.ts
import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { CartItem, Product } from '@/types'

interface CartState {
  items: CartItem[]
  addItem: (product: Product, quantity?: number) => void
  removeItem: (productId: string) => void
  updateQuantity: (productId: string, quantity: number) => void
  clearCart: () => void
  getTotal: () => number
  getItemCount: () => number
}

export const useCartStore = create<CartState>()(
  persist(
    (set, get) => ({
      items: [],

      addItem: (product, quantity = 1) => {
        const existingItem = get().items.find(item => item.productId === product.id)
        
        if (existingItem) {
          set((state) => ({
            items: state.items.map(item =>
              item.productId === product.id
                ? { ...item, quantity: item.quantity + quantity }
                : item
            )
          }))
        } else {
          const cartItem: CartItem = {
            productId: product.id,
            title: product.title,
            price: product.buyNowPrice || product.startingPrice,
            image: product.images?.[0] || '/placeholder-product.jpg',
            quantity,
            sellerId: product.sellerId
          }
          set((state) => ({ items: [...state.items, cartItem] }))
        }
      },

      removeItem: (productId) => {
        set((state) => ({
          items: state.items.filter(item => item.productId !== productId)
        }))
      },

      updateQuantity: (productId, quantity) => {
        if (quantity <= 0) {
          get().removeItem(productId)
          return
        }

        set((state) => ({
          items: state.items.map(item =>
            item.productId === productId
              ? { ...item, quantity }
              : item
          )
        }))
      },

      clearCart: () => {
        set({ items: [] })
      },

      getTotal: () => {
        return get().items.reduce((total, item) => total + (item.price * item.quantity), 0)
      },

      getItemCount: () => {
        return get().items.reduce((total, item) => total + item.quantity, 0)
      },
    }),
    {
      name: 'cart-storage',
    }
  )
)
```

## 🎣 Custom Hooks

### Authentication Hook
```typescript
// src/hooks/useAuth.ts
import { useAuthStore } from '@/store/authStore'
import { useRouter } from 'next/navigation'
import { authService } from '@/services/auth'

export const useAuth = () => {
  const { user, token, isAuthenticated, login, logout, updateUser } = useAuthStore()
  const router = useRouter()
  const [loading, setLoading] = React.useState(false)

  const handleLogin = async (email: string, password: string) => {
    setLoading(true)
    try {
      const response = await authService.login(email, password)
      login(response.user, response.access_token)
      router.push('/dashboard')
    } catch (error) {
      throw error
    } finally {
      setLoading(false)
    }
  }

  const handleLogout = async () => {
    try {
      await authService.logout()
    } finally {
      logout()
      router.push('/')
    }
  }

  return {
    user,
    token,
    isAuthenticated,
    loading,
    login: handleLogin,
    logout: handleLogout,
    updateUser,
  }
}
```

### Product Hook with React Query
```typescript
// src/hooks/useProducts.ts
import { useQuery, useInfiniteQuery } from '@tanstack/react-query'
import { productService } from '@/services/products'
import { Product, ProductFilter } from '@/types/product'

export const useProducts = (params?: ProductFilter) => {
  return useQuery({
    queryKey: ['products', params],
    queryFn: () => productService.getProducts(params),
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}

export const useInfiniteProducts = (params?: ProductFilter) => {
  return useInfiniteQuery({
    queryKey: ['products', 'infinite', params],
    queryFn: ({ pageParam = 1 }) => 
      productService.getProducts({ ...params, page: pageParam }),
    getNextPageParam: (lastPage, allPages) => {
      if (lastPage.data.length < (params?.limit || 20)) {
        return undefined
      }
      return allPages.length + 1
    },
    initialPageParam: 1,
  })
}

export const useProduct = (id: string) => {
  return useQuery({
    queryKey: ['product', id],
    queryFn: () => productService.getProduct(id),
    enabled: !!id,
  })
}

export const useFeaturedProducts = (limit = 10) => {
  return useQuery({
    queryKey: ['featured-products', limit],
    queryFn: () => productService.getFeaturedProducts(limit),
    staleTime: 10 * 60 * 1000, // 10 minutes
  })
}
```

## 📄 Page Implementation

### Homepage Implementation
```typescript
// src/app/page.tsx
import { Button } from '@/components/ui/Button'
import { ProductCard } from '@/components/product/ProductCard'
import { useFeaturedProducts } from '@/hooks/useProducts'
import { Animated } from '@/components/ui/Animated'

export default function HomePage() {
  const { data: featuredProducts, isLoading, error } = useFeaturedProducts(8)

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Hero Section */}
      <section className="bg-gradient-to-r from-primary-600 to-primary-800 text-white">
        <div className="mx-auto max-w-7xl px-4 py-16 sm:px-6 lg:px-8">
          <div className="text-center">
            <h1 className="text-4xl font-bold sm:text-5xl lg:text-6xl">
              Buy and Sell Anything
              <span className="block text-primary-200">Locally</span>
            </h1>
            <p className="mt-6 text-xl text-primary-100">
              The trusted marketplace for your community
            </p>
            <div className="mt-8 flex justify-center space-x-4">
              <Button size="lg" variant="secondary">
                Start Shopping
              </Button>
              <Button size="lg">
                Start Selling
              </Button>
            </div>
          </div>
        </div>
      </section>

      {/* Categories Section */}
      <section className="py-16 bg-white">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <h2 className="text-3xl font-bold text-gray-900 text-center">
            Browse Categories
          </h2>
          {/* Category Grid Implementation */}
        </div>
      </section>

      {/* Featured Products */}
      <section className="py-16">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <h2 className="text-3xl font-bold text-gray-900 text-center mb-8">
            Featured Products
          </h2>
          
          {isLoading && (
            <div className="text-center">Loading...</div>
          )}
          
          {error && (
            <div className="text-center text-red-600">
              Failed to load products
            </div>
          )}
          
          {featuredProducts && (
            <Animated type="slideIn">
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
                {featuredProducts.map((product) => (
                  <ProductCard key={product.id} product={product} />
                ))}
              </div>
            </Animated>
          )}
        </div>
      </section>
    </div>
  )
}
```

## 🧪 Testing Setup

### Vitest Configuration
```typescript
// vitest.config.ts
import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig({
  plugins: [react()],
  test: {
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'],
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
})
```

### Component Test Example
```typescript
// src/components/__tests__/Button.test.tsx
import { render, screen } from '@testing-library/react'
import { Button } from '../ui/Button'

describe('Button', () => {
  it('renders with default props', () => {
    render(<Button>Click me</Button>)
    const button = screen.getByRole('button', { name: 'Click me' })
    
    expect(button).toBeInTheDocument()
    expect(button).toHaveClass('bg-primary-600', 'text-white', 'h-10', 'px-4')
  })

  it('renders with different variants', () => {
    render(<Button variant="secondary">Secondary</Button>)
    const button = screen.getByRole('button', { name: 'Secondary' })
    
    expect(button).toHaveClass('bg-gray-100', 'text-gray-900')
  })

  it('renders in loading state', () => {
    render(<Button loading>Loading</Button>)
    expect(screen.getByText('Loading')).toBeInTheDocument()
    expect(screen.getByRole('button')).toBeDisabled()
  })

  it('calls onClick handler', async () => {
    const handleClick = vi.fn()
    render(<Button onClick={handleClick}>Click me</Button>)
    
    await user.click(screen.getByRole('button', { name: 'Click me' }))
    expect(handleClick).toHaveBeenCalledTimes(1)
  })
})
```

## 🚀 Deployment

### Build Commands
```bash
# Development
npm run dev

# Build for Production
npm run build

# Start Production Server
npm run start

# Test Build
npm run test

# Lint Code
npm run lint

# Type Check
npm run type-check
```

### Environment Configuration
```bash
# Production Environment Variables
NEXT_PUBLIC_API_URL=https://api.blytz.app/v1
NEXT_PUBLIC_WS_URL=wss://api.blytz.app
NEXT_PUBLIC_SENTRY_DSN=your-sentry-dsn
```

This comprehensive development guide provides everything needed to build the Blytz.Live marketplace frontend! 🚀
