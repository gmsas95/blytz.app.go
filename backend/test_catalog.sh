#!/bin/bash

# Blytz.live - Catalog Management Test Script
# Tests comprehensive catalog functionality including categories, variants, collections, and inventory

API_BASE="http://localhost:8080/api/v1"
TOKEN=""
USER_ID=""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test results
TESTS_PASSED=0
TESTS_FAILED=0

# Helper functions
log() {
    echo -e "${BLUE}[$(date '+%Y-%m-%d %H:%M:%S')]${NC} $1"
}

success() {
    echo -e "${GREEN}✅ $1${NC}"
    ((TESTS_PASSED++))
}

error() {
    echo -e "${RED}❌ $1${NC}"
    ((TESTS_FAILED++))
}

warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Make API request
make_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    local token=$4
    
    curl -s -w "%{http_code}" \
        -X "$method" \
        -H "Content-Type: application/json" \
        ${token:+-H "Authorization: Bearer $token"} \
        -d "$data" \
        "$API_BASE$endpoint"
}

# Register test user
setup_test_user() {
    log "Setting up test user..."
    
    # Register user
    response=$(make_request "POST" "/auth/register" '{
        "email": "catalogtest@example.com",
        "password": "testpassword123",
        "first_name": "Catalog",
        "last_name": "Tester",
        "role": "seller"
    }')
    
    http_code="${response: -3}"
    response_body="${response%???}"
    
    if [ "$http_code" = "201" ] || [ "$http_code" = "409" ]; then
        success "Test user setup complete"
    else
        warning "User setup response: $response_body"
    fi
    
    # Login to get token
    response=$(make_request "POST" "/auth/login" '{
        "email": "catalogtest@example.com",
        "password": "testpassword123"
    }')
    
    http_code="${response: -3}"
    response_body="${response%???}"
    
    if [ "$http_code" = "200" ]; then
        TOKEN=$(echo "$response_body" | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
        USER_ID=$(echo "$response_body" | grep -o '"id":"[^"]*' | cut -d'"' -f4)
        success "User authenticated successfully"
    else
        error "Failed to authenticate: $response_body"
        exit 1
    fi
}

# Test Category Management
test_categories() {
    log "\n=== Testing Category Management ==="
    
    # Create main category
    log "Creating main category..."
    response=$(make_request "POST" "/catalog/categories" '{
        "name": "Electronics",
        "description": "Electronic devices and accessories"
    }' "$TOKEN")
    
    http_code="${response: -3}"
    response_body="${response%???}"
    
    if [ "$http_code" = "201" ]; then
        success "Main category created"
        CATEGORY_ID=$(echo "$response_body" | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    else
        error "Failed to create main category: $response_body"
        return
    fi
    
    # Create subcategory
    log "Creating subcategory..."
    response=$(make_request "POST" "/catalog/categories" "{
        \"name\": \"Smartphones\",
        \"description\": \"Mobile phones and accessories\",
        \"parent_id\": \"$CATEGORY_ID\"
    }" "$TOKEN")
    
    http_code="${response: -3}"
    response_body="${response%???}"
    
    if [ "$http_code" = "201" ]; then
        success "Subcategory created"
        SUBCATEGORY_ID=$(echo "$response_body" | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    else
        error "Failed to create subcategory: $response_body"
    fi
    
    # Get category tree
    log "Getting category tree..."
    response=$(make_request "GET" "/catalog/categories?include_product_count=true" "" "")
    
    http_code="${response: -3}"
    response_body="${response%???}"
    
    if [ "$http_code" = "200" ]; then
        success "Category tree retrieved"
    else
        error "Failed to get category tree: $response_body"
    fi
    
    # Create category attribute
    log "Creating category attribute..."
    response=$(make_request "POST" "/catalog/categories/$CATEGORY_ID/attributes" '{
        "name": "Warranty Period",
        "type": "select",
        "required": false,
        "options": ["1 Year", "2 Years", "3 Years"],
        "default_value": "1 Year"
    }' "$TOKEN")
    
    http_code="${response: -3}"
    response_body="${response%???}"
    
    if [ "$http_code" = "201" ]; then
        success "Category attribute created"
        ATTRIBUTE_ID=$(echo "$response_body" | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    else
        error "Failed to create category attribute: $response_body"
    fi
    
    # Get category attributes
    log "Getting category attributes..."
    response=$(make_request "GET" "/catalog/categories/$CATEGORY_ID/attributes" "" "")
    
    http_code="${response: -3}"
    response_body="${response%???}"
    
    if [ "$http_code" = "200" ]; then
        success "Category attributes retrieved"
    else
        error "Failed to get category attributes: $response_body"
    fi
}

# Test Product Variants
test_variants() {
    log "\n=== Testing Product Variants ==="
    
    # Create a test product first
    log "Creating test product for variants..."
    response=$(make_request "POST" "/products" "{
        \"category_id\": \"$CATEGORY_ID\",
        \"title\": \"Premium Smartphone\",
        \"description\": \"A high-quality smartphone with advanced features\",
        \"condition\": \"new\",
        \"starting_price\": 599.99,
        \"buy_now_price\": 699.99,
        \"images\": [\"https://example.com/phone1.jpg\"],
        \"status\": \"active\"
    }" "$TOKEN")
    
    http_code="${response: -3}"
    response_body="${response%???}"
    
    if [ "$http_code" = "201" ]; then
        success "Test product created"
        PRODUCT_ID=$(echo "$response_body" | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    else
        error "Failed to create test product: $response_body"
        return
    fi
    
    # Create product variants
    log "Creating product variants..."
    response=$(make_request "POST" "/catalog/variants" "{
        \"product_id\": \"$PRODUCT_ID\",
        \"variant\": {
            \"title\": \"Premium Smartphone - Black 128GB\",
            \"sku\": \"PHONE-BLK-128\",
            \"price\": 599.99,
            \"inventory\": 50,
            \"attributes\": {
                \"color\": \"Black\",
                \"storage\": \"128GB\"
            }
        }
    }" "$TOKEN")
    
    http_code="${response: -3}"
    response_body="${response%???}"
    
    if [ "$http_code" = "201" ]; then
        success "Product variant created"
        VARIANT_ID=$(echo "$response_body" | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    else
        error "Failed to create product variant: $response_body"
    fi
    
    # Create another variant
    response=$(make_request "POST" "/catalog/variants" "{
        \"product_id\": \"$PRODUCT_ID\",
        \"variant\": {
            \"title\": \"Premium Smartphone - White 256GB\",
            \"sku\": \"PHONE-WHT-256\",
            \"price\": 699.99,
            \"inventory\": 30,
            \"attributes\": {
                \"color\": \"White\",
                \"storage\": \"256GB\"
            }
        }
    }" "$TOKEN")
    
    http_code="${response: -3}"
    
    if [ "$http_code" = "201" ]; then
        success "Second product variant created"
    else
        error "Failed to create second product variant"
    fi
    
    # Get product variants
    log "Getting product variants..."
    response=$(make_request "GET" "/catalog/variants/products/$PRODUCT_ID" "" "")
    
    http_code="${response: -3}"
    response_body="${response%???}"
    
    if [ "$http_code" = "200" ]; then
        success "Product variants retrieved"
    else
        error "Failed to get product variants: $response_body"
    fi
    
    # Update variant
    log "Updating product variant..."
    response=$(make_request "PUT" "/catalog/variants/$VARIANT_ID" "{
        \"title\": \"Premium Smartphone - Black 128GB (Updated)\",
        \"price\": 579.99,
        \"inventory\": 45
    }" "$TOKEN")
    
    http_code="${response: -3}"
    
    if [ "$http_code" = "200" ]; then
        success "Product variant updated"
    else
        error "Failed to update product variant"
    fi
}

# Test Product Collections
test_collections() {
    log "\n=== Testing Product Collections ==="
    
    # Create collection
    log "Creating product collection..."
    response=$(make_request "POST" "/catalog/collections" '{
        "name": "Featured Electronics",
        "description": "Hand-picked electronics for our customers",
        "is_active": true,
        "product_ids": []
    }' "$TOKEN")
    
    http_code="${response: -3}"
    response_body="${response%???}"
    
    if [ "$http_code" = "201" ]; then
        success "Product collection created"
        COLLECTION_ID=$(echo "$response_body" | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    else
        error "Failed to create collection: $response_body"
        return
    fi
    
    # Add products to collection
    log "Adding products to collection..."
    response=$(make_request "POST" "/catalog/collections/$COLLECTION_ID/products" "{
        \"product_ids\": [\"$PRODUCT_ID\"]
    }" "$TOKEN")
    
    http_code="${response: -3}"
    
    if [ "$http_code" = "200" ]; then
        success "Products added to collection"
    else
        error "Failed to add products to collection"
    fi
    
    # Get collections
    log "Getting product collections..."
    response=$(make_request "GET" "/catalog/collections" "" "")
    
    http_code="${response: -3}"
    response_body="${response%???}"
    
    if [ "$http_code" = "200" ]; then
        success "Product collections retrieved"
    else
        error "Failed to get collections: $response_body"
    fi
    
    # Get specific collection
    log "Getting specific collection..."
    response=$(make_request "GET" "/catalog/collections/$COLLECTION_ID" "" "")
    
    http_code="${response: -3}"
    response_body="${response%???}"
    
    if [ "$http_code" = "200" ]; then
        success "Specific collection retrieved"
    else
        error "Failed to get specific collection: $response_body"
    fi
}

# Test Inventory Management
test_inventory() {
    log "\n=== Testing Inventory Management ==="
    
    # Get product inventory
    log "Getting product inventory..."
    response=$(make_request "GET" "/catalog/inventory/products/$PRODUCT_ID" "" "")
    
    http_code="${response: -3}"
    response_body="${response%???}"
    
    if [ "$http_code" = "200" ]; then
        success "Product inventory retrieved"
    else
        error "Failed to get product inventory: $response_body"
    fi
    
    # Update inventory
    log "Updating product inventory..."
    response=$(make_request "PUT" "/catalog/inventory/products/$PRODUCT_ID" '{
        "quantity": 100,
        "low_stock_alert": 15,
        "track_inventory": true,
        "allow_backorder": false
    }' "$TOKEN")
    
    http_code="${response: -3}"
    response_body="${response%???}"
    
    if [ "$http_code" = "200" ]; then
        success "Product inventory updated"
    else
        error "Failed to update product inventory: $response_body"
    fi
    
    # Create stock movement
    log "Creating stock movement..."
    response=$(make_request "POST" "/catalog/inventory/products/$PRODUCT_ID/movements" '{
        "movement_type": "in",
        "quantity": 25,
        "reference": "Stock replenishment",
        "notes": "New inventory batch received"
    }' "$TOKEN")
    
    http_code="${response: -3}"
    
    if [ "$http_code" = "201" ]; then
        success "Stock movement created"
    else
        error "Failed to create stock movement"
    fi
    
    # Get stock movements
    log "Getting stock movements..."
    response=$(make_request "GET" "/catalog/inventory/products/$PRODUCT_ID/movements?limit=5" "" "")
    
    http_code="${response: -3}"
    response_body="${response%???}"
    
    if [ "$http_code" = "200" ]; then
        success "Stock movements retrieved"
    else
        error "Failed to get stock movements: $response_body"
    fi
}

# Test Search and Discovery
test_search() {
    log "\n=== Testing Search and Discovery ==="
    
    # Search products
    log "Searching products..."
    response=$(make_request "GET" "/catalog/search/products?q=Smartphone&limit=10" "" "")
    
    http_code="${response: -3}"
    response_body="${response%???}"
    
    if [ "$http_code" = "200" ]; then
        success "Product search completed"
    else
        error "Failed to search products: $response_body"
    fi
    
    # Get featured products
    log "Getting featured products..."
    response=$(make_request "GET" "/catalog/search/products/featured?limit=5" "" "")
    
    http_code="${response: -3}"
    response_body="${response%???}"
    
    if [ "$http_code" = "200" ]; then
        success "Featured products retrieved"
    else
        error "Failed to get featured products: $response_body"
    fi
    
    # Get related products
    log "Getting related products..."
    response=$(make_request "GET" "/catalog/search/products/$PRODUCT_ID/related?limit=3" "" "")
    
    http_code="${response: -3}"
    response_body="${response%???}"
    
    if [ "$http_code" = "200" ]; then
        success "Related products retrieved"
    else
        error "Failed to get related products: $response_body"
    fi
}

# Test Catalog Statistics
test_statistics() {
    log "\n=== Testing Catalog Statistics ==="
    
    # Get catalog stats
    log "Getting catalog statistics..."
    response=$(make_request "GET" "/catalog/stats/catalog" "" "")
    
    http_code="${response: -3}"
    response_body="${response%???}"
    
    if [ "$http_code" = "200" ]; then
        success "Catalog statistics retrieved"
    else
        error "Failed to get catalog statistics: $response_body"
    fi
    
    # Get category stats
    log "Getting category statistics..."
    response=$(make_request "GET" "/catalog/stats/categories/$CATEGORY_ID" "" "")
    
    http_code="${response: -3}"
    response_body="${response%???}"
    
    if [ "$http_code" = "200" ]; then
        success "Category statistics retrieved"
    else
        error "Failed to get category statistics: $response_body"
    fi
}

# Test Bulk Operations
test_bulk_operations() {
    log "\n=== Testing Bulk Operations ==="
    
    # Create another category for bulk operations
    response=$(make_request "POST" "/catalog/categories" '{
        "name": "Test Category 2",
        "description": "Second test category"
    }' "$TOKEN")
    
    http_code="${response: -3}"
    response_body="${response%???}"
    
    if [ "$http_code" = "201" ]; then
        TEST_CATEGORY2_ID=$(echo "$response_body" | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    else
        error "Failed to create second test category"
        return
    fi
    
    # Create third category
    response=$(make_request "POST" "/catalog/categories" '{
        "name": "Test Category 3",
        "description": "Third test category"
    }' "$TOKEN")
    
    http_code="${response: -3}"
    response_body="${response%???}"
    
    if [ "$http_code" = "201" ]; then
        TEST_CATEGORY3_ID=$(echo "$response_body" | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    else
        error "Failed to create third test category"
        return
    fi
    
    # Perform bulk operation - deactivate categories
    log "Performing bulk category operation (deactivate)..."
    response=$(make_request "POST" "/catalog/categories/bulk" "{
        \"category_ids\": [\"$TEST_CATEGORY2_ID\", \"$TEST_CATEGORY3_ID\"],
        \"operation\": \"deactivate\"
    }" "$TOKEN")
    
    http_code="${response: -3}"
    response_body="${response%???}"
    
    if [ "$http_code" = "200" ]; then
        success "Bulk category operation completed"
    else
        error "Failed to perform bulk category operation: $response_body"
    fi
    
    # Test bulk variant creation
    log "Testing bulk variant creation..."
    response=$(make_request "POST" "/catalog/variants/bulk" "{
        \"product_id\": \"$PRODUCT_ID\",
        \"variants\": [
            {
                \"title\": \"Phone - Red 128GB\",
                \"sku\": \"PHONE-RED-128\",
                \"price\": 599.99,
                \"inventory\": 20,
                \"attributes\": {
                    \"color\": \"Red\",
                    \"storage\": \"128GB\"
                }
            },
            {
                \"title\": \"Phone - Blue 256GB\",
                \"sku\": \"PHONE-BLU-256\",
                \"price\": 699.99,
                \"inventory\": 15,
                \"attributes\": {
                    \"color\": \"Blue\",
                    \"storage\": \"256GB\"
                }
            }
        ]
    }" "$TOKEN")
    
    http_code="${response: -3}"
    response_body="${response%???}"
    
    if [ "$http_code" = "201" ]; then
        success "Bulk variant creation completed"
    else
        error "Failed to create bulk variants: $response_body"
    fi
}

# Cleanup test data
cleanup() {
    log "\n=== Cleaning Up Test Data ==="
    
    # Delete variants
    if [ ! -z "$VARIANT_ID" ]; then
        make_request "DELETE" "/catalog/variants/$VARIANT_ID" "" "$TOKEN" > /dev/null
        log "Deleted test variant"
    fi
    
    # Delete collection
    if [ ! -z "$COLLECTION_ID" ]; then
        make_request "DELETE" "/catalog/collections/$COLLECTION_ID" "" "$TOKEN" > /dev/null
        log "Deleted test collection"
    fi
    
    # Delete test categories
    if [ ! -z "$TEST_CATEGORY3_ID" ]; then
        make_request "DELETE" "/catalog/categories/$TEST_CATEGORY3_ID" "" "$TOKEN" > /dev/null
        log "Deleted test category 3"
    fi
    
    if [ ! -z "$TEST_CATEGORY2_ID" ]; then
        make_request "DELETE" "/catalog/categories/$TEST_CATEGORY2_ID" "" "$TOKEN" > /dev/null
        log "Deleted test category 2"
    fi
    
    if [ ! -z "$SUBCATEGORY_ID" ]; then
        make_request "DELETE" "/catalog/categories/$SUBCATEGORY_ID" "" "$TOKEN" > /dev/null
        log "Deleted test subcategory"
    fi
    
    if [ ! -z "$CATEGORY_ID" ]; then
        make_request "DELETE" "/catalog/categories/$CATEGORY_ID" "" "$TOKEN" > /dev/null
        log "Deleted test main category"
    fi
    
    # Delete test product
    if [ ! -z "$PRODUCT_ID" ]; then
        make_request "DELETE" "/products/$PRODUCT_ID" "" "$TOKEN" > /dev/null
        log "Deleted test product"
    fi
    
    # Delete user (logout)
    if [ ! -z "$TOKEN" ]; then
        make_request "POST" "/auth/logout" "" "$TOKEN" > /dev/null
        log "Logged out test user"
    fi
    
    success "Cleanup completed"
}

# Main execution
main() {
    log "Starting Blytz.live Catalog Management Tests"
    log "=========================================="
    
    # Setup
    setup_test_user
    
    # Run tests
    test_categories
    test_variants
    test_collections
    test_inventory
    test_search
    test_statistics
    test_bulk_operations
    
    # Cleanup
    cleanup
    
    # Summary
    log "\n=========================================="
    log "Test Summary:"
    log "Tests Passed: $TESTS_PASSED"
    log "Tests Failed: $TESTS_FAILED"
    log "=========================================="
    
    if [ $TESTS_FAILED -eq 0 ]; then
        success "All catalog management tests passed! 🎉"
        exit 0
    else
        error "Some tests failed. Check the output above."
        exit 1
    fi
}

# Trap to ensure cleanup on script exit
trap cleanup EXIT

# Run main function
main