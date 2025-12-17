package catalog

import (
	"net/http"
	"strconv"

	"github.com/blytz.live.remake/backend/internal/auth"
	pkghttp "github.com/blytz.live.remake/backend/pkg/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles catalog-related HTTP requests
type Handler struct {
	service Service
}

// NewHandler creates a new catalog handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers catalog routes with router
func (h *Handler) RegisterRoutes(router *gin.RouterGroup, authHandler *auth.Handler) {
	// Public routes
	router.GET("/categories", h.GetCategoriesTree)
	router.GET("/categories/:id", h.GetCategoryByID)
	router.GET("/collections", h.GetProductCollections)
	router.GET("/collections/:id", h.GetProductCollectionByID)
	router.GET("/search/products", h.SearchProducts)
	router.GET("/search/products/featured", h.GetFeaturedProducts)
	router.GET("/search/products/:id/related", h.GetRelatedProducts)
	router.GET("/stats/catalog", h.GetCatalogStats)
	router.GET("/stats/categories/:id", h.GetCategoryStats)

	// Category Management
	categories := router.Group("/categories")
	{
		categories.POST("", authHandler.RequireAuth(), authHandler.RequireSellerOrAdmin(), h.CreateCategory)
		categories.PUT("/:id", authHandler.RequireAuth(), authHandler.RequireSellerOrAdmin(), h.UpdateCategory)
		categories.DELETE("/:id", authHandler.RequireAuth(), authHandler.RequireSellerOrAdmin(), h.DeleteCategory)
		categories.PUT("/:id/move", authHandler.RequireAuth(), authHandler.RequireSellerOrAdmin(), h.MoveCategory)
		categories.POST("/bulk", authHandler.RequireAuth(), authHandler.RequireRole("admin"), h.BulkCategoryOperation)

		// Category Attributes
		categories.GET("/:id/attributes", h.GetCategoryAttributes)
		categories.POST("/:id/attributes", authHandler.RequireAuth(), authHandler.RequireSellerOrAdmin(), h.CreateCategoryAttribute)
		categories.PUT("/attributes/:attributeId", authHandler.RequireAuth(), authHandler.RequireSellerOrAdmin(), h.UpdateCategoryAttribute)
		categories.DELETE("/attributes/:attributeId", authHandler.RequireAuth(), authHandler.RequireSellerOrAdmin(), h.DeleteCategoryAttribute)
	}

	// Product Variants
	variants := router.Group("/variants")
	{
		variants.GET("/products/:productId", h.GetProductVariants)
		variants.POST("", authHandler.RequireAuth(), authHandler.RequireSellerOrAdmin(), h.CreateProductVariant)
		variants.PUT("/:id", authHandler.RequireAuth(), authHandler.RequireSellerOrAdmin(), h.UpdateProductVariant)
		variants.DELETE("/:id", authHandler.RequireAuth(), authHandler.RequireSellerOrAdmin(), h.DeleteProductVariant)
		variants.POST("/bulk", authHandler.RequireAuth(), authHandler.RequireSellerOrAdmin(), h.BulkCreateProductVariants)
	}

	// Product Collections
	collections := router.Group("/collections")
	{
		collections.POST("", authHandler.RequireAuth(), authHandler.RequireSellerOrAdmin(), h.CreateProductCollection)
		collections.PUT("/:id", authHandler.RequireAuth(), authHandler.RequireSellerOrAdmin(), h.UpdateProductCollection)
		collections.DELETE("/:id", authHandler.RequireAuth(), authHandler.RequireSellerOrAdmin(), h.DeleteProductCollection)
		collections.POST("/:id/products", authHandler.RequireAuth(), authHandler.RequireSellerOrAdmin(), h.AddProductsToCollection)
		collections.DELETE("/:id/products", authHandler.RequireAuth(), authHandler.RequireSellerOrAdmin(), h.RemoveProductsFromCollection)
	}

	// Inventory Management
	inventory := router.Group("/inventory")
	{
		inventory.GET("/products/:productId", h.GetInventoryByProduct)
		inventory.GET("/variants/:variantId", h.GetInventoryByVariant)
		inventory.PUT("/products/:productId", authHandler.RequireAuth(), authHandler.RequireSellerOrAdmin(), h.UpdateInventory)
		inventory.GET("/products/:productId/movements", h.GetStockMovements)
		inventory.POST("/products/:productId/movements", authHandler.RequireAuth(), authHandler.RequireSellerOrAdmin(), h.CreateStockMovement)
		inventory.GET("/low-stock", authHandler.RequireAuth(), authHandler.RequireSellerOrAdmin(), h.GetLowStockProducts)
		inventory.GET("/out-of-stock", authHandler.RequireAuth(), authHandler.RequireSellerOrAdmin(), h.GetOutOfStockProducts)
	}
}

// Category Handlers

func (h *Handler) CreateCategory(c *gin.Context) {
	var req CategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		pkghttp.Error(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	id, err := uuid.Parse(userID.(string))
	if err != nil {
		pkghttp.Error(c, http.StatusInternalServerError, "Invalid user ID", err)
		return
	}

	category, err := h.service.CreateCategory(&req, id)
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Failed to create category", err)
		return
	}

	pkghttp.Success(c, http.StatusCreated, category)
}

func (h *Handler) GetCategoryByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid category ID", err)
		return
	}

	category, err := h.service.GetCategoryByID(id)
	if err != nil {
		pkghttp.Error(c, http.StatusNotFound, "Category not found", err)
		return
	}

	pkghttp.Success(c, http.StatusOK, category)
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid category ID", err)
		return
	}

	var req CategoryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	category, err := h.service.UpdateCategory(id, &req)
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Failed to update category", err)
		return
	}

	pkghttp.Success(c, http.StatusOK, category)
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid category ID", err)
		return
	}

	if err := h.service.DeleteCategory(id); err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Failed to delete category", err)
		return
	}

	pkghttp.Success(c, http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

func (h *Handler) GetCategoriesTree(c *gin.Context) {
	var req CategoryTreeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid query parameters", err)
		return
	}

	categories, err := h.service.GetCategoriesTree(&req)
	if err != nil {
		pkghttp.Error(c, http.StatusInternalServerError, "Failed to get categories", err)
		return
	}

	pkghttp.Success(c, http.StatusOK, categories)
}

