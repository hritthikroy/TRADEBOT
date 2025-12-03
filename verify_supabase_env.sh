#!/bin/bash

echo "🔍 Verifying Supabase Environment Variables"
echo "==========================================="
echo ""

# Check if .env file exists
if [ ! -f "backend/.env" ]; then
    echo "❌ backend/.env file not found!"
    exit 1
fi

echo "✅ Found backend/.env file"
echo ""

# Load and check SUPABASE_URL
SUPABASE_URL=$(grep "^SUPABASE_URL=" backend/.env | cut -d '=' -f2)
if [ -z "$SUPABASE_URL" ]; then
    echo "❌ SUPABASE_URL is empty or not set!"
else
    echo "✅ SUPABASE_URL: $SUPABASE_URL"
fi

# Load and check SUPABASE_KEY
SUPABASE_KEY=$(grep "^SUPABASE_KEY=" backend/.env | cut -d '=' -f2)
if [ -z "$SUPABASE_KEY" ]; then
    echo "❌ SUPABASE_KEY is empty or not set!"
else
    echo "✅ SUPABASE_KEY: ${SUPABASE_KEY:0:30}..."
fi

echo ""
echo "Testing connection to Supabase..."
echo ""

# Test connection
if [ -n "$SUPABASE_URL" ] && [ -n "$SUPABASE_KEY" ]; then
    RESPONSE=$(curl -s -w "\n%{http_code}" "$SUPABASE_URL/rest/v1/trading_signals?limit=1" \
      -H "apikey: $SUPABASE_KEY" \
      -H "Authorization: Bearer $SUPABASE_KEY")
    
    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    BODY=$(echo "$RESPONSE" | head -n-1)
    
    if [ "$HTTP_CODE" = "200" ]; then
        echo "✅ Connection successful! (HTTP 200)"
        echo "Response: $BODY"
    else
        echo "❌ Connection failed! (HTTP $HTTP_CODE)"
        echo "Response: $BODY"
    fi
else
    echo "❌ Cannot test connection - missing URL or KEY"
fi

echo ""
echo "==========================================="
echo ""
echo "Summary:"
if [ -n "$SUPABASE_URL" ] && [ -n "$SUPABASE_KEY" ] && [ "$HTTP_CODE" = "200" ]; then
    echo "✅ Everything looks good!"
    echo ""
    echo "If signals still don't save, check backend logs for:"
    echo "  - '🔍 Supabase URL from env:' message"
    echo "  - '❌ Supabase error' messages"
    echo "  - Any error details"
else
    echo "❌ Issues found! Fix the problems above."
fi
echo ""
