package catalog

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/blytz.live.remake/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service interface for catalog operations
type Service interface {
	// Category Management
	CreateCategory(req *CategoryCreateRequest, sellerID uuid.UUID) (*models.Category, error)
	GetCategoryByID(id uuid.UUID) (*models.Category, error)
	UpdateCategory(id uuid.UUID, req *CategoryUpdateRequest) (*models.Category, error)
	DeleteCategory(id uuid.UUID) error
	GetCategoriesTree(req *CategoryTreeRequest) ([]CategoryResponse, error)
	MoveCategory(id uuid.UUID, req *CategoryMoveRequest) error
	BulkCategoryOperation(req *BulkCategoryOperation) (*BulkOperationResponse, error)

	// Category Attributes
	CreateCategoryAttribute(categoryID uuid.UUID, req *CategoryAttributeRequest) (*CategoryAttribute, error)
	GetCategoryAttributes(categoryID uuid.UUID) ([]CategoryAttribute, error)
	UpdateCategoryAttribute(id uuid.UUID, req *CategoryAttributeRequest) (*CategoryAttribute, error)
	DeleteCategoryAttribute(id uuid.UUID) error

	// Product Variants
	CreateProductVariant(productID uuid.UUID, req *ProductVariantRequest) (*ProductVariant, error)
	GetProductVariants(productID uuid.UUID) ([]ProductVariant, error)
	UpdateProductVariant(id uuid.UUID, req *ProductVariantRequest) (*ProductVariant, error)
	DeleteProductVariant(id uuid.UUID) error
	BulkCreateProductVariants(productID uuid.UUID, req []ProductVariantRequest) ([]ProductVariant, error)

	// Product Collections
	CreateProductCollection(req *ProductCollectionRequest) (*ProductCollection, error)
	GetProductCollectionByID(id uuid.UUID) (*ProductCollection, error)
	GetProductCollections(req *map[string]interface{}) ([]ProductCollection, error)
	UpdateProductCollection(id uuid.UUID, req *ProductCollectionRequest) (*ProductCollection, error)
	DeleteProductCollection(id uuid.UUID) error
	AddProductsToCollection(collectionID uuid.UUID, productIDs []uuid.UUID) error
	RemoveProductsFromCollection(collectionID uuid.UUID, productIDs []uuid.UUID) error

	// Inventory Management
	GetInventoryByProduct(productID uuid.UUID) (*InventoryStock, error)
	GetInventoryByVariant(variantID uuid.UUID) (*InventoryStock, error)
	UpdateInventory(productID uuid.UUID, req *InventoryStockRequest) (*InventoryStock, error)
	CreateStockMovement(productID uuid.UUID, req *StockMovementRequest) (*StockMovement, error)
	GetStockMovements(productID uuid.UUID, limit int) ([]StockMovement, error)
	GetLowStockProducts() ([]InventoryStock, error)
	GetOutOfStockProducts() ([]InventoryStock, error)

	// Catalog Statistics
	GetCatalogStats() (*CatalogStats, error)
	GetCategoryStats(categoryID uuid.UUID) (*CategoryStats, error)

	// Search and Filtering
	SearchProducts(query string, filters map[string]interface{}) ([]models.Product, int64, error)
	GetFeaturedProducts(limit int) ([]models.Product, error)
	GetRelatedProducts(productID uuid.UUID, limit int) ([]models.Product, error)
}

type service struct {
	db *gorm.DB
}

// NewService creates a new catalog service
func NewService(db *gorm.DB) Service {
	return &service{db: db}
}

// Category Management Implementation

func (s *service) CreateCategory(req *CategoryCreateRequest, sellerID uuid.UUID) (*models.Category, error) {
	// Validate parent category if provided
	if req.ParentID != nil {
		var parent models.Category
		if err := s.db.First(&parent, "id = ?", req.ParentID).Error; err != nil {
			return nil, fmt.Errorf("invalid parent category: %w", err)
		}
	}

	// Generate slug from name
	slug := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(req.Name), " ", "-"))
	
	// Check for duplicate slug
	var existing models.Category
	if err := s.db.Where("slug = ?", slug).First(&existing).Error; err == nil {
		return nil, errors.New("category with this name already exists")
	}

	category := &models.Category{
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		ParentID:    req.ParentID,
		SortOrder:   0,
		IsActive:    true,
	}

	if req.SortOrder != nil {
		category.SortOrder = *req.SortOrder
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := s.db.Create(category).Error; err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

func (s *service) GetCategoryByID(id uuid.UUID) (*models.Category, error) {
	var category models.Category
	if err := s.db.Preload("Parent").First(&category, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("category not found")
		}
		return nil, err
	}
	return &category, nil
}

