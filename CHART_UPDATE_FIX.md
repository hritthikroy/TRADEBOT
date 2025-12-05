# Chart Not Updating Fix

## Problem
Charts were showing the same data (95 trades) regardless of date range changes. The issue was caused by:

1. **Browser caching** - API responses were being cached
2. **Duplicate chart creation** - `createEquityChart()` was called twice in `displayResults()`
3. **Stale data** - `currentResults` variable wasn't cleared before new requests
4. **Chart not properly destroyed** - Previous chart instance wasn't fully cleaned up

## Solution Applied

### 1. Added Cache Busting
```javascript
// Added timestamp to URL and cache-control headers
const cacheBuster = Date.now();
const response = await fetch(`${API_URL}/backtest/test-all-strategies?_=${cacheBuster}`, {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
        'Cache-Control': 'no-cache, no-store, must-revalidate',
        'Pragma': 'no-cache'
    },
    body: JSON.stringify(requestBody)
});
```

### 2. Clear Previous Results
```javascript
async function runBacktest() {
    // Clear previous results immediately to prevent stale data
    currentResults = null;
    // ... rest of function
}
```

### 3. Force Chart Destruction
```javascript
// Force chart destruction before displaying new results
if (equityChart) {
    console.log('🗑️ Destroying previous chart');
    equityChart.destroy();
    equityChart = null;
}
```

### 4. Removed Duplicate Chart Creation
```javascript
// BEFORE (lines 1544-1545):
createEquityChart(results);
createEquityChart(results);  // ❌ Duplicate!

// AFTER:
createEquityChart(results);  // ✅ Single call
```

### 5. Added Debug Logging
```javascript
console.log(`🔄 Starting backtest request at ${new Date(cacheBuster).toLocaleTimeString()}`);
console.log('Request parameters:', requestBody);
console.log(`✅ Received response with ${data.results?.length || 0} strategies`);
console.log(`📊 ${strategy} results:`, {
    totalTrades: selectedStrategy.totalTrades,
    winRate: selectedStrategy.winRate,
    returnPercent: selectedStrategy.returnPercent
});
```

## How to Test

1. Open browser console (F12)
2. Run a backtest with 15 days
3. Note the trade count in console and chart
4. Change to 30 days and run again
5. Verify:
   - Console shows new request timestamp
   - Console shows different trade count
   - Chart updates with new data
   - No "95 trades" stuck issue

## Expected Behavior

- **5 days**: ~185 trades for Session Trader (15m timeframe)
- **15 days**: ~427 trades for Session Trader (hits 1000 candle limit)
- **30 days**: ~427 trades for Session Trader (also hits 1000 candle limit - same as 15 days)
- **Custom dates**: Should show trades only within selected date range
- **Chart**: Should update immediately with new data, no caching

### Important: Binance API Limit

The backend fetches data from Binance API which has a **1000 candle limit per request**. This means:

- **15m timeframe**: 
  - 5 days = 480 candles ✅ Works fine
  - 10 days = 960 candles ✅ Works fine
  - 15 days = 1440 candles ⚠️ Capped at 1000 (only gets ~10.4 days)
  - 30 days = 2880 candles ⚠️ Capped at 1000 (only gets ~10.4 days)

- **5m timeframe**:
  - 3 days = 864 candles ✅ Works fine
  - 5 days = 1440 candles ⚠️ Capped at 1000 (only gets ~3.5 days)

- **4h timeframe**:
  - 100 days = 600 candles ✅ Works fine
  - 150 days = 900 candles ✅ Works fine
  - 200 days = 1200 candles ⚠️ Capped at 1000 (only gets ~166 days)

**This is why 15 days and 30 days show the same results** - both exceed the limit and get capped at 1000 candles.

**Solution**: Use custom date ranges with the calendar feature to test specific historical periods without hitting the limit.

## Files Modified

- `public/index.html`:
  - `runBacktest()` function - added cache busting and result clearing
  - `testAllStrategies()` function - added cache busting and result clearing
  - `displayResults()` function - removed duplicate chart creation
  - Added debug logging throughout

## Technical Details

### Cache Busting
The `?_=${cacheBuster}` query parameter forces the browser to treat each request as unique, preventing cached responses.

### Chart Lifecycle
1. Clear `currentResults` → 2. Destroy old chart → 3. Fetch new data → 4. Create new chart

This ensures no stale data or chart instances interfere with new results.
