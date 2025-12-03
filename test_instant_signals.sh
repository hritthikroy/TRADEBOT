#!/bin/bash

echo "🧪 Testing Instant Signal System"
echo "================================"
echo ""

# Test 1: Generate a live signal
echo "📡 Test 1: Generating live signal..."
curl -X POST http://localhost:8080/api/v1/live-signal \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "strategy": "session_trader"
  }' | jq '.'

echo ""
echo "✅ Signal generated! Check:"
echo "   1. Your Telegram for instant notification"
echo "   2. Supabase dashboard for saved signal"
echo ""

# Test 2: Check Telegram bot status
echo "📱 Test 2: Checking Telegram bot status..."
curl -s http://localhost:8080/api/v1/telegram/status | jq '.'

echo ""

# Test 3: Get recent signals from Supabase
echo "💾 Test 3: Fetching recent signals from Supabase..."
curl -s http://localhost:8080/api/v1/signals/recent?limit=5 | jq '.'

echo ""
echo "================================"
echo "✅ All tests complete!"
echo ""
echo "Expected Results:"
echo "  ✅ Signal generated and returned"
echo "  ✅ Telegram bot status: running=true"
echo "  ✅ Recent signals retrieved from Supabase"
echo "  ✅ Telegram message received (check your phone)"
echo ""