func (s *service) UpdateCategory(id uuid.UUID, req *CategoryUpdateRequest) (*models.Category, error) {
	var category models.Category
	if err := s.db.First(&category, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	// Validate parent category if provided and different
	if req.ParentID != nil && (category.ParentID == nil || *req.ParentID != *category.ParentID) {
		if *req.ParentID == id {
			return nil, errors.New("category cannot be its own parent")
		}
		
		var parent models.Category
		if err := s.db.First(&parent, "id = ?", req.ParentID).Error; err != nil {
			return nil, fmt.Errorf("invalid parent category: %w", err)
		}
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
		// Update slug if name changed
		slug := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(*req.Name), " ", "-"))
		updates["slug"] = slug
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.ImageURL != nil {
		updates["image_url"] = *req.ImageURL
	}
	if req.ParentID != nil {
		updates["parent_id"] = *req.ParentID
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if err := s.db.Model(&category).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	// Reload the updated category
	return s.GetCategoryByID(id)
}

func (s *service) DeleteCategory(id uuid.UUID) error {
	// Check if category has children
	var childCount int64
	if err := s.db.Model(&models.Category{}).Where("parent_id = ?", id).Count(&childCount).Error; err != nil {
		return err
	}
	if childCount > 0 {
		return errors.New("cannot delete category with subcategories")
	}

	// Check if category has products
	var productCount int64
	if err := s.db.Model(&models.Product{}).Where("category_id = ?", id).Count(&productCount).Error; err != nil {
		return err
	}
	if productCount > 0 {
		return errors.New("cannot delete category with associated products")
	}

	if err := s.db.Delete(&models.Category{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

func (s *service) GetCategoriesTree(req *CategoryTreeRequest) ([]CategoryResponse, error) {
	var categories []models.Category
	query := s.db.Preload("Parent")

	if req.ParentID != nil {
		query = query.Where("parent_id = ?", req.ParentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}

	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}

	// Apply sorting
	sortBy := "sort_order"
	if req.SortBy != "" {
		sortBy = req.SortBy
	}
	sortDirection := "asc"
	if req.SortDirection != "" {
		sortDirection = req.SortDirection
	}
	query = query.Order(fmt.Sprintf("%s %s", sortBy, sortDirection))

	if err := query.Find(&categories).Error; err != nil {
		return nil, err
	}

	// Build response tree
	var result []CategoryResponse
	for _, category := range categories {
		response := s.buildCategoryResponse(&category, &req.IncludeProductCount)
		result = append(result, *response)

		// Recursively get children if needed
		if req.MaxDepth == nil || *req.MaxDepth > 1 {
			children, err := s.getCategoryChildren(category.ID, req, 1)
			if err == nil {
				response.Children = children
			}
		}
	}

	return result, nil
}

func (s *service) getCategoryChildren(parentID uuid.UUID, req *CategoryTreeRequest, currentDepth int) ([]CategoryResponse, error) {
	if req.MaxDepth != nil && currentDepth >= *req.MaxDepth {
		return nil, nil
	}

	var categories []models.Category
	query := s.db.Where("parent_id = ?", parentID).Preload("Parent")

	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}

	sortBy := "sort_order"
	if req.SortBy != "" {
		sortBy = req.SortBy
	}
	sortDirection := "asc"
	if req.SortDirection != "" {
		sortDirection = req.SortDirection
	}
	query = query.Order(fmt.Sprintf("%s %s", sortBy, sortDirection))

	if err := query.Find(&categories).Error; err != nil {
		return nil, err
	}

	var result []CategoryResponse
	for _, category := range categories {
		response := s.buildCategoryResponse(&category, &req.IncludeProductCount)
		
		// Recursively get children
		children, err := s.getCategoryChildren(category.ID, req, currentDepth+1)
		if err == nil {
			response.Children = children
		}
		
		result = append(result, *response)
	}

	return result, nil
}

func (s *service) buildCategoryResponse(category *models.Category, includeProductCount *bool) *CategoryResponse {
	response := &CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		ImageURL:    category.ImageURL,
		ParentID:    category.ParentID,
		SortOrder:   category.SortOrder,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
		Level:       s.calculateCategoryLevel(category),
		Path:        s.getCategoryPath(category),
	}

	if category.Parent != nil {
		response.Parent = &CategoryResponse{
			ID:   category.Parent.ID,
			Name: category.Parent.Name,
			Slug: category.Parent.Slug,
		}
	}

	// Include product count if requested
	if includeProductCount != nil && *includeProductCount {
		var count int64
		s.db.Model(&models.Product{}).Where("category_id = ?", category.ID).Count(&count)
		response.ProductCount = int(count)
	}

	return response
}

func (s *service) calculateCategoryLevel(category *models.Category) int {
	level := 0
	current := category
	for current.ParentID != nil {
		level++
		var parent models.Category
		if err := s.db.First(&parent, "id = ?", current.ParentID).Error; err != nil {
			break
		}
		current = &parent
	}
	return level
}

func (s *service) getCategoryPath(category *models.Category) []string {
	var path []string
	current := category
	path = append([]string{current.Name}, path...)
	
	for current.ParentID != nil {
		var parent models.Category
		if err := s.db.First(&parent, "id = ?", current.ParentID).Error; err != nil {
			break
		}
		path = append([]string{parent.Name}, path...)
		current = &parent
	}
	
	return path
}

func (s *service) MoveCategory(id uuid.UUID, req *CategoryMoveRequest) error {
	var category models.Category
	if err := s.db.First(&category, "id = ?", id).Error; err != nil {
		return fmt.Errorf("category not found: %w", err)
	}

	updates := make(map[string]interface{})
	
	if req.NewParentID != nil {
		if *req.NewParentID == id {
			return errors.New("category cannot be its own parent")
		}
		updates["parent_id"] = req.NewParentID
	}
	
	if req.NewSortOrder != nil {
		updates["sort_order"] = *req.NewSortOrder
	}

	if len(updates) > 0 {
		if err := s.db.Model(&category).Updates(updates).Error; err != nil {
			return fmt.Errorf("failed to move category: %w", err)
		}
	}

	return nil
}

func (s *service) BulkCategoryOperation(req *BulkCategoryOperation) (*BulkOperationResponse, error) {
	response := &BulkOperationResponse{
		TotalProcessed: len(req.CategoryIDs),
		SuccessCount:   0,
		FailureCount:   0,
		Errors:         []string{},
	}

	for _, categoryID := range req.CategoryIDs {
		var err error
		
		switch req.Operation {
		case "activate":
			err = s.db.Model(&models.Category{}).Where("id = ?", categoryID).Update("is_active", true).Error
		case "deactivate":
			err = s.db.Model(&models.Category{}).Where("id = ?", categoryID).Update("is_active", false).Error
		case "delete":
			err = s.DeleteCategory(categoryID)
		default:
			err = fmt.Errorf("invalid operation: %s", req.Operation)
		}

		if err != nil {
			response.FailureCount++
			response.Errors = append(response.Errors, fmt.Sprintf("Failed to %s category %s: %v", req.Operation, categoryID, err))
		} else {
			response.SuccessCount++
		}
	}

	return response, nil
}

// Category Attributes Implementation

func (s *service) CreateCategoryAttribute(categoryID uuid.UUID, req *CategoryAttributeRequest) (*CategoryAttribute, error) {
	// Validate category exists
	var category models.Category
	if err := s.db.First(&category, "id = ?", categoryID).Error; err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	attribute := &CategoryAttribute{
		CategoryID:  categoryID,
		Name:        req.Name,
		Type:        req.Type,
		Required:    req.Required != nil && *req.Required,
		Options:     req.Options,
		DefaultValue: req.DefaultValue,
		SortOrder:   0,
	}

	if req.SortOrder != nil {
		attribute.SortOrder = *req.SortOrder
	}

	// Validate select options
	if (req.Type == "select" || req.Type == "multiselect") && len(req.Options) == 0 {
		return nil, errors.New("select and multiselect types require options")
	}

	if err := s.db.Create(attribute).Error; err != nil {
		return nil, fmt.Errorf("failed to create category attribute: %w", err)
	}

	return attribute, nil
}

func (s *service) GetCategoryAttributes(categoryID uuid.UUID) ([]CategoryAttribute, error) {
	var attributes []CategoryAttribute
	if err := s.db.Where("category_id = ?", categoryID).Order("sort_order, name").Find(&attributes).Error; err != nil {
		return nil, fmt.Errorf("failed to get category attributes: %w", err)
	}
	return attributes, nil
}

func (s *service) UpdateCategoryAttribute(id uuid.UUID, req *CategoryAttributeRequest) (*CategoryAttribute, error) {
	var attribute CategoryAttribute
	if err := s.db.First(&attribute, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("attribute not found: %w", err)
	}

	// Validate select options
	if (req.Type == "select" || req.Type == "multiselect") && len(req.Options) == 0 {
		return nil, errors.New("select and multiselect types require options")
	}

	updates := make(map[string]interface{})
	updates["name"] = req.Name
	updates["type"] = req.Type
	updates["options"] = req.Options
	updates["default_value"] = req.DefaultValue
	updates["updated_at"] = time.Now()

	if req.Required != nil {
		updates["required"] = *req.Required
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}

	if err := s.db.Model(&attribute).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update attribute: %w", err)
	}

	// Reload the updated attribute
	if err := s.db.First(&attribute, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &attribute, nil
}

func (s *service) DeleteCategoryAttribute(id uuid.UUID) error {
	if err := s.db.Delete(&CategoryAttribute{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete attribute: %w", err)
	}
	return nil
}

// Product Variants Implementation

func (s *service) CreateProductVariant(productID uuid.UUID, req *ProductVariantRequest) (*ProductVariant, error) {
	// Validate product exists
	var product models.Product
	if err := s.db.First(&product, "id = ?", productID).Error; err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Check for duplicate SKU
	var existing ProductVariant
	if err := s.db.Where("sku = ?", req.Sku).First(&existing).Error; err == nil {
		return nil, errors.New("variant with this SKU already exists")
	}

	attributesJSON, _ := json.Marshal(req.Attributes)

	variant := &ProductVariant{
		ID:           uuid.New(),
		ProductID:    productID,
		Sku:          req.Sku,
		Title:        req.Title,
		Price:        req.Price,
		ComparePrice: req.ComparePrice,
		CostPrice:    req.CostPrice,
		Weight:       req.Weight,
		Barcode:      req.Barcode,
		Inventory:    0,
		IsActive:     true,
		Attributes:   string(attributesJSON),
		ImageURL:     req.ImageURL,
		Position:     0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if req.Inventory != nil {
		variant.Inventory = *req.Inventory
	}
	if req.IsActive != nil {
		variant.IsActive = *req.IsActive
	}
	if req.Position != nil {
		variant.Position = *req.Position
	}

	if err := s.db.Create(variant).Error; err != nil {
		return nil, fmt.Errorf("failed to create product variant: %w", err)
	}

	// Create inventory record for variant
	inventory := &InventoryStock{
		ID:              uuid.New(),
		ProductID:       productID,
		VariantID:       &variant.ID,
		Quantity:        variant.Inventory,
		Reserved:        0,
		Available:       variant.Inventory,
		LowStockAlert:   10,
		TrackInventory:  true,
		AllowBackorder:  false,
		LastUpdated:     time.Now(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	inventory.UpdateAvailability()
	s.db.Create(inventory)

	return variant, nil
}

func (s *service) GetProductVariants(productID uuid.UUID) ([]ProductVariant, error) {
	var variants []ProductVariant
	if err := s.db.Where("product_id = ?", productID).Order("position, title").Find(&variants).Error; err != nil {
		return nil, fmt.Errorf("failed to get product variants: %w", err)
	}
	return variants, nil
}

func (s *service) UpdateProductVariant(id uuid.UUID, req *ProductVariantRequest) (*ProductVariant, error) {
	var variant ProductVariant
	if err := s.db.First(&variant, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("variant not found: %w", err)
	}

	// Check for duplicate SKU (excluding current)
	var existing ProductVariant
	if err := s.db.Where("sku = ? AND id != ?", req.Sku, id).First(&existing).Error; err == nil {
		return nil, errors.New("variant with this SKU already exists")
	}

	attributesJSON, _ := json.Marshal(req.Attributes)

	updates := make(map[string]interface{})
	updates["title"] = req.Title
	updates["sku"] = req.Sku
	updates["price"] = req.Price
	updates["compare_price"] = req.ComparePrice
	updates["cost_price"] = req.CostPrice
	updates["weight"] = req.Weight
	updates["barcode"] = req.Barcode
	updates["attributes"] = string(attributesJSON)
	updates["image_url"] = req.ImageURL
	updates["updated_at"] = time.Now()

	if req.Inventory != nil {
		updates["inventory"] = *req.Inventory
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if req.Position != nil {
		updates["position"] = *req.Position
	}

	if err := s.db.Model(&variant).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update variant: %w", err)
	}

	// Update inventory if changed
	if req.Inventory != nil {
		s.db.Model(&InventoryStock{}).Where("variant_id = ?", id).Updates(map[string]interface{}{
			"quantity":    *req.Inventory,
			"available":   *req.Inventory,
			"last_updated": time.Now(),
		})
	}

	// Reload the updated variant
	if err := s.db.First(&variant, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &variant, nil
}

func (s *service) DeleteProductVariant(id uuid.UUID) error {
	var variant ProductVariant
	if err := s.db.First(&variant, "id = ?", id).Error; err != nil {
		return fmt.Errorf("variant not found: %w", err)
	}

	// Delete associated inventory record
	s.db.Delete(&InventoryStock{}, "variant_id = ?", id)

	// Delete variant
	if err := s.db.Delete(&variant).Error; err != nil {
		return fmt.Errorf("failed to delete variant: %w", err)
	}

	return nil
}

func (s *service) BulkCreateProductVariants(productID uuid.UUID, req []ProductVariantRequest) ([]ProductVariant, error) {
	var variants []ProductVariant
	for _, variantReq := range req {
		variant, err := s.CreateProductVariant(productID, &variantReq)
		if err != nil {
			return nil, fmt.Errorf("failed to create variant %s: %w", variantReq.Sku, err)
		}
		variants = append(variants, *variant)
	}
	return variants, nil
}

// Product Collections Implementation

func (s *service) CreateProductCollection(req *ProductCollectionRequest) (*ProductCollection, error) {
	// Generate slug from name
	slug := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(req.Name), " ", "-"))
	
	// Check for duplicate slug
	var existing ProductCollection
	if err := s.db.Where("slug = ?", slug).First(&existing).Error; err == nil {
		return nil, errors.New("collection with this name already exists")
	}

	productIDsJSON, _ := json.Marshal(req.ProductIDs)

	collection := &ProductCollection{
		ID:          uuid.New(),
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		IsActive:    true,
		SortOrder:   0,
		ProductIDs:  string(productIDsJSON),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if req.IsActive != nil {
		collection.IsActive = *req.IsActive
	}
	if req.SortOrder != nil {
		collection.SortOrder = *req.SortOrder
	}

	if err := s.db.Create(collection).Error; err != nil {
		return nil, fmt.Errorf("failed to create collection: %w", err)
	}

	return collection, nil
}

func (s *service) GetProductCollectionByID(id uuid.UUID) (*ProductCollection, error) {
	var collection ProductCollection
	if err := s.db.First(&collection, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("collection not found")
		}
		return nil, err
	}
	return &collection, nil
}

func (s *service) GetProductCollections(req *map[string]interface{}) ([]ProductCollection, error) {
	var collections []ProductCollection
	query := s.db.Where("is_active = ?", true)

	// Add filters from request parameters
	if params := *req; len(params) > 0 {
		if isActive, ok := params["is_active"]; ok {
			query = query.Where("is_active = ?", isActive)
		}
	}

	if err := query.Order("sort_order, name").Find(&collections).Error; err != nil {
		return nil, fmt.Errorf("failed to get collections: %w", err)
	}
	return collections, nil
}

func (s *service) UpdateProductCollection(id uuid.UUID, req *ProductCollectionRequest) (*ProductCollection, error) {
	var collection ProductCollection
	if err := s.db.First(&collection, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("collection not found: %w", err)
	}

	productIDsJSON, _ := json.Marshal(req.ProductIDs)

	updates := make(map[string]interface{})
	updates["name"] = req.Name
	updates["description"] = req.Description
	updates["image_url"] = req.ImageURL
	updates["product_ids"] = string(productIDsJSON)
	updates["updated_at"] = time.Now()

	// Update slug if name changed
	if req.Name != collection.Name {
		slug := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(req.Name), " ", "-"))
		updates["slug"] = slug
	}

	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}

	if err := s.db.Model(&collection).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update collection: %w", err)
	}

	// Reload the updated collection
	return s.GetProductCollectionByID(id)
}

func (s *service) DeleteProductCollection(id uuid.UUID) error {
	if err := s.db.Delete(&ProductCollection{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete collection: %w", err)
	}
	return nil
}

func (s *service) AddProductsToCollection(collectionID uuid.UUID, productIDs []uuid.UUID) error {
	var collection ProductCollection
	if err := s.db.First(&collection, "id = ?", collectionID).Error; err != nil {
		return fmt.Errorf("collection not found: %w", err)
	}

	// Get current product IDs
	var currentIDs []uuid.UUID
	json.Unmarshal([]byte(collection.ProductIDs), &currentIDs)

	// Add new IDs (avoid duplicates)
	idMap := make(map[uuid.UUID]bool)
	for _, id := range currentIDs {
		idMap[id] = true
	}
	for _, id := range productIDs {
		idMap[id] = true
	}

	// Convert back to slice
	var updatedIDs []uuid.UUID
	for id := range idMap {
		updatedIDs = append(updatedIDs, id)
	}

	productIDsJSON, _ := json.Marshal(updatedIDs)
	if err := s.db.Model(&collection).Updates(map[string]interface{}{
		"product_ids": string(productIDsJSON),
		"updated_at":  time.Now(),
	}).Error; err != nil {
		return fmt.Errorf("failed to add products to collection: %w", err)
	}

	return nil
}

func (s *service) RemoveProductsFromCollection(collectionID uuid.UUID, productIDs []uuid.UUID) error {
	var collection ProductCollection
	if err := s.db.First(&collection, "id = ?", collectionID).Error; err != nil {
		return fmt.Errorf("collection not found: %w", err)
	}

	// Get current product IDs
	var currentIDs []uuid.UUID
	json.Unmarshal([]byte(collection.ProductIDs), &currentIDs)

	// Create map of IDs to remove
	removeMap := make(map[uuid.UUID]bool)
	for _, id := range productIDs {
		removeMap[id] = true
	}

	// Filter out IDs to remove
	var updatedIDs []uuid.UUID
	for _, id := range currentIDs {
		if !removeMap[id] {
			updatedIDs = append(updatedIDs, id)
		}
	}

	productIDsJSON, _ := json.Marshal(updatedIDs)
	if err := s.db.Model(&collection).Updates(map[string]interface{}{
		"product_ids": string(productIDsJSON),
		"updated_at":  time.Now(),
	}).Error; err != nil {
		return fmt.Errorf("failed to remove products from collection: %w", err)
	}

	return nil
}

// Inventory Management Implementation

func (s *service) GetInventoryByProduct(productID uuid.UUID) (*InventoryStock, error) {
	var inventory InventoryStock
	if err := s.db.Where("product_id = ? AND variant_id IS NULL", productID).First(&inventory).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create default inventory record
			inventory = InventoryStock{
				ID:              uuid.New(),
				ProductID:       productID,
				Quantity:        0,
				Reserved:        0,
				Available:       0,
				LowStockAlert:   10,
				TrackInventory:  true,
				AllowBackorder:  false,
				LastUpdated:     time.Now(),
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			}
			inventory.UpdateAvailability()
			s.db.Create(&inventory)
		} else {
			return nil, err
		}
	}
	return &inventory, nil
}

