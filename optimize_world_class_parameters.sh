#!/bin/bash

# World-Class Strategy Optimization
# Tests extensive parameter combinations to find the absolute best settings
# Focus: Low Drawdown + High Win Rate + High Profit Factor

echo "🌍 WORLD-CLASS STRATEGY OPTIMIZATION"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Testing Parameters:"
echo "  • Stop Loss: 0.5, 0.75, 1.0, 1.25, 1.5, 2.0 ATR"
echo "  • TP1: 2.0, 2.5, 3.0, 3.5, 4.0, 5.0 ATR"
echo "  • TP2: 3.0, 4.0, 4.5, 5.0, 6.0, 7.5 ATR"
echo "  • TP3: 5.0, 6.0, 7.5, 10.0, 12.5, 15.0 ATR"
echo "  • Risk: 0.5%, 1.0%, 1.5%, 2.0%, 2.5% per trade"
echo ""
echo "Optimization Goals:"
echo "  1. Win Rate > 60%"
echo "  2. Profit Factor > 3.0"
echo "  3. Max Drawdown < 15%"
echo "  4. Return > 500%"
echo "  5. Total Trades > 20"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Configuration
SYMBOL="BTCUSDT"
DAYS=180
START_BALANCE=1000
OUTPUT_FILE="WORLD_CLASS_OPTIMIZATION_RESULTS.json"

# Strategies to optimize
STRATEGIES=(
    "session_trader"
    "breakout_master"
    "liquidity_hunter"
    "trend_rider"
    "range_master"
    "smart_money_tracker"
    "institutional_follower"
    "reversal_sniper"
    "momentum_beast"
    "scalper_pro"
)

# Parameter ranges to test
STOP_LOSS_VALUES=(0.5 0.75 1.0 1.25 1.5 2.0)
TP1_VALUES=(2.0 2.5 3.0 3.5 4.0 5.0)
TP2_VALUES=(3.0 4.0 4.5 5.0 6.0 7.5)
TP3_VALUES=(5.0 6.0 7.5 10.0 12.5 15.0)
RISK_VALUES=(0.5 1.0 1.5 2.0 2.5)

# Initialize results
echo "{" > $OUTPUT_FILE
echo "  \"optimizationDate\": \"$(date -u +"%Y-%m-%dT%H:%M:%SZ")\"," >> $OUTPUT_FILE
echo "  \"symbol\": \"$SYMBOL\"," >> $OUTPUT_FILE
echo "  \"days\": $DAYS," >> $OUTPUT_FILE
echo "  \"startBalance\": $START_BALANCE," >> $OUTPUT_FILE
echo "  \"strategies\": {" >> $OUTPUT_FILE

FIRST_STRATEGY=true

