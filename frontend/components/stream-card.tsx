'use client';

import Link from 'next/link';
import Image from 'next/image';
import { Play, Users, Package, Star } from 'lucide-react';

interface StreamCardProps {
  stream: {
    id: string;
    title: string;
    seller: {
      name: string;
      avatar: string;
      rating: number;
    };
    thumbnail: string;
    viewerCount: number;
    productCount: number;
    startedAt: string;
  };
}

export function StreamCard({ stream }: StreamCardProps) {
  const formatDuration = (startedAt: string) => {
    const started = new Date(startedAt);
    const now = new Date();
    const diff = Math.floor((now.getTime() - started.getTime()) / 60000);
    if (diff < 60) return `${diff}m`;
    const hours = Math.floor(diff / 60);
    const mins = diff % 60;
    return `${hours}h ${mins}m`;
  };

  return (
    <Link href={`/streams/${stream.id}`}>
      <div className="group relative bg-[#0a0a0a] rounded-3xl overflow-hidden border border-neutral-800 hover:border-blytz-yellow/50 transition-all duration-300">
        {/* Thumbnail */}
        <div className="relative aspect-video overflow-hidden">
          <Image
            src={stream.thumbnail}
            alt={stream.title}
            fill
            className="object-cover group-hover:scale-110 transition-transform duration-500"
          />
          
          {/* Overlay */}
          <div className="absolute inset-0 bg-gradient-to-t from-black via-black/20 to-transparent" />
          
          {/* Play Button */}
          <div className="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
            <div className="w-16 h-16 rounded-full bg-blytz-yellow flex items-center justify-center transform scale-75 group-hover:scale-100 transition-transform">
              <Play className="w-6 h-6 text-black fill-current ml-1" />
            </div>
          </div>
          
          {/* Live Badge */}
          <div className="absolute top-4 left-4 px-3 py-1 bg-red-500 text-white text-xs font-bold rounded-full flex items-center gap-1.5">
            <span className="w-1.5 h-1.5 bg-white rounded-full animate-pulse" />
            LIVE
          </div>
          
          {/* Viewers */}
          <div className="absolute top-4 right-4 px-3 py-1 bg-black/60 backdrop-blur-sm text-white text-xs font-semibold rounded-full flex items-center gap-1.5">
            <Users className="w-3.5 h-3.5" />
            {stream.viewerCount.toLocaleString()}
          </div>
          
          {/* Duration */}
          <div className="absolute bottom-4 left-4 px-2 py-1 bg-black/60 text-white text-xs rounded">
            {formatDuration(stream.startedAt)}
          </div>
        </div>
        
        {/* Content */}
        <div className="p-5">
          <h3 className="font-bold text-white text-lg mb-4 line-clamp-1 group-hover:text-blytz-yellow transition-colors">
            {stream.title}
          </h3>
          
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="relative">
                <Image
                  src={stream.seller.avatar}
                  alt={stream.seller.name}
                  width={40}
                  height={40}
                  className="rounded-full border-2 border-neutral-800"
                />
                <div className="absolute -bottom-1 -right-1 w-4 h-4 bg-blytz-yellow rounded-full flex items-center justify-center">
                  <Star className="w-2.5 h-2.5 text-black fill-current" />
                </div>
              </div>
              <div>
                <p className="font-semibold text-white text-sm">{stream.seller.name}</p>
                <p className="text-blytz-yellow text-xs font-medium">{stream.seller.rating} Rating</p>
              </div>
            </div>
            
            <div className="flex items-center gap-1.5 text-gray-400 text-sm">
              <Package className="w-4 h-4" />
              <span>{stream.productCount}</span>
            </div>
          </div>
        </div>
      </div>
    </Link>
  );
}
