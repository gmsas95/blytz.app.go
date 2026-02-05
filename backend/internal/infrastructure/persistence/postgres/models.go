package postgres

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type JSONMap map[string]interface{}

func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, j)
	case string:
		return json.Unmarshal([]byte(v), j)
	default:
		return json.Unmarshal([]byte(fmt.Sprintf("%v", v)), j)
	}
}

func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

type StringArray []string

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, a)
	case string:
		return json.Unmarshal([]byte(v), a)
	default:
		return json.Unmarshal([]byte(fmt.Sprintf("%v", v)), a)
	}
}

func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

type User struct {
	BaseModel
	Email         string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash  string    `gorm:"not null" json:"-"`
	Role          string    `gorm:"not null;default:'buyer'" json:"role"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	AvatarURL     string    `json:"avatar_url"`
	Phone         string    `json:"phone"`
	EmailVerified bool      `gorm:"default:false" json:"email_verified"`
	LastLoginAt   *time.Time `json:"last_login_at"`
}

type Category struct {
	BaseModel
	Name        string     `gorm:"not null" json:"name"`
	Slug        string     `gorm:"uniqueIndex;not null" json:"slug"`
	Description string     `json:"description"`
	ImageURL    string     `json:"image_url"`
	ParentID    *uuid.UUID `gorm:"index" json:"parent_id"`
	SortOrder   int        `gorm:"default:0" json:"sort_order"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
}

type Product struct {
	BaseModel
	SellerID       uuid.UUID `gorm:"not null;index" json:"seller_id"`
	CategoryID     uuid.UUID `gorm:"not null;index" json:"category_id"`
	Title          string    `gorm:"not null" json:"title"`
	Description    string    `json:"description"`
	Condition      string    `json:"condition"`
	StartingPrice  float64   `gorm:"not null" json:"starting_price"`
	ReservePrice   *float64  `json:"reserve_price"`
	BuyNowPrice    *float64  `json:"buy_now_price"`
	Images         StringArray `gorm:"type:jsonb" json:"images"`
	VideoURL       string    `json:"video_url"`
	Specifications JSONMap   `gorm:"type:jsonb" json:"specifications"`
	ShippingInfo   JSONMap   `gorm:"type:jsonb" json:"shipping_info"`
	Status         string    `gorm:"default:'draft'" json:"status"`
	Featured       bool      `gorm:"default:false" json:"featured"`
	ViewCount      int       `gorm:"default:0" json:"view_count"`
	Rating         float64   `gorm:"default:0" json:"rating"`
	ReviewCount    int       `gorm:"default:0" json:"review_count"`
	IsFlash        bool      `gorm:"default:false" json:"is_flash"`
	IsHot          bool      `gorm:"default:false" json:"is_hot"`
	FlashEnd       *time.Time `json:"flash_end"`
}

type Auction struct {
	BaseModel
	ProductID    uuid.UUID  `gorm:"not null;index" json:"product_id"`
	SellerID     uuid.UUID  `gorm:"not null;index" json:"seller_id"`
	Title        string     `gorm:"not null" json:"title"`
	Description  string     `json:"description"`
	StartTime    time.Time  `gorm:"not null" json:"start_time"`
	EndTime      time.Time  `gorm:"not null" json:"end_time"`
	Status       string     `gorm:"default:'scheduled'" json:"status"`
	StartPrice   float64    `gorm:"not null" json:"start_price"`
	ReservePrice *float64   `json:"reserve_price"`
	BuyNowPrice  *float64   `json:"buy_now_price"`
	CurrentBidID *uuid.UUID `gorm:"index" json:"-"`
	BidCount     int        `gorm:"default:0" json:"bid_count"`
	WinnerID     *uuid.UUID `gorm:"index" json:"winner_id"`
	LiveKitRoom  string     `gorm:"not null;uniqueIndex" json:"livekit_room"`
	StreamKey    string     `json:"stream_key"`
	AutoExtend   bool       `gorm:"default:true" json:"auto_extend"`
	ExtendTime   int        `gorm:"default:300" json:"extend_time"`
	IsFeatured   bool       `gorm:"default:false" json:"is_featured"`
}

type Bid struct {
	BaseModel
	AuctionID uuid.UUID `gorm:"not null;index:idx_bid_auction_user" json:"auction_id"`
	UserID    uuid.UUID `gorm:"not null;index:idx_bid_auction_user" json:"user_id"`
	Amount    float64   `gorm:"not null" json:"amount"`
	IsAutoBid bool      `gorm:"default:false" json:"is_auto_bid"`
	IsWinning bool      `gorm:"default:false" json:"is_winning"`
	BidTime   time.Time `gorm:"not null" json:"bid_time"`
}

