import React, { useState } from 'react';
import { Package, FileSpreadsheet, Upload, Search, TrendingUp, Map, Users, Bell, Megaphone, Settings, BarChart3, CheckCircle } from 'lucide-react';
import { Button, Badge, Input } from '../UI';

export const DashboardInventory: React.FC = () => {
  const [products, setProducts] = useState(PRODUCTS.slice(0, 12));

  return (
    <div className="animate-in slide-in-from-bottom-4 h-full flex flex-col">
      <h2 className="text-2xl font-display font-bold text-white">Inventory Management</h2>

      {/* Search Bar */}
      <div className="flex justify-between items-center mb-6">
        <div className="relative">
          <input
            type="text"
            placeholder="Search SKU..."
            className="bg-black border border-white/10 pl-8 pr-2 text-sm text-white rounded focus:border-blytz-neon outline-none w-48"
          />
          <Search className="absolute left-2 top-1/2 -translate-y-1/2 text-gray-500" />
        </div>
        <Button size="sm" onClick={() => alert('Add Product functionality coming soon!')}>
          <Plus className="w-4 h-4 mr-1" /> Add Product
        </Button>
      </div>

      {/* Inventory Table */}
      <div className="bg-black/40 border border-white/10 rounded-lg overflow-hidden flex-1">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="border-b border-white/10 text-xs text-gray-500 uppercase tracking-widest bg-white/5">
              <th className="p-4 font-medium">Product</th>
              <th className="p-4 font-medium">Price</th>
              <th className="p-4 font-medium">Stock Level</th>
              <th className="p-4 font-medium">Status</th>
              <th className="p-4 font-medium text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-white/5">
            {products.map(p => (
              <tr key={p.id} className="hover:bg-white/5 transition-colors group">
                <td className="p-4">
                  <div className="flex items-center gap-3">
                    <img src={p.image} className="w-10 h-10 rounded bg-gray-800 object-cover" />
                    <div>
                      <div className="font-bold text-white text-sm">{p.title}</div>
                      <div className="text-xs text-gray-500 font-mono">SKU-{p.id.padStart(4, '0')}</div>
                    </div>
                  </td>
                  <td className="p-4 text-white font-mono">${p.price}</td>
                  <td className="p-4 text-sm text-gray-400">
                     <div className="w-24 bg-gray-800 h-1.5 rounded-full overflow-hidden mb-1">
                       <div className="h-full bg-blytz-neon w-full" style={{ width: `${Math.random() * 100}%`}}"></div>
                     {Math.floor(Math.random() * 50)} units
                  </div>
                  <td className="p-4">
                    <span className={`inline-flex items-center gap-1 px-2 py-1 rounded-full text-[10px] font-bold ${
                      p.stock_quantity > 10 ? 'bg-yellow-500/10 text-yellow-500' : 'bg-green-500/10 text-green-500'
                    }`}>
                      <span className="w-1 h-1 rounded-full bg-white"></span>
                      <span className="text-xs">{p.stock_quantity} units</span>
                    </span>
                  </td>
                  <td className="p-4">
                    <span className="inline-flex items-center gap-1 px-2 py-1 rounded-full bg-white">
                      <span className="w-1 h-1 rounded-full bg-green-500"></span>
                      <span className="text-xs">Active</span>
                    </span>
                  </td>
                  <td className="p-4 text-right">
                    <div className="flex justify-end gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                      <button className="p-1.5 hover:bg-white/10 rounded text-gray-400 hover:text-white">
                        <Edit className="w-4 h-4" />
                      </button>
                      <button className="p-1.5 hover:bg-red-500/10 text-red-500 hover:text-red-500">
                        <Copy className="w-4 h-4" />
                      </button>
                      <button className="p-1.5 hover:bg-red-500/10 text-red-500 hover:text-red-500">
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </div>
                  </td>
              </tr>
            ))}
          </tbody>
        </div>
    </div>
  );
};
