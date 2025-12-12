#!/bin/bash

# ============================================
# START PAPER TRADING - SESSION TRADER
# ============================================
# Starts paper trading with Session Trader strategy

echo "🚀 STARTING PAPER TRADING - SESSION TRADER"
echo "=========================================="
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

echo "📊 Checking current performance..."
echo ""

# Run quick check
result=$(curl -s -X POST "http://localhost:8080/api/v1/backtest/run" \
    -H "Content-Type: application/json" \
    -d '{"symbol":"BTCUSDT","interval":"15m","days":7,"strategy":"session_trader","startBalance":1000}')

win_rate=$(echo "$result" | jq -r '.winRate')
pf=$(echo "$result" | jq -r '.profitFactor')
return=$(echo "$result" | jq -r '.returnPercent')

printf "  7-Day Win Rate:     %.2f%%\n" "$win_rate"
printf "  7-Day Profit Factor: %.2f\n" "$pf"
printf "  7-Day Return:       %.2f%%\n" "$return"
echo ""

# Check if profitable
is_profitable=$(echo "$return > 0" | bc -l)
good_pf=$(echo "$pf >= 1.0" | bc -l)

if [ "$is_profitable" != "1" ] || [ "$good_pf" != "1" ]; then
    echo -e "${YELLOW}⚠️  WARNING: Recent performance is not profitable${NC}"
    echo ""
    echo "Do you still want to start paper trading? (y/n)"
    read -r response
    if [ "$response" != "y" ]; then
        echo "Paper trading cancelled."
        exit 0
    fi
fi

echo -e "${GREEN}✅ Starting paper trading...${NC}"
echo ""

# Start auto paper trading
response=$(curl -s -X POST "http://localhost:8080/api/v1/paper-trading/start-auto" \
    -H "Content-Type: application/json" \
    -d '{
        "symbol": "BTCUSDT",
        "strategy": "session_trader",
        "interval": "15m",
        "startBalance": 1000,
        "riskPercent": 1.0
    }')

echo "$response" | jq '.'
echo ""

# Check status
status=$(curl -s "http://localhost:8080/api/v1/paper-trading/stats")
is_running=$(echo "$status" | jq -r '.isRunning')

if [ "$is_running" = "true" ]; then
    echo -e "${GREEN}✅ Paper trading started successfully!${NC}"
    echo ""
    echo "📊 Dashboard: http://localhost:8080/paper-trading"
    echo ""
    echo "Monitor your trades:"
    echo "  • View dashboard in browser"
    echo "  • Check stats: ./check_paper_trading_stats.sh"
    echo "  • Stop trading: ./stop_paper_trading.sh"
    echo ""
    echo -e "${YELLOW}⚠️  IMPORTANT REMINDERS:${NC}"
    echo "  • This is PAPER TRADING (not real money)"
    echo "  • Monitor performance daily"
    echo "  • Stop if performance degrades"
    echo "  • Use small sizes when going live"
    echo ""
else
    echo -e "${RED}❌ Failed to start paper trading${NC}"
    echo "Check backend logs for errors"
fi
