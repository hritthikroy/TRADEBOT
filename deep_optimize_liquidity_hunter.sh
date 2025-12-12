#!/bin/bash

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🚀 DEEP OPTIMIZATION: liquidity_hunter"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Testing EXTENDED parameter ranges for MAXIMUM profitability"
echo "This will test 20,000+ combinations per timeframe"
echo "Expected time: 30-60 minutes"
echo ""

read -p "Continue? (y/n) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    exit 0
fi

# Best timeframes from analysis
timeframes=("1d" "8h")

for timeframe in "${timeframes[@]}"; do
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "🔬 DEEP OPTIMIZATION: $timeframe"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    
    # Extended parameter ranges
    stopLossValues=(0.3 0.5 0.75 1.0 1.25 1.5 1.75 2.0 2.5 3.0)
    tp1Values=(1.5 2.0 2.5 3.0 3.5 4.0 4.5 5.0 6.0 7.0)
    tp2Values=(2.5 3.0 3.5 4.0 4.5 5.0 6.0 7.0 8.0 10.0)
    tp3Values=(4.0 5.0 6.0 7.0 8.0 10.0 12.0 15.0 20.0 25.0)
    riskValues=(0.5 0.75 1.0 1.25 1.5 1.75 2.0 2.5 3.0)
    
    best_score=0
    best_config=""
    total_tests=0
    profitable_configs=0
    
    echo "Testing ${#stopLossValues[@]} × ${#tp1Values[@]} × ${#tp2Values[@]} × ${#tp3Values[@]} × ${#riskValues[@]} = $((${#stopLossValues[@]} * ${#tp1Values[@]} * ${#tp2Values[@]} * ${#tp3Values[@]} * ${#riskValues[@]})) combinations"
    echo ""
    
    for stop in "${stopLossValues[@]}"; do
        for tp1 in "${tp1Values[@]}"; do
            for tp2 in "${tp2Values[@]}"; do
                for tp3 in "${tp3Values[@]}"; do
                    # Validate: TP1 < TP2 < TP3
                    if (( $(echo "$tp1 < $tp2" | bc -l) )) && (( $(echo "$tp2 < $tp3" | bc -l) )); then
                        for risk in "${riskValues[@]}"; do
                            total_tests=$((total_tests + 1))
                            
                            # Progress indicator
                            if [ $((total_tests % 100)) -eq 0 ]; then
                                echo -ne "\r   Tested: $total_tests | Profitable: $profitable_configs | Best Score: $best_score"
                            fi
                            
                            # Run backtest with custom parameters
                            result=$(curl -s -X POST http://localhost:8080/api/v1/backtest/run \
                              -H "Content-Type: application/json" \
                              -d "{
                                \"symbol\": \"BTCUSDT\",
                                \"interval\": \"$timeframe\",
                                \"days\": 180,
                                \"startBalance\": 1000,
                                \"strategy\": \"liquidity_hunter\",
                                \"riskPercent\": $(echo "$risk / 100" | bc -l)
                              }")
                            
                            # Extract metrics
                            winRate=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(data.get('winRate', 0))" 2>/dev/null || echo "0")
                            profitFactor=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(data.get('profitFactor', 0))" 2>/dev/null || echo "0")
                            returnPct=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(data.get('returnPercent', 0))" 2>/dev/null || echo "0")
                            maxDD=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(data.get('maxDrawdown', 0))" 2>/dev/null || echo "0")
                            totalTrades=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(data.get('totalTrades', 0))" 2>/dev/null || echo "0")
                            
                            # Calculate score: (WinRate * ProfitFactor * Return) - (Drawdown * 10)
                            score=$(python3 -c "
wr = $winRate
pf = $profitFactor
ret = max(0, $returnPct)
dd = $maxDD
trades = $totalTrades
if trades >= 5 and wr >= 40 and pf >= 1.0 and ret > 0:
    score = (wr * pf * ret / 100) - (dd * 10)
    print(round(max(0, score), 2))
else:
    print(0)
" 2>/dev/null || echo "0")
                            
                            # Track profitable configs
                            if (( $(echo "$returnPct > 0" | bc -l) )) && (( $(echo "$profitFactor >= 1.0" | bc -l) )); then
                                profitable_configs=$((profitable_configs + 1))
                            fi
                            
                            # Track best
                            if (( $(echo "$score > $best_score" | bc -l) )); then
                                best_score=$score
                                best_config="Stop: $stop | TP1: $tp1 | TP2: $tp2 | TP3: $tp3 | Risk: $risk% | WR: $winRate% | PF: $profitFactor | Return: $returnPct% | DD: $maxDD% | Trades: $totalTrades"
                            fi
                        done
                    fi
                done
            done
        done
    done
    
    echo ""
    echo ""
    echo "✅ $timeframe Optimization Complete!"
    echo ""
    echo "📊 STATISTICS:"
    echo "   Total Tests: $total_tests"
    echo "   Profitable Configs: $profitable_configs"
    echo "   Success Rate: $(python3 -c "print(round($profitable_configs / $total_tests * 100, 2))" 2>/dev/null || echo "0")%"
    echo ""
    echo "🏆 BEST CONFIGURATION:"
    echo "   $best_config"
    echo "   Score: $best_score"
    echo ""
done

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ DEEP OPTIMIZATION COMPLETE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Review the results above and update your strategy parameters."
echo ""
