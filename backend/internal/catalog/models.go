package catalog

import (
	"encoding/json"
	"time"

	"github.com/blytz.live.remake/backend/internal/common"
	"github.com/google/uuid"
)

// CategoryCreateRequest represents category creation request
type CategoryCreateRequest struct {
	Name        string     `json:"name" binding:"required,min=2,max=100"`
	Description *string    `json:"description,omitempty"`
	ImageURL    *string    `json:"image_url,omitempty"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	SortOrder   *int       `json:"sort_order,omitempty"`
	IsActive    *bool      `json:"is_active,omitempty"`
}

// CategoryUpdateRequest represents category update request
type CategoryUpdateRequest struct {
	Name        *string    `json:"name,omitempty" binding:"omitempty,min=2,max=100"`
	Description *string    `json:"description,omitempty"`
	ImageURL    *string    `json:"image_url,omitempty"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	SortOrder   *int       `json:"sort_order,omitempty"`
	IsActive    *bool      `json:"is_active,omitempty"`
}

// CategoryResponse represents category response with hierarchy
type CategoryResponse struct {
	ID            uuid.UUID            `json:"id"`
	Name          string               `json:"name"`
	Slug          string               `json:"slug"`
	Description   *string              `json:"description"`
	ImageURL      *string              `json:"image_url"`
	ParentID      *uuid.UUID           `json:"parent_id"`
	SortOrder     int                  `json:"sort_order"`
	IsActive      bool                 `json:"is_active"`
	ProductCount  int                  `json:"product_count"`
	Level         int                  `json:"level"`
	Path          []string             `json:"path"`
	Children      []CategoryResponse   `json:"children,omitempty"`
	Parent        *CategoryResponse    `json:"parent,omitempty"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
}

// CategoryTreeRequest represents category tree query
type CategoryTreeRequest struct {
	ParentID   *uuid.UUID `form:"parent_id"`
	IsActive   *bool      `form:"is_active"`
	MaxDepth   *int       `form:"max_depth"`
	IncludeProductCount bool `form:"include_product_count"`
	SortBy     string     `form:"sort_by" binding:"omitempty,oneof=name sort_order created_at"`
	SortDirection string  `form:"sort_direction" binding:"omitempty,oneof=asc desc"`
}

// CategoryMoveRequest represents moving a category to a new parent
type CategoryMoveRequest struct {
	NewParentID *uuid.UUID `json:"new_parent_id,omitempty"`
	NewSortOrder *int       `json:"new_sort_order,omitempty"`
}

// BulkCategoryOperation represents bulk category operations
type BulkCategoryOperation struct {
	CategoryIDs []uuid.UUID `json:"category_ids" binding:"required"`
	Operation   string      `json:"operation" binding:"required,oneof=activate deactivate delete"`
	NewParentID *uuid.UUID  `json:"new_parent_id,omitempty"`
}

// CategoryAttribute represents custom attributes for categories
type CategoryAttribute struct {
	common.BaseModel
	CategoryID  uuid.UUID `gorm:"not null;references:ID" json:"category_id"`
	Name        string    `gorm:"not null" json:"name"`
	Type        string    `gorm:"not null" json:"type"` // text, number, boolean, select, multiselect
	Required    bool      `gorm:"default:false" json:"required"`
	Options     []string  `gorm:"type:jsonb" json:"options,omitempty"` // For select/multiselect
	DefaultValue *string   `json:"default_value,omitempty"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
}

// CategoryAttributeRequest represents category attribute request
type CategoryAttributeRequest struct {
	Name         string     `json:"name" binding:"required,min=2,max=50"`
	Type         string     `json:"type" binding:"required,oneof=text number boolean select multiselect"`
	Required     *bool      `json:"required,omitempty"`
	Options      []string   `json:"options,omitempty"`
	DefaultValue *string    `json:"default_value,omitempty"`
	SortOrder    *int       `json:"sort_order,omitempty"`
}

