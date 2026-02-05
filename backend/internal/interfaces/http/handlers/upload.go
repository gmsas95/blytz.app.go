package handlers

import (
	"net/http"

	"github.com/blytz/live/backend/internal/application/upload"
	"github.com/gin-gonic/gin"
)

// UploadHandler handles file upload requests
type UploadHandler struct {
	service *upload.Service
}

// NewUploadHandler creates a new upload handler
func NewUploadHandler(service *upload.Service) *UploadHandler {
	return &UploadHandler{service: service}
}

// UploadProductImage handles product image upload
func (h *UploadHandler) UploadProductImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "no file provided"})
		return
	}
	defer file.Close()

	result, err := h.service.UploadProductImage(c.Request.Context(), file, header.Filename, header.Size)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: result})
}

// UploadAvatar handles avatar upload
func (h *UploadHandler) UploadAvatar(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "no file provided"})
		return
	}
	defer file.Close()

	result, err := h.service.UploadAvatar(c.Request.Context(), file, header.Filename, header.Size)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: result})
}

// UploadStreamThumbnail handles stream thumbnail upload
func (h *UploadHandler) UploadStreamThumbnail(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "no file provided"})
		return
	}
	defer file.Close()

	result, err := h.service.UploadStreamThumbnail(c.Request.Context(), file, header.Filename, header.Size)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: result})
}

// UploadGeneric handles generic file upload
func (h *UploadHandler) UploadGeneric(c *gin.Context) {
	folder := c.Param("folder")
	if folder == "" {
		folder = "uploads"
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "no file provided"})
		return
	}
	defer file.Close()

	result, err := h.service.UploadFromMultipart(c.Request.Context(), header, folder)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: result})
}

// DeleteFile handles file deletion
func (h *UploadHandler) DeleteFile(c *gin.Context) {
	var req struct {
		Key string `json:"key" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.service.DeleteFile(c.Request.Context(), req.Key); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: map[string]string{"message": "file deleted"}})
}
