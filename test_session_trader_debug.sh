#!/bin/bash

echo "🔍 SESSION TRADER DEBUG TEST"
echo "============================="
echo ""

# Test with verbose output
echo "Testing Session Trader with 30 days..."
result=$(curl -s -X POST "http://localhost:8080/api/v1/backtest/run" \
    -H "Content-Type: application/json" \
    -d '{"symbol":"BTCUSDT","interval":"15m","days":30,"strategy":"session_trader","startBalance":1000}')

echo "Full Response:"
echo "$result" | jq '.'
echo ""

# Check specific fields
total_trades=$(echo "$result" | jq -r '.totalTrades')
duration=$(echo "$result" | jq -r '.duration')

echo "Summary:"
echo "  Total Trades: $total_trades"
echo "  Duration: $duration"
echo ""

# Test Breakout Master for comparison
echo "Testing Breakout Master (for comparison)..."
result2=$(curl -s -X POST "http://localhost:8080/api/v1/backtest/run" \
    -H "Content-Type: application/json" \
    -d '{"symbol":"BTCUSDT","interval":"15m","days":30,"strategy":"breakout_master","startBalance":1000}')

total_trades2=$(echo "$result2" | jq -r '.totalTrades')
echo "  Breakout Master Trades: $total_trades2"
echo ""

# Test if it's a data issue
echo "Checking data availability..."
echo "  Requesting 30 days of 15m candles = ~2880 candles"
echo "  Binance limit per request = 1000 candles"
echo "  Multiple requests needed = YES"
echo ""

echo "If Session Trader shows 0 trades but Breakout Master shows >0,"
echo "then the issue is in Session Trader signal generation logic."
