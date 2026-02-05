'use client';

import { Video, Gavel, Shield, Truck } from 'lucide-react';

const features = [
  {
    icon: Video,
    title: 'Live Streaming',
    description: 'Watch sellers showcase products in real-time with interactive Q&A',
  },
  {
    icon: Gavel,
    title: 'Real-time Bidding',
    description: 'Bid instantly during live streams with automatic bid placement',
  },
  {
    icon: Shield,
    title: 'Secure Payments',
    description: 'Protected transactions with Stripe integration and buyer protection',
  },
  {
    icon: Truck,
    title: 'Fast Delivery',
    description: 'Integrated logistics with NinjaVan for reliable shipping',
  },
];

export function FeaturesSection() {
  return (
    <section className="py-12">
      <div className="text-center mb-10">
        <h2 className="text-2xl font-bold mb-2">Why Choose Blytz</h2>
        <p className="text-muted-foreground">
          The most exciting way to shop and sell online
        </p>
      </div>
      
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        {features.map((feature) => (
          <div
            key={feature.title}
            className="group p-6 rounded-xl bg-card border hover:shadow-lg transition-all duration-300"
          >
            <div className="h-12 w-12 rounded-lg bg-primary/10 flex items-center justify-center mb-4 group-hover:bg-primary group-hover:text-primary-foreground transition-colors">
              <feature.icon className="h-6 w-6 text-primary group-hover:text-primary-foreground" />
            </div>
            <h3 className="font-semibold mb-2">{feature.title}</h3>
            <p className="text-sm text-muted-foreground">{feature.description}</p>
          </div>
        ))}
      </div>
    </section>
  );
}
