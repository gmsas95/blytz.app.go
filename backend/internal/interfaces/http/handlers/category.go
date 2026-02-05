package handlers

import (
	"net/http"

	"github.com/blytz/live/backend/internal/application/category"
	categoryDomain "github.com/blytz/live/backend/internal/domain/category"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CategoryHandler handles category HTTP requests
type CategoryHandler struct {
	service *category.Service
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(service *category.Service) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// CreateCategoryRequest represents create category request
type CreateCategoryRequest struct {
	Name        string     `json:"name" binding:"required"`
	Description *string    `json:"description"`
	ImageURL    *string    `json:"image_url"`
	ParentID    *uuid.UUID `json:"parent_id"`
	SortOrder   int        `json:"sort_order"`
}

// UpdateCategoryRequest represents update category request
type UpdateCategoryRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	ImageURL    *string `json:"image_url"`
	SortOrder   *int    `json:"sort_order"`
	IsActive    *bool   `json:"is_active"`
}

// CategoryResponse represents category response
type CategoryResponse struct {
	ID           uuid.UUID          `json:"id"`
	Name         string             `json:"name"`
	Slug         string             `json:"slug"`
	Description  *string            `json:"description,omitempty"`
	ImageURL     *string            `json:"image_url,omitempty"`
	ParentID     *uuid.UUID         `json:"parent_id,omitempty"`
	Children     []CategoryResponse `json:"children,omitempty"`
	SortOrder    int                `json:"sort_order"`
	IsActive     bool               `json:"is_active"`
	ProductCount int                `json:"product_count"`
	CreatedAt    string             `json:"created_at"`
	UpdatedAt    string             `json:"updated_at"`
}

// Create creates a new category (admin only)
func (h *CategoryHandler) Create(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	
	dto := category.CreateCategoryDTO{
		Name:        req.Name,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		ParentID:    req.ParentID,
		SortOrder:   req.SortOrder,
	}
	
	cat, err := h.service.CreateCategory(c.Request.Context(), dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, SuccessResponse{Data: toCategoryResponse(cat)})
}

// Get retrieves a category by ID
func (h *CategoryHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid category id"})
		return
	}
	
	cat, err := h.service.GetCategory(c.Request.Context(), id)
	if err != nil {
		if err == categoryDomain.ErrCategoryNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, SuccessResponse{Data: toCategoryResponse(cat)})
}

// GetBySlug retrieves a category by slug
func (h *CategoryHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	
	cat, err := h.service.GetCategoryBySlug(c.Request.Context(), slug)
	if err != nil {
		if err == categoryDomain.ErrCategoryNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, SuccessResponse{Data: toCategoryResponse(cat)})
}

// GetTree retrieves the full category tree
func (h *CategoryHandler) GetTree(c *gin.Context) {
	tree, err := h.service.GetCategoryTree(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	
	response := make([]CategoryResponse, len(tree))
	for i, cat := range tree {
		response[i] = toCategoryResponse(cat)
	}
	
	c.JSON(http.StatusOK, SuccessResponse{Data: response})
}

// List retrieves all categories
func (h *CategoryHandler) List(c *gin.Context) {
	onlyActive := c.DefaultQuery("active_only", "true") == "true"
	
	cats, err := h.service.ListCategories(c.Request.Context(), onlyActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	
	response := make([]CategoryResponse, len(cats))
	for i, cat := range cats {
		response[i] = toCategoryResponse(cat)
	}
	
	c.JSON(http.StatusOK, SuccessResponse{Data: response})
}

// Update updates a category (admin only)
func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid category id"})
		return
	}
	
	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	
	dto := category.UpdateCategoryDTO{
		Name:        req.Name,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		SortOrder:   req.SortOrder,
		IsActive:    req.IsActive,
	}
	
	cat, err := h.service.UpdateCategory(c.Request.Context(), id, dto)
	if err != nil {
		switch err {
		case categoryDomain.ErrCategoryNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "category not found"})
		default:
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		}
		return
	}
	
	c.JSON(http.StatusOK, SuccessResponse{Data: toCategoryResponse(cat)})
}

// Delete deletes a category (admin only)
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid category id"})
		return
	}
	
	err = h.service.DeleteCategory(c.Request.Context(), id)
	if err != nil {
		switch err {
		case categoryDomain.ErrCategoryNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "category not found"})
		case categoryDomain.ErrHasProducts:
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "category has products"})
		case categoryDomain.ErrHasSubcategories:
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "category has subcategories"})
		default:
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		}
		return
	}
	
	c.JSON(http.StatusOK, SuccessResponse{Data: map[string]string{"message": "category deleted"}})
}

// Helper function
func toCategoryResponse(cat *categoryDomain.Category) CategoryResponse {
	resp := CategoryResponse{
		ID:           cat.ID,
		Name:         cat.Name,
		Slug:         cat.Slug,
		Description:  cat.Description,
		ImageURL:     cat.ImageURL,
		ParentID:     cat.ParentID,
		SortOrder:    cat.SortOrder,
		IsActive:     cat.IsActive,
		ProductCount: cat.ProductCount,
		CreatedAt:    cat.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    cat.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		Children:     make([]CategoryResponse, len(cat.Children)),
	}
	
	for i, child := range cat.Children {
		resp.Children[i] = toCategoryResponse(child)
	}
	
	return resp
}
