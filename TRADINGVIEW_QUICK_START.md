# ⚡ TRADINGVIEW - QUICK START (1 MINUTE)

## 🎯 FASTEST WAY TO START

---

## 1️⃣ GET SIGNAL (10 seconds)

```bash
curl -X POST http://localhost:8080/api/v1/backtest/live-signal \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTCUSDT","interval":"15m","strategy":"session_trader"}' | jq '{signal,entry,stopLoss,tp3}'
```

**Copy these 4 numbers:**
- signal (BUY or SELL)
- entry
- stopLoss  
- tp3

---

## 2️⃣ OPEN TRADINGVIEW (20 seconds)

1. Go to: **https://www.tradingview.com/chart/**
2. Click **"Trading Panel"** (bottom)
3. Select **"Paper Trading"**
4. Search **"BTCUSDT"**
5. Click **"15m"** timeframe

---

## 3️⃣ PLACE TRADE (30 seconds)

1. Click **"BUY"** (or "SELL" if signal says SELL)
2. Enter:
   - **Price**: (entry from API)
   - **Stop Loss**: (stopLoss from API)
   - **Take Profit**: (tp3 from API)
   - **Amount**: 0.00014 BTC
3. Click **"Place Order"**

---

## ✅ DONE!

Trade is now active. It will close automatically at TP or SL.

---

## 🔄 REPEAT

Check for new signal every 15 minutes:
```bash
curl -X POST http://localhost:8080/api/v1/backtest/live-signal \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTCUSDT","interval":"15m","strategy":"session_trader"}' | jq '{signal,entry,stopLoss,tp3}'
```

---

## 📱 EVEN FASTER: AUTO MODE

Don't want to do this manually? Use auto paper trading:

```bash
# Start auto trading
curl -X POST http://localhost:8080/api/v1/paper-trading/start-auto

# Check results (once per day)
curl http://localhost:8080/api/v1/paper-trading/stats | jq '.stats'
```

**This does everything automatically!**

---

## 🎯 THAT'S IT!

You're now connected to TradingView!

**Full guides:**
- `TRADINGVIEW_STEP_BY_STEP.md` - Detailed with pictures
- `CONNECT_TRADINGVIEW_SIMPLE.md` - Complete explanation
- `PAPER_TRADING_READY.md` - Auto trading guide
