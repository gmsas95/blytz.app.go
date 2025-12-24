package payments

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/blytz.live.remake/backend/internal/models"
	"github.com/blytz.live.remake/backend/pkg/logging"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/refund"
	"github.com/stripe/stripe-go/v76/webhook"
)

// Service provides payment processing services
type Service struct {
	db     *gorm.DB
	logger *logging.Logger
	apiKey string
}

// NewService creates a new payment service
func NewService(db *gorm.DB, apiKey string) *Service {
	logger := logging.NewLogger()
	
	// Initialize Stripe
	stripe.Key = apiKey
	
	return &Service{
		db:     db,
		logger: logger,
		apiKey: apiKey,
	}
}

// CreatePaymentIntent creates a new payment intent
func (s *Service) CreatePaymentIntent(ctx context.Context, userID uuid.UUID, amount float64, currency string, metadata map[string]string) (*models.PaymentIntent, error) {
	// Convert amount to cents (Stripe uses smallest currency unit)
	amountCents := int64(amount * 100)
	
	// Create payment intent with Stripe
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amountCents),
		Currency: stripe.String(currency),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}
	
	// Set expiration to 30 minutes
	expiresAt := time.Now().Add(30 * time.Minute)
	
	// Use metadata for expiration
	if params.Metadata == nil {
		params.Metadata = make(map[string]string)
	}
	for k, v := range metadata {
		params.Metadata[k] = v
	}
	params.Metadata["expires_at"] = expiresAt.Format(time.RFC3339)
	
	pi, err := paymentintent.New(params)
	if err != nil {
		s.logger.Error("Failed to create Stripe payment intent", map[string]interface{}{
			"error":  err.Error(),
			"amount": amount,
			"user_id": userID,
		})
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}
	
	// Save to database
	paymentIntent := &models.PaymentIntent{
		UserID:       userID,
		Amount:       amount,
		Currency:     currency,
		Status:       string(pi.Status),
		ClientSecret: pi.ClientSecret,
		GatewayRef:   pi.ID,
		GatewayType:  "stripe",
		PaymentMethods: []string{"card"},
		ExpiresAt:    expiresAt,
		Metadata:     s.mapToJSON(metadata),
	}
	
	if err := s.db.WithContext(ctx).Create(paymentIntent).Error; err != nil {
		return nil, fmt.Errorf("failed to save payment intent: %w", err)
	}
	
	s.logger.Info("Payment intent created", map[string]interface{}{
		"payment_intent_id": paymentIntent.ID,
		"stripe_id":        pi.ID,
		"amount":           amount,
		"user_id":          userID,
	})
	
	return paymentIntent, nil
}

// ConfirmPayment confirms and processes a payment
func (s *Service) ConfirmPayment(ctx context.Context, paymentIntentID uuid.UUID) (*models.Payment, error) {
	// Get payment intent from database
	var paymentIntent models.PaymentIntent
	if err := s.db.WithContext(ctx).First(&paymentIntent, "id = ?", paymentIntentID).Error; err != nil {
		return nil, fmt.Errorf("payment intent not found: %w", err)
	}
	
	// Get status from Stripe
	pi, err := paymentintent.Get(paymentIntent.GatewayRef, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve payment intent from Stripe: %w", err)
	}
	
	// Update payment intent status
	paymentIntent.Status = string(pi.Status)
	s.db.WithContext(ctx).Save(&paymentIntent)
	
	if pi.Status != stripe.PaymentIntentStatusSucceeded {
		return nil, fmt.Errorf("payment not successful: %s", pi.Status)
	}
	
	// Create payment record
	payment := &models.Payment{
		UserID:        paymentIntent.UserID,
		PaymentMethod: "stripe",
		Amount:        paymentIntent.Amount,
		Currency:      paymentIntent.Currency,
		Status:        "completed",
		TransactionID: pi.ID,
		GatewayRef:    pi.ID,
		GatewayType:   "stripe",
		ProcessedAt:   &[]time.Time{time.Now()}[0],
		Metadata:      paymentIntent.Metadata,
	}
	
	if paymentIntent.OrderID != nil {
		payment.OrderID = *paymentIntent.OrderID
	}
	
	if paymentIntent.OrderID != nil {
		payment.OrderID = *paymentIntent.OrderID
	}
	
	if err := s.db.WithContext(ctx).Create(payment).Error; err != nil {
		return nil, fmt.Errorf("failed to create payment record: %w", err)
	}
	
	s.logger.Info("Payment confirmed", map[string]interface{}{
		"payment_id":        payment.ID,
		"payment_intent_id": paymentIntentID,
		"stripe_id":         pi.ID,
		"amount":           payment.Amount,
	})
	
	return payment, nil
}

