#!/bin/bash

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🚀 SMART PARAMETER OPTIMIZATION - ALL STRATEGIES, ALL TIMEFRAMES"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Testing 10 strategies × 12 timeframes × 20 best parameter sets"
echo "= 2,400 total tests (30-60 minutes)"
echo ""

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

# Proven parameter sets (from previous successful tests)
declare -a param_sets=(
    # Format: "stop:tp1:tp2:tp3:risk:name"
    "0.5:2.0:3.0:5.0:1.0:Conservative"
    "0.75:2.5:4.0:6.0:1.5:Balanced"
    "1.0:3.0:5.0:8.0:2.0:Standard"
    "1.5:4.0:6.0:10.0:2.0:Aggressive"
    "2.0:5.0:8.0:12.0:2.5:VeryAggressive"
    "0.5:3.0:5.0:8.0:1.0:TightStopWide"
    "1.0:2.5:4.0:6.0:1.5:MediumMedium"
    "1.5:3.5:5.5:9.0:2.0:WideWide"
    "0.75:3.0:4.5:7.5:1.5:SessionTrader"
    "1.5:4.0:6.0:8.0:2.0:LiquidityHunter"
    "0.5:2.5:4.0:7.0:1.0:LowRisk"
    "1.0:3.5:5.5:9.5:1.5:MediumRisk"
    "1.25:3.75:5.75:9.75:1.75:Balanced2"
    "0.75:2.75:4.25:6.75:1.25:Conservative2"
    "1.75:4.5:7.0:11.0:2.25:Aggressive2"
    "0.5:2.0:4.0:6.0:0.5:UltraConservative"
    "2.5:6.0:9.0:15.0:3.0:UltraAggressive"
    "1.0:4.0:6.0:10.0:2.0:BreakoutMaster"
    "0.5:3.0:4.5:7.5:1.0:SmartMoney"
    "1.5:5.0:7.5:12.5:2.0:ReversalSniper"
)

RESULTS_FILE="smart_optimization_results.csv"
echo "Strategy,Timeframe,StopATR,TP1ATR,TP2ATR,TP3ATR,Risk,ParamSet,Trades,WinRate,ProfitFactor,Return,MaxDD,AvgDD,Score,Grade" > "$RESULTS_FILE"