# Optimize each strategy
for STRATEGY in "${STRATEGIES[@]}"; do
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "🎯 Optimizing: $STRATEGY"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    BEST_SCORE=0
    BEST_PARAMS=""
    BEST_RESULT=""
    TOTAL_TESTS=0
    
    # Test all parameter combinations
    for STOP in "${STOP_LOSS_VALUES[@]}"; do
        for TP1 in "${TP1_VALUES[@]}"; do
            for TP2 in "${TP2_VALUES[@]}"; do
                for TP3 in "${TP3_VALUES[@]}"; do
                    for RISK in "${RISK_VALUES[@]}"; do
                        # Validate: TP1 < TP2 < TP3
                        if (( $(echo "$TP1 < $TP2" | bc -l) )) && (( $(echo "$TP2 < $TP3" | bc -l) )); then
                            TOTAL_TESTS=$((TOTAL_TESTS + 1))
                            
                            # Run backtest
                            RESULT=$(curl -s -X POST http://localhost:8080/api/v1/backtest/run \
                                -H "Content-Type: application/json" \
                                -d "{
                                    \"symbol\": \"$SYMBOL\",
                                    \"interval\": \"15m\",
                                    \"days\": $DAYS,
                                    \"startBalance\": $START_BALANCE,
                                    \"riskPercent\": $(echo "$RISK / 100" | bc -l),
                                    \"strategy\": \"$STRATEGY\",
                                    \"stopATR\": $STOP,
                                    \"tp1ATR\": $TP1,
                                    \"tp2ATR\": $TP2,
                                    \"tp3ATR\": $TP3
                                }")
                            
                            # Extract metrics
                            WIN_RATE=$(echo $RESULT | jq -r '.winRate // 0')
                            PROFIT_FACTOR=$(echo $RESULT | jq -r '.profitFactor // 0')
                            RETURN=$(echo $RESULT | jq -r '.returnPercent // 0')
                            MAX_DD=$(echo $RESULT | jq -r '.maxDrawdown // 100')
                            TOTAL_TRADES=$(echo $RESULT | jq -r '.totalTrades // 0')
                            
                            # Calculate score (weighted formula)
                            # Score = (WinRate * 2) + (ProfitFactor * 10) + (Return / 10) - (MaxDD * 2) + (Trades / 2)
                            SCORE=$(echo "scale=2; ($WIN_RATE * 2) + ($PROFIT_FACTOR * 10) + ($RETURN / 10) - ($MAX_DD * 2) + ($TOTAL_TRADES / 2)" | bc -l)
                            
                            # Check if this is the best so far
                            if (( $(echo "$SCORE > $BEST_SCORE" | bc -l) )) && \
                               (( $(echo "$WIN_RATE >= 55" | bc -l) )) && \
                               (( $(echo "$PROFIT_FACTOR >= 2.5" | bc -l) )) && \
                               (( $(echo "$MAX_DD <= 20" | bc -l) )) && \
                               (( $(echo "$TOTAL_TRADES >= 15" | bc -l) )); then
                                BEST_SCORE=$SCORE
                                BEST_PARAMS="Stop:$STOP TP1:$TP1 TP2:$TP2 TP3:$TP3 Risk:$RISK%"
                                BEST_RESULT=$RESULT
                                
                                echo "  ✨ New Best! Score: $SCORE | WR: $WIN_RATE% | PF: $PROFIT_FACTOR | DD: $MAX_DD% | Trades: $TOTAL_TRADES"
                            fi
                            
                            # Progress indicator
                            if [ $((TOTAL_TESTS % 10)) -eq 0 ]; then
                                echo "  ⏳ Tested $TOTAL_TESTS combinations..."
                            fi
                        fi
                    done
                done
            done
        done
    done
    
    echo ""
    echo "✅ Optimization Complete for $STRATEGY"
    echo "   Total Tests: $TOTAL_TESTS"
    echo "   Best Score: $BEST_SCORE"
    echo "   Best Params: $BEST_PARAMS"
    echo ""
    
    # Add to results file
    if [ "$FIRST_STRATEGY" = true ]; then
        FIRST_STRATEGY=false
    else
        echo "," >> $OUTPUT_FILE
    fi
    
    echo "    \"$STRATEGY\": {" >> $OUTPUT_FILE
    echo "      \"totalTests\": $TOTAL_TESTS," >> $OUTPUT_FILE
    echo "      \"bestScore\": $BEST_SCORE," >> $OUTPUT_FILE
    echo "      \"bestParameters\": \"$BEST_PARAMS\"," >> $OUTPUT_FILE
    echo "      \"result\": $BEST_RESULT" >> $OUTPUT_FILE
    echo -n "    }" >> $OUTPUT_FILE
done

# Close JSON
echo "" >> $OUTPUT_FILE
echo "  }" >> $OUTPUT_FILE
echo "}" >> $OUTPUT_FILE

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🎉 WORLD-CLASS OPTIMIZATION COMPLETE!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Results saved to: $OUTPUT_FILE"
echo ""
echo "Next Steps:"
echo "  1. Review the results in $OUTPUT_FILE"
echo "  2. Apply the best parameters to your strategies"
echo "  3. Run live testing with the optimized settings"
echo ""
