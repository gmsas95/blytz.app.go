package models

import (
	"time"

	"github.com/blytz.live.remake/backend/internal/common"
	"github.com/google/uuid"
)

// User represents a user in system
type User struct {
	common.BaseModel
	Email         string  `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash  string  `gorm:"not null" json:"-"`
	Role          string  `gorm:"not null;default:'buyer'" json:"role"` // 'buyer', 'seller', 'admin'
	FirstName     *string `json:"first_name"`
	LastName      *string `json:"last_name"`
	AvatarURL     *string `json:"avatar_url"`
	Phone         *string `json:"phone"`
	EmailVerified bool      `gorm:"default:false" json:"email_verified"`
	LastLoginAt   *time.Time `json:"last_login_at"`
}

// Category represents a product category
type Category struct {
	common.BaseModel
	Name        string  `gorm:"not null" json:"name"`
	Slug        string  `gorm:"uniqueIndex;not null" json:"slug"`
	Description *string `json:"description"`
	ImageURL    *string `json:"image_url"`
	ParentID    *uuid.UUID `gorm:"references:ID" json:"parent_id"`
	Parent      *Category `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	SortOrder   int     `gorm:"default:0" json:"sort_order"`
	IsActive    bool    `gorm:"default:true" json:"is_active"`
	Categories  []Category `gorm:"foreignKey:ParentID" json:"categories,omitempty"`
}

// Product represents a product in marketplace
type Product struct {
	common.BaseModel
	SellerID      uuid.UUID `gorm:"not null;references:ID" json:"seller_id"`
	Seller        User      `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	CategoryID    uuid.UUID `gorm:"not null;references:ID" json:"category_id"`
	Category      Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Title         string    `gorm:"not null" json:"title"`
	Description   *string   `json:"description"`
	Condition     *string   `json:"condition"` // 'new', 'like_new', 'good', 'fair'
	StartingPrice float64   `gorm:"not null" json:"starting_price"`
	ReservePrice  *float64  `json:"reserve_price"`
	BuyNowPrice   *float64  `json:"buy_now_price"`
	Images        *string   `gorm:"type:jsonb" json:"images"`       // JSON array of image URLs
	VideoURL      *string   `json:"video_url"`
	Specifications *string  `gorm:"type:jsonb" json:"specifications"` // JSON object
	ShippingInfo  *string   `gorm:"type:jsonb" json:"shipping_info"`   // JSON object
	Status        string    `gorm:"default:'draft'" json:"status"`     // 'draft', 'active', 'sold', 'cancelled'
	Featured      bool      `gorm:"default:false" json:"featured"`
	ViewCount     int       `gorm:"default:0" json:"view_count"`
}

// Order represents a customer order
type Order struct {
	common.BaseModel
	UserID          uuid.UUID  `gorm:"not null;references:ID" json:"user_id"`
	Status          string     `gorm:"not null;default:'pending'" json:"status"` // pending, processing, shipped, delivered, cancelled
	TotalAmount     float64    `gorm:"not null" json:"total_amount"`
	Subtotal        float64    `gorm:"not null" json:"subtotal"`
	TaxAmount       float64    `gorm:"default:0" json:"tax_amount"`
	ShippingCost    float64    `gorm:"default:0" json:"shipping_cost"`
	DiscountAmount  float64    `gorm:"default:0" json:"discount_amount"`
	ShippingAddress *Address   `gorm:"embedded;embeddedPrefix:shipping_" json:"shipping_address"`
	BillingAddress  *Address   `gorm:"embedded;embeddedPrefix:billing_" json:"billing_address"`
	PaymentID       *uuid.UUID `gorm:"references:ID" json:"payment_id"`
	TrackingNumber  *string    `json:"tracking_number"`
	Notes           *string    `json:"notes"`
	Items          []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

// OrderItem represents items in an order
type OrderItem struct {
	common.BaseModel
	OrderID     uuid.UUID `gorm:"not null;references:ID" json:"order_id"`
	ProductID   uuid.UUID `gorm:"not null;references:ID" json:"product_id"`
	Quantity    int        `gorm:"not null" json:"quantity"`
	UnitPrice   float64    `gorm:"not null" json:"unit_price"`
	Total       float64    `gorm:"not null" json:"total"`
	Product     Product    `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// Address represents shipping/billing address (embedded)
type Address struct {
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Company      *string `json:"company"`
	AddressLine1 string  `json:"address_line1"`
	AddressLine2 *string `json:"address_line2"`
	City         string  `json:"city"`
	State        string  `json:"state"`
	PostalCode   string  `json:"postal_code"`
	Country      string  `json:"country"`
	Phone        *string `json:"phone"`
}

// Cart represents a shopping cart
type Cart struct {
	common.BaseModel
	UserID    *uuid.UUID `gorm:"references:ID" json:"user_id,omitempty"`
	Token      string     `gorm:"uniqueIndex" json:"token"`
	ExpiresAt  time.Time  `gorm:"not null" json:"expires_at"`
	Items      []CartItem `gorm:"foreignKey:CartID" json:"items,omitempty"`
}

// CartItem represents items in a cart
type CartItem struct {
	common.BaseModel
	CartID    uuid.UUID `gorm:"not null;references:ID" json:"cart_id"`
	ProductID uuid.UUID `gorm:"not null;references:ID" json:"product_id"`
	Quantity   int        `gorm:"not null" json:"quantity"`
	AddedAt    time.Time  `gorm:"autoCreateTime" json:"added_at"`
	Product    Product    `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}