func (h *Handler) MoveCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid category ID", err)
		return
	}

	var req CategoryMoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := h.service.MoveCategory(id, &req); err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Failed to move category", err)
		return
	}

	pkghttp.Success(c, http.StatusOK, gin.H{"message": "Category moved successfully"})
}

func (h *Handler) BulkCategoryOperation(c *gin.Context) {
	var req BulkCategoryOperation
	if err := c.ShouldBindJSON(&req); err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	result, err := h.service.BulkCategoryOperation(&req)
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Failed to perform bulk operation", err)
		return
	}

	pkghttp.Success(c, http.StatusOK, result)
}

// Category Attribute Handlers

func (h *Handler) CreateCategoryAttribute(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid category ID", err)
		return
	}

	var req CategoryAttributeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	attribute, err := h.service.CreateCategoryAttribute(categoryID, &req)
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Failed to create category attribute", err)
		return
	}

	pkghttp.Success(c, http.StatusCreated, attribute)
}

func (h *Handler) GetCategoryAttributes(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid category ID", err)
		return
	}

	attributes, err := h.service.GetCategoryAttributes(categoryID)
	if err != nil {
		pkghttp.Error(c, http.StatusInternalServerError, "Failed to get category attributes", err)
		return
	}

	pkghttp.Success(c, http.StatusOK, attributes)
}

func (h *Handler) UpdateCategoryAttribute(c *gin.Context) {
	attributeID, err := uuid.Parse(c.Param("attributeId"))
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid attribute ID", err)
		return
	}

	var req CategoryAttributeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	attribute, err := h.service.UpdateCategoryAttribute(attributeID, &req)
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Failed to update attribute", err)
		return
	}

	pkghttp.Success(c, http.StatusOK, attribute)
}

func (h *Handler) DeleteCategoryAttribute(c *gin.Context) {
	attributeID, err := uuid.Parse(c.Param("attributeId"))
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid attribute ID", err)
		return
	}

	if err := h.service.DeleteCategoryAttribute(attributeID); err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Failed to delete attribute", err)
		return
	}

	pkghttp.Success(c, http.StatusOK, gin.H{"message": "Attribute deleted successfully"})
}

// Product Variant Handlers

