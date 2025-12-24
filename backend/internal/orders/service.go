package orders

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/blytz.live.remake/backend/internal/cart"
	"github.com/blytz.live.remake/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func stringPtr(s string) *string {
	return &s
}

// Service handles order business logic
// NOTE: This is Phase 5A - Order Management implementation
// TODO: Add more sophisticated tax calculation based on jurisdiction
// TODO: Add shipping carrier integration
// TODO: Implement order cancellation policies
type Service struct {
	db          *gorm.DB
	cartService *cart.Service
}

// NewService creates a new order service
func NewService(db *gorm.DB, cartService *cart.Service) *Service {
	return &Service{
		db:          db,
		cartService: cartService,
	}
}

// CreateOrder creates an order from cart
func (s *Service) CreateOrder(userID uuid.UUID, req OrderCreateRequest) (*OrderResponse, error) {
	// Get cart details
	cart, err := s.cartService.GetCartWithDetails(req.CartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Validate cart has items
	if len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Validate all products are still available and prices haven't changed
	for _, item := range cart.Items {
		var product models.Product
		if err := s.db.First(&product, item.ProductID).Error; err != nil {
			return nil, fmt.Errorf("product %s not found: %w", item.ProductID, err)
		}

		if product.Status != "active" {
			return nil, fmt.Errorf("product %s is no longer available", product.Title)
		}

		// Check if price has changed significantly (more than 5%)
		priceDifference := (product.StartingPrice - item.Product.StartingPrice) / item.Product.StartingPrice
		if priceDifference > 0.05 || priceDifference < -0.05 {
			return nil, fmt.Errorf("price for product %s has changed", product.Title)
		}
	}

	// Calculate totals
	subtotal := cart.Subtotal
	taxAmount := s.calculateTax(subtotal, req.ShippingAddress)
	shippingCost := s.calculateShippingCost(req.ShippingAddress, cart.TotalItems)
	totalAmount := subtotal + taxAmount + shippingCost

	// Create order
	order := Order{
		ID:              uuid.New(),
		UserID:          userID,
		Status:          "pending",
		TotalAmount:     totalAmount,
		Subtotal:        subtotal,
		TaxAmount:       taxAmount,
		ShippingCost:    shippingCost,
		DiscountAmount:  0,
		ShippingAddress: &req.ShippingAddress,
		BillingAddress:  &req.BillingAddress,
		Notes:           req.Notes,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Start transaction
	tx := s.db.Begin()

	// Save order
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Create order items
	orderItems := make([]OrderItem, len(cart.Items))
	for i, cartItem := range cart.Items {
		orderItem := OrderItem{
			ID:        uuid.New(),
			OrderID:   order.ID,
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
			UnitPrice: cartItem.Product.StartingPrice,
			Total:     cartItem.LineTotal,
			CreatedAt: time.Now(),
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}

		orderItems[i] = orderItem
	}

	// Update product stock (reserve stock)
	for _, cartItem := range cart.Items {
		if err := s.reserveStock(tx, cartItem.ProductID, cartItem.Quantity); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to reserve stock: %w", err)
		}
	}

	// Clear cart
	if err := s.cartService.ClearCart(req.CartID); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to clear cart: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Convert to response
	return s.orderToResponse(&order, orderItems)
}

// GetOrder gets order by ID with user validation
func (s *Service) GetOrder(orderID, userID uuid.UUID) (*OrderResponse, error) {
	var order Order
	if err := s.db.Preload("Items").First(&order, orderID).Error; err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// Check ownership
	if order.UserID != userID {
		return nil, errors.New("order not found")
	}

	// Get order items with products
	var orderItems []OrderItem
	if err := s.db.Preload("Items.Product").Find(&orderItems, "order_id = ?", orderID).Error; err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}

	return s.orderToResponse(&order, orderItems)
}

// ListOrders lists user orders with filtering
func (s *Service) ListOrders(userID uuid.UUID, req OrderListRequest, pagination PaginationRequest) (*PaginatedResponse, error) {
	query := s.db.Model(&Order{}).Where("user_id = ?", userID)

	// Apply filters
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if req.Search != "" {
		searchPattern := "%" + req.Search + "%"
		query = query.Where("id ILIKE ? OR notes ILIKE ?", searchPattern, searchPattern)
	}

	// Count total records
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count orders: %w", err)
	}

	// Apply sorting
	sortBy := "created_at"
	if req.SortBy != "" {
		sortBy = req.SortBy
	}

	sortDirection := "desc"
	if req.SortDirection != "" {
		sortDirection = req.SortDirection
	}

	query = query.Order(sortBy + " " + sortDirection)

	// Fetch paginated results
	var orders []Order
	if err := query.Offset(pagination.GetOffset()).
		Limit(pagination.PageSize).
		Preload("Items").
		Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}

	// Convert to responses
	orderResponses := make([]*OrderResponse, len(orders))
	for i, order := range orders {
		var orderItems []OrderItem
		if err := s.db.Preload("Items.Product").Find(&orderItems, "order_id = ?", order.ID).Error; err != nil {
			return nil, fmt.Errorf("failed to get order items: %w", err)
		}

		response, err := s.orderToResponse(&order, orderItems)
		if err != nil {
			return nil, fmt.Errorf("failed to convert order: %w", err)
		}

		orderResponses[i] = response
	}

	return NewPaginatedResponse(orderResponses, total, pagination.Page, pagination.PageSize), nil
}

