#!/bin/bash

echo "📊 BEFORE vs AFTER AMD Phases Comparison"
echo "========================================="
echo ""

# Test function
test_strategy() {
    local days=$1
    local label=$2
    
    echo "Testing $label ($days days)..."
    curl -s -X POST http://localhost:8080/api/v1/backtest/test-all-strategies \
      -H "Content-Type: application/json" \
      -d "{
        \"symbol\": \"BTCUSDT\",
        \"days\": $days,
        \"startBalance\": 1000,
        \"filterBuy\": false,
        \"filterSell\": true
      }" | jq -r '.strategies[] | select(.name == "Session Trader") | 
        "Trades: \(.trades) | WR: \(.winRate)% | PF: \(.profitFactor) | DD: \(.maxDrawdown)% | Return: \(.totalReturn)%"'
}

# Check if backend is running
if ! pgrep -x "tradebot" > /dev/null; then
    echo "⚠️  Backend not running. Starting..."
    ./tradebot &
    sleep 3
fi

echo "🔴 TESTING ORIGINAL (from backup)..."
echo "===================================="
echo ""

# Restore original
cp backend/unified_signal_generator.go.backup backend/unified_signal_generator_amd.go.temp
cp backend/unified_signal_generator.go.backup backend/unified_signal_generator.go

# Rebuild
cd backend && go build -o ../tradebot 2>/dev/null && cd ..

# Restart backend
pkill tradebot
./tradebot &
sleep 3

echo "📊 Original Results:"
echo "-------------------"
test_strategy 30 "30 days"
echo ""
test_strategy 7 "7 days"
echo ""
test_strategy 5 "Bad period (Nov 30-Dec 4)"
echo ""

# Save original results
echo "Saving original results..."
curl -s -X POST http://localhost:8080/api/v1/backtest/test-all-strategies \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTCUSDT","days":30,"startBalance":1000,"filterBuy":false,"filterSell":true}' \
  > /tmp/original_results.json

echo ""
echo "🟢 TESTING WITH AMD PHASES..."
echo "============================="
echo ""

# Restore AMD version
cp backend/unified_signal_generator_amd.go.temp backend/unified_signal_generator.go

# Rebuild
cd backend && go build -o ../tradebot 2>/dev/null && cd ..

# Restart backend
pkill tradebot
./tradebot &
sleep 3

echo "📊 AMD Phase Results:"
echo "--------------------"
test_strategy 30 "30 days"
echo ""
test_strategy 7 "7 days"
echo ""
test_strategy 5 "Bad period (Nov 30-Dec 4)"
echo ""

# Save AMD results
echo "Saving AMD results..."
curl -s -X POST http://localhost:8080/api/v1/backtest/test-all-strategies \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTCUSDT","days":30,"startBalance":1000,"filterBuy":false,"filterSell":true}' \
  > /tmp/amd_results.json

echo ""
echo "📈 COMPARISON (30-day period)"
echo "============================="
echo ""

# Extract and compare
original_trades=$(jq -r '.strategies[] | select(.name == "Session Trader") | .trades' /tmp/original_results.json)
original_wr=$(jq -r '.strategies[] | select(.name == "Session Trader") | .winRate' /tmp/original_results.json)
original_pf=$(jq -r '.strategies[] | select(.name == "Session Trader") | .profitFactor' /tmp/original_results.json)
original_dd=$(jq -r '.strategies[] | select(.name == "Session Trader") | .maxDrawdown' /tmp/original_results.json)
original_return=$(jq -r '.strategies[] | select(.name == "Session Trader") | .totalReturn' /tmp/original_results.json)

amd_trades=$(jq -r '.strategies[] | select(.name == "Session Trader") | .trades' /tmp/amd_results.json)
amd_wr=$(jq -r '.strategies[] | select(.name == "Session Trader") | .winRate' /tmp/amd_results.json)
amd_pf=$(jq -r '.strategies[] | select(.name == "Session Trader") | .profitFactor' /tmp/amd_results.json)
amd_dd=$(jq -r '.strategies[] | select(.name == "Session Trader") | .maxDrawdown' /tmp/amd_results.json)
amd_return=$(jq -r '.strategies[] | select(.name == "Session Trader") | .totalReturn' /tmp/amd_results.json)

echo "Metric          | Original | AMD Phases | Change"
echo "----------------|----------|------------|--------"
echo "Trades          | $original_trades | $amd_trades | $(echo "$amd_trades - $original_trades" | bc)"
echo "Win Rate        | $original_wr% | $amd_wr% | $(echo "$amd_wr - $original_wr" | bc)%"
echo "Profit Factor   | $original_pf | $amd_pf | $(echo "$amd_pf - $original_pf" | bc)"
echo "Max Drawdown    | $original_dd% | $amd_dd% | $(echo "$amd_dd - $original_dd" | bc)%"
echo "Total Return    | $original_return% | $amd_return% | $(echo "$amd_return - $original_return" | bc)%"

echo ""
echo "🎯 RECOMMENDATION"
echo "================="
echo ""

# Calculate improvements
wr_improved=$(echo "$amd_wr > $original_wr + 3" | bc)
pf_improved=$(echo "$amd_pf > $original_pf * 1.15" | bc)
dd_improved=$(echo "$amd_dd < $original_dd * 0.9" | bc)

improvements=0
if [ "$wr_improved" -eq 1 ]; then improvements=$((improvements + 1)); fi
if [ "$pf_improved" -eq 1 ]; then improvements=$((improvements + 1)); fi
if [ "$dd_improved" -eq 1 ]; then improvements=$((improvements + 1)); fi

if [ $improvements -ge 2 ]; then
    echo "✅ KEEP AMD PHASES"
    echo ""
    echo "Reasons:"
    [ "$wr_improved" -eq 1 ] && echo "  ✅ Win rate improved by $(echo "$amd_wr - $original_wr" | bc)%"
    [ "$pf_improved" -eq 1 ] && echo "  ✅ Profit factor improved by $(echo "scale=2; ($amd_pf - $original_pf) / $original_pf * 100" | bc)%"
    [ "$dd_improved" -eq 1 ] && echo "  ✅ Drawdown reduced by $(echo "scale=2; ($original_dd - $amd_dd) / $original_dd * 100" | bc)%"
    echo ""
    echo "The AMD phase detection is working! Keep this version."
else
    echo "❌ ROLLBACK TO ORIGINAL"
    echo ""
    echo "Reasons:"
    [ "$wr_improved" -eq 0 ] && echo "  ❌ Win rate not significantly improved"
    [ "$pf_improved" -eq 0 ] && echo "  ❌ Profit factor not significantly improved"
    [ "$dd_improved" -eq 0 ] && echo "  ❌ Drawdown not significantly reduced"
    echo ""
    echo "AMD phases didn't improve results. Restore original:"
    echo "  cp backend/unified_signal_generator.go.backup backend/unified_signal_generator.go"
fi

echo ""
echo "📋 Detailed results saved:"
echo "  Original: /tmp/original_results.json"
echo "  AMD:      /tmp/amd_results.json"

# Cleanup
rm -f backend/unified_signal_generator_amd.go.temp
