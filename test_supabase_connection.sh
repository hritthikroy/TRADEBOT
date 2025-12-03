#!/bin/bash

echo "🔍 Testing Supabase Connection"
echo "================================"
echo ""

# Load environment variables
source backend/.env

echo "📋 Configuration:"
echo "SUPABASE_URL: $SUPABASE_URL"
echo "SUPABASE_KEY: ${SUPABASE_KEY:0:20}..."
echo ""

# Test 1: Check if table exists
echo "🧪 Test 1: Checking if trading_signals table exists..."
curl -s "$SUPABASE_URL/rest/v1/trading_signals?limit=1" \
  -H "apikey: $SUPABASE_KEY" \
  -H "Authorization: Bearer $SUPABASE_KEY" | jq '.'

echo ""
echo "================================"
echo ""

# Test 2: Try to insert a test signal
echo "🧪 Test 2: Attempting to insert test signal..."
curl -X POST "$SUPABASE_URL/rest/v1/trading_signals" \
  -H "apikey: $SUPABASE_KEY" \
  -H "Authorization: Bearer $SUPABASE_KEY" \
  -H "Content-Type: application/json" \
  -H "Prefer: return=representation" \
  -d '{
    "symbol": "BTCUSDT",
    "strategy": "test_strategy",
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
  }' | jq '.'

echo ""
echo "================================"
echo ""

# Test 3: Check what was inserted
echo "🧪 Test 3: Checking recent signals..."
curl -s "$SUPABASE_URL/rest/v1/trading_signals?order=created_at.desc&limit=5" \
  -H "apikey: $SUPABASE_KEY" \
  -H "Authorization: Bearer $SUPABASE_KEY" | jq '.'

echo ""
echo "================================"
echo "✅ Tests complete!"
echo ""
echo "If you see errors above, the table might not exist."
echo "Run the SQL in supabase-setup.sql in your Supabase SQL Editor."
