#!/bin/bash

echo "🧪 Testing Session Trader with AMD Phase Detection"
echo "=================================================="
echo ""

# Rebuild backend
echo "📦 Rebuilding backend with AMD phases..."
cd backend
go build -o ../tradebot
cd ..

# Start backend in background
echo "🚀 Starting backend..."
./tradebot &
BACKEND_PID=$!

# Wait for backend to start
echo "⏳ Waiting for backend to start..."
sleep 3

echo ""
echo "📊 Testing Session Trader Strategy with AMD Phases"
echo "=================================================="

# Test 30 days
echo ""
echo "1️⃣ Testing 30-day period..."
curl -X POST http://localhost:8080/api/v1/backtest/test-all-strategies \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "days": 30,
    "startBalance": 1000,
    "filterBuy": false,
    "filterSell": true
  }' | jq '.strategies[] | select(.name == "Session Trader") | {
    name,
    trades,
    winRate,
    profitFactor,
    maxDrawdown,
    totalReturn,
    wins,
    losses
  }'

echo ""
echo "2️⃣ Testing 7-day period..."
curl -X POST http://localhost:8080/api/v1/backtest/test-all-strategies \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "days": 7,
    "startBalance": 1000,
    "filterBuy": false,
    "filterSell": true
  }' | jq '.strategies[] | select(.name == "Session Trader") | {
    name,
    trades,
    winRate,
    profitFactor,
    maxDrawdown,
    totalReturn
  }'

echo ""
echo "3️⃣ Testing bad period (Nov 30 - Dec 4)..."
curl -X POST http://localhost:8080/api/v1/backtest/test-all-strategies \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "days": 5,
    "startBalance": 1000,
    "filterBuy": false,
    "filterSell": true
  }' | jq '.strategies[] | select(.name == "Session Trader") | {
    name,
    trades,
    winRate,
    profitFactor,
    maxDrawdown,
    totalReturn,
    wins,
    losses
  }'

# Stop backend
echo ""
echo "🛑 Stopping backend..."
kill $BACKEND_PID

echo ""
echo "✅ Testing complete!"
echo ""
echo "📋 What to look for:"
echo "  • Fewer trades (more selective)"
echo "  • Higher win rate (better quality)"
echo "  • Better profit factor"
echo "  • Lower drawdown"
echo "  • Signals should show AMD phase indicators (🟢 ACCUMULATION, 📈 MARKUP, 🔴 DISTRIBUTION, 📉 MARKDOWN)"
