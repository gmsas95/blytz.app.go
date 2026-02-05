'use client';

import Link from 'next/link';
import { AuctionCard } from '@/components/auction-card';
import { Button } from '@/components/ui/button';
import { ArrowRight } from 'lucide-react';

const featuredAuctions = [
  {
    id: '1',
    title: 'Vintage Film Camera Collection',
    description: 'Rare 35mm cameras from the 1960s-80s in excellent condition',
    currentPrice: 1250,
    startingPrice: 500,
    image: 'https://images.unsplash.com/photo-1526170375885-4d8ecf77b99f?w=400',
    seller: { name: 'CameraVault', avatar: 'https://i.pravatar.cc/150?u=camera' },
    endTime: new Date(Date.now() + 86400000 * 2).toISOString(),
    bidCount: 24,
    isLive: true,
  },
  {
    id: '2',
    title: 'Limited Edition Sneakers',
    description: 'Nike Air Jordan 1 Retro High - Size US 10',
    currentPrice: 380,
    startingPrice: 200,
    image: 'https://images.unsplash.com/photo-1552346154-21d32810aba3?w=400',
    seller: { name: 'SneakerHead', avatar: 'https://i.pravatar.cc/150?u=sneaker' },
    endTime: new Date(Date.now() + 86400000 * 3).toISOString(),
    bidCount: 56,
    isLive: false,
  },
  {
    id: '3',
    title: 'Antique Pocket Watch',
    description: '18K Gold Swiss pocket watch from 1890s, fully restored',
    currentPrice: 2800,
    startingPrice: 1500,
    image: 'https://images.unsplash.com/photo-1509048191080-d2984bad6ae5?w=400',
    seller: { name: 'Timeless Treasures', avatar: 'https://i.pravatar.cc/150?u=time' },
    endTime: new Date(Date.now() + 86400000).toISOString(),
    bidCount: 12,
    isLive: true,
  },
  {
    id: '4',
    title: 'Gaming Laptop RTX 4080',
    description: 'ASUS ROG Strix, barely used, warranty until 2026',
    currentPrice: 1850,
    startingPrice: 1200,
    image: 'https://images.unsplash.com/photo-1603302576837-37561b2e2302?w=400',
    seller: { name: 'TechDeals', avatar: 'https://i.pravatar.cc/150?u=tech' },
    endTime: new Date(Date.now() + 86400000 * 4).toISOString(),
    bidCount: 18,
    isLive: false,
  },
];

export function FeaturedAuctionsSection() {
  return (
    <section className="py-12">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-xl font-bold">Featured Auctions</h2>
          <p className="text-sm text-muted-foreground">Hot items ending soon</p>
        </div>
        <Link href="/auctions">
          <Button variant="ghost" className="hidden sm:flex">
            View All
            <ArrowRight className="ml-2 h-4 w-4" />
          </Button>
        </Link>
      </div>
      
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        {featuredAuctions.map((auction) => (
          <AuctionCard key={auction.id} auction={auction} />
        ))}
      </div>
      
      <div className="mt-6 text-center sm:hidden">
        <Link href="/auctions">
          <Button variant="outline" className="w-full">
            View All Auctions
            <ArrowRight className="ml-2 h-4 w-4" />
          </Button>
        </Link>
      </div>
    </section>
  );
}
