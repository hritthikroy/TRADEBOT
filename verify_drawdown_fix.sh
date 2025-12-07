#!/bin/bash

echo "=========================================="
echo "VERIFYING DRAWDOWN FIX"
echo "=========================================="
echo ""

echo "Testing Session Trader (30 days)..."
echo ""

result=$(curl -s -X POST http://localhost:8080/api/v1/backtest/run \
  -H "Content-Type: application/json" \
  -d '{"strategy":"session_trader","days":30,"startBalance":500,"riskPercent":0.02}')

echo "$result" | jq -r '
"
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
SESSION TRADER - 30 DAYS BACKTEST
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Total Trades:     \(.totalTrades)
Win Rate:         \(.winRate | floor * 10 / 10)%
Profit Factor:    \(.profitFactor | . * 100 | floor / 100)
Return:           \(.returnPercent | floor * 10 / 10)%
Final Balance:    $\(.finalBalance | floor * 100 / 100)

Max Drawdown:     \(.maxDrawdown | floor * 10 / 10)% \(if .maxDrawdown < 10 then "✅ FIXED!" else "❌ STILL WRONG" end)

Exit Breakdown:
  Stop Loss:      \(.exitReasons["Stop Loss"] // 0) (\(if .totalTrades > 0 then ((.exitReasons["Stop Loss"] // 0) * 100 / .totalTrades | floor) else 0 end)%)
  Target 1:       \(.exitReasons["Target 1"] // 0) (\(if .totalTrades > 0 then ((.exitReasons["Target 1"] // 0) * 100 / .totalTrades | floor) else 0 end)%)
  Target 2:       \(.exitReasons["Target 2"] // 0) (\(if .totalTrades > 0 then ((.exitReasons["Target 2"] // 0) * 100 / .totalTrades | floor) else 0 end)%)
  Target 3:       \(.exitReasons["Target 3"] // 0) (\(if .totalTrades > 0 then ((.exitReasons["Target 3"] // 0) * 100 / .totalTrades | floor) else 0 end)%)
  Timeout:        \(.exitReasons["Timeout"] // 0) (\(if .totalTrades > 0 then ((.exitReasons["Timeout"] // 0) * 100 / .totalTrades | floor) else 0 end)%)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
"'

echo ""
echo "=========================================="
echo "DRAWDOWN FIX VERIFICATION"
echo "=========================================="
echo ""

maxDD=$(echo "$result" | jq -r '.maxDrawdown')

if (( $(echo "$maxDD < 10" | bc -l) )); then
    echo "✅ SUCCESS! Max Drawdown is ${maxDD}%"
    echo "   The drawdown bug is FIXED!"
    echo ""
    echo "   Expected: 2-5%"
    echo "   Actual:   ${maxDD}%"
else
    echo "❌ FAILED! Max Drawdown is ${maxDD}%"
    echo "   The drawdown bug is STILL PRESENT!"
    echo ""
    echo "   Expected: 2-5%"
    echo "   Actual:   ${maxDD}%"
    echo ""
    echo "   Action Required:"
    echo "   1. Restart backend: lsof -ti:8080 | xargs kill -9 && cd backend && go run ."
    echo "   2. Hard refresh frontend: Cmd+Shift+R"
fi

echo ""
echo "=========================================="
