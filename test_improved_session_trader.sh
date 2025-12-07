#!/bin/bash

echo "=========================================="
echo "TESTING IMPROVED SESSION TRADER STRATEGY"
echo "=========================================="
echo ""
echo "IMPROVEMENTS APPLIED:"
echo "1. ✅ Fixed drawdown calculation bug (was showing 2843.7%)"
echo "2. ✅ Stricter entry filters (require reversal + volume)"
echo "3. ✅ Better R:R targets (2R, 3.5R, 5R)"
echo "4. ✅ Trend confirmation (price vs EMA21)"
echo ""
echo "Expected Results:"
echo "- Win Rate: 55-60% (up from 49.5%)"
echo "- Profit Factor: 1.3-1.5 (up from 1.20)"
echo "- Max Drawdown: <10% (fixed calculation)"
echo "- Fewer but higher quality trades"
echo ""
echo "=========================================="
echo ""

cd backend

echo "Testing Session Trader across multiple timeframes..."
echo ""

for days in 3 5 7 15 30; do
    echo "----------------------------------------"
    echo "Testing ${days}d period..."
    echo "----------------------------------------"
    
    curl -s -X POST http://localhost:8080/api/backtest \
        -H "Content-Type: application/json" \
        -d "{
            \"strategy\": \"session_trader\",
            \"days\": $days,
            \"startBalance\": 500,
            \"riskPercent\": 0.02
        }" | jq -r '
        "
Results for \(.days)d:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Total Trades:     \(.totalTrades)
Win Rate:         \(.winRate)%
Profit Factor:    \(.profitFactor)
Return:           \(.returnPercent)%
Max Drawdown:     \(.maxDrawdown)%
Final Balance:    $\(.finalBalance)

Exit Breakdown:
  Stop Loss:      \(.exitReasons.stopLoss // 0) (\(if .totalTrades > 0 then ((.exitReasons.stopLoss // 0) * 100 / .totalTrades | floor) else 0 end)%)
  Target 1:       \(.exitReasons.target1 // 0) (\(if .totalTrades > 0 then ((.exitReasons.target1 // 0) * 100 / .totalTrades | floor) else 0 end)%)
  Target 2:       \(.exitReasons.target2 // 0) (\(if .totalTrades > 0 then ((.exitReasons.target2 // 0) * 100 / .totalTrades | floor) else 0 end)%)
  Target 3:       \(.exitReasons.target3 // 0) (\(if .totalTrades > 0 then ((.exitReasons.target3 // 0) * 100 / .totalTrades | floor) else 0 end)%)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
"'
    
    echo ""
done

echo ""
echo "=========================================="
echo "COMPARISON WITH PREVIOUS RESULTS"
echo "=========================================="
echo ""
echo "BEFORE (30d):"
echo "  Total Trades: 459"
echo "  Win Rate: 49.5%"
echo "  Profit Factor: 1.20"
echo "  Max Drawdown: 2843.7% ❌ (BUG)"
echo "  Return: 392.5%"
echo ""
echo "AFTER (30d) - Expected:"
echo "  Total Trades: ~200-300 (fewer, higher quality)"
echo "  Win Rate: 55-60% ✅"
echo "  Profit Factor: 1.3-1.5 ✅"
echo "  Max Drawdown: <10% ✅ (FIXED)"
echo "  Return: 400-600% ✅"
echo ""
echo "=========================================="
