#!/bin/bash

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🧪 TESTING EACH STRATEGY INDIVIDUALLY"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check if server is running
if ! curl -s http://localhost:8080/api/v1/health > /dev/null 2>&1; then
    echo "❌ Server is not running!"
    echo "Please start the server first: cd backend && go run ."
    exit 1
fi

echo "✅ Server is running"
echo ""

# All 10 strategies
strategies=(
    "liquidity_hunter"
    "session_trader"
    "breakout_master"
    "trend_rider"
    "range_master"
    "smart_money_tracker"
    "institutional_follower"
    "reversal_sniper"
    "momentum_beast"
    "scalper_pro"
)

echo "Testing each strategy with a simple 30-day backtest..."
echo ""

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
    
    # Extract metrics
    totalTrades=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(data.get('totalTrades', 0))" 2>/dev/null || echo "0")
    winRate=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(f\"{data.get('winRate', 0):.2f}\")" 2>/dev/null || echo "0")
    profitFactor=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(f\"{data.get('profitFactor', 0):.2f}\")" 2>/dev/null || echo "0")
    returnPct=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(f\"{data.get('returnPercent', 0):.2f}\")" 2>/dev/null || echo "0")
    
    if [ "$totalTrades" -gt "0" ]; then
        echo "✅ WORKING - Generated $totalTrades trades"
        echo "   Win Rate: ${winRate}%"
        echo "   Profit Factor: ${profitFactor}"
        echo "   Return: ${returnPct}%"
    else
        echo "❌ NOT WORKING - No trades generated"
        echo "   This strategy needs signal generation logic in unified_signal_generator.go"
    fi
    echo ""
done

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🎯 SUMMARY"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Strategies that are working will show trades above."
echo "Strategies that show 0 trades need their signal generation logic implemented."
echo ""
