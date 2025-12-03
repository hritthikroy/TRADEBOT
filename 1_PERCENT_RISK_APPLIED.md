# ✅ 1% Risk Applied to All Strategies

## 🎯 What Changed

All trading strategies now use **1% risk per trade** to match your backtest settings.

## 📊 Before vs After

### Before (Example):
```
Entry: $91,657.88
Stop Loss: $88,288.31 (3.67% risk) ❌
Take Profit: $100,081.81
Risk: $3,369.57
```

### After (Same Signal):
```
Entry: $91,657.88
Stop Loss: $90,741.30 (1% risk) ✅
Take Profit: $93,574.46
Risk: $916.58
```

## 🔧 Changes Made

### All Strategies Updated:
1. **Session Trader** - 1% stop, 2.5:1 R:R
2. **Breakout Master** - 1% stop, 3:1 R:R
3. **Liquidity Hunter** - 1% stop, 3:1 R:R
4. **Trend Rider** - 1% stop, 3:1 R:R
5. **Range Master** - 1% stop, dynamic target
6. **Scalper Pro** - 1% stop, 2:1 R:R

### Stop Loss Formula:
**BUY Signals:**
```go
StopLoss = CurrentPrice * 0.99  // 1% below entry
```

**SELL Signals:**
```go
StopLoss = CurrentPrice * 1.01  // 1% above entry
```

## 💰 Risk Calculation Examples

### Example 1: BTC at $90,000
**BUY Signal:**
- Entry: $90,000
- Stop Loss: $89,100 (1% risk)
- Risk: $900

**Position Sizing:**
- $10,000 account, 1% risk = $100 risk
- Position size: $100 / $900 × $90,000 = $10,000
- Or: 0.111 BTC

### Example 2: BTC at $50,000
**BUY Signal:**
- Entry: $50,000
- Stop Loss: $49,500 (1% risk)
- Risk: $500

**Position Sizing:**
- $10,000 account, 1% risk = $100 risk
- Position size: $100 / $500 × $50,000 = $10,000
- Or: 0.2 BTC

## 🎯 Risk/Reward Ratios

With 1% risk, here are the typical R:R ratios:

| Strategy | Risk | Reward | Ratio |
|----------|------|--------|-------|
| Session Trader | 1% | 2.5% | 2.5:1 |
| Breakout Master | 1% | 3% | 3:1 |
| Liquidity Hunter | 1% | 3% | 3:1 |
| Trend Rider | 1% | 3% | 3:1 |
| Range Master | 1% | Variable | 2-3:1 |
| Scalper Pro | 1% | 2% | 2:1 |

## 📈 Expected Results

### With 1% Risk Per Trade:

**Scenario 1: 50% Win Rate, 2.5:1 R:R**
- 10 trades: 5 wins, 5 losses
- Wins: 5 × 2.5% = +12.5%
- Losses: 5 × 1% = -5%
- Net: +7.5% profit ✅

**Scenario 2: 40% Win Rate, 3:1 R:R**
- 10 trades: 4 wins, 6 losses
- Wins: 4 × 3% = +12%
- Losses: 6 × 1% = -6%
- Net: +6% profit ✅

**Scenario 3: 60% Win Rate, 2:1 R:R**
- 10 trades: 6 wins, 4 losses
- Wins: 6 × 2% = +12%
- Losses: 4 × 1% = -4%
- Net: +8% profit ✅

## 💡 Position Sizing Guide

### For 1% Account Risk:

| Account Size | Max Risk Per Trade | Example Position |
|--------------|-------------------|------------------|
| $1,000 | $10 | $1,000 |
| $5,000 | $50 | $5,000 |
| $10,000 | $100 | $10,000 |
| $50,000 | $500 | $50,000 |
| $100,000 | $1,000 | $100,000 |

**Formula:**
```
Position Size = (Account Size × Risk %) / Stop Loss %
Position Size = (Account Size × 1%) / 1%
Position Size = Account Size
```

With 1% risk and 1% stop loss, you can use your full account size!

## 🎯 Benefits of 1% Risk

1. **Consistent with Backtest** - Matches your historical testing
2. **Sustainable** - Can survive 100 consecutive losses
3. **Psychological** - Easier to handle losses
4. **Scalable** - Works for any account size
5. **Professional** - Industry standard for risk management

## 🚀 How to Test

### Step 1: Restart Backend
```bash
cd backend
go run .
```

### Step 2: Generate Signal
Wait for Telegram bot or manually generate:
```bash
curl -X POST http://localhost:8080/api/v1/backtest/live-signal \
  -H "Content-Type: application/json" \
  -d '{"symbol": "BTCUSDT", "strategy": "session_trader"}'
```

### Step 3: Verify Risk
Check the signal:
```
Entry: $91,657.88
Stop Loss: $90,741.30
Risk: $916.58 (1% of entry) ✅
```

## 📊 Comparison

### Old Signal (3.67% risk):
```
Entry: $91,657.88
Stop Loss: $88,288.31
Risk: $3,369.57
Reward: $8,423.93
R:R: 2.50:1

$10,000 account:
- Risk: $367
- Reward: $919
```

### New Signal (1% risk):
```
Entry: $91,657.88
Stop Loss: $90,741.30
Risk: $916.58
Reward: $2,291.45
R:R: 2.50:1

$10,000 account:
- Risk: $100 ✅
- Reward: $250
```

## ✅ Summary

**Changed:** All strategies now use 1% risk
**Stop Loss:** Always 1% from entry price
**Risk/Reward:** Maintained (2:1 to 3:1)
**Position Sizing:** Simplified (1% risk = full account)
**Consistency:** Matches backtest settings ✅

Your live signals will now match your backtest risk profile! 🎯
