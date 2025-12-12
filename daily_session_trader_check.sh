#!/bin/bash

# ============================================
# DAILY SESSION TRADER PERFORMANCE CHECK
# ============================================
# Run this daily to check if strategy is profitable
# Only trade when this shows GREEN status

echo "📊 SESSION TRADER - DAILY PERFORMANCE CHECK"
echo "==========================================="
echo ""
echo "Date: $(date '+%Y-%m-%d %H:%M:%S')"
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
    echo "Start backend: cd backend && ./tradebot"
    exit 1
fi

# Run 7-day backtest
echo "Running 7-day backtest..."
result=$(curl -s -X POST "http://localhost:8080/api/v1/backtest/run" \
    -H "Content-Type: application/json" \
    -d '{"symbol":"BTCUSDT","interval":"15m","days":7,"strategy":"session_trader","startBalance":1000}')

# Extract metrics
trades=$(echo "$result" | jq -r '.totalTrades')
wins=$(echo "$result" | jq -r '.winningTrades')
losses=$(echo "$result" | jq -r '.losingTrades')
win_rate=$(echo "$result" | jq -r '.winRate')
pf=$(echo "$result" | jq -r '.profitFactor')
return=$(echo "$result" | jq -r '.returnPercent')
dd=$(echo "$result" | jq -r '.maxDrawdown')
final=$(echo "$result" | jq -r '.finalBalance')

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}📈 7-DAY PERFORMANCE${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
printf "  Total Trades:    %d\n" "$trades"
printf "  Wins/Losses:     %dW / %dL\n" "$wins" "$losses"
printf "  Win Rate:        %.2f%%\n" "$win_rate"
printf "  Profit Factor:   %.2f\n" "$pf"
printf "  Return:          %.2f%%\n" "$return"
printf "  Max Drawdown:    %.2f%%\n" "$dd"
printf "  Final Balance:   \$%.2f\n" "$final"
echo ""

# Determine trading status
profit=$(echo "$final - 1000" | bc)
is_profitable=$(echo "$return > 0" | bc -l)
good_wr=$(echo "$win_rate >= 45" | bc -l)
good_pf=$(echo "$pf >= 1.0" | bc -l)

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}🎯 TRADING DECISION${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Check all criteria
if [ "$is_profitable" = "1" ] && [ "$good_wr" = "1" ] && [ "$good_pf" = "1" ]; then
    echo -e "${GREEN}✅ STATUS: TRADE TODAY${NC}"
    echo ""
    echo "All criteria met:"
    echo -e "  ✅ Positive return (%.2f%%)\n" "$return"
    echo -e "  ✅ Win rate ≥ 45%% (%.2f%%)\n" "$win_rate"
    echo -e "  ✅ Profit factor ≥ 1.0 (%.2f)\n" "$pf"
    echo ""
    echo -e "${GREEN}🟢 RECOMMENDATION: Safe to trade today${NC}"
    echo ""
    echo "Suggested actions:"
    echo "  1. Start paper trading if not already running"
    echo "  2. Monitor trades closely"
    echo "  3. Use small position sizes (0.5-1% risk)"
    echo "  4. Set stop losses as per strategy"
    
elif [ "$is_profitable" = "1" ] && [ "$good_pf" = "1" ]; then
    echo -e "${YELLOW}⚠️  STATUS: TRADE WITH CAUTION${NC}"
    echo ""
    echo "Partially met:"
    echo -e "  ✅ Positive return (%.2f%%)\n" "$return"
    echo -e "  ✅ Profit factor ≥ 1.0 (%.2f)\n" "$pf"
    echo -e "  ⚠️  Win rate < 45%% (%.2f%%)\n" "$win_rate"
    echo ""
    echo -e "${YELLOW}🟡 RECOMMENDATION: Trade with reduced size${NC}"
    echo ""
    echo "Suggested actions:"
    echo "  1. Reduce position size by 50%"
    echo "  2. Be extra selective with entries"
    echo "  3. Monitor closely"
    
else
    echo -e "${RED}❌ STATUS: DO NOT TRADE${NC}"
    echo ""
    echo "Criteria not met:"
    if [ "$is_profitable" != "1" ]; then
        echo -e "  ❌ Negative return (%.2f%%)\n" "$return"
    fi
    if [ "$good_wr" != "1" ]; then
        echo -e "  ❌ Win rate < 45%% (%.2f%%)\n" "$win_rate"
    fi
    if [ "$good_pf" != "1" ]; then
        echo -e "  ❌ Profit factor < 1.0 (%.2f)\n" "$pf"
    fi
    echo ""
    echo -e "${RED}🔴 RECOMMENDATION: Pause trading today${NC}"
    echo ""
    echo "Suggested actions:"
    echo "  1. Stop paper trading if running"
    echo "  2. Wait for better market conditions"
    echo "  3. Check again tomorrow"
    echo "  4. Review recent trades for issues"
fi

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}📝 NOTES${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo "• Run this check every morning before trading"
echo "• Results are based on last 7 days of data"
echo "• Market conditions change - always monitor"
echo "• Start with paper trading to verify"
echo ""

# Save to log
log_file="session_trader_daily_checks.log"
echo "$(date '+%Y-%m-%d %H:%M:%S'),${trades},${win_rate},${pf},${return},${final}" >> "$log_file"
echo "Results saved to: $log_file"
echo ""
