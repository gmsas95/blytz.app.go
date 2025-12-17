#!/usr/bin/env python3

import requests
import json

API_BASE = "http://localhost:8080/api/v1"

# Test login response structure
login_data = {
    "email": "catalogtest@example.com",
    "password": "testpassword123"
}

print("Testing login...")
response = requests.post(f"{API_BASE}/auth/login", json=login_data)
print(f"Status: {response.status_code}")
print(f"Response: {response.text}")

if response.status_code == 200:
    data = response.json()
    print(f"JSON keys: {data.keys()}")
    if 'data' in data:
        print(f"Data keys: {data['data'].keys()}")