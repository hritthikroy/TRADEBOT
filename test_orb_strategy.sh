#!/bin/bash

# Test script for Academic ORB Strategy
# Based on research by Zarattini, Barbon, and Aziz (2024)

echo "=========================================="
echo "Testing Academic ORB Strategy"
echo "=========================================="
echo ""

BASE_URL="http://localhost:8080/api/v1"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test 1: Run 5-minute ORB backtest
echo -e "${YELLOW}Test 1: Running 5-minute ORB Backtest${NC}"
echo "Expected: 1,637% return, Sharpe 2.81, Alpha 36%"
echo ""

curl -X POST "${BASE_URL}/orb/backtest" \
  -H "Content-Type: application/json" \
  -d '{
    "timeFrame": 5,
    "startDate": "2016-01-01",
    "endDate": "2023-12-31",
    "initialCapital": 25000,
    "topNStocks": 20,
    "minRelativeVol": 1.0
  }' | jq '.'

echo ""
echo "=========================================="
echo ""

# Test 2: Compare all timeframes
echo -e "${YELLOW}Test 2: Comparing All Timeframes (5m, 15m, 30m, 60m)${NC}"
echo "Expected: 5m performs best"
echo ""

curl -X POST "${BASE_URL}/orb/compare" \
  -H "Content-Type: application/json" \
  -d '{
    "startDate": "2016-01-01",
    "endDate": "2023-12-31",
    "initialCapital": 25000
  }' | jq '.'

echo ""
echo "=========================================="
echo ""

# Test 3: Get top performers for 5-minute ORB
echo -e "${YELLOW}Test 3: Top Performers (5-minute ORB)${NC}"
echo "Expected: DDD, FSLR, NVDA, SWBI, RCL"
echo ""

curl -X GET "${BASE_URL}/orb/top-performers?timeframe=5" | jq '.'

echo ""
echo "=========================================="
echo ""

# Test 4: Get live signals (if market is open)
echo -e "${YELLOW}Test 4: Getting Live ORB Signals${NC}"
echo ""

curl -X GET "${BASE_URL}/orb/live-signals?timeframe=5" | jq '.'

echo ""
echo "=========================================="
echo ""

# Test 5: Test different timeframes
echo -e "${YELLOW}Test 5: Testing 15-minute ORB${NC}"
echo "Expected: 272% return, Sharpe 1.43"
echo ""

curl -X POST "${BASE_URL}/orb/backtest" \
  -H "Content-Type: application/json" \
  -d '{
    "timeFrame": 15,
    "startDate": "2016-01-01",
    "endDate": "2023-12-31",
    "initialCapital": 25000,
    "topNStocks": 20,
    "minRelativeVol": 1.0
  }' | jq '.summary'

echo ""
echo "=========================================="
echo ""

echo -e "${GREEN}All tests completed!${NC}"
echo ""
echo "Key Findings from Research:"
echo "  • 5-minute ORB: 1,637% return (best)"
echo "  • 15-minute ORB: 272% return"
echo "  • 30-minute ORB: 21% return"
echo "  • 60-minute ORB: 39% return"
echo "  • S&P 500: 198% return"
echo ""
echo "Critical Success Factor:"
echo "  • Relative Volume > 100% is ESSENTIAL"
echo "  • Without it: only 29% return"
echo "  • With it: 1,637% return (56x improvement!)"
echo ""
echo "Access web interface at:"
echo "  http://localhost:8080/orb_academic.html"
echo ""
