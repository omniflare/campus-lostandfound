#!/bin/bash

# Test token extraction
test_response='{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test.signature"}'
echo "Test response: $test_response"

# Using grep method
token_grep=$(echo $test_response | grep -o '"token":"[^"]*' | cut -d'"' -f4)
echo "Token extracted with grep: $token_grep"

# Using sed method
token_sed=$(echo "$test_response" | sed 's/.*"token":"\([^"]*\)".*/\1/')
echo "Token extracted with sed: $token_sed"
