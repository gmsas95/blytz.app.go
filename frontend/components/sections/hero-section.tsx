'use client';

import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Play, TrendingUp } from 'lucide-react';

export function HeroSection() {
  return (
    <section className="relative overflow-hidden rounded-2xl bg-gradient-to-br from-primary via-primary/90 to-primary/80 text-primary-foreground">
      <div className="absolute inset-0 bg-[url('/grid.svg')] opacity-10" />
      
      <div className="relative px-6 py-12 sm:px-12 sm:py-16 lg:py-20">
        <div className="max-w-2xl">
          <div className="flex items-center space-x-2 mb-4">
            <span className="inline-flex items-center rounded-full bg-blytz-red/20 px-3 py-1 text-xs font-medium text-white">
              <span className="mr-1.5 h-2 w-2 rounded-full bg-blytz-red animate-pulse-live" />
              Live Now
            </span>
            <span className="inline-flex items-center rounded-full bg-white/20 px-3 py-1 text-xs font-medium">
              <TrendingUp className="mr-1 h-3 w-3" />
              Trending
            </span>
          </div>
          
          <h1 className="text-3xl font-bold tracking-tight sm:text-4xl lg:text-5xl mb-4">
            Discover Unique Finds Through{' '}
            <span className="text-blytz-gold">Live Auctions</span>
          </h1>
          
          <p className="text-lg text-primary-foreground/80 mb-8 max-w-lg">
            Join live streams, bid in real-time, and win exclusive products from trusted sellers across Malaysia.
          </p>
          
          <div className="flex flex-col sm:flex-row gap-3">
            <Link href="/streams">
              <Button size="lg" className="bg-blytz-red hover:bg-blytz-red-dark text-white w-full sm:w-auto">
                <Play className="mr-2 h-4 w-4" />
                Watch Live
              </Button>
            </Link>
            <Link href="/auctions">
              <Button
                size="lg"
                variant="secondary"
                className="bg-white/10 hover:bg-white/20 text-white border-0 w-full sm:w-auto"
              >
                Browse Auctions
              </Button>
            </Link>
          </div>
        </div>
      </div>
      
      {/* Decorative elements */}
      <div className="absolute right-0 top-0 h-full w-1/3 hidden lg:block">
        <div className="absolute right-12 top-1/2 -translate-y-1/2 w-64 h-64 rounded-full bg-blytz-gold/20 blur-3xl" />
        <div className="absolute right-24 bottom-12 w-48 h-48 rounded-full bg-blytz-red/20 blur-2xl" />
      </div>
    </section>
  );
}
