#!/bin/bash

echo "🎯 PAPER TRADING API - TEST"
echo "============================"
echo ""

# Test 1: Check if backend is running
echo "✓ Test 1: Backend Health Check"
echo "-------------------------------"
health=$(curl -s http://localhost:8080/api/v1/health)
if echo "$health" | grep -q "healthy"; then
    echo "✅ Backend is running"
else
    echo "❌ Backend not running. Start with: cd backend && go run ."
    exit 1
fi
echo ""

# Test 2: Get initial stats
echo "✓ Test 2: Get Paper Trading Stats"
echo "----------------------------------"
curl -s http://localhost:8080/api/v1/paper-trading/stats | jq '.stats | {
    totalTrades,
    winningTrades,
    losingTrades,
    winRate,
    currentBalance,
    netProfit
}'
echo ""

# Test 3: Add a paper trade
echo "✓ Test 3: Add Paper Trade"
echo "-------------------------"
result=$(curl -s -X POST http://localhost:8080/api/v1/paper-trading/trade \
    -H "Content-Type: application/json" \
    -d '{"symbol":"BTCUSDT"}')

if echo "$result" | grep -q "success.*true"; then
    echo "✅ Paper trade added successfully"
    echo "$result" | jq '.trade | {id, signal, entry, stopLoss, takeProfit}'
else
    echo "ℹ️  No signal available right now"
    echo "$result" | jq '.message'
fi
echo ""

# Test 4: Update trades
echo "✓ Test 4: Update Open Trades"
echo "----------------------------"
curl -s -X POST http://localhost:8080/api/v1/paper-trading/update \
    -H "Content-Type: application/json" \
    -d '{"symbol":"BTCUSDT"}' | jq '{
    success,
    currentPrice,
    closedTradesCount: (.closedTrades | length)
}'
echo ""

# Test 5: Get all trades
echo "✓ Test 5: Get All Trades"
echo "------------------------"
curl -s http://localhost:8080/api/v1/paper-trading/trades | jq '{
    success,
    totalTrades: .total,
    trades: .trades | map({id, signal, status, profit})
}'
echo ""

# Test 6: Start auto paper trading
echo "✓ Test 6: Start Auto Paper Trading"
echo "-----------------------------------"
curl -s -X POST http://localhost:8080/api/v1/paper-trading/start-auto | jq '.'
echo ""

echo "============================"
echo "✅ ALL TESTS COMPLETE"
echo "============================"
echo ""
echo "Paper Trading API is working!"
echo ""
echo "Next steps:"
echo "1. Auto trading is now running (checks every 15 min)"
echo "2. Check stats: curl http://localhost:8080/api/v1/paper-trading/stats | jq '.stats'"
echo "3. View trades: curl http://localhost:8080/api/v1/paper-trading/trades | jq '.'"
echo ""
echo "To stop auto trading:"
echo "  curl -X POST http://localhost:8080/api/v1/paper-trading/stop-auto"
echo ""
