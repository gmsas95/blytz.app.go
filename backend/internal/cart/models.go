package cart

import (
	"time"
	
	"github.com/google/uuid"
)

// Cart represents a shopping cart
type Cart struct {
	ID        uuid.UUID  `json:"id"`
	UserID    *uuid.UUID `json:"user_id,omitempty"` // nullable for guest carts
	Token      string     `json:"token"`      // for guest carts
	ExpiresAt  time.Time  `json:"expires_at"`
	Items      []CartItem `json:"items"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// CartItem represents items in a cart
type CartItem struct {
	ID        uuid.UUID  `json:"id"`
	CartID    uuid.UUID  `json:"cart_id"`
	ProductID uuid.UUID  `json:"product_id"`
	Quantity  int        `json:"quantity"`
	AddedAt   time.Time  `json:"added_at"`
}

// CartCreateRequest represents cart creation request
type CartCreateRequest struct {
	UserID *uuid.UUID `json:"user_id,omitempty"`
}

// CartResponse represents cart response with product details
type CartResponse struct {
	ID             uuid.UUID              `json:"id"`
	UserID         *uuid.UUID             `json:"user_id,omitempty"`
	Token          string                 `json:"token"`
	ExpiresAt      time.Time              `json:"expires_at"`
	Items          []CartItemResponse     `json:"items"`
	ItemCount      int                    `json:"item_count"`
	Subtotal       float64                `json:"subtotal"`
	TotalItems     int                    `json:"total_items"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// CartItemResponse represents cart item with product details
type CartItemResponse struct {
	ID        uuid.UUID                `json:"id"`
	ProductID uuid.UUID                `json:"product_id"`
	Quantity  int                      `json:"quantity"`
	AddedAt   time.Time                `json:"added_at"`
	Product   ProductResponse          `json:"product"`
	LineTotal float64                  `json:"line_total"`
}

// ProductResponse represents product information in cart context
type ProductResponse struct {
	ID            uuid.UUID  `json:"id"`
	Title         string     `json:"title"`
	Description   *string    `json:"description"`
	Condition     *string    `json:"condition"`
	StartingPrice float64    `json:"starting_price"`
	ReservePrice  *float64   `json:"reserve_price"`
	BuyNowPrice   *float64   `json:"buy_now_price"`
	Images        []string   `json:"images"`
	Status        string     `json:"status"`
}

// AddItemRequest represents add item to cart request
type AddItemRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,min=1,max=10"`
}

// UpdateItemRequest represents update cart item request
type UpdateItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1,max=10"`
}

// MergeCartRequest represents merge guest cart to user cart request
type MergeCartRequest struct {
	GuestToken string `json:"guest_token" binding:"required"`
}