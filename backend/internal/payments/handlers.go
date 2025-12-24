package payments

import (
	"net/http"
	"strconv"

	"github.com/blytz.live.remake/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Handler provides payment HTTP handlers
type Handler struct {
	service *Service
}

// NewHandler creates a new payment handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetDB provides access to GORM DB instance for internal use
func (h *Handler) GetDB() *gorm.DB {
	return h.service.db
}

// CreatePaymentIntentRequest represents request body for creating payment intent
type CreatePaymentIntentRequest struct {
	Amount   float64            `json:"amount" binding:"required,gt=0"`
	Currency string             `json:"currency" binding:"required"`
	OrderID  *uuid.UUID         `json:"order_id,omitempty"`
	Metadata map[string]string  `json:"metadata,omitempty"`
}

// ConfirmPaymentRequest represents request body for confirming payment
type ConfirmPaymentRequest struct {
	PaymentIntentID uuid.UUID `json:"payment_intent_id" binding:"required"`
}

// RefundPaymentRequest represents request body for refunding payment
type RefundPaymentRequest struct {
	PaymentID uuid.UUID `json:"payment_id" binding:"required"`
	Amount    float64   `json:"amount" binding:"required,gt=0"`
	Reason    string    `json:"reason" binding:"required,oneof=duplicate fraudulent requested_by_customer"`
	Notes     *string   `json:"notes,omitempty"`
}

// SavePaymentMethodRequest represents request body for saving payment method
type SavePaymentMethodRequest struct {
	PaymentMethodID string `json:"payment_method_id" binding:"required"`
	IsDefault      bool   `json:"is_default"`
}

// CreatePaymentIntent creates a new payment intent
func (h *Handler) CreatePaymentIntent(c *gin.Context) {
	var req CreatePaymentIntentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Set default currency if not provided
	if req.Currency == "" {
		req.Currency = "USD"
	}

	// Create payment intent
	paymentIntent, err := h.service.CreatePaymentIntent(
		c.Request.Context(),
		userID.(uuid.UUID),
		req.Amount,
		req.Currency,
		req.Metadata,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, paymentIntent)
}

// ConfirmPayment confirms and processes a payment
func (h *Handler) ConfirmPayment(c *gin.Context) {
	var req ConfirmPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.service.ConfirmPayment(c.Request.Context(), req.PaymentIntentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// GetPaymentIntent gets a payment intent by ID
func (h *Handler) GetPaymentIntent(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment intent ID"})
		return
	}

	paymentIntent, err := h.service.GetPaymentIntent(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment intent not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, paymentIntent)
}

// GetUserPaymentMethods gets all payment methods for the current user
func (h *Handler) GetUserPaymentMethods(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	paymentMethods, err := h.service.GetUserPaymentMethods(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payment_methods": paymentMethods})
}

// SavePaymentMethod saves a payment method for the current user
func (h *Handler) SavePaymentMethod(c *gin.Context) {
	var req SavePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	paymentMethod, err := h.service.SavePaymentMethod(
		c.Request.Context(),
		userID.(uuid.UUID),
		req.PaymentMethodID,
		req.IsDefault,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, paymentMethod)
}

// GetPaymentMethods gets available payment methods
func (h *Handler) GetPaymentMethods(c *gin.Context) {
	paymentMethods, err := h.service.GetPaymentMethods(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payment_methods": paymentMethods})
}

// RefundPayment processes a refund (admin only)
func (h *Handler) RefundPayment(c *gin.Context) {
	var req RefundPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (must be admin)
	processedBy, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check if user is admin
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	refund, err := h.service.RefundPayment(
		c.Request.Context(),
		req.PaymentID,
		req.Amount,
		req.Reason,
		processedBy.(uuid.UUID),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Add notes if provided
	if req.Notes != nil {
		refund.Notes = req.Notes
		h.service.db.WithContext(c.Request.Context()).Save(refund)
	}

	c.JSON(http.StatusCreated, refund)
}

// ProcessWebhook processes Stripe webhooks
func (h *Handler) ProcessWebhook(c *gin.Context) {
	// Get webhook endpoint secret from environment (temporarily hardcode for demo)
	endpointSecret := "whsec_your_webhook_secret" // This should come from config
	
	// In development, skip webhook signature verification
	if gin.Mode() == gin.DebugMode {
		c.JSON(http.StatusOK, gin.H{"status": "received (debug mode - verification skipped)"})
		return
	}
	
	if endpointSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Webhook endpoint secret not configured"})
		return
	}

	// Read the request body
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Get Stripe signature header
	signatureHeader := c.GetHeader("Stripe-Signature")
	if signatureHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stripe signature header missing"})
		return
	}

	// Process webhook
	err = h.service.ProcessWebhook(c.Request.Context(), body, signatureHeader, endpointSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

// ListPayments gets paginated list of payments (admin only)
func (h *Handler) ListPayments(c *gin.Context) {
	// Check if user is admin
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Build query
	db := h.service.db.WithContext(c.Request.Context()).Model(&models.Payment{}).Preload("User").Preload("Order")

	if status != "" {
		db = db.Where("status = ?", status)
	}

	// Count total records
	var total int64
	if err := db.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get paginated results
	var payments []models.Payment
	offset := (page - 1) * limit
	err := db.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&payments).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payments": payments,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

// GetPayment gets a payment by ID (admin only)
func (h *Handler) GetPayment(c *gin.Context) {
	// Check if user is admin
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	var payment models.Payment
	err = h.service.db.WithContext(c.Request.Context()).
		Preload("User").
		Preload("Order").
		First(&payment, "id = ?", id).Error

	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}