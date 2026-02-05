# Database Schema

## Overview

PostgreSQL 17+ database schema for Blytz livestream ecommerce platform.

## Entity Relationship Diagram

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│     users       │────▶│    sellers      │     │    buyers       │
├─────────────────┤     ├─────────────────┤     ├─────────────────┤
│ id (PK)         │     │ id (PK)         │     │ id (PK)         │
│ email           │     │ user_id (FK)    │     │ user_id (FK)    │
│ password_hash   │     │ store_name      │     │ shipping_addr   │
│ role            │     │ description     │     │ preferences     │
│ status          │     │ rating          │     └─────────────────┘
│ created_at      │     │ verified        │
└─────────────────┘     └─────────────────┘
         │                       │
         │              ┌────────┴────────┐
         │              │                 │
         ▼              ▼                 ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   followers     │  │    products     │  │    streams      │
├─────────────────┤  ├─────────────────┤  ├─────────────────┤
│ id (PK)         │  │ id (PK)         │  │ id (PK)         │
│ user_id (FK)    │  │ seller_id (FK)  │  │ seller_id (FK)  │
│ seller_id (FK)  │  │ category_id(FK) │  │ title           │
│ created_at      │  │ name            │  │ status          │
└─────────────────┘  │ description     │  │ started_at      │
                     │ price           │  │ ended_at        │
                     │ stock_qty       │  │ recording_url   │
                     │ status          │  └─────────────────┘
                     └────────┬────────┘           │
                              │                    │
              ┌───────────────┼────────────────────┤
              │               │                    │
              ▼               ▼                    ▼
       ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
       │ product_images  │ │    auctions     │ │  stream_views   │
       ├─────────────────┤ ├─────────────────┤ ├─────────────────┤
       │ id (PK)         │ │ id (PK)         │ │ id (PK)         │
       │ product_id (FK) │ │ product_id (FK) │ │ stream_id (FK)  │
       │ url             │ │ stream_id (FK)  │ │ user_id (FK)    │
       │ order           │ │ start_price     │ │ joined_at       │
       └─────────────────┘ │ reserve_price   │ │ left_at         │
                           │ buy_now_price   │ └─────────────────┘
                           │ start_time      │
                           │ end_time        │
                           │ status          │
                           └────────┬────────┘
                                    │
                                    ▼
                           ┌─────────────────┐
                           │      bids       │
                           ├─────────────────┤
                           │ id (PK)         │
                           │ auction_id (FK) │
                           │ bidder_id (FK)  │
                           │ amount          │
                           │ is_auto_bid     │
                           │ max_auto_amount │
                           │ created_at      │
                           └─────────────────┘
                                    │
                                    ▼
                           ┌─────────────────┐
                           │     orders      │
                           ├─────────────────┤
                           │ id (PK)         │
                           │ buyer_id (FK)   │
                           │ seller_id (FK)  │
                           │ auction_id (FK) │
                           │ total_amount    │
                           │ status          │
                           │ created_at      │
                           └────────┬────────┘
                                    │
              ┌─────────────────────┼─────────────────────┐
              │                     │                     │
              ▼                     ▼                     ▼
       ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
       │   order_items   │ │    payments     │ │   shipments     │
       ├─────────────────┤ ├─────────────────┤ ├─────────────────┤
       │ id (PK)         │ │ id (PK)         │ │ id (PK)         │
       │ order_id (FK)   │ │ order_id (FK)   │ │ order_id (FK)   │
       │ product_id (FK) │ │ method          │ │ tracking_num    │
       │ quantity        │ │ amount          │ │ carrier         │
       │ unit_price      │ │ status          │ │ status          │
       └─────────────────┘ │ stripe_intent_id│ │ shipped_at      │
                           │ created_at      │ │ delivered_at    │
                           └─────────────────┘ └─────────────────┘

┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  categories     │     │ chat_messages   │     │  notifications  │
├─────────────────┤     ├─────────────────┤     ├─────────────────┤
│ id (PK)         │     │ id (PK)         │     │ id (PK)         │
│ name            │     │ stream_id (FK)  │     │ user_id (FK)    │
│ slug            │     │ user_id (FK)    │     │ type            │
│ parent_id (FK)  │     │ content         │     │ title           │
│ image_url       │     │ message_type    │     │ message         │
│ sort_order      │     │ created_at      │     │ is_read         │
└─────────────────┘     └─────────────────┘     │ created_at      │
                                                └─────────────────┘
