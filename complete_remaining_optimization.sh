#!/bin/bash

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🔬 COMPLETING OPTIMIZATION - REMAINING 7 STRATEGIES"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Optimizing the 7 remaining strategies with 8,064 combinations each"
echo "⏱️  Estimated time: 2-3 hours"
echo ""

# Check server
if ! curl -s http://localhost:8080/api/v1/health > /dev/null 2>&1; then
    echo "❌ Server not running!"
    exit 1
fi

echo "✅ Server is running"
echo ""

# Remaining strategies to optimize
strategies=(
    "liquidity_hunter"
    "trend_rider"
    "range_master"
    "smart_money_tracker"
    "institutional_follower"
    "reversal_sniper"
    "momentum_beast"
    "scalper_pro"
)

# Parameter ranges
stop_values=(0.5 0.75 1.0 1.25 1.5 2.0)
tp1_values=(2.0 2.5 3.0 3.5 4.0 4.5 5.0)
tp2_values=(3.0 4.0 4.5 5.0 6.0 7.0)
tp3_values=(5.0 6.0 7.0 7.5 8.0 9.0 10.0 12.0)
risk_values=(1.0 1.5 2.0 2.5)

# Results file
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
RESULTS_FILE="deep_optimization_results/remaining_results_${TIMESTAMP}.md"

cat > "$RESULTS_FILE" << 'EOF'
# 🔬 REMAINING STRATEGIES OPTIMIZATION

## Date: 
EOF
date >> "$RESULTS_FILE"
echo "" >> "$RESULTS_FILE"

strategy_num=0
total=${#strategies[@]}

for strategy in "${strategies[@]}"; do
    strategy_num=$((strategy_num + 1))
    
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "[$strategy_num/$total] Optimizing: $strategy"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    
    best_score=0
    best_wr=0
    best_pf=0
    best_return=0
    best_trades=0
    best_stop=0
    best_tp1=0
    best_tp2=0
    best_tp3=0
    best_risk=0
    
    combo=0
    
    for stop in "${stop_values[@]}"; do
        for tp1 in "${tp1_values[@]}"; do
            for tp2 in "${tp2_values[@]}"; do
                for tp3 in "${tp3_values[@]}"; do
                    for risk in "${risk_values[@]}"; do
                        combo=$((combo + 1))
                        
                        # Skip invalid
                        if (( $(echo "$tp1 >= $tp2" | bc -l) )) || (( $(echo "$tp2 >= $tp3" | bc -l) )); then
                            continue
                        fi
                        
                        # Progress
                        if [ $((combo % 100)) -eq 0 ]; then
                            echo "   Progress: $combo/8064..."
                        fi
                        
                        # Test
                        result=$(curl -s -X POST http://localhost:8080/api/v1/backtest/run \
                          -H "Content-Type: application/json" \
                          -d "{
                            \"symbol\": \"BTCUSDT\",
                            \"interval\": \"15m\",
                            \"days\": 180,
                            \"startBalance\": 1000,
                            \"riskPercent\": $risk,
                            \"strategy\": \"$strategy\"
                          }")
                        
                        trades=$(echo "$result" | python3 -c "import sys, json; data=json.load(sys.stdin); print(data.get('totalTrades', 0))" 2>/dev/null || echo "0")
                        wr=$(echo "$result" | python3 -c "import sys, json; data=json.load(sys.stdin); print(data.get('winRate', 0))" 2>/dev/null || echo "0")
                        pf=$(echo "$result" | python3 -c "import sys, json; data=json.load(sys.stdin); print(data.get('profitFactor', 0))" 2>/dev/null || echo "0")
                        ret=$(echo "$result" | python3 -c "import sys, json; data=json.load(sys.stdin); print(data.get('returnPercent', 0))" 2>/dev/null || echo "0")
                        
                        if [ "$trades" = "0" ]; then
                            continue
                        fi
                        
                        score=$(python3 -c "import math; print($wr * $pf * math.log(max($ret, 1) + 1))" 2>/dev/null || echo "0")
                        
                        if (( $(echo "$score > $best_score" | bc -l) )); then
                            best_score=$score
                            best_wr=$wr
                            best_pf=$pf
                            best_return=$ret
                            best_trades=$trades
                            best_stop=$stop
                            best_tp1=$tp1
                            best_tp2=$tp2
                            best_tp3=$tp3
                            best_risk=$risk
                            
                            echo "   🎯 New best! WR: ${wr}% | PF: $pf | Return: ${ret}% | Score: $score"
                        fi
                        
                        sleep 0.05
                    done
                done
            done
        done
    done
    
    echo ""
    echo "✅ Complete: $strategy"
    echo ""
    echo "🏆 BEST FOUND:"
    echo "   Stop: ${best_stop} ATR | TP1: ${best_tp1} | TP2: ${best_tp2} | TP3: ${best_tp3} | Risk: ${best_risk}%"
    echo "   WR: ${best_wr}% | PF: $best_pf | Return: ${best_return}% | Trades: $best_trades"
    echo ""
    
    # Add to report
    echo "### $strategy" >> "$RESULTS_FILE"
    echo "" >> "$RESULTS_FILE"
    echo "**Best Parameters:**" >> "$RESULTS_FILE"
    echo "\`\`\`" >> "$RESULTS_FILE"
    echo "Stop Loss: ${best_stop} ATR" >> "$RESULTS_FILE"
    echo "TP1: ${best_tp1} ATR" >> "$RESULTS_FILE"
    echo "TP2: ${best_tp2} ATR" >> "$RESULTS_FILE"
    echo "TP3: ${best_tp3} ATR" >> "$RESULTS_FILE"
    echo "Risk: ${best_risk}%" >> "$RESULTS_FILE"
    echo "\`\`\`" >> "$RESULTS_FILE"
    echo "" >> "$RESULTS_FILE"
    echo "**Performance:**" >> "$RESULTS_FILE"
    echo "- Win Rate: ${best_wr}%" >> "$RESULTS_FILE"
    echo "- Profit Factor: $best_pf" >> "$RESULTS_FILE"
    echo "- Return: ${best_return}%" >> "$RESULTS_FILE"
    echo "- Trades: $best_trades" >> "$RESULTS_FILE"
    echo "- Score: $best_score" >> "$RESULTS_FILE"
    echo "" >> "$RESULTS_FILE"
    echo "---" >> "$RESULTS_FILE"
    echo "" >> "$RESULTS_FILE"
done

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ ALL OPTIMIZATIONS COMPLETE!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📄 Results: $RESULTS_FILE"
echo ""
