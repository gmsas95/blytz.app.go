'use client';

import Link from 'next/link';
import Image from 'next/image';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
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
    <Card className="group overflow-hidden hover:shadow-lg transition-all duration-300">
      <Link href={`/streams/${stream.id}`}>
        <div className="relative aspect-video overflow-hidden bg-muted">
          <Image
            src={stream.thumbnail}
            alt={stream.title}
            fill
            className="object-cover group-hover:scale-105 transition-transform duration-500"
          />
          
          <div className="absolute inset-0 bg-black/30 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
            <div className="h-16 w-16 rounded-full bg-blytz-red flex items-center justify-center">
              <Play className="h-8 w-8 text-white fill-current ml-1" />
            </div>
          </div>
          
          <Badge className="absolute top-3 left-3 bg-blytz-red hover:bg-blytz-red text-white animate-pulse-live">
            LIVE
          </Badge>
          
          <div className="absolute top-3 right-3 bg-black/60 text-white text-xs px-2 py-1 rounded-full flex items-center">
            <Users className="mr-1 h-3 w-3" />
            {stream.viewerCount.toLocaleString()}
          </div>
          
          <div className="absolute bottom-3 left-3 bg-black/60 text-white text-xs px-2 py-1 rounded">
            {formatDuration(stream.startedAt)}
          </div>
        </div>
      </Link>
      
      <CardContent className="p-4">
        <Link href={`/streams/${stream.id}`}>
          <h3 className="font-semibold line-clamp-1 group-hover:text-primary transition-colors">
            {stream.title}
          </h3>
        </Link>
        
        <div className="flex items-center justify-between mt-3">
          <div className="flex items-center">
            <Avatar className="h-8 w-8 mr-2">
              <AvatarImage src={stream.seller.avatar} />
              <AvatarFallback>{stream.seller.name[0]}</AvatarFallback>
            </Avatar>
            <div>
              <p className="text-sm font-medium leading-none">{stream.seller.name}</p>
              <div className="flex items-center text-xs text-muted-foreground mt-0.5">
                <Star className="mr-0.5 h-3 w-3 fill-yellow-400 text-yellow-400" />
                {stream.seller.rating}
              </div>
            </div>
          </div>
          
          <div className="flex items-center text-xs text-muted-foreground">
            <Package className="mr-1 h-3 w-3" />
            {stream.productCount} items
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
