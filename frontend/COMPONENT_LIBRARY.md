# Blytz.live Component Library & Design System

## 🎨 Design System Foundation

### Brand Identity
```
Primary Color: #6366F1 (Indigo)
Secondary Color: #10B981 (Green)
Accent Color: #F59E0B (Amber)
Error Color: #EF4444 (Red)
Background Color: #FFFFFF (White)
Surface Color: #F9FAFB (Light Gray)
Text Primary: #1F2937 (Dark Gray)
Text Secondary: #6B7280 (Medium Gray)
Border Color: #E5E7EB (Light Gray)
```

### Typography Scale
```css
.text-xs    { font-size: 0.75rem; line-height: 1rem; }   /* 12px */
.text-sm    { font-size: 0.875rem; line-height: 1.25rem; } /* 14px */
.text-base  { font-size: 1rem; line-height: 1.5rem; }   /* 16px */
.text-lg    { font-size: 1.125rem; line-height: 1.75rem; } /* 18px */
.text-xl    { font-size: 1.25rem; line-height: 1.75rem; }  /* 20px */
.text-2xl   { font-size: 1.5rem; line-height: 2rem; }    /* 24px */
.text-3xl   { font-size: 1.875rem; line-height: 2.25rem; } /* 30px */
.text-4xl   { font-size: 2.25rem; line-height: 2.5rem; }  /* 36px */
```

### Spacing System
```css
.space-0    { padding: 0; }
.space-1    { padding: 0.25rem; }  /* 4px */
.space-2    { padding: 0.5rem; }   /* 8px */
.space-3    { padding: 0.75rem; }  /* 12px */
.space-4    { padding: 1rem; }     /* 16px */
.space-5    { padding: 1.25rem; }  /* 20px */
.space-6    { padding: 1.5rem; }   /* 24px */
.space-8    { padding: 2rem; }     /* 32px */
.space-10   { padding: 2.5rem; }  /* 40px */
.space-12   { padding: 3rem; }     /* 48px */
```

### Border Radius
```css
.radius-none { border-radius: 0; }
.radius-sm    { border-radius: 0.125rem; }  /* 2px */
.radius-md    { border-radius: 0.375rem; }  /* 6px */
.radius-lg    { border-radius: 0.5rem; }    /* 8px */
.radius-xl    { border-radius: 0.75rem; }   /* 12px */
.radius-2xl   { border-radius: 1rem; }     /* 16px */
.radius-full   { border-radius: 9999px; }
```

## 🧩 Core Components

### 1. Button Component
```typescript
// components/ui/Button.tsx
import React from 'react';
import { twMerge } from 'tailwind-merge';
import { cva, type VariantProps } from 'class-variance-authority';

const buttonVariants = cva(
  'inline-flex items-center justify-center rounded-md font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:pointer-events-none disabled:opacity-50',
  {
    variants: {
      variant: {
        primary: 'bg-indigo-600 text-white hover:bg-indigo-700 focus:ring-indigo-500',
        secondary: 'bg-green-600 text-white hover:bg-green-700 focus:ring-green-500',
        outline: 'border border-gray-300 bg-white text-gray-700 hover:bg-gray-50 focus:ring-indigo-500',
        ghost: 'text-gray-700 hover:bg-gray-100 focus:ring-indigo-500',
        danger: 'bg-red-600 text-white hover:bg-red-700 focus:ring-red-500'
      },
      size: {
        sm: 'h-9 px-3 text-sm',
        md: 'h-10 px-4 py-2',
        lg: 'h-11 px-8',
        xl: 'h-12 px-8 text-lg'
      }
    },
    defaultVariants: {
      variant: 'primary',
      size: 'md'
    }
  }
);

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement>,
  VariantProps<typeof buttonVariants> {
  loading?: boolean;
  fullWidth?: boolean;
  children: React.ReactNode;
}

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, loading, fullWidth, children, disabled, ...props }, ref) => {
    return (
      <button
        className={twMerge(
          buttonVariants({ variant, size, className }),
          fullWidth && 'w-full'
        )}
        ref={ref}
        disabled={disabled || loading}
        {...props}
      >
        {loading && (
          <svg className="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
          </svg>
        )}
        {children}
      </button>
    );
  }
);

Button.displayName = 'Button';
```

