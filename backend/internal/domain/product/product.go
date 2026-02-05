package product

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Errors
var (
	ErrProductNotFound      = errors.New("product not found")
	ErrInvalidPrice         = errors.New("price must be greater than zero")
	ErrInvalidStock         = errors.New("stock quantity cannot be negative")
	ErrProductAlreadySold   = errors.New("product already sold")
	ErrUnauthorized         = errors.New("unauthorized action")
	ErrProductHasActiveAuction = errors.New("product has active auction")
)

// Condition represents the condition of a product
type Condition string

const (
	ConditionNew        Condition = "new"
	ConditionUsed       Condition = "used"
	ConditionRefurbished Condition = "refurbished"
)

// Status represents the listing status of a product
type Status string

const (
	StatusDraft    Status = "draft"
	StatusActive   Status = "active"
	StatusSold     Status = "sold"
	StatusArchived Status = "archived"
)

// Product represents a product listing
type Product struct {
	ID              uuid.UUID
	SellerID        uuid.UUID
	CategoryID      *uuid.UUID
	Name            string
	Slug            string
	Description     string
	Condition       Condition
	BasePrice       float64
	CompareAtPrice  *float64 // Original price for sales
	StockQuantity   int
	SKU             *string
	WeightGrams     *int
	DimensionsCm    *Dimensions
	Attributes      map[string]string // e.g., {"color": "red", "size": "L"}
	Images          []ProductImage
	Status          Status
	ViewCount       int
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}