func (s *service) GetInventoryByVariant(variantID uuid.UUID) (*InventoryStock, error) {
	var inventory InventoryStock
	if err := s.db.Where("variant_id = ?", variantID).First(&inventory).Error; err != nil {
		return nil, fmt.Errorf("inventory not found for variant: %w", err)
	}
	return &inventory, nil
}

func (s *service) UpdateInventory(productID uuid.UUID, req *InventoryStockRequest) (*InventoryStock, error) {
	inventory, err := s.GetInventoryByProduct(productID)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	updates["quantity"] = req.Quantity
	updates["last_updated"] = time.Now()
	updates["updated_at"] = time.Now()

	if req.LowStockAlert != nil {
		updates["low_stock_alert"] = *req.LowStockAlert
	}
	if req.TrackInventory != nil {
		updates["track_inventory"] = *req.TrackInventory
	}
	if req.AllowBackorder != nil {
		updates["allow_backorder"] = *req.AllowBackorder
	}
	if req.WarehouseID != nil {
		updates["warehouse_id"] = *req.WarehouseID
	}

	if err := s.db.Model(inventory).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update inventory: %w", err)
	}

	// Recalculate availability
	inventory.Quantity = req.Quantity
	inventory.UpdateAvailability()
	s.db.Model(inventory).Updates(map[string]interface{}{
		"available": inventory.Available,
	})

	// Reload updated inventory
	return s.GetInventoryByProduct(productID)
}

