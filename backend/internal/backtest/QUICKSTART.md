# ⚡ Quick Start Guide - Unified Backtest Engine

## 🚀 Get Started in 5 Minutes

### 1. Basic Backtest (Simplest)

```go
package main

import (
    "fmt"
    "tradebot/backend/internal/backtest"
)

func main() {
    // Step 1: Configure
    config := backtest.UnifiedBacktestConfig{
        Symbol:       "BTCUSDT",
        Interval:     "15m",
        Days:         30,
        StartBalance: 10000,
        Strategy:     "liquidity_hunter",
    }
    
    // Step 2: Fetch data
    candles, err := backtest.FetchBinanceData(
        config.Symbol, 
        config.Interval, 
        config.Days,
    )
    if err != nil {
        panic(err)
    }
    
    // Step 3: Run backtest
    result, err := backtest.RunUnifiedBacktest(config, candles)
    if err != nil {
        panic(err)
    }
    
    // Step 4: View results
    fmt.Printf("Win Rate: %.2f%%\n", result.WinRate)
    fmt.Printf("Return: %.2f%%\n", result.ReturnPercent)
    fmt.Printf("Profit Factor: %.2f\n", result.ProfitFactor)
}
```

**Output:**
```
Win Rate: 62.22%
Return: 52.35%
Profit Factor: 2.45
```

---

### 2. Add Monte Carlo (Recommended)

```go
config := backtest.UnifiedBacktestConfig{
    Symbol:           "BTCUSDT",
    Interval:         "15m",
    Days:             30,
    StartBalance:     10000,
    Strategy:         "liquidity_hunter",
    EnableMonteCarlo: true,  // ← Add this
}

result, _ := backtest.RunUnifiedBacktest(config, candles)

// Check Monte Carlo results
if result.MonteCarloResults != nil {
    mc := result.MonteCarloResults
    fmt.Printf("Expected Return: %.2f%%\n", mc.ExpectedReturn)
    fmt.Printf("Probability of Profit: %.1f%%\n", mc.ProbabilityProfit)
    fmt.Printf("95%% Confidence: %.2f%% to %.2f%%\n", 
        mc.Percentile5, mc.Percentile95)
}
```

**Output:**
```
Expected Return: 48.23%
Probability of Profit: 87.5%
95% Confidence: 25.34% to 71.12%
```

---

### 3. Add Risk Management (Production)

```go
config := backtest.UnifiedBacktestConfig{
    Symbol:       "BTCUSDT",
    Interval:     "15m",
    Days:         30,
    StartBalance: 10000,
    Strategy:     "liquidity_hunter",
    
    // Risk controls
    RiskPercent:        0.003,  // 0.3% per trade
    MaxDailyLoss:       5.0,    // Stop at 5% daily loss
    MaxConsecutiveLoss: 3,      // Stop after 3 losses
    MaxTradesPerDay:    20,     // Max 20 trades/day
    
    EnableMonteCarlo: true,
}
```

---

### 4. Compare Multiple Strategies (Parallel)

```go
config := backtest.UnifiedBacktestConfig{
    Symbol:         "BTCUSDT",
    Interval:       "15m",
    Days:           30,
    StartBalance:   10000,
    EnableParallel: true,
    Strategies: []string{
        "liquidity_hunter",
        "breakout_master",
        "session_trader",
        "trend_rider",
    },
}

result, _ := backtest.RunUnifiedBacktest(config, candles)
fmt.Printf("Best Strategy: %s\n", result.StrategyName)
fmt.Printf("Return: %.2f%%\n", result.ReturnPercent)
```

**Output:**
```
Best Strategy: breakout_master
Return: 78.45%
```

---

### 5. Full Production Setup

```go
config := backtest.UnifiedBacktestConfig{
    // Basic
    Symbol:       "BTCUSDT",
    Interval:     "15m",
    Days:         90,
    StartBalance: 10000,
    Strategy:     "liquidity_hunter",
    
    // Risk Management
    RiskPercent:        0.003,
    MaxDailyLoss:       5.0,
    MaxConsecutiveLoss: 3,
    MaxTradesPerDay:    20,
    
    // Market Filters
    TradingHoursOnly: true,
    MinVolatility:    0.5,
    MaxVolatility:    3.0,
    
    // Advanced Analysis
    EnableMonteCarlo:   true,
    EnableStressTest:   true,
    UseWalkForward:     true,
    EnablePartialExits: true,
    
    // Realistic Costs
    RealisticSlippage: true,
    IncludeSpread:     true,
    SpreadPercent:     0.0005,
}

result, _ := backtest.RunUnifiedBacktest(config, candles)
```