// ProductVariant represents product variants (size, color, etc.)
type ProductVariant struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`
	ProductID   uuid.UUID `gorm:"not null;references:ID" json:"product_id"`
	Sku         string    `gorm:"uniqueIndex;not null" json:"sku"`
	Title       string    `gorm:"not null" json:"title"`
	Price       float64   `gorm:"not null" json:"price"`
	ComparePrice *float64 `json:"compare_price"`
	CostPrice   *float64  `json:"cost_price"`
	Weight      *float64  `json:"weight,omitempty"`
	Barcode     *string   `json:"barcode,omitempty"`
	Inventory   int       `gorm:"default:0" json:"inventory"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	Attributes  string    `gorm:"type:jsonb" json:"attributes"` // JSON: {"color": "Red", "size": "M"}
	ImageURL    *string   `json:"image_url,omitempty"`
	Position    int       `gorm:"default:0" json:"position"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProductVariantRequest represents product variant request
type ProductVariantRequest struct {
	Title        string             `json:"title" binding:"required,min=2,max=100"`
	Sku          string             `json:"sku" binding:"required,min=2,max=50"`
	Price        float64            `json:"price" binding:"required,gt=0"`
	ComparePrice *float64           `json:"compare_price,omitempty" binding:"omitempty,gt=0"`
	CostPrice    *float64           `json:"cost_price,omitempty" binding:"omitempty,gt=0"`
	Weight       *float64           `json:"weight,omitempty" binding:"omitempty,gt=0"`
	Barcode      *string            `json:"barcode,omitempty"`
	Inventory    *int               `json:"inventory,omitempty" binding:"omitempty,gte=0"`
	IsActive     *bool              `json:"is_active,omitempty"`
	Attributes   map[string]string  `json:"attributes"`
	ImageURL     *string            `json:"image_url,omitempty"`
	Position     *int               `json:"position,omitempty" binding:"omitempty,gte=0"`
}

// ProductVariantResponse represents product variant response
type ProductVariantResponse struct {
	ID           uuid.UUID          `json:"id"`
	ProductID    uuid.UUID          `json:"product_id"`
	Sku          string             `json:"sku"`
	Title        string             `json:"title"`
	Price        float64            `json:"price"`
	ComparePrice *float64           `json:"compare_price"`
	CostPrice    *float64           `json:"cost_price"`
	Weight       *float64           `json:"weight"`
	Barcode      *string            `json:"barcode"`
	Inventory    int                `json:"inventory"`
	IsActive     bool               `json:"is_active"`
	Attributes   map[string]string  `json:"attributes"`
	ImageURL     *string            `json:"image_url"`
	Position     int                `json:"position"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

