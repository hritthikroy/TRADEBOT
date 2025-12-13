# 🚀 Unified Backtest Engine

## Overview

The **Unified Backtest Engine** is a powerful, all-in-one backtesting solution that merges the best features from all previous backtest engines into a single, comprehensive system.

## 🎯 Key Features

### 1. **Multiple Execution Modes**
- **Standard Backtest**: Fast, efficient backtesting with all essential features
- **Partial Exits**: Professional 3-stage exit strategy (30%, 30%, 40%)
- **Walk-Forward Analysis**: Test strategy robustness across time periods
- **Parallel Testing**: Test multiple strategies simultaneously

### 2. **Advanced Risk Management**
- Dynamic position sizing based on equity
- Maximum daily/weekly loss limits
- Consecutive loss protection
- Maximum trades per day limit
- Realistic slippage and fees

### 3. **Market Condition Filters**
- Volatility filters (min/max ATR)
- Volume filters
- Trading hours restrictions
- Weekend/holiday filtering

### 4. **Comprehensive Analysis**
- **Monte Carlo Simulation**: 1000+ runs to assess probability
- **Stress Testing**: Test under crash/rally scenarios
- **Walk-Forward Validation**: Prevent overfitting
- **Multi-Timeframe Analysis**: Confluence-based signals

### 5. **Professional Metrics**
- Sharpe Ratio (risk-adjusted returns)
- Sortino Ratio (downside risk)
- Calmar Ratio (return vs drawdown)
- Recovery Factor
- Maximum consecutive losses
- Win/loss streaks
- Expectancy per trade

## 📊 Usage

### Basic Usage

```go
config := backtest.UnifiedBacktestConfig{
    Symbol:       "BTCUSDT",
    Interval:     "15m",
    Days:         30,
    StartBalance: 10000,
    Strategy:     "liquidity_hunter",
}

candles, _ := backtest.FetchBinanceData(config.Symbol, config.Interval, config.Days)
result, _ := backtest.RunUnifiedBacktest(config, candles)
```

### Advanced Usage with All Features

```go
config := backtest.UnifiedBacktestConfig{
    // Basic
    Symbol:       "BTCUSDT",
    Interval:     "15m",
    Days:         90,
    StartBalance: 10000,
    Strategy:     "liquidity_hunter",
    
    // Risk Management
    RiskPercent:        0.003,  // 0.3% per trade
    MaxDailyLoss:       5.0,    // Stop if lose 5% in a day
    MaxConsecutiveLoss: 3,      // Stop after 3 losses
    MaxTradesPerDay:    20,     // Max 20 trades/day
    
    // Market Filters
    TradingHoursOnly:   true,
    MinVolatility:      0.5,
    MaxVolatility:      3.0,
    
    // Advanced Analysis
    EnableMonteCarlo:   true,
    EnableStressTest:   true,
    EnableWalkForward:  true,
    EnablePartialExits: true,
    
    // Simulation
    WindowType:         "expanding",
    RealisticSlippage:  true,
}

result, _ := backtest.RunUnifiedBacktest(config, candles)
```