func (s *service) CreateStockMovement(productID uuid.UUID, req *StockMovementRequest) (*StockMovement, error) {
	// Validate product exists
	var product models.Product
	if err := s.db.First(&product, "id = ?", productID).Error; err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	movement := &StockMovement{
		ID:           uuid.New(),
		ProductID:    productID,
		MovementType: req.MovementType,
		Quantity:     req.Quantity,
		Reference:    req.Reference,
		Notes:        req.Notes,
		WarehouseID:  req.WarehouseID,
		CreatedAt:    time.Now(),
	}

	if err := s.db.Create(movement).Error; err != nil {
		return nil, fmt.Errorf("failed to create stock movement: %w", err)
	}

	// Update inventory based on movement type
	inventory, err := s.GetInventoryByProduct(productID)
	if err != nil {
		return movement, err
	}

	switch req.MovementType {
	case "in", "release":
		inventory.Quantity += req.Quantity
	case "out", "reserve":
		inventory.Reserved += req.Quantity
	case "adjustment":
		inventory.Quantity = req.Quantity
		inventory.Reserved = 0
	}

	inventory.UpdateAvailability()
	s.db.Model(inventory).Updates(map[string]interface{}{
		"quantity":     inventory.Quantity,
		"reserved":     inventory.Reserved,
		"available":    inventory.Available,
		"last_updated": time.Now(),
	})

	return movement, nil
}

