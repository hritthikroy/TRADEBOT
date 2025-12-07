#!/bin/bash

# Comprehensive Backtest Analysis
# Tests all timeframes and periods to find optimal settings

echo "🚀 COMPREHENSIVE BACKTEST ANALYSIS"
echo "=================================="
echo ""

# Output file
OUTPUT_FILE="comprehensive_backtest_results.json"
SUMMARY_FILE="COMPREHENSIVE_BACKTEST_SUMMARY.md"

echo "[]" > $OUTPUT_FILE

# Test periods (in days)
PERIODS=(2 5 7 15 20 30 60 90 180)

# Test timeframes
TIMEFRAMES=("1m" "5m" "15m" "30m" "1h" "4h" "1d")

echo "# 📊 COMPREHENSIVE BACKTEST ANALYSIS" > $SUMMARY_FILE
echo "" >> $SUMMARY_FILE
echo "Testing Session Trader strategy across all timeframes and periods." >> $SUMMARY_FILE
echo "" >> $SUMMARY_FILE
echo "## 🎯 Test Configuration" >> $SUMMARY_FILE
echo "" >> $SUMMARY_FILE
echo "- **Strategy:** session_trader (Balanced World-Class)" >> $SUMMARY_FILE
echo "- **Symbol:** BTCUSDT" >> $SUMMARY_FILE
echo "- **Start Balance:** \$500" >> $SUMMARY_FILE
echo "- **Risk per Trade:** 1%" >> $SUMMARY_FILE
echo "- **Periods Tested:** 2, 5, 7, 15, 20, 30, 60, 90, 180 days" >> $SUMMARY_FILE
echo "- **Timeframes Tested:** 1m, 5m, 15m, 30m, 1h, 4h, 1d" >> $SUMMARY_FILE
echo "" >> $SUMMARY_FILE
echo "---" >> $SUMMARY_FILE
echo "" >> $SUMMARY_FILE

