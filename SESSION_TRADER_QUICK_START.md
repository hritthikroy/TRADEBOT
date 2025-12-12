# 🚀 SESSION TRADER - QUICK START GUIDE

**Last Updated:** December 7, 2024  
**Status:** ✅ Ready for Paper Trading

---

## 📋 Daily Workflow

### Step 1: Morning Check (Every Day)

```bash
./daily_session_trader_check.sh
```

**What it does:**
- Runs 7-day backtest
- Shows current performance
- Tells you if it's safe to trade today

**Decision Matrix:**
- 🟢 **GREEN** = Trade today (all criteria met)
- 🟡 **YELLOW** = Trade with caution (reduced size)
- 🔴 **RED** = Don't trade today (wait for better conditions)

### Step 2: Start Paper Trading (If Green/Yellow)

```bash
./start_paper_trading_session_trader.sh
```

**What it does:**
- Checks recent performance
- Starts automated paper trading
- Uses Session Trader strategy
- Trades BTCUSDT on 15m timeframe

### Step 3: Monitor Performance (Throughout Day)

```bash
./check_paper_trading_stats.sh
```

**What it does:**
- Shows current balance
- Displays win rate and profit factor
- Lists recent trades
- Gives performance rating

**Check frequency:**
- Morning: Before starting
- Midday: Check progress
- Evening: Review results

### Step 4: Stop Trading (End of Day or If Losing)

```bash
./stop_paper_trading.sh
```

**What it does:**
- Shows final statistics
- Stops paper trading
- Saves results

---

## 🎯 Trading Criteria

### ✅ Safe to Trade When:
- ✅ 7-day return > 0%
- ✅ 7-day win rate ≥ 45%
- ✅ 7-day profit factor ≥ 1.0

### ⚠️ Trade with Caution When:
- ✅ 7-day return > 0%
- ✅ 7-day profit factor ≥ 1.0
- ⚠️ 7-day win rate < 45%

**Action:** Reduce position size by 50%

### ❌ Don't Trade When:
- ❌ 7-day return < 0%
- ❌ 7-day profit factor < 1.0

**Action:** Wait for better market conditions

---

## 📊 Expected Performance

### Short-Term (7 days)
```
Win Rate:        45-55%
Profit Factor:   1.0-1.2
Return:          0-2%
Trades:          20-30
Status:          ✅ Can be profitable
```

### Medium-Term (14-30 days)
```
Win Rate:        35-45%
Profit Factor:   0.8-1.1
Return:          -1% to +2%
Trades:          50-120
Status:          ⚠️ Variable
```

---

## 🛡️ Risk Management

### Position Sizing
```
Paper Trading:   $1000 starting balance
Risk per Trade:  1% ($10 per trade)
Stop Loss:       1.5 × ATR
Max Open Trades: 3
```

### When to Stop Trading
1. ❌ 3 consecutive losses
2. ❌ Daily loss > 3%
3. ❌ Win rate drops below 30%
4. ❌ Profit factor drops below 0.7

---

## 📈 Going Live Checklist

### Before Trading Real Money:

1. **Paper Trade Successfully**
   - [ ] At least 2 weeks of paper trading
   - [ ] Positive overall return
   - [ ] Win rate ≥ 40%
   - [ ] Profit factor ≥ 1.0
   - [ ] At least 30 trades completed

2. **Understand the Strategy**
   - [ ] Know entry conditions
   - [ ] Know exit conditions
   - [ ] Understand risk management
   - [ ] Can identify good vs bad setups

3. **Risk Management Ready**
   - [ ] Position size calculated (0.5-1% risk)
   - [ ] Stop losses understood
   - [ ] Max daily loss limit set
   - [ ] Emergency stop plan ready

4. **Mental Preparation**
   - [ ] Comfortable with losses
   - [ ] Can follow rules without emotion
   - [ ] Have realistic expectations
   - [ ] Won't overtrade

### Starting Live Trading:

**Week 1-2: Micro Size**
- Risk: 0.25% per trade
- Max trades: 2 per day
- Goal: Get comfortable with real money

**Week 3-4: Small Size**
- Risk: 0.5% per trade
- Max trades: 3 per day
- Goal: Build confidence