func (h *Handler) CreateProductVariant(c *gin.Context) {
	var req struct {
		ProductID uuid.UUID             `json:"product_id" binding:"required"`
		Variant   ProductVariantRequest `json:"variant" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	variant, err := h.service.CreateProductVariant(req.ProductID, &req.Variant)
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Failed to create product variant", err)
		return
	}

	pkghttp.Success(c, http.StatusCreated, variant)
}

func (h *Handler) GetProductVariants(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid product ID", err)
		return
	}

	variants, err := h.service.GetProductVariants(productID)
	if err != nil {
		pkghttp.Error(c, http.StatusInternalServerError, "Failed to get product variants", err)
		return
	}

	pkghttp.Success(c, http.StatusOK, variants)
}

func (h *Handler) UpdateProductVariant(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid variant ID", err)
		return
	}

	var req ProductVariantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	variant, err := h.service.UpdateProductVariant(id, &req)
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Failed to update variant", err)
		return
	}

	pkghttp.Success(c, http.StatusOK, variant)
}

func (h *Handler) DeleteProductVariant(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid variant ID", err)
		return
	}

	if err := h.service.DeleteProductVariant(id); err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Failed to delete variant", err)
		return
	}

	pkghttp.Success(c, http.StatusOK, gin.H{"message": "Variant deleted successfully"})
}

func (h *Handler) BulkCreateProductVariants(c *gin.Context) {
	var req struct {
		ProductID uuid.UUID             `json:"product_id" binding:"required"`
		Variants  []ProductVariantRequest `json:"variants" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	variants, err := h.service.BulkCreateProductVariants(req.ProductID, req.Variants)
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Failed to create product variants", err)
		return
	}

	pkghttp.Success(c, http.StatusCreated, variants)
}

// Product Collection Handlers

func (h *Handler) CreateProductCollection(c *gin.Context) {
	var req ProductCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	collection, err := h.service.CreateProductCollection(&req)
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Failed to create collection", err)
		return
	}

	pkghttp.Success(c, http.StatusCreated, collection)
}

func (h *Handler) GetProductCollectionByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid collection ID", err)
		return
	}

	collection, err := h.service.GetProductCollectionByID(id)
	if err != nil {
		pkghttp.Error(c, http.StatusNotFound, "Collection not found", err)
		return
	}

	pkghttp.Success(c, http.StatusOK, collection)
}

func (h *Handler) GetProductCollections(c *gin.Context) {
	// Convert query parameters to map for filtering
	params := make(map[string]interface{})
	if isActive := c.Query("is_active"); isActive != "" {
		params["is_active"] = isActive == "true"
	}

	collections, err := h.service.GetProductCollections(&params)
	if err != nil {
		pkghttp.Error(c, http.StatusInternalServerError, "Failed to get collections", err)
		return
	}

	pkghttp.Success(c, http.StatusOK, collections)
}

func (h *Handler) UpdateProductCollection(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid collection ID", err)
		return
	}

	var req ProductCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	collection, err := h.service.UpdateProductCollection(id, &req)
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Failed to update collection", err)
		return
	}

	pkghttp.Success(c, http.StatusOK, collection)
}

func (h *Handler) DeleteProductCollection(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid collection ID", err)
		return
	}

	if err := h.service.DeleteProductCollection(id); err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Failed to delete collection", err)
		return
	}

	pkghttp.Success(c, http.StatusOK, gin.H{"message": "Collection deleted successfully"})
}

func (h *Handler) AddProductsToCollection(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid collection ID", err)
		return
	}

	var req struct {
		ProductIDs []uuid.UUID `json:"product_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := h.service.AddProductsToCollection(id, req.ProductIDs); err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Failed to add products to collection", err)
		return
	}

	pkghttp.Success(c, http.StatusOK, gin.H{"message": "Products added to collection successfully"})
}

func (h *Handler) RemoveProductsFromCollection(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid collection ID", err)
		return
	}

	var req struct {
		ProductIDs []uuid.UUID `json:"product_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := h.service.RemoveProductsFromCollection(id, req.ProductIDs); err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Failed to remove products from collection", err)
		return
	}

	pkghttp.Success(c, http.StatusOK, gin.H{"message": "Products removed from collection successfully"})
}

// Inventory Handlers

func (h *Handler) GetInventoryByProduct(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid product ID", err)
		return
	}

	inventory, err := h.service.GetInventoryByProduct(productID)
	if err != nil {
		pkghttp.Error(c, http.StatusInternalServerError, "Failed to get inventory", err)
		return
	}

	pkghttp.Success(c, http.StatusOK, inventory)
}

func (h *Handler) GetInventoryByVariant(c *gin.Context) {
	variantID, err := uuid.Parse(c.Param("variantId"))
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid variant ID", err)
		return
	}

	inventory, err := h.service.GetInventoryByVariant(variantID)
	if err != nil {
		pkghttp.Error(c, http.StatusInternalServerError, "Failed to get inventory", err)
		return
	}

	pkghttp.Success(c, http.StatusOK, inventory)
}

