#!/bin/bash

# ============================================
# SESSION TRADER - SIMPLE BACKTEST
# ============================================

echo "🚀 SESSION TRADER - SIMPLE BACKTEST"
echo "====================================="
echo ""

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Check backend
echo "📡 Checking backend..."
if ! curl -s http://localhost:8080/api/v1/health > /dev/null 2>&1; then
    echo -e "${RED}❌ Backend not running${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Backend running${NC}"
echo ""

# Function to run backtest
run_test() {
    local days=$1
    local desc=$2
    
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${YELLOW}📊 $desc ($days days)${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    
    result=$(curl -s -X POST "http://localhost:8080/api/v1/backtest/run" \
        -H "Content-Type: application/json" \
        -d "{
            \"symbol\": \"BTCUSDT\",
            \"interval\": \"15m\",
            \"days\": $days,
            \"strategy\": \"session_trader\",
            \"startBalance\": 1000
        }")
    
    # Check if we got valid data
    total_trades=$(echo "$result" | jq -r '.totalTrades // 0')
    
    if [ "$total_trades" = "0" ]; then
        echo -e "${RED}❌ No trades generated${NC}"
        echo "Response: $result" | jq '.'
        return
    fi
    
    # Extract metrics
    wins=$(echo "$result" | jq -r '.winningTrades')
    losses=$(echo "$result" | jq -r '.losingTrades')
    win_rate=$(echo "$result" | jq -r '.winRate')
    pf=$(echo "$result" | jq -r '.profitFactor')
    dd=$(echo "$result" | jq -r '.maxDrawdown')
    ret=$(echo "$result" | jq -r '.returnPercent')
    final=$(echo "$result" | jq -r '.finalBalance')
    
    # Display
    echo ""
    echo "📈 RESULTS:"
    echo "  Total Trades:    $total_trades"
    echo "  Wins/Losses:     $wins W / $losses L"
    printf "  Win Rate:        %.2f%%\n" "$win_rate"
    printf "  Profit Factor:   %.2f\n" "$pf"
    printf "  Max Drawdown:    %.2f%%\n" "$dd"
    printf "  Return:          %.2f%%\n" "$ret"
    printf "  Final Balance:   \$%.2f\n" "$final"
    
    # Rating
    echo ""
    if (( $(echo "$win_rate >= 50" | bc -l) )) && (( $(echo "$pf >= 2.5" | bc -l) )); then
        echo -e "${GREEN}✅ EXCELLENT${NC}"
    elif (( $(echo "$win_rate >= 45" | bc -l) )) && (( $(echo "$pf >= 2.0" | bc -l) )); then
        echo -e "${YELLOW}⚠️  GOOD${NC}"
    else
        echo -e "${RED}❌ NEEDS WORK${NC}"
    fi
    echo ""
}

# Run tests
run_test 7 "Last Week"
run_test 14 "Last 2 Weeks"
run_test 30 "Last Month"
run_test 60 "Last 2 Months"
run_test 90 "Last 3 Months"

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}📋 SESSION TRADER OVERVIEW${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo "Strategy:          Session Trader"
echo "Timeframe:         15m"
echo "Type:              Multi-Timeframe + Smart Money"
echo ""
echo "Features:"
echo "  ✓ AMD Phase Detection"
echo "  ✓ Market Regime Adaptive"
echo "  ✓ Smart Money Concepts"
echo "  ✓ 7 BUY + 7 SELL Strategies"
echo ""
echo "Expected (30d):"
echo "  Trades:          ~81"
echo "  Win Rate:        ~49.4%"
echo "  Profit Factor:   ~2.82"
echo "  Max Drawdown:    ~34.6%"
echo ""
echo -e "${GREEN}✅ Complete!${NC}"