func (s *service) GetStockMovements(productID uuid.UUID, limit int) ([]StockMovement, error) {
	var movements []StockMovement
	query := s.db.Where("product_id = ?", productID).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&movements).Error; err != nil {
		return nil, fmt.Errorf("failed to get stock movements: %w", err)
	}
	return movements, nil
}

func (s *service) GetLowStockProducts() ([]InventoryStock, error) {
	var inventory []InventoryStock
	if err := s.db.Where("track_inventory = ? AND available <= low_stock_alert AND available > 0", true).
		Preload("Product").Find(&inventory).Error; err != nil {
		return nil, fmt.Errorf("failed to get low stock products: %w", err)
	}
	return inventory, nil
}

func (s *service) GetOutOfStockProducts() ([]InventoryStock, error) {
	var inventory []InventoryStock
	if err := s.db.Where("track_inventory = ? AND available = ? AND allow_backorder = ?", true, 0, false).
		Preload("Product").Find(&inventory).Error; err != nil {
		return nil, fmt.Errorf("failed to get out of stock products: %w", err)
	}
	return inventory, nil
}

// Catalog Statistics Implementation

func (s *service) GetCatalogStats() (*CatalogStats, error) {
	stats := &CatalogStats{}

	var totalCategories int64
	var totalProducts int64
	var activeProducts int64
	var totalVariants int64
	
	// Category count
	s.db.Model(&models.Category{}).Count(&totalCategories)
	stats.TotalCategories = int(totalCategories)

	// Product counts
	s.db.Model(&models.Product{}).Count(&totalProducts)
	stats.TotalProducts = int(totalProducts)
	s.db.Model(&models.Product{}).Where("status = 'active'").Count(&activeProducts)
	stats.ActiveProducts = int(activeProducts)

	// Variant count
	s.db.Model(&ProductVariant{}).Count(&totalVariants)
	stats.TotalVariants = int(totalVariants)

	// Low stock and out of stock products
	lowStock, _ := s.GetLowStockProducts()
	stats.LowStockProducts = len(lowStock)
	outOfStock, _ := s.GetOutOfStockProducts()
	stats.OutofStockProducts = len(outOfStock)

	// Products by category
	var results []struct {
		CategoryName string
		Count        int64
	}
	s.db.Model(&models.Product{}).
		Select("categories.name as category_name, COUNT(*) as count").
		Joins("LEFT JOIN categories ON products.category_id = categories.id").
		Group("categories.name").
		Scan(&results)

	stats.ProductsByCategory = make(map[string]int)
	for _, result := range results {
		stats.ProductsByCategory[result.CategoryName] = int(result.Count)
	}

	// Top categories
	var topCategories []struct {
		ID           uuid.UUID
		Name         string
		ProductCount int64
	}
	s.db.Model(&models.Category{}).
		Select("categories.id, categories.name, COUNT(products.id) as product_count").
		Joins("LEFT JOIN products ON categories.id = products.category_id").
		Group("categories.id, categories.name").
		Order("product_count DESC").
		Limit(5).
		Scan(&topCategories)

	for _, cat := range topCategories {
		stats.TopCategories = append(stats.TopCategories, CategoryStats{
			ID:           cat.ID,
			Name:         cat.Name,
			ProductCount: int(cat.ProductCount),
		})
	}

	// Recent products
	var recentProducts []models.Product
	s.db.Model(&models.Product{}).
		Where("status = 'active'").
		Order("created_at DESC").
		Limit(5).
		Find(&recentProducts)

	for _, product := range recentProducts {
		stats.RecentProducts = append(stats.RecentProducts, ProductSummary{
			ID:        product.ID,
			Title:     product.Title,
			Price:     product.StartingPrice,
			CreatedAt: product.CreatedAt,
		})
	}

	return stats, nil
}

