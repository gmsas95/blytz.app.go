package cart

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/blytz.live.remake/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service handles cart business logic
type Service struct {
	db *gorm.DB
}

// NewService creates a new cart service
func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

// GetOrCreateCart gets existing cart or creates new one
func (s *Service) GetOrCreateCart(userID *uuid.UUID, token *string) (*CartResponse, error) {
	var cart Cart

	// Try to find existing cart
	if userID != nil {
		// User cart
		err := s.db.Where("user_id = ? AND expires_at > ?", *userID, time.Now()).
			Preload("Items").
			First(&cart).Error
		if err == nil {
			return s.cartToResponse(&cart)
		}
	} else if token != nil {
		// Guest cart
		err := s.db.Where("token = ? AND expires_at > ?", *token, time.Now()).
			Preload("Items").
			First(&cart).Error
		if err == nil {
			return s.cartToResponse(&cart)
		}
	}

	// Create new cart
	newCart := Cart{
		ID:       uuid.New(),
		UserID:   userID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if userID == nil {
		// Generate guest token
		newCart.Token = s.generateGuestToken()
	}

	if err := s.db.Create(&newCart).Error; err != nil {
		return nil, fmt.Errorf("failed to create cart: %w", err)
	}

	return s.cartToResponse(&newCart)
}

// AddItem adds item to cart
func (s *Service) AddItem(cartID uuid.UUID, req AddItemRequest) (*CartResponse, error) {
	// Get cart
	var cart Cart
	if err := s.db.Preload("Items").First(&cart, cartID).Error; err != nil {
		return nil, fmt.Errorf("cart not found: %w", err)
	}

	// Check if cart is expired
	if cart.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("cart has expired")
	}

	// Validate product exists and is active
	var product models.Product
	if err := s.db.First(&product, req.ProductID).Error; err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	if product.Status != "active" {
		return nil, errors.New("product is not available for purchase")
	}

	// Check if item already exists in cart
	for i, item := range cart.Items {
		if item.ProductID == req.ProductID {
			// Update quantity
			newQuantity := item.Quantity + req.Quantity
			if newQuantity > 10 {
				return nil, errors.New("maximum quantity per item is 10")
			}

			cart.Items[i].Quantity = newQuantity
			break
		}
	}

	// If item doesn't exist, add new item
	itemExists := false
	for _, item := range cart.Items {
		if item.ProductID == req.ProductID {
			itemExists = true
			break
		}
	}

	if !itemExists {
		cartItem := CartItem{
			ID:        uuid.New(),
			CartID:    cartID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
			AddedAt:   time.Now(),
		}
		cart.Items = append(cart.Items, cartItem)
	}

	// Update cart
	cart.UpdatedAt = time.Now()
	if err := s.db.Save(&cart).Error; err != nil {
		return nil, fmt.Errorf("failed to update cart: %w", err)
	}

	return s.cartToResponse(&cart)
}

// UpdateItem updates cart item quantity
func (s *Service) UpdateItem(cartID, itemID uuid.UUID, req UpdateItemRequest) (*CartResponse, error) {
	// Get cart
	var cart Cart
	if err := s.db.Preload("Items").First(&cart, cartID).Error; err != nil {
		return nil, fmt.Errorf("cart not found: %w", err)
	}

	// Check if cart is expired
	if cart.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("cart has expired")
	}

	// Find and update item
	itemFound := false
	for i, item := range cart.Items {
		if item.ID == itemID {
			cart.Items[i].Quantity = req.Quantity
			itemFound = true
			break
		}
	}

	if !itemFound {
		return nil, errors.New("item not found in cart")
	}

	// Update cart
	cart.UpdatedAt = time.Now()
	if err := s.db.Save(&cart).Error; err != nil {
		return nil, fmt.Errorf("failed to update cart: %w", err)
	}

	return s.cartToResponse(&cart)
}

// RemoveItem removes item from cart
func (s *Service) RemoveItem(cartID, itemID uuid.UUID) (*CartResponse, error) {
	// Get cart
	var cart Cart
	if err := s.db.Preload("Items").First(&cart, cartID).Error; err != nil {
		return nil, fmt.Errorf("cart not found: %w", err)
	}

	// Check if cart is expired
	if cart.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("cart has expired")
	}

	// Remove item
	itemFound := false
	for i, item := range cart.Items {
		if item.ID == itemID {
			cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			itemFound = true
			break
		}
	}

	if !itemFound {
		return nil, errors.New("item not found in cart")
	}

	// Update cart
	cart.UpdatedAt = time.Now()
	if err := s.db.Save(&cart).Error; err != nil {
		return nil, fmt.Errorf("failed to update cart: %w", err)
	}

	return s.cartToResponse(&cart)
}

// ClearCart clears all items from cart
func (s *Service) ClearCart(cartID uuid.UUID) error {
	// Delete all cart items
	if err := s.db.Where("cart_id = ?", cartID).Delete(&CartItem{}).Error; err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}

	// Update cart timestamp
	if err := s.db.Model(&Cart{}).Where("id = ?", cartID).
		Update("updated_at", time.Now()).Error; err != nil {
		return fmt.Errorf("failed to update cart timestamp: %w", err)
	}

	return nil
}

