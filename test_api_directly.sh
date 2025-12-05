#!/bin/bash

echo "🧪 Testing API Directly - Does Backend Return Different Data?"
echo "================================================================"
echo ""

# Test 5 days
echo "📊 Test 1: Requesting 5 days..."
result5=$(curl -s -X POST http://localhost:8080/api/v1/backtest/test-all-strategies \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTCUSDT","days":5,"startBalance":10000,"filterBuy":true,"filterSell":true}' \
  | jq -r '.results[] | select(.strategyName == "session_trader") | .totalTrades')

echo "✅ 5 days result: $result5 trades"
echo ""

# Wait a moment
sleep 1

# Test 10 days
echo "📊 Test 2: Requesting 10 days..."
result10=$(curl -s -X POST http://localhost:8080/api/v1/backtest/test-all-strategies \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTCUSDT","days":10,"startBalance":10000,"filterBuy":true,"filterSell":true}' \
  | jq -r '.results[] | select(.strategyName == "session_trader") | .totalTrades')

echo "✅ 10 days result: $result10 trades"
echo ""

# Wait a moment
sleep 1

# Test 15 days
echo "📊 Test 3: Requesting 15 days..."
result15=$(curl -s -X POST http://localhost:8080/api/v1/backtest/test-all-strategies \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTCUSDT","days":15,"startBalance":10000,"filterBuy":true,"filterSell":true}' \
  | jq -r '.results[] | select(.strategyName == "session_trader") | .totalTrades')

echo "✅ 15 days result: $result15 trades"
echo ""

# Compare results
echo "================================================================"
echo "📈 COMPARISON:"
echo "================================================================"
echo "5 days:  $result5 trades"
echo "10 days: $result10 trades"
echo "15 days: $result15 trades"
echo ""

if [ "$result5" != "$result10" ]; then
    echo "✅ SUCCESS: 5 days and 10 days are DIFFERENT!"
    echo "   Backend is working correctly."
else
    echo "⚠️  WARNING: 5 days and 10 days are the SAME"
    echo "   This might indicate an issue."
fi

if [ "$result10" == "$result15" ]; then
    echo "ℹ️  NOTE: 10 days and 15 days are the same"
    echo "   This is expected due to Binance 1000 candle limit."
fi

echo ""
echo "================================================================"
echo "🎯 CONCLUSION:"
echo "================================================================"

if [ "$result5" != "$result10" ]; then
    echo "✅ Backend API is working correctly!"
    echo "✅ Different day values return different results"
    echo ""
    echo "If your browser shows the same results, the issue is:"
    echo "  1. Browser cache (force refresh with Ctrl+Shift+R)"
    echo "  2. JavaScript not executing properly"
    echo "  3. Display not updating (but data is changing)"
    echo ""
    echo "Try opening SIMPLE_TEST.html to test the frontend."
else
    echo "❌ Backend might have an issue"
    echo "   Both requests returned the same number of trades"
fi

echo ""
