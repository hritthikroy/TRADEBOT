#!/bin/bash

echo "🔍 TESTING DATE RANGE ISSUE"
echo "============================"
echo ""
echo "Testing if different 'days' parameters return different results..."
echo ""

# Test multiple periods
for days in 15 30 60 90; do
    echo "Testing ${days} days..."
    
    result=$(curl -s -X POST http://localhost:8080/api/v1/backtest/test-all-strategies \
      -H "Content-Type: application/json" \
      -d "{\"symbol\":\"BTCUSDT\",\"days\":${days},\"startBalance\":1000,\"filterBuy\":false,\"filterSell\":true}" | \
    python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    for r in data['results']:
        if r['strategyName'] == 'session_trader':
            print(f\"{r['totalTrades']} trades, {r['winningTrades']}W/{r['losingTrades']}L, {r['winRate']:.1f}% WR\")
            break
except:
    print('Error')
" 2>&1)
    
    echo "  ${days}d: ${result}"
    sleep 1
done

echo ""
echo "📊 ANALYSIS:"
echo "============"
echo ""
echo "If all periods show IDENTICAL results (same trade count), there's a bug!"
echo ""
echo "Expected behavior:"
echo "  15d: ~94 trades"
echo "  30d: ~150-180 trades (MORE than 15d)"
echo "  60d: ~250-300 trades (MORE than 30d)"
echo "  90d: ~350-400 trades (MORE than 60d)"
echo ""
echo "⚠️  If you see the same trade count for all periods, the backend"
echo "    is NOT fetching different date ranges correctly!"
echo ""