// GetCartWithDetails gets cart with product details
func (s *Service) GetCartWithDetails(cartID uuid.UUID) (*CartResponse, error) {
	// Get cart with items and products
	var cart Cart
	if err := s.db.Preload("Items").First(&cart, cartID).Error; err != nil {
		return nil, fmt.Errorf("cart not found: %w", err)
	}

	// Check if cart is expired
	if cart.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("cart has expired")
	}

	// Get product details for each item
	itemResponses := make([]CartItemResponse, len(cart.Items))
	var subtotal float64
	var totalItems int

	for i, item := range cart.Items {
		var product models.Product
		if err := s.db.First(&product, item.ProductID).Error; err != nil {
			return nil, fmt.Errorf("product not found: %w", err)
		}

		// Parse images JSON
		var images []string
		if product.Images != nil {
			json.Unmarshal([]byte(*product.Images), &images)
		}

		productResponse := ProductResponse{
			ID:            product.ID,
			Title:         product.Title,
			Description:   product.Description,
			Condition:     product.Condition,
			StartingPrice: product.StartingPrice,
			ReservePrice:  product.ReservePrice,
			BuyNowPrice:   product.BuyNowPrice,
			Images:        images,
			Status:        product.Status,
		}

		lineTotal := product.StartingPrice * float64(item.Quantity)

		itemResponses[i] = CartItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			AddedAt:   item.AddedAt,
			Product:   productResponse,
			LineTotal: lineTotal,
		}

		subtotal += lineTotal
		totalItems += item.Quantity
	}

	cartResponse := CartResponse{
		ID:         cart.ID,
		UserID:     cart.UserID,
		Token:      cart.Token,
		ExpiresAt:  cart.ExpiresAt,
		Items:      itemResponses,
		ItemCount:  len(itemResponses),
		Subtotal:   subtotal,
		TotalItems: totalItems,
		CreatedAt:  cart.CreatedAt,
		UpdatedAt:  cart.UpdatedAt,
	}

	return &cartResponse, nil
}

// MergeGuestCart merges guest cart to user cart
func (s *Service) MergeGuestCart(guestToken string, userID uuid.UUID) (*CartResponse, error) {
	// Get guest cart
	var guestCart Cart
	if err := s.db.Preload("Items").Where("token = ? AND expires_at > ?", guestToken, time.Now()).
		First(&guestCart).Error; err != nil {
		return nil, fmt.Errorf("guest cart not found: %w", err)
	}

	// Get user cart
	var userCart Cart
	if err := s.db.Where("user_id = ?", userID).First(&userCart).Error; err != nil {
		return nil, fmt.Errorf("failed to get user cart: %w", err)
	}

	// Merge items
	for _, guestItem := range guestCart.Items {
		itemExists := false
		for _, userItem := range userCart.Items {
			if userItem.ProductID == guestItem.ProductID {
				// Update quantity (max 10 per item)
				newQuantity := userItem.Quantity + guestItem.Quantity
				if newQuantity > 10 {
					newQuantity = 10
				}

				// Update in database
				if err := s.db.Model(&CartItem{}).
					Where("id = ?", userItem.ID).
					Update("quantity", newQuantity).Error; err != nil {
					return nil, fmt.Errorf("failed to merge item: %w", err)
				}
				itemExists = true
				break
			}
		}

		if !itemExists {
			// Add new item to user cart
			newItem := CartItem{
				ID:        uuid.New(),
				CartID:    userCart.ID,
				ProductID: guestItem.ProductID,
				Quantity:  guestItem.Quantity,
				AddedAt:   time.Now(),
			}
			if err := s.db.Create(&newItem).Error; err != nil {
				return nil, fmt.Errorf("failed to add merged item: %w", err)
			}
		}
	}

	// Delete guest cart
	if err := s.db.Delete(&guestCart).Error; err != nil {
		return nil, fmt.Errorf("failed to delete guest cart: %w", err)
	}

	// Get updated user cart with details
	return s.cartToResponse(&userCart)
}

// Helper function to convert Cart to CartResponse
func (s *Service) cartToResponse(cart *Cart) (*CartResponse, error) {
	// Preload items with products
	var cartWithItems Cart
	if err := s.db.Preload("Items").First(&cartWithItems, cart.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load cart items: %w", err)
	}

	// Get product details for each item
	itemResponses := make([]CartItemResponse, len(cartWithItems.Items))
	var subtotal float64
	var totalItems int

	for i, item := range cartWithItems.Items {
		// Get product for this item
		var product models.Product
		if err := s.db.First(&product, item.ProductID).Error; err != nil {
			return nil, fmt.Errorf("product not found: %w", err)
		}

		// Parse images JSON
		var images []string
		if product.Images != nil {
			json.Unmarshal([]byte(*product.Images), &images)
		}

		productResponse := ProductResponse{
			ID:            product.ID,
			Title:         product.Title,
			Description:   product.Description,
			Condition:     product.Condition,
			StartingPrice: product.StartingPrice,
			ReservePrice:  product.ReservePrice,
			BuyNowPrice:   product.BuyNowPrice,
			Images:        images,
			Status:        product.Status,
		}

		lineTotal := product.StartingPrice * float64(item.Quantity)

		itemResponses[i] = CartItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			AddedAt:   item.AddedAt,
			Product:   productResponse,
			LineTotal: lineTotal,
		}

		subtotal += lineTotal
		totalItems += item.Quantity
	}

	cartResponse := CartResponse{
		ID:         cart.ID,
		UserID:     cart.UserID,
		Token:      cart.Token,
		ExpiresAt:  cart.ExpiresAt,
		Items:      itemResponses,
		ItemCount:  len(itemResponses),
		Subtotal:   subtotal,
		TotalItems: totalItems,
		CreatedAt:  cart.CreatedAt,
		UpdatedAt:  cart.UpdatedAt,
	}

	return &cartResponse, nil
}

// CleanupExpiredCarts removes expired carts
func (s *Service) CleanupExpiredCarts() error {
	return s.db.Where("expires_at < ?", time.Now()).Delete(&Cart{}).Error
}

// Helper function to generate guest token
func (s *Service) generateGuestToken() string {
	return uuid.New().String()
}