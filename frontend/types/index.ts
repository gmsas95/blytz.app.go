// User types
export interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  avatar_url?: string;
  phone?: string;
  role: 'buyer' | 'seller' | 'admin';
  created_at: string;
}

export interface Seller {
  id: string;
  user_id: string;
  store_name: string;
  store_slug: string;
  description?: string;
  logo_url?: string;
  banner_url?: string;
  rating: number;
  review_count: number;
  verified: boolean;
  total_sales: number;
}

// Product types
export interface Product {
  id: string;
  seller_id: string;
  category_id?: string;
  name: string;
  slug: string;
  description: string;
  condition: 'new' | 'used' | 'refurbished';
  base_price: number;
  compare_at_price?: number;
  stock_quantity: number;
  sku?: string;
  weight_grams?: number;
  dimensions_cm?: {
    length: number;
    width: number;
    height: number;
  };
  attributes: Record<string, string>;
  images: ProductImage[];
  status: 'draft' | 'active' | 'sold' | 'archived';
  view_count: number;
  created_at: string;
  updated_at: string;
  seller?: Seller;
  category?: Category;
}

export interface ProductImage {
  id: string;
  url: string;
  thumbnail_url?: string;
  alt_text?: string;
  sort_order: number;
  is_primary: boolean;
}

export interface Category {
  id: string;
  name: string;
  slug: string;
  description?: string;
  image_url?: string;
  parent_id?: string;
  children?: Category[];
  sort_order: number;
  is_active: boolean;
  product_count: number;
}

// Auction types
export interface Auction {
  id: string;
  product_id: string;
  seller_id: string;
  title: string;
  description?: string;
  start_time: string;
  end_time: string;
  status: 'scheduled' | 'live' | 'ended' | 'cancelled';
  start_price: number;
  reserve_price?: number;
  buy_now_price?: number;
  current_bid?: Bid;
  bid_count: number;
  winner_id?: string;
  product?: Product;
  seller?: Seller;
  is_live?: boolean;
  viewer_count?: number;
}

export interface Bid {
  id: string;
  auction_id: string;
  bidder_id: string;
  amount: number;
  is_auto_bid: boolean;
  is_winning: boolean;
  created_at: string;
  bidder?: User;
}

// Order types
export interface Order {
  id: string;
  order_number: string;
  buyer_id: string;
  seller_id: string;
  auction_id?: string;
  status: 'pending' | 'confirmed' | 'processing' | 'shipped' | 'delivered' | 'cancelled' | 'refunded';
  subtotal: number;
  shipping_cost: number;
  tax_amount: number;
  total_amount: number;
  shipping_address: Address;
  items: OrderItem[];
  payment?: Payment;
  shipment?: Shipment;
  created_at: string;
}

export interface OrderItem {
  id: string;
  product_id: string;
  product_name: string;
  product_image?: string;
  quantity: number;
  unit_price: number;
  total: number;
}

export interface Address {
  id?: string;
  label: string;
  recipient_name: string;
  phone: string;
  address_line1: string;
  address_line2?: string;
  city: string;
  state: string;
  postal_code: string;
  country: string;
  is_default?: boolean;
}

export interface Payment {
  id: string;
  order_id: string;
  method: string;
  amount: number;
  currency: string;
  status: 'pending' | 'processing' | 'completed' | 'failed' | 'refunded';
  stripe_payment_intent_id?: string;
  paid_at?: string;
}

export interface Shipment {
  id: string;
  tracking_number?: string;
  carrier: string;
  status: string;
  label_url?: string;
  shipped_at?: string;
  delivered_at?: string;
  tracking_history: TrackingEvent[];
}

export interface TrackingEvent {
  status: string;
  location: string;
  timestamp: string;
}

// Stream types
export interface Stream {
  id: string;
  seller_id: string;
  title: string;
  description?: string;
  status: 'scheduled' | 'live' | 'ended' | 'cancelled';
  thumbnail_url?: string;
  viewer_count: number;
  max_viewers: number;
  started_at?: string;
  ended_at?: string;
  recording_url?: string;
  products?: Product[];
  seller?: Seller;
}

// Chat types
export interface ChatMessage {
  id: string;
  stream_id: string;
  user_id: string;
  user_name: string;
  user_avatar?: string;
  content: string;
  message_type: 'text' | 'system';
  created_at: string;
}

// WebSocket types
export interface WSMessage {
  type: 'bid' | 'auction_started' | 'auction_ended' | 'auction_extended' | 'viewer_count' | 'chat';
  auction_id: string;
  data: any;
  timestamp: string;
}

// Pagination
export interface PaginatedResponse<T> {
  data: T[];
  total_count: number;
  page: number;
  page_size: number;
}

// Frontend Product type (simplified for UI)
export interface FrontendProduct {
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
  sellerId?: string;
  seller?: {
    id: string;
    first_name: string;
    last_name: string;
    email: string;
  };
}

// Filter types
export interface ProductFilter {
  category_id?: string;
  seller_id?: string;
  status?: string;
  condition?: string;
  min_price?: number;
  max_price?: number;
  q?: string;
  sort?: string;
  page?: number;
  per_page?: number;
}

export interface AuctionFilter {
  status?: string;
  category_id?: string;
  seller_id?: string;
  sort?: string;
  page?: number;
  per_page?: number;
}

// Cart types
export interface CartItem {
  id: string;
  title: string;
  price: number;
  quantity: number;
  image: string;
}

export interface Cart {
  items: CartItem[];
  total: number;
  itemCount: number;
}