type AutoBid struct {
	BaseModel
	AuctionID    uuid.UUID  `gorm:"not null;index" json:"auction_id"`
	UserID       uuid.UUID  `gorm:"not null;index" json:"user_id"`
	MaxAmount    float64    `gorm:"not null" json:"max_amount"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	CurrentBid   *float64   `json:"current_bid"`
	BidIncrement float64    `gorm:"not null;default:5.0" json:"bid_increment"`
	LastBidTime  *time.Time `json:"last_bid_time"`
}

type Order struct {
	BaseModel
	UserID          uuid.UUID `gorm:"not null;index" json:"user_id"`
	AuctionID       *uuid.UUID `gorm:"index" json:"auction_id"`
	Status          string    `gorm:"default:'pending'" json:"status"`
	TotalAmount     float64   `gorm:"not null" json:"total_amount"`
	Subtotal        float64   `gorm:"not null" json:"subtotal"`
	TaxAmount       float64   `gorm:"default:0" json:"tax_amount"`
	ShippingCost    float64   `gorm:"default:0" json:"shipping_cost"`
	DiscountAmount  float64   `gorm:"default:0" json:"discount_amount"`
	ShippingAddress JSONMap   `gorm:"type:jsonb" json:"shipping_address"`
	BillingAddress  JSONMap   `gorm:"type:jsonb" json:"billing_address"`
	PaymentID       *uuid.UUID `gorm:"index" json:"payment_id"`
	TrackingNumber  string    `json:"tracking_number"`
	Notes           string    `json:"notes"`
}

type OrderItem struct {
	BaseModel
	OrderID   uuid.UUID `gorm:"not null;index" json:"order_id"`
	ProductID uuid.UUID `gorm:"not null" json:"product_id"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	UnitPrice float64   `gorm:"not null" json:"unit_price"`
	Total     float64   `gorm:"not null" json:"total"`
}

type Cart struct {
	BaseModel
	UserID    *uuid.UUID `gorm:"index" json:"user_id"`
	Token     string     `gorm:"uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
}

type CartItem struct {
	BaseModel
	CartID    uuid.UUID `gorm:"not null;index" json:"cart_id"`
	ProductID uuid.UUID `gorm:"not null" json:"product_id"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	AddedAt   time.Time `gorm:"autoCreateTime" json:"added_at"`
}

type Payment struct {
	BaseModel
	OrderID           uuid.UUID  `gorm:"not null;index" json:"order_id"`
	UserID            uuid.UUID  `gorm:"not null;index" json:"user_id"`
	Amount            float64    `gorm:"not null" json:"amount"`
	Currency          string     `gorm:"default:'USD'" json:"currency"`
	Status            string     `gorm:"default:'pending'" json:"status"`
	Method            string     `json:"method"`
	TransactionID     string     `gorm:"uniqueIndex" json:"transaction_id"`
	GatewayReference  string     `json:"gateway_reference"`
	FailureReason     string     `json:"failure_reason"`
	RefundedAmount    float64    `gorm:"default:0" json:"refunded_amount"`
	RefundedAt        *time.Time `json:"refunded_at"`
	Metadata          JSONMap    `gorm:"type:jsonb" json:"metadata"`
}

type PaymentMethod struct {
	BaseModel
	UserID     uuid.UUID  `gorm:"not null;index" json:"user_id"`
	Type       string     `gorm:"not null" json:"type"`
	Provider   string     `gorm:"not null" json:"provider"`
	MethodRef  string     `gorm:"not null" json:"method_ref"`
	IsDefault  bool       `gorm:"default:false" json:"is_default"`
	Last4      string     `json:"last4"`
	ExpiryDate *time.Time `json:"expiry_date"`
}

func AutoMigrate(db *gorm.DB) error {
	// Migrate legacy models
	if err := db.AutoMigrate(
		&User{},
		&Auction{},
		&Bid{},
		&AutoBid{},
		&Order{},
		&OrderItem{},
		&Cart{},
		&CartItem{},
		&Payment{},
		&PaymentMethod{},
	); err != nil {
		return err
	}
	
	// Migrate new product/category models
	if err := AutoMigrateProduct(db); err != nil {
		return err
	}
	
	if err := AutoMigrateCategory(db); err != nil {
		return err
	}
	
	// Seed default categories
	if err := SeedCategories(db); err != nil {
		return err
	}
	
	return nil
}