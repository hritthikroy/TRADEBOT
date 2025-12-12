# Academic ORB Strategy Implementation

## Overview

This implementation is based on the academic research paper:
**"A Profitable Day Trading Strategy For The U.S. Equity Market"**
by Carlo Zarattini, Andrea Barbon, and Andrew Aziz (2024)

Published in SSRN: https://papers.ssrn.com/sol3/papers.cfm?abstract_id=4729284

## Key Research Findings

### Performance Results (2016-2023)

| Strategy | Total Return | Annual Return | Sharpe Ratio | Max DD | Win Rate | Alpha | Beta |
|----------|-------------|---------------|--------------|--------|----------|-------|------|
| 5m ORB + RelVol | **1,637%** | **41.6%** | **2.81** | 12% | 48.4% | **36%** | 0.00 |
| 15m ORB + RelVol | 272% | 17.4% | 1.43 | 11% | 44.7% | 17% | -0.01 |
| 30m ORB + RelVol | 21% | 2.3% | 0.21 | 35% | 42.4% | 3% | 0.01 |
| 60m ORB + RelVol | 39% | 4.1% | 0.40 | 21% | 42.3% | 4% | 0.01 |
| S&P 500 (Buy & Hold) | 198% | 14.2% | 0.78 | 34% | 54.9% | 0% | 1.00 |

### Critical Success Factor: Relative Volume

The research found that **Relative Volume is the most important filter**:

- Without Relative Volume filter: Only **29% total return** (3.2% annual)
- With Relative Volume > 100%: **1,637% total return** (41.6% annual)

**This is a 56x improvement in performance!**

## Strategy Rules

### 1. Stock Selection Filters

All stocks must meet these criteria:

```
✓ Opening price > $5
✓ 14-day average volume > 1,000,000 shares
✓ 14-day ATR > $0.50
✓ Relative Volume > 100% (CRITICAL!)
✓ Trade only top 20 stocks by Relative Volume
```

### 2. Relative Volume Calculation

```
Relative Volume = Current OR Volume / Average OR Volume (14 days)

Where:
- Current OR Volume = Volume in first N minutes today
- Average OR Volume = Average volume in first N minutes over past 14 days
```

**Example:**
- Stock XYZ typically trades 500K shares in first 5 minutes
- Today it traded 2M shares in first 5 minutes
- Relative Volume = 2M / 500K = 4.0 (400%)
- This stock is "in play" and qualifies for trading

### 3. Direction Determination

Based on the opening range candle:

```
IF Opening Range Close > Opening Range Open:
    → LONG ONLY (place stop buy order at OR high)
    
IF Opening Range Close < Opening Range Open:
    → SHORT ONLY (place stop sell order at OR low)
    
IF Opening Range Close = Opening Range Open (Doji):
    → NO TRADE
```

### 4. Entry Rules

- **Long Entry:** Price breaks above Opening Range High
- **Short Entry:** Price breaks below Opening Range Low
- Entry is triggered only AFTER the opening range period ends

### 5. Risk Management

**Stop Loss:**
```
Stop Loss Distance = 10% of 14-day ATR

For LONG: Stop Loss = Entry Price - (ATR × 0.10)
For SHORT: Stop Loss = Entry Price + (ATR × 0.10)
```

**Position Sizing:**
```
Risk per trade = 1% of capital
Position Size = (Capital × 0.01) / (Entry Price - Stop Loss)

Maximum Leverage = 4x (FINRA regulation)
```

**Exit:**
- Exit at End of Day (4:00 PM ET) if not stopped out
- Exit immediately if stop loss is hit

### 6. Commission Costs

The research used **$0.0035 per share** (Interactive Brokers Pro - Tiered pricing as of Dec 2023)

All results are **NET of commissions** (entry + exit).

## Why This Strategy Works

### 1. Stocks in Play Concept

The strategy focuses on **Stocks in Play** - stocks with abnormal trading activity due to:

