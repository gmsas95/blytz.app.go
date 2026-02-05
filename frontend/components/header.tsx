'use client';

import Link from 'next/link';
import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Separator } from '@/components/ui/separator';
import {
  Search,
  Menu,
  Plus,
  Bell,
  Heart,
  User,
  ShoppingBag,
  Video,
  Gavel,
  LogOut,
  Settings,
} from 'lucide-react';
import { useAuthStore } from '@/stores/auth-store';

const navLinks = [
  { name: 'Browse', href: '/auctions', icon: Gavel },
  { name: 'Live Streams', href: '/streams', icon: Video },
  { name: 'Sell', href: '/sell', icon: Plus },
];

export function Header() {
  const [isSearchOpen, setIsSearchOpen] = useState(false);
  const { user, isAuthenticated, logout } = useAuthStore();

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container mx-auto px-4">
        <div className="flex h-16 items-center justify-between">
          {/* Logo */}
          <Link href="/" className="flex items-center space-x-2">
            <div className="h-8 w-8 rounded-lg bg-primary flex items-center justify-center text-primary-foreground font-bold">
              B
            </div>
            <span className="text-xl font-bold hidden sm:inline">Blytz</span>
          </Link>

          {/* Desktop Navigation */}
          <nav className="hidden md:flex items-center space-x-6">
            {navLinks.map((link) => (
              <Link
                key={link.name}
                href={link.href}
                className="flex items-center space-x-1 text-sm font-medium text-muted-foreground hover:text-foreground transition-colors"
              >
                <link.icon className="h-4 w-4" />
                <span>{link.name}</span>
              </Link>
            ))}
          </nav>

          {/* Search Bar - Desktop */}
          <div className="hidden md:flex flex-1 max-w-md mx-6">
            <div className="relative w-full">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                type="search"
                placeholder="Search auctions..."
                className="pl-9 w-full"
              />
            </div>
          </div>

          {/* Right Side Actions */}
          <div className="flex items-center space-x-2">
            {/* Mobile Search Toggle */}
            <Button
              variant="ghost"
              size="icon"
              className="md:hidden"
              onClick={() => setIsSearchOpen(!isSearchOpen)}
            >
              <Search className="h-5 w-5" />
            </Button>

            {/* Notifications */}
            <Button variant="ghost" size="icon" className="relative">
              <Bell className="h-5 w-5" />
              <span className="absolute top-1 right-1 h-2 w-2 rounded-full bg-blytz-red" />
            </Button>

            {/* Watchlist */}
            <Link href="/watchlist">
              <Button variant="ghost" size="icon" className="hidden sm:flex">
                <Heart className="h-5 w-5" />
              </Button>
            </Link>

            {/* User Menu */}
            {isAuthenticated ? (
              <Sheet>
                <SheetTrigger asChild>
                  <Button variant="ghost" size="icon" className="rounded-full">
                    <Avatar className="h-8 w-8">
                      <AvatarImage src={user?.avatar_url} />
                      <AvatarFallback>{user?.first_name?.[0] || 'U'}</AvatarFallback>
                    </Avatar>
                  </Button>
                </SheetTrigger>
                <SheetContent>
                  <SheetHeader>
                    <SheetTitle>My Account</SheetTitle>
                  </SheetHeader>
                  <div className="mt-6 space-y-4">
                    <div className="flex items-center space-x-3">
                      <Avatar className="h-12 w-12">
                        <AvatarImage src={user?.avatar_url} />
                        <AvatarFallback>{user?.first_name?.[0] || 'U'}</AvatarFallback>
                      </Avatar>
                      <div>
                        <p className="font-medium">{user?.first_name || 'User'}</p>
                        <p className="text-sm text-muted-foreground">{user?.email}</p>
                      </div>
                    </div>
                    <Separator />
                    <nav className="space-y-2">
                      <Link href="/profile">
                        <Button variant="ghost" className="w-full justify-start">
                          <User className="mr-2 h-4 w-4" />
                          Profile
                        </Button>
                      </Link>
                      <Link href="/my-auctions">
                        <Button variant="ghost" className="w-full justify-start">
                          <Gavel className="mr-2 h-4 w-4" />
                          My Auctions
                        </Button>
                      </Link>
                      <Link href="/orders">
                        <Button variant="ghost" className="w-full justify-start">
                          <ShoppingBag className="mr-2 h-4 w-4" />
                          Orders
                        </Button>
                      </Link>
                      <Link href="/settings">
                        <Button variant="ghost" className="w-full justify-start">
                          <Settings className="mr-2 h-4 w-4" />
                          Settings
                        </Button>
                      </Link>
                    </nav>
                    <Separator />
                    <Button
                      variant="ghost"
                      className="w-full justify-start text-destructive"
                      onClick={logout}
                    >
                      <LogOut className="mr-2 h-4 w-4" />
                      Logout
                    </Button>
                  </div>
                </SheetContent>
              </Sheet>
            ) : (
              <div className="flex items-center space-x-2">
                <Link href="/login" className="hidden sm:block">
                  <Button variant="ghost">Login</Button>
                </Link>
                <Link href="/register">
                  <Button>Sign Up</Button>
                </Link>
              </div>
            )}

            {/* Mobile Menu */}
            <Sheet>
              <SheetTrigger asChild>
                <Button variant="ghost" size="icon" className="md:hidden">
                  <Menu className="h-5 w-5" />
                </Button>
              </SheetTrigger>
              <SheetContent side="right">
                <SheetHeader>
                  <SheetTitle>Menu</SheetTitle>
                </SheetHeader>
                <nav className="mt-6 space-y-2">
                  {navLinks.map((link) => (
                    <Link key={link.name} href={link.href}>
                      <Button variant="ghost" className="w-full justify-start">
                        <link.icon className="mr-2 h-4 w-4" />
                        {link.name}
                      </Button>
                    </Link>
                  ))}
                  <Separator className="my-4" />
                  <Link href="/watchlist">
                    <Button variant="ghost" className="w-full justify-start">
                      <Heart className="mr-2 h-4 w-4" />
                      Watchlist
                    </Button>
                  </Link>
                  <Link href="/notifications">
                    <Button variant="ghost" className="w-full justify-start">
                      <Bell className="mr-2 h-4 w-4" />
                      Notifications
                    </Button>
                  </Link>
                </nav>
              </SheetContent>
            </Sheet>
          </div>
        </div>

        {/* Mobile Search Bar */}
        {isSearchOpen && (
          <div className="md:hidden py-3">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                type="search"
                placeholder="Search auctions..."
                className="pl-9 w-full"
                autoFocus
              />
            </div>
          </div>
        )}
      </div>
    </header>
  );
}
