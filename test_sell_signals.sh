#!/bin/bash

# Test SELL signal generation for different strategies
echo "🔍 Testing SELL Signal Generation"
echo "=================================="
echo ""

# Test Session Trader
echo "📊 Testing Session Trader..."
curl -s -X POST http://localhost:8080/api/live-signal \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTCUSDT","strategy":"session_trader"}' | jq '.'

echo ""
echo "---"
echo ""

# Test Breakout Master
echo "📊 Testing Breakout Master..."
curl -s -X POST http://localhost:8080/api/live-signal \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTCUSDT","strategy":"breakout_master"}' | jq '.'

echo ""
echo "---"
echo ""

# Test Liquidity Hunter
echo "📊 Testing Liquidity Hunter..."
curl -s -X POST http://localhost:8080/api/live-signal \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTCUSDT","strategy":"liquidity_hunter"}' | jq '.'

echo ""
echo "---"
echo ""

# Test Range Master
echo "📊 Testing Range Master..."
curl -s -X POST http://localhost:8080/api/live-signal \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTCUSDT","strategy":"range_master"}' | jq '.'

echo ""
echo "---"
echo ""

# Test Trend Rider
echo "📊 Testing Trend Rider..."
curl -s -X POST http://localhost:8080/api/live-signal \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTCUSDT","strategy":"trend_rider"}' | jq '.'

echo ""
echo "=================================="
echo "✅ Test Complete"
echo ""
echo "💡 Note: If all signals are NONE or BUY, it means:"
echo "   1. Market is in an uptrend (no SELL conditions met)"
echo "   2. SELL conditions are too strict"
echo "   3. Try different symbols (ETHUSDT, SOLUSDT, etc.)"
echo ""
echo "🔧 To test with different symbols:"
echo "   ./test_sell_signals.sh ETHUSDT"
