#!/bin/bash

# TradingBot Test Script
# Quick testing for all strategies

echo "🧪 TradingBot Test Suite"
echo "========================"
echo ""

# Check if backend is running
if ! curl -s http://localhost:8080/api/v1/health > /dev/null 2>&1; then
    echo "❌ Backend is not running!"
    echo "Start it with: cd backend && go run ."
    exit 1
fi

echo "✅ Backend is running"
echo ""

# Test parameters
SYMBOL="${1:-BTCUSDT}"
INTERVAL="${2:-15m}"
DAYS="${3:-30}"
STRATEGY="${4:-session_trader}"

echo "📊 Test Parameters:"
echo "   Symbol: $SYMBOL"
echo "   Interval: $INTERVAL"
echo "   Days: $DAYS"
echo "   Strategy: $STRATEGY"
echo ""

# Run backtest
echo "🚀 Running backtest..."
echo ""

curl -s -X POST "http://localhost:8080/api/v1/backtest/run" \
  -H "Content-Type: application/json" \
  -d "{\"symbol\":\"$SYMBOL\",\"interval\":\"$INTERVAL\",\"days\":$DAYS,\"strategy\":\"$STRATEGY\",\"startBalance\":1000}" \
  | jq '{
    strategy,
    totalTrades,
    winningTrades,
    losingTrades,
    winRate,
    profitFactor,
    maxDrawdown,
    finalBalance,
    returnPercent: ((.finalBalance - 1000) / 1000 * 100)
  }'

echo ""
echo "✅ Test complete!"
echo ""
echo "💡 Usage: ./test.sh [SYMBOL] [INTERVAL] [DAYS] [STRATEGY]"
echo "   Example: ./test.sh BTCUSDT 15m 30 session_trader"
echo ""
echo "📚 Available strategies:"
echo "   - session_trader (recommended)"
echo "   - liquidity_hunter"
echo "   - breakout_master"
echo "   - trend_rider"
echo "   - range_master"
echo "   - smart_money_tracker"
echo "   - institutional_follower"
echo "   - reversal_sniper"
echo "   - momentum_beast"
echo "   - scalper_pro"