// ProductCollection represents product collections/groupings
type ProductCollection struct {
	ID          uuid.UUID  `gorm:"primaryKey;type:uuid" json:"id"`
	Name        string     `gorm:"not null" json:"name"`
	Slug        string     `gorm:"uniqueIndex;not null" json:"slug"`
	Description *string    `json:"description"`
	ImageURL    *string    `json:"image_url"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	SortOrder   int        `gorm:"default:0" json:"sort_order"`
	ProductIDs  string     `gorm:"type:jsonb" json:"product_ids"` // Array of product UUIDs
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ProductCollectionRequest represents product collection request
type ProductCollectionRequest struct {
	Name        string      `json:"name" binding:"required,min=2,max=100"`
	Description *string     `json:"description,omitempty"`
	ImageURL    *string     `json:"image_url,omitempty"`
	IsActive    *bool       `json:"is_active,omitempty"`
	SortOrder   *int        `json:"sort_order,omitempty"`
	ProductIDs  []uuid.UUID `json:"product_ids,omitempty"`
}

// ProductCollectionResponse represents product collection response
type ProductCollectionResponse struct {
	ID          uuid.UUID    `json:"id"`
	Name        string       `json:"name"`
	Slug        string       `json:"slug"`
	Description *string      `json:"description"`
	ImageURL    *string      `json:"image_url"`
	IsActive    bool         `json:"is_active"`
	SortOrder   int          `json:"sort_order"`
	ProductIDs  []uuid.UUID  `json:"product_ids"`
	ProductCount int         `json:"product_count"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// InventoryStock represents inventory tracking
type InventoryStock struct {
	ID              uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`
	ProductID       uuid.UUID `gorm:"not null;uniqueIndex;references:ID" json:"product_id"`
	VariantID       *uuid.UUID `gorm:"references:ID" json:"variant_id,omitempty"`
	Quantity        int        `gorm:"not null;default:0" json:"quantity"`
	Reserved        int        `gorm:"not null;default:0" json:"reserved"`
	Available       int        `gorm:"not null;default:0" json:"available"`
	LowStockAlert   int        `gorm:"default:10" json:"low_stock_alert"`
	TrackInventory bool       `gorm:"default:true" json:"track_inventory"`
	AllowBackorder  bool       `gorm:"default:false" json:"allow_backorder"`
	WarehouseID     *uuid.UUID `json:"warehouse_id,omitempty"`
	LastUpdated     time.Time  `gorm:"autoUpdateTime" json:"last_updated"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// InventoryStockRequest represents inventory stock request
type InventoryStockRequest struct {
	Quantity       int        `json:"quantity" binding:"required,gte=0"`
	LowStockAlert  *int       `json:"low_stock_alert,omitempty" binding:"omitempty,gte=0"`
	TrackInventory *bool      `json:"track_inventory,omitempty"`
	AllowBackorder *bool      `json:"allow_backorder,omitempty"`
	WarehouseID    *uuid.UUID `json:"warehouse_id,omitempty"`
}

// InventoryStockResponse represents inventory stock response
type InventoryStockResponse struct {
	ID              uuid.UUID  `json:"id"`
	ProductID       uuid.UUID  `json:"product_id"`
	VariantID       *uuid.UUID `json:"variant_id,omitempty"`
	Quantity        int        `json:"quantity"`
	Reserved        int        `json:"reserved"`
	Available       int        `json:"available"`
	LowStockAlert   int        `json:"low_stock_alert"`
	TrackInventory bool       `json:"track_inventory"`
	AllowBackorder  bool       `json:"allow_backorder"`
	WarehouseID     *uuid.UUID `json:"warehouse_id,omitempty"`
	LastUpdated     time.Time  `json:"last_updated"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// StockMovement represents stock movement history
type StockMovement struct {
	ID          uuid.UUID  `gorm:"primaryKey;type:uuid" json:"id"`
	ProductID   uuid.UUID  `gorm:"not null;references:ID" json:"product_id"`
	VariantID   *uuid.UUID `gorm:"references:ID" json:"variant_id,omitempty"`
	MovementType string     `gorm:"not null" json:"movement_type"` // in, out, adjustment, reserve, release
	Quantity    int        `gorm:"not null" json:"quantity"`
	Reference   *string    `json:"reference,omitempty"`
	Notes       *string    `json:"notes,omitempty"`
	WarehouseID *uuid.UUID `json:"warehouse_id,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// StockMovementRequest represents stock movement request
type StockMovementRequest struct {
	MovementType string     `json:"movement_type" binding:"required,oneof=in out adjustment reserve release"`
	Quantity     int        `json:"quantity" binding:"required,gt=0"`
	Reference    *string    `json:"reference,omitempty"`
	Notes        *string    `json:"notes,omitempty"`
	WarehouseID  *uuid.UUID `json:"warehouse_id,omitempty"`
}

// CatalogStats represents catalog statistics
type CatalogStats struct {
	TotalCategories    int                    `json:"total_categories"`
	TotalProducts      int                    `json:"total_products"`
	TotalVariants      int                    `json:"total_variants"`
	ActiveProducts     int                    `json:"active_products"`
	LowStockProducts   int                    `json:"low_stock_products"`
	OutofStockProducts int                    `json:"out_of_stock_products"`
	ProductsByCategory map[string]int        `json:"products_by_category"`
	TopCategories      []CategoryStats        `json:"top_categories"`
	RecentProducts     []ProductSummary       `json:"recent_products"`
}

// CategoryStats represents category statistics
type CategoryStats struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	ProductCount  int       `json:"product_count"`
	ActiveCount   int       `json:"active_count"`
	TotalRevenue  float64   `json:"total_revenue"`
}

// ProductSummary represents product summary for stats
type ProductSummary struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

// BulkOperationResponse represents bulk operation result
type BulkOperationResponse struct {
	TotalProcessed int      `json:"total_processed"`
	SuccessCount   int      `json:"success_count"`
	FailureCount   int      `json:"failure_count"`
	Errors         []string `json:"errors,omitempty"`
}

// Helper function to get inventory availability
func (i *InventoryStock) UpdateAvailability() {
	i.Available = i.Quantity - i.Reserved
	if i.Available < 0 {
		i.Available = 0
	}
}

// Helper function to check if product is low stock
func (i *InventoryStock) IsLowStock() bool {
	return i.TrackInventory && i.Available <= i.LowStockAlert
}

// Helper function to check if product is out of stock
func (i *InventoryStock) IsOutOfStock() bool {
	return i.TrackInventory && i.Available == 0 && !i.AllowBackorder
}

// Helper function to check if product is available
func (i *InventoryStock) IsAvailable(quantity int) bool {
	if !i.TrackInventory {
		return true
	}
	return i.AllowBackorder || i.Available >= quantity
}

// Helper functions for variant attributes
func (pv *ProductVariant) GetAttribute(key string) (string, bool) {
	var attrs map[string]string
	if err := json.Unmarshal([]byte(pv.Attributes), &attrs); err != nil {
		return "", false
	}
	val, exists := attrs[key]
	return val, exists
}

func (pv *ProductVariant) SetAttribute(key, value string) {
	var attrs map[string]string
	json.Unmarshal([]byte(pv.Attributes), &attrs)
	if attrs == nil {
		attrs = make(map[string]string)
	}
	attrs[key] = value
	data, _ := json.Marshal(attrs)
	pv.Attributes = string(data)
}

// Helper function to get collection product count
func (pc *ProductCollection) GetProductCount() int {
	var ids []uuid.UUID
	if err := json.Unmarshal([]byte(pc.ProductIDs), &ids); err != nil {
		return 0
	}
	return len(ids)
}