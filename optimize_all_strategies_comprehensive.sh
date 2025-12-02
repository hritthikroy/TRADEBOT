#!/bin/bash

echo "🔬 COMPREHENSIVE STRATEGY OPTIMIZATION"
echo "======================================"
echo ""
echo "This will optimize parameters for ALL strategies to find the best settings."
echo "Testing combinations of:"
echo "  - Confluence levels: 4, 5, 6, 7, 8"
echo "  - Stop Loss ATR: 0.5, 1.0, 1.5, 2.0"
echo "  - Take Profit ATR: 2.0, 3.0, 4.0, 5.0"
echo "  - Risk per trade: 1%, 1.5%, 2%, 2.5%"
echo ""
echo "⏱️  This may take 5-10 minutes..."
echo ""

# Configuration
SYMBOL="BTCUSDT"
START_BALANCE=500
DAYS=90

# Run optimization
curl -X POST http://localhost:8080/api/v1/backtest/optimize-all \
  -H "Content-Type: application/json" \
  -d "{
    \"symbol\": \"$SYMBOL\",
    \"startBalance\": $START_BALANCE,
    \"days\": $DAYS
  }" | jq '.'

echo ""
echo "✅ Optimization complete!"
echo ""
echo "📊 Check the backend logs for detailed results"
echo "💡 The best parameters for each strategy are shown above"
