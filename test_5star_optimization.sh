#!/bin/bash

# ============================================
# TEST 5-STAR OPTIMIZATION
# ============================================

echo "🌟 TESTING 5-STAR OPTIMIZATION"
echo "========================================"
echo ""

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}⚠️  IMPORTANT: You need to restart the backend for changes to take effect!${NC}"
echo ""
echo "Steps:"
echo "1. Stop the current backend (Ctrl+C in the terminal running it)"
echo "2. Restart: cd backend && go run ."
echo "3. Wait for it to start (you'll see 'Server starting on port 8080')"
echo "4. Run this script again"
echo ""
read -p "Have you restarted the backend? (y/n) " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    echo -e "${RED}Please restart the backend first!${NC}"
    exit 1
fi

echo -e "${BLUE}Testing optimized Session Trader...${NC}"
echo ""

# Test 30-day backtest
echo "Running 30-day backtest..."
RESULT=$(curl -s -X POST "http://localhost:8080/api/v1/backtest/run" \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTCUSDT","interval":"15m","days":30,"strategy":"session_trader","startBalance":1000}')

# Extract metrics
TRADES=$(echo "$RESULT" | jq -r '.totalTrades')
WR=$(echo "$RESULT" | jq -r '.winRate')
PF=$(echo "$RESULT" | jq -r '.profitFactor')
BALANCE=$(echo "$RESULT" | jq -r '.finalBalance')
DD=$(echo "$RESULT" | jq -r '.maxDrawdown')

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 RESULTS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
printf "Total Trades:    %d\n" "$TRADES"
printf "Win Rate:        %.2f%%\n" "$WR"
printf "Profit Factor:   %.2f\n" "$PF"
printf "Final Balance:   \$%.2f\n" "$BALANCE"
printf "Max Drawdown:    %.2f%%\n" "$DD"
printf "Return:          %.2f%%\n" "$(echo "($BALANCE - 1000) / 10" | bc -l)"
echo ""

# Check if optimizations are working
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ OPTIMIZATION CHECK"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check trades (should be 40-80 if optimized)
if [ "$TRADES" -lt 100 ]; then
    echo -e "${GREEN}✅ Trade count looks good ($TRADES < 100)${NC}"
    echo "   ADX filter and cooldown are working!"
else
    echo -e "${RED}❌ Too many trades ($TRADES)${NC}"
    echo "   ADX filter or cooldown may not be working"
    echo "   Expected: 40-80 trades"
fi
echo ""

# Check win rate (should be 50%+ if optimized)
if (( $(echo "$WR >= 50" | bc -l) )); then
    echo -e "${GREEN}✅ Win rate is good (${WR}% >= 50%)${NC}"
else
    echo -e "${YELLOW}⚠️  Win rate needs improvement (${WR}% < 50%)${NC}"
    echo "   Expected: 55-65%"
fi
echo ""

# Check profit factor (should be 2.0+ if optimized)
if (( $(echo "$PF >= 2.0" | bc -l) )); then
    echo -e "${GREEN}✅ Profit factor is excellent (${PF} >= 2.0)${NC}"
elif (( $(echo "$PF >= 1.5" | bc -l) )); then
    echo -e "${YELLOW}⚠️  Profit factor is good (${PF} >= 1.5)${NC}"
else
    echo -e "${RED}❌ Profit factor needs improvement (${PF} < 1.5)${NC}"
    echo "   Expected: 3.5-5.0"
fi
echo ""

# Overall rating
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🏆 RATING"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

if [ "$TRADES" -lt 100 ] && (( $(echo "$WR >= 55" | bc -l) )) && (( $(echo "$PF >= 3.5" | bc -l) )); then
    echo -e "${GREEN}⭐⭐⭐⭐⭐ (5/5) - EXCELLENT!${NC}"
    echo "Session Trader is now world-class!"
elif [ "$TRADES" -lt 100 ] && (( $(echo "$WR >= 50" | bc -l) )) && (( $(echo "$PF >= 2.5" | bc -l) )); then
    echo -e "${GREEN}⭐⭐⭐⭐ (4/5) - VERY GOOD!${NC}"
    echo "Session Trader is competitive with professional bots"
elif [ "$TRADES" -lt 120 ] && (( $(echo "$WR >= 45" | bc -l) )) && (( $(echo "$PF >= 1.5" | bc -l) )); then
    echo -e "${YELLOW}⭐⭐⭐ (3/5) - GOOD${NC}"
    echo "Session Trader is improving but needs more work"
elif (( $(echo "$WR >= 40" | bc -l) )) && (( $(echo "$PF >= 1.0" | bc -l) )); then
    echo -e "${YELLOW}⭐⭐ (2/5) - FAIR${NC}"
    echo "Session Trader is break-even to slightly profitable"
else
    echo -e "${RED}⭐ (1/5) - POOR${NC}"
    echo "Session Trader still needs optimization"
fi
echo ""

# Next steps
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🚀 NEXT STEPS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

if [ "$TRADES" -gt 100 ]; then
    echo "1. ADX filter may not be working - check if calculateADX is being called"
    echo "2. Cooldown may not be working - check if lastSessionTraderIndex is being set"
    echo "3. Try lowering ADX threshold from 25 to 20"
    echo "4. Try increasing cooldown from 30 to 40 candles"
fi

if (( $(echo "$WR < 50" | bc -l) )); then
    echo "1. Win rate is low - need stricter entry conditions"
    echo "2. Add pullback entry requirement"
    echo "3. Increase confluence requirement from 6 to 8"
    echo "4. Add volume profile analysis"
fi

if (( $(echo "$PF < 2.0" | bc -l) )); then
    echo "1. Profit factor is low - need better risk/reward"
    echo "2. Tighten stop loss from 1.5 ATR to 1.0 ATR"
    echo "3. Increase take profit targets"
    echo "4. Add trailing stops"
fi

echo ""
echo -e "${GREEN}✅ Test complete!${NC}"
echo ""