### Parallel Strategy Testing

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
// Returns best performing strategy
```

## 🎨 Configuration Options

### Basic Configuration
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `Symbol` | string | "BTCUSDT" | Trading pair |
| `Interval` | string | "15m" | Timeframe |
| `Days` | int | 30 | Historical days |
| `StartBalance` | float64 | 10000 | Starting capital |
| `Strategy` | string | "liquidity_hunter" | Strategy name |

### Risk Management
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `RiskPercent` | float64 | 0.003 | Risk per trade (0.3%) |
| `MaxPositionCap` | float64 | StartBalance*10 | Max position size |
| `MaxDailyLoss` | float64 | 0 | Stop if daily loss exceeds % |
| `MaxWeeklyLoss` | float64 | 0 | Stop if weekly loss exceeds % |
| `MaxConsecutiveLoss` | int | 0 | Stop after N losses |
| `DynamicPositionSize` | bool | false | Adjust size with equity |
| `MaxTradesPerDay` | int | 20 | Max trades per day |

### Trading Costs
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `SlippagePercent` | float64 | 0.0015 | Slippage (0.15%) |
| `FeePercent` | float64 | 0.001 | Trading fees (0.1%) |
| `RealisticSlippage` | bool | false | Variable slippage |
| `IncludeSpread` | bool | false | Include bid-ask spread |
| `SpreadPercent` | float64 | 0 | Spread % |

### Market Filters
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `MinVolatility` | float64 | 0 | Min ATR to trade |
| `MaxVolatility` | float64 | 0 | Max ATR to trade |
| `MinVolume` | float64 | 0 | Min volume multiplier |
| `TradingHoursOnly` | bool | false | Trade 8am-8pm UTC only |

### Simulation Methods
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `WindowType` | string | "expanding" | "expanding", "rolling", "fixed" |
| `MinWindow` | int | 100 | Min candles needed |
| `MaxWindow` | int | 200 | Max window size |
| `UseWalkForward` | bool | false | Walk-forward analysis |
| `TrainingDays` | int | 60 | Training period days |
| `TestingDays` | int | 30 | Testing period days |

### Advanced Analysis
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `EnableMonteCarlo` | bool | false | Monte Carlo simulation |
| `MonteCarloRuns` | int | 1000 | Number of MC runs |
| `EnableStressTest` | bool | false | Stress testing |
| `EnablePartialExits` | bool | false | 3-stage exits |
| `EnableParallel` | bool | false | Parallel strategies |
| `Strategies` | []string | nil | Strategy list for parallel |

## 📈 Result Metrics

### Basic Metrics
- Total Trades
- Winning/Losing Trades
- Win Rate
- Total Profit/Loss
- Net Profit
- Return Percent
- Profit Factor
- Average RR
- Max Drawdown

### Advanced Risk Metrics
- **Sharpe Ratio**: Risk-adjusted returns (>1.0 is good, >2.0 is excellent)
- **Sortino Ratio**: Downside risk-adjusted returns
- **Calmar Ratio**: Annual return / max drawdown
- **Recovery Factor**: Net profit / max drawdown
- **Max Consecutive Losses**: Longest losing streak

### Performance Breakdown
- Win/Loss Streak Max
- Average Win/Loss
- Largest Win/Loss
- Expectancy Per Trade
- Average Trade Hours
- Best/Worst Trading Hour

### Advanced Analysis Results
- **Monte Carlo**: Probability distributions, confidence intervals
- **Walk-Forward**: In-sample vs out-sample performance
- **Stress Test**: Performance under extreme conditions

## 🏆 Performance Ratings

The engine automatically rates strategy performance:

- **🔥 EXCEPTIONAL**: Score >100, WR ≥70%, PF ≥2.5
- **⭐ EXCELLENT**: Score >50, WR ≥60%, PF ≥2.0
- **✅ GOOD**: Score >20, WR ≥55%, PF ≥1.5
- **⚠️ MODERATE**: Score >0, WR ≥50%
- **❌ POOR**: Needs improvement

## 🔧 API Endpoint

### POST `/api/backtest/unified`

```json
{
  "symbol": "BTCUSDT",
  "interval": "15m",
  "days": 30,
  "startBalance": 10000,
  "strategy": "liquidity_hunter",
  "enableMonteCarlo": true,
  "enableStressTest": true,
  "enableWalkForward": false,
  "enablePartialExits": true,
  "riskPercent": 0.003,
  "maxTradesPerDay": 20,
  "tradingHoursOnly": true
}
```

## 🎯 Best Practices

1. **Start Simple**: Begin with standard backtest, add features as needed
2. **Use Walk-Forward**: Validate strategy robustness across time
3. **Enable Monte Carlo**: Understand probability of outcomes
4. **Set Risk Limits**: Use MaxDailyLoss and MaxConsecutiveLoss
5. **Filter Market Conditions**: Use volatility and time filters
6. **Test Multiple Strategies**: Use parallel mode to compare
7. **Realistic Costs**: Enable RealisticSlippage for production estimates

## 🚀 Performance

- **Standard Backtest**: ~100-500ms for 30 days
- **With Monte Carlo**: ~2-5 seconds (1000 runs)
- **Walk-Forward**: ~5-10 seconds (5 periods)
- **Parallel Strategies**: ~1-3 seconds (4 strategies)

## 📝 Example Output

```
🚀 Starting Unified Backtest Engine
📊 Symbol: BTCUSDT | Interval: 15m | Days: 30 | Strategy: liquidity_hunter
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📊 UNIFIED BACKTEST RESULTS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🏆 STRATEGY: liquidity_hunter
⏱️  Duration: 1.234s

