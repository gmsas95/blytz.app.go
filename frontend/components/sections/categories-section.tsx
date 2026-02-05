'use client';

import Link from 'next/link';
import { Smartphone, Shirt, Watch, Home, Gamepad2, BookOpen, Gem, Camera } from 'lucide-react';

const categories = [
  { name: 'Electronics', icon: Smartphone, slug: 'electronics', count: 128, color: 'from-yellow-400 to-yellow-600' },
  { name: 'Fashion', icon: Shirt, slug: 'fashion', count: 256, color: 'from-gray-700 to-black' },
  { name: 'Watches', icon: Watch, slug: 'watches', count: 89, color: 'from-yellow-500 to-yellow-700' },
  { name: 'Home', icon: Home, slug: 'home', count: 167, color: 'from-neutral-600 to-neutral-800' },
  { name: 'Gaming', icon: Gamepad2, slug: 'gaming', count: 74, color: 'from-purple-600 to-purple-800' },
  { name: 'Collectibles', icon: Gem, slug: 'collectibles', count: 145, color: 'from-yellow-600 to-yellow-800' },
  { name: 'Books', icon: BookOpen, slug: 'books', count: 45, color: 'from-amber-600 to-amber-800' },
  { name: 'Cameras', icon: Camera, slug: 'cameras', count: 67, color: 'from-slate-600 to-slate-800' },
];

export function CategoriesSection() {
  return (
    <section className="section-padding bg-black">
      <div className="container-modern">
        {/* Header */}
        <div className="text-center mb-16">
          <p className="text-blytz-yellow font-semibold mb-4 tracking-wider uppercase">Categories</p>
          <h2 className="text-4xl md:text-5xl font-black text-white mb-6">
            BROWSE BY <span className="text-blytz-yellow">INTEREST</span>
          </h2>
          <p className="text-gray-400 max-w-2xl mx-auto text-lg">
            Find exactly what you&apos;re looking for across our diverse categories
          </p>
        </div>

        {/* Categories Grid */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          {categories.map((category) => (
            <Link key={category.slug} href={`/categories/${category.slug}`}>
              <div className="group relative p-6 rounded-3xl bg-[#0a0a0a] border border-neutral-800 hover:border-blytz-yellow/30 transition-all duration-300 hover:transform hover:-translate-y-1 overflow-hidden">
                {/* Background Gradient */}
                <div className={`absolute inset-0 bg-gradient-to-br ${category.color} opacity-0 group-hover:opacity-10 transition-opacity`} />
                
                {/* Icon */}
                <div className="relative w-14 h-14 rounded-2xl bg-neutral-800 flex items-center justify-center mb-4 group-hover:bg-blytz-yellow transition-colors">
                  <category.icon className="w-7 h-7 text-gray-400 group-hover:text-black transition-colors" />
                </div>
                
                {/* Content */}
                <h3 className="relative font-bold text-white text-lg mb-1 group-hover:text-blytz-yellow transition-colors">
                  {category.name}
                </h3>
                <p className="relative text-gray-500 text-sm">
                  {category.count} items
                </p>
              </div>
            </Link>
          ))}
        </div>
      </div>
    </section>
  );
}
