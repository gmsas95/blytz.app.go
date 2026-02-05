'use client';

import { Button } from '@/components/ui/button';
import { Apple, Play, Bell, Zap, Shield } from 'lucide-react';

export function AppDownloadSection() {
  return (
    <section className="section-padding bg-black relative overflow-hidden">
      {/* Background Elements */}
      <div className="absolute inset-0">
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-blytz-yellow/10 rounded-full blur-[100px]" />
      </div>

      <div className="container-modern relative z-10">
        <div className="relative bg-gradient-to-br from-neutral-900 to-black rounded-[3rem] p-8 md:p-16 border border-neutral-800 overflow-hidden">
          {/* Yellow Accent */}
          <div className="absolute top-0 right-0 w-64 h-64 bg-blytz-yellow/20 rounded-full blur-[80px] -translate-y-1/2 translate-x-1/2" />
          
          <div className="grid lg:grid-cols-2 gap-12 items-center">
            {/* Left Content */}
            <div className="space-y-8">
              <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-blytz-yellow/10 border border-blytz-yellow/30">
                <Zap className="w-4 h-4 text-blytz-yellow" />
                <span className="text-blytz-yellow text-sm font-semibold">Download Now</span>
              </div>

              <div>
                <h2 className="text-4xl md:text-5xl font-black text-white mb-6 leading-tight">
                  BID ON THE<br />
                  <span className="text-blytz-yellow">GO</span>
                </h2>
                <p className="text-gray-400 text-lg max-w-md">
                  Never miss a winning bid. Get real-time notifications, track your auctions, and bid from anywhere with our mobile app.
                </p>
              </div>

              {/* Features */}
              <div className="flex flex-wrap gap-4">
                <div className="flex items-center gap-2 text-gray-300">
                  <Bell className="w-5 h-5 text-blytz-yellow" />
                  <span className="text-sm">Push Notifications</span>
                </div>
                <div className="flex items-center gap-2 text-gray-300">
                  <Zap className="w-5 h-5 text-blytz-yellow" />
                  <span className="text-sm">Real-time Bidding</span>
                </div>
                <div className="flex items-center gap-2 text-gray-300">
                  <Shield className="w-5 h-5 text-blytz-yellow" />
                  <span className="text-sm">Secure Payments</span>
                </div>
              </div>

              {/* Download Buttons */}
              <div className="flex flex-wrap gap-4">
                <Button size="lg" className="bg-white text-black hover:bg-gray-100 rounded-full px-6 h-14">
                  <Apple className="w-6 h-6 mr-2" />
                  <div className="text-left">
                    <div className="text-xs text-gray-600">Download on the</div>
                    <div className="text-sm font-bold">App Store</div>
                  </div>
                </Button>
                <Button size="lg" className="bg-white text-black hover:bg-gray-100 rounded-full px-6 h-14">
                  <Play className="w-6 h-6 mr-2 fill-current" />
                  <div className="text-left">
                    <div className="text-xs text-gray-600">Get it on</div>
                    <div className="text-sm font-bold">Google Play</div>
                  </div>
                </Button>
              </div>
            </div>

            {/* Right - Phone Mockup */}
            <div className="hidden lg:flex justify-center">
              <div className="relative">
                {/* Phone Frame */}
                <div className="w-72 h-[580px] bg-black rounded-[3rem] border-4 border-neutral-800 p-3 shadow-2xl">
                  {/* Screen */}
                  <div className="w-full h-full bg-neutral-900 rounded-[2.5rem] overflow-hidden relative">
                    {/* App Header */}
                    <div className="bg-blytz-yellow p-6 pt-12">
                      <div className="flex items-center justify-between">
                        <div className="w-10 h-10 bg-black rounded-full flex items-center justify-center">
                          <span className="text-blytz-yellow font-black text-lg">B</span>
                        </div>
                        <Bell className="w-5 h-5 text-black" />
                      </div>
                      <p className="text-black font-bold text-xl mt-4">Live Auctions</p>
                    </div>
                    
                    {/* App Content */}
                    <div className="p-4 space-y-3">
                      {[1, 2, 3].map((i) => (
                        <div key={i} className="bg-neutral-800 rounded-xl p-3 flex gap-3">
                          <div className="w-16 h-16 bg-neutral-700 rounded-lg" />
                          <div className="flex-1 space-y-2">
                            <div className="h-3 bg-neutral-700 rounded w-3/4" />
                            <div className="h-3 bg-neutral-700 rounded w-1/2" />
                            <div className="h-4 bg-blytz-yellow rounded w-1/3" />
                          </div>
                        </div>
                      ))}
                    </div>

                    {/* Floating Badge */}
                    <div className="absolute bottom-6 left-4 right-4 bg-blytz-yellow rounded-xl p-3 flex items-center justify-between">
                      <span className="text-black font-bold text-sm">You won!</span>
                      <span className="text-black text-xs">RM 450</span>
                    </div>
                  </div>
                </div>

                {/* Decorative Elements */}
                <div className="absolute -top-4 -right-4 w-20 h-20 bg-blytz-yellow rounded-full flex items-center justify-center animate-float">
                  <span className="text-black font-black text-2xl">B</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