func (h *Handler) UpdateInventory(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid product ID", err)
		return
	}

	var req InventoryStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	inventory, err := h.service.UpdateInventory(productID, &req)
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Failed to update inventory", err)
		return
	}

	pkghttp.Success(c, http.StatusOK, inventory)
}

func (h *Handler) GetStockMovements(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid product ID", err)
		return
	}

	limit := 50 // default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	movements, err := h.service.GetStockMovements(productID, limit)
	if err != nil {
		pkghttp.Error(c, http.StatusInternalServerError, "Failed to get stock movements", err)
		return
	}

	pkghttp.Success(c, http.StatusOK, movements)
}

func (h *Handler) CreateStockMovement(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid product ID", err)
		return
	}

	var req StockMovementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	movement, err := h.service.CreateStockMovement(productID, &req)
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Failed to create stock movement", err)
		return
	}

	pkghttp.Success(c, http.StatusCreated, movement)
}

func (h *Handler) GetLowStockProducts(c *gin.Context) {
	inventory, err := h.service.GetLowStockProducts()
	if err != nil {
		pkghttp.Error(c, http.StatusInternalServerError, "Failed to get low stock products", err)
		return
	}

	pkghttp.Success(c, http.StatusOK, inventory)
}

func (h *Handler) GetOutOfStockProducts(c *gin.Context) {
	inventory, err := h.service.GetOutOfStockProducts()
	if err != nil {
		pkghttp.Error(c, http.StatusInternalServerError, "Failed to get out of stock products", err)
		return
	}

	pkghttp.Success(c, http.StatusOK, inventory)
}

// Statistics Handlers

func (h *Handler) GetCatalogStats(c *gin.Context) {
	stats, err := h.service.GetCatalogStats()
	if err != nil {
		pkghttp.Error(c, http.StatusInternalServerError, "Failed to get catalog statistics", err)
		return
	}

	pkghttp.Success(c, http.StatusOK, stats)
}

func (h *Handler) GetCategoryStats(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid category ID", err)
		return
	}

	stats, err := h.service.GetCategoryStats(id)
	if err != nil {
		pkghttp.Error(c, http.StatusInternalServerError, "Failed to get category statistics", err)
		return
	}

	pkghttp.Success(c, http.StatusOK, stats)
}

// Search and Discovery Handlers

func (h *Handler) SearchProducts(c *gin.Context) {
	query := c.Query("q")
	
	// Build filters from query parameters
	filters := make(map[string]interface{})
	
	if categoryID := c.Query("category_id"); categoryID != "" {
		if id, err := uuid.Parse(categoryID); err == nil {
			filters["category_id"] = id
		}
	}
	
	if sellerID := c.Query("seller_id"); sellerID != "" {
		if id, err := uuid.Parse(sellerID); err == nil {
			filters["seller_id"] = id
		}
	}
	
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	
	if featured := c.Query("featured"); featured != "" {
		filters["featured"] = featured == "true"
	}
	
	if condition := c.Query("condition"); condition != "" {
		filters["condition"] = condition
	}
	
	if minPrice := c.Query("min_price"); minPrice != "" {
		if price, err := strconv.ParseFloat(minPrice, 64); err == nil {
			filters["min_price"] = price
		}
	}
	
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if price, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			filters["max_price"] = price
		}
	}
	
	if sortBy := c.Query("sort_by"); sortBy != "" {
		filters["sort_by"] = sortBy
	}
	
	if sortDirection := c.Query("sort_direction"); sortDirection != "" {
		filters["sort_direction"] = sortDirection
	}
	
	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			filters["page"] = p
		}
	}
	
	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			filters["limit"] = l
		}
	}

	products, total, err := h.service.SearchProducts(query, filters)
	if err != nil {
		pkghttp.Error(c, http.StatusInternalServerError, "Failed to search products", err)
		return
	}

	pkghttp.SuccessWithPagination(c, http.StatusOK, products, total, len(products))
}

func (h *Handler) GetFeaturedProducts(c *gin.Context) {
	limit := 10 // default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	products, err := h.service.GetFeaturedProducts(limit)
	if err != nil {
		pkghttp.Error(c, http.StatusInternalServerError, "Failed to get featured products", err)
		return
	}

	pkghttp.Success(c, http.StatusOK, products)
}

func (h *Handler) GetRelatedProducts(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkghttp.Error(c, http.StatusBadRequest, "Invalid product ID", err)
		return
	}

	limit := 10 // default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	products, err := h.service.GetRelatedProducts(id, limit)
	if err != nil {
		pkghttp.Error(c, http.StatusInternalServerError, "Failed to get related products", err)
		return
	}

	pkghttp.Success(c, http.StatusOK, products)
}