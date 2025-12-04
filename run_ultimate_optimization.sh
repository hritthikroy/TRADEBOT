#!/bin/bash

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🏆 ULTIMATE WORLD-CLASS PARAMETER OPTIMIZATION"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Testing ALL possible parameter combinations for ALL 10 strategies"
echo "This will test THOUSANDS of combinations to find the BEST parameters"
echo ""
echo "⚠️  WARNING: This will take 2-4 HOURS to complete!"
echo "⚠️  Make sure your computer won't sleep during this time"
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

read -p "Continue with ultimate optimization? (y/n) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Optimization cancelled"
    exit 0
fi

echo ""
echo "🚀 Starting ultimate optimization..."
echo ""

# Create results directory
mkdir -p ultimate_optimization_results
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
RESULTS_FILE="ultimate_optimization_results/ultimate_results_${TIMESTAMP}.json"
REPORT_FILE="ultimate_optimization_results/ultimate_report_${TIMESTAMP}.md"

# Initialize results
echo "{" > "$RESULTS_FILE"
echo "  \"timestamp\": \"$(date)\"," >> "$RESULTS_FILE"
echo "  \"strategies\": {" >> "$RESULTS_FILE"

# Initialize report
cat > "$REPORT_FILE" << 'EOF'
# 🏆 ULTIMATE WORLD-CLASS OPTIMIZATION RESULTS

## Optimization Date: 
EOF
date >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"
echo "## 🎯 OBJECTIVE" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"
echo "Find the ABSOLUTE BEST parameters for each strategy through exhaustive testing." >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"
echo "## 📊 TESTING METHODOLOGY" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"
echo "- **Parameter Combinations**: 10,000+ per strategy" >> "$REPORT_FILE"
echo "- **Test Period**: 180 days of BTCUSDT data" >> "$REPORT_FILE"
echo "- **Optimization Goal**: Maximize (Win Rate × Profit Factor × Return)" >> "$REPORT_FILE"
echo "- **Minimum Requirements**: 40% WR, 1.5 PF, 20+ trades" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"
echo "## 🏆 RESULTS" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"

# All 10 strategies
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

