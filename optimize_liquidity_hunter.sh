#!/bin/bash

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🎯 OPTIMIZING LIQUIDITY_HUNTER FOR MAXIMUM PROFITABILITY"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Goal: High Win Rate + High Return + Low Drawdown + High Profit Factor"
echo ""
echo "Testing on best timeframes: 1d, 8h, 2h"
echo "This will take 10-15 minutes..."
echo ""

# Check if server is running
if ! curl -s http://localhost:8080/api/v1/health > /dev/null 2>&1; then
    echo "❌ Server is not running!"
    exit 1
fi

timeframes=("1d" "8h" "2h")
best_overall_score=0
best_overall_config=""

for timeframe in "${timeframes[@]}"; do
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "🔍 OPTIMIZING: liquidity_hunter on $timeframe"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    
    result=$(curl -s -X POST http://localhost:8080/api/v1/backtest/world-class-optimize \
      -H "Content-Type: application/json" \
      -d "{
        \"symbol\": \"BTCUSDT\",
        \"days\": 180,
        \"startBalance\": 1000,
        \"strategies\": [\"liquidity_hunter\"]
      }")
    
    # Extract best parameters
    stopATR=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    result = data['results']['liquidity_hunter']
    print(result['bestParams']['StopATR'])
except:
    print(0)
" 2>/dev/null || echo "0")
    
    tp1ATR=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    result = data['results']['liquidity_hunter']
    print(result['bestParams']['TP1ATR'])
except:
    print(0)
" 2>/dev/null || echo "0")
    
    tp2ATR=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    result = data['results']['liquidity_hunter']
    print(result['bestParams']['TP2ATR'])
except:
    print(0)
" 2>/dev/null || echo "0")
    
    tp3ATR=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    result = data['results']['liquidity_hunter']
    print(result['bestParams']['TP3ATR'])
except:
    print(0)
" 2>/dev/null || echo "0")
    
    riskPct=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    result = data['results']['liquidity_hunter']
    print(result['bestParams']['RiskPercent'])
except:
    print(0)
" 2>/dev/null || echo "0")
    
    # Extract performance metrics
    winRate=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    result = data['results']['liquidity_hunter']
    if result['backtestResult']:
        print(f\"{result['backtestResult']['winRate']:.2f}\")
    else:
        print(0)
except:
    print(0)
" 2>/dev/null || echo "0")
    
    profitFactor=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    result = data['results']['liquidity_hunter']
    if result['backtestResult']:
        print(f\"{result['backtestResult']['profitFactor']:.2f}\")
    else:
        print(0)
except:
    print(0)
" 2>/dev/null || echo "0")
    
    returnPct=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    result = data['results']['liquidity_hunter']
    if result['backtestResult']:
        print(f\"{result['backtestResult']['returnPercent']:.2f}\")
    else:
        print(0)
except:
    print(0)
" 2>/dev/null || echo "0")
    
    maxDD=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    result = data['results']['liquidity_hunter']
    if result['backtestResult']:
        print(f\"{result['backtestResult']['maxDrawdown']:.2f}\")
    else:
        print(0)
except:
    print(0)
" 2>/dev/null || echo "0")
    
    totalTrades=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    result = data['results']['liquidity_hunter']
    if result['backtestResult']:
        print(result['backtestResult']['totalTrades'])
    else:
        print(0)
except:
    print(0)
" 2>/dev/null || echo "0")
    
    score=$(echo "$result" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    result = data['results']['liquidity_hunter']
    print(f\"{result['bestScore']:.2f}\")
except:
    print(0)
" 2>/dev/null || echo "0")
    
    echo "✅ Optimization Complete!"
    echo ""
    echo "📊 BEST PARAMETERS FOUND:"
    echo "   Stop Loss: ${stopATR} ATR"
    echo "   TP1: ${tp1ATR} ATR (33%)"
    echo "   TP2: ${tp2ATR} ATR (33%)"
    echo "   TP3: ${tp3ATR} ATR (34%)"
    echo "   Risk: ${riskPct}%"
    echo ""
    echo "📈 PERFORMANCE:"
    echo "   Win Rate: ${winRate}%"
    echo "   Profit Factor: ${profitFactor}"
    echo "   Return: ${returnPct}%"
    echo "   Max Drawdown: ${maxDD}%"
    echo "   Total Trades: ${totalTrades}"
    echo "   Score: ${score}"
    echo ""
    
    # Track best overall
    if (( $(echo "$score > $best_overall_score" | bc -l) )); then
        best_overall_score=$score
        best_overall_config="Timeframe: $timeframe | Stop: ${stopATR} | TP1: ${tp1ATR} | TP2: ${tp2ATR} | TP3: ${tp3ATR} | Risk: ${riskPct}% | WR: ${winRate}% | PF: ${profitFactor} | Return: ${returnPct}% | DD: ${maxDD}% | Trades: ${totalTrades}"
    fi
done

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🏆 BEST CONFIGURATION FOUND"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "$best_overall_config"
echo ""
echo "Score: $best_overall_score"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📝 NEXT STEPS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "1. Update unified_signal_generator.go with these parameters"
echo "2. Test with paper trading for 30 days"
echo "3. Monitor performance closely"
echo "4. Only go live if paper trading is successful"
echo ""
echo "⚠️  Remember: Past performance does not guarantee future results!"
echo ""
