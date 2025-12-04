# 🧪 Test: Days Parameter Fix

## Quick Test Guide

### Test 1: Run Backtest with Different Days

```
Step 1: Open http://localhost:8080

Step 2: Test with 7 days
- Change "Days to Test" to 7
- Click "Run Backtest"
- Note the number of trades

Step 3: Test with 30 days
- Change "Days to Test" to 30
- Click "Run Backtest"
- Number of trades should be HIGHER

Step 4: Test with 90 days
- Change "Days to Test" to 90
- Click "Run Backtest"
- Number of trades should be MUCH HIGHER
```

### Expected Results:

| Days | Expected Trades | Expected Behavior |
|------|----------------|-------------------|
| 1    | Very few (0-5) | Recent data only |
| 7    | Few (10-30)    | Last week |
| 15   | Moderate (20-60) | Last 2 weeks |
| 30   | Normal (50-200) | Last month (default) |
| 90   | Many (150-600) | Last 3 months |
| 180  | Lots (300-1200) | Last 6 months |
| 365  | Maximum (600-2400) | Full year |

---

## Test 2: Test All Strategies with Different Days

```
Step 1: Open http://localhost:8080

Step 2: Test with 7 days
- Change "Days to Test" to 7
- Click "🏆 Test All Strategies"
- Note Session Trader trades

Step 3: Test with 30 days
- Change "Days to Test" to 30
- Click "🏆 Test All Strategies"
- Session Trader should have MORE trades

Step 4: Compare
- 7 days: ~50-70 trades
- 30 days: ~192 trades (as we saw before)
- 90 days: ~500+ trades
```

---

## What to Look For

### ✅ Signs the Fix is Working:

1. **Trade Count Changes**
   - Fewer days = Fewer trades
   - More days = More trades

2. **Win Rate May Change**
   - Different time periods = Different market conditions
   - Win rate may vary slightly

3. **Return Changes**
   - Different number of trades = Different returns
   - More trades usually = Higher returns

4. **Chart Updates**
   - Equity curve should show different time periods
   - Shorter days = Shorter curve
   - Longer days = Longer curve

5. **Trade List Updates**
   - Trade dates should match the time period
   - 7 days = Trades from last 7 days only
   - 30 days = Trades from last 30 days

### ❌ Signs the Fix is NOT Working:

1. **Trade Count Stays Same**
   - Always ~192 trades regardless of days
   - This means days parameter still not working

2. **No Chart Changes**
   - Equity curve looks identical
   - Same time period shown

3. **Same Results Every Time**
   - Win rate never changes
   - Return never changes
   - Exactly same numbers

---

## Quick API Test

### Test via curl:

```bash
# Test with 7 days
curl -X POST http://localhost:8080/api/v1/backtest/test-all-strategies \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "days": 7,
    "startBalance": 1000,
    "filterBuy": false,
    "filterSell": true
  }' | python3 -c "
import sys, json
data = json.load(sys.stdin)
for r in data['results']:
    if r['strategyName'] == 'session_trader':
        print(f\"7 days: {r['totalTrades']} trades, {r['winRate']:.1f}% WR\")
"

# Test with 30 days
curl -X POST http://localhost:8080/api/v1/backtest/test-all-strategies \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "days": 30,
    "startBalance": 1000,
    "filterBuy": false,
    "filterSell": true
  }' | python3 -c "
import sys, json
data = json.load(sys.stdin)
for r in data['results']:
    if r['strategyName'] == 'session_trader':
        print(f\"30 days: {r['totalTrades']} trades, {r['winRate']:.1f}% WR\")
"
```

### Expected Output:
```
7 days: 50-70 trades, 52.6% WR
30 days: 180-200 trades, 52.6% WR
```

If both show the same number of trades, the fix didn't work!

---

## Summary

### What Was Fixed:
- ✅ Added `days` parameter to runBacktest() API call
- ✅ Added `days` parameter to testAllStrategies() API call
- ✅ Export function already had it

### What to Test:
1. Change days in UI
2. Run backtest
3. Verify trade count changes
4. Verify charts update

### Success Criteria:
- ✅ Different days = Different trade counts
- ✅ Charts update with new data
- ✅ Results change based on time period

---

**Status**: ✅ Fix applied, ready to test  
**Action**: Open http://localhost:8080 and try different day values!
