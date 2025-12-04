#!/bin/bash

echo "⏳ Waiting for server to be ready..."
echo ""

# Wait for server to respond
max_attempts=30
attempt=0

while [ $attempt -lt $max_attempts ]; do
    if curl -s http://localhost:8080/api/v1/health > /dev/null 2>&1; then
        echo "✅ Server is ready!"
        echo ""
        break
    fi
    attempt=$((attempt + 1))
    echo "   Attempt $attempt/$max_attempts..."
    sleep 1
done

if [ $attempt -eq $max_attempts ]; then
    echo "❌ Server is not responding after 30 seconds"
    echo ""
    echo "Please make sure the server is running:"
    echo "   cd backend && go run ."
    exit 1
fi

echo "🚀 Starting optimization with 30 days of data..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Run the optimization
curl -X POST http://localhost:8080/api/v1/backtest/world-class-optimize \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "days": 30,
    "startBalance": 1000
  }' 2>/dev/null | python3 -m json.tool > optimization_results.json

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Optimization complete! Results saved to optimization_results.json"
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "📊 QUICK SUMMARY"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    
    # Show summary for each strategy
    for strategy in session_trader breakout_master liquidity_hunter trend_rider range_master smart_money_tracker institutional_follower reversal_sniper momentum_beast scalper_pro; do
        score=$(cat optimization_results.json | grep -A 20 "\"$strategy\"" | grep "bestScore" | head -1 | grep -o '[0-9.]*' | head -1)
        trades=$(cat optimization_results.json | grep -A 30 "\"$strategy\"" | grep "totalTrades" | head -1 | grep -o '[0-9]*' | head -1)
        
        if [ ! -z "$score" ] && [ ! -z "$trades" ]; then
            printf "%-25s Score: %-8s Trades: %s\n" "$strategy" "$score" "$trades"
        fi
    done
    
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "📄 View full results: cat optimization_results.json | python3 -m json.tool"
    echo ""
else
    echo ""
    echo "❌ Optimization failed. Check server logs for errors."
    echo ""
fi