### 2. Input Component
```typescript
// components/ui/Input.tsx
import React from 'react';
import { twMerge } from 'tailwind-merge';
import { cva, type VariantProps } from 'class-variance-authority';

const inputVariants = cva(
  'flex h-10 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm placeholder:text-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent disabled:cursor-not-allowed disabled:opacity-50',
  {
    variants: {
      variant: {
        default: 'border-gray-300',
        error: 'border-red-500 focus:ring-red-500',
        success: 'border-green-500 focus:ring-green-500'
      },
      size: {
        sm: 'h-9 text-xs',
        md: 'h-10 text-sm',
        lg: 'h-11 text-base'
      }
    },
    defaultVariants: {
      variant: 'default',
      size: 'md'
    }
  }
);

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement>,
  VariantProps<typeof inputVariants> {
  label?: string;
  error?: string;
  helperText?: string;
  leftIcon?: React.ReactNode;
  rightIcon?: React.ReactNode;
}

export const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ className, variant, size, label, error, helperText, leftIcon, rightIcon, id, ...props }, ref) => {
    const inputId = id || `input-${React.useId()}`;
    
    return (
      <div className="w-full">
        {label && (
          <label htmlFor={inputId} className="mb-2 block text-sm font-medium text-gray-700">
            {label}
          </label>
        )}
        <div className="relative">
          {leftIcon && (
            <div className="absolute inset-y-0 left-0 flex items-center pl-3">
              {leftIcon}
            </div>
          )}
          <input
            className={twMerge(
              inputVariants({ variant, size, className }),
              leftIcon && 'pl-10',
              rightIcon && 'pr-10'
            )}
            id={inputId}
            ref={ref}
            {...props}
          />
          {rightIcon && (
            <div className="absolute inset-y-0 right-0 flex items-center pr-3">
              {rightIcon}
            </div>
          )}
        </div>
        {(error || helperText) && (
          <p className={twMerge(
            'mt-1 text-xs',
            error ? 'text-red-600' : 'text-gray-500'
          )}>
            {error || helperText}
          </p>
        )}
      </div>
    );
  }
);

Input.displayName = 'Input';
```

### 3. Card Component
```typescript
// components/ui/Card.tsx
import React from 'react';
import { twMerge } from 'tailwind-merge';
import { cva, type VariantProps } from 'class-variance-authority';

const cardVariants = cva(
  'rounded-lg border bg-white shadow-sm',
  {
    variants: {
      variant: {
        default: 'border-gray-200',
        elevated: 'border-gray-200 shadow-lg',
        outlined: 'border-2 border-gray-300'
      },
      padding: {
        none: 'p-0',
        sm: 'p-4',
        md: 'p-6',
        lg: 'p-8'
      }
    },
    defaultVariants: {
      variant: 'default',
      padding: 'md'
    }
  }
);

interface CardProps extends React.HTMLAttributes<HTMLDivElement>,
  VariantProps<typeof cardVariants> {
  children: React.ReactNode;
}

export const Card = React.forwardRef<HTMLDivElement, CardProps>(
  ({ className, variant, padding, children, ...props }, ref) => {
    return (
      <div
        className={twMerge(cardVariants({ variant, padding, className }))}
        ref={ref}
        {...props}
      >
        {children}
      </div>
    );
  }
);

Card.displayName = 'Card';
```

