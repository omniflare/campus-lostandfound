#!/bin/bash

# Test script for Campus Lost and Found API
# This script will test all API endpoints and log results

# Set the base URL
BASE_URL="http://localhost:3000"

# Create tests directory if it doesn't exist
mkdir -p tests

# Create or clear log files
> tests/success.txt
> tests/error.txt
> tests/success_details.txt
> tests/error_details.txt

echo "Starting API tests for Campus Lost and Found..."

# Helper function to run tests
run_test() {
    method=$1
    endpoint=$2
    payload=$3
    token=$4
    description=$5
    
    echo "Testing $description ($method $endpoint)"
    
    # Prepare curl command
    if [ "$method" = "GET" ]; then
        if [ -z "$token" ]; then
            result=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$endpoint")
        else
            result=$(curl -s -w "\n%{http_code}" -X $method -H "Authorization: Bearer $token" "$BASE_URL$endpoint")
        fi
    else
        if [ -z "$token" ]; then
            result=$(curl -s -w "\n%{http_code}" -X $method -H "Content-Type: application/json" -d "$payload" "$BASE_URL$endpoint")
        else
            result=$(curl -s -w "\n%{http_code}" -X $method -H "Content-Type: application/json" -H "Authorization: Bearer $token" -d "$payload" "$BASE_URL$endpoint")
        fi
    fi
    
    # Extract status code from response
    http_code=$(echo "$result" | tail -n1)
    response=$(echo "$result" | sed '$d')
    
    # Check if status code indicates success (2xx)
    if [[ $http_code -ge 200 && $http_code -lt 300 ]]; then
        echo "$description: SUCCESS ($http_code)" >> tests/success.txt
        echo "$description: $response" >> tests/success_details.txt
    else
        echo "$description: FAILED ($http_code)" >> tests/error.txt
        echo "$description: $response" >> tests/error_details.txt
    fi
}

# Test health check endpoint
run_test "GET" "/health" "" "" "Health Check"

# Test public endpoints
run_test "GET" "/api/v1/items" "" "" "Get All Items (Public)"
run_test "GET" "/api/v1/items/search?q=phone" "" "" "Search Items (Public)"
run_test "GET" "/api/v1/items/1" "" "" "Get Item Details (Public)"

# First register a test user
echo "Registering test user..."
register_payload='{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123",
  "first_name": "Test",
  "last_name": "User",
  "phone": "1234567890"
}'
register_response=$(curl -s -X POST -H "Content-Type: application/json" -d "$register_payload" "$BASE_URL/api/v1/auth/register")
echo "Register response: $register_response"

# Register an admin user
echo "Registering admin user..."
admin_register_payload='{
  "username": "admin",
  "email": "admin@example.com",
  "password": "admin123",
  "first_name": "Admin",
  "last_name": "User",
  "phone": "9876543210"
}'
admin_register_response=$(curl -s -X POST -H "Content-Type: application/json" -d "$admin_register_payload" "$BASE_URL/api/v1/auth/register")
echo "Admin register response: $admin_register_response"

# Login to get token
echo "Logging in to get authentication token..."
login_payload='{
  "username": "testuser",
  "password": "password123"
}'
login_response=$(curl -s -X POST -H "Content-Type: application/json" -d "$login_payload" "$BASE_URL/api/v1/auth/login")
echo "Login response: $login_response"

# Extract token from login response - improved parsing
if [[ "$login_response" == *"token"* ]]; then
  token=$(echo "$login_response" | sed 's/.*"token":"\([^"]*\)".*/\1/')
  echo "Auth token: $token"
else
  echo "Could not find token in response"
  token=""
fi

# Check if we got a valid token
if [ -z "$token" ] || [ "$token" == "$login_response" ]; then
    echo "Failed to extract authentication token. Skipping authenticated tests."
    exit 1
fi

# Login admin to get admin token
echo "Logging in as admin to get admin token..."
admin_login_payload='{
  "username": "admin",
  "password": "admin123"
}'
admin_login_response=$(curl -s -X POST -H "Content-Type: application/json" -d "$admin_login_payload" "$BASE_URL/api/v1/auth/login")
echo "Admin login response: $admin_login_response"

