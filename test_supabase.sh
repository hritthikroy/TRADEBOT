#!/bin/bash

echo "🔍 SUPABASE CONNECTION TEST"
echo "============================"
echo ""

# Check if .env file exists
if [ ! -f .env ]; then
    echo "⚠️  No .env file found"
    echo ""
    echo "To configure Supabase:"
    echo "1. Copy template: cp .env.example .env"
    echo "2. Add your Supabase credentials:"
    echo "   SUPABASE_URL=https://your-project.supabase.co"
    echo "   SUPABASE_KEY=your-anon-key"
    echo ""
    exit 1
fi

# Load environment variables
source .env 2>/dev/null || true

# Check if Supabase is configured
if [ -z "$SUPABASE_URL" ] || [ -z "$SUPABASE_KEY" ]; then
    echo "❌ Supabase not configured in .env"
    echo ""
    echo "Add these to your .env file:"
    echo "SUPABASE_URL=https://your-project.supabase.co"
    echo "SUPABASE_KEY=your-anon-key"
    echo ""
    echo "Get your credentials from:"
    echo "https://app.supabase.com/project/_/settings/api"
    echo ""
    exit 1
fi

echo "✅ Environment variables found"
echo "   URL: $SUPABASE_URL"
echo "   Key: ${SUPABASE_KEY:0:20}..."
echo ""

# Test 1: Health Check
echo "📡 Test 1: Health Check"
echo "----------------------"
HEALTH_RESPONSE=$(curl -s -w "\n%{http_code}" \
    -H "apikey: $SUPABASE_KEY" \
    -H "Authorization: Bearer $SUPABASE_KEY" \
    "$SUPABASE_URL/rest/v1/")

HTTP_CODE=$(echo "$HEALTH_RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$HEALTH_RESPONSE" | head -n-1)

if [ "$HTTP_CODE" = "200" ]; then
    echo "✅ Supabase is reachable"
else
    echo "❌ Supabase connection failed (HTTP $HTTP_CODE)"
    echo "Response: $RESPONSE_BODY"
    exit 1
fi
echo ""

# Test 2: Check Tables
echo "📊 Test 2: Check Tables"
echo "----------------------"

# Check user_settings table
echo "Checking user_settings table..."
SETTINGS_RESPONSE=$(curl -s -w "\n%{http_code}" \
    -H "apikey: $SUPABASE_KEY" \
    -H "Authorization: Bearer $SUPABASE_KEY" \
    "$SUPABASE_URL/rest/v1/user_settings?select=*&limit=1")

HTTP_CODE=$(echo "$SETTINGS_RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$SETTINGS_RESPONSE" | head -n-1)

if [ "$HTTP_CODE" = "200" ]; then
    echo "✅ user_settings table exists"
    echo "   Data: $RESPONSE_BODY"
else
    echo "⚠️  user_settings table not found or not accessible"
    echo "   You may need to create it. See database_setup.sql"
fi
echo ""

# Check signals table
echo "Checking signals table..."
SIGNALS_RESPONSE=$(curl -s -w "\n%{http_code}" \
    -H "apikey: $SUPABASE_KEY" \
    -H "Authorization: Bearer $SUPABASE_KEY" \
    "$SUPABASE_URL/rest/v1/signals?select=*&limit=1")

HTTP_CODE=$(echo "$SIGNALS_RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$SIGNALS_RESPONSE" | head -n-1)

if [ "$HTTP_CODE" = "200" ]; then
    echo "✅ signals table exists"
    SIGNAL_COUNT=$(echo "$RESPONSE_BODY" | jq '. | length' 2>/dev/null || echo "0")
    echo "   Signals found: $SIGNAL_COUNT"
else
    echo "⚠️  signals table not found or not accessible"
    echo "   You may need to create it. See database_setup.sql"
fi
echo ""

# Test 3: Write Test
echo "✏️  Test 3: Write Test"
echo "----------------------"
echo "Testing write permissions..."

TEST_DATA='{"id":999,"filter_buy":true,"filter_sell":true,"min_confluence":5,"max_risk_percent":2.0,"enabled_strategies":["test"]}'

WRITE_RESPONSE=$(curl -s -w "\n%{http_code}" \
    -X POST \
    -H "apikey: $SUPABASE_KEY" \
    -H "Authorization: Bearer $SUPABASE_KEY" \
    -H "Content-Type: application/json" \
    -H "Prefer: return=minimal" \
    -d "$TEST_DATA" \
    "$SUPABASE_URL/rest/v1/user_settings")

HTTP_CODE=$(echo "$WRITE_RESPONSE" | tail -n1)

if [ "$HTTP_CODE" = "201" ] || [ "$HTTP_CODE" = "200" ]; then
    echo "✅ Write permissions working"
    
    # Clean up test data
    curl -s -X DELETE \
        -H "apikey: $SUPABASE_KEY" \
        -H "Authorization: Bearer $SUPABASE_KEY" \
        "$SUPABASE_URL/rest/v1/user_settings?id=eq.999" > /dev/null
    
    echo "   (Test data cleaned up)"
else
    echo "⚠️  Write permissions may be restricted (HTTP $HTTP_CODE)"
    echo "   This is OK if you only need read access"
fi
echo ""

# Test 4: Backend Integration
echo "🔗 Test 4: Backend Integration"
echo "------------------------------"

# Check if backend is running
if curl -s http://localhost:8080/api/v1/health > /dev/null 2>&1; then
    echo "✅ Backend is running"
    
    # Test user settings endpoint
    echo "Testing /api/v1/settings endpoint..."
    BACKEND_RESPONSE=$(curl -s http://localhost:8080/api/v1/settings)
    
    if echo "$BACKEND_RESPONSE" | jq . > /dev/null 2>&1; then
        echo "✅ Settings endpoint working"
        echo "   Response: $BACKEND_RESPONSE" | head -c 100
        echo "..."
    else
        echo "⚠️  Settings endpoint returned invalid JSON"
    fi
else
    echo "⚠️  Backend not running"
    echo "   Start with: cd backend && go run ."
fi
echo ""

# Summary
echo "📋 SUMMARY"
echo "=========="
echo ""
echo "Supabase Status:"
echo "  Connection: ✅ Working"
echo "  URL: $SUPABASE_URL"
echo ""
echo "Tables:"
echo "  user_settings: Check above"
echo "  signals: Check above"
echo ""
echo "Next Steps:"
echo "1. If tables don't exist, run: psql < database_setup.sql"
echo "2. Or create tables in Supabase SQL Editor"
echo "3. Start backend: cd backend && go run ."
echo "4. Test endpoints: curl http://localhost:8080/api/v1/settings"
echo ""
echo "Documentation:"
echo "  Supabase Dashboard: https://app.supabase.com"
echo "  API Docs: $SUPABASE_URL/rest/v1/"
echo ""