// RefundPayment processes a refund
func (s *Service) RefundPayment(ctx context.Context, paymentID uuid.UUID, amount float64, reason string, processedBy uuid.UUID) (*models.Refund, error) {
	// Get payment from database
	var payment models.Payment
	if err := s.db.WithContext(ctx).First(&payment, "id = ?", paymentID).Error; err != nil {
		return nil, fmt.Errorf("payment not found: %w", err)
	}
	
	if payment.Status != "completed" {
		return nil, fmt.Errorf("cannot refund payment with status: %s", payment.Status)
	}
	
	// Convert amount to cents
	amountCents := int64(amount * 100)
	
	// Create refund with Stripe
	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(payment.GatewayRef),
		Amount:       stripe.Int64(amountCents),
		Reason:       stripe.String(reason),
		Metadata: map[string]string{
			"refund_id": uuid.New().String(),
		},
	}
	
	refundObj, err := refund.New(params)
	if err != nil {
		s.logger.Error("Failed to create Stripe refund", map[string]interface{}{
			"error":      err.Error(),
			"payment_id": paymentID,
			"amount":     amount,
		})
		return nil, fmt.Errorf("failed to process refund: %w", err)
	}
	
	// Create refund record
	refundRecord := &models.Refund{
		PaymentID:   paymentID,
		Amount:      amount,
		Reason:      reason,
		Status:      string(refundObj.Status),
		GatewayRef:  refundObj.ID,
		GatewayType: "stripe",
		ProcessedBy: processedBy,
		ProcessedAt: &[]time.Time{time.Now()}[0],
		Metadata:    s.mapToJSON(refundObj.Metadata),
	}
	
	if err := s.db.WithContext(ctx).Create(refundRecord).Error; err != nil {
		return nil, fmt.Errorf("failed to save refund record: %w", err)
	}
	
	// Update payment refund amount
	payment.RefundedAmount += amount
	s.db.WithContext(ctx).Save(&payment)
	
	s.logger.Info("Refund processed", map[string]interface{}{
		"refund_id":  refundRecord.ID,
		"payment_id": paymentID,
		"stripe_id":  refundObj.ID,
		"amount":     amount,
	})
	
	return refundRecord, nil
}

// SavePaymentMethod saves a payment method for a user
func (s *Service) SavePaymentMethod(ctx context.Context, userID uuid.UUID, paymentMethodID string, isDefault bool) (*models.PaymentMethod, error) {
	// Get payment method details from Stripe
	// Note: In a real implementation, you would retrieve payment method details from Stripe API
	// For now, we'll create a basic record
	
	paymentMethod := &models.PaymentMethod{
		UserID:    userID,
		Type:      "card",
		Provider:  "stripe",
		MethodRef: paymentMethodID,
		IsDefault: isDefault,
		IsVerified: true,
	}
	
	// If setting as default, unset other defaults
	if isDefault {
		s.db.WithContext(ctx).Model(&models.PaymentMethod{}).
			Where("user_id = ? AND is_default = ?", userID, true).
			Update("is_default", false)
	}
	
	if err := s.db.WithContext(ctx).Create(paymentMethod).Error; err != nil {
		return nil, fmt.Errorf("failed to save payment method: %w", err)
	}
	
	s.logger.Info("Payment method saved", map[string]interface{}{
		"payment_method_id": paymentMethod.ID,
		"user_id":          userID,
		"stripe_id":        paymentMethodID,
	})
	
	return paymentMethod, nil
}

// GetUserPaymentMethods gets all payment methods for a user
func (s *Service) GetUserPaymentMethods(ctx context.Context, userID uuid.UUID) ([]models.PaymentMethod, error) {
	var paymentMethods []models.PaymentMethod
	err := s.db.WithContext(ctx).
		Where("user_id = ? AND is_verified = ?", userID, true).
		Order("is_default DESC, created_at DESC").
		Find(&paymentMethods).Error
	
	return paymentMethods, err
}

// ProcessWebhook processes Stripe webhook events
func (s *Service) ProcessWebhook(ctx context.Context, payload []byte, signatureHeader string, endpointSecret string) error {
	event, err := webhook.ConstructEvent(payload, signatureHeader, endpointSecret)
	if err != nil {
		return fmt.Errorf("webhook signature verification failed: %w", err)
	}
	
	switch event.Type {
	case "payment_intent.succeeded":
		return s.handlePaymentSucceeded(ctx, event)
	case "payment_intent.payment_failed":
		return s.handlePaymentFailed(ctx, event)
	case "payment_intent.canceled":
		return s.handlePaymentCanceled(ctx, event)
	case "charge.dispute.created":
		return s.handleDisputeCreated(ctx, event)
	default:
		s.logger.Info("Unhandled webhook event", map[string]interface{}{
			"type": event.Type,
		})
	}
	
	return nil
}

