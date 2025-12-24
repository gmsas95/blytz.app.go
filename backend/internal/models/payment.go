package models

import (
	"time"

	"github.com/blytz.live.remake/backend/internal/common"
	"github.com/google/uuid"
)

// Payment represents a payment transaction
type Payment struct {
	common.BaseModel
	OrderID        uuid.UUID  `gorm:"not null;references:ID" json:"order_id"`
	Order          Order      `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	UserID         uuid.UUID  `gorm:"not null;references:ID" json:"user_id"`
	User           User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	PaymentMethod  string     `gorm:"not null" json:"payment_method"` // stripe, paypal, credit_card
	Amount         float64    `gorm:"not null" json:"amount"`
	Currency       string     `gorm:"not null;default:'USD'" json:"currency"`
	Status         string     `gorm:"default:'pending'" json:"status"` // pending, processing, completed, failed, refunded, cancelled
	TransactionID  string     `gorm:"uniqueIndex" json:"transaction_id"`
	GatewayRef     string     `json:"gateway_ref"` // stripe_payment_id, paypal_id, etc.
	GatewayType    string     `gorm:"not null" json:"gateway_type"` // stripe, paypal, apple_pay, google_pay
	FailureReason  *string    `json:"failure_reason"`
	RefundedAmount float64    `gorm:"default:0" json:"refunded_amount"`
	RefundReason   *string    `json:"refund_reason"`
	ProcessedAt    *time.Time `json:"processed_at"`
	ExpiresAt      *time.Time `json:"expires_at"`
	Metadata       string     `gorm:"type:jsonb" json:"metadata"` // JSON object for additional data
}

// PaymentMethod represents a user's saved payment method
type PaymentMethod struct {
	common.BaseModel
	UserID      uuid.UUID  `gorm:"not null;references:ID" json:"user_id"`
	User        User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Type        string     `gorm:"not null" json:"type"` // credit_card, debit_card, paypal, bank_account
	Provider    string     `gorm:"not null" json:"provider"` // stripe, paypal, etc.
	MethodRef   string     `gorm:"not null" json:"method_ref"` // tokenized reference
	IsDefault   bool       `gorm:"default:false" json:"is_default"`
	Last4       *string    `json:"last4"`
	Brand       *string    `json:"brand"` // visa, mastercard, etc.
	ExpiryMonth *int       `json:"expiry_month"`
	ExpiryYear  *int       `json:"expiry_year"`
	Name        *string    `json:"name"`
	Email       *string    `json:"email"`
	Phone       *string    `json:"phone"`
	Address     *Address   `gorm:"embedded;embeddedPrefix:address_" json:"address"`
	IsVerified  bool       `gorm:"default:false" json:"is_verified"`
	Metadata    string     `gorm:"type:jsonb" json:"metadata"`
}

// PaymentIntent represents a payment intent for processing
type PaymentIntent struct {
	common.BaseModel
	UserID       uuid.UUID  `gorm:"not null;references:ID" json:"user_id"`
	User         User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Amount       float64    `gorm:"not null" json:"amount"`
	Currency     string     `gorm:"not null;default:'USD'" json:"currency"`
	Status       string     `gorm:"default:'requires_payment_method'" json:"status"` // requires_payment_method, requires_confirmation, requires_action, processing, succeeded, canceled
	ClientSecret string     `gorm:"not null;json:"client_secret"`
	GatewayRef   string     `gorm:"not null;json:"gateway_ref"` // stripe_payment_intent_id
	GatewayType  string     `gorm:"not null;json:"gateway_type"` // stripe, paypal, etc.
	PaymentMethods []string `gorm:"type:jsonb" json:"payment_methods"` // card, ideal, sepa_debit, etc.
	ExpiresAt    time.Time  `gorm:"not null" json:"expires_at"`
	OrderID      *uuid.UUID `gorm:"references:ID" json:"order_id"`
	Order        *Order     `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Metadata     string     `gorm:"type:jsonb" json:"metadata"`
}

// Refund represents a refund transaction
type Refund struct {
	common.BaseModel
	PaymentID     uuid.UUID `gorm:"not null;references:ID" json:"payment_id"`
	Payment       Payment   `gorm:"foreignKey:PaymentID" json:"payment,omitempty"`
	Amount        float64   `gorm:"not null" json:"amount"`
	Reason        string    `gorm:"not null" json:"reason"` // duplicate, fraudulent, requested_by_customer
	Status        string    `gorm:"default:'pending'" json:"status"` // pending, succeeded, failed, cancelled
	GatewayRef    string    `json:"gateway_ref"` // stripe_refund_id
	GatewayType   string    `gorm:"not null" json:"gateway_type"`
	ProcessedBy   uuid.UUID `gorm:"not null;references:ID" json:"processed_by"`
	Processor     User      `gorm:"foreignKey:ProcessedBy" json:"processor,omitempty"`
	Notes         *string   `json:"notes"`
	ProcessedAt   *time.Time `json:"processed_at"`
	Metadata      string    `gorm:"type:jsonb" json:"metadata"`
}

// Transaction represents a generic financial transaction
type Transaction struct {
	common.BaseModel
	Type         string     `gorm:"not null" json:"type"` // payment, refund, payout, fee
	Status       string     `gorm:"not null;default:'pending'" json:"status"`
	Amount       float64    `gorm:"not null" json:"amount"`
	Currency     string     `gorm:"not null;default:'USD'" json:"currency"`
	UserID       uuid.UUID  `gorm:"not null;references:ID" json:"user_id"`
	User         User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	OrderID      *uuid.UUID `gorm:"references:ID" json:"order_id"`
	Order        *Order     `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	PaymentID    *uuid.UUID `gorm:"references:ID" json:"payment_id"`
	Payment      *Payment   `gorm:"foreignKey:PaymentID" json:"payment,omitempty"`
	GatewayRef   string     `json:"gateway_ref"`
	GatewayType  string     `gorm:"not null" json:"gateway_type"`
	Description  string     `gorm:"not null" json:"description"`
	Metadata     string     `gorm:"type:jsonb" json:"metadata"`
	Balance      float64    `gorm:"not null;default:0" json:"balance"` // Running balance
	ProcessedAt  *time.Time `json:"processed_at"`
}

