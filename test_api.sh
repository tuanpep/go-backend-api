#!/bin/bash

# Test script for Go Backend API
echo "🧪 Testing Go Backend API..."

# Test health endpoint
echo "1. Testing health endpoint..."
HEALTH_RESPONSE=$(curl -s http://localhost:8080/health)
if echo "$HEALTH_RESPONSE" | grep -q "success"; then
    echo "✅ Health check passed"
    echo "$HEALTH_RESPONSE"
else
    echo "❌ Health check failed"
    echo "$HEALTH_RESPONSE"
fi

echo -e "\n2. Testing user registration..."
REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "SecurePass123!"
  }')

if echo "$REGISTER_RESPONSE" | grep -q "success"; then
    echo "✅ Registration successful"
    echo "$REGISTER_RESPONSE"
else
    echo "❌ Registration failed"
    echo "$REGISTER_RESPONSE"
fi

echo -e "\n3. Testing user login..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecurePass123!"
  }')

if echo "$LOGIN_RESPONSE" | grep -q "access_token"; then
    echo "✅ Login successful"
    echo "$LOGIN_RESPONSE"
    
    # Extract token using basic text processing (no jq needed)
    TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
    
    if [ -n "$TOKEN" ]; then
        echo -e "\n4. Testing create post..."
        CREATE_POST_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/posts \
          -H "Content-Type: application/json" \
          -H "Authorization: Bearer $TOKEN" \
          -d '{
            "title": "Test Post",
            "content": "This is a test post created via API"
          }')
        
        if echo "$CREATE_POST_RESPONSE" | grep -q "success"; then
            echo "✅ Create post successful"
            echo "$CREATE_POST_RESPONSE"
        else
            echo "❌ Create post failed"
            echo "$CREATE_POST_RESPONSE"
        fi
        
        echo -e "\n5. Testing get posts..."
        GET_POSTS_RESPONSE=$(curl -s -X GET http://localhost:8080/api/v1/posts \
          -H "Authorization: Bearer $TOKEN")
        
        if echo "$GET_POSTS_RESPONSE" | grep -q "success"; then
            echo "✅ Get posts successful"
            echo "$GET_POSTS_RESPONSE"
        else
            echo "❌ Get posts failed"
            echo "$GET_POSTS_RESPONSE"
        fi
        
        echo -e "\n6. Testing get profile..."
        GET_PROFILE_RESPONSE=$(curl -s -X GET http://localhost:8080/api/v1/users/profile \
          -H "Authorization: Bearer $TOKEN")
        
        if echo "$GET_PROFILE_RESPONSE" | grep -q "success"; then
            echo "✅ Get profile successful"
            echo "$GET_PROFILE_RESPONSE"
        else
            echo "❌ Get profile failed"
            echo "$GET_PROFILE_RESPONSE"
        fi
    else
        echo "❌ Could not extract token from login response"
    fi
else
    echo "❌ Login failed"
    echo "$LOGIN_RESPONSE"
fi

echo -e "\n🎉 API testing completed!"