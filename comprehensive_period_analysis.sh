#!/bin/bash

# Comprehensive Period Analysis
# Tests all periods to find optimal trading windows

echo "🚀 COMPREHENSIVE PERIOD ANALYSIS"
echo "================================"
echo ""

# Output files
OUTPUT_FILE="period_analysis_results.json"
SUMMARY_FILE="COMPREHENSIVE_PERIOD_ANALYSIS.md"

echo "[]" > $OUTPUT_FILE

# Test periods (in days)
PERIODS=(2 3 5 7 10 15 20 30 45 60 90 120 150 180)

echo "# 📊 COMPREHENSIVE PERIOD ANALYSIS" > $SUMMARY_FILE
echo "" >> $SUMMARY_FILE
echo "**Strategy:** Session Trader (Balanced World-Class)" >> $SUMMARY_FILE
echo "**Symbol:** BTCUSDT" >> $SUMMARY_FILE
echo "**Timeframe:** 15m" >> $SUMMARY_FILE
echo "**Start Balance:** \$500" >> $SUMMARY_FILE
echo "**Risk per Trade:** 1%" >> $SUMMARY_FILE
echo "" >> $SUMMARY_FILE
echo "---" >> $SUMMARY_FILE
echo "" >> $SUMMARY_FILE

echo "## 📈 Results by Period" >> $SUMMARY_FILE
echo "" >> $SUMMARY_FILE
echo "| Days | Trades | Win Rate | Profit Factor | Return % | Max DD % | Final Balance | BUY Trades | SELL Trades | BUY WR | SELL WR | Sharpe |" >> $SUMMARY_FILE
echo "|------|--------|----------|---------------|----------|----------|---------------|------------|-------------|--------|---------|--------|" >> $SUMMARY_FILE

