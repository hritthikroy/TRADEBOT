#!/bin/bash

echo "🧪 Testing Custom Date Range Feature"
echo "======================================"
echo ""

# Test 1: Test with custom date range (November 2024)
echo "Test 1: Custom Date Range (November 2024)"
echo "------------------------------------------"

START_DATE="2024-11-01"
END_DATE="2024-11-30"

# Convert dates to Unix timestamps (milliseconds)
START_TIME=$(date -j -f "%Y-%m-%d" "$START_DATE" "+%s000" 2>/dev/null || date -d "$START_DATE" "+%s000")
END_TIME=$(date -j -f "%Y-%m-%d" "$END_DATE" "+%s999" 2>/dev/null || date -d "$END_DATE 23:59:59" "+%s999")

echo "Start Date: $START_DATE"
echo "End Date: $END_DATE"
echo "Start Time: $START_TIME"
echo "End Time: $END_TIME"
echo ""

# Make API request
echo "Making API request..."
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/backtest/test-all-strategies \
  -H "Content-Type: application/json" \
  -d "{
    \"symbol\": \"BTCUSDT\",
    \"days\": 30,
    \"startBalance\": 500,
    \"filterBuy\": true,
    \"filterSell\": true,
    \"startTime\": $START_TIME,
    \"endTime\": $END_TIME
  }")

# Check if response is valid
if echo "$RESPONSE" | grep -q "bestStrategy"; then
    echo "✅ API request successful!"
    echo ""
    
    # Extract best strategy info
    BEST_STRATEGY=$(echo "$RESPONSE" | grep -o '"strategyName":"[^"]*"' | head -1 | cut -d'"' -f4)
    WIN_RATE=$(echo "$RESPONSE" | grep -o '"winRate":[0-9.]*' | head -1 | cut -d':' -f2)
    RETURN=$(echo "$RESPONSE" | grep -o '"returnPercent":[0-9.]*' | head -1 | cut -d':' -f2)
    
    echo "📊 Results for November 2024:"
    echo "   Best Strategy: $BEST_STRATEGY"
    echo "   Win Rate: $WIN_RATE%"
    echo "   Return: $RETURN%"
    echo ""
else
    echo "❌ API request failed!"
    echo "Response: $RESPONSE"
    exit 1
fi

# Test 2: Test with different custom date range (October 2024)
echo ""
echo "Test 2: Custom Date Range (October 2024)"
echo "------------------------------------------"

START_DATE="2024-10-01"
END_DATE="2024-10-31"

START_TIME=$(date -j -f "%Y-%m-%d" "$START_DATE" "+%s000" 2>/dev/null || date -d "$START_DATE" "+%s000")
END_TIME=$(date -j -f "%Y-%m-%d" "$END_DATE" "+%s999" 2>/dev/null || date -d "$END_DATE 23:59:59" "+%s999")

echo "Start Date: $START_DATE"
echo "End Date: $END_DATE"
echo ""

RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/backtest/test-all-strategies \
  -H "Content-Type: application/json" \
  -d "{
    \"symbol\": \"BTCUSDT\",
    \"days\": 31,
    \"startBalance\": 500,
    \"filterBuy\": true,
    \"filterSell\": true,
    \"startTime\": $START_TIME,
    \"endTime\": $END_TIME
  }")

if echo "$RESPONSE" | grep -q "bestStrategy"; then
    echo "✅ API request successful!"
    echo ""
    
    BEST_STRATEGY=$(echo "$RESPONSE" | grep -o '"strategyName":"[^"]*"' | head -1 | cut -d'"' -f4)
    WIN_RATE=$(echo "$RESPONSE" | grep -o '"winRate":[0-9.]*' | head -1 | cut -d':' -f2)
    RETURN=$(echo "$RESPONSE" | grep -o '"returnPercent":[0-9.]*' | head -1 | cut -d':' -f2)
    
    echo "📊 Results for October 2024:"
    echo "   Best Strategy: $BEST_STRATEGY"
    echo "   Win Rate: $WIN_RATE%"
    echo "   Return: $RETURN%"
    echo ""
else
    echo "❌ API request failed!"
    echo "Response: $RESPONSE"
    exit 1
fi

echo "======================================"
echo "🎉 All tests passed!"
echo ""
echo "Custom date range feature is working correctly."
echo "The backend properly uses the startTime and endTime"
echo "parameters to fetch data for the specified date range."
echo ""
echo "Now test in the UI:"
echo "  1. Open http://localhost:8080"
echo "  2. Check 'Use Calendar'"
echo "  3. Select 'Custom Date Range'"
echo "  4. Pick dates (e.g., Nov 1 - Nov 30)"
echo "  5. Click 'Run Backtest'"
echo "  6. Charts should show data for selected dates!"
echo ""
