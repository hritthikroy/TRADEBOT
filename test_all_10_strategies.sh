#!/bin/bash

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🏆 COMPLETE BACKTEST - ALL 10 STRATEGIES"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Testing all 10 strategies to find the BEST 3"
echo "⏱️  Estimated time: 5-10 minutes"
echo ""

# Check server
if ! curl -s http://localhost:8080/api/v1/health > /dev/null 2>&1; then
    echo "❌ Server not running!"
    exit 1
fi

echo "✅ Server is running"
echo ""

# Create results directory
mkdir -p complete_backtest_results
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
RESULTS_FILE="complete_backtest_results/complete_results_${TIMESTAMP}.md"

# Initialize report
cat > "$RESULTS_FILE" << 'EOF'
# 🏆 COMPLETE BACKTEST RESULTS - ALL 10 STRATEGIES

## Test Date: 
EOF
date >> "$RESULTS_FILE"
echo "" >> "$RESULTS_FILE"
echo "## 📊 RESULTS" >> "$RESULTS_FILE"
echo "" >> "$RESULTS_FILE"
echo "| Rank | Strategy | Win Rate | Profit Factor | Return % | Trades | Score |" >> "$RESULTS_FILE"
echo "|------|----------|----------|---------------|----------|--------|-------|" >> "$RESULTS_FILE"

# All 10 strategies
strategies=(
    "session_trader"
    "breakout_master"
    "liquidity_hunter"
    "trend_rider"
    "range_master"
    "smart_money_tracker"
    "institutional_follower"
    "reversal_sniper"
    "momentum_beast"
    "scalper_pro"
)

# Store results for ranking
declare -A results_wr
declare -A results_pf
declare -A results_return
declare -A results_trades
declare -A results_score

# Test each strategy
for strategy in "${strategies[@]}"; do
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "Testing: $strategy"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
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
    
    # Extract metrics
    trades=$(echo "$result" | python3 -c "import sys, json; data=json.load(sys.stdin); print(data.get('totalTrades', 0))" 2>/dev/null || echo "0")
    wr=$(echo "$result" | python3 -c "import sys, json; data=json.load(sys.stdin); print(f\"{data.get('winRate', 0):.2f}\")" 2>/dev/null || echo "0")
    pf=$(echo "$result" | python3 -c "import sys, json; data=json.load(sys.stdin); print(f\"{data.get('profitFactor', 0):.2f}\")" 2>/dev/null || echo "0")
    ret=$(echo "$result" | python3 -c "import sys, json; data=json.load(sys.stdin); print(f\"{data.get('returnPercent', 0):.0f}\")" 2>/dev/null || echo "0")
    
    # Calculate score: WR * PF * log(Return + 1)
    score=$(python3 -c "import math; print(f\"{float('$wr') * float('$pf') * math.log(max(float('$ret'), 1) + 1):.2f}\")" 2>/dev/null || echo "0")
    
    # Store results
    results_wr[$strategy]=$wr
    results_pf[$strategy]=$pf
    results_return[$strategy]=$ret
    results_trades[$strategy]=$trades
    results_score[$strategy]=$score
    
    echo "   Trades: $trades | WR: ${wr}% | PF: $pf | Return: ${ret}% | Score: $score"
    echo ""
done

# Sort strategies by score
sorted_strategies=($(for strategy in "${strategies[@]}"; do
    echo "${results_score[$strategy]} $strategy"
done | sort -rn | awk '{print $2}'))

# Add ranked results to report
rank=1
for strategy in "${sorted_strategies[@]}"; do
    wr=${results_wr[$strategy]}
    pf=${results_pf[$strategy]}
    ret=${results_return[$strategy]}
    trades=${results_trades[$strategy]}
    score=${results_score[$strategy]}
    
    # Add medal for top 3
    if [ $rank -eq 1 ]; then
        medal="🥇"
    elif [ $rank -eq 2 ]; then
        medal="🥈"
    elif [ $rank -eq 3 ]; then
        medal="🥉"
    else
        medal="  "
    fi
    
    echo "| $medal $rank | $strategy | ${wr}% | $pf | ${ret}% | $trades | $score |" >> "$RESULTS_FILE"
    rank=$((rank + 1))
done

