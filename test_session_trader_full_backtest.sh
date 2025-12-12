#!/bin/bash

# ============================================
# SESSION TRADER - COMPREHENSIVE BACKTEST
# ============================================
# Tests Session Trader strategy across multiple timeframes
# and provides detailed performance analysis

echo "🚀 SESSION TRADER - FULL BACKTEST ANALYSIS"
echo "=========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# API endpoint
API_URL="http://localhost:8080/api/v1/backtest/test-all-strategies"

# Check if backend is running
echo "📡 Checking backend status..."
if ! curl -s http://localhost:8080/api/v1/health > /dev/null 2>&1; then
    echo -e "${RED}❌ Backend is not running!${NC}"
    echo "Please start the backend first:"
    echo "  cd backend && go run ."
    exit 1
fi
echo -e "${GREEN}✅ Backend is running${NC}"
echo ""

# Function to run backtest
run_backtest() {
    local days=$1
    local description=$2
    
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${YELLOW}📊 Testing: $description ($days days)${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    
    response=$(curl -s -X POST "$API_URL" \
        -H "Content-Type: application/json" \
        -d "{
            \"symbol\": \"BTCUSDT\",
            \"days\": $days,
            \"startBalance\": 1000,
            \"filterBuy\": false,
            \"filterSell\": false,
            \"selectedStrategies\": [\"session_trader\"]
        }")
    
    # Extract Session Trader results
    strategy_data=$(echo "$response" | jq -r '.strategies[] | select(.strategy == "session_trader")')
    
    if [ -z "$strategy_data" ]; then
        echo -e "${RED}❌ No data returned${NC}"
        return
    fi
    
    # Extract metrics
    trades=$(echo "$strategy_data" | jq -r '.totalTrades')
    wins=$(echo "$strategy_data" | jq -r '.wins')
    losses=$(echo "$strategy_data" | jq -r '.losses')
    win_rate=$(echo "$strategy_data" | jq -r '.winRate')
    profit_factor=$(echo "$strategy_data" | jq -r '.profitFactor')
    max_dd=$(echo "$strategy_data" | jq -r '.maxDrawdown')
    final_balance=$(echo "$strategy_data" | jq -r '.finalBalance')
    total_return=$(echo "$strategy_data" | jq -r '.totalReturn')
    
    # Calculate profit/loss
    profit=$(echo "$final_balance - 1000" | bc)
    
    # Display results
    echo ""
    echo "📈 PERFORMANCE METRICS:"
    echo "  Total Trades:    $trades"
    echo "  Wins/Losses:     $wins W / $losses L"
    printf "  Win Rate:        %.2f%%\n" "$win_rate"
    printf "  Profit Factor:   %.2f\n" "$profit_factor"
    printf "  Max Drawdown:    %.2f%%\n" "$max_dd"
    printf "  Total Return:    %.2f%%\n" "$total_return"
    printf "  Final Balance:   \$%.2f\n" "$final_balance"
    printf "  Profit/Loss:     \$%.2f\n" "$profit"
    
    # Performance rating
    echo ""
    if (( $(echo "$win_rate >= 50" | bc -l) )) && (( $(echo "$profit_factor >= 2.5" | bc -l) )); then
        echo -e "${GREEN}✅ EXCELLENT PERFORMANCE${NC}"
    elif (( $(echo "$win_rate >= 45" | bc -l) )) && (( $(echo "$profit_factor >= 2.0" | bc -l) )); then
        echo -e "${YELLOW}⚠️  GOOD PERFORMANCE${NC}"
    else
        echo -e "${RED}❌ NEEDS IMPROVEMENT${NC}"
    fi
    
    echo ""
}

# ============================================
# RUN COMPREHENSIVE BACKTESTS
# ============================================

echo "🔍 Running comprehensive backtests..."
echo ""

# Test 1: Last 7 days (Recent performance)
run_backtest 7 "Recent Week"

# Test 2: Last 14 days (Two weeks)
run_backtest 14 "Two Weeks"

# Test 3: Last 30 days (One month - Standard)
run_backtest 30 "One Month (Standard)"

# Test 4: Last 60 days (Two months)
run_backtest 60 "Two Months"

# Test 5: Last 90 days (Three months - Long term)
run_backtest 90 "Three Months (Long Term)"

# ============================================
# DETAILED ANALYSIS
# ============================================

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}📊 DETAILED STRATEGY ANALYSIS${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

echo "🎯 SESSION TRADER STRATEGY OVERVIEW:"
echo ""
echo "Strategy Type:     Multi-Timeframe + Smart Money Concepts"
echo "Timeframe:         15m"
echo "Target Win Rate:   55-65%"
echo "Target PF:         3.5-5.0"
echo "Target Drawdown:   <12%"
echo ""

echo "🔧 KEY FEATURES:"
echo "  ✓ AMD Phase Detection (Manipulation filtering)"
echo "  ✓ Market Regime Adaptive (Bull/Bear/Sideways)"
echo "  ✓ Smart Money Concepts (Order Blocks)"
echo "  ✓ Multi-EMA Trend Analysis"
echo "  ✓ Volume Analysis (Institutional detection)"
echo "  ✓ 7 BUY Strategies (Bull/Sideways markets)"
echo "  ✓ 7 SELL Strategies (Bear/Sideways markets)"
echo ""

echo "📈 BUY STRATEGIES:"
echo "  1. Strong Trend Following"
echo "  2. Order Block Bounce"
echo "  3. Momentum Breakout"
echo "  4. Pullback Entry"
echo "  5. EMA Bounce"
echo "  6. Volume Spike Reversal"
echo "  7. Simple Trend + RSI"
echo ""

echo "📉 SELL STRATEGIES:"
echo "  1. Perfect Trend Following"
echo "  2. Order Block Rejection"
echo "  3. Momentum Breakdown"
echo "  4. Conservative Pullback"
echo "  5. Strong Downtrend + Volume"
echo "  6. EMA Rejection"
echo "  7. Volume Spike + Reversal"
echo ""

echo "🎲 RISK MANAGEMENT:"
echo "  Stop Loss:     1.0-1.5 × ATR"
echo "  Take Profit 1: 2.0-3.0 × ATR"
echo "  Take Profit 2: 3.0-4.5 × ATR"
echo "  Take Profit 3: 4.5-7.5 × ATR"
echo "  Risk/Reward:   2.5:1 to 6:1"
echo ""

# ============================================
# COMPARISON WITH DOCUMENTATION
# ============================================

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}📋 COMPARISON WITH DOCUMENTED RESULTS${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

echo "Expected Performance (30 days):"
echo "  Trades:          81"
echo "  Win Rate:        49.4%"
echo "  Profit Factor:   2.82"
echo "  Max Drawdown:    34.6%"
echo "  Wins/Losses:     40W/41L"
echo ""

echo "Key Improvements from Original:"
echo "  ✓ 58% reduction in bad trades (192 → 81)"
echo "  ✓ 38% better profit factor (2.05 → 2.82)"
echo "  ✓ 13% lower drawdown (39.9% → 34.6%)"
echo "  ✓ Better trade quality (more selective)"
echo ""

# ============================================
# RECOMMENDATIONS
# ============================================

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}💡 RECOMMENDATIONS${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

echo "✅ READY FOR LIVE TRADING IF:"
echo "  • Win Rate ≥ 45%"
echo "  • Profit Factor ≥ 2.0"
echo "  • Max Drawdown ≤ 40%"
echo "  • Consistent across timeframes"
echo ""

echo "⚠️  NEEDS OPTIMIZATION IF:"
echo "  • Win Rate < 40%"
echo "  • Profit Factor < 1.5"
echo "  • Max Drawdown > 50%"
echo "  • Inconsistent results"
echo ""

echo "🔧 OPTIMIZATION OPTIONS:"
echo "  1. Adjust market regime thresholds (70% → 65%)"
echo "  2. Fine-tune AMD phase detection"
echo "  3. Modify volume filters"
echo "  4. Adjust RSI ranges"
echo "  5. Test different ATR multipliers"
echo ""

# ============================================
# NEXT STEPS
# ============================================

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}🚀 NEXT STEPS${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

echo "1. Review Results:"
echo "   • Check if performance meets expectations"
echo "   • Compare with documented results"
echo "   • Identify any anomalies"
echo ""

echo "2. Test Specific Periods:"
echo "   • Test bad period (Nov 30 - Dec 4):"
echo "     ./test_nov30_dec4_period.sh"
echo ""

echo "3. Compare with Other Strategies:"
echo "   • Test all 10 strategies:"
echo "     ./test_all_10_strategies.sh"
echo ""

echo "4. Start Paper Trading:"
echo "   • If results are good, start paper trading:"
echo "     ./start_paper_trading.sh"
echo ""

echo "5. Optimize Further:"
echo "   • Run comprehensive optimization:"
echo "     ./run_ultimate_optimization.sh"
echo ""

echo -e "${GREEN}✅ Backtest complete!${NC}"
echo ""
