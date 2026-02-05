'use client';

import Link from 'next/link';
import Image from 'next/image';
import { useCountdown } from '@/hooks/use-countdown';
import { Play, Users, ArrowUpRight } from 'lucide-react';

interface AuctionCardProps {
  auction: {
    id: string;
    title: string;
    description: string;
    currentPrice: number;
    startingPrice: number;
    image: string;
    seller: {
      name: string;
      avatar: string;
    };
    endTime: string;
    bidCount: number;
    isLive?: boolean;
  };
}

export function AuctionCard({ auction }: AuctionCardProps) {
  const { days, hours, minutes, isExpired } = useCountdown(auction.endTime);

  const formatTime = () => {
    if (isExpired) return 'Ended';
    if (days > 0) return `${days}d ${hours}h`;
    if (hours > 0) return `${hours}h ${minutes}m`;
    return `${minutes}m`;
  };

  return (
    <Link href={`/auctions/${auction.id}`}>
      <div className="group relative bg-neutral-900 rounded-3xl overflow-hidden border border-neutral-800 hover:border-blytz-yellow/50 transition-all duration-300">
        {/* Image Container */}
        <div className="relative aspect-[4/3] overflow-hidden">
          <Image
            src={auction.image}
            alt={auction.title}
            fill
            className="object-cover group-hover:scale-110 transition-transform duration-500"
          />
          
          {/* Gradient Overlay */}
          <div className="absolute inset-0 bg-gradient-to-t from-black via-transparent to-transparent" />
          
          {/* Live Badge */}
          {auction.isLive && (
            <div className="absolute top-4 left-4 px-3 py-1.5 bg-red-500 text-white text-xs font-bold rounded-full flex items-center gap-1.5">
              <span className="w-1.5 h-1.5 bg-white rounded-full animate-pulse" />
              LIVE
            </div>
          )}
          
          {/* Arrow Icon */}
          <div className="absolute top-4 right-4 w-10 h-10 bg-white/10 backdrop-blur-sm rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
            <ArrowUpRight className="w-5 h-5 text-white" />
          </div>
          
          {/* Price Tag */}
          <div className="absolute bottom-4 left-4 right-4">
            <div className="flex items-end justify-between">
              <div>
                <p className="text-gray-400 text-xs mb-1">Current Bid</p>
                <p className="text-2xl font-black text-blytz-yellow">
                  RM {auction.currentPrice.toLocaleString()}
                </p>
              </div>
              <div className="text-right">
                <p className="text-gray-400 text-xs mb-1">{isExpired ? 'Status' : 'Ends in'}</p>
                <p className={`text-sm font-bold ${isExpired ? 'text-red-500' : 'text-white'}`}>
                  {formatTime()}
                </p>
              </div>
            </div>
          </div>
        </div>
        
        {/* Content */}
        <div className="p-5">
          <h3 className="font-bold text-white text-lg mb-2 line-clamp-1 group-hover:text-blytz-yellow transition-colors">
            {auction.title}
          </h3>
          <p className="text-gray-500 text-sm mb-4 line-clamp-2">
            {auction.description}
          </p>
          
          {/* Seller & Bids */}
          <div className="flex items-center justify-between pt-4 border-t border-neutral-800">
            <div className="flex items-center gap-2">
              <Image
                src={auction.seller.avatar}
                alt={auction.seller.name}
                width={28}
                height={28}
                className="rounded-full"
              />
              <span className="text-gray-400 text-sm">{auction.seller.name}</span>
            </div>
            <div className="flex items-center gap-1.5 text-gray-400 text-sm">
              <Users className="w-4 h-4" />
              <span>{auction.bidCount} bids</span>
            </div>
          </div>
        </div>
        
        {/* Hover Glow */}
        <div className="absolute inset-0 rounded-3xl bg-blytz-yellow/5 opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none" />
      </div>
    </Link>
  );
}
