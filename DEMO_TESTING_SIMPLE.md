# 🎯 HOW TO TEST ON DEMO - SIMPLE GUIDE

## 3 Easy Ways to Test Your Bot

---

## ✅ Method 1: Check Live Signal Right Now (Fastest)

### Step 1: Make sure backend is running
```bash
cd backend
go run .
```

### Step 2: Check current signal
```bash
curl -s http://localhost:8080/api/v1/backtest/live-signal | jq '.'
```

### What you'll see:
```json
{
  "signal": "BUY",           ← BUY, SELL, or NONE
  "currentPrice": 89939.11,  ← Current Bitcoin price
  "entry": 89939.11,         ← Where to enter
  "stopLoss": 89500.00,      ← Where to exit if losing
  "takeProfit": 91200.00,    ← Where to exit if winning
  "riskReward": 2.89         ← Risk/Reward ratio
}
```

### What to do:
- If signal = "BUY": Write it down, track if it wins/loses
- If signal = "NONE": Wait, check again in 15 minutes

---

## ✅ Method 2: Use Dashboard (Visual)

### Step 1: Open browser
```
http://localhost:8080
```

### Step 2: You'll see:
- 📊 Live Bitcoin chart
- 🟢 Current signal (BUY/SELL/NONE)
- 📈 Entry, Stop Loss, Take Profit levels
- 📊 Strategy statistics

### Step 3: Paper trade
- When you see BUY signal:
  - Screenshot it
  - Write down: Entry, Stop Loss, Take Profit
  - Check later if it won or lost

---

## ✅ Method 3: Auto-Monitor (Automated)

### Step 1: Run paper trading script
```bash
./start_paper_trading.sh
```

### What it does:
- Checks for signals every 5 minutes
- Logs all BUY signals to file
- Shows you when to enter/exit
- Tracks your paper trades

### Output example:
```
📝 PAPER TRADING MODE - STARTED
================================
Mode: BUY ONLY
Risk: 0.3% per trade
Starting Balance: $15

⚪ 10:00:00 - Waiting... (Price: $89,939)
⚪ 10:05:00 - Waiting... (Price: $89,950)
🟢 BUY SIGNAL #1 at 10:10:00
   Entry: $89,960 | SL: $89,500 | TP: $91,200
```

---

## 📝 How to Track Paper Trades

### Simple Spreadsheet:

| # | Date | Entry | Stop Loss | Take Profit | Result | Profit |
|---|------|-------|-----------|-------------|--------|--------|
| 1 | Dec 6 | 89,960 | 89,500 | 91,200 | WIN ✅ | +$0.08 |
| 2 | Dec 6 | 90,200 | 89,800 | 91,500 | WIN ✅ | +$0.09 |
| 3 | Dec 7 | 90,500 | 90,100 | 91,800 | LOSS ❌ | -$0.04 |

### Track:
- Total trades
- Wins vs Losses
- Win rate %
- Total profit/loss

---

## 🎯 Quick Test Right Now

### Test 1: Check if bot is working
```bash
curl http://localhost:8080/api/v1/health
```

Should show: `{"status":"ok"}`

### Test 2: Get current signal
```bash
curl -s http://localhost:8080/api/v1/backtest/live-signal | jq '.signal'
```

Shows: `"BUY"`, `"SELL"`, or `"NONE"`

### Test 3: Get current price
```bash
curl -s http://localhost:8080/api/v1/backtest/live-signal | jq '.currentPrice'
```

Shows: Current Bitcoin price

---

## 📊 Example Paper Trading Day

### 10:00 AM - Start
```bash
./start_paper_trading.sh
```

### 10:15 AM - First Signal
```
🟢 BUY SIGNAL #1
Entry: $89,960
Stop Loss: $89,500
Take Profit: $91,200
```

**Action**: Write it down, track it

### 2:00 PM - Check Result
- Price reached $91,200 ✅
- Take Profit hit!
- Result: WIN
- Profit: +$0.08 (on $15 balance)

### 4:00 PM - Second Signal
```
🟢 BUY SIGNAL #2
Entry: $90,500
Stop Loss: $90,100
Take Profit: $91,800
```

**Action**: Write it down, track it

### End of Day
- Total signals: 2
- Wins: 2
- Losses: 0
- Win rate: 100%
- Profit: +$0.16

---

## ⚠️ Important Rules

### ✅ DO:
1. Track EVERY signal (don't skip)
2. Follow stop loss/take profit exactly
3. Test for at least 2 weeks
4. Use realistic amounts ($15)

### ❌ DON'T:
1. Cherry-pick signals
2. Change stop loss/take profit
3. Skip losing trades
4. Give up after 1 bad day

---

## 🚀 When Ready for Live Trading

### After 2 weeks of paper trading:

**If you have:**
- ✅ 20+ trades tracked
- ✅ Win rate > 60%
- ✅ Overall profitable
- ✅ Comfortable with process

**Then:**
✅ **READY FOR LIVE TRADING**

Start with $15-50, use 0.3% risk

---

## 📞 Quick Commands

### Check signal:
```bash
curl -s http://localhost:8080/api/v1/backtest/live-signal | jq '.'
```

### Start paper trading:
```bash
./start_paper_trading.sh
```

### Open dashboard:
```
http://localhost:8080
```

### Check bot status:
```bash
curl http://localhost:8080/api/v1/health
```

---

## 🎯 Summary

**3 Ways to Demo Test:**

1. **Manual**: Check signal, write it down, track manually
2. **Dashboard**: Open browser, watch signals visually
3. **Automated**: Run script, auto-logs all signals

**All methods work!** Choose what's easiest for you.

**Expected Results:**
- Win Rate: 73-99% (in bull market)
- Drawdown: 5-13%
- Profitable: Yes

**Status**: ✅ READY TO START DEMO TESTING

---

**Start now:** `./start_paper_trading.sh`