# Counter
TOTAL_TESTS=$((${#PERIODS[@]} * ${#TIMEFRAMES[@]}))
CURRENT_TEST=0

echo "📋 Running $TOTAL_TESTS backtests..."
echo ""

# Test each timeframe
for TIMEFRAME in "${TIMEFRAMES[@]}"; do
    echo "## 📈 Timeframe: $TIMEFRAME" >> $SUMMARY_FILE
    echo "" >> $SUMMARY_FILE
    echo "| Days | Trades | Win Rate | Profit Factor | Return % | Max DD % | Final Balance | BUY WR | SELL WR |" >> $SUMMARY_FILE
    echo "|------|--------|----------|---------------|----------|----------|---------------|--------|---------|" >> $SUMMARY_FILE
    
    for DAYS in "${PERIODS[@]}"; do
        CURRENT_TEST=$((CURRENT_TEST + 1))
        echo "[$CURRENT_TEST/$TOTAL_TESTS] Testing $TIMEFRAME timeframe with $DAYS days..."
        
        # Run backtest
        RESULT=$(curl -s -X POST http://localhost:8080/api/v1/backtest/test-all-strategies \
            -H "Content-Type: application/json" \
            -d "{\"symbol\":\"BTCUSDT\",\"days\":$DAYS,\"startBalance\":500,\"riskPercent\":0.01,\"timeframe\":\"$TIMEFRAME\"}" \
            2>/dev/null | jq '.results[] | select(.strategyName == "session_trader")')
        
        if [ -z "$RESULT" ] || [ "$RESULT" == "null" ]; then
            echo "| $DAYS | N/A | N/A | N/A | N/A | N/A | N/A | N/A | N/A |" >> $SUMMARY_FILE
            echo "  ⚠️  No data available"
            continue
        fi
        
        # Extract metrics
        TRADES=$(echo $RESULT | jq -r '.totalTrades // 0')
        WIN_RATE=$(echo $RESULT | jq -r '.winRate // 0' | awk '{printf "%.1f", $1}')
        PROFIT_FACTOR=$(echo $RESULT | jq -r '.profitFactor // 0' | awk '{printf "%.2f", $1}')
        RETURN=$(echo $RESULT | jq -r '.returnPercent // 0' | awk '{printf "%.0f", $1}')
        MAX_DD=$(echo $RESULT | jq -r '.maxDrawdown // 0' | awk '{printf "%.1f", $1}')
        BALANCE=$(echo $RESULT | jq -r '.finalBalance // 0' | awk '{printf "%.0f", $1}')
        BUY_WR=$(echo $RESULT | jq -r '.buyWinRate // 0' | awk '{printf "%.0f", $1}')
        SELL_WR=$(echo $RESULT | jq -r '.sellWinRate // 0' | awk '{printf "%.0f", $1}')
        
        # Add to summary
        echo "| $DAYS | $TRADES | $WIN_RATE% | $PROFIT_FACTOR | $RETURN% | $MAX_DD% | \$$BALANCE | $BUY_WR% | $SELL_WR% |" >> $SUMMARY_FILE
        
        # Print progress
        echo "  ✅ Trades: $TRADES | WR: $WIN_RATE% | PF: $PROFIT_FACTOR | Return: $RETURN% | DD: $MAX_DD%"
        
        # Save to JSON
        echo $RESULT | jq ". + {timeframe: \"$TIMEFRAME\", days: $DAYS}" >> ${OUTPUT_FILE}.tmp
    done
    
    echo "" >> $SUMMARY_FILE
    echo "" >> $SUMMARY_FILE
done

# Combine JSON results
if [ -f "${OUTPUT_FILE}.tmp" ]; then
    jq -s '.' ${OUTPUT_FILE}.tmp > $OUTPUT_FILE
    rm ${OUTPUT_FILE}.tmp
fi

echo ""
echo "✅ Backtest complete!"
echo ""
echo "📊 Results saved to:"
echo "  - $OUTPUT_FILE (raw data)"
echo "  - $SUMMARY_FILE (formatted summary)"
echo ""

# Find best performers
echo "## 🏆 BEST PERFORMERS" >> $SUMMARY_FILE
echo "" >> $SUMMARY_FILE

# Best by profit factor
echo "### 🥇 Highest Profit Factor" >> $SUMMARY_FILE
echo "" >> $SUMMARY_FILE
BEST_PF=$(jq -r 'sort_by(-.profitFactor) | .[0] | "**Timeframe:** \(.timeframe) | **Days:** \(.days) | **PF:** \(.profitFactor | tonumber | . * 100 | floor / 100) | **WR:** \(.winRate | floor)% | **Return:** \(.returnPercent | floor)% | **DD:** \(.maxDrawdown | floor)%"' $OUTPUT_FILE 2>/dev/null)
if [ ! -z "$BEST_PF" ]; then
    echo "$BEST_PF" >> $SUMMARY_FILE
fi
echo "" >> $SUMMARY_FILE

# Best by return
echo "### 💰 Highest Return" >> $SUMMARY_FILE
echo "" >> $SUMMARY_FILE
BEST_RETURN=$(jq -r 'sort_by(-.returnPercent) | .[0] | "**Timeframe:** \(.timeframe) | **Days:** \(.days) | **Return:** \(.returnPercent | floor)% | **PF:** \(.profitFactor | tonumber | . * 100 | floor / 100) | **WR:** \(.winRate | floor)% | **DD:** \(.maxDrawdown | floor)%"' $OUTPUT_FILE 2>/dev/null)
if [ ! -z "$BEST_RETURN" ]; then
    echo "$BEST_RETURN" >> $SUMMARY_FILE
fi
echo "" >> $SUMMARY_FILE

# Best by win rate
echo "### 🎯 Highest Win Rate" >> $SUMMARY_FILE
echo "" >> $SUMMARY_FILE
BEST_WR=$(jq -r 'sort_by(-.winRate) | .[0] | "**Timeframe:** \(.timeframe) | **Days:** \(.days) | **WR:** \(.winRate | floor)% | **PF:** \(.profitFactor | tonumber | . * 100 | floor / 100) | **Return:** \(.returnPercent | floor)% | **DD:** \(.maxDrawdown | floor)%"' $OUTPUT_FILE 2>/dev/null)
if [ ! -z "$BEST_WR" ]; then
    echo "$BEST_WR" >> $SUMMARY_FILE
fi
echo "" >> $SUMMARY_FILE

# Lowest drawdown
echo "### 🛡️ Lowest Drawdown" >> $SUMMARY_FILE
echo "" >> $SUMMARY_FILE
BEST_DD=$(jq -r 'sort_by(.maxDrawdown) | .[0] | "**Timeframe:** \(.timeframe) | **Days:** \(.days) | **DD:** \(.maxDrawdown | floor)% | **PF:** \(.profitFactor | tonumber | . * 100 | floor / 100) | **WR:** \(.winRate | floor)% | **Return:** \(.returnPercent | floor)%"' $OUTPUT_FILE 2>/dev/null)
if [ ! -z "$BEST_DD" ]; then
    echo "$BEST_DD" >> $SUMMARY_FILE
fi
echo "" >> $SUMMARY_FILE

echo "---" >> $SUMMARY_FILE
echo "" >> $SUMMARY_FILE
echo "**Analysis complete!** Check the tables above for detailed results." >> $SUMMARY_FILE

echo "🎉 Analysis complete! Check $SUMMARY_FILE for results."
