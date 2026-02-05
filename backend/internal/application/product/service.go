package product

import (
	"context"
	"fmt"

	"github.com/blytz/live/backend/internal/domain/category"
	"github.com/blytz/live/backend/internal/domain/product"
	"github.com/google/uuid"
)

// Service handles product-related use cases
type Service struct {
	productRepo  product.Repository
	categoryRepo category.Repository
}

// NewService creates a new product service
func NewService(productRepo product.Repository, categoryRepo category.Repository) *Service {
	return &Service{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

// CreateProductDTO represents data for creating a product
type CreateProductDTO struct {
	SellerID       uuid.UUID
	CategoryID     *uuid.UUID
	Name           string
	Description    string
	Condition      product.Condition
	BasePrice      float64
	CompareAtPrice *float64
	StockQuantity  int
	SKU            *string
	WeightGrams    *int
	DimensionsCm   *product.Dimensions
	Attributes     map[string]string
	Images         []CreateImageDTO
}

// CreateImageDTO represents image data for creating a product
type CreateImageDTO struct {
	URL          string
	ThumbnailURL *string
	AltText      *string
	IsPrimary    bool
}

// UpdateProductDTO represents data for updating a product
type UpdateProductDTO struct {
	CategoryID     *uuid.UUID
	Name           *string
	Description    *string
	Condition      *product.Condition
	BasePrice      *float64
	CompareAtPrice *float64
	StockQuantity  *int
	SKU            *string
	WeightGrams    *int
	DimensionsCm   *product.Dimensions
	Attributes     map[string]string
}

// ListProductsDTO represents filter criteria for listing products
type ListProductsDTO struct {
	SellerID    *uuid.UUID
	CategoryID  *uuid.UUID
	Status      *product.Status
	Condition   *product.Condition
	MinPrice    *float64
	MaxPrice    *float64
	Query       string
	SortBy      string
	Page        int
	PageSize    int
}

// ListResult represents the result of listing products
type ListResult struct {
	Products   []*product.Product
	TotalCount int64
	Page       int
	PageSize   int
}

// CreateProduct creates a new product
func (s *Service) CreateProduct(ctx context.Context, dto CreateProductDTO) (*product.Product, error) {
	// Validate category if provided
	if dto.CategoryID != nil {
		_, err := s.categoryRepo.GetByID(ctx, *dto.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("invalid category: %w", err)
		}
	}
	
	// Create product
	p := product.NewProduct(
		dto.SellerID,
		dto.Name,
		dto.Description,
		dto.Condition,
		dto.BasePrice,
		dto.StockQuantity,
	)
	
	p.CategoryID = dto.CategoryID
	p.CompareAtPrice = dto.CompareAtPrice
	p.SKU = dto.SKU
	p.WeightGrams = dto.WeightGrams
	p.DimensionsCm = dto.DimensionsCm
	
	if dto.Attributes != nil {
		p.Attributes = dto.Attributes
	}
	
	// Validate
	if err := p.Validate(); err != nil {
		return nil, err
	}
	
	// Add images
	for i, imgDTO := range dto.Images {
		p.AddImage(imgDTO.URL, i == 0 || imgDTO.IsPrimary)
	}
	
	// Save
	if err := s.productRepo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	
	return p, nil
}

// GetProduct retrieves a product by ID
func (s *Service) GetProduct(ctx context.Context, id uuid.UUID) (*product.Product, error) {
	p, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Increment view count asynchronously
	go s.productRepo.IncrementViewCount(context.Background(), id)
	
	return p, nil
}

// GetProductBySlug retrieves a product by its slug
func (s *Service) GetProductBySlug(ctx context.Context, slug string) (*product.Product, error) {
	p, err := s.productRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	
	// Increment view count asynchronously
	go s.productRepo.IncrementViewCount(context.Background(), p.ID)
	
	return p, nil
}

// UpdateProduct updates an existing product
func (s *Service) UpdateProduct(ctx context.Context, productID uuid.UUID, sellerID uuid.UUID, dto UpdateProductDTO) (*product.Product, error) {
	// Get existing product
	p, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	
	// Verify ownership
	if !p.CanEdit(sellerID) {
		return nil, product.ErrUnauthorized
	}
	
	// Validate category if changed
	if dto.CategoryID != nil && (p.CategoryID == nil || *dto.CategoryID != *p.CategoryID) {
		_, err := s.categoryRepo.GetByID(ctx, *dto.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("invalid category: %w", err)
		}
		p.CategoryID = dto.CategoryID
	}
	
	// Update fields
	if dto.Name != nil {
		p.Name = *dto.Name
	}
	if dto.Description != nil {
		p.Description = *dto.Description
	}
	if dto.Condition != nil {
		p.Condition = *dto.Condition
	}
	if dto.BasePrice != nil {
		p.BasePrice = *dto.BasePrice
	}
	if dto.CompareAtPrice != nil {
		p.CompareAtPrice = dto.CompareAtPrice
	}
	if dto.StockQuantity != nil {
		p.StockQuantity = *dto.StockQuantity
	}
	if dto.SKU != nil {
		p.SKU = dto.SKU
	}
	if dto.WeightGrams != nil {
		p.WeightGrams = dto.WeightGrams
	}
	if dto.DimensionsCm != nil {
		p.DimensionsCm = dto.DimensionsCm
	}
	if dto.Attributes != nil {
		p.Attributes = dto.Attributes
	}
	
	// Validate and save
	if err := p.Validate(); err != nil {
		return nil, err
	}
	
	if err := s.productRepo.Update(ctx, p); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}
	
	return p, nil
}

// DeleteProduct soft-deletes a product
func (s *Service) DeleteProduct(ctx context.Context, productID uuid.UUID, sellerID uuid.UUID) error {
	p, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return err
	}
	
	// Verify ownership
	if p.SellerID != sellerID {
		return product.ErrUnauthorized
	}
	
	// Can't delete sold products
	if p.Status == product.StatusSold {
		return product.ErrProductAlreadySold
	}
	
	return s.productRepo.Delete(ctx, productID)
}

