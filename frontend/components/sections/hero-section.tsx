'use client';

import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Play, ArrowRight, TrendingUp, Users } from 'lucide-react';

export function HeroSection() {
  return (
    <section className="relative min-h-[90vh] flex items-center overflow-hidden bg-black">
      {/* Background Pattern */}
      <div className="absolute inset-0">
        <div className="absolute inset-0 bg-[linear-gradient(rgba(255,215,0,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(255,215,0,0.03)_1px,transparent_1px)] bg-[size:60px_60px]" />
        <div className="absolute top-0 left-1/4 w-96 h-96 bg-blytz-yellow/20 rounded-full blur-[150px]" />
        <div className="absolute bottom-0 right-1/4 w-96 h-96 bg-blytz-yellow/10 rounded-full blur-[150px]" />
      </div>

      <div className="container-modern relative z-10">
        <div className="grid lg:grid-cols-2 gap-12 items-center">
          {/* Left Content */}
          <div className="space-y-8 animate-slide-up">
            {/* Badge */}
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-blytz-yellow/10 border border-blytz-yellow/30">
              <span className="w-2 h-2 rounded-full bg-blytz-yellow animate-pulse-live" />
              <span className="text-blytz-yellow text-sm font-medium">Live Auctions Happening Now</span>
            </div>

            {/* Headline */}
            <div className="space-y-4">
              <h1 className="text-5xl md:text-7xl font-black text-white leading-[1.1]">
                BID.<br />
                <span className="text-blytz-yellow">WIN.</span><br />
                REPEAT.
              </h1>
              <p className="text-xl text-gray-400 max-w-md leading-relaxed">
                Malaysia&apos;s most exciting live auction platform. Real-time bidding, authentic products, unbeatable deals.
              </p>
            </div>

            {/* Stats */}
            <div className="flex gap-8">
              <div>
                <p className="text-3xl font-black text-blytz-yellow">50K+</p>
                <p className="text-sm text-gray-500">Active Users</p>
              </div>
              <div>
                <p className="text-3xl font-black text-blytz-yellow">10K+</p>
                <p className="text-sm text-gray-500">Auctions Won</p>
              </div>
              <div>
                <p className="text-3xl font-black text-blytz-yellow">RM2M+</p>
                <p className="text-sm text-gray-500">Value Sold</p>
              </div>
            </div>

            {/* CTA Buttons */}
            <div className="flex flex-wrap gap-4">
              <Link href="/streams">
                <Button size="lg" className="btn-primary gap-2 text-base">
                  <Play className="w-5 h-5 fill-current" />
                  Watch Live
                </Button>
              </Link>
              <Link href="/auctions">
                <Button size="lg" className="btn-outline text-white border-white hover:bg-white hover:text-black">
                  Browse Auctions
                  <ArrowRight className="w-5 h-5 ml-2" />
                </Button>
              </Link>
            </div>
          </div>

          {/* Right Content - Featured Card */}
          <div className="relative hidden lg:block">
            {/* Floating Cards */}
            <div className="relative w-full aspect-square max-w-lg mx-auto">
              {/* Main Card */}
              <div className="absolute inset-0 bg-gradient-to-br from-blytz-yellow to-yellow-600 rounded-3xl transform rotate-3 animate-float" />
              <div className="absolute inset-0 bg-neutral-900 rounded-3xl p-6 flex flex-col">
                <div className="relative flex-1 rounded-2xl overflow-hidden bg-neutral-800">
                  <img 
                    src="https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=600" 
                    alt="Featured Auction"
                    className="w-full h-full object-cover"
                  />
                  <div className="absolute top-4 left-4 px-3 py-1 bg-blytz-yellow text-black text-sm font-bold rounded-full">
                    LIVE NOW
                  </div>
                  <div className="absolute bottom-4 left-4 right-4">
                    <div className="glass-dark rounded-xl p-4">
                      <p className="text-white font-semibold">Limited Edition Sneakers</p>
                      <div className="flex items-center justify-between mt-2">
                        <span className="text-blytz-yellow font-black text-2xl">RM 450</span>
                        <div className="flex items-center gap-1 text-gray-400 text-sm">
                          <Users className="w-4 h-4" />
                          <span>234 bidding</span>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              {/* Floating Elements */}
              <div className="absolute -top-6 -right-6 bg-white rounded-2xl p-4 shadow-2xl animate-float" style={{ animationDelay: '0.5s' }}>
                <div className="flex items-center gap-3">
                  <div className="w-12 h-12 bg-blytz-yellow rounded-full flex items-center justify-center">
                    <TrendingUp className="w-6 h-6 text-black" />
                  </div>
                  <div>
                    <p className="text-sm text-gray-500">Trending Now</p>
                    <p className="font-bold text-black">Vintage Watches</p>
                  </div>
                </div>
              </div>

              <div className="absolute -bottom-4 -left-4 bg-neutral-900 rounded-2xl p-4 border border-neutral-800 animate-float" style={{ animationDelay: '1s' }}>
                <div className="flex items-center gap-2">
                  <div className="flex -space-x-2">
                    {[1,2,3].map(i => (
                      <div key={i} className="w-8 h-8 rounded-full bg-neutral-700 border-2 border-neutral-900" />
                    ))}
                  </div>
                  <p className="text-sm text-gray-400">+1.2k online</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Marquee */}
        <div className="mt-20 border-t border-neutral-800 pt-8 overflow-hidden">
          <div className="flex animate-marquee whitespace-nowrap">
            {[...Array(2)].map((_, i) => (
              <div key={i} className="flex items-center gap-12 mr-12">
                <span className="text-neutral-600 text-lg font-medium">ELECTRONICS</span>
                <span className="text-blytz-yellow">★</span>
                <span className="text-neutral-600 text-lg font-medium">FASHION</span>
                <span className="text-blytz-yellow">★</span>
                <span className="text-neutral-600 text-lg font-medium">COLLECTIBLES</span>
                <span className="text-blytz-yellow">★</span>
                <span className="text-neutral-600 text-lg font-medium">LUXURY</span>
                <span className="text-blytz-yellow">★</span>
                <span className="text-neutral-600 text-lg font-medium">SNEAKERS</span>
                <span className="text-blytz-yellow">★</span>
                <span className="text-neutral-600 text-lg font-medium">WATCHES</span>
                <span className="text-blytz-yellow">★</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
