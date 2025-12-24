import React from 'react';
import { ShoppingBag, X, Trash2, Minus, Plus } from 'lucide-react';
import { useCartStore } from '../../store/cartStore';
import { Button } from './UI';
import { ViewState } from '../types';

interface CartDrawerProps {
  view: ViewState;
  onViewChange: (view: ViewState) => void;
}

export const CartDrawer: React.FC<CartDrawerProps> = ({ view, onViewChange }) => {
  const cart = useCartStore();

  return (
    <>
      <div
        className="fixed inset-0 bg-black/80 backdrop-blur-sm z-50 transition-opacity"
        onClick={() => cart.setIsCartOpen(false)}
      />
      <div className="fixed inset-y-0 right-0 w-full md:w-[450px] bg-blytz-dark border-l border-white/10 z-[60] shadow-2xl flex flex-col animate-in slide-in-from-right duration-300">
        <div className="p-6 border-b border-white/10 flex items-center justify-between bg-blytz-black">
          <h2 className="text-2xl font-display font-bold italic text-white">YOUR HAUL</h2>
          <button onClick={() => cart.setIsCartOpen(false)} className="text-gray-400 hover:text-white">
            <X />
          </button>
        </div>

        <div className="flex-1 overflow-y-auto p-6 space-y-6">
          {cart.isLoading ? (
            <div className="h-full flex items-center justify-center text-gray-500">
              <div className="inline-block w-8 h-8 border-2 border-blytz-neon border-t-transparent rounded-full animate-spin mb-4 mx-auto"></div>
              <p>Loading your haul...</p>
            </div>
          ) : cart.items.length === 0 ? (
            <div className="h-full flex flex-col items-center justify-center text-gray-500">
              <ShoppingBag className="w-16 h-16 mb-4 opacity-20" />
              <p className="text-lg">Your cart is empty.</p>
              <Button
                variant="outline"
                className="mt-4"
                onClick={() => {
                  cart.setIsCartOpen(false);
                  onViewChange('HOME');
                }}
              >
                Start Shopping
              </Button>
            </div>
          ) : (
            cart.items.map(item => (
              <div key={item.id} className="flex gap-4 bg-white/5 p-4 rounded border border-white/5">
                <img src={item.image} alt={item.title} className="w-20 h-20 object-cover rounded bg-black" />
                <div className="flex-1">
                  <div className="flex justify-between items-start mb-1">
                    <h4 className="font-bold text-white line-clamp-1">{item.title}</h4>
                    <button onClick={() => cart.removeItem(item.id)} className="text-gray-500 hover:text-red-500">
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </div>
                  <p className="text-blytz-neon font-mono mb-3">${item.price.toFixed(2)}</p>

                  <div className="flex items-center gap-3">
                    <button
                      className="w-6 h-6 rounded bg-black border border-white/20 flex items-center justify-center hover:border-blytz-neon"
                      onClick={() => cart.updateQuantity(item.id, -1)}
                    >
                      <Minus className="w-3 h-3" />
                    </button>
                    <span className="text-sm font-bold w-4 text-center">{item.quantity}</span>
                    <button
                      className="w-6 h-6 rounded bg-black border border-white/20 flex items-center justify-center hover:border-blytz-neon"
                      onClick={() => cart.updateQuantity(item.id, 1)}
                    >
                      <Plus className="w-3 h-3" />
                    </button>
                  </div>
                </div>
              </div>
            ))
          )}
        </div>

        {cart.items.length > 0 && (
          <div className="p-6 bg-blytz-black border-t border-white/10">
            <div className="flex justify-between items-center mb-2 text-gray-400">
              <span>Subtotal</span>
              <span>${cart.getTotal().toFixed(2)}</span>
            </div>
            <div className="flex justify-between items-center mb-2 text-gray-400">
              <span>Shipping</span>
              <span className="text-blytz-neon">FREE (Blytz Prime)</span>
            </div>
            <div className="flex justify-between items-center mb-6 text-xl font-bold text-white">
              <span>Total</span>
              <span>${cart.getTotal().toFixed(2)}</span>
            </div>
            <Button
              className="w-full h-14 text-lg"
              onClick={() => {
                cart.setIsCartOpen(false);
                onViewChange('CHECKOUT');
              }}
            >
              SECURE CHECKOUT
            </Button>
          </div>
        )}
      </div>
    </>
  );
};
