import React from 'react';
import { Button } from './UI';
import { useAuthStore } from '../../store/authStore';

interface SellerDashboardProps {
  user: User | null;
}

export const SellerDashboard: React.FC<SellerDashboardProps> = ({ user }) => {
  return (
    <div className="container mx-auto px-4 py-8">
      <div className="text-center mb-12">
        <h1 className="text-4xl font-display font-bold text-white mb-4">SELLER DASHBOARD</h1>
        {user ? (
          <p className="text-gray-400">Welcome back, {user.first_name}</p>
        ) : (
          <Button onClick={() => alert('Login functionality coming soon!')}>Login to Start Selling</Button>
        )}
      </div>

      <div className="max-w-2xl">
        <div className="bg-blytz-dark border border-white/10 p-8 rounded-lg">
          <h2 className="text-2xl font-bold text-white mb-6">List New Product</h2>
          <div className="p-4 text-gray-400 mb-4">
            Listing functionality coming soon!
          </div>
        </div>
      </div>
    </div>
  );
};
