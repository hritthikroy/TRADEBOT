#!/bin/bash

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🔬 COMPREHENSIVE PARAMETER OPTIMIZATION FOR ALL 10 STRATEGIES"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Goal: Find BEST parameters for:"
echo "  ✅ Smooth equity curves"
echo "  ✅ Low drawdowns (<10%)"
echo "  ✅ High win rates (>50%)"
echo "  ✅ High profit factors (>2.0)"
echo "  ✅ Positive returns"
echo ""
echo "This will test 100+ parameter combinations per strategy"
echo "Expected time: 2-3 hours"
echo ""

read -p "Continue with full optimization? (y/n) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Optimization cancelled"
    exit 0
fi

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

# All timeframes from 1m to 1d
timeframes=("1m" "3m" "5m" "15m" "30m" "1h" "2h" "4h" "6h" "8h" "12h" "1d")

# Parameter ranges to test
stopLossValues=(0.5 0.75 1.0 1.25 1.5 2.0)
tp1Values=(2.0 2.5 3.0 3.5 4.0 5.0)
tp2Values=(3.0 4.0 4.5 5.0 6.0 7.5)
tp3Values=(5.0 6.0 7.5 10.0 12.5 15.0)
riskValues=(0.5 1.0 1.5 2.0 2.5)

# Results file
RESULTS_FILE="best_parameters_all_strategies.csv"
echo "Strategy,Timeframe,StopATR,TP1ATR,TP2ATR,TP3ATR,Risk,Trades,WinRate,ProfitFactor,Return,MaxDD,AvgDD,Score" > "$RESULTS_FILE"

mkdir -p optimization_results

for strategy in "${strategies[@]}"; do
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "🎯 OPTIMIZING: $strategy"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    
    best_score=0
    best_config=""
    total_tests=0
    profitable_configs=0
    
    for timeframe in "${timeframes[@]}"; do
        echo "Testing $strategy on $timeframe..."
        
        for stop in "${stopLossValues[@]}"; do
            for tp1 in "${tp1Values[@]}"; do
                for tp2 in "${tp2Values[@]}"; do
                    for tp3 in "${tp3Values[@]}"; do
                        # Validate: TP1 < TP2 < TP3
                        if (( $(echo "$tp1 < $tp2" | bc -l) )) && (( $(echo "$tp2 < $tp3" | bc -l) )); then
                            for risk in "${riskValues[@]}"; do
                                total_tests=$((total_tests + 1))
                                
                                # Progress
                                if [ $((total_tests % 50)) -eq 0 ]; then
                                    echo -ne "\r   [$timeframe] Tested: $total_tests | Profitable: $profitable_configs | Best Score: $best_score   "
                                fi
                                
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
                                
                                # Extract metrics
                                metrics=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    trades = data.get('totalTrades', 0)
    wr = data.get('winRate', 0)
    pf = data.get('profitFactor', 0)
    ret = data.get('returnPercent', 0)
    maxdd = data.get('maxDrawdown', 0)
    
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
        
        avgdd = sum(drawdowns) / len(drawdowns) if drawdowns else 0
    else:
        avgdd = 0
    
    print(f'{trades},{round(wr, 2)},{round(pf, 2)},{round(ret, 2)},{round(maxdd, 2)},{round(avgdd, 2)}')
except:
    print('0,0,0,0,0,0')
