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
  X,
} from 'lucide-react';
import { useAuthStore } from '@/stores/auth-store';

const navLinks = [
  { name: 'Auctions', href: '/auctions', icon: Gavel },
  { name: 'Streams', href: '/streams', icon: Video },
  { name: 'Sell', href: '/sell', icon: Plus },
];

export function Header() {
  const [isSearchOpen, setIsSearchOpen] = useState(false);
  const { user, isAuthenticated, logout } = useAuthStore();

  return (
    <header className="sticky top-0 z-50 w-full bg-black/80 backdrop-blur-xl border-b border-neutral-800">
      <div className="container-modern">
        <div className="flex h-16 items-center justify-between">
          {/* Logo */}
          <Link href="/" className="flex items-center gap-2">
            <div className="w-10 h-10 bg-blytz-yellow rounded-lg flex items-center justify-center">
              <span className="text-black font-black text-xl">B</span>
            </div>
            <span className="text-xl font-black text-white hidden sm:block">BLYTZ</span>
          </Link>

          {/* Desktop Navigation */}
          <nav className="hidden md:flex items-center gap-1">
            {navLinks.map((link) => (
              <Link
                key={link.name}
                href={link.href}
                className="flex items-center gap-2 px-4 py-2 text-gray-300 hover:text-white rounded-full hover:bg-neutral-900 transition-all"
              >
                <link.icon className="w-4 h-4" />
                <span className="text-sm font-medium">{link.name}</span>
              </Link>
            ))}
          </nav>

          {/* Search Bar - Desktop */}
          <div className="hidden lg:flex flex-1 max-w-md mx-8">
            <div className="relative w-full group">
              <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-500 group-focus-within:text-blytz-yellow transition-colors" />
              <Input
                type="search"
                placeholder="Search auctions..."
                className="w-full pl-11 pr-4 py-2 bg-neutral-900 border-neutral-800 rounded-full text-white placeholder:text-gray-500 focus:border-blytz-yellow focus:ring-blytz-yellow/20"
              />
            </div>
          </div>

          {/* Right Side Actions */}
          <div className="flex items-center gap-2">
            {/* Mobile Search Toggle */}
            <Button
              variant="ghost"
              size="icon"
              className="lg:hidden text-gray-300 hover:text-white hover:bg-neutral-900"
              onClick={() => setIsSearchOpen(!isSearchOpen)}
            >
              {isSearchOpen ? <X className="w-5 h-5" /> : <Search className="w-5 h-5" />}
            </Button>

            {/* Notifications */}
            <Button variant="ghost" size="icon" className="relative text-gray-300 hover:text-white hover:bg-neutral-900">
              <Bell className="w-5 h-5" />
              <span className="absolute top-1.5 right-1.5 w-2 h-2 bg-blytz-yellow rounded-full" />
            </Button>

            {/* Watchlist */}
            <Link href="/watchlist" className="hidden sm:block">
              <Button variant="ghost" size="icon" className="text-gray-300 hover:text-white hover:bg-neutral-900">
                <Heart className="w-5 h-5" />
              </Button>
            </Link>

            {/* User Menu */}
            {isAuthenticated ? (
              <Sheet>
                <SheetTrigger asChild>
                  <Button variant="ghost" size="icon" className="rounded-full hover:bg-neutral-900">
                    <Avatar className="w-8 h-8 border-2 border-neutral-800">
                      <AvatarImage src={user?.avatar_url} />
                      <AvatarFallback className="bg-blytz-yellow text-black font-bold">
                        {user?.first_name?.[0] || 'U'}
                      </AvatarFallback>
                    </Avatar>
                  </Button>
                </SheetTrigger>
                <SheetContent className="bg-neutral-900 border-neutral-800">
                  <SheetHeader>
                    <SheetTitle className="text-white">My Account</SheetTitle>
                  </SheetHeader>
                  <div className="mt-6 space-y-4">
                    <div className="flex items-center gap-3 p-3 rounded-2xl bg-neutral-800">
                      <Avatar className="w-12 h-12">
                        <AvatarImage src={user?.avatar_url} />
                        <AvatarFallback className="bg-blytz-yellow text-black font-bold">
                          {user?.first_name?.[0] || 'U'}
                        </AvatarFallback>
                      </Avatar>
                      <div>
                        <p className="font-bold text-white">{user?.first_name || 'User'}</p>
                        <p className="text-sm text-gray-400">{user?.email}</p>
                      </div>
                    </div>
                    <Separator className="bg-neutral-800" />
                    <nav className="space-y-1">
                      <Link href="/profile">
                        <Button variant="ghost" className="w-full justify-start text-gray-300 hover:text-white hover:bg-neutral-800">
                          <User className="mr-3 w-4 h-4" />
                          Profile
                        </Button>
                      </Link>
                      <Link href="/my-auctions">
                        <Button variant="ghost" className="w-full justify-start text-gray-300 hover:text-white hover:bg-neutral-800">
                          <Gavel className="mr-3 w-4 h-4" />
                          My Auctions
                        </Button>
                      </Link>
                      <Link href="/orders">
                        <Button variant="ghost" className="w-full justify-start text-gray-300 hover:text-white hover:bg-neutral-800">
                          <ShoppingBag className="mr-3 w-4 h-4" />
                          Orders
                        </Button>
                      </Link>
                      <Link href="/settings">
                        <Button variant="ghost" className="w-full justify-start text-gray-300 hover:text-white hover:bg-neutral-800">
                          <Settings className="mr-3 w-4 h-4" />
                          Settings
                        </Button>
                      </Link>
                    </nav>
                    <Separator className="bg-neutral-800" />
                    <Button
                      variant="ghost"
                      className="w-full justify-start text-red-500 hover:text-red-400 hover:bg-red-500/10"
                      onClick={logout}
                    >
                      <LogOut className="mr-3 w-4 h-4" />
                      Logout
                    </Button>
                  </div>
                </SheetContent>
              </Sheet>
            ) : (
              <div className="flex items-center gap-2">
                <Link href="/login" className="hidden sm:block">
                  <Button variant="ghost" className="text-gray-300 hover:text-white hover:bg-neutral-900">
                    Login
                  </Button>
                </Link>
                <Link href="/register">
                  <Button className="bg-blytz-yellow text-black hover:brightness-110 rounded-full font-semibold">
                    Sign Up
                  </Button>
                </Link>
              </div>
            )}

            {/* Mobile Menu */}
            <Sheet>
              <SheetTrigger asChild>
                <Button variant="ghost" size="icon" className="md:hidden text-gray-300 hover:text-white hover:bg-neutral-900">
                  <Menu className="w-5 h-5" />
                </Button>
              </SheetTrigger>
              <SheetContent side="right" className="bg-neutral-900 border-neutral-800">
                <SheetHeader>
                  <SheetTitle className="text-white">Menu</SheetTitle>
                </SheetHeader>
                <nav className="mt-6 space-y-1">
                  {navLinks.map((link) => (
                    <Link key={link.name} href={link.href}>
                      <Button variant="ghost" className="w-full justify-start text-gray-300 hover:text-white hover:bg-neutral-800">
                        <link.icon className="mr-3 w-4 h-4" />
                        {link.name}
                      </Button>
                    </Link>
                  ))}
                  <Separator className="bg-neutral-800 my-4" />
                  <Link href="/watchlist">
                    <Button variant="ghost" className="w-full justify-start text-gray-300 hover:text-white hover:bg-neutral-800">
                      <Heart className="mr-3 w-4 h-4" />
                      Watchlist
                    </Button>
                  </Link>
                  <Link href="/notifications">
                    <Button variant="ghost" className="w-full justify-start text-gray-300 hover:text-white hover:bg-neutral-800">
                      <Bell className="mr-3 w-4 h-4" />
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
          <div className="lg:hidden py-4 border-t border-neutral-800">
            <div className="relative">
              <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-500" />
              <Input
                type="search"
                placeholder="Search auctions..."
                className="w-full pl-11 pr-4 py-2 bg-neutral-900 border-neutral-800 rounded-full text-white placeholder:text-gray-500"
                autoFocus
              />
            </div>
          </div>
        )}
      </div>
    </header>
  );
}
