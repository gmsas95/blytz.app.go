import React from 'react';
import { LayoutDashboard, FileSpreadsheet, BarChart3, Search, TrendingUp, Map, Users, Bell, Megaphone, Settings, Package, DollarSign, Zap, CheckCircle } from 'lucide-react';
import { Button, Badge, Input } from '../UI';

interface DashboardOverviewProps {
  activeTab: string;
  onTabChange: (tab: string) => void;
}

export const DashboardOverview: React.FC<DashboardOverviewProps> = ({ activeTab, onTabChange }) => {
  return (
    <div className="space-y-6 animate-in slide-in-from-right-4">
      <h2 className="text-2xl font-display font-bold text-white mb-6">Dashboard Overview</h2>

      {/* Stats Row */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-black/40 border border-white/10 p-6 rounded-lg relative overflow-hidden group hover:border-blytz-neon/50 transition-all">
          <div className="flex justify-between items-start mb-4">
            <div>
              <p className="text-gray-400 text-xs font-bold uppercase tracking-widest">Total Revenue</p>
              <h3 className="text-3xl font-bold text-white mt-1">$12,450.00</h3>
            </div>
            <Zap className="text-blytz-neon w-6 h-6" />
          </div>
          <div className="text-xs text-green-400 flex items-center gap-1">
            <TrendingUp className="w-3 h-3" /> +12% this week
          </div>
        </div>

        <div className="bg-black/40 border border-white/10 p-6 rounded-lg relative overflow-hidden group hover:border-blytz-neon/50 transition-all">
          <div className="flex justify-between items-start mb-4">
            <div>
              <p className="text-gray-400 text-xs font-bold uppercase tracking-widest">Active Orders</p>
              <h3 className="text-3xl font-bold text-white mt-1">24</h3>
            </div>
            <Package className="text-blue-500 w-6 h-6" />
          </div>
          <div className="text-xs text-gray-400">
            5 items pending dispatch
          </div>
        </div>

        <div className="bg-black/40 border border-white/10 p-6 rounded-lg relative overflow-hidden group hover:border-blytz-neon/50 transition-all">
          <div className="flex justify-between items-start mb-4">
            <div>
              <p className="text-gray-400 text-xs font-bold uppercase tracking-widest">Seller Rating</p>
              <h3 className="text-3xl font-bold text-white mt-1">4.9<span className="text-sm text-gray-500">/5.0</span></h3>
            </div>
            <ShieldCheck className="text-purple-500 w-6 h-6" />
          </div>
          <div className="text-xs text-purple-400">
            Top Rated Seller
          </div>
        </div>
      </div>

      {/* Recent Activity Graph Placeholder */}
      <div className="bg-black/40 border border-white/10 rounded-lg p-6">
        <div className="flex justify-between items-center mb-6">
          <h3 className="font-bold text-white flex items-center gap-2"><BarChart3 className="w-4 h-4 text-blytz-neon" /> Sales Performance</h3>
          <select className="bg-black border border-white/10 text-xs text-white p-1 rounded">
            <option>Last 7 Days</option>
            <option>Last 30 Days</option>
          </select>
        </div>
        <div className="h-48 flex items-end gap-2 justify-between px-2">
           {[30, 45, 25, 60, 75, 50, 80, 40, 55, 70, 65, 90].map((h, i) => (
            <div key={i} className="w-full bg-white/5 hover:bg-blytz-neon transition-colors rounded-t" style={{height: `${h}%`}}></div>
           ))}
        </div>
        <div className="flex justify-between mt-2 text-xs text-gray-500 font-mono">
           <span>Mon</span>
           <span>Tue</span>
           <span>Wed</span>
           <span>Thu</span>
           <span>Fri</span>
           <span>Sat</span>
           <span>Sun</span>
         </div>
      </div>
    </div>
  );
};
