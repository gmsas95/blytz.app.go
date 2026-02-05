package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/blytz/live/backend/internal/domain/product"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ProductModel represents the product database model
type ProductModel struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SellerID       uuid.UUID      `gorm:"type:uuid;not null;index"`
	CategoryID     *uuid.UUID     `gorm:"type:uuid;index"`
	Name           string         `gorm:"not null"`
	Slug           string         `gorm:"uniqueIndex;not null"`
	Description    string         `gorm:"type:text"`
	Condition      string         `gorm:"not null;default:'new'"`
	BasePrice      float64        `gorm:"type:decimal(12,2);not null"`
	CompareAtPrice *float64       `gorm:"type:decimal(12,2)"`
	StockQuantity  int            `gorm:"not null;default:0"`
	SKU            *string
	WeightGrams    *int
	DimensionsCm   json.RawMessage `gorm:"type:jsonb"`
	Attributes     json.RawMessage `gorm:"type:jsonb;default:'{}'"`
	Status         string         `gorm:"not null;default:'draft'"`
	ViewCount      int            `gorm:"not null;default:0"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	
	// Associations
	Images []ProductImageModel `gorm:"foreignKey:ProductID;order:sort_order ASC"`
}

func (ProductModel) TableName() string {
	return "products"
}

// ProductImageModel represents the product image database model
type ProductImageModel struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID    uuid.UUID `gorm:"type:uuid;not null;index"`
	URL          string    `gorm:"not null"`
	ThumbnailURL *string
	AltText      *string
	SortOrder    int       `gorm:"not null;default:0"`
	IsPrimary    bool      `gorm:"not null;default:false"`
	CreatedAt    time.Time
}

func (ProductImageModel) TableName() string {
	return "product_images"
}

// ProductRepository implements product.Repository
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// GetByID retrieves a product by ID
func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*product.Product, error) {
	var model ProductModel
	
	err := r.db.WithContext(ctx).
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		First(&model, "id = ?", id).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, product.ErrProductNotFound
		}
		return nil, err
	}
	
	return r.toDomain(&model), nil
}

// GetBySlug retrieves a product by its slug
func (r *ProductRepository) GetBySlug(ctx context.Context, slug string) (*product.Product, error) {
	var model ProductModel
	
	err := r.db.WithContext(ctx).
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		First(&model, "slug = ?", slug).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, product.ErrProductNotFound
		}
		return nil, err
	}
	
	return r.toDomain(&model), nil
}

// List retrieves products with filtering
func (r *ProductRepository) List(ctx context.Context, filter product.Filter) ([]*product.Product, int64, error) {
	var models []ProductModel
	var total int64
	
	query := r.db.WithContext(ctx).Model(&ProductModel{}).Preload("Images")
	
	// Apply filters
	if filter.SellerID != nil {
		query = query.Where("seller_id = ?", *filter.SellerID)
	}
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Condition != nil {
		query = query.Where("condition = ?", *filter.Condition)
	}
	if filter.MinPrice != nil {
		query = query.Where("base_price >= ?", *filter.MinPrice)
	}
	if filter.MaxPrice != nil {
		query = query.Where("base_price <= ?", *filter.MaxPrice)
	}
	if filter.Query != "" {
		query = query.Where(
			"name ILIKE ? OR description ILIKE ?",
			fmt.Sprintf("%%%s%%", filter.Query),
			fmt.Sprintf("%%%s%%", filter.Query),
		)
	}
	
	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Apply sorting
	switch filter.SortBy {
	case "price_asc":
		query = query.Order("base_price ASC")
	case "price_desc":
		query = query.Order("base_price DESC")
	case "popular":
		query = query.Order("view_count DESC")
	case "newest":
		query = query.Order("created_at DESC")
	default:
		query = query.Order("created_at DESC")
	}
	
	// Apply pagination
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)
	
	if err := query.Find(&models).Error; err != nil {
		return nil, 0, err
	}
	
	// Convert to domain
	products := make([]*product.Product, len(models))
	for i, model := range models {
		products[i] = r.toDomain(&model)
	}
	
	return products, total, nil
}

// Create creates a new product
func (r *ProductRepository) Create(ctx context.Context, p *product.Product) error {
	model := r.toModel(p)
	
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&model).Error; err != nil {
			return err
		}
		
		// Create images
		for i := range p.Images {
			imgModel := r.toImageModel(&p.Images[i])
			if err := tx.Create(&imgModel).Error; err != nil {
				return err
			}
		}
		
		return nil
	})
}

// Update updates an existing product
func (r *ProductRepository) Update(ctx context.Context, p *product.Product) error {
	model := r.toModel(p)
	
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update product
		if err := tx.Model(&ProductModel{}).Where("id = ?", p.ID).Updates(map[string]interface{}{
			"category_id":      model.CategoryID,
			"name":             model.Name,
			"slug":             model.Slug,
			"description":      model.Description,
			"condition":        model.Condition,
			"base_price":       model.BasePrice,
			"compare_at_price": model.CompareAtPrice,
			"stock_quantity":   model.StockQuantity,
			"sku":              model.SKU,
			"weight_grams":     model.WeightGrams,
			"dimensions_cm":    model.DimensionsCm,
			"attributes":       model.Attributes,
			"status":           model.Status,
			"updated_at":       time.Now(),
		}).Error; err != nil {
			return err
		}
		
		// Update images - delete old, create new
		if err := tx.Where("product_id = ?", p.ID).Delete(&ProductImageModel{}).Error; err != nil {
			return err
		}
		
		for i := range p.Images {
			imgModel := r.toImageModel(&p.Images[i])
			if err := tx.Create(&imgModel).Error; err != nil {
				return err
			}
		}
		
		return nil
	})
}

// Delete soft-deletes a product
func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&ProductModel{}, "id = ?", id).Error
}

// IncrementViewCount increments the view count
func (r *ProductRepository) IncrementViewCount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&ProductModel{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

// UpdateStock updates the stock quantity
func (r *ProductRepository) UpdateStock(ctx context.Context, id uuid.UUID, quantity int) error {
	return r.db.WithContext(ctx).Model(&ProductModel{}).
		Where("id = ?", id).
		Update("stock_quantity", quantity).Error
}

// toDomain converts a database model to domain entity
func (r *ProductRepository) toDomain(model *ProductModel) *product.Product {
	p := &product.Product{
		ID:            model.ID,
		SellerID:      model.SellerID,
		CategoryID:    model.CategoryID,
		Name:          model.Name,
		Slug:          model.Slug,
		Description:   model.Description,
		Condition:     product.Condition(model.Condition),
		BasePrice:     model.BasePrice,
		CompareAtPrice: model.CompareAtPrice,
		StockQuantity: model.StockQuantity,
		SKU:           model.SKU,
		WeightGrams:   model.WeightGrams,
		Status:        product.Status(model.Status),
		ViewCount:     model.ViewCount,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}
	
	if model.DeletedAt.Valid {
		p.DeletedAt = &model.DeletedAt.Time
	}
	
	// Parse dimensions
	if len(model.DimensionsCm) > 0 && string(model.DimensionsCm) != "null" {
		var dims product.Dimensions
		json.Unmarshal(model.DimensionsCm, &dims)
		p.DimensionsCm = &dims
	}
	
	// Parse attributes
	if len(model.Attributes) > 0 && string(model.Attributes) != "null" {
		json.Unmarshal(model.Attributes, &p.Attributes)
	} else {
		p.Attributes = make(map[string]string)
	}
	
	// Convert images
	p.Images = make([]product.ProductImage, len(model.Images))
	for i, imgModel := range model.Images {
		p.Images[i] = product.ProductImage{
			ID:           imgModel.ID,
			ProductID:    imgModel.ProductID,
			URL:          imgModel.URL,
			ThumbnailURL: imgModel.ThumbnailURL,
			AltText:      imgModel.AltText,
			SortOrder:    imgModel.SortOrder,
			IsPrimary:    imgModel.IsPrimary,
			CreatedAt:    imgModel.CreatedAt,
		}
	}
	
	return p
}

// toModel converts a domain entity to database model
func (r *ProductRepository) toModel(p *product.Product) ProductModel {
	model := ProductModel{
		ID:             p.ID,
		SellerID:       p.SellerID,
		CategoryID:     p.CategoryID,
		Name:           p.Name,
		Slug:           p.Slug,
		Description:    p.Description,
		Condition:      string(p.Condition),
		BasePrice:      p.BasePrice,
		CompareAtPrice: p.CompareAtPrice,
		StockQuantity:  p.StockQuantity,
		SKU:            p.SKU,
		WeightGrams:    p.WeightGrams,
		Status:         string(p.Status),
		ViewCount:      p.ViewCount,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
	
	// Serialize dimensions
	if p.DimensionsCm != nil {
		dimsJSON, _ := json.Marshal(p.DimensionsCm)
		model.DimensionsCm = dimsJSON
	}
	
	// Serialize attributes
	if p.Attributes != nil {
		attrsJSON, _ := json.Marshal(p.Attributes)
		model.Attributes = attrsJSON
	}
	
	return model
}

// toImageModel converts a domain image to database model
func (r *ProductRepository) toImageModel(img *product.ProductImage) ProductImageModel {
	return ProductImageModel{
		ID:           img.ID,
		ProductID:    img.ProductID,
		URL:          img.URL,
		ThumbnailURL: img.ThumbnailURL,
		AltText:      img.AltText,
		SortOrder:    img.SortOrder,
		IsPrimary:    img.IsPrimary,
		CreatedAt:    img.CreatedAt,
	}
}

// Ensure ProductRepository implements product.Repository
var _ product.Repository = (*ProductRepository)(nil)

// AutoMigrateProduct creates product tables
func AutoMigrateProduct(db *gorm.DB) error {
	return db.AutoMigrate(&ProductModel{}, &ProductImageModel{})
}

// CreateProductIndexes creates indexes for product search
func CreateProductIndexes(db *gorm.DB) error {
	// Create full-text search index for PostgreSQL
	sql := `
	CREATE INDEX IF NOT EXISTS idx_products_search ON products 
	USING gin(to_tsvector('english', name || ' ' || COALESCE(description, '')));
	`
	return db.Exec(sql).Error
}

// SearchProducts performs full-text search
func (r *ProductRepository) SearchProducts(ctx context.Context, query string, page, pageSize int) ([]*product.Product, int64, error) {
	var models []ProductModel
	var total int64
	
	// Use PostgreSQL full-text search
	searchQuery := strings.Join(strings.Fields(query), " & ")
	
	sql := `
		SELECT *, ts_rank(to_tsvector('english', name || ' ' || COALESCE(description, '')), 
		plainto_tsquery('english', ?)) as rank
		FROM products
		WHERE to_tsvector('english', name || ' ' || COALESCE(description, '')) @@ plainto_tsquery('english', ?)
		AND status = 'active'
		AND deleted_at IS NULL
		ORDER BY rank DESC, created_at DESC
		LIMIT ? OFFSET ?
	`
	
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Raw(sql, query, query, pageSize, offset).Scan(&models).Error; err != nil {
		return nil, 0, err
	}
	
	// Get total
	countSql := `
		SELECT COUNT(*) FROM products
		WHERE to_tsvector('english', name || ' ' || COALESCE(description, '')) @@ plainto_tsquery('english', ?)
		AND status = 'active'
		AND deleted_at IS NULL
	`
	if err := r.db.WithContext(ctx).Raw(countSql, query).Scan(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Convert to domain
	products := make([]*product.Product, len(models))
	for i, model := range models {
		products[i] = r.toDomain(&model)
	}
	
	return products, total, nil
}
