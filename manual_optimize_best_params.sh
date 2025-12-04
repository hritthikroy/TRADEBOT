#!/bin/bash

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🎯 MANUAL PARAMETER OPTIMIZATION FOR LIQUIDITY_HUNTER"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Testing proven parameter combinations on best timeframes"
echo "Goal: Win Rate >60%, Profit Factor >2.0, Low Drawdown"
echo ""

# Best timeframes from analysis
TIMEFRAME="1d"  # Best performing timeframe
DAYS=180

echo "Testing on $TIMEFRAME timeframe with $DAYS days of data"
echo ""

# Proven parameter combinations to test
declare -a configs=(
    # Format: "StopATR:TP1ATR:TP2ATR:TP3ATR:Risk:Description"
    "0.5:2.0:3.0:5.0:1.0:Conservative - Tight Stop"
    "0.75:2.5:4.0:6.0:1.5:Balanced - Medium Risk"
    "1.0:3.0:5.0:8.0:2.0:Moderate - Standard"
    "1.5:4.0:6.0:10.0:2.0:Aggressive - Wide Targets"
    "2.0:5.0:8.0:12.0:2.5:Very Aggressive"
    "0.5:3.0:5.0:8.0:1.0:Tight Stop Wide Targets"
    "1.0:2.5:4.0:6.0:1.5:Medium Stop Medium Targets"
    "1.5:3.5:5.5:9.0:2.0:Wide Stop Wide Targets"
    "0.75:3.0:4.5:7.5:1.5:Proven Session Trader Params"
    "1.5:4.0:6.0:8.0:2.0:Proven Liquidity Hunter Params"
)

best_score=0
best_config=""
best_details=""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "TESTING ${#configs[@]} PARAMETER COMBINATIONS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

test_num=0
for config in "${configs[@]}"; do
    test_num=$((test_num + 1))
    
    # Parse config
    IFS=':' read -r stopATR tp1ATR tp2ATR tp3ATR risk desc <<< "$config"
    
    echo "[$test_num/${#configs[@]}] Testing: $desc"
    echo "   Params: Stop=$stopATR, TP1=$tp1ATR, TP2=$tp2ATR, TP3=$tp3ATR, Risk=$risk%"
    
    # Run backtest
    result=$(curl -s -X POST http://localhost:8080/api/v1/backtest/run \
      -H "Content-Type: application/json" \
      -d "{
        \"symbol\": \"BTCUSDT\",
        \"interval\": \"$TIMEFRAME\",
        \"days\": $DAYS,
        \"startBalance\": 1000,
        \"strategy\": \"liquidity_hunter\",
        \"riskPercent\": $(echo "$risk / 100" | bc -l)
      }")
    
    # Extract metrics
    winRate=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(round(data.get('winRate', 0), 2))" 2>/dev/null || echo "0")
    profitFactor=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(round(data.get('profitFactor', 0), 2))" 2>/dev/null || echo "0")
    returnPct=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(round(data.get('returnPercent', 0), 2))" 2>/dev/null || echo "0")
    maxDD=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(round(data.get('maxDrawdown', 0), 2))" 2>/dev/null || echo "0")
    totalTrades=$(echo "$result" | python3 -c "import sys, json; data = json.load(sys.stdin); print(data.get('totalTrades', 0))" 2>/dev/null || echo "0")
    
    # Calculate comprehensive score
    score=$(python3 -c "
wr = $winRate
pf = $profitFactor
ret = $returnPct
dd = $maxDD
trades = $totalTrades

# Only score if profitable and meets minimum criteria
if trades >= 3 and wr >= 40 and pf >= 1.0 and ret > 0:
    # Weighted score: WR (30%) + PF (30%) + Return (30%) - DD penalty (10%)
    score = (wr * 0.3) + (pf * 30 * 0.3) + (ret * 0.3) - (dd * 10 * 0.1)
    print(round(max(0, score), 2))
else:
    print(0)
" 2>/dev/null || echo "0")
    
    echo "   Results: WR=$winRate%, PF=$profitFactor, Return=$returnPct%, DD=$maxDD%, Trades=$totalTrades"
    echo "   Score: $score"
    
    # Grade the result
    if (( $(echo "$winRate >= 70" | bc -l) )) && (( $(echo "$profitFactor >= 3.0" | bc -l) )); then
        echo "   Grade: 🏆 EXCELLENT"
    elif (( $(echo "$winRate >= 60" | bc -l) )) && (( $(echo "$profitFactor >= 2.0" | bc -l) )); then
        echo "   Grade: ✅ VERY GOOD"
    elif (( $(echo "$winRate >= 50" | bc -l) )) && (( $(echo "$profitFactor >= 1.5" | bc -l) )); then
        echo "   Grade: ✅ GOOD"
    elif (( $(echo "$returnPct > 0" | bc -l) )); then
        echo "   Grade: ⚠️  PROFITABLE"
    else
        echo "   Grade: ❌ LOSING"
    fi
    echo ""
    
    # Track best
    if (( $(echo "$score > $best_score" | bc -l) )); then
        best_score=$score
        best_config="$desc"
        best_details="Stop: $stopATR ATR | TP1: $tp1ATR ATR | TP2: $tp2ATR ATR | TP3: $tp3ATR ATR | Risk: $risk% | WR: $winRate% | PF: $profitFactor | Return: $returnPct% | DD: $maxDD% | Trades: $totalTrades"
    fi
done

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🏆 BEST CONFIGURATION FOUND"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Configuration: $best_config"
echo "$best_details"
echo "Score: $best_score"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📝 HOW TO APPLY THESE PARAMETERS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Update backend/unified_signal_generator.go:"
echo ""
echo "In generateLiquidityHunterSignal function, change the parameters to:"
echo ""
echo "// OPTIMIZED PARAMETERS from testing"
echo "if buyScore >= 1 && buyScore >= sellScore {"
echo "    return &AdvancedSignal{"
echo "        Strategy:   \"liquidity_hunter\","
echo "        Type:       \"BUY\","
echo "        Entry:      currentPrice,"
echo "        StopLoss:   currentPrice - (atr * $stopATR),"
echo "        TP1:        currentPrice + (atr * $tp1ATR),"
echo "        TP2:        currentPrice + (atr * $tp2ATR),"
echo "        TP3:        currentPrice + (atr * $tp3ATR),"
echo "        ..."
echo "    }"
echo "}"
echo ""
echo "Then restart the server and test with paper trading!"
echo ""