TOTAL_TESTS=${#PERIODS[@]}
CURRENT_TEST=0

for DAYS in "${PERIODS[@]}"; do
    CURRENT_TEST=$((CURRENT_TEST + 1))
    echo "[$CURRENT_TEST/$TOTAL_TESTS] Testing $DAYS days..."
    
    # Run backtest
    RESULT=$(curl -s -X POST http://localhost:8080/api/v1/backtest/test-all-strategies \
        -H "Content-Type: application/json" \
        -d "{\"symbol\":\"BTCUSDT\",\"days\":$DAYS,\"startBalance\":500,\"riskPercent\":0.01}" \
        2>/dev/null | jq '.results[] | select(.strategyName == "session_trader")')
    
    if [ -z "$RESULT" ] || [ "$RESULT" == "null" ]; then
        echo "| $DAYS | N/A | N/A | N/A | N/A | N/A | N/A | N/A | N/A | N/A | N/A | N/A |" >> $SUMMARY_FILE
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
    BUY_TRADES=$(echo $RESULT | jq -r '.buyTrades // 0')
    SELL_TRADES=$(echo $RESULT | jq -r '.sellTrades // 0')
    BUY_WR=$(echo $RESULT | jq -r '.buyWinRate // 0' | awk '{printf "%.0f", $1}')
    SELL_WR=$(echo $RESULT | jq -r '.sellWinRate // 0' | awk '{printf "%.0f", $1}')
    
    # Calculate Sharpe-like ratio (Return / DD)
    if [ "$MAX_DD" != "0" ] && [ "$MAX_DD" != "0.0" ]; then
        SHARPE=$(echo "scale=2; $RETURN / $MAX_DD" | bc)
    else
        SHARPE="N/A"
    fi
    
    # Add to summary
    echo "| $DAYS | $TRADES | $WIN_RATE% | $PROFIT_FACTOR | $RETURN% | $MAX_DD% | \$$BALANCE | $BUY_TRADES | $SELL_TRADES | $BUY_WR% | $SELL_WR% | $SHARPE |" >> $SUMMARY_FILE
    
    # Print progress
    echo "  ✅ Trades: $TRADES | WR: $WIN_RATE% | PF: $PROFIT_FACTOR | Return: $RETURN% | DD: $MAX_DD% | Sharpe: $SHARPE"
    
    # Save to JSON
    echo $RESULT | jq ". + {days: $DAYS, sharpe: \"$SHARPE\"}" >> ${OUTPUT_FILE}.tmp
    
    # Small delay to avoid overwhelming the server
    sleep 1
done

# Combine JSON results
if [ -f "${OUTPUT_FILE}.tmp" ]; then
    jq -s '.' ${OUTPUT_FILE}.tmp > $OUTPUT_FILE
    rm ${OUTPUT_FILE}.tmp
fi

echo "" >> $SUMMARY_FILE
echo "" >> $SUMMARY_FILE

# Analysis section
echo "## 🏆 BEST PERFORMERS" >> $SUMMARY_FILE
echo "" >> $SUMMARY_FILE

# Best by profit factor
echo "### 🥇 Highest Profit Factor" >> $SUMMARY_FILE
echo "" >> $SUMMARY_FILE
BEST_PF=$(jq -r 'sort_by(-.profitFactor) | .[0] | "- **Period:** \(.days) days\n- **Profit Factor:** \(.profitFactor | tonumber | . * 100 | floor / 100)\n- **Win Rate:** \(.winRate | floor)%\n- **Return:** \(.returnPercent | floor)%\n- **Max DD:** \(.maxDrawdown | floor)%\n- **Trades:** \(.totalTrades)"' $OUTPUT_FILE 2>/dev/null)
if [ ! -z "$BEST_PF" ]; then
    echo "$BEST_PF" >> $SUMMARY_FILE
fi
echo "" >> $SUMMARY_FILE

# Best by return
echo "### 💰 Highest Return" >> $SUMMARY_FILE
echo "" >> $SUMMARY_FILE
BEST_RETURN=$(jq -r 'sort_by(-.returnPercent) | .[0] | "- **Period:** \(.days) days\n- **Return:** \(.returnPercent | floor)%\n- **Profit Factor:** \(.profitFactor | tonumber | . * 100 | floor / 100)\n- **Win Rate:** \(.winRate | floor)%\n- **Max DD:** \(.maxDrawdown | floor)%\n- **Trades:** \(.totalTrades)"' $OUTPUT_FILE 2>/dev/null)
if [ ! -z "$BEST_RETURN" ]; then
    echo "$BEST_RETURN" >> $SUMMARY_FILE
fi
echo "" >> $SUMMARY_FILE

# Best by win rate
echo "### 🎯 Highest Win Rate" >> $SUMMARY_FILE
echo "" >> $SUMMARY_FILE
BEST_WR=$(jq -r 'sort_by(-.winRate) | .[0] | "- **Period:** \(.days) days\n- **Win Rate:** \(.winRate | floor)%\n- **Profit Factor:** \(.profitFactor | tonumber | . * 100 | floor / 100)\n- **Return:** \(.returnPercent | floor)%\n- **Max DD:** \(.maxDrawdown | floor)%\n- **Trades:** \(.totalTrades)"' $OUTPUT_FILE 2>/dev/null)
if [ ! -z "$BEST_WR" ]; then
    echo "$BEST_WR" >> $SUMMARY_FILE
fi
echo "" >> $SUMMARY_FILE

# Lowest drawdown
echo "### 🛡️ Lowest Drawdown" >> $SUMMARY_FILE
echo "" >> $SUMMARY_FILE
BEST_DD=$(jq -r 'sort_by(.maxDrawdown) | .[0] | "- **Period:** \(.days) days\n- **Max DD:** \(.maxDrawdown | floor)%\n- **Profit Factor:** \(.profitFactor | tonumber | . * 100 | floor / 100)\n- **Win Rate:** \(.winRate | floor)%\n- **Return:** \(.returnPercent | floor)%\n- **Trades:** \(.totalTrades)"' $OUTPUT_FILE 2>/dev/null)
if [ ! -z "$BEST_DD" ]; then
    echo "$BEST_DD" >> $SUMMARY_FILE
fi
echo "" >> $SUMMARY_FILE

# Best Sharpe ratio
echo "### ⚡ Best Risk-Adjusted Return (Sharpe)" >> $SUMMARY_FILE
echo "" >> $SUMMARY_FILE
echo "Sharpe Ratio = Return / Drawdown (higher is better)" >> $SUMMARY_FILE
echo "" >> $SUMMARY_FILE
BEST_SHARPE=$(jq -r '[.[] | select(.sharpe != "N/A")] | sort_by(-.sharpe | tonumber) | .[0] | "- **Period:** \(.days) days\n- **Sharpe Ratio:** \(.sharpe)\n- **Return:** \(.returnPercent | floor)%\n- **Max DD:** \(.maxDrawdown | floor)%\n- **Profit Factor:** \(.profitFactor | tonumber | . * 100 | floor / 100)\n- **Win Rate:** \(.winRate | floor)%"' $OUTPUT_FILE 2>/dev/null)
if [ ! -z "$BEST_SHARPE" ]; then
    echo "$BEST_SHARPE" >> $SUMMARY_FILE
fi
echo "" >> $SUMMARY_FILE

# Most trades
echo "### 📊 Most Active Period" >> $SUMMARY_FILE
echo "" >> $SUMMARY_FILE
MOST_TRADES=$(jq -r 'sort_by(-.totalTrades) | .[0] | "- **Period:** \(.days) days\n- **Total Trades:** \(.totalTrades)\n- **Trades/Day:** \((.totalTrades / .days) | floor)\n- **Win Rate:** \(.winRate | floor)%\n- **Profit Factor:** \(.profitFactor | tonumber | . * 100 | floor / 100)\n- **Return:** \(.returnPercent | floor)%"' $OUTPUT_FILE 2>/dev/null)
if [ ! -z "$MOST_TRADES" ]; then
    echo "$MOST_TRADES" >> $SUMMARY_FILE
fi
echo "" >> $SUMMARY_FILE

echo "---" >> $SUMMARY_FILE
echo "" >> $SUMMARY_FILE

# Recommendations
echo "## 💡 RECOMMENDATIONS" >> $SUMMARY_FILE
echo "" >> $SUMMARY_FILE
echo "### For Different Trading Styles:" >> $SUMMARY_FILE
echo "" >> $SUMMARY_FILE
echo "**Conservative (Low Risk):**" >> $SUMMARY_FILE
CONSERVATIVE=$(jq -r '[.[] | select(.maxDrawdown < 15)] | sort_by(-.profitFactor) | .[0] | "- Use **\(.days) days** period\n- Expected DD: \(.maxDrawdown | floor)%\n- Expected PF: \(.profitFactor | tonumber | . * 100 | floor / 100)\n- Expected WR: \(.winRate | floor)%"' $OUTPUT_FILE 2>/dev/null)
if [ ! -z "$CONSERVATIVE" ]; then
    echo "$CONSERVATIVE" >> $SUMMARY_FILE
fi
echo "" >> $SUMMARY_FILE

echo "**Balanced (Best Overall):**" >> $SUMMARY_FILE
BALANCED=$(jq -r '[.[] | select(.sharpe != "N/A")] | sort_by(-.sharpe | tonumber) | .[0] | "- Use **\(.days) days** period\n- Best risk-adjusted returns (Sharpe: \(.sharpe))\n- Expected DD: \(.maxDrawdown | floor)%\n- Expected PF: \(.profitFactor | tonumber | . * 100 | floor / 100)\n- Expected WR: \(.winRate | floor)%"' $OUTPUT_FILE 2>/dev/null)
if [ ! -z "$BALANCED" ]; then
    echo "$BALANCED" >> $SUMMARY_FILE
fi
echo "" >> $SUMMARY_FILE

echo "**Aggressive (Maximum Returns):**" >> $SUMMARY_FILE
AGGRESSIVE=$(jq -r 'sort_by(-.returnPercent) | .[0] | "- Use **\(.days) days** period\n- Maximum returns: \(.returnPercent | floor)%\n- Expected DD: \(.maxDrawdown | floor)%\n- Expected PF: \(.profitFactor | tonumber | . * 100 | floor / 100)\n- Expected WR: \(.winRate | floor)%"' $OUTPUT_FILE 2>/dev/null)
if [ ! -z "$AGGRESSIVE" ]; then
    echo "$AGGRESSIVE" >> $SUMMARY_FILE
fi
echo "" >> $SUMMARY_FILE

echo ""
echo "✅ Analysis complete!"
echo ""
echo "📊 Results saved to:"
echo "  - $OUTPUT_FILE (raw data)"
echo "  - $SUMMARY_FILE (formatted summary)"
echo ""
echo "🎉 Check $SUMMARY_FILE for detailed analysis and recommendations!"
