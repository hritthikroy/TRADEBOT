#!/bin/bash

echo "🌍 WORLD-CLASS PARAMETER OPTIMIZATION"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "This will test THOUSANDS of parameter combinations to find the absolute"
echo "best settings for all 10 strategies."
echo ""
echo "⏱️  Estimated Time: 30-60 minutes"
echo "🎯 Goal: Win Rate > 60%, Profit Factor > 3.0, Max DD < 15%"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check if server is running
if ! curl -s http://localhost:8080/api/v1/health > /dev/null; then
    echo "❌ Server is not running!"
    echo "Please start the server first:"
    echo "  cd backend && go run ."
    exit 1
fi

echo "✅ Server is running"
echo ""
echo "🚀 Starting optimization..."
echo ""

# Run optimization
curl -X POST http://localhost:8080/api/v1/backtest/world-class-optimize \
    -H "Content-Type: application/json" \
    -o WORLD_CLASS_OPTIMIZATION_RESULTS.json

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ OPTIMIZATION COMPLETE!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Results saved to: WORLD_CLASS_OPTIMIZATION_RESULTS.json"
echo ""
echo "View results:"
echo "  cat WORLD_CLASS_OPTIMIZATION_RESULTS.json | jq ."
echo ""
echo "Next steps:"
echo "  1. Review the best parameters for each strategy"
echo "  2. Apply them to your live_signal_handler.go"
echo "  3. Test with live signals"
echo ""