// Payout represents a payout to a seller
type Payout struct {
	common.BaseModel
	SellerID     uuid.UUID `gorm:"not null;references:ID" json:"seller_id"`
	Seller       User      `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	Amount       float64   `gorm:"not null" json:"amount"`
	Currency     string    `gorm:"not null;default:'USD'" json:"currency"`
	Status       string    `gorm:"not null;default:'pending'" json:"status"` // pending, in_transit, paid, failed, cancelled
	GatewayRef   string    `json:"gateway_ref"`
	GatewayType  string    `gorm:"not null" json:"gateway_type"`
	Method       string    `gorm:"not null" json:"method"` // bank_transfer, stripe_connect, paypal
	Destination  string    `gorm:"not null" json:"destination"` // Bank account or email
	Fee          float64   `gorm:"default:0" json:"fee"`
	NetAmount    float64   `gorm:"not null" json:"net_amount"`
	OrderIDs     string    `gorm:"type:jsonb" json:"order_ids"` // Array of order UUIDs being paid out
	ProcessedBy  uuid.UUID `gorm:"not null;references:ID" json:"processed_by"`
	Processor    User      `gorm:"foreignKey:ProcessedBy" json:"processor,omitempty"`
	Notes        *string   `json:"notes"`
	ScheduledAt  *time.Time `json:"scheduled_at"`
	ProcessedAt  *time.Time `json:"processed_at"`
	Metadata     string    `gorm:"type:jsonb" json:"metadata"`
}

// Subscription represents a user subscription
type Subscription struct {
	common.BaseModel
	UserID          uuid.UUID  `gorm:"not null;references:ID" json:"user_id"`
	User            User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	PlanID          string     `gorm:"not null" json:"plan_id"`
	PlanName        string     `gorm:"not null" json:"plan_name"`
	Status          string     `gorm:"not null;default:'active'" json:"status"` // active, past_due, canceled, unpaid
	CurrentPeriodStart time.Time `gorm:"not null" json:"current_period_start"`
	CurrentPeriodEnd   time.Time `gorm:"not null" json:"current_period_end"`
	CancelAtPeriodEnd bool       `gorm:"default:false" json:"cancel_at_period_end"`
	GatewayRef      string     `gorm:"not null" json:"gateway_ref"`
	GatewayType     string     `gorm:"not null" json:"gateway_type"`
	TrialStart      *time.Time `json:"trial_start"`
	TrialEnd        *time.Time `json:"trial_end"`
	CanceledAt      *time.Time `json:"canceled_at"`
	EndedAt         *time.Time `json:"ended_at"`
	Amount          float64    `gorm:"not null" json:"amount"`
	Currency        string     `gorm:"not null;default:'USD'" json:"currency"`
	Interval        string     `gorm:"not null" json:"interval"` // month, year
	IntervalCount   int        `gorm:"not null;default:1" json:"interval_count"`
	Metadata        string     `gorm:"type:jsonb" json:"metadata"`
}