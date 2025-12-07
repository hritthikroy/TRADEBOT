# ⚡ TradingView Paper Trading - Quick Reference Card

## 🎯 Print This and Keep Next to Your Computer!

---

## ✅ ONE-TIME SETUP (5 Minutes)

### 1. Start Your Bot
```bash
cd backend
go run .
```
Open: http://localhost:8080
Click: "Live Signals" tab
Configure: Session Trader + BUY only
Click: "💾 Save All Settings"

### 2. Open TradingView
Go to: https://www.tradingview.com/chart/
Search: BTCUSDT
Timeframe: 15 minutes

### 3. Enable Paper Trading
Click: "Trading Panel" (bottom of screen)
Select: "Paper Trading"
Balance: $500
Click: "Connect"

---

## 🔄 DAILY ROUTINE

### When Signal Appears on Your Bot:

**Example Signal:**
```
🟢 BUY SIGNAL
Entry: $43,250
Stop Loss: $42,820
Take Profit: $44,110
```

### Execute in TradingView (30 seconds):

**Step 1:** Click "Buy/Long" button

**Step 2:** Fill in order:
- Order Type: Market
- Risk %: 0.5
- Entry: $43,250
- Stop Loss: $42,820

**Step 3:** Click "Buy/Long"

**Step 4:** After order fills:
- Click "Add Stop Loss" → Enter $42,820
- Click "Add Take Profit" → Enter $44,110

**Done!** TradingView manages the rest.

---

## 📊 POSITION SIZE CALCULATOR

### Quick Formula:
```
Risk Amount = Account × 0.5%
Risk Amount = $500 × 0.5% = $2.50

Distance to SL = Entry - Stop Loss
Distance to SL = $43,250 - $42,820 = $430

Position Size = Risk Amount / Distance to SL
Position Size = $2.50 / $430 = 0.0058 BTC
```

### Or Use TradingView's Calculator:
1. Click "Buy/Long"
2. Click "Risk" tab
3. Enter: 0.5% risk, SL price
4. TradingView calculates automatically!

---

## 🎯 WHAT TO EXPECT

### Per Trade:
- Risk: 0.5% ($2.50)
- Reward: 1.0% ($5.00)
- Risk/Reward: 2:1

### Per Day:
- Signals: 3-5
- Trades: 3-5
- Time spent: 10-15 minutes

### Per Week:
- Trades: 15-25
- Expected win rate: 48-50%
- Expected return: 3-5%

### After 50 Trades (2-4 weeks):
- Win rate: 48-50%
- Profit factor: 2-3
- Max drawdown: 6-8%
- Total return: 10-15%
- **Ready for live trading!** 🚀

---

## ✅ CHECKLIST (Print and Check Off)

### Before Each Trade:
- [ ] Bot showing clear signal
- [ ] Entry, SL, TP written down
- [ ] TradingView Paper Trading active
- [ ] Ready to execute

### During Execution:
- [ ] Clicked "Buy/Long"
- [ ] Entered 0.5% risk
- [ ] Set Stop Loss correctly
- [ ] Set Take Profit correctly
- [ ] Order confirmed

### After Trade Closes:
- [ ] Result noted (win/loss)
- [ ] P/L recorded
- [ ] Balance updated
- [ ] Notes written

---

## 🆘 QUICK TROUBLESHOOTING

### Problem: Can't find Trading Panel
**Solution:** Look at bottom of TradingView chart, click "Trading Panel" button

### Problem: Don't see Paper Trading option
**Solution:** Click dropdown in Trading Panel, select "Paper Trading"

### Problem: Don't know position size
**Solution:** Use TradingView's risk calculator (0.5% risk)

### Problem: Trade closed too early
**Solution:** Check if SL/TP were set correctly

### Problem: Results don't match backtest
**Solution:** Continue paper trading, check if following signals exactly

---

## 📝 DAILY LOG TEMPLATE

```
Date: ___________

Trade 1: ___:___ AM/PM
Signal: BUY / SELL
Entry: $_______
SL: $_______
TP: $_______
Result: TP Hit / SL Hit
P/L: $_______ (+/-)
Balance: $_______
Notes: _________________

Trade 2: ___:___ AM/PM
Signal: BUY / SELL
Entry: $_______
SL: $_______
TP: $_______
Result: TP Hit / SL Hit
P/L: $_______ (+/-)
Balance: $_______
Notes: _________________

Trade 3: ___:___ AM/PM
Signal: BUY / SELL
Entry: $_______
SL: $_______
TP: $_______
Result: TP Hit / SL Hit
P/L: $_______ (+/-)
Balance: $_______
Notes: _________________

Daily Summary:
Total Trades: ___
Wins: ___ (___%)
Losses: ___ (___%)
Total P/L: $_______ (___%)
Balance: $_______
```

---

## 🎯 PERFORMANCE TARGETS

### After 10 Trades (Week 1):
- [ ] Win rate: 40-55% (learning)
- [ ] Following signals: 100%
- [ ] Comfortable with platform

### After 25 Trades (Week 2):
- [ ] Win rate: 45-52%
- [ ] Profit factor: > 1.5
- [ ] Max drawdown: < 15%

### After 50 Trades (Week 3-4):
- [ ] Win rate: 48-50%
- [ ] Profit factor: 2-3
- [ ] Max drawdown: 6-8%
- [ ] **READY FOR LIVE!** 🚀

---

## 💡 QUICK TIPS

1. **Two Screens:** Bot on left, TradingView on right
2. **Write It Down:** Don't trust memory, write every signal
3. **Follow Rules:** No exceptions, no cherry-picking
4. **Be Patient:** Wait for high-quality setups
5. **Review Daily:** Check TradingView performance tab
6. **Stay Calm:** Losses are normal (48% WR = 52% losses)
7. **Learn:** Write why each trade won/lost

---

## 📞 HELP RESOURCES

### Detailed Guides:
- **HOW_IT_WORKS_TRADINGVIEW.md** - Step-by-step examples
- **TRADINGVIEW_PAPER_TRADING_GUIDE.md** - Complete guide
- **SIMPLE_WORKFLOW_DIAGRAM.md** - Visual workflow

### Quick Start:
- **PAPER_TRADING_START_HERE.md** - Overview
- **QUICK_START_PAPER_TRADING.md** - 2-minute setup

### Comparison:
- **PAPER_TRADING_METHODS_COMPARISON.md** - All methods

---

## 🚀 START NOW!

```bash
# 1. Start bot
./start_paper_trading.sh

# 2. Open TradingView
https://www.tradingview.com/chart/

# 3. Enable Paper Trading

# 4. Wait for signal

# 5. Execute!
```

---

## 📊 KEYBOARD SHORTCUTS (TradingView)

- **B** - Buy/Long
- **S** - Sell/Short
- **Esc** - Cancel order
- **Alt+T** - Open Trading Panel
- **Alt+W** - Open Watchlist

---

## ✅ SUCCESS FORMULA

```
Your Bot (Signals)
    +
TradingView (Execution)
    +
Your Discipline (Follow Rules)
    =
PROFITABLE TRADING! 🎉
```

---

**Print this page and keep it next to your computer!**

**Good luck! 🚀**
