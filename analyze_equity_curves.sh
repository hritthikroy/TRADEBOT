#!/bin/bash

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📈 EQUITY CURVE & DRAWDOWN ANALYSIS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Finding strategies with:"
echo "  ✅ Smooth equity curves (consistent growth)"
echo "  ✅ Low drawdowns (<10% ideal, <20% acceptable)"
echo "  ✅ High win rates (>50%)"
echo "  ✅ Positive returns"
echo ""

# Check server
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

timeframes=("1d" "8h" "4h" "2h" "1h")

# Results file
RESULTS_FILE="equity_curve_analysis.csv"
echo "Strategy,Timeframe,Trades,WinRate,ProfitFactor,Return,MaxDrawdown,AvgDrawdown,DrawdownDuration,EquityCurveScore,Grade" > "$RESULTS_FILE"

declare -a excellent_configs=()
declare -a good_configs=()
declare -a acceptable_configs=()

total_tests=$((${#strategies[@]} * ${#timeframes[@]}))
current_test=0

echo "Testing $total_tests combinations..."
echo ""

for strategy in "${strategies[@]}"; do
    for timeframe in "${timeframes[@]}"; do
        current_test=$((current_test + 1))
        echo -ne "\r[$current_test/$total_tests] Testing $strategy on $timeframe...                    "
        
        result=$(curl -s -X POST http://localhost:8080/api/v1/backtest/run \
          -H "Content-Type: application/json" \
          -d "{
            \"symbol\": \"BTCUSDT\",
            \"interval\": \"$timeframe\",
            \"days\": 90,
            \"startBalance\": 1000,
            \"strategy\": \"$strategy\"
          }")
        
        # Extract comprehensive metrics
        totalTrades=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(data.get('totalTrades', 0))" 2>/dev/null || echo "0")
        winRate=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(round(data.get('winRate', 0), 2))" 2>/dev/null || echo "0")
        profitFactor=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(round(data.get('profitFactor', 0), 2))" 2>/dev/null || echo "0")
        returnPct=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(round(data.get('returnPercent', 0), 2))" 2>/dev/null || echo "0")
        maxDD=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(round(data.get('maxDrawdown', 0), 2))" 2>/dev/null || echo "0")
        
        # Calculate equity curve metrics
        metrics=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    trades = data.get('trades', [])
    
    if len(trades) == 0:
        print('0,0,0')
        sys.exit(0)
    
    # Calculate equity curve
    equity = [1000]  # Start balance
    for trade in trades:
        equity.append(trade.get('balanceAfter', equity[-1]))
    
    # Calculate drawdowns
    peak = equity[0]
    drawdowns = []
    dd_durations = []
    current_dd_duration = 0
    
    for i, balance in enumerate(equity):
        if balance > peak:
            peak = balance
            if current_dd_duration > 0:
                dd_durations.append(current_dd_duration)
                current_dd_duration = 0
        else:
            dd = (peak - balance) / peak * 100
            drawdowns.append(dd)
            current_dd_duration += 1
    
    # Metrics
    avg_dd = sum(drawdowns) / len(drawdowns) if drawdowns else 0
    max_dd_duration = max(dd_durations) if dd_durations else 0
    
    # Equity curve smoothness (lower std dev = smoother)
    if len(equity) > 1:
        returns = [(equity[i] - equity[i-1]) / equity[i-1] * 100 for i in range(1, len(equity))]
        import statistics
        smoothness = statistics.stdev(returns) if len(returns) > 1 else 0
    else:
        smoothness = 0
    
    print(f'{round(avg_dd, 2)},{max_dd_duration},{round(smoothness, 2)}')
except Exception as e:
    print('0,0,0')
" 2>/dev/null || echo "0,0,0")
        
        IFS=',' read -r avgDD ddDuration smoothness <<< "$metrics"
        
        # Calculate Equity Curve Score
        # Higher is better: Positive return, low drawdown, high win rate, smooth curve
        equityScore=$(python3 -c "
ret = $returnPct
wr = $winRate
pf = $profitFactor
max_dd = $maxDD
avg_dd = $avgDD
smooth = $smoothness
trades = $totalTrades

if trades < 5 or ret <= 0:
    print(0)
else:
    # Score components (0-100 scale)
    return_score = min(ret, 100)  # Cap at 100%
    wr_score = wr  # Already 0-100
    pf_score = min(pf * 20, 100)  # 5.0 PF = 100 points
    dd_penalty = max_dd * 2  # Each 1% DD = -2 points
    smooth_bonus = max(0, 50 - smooth)  # Smoother = higher bonus
    
    # Weighted score
    score = (return_score * 0.3) + (wr_score * 0.25) + (pf_score * 0.25) + smooth_bonus * 0.1 - dd_penalty * 0.1
    print(round(max(0, score), 2))
" 2>/dev/null || echo "0")
        
        # Grade the configuration
        grade="❌ Poor"
        if (( $(echo "$returnPct > 0 && $maxDD < 10 && $winRate >= 60" | bc -l) )); then
            grade="🏆 Excellent"
            excellent_configs+=("$strategy|$timeframe|$totalTrades|$winRate|$profitFactor|$returnPct|$maxDD|$avgDD|$equityScore")
        elif (( $(echo "$returnPct > 0 && $maxDD < 15 && $winRate >= 50" | bc -l) )); then
            grade="✅ Good"
            good_configs+=("$strategy|$timeframe|$totalTrades|$winRate|$profitFactor|$returnPct|$maxDD|$avgDD|$equityScore")
        elif (( $(echo "$returnPct > 0 && $maxDD < 20 && $winRate >= 45" | bc -l) )); then
            grade="⚠️  Acceptable"
            acceptable_configs+=("$strategy|$timeframe|$totalTrades|$winRate|$profitFactor|$returnPct|$maxDD|$avgDD|$equityScore")
        fi
        
        # Save to CSV
        echo "$strategy,$timeframe,$totalTrades,$winRate,$profitFactor,$returnPct,$maxDD,$avgDD,$ddDuration,$equityScore,$grade" >> "$RESULTS_FILE"
    done
done

echo ""
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🏆 EXCELLENT CONFIGURATIONS (Low DD <10%, High WR >60%)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

if [ ${#excellent_configs[@]} -eq 0 ]; then
    echo "No configurations met excellent criteria"
else
    for config in "${excellent_configs[@]}"; do
        IFS='|' read -r strat tf trades wr pf ret maxdd avgdd score <<< "$config"
        echo "🎯 $strat ($tf)"
        echo "   Trades: $trades | WR: $wr% | PF: $pf | Return: $ret%"
        echo "   Max DD: $maxdd% | Avg DD: $avgdd% | Score: $score"
        echo ""
    done
fi

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ GOOD CONFIGURATIONS (Low DD <15%, High WR >50%)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

if [ ${#good_configs[@]} -eq 0 ]; then
    echo "No configurations met good criteria"
else
    for config in "${good_configs[@]}"; do
        IFS='|' read -r strat tf trades wr pf ret maxdd avgdd score <<< "$config"
        echo "🎯 $strat ($tf)"
        echo "   Trades: $trades | WR: $wr% | PF: $pf | Return: $ret%"
        echo "   Max DD: $maxdd% | Avg DD: $avgdd% | Score: $score"
        echo ""
    done
fi

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "⚠️  ACCEPTABLE CONFIGURATIONS (DD <20%, WR >45%)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

if [ ${#acceptable_configs[@]} -eq 0 ]; then
    echo "No configurations met acceptable criteria"
else
    # Show top 5 only
    count=0
    for config in "${acceptable_configs[@]}"; do
        if [ $count -ge 5 ]; then break; fi
        IFS='|' read -r strat tf trades wr pf ret maxdd avgdd score <<< "$config"
        echo "🎯 $strat ($tf)"
        echo "   Trades: $trades | WR: $wr% | PF: $pf | Return: $ret%"
        echo "   Max DD: $maxdd% | Avg DD: $avgdd% | Score: $score"
        echo ""
        count=$((count + 1))
    done
    
    if [ ${#acceptable_configs[@]} -gt 5 ]; then
        echo "... and $((${#acceptable_configs[@]} - 5)) more (see CSV file)"
        echo ""
    fi
fi

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 TOP 5 BY EQUITY CURVE SCORE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

sort -t',' -k10 -rn "$RESULTS_FILE" | head -6 | tail -5 | while IFS=',' read strat tf trades wr pf ret maxdd avgdd dddur score grade; do
    echo "🎯 $strat ($tf) - Score: $score"
    echo "   Trades: $trades | WR: $wr% | PF: $pf | Return: $ret%"
    echo "   Max DD: $maxdd% | Avg DD: $avgdd% | Grade: $grade"
    echo ""
done

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📉 LOWEST DRAWDOWN CONFIGURATIONS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

awk -F',' 'NR>1 && $6>0 {print $0}' "$RESULTS_FILE" | sort -t',' -k7 -n | head -5 | while IFS=',' read strat tf trades wr pf ret maxdd avgdd dddur score grade; do
    echo "🎯 $strat ($tf) - Max DD: $maxdd%"
    echo "   Trades: $trades | WR: $wr% | PF: $pf | Return: $ret%"
    echo "   Avg DD: $avgdd% | Score: $score"
    echo ""
done

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "💡 RECOMMENDATIONS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "For SMOOTH EQUITY CURVE with LOW DRAWDOWN:"
echo ""

# Find best overall
best_config=$(awk -F',' 'NR>1 && $6>0 && $7<15 && $4>=50 {print $0}' "$RESULTS_FILE" | sort -t',' -k10 -rn | head -1)

if [ -n "$best_config" ]; then
    IFS=',' read -r strat tf trades wr pf ret maxdd avgdd dddur score grade <<< "$best_config"
    echo "🏆 BEST OVERALL: $strat on $tf timeframe"
    echo ""
    echo "   Performance:"
    echo "   - Win Rate: $wr%"
    echo "   - Profit Factor: $pf"
    echo "   - Return: $ret%"
    echo "   - Max Drawdown: $maxdd%"
    echo "   - Avg Drawdown: $avgdd%"
    echo "   - Total Trades: $trades"
    echo "   - Equity Score: $score"
    echo ""
    echo "   Why this is good:"
    echo "   ✅ Positive returns"
    echo "   ✅ Low drawdown (<15%)"
    echo "   ✅ High win rate (>50%)"
    echo "   ✅ Smooth equity curve"
else
    echo "⚠️  No configuration met all criteria (Return >0, DD <15%, WR >50%)"
    echo ""
    echo "   Best alternative:"
    best_alt=$(awk -F',' 'NR>1 && $6>0 {print $0}' "$RESULTS_FILE" | sort -t',' -k10 -rn | head -1)
    if [ -n "$best_alt" ]; then
        IFS=',' read -r strat tf trades wr pf ret maxdd avgdd dddur score grade <<< "$best_alt"
        echo "   $strat on $tf: Return $ret%, DD $maxdd%, WR $wr%"
    fi
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📄 DETAILED RESULTS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Full analysis saved to: $RESULTS_FILE"
echo ""
echo "To view:"
echo "  cat $RESULTS_FILE | column -t -s','"
echo ""
echo "To find specific criteria:"
echo "  # Low drawdown (<10%)"
echo "  awk -F',' 'NR>1 && \$7<10 && \$6>0' $RESULTS_FILE | column -t -s','"
echo ""
echo "  # High win rate (>60%)"
echo "  awk -F',' 'NR>1 && \$4>60 && \$6>0' $RESULTS_FILE | column -t -s','"
echo ""