💰 PERFORMANCE:
  Start Balance:    $10000.00
  Final Balance:    $15234.56
  Net Profit:       $5234.56
  Return:           52.35%
  Profit Factor:    2.45

📈 TRADE STATISTICS:
  Total Trades:     45
  Winning Trades:   28
  Losing Trades:    17
  Win Rate:         62.22%
  Average RR:       1.85

⚠️  RISK METRICS:
  Max Drawdown:     8.45%
  Sharpe Ratio:     2.15
  Sortino Ratio:    3.42
  Calmar Ratio:     6.19
  Recovery Factor:  6.19
  Max Consecutive Losses: 3

🎲 MONTE CARLO ANALYSIS:
  Runs:             1000
  Mean Return:      48.23%
  95% Confidence:   25.34% to 71.12%
  Probability Profit: 87.5%
  Probability Ruin:   0.2%

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🏆 OVERALL RATING: ⭐ EXCELLENT - Professional-grade strategy
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

## 🔄 Migration from Old Engines

### From `RunBacktest`
```go
// Old
result, _ := backtest.RunBacktest(config, candles)

// New
unifiedConfig := backtest.UnifiedBacktestConfig{
    Symbol:       config.Symbol,
    Interval:     config.Interval,
    Days:         config.Days,
    StartBalance: config.StartBalance,
    Strategy:     config.Strategy,
}
result, _ := backtest.RunUnifiedBacktest(unifiedConfig, candles)
```

### From `RunProfessionalBacktest`
```go
// Old
result, _ := backtest.RunProfessionalBacktest(config, candles)

// New
unifiedConfig := backtest.UnifiedBacktestConfig{
    Symbol:             config.Symbol,
    Interval:           config.Interval,
    Days:               config.Days,
    StartBalance:       config.StartBalance,
    Strategy:           config.Strategy,
    EnablePartialExits: true,  // Key difference
}
result, _ := backtest.RunUnifiedBacktest(unifiedConfig, candles)
```

### From `RunWorldClassBacktest`
```go
// Old
wcConfig := backtest.WorldClassBacktestConfig{...}
result, _ := backtest.RunWorldClassBacktest(wcConfig, candles)

// New - All features included!
unifiedConfig := backtest.UnifiedBacktestConfig{
    Symbol:           wcConfig.Symbol,
    EnableMonteCarlo: true,
    EnableStressTest: true,
    // ... all other features
}
result, _ := backtest.RunUnifiedBacktest(unifiedConfig, candles)
```

## 🎉 Benefits

1. **Single Engine**: One codebase, easier maintenance
2. **Modular**: Enable only features you need
3. **Fast**: Optimized for performance
4. **Comprehensive**: All metrics in one place
5. **Professional**: Industry-standard risk metrics
6. **Flexible**: Works with any strategy
7. **Validated**: Walk-forward and Monte Carlo built-in

---

**The Unified Backtest Engine - One Engine to Rule Them All! 🚀**
