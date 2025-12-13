#!/bin/bash

echo "🤖 Testing Grok AI Integration"
echo "================================"

BASE_URL="http://localhost:8080/api/v1"

echo ""
echo "1️⃣ Testing AI Connection..."
curl -s "$BASE_URL/ai/test" | jq '.'

echo ""
echo "2️⃣ Getting AI Stats..."
curl -s "$BASE_URL/ai/stats" | jq '.'

echo ""
echo "3️⃣ Analyzing BTC Sentiment..."
curl -s -X POST "$BASE_URL/ai/sentiment" \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "signal_type": "BUY",
    "current_price": 95000,
    "strength": 75
  }' | jq '.'

echo ""
echo "4️⃣ Toggling AI Filter OFF..."
curl -s -X POST "$BASE_URL/ai/toggle" \
  -H "Content-Type: application/json" \
  -d '{"enabled": false}' | jq '.'

echo ""
echo "5️⃣ Toggling AI Filter ON..."
curl -s -X POST "$BASE_URL/ai/toggle" \
  -H "Content-Type: application/json" \
  -d '{"enabled": true}' | jq '.'

echo ""
echo "6️⃣ Getting Updated AI Stats..."
curl -s "$BASE_URL/ai/stats" | jq '.'

echo ""
echo "✅ Grok AI Integration Test Complete!"
