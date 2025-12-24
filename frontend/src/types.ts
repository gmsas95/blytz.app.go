import React from 'react';

export interface Product {
  id: string;
  title: string;
  price: number;
  originalPrice?: number;
  rating: number;
  reviews: number;
  image: string;
  category: string;
  isFlash?: boolean;
  isHot?: boolean;
  timeLeft?: string;
  description?: string;
  dropDate?: string;
}

export interface ProductFilter {
  category?: string;
  min_price?: number;
  max_price?: number;
  condition?: string;
  status?: string;
  seller_id?: string;
  featured?: boolean;
  search?: string;
  page?: number;
  limit?: number;
  sort_by?: string;
  sort_order?: 'asc' | 'desc';
}

export interface CartItem extends Product {
  quantity: number;
}

export interface Category {
  id: string;
  name: string;
  icon: React.ReactNode;
}

export type ViewState = 'HOME' | 'PRODUCT_DETAIL' | 'CHECKOUT' | 'DROPS' | 'SELL' | 'ACCOUNT';

export interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  role: string;
  avatar_url?: string;
  phone?: string;
  email_verified: boolean;
  last_login_at?: string;
  created_at: string;
  updated_at: string;
}