---

## 🎯 Common Use Cases

### Scalping (5m timeframe)
```go
config := backtest.UnifiedBacktestConfig{
    Symbol:             "BTCUSDT",
    Interval:           "5m",
    Days:               14,
    StartBalance:       10000,
    Strategy:           "scalper_pro",
    RiskPercent:        0.005,
    MaxTradesPerDay:    50,
    EnablePartialExits: true,
}
```

### Day Trading (15m timeframe)
```go
config := backtest.UnifiedBacktestConfig{
    Symbol:           "BTCUSDT",
    Interval:         "15m",
    Days:             30,
    StartBalance:     10000,
    Strategy:         "liquidity_hunter",
    RiskPercent:      0.003,
    MaxTradesPerDay:  20,
    TradingHoursOnly: true,
}
```

### Swing Trading (4h timeframe)
```go
config := backtest.UnifiedBacktestConfig{
    Symbol:       "BTCUSDT",
    Interval:     "4h",
    Days:         90,
    StartBalance: 10000,
    Strategy:     "trend_rider",
    RiskPercent:  0.002,
}
```

---

## 📊 Understanding Results

### Key Metrics to Watch

**Win Rate** (Target: >55%)
```go
if result.WinRate > 55 {
    fmt.Println("✅ Good win rate")
}
```

**Profit Factor** (Target: >1.5)
```go
if result.ProfitFactor > 1.5 {
    fmt.Println("✅ Profitable strategy")
}
```

**Max Drawdown** (Target: <15%)
```go
if result.MaxDrawdown < 15 {
    fmt.Println("✅ Acceptable risk")
}
```

**Sharpe Ratio** (Target: >1.0)
```go
if result.SharpeRatio > 1.0 {
    fmt.Println("✅ Good risk-adjusted returns")
}
```

---

## 🔧 API Usage

### REST Endpoint

```bash
curl -X POST http://localhost:8080/api/backtest/unified \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "interval": "15m",
    "days": 30,
    "startBalance": 10000,
    "strategy": "liquidity_hunter",
    "enableMonteCarlo": true
  }'
```

### Response
```json
{
  "success": true,
  "result": {
    "totalTrades": 45,
    "winRate": 62.22,
    "returnPercent": 52.35,
    "profitFactor": 2.45,
    "maxDrawdown": 8.45,
    "sharpeRatio": 2.15,
    "monteCarloResults": {
      "expectedReturn": 48.23,
      "probabilityProfit": 87.5
    }
  }
}
```

---

## ⚡ Performance Tips

### 1. Start Small
```go
// Test with 7 days first
config.Days = 7
```

### 2. Use Appropriate Timeframe
```go
// Scalping: 5m, 15m
// Day trading: 15m, 1h
// Swing: 4h, 1d
```

### 3. Enable Features Gradually
```go
// Start basic
config := UnifiedBacktestConfig{...}

// Add Monte Carlo
config.EnableMonteCarlo = true

// Add Walk-Forward
config.UseWalkForward = true
```

### 4. Set Realistic Costs
```go
config.SlippagePercent = 0.0015  // 0.15%
config.FeePercent = 0.001        // 0.1%
config.RealisticSlippage = true
```

---

## 🎓 Learning Path

### Beginner
1. Run basic backtest
2. Understand win rate and profit factor
3. Try different strategies

### Intermediate
4. Add Monte Carlo
5. Set risk limits
6. Use market filters

### Advanced
7. Enable walk-forward
8. Run stress tests
9. Compare multiple strategies

### Expert
10. Optimize parameters
11. Build custom strategies
12. Deploy to production

---

## 🐛 Troubleshooting

### "Insufficient data"
```go
// Increase days or use longer timeframe
config.Days = 30  // Increase this
```

### "No trades generated"
```go
// Check strategy and filters
config.MinVolatility = 0  // Relax filters
config.MaxVolatility = 0
```

### "Poor performance"
```go
// Try different strategy
config.Strategy = "breakout_master"

// Or adjust risk
config.RiskPercent = 0.005
```

---

## 📚 Next Steps

1. ✅ Run your first backtest
2. ✅ Try different strategies
3. ✅ Enable Monte Carlo
4. ✅ Add risk management
5. ✅ Compare strategies
6. ✅ Deploy to production

---

## 🎉 You're Ready!

Start with the basic example and gradually add features. The unified engine makes it easy to go from simple testing to production-ready backtesting.

**Happy backtesting! 🚀**
