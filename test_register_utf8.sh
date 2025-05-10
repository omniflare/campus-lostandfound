#!/bin/bash

echo "Testing user registration with explicit Content-Type charset..."

curl -v -X POST \
  http://localhost:3000/api/v1/auth/register \
  -H "Content-Type: application/json; charset=utf-8" \
  -d '{
  "username": "testuser3",
  "email": "test3@example.com",
  "password": "password123",
  "first_name": "Test",
  "last_name": "User",
  "phone": "1234567890"
}'
