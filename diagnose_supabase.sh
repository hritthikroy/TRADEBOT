#!/bin/bash

echo "🔍 SUPABASE DIAGNOSTIC TOOL"
echo "============================"
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Load environment variables
if [ -f backend/.env ]; then
    export $(cat backend/.env | grep -v '^#' | xargs)
else
    echo -e "${RED}❌ backend/.env file not found${NC}"
    exit 1
fi

echo "Step 1: Checking Environment Variables"
echo "---------------------------------------"
if [ -z "$SUPABASE_URL" ]; then
    echo -e "${RED}❌ SUPABASE_URL is not set${NC}"
    exit 1
else
    echo -e "${GREEN}✅ SUPABASE_URL: $SUPABASE_URL${NC}"
fi

if [ -z "$SUPABASE_KEY" ]; then
    echo -e "${RED}❌ SUPABASE_KEY is not set${NC}"
    exit 1
else
    echo -e "${GREEN}✅ SUPABASE_KEY: ${SUPABASE_KEY:0:30}...${NC}"
fi

echo ""
echo "Step 2: Testing Supabase Connection"
echo "------------------------------------"

# Test connection
RESPONSE=$(curl -s -w "\n%{http_code}" "$SUPABASE_URL/rest/v1/" \
  -H "apikey: $SUPABASE_KEY" \
  -H "Authorization: Bearer $SUPABASE_KEY")

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | head -n-1)

if [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✅ Supabase connection successful${NC}"
else
    echo -e "${RED}❌ Supabase connection failed (HTTP $HTTP_CODE)${NC}"
    echo "Response: $BODY"
    exit 1
fi

echo ""
echo "Step 3: Checking if trading_signals table exists"
echo "-------------------------------------------------"

RESPONSE=$(curl -s -w "\n%{http_code}" "$SUPABASE_URL/rest/v1/trading_signals?limit=1" \
  -H "apikey: $SUPABASE_KEY" \
  -H "Authorization: Bearer $SUPABASE_KEY")

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | head -n-1)

if [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✅ Table 'trading_signals' exists${NC}"
    
    # Check if table has data
    COUNT=$(echo "$BODY" | jq '. | length')
    if [ "$COUNT" -gt 0 ]; then
        echo -e "${GREEN}✅ Table has $COUNT signal(s)${NC}"
    else
        echo -e "${YELLOW}⚠️  Table exists but is empty${NC}"
    fi
else
    echo -e "${RED}❌ Table 'trading_signals' does NOT exist${NC}"
    echo ""
    echo "🔧 FIX: Run this SQL in Supabase SQL Editor:"
    echo "   1. Go to: https://supabase.com/dashboard/project/elqhqhjevajzjoghiiss/sql"
    echo "   2. Copy contents of: supabase-setup.sql"
    echo "   3. Paste and click 'Run'"
    echo ""
    exit 1
fi

echo ""
echo "Step 4: Testing INSERT permission"
echo "----------------------------------"

TEST_SIGNAL='{
  "symbol": "BTCUSDT",
  "strategy": "diagnostic_test",
  "signal_type": "BUY",
  "entry_price": 50000.00,
  "stop_loss": 49500.00,
  "take_profit": 51000.00,
  "current_price": 50000.00,
  "risk_reward": 2.0,
  "status": "ACTIVE",
  "progress": 0,
  "filter_buy": true,
  "filter_sell": true,
  "signal_time": "'$(date -u +"%Y-%m-%dT%H:%M:%SZ")'"
}'

RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$SUPABASE_URL/rest/v1/trading_signals" \
  -H "apikey: $SUPABASE_KEY" \
  -H "Authorization: Bearer $SUPABASE_KEY" \
  -H "Content-Type: application/json" \
  -H "Prefer: return=representation" \
  -d "$TEST_SIGNAL")

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | head -n-1)

if [ "$HTTP_CODE" = "201" ]; then
    echo -e "${GREEN}✅ INSERT permission works!${NC}"
    echo "Test signal created successfully"
    
    # Clean up test signal
    TEST_ID=$(echo "$BODY" | jq -r '.[0].id')
    if [ "$TEST_ID" != "null" ]; then
        curl -s -X DELETE "$SUPABASE_URL/rest/v1/trading_signals?id=eq.$TEST_ID" \
          -H "apikey: $SUPABASE_KEY" \
          -H "Authorization: Bearer $SUPABASE_KEY" > /dev/null
        echo "Test signal cleaned up"
    fi
else
    echo -e "${RED}❌ INSERT permission failed (HTTP $HTTP_CODE)${NC}"
    echo "Error: $BODY"
    echo ""
    echo "🔧 FIX: Check RLS policies in Supabase"
    echo "   The supabase-setup.sql should have created the policies"
    echo "   You may need to re-run it"
    exit 1
fi

echo ""
echo "Step 5: Testing SELECT permission"
echo "----------------------------------"

RESPONSE=$(curl -s -w "\n%{http_code}" "$SUPABASE_URL/rest/v1/trading_signals?limit=5&order=created_at.desc" \
  -H "apikey: $SUPABASE_KEY" \
  -H "Authorization: Bearer $SUPABASE_KEY")

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)

if [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✅ SELECT permission works!${NC}"
else
    echo -e "${RED}❌ SELECT permission failed (HTTP $HTTP_CODE)${NC}"
    exit 1
fi

echo ""
echo "============================"
echo -e "${GREEN}🎉 ALL CHECKS PASSED!${NC}"
echo "============================"
echo ""
echo "Your Supabase is configured correctly!"
echo ""
echo "Next steps:"
echo "1. Start your backend: cd backend && go run ."
echo "2. Generate a signal from the UI"
echo "3. Check Supabase Table Editor to see the signal"
echo ""