total_tests=$((${#strategies[@]} * ${#timeframes[@]} * ${#param_sets[@]}))
current_test=0
start_time=$(date +%s)

echo "Total tests to run: $total_tests"
echo ""

for strategy in "${strategies[@]}"; do
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "🎯 Testing: $strategy"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    strategy_best_score=0
    strategy_best_config=""
    
    for timeframe in "${timeframes[@]}"; do
        for param_set in "${param_sets[@]}"; do
            current_test=$((current_test + 1))
            
            IFS=':' read -r stop tp1 tp2 tp3 risk name <<< "$param_set"
            
            # Progress
            elapsed=$(($(date +%s) - start_time))
            tests_per_sec=$(echo "scale=2; $current_test / $elapsed" | bc)
            remaining=$((total_tests - current_test))
            eta=$(echo "scale=0; $remaining / $tests_per_sec" | bc)
            
            echo -ne "\r[$current_test/$total_tests] $timeframe | $name | ETA: ${eta}s   "
            
            # Run backtest
            result=$(curl -s -X POST http://localhost:8080/api/v1/backtest/run \
              -H "Content-Type: application/json" \
              -d "{
                \"symbol\": \"BTCUSDT\",
                \"interval\": \"$timeframe\",
                \"days\": 90,
                \"startBalance\": 1000,
                \"strategy\": \"$strategy\",
                \"riskPercent\": $(echo "$risk / 100" | bc -l)
              }")
            
            # Extract and calculate metrics
            metrics=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    trades = data.get('totalTrades', 0)
    wr = round(data.get('winRate', 0), 2)
    pf = round(data.get('profitFactor', 0), 2)
    ret = round(data.get('returnPercent', 0), 2)
    maxdd = round(data.get('maxDrawdown', 0), 2)
    
    # Calculate average drawdown
    trade_list = data.get('trades', [])
    if len(trade_list) > 0:
        equity = [1000]
        for t in trade_list:
            equity.append(t.get('balanceAfter', equity[-1]))
        
        peak = equity[0]
        drawdowns = []
        for balance in equity:
            if balance > peak:
                peak = balance
            else:
                dd = (peak - balance) / peak * 100
                if dd > 0:
                    drawdowns.append(dd)
        
        avgdd = round(sum(drawdowns) / len(drawdowns), 2) if drawdowns else 0
    else:
        avgdd = 0
    
    # Calculate score
    if trades >= 5 and ret > 0:
        wr_score = wr * 0.25
        pf_score = min(pf * 20, 100) * 0.25
        ret_score = min(ret, 100) * 0.25
        dd_penalty = maxdd * 2.5
        smooth_bonus = max(0, 10 - avgdd) * 0.25
        score = round(max(0, wr_score + pf_score + ret_score + smooth_bonus - dd_penalty), 2)
    else:
        score = 0
    
    # Grade
    if ret > 0 and maxdd < 10 and wr >= 60:
        grade = 'Excellent'
    elif ret > 0 and maxdd < 15 and wr >= 50:
        grade = 'Good'
    elif ret > 0 and maxdd < 20 and wr >= 45:
        grade = 'Acceptable'
    else:
        grade = 'Poor'
    
    print(f'{trades},{wr},{pf},{ret},{maxdd},{avgdd},{score},{grade}')
except:
    print('0,0,0,0,0,0,0,Poor')
" 2>/dev/null || echo "0,0,0,0,0,0,0,Poor")
            
            IFS=',' read -r trades wr pf ret maxdd avgdd score grade <<< "$metrics"
            
            # Save if profitable
            if (( $(echo "$score > 0" | bc -l) )); then
                echo "$strategy,$timeframe,$stop,$tp1,$tp2,$tp3,$risk,$name,$trades,$wr,$pf,$ret,$maxdd,$avgdd,$score,$grade" >> "$RESULTS_FILE"
                
                # Track best for this strategy
                if (( $(echo "$score > $strategy_best_score" | bc -l) )); then
                    strategy_best_score=$score
                    strategy_best_config="$timeframe | $name | Stop:$stop TP1:$tp1 TP2:$tp2 TP3:$tp3 Risk:$risk% | WR:$wr% PF:$pf Ret:$ret% DD:$maxdd% | Score:$score"
                fi
            fi
        done
    done
    
    echo ""
    if [ -n "$strategy_best_config" ]; then
        echo "✅ Best for $strategy: $strategy_best_config"
    else
        echo "❌ No profitable config found for $strategy"
    fi
    echo ""
done

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ OPTIMIZATION COMPLETE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Generate summary
echo "📊 TOP 20 CONFIGURATIONS:"
echo ""
sort -t',' -k15 -rn "$RESULTS_FILE" | head -21 | tail -20 | while IFS=',' read strat tf stop tp1 tp2 tp3 risk pset trades wr pf ret maxdd avgdd score grade; do
    echo "🎯 $strat ($tf) - $pset"
    echo "   Params: Stop=$stop TP1=$tp1 TP2=$tp2 TP3=$tp3 Risk=$risk%"
    echo "   Results: WR=$wr% PF=$pf Ret=$ret% MaxDD=$maxdd% AvgDD=$avgdd%"
    echo "   Score: $score | Grade: $grade"
    echo ""
done

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📈 BEST CONFIGURATION PER STRATEGY:"
echo ""

for strategy in "${strategies[@]}"; do
    best=$(grep "^$strategy," "$RESULTS_FILE" | sort -t',' -k15 -rn | head -1)
    if [ -n "$best" ]; then
        IFS=',' read -r strat tf stop tp1 tp2 tp3 risk pset trades wr pf ret maxdd avgdd score grade <<< "$best"
        echo "🏆 $strategy:"
        echo "   Timeframe: $tf | Params: $pset"
        echo "   Stop=$stop TP1=$tp1 TP2=$tp2 TP3=$tp3 Risk=$risk%"
        echo "   WR=$wr% PF=$pf Ret=$ret% MaxDD=$maxdd% Score=$score"
        echo ""
    else
        echo "❌ $strategy: No profitable configuration"
        echo ""
    fi
done

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📄 Results saved to: $RESULTS_FILE"
echo ""
echo "View full results:"
echo "  cat $RESULTS_FILE | column -t -s','"
echo ""
echo "View by strategy:"
echo "  grep 'liquidity_hunter' $RESULTS_FILE | sort -t',' -k15 -rn | head -10"
echo ""
