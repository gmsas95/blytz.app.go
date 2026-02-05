package category

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Errors
var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrCircularReference = errors.New("circular reference detected")
	ErrHasProducts      = errors.New("category has products")
	ErrHasSubcategories = errors.New("category has subcategories")
)

// Category represents a product category
type Category struct {
	ID          uuid.UUID
	Name        string
	Slug        string
	Description *string
	ImageURL    *string
	ParentID    *uuid.UUID
	Parent      *Category
	Children    []*Category
	SortOrder   int
	IsActive    bool
	ProductCount int // Cached count
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// IsRoot returns true if this is a top-level category
func (c *Category) IsRoot() bool {
	return c.ParentID == nil
}

// GetFullPath returns the full path from root to this category
func (c *Category) GetFullPath() []string {
	path := []string{c.Name}
	current := c
	
	for current.Parent != nil {
		path = append([]string{current.Parent.Name}, path...)
		current = current.Parent
	}
	
	return path
}

// CanDelete checks if the category can be deleted
func (c *Category) CanDelete() error {
	if c.ProductCount > 0 {
		return ErrHasProducts
	}
	
	if len(c.Children) > 0 {
		return ErrHasSubcategories
	}
	
	return nil
}

// Validate performs validation
func (c *Category) Validate() error {
	if c.Name == "" {
		return errors.New("category name is required")
	}
	
	if c.Slug == "" {
		return errors.New("category slug is required")
	}
	
	// Check for circular reference
	if c.ParentID != nil && *c.ParentID == c.ID {
		return ErrCircularReference
	}
	
	return nil
}

// Repository defines the interface for category data access
type Repository interface {
	// GetByID retrieves a category by ID
	GetByID(ctx context.Context, id uuid.UUID) (*Category, error)
	
	// GetBySlug retrieves a category by its slug
	GetBySlug(ctx context.Context, slug string) (*Category, error)
	
	// GetTree retrieves the full category tree
	GetTree(ctx context.Context) ([]*Category, error)
	
	// GetChildren retrieves immediate children of a category
	GetChildren(ctx context.Context, parentID *uuid.UUID) ([]*Category, error)
	
	// List retrieves all categories (flat list)
	List(ctx context.Context, onlyActive bool) ([]*Category, error)
	
	// Create creates a new category
	Create(ctx context.Context, category *Category) error
	
	// Update updates an existing category
	Update(ctx context.Context, category *Category) error
	
	// Delete soft-deletes a category
	Delete(ctx context.Context, id uuid.UUID) error
	
	// UpdateProductCount updates the cached product count
	UpdateProductCount(ctx context.Context, id uuid.UUID) error
}

// NewCategory creates a new category
func NewCategory(name, slug string, parentID *uuid.UUID) *Category {
	now := time.Now()
	return &Category{
		ID:        uuid.New(),
		Name:      name,
		Slug:      slug,
		ParentID:  parentID,
		Children:  make([]*Category, 0),
		SortOrder: 0,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// FlattenTree flattens a category tree to a slice
func FlattenTree(categories []*Category) []*Category {
	result := make([]*Category, 0)
	
	var flatten func(cats []*Category)
	flatten = func(cats []*Category) {
		for _, cat := range cats {
			result = append(result, cat)
			if len(cat.Children) > 0 {
				flatten(cat.Children)
			}
		}
	}
	
	flatten(categories)
	return result
}

// BuildTree builds a tree structure from flat categories
func BuildTree(categories []*Category) []*Category {
	// Map for quick lookup
	categoryMap := make(map[uuid.UUID]*Category)
	for _, cat := range categories {
		categoryMap[cat.ID] = cat
		cat.Children = make([]*Category, 0) // Reset children
	}
	
	// Build tree
	roots := make([]*Category, 0)
	for _, cat := range categories {
		if cat.ParentID == nil {
			roots = append(roots, cat)
		} else {
			if parent, ok := categoryMap[*cat.ParentID]; ok {
				parent.Children = append(parent.Children, cat)
				cat.Parent = parent
			}
		}
	}
	
	return roots
}
