# 🎯 CONNECT WITH TRADINGVIEW - SIMPLE GUIDE

## ✅ 3 EASY STEPS TO CONNECT

---

## STEP 1: GET YOUR SIGNAL FROM API

Open terminal and run:
```bash
curl http://localhost:8080/api/v1/backtest/live-signal \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "interval": "15m",
    "strategy": "session_trader"
  }' | jq '.'
```

**You will see:**
```json
{
  "signal": "BUY",
  "entry": 91420.50,
  "stopLoss": 91100.00,
  "tp1": 91700.00,
  "tp2": 92000.00,
  "tp3": 92500.00
}
```

---

## STEP 2: OPEN TRADINGVIEW

### A. Go to TradingView Paper Trading
1. Open: https://www.tradingview.com/chart/
2. Click **"Trading Panel"** at bottom
3. Select **"Paper Trading"** account
4. Search for **"BTCUSDT"** (or your symbol)

### B. Set Your Chart
1. Click **"15m"** timeframe (bottom toolbar)
2. Make sure you see **BTC/USDT** chart

---

## STEP 3: PLACE THE TRADE

### If Signal is BUY:

1. **Click "BUY" button** (green button in trading panel)

2. **Enter these values:**
   - **Price**: `91420.50` (from API entry)
   - **Amount**: Calculate based on your balance
   - **Stop Loss**: `91100.00` (from API stopLoss)
   - **Take Profit**: `92500.00` (from API tp3)

3. **Click "Place Order"**

### If Signal is SELL:

1. **Click "SELL" button** (red button in trading panel)

2. **Enter these values:**
   - **Price**: Entry from API
   - **Amount**: Calculate based on your balance
   - **Stop Loss**: From API
   - **Take Profit**: From API tp3

3. **Click "Place Order"**

---

## 📊 VISUAL EXAMPLE

```
┌─────────────────────────────────────────┐
│  TradingView Chart - BTCUSDT 15m        │
├─────────────────────────────────────────┤
│                                         │
│         📈 Price Chart Here             │
│                                         │
├─────────────────────────────────────────┤
│  Trading Panel (Paper Trading)          │
│                                         │
│  Symbol: BTCUSDT                        │
│  ┌─────────┐  ┌─────────┐              │
│  │   BUY   │  │  SELL   │              │
│  └─────────┘  └─────────┘              │
│                                         │
│  Entry Price:    91420.50               │
│  Stop Loss:      91100.00               │
│  Take Profit:    92500.00               │
│  Amount:         0.001 BTC              │
│                                         │
│  [Place Order]                          │
└─────────────────────────────────────────┘
```

---

## 💰 CALCULATE POSITION SIZE

### For $15 Balance with 0.3% Risk:

**Risk Amount** = $15 × 0.003 = $0.045

**Example:**
- Entry: $91,420.50
- Stop Loss: $91,100.00
- Risk per BTC: $320.50
- Position Size: $0.045 ÷ $320.50 = 0.00014 BTC

**In TradingView:**
- Enter Amount: `0.00014` BTC
- Or in USD: `$13` (use ~87% of balance)

---

## 🔄 COMPLETE WORKFLOW

### Every 15 Minutes:

1. **Check for new signal:**
   ```bash
   curl -X POST http://localhost:8080/api/v1/backtest/live-signal \
     -H "Content-Type: application/json" \
     -d '{"symbol":"BTCUSDT","interval":"15m","strategy":"session_trader"}' \
     | jq '{signal, entry, stopLoss, tp1, tp2, tp3}'
   ```

2. **If signal is BUY or SELL:**
   - Go to TradingView
   - Place the trade with values from API
   - Set Stop Loss and Take Profit

3. **If signal is NONE:**
   - Wait for next check (15 minutes)

---

## 🎯 QUICK REFERENCE

### Get Signal:
```bash
curl -X POST http://localhost:8080/api/v1/backtest/live-signal \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTCUSDT","interval":"15m","strategy":"session_trader"}' | jq '.'
```

### TradingView Paper Trading:
- URL: https://www.tradingview.com/chart/
- Account: Paper Trading
- Timeframe: 15m
- Symbol: BTCUSDT

### Position Size (for $15):
- Risk: 0.3% = $0.045
- Use ~87% of balance per trade
- Always set Stop Loss and Take Profit

---

## 📱 MOBILE APP (Optional)

### TradingView Mobile:
1. Download TradingView app
2. Login to your account
3. Enable Paper Trading
4. Get signals from API on computer
5. Place trades on mobile app

---

## ⚠️ IMPORTANT NOTES

1. **Always use Paper Trading first** - Don't use real money yet
2. **Check signal every 15 minutes** - Or use auto paper trading API
3. **Set Stop Loss** - Never trade without stop loss
4. **BUY ONLY mode** - Only take BUY signals in bull market
5. **Test for 2 weeks** - Verify 75%+ win rate before going live

---

## 🚀 AUTOMATED OPTION

Instead of manual TradingView, use our **Paper Trading API**:

```bash
# Start auto paper trading
curl -X POST http://localhost:8080/api/v1/paper-trading/start-auto

# Check results
curl http://localhost:8080/api/v1/paper-trading/stats | jq '.stats'
```

This automatically:
- Checks for signals every 15 minutes
- Tracks all trades
- Calculates win rate and profit
- No manual work needed!

---

## ✅ READY TO START!

**Option 1: Manual TradingView**
- Check signal every 15 min
- Place trades manually
- Good for learning

**Option 2: Auto Paper Trading API**
- Fully automated
- Just check stats daily
- Recommended for testing

Choose what works best for you!

---

## 📞 NEED HELP?

**Test your connection:**
```bash
# Get current signal
curl -X POST http://localhost:8080/api/v1/backtest/live-signal \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTCUSDT","interval":"15m","strategy":"session_trader"}' | jq '.'
```

**If you see a signal** → Go to TradingView and place the trade!
**If you see "NONE"** → Wait 15 minutes and check again!