- Earnings reports/surprises
- FDA approvals/disapprovals
- Mergers & acquisitions
- Major contract wins/losses
- Management changes
- Technical breakouts

### 2. Institutional Imbalance

The opening range captures the initial supply/demand imbalance from institutional investors reacting to news. High relative volume confirms this imbalance is significant.

### 3. Trend Continuation

When a stock breaks out of its opening range with high volume, it tends to continue in that direction throughout the day.

### 4. Uncorrelated Returns

- **Beta ≈ 0.00** - No correlation with S&P 500
- **Alpha = 36%** - Pure strategy returns, not market exposure
- Works in bull and bear markets

## Top Performing Stocks (2016-2023)

### 5-Minute ORB Top 5:

| Symbol | Cumulative PnL (R) | Win Rate |
|--------|-------------------|----------|
| DDD | 385R | 21% |
| FSLR | 370R | 20% |
| NVDA | 309R | 19% |
| SWBI | 272R | 24% |
| RCL | 271R | 20% |

### Notable Patterns:

- Popular retail stocks (TSLA, NVDA, AMD) perform well
- High volume stocks show better results
- Win rates typically 17-24% (but large R-multiples on winners)

## Implementation Guide

### API Endpoints

#### 1. Run Backtest
```bash
POST /api/v1/orb/backtest
Content-Type: application/json

{
  "timeFrame": 5,
  "startDate": "2016-01-01",
  "endDate": "2023-12-31",
  "initialCapital": 25000,
  "topNStocks": 20,
  "minRelativeVol": 1.0
}
```

#### 2. Get Live Signals
```bash
GET /api/v1/orb/live-signals?timeframe=5
```

#### 3. Compare Timeframes
```bash
POST /api/v1/orb/compare
Content-Type: application/json

{
  "startDate": "2016-01-01",
  "endDate": "2023-12-31",
  "initialCapital": 25000
}
```

#### 4. Get Top Performers
```bash
GET /api/v1/orb/top-performers?timeframe=5
```

### Web Interface

Access the web interface at:
```
http://localhost:8080/orb_academic.html
```

Features:
- Run backtests with custom parameters
- Compare all timeframes (5m, 15m, 30m, 60m)
- View top performing stocks from research
- Interactive results visualization

## Code Structure

```
backend/
├── orb_academic_strategy.go    # Core strategy logic
├── orb_backtest_engine.go      # Backtesting engine
├── orb_handlers.go             # HTTP handlers
└── routes.go                   # API routes

public/
└── orb_academic.html           # Web interface
```

## Key Classes

### ORBAcademicStrategy
```go
type ORBAcademicStrategy struct {
    TimeFrame        int     // 5, 15, 30, or 60 minutes
    MinPrice         float64 // $5
    MinAvgVolume     float64 // 1,000,000 shares
    MinATR           float64 // $0.50
    MinRelativeVol   float64 // 1.0 (100%)
    TopNStocks       int     // 20
    StopLossATRPct   float64 // 0.10 (10%)
    RiskPerTrade     float64 // 0.01 (1%)
    MaxLeverage      float64 // 4.0
    CommissionPerShr float64 // $0.0035
}
```

### ORBSignal
```go
type ORBSignal struct {
    Symbol        string
    Direction     string    // "LONG" or "SHORT"
    EntryPrice    float64
    StopLoss      float64
    ProfitTarget  float64   // EOD close
    PositionSize  int
    RiskAmount    float64
    RelativeVol   float64
    ATR           float64
    PnL           float64
    PnLInR        float64   // Profit in R multiples
}
```

## Usage Examples

### Example 1: Basic Backtest

