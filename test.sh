#!/bin/bash

# TradingBot Test Script
# Quick testing for all strategies

echo "🧪 TradingBot Test Suite"
echo "========================"
echo ""

# Check if backend is running
if ! curl -s http://localhost:8080/api/v1/health > /dev/null 2>&1; then
    echo "❌ Backend is not running!"
    echo "Start it with: cd backend && go run ."
    exit 1
fi

echo "✅ Backend is running"
echo ""

# Test parameters
SYMBOL="${1:-BTCUSDT}"
INTERVAL="${2:-15m}"
DAYS="${3:-30}"
STRATEGY="${4:-session_trader}"
MODE="${5:-standard}"

echo "📊 Test Parameters:"
echo "   Symbol: $SYMBOL"
echo "   Interval: $INTERVAL"
echo "   Days: $DAYS"
echo "   Strategy: $STRATEGY"
echo "   Mode: $MODE"
echo ""

# Run backtest based on mode
if [ "$MODE" = "world-class" ]; then
    echo "🌟 Running WORLD-CLASS backtest (with advanced metrics)..."
    echo ""
    
    curl -s -X POST "http://localhost:8080/api/v1/backtest/world-class" \
      -H "Content-Type: application/json" \
      -d "{\"symbol\":\"$SYMBOL\",\"interval\":\"$INTERVAL\",\"days\":$DAYS,\"strategy\":\"$STRATEGY\",\"startBalance\":1000,\"enableMonteCarlo\":true,\"monteCarloRuns\":1000,\"enableWalkForward\":true,\"walkForwardPeriods\":5,\"enableStressTest\":true}" \
      | jq '{
        strategy,
        totalTrades,
        winRate,
        profitFactor,
        maxDrawdown,
        finalBalance,
        returnPercent: ((.finalBalance - 1000) / 1000 * 100),
        sharpeRatio,
        sortinoRatio,
        calmarRatio,
        recoveryFactor,
        expectancyPerTrade,
        winStreakMax,
        lossStreakMax,
        monteCarloResults: {
          probabilityProfit,
          expectedReturn,
          worstCase,
          bestCase,
          confidence95Low: .percentile5,
          confidence95High: .percentile95
        },
        walkForwardAnalysis: {
          consistency,
          overfittingScore,
          inSampleWinRate,
          outOfSampleWinRate
        }
      }'
      
elif [ "$MODE" = "compare" ]; then
    echo "⚖️  Running COMPARISON (standard vs world-class)..."
    echo ""
    
    curl -s -X POST "http://localhost:8080/api/v1/backtest/compare" \
      -H "Content-Type: application/json" \
      -d "{\"symbol\":\"$SYMBOL\",\"interval\":\"$INTERVAL\",\"days\":$DAYS,\"strategy\":\"$STRATEGY\",\"startBalance\":1000}" \
      | jq
      
else
    echo "🚀 Running STANDARD backtest..."
    echo ""
    
    curl -s -X POST "http://localhost:8080/api/v1/backtest/run" \
      -H "Content-Type: application/json" \
      -d "{\"symbol\":\"$SYMBOL\",\"interval\":\"$INTERVAL\",\"days\":$DAYS,\"strategy\":\"$STRATEGY\",\"startBalance\":1000}" \
      | jq '{
        strategy,
        totalTrades,
        winningTrades,
        losingTrades,
        winRate,
        profitFactor,
        maxDrawdown,
        finalBalance,
        returnPercent: ((.finalBalance - 1000) / 1000 * 100)
      }'
fi

echo ""
echo "✅ Test complete!"
echo ""
echo "💡 Usage: ./test.sh [SYMBOL] [INTERVAL] [DAYS] [STRATEGY] [MODE]"
echo "   Example: ./test.sh BTCUSDT 15m 30 session_trader standard"
echo "   Example: ./test.sh BTCUSDT 15m 30 session_trader world-class"
echo ""
echo "📚 Available strategies:"
echo "   - session_trader (recommended)"
echo "   - liquidity_hunter"
echo "   - breakout_master"
echo "   - trend_rider"
echo "   - range_master"
echo "   - smart_money_tracker"
echo "   - institutional_follower"
echo "   - reversal_sniper"
echo "   - momentum_beast"
echo "   - scalper_pro"
echo ""
echo "🌟 Test Modes:"
echo "   - standard: Basic backtest (fast)"
echo "   - world-class: Advanced metrics + Monte Carlo + Walk Forward (slower)"
echo "   - compare: Compare both modes side-by-side"
