#!/bin/bash

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🔬 COMPREHENSIVE STRATEGY VALIDATION"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Testing all 10 strategies across multiple market conditions"
echo "This will take 10-15 minutes to complete..."
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

# Create results directory
mkdir -p validation_results
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
REPORT_FILE="validation_results/validation_report_${TIMESTAMP}.md"

# Initialize report
cat > "$REPORT_FILE" << 'EOF'
# 🔬 COMPREHENSIVE VALIDATION REPORT

## Test Date: 
EOF
date >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"

# Strategies to test
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

# Market periods to test
declare -A periods
periods["2024_bull"]="2024-01-01:2024-03-31:Bull Market +74%"
periods["2023_bull"]="2023-10-01:2023-12-31:Bull Market +63%"
periods["2023_range"]="2023-04-01:2023-09-30:Range Market"
periods["recent"]="recent:180:Recent 180 Days"

echo "📊 TESTING MATRIX" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"
echo "| Strategy | Recent | 2024 Bull | 2023 Bull | 2023 Range | Score |" >> "$REPORT_FILE"
echo "|----------|--------|-----------|-----------|------------|-------|" >> "$REPORT_FILE"

# Test each strategy
for strategy in "${strategies[@]}"; do
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "Testing: $strategy"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    
    # Start row in report
    echo -n "| $strategy | " >> "$REPORT_FILE"
    
    total_score=0
    period_count=0
    
    # Test each period
    for period_key in "recent" "2024_bull" "2023_bull" "2023_range"; do
        period_info="${periods[$period_key]}"
        IFS=':' read -r start end desc <<< "$period_info"
        
        echo "  Testing period: $desc"
        
        # Prepare request based on period type
        if [ "$start" = "recent" ]; then
            # Recent period - use days
            request_body="{
                \"symbol\": \"BTCUSDT\",
                \"interval\": \"15m\",
                \"days\": $end,
                \"startBalance\": 1000,
                \"riskPercent\": 2,
                \"strategy\": \"$strategy\"
            }"
        else
            # Historical period - use dates
            request_body="{
                \"symbol\": \"BTCUSDT\",
                \"interval\": \"15m\",
                \"startDate\": \"$start\",
                \"endDate\": \"$end\",
                \"startBalance\": 1000,
                \"riskPercent\": 2,
                \"strategy\": \"$strategy\"
            }"
        fi
        
        # Run backtest
        result=$(curl -s -X POST http://localhost:8080/api/v1/backtest/run \
          -H "Content-Type: application/json" \
          -d "$request_body")
        
        # Extract metrics
        trades=$(echo "$result" | python3 -c "import sys, json; data=json.load(sys.stdin); print(data.get('totalTrades', 0))" 2>/dev/null || echo "0")
        winRate=$(echo "$result" | python3 -c "import sys, json; data=json.load(sys.stdin); print(f\"{data.get('winRate', 0):.1f}\")" 2>/dev/null || echo "0")
        profitFactor=$(echo "$result" | python3 -c "import sys, json; data=json.load(sys.stdin); print(f\"{data.get('profitFactor', 0):.2f}\")" 2>/dev/null || echo "0")
        returnPct=$(echo "$result" | python3 -c "import sys, json; data=json.load(sys.stdin); print(f\"{data.get('returnPercent', 0):.0f}\")" 2>/dev/null || echo "0")
        
        # Calculate period score
        period_score=0
        
        # Win rate score (0-25 points)
        if (( $(echo "$winRate >= 40" | bc -l) )); then
            period_score=$((period_score + 25))
        elif (( $(echo "$winRate >= 35" | bc -l) )); then
            period_score=$((period_score + 15))
        fi
        
        # Profit factor score (0-25 points)
        if (( $(echo "$profitFactor >= 2.0" | bc -l) )); then
            period_score=$((period_score + 25))
        elif (( $(echo "$profitFactor >= 1.5" | bc -l) )); then
            period_score=$((period_score + 18))
        elif (( $(echo "$profitFactor >= 1.2" | bc -l) )); then
            period_score=$((period_score + 10))
        fi
        
        # Return score (0-25 points)
        if (( $(echo "$returnPct >= 50" | bc -l) )); then
            period_score=$((period_score + 25))
        elif (( $(echo "$returnPct >= 20" | bc -l) )); then
            period_score=$((period_score + 18))
        elif (( $(echo "$returnPct >= 0" | bc -l) )); then
            period_score=$((period_score + 10))
        fi
        
        # Trade frequency score (0-25 points)
        if [ "$trades" -ge 10 ] && [ "$trades" -le 100 ]; then
            period_score=$((period_score + 25))
        elif [ "$trades" -ge 5 ]; then
            period_score=$((period_score + 15))
        fi
        
        total_score=$((total_score + period_score))
        period_count=$((period_count + 1))
        
        # Add to report
        if [ "$returnPct" -gt 0 ]; then
            echo -n "✅ ${returnPct}% | " >> "$REPORT_FILE"
        else
            echo -n "❌ ${returnPct}% | " >> "$REPORT_FILE"
        fi
        
        echo "    Trades: $trades | WR: ${winRate}% | PF: $profitFactor | Return: ${returnPct}%"
        
        # Small delay to avoid overwhelming the server
        sleep 1
    done
    
    # Calculate average score
    avg_score=$((total_score / period_count))
    
    # Add final score to report
    echo "$avg_score |" >> "$REPORT_FILE"
    
    echo ""
    echo "  Average Score: $avg_score/100"
    
    # Grade the strategy
    if [ "$avg_score" -ge 80 ]; then
        echo "  Grade: ✅ EXCELLENT - Ready for real trading"
    elif [ "$avg_score" -ge 70 ]; then
        echo "  Grade: ✅ GOOD - Ready with caution"
    elif [ "$avg_score" -ge 60 ]; then
        echo "  Grade: ⚠️ ACCEPTABLE - Needs monitoring"
    else
        echo "  Grade: ❌ FAIL - Not ready for real trading"
    fi
    
    echo ""
done

# Add summary to report
echo "" >> "$REPORT_FILE"
echo "## 📊 SUMMARY" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"
echo "### Grading Scale:" >> "$REPORT_FILE"
echo "- **80-100**: ✅ EXCELLENT - Ready for real trading" >> "$REPORT_FILE"
echo "- **70-79**: ✅ GOOD - Ready with caution" >> "$REPORT_FILE"
echo "- **60-69**: ⚠️ ACCEPTABLE - Needs monitoring" >> "$REPORT_FILE"
echo "- **<60**: ❌ FAIL - Not ready for real trading" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"
echo "### Recommendations:" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"
echo "1. **Only trade strategies scoring 70+**" >> "$REPORT_FILE"
echo "2. **Start with smallest position sizes**" >> "$REPORT_FILE"
echo "3. **Always use stop losses**" >> "$REPORT_FILE"
echo "4. **Monitor performance weekly**" >> "$REPORT_FILE"
echo "5. **Stop trading after 3 consecutive losses**" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"
echo "### Next Steps:" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"
echo "1. Review this report carefully" >> "$REPORT_FILE"
echo "2. Select top 3 strategies" >> "$REPORT_FILE"
echo "3. Run live paper trading for 30 days" >> "$REPORT_FILE"
echo "4. Only then consider real trading" >> "$REPORT_FILE"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ VALIDATION COMPLETE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📄 Report saved to: $REPORT_FILE"
echo ""
echo "To view the report:"
echo "   cat $REPORT_FILE"
echo ""
echo "⚠️  IMPORTANT: Review the report before trading real money!"
echo ""
