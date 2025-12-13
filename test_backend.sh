#!/bin/bash
# Backend Testing Script

BASE_URL="http://localhost:8080"

echo "🧪 Testing Backend API..."
echo ""

# Health check
echo "1. Health Check"
curl -s "$BASE_URL/api/v1/health" | jq '.' || echo "❌ Health check failed"
echo ""

# Test signal generation
echo "2. Signal Generation"
curl -s -X POST "$BASE_URL/api/v1/signals" \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTCUSDT","strategy":"ict_smc"}' | jq '.' || echo "❌ Signal generation failed"
echo ""

# Test backtest
echo "3. Backtest"
curl -s -X POST "$BASE_URL/api/v1/backtest/run" \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTCUSDT","strategy":"ict_smc","days":7}' | jq '.' || echo "❌ Backtest failed"
echo ""

# Test paper trading status
echo "4. Paper Trading Status"
curl -s "$BASE_URL/api/v1/paper-trading/status" | jq '.' || echo "❌ Paper trading status failed"
echo ""

echo "✅ Tests complete!"
