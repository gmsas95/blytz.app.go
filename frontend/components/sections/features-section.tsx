'use client';

import { Video, Gavel, Shield, Truck, Zap, Trophy } from 'lucide-react';

const features = [
  {
    icon: Video,
    title: 'Live Streaming',
    description: 'Watch sellers showcase products in real-time with crystal clear video quality',
    color: 'bg-blytz-yellow',
  },
  {
    icon: Gavel,
    title: 'Real-time Bidding',
    description: 'Place bids instantly with our lightning-fast auction system',
    color: 'bg-white',
  },
  {
    icon: Shield,
    title: 'Secure Payments',
    description: 'Your transactions are protected with bank-level security',
    color: 'bg-blytz-yellow',
  },
  {
    icon: Truck,
    title: 'Fast Delivery',
    description: 'Get your winnings delivered quickly with our trusted partners',
    color: 'bg-white',
  },
  {
    icon: Zap,
    title: 'Instant Wins',
    description: 'Know immediately when you win with real-time notifications',
    color: 'bg-blytz-yellow',
  },
  {
    icon: Trophy,
    title: 'Rare Finds',
    description: 'Discover unique items you won\'t find anywhere else',
    color: 'bg-white',
  },
];

export function FeaturesSection() {
  return (
    <section className="section-padding bg-neutral-950">
      <div className="container-modern">
        {/* Header */}
        <div className="text-center mb-16">
          <p className="text-blytz-yellow font-semibold mb-4 tracking-wider uppercase">Why Blytz</p>
          <h2 className="text-4xl md:text-5xl font-black text-white mb-6">
            THE FUTURE OF<br />
            <span className="text-blytz-yellow">AUCTIONS</span>
          </h2>
          <p className="text-gray-400 max-w-2xl mx-auto text-lg">
            We&apos;ve revolutionized the auction experience with cutting-edge technology and user-first design
          </p>
        </div>

        {/* Features Grid */}
        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          {features.map((feature, index) => (
            <div
              key={feature.title}
              className="group relative p-8 rounded-3xl bg-neutral-900 border border-neutral-800 hover:border-blytz-yellow/50 transition-all duration-300 hover:transform hover:-translate-y-1"
              style={{ animationDelay: `${index * 0.1}s` }}
            >
              {/* Icon */}
              <div className={`w-14 h-14 rounded-2xl ${feature.color} flex items-center justify-center mb-6 group-hover:scale-110 transition-transform`}>
                <feature.icon className={`w-7 h-7 ${feature.color === 'bg-blytz-yellow' ? 'text-black' : 'text-black'}`} />
              </div>

              {/* Content */}
              <h3 className="text-xl font-bold text-white mb-3">{feature.title}</h3>
              <p className="text-gray-400 leading-relaxed">{feature.description}</p>

              {/* Hover Glow */}
              <div className="absolute inset-0 rounded-3xl bg-blytz-yellow/5 opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none" />
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
