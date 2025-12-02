#!/bin/bash

echo "🔬 Running Comprehensive Backtest"
echo "=================================="
echo ""
echo "Testing all strategies across multiple timeframes..."
echo ""

# Test different timeframes
TIMEFRAMES=("15m" "1h" "4h")
DAYS=30
BALANCE=500

echo "📊 Configuration:"
echo "   Symbol: BTCUSDT"
echo "   Days: $DAYS"
echo "   Starting Balance: \$$BALANCE"
echo "   Timeframes: ${TIMEFRAMES[*]}"
echo ""

# Run backtests via API
BASE_URL="http://localhost:8080/api/v1"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📈 BACKTEST RESULTS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

for tf in "${TIMEFRAMES[@]}"; do
    echo ""
    echo "⏱️ Timeframe: $tf"
    echo "-------------------"
    
    result=$(curl -s -X POST "$BASE_URL/backtest/run" \
        -H "Content-Type: application/json" \
        -d "{
            \"symbol\": \"BTCUSDT\",
            \"interval\": \"$tf\",
            \"days\": $DAYS,
            \"startBalance\": $BALANCE
        }")
    
    # Extract key metrics
    trades=$(echo $result | jq -r '.totalTrades // 0')
    winRate=$(echo $result | jq -r '.winRate // 0')
    returnPct=$(echo $result | jq -r '.returnPercent // 0')
    profitFactor=$(echo $result | jq -r '.profitFactor // 0')
    maxDD=$(echo $result | jq -r '.maxDrawdown // 0')
    finalBalance=$(echo $result | jq -r '.finalBalance // 0')
    
    echo "   Trades: $trades"
    echo "   Win Rate: ${winRate}%"
    echo "   Return: ${returnPct}%"
    echo "   Profit Factor: $profitFactor"
    echo "   Max Drawdown: ${maxDD}%"
    echo "   Final Balance: \$$finalBalance"
done

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Backtest Complete!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