// UpdateOrderStatus updates order status (admin/seller function)
func (s *Service) UpdateOrderStatus(orderID uuid.UUID, req UpdateOrderStatusRequest) (*OrderResponse, error) {
	var order Order
	if err := s.db.First(&order, orderID).Error; err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// Validate status transition
	if !s.isValidStatusTransition(order.Status, req.Status) {
		return nil, fmt.Errorf("invalid status transition from %s to %s", order.Status, req.Status)
	}

	// Update order
	updates := map[string]interface{}{
		"status":     req.Status,
		"updated_at": time.Now(),
	}

	if req.Notes != nil {
		updates["notes"] = req.Notes
	}

	if err := s.db.Model(&order).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	// Get updated order with items
	var orderItems []OrderItem
	if err := s.db.Preload("Items.Product").Find(&orderItems, "order_id = ?", orderID).Error; err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}

	return s.orderToResponse(&order, orderItems)
}

// CancelOrder cancels an order (user function)
func (s *Service) CancelOrder(orderID, userID uuid.UUID) error {
	var order Order
	if err := s.db.First(&order, orderID).Error; err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// Check ownership
	if order.UserID != userID {
		return errors.New("order not found")
	}

	// Check if order can be cancelled
	if order.Status != "pending" {
		return errors.New("order cannot be cancelled in current status")
	}

	// Start transaction
	tx := s.db.Begin()

	// Update order status
	if err := tx.Model(&order).Updates(map[string]interface{}{
		"status":     "cancelled",
		"updated_at": time.Now(),
	}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	// Release stock reservations
	if err := s.releaseStockReservations(tx, orderID); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to release stock: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetOrderStatistics gets order statistics
func (s *Service) GetOrderStatistics() (*OrderStatistics, error) {
	var stats OrderStatistics

	// Get total orders and revenue
	var totalRevenue float64
	if err := s.db.Model(&Order{}).
		Select("COUNT(*) as total_orders, COALESCE(SUM(total_amount), 0) as total_revenue").
		Where("status NOT IN (?)", []string{"cancelled"}).
		Row().
		Scan(&stats.TotalOrders, &totalRevenue); err != nil {
		return nil, fmt.Errorf("failed to get order statistics: %w", err)
	}
	stats.TotalRevenue = totalRevenue

	// Get item count
	if err := s.db.Model(&OrderItem{}).
		Select("COALESCE(SUM(quantity), 0) as total_items").
		Row().
		Scan(&stats.TotalItems); err != nil {
		return nil, fmt.Errorf("failed to get item statistics: %w", err)
	}

	// Calculate average order value
	if stats.TotalOrders > 0 {
		stats.AverageOrderValue = stats.TotalRevenue / float64(stats.TotalOrders)
	}

	// Get status counts
	statusCounts := []struct {
		Status string
		Count  int
	}{}

	if err := s.db.Model(&Order{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&statusCounts).Error; err != nil {
		return nil, fmt.Errorf("failed to get status counts: %w", err)
	}

	for _, sc := range statusCounts {
		switch sc.Status {
		case "pending":
			stats.PendingOrders = sc.Count
		case "processing":
			stats.ProcessingOrders = sc.Count
		case "shipped":
			stats.ShippedOrders = sc.Count
		case "delivered":
			stats.DeliveredOrders = sc.Count
		case "cancelled":
			stats.CancelledOrders = sc.Count
		}
	}

	return &stats, nil
}

// Helper functions

// calculateTax calculates tax based on address (simplified)
func (s *Service) calculateTax(subtotal float64, address Address) float64 {
	// Simplified tax calculation - in production, use tax service API
	taxRate := 0.08 // 8% default tax rate

	// Different tax rates based on location (simplified)
	switch address.Country {
	case "US":
		// Could be different by state in production
		taxRate = 0.08
	case "GB":
		taxRate = 0.20 // VAT
	case "DE":
		taxRate = 0.19 // VAT
	case "FR":
		taxRate = 0.20 // VAT
	default:
		taxRate = 0.10 // Default international rate
	}

	return subtotal * taxRate
}

// calculateShippingCost calculates shipping cost (simplified)
func (s *Service) calculateShippingCost(address Address, totalItems int) float64 {
	// Simplified shipping cost calculation
	baseCost := 5.99

	// Different rates based on location
	switch address.Country {
	case "US":
		if totalItems > 5 {
			baseCost = 12.99
		}
	case "CA":
		baseCost = 15.99
	default: // International
		baseCost = 25.99
		if totalItems > 3 {
			baseCost = 45.99
		}
	}

	return baseCost
}

// reserveStock reserves stock for order
func (s *Service) reserveStock(tx *gorm.DB, productID uuid.UUID, quantity int) error {
	var stock models.InventoryStock

	if err := tx.Where("product_id = ?", productID).
		First(&stock).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("product stock not found")
		}
		return fmt.Errorf("failed to get stock: %w", err)
	}

	// Check availability
	if stock.Available < quantity {
		return errors.New("insufficient stock")
	}

	// Update reservation
	newReserved := stock.Reserved + quantity
	newAvailable := stock.Quantity - newReserved

	if err := tx.Model(&stock).
		Updates(map[string]interface{}{
			"reserved":  newReserved,
			"available": newAvailable,
		}).Error; err != nil {
		return fmt.Errorf("failed to reserve stock: %w", err)
	}

	// Record stock movement
	movement := models.StockMovement{
		ProductID:    productID,
		MovementType: "reserve",
		Quantity:     quantity,
		Reference:    stringPtr("order_reservation"),
		Notes:        stringPtr(fmt.Sprintf("Reserved %d units for order", quantity)),
	}

	if err := tx.Create(&movement).Error; err != nil {
		return fmt.Errorf("failed to record stock movement: %w", err)
	}

	return nil
}

// releaseStockReservations releases stock for cancelled order
func (s *Service) releaseStockReservations(tx *gorm.DB, orderID uuid.UUID) error {
	// Get order items
	var items []OrderItem
	if err := tx.Where("order_id = ?", orderID).Find(&items).Error; err != nil {
		return fmt.Errorf("failed to get order items: %w", err)
	}

	// Release stock for each item
	for _, item := range items {
		var stock models.InventoryStock
		if err := tx.Where("product_id = ?", item.ProductID).First(&stock).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				continue
			}
			return fmt.Errorf("failed to get stock: %w", err)
		}

		// Calculate new values
		newReserved := stock.Reserved - item.Quantity
		newAvailable := stock.Quantity - newReserved

		if newReserved < 0 {
			newReserved = 0
		}

		// Update stock
		if err := tx.Model(&stock).
			Updates(map[string]interface{}{
				"reserved":  newReserved,
				"available": newAvailable,
			}).Error; err != nil {
			return fmt.Errorf("failed to release stock: %w", err)
		}

		// Record stock movement
		movement := models.StockMovement{
			ProductID:    item.ProductID,
			MovementType: "release",
			Quantity:     item.Quantity,
			Reference:    stringPtr("order_cancellation"),
			Notes:        stringPtr(fmt.Sprintf("Released %d units from cancelled order", item.Quantity)),
		}

		if err := tx.Create(&movement).Error; err != nil {
			return fmt.Errorf("failed to record stock movement: %w", err)
		}
	}

	return nil
}

