#!/bin/bash

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🧪 TESTING 5 PREVIOUSLY BROKEN STRATEGIES"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

strategies=(
    "session_trader"
    "trend_rider"
    "range_master"
    "reversal_sniper"
    "scalper_pro"
)

for strategy in "${strategies[@]}"; do
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "Testing: $strategy"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    result=$(curl -s -X POST http://localhost:8080/api/v1/backtest/run \
      -H "Content-Type: application/json" \
      -d "{
        \"symbol\": \"BTCUSDT\",
        \"interval\": \"15m\",
        \"days\": 30,
        \"startBalance\": 1000,
        \"strategy\": \"$strategy\"
      }")
    
    totalTrades=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(data.get('totalTrades', 0))" 2>/dev/null || echo "0")
    winRate=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(f\"{data.get('winRate', 0):.2f}\")" 2>/dev/null || echo "0")
    profitFactor=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(f\"{data.get('profitFactor', 0):.2f}\")" 2>/dev/null || echo "0")
    returnPct=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(f\"{data.get('returnPercent', 0):.2f}\")" 2>/dev/null || echo "0")
    
    if [ "$totalTrades" -gt "0" ]; then
        echo "✅ NOW WORKING - Generated $totalTrades trades"
        echo "   Win Rate: ${winRate}%"
        echo "   Profit Factor: ${profitFactor}"
        echo "   Return: ${returnPct}%"
    else
        echo "❌ STILL NOT WORKING - No trades generated"
    fi
    echo ""
done

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ TEST COMPLETE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
