'use client';

import Link from 'next/link';
import { Smartphone, Watch, Shirt, Home, Gamepad2, BookOpen, MoreHorizontal } from 'lucide-react';

const categories = [
  { name: 'Electronics', icon: Smartphone, slug: 'electronics', count: 128 },
  { name: 'Fashion', icon: Shirt, slug: 'fashion', count: 256 },
  { name: 'Accessories', icon: Watch, slug: 'accessories', count: 89 },
  { name: 'Home', icon: Home, slug: 'home', count: 167 },
  { name: 'Gaming', icon: Gamepad2, slug: 'gaming', count: 74 },
  { name: 'Books', icon: BookOpen, slug: 'books', count: 45 },
];

export function CategoriesSection() {
  return (
    <section className="py-12">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-xl font-bold">Browse Categories</h2>
        <Link
          href="/categories"
          className="text-sm text-primary hover:underline flex items-center"
        >
          View All
          <MoreHorizontal className="ml-1 h-4 w-4" />
        </Link>
      </div>
      
      <div className="grid grid-cols-3 sm:grid-cols-6 gap-4">
        {categories.map((category) => (
          <Link
            key={category.slug}
            href={`/categories/${category.slug}`}
            className="group flex flex-col items-center p-4 rounded-xl bg-card border hover:border-primary hover:shadow-md transition-all"
          >
            <div className="h-14 w-14 rounded-full bg-primary/10 flex items-center justify-center mb-3 group-hover:bg-primary group-hover:text-primary-foreground transition-colors">
              <category.icon className="h-6 w-6" />
            </div>
            <span className="text-sm font-medium text-center">{category.name}</span>
            <span className="text-xs text-muted-foreground">{category.count} items</span>
          </Link>
        ))}
      </div>
    </section>
  );
}
