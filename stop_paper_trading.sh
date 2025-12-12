#!/bin/bash

# ============================================
# STOP PAPER TRADING
# ============================================

echo "🛑 STOPPING PAPER TRADING"
echo "========================"
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# Check backend
if ! curl -s http://localhost:8080/api/v1/health > /dev/null 2>&1; then
    echo -e "${RED}❌ Backend not running${NC}"
    exit 1
fi

# Get final stats before stopping
echo "📊 Final Statistics:"
echo ""

stats=$(curl -s "http://localhost:8080/api/v1/paper-trading/stats")
total_trades=$(echo "$stats" | jq -r '.totalTrades')
win_rate=$(echo "$stats" | jq -r '.winRate')
net_profit=$(echo "$stats" | jq -r '.netProfit')
return_pct=$(echo "$stats" | jq -r '.returnPercent')

printf "  Total Trades:  %d\n" "$total_trades"
printf "  Win Rate:      %.2f%%\n" "$win_rate"
printf "  Net Profit:    \$%.2f\n" "$net_profit"
printf "  Return:        %.2f%%\n" "$return_pct"
echo ""

# Stop paper trading
response=$(curl -s -X POST "http://localhost:8080/api/v1/paper-trading/stop-auto")

success=$(echo "$response" | jq -r '.success')

if [ "$success" = "true" ]; then
    echo -e "${GREEN}✅ Paper trading stopped successfully${NC}"
    echo ""
    echo "Results have been saved."
    echo "You can restart anytime with: ./start_paper_trading_session_trader.sh"
else
    echo -e "${RED}❌ Failed to stop paper trading${NC}"
    echo "$response" | jq '.'
fi

echo ""
