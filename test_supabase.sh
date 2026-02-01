#!/bin/bash

# Test Supabase connection
# Usage: ./test_supabase.sh [password]

if [ -z "$1" ]; then
    echo "Usage: ./test_supabase.sh [password]"
    exit 1
fi

export APP_DATABASE_HOST=aws-1-ap-southeast-1.pooler.supabase.com
export APP_DATABASE_PORT=6543
export APP_DATABASE_USER=postgres.llpnxoxeogdrqyffjtjh
export APP_DATABASE_PASSWORD=$1
export APP_DATABASE_DBNAME=postgres
export APP_DATABASE_SSLMODE=require

echo "Testing Supabase connection..."
echo "Host: $APP_DATABASE_HOST"
echo "Port: $APP_DATABASE_PORT"
echo "User: $APP_DATABASE_USER"
echo "Database: $APP_DATABASE_DBNAME"
echo ""

go run ./cmd/api/ migrate
