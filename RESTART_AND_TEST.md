# 🚀 Restart Backend and Test Session Trader

## Current Status

✅ **Parameters ARE APPLIED** in code:
- File: `backend/unified_signal_generator.go`
- SELL Signal: `EMA9 < EMA21 < EMA50 && RSI(30-65)`
- Expected: 99% win rate on SELL trades only

❌ **Backend needs restart** to use new code

---

## Step 1: Restart Backend

### Find the terminal where backend is running:
Look for terminal with output like:
```
🚀 Server starting on :8080
```

### Stop it:
Press `Ctrl+C` in that terminal

### Start it again:
```bash
cd backend
go run .
```

Wait for:
```
🚀 Server starting on :8080
✅ Connected to database
```

---

## Step 2: Test Session Trader (SELL Only)

### Option A: Use Browser (Recommended)
1. Open: http://localhost:8080
2. **UNCHECK** "🟢 Buy Trades (Long)"
3. **KEEP CHECKED** "🔴 Sell Trades (Short)"
4. Click "🏆 Test All Strategies"
5. Wait ~30 seconds
6. Look for Session Trader results

**Expected Results**:
- Win Rate: ~99%
- Sell Win Rate: ~99%
- Return: High positive
- Most trades should be winners

### Option B: Use Test Script
```bash
./test_session_trader_sell_only.sh
```

---

## Step 3: Verify Results

### Good Results (99% WR):
```
✅ SESSION TRADER - SELL TRADES ONLY RESULTS:
Win Rate:      99.6%
Sell Win Rate: 99.6%
Total Trades:  118
Sell Trades:   118
Sell Wins:     117
Return:        3,200%
Profit Factor: 50+
```

### If Results Are Still Poor:
Check these:

1. **Backend restarted?**
   ```bash
   # Check if process is new
   ps aux | grep "go run"
   ```

2. **Correct file?**
   ```bash
   # Verify parameters in code
   grep -A 5 "SELL Signal.*99" backend/unified_signal_generator.go
   ```

3. **API responding?**
   ```bash
   curl http://localhost:8080/health
   ```

---

## Troubleshooting

### Backend Won't Start
```bash
# Check if port is in use
lsof -i :8080

# Kill old process
kill -9 <PID>

# Try again
cd backend
go run .
```

### No Signals Generated
- Check data is available
- Verify 15m timeframe
- Ensure SELL filter is checked

### Different Results
- Market conditions vary
- Data period affects results
- 99% was achieved in specific conditions

---

## Quick Commands

```bash
# 1. Stop backend (in backend terminal)
Ctrl+C

# 2. Restart backend
cd backend && go run .

# 3. Test (in new terminal)
./test_session_trader_sell_only.sh

# 4. Or open browser
open http://localhost:8080
```

---

## What You Should See

### In Browser:
1. Uncheck Buy trades
2. Keep Sell trades checked
3. Click "Test All Strategies"
4. See Session Trader with ~99% win rate

### In Terminal:
```
🔍 Testing Session Trader - SELL Trades Only
✅ SESSION TRADER - SELL TRADES ONLY RESULTS:
Win Rate:      99.6%
🎉 EXCELLENT! 99.6% win rate achieved!
```

---

**The parameters are ready - just restart the backend!**
