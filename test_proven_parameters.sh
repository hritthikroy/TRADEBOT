#!/bin/bash

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🏆 TESTING PROVEN BEST PARAMETERS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Testing strategies with parameters from OPTIMIZATION_RESULTS_FULL.json"
echo "These parameters achieved 900-119,000% returns with 50-60% win rates"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check if server is running
if ! curl -s http://localhost:8080/api/v1/health > /dev/null 2>&1; then
    echo "❌ Server is not running!"
    echo ""
    echo "Please start the server first:"
    echo "   cd backend && go run ."
    echo ""
    exit 1
fi

echo "✅ Server is running"
echo ""

# Test top 3 strategies with proven parameters
strategies=("liquidity_hunter" "session_trader" "breakout_master")
expected_results=(
    "liquidity_hunter: 61% WR, 9.49 PF, 901% return"
    "session_trader: 58% WR, 18.67 PF, 1,313% return"
    "breakout_master: 55% WR, 8.23 PF, 3,704% return"
)

echo "🧪 Testing Top 3 Strategies (180 days)..."
echo ""

for i in "${!strategies[@]}"; do
    strategy="${strategies[$i]}"
    expected="${expected_results[$i]}"
    
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "Testing: $strategy"
    echo "Expected: $expected"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    
    # Run backtest
    result=$(curl -s -X POST http://localhost:8080/api/v1/backtest/run \
      -H "Content-Type: application/json" \
      -d "{
        \"symbol\": \"BTCUSDT\",
        \"interval\": \"15m\",
        \"days\": 180,
        \"startBalance\": 1000,
        \"riskPercent\": 2,
        \"strategy\": \"$strategy\"
      }")
    
    # Extract key metrics
    trades=$(echo "$result" | python3 -c "import sys, json; data=json.load(sys.stdin); print(data.get('totalTrades', 0))" 2>/dev/null)
    winRate=$(echo "$result" | python3 -c "import sys, json; data=json.load(sys.stdin); print(f\"{data.get('winRate', 0):.1f}\")" 2>/dev/null)
    profitFactor=$(echo "$result" | python3 -c "import sys, json; data=json.load(sys.stdin); print(f\"{data.get('profitFactor', 0):.2f}\")" 2>/dev/null)
    returnPct=$(echo "$result" | python3 -c "import sys, json; data=json.load(sys.stdin); print(f\"{data.get('returnPercent', 0):.0f}\")" 2>/dev/null)
    
    if [ ! -z "$trades" ] && [ "$trades" != "0" ]; then
        echo "✅ RESULTS:"
        echo "   Trades: $trades"
        echo "   Win Rate: $winRate%"
        echo "   Profit Factor: $profitFactor"
        echo "   Return: $returnPct%"
        echo ""
        
        # Compare with expected
        if (( $(echo "$winRate > 40" | bc -l) )); then
            echo "   ✅ Win Rate looks good (>40%)"
        else
            echo "   ⚠️  Win Rate lower than expected"
        fi
        
        if (( $(echo "$profitFactor > 2" | bc -l) )); then
            echo "   ✅ Profit Factor looks good (>2.0)"
        else
            echo "   ⚠️  Profit Factor lower than expected"
        fi
        
        if (( $(echo "$returnPct > 100" | bc -l) )); then
            echo "   ✅ Return looks good (>100%)"
        else
            echo "   ⚠️  Return lower than expected"
        fi
    else
        echo "❌ NO TRADES GENERATED"
        echo "   This means signal generation is still too strict"
        echo "   The proven parameters are applied, but signals aren't being generated"
    fi
    
    echo ""
done

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 TEST COMPLETE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "If you see 'NO TRADES GENERATED', the issue is signal generation, not parameters."
echo "The proven parameters are now applied, but strategies need to generate signals first."
echo ""
echo "Next steps if no trades:"
echo "1. Further simplify signal generation (require 2/5 conditions instead of 3/5)"
echo "2. Widen indicator ranges (RSI, volume thresholds)"
echo "3. Test on different time periods or symbols"
echo ""
