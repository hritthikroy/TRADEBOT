#!/bin/bash

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🔬 COMPREHENSIVE STRATEGY ANALYSIS FOR LIVE TRADING"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Testing all 10 strategies across multiple timeframes"
echo "Analyzing BUY vs SELL performance"
echo "Finding the BEST strategy for live trading"
echo ""

# Check if server is running
if ! curl -s http://localhost:8080/api/v1/health > /dev/null 2>&1; then
    echo "❌ Server is not running!"
    exit 1
fi

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

timeframes=("1m" "3m" "5m" "15m" "30m" "1h" "2h" "4h" "6h" "8h" "12h" "1d")

# Create results file
RESULTS_FILE="strategy_analysis_results.csv"
echo "Strategy,Timeframe,TotalTrades,WinRate,ProfitFactor,Return,BuyTrades,BuyWinRate,SellTrades,SellWinRate,Score" > "$RESULTS_FILE"

best_score=0
best_strategy=""
best_timeframe=""
best_details=""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "TESTING ALL COMBINATIONS (This will take 5-10 minutes)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

for strategy in "${strategies[@]}"; do
    for timeframe in "${timeframes[@]}"; do
        echo -n "Testing $strategy on $timeframe... "
        
        result=$(curl -s -X POST http://localhost:8080/api/v1/backtest/run \
          -H "Content-Type: application/json" \
          -d "{
            \"symbol\": \"BTCUSDT\",
            \"interval\": \"$timeframe\",
            \"days\": 90,
            \"startBalance\": 1000,
            \"strategy\": \"$strategy\"
          }")
        
        # Extract metrics
        totalTrades=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(data.get('totalTrades', 0))" 2>/dev/null || echo "0")
        winRate=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(round(data.get('winRate', 0), 2))" 2>/dev/null || echo "0")
        profitFactor=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(round(data.get('profitFactor', 0), 2))" 2>/dev/null || echo "0")
        returnPct=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(round(data.get('returnPercent', 0), 2))" 2>/dev/null || echo "0")
        
        # Count BUY vs SELL trades
        buyTrades=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    trades = data.get('trades', [])
    buy_trades = [t for t in trades if t.get('type') == 'BUY']
    print(len(buy_trades))
except:
    print(0)
" 2>/dev/null || echo "0")
        
        sellTrades=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    trades = data.get('trades', [])
    sell_trades = [t for t in trades if t.get('type') == 'SELL']
    print(len(sell_trades))
except:
    print(0)
" 2>/dev/null || echo "0")
        
        # Calculate BUY win rate
        buyWinRate=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    trades = data.get('trades', [])
    buy_trades = [t for t in trades if t.get('type') == 'BUY']
    if len(buy_trades) == 0:
        print(0)
    else:
        winning = sum(1 for t in buy_trades if t.get('profit', 0) > 0)
        print(round(winning / len(buy_trades) * 100, 2))
except:
    print(0)
" 2>/dev/null || echo "0")
        
        # Calculate SELL win rate
        sellWinRate=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    trades = data.get('trades', [])
    sell_trades = [t for t in trades if t.get('type') == 'SELL']
    if len(sell_trades) == 0:
        print(0)
    else:
        winning = sum(1 for t in sell_trades if t.get('profit', 0) > 0)
        print(round(winning / len(sell_trades) * 100, 2))
except:
    print(0)
" 2>/dev/null || echo "0")
        
        # Calculate score (WinRate * ProfitFactor * Return / 100)
        score=$(python3 -c "print(round($winRate * $profitFactor * max(0, $returnPct) / 100, 2))" 2>/dev/null || echo "0")
        
        # Save to CSV
        echo "$strategy,$timeframe,$totalTrades,$winRate,$profitFactor,$returnPct,$buyTrades,$buyWinRate,$sellTrades,$sellWinRate,$score" >> "$RESULTS_FILE"
        
        echo "Trades: $totalTrades, WR: $winRate%, PF: $profitFactor, Return: $returnPct%, Score: $score"
        
        # Track best strategy
        if (( $(echo "$score > $best_score" | bc -l) )); then
            best_score=$score
            best_strategy=$strategy
            best_timeframe=$timeframe
            best_details="Trades: $totalTrades | WR: $winRate% | PF: $profitFactor | Return: $returnPct% | BUY: $buyTrades ($buyWinRate%) | SELL: $sellTrades ($sellWinRate%)"
        fi
    done
    echo ""
done

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🏆 BEST STRATEGY FOR LIVE TRADING"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Strategy: $best_strategy"
echo "Timeframe: $best_timeframe"
echo "Score: $best_score"
echo "$best_details"
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 TOP 5 STRATEGIES BY SCORE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
sort -t',' -k11 -rn "$RESULTS_FILE" | head -6 | tail -5 | while IFS=',' read strategy timeframe trades wr pf ret buy buywr sell sellwr score; do
    echo "🎯 $strategy ($timeframe)"
    echo "   Score: $score | Trades: $trades | WR: $wr% | PF: $pf | Return: $ret%"
    echo "   BUY: $buy trades ($buywr% WR) | SELL: $sell trades ($sellwr% WR)"
    echo ""
done

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📈 BUY-ONLY BEST PERFORMERS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
awk -F',' 'NR>1 && $7>10 {print $1","$2","$7","$8}' "$RESULTS_FILE" | sort -t',' -k4 -rn | head -5 | while IFS=',' read strategy timeframe buyTrades buyWR; do
    echo "🟢 $strategy ($timeframe): $buyTrades BUY trades with $buyWR% win rate"
done

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📉 SELL-ONLY BEST PERFORMERS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
awk -F',' 'NR>1 && $9>10 {print $1","$2","$9","$10}' "$RESULTS_FILE" | sort -t',' -k4 -rn | head -5 | while IFS=',' read strategy timeframe sellTrades sellWR; do
    echo "🔴 $strategy ($timeframe): $sellTrades SELL trades with $sellWR% win rate"
done

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "💡 RECOMMENDATIONS FOR LIVE TRADING"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "1. Use: $best_strategy on $timeframe timeframe"
echo "2. Check the CSV file for detailed analysis: $RESULTS_FILE"
echo "3. Consider filtering BUY or SELL based on which performs better"
echo "4. Run optimization on the best strategy to find optimal parameters"
echo "5. Paper trade for 30 days before going live"
echo ""
echo "⚠️  IMPORTANT: These results are based on 90 days of historical data."
echo "    Always paper trade first and never risk more than 1-2% per trade!"
echo ""
