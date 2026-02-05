'use client';

import Link from 'next/link';
import { StreamCard } from '@/components/stream-card';
import { Button } from '@/components/ui/button';
import { ArrowRight, Radio } from 'lucide-react';

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
    <section className="section-padding bg-black relative overflow-hidden">
      {/* Background Accent */}
      <div className="absolute top-0 right-0 w-1/2 h-full bg-gradient-to-l from-blytz-yellow/5 to-transparent" />
      
      <div className="container-modern relative z-10">
        {/* Header */}
        <div className="flex flex-col md:flex-row md:items-end md:justify-between gap-6 mb-12">
          <div>
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-red-500/10 border border-red-500/30 mb-4">
              <Radio className="w-4 h-4 text-red-500 animate-pulse" />
              <span className="text-red-500 text-sm font-semibold">LIVE NOW</span>
            </div>
            <h2 className="text-4xl md:text-5xl font-black text-white">
              LIVE <span className="text-blytz-yellow">STREAMS</span>
            </h2>
            <p className="text-gray-400 mt-3 text-lg">Join the action in real-time</p>
          </div>
          <Link href="/streams">
            <Button variant="outline" className="border-white/20 text-white hover:bg-white hover:text-black rounded-full px-6">
              View All Streams
              <ArrowRight className="w-4 h-4 ml-2" />
            </Button>
          </Link>
        </div>

        {/* Streams Grid */}
        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          {liveStreams.map((stream) => (
            <StreamCard key={stream.id} stream={stream} />
          ))}
        </div>
      </div>
    </section>
  );
}
