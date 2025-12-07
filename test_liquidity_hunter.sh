#!/bin/bash

echo "=========================================="
echo "LIQUIDITY HUNTER STRATEGY - COMPREHENSIVE TEST"
echo "Testing proven 61% WR strategy"
echo "=========================================="
echo ""

# Test all time periods
periods=("1" "3" "5" "7" "15" "30" "60" "90")

for days in "${periods[@]}"; do
    echo "----------------------------------------"
    echo "Testing ${days}d period..."
    echo "----------------------------------------"
    
    response=$(curl -s -X POST "http://localhost:8080/api/v1/backtest/run" \
        -H "Content-Type: application/json" \
        -d "{\"strategy\":\"liquidity_hunter\",\"days\":${days}}")
    
    # Extract key metrics using python
    metrics=$(echo "$response" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    print(f\"{data.get('totalTrades', 0)}|{data.get('winRate', 0):.1f}|{data.get('profitFactor', 0):.2f}|{data.get('maxDrawdown', 0):.1f}|{data.get('returnPercentage', 0):.1f}\")
except:
    print('0|0|0|0|0')
")
    
    IFS='|' read -r trades winRate profitFactor maxDrawdown returnPct <<< "$metrics"
    
    # Determine status
    status="✅"
    if (( $(echo "$winRate < 50" | bc -l 2>/dev/null || echo 1) )) || (( $(echo "$maxDrawdown > 30" | bc -l 2>/dev/null || echo 1) )); then
        status="⚠️"
    fi
    if (( $(echo "$winRate < 40" | bc -l 2>/dev/null || echo 1) )) || (( $(echo "$maxDrawdown > 40" | bc -l 2>/dev/null || echo 1) )); then
        status="❌"
    fi
    if [ "$trades" = "0" ]; then
        status="⚠️"
    fi
    
    echo "${status} ${days}d: Trades=$trades, WR=${winRate}%, PF=${profitFactor}, DD=${maxDrawdown}%, Return=${returnPct}%"
    echo ""
done

echo "=========================================="
echo "LEGEND:"
echo "✅ = Good (50%+ WR, <30% DD)"
echo "⚠️ = Moderate (40-50% WR or 30-40% DD)"
echo "❌ = Poor (<40% WR or >40% DD)"
echo "=========================================="
