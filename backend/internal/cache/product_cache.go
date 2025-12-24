package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/blytz.live.remake/backend/internal/models"
	"github.com/blytz.live.remake/backend/pkg/logging"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProductCache provides cached product operations
type ProductCache struct {
	db    *gorm.DB
	cache *Cache
	logger *logging.Logger
}

// NewProductCache creates a new product cache
func NewProductCache(db *gorm.DB, cache *Cache) *ProductCache {
	logger := logging.NewLogger()
	return &ProductCache{
		db:     db,
		cache:   cache,
		logger:  logger,
	}
}

// GetProduct retrieves a product with caching
func (pc *ProductCache) GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	cacheKey := fmt.Sprintf("product:%s", id.String())
	
	// Try cache first
	var product models.Product
	err := pc.cache.Get(ctx, cacheKey, &product)
	if err == nil {
		pc.logger.Debug("Product cache hit", map[string]interface{}{
			"product_id": id,
		})
		return &product, nil
	}
	
	// Cache miss, get from database
	err = pc.db.WithContext(ctx).
		Preload("Seller").
		Preload("Category").
		First(&product, "id = ?", id).Error
	
	if err != nil {
		return nil, err
	}
	
	// Cache for 15 minutes
	cacheErr := pc.cache.Set(ctx, cacheKey, &product, 15*time.Minute)
	if cacheErr != nil {
		pc.logger.Error("Failed to cache product", map[string]interface{}{
			"product_id": id,
			"error":      cacheErr.Error(),
		})
	}
	
	pc.logger.Debug("Product cached", map[string]interface{}{
		"product_id": id,
	})
	
	return &product, nil
}

// ListProducts retrieves paginated products with caching
func (pc *ProductCache) ListProducts(ctx context.Context, page, limit int, categoryID *uuid.UUID) ([]models.Product, int64, error) {
	cacheKey := fmt.Sprintf("products:list:page:%d:limit:%d", page, limit)
	if categoryID != nil {
		cacheKey += fmt.Sprintf(":category:%s", categoryID.String())
	}
	
	// Try cache first
	var cachedData struct {
		Products []models.Product `json:"products"`
		Total    int64           `json:"total"`
	}
	
	err := pc.cache.Get(ctx, cacheKey, &cachedData)
	if err == nil {
		pc.logger.Debug("Product list cache hit", map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"category_id": categoryID,
		})
		return cachedData.Products, cachedData.Total, nil
	}
	
	// Cache miss, get from database
	var products []models.Product
	var total int64
	
	query := pc.db.WithContext(ctx).Model(&models.Product{}).Preload("Seller").Preload("Category")
	
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}
	
	// Count total
	err = query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	// Get paginated results
	offset := (page - 1) * limit
	err = query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&products).Error
	
	if err != nil {
		return nil, 0, err
	}
	
	// Cache for 5 minutes
	data := struct {
		Products []models.Product `json:"products"`
		Total    int64           `json:"total"`
	}{
		Products: products,
		Total:    total,
	}
	
	cacheErr := pc.cache.Set(ctx, cacheKey, &data, 5*time.Minute)
	if cacheErr != nil {
		pc.logger.Error("Failed to cache product list", map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"category_id": categoryID,
			"error":       cacheErr.Error(),
		})
	}
	
	pc.logger.Debug("Product list cached", map[string]interface{}{
		"page":        page,
		"limit":       limit,
		"category_id": categoryID,
	})
	
	return products, total, nil
}

// InvalidateProduct invalidates product cache
func (pc *ProductCache) InvalidateProduct(ctx context.Context, id uuid.UUID) error {
	cacheKey := fmt.Sprintf("product:%s", id.String())
	
	err := pc.cache.Delete(ctx, cacheKey)
	if err != nil {
		pc.logger.Error("Failed to invalidate product cache", map[string]interface{}{
			"product_id": id,
			"error":      err.Error(),
		})
		return err
	}
	
	pc.logger.Debug("Product cache invalidated", map[string]interface{}{
		"product_id": id,
	})
	
	// Also invalidate product lists
	err = pc.cache.Clear(ctx, "products:list:*")
	if err != nil {
		pc.logger.Error("Failed to invalidate product list cache", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}
	
	return nil
}

// GetLiveAuctions retrieves live auctions with caching
func (pc *ProductCache) GetLiveAuctions(ctx context.Context) ([]models.Auction, error) {
	cacheKey := "auctions:live"
	
	// Try cache first
	var auctions []models.Auction
	err := pc.cache.Get(ctx, cacheKey, &auctions)
	if err == nil {
		pc.logger.Debug("Live auctions cache hit")
		return auctions, nil
	}
	
	// Cache miss, get from database
	err = pc.db.WithContext(ctx).
		Preload("Product").
		Preload("Seller").
		Where("status = ? AND end_time > ?", "live", time.Now()).
		Order("end_time ASC").
		Find(&auctions).Error
	
	if err != nil {
		return nil, err
	}
	
	// Cache for 2 minutes (live auctions change frequently)
	cacheErr := pc.cache.Set(ctx, cacheKey, &auctions, 2*time.Minute)
	if cacheErr != nil {
		pc.logger.Error("Failed to cache live auctions", map[string]interface{}{
			"error": cacheErr.Error(),
		})
	}
	
	pc.logger.Debug("Live auctions cached")
	
	return auctions, nil
}

// IncrementViewCount increments product view count with cache
func (pc *ProductCache) IncrementViewCount(ctx context.Context, productID uuid.UUID) error {
	// Update database
	err := pc.db.WithContext(ctx).
		Model(&models.Product{}).
		Where("id = ?", productID).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
	
	if err != nil {
		return err
	}
	
	// Increment in cache for real-time updates
	cacheKey := fmt.Sprintf("product:views:%s", productID.String())
	_, err = pc.cache.Increment(ctx, cacheKey)
	if err != nil {
		pc.logger.Error("Failed to increment view count in cache", map[string]interface{}{
			"product_id": productID,
			"error":      err.Error(),
		})
	}
	
	// Invalidate product cache to trigger refresh
	pc.InvalidateProduct(ctx, productID)
	
	return nil
}

// GetProductStats retrieves product statistics with caching
func (pc *ProductCache) GetProductStats(ctx context.Context, productID uuid.UUID) (map[string]interface{}, error) {
	cacheKey := fmt.Sprintf("product:stats:%s", productID.String())
	
	// Try cache first
	var stats map[string]interface{}
	err := pc.cache.Get(ctx, cacheKey, &stats)
	if err == nil {
		return stats, nil
	}
	
	// Cache miss, calculate stats
	var product models.Product
	err = pc.db.WithContext(ctx).
		First(&product, "id = ?", productID).Error
	
	if err != nil {
		return nil, err
	}
	
	// Get view count from cache or database
	viewCount := product.ViewCount
	viewCacheKey := fmt.Sprintf("product:views:%s", productID.String())
	if cachedViews, err := pc.cache.GetClient().Get(ctx, viewCacheKey).Int64(); err == nil {
		viewCount += int(cachedViews)
	}
	
	stats = map[string]interface{}{
		"view_count": viewCount,
		"status":     product.Status,
		"updated_at": time.Now(),
	}
	
	// Cache for 10 minutes
	cacheErr := pc.cache.Set(ctx, cacheKey, stats, 10*time.Minute)
	if cacheErr != nil {
		pc.logger.Error("Failed to cache product stats", map[string]interface{}{
			"product_id": productID,
			"error":      cacheErr.Error(),
		})
	}
	
	return stats, nil
}