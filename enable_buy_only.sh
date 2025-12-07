#!/bin/bash

echo "🎯 Enabling BUY ONLY Mode"
echo "========================"
echo ""

# Enable BUY trades, disable SELL trades
curl -s -X POST http://localhost:8080/api/v1/settings \
  -H "Content-Type: application/json" \
  -d '{
    "filterBuy": false,
    "filterSell": true,
    "strategies": ["session_trader"]
  }' | jq '.'

echo ""
echo "✅ BUY ONLY Mode Enabled!"
echo ""
echo "Settings:"
echo "- BUY trades: ENABLED ✅"
echo "- SELL trades: DISABLED ❌"
echo ""
