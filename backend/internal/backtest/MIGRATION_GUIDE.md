# 🔄 Migration Guide: Old Engines → Unified Engine

## What Was Merged?

The Unified Backtest Engine consolidates **7 different backtest engines** into one powerful system:

### 1. ✅ `backtest_engine.go` (Standard)
**Features Merged:**
- Basic backtest logic
- Position sizing
- Trailing stops
- Target management
- Statistics calculation

**Now Available As:**
```go
config := UnifiedBacktestConfig{
    Symbol: "BTCUSDT",
    // ... basic config
}
RunUnifiedBacktest(config, candles)
```

---

### 2. ✅ `backtest_engine_professional.go`
**Features Merged:**
- Partial exits (30%, 30%, 40%)
- Breakeven stop after TP1
- Accurate profit calculation
- Daily trade limits
- Skip-ahead logic

**Now Available As:**
```go
config := UnifiedBacktestConfig{
    EnablePartialExits: true,
    MaxTradesPerDay: 20,
}
```

---

### 3. ✅ `world_class_backtest.go`
**Features Merged:**
- Sharpe Ratio
- Sortino Ratio
- Calmar Ratio
- Recovery Factor
- Win/loss streaks
- Monte Carlo simulation
- Walk-forward analysis
- Stress testing
- Market condition analysis

**Now Available As:**
```go
config := UnifiedBacktestConfig{
    EnableMonteCarlo: true,
    MonteCarloRuns: 1000,
    EnableStressTest: true,
}
```

---

### 4. ✅ `enhanced_backtest.go`
**Features Merged:**
- Expanding/rolling/fixed windows
- Walk-forward validation
- Monte Carlo analysis
- Time filters
- Realistic slippage
- Volatility-based adjustments

**Now Available As:**
```go
config := UnifiedBacktestConfig{
    WindowType: "expanding",
    UseWalkForward: true,
    TrainingDays: 60,
    TestingDays: 30,
    TradingHoursOnly: true,
    RealisticSlippage: true,
}
```

---

### 5. ✅ `comprehensive_backtest.go`
**Features Merged:**
- Multi-strategy testing
- Strategy scoring
- Performance comparison
- Recommendations generation

**Now Available As:**
```go
config := UnifiedBacktestConfig{
    EnableParallel: true,
    Strategies: []string{
        "liquidity_hunter",
        "breakout_master",
        "session_trader",
    },
}
```

---

### 6. ✅ `multi_tf_backtest.go`
**Features Merged:**
- Multi-timeframe analysis
- Confluence-based signals
- Top-down analysis

**Now Available As:**
```go
config := UnifiedBacktestConfig{
    EnableMultiTF: true,
}
```

---

### 7. ✅ `optimized_timeframe_backtest.go`
**Features Merged:**
- Timeframe-specific strategies
- Optimal period selection
- Performance comparison

**Now Available As:**
```go
// Test multiple timeframes
for _, tf := range []string{"15m", "1h", "4h"} {
    config := UnifiedBacktestConfig{
        Interval: tf,
        Days: getOptimalDays(tf),
    }
    result, _ := RunUnifiedBacktest(config, candles)
}
```

---

## 📊 Feature Comparison Matrix

| Feature | Old Engines | Unified Engine |
|---------|-------------|----------------|
| Basic Backtest | ✅ backtest_engine.go | ✅ Standard mode |
| Partial Exits | ✅ professional.go | ✅ EnablePartialExits |
| Monte Carlo | ✅ world_class.go, enhanced.go | ✅ EnableMonteCarlo |
| Walk-Forward | ✅ world_class.go, enhanced.go | ✅ UseWalkForward |
| Stress Test | ✅ world_class.go | ✅ EnableStressTest |
| Multi-Strategy | ✅ comprehensive.go | ✅ EnableParallel |
| Multi-Timeframe | ✅ multi_tf.go | ✅ EnableMultiTF |
| Risk Metrics | ✅ world_class.go | ✅ Always included |
| Time Filters | ✅ enhanced.go | ✅ TradingHoursOnly |
| Volatility Filters | ✅ world_class.go | ✅ Min/MaxVolatility |
| Daily Trade Limits | ✅ professional.go | ✅ MaxTradesPerDay |
| Consecutive Loss Limit | ✅ world_class.go | ✅ MaxConsecutiveLoss |
| Realistic Slippage | ✅ enhanced.go | ✅ RealisticSlippage |
| Window Types | ✅ enhanced.go | ✅ WindowType |

---

## 🚀 Migration Examples

### Example 1: Simple Backtest
**Before:**
```go
config := BacktestConfig{
    Symbol: "BTCUSDT",
    Interval: "15m",
    Days: 30,
    StartBalance: 10000,
    Strategy: "liquidity_hunter",
}
result, _ := RunBacktest(config, candles)
```