strategy_count=0
total_strategies=${#strategies[@]}

# Test each strategy
for strategy in "${strategies[@]}"; do
    strategy_count=$((strategy_count + 1))
    
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "[$strategy_count/$total_strategies] Optimizing: $strategy"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "Testing thousands of parameter combinations..."
    echo "This will take 10-20 minutes per strategy..."
    echo ""
    
    # Run world-class optimization for this strategy
    result=$(curl -s -X POST http://localhost:8080/api/v1/backtest/world-class-optimize \
      -H "Content-Type: application/json" \
      -d "{
        \"symbol\": \"BTCUSDT\",
        \"days\": 180,
        \"startBalance\": 1000,
        \"strategies\": [\"$strategy\"]
      }")
    
    # Extract best result
    bestWR=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    if 'results' in data and '$strategy' in data['results']:
        result = data['results']['$strategy']
        if 'backtestResult' in result and result['backtestResult']:
            print(f\"{result['backtestResult'].get('winRate', 0):.2f}\")
        else:
            print('0')
    else:
        print('0')
except:
    print('0')
" 2>/dev/null || echo "0")
    
    bestPF=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    if 'results' in data and '$strategy' in data['results']:
        result = data['results']['$strategy']
        if 'backtestResult' in result and result['backtestResult']:
            print(f\"{result['backtestResult'].get('profitFactor', 0):.2f}\")
        else:
            print('0')
    else:
        print('0')
except:
    print('0')
" 2>/dev/null || echo "0")
    
    bestReturn=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    if 'results' in data and '$strategy' in data['results']:
        result = data['results']['$strategy']
        if 'backtestResult' in result and result['backtestResult']:
            print(f\"{result['backtestResult'].get('returnPercent', 0):.0f}\")
        else:
            print('0')
    else:
        print('0')
except:
    print('0')
" 2>/dev/null || echo "0")
    
    bestTrades=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    if 'results' in data and '$strategy' in data['results']:
        result = data['results']['$strategy']
        if 'backtestResult' in result and result['backtestResult']:
            print(result['backtestResult'].get('totalTrades', 0))
        else:
            print('0')
    else:
        print('0')
except:
    print('0')
" 2>/dev/null || echo "0")
    
    stopATR=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    if 'results' in data and '$strategy' in data['results']:
        result = data['results']['$strategy']
        if 'bestParams' in result:
            print(result['bestParams'].get('StopATR', 0))
        else:
            print('0')
    else:
        print('0')
except:
    print('0')
" 2>/dev/null || echo "0")
    
    tp1ATR=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    if 'results' in data and '$strategy' in data['results']:
        result = data['results']['$strategy']
        if 'bestParams' in result:
            print(result['bestParams'].get('TP1ATR', 0))
        else:
            print('0')
    else:
        print('0')
except:
    print('0')
" 2>/dev/null || echo "0")
    
    tp2ATR=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    if 'results' in data and '$strategy' in data['results']:
        result = data['results']['$strategy']
        if 'bestParams' in result:
            print(result['bestParams'].get('TP2ATR', 0))
        else:
            print('0')
    else:
        print('0')
except:
    print('0')
" 2>/dev/null || echo "0")
    
    tp3ATR=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    if 'results' in data and '$strategy' in data['results']:
        result = data['results']['$strategy']
        if 'bestParams' in result:
            print(result['bestParams'].get('TP3ATR', 0))
        else:
            print('0')
    else:
        print('0')
except:
    print('0')
" 2>/dev/null || echo "0")
    
    riskPct=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    if 'results' in data and '$strategy' in data['results']:
        result = data['results']['$strategy']
        if 'bestParams' in result:
            print(result['bestParams'].get('RiskPercent', 0))
        else:
            print('0')
    else:
        print('0')
except:
    print('0')
" 2>/dev/null || echo "0")
    
    # Display results
    echo "✅ Optimization Complete!"
    echo ""
    echo "📊 BEST PARAMETERS FOUND:"
    echo "   Stop Loss: ${stopATR} ATR"
    echo "   TP1: ${tp1ATR} ATR"
    echo "   TP2: ${tp2ATR} ATR"
    echo "   TP3: ${tp3ATR} ATR"
    echo "   Risk: ${riskPct}%"
    echo ""
    echo "📈 PERFORMANCE:"
    echo "   Win Rate: ${bestWR}%"
    echo "   Profit Factor: ${bestPF}"
    echo "   Return: ${bestReturn}%"
    echo "   Total Trades: ${bestTrades}"
    echo ""
    
    # Add to report
    echo "### $strategy" >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"
    echo "**Performance:**" >> "$REPORT_FILE"
    echo "- Win Rate: ${bestWR}%" >> "$REPORT_FILE"
    echo "- Profit Factor: ${bestPF}" >> "$REPORT_FILE"
    echo "- Return: ${bestReturn}%" >> "$REPORT_FILE"
    echo "- Total Trades: ${bestTrades}" >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"
    echo "**Best Parameters:**" >> "$REPORT_FILE"
    echo "\`\`\`" >> "$REPORT_FILE"
    echo "Stop Loss: ${stopATR} ATR" >> "$REPORT_FILE"
    echo "TP1: ${tp1ATR} ATR (33%)" >> "$REPORT_FILE"
    echo "TP2: ${tp2ATR} ATR (33%)" >> "$REPORT_FILE"
    echo "TP3: ${tp3ATR} ATR (34%)" >> "$REPORT_FILE"
    echo "Risk per Trade: ${riskPct}%" >> "$REPORT_FILE"
    echo "\`\`\`" >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"
    
    # Grade the strategy
    if (( $(echo "$bestWR >= 50" | bc -l) )) && (( $(echo "$bestPF >= 3.0" | bc -l) )); then
        echo "🏆 Grade: WORLD-CLASS - Exceptional performance!" >> "$REPORT_FILE"
    elif (( $(echo "$bestWR >= 45" | bc -l) )) && (( $(echo "$bestPF >= 2.0" | bc -l) )); then
        echo "✅ Grade: EXCELLENT - Ready for real trading" >> "$REPORT_FILE"
    elif (( $(echo "$bestWR >= 40" | bc -l) )) && (( $(echo "$bestPF >= 1.5" | bc -l) )); then
        echo "✅ Grade: GOOD - Ready with caution" >> "$REPORT_FILE"
    else
        echo "⚠️ Grade: NEEDS IMPROVEMENT" >> "$REPORT_FILE"
    fi
    echo "" >> "$REPORT_FILE"
    echo "---" >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"
    
    # Add to JSON results
    if [ $strategy_count -lt $total_strategies ]; then
        echo "    \"$strategy\": $result," >> "$RESULTS_FILE"
    else
        echo "    \"$strategy\": $result" >> "$RESULTS_FILE"
    fi
done

# Close JSON
echo "  }" >> "$RESULTS_FILE"
echo "}" >> "$RESULTS_FILE"

# Add summary to report
echo "## 📊 SUMMARY" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"
echo "### Top 3 Strategies (By Score)" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"
echo "Review the results above and select the top 3 strategies with:" >> "$REPORT_FILE"
echo "- Win Rate >45%" >> "$REPORT_FILE"
echo "- Profit Factor >2.0" >> "$REPORT_FILE"
echo "- Return >100%" >> "$REPORT_FILE"
echo "- Total Trades >20" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"
echo "### Recommended Next Steps:" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"
echo "1. **Review this report carefully**" >> "$REPORT_FILE"
echo "2. **Update unified_signal_generator.go with best parameters**" >> "$REPORT_FILE"
echo "3. **Run comprehensive validation**" >> "$REPORT_FILE"
echo "4. **Start paper trading for 30 days**" >> "$REPORT_FILE"
echo "5. **Only then consider real trading**" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"
echo "### ⚠️ IMPORTANT WARNINGS:" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"
echo "- These parameters are optimized on HISTORICAL data" >> "$REPORT_FILE"
echo "- Past performance does NOT guarantee future results" >> "$REPORT_FILE"
echo "- ALWAYS use proper risk management" >> "$REPORT_FILE"
echo "- NEVER risk more than 1-2% per trade" >> "$REPORT_FILE"
echo "- ALWAYS use stop losses" >> "$REPORT_FILE"
echo "- START with paper trading first" >> "$REPORT_FILE"

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🏆 ULTIMATE OPTIMIZATION COMPLETE!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📄 Results saved to:"
echo "   JSON: $RESULTS_FILE"
echo "   Report: $REPORT_FILE"
echo ""
echo "To view the report:"
echo "   cat $REPORT_FILE"
echo ""
echo "🎯 Next Steps:"
echo "   1. Review the report"
echo "   2. Update parameters in code"
echo "   3. Run comprehensive validation"
echo "   4. Start paper trading"
echo ""
echo "⚠️  DO NOT trade real money until paper trading is successful!"
echo ""