```

## Table Definitions

### users
Core user accounts table.

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    avatar_url VARCHAR(500),
    role VARCHAR(20) DEFAULT 'buyer' CHECK (role IN ('buyer', 'seller', 'admin')),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'suspended', 'deleted')),
    email_verified_at TIMESTAMP,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_status ON users(status);
```

### sellers
Extended profile for sellers.

```sql
CREATE TABLE sellers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    store_name VARCHAR(100) NOT NULL,
    store_slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    logo_url VARCHAR(500),
    banner_url VARCHAR(500),
    rating DECIMAL(2,1) DEFAULT 5.0 CHECK (rating >= 0 AND rating <= 5),
    review_count INTEGER DEFAULT 0,
    verified BOOLEAN DEFAULT FALSE,
    verified_at TIMESTAMP,
    commission_rate DECIMAL(4,2) DEFAULT 5.00, -- percentage
    total_sales INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sellers_user_id ON sellers(user_id);
CREATE INDEX idx_sellers_slug ON sellers(store_slug);
CREATE INDEX idx_sellers_verified ON sellers(verified);
```

### buyers
Extended profile for buyers with shipping preferences.

```sql
CREATE TABLE buyers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    default_shipping_address_id UUID,
    preferences JSONB DEFAULT '{}',
    total_purchases INTEGER DEFAULT 0,
    total_spent DECIMAL(12,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_buyers_user_id ON buyers(user_id);
```

### addresses
Shipping addresses for users.

```sql
CREATE TABLE addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    label VARCHAR(50) NOT NULL, -- 'Home', 'Office', etc.
    recipient_name VARCHAR(200) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    address_line1 VARCHAR(255) NOT NULL,
    address_line2 VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL, -- Malaysian states
    postal_code VARCHAR(10) NOT NULL,
    country VARCHAR(100) DEFAULT 'Malaysia',
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_addresses_user_id ON addresses(user_id);
CREATE INDEX idx_addresses_default ON addresses(user_id, is_default);
```

### categories
Product categories (hierarchical).

```sql
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    image_url VARCHAR(500),
    parent_id UUID REFERENCES categories(id),
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_categories_slug ON categories(slug);
CREATE INDEX idx_categories_parent ON categories(parent_id);
CREATE INDEX idx_categories_active ON categories(is_active);
```

### products
Product catalog.

```sql
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seller_id UUID NOT NULL REFERENCES sellers(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    condition VARCHAR(20) DEFAULT 'new' CHECK (condition IN ('new', 'used', 'refurbished')),
    base_price DECIMAL(12,2) NOT NULL,
    compare_at_price DECIMAL(12,2), -- original price for sales
    stock_quantity INTEGER DEFAULT 0,
    sku VARCHAR(100),
    weight_grams INTEGER, -- for shipping calculations
    dimensions_cm JSONB, -- {"length": 10, "width": 5, "height": 3}
    attributes JSONB DEFAULT '{}', -- {"color": "red", "size": "L"}
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'archived')),
    view_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_products_seller ON products(seller_id);
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_products_status ON products(status);
CREATE INDEX idx_products_slug ON products(slug);
CREATE INDEX idx_products_price ON products(base_price);
```

### product_images
Product images (multiple per product).

```sql
CREATE TABLE product_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    url VARCHAR(500) NOT NULL,
    thumbnail_url VARCHAR(500),
    alt_text VARCHAR(255),
    sort_order INTEGER DEFAULT 0,
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_product_images_product ON product_images(product_id);
CREATE INDEX idx_product_images_primary ON product_images(product_id, is_primary);
```

### streams
Live stream sessions.

```sql
CREATE TABLE streams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seller_id UUID NOT NULL REFERENCES sellers(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    livekit_room_name VARCHAR(100) UNIQUE,
    status VARCHAR(20) DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'live', 'ended', 'cancelled')),
    scheduled_at TIMESTAMP,
    started_at TIMESTAMP,
    ended_at TIMESTAMP,
    viewer_count INTEGER DEFAULT 0,
    max_viewers INTEGER DEFAULT 0,
    recording_url VARCHAR(500),
    thumbnail_url VARCHAR(500),
    is_recorded BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_streams_seller ON streams(seller_id);
CREATE INDEX idx_streams_status ON streams(status);
CREATE INDEX idx_streams_scheduled ON streams(scheduled_at);
```

### stream_products
Products featured in a stream (many-to-many).

