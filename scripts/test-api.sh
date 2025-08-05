#!/bin/bash

# API Testing Script for Finance Manager

BASE_URL="http://localhost:8080/api/v1"
CONTENT_TYPE="Content-Type: application/json"

echo "=== Finance Manager API Test Script ==="
echo

# Test 1: Health Check
echo "1. Testing Health Check..."
curl -s "$BASE_URL/../health" | jq .
echo

# Test 2: Register a new user
echo "2. Registering a new user..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "$CONTENT_TYPE" \
  -d '{
    "email": "test@example.com",
    "name": "Test User",
    "password": "password123",
    "role": "user"
  }')

echo "$REGISTER_RESPONSE" | jq .
TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.data.token // empty')
echo

# Test 3: Login
echo "3. Testing Login..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "$CONTENT_TYPE" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }')

echo "$LOGIN_RESPONSE" | jq .
if [ -z "$TOKEN" ]; then
  TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.token // empty')
fi
echo

if [ -n "$TOKEN" ] && [ "$TOKEN" != "null" ]; then
  # Test 4: Get Profile
  echo "4. Getting User Profile..."
  curl -s -X GET "$BASE_URL/auth/profile" \
    -H "$CONTENT_TYPE" \
    -H "Authorization: Bearer $TOKEN" | jq .
  echo

  # Test 5: Get Users
  echo "5. Getting All Users..."
  curl -s -X GET "$BASE_URL/users?page=1&limit=10" \
    -H "$CONTENT_TYPE" \
    -H "Authorization: Bearer $TOKEN" | jq .
  echo

else
  echo "No token available, skipping authenticated endpoints"
fi

echo "=== Test Complete ==="
echo
echo "Visit Swagger UI at: http://localhost:8080/swagger/index.html"
