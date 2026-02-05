package postgres

import (
	"context"
	"time"

	"github.com/blytz/live/backend/internal/domain/category"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CategoryModel represents the category database model
type CategoryModel struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name         string         `gorm:"not null"`
	Slug         string         `gorm:"uniqueIndex;not null"`
	Description  *string
	ImageURL     *string
	ParentID     *uuid.UUID     `gorm:"type:uuid;index"`
	SortOrder    int            `gorm:"not null;default:0"`
	IsActive     bool           `gorm:"not null;default:true"`
	ProductCount int            `gorm:"not null;default:0"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (CategoryModel) TableName() string {
	return "categories"
}

// CategoryRepository implements category.Repository
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// GetByID retrieves a category by ID
func (r *CategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*category.Category, error) {
	var model CategoryModel
	
	err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, category.ErrCategoryNotFound
		}
		return nil, err
	}
	
	return r.toDomain(&model), nil
}

// GetBySlug retrieves a category by its slug
func (r *CategoryRepository) GetBySlug(ctx context.Context, slug string) (*category.Category, error) {
	var model CategoryModel
	
	err := r.db.WithContext(ctx).First(&model, "slug = ?", slug).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, category.ErrCategoryNotFound
		}
		return nil, err
	}
	
	return r.toDomain(&model), nil
}

// GetTree retrieves the full category tree
func (r *CategoryRepository) GetTree(ctx context.Context) ([]*category.Category, error) {
	var models []CategoryModel
	
	// Get all categories ordered by parent_id and sort_order
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Where("deleted_at IS NULL").
		Order("COALESCE(parent_id::text, ''), sort_order ASC").
		Find(&models).Error
	
	if err != nil {
		return nil, err
	}
	
	// Convert to domain
	cats := make([]*category.Category, len(models))
	for i, model := range models {
		cats[i] = r.toDomain(&model)
	}
	
	// Build tree
	tree := category.BuildTree(cats)
	
	return tree, nil
}

// GetChildren retrieves immediate children of a category
func (r *CategoryRepository) GetChildren(ctx context.Context, parentID *uuid.UUID) ([]*category.Category, error) {
	var models []CategoryModel
	
	query := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Where("deleted_at IS NULL").
		Order("sort_order ASC")
	
	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}
	
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	
	cats := make([]*category.Category, len(models))
	for i, model := range models {
		cats[i] = r.toDomain(&model)
	}
	
	return cats, nil
}

// List retrieves all categories (flat list)
func (r *CategoryRepository) List(ctx context.Context, onlyActive bool) ([]*category.Category, error) {
	var models []CategoryModel
	
	query := r.db.WithContext(ctx).Where("deleted_at IS NULL")
	
	if onlyActive {
		query = query.Where("is_active = ?", true)
	}
	
	if err := query.Order("sort_order ASC, name ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	
	cats := make([]*category.Category, len(models))
	for i, model := range models {
		cats[i] = r.toDomain(&model)
	}
	
	return cats, nil
}

// Create creates a new category
func (r *CategoryRepository) Create(ctx context.Context, cat *category.Category) error {
	model := r.toModel(cat)
	return r.db.WithContext(ctx).Create(&model).Error
}

// Update updates an existing category
func (r *CategoryRepository) Update(ctx context.Context, cat *category.Category) error {
	model := r.toModel(cat)
	return r.db.WithContext(ctx).Model(&CategoryModel{}).Where("id = ?", cat.ID).Updates(map[string]interface{}{
		"name":          model.Name,
		"slug":          model.Slug,
		"description":   model.Description,
		"image_url":     model.ImageURL,
		"parent_id":     model.ParentID,
		"sort_order":    model.SortOrder,
		"is_active":     model.IsActive,
		"updated_at":    time.Now(),
	}).Error
}

// Delete soft-deletes a category
func (r *CategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&CategoryModel{}, "id = ?", id).Error
}

// UpdateProductCount updates the cached product count
func (r *CategoryRepository) UpdateProductCount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Exec(`
		UPDATE categories 
		SET product_count = (
			SELECT COUNT(*) FROM products 
			WHERE category_id = ? AND status = 'active' AND deleted_at IS NULL
		),
		updated_at = NOW()
		WHERE id = ?
	`, id, id).Error
}

// toDomain converts a database model to domain entity
func (r *CategoryRepository) toDomain(model *CategoryModel) *category.Category {
	return &category.Category{
		ID:           model.ID,
		Name:         model.Name,
		Slug:         model.Slug,
		Description:  model.Description,
		ImageURL:     model.ImageURL,
		ParentID:     model.ParentID,
		SortOrder:    model.SortOrder,
		IsActive:     model.IsActive,
		ProductCount: model.ProductCount,
		Children:     make([]*category.Category, 0),
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}
}

// toModel converts a domain entity to database model
func (r *CategoryRepository) toModel(cat *category.Category) CategoryModel {
	return CategoryModel{
		ID:           cat.ID,
		Name:         cat.Name,
		Slug:         cat.Slug,
		Description:  cat.Description,
		ImageURL:     cat.ImageURL,
		ParentID:     cat.ParentID,
		SortOrder:    cat.SortOrder,
		IsActive:     cat.IsActive,
		ProductCount: cat.ProductCount,
		CreatedAt:    cat.CreatedAt,
		UpdatedAt:    cat.UpdatedAt,
	}
}

// Ensure CategoryRepository implements category.Repository
var _ category.Repository = (*CategoryRepository)(nil)

// AutoMigrateCategory creates category tables
func AutoMigrateCategory(db *gorm.DB) error {
	return db.AutoMigrate(&CategoryModel{})
}

// SeedCategories inserts default categories
func SeedCategories(db *gorm.DB) error {
	categories := []CategoryModel{
		{Name: "Fashion", Slug: "fashion", IsActive: true, SortOrder: 1},
		{Name: "Electronics", Slug: "electronics", IsActive: true, SortOrder: 2},
		{Name: "Collectibles", Slug: "collectibles", IsActive: true, SortOrder: 3},
		{Name: "Home & Living", Slug: "home-living", IsActive: true, SortOrder: 4},
		{Name: "Beauty & Health", Slug: "beauty-health", IsActive: true, SortOrder: 5},
		{Name: "Sports & Outdoors", Slug: "sports-outdoors", IsActive: true, SortOrder: 6},
		{Name: "Toys & Games", Slug: "toys-games", IsActive: true, SortOrder: 7},
		{Name: "Automotive", Slug: "automotive", IsActive: true, SortOrder: 8},
	}
	
	for _, cat := range categories {
		// Check if exists
		var existing CategoryModel
		result := db.Where("slug = ?", cat.Slug).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&cat).Error; err != nil {
				return err
			}
		}
	}
	
	return nil
}
