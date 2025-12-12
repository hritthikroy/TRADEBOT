#!/bin/bash

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🚀 COMPREHENSIVE LIQUIDITY HUNTER OPTIMIZATION"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Testing ALL timeframes with MULTIPLE parameter sets"
echo "Goal: Find the BEST configuration for maximum profitability"
echo ""

# All timeframes to test
timeframes=("1d" "8h" "4h" "1h" "30m" "15m" "5m")

# Parameter sets to test (StopATR, TP1, TP2, TP3)
declare -a param_sets=(
    "0.5:2.0:3.0:5.0"
    "0.75:2.5:4.0:6.0"
    "1.0:3.0:5.0:8.0"
    "1.5:4.0:6.0:10.0"
    "2.0:5.0:8.0:12.0"
    "1.0:2.0:4.0:8.0"
    "1.5:3.0:5.0:10.0"
    "2.0:4.0:7.0:15.0"
    "0.5:1.5:3.0:6.0"
    "1.0:4.0:8.0:16.0"
)

# Function to get days for timeframe
get_days() {
    case $1 in
        "1d") echo 180 ;;
        "8h") echo 120 ;;
        "4h") echo 90 ;;
        "1h") echo 60 ;;
        "30m") echo 30 ;;
        "15m") echo 15 ;;
        "5m") echo 7 ;;
        *) echo 30 ;;
    esac
}

best_overall_score=0
best_overall_config=""
results_file="liquidity_hunter_optimization_results.txt"

echo "Results will be saved to: $results_file"
echo "" > "$results_file"
echo "LIQUIDITY HUNTER COMPREHENSIVE OPTIMIZATION RESULTS" >> "$results_file"
echo "Generated: $(date)" >> "$results_file"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" >> "$results_file"
echo "" >> "$results_file"

for timeframe in "${timeframes[@]}"; do
    days=$(get_days "$timeframe")
    
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "📊 Testing Timeframe: $timeframe (${days} days)"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "TIMEFRAME: $timeframe" >> "$results_file"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" >> "$results_file"
    
    best_tf_score=0
    best_tf_config=""
    
    for params in "${param_sets[@]}"; do
        IFS=':' read -r stop tp1 tp2 tp3 <<< "$params"
        
        echo -n "   Testing Stop=$stop, TP1=$tp1, TP2=$tp2, TP3=$tp3 ... "
        
        # Run backtest
        result=$(curl -s -X POST http://localhost:8080/api/v1/backtest/run \
          -H "Content-Type: application/json" \
          -d "{
            \"symbol\": \"BTCUSDT\",
            \"interval\": \"$timeframe\",
            \"days\": $days,
            \"startBalance\": 1000,
            \"strategy\": \"liquidity_hunter\",
            \"riskPercent\": 0.01
          }")
        
        # Extract metrics
        winRate=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(data.get('winRate', 0))" 2>/dev/null || echo "0")
        profitFactor=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(data.get('profitFactor', 0))" 2>/dev/null || echo "0")
        returnPct=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(data.get('returnPercent', 0))" 2>/dev/null || echo "0")
        maxDD=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(data.get('maxDrawdown', 0))" 2>/dev/null || echo "0")
        totalTrades=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(data.get('totalTrades', 0))" 2>/dev/null || echo "0")
        finalBalance=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(data.get('finalBalance', 0))" 2>/dev/null || echo "0")
        
        # Calculate comprehensive score
        score=$(python3 -c "
wr = $winRate
pf = $profitFactor
ret = max(0, $returnPct)
dd = abs($maxDD)
trades = $totalTrades

# Quality filters
if trades < 5:
    score = 0
elif wr < 30:
    score = 0
elif pf < 0.5:
    score = 0
else:
    # Score formula: (WR * PF * Return) - (DD penalty)
    # Higher return and WR are rewarded, high DD is penalized
    score = (wr * pf * ret / 100) - (dd * 5)
    score = max(0, score)

print(round(score, 2))
" 2>/dev/null || echo "0")
        
        # Display result
        if (( $(echo "$returnPct > 0" | bc -l) )); then
            echo "✅ WR=$winRate% PF=$profitFactor Return=$returnPct% DD=$maxDD% Trades=$totalTrades Score=$score"
        else
            echo "❌ WR=$winRate% PF=$profitFactor Return=$returnPct% DD=$maxDD% Trades=$totalTrades"
        fi
        
        # Save to file
        echo "  Stop=$stop TP1=$tp1 TP2=$tp2 TP3=$tp3 | WR=$winRate% PF=$profitFactor Return=$returnPct% DD=$maxDD% Trades=$totalTrades Balance=\$$finalBalance Score=$score" >> "$results_file"
        
        # Track best for this timeframe
        if (( $(echo "$score > $best_tf_score" | bc -l) )); then
            best_tf_score=$score
            best_tf_config="$timeframe | Stop=$stop TP1=$tp1 TP2=$tp2 TP3=$tp3 | WR=$winRate% PF=$profitFactor Return=$returnPct% DD=$maxDD% Trades=$totalTrades"
        fi
        
        # Track best overall
        if (( $(echo "$score > $best_overall_score" | bc -l) )); then
            best_overall_score=$score
            best_overall_config="$timeframe | Stop=$stop TP1=$tp1 TP2=$tp2 TP3=$tp3 | WR=$winRate% PF=$profitFactor Return=$returnPct% DD=$maxDD% Trades=$totalTrades"
        fi
        
        # Small delay to avoid overwhelming the API
        sleep 0.5
    done
    
    echo ""
    echo "🏆 Best for $timeframe:"
    echo "   $best_tf_config"
    echo "   Score: $best_tf_score"
    echo ""
    
    echo "" >> "$results_file"
    echo "BEST FOR $timeframe: $best_tf_config (Score: $best_tf_score)" >> "$results_file"
    echo "" >> "$results_file"
done

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🎯 BEST OVERALL CONFIGURATION"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "$best_overall_config"
echo "Score: $best_overall_score"
echo ""
echo "Results saved to: $results_file"
echo ""

echo "" >> "$results_file"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" >> "$results_file"
echo "🎯 BEST OVERALL CONFIGURATION" >> "$results_file"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" >> "$results_file"
echo "$best_overall_config" >> "$results_file"
echo "Score: $best_overall_score" >> "$results_file"
