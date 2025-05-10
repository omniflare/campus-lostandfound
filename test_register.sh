#!/bin/bash

echo "Testing user registration with correct Content-Type header..."

curl -v -X POST \
  http://localhost:3000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
  "username": "testuser2",
  "email": "test2@example.com",
  "password": "password123",
  "first_name": "Test",
  "last_name": "User",
  "phone": "1234567890"
}'

echo -e "\n\nTesting login with the same user..."

curl -v -X POST \
  http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
  "username": "testuser2",
  "password": "password123"
}'