### 4. Modal Component
```typescript
// components/ui/Modal.tsx
import React, { useEffect, useRef } from 'react';
import { createPortal } from 'react-dom';
import { Button } from './Button';

interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  title?: string;
  children: React.ReactNode;
  size?: 'sm' | 'md' | 'lg' | 'xl';
  showCloseButton?: boolean;
}

export const Modal: React.FC<ModalProps> = ({
  isOpen,
  onClose,
  title,
  children,
  size = 'md',
  showCloseButton = true
}) => {
  const modalRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose();
    };

    if (isOpen) {
      document.addEventListener('keydown', handleEscape);
      document.body.style.overflow = 'hidden';
    }

    return () => {
      document.removeEventListener('keydown', handleEscape);
      document.body.style.overflow = 'unset';
    };
  }, [isOpen, onClose]);

  const handleBackdropClick = (e: React.MouseEvent) => {
    if (e.target === e.currentTarget) onClose();
  };

  const sizeClasses = {
    sm: 'max-w-md',
    md: 'max-w-lg',
    lg: 'max-w-2xl',
    xl: 'max-w-4xl'
  };

  if (!isOpen) return null;

  return createPortal(
    <div className="fixed inset-0 z-50 overflow-y-auto">
      <div className="flex min-h-full items-center justify-center p-4">
        <div className="fixed inset-0 bg-black opacity-50" onClick={handleBackdropClick} />
        <div
          ref={modalRef}
          className={`relative w-full ${sizeClasses[size]} bg-white rounded-lg shadow-xl`}
          onClick={(e) => e.stopPropagation()}
        >
          {(title || showCloseButton) && (
            <div className="flex items-center justify-between border-b p-6">
              {title && (
                <h2 className="text-lg font-semibold text-gray-900">{title}</h2>
              )}
              {showCloseButton && (
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={onClose}
                  className="!p-2"
                >
                  <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </Button>
              )}
            </div>
          )}
          <div className="p-6">
            {children}
          </div>
        </div>
      </div>
    </div>,
    document.body
  );
};
```

### 5. Badge Component
```typescript
// components/ui/Badge.tsx
import React from 'react';
import { twMerge } from 'tailwind-merge';
import { cva, type VariantProps } from 'class-variance-authority';

const badgeVariants = cva(
  'inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium',
  {
    variants: {
      variant: {
        default: 'bg-gray-100 text-gray-800',
        primary: 'bg-indigo-100 text-indigo-800',
        success: 'bg-green-100 text-green-800',
        warning: 'bg-amber-100 text-amber-800',
        danger: 'bg-red-100 text-red-800',
        info: 'bg-blue-100 text-blue-800'
      },
      size: {
        sm: 'px-2 py-0.5 text-xs',
        md: 'px-2.5 py-0.5 text-xs',
        lg: 'px-3 py-1 text-sm'
      }
    },
    defaultVariants: {
      variant: 'default',
      size: 'md'
    }
  }
);

interface BadgeProps extends React.HTMLAttributes<HTMLSpanElement>,
  VariantProps<typeof badgeVariants> {
  children: React.ReactNode;
}

export const Badge = React.forwardRef<HTMLSpanElement, BadgeProps>(
  ({ className, variant, size, children, ...props }, ref) => {
    return (
      <span
        className={twMerge(badgeVariants({ variant, size, className }))}
        ref={ref}
        {...props}
      >
        {children}
      </span>
    );
  }
);

Badge.displayName = 'Badge';
```

## 📦 Product-Specific Components

