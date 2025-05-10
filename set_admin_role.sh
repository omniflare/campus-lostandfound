#!/bin/bash

# This script updates an existing user to have admin role
# Usage: ./set_admin_role.sh <username>

# Check if a username was provided
if [ -z "$1" ]; then
  echo "Usage: ./set_admin_role.sh <username>"
  exit 1
fi

USERNAME="$1"

# Get the DATABASE_URL from environment or use default Neon Tech URL
DATABASE_URL=${DATABASE_URL:-"postgresql://neondb_owner:npg_MVd1DFk6rxut@ep-wispy-rice-a4b0w7v3-pooler.us-east-1.aws.neon.tech/neondb?sslmode=require"}

echo "Attempting to update role for user: $USERNAME to 'admin'"

# Use psql to update the role
PGPASSWORD=$(echo $DATABASE_URL | sed -n 's/.*:\/\/[^:]*:\([^@]*\)@.*/\1/p') \
psql "$DATABASE_URL" -c "UPDATE users SET role = 'admin' WHERE username = '$USERNAME'"

# Check if the update was successful
if [ $? -eq 0 ]; then
  echo "Successfully updated role for $USERNAME to admin"
else
  echo "Failed to update role. Make sure psql is installed and the database is accessible."
fi
