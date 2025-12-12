#!/bin/bash

# Academic ORB Strategy - API Examples
# Demonstrates various ways to use the ORB API endpoints

echo "╔══════════════════════════════════════════════════════════════╗"
echo "║         Academic ORB Strategy - API Examples                 ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo ""

BASE_URL="http://localhost:8080/api/v1/orb"

# Example 1: Basic 5-minute backtest
echo "📊 Example 1: Basic 5-Minute Backtest"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
curl -s -X POST "${BASE_URL}/backtest" \
  -H "Content-Type: application/json" \
  -d '{
    "timeFrame": 5,
    "startDate": "2020-01-01",
    "endDate": "2023-12-31",
    "initialCapital": 10000
  }' | jq '.summary'
echo ""
echo ""

# Example 2: 15-minute backtest with custom parameters
echo "📊 Example 2: 15-Minute Backtest (Custom Parameters)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
curl -s -X POST "${BASE_URL}/backtest" \
  -H "Content-Type: application/json" \
  -d '{
    "timeFrame": 15,
    "startDate": "2021-01-01",
    "endDate": "2023-12-31",
    "initialCapital": 50000,
    "topNStocks": 10,
    "minRelativeVol": 1.5
  }' | jq '.summary'
echo ""
echo ""

# Example 3: Compare all timeframes
echo "📊 Example 3: Compare All Timeframes"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
curl -s -X POST "${BASE_URL}/compare" \
  -H "Content-Type: application/json" \
  -d '{
    "startDate": "2022-01-01",
    "endDate": "2023-12-31",
    "initialCapital": 25000
  }' | jq '.results | to_entries | .[] | {timeframe: .key, return: .value.totalReturn, sharpe: .value.sharpeRatio}'
echo ""
echo ""

# Example 4: Get top performers for different timeframes
echo "📊 Example 4: Top Performers (5-minute)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
curl -s "${BASE_URL}/top-performers?timeframe=5" | jq '.topPerformers[:3]'
echo ""

echo "📊 Example 5: Top Performers (15-minute)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
curl -s "${BASE_URL}/top-performers?timeframe=15" | jq '.topPerformers[:3]'
echo ""
echo ""

# Example 6: Get live signals
echo "📊 Example 6: Live Signals (5-minute)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
curl -s "${BASE_URL}/live-signals?timeframe=5" | jq '.'
echo ""
echo ""

# Example 7: Recent period backtest
echo "📊 Example 7: Recent Period (Last 6 Months)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
curl -s -X POST "${BASE_URL}/backtest" \
  -H "Content-Type: application/json" \
  -d '{
    "timeFrame": 5,
    "startDate": "2023-06-01",
    "endDate": "2023-12-31",
    "initialCapital": 25000
  }' | jq '{
    strategy: .summary.strategy,
    period: .summary.period,
    return: .summary.totalReturn,
    sharpe: .summary.sharpeRatio,
    winRate: .summary.winRate
  }'
echo ""
echo ""

echo "╔══════════════════════════════════════════════════════════════╗"
echo "║                    Examples Complete!                        ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo ""
echo "💡 Tips:"
echo "  • Use timeFrame: 5 for best results (1,637% return)"
echo "  • Set minRelativeVol > 1.0 to focus on Stocks in Play"
echo "  • Adjust topNStocks (1-50) to trade more/fewer stocks"
echo "  • Compare timeframes to see performance differences"
echo ""
