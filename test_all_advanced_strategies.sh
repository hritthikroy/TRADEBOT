#!/bin/bash

# Test All 10 Advanced Strategies
# Finds the best performing strategy

echo "🚀 TESTING ALL 10 ADVANCED STRATEGIES"
echo "======================================"
echo ""

BASE_URL="http://localhost:8080"

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}Symbol: BTCUSDT${NC}"
echo -e "${BLUE}Start Balance: \$500${NC}"
echo ""

echo "📊 Testing all strategies (this may take 2-3 minutes)..."
echo ""

# Call the API to test all strategies
RESULT=$(curl -s -X POST "$BASE_URL/api/v1/backtest/test-all-strategies" \
    -H "Content-Type: application/json" \
    -d '{"symbol":"BTCUSDT","startBalance":500}')

# Check if successful
if echo "$RESULT" | grep -q "bestStrategy"; then
    echo -e "${GREEN}✅ All strategies tested successfully!${NC}"
    echo ""
    
    # Parse and display results
    echo "$RESULT" | python3 << 'PYTHON'
import json
import sys

data = json.load(sys.stdin)
results = data.get('results', [])
best = data.get('bestStrategy', {})

print("┌────┬─────────────────────────┬──────────┬─────────┬──────────┬──────────┬──────────┬────────┐")
print("│Rank│ Strategy                │Timeframe │ Trades  │ Win Rate │ Return % │ Profit F │ Score  │")
print("├────┼─────────────────────────┼──────────┼─────────┼──────────┼──────────┼──────────┼────────┤")

for i, r in enumerate(results[:10], 1):
    status = "✅" if r['winRate'] >= r.get('targetWinRate', 60) * 0.9 else "❌"
    print(f"│ {i:2d} │ {r['strategyName'][:23]:23s} │ {r['timeframe']:8s} │ {r['totalTrades']:7d} │ {r['winRate']:7.1f}% │ {r['returnPercent']:7.1f}% │ {r['profitFactor']:8.2f} │ {r['score']:6.1f} │ {status}")

print("└────┴─────────────────────────┴──────────┴─────────┴──────────┴──────────┴──────────┴────────┘")
print("")
print("🏆 BEST STRATEGY:")
print(f"   Name: {best.get('strategyName', 'N/A')}")
print(f"   Timeframe: {best.get('timeframe', 'N/A')}")
print(f"   Win Rate: {best.get('winRate', 0):.1f}%")
print(f"   Return: {best.get('returnPercent', 0):.1f}%")
print(f"   Profit Factor: {best.get('profitFactor', 0):.2f}")
print(f"   Total Trades: {best.get('totalTrades', 0)}")
print("")
PYTHON
    
else
    echo -e "${RED}❌ Failed to test strategies${NC}"
    echo "Error: $RESULT"
fi

echo ""
echo -e "${GREEN}✅ Strategy testing complete!${NC}"