func (s *service) GetCategoryStats(categoryID uuid.UUID) (*CategoryStats, error) {
	stats := &CategoryStats{}

	// Get category info
	var category models.Category
	if err := s.db.First(&category, "id = ?", categoryID).Error; err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	stats.ID = category.ID
	stats.Name = category.Name

	// Product counts
	var productCount int64
	var activeCount int64
	s.db.Model(&models.Product{}).Where("category_id = ?", categoryID).Count(&productCount)
	s.db.Model(&models.Product{}).Where("category_id = ? AND status = 'active'", categoryID).Count(&activeCount)
	stats.ProductCount = int(productCount)
	stats.ActiveCount = int(activeCount)

	// Total revenue (would need to join with orders - placeholder for now)
	stats.TotalRevenue = 0.0

	return stats, nil
}

// Search and Filtering Implementation

func (s *service) SearchProducts(query string, filters map[string]interface{}) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	dbQuery := s.db.Model(&models.Product{}).Preload("Category").Preload("Seller")

	// Text search
	if query != "" {
		dbQuery = dbQuery.Where("title LIKE ? OR description LIKE ?", 
			"%"+query+"%", "%"+query+"%")
	}

	// Apply filters
	if categoryID, ok := filters["category_id"]; ok {
		dbQuery = dbQuery.Where("category_id = ?", categoryID)
	}
	if sellerID, ok := filters["seller_id"]; ok {
		dbQuery = dbQuery.Where("seller_id = ?", sellerID)
	}
	if status, ok := filters["status"]; ok {
		dbQuery = dbQuery.Where("status = ?", status)
	}
	if featured, ok := filters["featured"]; ok {
		dbQuery = dbQuery.Where("featured = ?", featured)
	}
	if condition, ok := filters["condition"]; ok {
		dbQuery = dbQuery.Where("condition = ?", condition)
	}
	if minPrice, ok := filters["min_price"]; ok {
		dbQuery = dbQuery.Where("starting_price >= ?", minPrice)
	}
	if maxPrice, ok := filters["max_price"]; ok {
		dbQuery = dbQuery.Where("starting_price <= ?", maxPrice)
	}

	// Count total
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if page, ok := filters["page"]; ok {
		if limit, ok := filters["limit"]; ok {
			offset := (page.(int) - 1) * limit.(int)
			dbQuery = dbQuery.Offset(offset).Limit(limit.(int))
		}
	}

	// Apply sorting
	sortBy := "created_at"
	if sort, ok := filters["sort_by"]; ok {
		sortBy = sort.(string)
	}
	sortDirection := "desc"
	if dir, ok := filters["sort_direction"]; ok {
		sortDirection = dir.(string)
	}
	dbQuery = dbQuery.Order(fmt.Sprintf("%s %s", sortBy, sortDirection))

	if err := dbQuery.Find(&products).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search products: %w", err)
	}

	return products, total, nil
}