```go
// Create strategy
strategy := NewORBAcademicStrategy(5) // 5-minute ORB

// Create backtest engine
engine := NewORBBacktestEngine(strategy, 25000)

// Run backtest
result, err := RunORBBacktest(
    5,                    // timeFrame
    startDate,
    endDate,
    stocksDataByDay,
    25000,               // initialCapital
)

// Print results
fmt.Printf("Total Return: %.2f%%\n", result.TotalReturn * 100)
fmt.Printf("Sharpe Ratio: %.2f\n", result.SharpeRatio)
fmt.Printf("Win Rate: %.2f%%\n", result.WinRate * 100)
```

### Example 2: Generate Live Signals

```go
// Get today's stock data
stocksData := fetchTodayStockData()

// Process trading day
engine.ProcessTradingDay(time.Now(), stocksData)

// Get active signals
for symbol, signal := range engine.ActiveSignals {
    fmt.Printf("%s: %s at $%.2f, Stop: $%.2f, RelVol: %.1fx\n",
        symbol,
        signal.Direction,
        signal.EntryPrice,
        signal.StopLoss,
        signal.RelativeVol,
    )
}
```

## Performance Metrics Explained

### R-Multiple
- **R** = Risk per trade (distance from entry to stop loss)
- **PnL in R** = Actual profit/loss divided by R
- Example: If R = $0.50 and profit = $5.00, then PnL = 10R

### Sharpe Ratio
- Measures risk-adjusted returns
- Formula: (Return - Risk-Free Rate) / Volatility
- **2.81** is exceptional (>2.0 is excellent)

### Alpha
- Returns not explained by market exposure
- **36%** means strategy generates 36% annual return independent of market

### Beta
- Correlation with market (S&P 500)
- **0.00** means zero correlation - pure strategy returns

## Important Notes

### 1. Data Requirements

To run this strategy, you need:
- **Intraday 1-minute data** for all US stocks
- **Daily OHLCV data** for 14-day calculations
- **Historical opening range volumes** (14 days)

### 2. Execution Considerations

- Strategy requires fast execution (stop orders)
- Best suited for liquid stocks (>1M avg volume)
- Commission costs matter - use low-cost broker
- Slippage not included in research (add buffer)

### 3. Risk Warnings

- Past performance doesn't guarantee future results
- Strategy requires discipline and automation
- Market conditions can change
- Always use proper risk management

### 4. Improvements to Consider

- Add slippage modeling (0.01-0.02% per trade)
- Implement real-time news filtering
- Add pre-market gap analysis
- Consider sector rotation
- Implement dynamic position sizing

## Comparison with Your Existing Strategies

### ORB vs ICT/SMC:

| Feature | ORB Academic | ICT/SMC |
|---------|-------------|---------|
| **Timeframe** | Intraday (5-60 min) | Multi-timeframe |
| **Entry** | Breakout-based | Liquidity-based |
| **Focus** | Stocks in Play | Order flow |
| **Holding** | Intraday only | Can be multi-day |
| **Filters** | Relative Volume | Liquidity sweeps, FVG |
| **Complexity** | Simple | Complex |
| **Backtested** | Yes (academic) | Discretionary |

### Potential Integration:

You could combine both approaches:
1. Use ORB to identify Stocks in Play
2. Use ICT/SMC for precise entry/exit
3. Use ORB filters for stock selection
4. Use ICT concepts for trade management

## References

1. Zarattini, C., Barbon, A., & Aziz, A. (2024). "A Profitable Day Trading Strategy For The U.S. Equity Market". SSRN Electronic Journal.

2. Crabel, T. (1990). "Day Trading with Short Term Price Patterns and Opening Range Breakout".

3. Wilder, J. W. (1978). "New Concepts in Technical Trading Systems".

4. Aziz, A. (2015). "How to Day Trade for a Living".

## License

This implementation is for educational and research purposes. The strategy is based on published academic research and is freely available for use.

## Support

For questions or issues:
- Check the web interface at `/orb_academic.html`
- Review the API documentation above
- Examine the code in `backend/orb_*.go` files

---

**Remember:** The key to this strategy's success is the **Relative Volume filter**. Without it, performance drops dramatically. Always trade only the top 20 stocks by Relative Volume each day!