### 1. ProductCard Component
```typescript
// components/product/ProductCard.tsx
import React from 'react';
import Image from 'next/image';
import { Button } from '../ui/Button';
import { Badge } from '../ui/Badge';
import { Product } from '@/types/product';

interface ProductCardProps {
  product: Product;
  onAddToCart?: (productId: string) => void;
  className?: string;
}

export const ProductCard: React.FC<ProductCardProps> = ({
  product,
  onAddToCart,
  className
}) => {
  const [isLoading, setIsLoading] = React.useState(false);

  const handleAddToCart = async () => {
    if (onAddToCart && !isLoading) {
      setIsLoading(true);
      try {
        await onAddToCart(product.id);
      } finally {
        setIsLoading(false);
      }
    }
  };

  return (
    <div className={`bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden hover:shadow-md transition-shadow ${className}`}>
      {/* Product Image */}
      <div className="relative aspect-square">
        <Image
          src={product.images?.[0] || '/placeholder-product.jpg'}
          alt={product.title}
          fill
          className="object-cover"
        />
        
        {/* Condition Badge */}
        <Badge 
          variant="default" 
          className="absolute top-2 left-2"
        >
          {product.condition}
        </Badge>

        {/* Featured Badge */}
        {product.featured && (
          <Badge 
            variant="primary" 
            className="absolute top-2 right-2"
          >
            Featured
          </Badge>
        )}
      </div>

      {/* Product Info */}
      <div className="p-4">
        {/* Title */}
        <h3 className="font-medium text-gray-900 line-clamp-2 h-12">
          {product.title}
        </h3>

        {/* Price */}
        <div className="mt-2">
          {product.buyNowPrice ? (
            <div className="flex items-center gap-2">
              <span className="text-lg font-bold text-green-600">
                ${product.buyNowPrice}
              </span>
              <span className="text-sm text-gray-500 line-through">
                ${product.startingPrice}
              </span>
            </div>
          ) : (
            <span className="text-lg font-bold text-gray-900">
              ${product.startingPrice}
            </span>
          )}
        </div>

        {/* Seller Info */}
        <div className="mt-2 flex items-center text-sm text-gray-500">
          <div className="w-6 h-6 bg-gray-200 rounded-full mr-2" />
          {product.seller?.firstName} {product.seller?.lastName}
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
  );
};
```

### 2. ProductGallery Component
```typescript
// components/product/ProductGallery.tsx
import React, { useState } from 'react';
import Image from 'next/image';

interface ProductGalleryProps {
  images: string[];
  alt: string;
}

export const ProductGallery: React.FC<ProductGalleryProps> = ({ images, alt }) => {
  const [selectedImage, setSelectedImage] = useState(0);
  const [isZoomed, setIsZoomed] = useState(false);

  if (!images || images.length === 0) {
    return (
      <div className="aspect-square bg-gray-200 rounded-lg flex items-center justify-center">
        <span className="text-gray-500">No images available</span>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Main Image */}
      <div className="relative aspect-square overflow-hidden rounded-lg">
        <Image
          src={images[selectedImage]}
          alt={`${alt} - Image ${selectedImage + 1}`}
          fill
          className={`object-cover cursor-pointer transition-transform ${isZoomed ? 'scale-150' : ''}`}
          onClick={() => setIsZoomed(!isZoomed)}
        />
        
        {/* Image Navigation */}
        {images.length > 1 && (
          <>
            <button
              onClick={() => setSelectedImage((prev) => (prev - 1 + images.length) % images.length)}
              className="absolute left-2 top-1/2 -translate-y-1/2 bg-white/80 rounded-full p-2 hover:bg-white"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
              </svg>
            </button>
            <button
              onClick={() => setSelectedImage((prev) => (prev + 1) % images.length)}
              className="absolute right-2 top-1/2 -translate-y-1/2 bg-white/80 rounded-full p-2 hover:bg-white"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
              </svg>
            </button>
          </>
        )}
      </div>

      {/* Thumbnail Grid */}
      {images.length > 1 && (
        <div className="grid grid-cols-4 gap-2">
          {images.map((image, index) => (
            <button
              key={index}
              onClick={() => setSelectedImage(index)}
              className={`relative aspect-square rounded-lg overflow-hidden border-2 ${
                selectedImage === index ? 'border-indigo-500' : 'border-gray-200'
              }`}
            >
              <Image
                src={image}
                alt={`${alt} - Thumbnail ${index + 1}`}
                fill
                className="object-cover"
              />
            </button>
          ))}
        </div>
      )}
    </div>
  );
};
```

