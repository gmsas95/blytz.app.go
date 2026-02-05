'use client';

import Link from 'next/link';
import { StreamCard } from '@/components/stream-card';
import { Button } from '@/components/ui/button';
import { ArrowRight } from 'lucide-react';

const liveStreams = [
  {
    id: '1',
    title: 'Sunday Funday Auction! 🔥',
    seller: {
      name: 'Collectibles King',
      avatar: 'https://i.pravatar.cc/150?u=collectibles',
      rating: 4.9,
    },
    thumbnail: 'https://images.unsplash.com/photo-1556742049-0cfed4f6a45d?w=400',
    viewerCount: 1247,
    productCount: 15,
    startedAt: new Date(Date.now() - 3600000).toISOString(),
  },
  {
    id: '2',
    title: 'Rare Sneakers Drop! 💎',
    seller: {
      name: 'Sneaker Vault MY',
      avatar: 'https://i.pravatar.cc/150?u=vault',
      rating: 4.8,
    },
    thumbnail: 'https://images.unsplash.com/photo-1560769629-975ec94e6a86?w=400',
    viewerCount: 892,
    productCount: 8,
    startedAt: new Date(Date.now() - 1800000).toISOString(),
  },
  {
    id: '3',
    title: 'Vintage Watch Collection',
    seller: {
      name: 'Timekeeper Pro',
      avatar: 'https://i.pravatar.cc/150?u=timekeeper',
      rating: 5.0,
    },
    thumbnail: 'https://images.unsplash.com/photo-1523170335258-f5ed11844a49?w=400',
    viewerCount: 456,
    productCount: 12,
    startedAt: new Date(Date.now() - 7200000).toISOString(),
  },
];

export function LiveStreamsSection() {
  return (
    <section className="py-12">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-xl font-bold flex items-center">
            Live Now
            <span className="ml-2 inline-flex h-2 w-2 rounded-full bg-blytz-red animate-pulse-live" />
          </h2>
          <p className="text-sm text-muted-foreground">Join the action in real-time</p>
        </div>
        <Link href="/streams">
          <Button variant="ghost" className="hidden sm:flex">
            View All
            <ArrowRight className="ml-2 h-4 w-4" />
          </Button>
        </Link>
      </div>
      
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
        {liveStreams.map((stream) => (
          <StreamCard key={stream.id} stream={stream} />
        ))}
      </div>
      
      <div className="mt-6 text-center sm:hidden">
        <Link href="/streams">
          <Button variant="outline" className="w-full">
            View All Streams
            <ArrowRight className="ml-2 h-4 w-4" />
          </Button>
        </Link>
      </div>
    </section>
  );
}
