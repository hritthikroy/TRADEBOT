#!/bin/bash

echo "🧪 Testing Optimization Fix..."
echo ""
echo "Running optimization with 30 days of data..."
echo ""

curl -X POST http://localhost:8080/api/v1/backtest/world-class-optimize \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "days": 30,
    "startBalance": 1000
  }' 2>/dev/null | python3 -m json.tool > test_results.json

echo ""
echo "✅ Results saved to test_results.json"
echo ""
echo "📊 Quick Summary:"
cat test_results.json | grep -A 3 '"bestScore"' | head -20
