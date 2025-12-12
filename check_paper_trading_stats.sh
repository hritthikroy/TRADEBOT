#!/bin/bash

# ============================================
# CHECK PAPER TRADING STATS
# ============================================

echo "📊 PAPER TRADING STATISTICS"
echo "==========================="
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

# Check backend
if ! curl -s http://localhost:8080/api/v1/health > /dev/null 2>&1; then
    echo -e "${RED}❌ Backend not running${NC}"
    exit 1
fi

# Get stats
stats=$(curl -s "http://localhost:8080/api/v1/paper-trading/stats")

# Extract data
is_running=$(echo "$stats" | jq -r '.isRunning')
total_trades=$(echo "$stats" | jq -r '.totalTrades')
open_trades=$(echo "$stats" | jq -r '.openTrades')
closed_trades=$(echo "$stats" | jq -r '.closedTrades')
wins=$(echo "$stats" | jq -r '.wins')
losses=$(echo "$stats" | jq -r '.losses')
win_rate=$(echo "$stats" | jq -r '.winRate')
total_profit=$(echo "$stats" | jq -r '.totalProfit')
total_loss=$(echo "$stats" | jq -r '.totalLoss')
net_profit=$(echo "$stats" | jq -r '.netProfit')
profit_factor=$(echo "$stats" | jq -r '.profitFactor')
current_balance=$(echo "$stats" | jq -r '.currentBalance')
start_balance=$(echo "$stats" | jq -r '.startBalance')
return_pct=$(echo "$stats" | jq -r '.returnPercent')

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}📈 OVERALL PERFORMANCE${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

if [ "$is_running" = "true" ]; then
    echo -e "Status:           ${GREEN}🟢 RUNNING${NC}"
else
    echo -e "Status:           ${RED}🔴 STOPPED${NC}"
fi

printf "Start Balance:    \$%.2f\n" "$start_balance"
printf "Current Balance:  \$%.2f\n" "$current_balance"
printf "Net Profit:       \$%.2f\n" "$net_profit"
printf "Return:           %.2f%%\n" "$return_pct"
echo ""

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}📊 TRADE STATISTICS${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

printf "Total Trades:     %d\n" "$total_trades"
printf "Open Trades:      %d\n" "$open_trades"
printf "Closed Trades:    %d\n" "$closed_trades"
printf "Wins:             %d\n" "$wins"
printf "Losses:           %d\n" "$losses"
printf "Win Rate:         %.2f%%\n" "$win_rate"
printf "Profit Factor:    %.2f\n" "$profit_factor"
echo ""

# Performance rating
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}🎯 PERFORMANCE RATING${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

is_profitable=$(echo "$return_pct > 0" | bc -l)
good_wr=$(echo "$win_rate >= 45" | bc -l)
good_pf=$(echo "$profit_factor >= 1.0" | bc -l)

if [ "$is_profitable" = "1" ] && [ "$good_wr" = "1" ] && [ "$good_pf" = "1" ]; then
    echo -e "${GREEN}✅ EXCELLENT - All targets met${NC}"
    echo ""
    echo "  ✅ Profitable"
    echo "  ✅ Win rate ≥ 45%"
    echo "  ✅ Profit factor ≥ 1.0"
    echo ""
    echo -e "${GREEN}Continue trading with current settings${NC}"
    
elif [ "$is_profitable" = "1" ] && [ "$good_pf" = "1" ]; then
    echo -e "${YELLOW}⚠️  GOOD - Profitable but low win rate${NC}"
    echo ""
    echo "  ✅ Profitable"
    echo "  ⚠️  Win rate < 45%"
    echo "  ✅ Profit factor ≥ 1.0"
    echo ""
    echo -e "${YELLOW}Monitor closely, consider reducing size${NC}"
    
elif [ "$closed_trades" -lt "10" ]; then
    echo -e "${BLUE}ℹ️  INSUFFICIENT DATA${NC}"
    echo ""
    echo "  Need at least 10 closed trades for accurate assessment"
    echo "  Current closed trades: $closed_trades"
    echo ""
    echo -e "${BLUE}Continue trading to gather more data${NC}"
    
else
    echo -e "${RED}❌ POOR - Not meeting targets${NC}"
    echo ""
    if [ "$is_profitable" != "1" ]; then
        echo "  ❌ Not profitable"
    fi
    if [ "$good_wr" != "1" ]; then
        echo "  ❌ Win rate < 45%"
    fi
    if [ "$good_pf" != "1" ]; then
        echo "  ❌ Profit factor < 1.0"
    fi
    echo ""
    echo -e "${RED}Consider stopping and reviewing strategy${NC}"
fi

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Get recent trades
echo "📋 Recent Trades (Last 5):"
echo ""

trades=$(curl -s "http://localhost:8080/api/v1/paper-trading/trades" | jq -r '.trades[-5:] | reverse | .[] | "  \(.type) @ $\(.entry | tonumber | floor) → \(.status) | P/L: $\(.profit | tonumber)"')

if [ -z "$trades" ]; then
    echo "  No trades yet"
else
    echo "$trades"
fi

echo ""
echo "View full dashboard: http://localhost:8080/paper-trading"
echo ""
