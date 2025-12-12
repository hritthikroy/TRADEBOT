#!/bin/bash

echo "Testing signal generation count for Liquidity Hunter..."
echo ""

result=$(curl -s -X POST http://localhost:8080/api/v1/backtest/run \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "interval": "4h",
    "days": 30,
    "startBalance": 1000,
    "strategy": "liquidity_hunter",
    "riskPercent": 0.01
  }')

trades=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(data.get('totalTrades', 0))")
winRate=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(data.get('winRate', 0))")
pf=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(data.get('profitFactor', 0))")

echo "30 days on 4h timeframe:"
echo "  Total Trades: $trades"
echo "  Win Rate: $winRate%"
echo "  Profit Factor: $pf"
echo ""

if [ "$trades" -lt 10 ]; then
    echo "✅ Signal count is good (< 10 trades in 30 days)"
elif [ "$trades" -lt 20 ]; then
    echo "⚠️  Signal count is moderate (10-20 trades in 30 days)"
else
    echo "❌ Too many signals ($trades trades in 30 days)"
    echo "   Expected: < 10 trades for 80-90% win rate strategy"
fi

echo ""
echo "Target for 80-90% WR: 5-10 trades per 30 days"