## 🎨 Layout Components

### 1. Header Component
```typescript
// components/layout/Header.tsx
import React, { useState } from 'react';
import Link from 'next/link';
import { Button } from '../ui/Button';
import { Input } from '../ui/Input';
import { Badge } from '../ui/Badge';
import { useAuth } from '@/hooks/useAuth';
import { useCart } from '@/hooks/useCart';

export const Header: React.FC = () => {
  const { user, logout } = useAuth();
  const { items } = useCart();
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');

  const cartItemCount = items.reduce((total, item) => total + item.quantity, 0);

  return (
    <header className="bg-white shadow-sm sticky top-0 z-40">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="flex h-16 items-center justify-between">
          {/* Logo */}
          <Link href="/" className="flex items-center">
            <span className="text-2xl font-bold text-indigo-600">Blytz</span>
            <span className="text-2xl font-light text-gray-600">.live</span>
          </Link>

          {/* Search Bar - Desktop */}
          <div className="hidden md:flex flex-1 max-w-lg mx-8">
            <Input
              type="search"
              placeholder="Search for anything..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              leftIcon={
                <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              }
            />
          </div>

          {/* Navigation Items */}
          <div className="flex items-center space-x-4">
            {/* Browse - Desktop */}
            <Link href="/products" className="hidden md:block text-gray-700 hover:text-indigo-600">
              Browse
            </Link>

            {/* Sell - Desktop */}
            <Link href="/sell" className="hidden md:block text-gray-700 hover:text-indigo-600">
              Sell
            </Link>

            {/* Cart */}
            <Link href="/cart" className="relative">
              <svg className="w-6 h-6 text-gray-700 hover:text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 002 2v2a2 2 0 01-2 2H5a2 2 0 01-2-2v-2a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m14 0V9a2 2 0 00-2-2" />
              </svg>
              {cartItemCount > 0 && (
                <Badge variant="danger" className="absolute -top-2 -right-2 !p-1">
                  {cartItemCount > 99 ? '99+' : cartItemCount}
                </Badge>
              )}
            </Link>

            {/* User Menu */}
            {user ? (
              <div className="relative">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
                >
                  <div className="w-6 h-6 bg-gray-300 rounded-full" />
                </Button>
                
                {/* Dropdown Menu */}
                {isMobileMenuOpen && (
                  <div className="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg border border-gray-200 py-1">
                    <Link href="/dashboard" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">
                      My Account
                    </Link>
                    <Link href="/orders" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">
                      My Orders
                    </Link>
                    {user.role === 'seller' && (
                      <Link href="/sell" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">
                        Seller Dashboard
                      </Link>
                    )}
                    <button
                      onClick={logout}
                      className="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                    >
                      Sign Out
                    </button>
                  </div>
                )}
              </div>
            ) : (
              <Link href="/auth/login">
                <Button size="sm">Sign In</Button>
              </Link>
            )}

            {/* Mobile Menu Toggle */}
            <button
              className="md:hidden"
              onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
            >
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
              </svg>
            </button>
          </div>
        </div>

        {/* Mobile Menu */}
        {isMobileMenuOpen && (
          <div className="md:hidden border-t border-gray-200 py-4">
            <div className="space-y-3">
              <Input
                type="search"
                placeholder="Search..."
                leftIcon={
                  <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                  </svg>
                }
              />
              <Link href="/products" className="block py-2 text-gray-700">Browse</Link>
              <Link href="/sell" className="block py-2 text-gray-700">Sell</Link>
            </div>
          </div>
        )}
      </div>
    </header>
  );
};
```

## 📐 Responsive Grid System

