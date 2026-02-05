package handlers

import (
	"net/http"
	"strconv"

	"github.com/blytz/live/backend/internal/application/product"
	productDomain "github.com/blytz/live/backend/internal/domain/product"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ProductHandler handles product HTTP requests
type ProductHandler struct {
	service *product.Service
}

// NewProductHandler creates a new product handler
func NewProductHandler(service *product.Service) *ProductHandler {
	return &ProductHandler{service: service}
}

// CreateProductRequest represents create product request
type CreateProductRequest struct {
	CategoryID     *uuid.UUID            `json:"category_id"`
	Name           string                `json:"name" binding:"required"`
	Description    string                `json:"description" binding:"required"`
	Condition      productDomain.Condition `json:"condition" binding:"required"`
	BasePrice      float64               `json:"base_price" binding:"required,gt=0"`
	CompareAtPrice *float64              `json:"compare_at_price"`
	StockQuantity  int                   `json:"stock_quantity" binding:"gte=0"`
	SKU            *string               `json:"sku"`
	WeightGrams    *int                  `json:"weight_grams"`
	DimensionsCm   *productDomain.Dimensions `json:"dimensions_cm"`
	Attributes     map[string]string     `json:"attributes"`
	Images         []ImageRequest        `json:"images" binding:"required,min=1"`
}

// ImageRequest represents image in request
type ImageRequest struct {
	URL       string `json:"url" binding:"required,url"`
	AltText   string `json:"alt_text"`
	IsPrimary bool   `json:"is_primary"`
}

// UpdateProductRequest represents update product request
type UpdateProductRequest struct {
	CategoryID     *uuid.UUID                `json:"category_id"`
	Name           *string                   `json:"name"`
	Description    *string                   `json:"description"`
	Condition      *productDomain.Condition  `json:"condition"`
	BasePrice      *float64                  `json:"base_price"`
	CompareAtPrice *float64                  `json:"compare_at_price"`
	StockQuantity  *int                      `json:"stock_quantity"`
	SKU            *string                   `json:"sku"`
	WeightGrams    *int                      `json:"weight_grams"`
	DimensionsCm   *productDomain.Dimensions `json:"dimensions_cm"`
	Attributes     map[string]string         `json:"attributes"`
}

// ProductResponse represents product response
type ProductResponse struct {
	ID             uuid.UUID           `json:"id"`
	SellerID       uuid.UUID           `json:"seller_id"`
	CategoryID     *uuid.UUID          `json:"category_id,omitempty"`
	Name           string              `json:"name"`
	Slug           string              `json:"slug"`
	Description    string              `json:"description"`
	Condition      string              `json:"condition"`
	BasePrice      float64             `json:"base_price"`
	CompareAtPrice *float64            `json:"compare_at_price,omitempty"`
	StockQuantity  int                 `json:"stock_quantity"`
	SKU            *string             `json:"sku,omitempty"`
	WeightGrams    *int                `json:"weight_grams,omitempty"`
	DimensionsCm   *DimensionsResponse `json:"dimensions_cm,omitempty"`
	Attributes     map[string]string   `json:"attributes,omitempty"`
	Images         []ImageResponse     `json:"images"`
	Status         string              `json:"status"`
	ViewCount      int                 `json:"view_count"`
	CreatedAt      string              `json:"created_at"`
	UpdatedAt      string              `json:"updated_at"`
}

// DimensionsResponse represents dimensions response
type DimensionsResponse struct {
	Length int `json:"length"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// ImageResponse represents image response
type ImageResponse struct {
	ID        uuid.UUID `json:"id"`
	URL       string    `json:"url"`
	AltText   *string   `json:"alt_text,omitempty"`
	SortOrder int       `json:"sort_order"`
	IsPrimary bool      `json:"is_primary"`
}

// ListProductsResponse represents list products response
type ListProductsResponse struct {
	Products   []ProductResponse `json:"products"`
	TotalCount int64             `json:"total_count"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
}

// Create creates a new product
func (h *ProductHandler) Create(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	
	// Get seller ID from context (set by auth middleware)
	sellerID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}
	
	// Convert images
	images := make([]product.CreateImageDTO, len(req.Images))
	for i, img := range req.Images {
		images[i] = product.CreateImageDTO{
			URL:       img.URL,
			AltText:   &img.AltText,
			IsPrimary: img.IsPrimary,
		}
	}
	
	dto := product.CreateProductDTO{
		SellerID:       sellerID.(uuid.UUID),
		CategoryID:     req.CategoryID,
		Name:           req.Name,
		Description:    req.Description,
		Condition:      req.Condition,
		BasePrice:      req.BasePrice,
		CompareAtPrice: req.CompareAtPrice,
		StockQuantity:  req.StockQuantity,
		SKU:            req.SKU,
		WeightGrams:    req.WeightGrams,
		DimensionsCm:   req.DimensionsCm,
		Attributes:     req.Attributes,
		Images:         images,
	}
	
	p, err := h.service.CreateProduct(c.Request.Context(), dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, SuccessResponse{Data: toProductResponse(p)})
}