# Add detailed analysis
echo "" >> "$RESULTS_FILE"
echo "## 🏆 TOP 3 STRATEGIES" >> "$RESULTS_FILE"
echo "" >> "$RESULTS_FILE"

# Top 3 details
for i in 0 1 2; do
    strategy="${sorted_strategies[$i]}"
    wr=${results_wr[$strategy]}
    pf=${results_pf[$strategy]}
    ret=${results_return[$strategy]}
    trades=${results_trades[$strategy]}
    score=${results_score[$strategy]}
    
    if [ $i -eq 0 ]; then
        echo "### 🥇 #1: $strategy (BEST)" >> "$RESULTS_FILE"
    elif [ $i -eq 1 ]; then
        echo "### 🥈 #2: $strategy" >> "$RESULTS_FILE"
    else
        echo "### 🥉 #3: $strategy" >> "$RESULTS_FILE"
    fi
    
    echo "" >> "$RESULTS_FILE"
    echo "**Performance:**" >> "$RESULTS_FILE"
    echo "- Win Rate: ${wr}%" >> "$RESULTS_FILE"
    echo "- Profit Factor: $pf" >> "$RESULTS_FILE"
    echo "- Return: ${ret}%" >> "$RESULTS_FILE"
    echo "- Total Trades: $trades" >> "$RESULTS_FILE"
    echo "- Score: $score" >> "$RESULTS_FILE"
    echo "" >> "$RESULTS_FILE"
    
    # Grade
    if (( $(echo "$wr >= 50" | bc -l) )) && (( $(echo "$pf >= 3" | bc -l) )); then
        echo "🏆 **Grade: WORLD-CLASS**" >> "$RESULTS_FILE"
    elif (( $(echo "$wr >= 45" | bc -l) )) && (( $(echo "$pf >= 2" | bc -l) )); then
        echo "✅ **Grade: EXCELLENT**" >> "$RESULTS_FILE"
    elif (( $(echo "$wr >= 40" | bc -l) )) && (( $(echo "$pf >= 1.5" | bc -l) )); then
        echo "✅ **Grade: GOOD**" >> "$RESULTS_FILE"
    else
        echo "⚠️ **Grade: NEEDS IMPROVEMENT**" >> "$RESULTS_FILE"
    fi
    echo "" >> "$RESULTS_FILE"
    echo "---" >> "$RESULTS_FILE"
    echo "" >> "$RESULTS_FILE"
done

# Add recommendations
echo "## 💡 RECOMMENDATIONS" >> "$RESULTS_FILE"
echo "" >> "$RESULTS_FILE"
echo "### Use These 3 Strategies:" >> "$RESULTS_FILE"
echo "" >> "$RESULTS_FILE"
echo "1. **${sorted_strategies[0]}** - Best overall performance" >> "$RESULTS_FILE"
echo "2. **${sorted_strategies[1]}** - Second best" >> "$RESULTS_FILE"
echo "3. **${sorted_strategies[2]}** - Third best" >> "$RESULTS_FILE"
echo "" >> "$RESULTS_FILE"
echo "### Portfolio Allocation:" >> "$RESULTS_FILE"
echo "" >> "$RESULTS_FILE"
echo "- 50% ${sorted_strategies[0]}" >> "$RESULTS_FILE"
echo "- 30% ${sorted_strategies[1]}" >> "$RESULTS_FILE"
echo "- 20% ${sorted_strategies[2]}" >> "$RESULTS_FILE"
echo "" >> "$RESULTS_FILE"
echo "### Next Steps:" >> "$RESULTS_FILE"
echo "" >> "$RESULTS_FILE"
echo "1. Update code with top 3 strategies" >> "$RESULTS_FILE"
echo "2. Run deep optimization on top 3 for best parameters" >> "$RESULTS_FILE"
echo "3. Start paper trading for 30 days" >> "$RESULTS_FILE"
echo "4. Go live (if successful)" >> "$RESULTS_FILE"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ COMPLETE BACKTEST FINISHED!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📄 Results saved to: $RESULTS_FILE"
echo ""
echo "🏆 TOP 3 STRATEGIES:"
echo "   1. ${sorted_strategies[0]}"
echo "   2. ${sorted_strategies[1]}"
echo "   3. ${sorted_strategies[2]}"
echo ""
echo "To view full report:"
echo "   cat $RESULTS_FILE"
echo ""