// ListProducts retrieves a list of products with filtering
func (s *Service) ListProducts(ctx context.Context, dto ListProductsDTO) (*ListResult, error) {
	// Set defaults
	if dto.Page <= 0 {
		dto.Page = 1
	}
	if dto.PageSize <= 0 {
		dto.PageSize = 20
	}
	if dto.PageSize > 100 {
		dto.PageSize = 100
	}
	
	filter := product.Filter{
		SellerID:   dto.SellerID,
		CategoryID: dto.CategoryID,
		Status:     dto.Status,
		Condition:  dto.Condition,
		MinPrice:   dto.MinPrice,
		MaxPrice:   dto.MaxPrice,
		Query:      dto.Query,
		SortBy:     dto.SortBy,
		Page:       dto.Page,
		PageSize:   dto.PageSize,
	}
	
	products, total, err := s.productRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	
	return &ListResult{
		Products:   products,
		TotalCount: total,
		Page:       dto.Page,
		PageSize:   dto.PageSize,
	}, nil
}

// PublishProduct changes product status from draft to active
func (s *Service) PublishProduct(ctx context.Context, productID uuid.UUID, sellerID uuid.UUID) error {
	p, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return err
	}
	
	if p.SellerID != sellerID {
		return product.ErrUnauthorized
	}
	
	if err := p.Publish(); err != nil {
		return err
	}
	
	return s.productRepo.Update(ctx, p)
}

// ArchiveProduct changes product status to archived
func (s *Service) ArchiveProduct(ctx context.Context, productID uuid.UUID, sellerID uuid.UUID) error {
	p, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return err
	}
	
	if p.SellerID != sellerID {
		return product.ErrUnauthorized
	}
	
	if err := p.Archive(); err != nil {
		return err
	}
	
	return s.productRepo.Update(ctx, p)
}

// UpdateStock updates the stock quantity
func (s *Service) UpdateStock(ctx context.Context, productID uuid.UUID, sellerID uuid.UUID, quantity int) error {
	p, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return err
	}
	
	if p.SellerID != sellerID {
		return product.ErrUnauthorized
	}
	
	if quantity < 0 {
		return product.ErrInvalidStock
	}
	
	return s.productRepo.UpdateStock(ctx, productID, quantity)
}

// GetSellerProducts retrieves all products for a seller
func (s *Service) GetSellerProducts(ctx context.Context, sellerID uuid.UUID, status *product.Status, page, pageSize int) (*ListResult, error) {
	return s.ListProducts(ctx, ListProductsDTO{
		SellerID: sellerID,
		Status:   status,
		Page:     page,
		PageSize: pageSize,
	})
}

// AddProductImage adds an image to a product
func (s *Service) AddProductImage(ctx context.Context, productID uuid.UUID, sellerID uuid.UUID, url string, isPrimary bool) (*product.ProductImage, error) {
	p, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	
	if !p.CanEdit(sellerID) {
		return nil, product.ErrUnauthorized
	}
	
	img := p.AddImage(url, isPrimary)
	
	if err := s.productRepo.Update(ctx, p); err != nil {
		return nil, err
	}
	
	return &img, nil
}

// SetPrimaryImage sets an image as the primary image
func (s *Service) SetPrimaryImage(ctx context.Context, productID uuid.UUID, sellerID uuid.UUID, imageID uuid.UUID) error {
	p, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return err
	}
	
	if !p.CanEdit(sellerID) {
		return product.ErrUnauthorized
	}
	
	if err := p.SetPrimaryImage(imageID); err != nil {
		return err
	}
	
	return s.productRepo.Update(ctx, p)
}