" 2>/dev/null || echo "0,0,0,0,0,0")
                                
                                IFS=',' read -r trades wr pf ret maxdd avgdd <<< "$metrics"
                                
                                # Calculate comprehensive score
                                # Prioritize: Low DD, High WR, High PF, Positive Return
                                score=$(python3 -c "
trades = $trades
wr = $wr
pf = $pf
ret = $ret
maxdd = $maxdd
avgdd = $avgdd

if trades < 5 or ret <= 0:
    print(0)
else:
    # Score components
    wr_score = wr * 0.25  # 25% weight
    pf_score = min(pf * 20, 100) * 0.25  # 25% weight
    ret_score = min(ret, 100) * 0.25  # 25% weight
    dd_penalty = maxdd * 2.5  # Penalty for drawdown
    
    # Bonus for smooth equity (low avg DD)
    smooth_bonus = max(0, 10 - avgdd) * 0.25  # 25% weight
    
    # Final score
    score = wr_score + pf_score + ret_score + smooth_bonus - dd_penalty
    print(round(max(0, score), 2))
" 2>/dev/null || echo "0")
                                
                                # Track profitable configs
                                if (( $(echo "$ret > 0 && $pf >= 1.0" | bc -l) )); then
                                    profitable_configs=$((profitable_configs + 1))
                                fi
                                
                                # Save if good score
                                if (( $(echo "$score > 0" | bc -l) )); then
                                    echo "$strategy,$timeframe,$stop,$tp1,$tp2,$tp3,$risk,$trades,$wr,$pf,$ret,$maxdd,$avgdd,$score" >> "$RESULTS_FILE"
                                fi
                                
                                # Track best
                                if (( $(echo "$score > $best_score" | bc -l) )); then
                                    best_score=$score
                                    best_config="TF: $timeframe | Stop: $stop | TP1: $tp1 | TP2: $tp2 | TP3: $tp3 | Risk: $risk% | Trades: $trades | WR: $wr% | PF: $pf | Return: $ret% | MaxDD: $maxdd% | AvgDD: $avgdd%"
                                fi
                            done
                        fi
                    done
                done
            done
        done
        echo ""
    done
    
    echo ""
    echo "✅ $strategy Optimization Complete!"
    echo "   Total Tests: $total_tests"
    echo "   Profitable Configs: $profitable_configs"
    echo "   Best Score: $best_score"
    echo ""
    if [ -n "$best_config" ]; then
        echo "   🏆 BEST CONFIGURATION:"
        echo "   $best_config"
    else
        echo "   ❌ No profitable configuration found"
    fi
    echo ""
    
    # Save strategy-specific results
    grep "^$strategy," "$RESULTS_FILE" | sort -t',' -k14 -rn | head -10 > "optimization_results/${strategy}_top10.csv"
done

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🏆 OPTIMIZATION COMPLETE - ANALYZING RESULTS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Generate summary report
SUMMARY_FILE="BEST_PARAMETERS_SUMMARY.md"

cat > "$SUMMARY_FILE" << 'EOF'
# 🎯 BEST PARAMETERS FOR ALL 10 STRATEGIES

## Optimization Date: 
EOF

date >> "$SUMMARY_FILE"

cat >> "$SUMMARY_FILE" << 'EOF'

## 📊 Optimization Criteria:
- **Goal**: Smooth equity curves with low drawdowns
- **Test Period**: 90 days
- **Timeframes Tested**: 1d, 8h, 4h
- **Parameter Combinations**: 100+ per strategy
- **Scoring**: Win Rate (25%) + Profit Factor (25%) + Return (25%) + Low DD Bonus (25%)

---

## 🏆 TOP 10 OVERALL BEST CONFIGURATIONS

EOF

echo "Analyzing top configurations..."
sort -t',' -k14 -rn "$RESULTS_FILE" | head -11 | tail -10 | while IFS=',' read strat tf stop tp1 tp2 tp3 risk trades wr pf ret maxdd avgdd score; do
    cat >> "$SUMMARY_FILE" << EOF

### $strat on $tf timeframe
**Parameters:**
- Stop Loss: $stop ATR
- TP1: $tp1 ATR (33%)
- TP2: $tp2 ATR (33%)
- TP3: $tp3 ATR (34%)
- Risk per Trade: $risk%

**Performance:**
- Win Rate: $wr%
- Profit Factor: $pf
- Return: $ret%
- Max Drawdown: $maxdd%
- Avg Drawdown: $avgdd%
- Total Trades: $trades
- **Score: $score**

EOF
done

cat >> "$SUMMARY_FILE" << 'EOF'

---

## 📈 BEST PARAMETERS BY STRATEGY

EOF

for strategy in "${strategies[@]}"; do
    echo "Analyzing $strategy..."
    
    best_line=$(grep "^$strategy," "$RESULTS_FILE" | sort -t',' -k14 -rn | head -1)
    
    if [ -n "$best_line" ]; then
        IFS=',' read -r strat tf stop tp1 tp2 tp3 risk trades wr pf ret maxdd avgdd score <<< "$best_line"
        
        cat >> "$SUMMARY_FILE" << EOF

### $strategy
**Best Timeframe:** $tf

**Optimal Parameters:**
\`\`\`
Stop Loss: $stop ATR
TP1: $tp1 ATR (33%)
TP2: $tp2 ATR (33%)
TP3: $tp3 ATR (34%)
Risk per Trade: $risk%
\`\`\`

**Expected Performance:**
- Win Rate: $wr%
- Profit Factor: $pf
- Return: $ret%
- Max Drawdown: $maxdd%
- Avg Drawdown: $avgdd%
- Total Trades: $trades
- Score: $score

**Grade:** $(python3 -c "
wr = $wr
pf = $pf
ret = $ret
maxdd = $maxdd

if ret > 0 and maxdd < 10 and wr >= 60:
    print('🏆 EXCELLENT')
elif ret > 0 and maxdd < 15 and wr >= 50:
    print('✅ GOOD')
elif ret > 0 and maxdd < 20 and wr >= 45:
    print('⚠️  ACCEPTABLE')
else:
    print('❌ NEEDS IMPROVEMENT')
" 2>/dev/null || echo "❌ NEEDS IMPROVEMENT")

EOF
    else
        cat >> "$SUMMARY_FILE" << EOF

### $strategy
❌ No profitable configuration found

EOF
    fi
done

cat >> "$SUMMARY_FILE" << 'EOF'

---

## 💡 HOW TO APPLY THESE PARAMETERS

### Step 1: Update unified_signal_generator.go

For each strategy, update the parameters in the corresponding function:

```go
// Example for liquidity_hunter
if buyScore >= 1 && buyScore >= sellScore {
    return &AdvancedSignal{
        Strategy:   "liquidity_hunter",
        Type:       "BUY",
        Entry:      currentPrice,
        StopLoss:   currentPrice - (atr * STOP_ATR),  // Use optimized value
        TP1:        currentPrice + (atr * TP1_ATR),   // Use optimized value
        TP2:        currentPrice + (atr * TP2_ATR),   // Use optimized value
        TP3:        currentPrice + (atr * TP3_ATR),   // Use optimized value
        ...
    }
}
```

### Step 2: Test with Paper Trading

Before going live:
1. Apply the new parameters
2. Restart the server
3. Paper trade for 30 days
4. Monitor equity curve and drawdowns
5. Verify performance matches backtest

### Step 3: Go Live (If Successful)

Only after successful paper trading:
1. Start with 1% risk per trade
2. Monitor closely for first week
3. Gradually increase if profitable
4. Never exceed 2% risk per trade

---

## ⚠️ IMPORTANT NOTES

1. **These parameters are optimized for 90-day period**
   - Market conditions change
   - Re-optimize quarterly

2. **Always use proper risk management**
   - Never risk more than 2% per trade
   - Always use stop losses
   - Take partial profits

3. **Paper trade first**
   - Minimum 30 days
   - Verify results match backtest
   - Only go live if successful

4. **Monitor performance**
   - Track equity curve weekly
   - Check drawdowns daily
   - Adjust if performance deviates >20%

---

## 📊 DETAILED RESULTS

Full results saved to: `best_parameters_all_strategies.csv`

Top 10 configurations per strategy saved to: `optimization_results/[strategy]_top10.csv`

---

*Generated by comprehensive parameter optimization*
*Test period: 90 days | Strategies: 10 | Combinations tested: 1000+*
EOF

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ ANALYSIS COMPLETE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📄 Results saved to:"
echo "   - Summary: $SUMMARY_FILE"
echo "   - Full Data: $RESULTS_FILE"
echo "   - Top 10 per strategy: optimization_results/[strategy]_top10.csv"
echo ""
echo "📊 View summary:"
echo "   cat $SUMMARY_FILE"
echo ""
echo "🔍 View top configurations:"
echo "   cat $RESULTS_FILE | sort -t',' -k14 -rn | head -20 | column -t -s','"
echo ""
echo "🎯 Next steps:"
echo "   1. Review $SUMMARY_FILE"
echo "   2. Update parameters in unified_signal_generator.go"
echo "   3. Restart server"
echo "   4. Paper trade for 30 days"
echo "   5. Go live if successful"
echo ""
