# 🎯 PAPER TRADING API - Complete Guide

## ✅ NEW! Automated Paper Trading API

I've implemented a complete Paper Trading API that automatically tracks and tests your trades!

---

## 🚀 API Endpoints

### 1. Get Paper Trading Statistics
```bash
curl http://localhost:8080/api/v1/paper-trading/stats | jq '.'
```

**Response:**
```json
{
  "success": true,
  "stats": {
    "totalTrades": 15,
    "openTrades": 2,
    "closedTrades": 13,
    "winningTrades": 11,
    "losingTrades": 2,
    "winRate": 84.6,
    "totalProfit": 2.45,
    "totalLoss": 0.32,
    "netProfit": 2.13,
    "profitFactor": 7.66,
    "startBalance": 15.00,
    "currentBalance": 17.13,
    "returnPercent": 14.2,
    "maxDrawdown": 2.1,
    "averageWin": 0.22,
    "averageLoss": 0.16
  }
}
```

---

### 2. Get All Paper Trades
```bash
curl http://localhost:8080/api/v1/paper-trading/trades | jq '.'
```

**Response:**
```json
{
  "success": true,
  "trades": [
    {
      "id": 1,
      "signal": "BUY",
      "entry": 89939.11,
      "stopLoss": 89500.00,
      "takeProfit": 91200.00,
      "tp1": 90400.00,
      "tp2": 90800.00,
      "tp3": 91200.00,
      "entryTime": "2025-12-06T10:15:00Z",
      "exitTime": "2025-12-06T14:30:00Z",
      "exitPrice": 90800.00,
      "exitReason": "TP2",
      "profit": 0.09,
      "profitPercent": 0.6,
      "status": "won",
      "riskAmount": 0.045
    }
  ],
  "total": 15
}
```

---

### 3. Add New Paper Trade (Manual)
```bash
curl -X POST http://localhost:8080/api/v1/paper-trading/trade \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTCUSDT"}' | jq '.'
```

**What it does:**
- Gets current live signal
- If BUY signal exists, adds it as paper trade
- Tracks entry, stop loss, take profits

---

### 4. Update Open Trades
```bash
curl -X POST http://localhost:8080/api/v1/paper-trading/update \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTCUSDT"}' | jq '.'
```

**What it does:**
- Checks current Bitcoin price
- Updates all open trades
- Closes trades that hit TP or SL
- Returns list of closed trades

**Response:**
```json
{
  "success": true,
  "currentPrice": 90800.00,
  "closedTrades": [
    {
      "id": 1,
      "exitReason": "TP2",
      "profit": 0.09,
      "status": "won"
    }
  ],
  "message": "Trades updated"
}
```

---

### 5. Start Auto Paper Trading
```bash
curl -X POST http://localhost:8080/api/v1/paper-trading/start-auto | jq '.'
```

**What it does:**
- Automatically checks for signals every 15 minutes
- Adds new trades when signals appear
- Updates open trades automatically
- Runs in background

**Response:**
```json
{
  "success": true,
  "message": "Auto paper trading started"
}
```

---

### 6. Stop Auto Paper Trading
```bash
curl -X POST http://localhost:8080/api/v1/paper-trading/stop-auto | jq '.'
```

---

### 7. Reset Paper Trading
```bash
curl -X POST http://localhost:8080/api/v1/paper-trading/reset | jq '.'
```

**Warning:** This deletes all paper trades and resets balance to $15!

---

## 🎯 Quick Start: Auto Paper Trading

### Step 1: Start Backend
```bash
cd backend
go run .
```

### Step 2: Start Auto Paper Trading
```bash
curl -X POST http://localhost:8080/api/v1/paper-trading/start-auto
```

### Step 3: Check Stats Anytime
```bash
curl http://localhost:8080/api/v1/paper-trading/stats | jq '.stats'
```

### Step 4: View All Trades
```bash
curl http://localhost:8080/api/v1/paper-trading/trades | jq '.trades'
```

---

## 📊 Example Usage