**After:**
```go
config := UnifiedBacktestConfig{
    Symbol: "BTCUSDT",
    Interval: "15m",
    Days: 30,
    StartBalance: 10000,
    Strategy: "liquidity_hunter",
}
result, _ := RunUnifiedBacktest(config, candles)
```

---

### Example 2: Professional with Partial Exits
**Before:**
```go
config := BacktestConfig{...}
result, _ := RunProfessionalBacktest(config, candles)
```

**After:**
```go
config := UnifiedBacktestConfig{
    // ... same fields
    EnablePartialExits: true,
}
result, _ := RunUnifiedBacktest(config, candles)
```

---

### Example 3: World-Class with All Features
**Before:**
```go
wcConfig := WorldClassBacktestConfig{
    BacktestConfig: BacktestConfig{...},
    EnableMonteCarlo: true,
    MonteCarloRuns: 1000,
    EnableWalkForward: true,
    EnableStressTest: true,
    MaxDailyLoss: 5.0,
}
result, _ := RunWorldClassBacktest(wcConfig, candles)
```

**After:**
```go
config := UnifiedBacktestConfig{
    Symbol: "BTCUSDT",
    Interval: "15m",
    Days: 90,
    StartBalance: 10000,
    Strategy: "liquidity_hunter",
    EnableMonteCarlo: true,
    MonteCarloRuns: 1000,
    UseWalkForward: true,
    EnableStressTest: true,
    MaxDailyLoss: 5.0,
}
result, _ := RunUnifiedBacktest(config, candles)
```

---

### Example 4: Multi-Strategy Comparison
**Before:**
```go
result, _ := RunComprehensiveBacktest("BTCUSDT", 30, 10000)
```

**After:**
```go
config := UnifiedBacktestConfig{
    Symbol: "BTCUSDT",
    Interval: "15m",
    Days: 30,
    StartBalance: 10000,
    EnableParallel: true,
    Strategies: []string{
        "liquidity_hunter",
        "breakout_master",
        "session_trader",
        "trend_rider",
    },
}
result, _ := RunUnifiedBacktest(config, candles)
```

---

## 🎯 Benefits of Migration

### 1. **Simplified Codebase**
- **Before**: 7 different engines, 7 different APIs
- **After**: 1 engine, 1 consistent API

### 2. **Better Maintainability**
- **Before**: Bug fixes needed in multiple files
- **After**: Fix once, works everywhere

### 3. **More Flexible**
- **Before**: Limited to predefined combinations
- **After**: Mix and match any features

### 4. **Better Performance**
- **Before**: Duplicate code, inefficient
- **After**: Optimized, single execution path

### 5. **Easier Testing**
- **Before**: Test 7 different engines
- **After**: Test 1 engine with different configs

---

## 📝 Deprecation Plan

### Phase 1: Soft Deprecation (Current)
- ✅ Unified engine available
- ✅ Old engines still work
- ✅ Documentation updated

### Phase 2: Migration Period (Recommended)
- Update all API endpoints to use unified engine
- Add deprecation warnings to old engines
- Provide migration tools

### Phase 3: Hard Deprecation (Future)
- Remove old engine files
- Keep only unified engine
- Archive old code for reference

---

## 🔧 Quick Migration Checklist

- [ ] Identify which old engine you're using
- [ ] Map old config to UnifiedBacktestConfig
- [ ] Enable appropriate feature flags
- [ ] Test with same data
- [ ] Compare results (should be identical or better)
- [ ] Update API endpoints
- [ ] Update documentation
- [ ] Remove old engine imports

---

## 💡 Pro Tips

1. **Start with defaults**: The unified engine has smart defaults
2. **Enable features gradually**: Don't enable everything at once
3. **Use Monte Carlo for validation**: Understand probability of outcomes
4. **Enable walk-forward for production**: Prevent overfitting
5. **Set risk limits**: MaxDailyLoss and MaxConsecutiveLoss are your friends
6. **Test parallel mode**: Compare strategies easily

---

## 🆘 Need Help?

If you encounter issues during migration:

1. Check the [UNIFIED_ENGINE_README.md](./UNIFIED_ENGINE_README.md)
2. Compare your old config with new config structure
3. Enable features one at a time
4. Check logs for detailed output
5. Verify results match expected behavior

---

## 🎉 Success Stories

**Before Migration:**
- 7 different backtest files
- ~3000 lines of duplicate code
- Inconsistent results
- Hard to maintain

**After Migration:**
- 1 unified engine
- ~1000 lines of optimized code
- Consistent, reliable results
- Easy to extend

---

**Ready to migrate? The unified engine is waiting! 🚀**