```sql
CREATE TABLE stream_products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stream_id UUID NOT NULL REFERENCES streams(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    featured_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    auction_started_at TIMESTAMP,
    auction_ended_at TIMESTAMP,
    UNIQUE(stream_id, product_id)
);

CREATE INDEX idx_stream_products_stream ON stream_products(stream_id);
CREATE INDEX idx_stream_products_product ON stream_products(product_id);
```

### auctions
Auction sessions for products.

```sql
CREATE TABLE auctions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stream_id UUID REFERENCES streams(id), -- optional: standalone auction
    product_id UUID NOT NULL REFERENCES products(id),
    seller_id UUID NOT NULL REFERENCES sellers(id),
    start_price DECIMAL(12,2) NOT NULL,
    reserve_price DECIMAL(12,2), -- minimum price to sell
    buy_now_price DECIMAL(12,2), -- instant purchase price
    min_bid_increment DECIMAL(12,2) DEFAULT 1.00,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'active', 'ended', 'cancelled')),
    winner_id UUID REFERENCES users(id),
    winning_bid_id UUID, -- set when auction ends
    total_bids INTEGER DEFAULT 0,
    view_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_auctions_stream ON auctions(stream_id);
CREATE INDEX idx_auctions_product ON auctions(product_id);
CREATE INDEX idx_auctions_seller ON auctions(seller_id);
CREATE INDEX idx_auctions_status ON auctions(status);
CREATE INDEX idx_auctions_time ON auctions(start_time, end_time);
```

### bids
Bid history.

```sql
CREATE TABLE bids (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    auction_id UUID NOT NULL REFERENCES auctions(id),
    bidder_id UUID NOT NULL REFERENCES users(id),
    amount DECIMAL(12,2) NOT NULL,
    is_auto_bid BOOLEAN DEFAULT FALSE,
    max_auto_amount DECIMAL(12,2), -- for proxy bidding
    is_winning BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bids_auction ON bids(auction_id);
CREATE INDEX idx_bids_bidder ON bids(bidder_id);
CREATE INDEX idx_bids_amount ON bids(auction_id, amount DESC);
CREATE INDEX idx_bids_created ON bids(created_at);
```

### orders
Customer orders.

```sql
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number VARCHAR(50) UNIQUE NOT NULL, -- human-readable
    buyer_id UUID NOT NULL REFERENCES buyers(id),
    seller_id UUID NOT NULL REFERENCES sellers(id),
    auction_id UUID REFERENCES auctions(id), -- if from auction
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'processing', 'shipped', 'delivered', 'cancelled', 'refunded')),
    subtotal DECIMAL(12,2) NOT NULL,
    shipping_cost DECIMAL(12,2) DEFAULT 0,
    tax_amount DECIMAL(12,2) DEFAULT 0,
    discount_amount DECIMAL(12,2) DEFAULT 0,
    total_amount DECIMAL(12,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'MYR',
    shipping_address JSONB NOT NULL, -- snapshot of address
    notes TEXT,
    paid_at TIMESTAMP,
    confirmed_at TIMESTAMP,
    shipped_at TIMESTAMP,
    delivered_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_orders_buyer ON orders(buyer_id);
CREATE INDEX idx_orders_seller ON orders(seller_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_auction ON orders(auction_id);
CREATE INDEX idx_orders_created ON orders(created_at);
```

### order_items
Line items in an order.

```sql
CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id),
    product_name VARCHAR(255) NOT NULL, -- snapshot
    product_image VARCHAR(500), -- snapshot
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(12,2) NOT NULL,
    total_price DECIMAL(12,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_order_items_order ON order_items(order_id);
```

### payments
Payment records (Stripe integration).

```sql
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id),
    user_id UUID NOT NULL REFERENCES users(id),
    amount DECIMAL(12,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'MYR',
    method VARCHAR(50) NOT NULL, -- 'card', 'fpx', 'grabpay', etc.
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'refunded')),
    stripe_payment_intent_id VARCHAR(100),
    stripe_charge_id VARCHAR(100),
    failure_reason TEXT,
    metadata JSONB DEFAULT '{}',
    processed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_payments_order ON payments(order_id);
CREATE INDEX idx_payments_user ON payments(user_id);
CREATE INDEX idx_payments_stripe ON payments(stripe_payment_intent_id);
CREATE INDEX idx_payments_status ON payments(status);
```

### shipments
Shipping records (NinjaVan integration).

