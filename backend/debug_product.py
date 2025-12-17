#!/usr/bin/env python3

import requests
import json

API_BASE = "http://localhost:8080/api/v1"

# Login first
login_data = {
    "email": "admin@example.com",
    "password": "adminpassword123"
}
response = requests.post(f"{API_BASE}/auth/login", json=login_data)
token = response.json().get('access_token', '')
headers = {"Authorization": f"Bearer {token}"}

# Get test category
response = requests.get(f"{API_BASE}/catalog/categories")
categories = response.json().get('data', [])
test_category_id = categories[0]['id']

# Create product
product_data = {
    "category_id": test_category_id,
    "title": "Premium Test Product",
    "description": "A high-quality test product with advanced features",
    "condition": "new",
    "starting_price": 599.99,
    "buy_now_price": 699.99,
    "images": ["https://example.com/product1.jpg"],
    "status": "active"
}

print("Creating product...")
response = requests.post(f"{API_BASE}/products", json=product_data, headers=headers)
print(f"Status: {response.status_code}")
print(f"Response: {response.text}")

if response.status_code == 201:
    data = response.json()
    print(f"Data structure: {data.keys()}")