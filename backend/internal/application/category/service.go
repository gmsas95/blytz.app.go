package category

import (
	"context"
	"fmt"
	"strings"

	"github.com/blytz/live/backend/internal/domain/category"
	"github.com/google/uuid"
)

// Service handles category-related use cases
type Service struct {
	repo category.Repository
}

// NewService creates a new category service
func NewService(repo category.Repository) *Service {
	return &Service{repo: repo}
}

// CreateCategoryDTO represents data for creating a category
type CreateCategoryDTO struct {
	Name        string
	Description *string
	ImageURL    *string
	ParentID    *uuid.UUID
	SortOrder   int
}

// UpdateCategoryDTO represents data for updating a category
type UpdateCategoryDTO struct {
	Name        *string
	Description *string
	ImageURL    *string
	SortOrder   *int
	IsActive    *bool
}

// CreateCategory creates a new category
func (s *Service) CreateCategory(ctx context.Context, dto CreateCategoryDTO) (*category.Category, error) {
	slug := generateSlug(dto.Name)
	
	cat := category.NewCategory(dto.Name, slug, dto.ParentID)
	cat.Description = dto.Description
	cat.ImageURL = dto.ImageURL
	cat.SortOrder = dto.SortOrder
	
	if err := cat.Validate(); err != nil {
		return nil, err
	}
	
	if err := s.repo.Create(ctx, cat); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}
	
	return cat, nil
}

// GetCategory retrieves a category by ID
func (s *Service) GetCategory(ctx context.Context, id uuid.UUID) (*category.Category, error) {
	return s.repo.GetByID(ctx, id)
}

// GetCategoryBySlug retrieves a category by its slug
func (s *Service) GetCategoryBySlug(ctx context.Context, slug string) (*category.Category, error) {
	return s.repo.GetBySlug(ctx, slug)
}

// GetCategoryTree retrieves the full category tree
func (s *Service) GetCategoryTree(ctx context.Context) ([]*category.Category, error) {
	return s.repo.GetTree(ctx)
}

// ListCategories retrieves all categories
func (s *Service) ListCategories(ctx context.Context, onlyActive bool) ([]*category.Category, error) {
	return s.repo.List(ctx, onlyActive)
}

// UpdateCategory updates an existing category
func (s *Service) UpdateCategory(ctx context.Context, id uuid.UUID, dto UpdateCategoryDTO) (*category.Category, error) {
	cat, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if dto.Name != nil {
		cat.Name = *dto.Name
		cat.Slug = generateSlug(*dto.Name)
	}
	if dto.Description != nil {
		cat.Description = dto.Description
	}
	if dto.ImageURL != nil {
		cat.ImageURL = dto.ImageURL
	}
	if dto.SortOrder != nil {
		cat.SortOrder = *dto.SortOrder
	}
	if dto.IsActive != nil {
		cat.IsActive = *dto.IsActive
	}
	
	if err := cat.Validate(); err != nil {
		return nil, err
	}
	
	if err := s.repo.Update(ctx, cat); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}
	
	return cat, nil
}

// DeleteCategory deletes a category
func (s *Service) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	cat, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	if err := cat.CanDelete(); err != nil {
		return err
	}
	
	return s.repo.Delete(ctx, id)
}

// GetCategoryPath returns the full path from root to this category
func (s *Service) GetCategoryPath(ctx context.Context, id uuid.UUID) ([]string, error) {
	// Get the full tree
	tree, err := s.repo.GetTree(ctx)
	if err != nil {
		return nil, err
	}
	
	// Flatten the tree
	flat := category.FlattenTree(tree)
	
	// Find the category and build path
	var target *category.Category
	for _, c := range flat {
		if c.ID == id {
			target = c
			break
		}
	}
	
	if target == nil {
		return nil, category.ErrCategoryNotFound
	}
	
	return target.GetFullPath(), nil
}

// generateSlug creates a URL-friendly slug
func generateSlug(name string) string {
	// Simple slug generation
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	
	// Remove special characters
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	
	return result.String()
}
