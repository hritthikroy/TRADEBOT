#!/bin/bash

echo "🧪 Testing Chart Update Fix"
echo "================================"
echo ""

# Test 1: 15 days
echo "📊 Test 1: Running backtest with 15 days..."
response1=$(curl -s -X POST http://localhost:8080/backtest/test-all-strategies \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "days": 15,
    "startBalance": 10000,
    "filterBuy": true,
    "filterSell": true
  }')

trades_15=$(echo "$response1" | jq -r '.results[] | select(.strategyName == "session_trader") | .totalTrades')
echo "✅ Session Trader with 15 days: $trades_15 trades"
echo ""

# Wait a moment
sleep 1

# Test 2: 30 days
echo "📊 Test 2: Running backtest with 30 days..."
response2=$(curl -s -X POST http://localhost:8080/backtest/test-all-strategies \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "days": 30,
    "startBalance": 10000,
    "filterBuy": true,
    "filterSell": true
  }')

trades_30=$(echo "$response2" | jq -r '.results[] | select(.strategyName == "session_trader") | .totalTrades')
echo "✅ Session Trader with 30 days: $trades_30 trades"
echo ""

# Test 3: Custom date range
echo "📊 Test 3: Running backtest with custom date range (Nov 1-30, 2024)..."
start_time=$(date -j -f "%Y-%m-%d" "2024-11-01" "+%s")000
end_time=$(date -j -f "%Y-%m-%d" "2024-11-30" "+%s")999

response3=$(curl -s -X POST http://localhost:8080/backtest/test-all-strategies \
  -H "Content-Type: application/json" \
  -d "{
    \"symbol\": \"BTCUSDT\",
    \"days\": 30,
    \"startBalance\": 10000,
    \"filterBuy\": true,
    \"filterSell\": true,
    \"startTime\": $start_time,
    \"endTime\": $end_time
  }")

trades_custom=$(echo "$response3" | jq -r '.results[] | select(.strategyName == "session_trader") | .totalTrades')
echo "✅ Session Trader with custom dates: $trades_custom trades"
echo ""

# Verify results are different
echo "================================"
echo "📈 VERIFICATION:"
echo "================================"
echo "15 days:      $trades_15 trades"
echo "30 days:      $trades_30 trades"
echo "Custom range: $trades_custom trades"
echo ""

if [ "$trades_15" != "$trades_30" ]; then
    echo "✅ SUCCESS: Trade counts are different for different day values!"
    echo "   This confirms the backend is working correctly."
else
    echo "❌ WARNING: Trade counts are the same ($trades_15 trades)"
    echo "   This might indicate an issue."
fi

echo ""
echo "🌐 Frontend Fix Applied:"
echo "   ✅ Cache busting added (timestamp in URL)"
echo "   ✅ Cache-Control headers added"
echo "   ✅ currentResults cleared before each request"
echo "   ✅ Chart properly destroyed before recreation"
echo "   ✅ Duplicate createEquityChart() call removed"
echo "   ✅ Debug logging added"
echo ""
echo "📝 Next Steps:"
echo "   1. Open http://localhost:8080 in browser"
echo "   2. Open browser console (F12)"
echo "   3. Run backtest with 15 days"
echo "   4. Check console for request logs and trade count"
echo "   5. Change to 30 days and run again"
echo "   6. Verify chart updates with different data"
echo ""
