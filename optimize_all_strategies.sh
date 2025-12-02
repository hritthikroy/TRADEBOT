#!/bin/bash

echo "╔══════════════════════════════════════════════════════════════╗"
echo "║     🔬 COMPREHENSIVE STRATEGY PARAMETER OPTIMIZATION 🔬      ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo ""

# Test parameters
SYMBOLS=("BTCUSDT" "ETHUSDT" "BNBUSDT")
BALANCES=(500 1000 2000)
DAYS=(30 60 90 180)

echo "📊 Testing Configuration:"
echo "   Symbols: ${SYMBOLS[@]}"
echo "   Start Balances: ${BALANCES[@]}"
echo "   Time Periods: ${DAYS[@]} days"
echo ""

BEST_SCORE=0
BEST_CONFIG=""
BEST_RESULT=""

for SYMBOL in "${SYMBOLS[@]}"; do
    for BALANCE in "${BALANCES[@]}"; do
        for DAY in "${DAYS[@]}"; do
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            echo "🧪 Testing: $SYMBOL | Balance: \$$BALANCE | Period: $DAY days"
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            
            RESULT=$(curl -s -X POST http://localhost:8080/api/v1/backtest/test-all-strategies \
                -H "Content-Type: application/json" \
                -d "{\"symbol\":\"$SYMBOL\",\"startBalance\":$BALANCE,\"days\":$DAY}")
            
            # Extract best strategy score
            SCORE=$(echo "$RESULT" | grep -o '"score":[0-9.]*' | head -1 | cut -d':' -f2)
            
            if [ ! -z "$SCORE" ]; then
                echo "   Score: $SCORE"
                
                # Check if this is the best score
                if (( $(echo "$SCORE > $BEST_SCORE" | bc -l) )); then
                    BEST_SCORE=$SCORE
                    BEST_CONFIG="Symbol: $SYMBOL | Balance: \$$BALANCE | Days: $DAY"
                    BEST_RESULT="$RESULT"
                    echo "   🏆 NEW BEST CONFIGURATION!"
                fi
            fi
            
            echo ""
            sleep 2
        done
    done
done

echo ""
echo "╔══════════════════════════════════════════════════════════════╗"
echo "║                  🏆 OPTIMIZATION RESULTS 🏆                  ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo ""
echo "🥇 BEST CONFIGURATION:"
echo "   $BEST_CONFIG"
echo "   Score: $BEST_SCORE"
echo ""
echo "📊 DETAILED RESULTS:"
echo "$BEST_RESULT" | python3 -m json.tool
echo ""
echo "✅ Optimization Complete!"
echo ""
