#!/bin/bash

echo "🔍 SESSION TRADER DIAGNOSTIC"
echo "============================"
echo ""

# Test 1: Check if backend is running
echo "1. Backend Status:"
if curl -s http://localhost:8080/api/v1/health > /dev/null 2>&1; then
    echo "   ✅ Backend is running"
else
    echo "   ❌ Backend is NOT running"
    exit 1
fi
echo ""

# Test 2: Check available strategies
echo "2. Available Strategies:"
curl -s http://localhost:8080/api/v1/backtest/run \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTCUSDT","interval":"15m","days":1,"strategy":"liquidity_hunter","startBalance":1000}' | jq -r '.totalTrades' > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "   ✅ Backtest endpoint working"
else
    echo "   ❌ Backtest endpoint error"
fi
echo ""

# Test 3: Test with different strategies
echo "3. Testing Different Strategies (7 days):"
for strategy in "session_trader" "liquidity_hunter" "breakout_master"; do
    echo -n "   Testing $strategy... "
    trades=$(curl -s -X POST "http://localhost:8080/api/v1/backtest/run" \
        -H "Content-Type: application/json" \
        -d "{\"symbol\":\"BTCUSDT\",\"interval\":\"15m\",\"days\":7,\"strategy\":\"$strategy\",\"startBalance\":1000}" | jq -r '.totalTrades')
    echo "$trades trades"
done
echo ""

# Test 4: Test different timeframes
echo "4. Testing Session Trader - Different Days:"
for days in 1 3 7 14 30; do
    echo -n "   $days days... "
    result=$(curl -s -X POST "http://localhost:8080/api/v1/backtest/run" \
        -H "Content-Type: application/json" \
        -d "{\"symbol\":\"BTCUSDT\",\"interval\":\"15m\",\"days\":$days,\"strategy\":\"session_trader\",\"startBalance\":1000}")
    trades=$(echo "$result" | jq -r '.totalTrades')
    duration=$(echo "$result" | jq -r '.duration')
    echo "$trades trades (took $duration)"
done
echo ""

# Test 5: Check if it's a data issue
echo "5. Testing Data Availability:"
echo -n "   Fetching candles... "
candles=$(curl -s "https://api.binance.com/api/v3/klines?symbol=BTCUSDT&interval=15m&limit=100" | jq '. | length')
echo "$candles candles available from Binance"
echo ""

# Test 6: Test with date range
echo "6. Testing with Custom Date Range:"
echo -n "   Nov 2024... "
result=$(curl -s -X POST "http://localhost:8080/api/v1/backtest/run" \
    -H "Content-Type: application/json" \
    -d '{"symbol":"BTCUSDT","interval":"15m","startDate":"2024-11-01","endDate":"2024-11-30","strategy":"session_trader","startBalance":1000}')
trades=$(echo "$result" | jq -r '.totalTrades')
echo "$trades trades"
echo ""

# Summary
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "DIAGNOSIS SUMMARY"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "If all tests show 0 trades, possible issues:"
echo "  1. Strategy logic is too strict (no signals generated)"
echo "  2. Data fetching issue from Binance"
echo "  3. Signal generation logic has a bug"
echo "  4. Market regime filters are blocking all trades"
echo ""
echo "Next steps:"
echo "  • Check backend logs for errors"
echo "  • Test with a simpler strategy"
echo "  • Review Session Trader signal conditions"
echo "  • Check if AMD phase detection is blocking all signals"
echo ""
