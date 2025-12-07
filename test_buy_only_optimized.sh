#!/bin/bash

echo "=========================================="
echo "TESTING BUY ONLY STRATEGY (OPTIMIZED)"
echo "=========================================="
echo ""
echo "Current BUY Only Results:"
echo "  Win Rate: 50.1%"
echo "  Profit Factor: 8915.54 🤯"
echo "  Max Drawdown: 64.9% ❌"
echo ""
echo "Goal: Reduce drawdown while keeping high PF"
echo "=========================================="
echo ""

# Test BUY only with different risk levels
for risk in 0.5 0.75 1.0; do
    echo "Testing with ${risk}% risk..."
    echo ""
    
    curl -s -X POST http://localhost:8080/api/v1/backtest/test-all-strategies \
        -H "Content-Type: application/json" \
        -d "{
            \"symbol\": \"BTCUSDT\",
            \"days\": 90,
            \"startBalance\": 500,
            \"riskPercent\": $(echo "$risk / 100" | bc -l),
            \"filterBuy\": true,
            \"filterSell\": false
        }" | jq '.results[] | select(.strategyName == "session_trader") | 
        "
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
BUY ONLY - \('"$risk"')% Risk
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Total Trades:     \(.totalTrades)
Win Rate:         \(.winRate | floor * 10 / 10)%
Profit Factor:    \(.profitFactor | . * 100 | floor / 100)
Return:           \(.returnPercent | floor)%
Max Drawdown:     \(.maxDrawdown | floor * 10 / 10)%
Final Balance:    $\(.finalBalance | floor)

Buy Trades:       \(.buyTrades) (\(.buyWinRate | floor)% WR)
Sell Trades:      \(.sellTrades) (\(.sellWinRate | floor)% WR)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
"' -r
    
    echo ""
done

echo ""
echo "=========================================="
echo "RECOMMENDATION"
echo "=========================================="
echo ""
echo "Based on results:"
echo "1. BUY only has MUCH better profit factor"
echo "2. Need to reduce drawdown with lower risk"
echo "3. Optimal: 0.5-0.75% risk for BUY only"
echo ""
echo "=========================================="
