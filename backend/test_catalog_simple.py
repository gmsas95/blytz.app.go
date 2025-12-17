#!/usr/bin/env python3

import requests
import json
import uuid
from datetime import datetime

API_BASE = "http://localhost:8080/api/v1"

# Test basic catalog endpoints
def test_catalog():
    print("Testing Blytz.live Catalog Management")
    print("=" * 40)
    
    # Test health endpoint
    try:
        response = requests.get(f"http://localhost:8080/health")
        if response.status_code == 200:
            print("✅ Server is healthy")
            print(f"Response: {response.json()}")
        else:
            print("❌ Server health check failed")
    except Exception as e:
        print(f"❌ Server connection error: {e}")
        return
    
    # Test search products
    try:
        response = requests.get(f"{API_BASE}/catalog/search/products?q=phone")
        if response.status_code == 200:
            data = response.json()
            print("✅ Product search works")
            print(f"Success: {data.get('success')}")
            if 'data' in data:
                print(f"Products found: {len(data['data'])}")
        else:
            print("❌ Product search failed")
            print(f"Status: {response.status_code}")
            print(f"Response: {response.text}")
    except Exception as e:
        print(f"❌ Search error: {e}")
    
    # Test categories
    try:
        response = requests.get(f"{API_BASE}/catalog/categories")
        if response.status_code == 200:
            data = response.json()
            print("✅ Categories endpoint works")
            print(f"Success: {data.get('success')}")
            if 'data' in data:
                print(f"Categories found: {len(data['data'])}")
        else:
            print("❌ Categories endpoint failed")
    except Exception as e:
        print(f"❌ Categories error: {e}")
    
    # Test catalog stats
    try:
        response = requests.get(f"{API_BASE}/catalog/stats/catalog")
        if response.status_code == 200:
            data = response.json()
            print("✅ Catalog stats works")
            print(f"Stats: {data}")
        else:
            print("❌ Catalog stats failed")
    except Exception as e:
        print(f"❌ Stats error: {e}")

    print("\n" + "=" * 40)
    print("Catalog tests completed!")

if __name__ == "__main__":
    test_catalog()