// Get retrieves a product by ID
func (h *ProductHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid product id"})
		return
	}
	
	p, err := h.service.GetProduct(c.Request.Context(), id)
	if err != nil {
		if err == productDomain.ErrProductNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, SuccessResponse{Data: toProductResponse(p)})
}

// GetBySlug retrieves a product by slug
func (h *ProductHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	
	p, err := h.service.GetProductBySlug(c.Request.Context(), slug)
	if err != nil {
		if err == productDomain.ErrProductNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, SuccessResponse{Data: toProductResponse(p)})
}

// List retrieves a list of products
func (h *ProductHandler) List(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	
	dto := product.ListProductsDTO{
		Query:    c.Query("q"),
		SortBy:   c.Query("sort"),
		Page:     page,
		PageSize: pageSize,
	}
	
	// Parse optional filters
	if categoryID := c.Query("category_id"); categoryID != "" {
		if id, err := uuid.Parse(categoryID); err == nil {
			dto.CategoryID = &id
		}
	}
	
	if sellerID := c.Query("seller_id"); sellerID != "" {
		if id, err := uuid.Parse(sellerID); err == nil {
			dto.SellerID = &id
		}
	}
	
	if status := c.Query("status"); status != "" {
		s := productDomain.Status(status)
		dto.Status = &s
	}
	
	if condition := c.Query("condition"); condition != "" {
		c := productDomain.Condition(condition)
		dto.Condition = &c
	}
	
	if minPrice := c.Query("min_price"); minPrice != "" {
		if price, err := strconv.ParseFloat(minPrice, 64); err == nil {
			dto.MinPrice = &price
		}
	}
	
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if price, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			dto.MaxPrice = &price
		}
	}
	
	result, err := h.service.ListProducts(c.Request.Context(), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	
	response := ListProductsResponse{
		Products:   make([]ProductResponse, len(result.Products)),
		TotalCount: result.TotalCount,
		Page:       result.Page,
		PageSize:   result.PageSize,
	}
	
	for i, p := range result.Products {
		response.Products[i] = toProductResponse(p)
	}
	
	c.JSON(http.StatusOK, SuccessResponse{Data: response})
}

// Update updates a product
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid product id"})
		return
	}
	
	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	
	sellerID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}
	
	dto := product.UpdateProductDTO{
		CategoryID:     req.CategoryID,
		Name:           req.Name,
		Description:    req.Description,
		Condition:      req.Condition,
		BasePrice:      req.BasePrice,
		CompareAtPrice: req.CompareAtPrice,
		StockQuantity:  req.StockQuantity,
		SKU:            req.SKU,
		WeightGrams:    req.WeightGrams,
		DimensionsCm:   req.DimensionsCm,
		Attributes:     req.Attributes,
	}
	
	p, err := h.service.UpdateProduct(c.Request.Context(), id, sellerID.(uuid.UUID), dto)
	if err != nil {
		switch err {
		case productDomain.ErrProductNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "product not found"})
		case productDomain.ErrUnauthorized:
			c.JSON(http.StatusForbidden, ErrorResponse{Error: "forbidden"})
		default:
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		}
		return
	}
	
	c.JSON(http.StatusOK, SuccessResponse{Data: toProductResponse(p)})
}

