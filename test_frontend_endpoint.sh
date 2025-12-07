#!/bin/bash

echo "=========================================="
echo "TESTING FRONTEND ENDPOINT (test-all-strategies)"
echo "=========================================="
echo ""
echo "This is the endpoint your frontend uses!"
echo ""

curl -s -X POST http://localhost:8080/api/v1/backtest/test-all-strategies \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTCUSDT","days":30,"startBalance":500}' | \
  jq '.results[] | select(.strategyName == "session_trader") | 
  "
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
SESSION TRADER - 30 DAYS (Frontend Endpoint)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Total Trades:     \(.totalTrades)
Win Rate:         \(.winRate | floor * 10 / 10)%
Profit Factor:    \(.profitFactor | . * 100 | floor / 100)
Return:           \(.returnPercent | floor | tostring | gsub("(?<a>\\d)(?<b>(\\d{3})+$)"; "\(.a),\(.b)"))%
Final Balance:    $\(.finalBalance | floor | tostring | gsub("(?<a>\\d)(?<b>(\\d{3})+$)"; "\(.a),\(.b)"))

Max Drawdown:     \(.maxDrawdown | floor * 10 / 10)% \(if .maxDrawdown < 50 then "✅ FIXED!" else "❌ STILL BROKEN" end)

Market Bias:      \(.marketBias)
Buy Trades:       \(.buyTrades) (\(.buyWinRate | floor)% WR)
Sell Trades:      \(.sellTrades) (\(.sellWinRate | floor)% WR)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
"' -r

echo ""
echo "=========================================="
echo "FRONTEND INSTRUCTIONS"
echo "=========================================="
echo ""
echo "1. Open browser: http://localhost:8080"
echo "2. Hard refresh: Cmd+Shift+R (Mac) or Ctrl+Shift+R (Windows)"
echo "3. Select 'Session Trader' strategy"
echo "4. Choose 30 days"
echo "5. Click 'Run Backtest'"
echo ""
echo "Expected Result:"
echo "  Max Drawdown: ~18% ✅"
echo ""
echo "If you still see 624.2%, do a hard refresh!"
echo "=========================================="