**Month 2+: Normal Size**
- Risk: 1% per trade
- Max trades: 5 per day
- Goal: Consistent profitability

---

## 🔧 Troubleshooting

### Problem: Backend not running
```bash
cd backend
./tradebot
```

### Problem: No trades being generated
```bash
./diagnose_session_trader.sh
```

### Problem: Poor performance
```bash
# Check if market conditions changed
./daily_session_trader_check.sh

# If RED status, stop trading
./stop_paper_trading.sh
```

### Problem: Too many losses
```bash
# Stop immediately
./stop_paper_trading.sh

# Review recent performance
./check_paper_trading_stats.sh

# Wait for better conditions
# Check again tomorrow
```

---

## 📁 Important Files

### Scripts
- `daily_session_trader_check.sh` - Daily performance check
- `start_paper_trading_session_trader.sh` - Start paper trading
- `check_paper_trading_stats.sh` - Check current stats
- `stop_paper_trading.sh` - Stop paper trading
- `test_session_trader_simple.sh` - Run full backtest

### Documentation
- `SESSION_TRADER_PROFITABLE_VERSION.md` - Full strategy guide
- `SESSION_TRADER_FIXED_STATUS.md` - Technical details
- `SESSION_TRADER_QUICK_START.md` - This file

### Logs
- `session_trader_daily_checks.log` - Daily check history

---

## 💡 Tips for Success

### 1. Be Disciplined
- ✅ Follow the daily check routine
- ✅ Only trade when criteria are met
- ✅ Stop when performance degrades
- ❌ Don't trade on emotions
- ❌ Don't overtrade

### 2. Start Small
- Begin with paper trading
- Use micro sizes when going live
- Scale up slowly over weeks
- Never risk more than 1% per trade

### 3. Monitor Daily
- Check performance every morning
- Review trades every evening
- Keep a trading journal
- Learn from mistakes

### 4. Be Patient
- Strategy works best in trending markets
- Some days will have no trades
- Some weeks will be break-even
- Focus on long-term consistency

### 5. Know When to Stop
- Stop if 3 consecutive losses
- Stop if daily loss > 3%
- Stop if performance degrades
- Take breaks when needed

---

## 🎯 Success Metrics

### Daily
- [ ] Ran morning check
- [ ] Followed trading decision
- [ ] Monitored trades
- [ ] Reviewed evening performance

### Weekly
- [ ] Positive return
- [ ] Win rate ≥ 40%
- [ ] Profit factor ≥ 1.0
- [ ] No major drawdowns

### Monthly
- [ ] Consistent profitability
- [ ] Improving win rate
- [ ] Good risk management
- [ ] Following all rules

---

## 📞 Quick Commands

```bash
# Daily routine
./daily_session_trader_check.sh          # Morning check
./start_paper_trading_session_trader.sh  # Start trading
./check_paper_trading_stats.sh           # Check stats
./stop_paper_trading.sh                  # Stop trading

# Testing
./test_session_trader_simple.sh          # Full backtest
./diagnose_session_trader.sh             # Diagnostic

# Dashboard
open http://localhost:8080/paper-trading # View in browser
```

---

## ⚠️ Important Reminders

1. **This is a tool, not a guarantee**
   - Past performance ≠ future results
   - Markets change constantly
   - Always use risk management

2. **Start with paper trading**
   - No real money at risk
   - Learn the strategy
   - Build confidence

3. **Use small sizes when live**
   - Start with 0.25-0.5% risk
   - Scale up slowly
   - Never risk more than you can afford to lose

4. **Monitor daily**
   - Check performance every day
   - Stop if conditions change
   - Be disciplined

5. **Have realistic expectations**
   - Not every day will be profitable
   - Some weeks will be break-even
   - Focus on consistency, not home runs

---

## 🚀 Ready to Start?

1. **Run morning check:**
   ```bash
   ./daily_session_trader_check.sh
   ```

2. **If GREEN, start paper trading:**
   ```bash
   ./start_paper_trading_session_trader.sh
   ```

3. **Monitor throughout the day:**
   ```bash
   ./check_paper_trading_stats.sh
   ```

4. **Review and stop at end of day:**
   ```bash
   ./stop_paper_trading.sh
   ```

---

**Good luck and trade safely! 🎯**

