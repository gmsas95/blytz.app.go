import React from 'react';
import { X, Minus, Plus, Trash2, ArrowRight, Badge } from 'lucide-react';
import { Product } from '../types';
import { Button, Input } from './UI';

interface ProductDetailProps {
  product: Product;
  onBack: () => void;
  onAddToCart: (product: Product) => void;
}

export const ProductDetail: React.FC<ProductDetailProps> = ({ product, onBack, onAddToCart }) => {
  return (
    <div className="container mx-auto px-4 py-8 animate-in slide-in-from-right-8 duration-300">
      <button
        onClick={onBack}
        className="mb-6 text-gray-400 hover:text-white flex items-center gap-2"
      >
        ← Back to Marketplace
      </button>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-12">
        <div className="space-y-4">
          <div className="aspect-square bg-blytz-dark rounded-lg overflow-hidden border border-white/10 relative group">
            <img
              src={product.image}
              alt={product.title}
              className="w-full h-full object-cover"
            />
            {product.isFlash && (
              <div className="absolute top-4 left-4">
                <Badge variant="flash">Flash Deal Ends {product.timeLeft}</Badge>
              </div>
            )}
          </div>

          <div className="grid grid-cols-4 gap-4">
            {[1, 2, 3, 4].map((i) => (
              <div key={i} className="aspect-square bg-blytz-dark rounded border border-white/10 cursor-pointer hover:border-blytz-neon"></div>
            ))}
          </div>
        </div>

        <div className="flex flex-col h-full">
          <div className="mb-auto">
            <h1 className="text-4xl md:text-5xl font-display font-bold text-white mb-2 italic">
              {product.title}
            </h1>
            <div className="flex items-center gap-4 mb-6">
              <span className="text-3xl font-bold text-blytz-neon">
                ${product.price.toFixed(2)}
              </span>
              {product.originalPrice && (
                <span className="text-xl text-gray-500 line-through">
                  ${product.originalPrice.toFixed(2)}
                </span>
              )}
              <div className="flex items-center gap-1 text-yellow-400 ml-4">
                <span className="font-bold">{product.rating}</span>
                <div className="flex">
                  {[...Array(5)].map((_, i) => (
                    <span key={i} className="w-4 h-4">★</span>
                  ))}
                </div>
                <span className="text-gray-400 text-sm ml-2">({product.reviews} verified)</span>
              </div>
            </div>
          </div>

          <p className="text-gray-300 text-lg leading-relaxed mb-8">
            {product.description}
          </p>

          <div className="space-y-6 mb-8">
            <div className="p-4 bg-white/5 border border-white/10 rounded">
              <h3 className="text-white font-bold mb-2">Delivery Information</h3>
              <p className="text-sm text-gray-400">Order in the next 2 hours for delivery by tomorrow, 10 AM.</p>
            </div>
          </div>
        </div>

        <div className="flex gap-4 mt-8 pt-8 border-t border-white/10">
          <Button
            variant="primary"
            size="lg"
            className="flex-1 text-lg"
            onClick={() => onAddToCart(product)}
          >
            Add To Cart
          </Button>
          <Button variant="outline" size="lg" className="px-4">
            <Trash2 className="w-6 h-6" />
          </Button>
        </div>
      </div>
    </div>
  );
};