// Delete deletes a product
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid product id"})
		return
	}
	
	sellerID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}
	
	err = h.service.DeleteProduct(c.Request.Context(), id, sellerID.(uuid.UUID))
	if err != nil {
		switch err {
		case productDomain.ErrProductNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "product not found"})
		case productDomain.ErrUnauthorized:
			c.JSON(http.StatusForbidden, ErrorResponse{Error: "forbidden"})
		default:
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		}
		return
	}
	
	c.JSON(http.StatusOK, SuccessResponse{Data: map[string]string{"message": "product deleted"}})
}

// Publish publishes a product
func (h *ProductHandler) Publish(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid product id"})
		return
	}
	
	sellerID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}
	
	err = h.service.PublishProduct(c.Request.Context(), id, sellerID.(uuid.UUID))
	if err != nil {
		switch err {
		case productDomain.ErrProductNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "product not found"})
		case productDomain.ErrUnauthorized:
			c.JSON(http.StatusForbidden, ErrorResponse{Error: "forbidden"})
		default:
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		}
		return
	}
	
	c.JSON(http.StatusOK, SuccessResponse{Data: map[string]string{"message": "product published"}})
}

// Archive archives a product
func (h *ProductHandler) Archive(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid product id"})
		return
	}
	
	sellerID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}
	
	err = h.service.ArchiveProduct(c.Request.Context(), id, sellerID.(uuid.UUID))
	if err != nil {
		switch err {
		case productDomain.ErrProductNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "product not found"})
		case productDomain.ErrUnauthorized:
			c.JSON(http.StatusForbidden, ErrorResponse{Error: "forbidden"})
		default:
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		}
		return
	}
	
	c.JSON(http.StatusOK, SuccessResponse{Data: map[string]string{"message": "product archived"}})
}

// GetMyProducts retrieves products for the authenticated seller
func (h *ProductHandler) GetMyProducts(c *gin.Context) {
	sellerID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	
	var status *productDomain.Status
	if s := c.Query("status"); s != "" {
		st := productDomain.Status(s)
		status = &st
	}
	
	result, err := h.service.GetSellerProducts(c.Request.Context(), sellerID.(uuid.UUID), status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	
	response := ListProductsResponse{
		Products:   make([]ProductResponse, len(result.Products)),
		TotalCount: result.TotalCount,
		Page:       result.Page,
		PageSize:   result.PageSize,
	}
	
	for i, p := range result.Products {
		response.Products[i] = toProductResponse(p)
	}
	
	c.JSON(http.StatusOK, SuccessResponse{Data: response})
}

// Helper functions

func toProductResponse(p *productDomain.Product) ProductResponse {
	resp := ProductResponse{
		ID:            p.ID,
		SellerID:      p.SellerID,
		CategoryID:    p.CategoryID,
		Name:          p.Name,
		Slug:          p.Slug,
		Description:   p.Description,
		Condition:     string(p.Condition),
		BasePrice:     p.BasePrice,
		StockQuantity: p.StockQuantity,
		SKU:           p.SKU,
		WeightGrams:   p.WeightGrams,
		Attributes:    p.Attributes,
		Images:        make([]ImageResponse, len(p.Images)),
		Status:        string(p.Status),
		ViewCount:     p.ViewCount,
		CreatedAt:     p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
	
	if p.CompareAtPrice != nil {
		resp.CompareAtPrice = p.CompareAtPrice
	}
	
	if p.DimensionsCm != nil {
		resp.DimensionsCm = &DimensionsResponse{
			Length: p.DimensionsCm.Length,
			Width:  p.DimensionsCm.Width,
			Height: p.DimensionsCm.Height,
		}
	}
	
	for i, img := range p.Images {
		resp.Images[i] = ImageResponse{
			ID:        img.ID,
			URL:       img.URL,
			AltText:   img.AltText,
			SortOrder: img.SortOrder,
			IsPrimary: img.IsPrimary,
		}
	}
	
	return resp
}

// SuccessResponse represents a successful response
type SuccessResponse struct {
	Data interface{} `json:"data"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}