func (s *service) GetFeaturedProducts(limit int) ([]models.Product, error) {
	var products []models.Product
	if err := s.db.Where("featured = ? AND status = 'active'", true).
		Preload("Category").Preload("Seller").
		Order("created_at DESC").
		Limit(limit).
		Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to get featured products: %w", err)
	}
	return products, nil
}

func (s *service) GetRelatedProducts(productID uuid.UUID, limit int) ([]models.Product, error) {
	// Get product info
	var product models.Product
	if err := s.db.First(&product, "id = ?", productID).Error; err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	var products []models.Product

	// Find related products by category (excluding current product)
	if err := s.db.Where("category_id = ? AND id != ? AND status = 'active'", 
		product.CategoryID, productID).
		Preload("Category").Preload("Seller").
		Order("created_at DESC").
		Limit(limit).
		Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to get related products: %w", err)
	}

	// If not enough products by category, get more by similar price range
	if len(products) < limit {
		remaining := limit - len(products)
		priceRange := product.StartingPrice * 0.2 // 20% price range

		var additionalProducts []models.Product
		if err := s.db.Where("id NOT IN (?) AND id != ? AND status = 'active' AND starting_price BETWEEN ? AND ?", 
			append(getProductIDs(products), productID), 
			product.StartingPrice-priceRange, product.StartingPrice+priceRange).
			Preload("Category").Preload("Seller").
			Order("created_at DESC").
			Limit(remaining).
			Find(&additionalProducts).Error; err != nil {
			return nil, err
		}

		products = append(products, additionalProducts...)
	}

	return products, nil
}

// Helper functions
func getProductIDs(products []models.Product) []uuid.UUID {
	var ids []uuid.UUID
	for _, product := range products {
		ids = append(ids, product.ID)
	}
	return ids
}