```sql
CREATE TABLE shipments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id),
    carrier VARCHAR(50) DEFAULT 'NinjaVan',
    tracking_number VARCHAR(100),
    ninjavan_order_id VARCHAR(100), -- NinjaVan's order ID
    label_url VARCHAR(500), -- shipping label PDF
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'label_created', 'picked_up', 'in_transit', 'out_for_delivery', 'delivered', 'failed')),
    shipping_cost DECIMAL(12,2),
    weight_grams INTEGER,
    dimensions_cm JSONB,
    estimated_delivery TIMESTAMP,
    shipped_at TIMESTAMP,
    delivered_at TIMESTAMP,
    tracking_history JSONB DEFAULT '[]',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_shipments_order ON shipments(order_id);
CREATE INDEX idx_shipments_tracking ON shipments(tracking_number);
CREATE INDEX idx_shipments_status ON shipments(status);
```

### chat_messages
Stream chat messages.

```sql
CREATE TABLE chat_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stream_id UUID NOT NULL REFERENCES streams(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    message_type VARCHAR(20) DEFAULT 'text' CHECK (message_type IN ('text', 'system', 'notification')),
    reply_to_id UUID REFERENCES chat_messages(id),
    is_deleted BOOLEAN DEFAULT FALSE,
    deleted_by UUID REFERENCES users(id),
    deleted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_chat_messages_stream ON chat_messages(stream_id);
CREATE INDEX idx_chat_messages_user ON chat_messages(user_id);
CREATE INDEX idx_chat_messages_created ON chat_messages(created_at);
```

### followers
User follows for sellers.

```sql
CREATE TABLE followers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    seller_id UUID NOT NULL REFERENCES sellers(id) ON DELETE CASCADE,
    notifications_enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, seller_id)
);

CREATE INDEX idx_followers_user ON followers(user_id);
CREATE INDEX idx_followers_seller ON followers(seller_id);
```

### notifications
User notifications.

```sql
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL, -- 'bid_outbid', 'auction_won', 'order_shipped', etc.
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    data JSONB DEFAULT '{}', -- additional context
    action_url VARCHAR(500),
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notifications_user ON notifications(user_id);
CREATE INDEX idx_notifications_read ON notifications(user_id, is_read);
CREATE INDEX idx_notifications_type ON notifications(type);
```

### reviews
Seller reviews from buyers.

```sql
CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id),
    reviewer_id UUID NOT NULL REFERENCES users(id),
    seller_id UUID NOT NULL REFERENCES sellers(id),
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    title VARCHAR(255),
    content TEXT,
    images JSONB DEFAULT '[]',
    is_verified_purchase BOOLEAN DEFAULT TRUE,
    helpful_count INTEGER DEFAULT 0,
    seller_response TEXT,
    seller_responded_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_reviews_seller ON reviews(seller_id);
CREATE INDEX idx_reviews_reviewer ON reviews(reviewer_id);
CREATE INDEX idx_reviews_rating ON reviews(seller_id, rating);
```

## Indexes Summary

### Performance Critical Indexes
```sql
-- Fast user lookup
CREATE INDEX CONCURRENTLY idx_users_email ON users(email);

-- Fast auction queries
CREATE INDEX CONCURRENTLY idx_auctions_status_time ON auctions(status, start_time, end_time);

-- Fast bid history
CREATE INDEX CONCURRENTLY idx_bids_auction_amount ON bids(auction_id, amount DESC);

-- Fast order lookup
CREATE INDEX CONCURRENTLY idx_orders_buyer_status ON orders(buyer_id, status);

-- Fast product search
CREATE INDEX CONCURRENTLY idx_products_status_category ON products(status, category_id);
```

## Triggers

### Auto-update timestamps
```sql
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply to all tables with updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_sellers_updated_at BEFORE UPDATE ON sellers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Repeat for other tables...
```

### Update seller rating on new review
```sql
CREATE OR REPLACE FUNCTION update_seller_rating()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE sellers
    SET 
        rating = (SELECT AVG(rating)::DECIMAL(2,1) FROM reviews WHERE seller_id = NEW.seller_id),
        review_count = (SELECT COUNT(*) FROM reviews WHERE seller_id = NEW.seller_id)
    WHERE id = NEW.seller_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_seller_rating_trigger
    AFTER INSERT OR UPDATE ON reviews
    FOR EACH ROW EXECUTE FUNCTION update_seller_rating();
```

---

*Last updated: 2025-02-05*