### Complete Paper Trading Session:

```bash
# 1. Start auto paper trading
curl -X POST http://localhost:8080/api/v1/paper-trading/start-auto

# 2. Wait 15 minutes...

# 3. Check stats
curl http://localhost:8080/api/v1/paper-trading/stats | jq '.stats | {
  trades: .totalTrades,
  wins: .winningTrades,
  losses: .losingTrades,
  winRate: .winRate,
  profit: .netProfit,
  balance: .currentBalance
}'

# Output:
{
  "trades": 5,
  "wins": 4,
  "losses": 1,
  "winRate": 80,
  "profit": 0.32,
  "balance": 15.32
}

# 4. View trade details
curl http://localhost:8080/api/v1/paper-trading/trades | jq '.trades[] | {
  id: .id,
  signal: .signal,
  entry: .entry,
  exit: .exitPrice,
  result: .status,
  profit: .profit
}'
```

---

## 🔄 How It Works

### Auto Paper Trading Flow:

1. **Every 15 minutes:**
   - Checks for new BUY signals
   - If signal exists, creates paper trade
   - Updates all open trades with current price

2. **Trade Tracking:**
   - Monitors if price hits TP1, TP2, TP3, or Stop Loss
   - Automatically closes trades
   - Calculates profit/loss
   - Updates balance

3. **Statistics:**
   - Tracks win rate, profit factor, drawdown
   - Calculates returns
   - Shows average win/loss

---

## 📝 Data Storage

All paper trades are saved to: `backend/paper_trades.json`

**Format:**
```json
{
  "trades": [...],
  "startBalance": 15.00,
  "currentBalance": 17.13
}
```

**Persists across restarts!**

---

## 🎯 Testing Script

I've created a test script for you:

```bash
./test_paper_trading_api.sh
```

This will:
- Test all API endpoints
- Add sample trades
- Show statistics
- Verify everything works

---

## 📊 Dashboard Integration

The paper trading data can be displayed on your dashboard at:
`http://localhost:8080`

Shows:
- Current balance
- Win rate
- Open trades
- Closed trades
- Profit/loss chart

---

## ⚠️ Important Notes

### Auto Paper Trading:
- ✅ Runs in background
- ✅ Checks every 15 minutes
- ✅ Automatically tracks trades
- ✅ Persists data to file
- ⚠️ Stops when backend restarts (restart with API call)

### Manual Paper Trading:
- ✅ Full control
- ✅ Add trades manually
- ✅ Update when you want
- ✅ Good for testing

---

## 🚀 Recommended Workflow

### Week 1-2: Auto Paper Trading
```bash
# Start auto trading
curl -X POST http://localhost:8080/api/v1/paper-trading/start-auto

# Check daily
curl http://localhost:8080/api/v1/paper-trading/stats | jq '.stats'
```

### After 2 Weeks:
```bash
# Review results
curl http://localhost:8080/api/v1/paper-trading/stats | jq '.stats | {
  totalTrades,
  winRate,
  profitFactor,
  returnPercent,
  maxDrawdown
}'

# If good results (>60% WR, profitable):
# → Ready for live trading!
```

---

## 📞 Quick Commands

### Start auto trading:
```bash
curl -X POST http://localhost:8080/api/v1/paper-trading/start-auto
```

### Check stats:
```bash
curl http://localhost:8080/api/v1/paper-trading/stats | jq '.stats'
```

### View trades:
```bash
curl http://localhost:8080/api/v1/paper-trading/trades | jq '.trades'
```

### Reset (start over):
```bash
curl -X POST http://localhost:8080/api/v1/paper-trading/reset
```

---

## ✅ Status

**Paper Trading API**: ✅ IMPLEMENTED
**Auto Trading**: ✅ WORKING
**Data Persistence**: ✅ WORKING
**Statistics**: ✅ WORKING

**Ready to use!** 🚀

---

**Start now:**
```bash
curl -X POST http://localhost:8080/api/v1/paper-trading/start-auto
```