### Tailwind Grid Configuration
```typescript
// tailwind.config.js
module.exports = {
  theme: {
    extend: {
      container: {
        center: true,
        padding: {
          DEFAULT: '1rem',
          sm: '2rem',
          lg: '4rem',
          xl: '5rem',
          '2xl': '6rem',
        },
      },
      gridTemplateColumns: {
        '13': 'repeat(13, minmax(0, 1fr))',
        '14': 'repeat(14, minmax(0, 1fr))',
        '15': 'repeat(15, minmax(0, 1fr))',
        '16': 'repeat(16, minmax(0, 1fr))',
      }
    }
  }
}
```

### Responsive Components
```typescript
// components/layout/Grid.tsx
import React from 'react';
import { twMerge } from 'tailwind-merge';

interface GridProps {
  children: React.ReactNode;
  cols?: {
    xs?: 1 | 2 | 3 | 4;
    sm?: 1 | 2 | 3 | 4 | 6;
    md?: 1 | 2 | 3 | 4 | 6 | 8;
    lg?: 1 | 2 | 3 | 4 | 6 | 8 | 12;
    xl?: 1 | 2 | 3 | 4 | 6 | 8 | 12;
  };
  gap?: number;
  className?: string;
}

export const Grid: React.FC<GridProps> = ({ 
  children, 
  cols = {}, 
  gap = 4,
  className 
}) => {
  const gridClasses = [
    'grid',
    cols.xs && `grid-cols-${cols.xs}`,
    cols.sm && `sm:grid-cols-${cols.sm}`,
    cols.md && `md:grid-cols-${cols.md}`,
    cols.lg && `lg:grid-cols-${cols.lg}`,
    cols.xl && `xl:grid-cols-${cols.xl}`,
    `gap-${gap}`
  ].filter(Boolean).join(' ');

  return (
    <div className={twMerge(gridClasses, className)}>
      {children}
    </div>
  );
};

// Usage Example:
export const ProductGrid = ({ products }) => {
  return (
    <Grid 
      cols={{ 
        xs: 2, 
        sm: 2, 
        md: 3, 
        lg: 4, 
        xl: 4 
      }}
      gap={6}
    >
      {products.map(product => (
        <ProductCard key={product.id} product={product} />
      ))}
    </Grid>
  );
};
```

## 🎭 Animation & Transitions

### Animation Presets
```typescript
// components/ui/Animated.tsx
import React from 'react';
import { motion, AnimatePresence } from 'framer-motion';

const fadeIn = {
  initial: { opacity: 0 },
  animate: { opacity: 1 },
  exit: { opacity: 0 }
};

const slideInFromBottom = {
  initial: { y: 20, opacity: 0 },
  animate: { y: 0, opacity: 1 },
  exit: { y: -20, opacity: 0 }
};

const scaleIn = {
  initial: { scale: 0.9, opacity: 0 },
  animate: { scale: 1, opacity: 1 },
  exit: { scale: 0.9, opacity: 0 }
};

interface AnimatedProps {
  children: React.ReactNode;
  type?: 'fadeIn' | 'slideIn' | 'scaleIn';
  className?: string;
}

export const Animated: React.FC<AnimatedProps> = ({ 
  children, 
  type = 'fadeIn',
  className 
}) => {
  const variants = {
    fadeIn,
    slideIn: slideInFromBottom,
    scaleIn
  };

  return (
    <motion.div
      variants={variants[type]}
      initial="initial"
      animate="animate"
      exit="exit"
      transition={{ duration: 0.3 }}
      className={className}
    >
      {children}
    </motion.div>
  );
};

export const FadeIn = ({ children, className }) => (
  <Animated type="fadeIn" className={className}>
    {children}
  </Animated>
);

export const SlideIn = ({ children, className }) => (
  <Animated type="slideIn" className={className}>
    {children}
  </Animated>
);
```

This comprehensive component library provides everything needed to build a **professional, consistent, and scalable** frontend for Blytz.live marketplace! 🚀