// isValidStatusTransition validates order status transitions
func (s *Service) isValidStatusTransition(currentStatus, newStatus string) bool {
	validTransitions := map[string][]string{
		"pending":    {"processing", "cancelled"},
		"processing": {"shipped", "cancelled"},
		"shipped":    {"delivered"},
		"delivered":  {}, // Terminal state
		"cancelled":  {}, // Terminal state
	}

	allowedStatuses, exists := validTransitions[currentStatus]
	if !exists {
		return false
	}

	for _, status := range allowedStatuses {
		if status == newStatus {
			return true
		}
	}

	return false
}

// orderToResponse converts order to response
func (s *Service) orderToResponse(order *Order, items []OrderItem) (*OrderResponse, error) {
	// Convert items to responses
	itemResponses := make([]OrderItemResponse, len(items))
	var totalQuantity int

	for i, item := range items {
		// Parse images JSON from product model
		var images []string
		if item.Product.Images != nil {
			json.Unmarshal([]byte(*item.Product.Images), &images)
		}

		productResponse := ProductResponse{
			ID:            item.Product.ID,
			Title:         item.Product.Title,
			Description:   item.Product.Description,
			Condition:     item.Product.Condition,
			StartingPrice: item.Product.StartingPrice,
			ReservePrice:  item.Product.ReservePrice,
			BuyNowPrice:   item.Product.BuyNowPrice,
			Images:        images,
			Status:        item.Product.Status,
		}

		itemResponses[i] = OrderItemResponse{
			ID:        item.ID,
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			Total:     item.Total,
			Product:   productResponse,
			CreatedAt: item.CreatedAt,
		}

		totalQuantity += item.Quantity
	}

	return &OrderResponse{
		ID:              order.ID,
		UserID:          order.UserID,
		Status:          order.Status,
		TotalAmount:     order.TotalAmount,
		Subtotal:        order.Subtotal,
		TaxAmount:       order.TaxAmount,
		ShippingCost:    order.ShippingCost,
		DiscountAmount:  order.DiscountAmount,
		ShippingAddress: order.ShippingAddress,
		BillingAddress:  order.BillingAddress,
		PaymentID:       order.PaymentID,
		TrackingNumber:  order.TrackingNumber,
		Notes:           order.Notes,
		Items:           itemResponses,
		ItemCount:       len(itemResponses),
		TotalQuantity:   totalQuantity,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}, nil
}
