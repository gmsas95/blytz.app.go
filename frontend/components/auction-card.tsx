'use client';

import Link from 'next/link';
import Image from 'next/image';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Play, Clock, Users, Heart } from 'lucide-react';
import { useCountdown } from '@/hooks/use-countdown';

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
  const { hours, minutes, seconds, isExpired } = useCountdown(auction.endTime);

  return (
    <Card className="group overflow-hidden hover:shadow-lg transition-all duration-300">
      <div className="relative aspect-[4/3] overflow-hidden bg-muted">
        <Image
          src={auction.image}
          alt={auction.title}
          fill
          className="object-cover group-hover:scale-105 transition-transform duration-500"
        />
        
        {auction.isLive && (
          <Badge className="absolute top-3 left-3 bg-blytz-red hover:bg-blytz-red text-white animate-pulse-live">
            <Play className="mr-1 h-3 w-3 fill-current" />
            LIVE
          </Badge>
        )}
        
        <Button
          size="icon"
          variant="secondary"
          className="absolute top-3 right-3 h-8 w-8 opacity-0 group-hover:opacity-100 transition-opacity"
        >
          <Heart className="h-4 w-4" />
        </Button>
        
        <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/60 to-transparent p-3">
          <div className="flex items-center text-white text-sm">
            <Clock className="mr-1 h-4 w-4" />
            {isExpired ? (
              'Ended'
            ) : (
              <span className="font-mono">
                {hours}h {minutes}m {seconds}s
              </span>
            )}
          </div>
        </div>
      </div>
      
      <CardContent className="p-4">
        <Link href={`/auctions/${auction.id}`}>
          <h3 className="font-semibold line-clamp-1 group-hover:text-primary transition-colors">
            {auction.title}
          </h3>
        </Link>
        
        <p className="text-sm text-muted-foreground line-clamp-2 mt-1">
          {auction.description}
        </p>
        
        <div className="flex items-center justify-between mt-3">
          <div>
            <p className="text-xs text-muted-foreground">Current Bid</p>
            <p className="text-lg font-bold text-primary">
              RM {auction.currentPrice.toLocaleString()}
            </p>
          </div>
          
          <div className="text-right">
            <p className="text-xs text-muted-foreground">Bids</p>
            <p className="text-sm font-medium flex items-center">
              <Users className="mr-1 h-3 w-3" />
              {auction.bidCount}
            </p>
          </div>
        </div>
        
        <div className="flex items-center mt-4 pt-3 border-t">
          <Avatar className="h-6 w-6 mr-2">
            <AvatarImage src={auction.seller.avatar} />
            <AvatarFallback>{auction.seller.name[0]}</AvatarFallback>
          </Avatar>
          <span className="text-xs text-muted-foreground truncate">
            {auction.seller.name}
          </span>
        </div>
      </CardContent>
    </Card>
  );
}
