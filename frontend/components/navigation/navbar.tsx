'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Search, ShoppingCart, User, Menu, Bell } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { useAuthStore } from '@/stores/auth-store';
import { cn } from '@/lib/utils';

export function Navbar() {
  const pathname = usePathname();
  const { isAuthenticated, user } = useAuthStore();

  const isActive = (path: string) => pathname === path;

  return (
    <header className="fixed top-0 left-0 right-0 z-50 bg-background/80 backdrop-blur-md border-b">
      <div className="container mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="flex h-16 items-center justify-between">
          {/* Logo */}
          <Link href="/" className="flex items-center space-x-2">
            <span className="text-2xl font-bold tracking-tight">
              <span className="text-blytz-red">Blytz</span>
            </span>
          </Link>

          {/* Search - Hidden on mobile */}
          <div className="hidden md:flex flex-1 max-w-md mx-8">
            <div className="relative w-full">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <input
                type="text"
                placeholder="Search auctions, products..."
                className="w-full rounded-full border bg-background pl-10 pr-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              />
            </div>
          </div>

          {/* Actions */}
          <div className="flex items-center space-x-2">
            {/* Mobile Search */}
            <Button
              variant="ghost"
              size="icon"
              className="md:hidden"
              aria-label="Search"
            >
              <Search className="h-5 w-5" />
            </Button>

            {/* Notifications */}
            {isAuthenticated && (
              <Button
                variant="ghost"
                size="icon"
                className="relative"
                aria-label="Notifications"
              >
                <Bell className="h-5 w-5" />
                <span className="absolute top-1 right-1 h-2 w-2 rounded-full bg-blytz-red" />
              </Button>
            )}

            {/* Cart */}
            <Button
              variant="ghost"
              size="icon"
              className="relative"
              aria-label="Cart"
            >
              <ShoppingCart className="h-5 w-5" />
            </Button>

            {/* Auth */}
            {isAuthenticated ? (
              <Link href="/profile">
                <Button
                  variant="ghost"
                  size="icon"
                  className={cn(
                    'rounded-full',
                    isActive('/profile') && 'bg-accent'
                  )}
                >
                  <User className="h-5 w-5" />
                </Button>
              </Link>
            ) : (
              <Link href="/login">
                <Button size="sm">Sign In</Button>
              </Link>
            )}
          </div>
        </div>
      </div>
    </header>
  );
}