# Extract admin token from login response - improved parsing
if [[ "$admin_login_response" == *"token"* ]]; then
  admin_token=$(echo "$admin_login_response" | sed 's/.*"token":"\([^"]*\)".*/\1/')
  echo "Admin token: $admin_token"
else
  echo "Could not find admin token in response"
  admin_token=""
fi

# Check if admin token is valid
if [ -z "$admin_token" ] || [ "$admin_token" == "$admin_login_response" ]; then
  echo "Failed to extract admin token. Will skip admin tests."
else
  echo "Successfully obtained admin token."
fi

# Test authentication endpoints
run_test "POST" "/api/v1/auth/register" '{"username":"newuser","email":"new@example.com","password":"password123","first_name":"New","last_name":"User","phone":"0987654321"}' "" "Register New User"
run_test "POST" "/api/v1/auth/login" '{"username":"testuser","password":"password123"}' "" "Login User"

# Test authenticated user endpoints
echo "Testing authenticated endpoints with token: $token"
run_test "GET" "/api/v1/user/profile" "" "$token" "Get User Profile"
run_test "PUT" "/api/v1/user/profile" '{"first_name":"Updated","last_name":"User","email":"test@example.com","phone":"1234567890"}' "$token" "Update User Profile"
run_test "PUT" "/api/v1/user/password" '{"current_password":"password123","new_password":"newpassword123"}' "$token" "Change Password"
run_test "GET" "/api/v1/user/items" "" "$token" "Get User Items"
run_test "GET" "/api/v1/user/messages/unread" "" "$token" "Get Unread Message Count"
run_test "GET" "/api/v1/user/messages/conversations" "" "$token" "Get Conversations"

# Test item endpoints
run_test "POST" "/api/v1/items/lost" '{"title":"Lost Laptop","description":"MacBook Pro 16-inch","category":"Electronics","location":"Library","lost_time":"2023-05-10T15:00:00Z"}' "$token" "Report Lost Item"
run_test "POST" "/api/v1/items/found" '{"title":"Found Phone","description":"iPhone 13 Pro","category":"Electronics","location":"Cafeteria"}' "$token" "Report Found Item"
run_test "PUT" "/api/v1/items/1/status" '{"status":"claimed"}' "$token" "Update Item Status"

# Note: Image upload requires multipart/form-data which is more complex
# You would need curl -F "image=@path/to/image.jpg" which we're skipping in this example

# Test messaging endpoints
run_test "POST" "/api/v1/user/messages" '{"receiver_id":2,"item_id":1,"content":"Hi, I think I found your lost item"}' "$token" "Send Message"
run_test "GET" "/api/v1/user/messages/2" "" "$token" "Get Messages with User"

# Test reporting
run_test "POST" "/api/v1/user/reports" '{"reported_id":2,"reason":"Suspicious behavior"}' "$token" "Create Report"

# Update user role to admin (requires manual database update or separate functionality)
# For testing purposes, we'd need to manually update the role in the database

# Test guard endpoints (would need guard role)
# run_test "GET" "/api/v1/guard/items" "" "$token" "Get All Items (Guard)"

# Test admin endpoints
if [ ! -z "$admin_token" ] && [ "$admin_token" != "$admin_login_response" ]; then
    echo "Testing admin endpoints with admin token: $admin_token"
    
    run_test "GET" "/api/v1/admin/users" "" "$admin_token" "Get All Users (Admin)"
    run_test "PUT" "/api/v1/admin/users/2/role" '{"role":"guard"}' "$admin_token" "Update User Role (Admin)"
    run_test "GET" "/api/v1/admin/reports" "" "$admin_token" "Get All Reports (Admin)"
    run_test "PUT" "/api/v1/admin/reports/1/status" '{"status":"resolved","comment":"No issues found"}' "$admin_token" "Update Report Status (Admin)"
    run_test "GET" "/api/v1/admin/stats" "" "$admin_token" "Get Stats (Admin)"
else
    echo "No valid admin token available. Skipping admin tests."
fi

echo "API tests completed. Check tests/success.txt and tests/error.txt for results."
echo "Detailed responses are in tests/success_details.txt and tests/error_details.txt"
