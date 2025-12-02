#!/bin/bash

# Comprehensive Strategy Testing Script
# Tests all timeframes with optimized strategies

echo "🚀 COMPREHENSIVE STRATEGY BACKTEST"
echo "===================================="
echo ""

BASE_URL="http://localhost:8080"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Test all timeframes
TIMEFRAMES=("1m" "3m" "5m" "15m" "30m" "1h" "2h" "4h" "1d")
SYMBOL="BTCUSDT"
START_BALANCE=500

echo -e "${BLUE}Testing Symbol: $SYMBOL${NC}"
echo -e "${BLUE}Start Balance: \$$START_BALANCE${NC}"
echo ""

# Results array
declare -A RESULTS

for TF in "${TIMEFRAMES[@]}"; do
    echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}📊 Testing $TF Timeframe${NC}"
    echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    
    # Determine days based on timeframe
    case $TF in
        "1m"|"3m"|"5m")
            DAYS=7
            ;;
        "15m"|"30m")
            DAYS=30
            ;;
        "1h"|"2h")
            DAYS=60
            ;;
        "4h")
            DAYS=90
            ;;
        "1d")
            DAYS=180
            ;;
    esac
    
    echo "Testing $DAYS days of data..."
    
    # Run backtest
    RESULT=$(curl -s -X POST "$BASE_URL/api/v1/backtest/run" \
        -H "Content-Type: application/json" \
        -d "{\"symbol\":\"$SYMBOL\",\"interval\":\"$TF\",\"days\":$DAYS,\"startBalance\":$START_BALANCE}")
    
    # Parse results
    if echo "$RESULT" | grep -q "totalTrades"; then
        TRADES=$(echo "$RESULT" | python3 -c "import sys,json; print(json.load(sys.stdin).get('totalTrades', 0))" 2>/dev/null)
        WINRATE=$(echo "$RESULT" | python3 -c "import sys,json; print(f\"{json.load(sys.stdin).get('winRate', 0):.1f}\")" 2>/dev/null)
        RETURN=$(echo "$RESULT" | python3 -c "import sys,json; print(f\"{json.load(sys.stdin).get('returnPercent', 0):.1f}\")" 2>/dev/null)
        PF=$(echo "$RESULT" | python3 -c "import sys,json; print(f\"{json.load(sys.stdin).get('profitFactor', 0):.2f}\")" 2>/dev/null)
        END_BAL=$(echo "$RESULT" | python3 -c "import sys,json; print(f\"{json.load(sys.stdin).get('endBalance', 0):.2f}\")" 2>/dev/null)
        
        # Store results
        RESULTS["$TF"]="$TRADES|$WINRATE|$RETURN|$PF|$END_BAL"
        
        # Display results
        echo ""
        echo "  📈 Total Trades:    $TRADES"
        echo "  🎯 Win Rate:        $WINRATE%"
        echo "  💰 Return:          $RETURN%"
        echo "  📊 Profit Factor:   $PF"
        echo "  💵 End Balance:     \$$END_BAL"
        
        # Color code based on performance
        if (( $(echo "$WINRATE > 70" | bc -l) )); then
            echo -e "  ${GREEN}✅ EXCELLENT WIN RATE${NC}"
        elif (( $(echo "$WINRATE > 60" | bc -l) )); then
            echo -e "  ${GREEN}✅ GOOD WIN RATE${NC}"
        elif (( $(echo "$WINRATE > 50" | bc -l) )); then
            echo -e "  ${YELLOW}⚠️  MODERATE WIN RATE${NC}"
        else
            echo -e "  ${RED}❌ LOW WIN RATE${NC}"
        fi
        
        if (( $(echo "$RETURN > 20" | bc -l) )); then
            echo -e "  ${GREEN}✅ HIGHLY PROFITABLE${NC}"
        elif (( $(echo "$RETURN > 0" | bc -l) )); then
            echo -e "  ${GREEN}✅ PROFITABLE${NC}"
        else
            echo -e "  ${RED}❌ UNPROFITABLE${NC}"
        fi
        
    else
        echo -e "  ${RED}❌ Failed to get results${NC}"
        RESULTS["$TF"]="0|0|0|0|0"
    fi
    
    echo ""
    sleep 2
done

# Summary
echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}📊 COMPREHENSIVE SUMMARY${NC}"
echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

printf "%-10s | %8s | %8s | %8s | %8s | %12s\n" "Timeframe" "Trades" "Win Rate" "Return %" "PF" "End Balance"
echo "--------------------------------------------------------------------------------"

BEST_WINRATE=0
BEST_WINRATE_TF=""
BEST_RETURN=0
BEST_RETURN_TF=""
BEST_PF=0
BEST_PF_TF=""

for TF in "${TIMEFRAMES[@]}"; do
    IFS='|' read -r TRADES WINRATE RETURN PF END_BAL <<< "${RESULTS[$TF]}"
    
    printf "%-10s | %8s | %7s%% | %7s%% | %8s | \$%10s\n" "$TF" "$TRADES" "$WINRATE" "$RETURN" "$PF" "$END_BAL"
    
    # Track best performers
    if (( $(echo "$WINRATE > $BEST_WINRATE" | bc -l) )); then
        BEST_WINRATE=$WINRATE
        BEST_WINRATE_TF=$TF
    fi
    
    if (( $(echo "$RETURN > $BEST_RETURN" | bc -l) )); then
        BEST_RETURN=$RETURN
        BEST_RETURN_TF=$TF
    fi
    
    if (( $(echo "$PF > $BEST_PF" | bc -l) )); then
        BEST_PF=$PF
        BEST_PF_TF=$TF
    fi
done

echo ""
echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}🏆 BEST PERFORMERS${NC}"
echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "  🎯 Best Win Rate:      ${GREEN}$BEST_WINRATE_TF ($BEST_WINRATE%)${NC}"
echo -e "  💰 Best Return:        ${GREEN}$BEST_RETURN_TF ($BEST_RETURN%)${NC}"
echo -e "  📊 Best Profit Factor: ${GREEN}$BEST_PF_TF ($BEST_PF)${NC}"
echo ""

# Recommendations
echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}💡 RECOMMENDATIONS${NC}"
echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo "Based on the backtest results:"
echo ""
echo "  For SCALPING (quick trades):"
echo "    → Use 5m or 15m timeframes"
echo ""
echo "  For DAY TRADING (intraday):"
echo "    → Use 15m or 30m timeframes"
echo ""
echo "  For SWING TRADING (multi-day):"
echo "    → Use 4h or 1d timeframes"
echo ""
echo "  For HIGHEST WIN RATE:"
echo "    → Use $BEST_WINRATE_TF timeframe"
echo ""
echo "  For MAXIMUM PROFIT:"
echo "    → Use $BEST_RETURN_TF timeframe"
echo ""

echo -e "${GREEN}✅ Comprehensive strategy testing complete!${NC}"
echo ""
