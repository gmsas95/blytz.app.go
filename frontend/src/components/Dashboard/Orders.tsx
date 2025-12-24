import React, { useState } from 'react';
import { FileSpreadsheet, BarChart3, Search, TrendingUp, Map, Users, Bell, Megaphone, Settings, Package, DollarSign, CheckCircle, Trash2, Plus } from 'lucide-react';
import { Button, Badge, Input } from '../UI';
import { useCartStore } from '../store/cartStore';
import { PRODUCTS } from '../constants';

interface DashboardOrdersProps {
  onBack: () => void;
}

export const DashboardOrders: React.FC<DashboardOrdersProps> = ({ onBack }) => {
  return (
    <div className="animate-in slide-in-from-bottom-4">
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-display font-bold text-white mb-4">Order Management</h2>
        <button onClick={onBack} className="text-gray-400 hover:text-white">
          ← Back to Dashboard
        </button>
      </div>

      <div className="bg-blytz-dark border border-white/10 rounded-lg overflow-hidden flex flex-col h-[600px]">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="border-b border-white/10 text-xs text-gray-500 uppercase tracking-widest bg-white/5">
              <th className="p-4 font-medium">Order ID</th>
              <th className="p-4 font-medium">Customer</th>
              <th className="p-4 font-medium">Items</th>
              <th className="p-4 font-medium">Total</th>
              <th className="p-4 font-medium">Status</th>
              <th className="p-4 text-right">Action</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-white/5">
            {[1,2,3,4,5].map(i => (
              <tr key={i} className="hover:bg-white/5 transition-colors">
                <td className="p-4 text-sm font-mono text-blytz-neon">#ORD-{9000+i}</td>
                <td className="p-4 text-sm text-white">Customer_{i}</td>
                <td className="p-4 text-sm text-gray-400">{i + 1} items</td>
                <td className="text-white font-bold">${(Math.random() * 500).toFixed(2)}</td>
                <td>
                  {i % 2 === 0 ? (
                    <span className="px-2 py-1 rounded bg-yellow-500/20 text-yellow-900 text-xs font-bold border border-yellow-500/30">Pending</span>
                  ) : (
                    <span className="px-2 py-1 rounded bg-green-500/20 text-green-900 text-xs font-bold border-green-500/30">Shipped</span>
                  )}
                </td>
                <td className={`text-xs font-bold ${i % 2 === 0 ? 'text-yellow-500' : 'text-green-500'}`}>
                  {i % 2 === 0 ? 'Pending' : 'Shipped'}
                </td>
                <td className="text-right">
                  {i % 2 === 0 && (
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => alert('Order #' + i + ' fulfillment functionality coming soon!')}
                    >
                      Fulfill
                    </Button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};