// Dimensions represents product dimensions
type Dimensions struct {
	Length int `json:"length"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// ProductImage represents a product image
type ProductImage struct {
	ID            uuid.UUID
	ProductID     uuid.UUID
	URL           string
	ThumbnailURL  *string
	AltText       *string
	SortOrder     int
	IsPrimary     bool
	CreatedAt     time.Time
}

// IsAvailable returns true if the product is available for purchase
func (p *Product) IsAvailable() bool {
	return p.Status == StatusActive && p.StockQuantity > 0
}

// CanEdit returns true if the product can be edited by the seller
func (p *Product) CanEdit(sellerID uuid.UUID) bool {
	if p.SellerID != sellerID {
		return false
	}
	// Can't edit sold products
	if p.Status == StatusSold {
		return false
	}
	return true
}

// DecrementStock reduces stock quantity
func (p *Product) DecrementStock(quantity int) error {
	if p.StockQuantity < quantity {
		return errors.New("insufficient stock")
	}
	p.StockQuantity -= quantity
	
	// Auto-mark as sold if stock reaches 0
	if p.StockQuantity == 0 {
		p.Status = StatusSold
	}
	
	return nil
}

// IncrementStock increases stock quantity
func (p *Product) IncrementStock(quantity int) {
	p.StockQuantity += quantity
	// If was sold and now has stock, make active again
	if p.Status == StatusSold && p.StockQuantity > 0 {
		p.Status = StatusActive
	}
}

// Validate performs basic validation on the product
func (p *Product) Validate() error {
	if p.Name == "" {
		return errors.New("product name is required")
	}
	
	if p.BasePrice <= 0 {
		return ErrInvalidPrice
	}
	
	if p.StockQuantity < 0 {
		return ErrInvalidStock
	}
	
	if p.Condition != ConditionNew && p.Condition != ConditionUsed && p.Condition != ConditionRefurbished {
		return errors.New("invalid condition")
	}
	
	return nil
}

// GetPrimaryImage returns the primary image or first image
func (p *Product) GetPrimaryImage() *ProductImage {
	for _, img := range p.Images {
		if img.IsPrimary {
			return &img
		}
	}
	if len(p.Images) > 0 {
		return &p.Images[0]
	}
	return nil
}

// Repository defines the interface for product data access
type Repository interface {
	// GetByID retrieves a product by ID
	GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
	
	// GetBySlug retrieves a product by its slug
	GetBySlug(ctx context.Context, slug string) (*Product, error)
	
	// List retrieves a list of products with filtering
	List(ctx context.Context, filter Filter) ([]*Product, int64, error)
	
	// Create creates a new product
	Create(ctx context.Context, product *Product) error
	
	// Update updates an existing product
	Update(ctx context.Context, product *Product) error
	
	// Delete soft-deletes a product
	Delete(ctx context.Context, id uuid.UUID) error
	
	// IncrementViewCount increments the view count
	IncrementViewCount(ctx context.Context, id uuid.UUID) error
	
	// UpdateStock updates the stock quantity
	UpdateStock(ctx context.Context, id uuid.UUID, quantity int) error
}

// Filter represents filter criteria for listing products
type Filter struct {
	SellerID       *uuid.UUID
	CategoryID     *uuid.UUID
	Status         *Status
	Condition      *Condition
	MinPrice       *float64
	MaxPrice       *float64
	Query          string // Search query
	SortBy         string // e.g., "newest", "price_asc", "price_desc", "popular"
	Page           int
	PageSize       int
}

// NewProduct creates a new product with default values
func NewProduct(sellerID uuid.UUID, name, description string, condition Condition, basePrice float64, stockQty int) *Product {
	now := time.Now()
	return &Product{
		ID:            uuid.New(),
		SellerID:      sellerID,
		Name:          name,
		Slug:          generateSlug(name),
		Description:   description,
		Condition:     condition,
		BasePrice:     basePrice,
		StockQuantity: stockQty,
		Attributes:    make(map[string]string),
		Images:        make([]ProductImage, 0),
		Status:        StatusDraft,
		ViewCount:     0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// generateSlug creates a URL-friendly slug from the name
func generateSlug(name string) string {
	// Simple slug generation - in production, use a proper library
	// This is a placeholder implementation
	slug := ""
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			slug += string(r)
		} else if r == ' ' || r == '-' {
			slug += "-"
		}
	}
	// Add unique suffix
	slug += "-" + uuid.New().String()[:8]
	return slug
}

// AddImage adds an image to the product
func (p *Product) AddImage(url string, isPrimary bool) ProductImage {
	img := ProductImage{
		ID:        uuid.New(),
		ProductID: p.ID,
		URL:       url,
		SortOrder: len(p.Images),
		IsPrimary: isPrimary,
		CreatedAt: time.Now(),
	}
	
	// If this is primary, unset other primaries
	if isPrimary {
		for i := range p.Images {
			p.Images[i].IsPrimary = false
		}
	}
	
	p.Images = append(p.Images, img)
	return img
}

// SetPrimaryImage sets an image as the primary image
func (p *Product) SetPrimaryImage(imageID uuid.UUID) error {
	found := false
	for i := range p.Images {
		if p.Images[i].ID == imageID {
			p.Images[i].IsPrimary = true
			found = true
		} else {
			p.Images[i].IsPrimary = false
		}
	}
	
	if !found {
		return errors.New("image not found")
	}
	
	return nil
}

// Update updates the product fields
func (p *Product) Update(name, description string, basePrice float64, stockQty int) error {
	if name != "" {
		p.Name = name
	}
	if description != "" {
		p.Description = description
	}
	if basePrice > 0 {
		p.BasePrice = basePrice
	}
	if stockQty >= 0 {
		p.StockQuantity = stockQty
	}
	
	p.UpdatedAt = time.Now()
	return p.Validate()
}

// Publish changes status from draft to active
func (p *Product) Publish() error {
	if p.Status != StatusDraft && p.Status != StatusArchived {
		return errors.New("can only publish draft or archived products")
	}
	
	if len(p.Images) == 0 {
		return errors.New("product must have at least one image")
	}
	
	if err := p.Validate(); err != nil {
		return err
	}
	
	p.Status = StatusActive
	p.UpdatedAt = time.Now()
	return nil
}

// Archive changes status to archived
func (p *Product) Archive() error {
	if p.Status == StatusSold {
		return errors.New("cannot archive sold product")
	}
	p.Status = StatusArchived
	p.UpdatedAt = time.Now()
	return nil
}
