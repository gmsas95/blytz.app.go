import React from 'react';
import { Button } from './UI';
import { DROPS } from '../constants';

export const Drops: React.FC = () => {
  return (
    <div className="container mx-auto px-4 py-12">
      <div className="text-center mb-16">
        <Button variant="flash">INCOMING TRANSMISSIONS</Button>
        <h1 className="text-5xl md:text-7xl font-display font-bold text-white mt-4 mb-4 italic">FUTURE DROPS</h1>
        <p className="text-gray-400 max-w-2xl mx-auto">Get notified before the masses. These limited run items will sell out in seconds.</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
        {DROPS.map(product => (
          <div key={product.id} className="group relative border border-white/10 bg-blytz-dark overflow-hidden hover:border-blytz-neon transition-all">
            <div className="aspect-[4/5] relative">
              <img src={product.image} alt={product.title} className="w-full h-full object-cover opacity-60 group-hover:opacity-80 transition-opacity grayscale group-hover:grayscale-0" />
              <div className="absolute inset-0 flex items-center justify-center">
                <div className="bg-black/80 backdrop-blur border border-white/20 px-6 py-3 transform rotate-[-5deg]">
                  <span className="font-mono text-2xl font-bold text-blytz-neon block text-center">{product.dropDate}</span>
                  <span className="text-xs text-gray-400 uppercase tracking-widest block text-center">Launch Date</span>
                </div>
              </div>
            </div>
            <div className="p-6">
              <h3 className="text-2xl font-bold text-white mb-2">{product.title}</h3>
              <p className="text-gray-400 mb-6 line-clamp-2">{product.description}</p>
              <div className="flex items-center justify-between">
                <span className="text-xl font-mono text-gray-300">${product.price}</span>
                <Button variant="outline" size="sm">Notify Me</Button>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
