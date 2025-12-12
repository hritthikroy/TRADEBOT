# 🎯 AMD Phases - Quick Reference Card

## One-Line Summary
**Wyckoff AMD detection: Trade accumulation breakouts (🟢) and distribution breakdowns (🔴), avoid manipulation (⚠️)**

---

## 🚀 Quick Start
```bash
./compare_before_after_amd.sh  # Test and compare
```

---

## 📊 The 5 Phases

| Phase | Icon | Action | Quality | R:R |
|-------|------|--------|---------|-----|
| Accumulation | 🟢 | **BUY** | 95/100 | 8:1 |
| Markup | 📈 | BUY | 88/100 | 6:1 |
| Distribution | 🔴 | **SELL** | 95/100 | 8:1 |
| Markdown | 📉 | SELL | 90/100 | 6:1 |
| Manipulation | ⚠️ | **SKIP** | 0/100 | - |

---

## 🎯 Trading Rules

### BUY Priority
1. 🟢 Accumulation breakout → **BEST BUY**
2. 📈 Markup continuation → Good BUY
3. 🔴 Distribution → **SKIP**
4. ⚠️ Manipulation → **SKIP**

### SELL Priority
1. 🔴 Distribution breakdown → **BEST SELL**
2. 📉 Markdown continuation → Good SELL
3. 🟢 Accumulation → **SKIP**
4. ⚠️ Manipulation → **SKIP**

---

## 🔍 Phase Identification

### 🟢 Accumulation
- Price at bottom (lower 40%)
- Tight consolidation
- 2-4 volume spikes
- Narrow candles

### 🔴 Distribution
- Price at top (upper 40%)
- Tight consolidation
- 3+ volume spikes
- RSI > 55

### ⚠️ Manipulation
- 4+ volume spikes
- 5+ wide candles
- Whipsawing
- Middle of range

### 📈 Markup
- EMA9 > EMA21 > EMA50
- Price > EMA50
- Uptrend

### 📉 Markdown
- EMA9 < EMA21 < EMA50
- Price < EMA50
- Downtrend

---

## 📈 Expected Results

| Metric | Before | After | Target |
|--------|--------|-------|--------|
| Trades | 81 | 40-60 | -30% |
| Win Rate | 49% | 55-65% | +10% |
| Profit Factor | 2.82 | 3.5-5.0 | +40% |
| Drawdown | 35% | 20-30% | -30% |

---

## ✅ Keep AMD If

- ✅ Win rate > 55%
- ✅ Profit factor > 3.5
- ✅ Drawdown < 30%
- ✅ Better signal quality

---

## ❌ Rollback If

- ❌ Win rate < 45%
- ❌ Profit factor < 2.5
- ❌ Too few trades
- ❌ No improvement

---

## 🔄 Rollback Command

```bash
cp backend/unified_signal_generator.go.backup backend/unified_signal_generator.go
cd backend && go build -o ../tradebot && cd ..
pkill tradebot && ./tradebot &
```

---

## 📚 Full Documentation

- **START_HERE_AMD_PHASES.md** - Complete guide
- **AMD_PHASES_IMPLEMENTATION.md** - Technical details
- **AMD_PHASES_VISUAL_GUIDE.md** - Visual examples

---

## 🎓 Remember

**Best trades = Accumulation breakouts (🟢) and Distribution breakdowns (🔴)**

**Worst trades = Manipulation phase (⚠️) - always skip!**

---

**Test now:** `./compare_before_after_amd.sh`