// handlePaymentSucceeded handles payment_intent.succeeded webhook
func (s *Service) handlePaymentSucceeded(ctx context.Context, event stripe.Event) error {
	var paymentIntent stripe.PaymentIntent
	if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
		return fmt.Errorf("failed to parse payment intent: %w", err)
	}
	
	// Update payment intent in database
	var dbPaymentIntent models.PaymentIntent
	err := s.db.WithContext(ctx).
		Where("gateway_ref = ?", paymentIntent.ID).
		First(&dbPaymentIntent).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Payment intent not found in our database, might be from external source
			s.logger.Warn("Payment intent not found in database", map[string]interface{}{
				"stripe_id": paymentIntent.ID,
			})
			return nil
		}
		return fmt.Errorf("failed to find payment intent: %w", err)
	}
	
	// Update status
	dbPaymentIntent.Status = string(paymentIntent.Status)
	s.db.WithContext(ctx).Save(&dbPaymentIntent)
	
	// Create payment record if status is succeeded
	if paymentIntent.Status == stripe.PaymentIntentStatusSucceeded {
		payment := &models.Payment{
			UserID:        dbPaymentIntent.UserID,
			PaymentMethod: "stripe",
			Amount:        dbPaymentIntent.Amount,
			Currency:      dbPaymentIntent.Currency,
			Status:        "completed",
			TransactionID: paymentIntent.ID,
			GatewayRef:    paymentIntent.ID,
			GatewayType:   "stripe",
			ProcessedAt:   &[]time.Time{time.Now()}[0],
			Metadata:      dbPaymentIntent.Metadata,
		}
		
		if dbPaymentIntent.OrderID != nil {
			payment.OrderID = *dbPaymentIntent.OrderID
		}
		
		if err := s.db.WithContext(ctx).Create(payment).Error; err != nil {
			return fmt.Errorf("failed to create payment record: %w", err)
		}
		
		s.logger.Info("Payment completed via webhook", map[string]interface{}{
			"payment_id":        payment.ID,
			"payment_intent_id": dbPaymentIntent.ID,
			"stripe_id":         paymentIntent.ID,
		})
	}
	
	return nil
}

// handlePaymentFailed handles payment_intent.payment_failed webhook
func (s *Service) handlePaymentFailed(ctx context.Context, event stripe.Event) error {
	var paymentIntent stripe.PaymentIntent
	if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
		return fmt.Errorf("failed to parse payment intent: %w", err)
	}
	
	// Update payment intent in database
	return s.db.WithContext(ctx).
		Model(&models.PaymentIntent{}).
		Where("gateway_ref = ?", paymentIntent.ID).
		Updates(map[string]interface{}{
			"status": string(paymentIntent.Status),
		}).Error
}

// handlePaymentCanceled handles payment_intent.canceled webhook
func (s *Service) handlePaymentCanceled(ctx context.Context, event stripe.Event) error {
	var paymentIntent stripe.PaymentIntent
	if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
		return fmt.Errorf("failed to parse payment intent: %w", err)
	}
	
	// Update payment intent in database
	return s.db.WithContext(ctx).
		Model(&models.PaymentIntent{}).
		Where("gateway_ref = ?", paymentIntent.ID).
		Updates(map[string]interface{}{
			"status": string(paymentIntent.Status),
		}).Error
}

// handleDisputeCreated handles charge.dispute.created webhook
func (s *Service) handleDisputeCreated(ctx context.Context, event stripe.Event) error {
	// Log the dispute for review
	s.logger.Warn("Payment dispute created", map[string]interface{}{
		"event_id": event.ID,
		"type":     event.Type,
	})
	
	// In a real implementation, you would:
	// 1. Find the associated payment
	// 2. Update its status to 'disputed'
	// 3. Notify the merchant/admin
	// 4. Create a dispute record in the database
	
	return nil
}

// mapToJSON converts a map to JSON string
func (s *Service) mapToJSON(data map[string]string) string {
	if data == nil {
		return "{}"
	}
	
	jsonData, err := json.Marshal(data)
	if err != nil {
		s.logger.Error("Failed to marshal metadata", map[string]interface{}{
			"error": err.Error(),
			"data":  data,
		})
		return "{}"
	}
	
	return string(jsonData)
}

// GetPaymentIntent retrieves a payment intent by ID
func (s *Service) GetPaymentIntent(ctx context.Context, id uuid.UUID) (*models.PaymentIntent, error) {
	var paymentIntent models.PaymentIntent
	err := s.db.WithContext(ctx).First(&paymentIntent, "id = ?", id).Error
	return &paymentIntent, err
}

// GetPaymentMethods gets available payment methods
func (s *Service) GetPaymentMethods(ctx context.Context) ([]string, error) {
	// Return supported payment methods based on Stripe configuration
	return []string{
		"card",
		"apple_pay",
		"google_pay",
	}, nil
}