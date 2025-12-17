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
print(f"Available categories: {len(categories)}")
for i, cat in enumerate(categories):
    print(f"  {i}: {cat['name']} ({cat['id']})")

test_category_id = categories[1]['id'] if len(categories) > 1 else categories[0]['id']

print(f"Using category ID: {test_category_id}")

# Create category attribute
attribute_data = {
    "name": "Warranty Period",
    "type": "select",
    "required": False,
    "options": ["1 Year", "2 Years", "3 Years"],
    "default_value": "1 Year"
}

print("Creating category attribute...")
print(f"Request data: {json.dumps(attribute_data, indent=2)}")

response = requests.post(f"{API_BASE}/catalog/categories/{test_category_id}/attributes",
                      json=attribute_data, headers=headers)

print(f"Status: {response.status_code}")
print(f"Response: {response.text}")