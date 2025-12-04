#!/bin/bash

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🚀 QUICK TEST - Optimized Parameters (All 10 Strategies)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Testing 10 strategies with their optimal timeframes"
echo "Using current optimized parameters (MinConfluence 4-5)"
echo ""

# Strategy configurations (strategy:timeframe)
declare -a configs=(
    "liquidity_hunter:15m"
    "smart_money_tracker:1h"
    "breakout_master:15m"
    "trend_rider:4h"
    "scalper_pro:5m"
    "reversal_sniper:1h"
    "session_trader:15m"
    "momentum_beast:15m"
    "range_master:1h"
    "institutional_follower:4h"
)

RESULTS_FILE="quick_test_results_$(date +%Y%m%d_%H%M%S).json"
echo "[" > "$RESULTS_FILE"

total=${#configs[@]}
current=0

for config in "${configs[@]}"; do
    current=$((current + 1))
    IFS=':' read -r strategy timeframe <<< "$config"
    
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "[$current/$total] Testing: $strategy ($timeframe)"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    # Run backtest with optimized parameters
    result=$(curl -s -X POST http://localhost:8080/api/v1/backtest/run \
      -H "Content-Type: application/json" \
      -d "{
        \"symbol\": \"BTCUSDT\",
        \"interval\": \"$timeframe\",
        \"days\": 90,
        \"startBalance\": 1000,
        \"strategy\": \"$strategy\",
        \"riskPercent\": 0.02
      }")
    
    # Check if successful
    if echo "$result" | grep -q "totalTrades"; then
        # Extract key metrics
        trades=$(echo "$result" | python3 -c "import sys,json; print(json.load(sys.stdin).get('totalTrades', 0))" 2>/dev/null || echo "0")
        winrate=$(echo "$result" | python3 -c "import sys,json; print(round(json.load(sys.stdin).get('winRate', 0), 2))" 2>/dev/null || echo "0")
        pf=$(echo "$result" | python3 -c "import sys,json; print(round(json.load(sys.stdin).get('profitFactor', 0), 2))" 2>/dev/null || echo "0")
        ret=$(echo "$result" | python3 -c "import sys,json; print(round(json.load(sys.stdin).get('returnPercent', 0), 2))" 2>/dev/null || echo "0")
        maxdd=$(echo "$result" | python3 -c "import sys,json; print(round(json.load(sys.stdin).get('maxDrawdown', 0), 2))" 2>/dev/null || echo "0")
        final=$(echo "$result" | python3 -c "import sys,json; print(round(json.load(sys.stdin).get('finalBalance', 0), 2))" 2>/dev/null || echo "0")
        
        echo "✅ Results:"
        echo "   Trades: $trades"
        echo "   Win Rate: $winrate%"
        echo "   Profit Factor: $pf"
        echo "   Return: $ret%"
        echo "   Max Drawdown: $maxdd%"
        echo "   Final Balance: \$$final"
        
        # Save to file
        if [ $current -gt 1 ]; then
            echo "," >> "$RESULTS_FILE"
        fi
        echo "  {" >> "$RESULTS_FILE"
        echo "    \"strategy\": \"$strategy\"," >> "$RESULTS_FILE"
        echo "    \"timeframe\": \"$timeframe\"," >> "$RESULTS_FILE"
        echo "    \"trades\": $trades," >> "$RESULTS_FILE"
        echo "    \"winRate\": $winrate," >> "$RESULTS_FILE"
        echo "    \"profitFactor\": $pf," >> "$RESULTS_FILE"
        echo "    \"returnPercent\": $ret," >> "$RESULTS_FILE"
        echo "    \"maxDrawdown\": $maxdd," >> "$RESULTS_FILE"
        echo "    \"finalBalance\": $final" >> "$RESULTS_FILE"
        echo -n "  }" >> "$RESULTS_FILE"
    else
        echo "❌ Failed to get results"
        echo "$result" | head -3
    fi
    
    echo ""
    
    # Small delay to avoid rate limiting
    sleep 2
done

echo "]" >> "$RESULTS_FILE"

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ TESTING COMPLETE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📊 SUMMARY:"
echo ""

# Generate summary
python3 << 'PYTHON_SCRIPT'
import json
import sys

try:
    with open([f for f in __import__('os').listdir('.') if f.startswith('quick_test_results_') and f.endswith('.json')][-1], 'r') as f:
        results = json.load(f)
    
    # Sort by return
    sorted_results = sorted(results, key=lambda x: x.get('returnPercent', 0), reverse=True)
    
    print("🏆 TOP 5 STRATEGIES BY RETURN:")
    print("")
    for i, r in enumerate(sorted_results[:5], 1):
        print(f"{i}. {r['strategy']} ({r['timeframe']})")
        print(f"   Return: {r['returnPercent']}% | WR: {r['winRate']}% | PF: {r['profitFactor']} | Trades: {r['trades']}")
        print("")
    
    print("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
    print("📈 BEST BY METRIC:")
    print("")
    
    # Best win rate
    best_wr = max(results, key=lambda x: x.get('winRate', 0))
    print(f"🎯 Highest Win Rate: {best_wr['strategy']} ({best_wr['timeframe']}) - {best_wr['winRate']}%")
    
    # Best profit factor
    best_pf = max(results, key=lambda x: x.get('profitFactor', 0))
    print(f"💰 Best Profit Factor: {best_pf['strategy']} ({best_pf['timeframe']}) - {best_pf['profitFactor']}")
    
    # Most trades
    best_trades = max(results, key=lambda x: x.get('trades', 0))
    print(f"⚡ Most Active: {best_trades['strategy']} ({best_trades['timeframe']}) - {best_trades['trades']} trades")
    
    # Best return
    best_ret = max(results, key=lambda x: x.get('returnPercent', 0))
    print(f"🚀 Highest Return: {best_ret['strategy']} ({best_ret['timeframe']}) - {best_ret['returnPercent']}%")
    
    print("")
    print("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
    
except Exception as e:
    print(f"Error generating summary: {e}")
    sys.exit(1)
PYTHON_SCRIPT

echo ""
echo "📄 Full results saved to: $RESULTS_FILE"
echo ""
echo "View results:"
echo "  cat $RESULTS_FILE | python3 -m json.tool"
echo ""
