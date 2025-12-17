#!/usr/bin/env python3

import requests
import json
import uuid
from datetime import datetime

API_BASE = "http://localhost:8080/api/v1"

class CatalogTest:
    def __init__(self):
        self.token = ""
        self.user_id = ""
        self.test_category_id = ""
        self.test_product_id = ""
        self.test_collection_id = ""
        self.test_attribute_id = ""
        self.test_variant_id = ""
    
    def log(self, message, prefix="INFO"):
        timestamp = datetime.now().strftime('%Y-%m-%d %H:%M:%S')
        print(f"[{timestamp}] {prefix}: {message}")
    
    def success(self, message):
        self.log(f"✅ {message}", "SUCCESS")
    
    def error(self, message):
        self.log(f"❌ {message}", "ERROR")
    
    def warning(self, message):
        self.log(f"⚠️  {message}", "WARNING")
    
    def setup_test_user(self):
        """Register and login test user"""
        self.log("Setting up test user...")
        
        # Register user
        user_data = {
            "email": "admin@example.com",
            "password": "adminpassword123",
            "first_name": "Admin",
            "last_name": "User",
            "role": "admin"
        }
        
        response = requests.post(f"{API_BASE}/auth/register", json=user_data)
        if response.status_code not in [201, 409]:  # 409 = already exists
            self.error(f"User registration failed: {response.status_code} - {response.text}")
            return False
        
        self.success("User setup complete")
        
        # Login
        login_data = {
            "email": "admin@example.com",
            "password": "adminpassword123"
        }
        
        response = requests.post(f"{API_BASE}/auth/login", json=login_data)
        if response.status_code != 200:
            self.error(f"Login failed: {response.status_code} - {response.text}")
            return False
        
        data = response.json()
        self.token = data.get('access_token', '')
        self.user_id = data.get('user', {}).get('id', '')
        user_role = data.get('user', {}).get('role', '')
        
        if self.token:
            self.success(f"User authenticated successfully as {user_role}")
            return True
        else:
            self.error("Failed to get authentication token")
            return False
    
    def create_test_product(self):
        """Create a test product for testing variants and collections"""
        self.log("Creating test product...")
        
        # First we need to use existing test category
        response = requests.get(f"{API_BASE}/catalog/categories")
        if response.status_code == 200:
            data = response.json()
            categories = data.get('data', [])
            if categories:
                self.test_category_id = categories[0]['id']
            else:
                self.error("No categories found")
                return False
        
        product_data = {
            "category_id": self.test_category_id,
            "title": "Premium Test Product",
            "description": "A high-quality test product with advanced features",
            "condition": "new",
            "starting_price": 599.99,
            "buy_now_price": 699.99,
            "images": ["https://example.com/product1.jpg"],
            "status": "active"
        }
        
        response = requests.post(f"{API_BASE}/products", json=product_data, headers=self.get_headers())
        if response.status_code != 201:
            self.error(f"Product creation failed: {response.status_code} - {response.text}")
            return False
        
        data = response.json()
        self.test_product_id = data.get('product', {}).get('id', '')
        if self.test_product_id:
            self.success("Test product created successfully")
            return True
        else:
            self.error("Failed to get product ID")
            return False
    
    def get_headers(self):
        """Get authentication headers"""
        return {"Authorization": f"Bearer {self.token}"}
    
    def test_categories(self):
        """Test category management"""
        self.log("\n=== Testing Category Management ===")
        
        # Create main category
        self.log("Creating main category...")
        category_data = {
            "name": "Electronics Test",
            "description": "Electronic devices and accessories for testing"
        }
        
        response = requests.post(f"{API_BASE}/catalog/categories", 
                              json=category_data, headers=self.get_headers())
        
        if response.status_code == 201:
            self.success("Main category created")
            data = response.json()
            self.test_category_id = data.get('data', {}).get('id', '')
        else:
            self.error(f"Failed to create main category: {response.status_code} - {response.text}")
            return False
        
        # Get categories tree
        self.log("Getting categories tree...")
        response = requests.get(f"{API_BASE}/catalog/categories?include_product_count=true")
        if response.status_code == 200:
            self.success("Categories tree retrieved")
        else:
            self.error(f"Failed to get categories: {response.status_code}")
        
        # Create category attribute
        self.log("Creating category attribute...")
        attribute_data = {
            "name": "Warranty Period",
            "type": "select",
            "required": False,
            "options": ["1 Year", "2 Years", "3 Years"],
            "default_value": "1 Year"
        }
        
        response = requests.post(f"{API_BASE}/catalog/categories/{self.test_category_id}/attributes",
                              json=attribute_data, headers=self.get_headers())
        
        if response.status_code == 201:
            self.success("Category attribute created")
            data = response.json()
            self.test_attribute_id = data.get('data', {}).get('id', '')
        else:
            self.error(f"Failed to create category attribute: {response.status_code} - {response.text}")
        
        return True
    
    def test_variants(self):
        """Test product variants"""
        self.log("\n=== Testing Product Variants ===")
        
        if not self.test_product_id:
            self.error("No test product available")
            return False
        
        # Create product variant
        self.log("Creating product variant...")
        variant_data = {
            "product_id": self.test_product_id,
            "variant": {
                "title": "Premium Product - Black 128GB",
                "sku": "PROD-BLK-128",
                "price": 599.99,
                "inventory": 50,
                "attributes": {
                    "color": "Black",
                    "storage": "128GB"
                }
            }
        }
        
        response = requests.post(f"{API_BASE}/catalog/variants", 
                              json=variant_data, headers=self.get_headers())
        
        if response.status_code == 201:
            self.success("Product variant created")
            data = response.json()
            self.test_variant_id = data.get('data', {}).get('id', '')
        else:
            self.error(f"Failed to create product variant: {response.status_code} - {response.text}")
            return False
        
        # Get product variants
        self.log("Getting product variants...")
        response = requests.get(f"{API_BASE}/catalog/variants/products/{self.test_product_id}")
        if response.status_code == 200:
            self.success("Product variants retrieved")
            data = response.json()
            variants = data.get('data', [])
            self.log(f"Found {len(variants)} variants")
        else:
            self.error(f"Failed to get product variants: {response.status_code}")
        
        # Test bulk variant creation
        self.log("Testing bulk variant creation...")
        bulk_data = {
            "product_id": self.test_product_id,
            "variants": [
                {
                    "title": "Product - White 256GB",
                    "sku": "PROD-WHT-256",
                    "price": 699.99,
                    "inventory": 30,
                    "attributes": {
                        "color": "White",
                        "storage": "256GB"
                    }
                },
                {
                    "title": "Product - Blue 512GB",
                    "sku": "PROD-BLU-512",
                    "price": 799.99,
                    "inventory": 20,
                    "attributes": {
                        "color": "Blue",
                        "storage": "512GB"
                    }
                }
            ]
        }
        
        response = requests.post(f"{API_BASE}/catalog/variants/bulk",
                              json=bulk_data, headers=self.get_headers())
        
        if response.status_code == 201:
            self.success("Bulk variant creation completed")
            data = response.json()
            variants = data.get('data', [])
            self.log(f"Created {len(variants)} variants in bulk")
        else:
            self.error(f"Failed to create bulk variants: {response.status_code} - {response.text}")
        
        return True
    
    def test_collections(self):
        """Test product collections"""
        self.log("\n=== Testing Product Collections ===")
        
        # Create collection
        self.log("Creating product collection...")
        collection_data = {
            "name": "Featured Electronics",
            "description": "Hand-picked electronics for testing",
            "is_active": True
        }
        
        response = requests.post(f"{API_BASE}/catalog/collections",
                              json=collection_data, headers=self.get_headers())
        
        if response.status_code == 201:
            self.success("Product collection created")
            data = response.json()
            self.test_collection_id = data.get('data', {}).get('id', '')
        else:
            self.error(f"Failed to create collection: {response.status_code} - {response.text}")
            return False
        
        # Add products to collection
        self.log("Adding products to collection...")
        add_data = {
            "product_ids": [self.test_product_id]
        }
        
        response = requests.post(f"{API_BASE}/catalog/collections/{self.test_collection_id}/products",
                              json=add_data, headers=self.get_headers())
        
        if response.status_code == 200:
            self.success("Products added to collection")
        else:
            self.error(f"Failed to add products to collection: {response.status_code} - {response.text}")
        
        # Get collections
        self.log("Getting product collections...")
        response = requests.get(f"{API_BASE}/catalog/collections")
        if response.status_code == 200:
            self.success("Product collections retrieved")
            data = response.json()
            collections = data.get('data', [])
            self.log(f"Found {len(collections)} collections")
        else:
            self.error(f"Failed to get collections: {response.status_code}")
        
        return True
    
    def test_inventory(self):
        """Test inventory management"""
        self.log("\n=== Testing Inventory Management ===")
        
        if not self.test_product_id:
            self.error("No test product available")
            return False
        
        # Get product inventory
        self.log("Getting product inventory...")
        response = requests.get(f"{API_BASE}/catalog/inventory/products/{self.test_product_id}")
        if response.status_code == 200:
            self.success("Product inventory retrieved")
        else:
            self.error(f"Failed to get product inventory: {response.status_code}")
        
        # Update inventory
        self.log("Updating product inventory...")
        inventory_data = {
            "quantity": 100,
            "low_stock_alert": 15,
            "track_inventory": True,
            "allow_backorder": False
        }
        
        response = requests.put(f"{API_BASE}/catalog/inventory/products/{self.test_product_id}",
                             json=inventory_data, headers=self.get_headers())
        
        if response.status_code == 200:
            self.success("Product inventory updated")
        else:
            self.error(f"Failed to update product inventory: {response.status_code} - {response.text}")
        
        # Create stock movement
        self.log("Creating stock movement...")
        movement_data = {
            "movement_type": "in",
            "quantity": 25,
            "reference": "Stock replenishment",
            "notes": "New inventory batch received"
        }
        
        response = requests.post(f"{API_BASE}/catalog/inventory/products/{self.test_product_id}/movements",
                              json=movement_data, headers=self.get_headers())
        
        if response.status_code == 201:
            self.success("Stock movement created")
        else:
            self.error(f"Failed to create stock movement: {response.status_code} - {response.text}")
        
        return True
    
    def test_search_and_discovery(self):
        """Test search and discovery features"""
        self.log("\n=== Testing Search and Discovery ===")
        
        # Search products
        self.log("Searching products...")
        response = requests.get(f"{API_BASE}/catalog/search/products?q=Premium&limit=10")
        if response.status_code == 200:
            self.success("Product search completed")
            data = response.json()
            products = data.get('data', [])
            self.log(f"Found {len(products)} products")
        else:
            self.error(f"Failed to search products: {response.status_code}")
        
        # Get featured products
        self.log("Getting featured products...")
        response = requests.get(f"{API_BASE}/catalog/search/products/featured?limit=5")
        if response.status_code == 200:
            self.success("Featured products retrieved")
            data = response.json()
            products = data.get('data', [])
            self.log(f"Found {len(products)} featured products")
        else:
            self.error(f"Failed to get featured products: {response.status_code}")
        
        # Get related products
        if self.test_product_id:
            self.log("Getting related products...")
            response = requests.get(f"{API_BASE}/catalog/search/products/{self.test_product_id}/related?limit=3")
            if response.status_code == 200:
                self.success("Related products retrieved")
                data = response.json()
                products = data.get('data', [])
                self.log(f"Found {len(products)} related products")
            else:
                self.error(f"Failed to get related products: {response.status_code}")
        
        return True
    
    def test_statistics(self):
        """Test catalog statistics"""
        self.log("\n=== Testing Catalog Statistics ===")
        
        # Get catalog stats
        self.log("Getting catalog statistics...")
        response = requests.get(f"{API_BASE}/catalog/stats/catalog")
        if response.status_code == 200:
            self.success("Catalog statistics retrieved")
            data = response.json()
            stats = data.get('data', {})
            self.log(f"Stats: {json.dumps(stats, indent=2)}")
        else:
            self.error(f"Failed to get catalog statistics: {response.status_code}")
        
        # Get category stats
        if self.test_category_id:
            self.log("Getting category statistics...")
            response = requests.get(f"{API_BASE}/catalog/stats/categories/{self.test_category_id}")
            if response.status_code == 200:
                self.success("Category statistics retrieved")
                data = response.json()
                stats = data.get('data', {})
                self.log(f"Category stats: {json.dumps(stats, indent=2)}")
            else:
                self.error(f"Failed to get category statistics: {response.status_code}")
        
        return True
    
    def cleanup(self):
        """Clean up test data"""
        self.log("\n=== Cleaning Up Test Data ===")
        
        # Logout
        if self.token:
            response = requests.post(f"{API_BASE}/auth/logout", headers=self.get_headers())
            if response.status_code == 200:
                self.success("Logged out test user")
            else:
                self.warning("Failed to logout")
        
        self.success("Cleanup completed")
    
    def run_all_tests(self):
        """Run all catalog management tests"""
        self.log("Starting Comprehensive Blytz.live Catalog Management Tests")
        self.log("=" * 60)
        
        tests_passed = 0
        tests_failed = 0
        
        # Setup
        if not self.setup_test_user():
            return False
        
        # Test all features
        if self.test_categories():
            tests_passed += 1
        else:
            tests_failed += 1
        
        if self.create_test_product():
            if self.test_variants():
                tests_passed += 1
            else:
                tests_failed += 1
        else:
            tests_failed += 1
        
        if self.test_collections():
            tests_passed += 1
        else:
            tests_failed += 1
        
        if self.test_inventory():
            tests_passed += 1
        else:
            tests_failed += 1
        
        if self.test_search_and_discovery():
            tests_passed += 1
        else:
            tests_failed += 1
        
        if self.test_statistics():
            tests_passed += 1
        else:
            tests_failed += 1
        
        # Cleanup
        self.cleanup()
        
        # Summary
        self.log("\n" + "=" * 60)
        self.log(f"Test Summary: {tests_passed} passed, {tests_failed} failed")
        self.log("=" * 60)
        
        if tests_failed == 0:
            self.success("All catalog management tests passed! 🎉")
            return True
        else:
            self.error(f"{tests_failed} test(s) failed. Check the output above.")
            return False

# Run the tests
if __name__ == "__main__":
    tester = CatalogTest()
    success = tester.run_all_tests()
    exit(0 if success else 1)