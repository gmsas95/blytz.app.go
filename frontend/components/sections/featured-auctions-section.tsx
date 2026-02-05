'use client';

import Link from 'next/link';
import { AuctionCard } from '@/components/auction-card';
import { Button } from '@/components/ui/button';
import { ArrowRight, Flame } from 'lucide-react';

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
    <section className="section-padding bg-neutral-950">
      <div className="container-modern">
        {/* Header */}
        <div className="flex flex-col md:flex-row md:items-end md:justify-between gap-6 mb-12">
          <div>
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-blytz-yellow/10 border border-blytz-yellow/30 mb-4">
              <Flame className="w-4 h-4 text-blytz-yellow" />
              <span className="text-blytz-yellow text-sm font-semibold">HOT AUCTIONS</span>
            </div>
            <h2 className="text-4xl md:text-5xl font-black text-white">
              FEATURED <span className="text-blytz-yellow">ITEMS</span>
            </h2>
            <p className="text-gray-400 mt-3 text-lg">Don&apos;t miss out on these incredible deals</p>
          </div>
          <Link href="/auctions">
            <Button className="btn-primary gap-2">
              View All
              <ArrowRight className="w-4 h-4" />
            </Button>
          </Link>
        </div>

        {/* Auctions Grid */}
        <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
          {featuredAuctions.map((auction) => (
            <AuctionCard key={auction.id} auction={auction} />
          ))}
        </div>
      </div>
    </section>
  );
}
