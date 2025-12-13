# 🚀 Academic ORB Strategy - Quick Start Guide

## What is This?

This is an implementation of a **research-proven day trading strategy** that achieved:

- **1,637% return** over 8 years (2016-2023)
- **41.6% annualized return**
- **Sharpe Ratio of 2.81** (exceptional risk-adjusted returns)
- **36% alpha** (uncorrelated to market)
- **Only 12% max drawdown**

Compare this to S&P 500's 198% return during the same period!

## The Secret Sauce: Relative Volume

The strategy's success hinges on one critical filter: **Relative Volume**

```
Without Relative Volume filter:  29% total return  ❌
With Relative Volume > 100%:     1,637% total return ✅

That's a 56x improvement!
```

## How It Works (Simple Version)

1. **Find Stocks in Play** - Stocks with unusually high volume (news, earnings, etc.)
2. **Wait for Opening Range** - First 5 minutes of trading
3. **Trade the Breakout** - Enter when price breaks the opening range
4. **Manage Risk** - Use tight stop loss (10% of ATR)
5. **Exit at Close** - Close all positions at end of day

## Quick Start

### 1. Start the Server

```bash
cd backend
go run .
```

### 2. Open Web Interface

Navigate to: `http://localhost:8080/orb_academic.html`

### 3. Run Your First Backtest

Click "Run Backtest" with default settings:
- Timeframe: 5 minutes
- Period: 2016-2023
- Capital: $25,000

### 4. Compare Timeframes

Click "Compare 5m, 15m, 30m, 60m" to see which performs best.

## API Examples

### Run a Backtest

```bash
curl -X POST http://localhost:8080/api/v1/orb/backtest \
  -H "Content-Type: application/json" \
  -d '{
    "timeFrame": 5,
    "startDate": "2016-01-01",
    "endDate": "2023-12-31",
    "initialCapital": 25000,
    "topNStocks": 20,
    "minRelativeVol": 1.0
  }'
```

### Get Live Signals

```bash
curl http://localhost:8080/api/v1/orb/live-signals?timeframe=5
```

### Compare All Timeframes

```bash
curl -X POST http://localhost:8080/api/v1/orb/compare \
  -H "Content-Type: application/json" \
  -d '{
    "startDate": "2016-01-01",
    "endDate": "2023-12-31",
    "initialCapital": 25000
  }'
```

## Strategy Rules (Cheat Sheet)

### Stock Filters
```
✓ Price > $5
✓ Avg Volume (14d) > 1M shares
✓ ATR (14d) > $0.50
✓ Relative Volume > 100%  ← MOST IMPORTANT!
✓ Trade only top 20 stocks
```

### Entry Rules
```
Bullish Opening Range → LONG at OR high breakout
Bearish Opening Range → SHORT at OR low breakout
```

### Risk Management
```
Stop Loss: 10% of ATR from entry
Position Size: 1% risk per trade
Max Leverage: 4x
Exit: End of day or stop loss
```

## Top Performing Stocks (From Research)

### 5-Minute ORB Champions:
1. **DDD** - 385R profit, 21% win rate
2. **FSLR** - 370R profit, 20% win rate
3. **NVDA** - 309R profit, 19% win rate
4. **SWBI** - 272R profit, 24% win rate
5. **RCL** - 271R profit, 20% win rate

*Note: Low win rate but huge R-multiples = profitable!*

## Why This Works

### 1. Stocks in Play
Focus on stocks with news/catalysts that create abnormal volume and volatility.

### 2. Institutional Imbalance
Opening range captures institutional supply/demand imbalance.

### 3. Trend Continuation
High relative volume breakouts tend to continue throughout the day.

### 4. Uncorrelated Returns
Beta ≈ 0 means returns are independent of market direction.

## Performance by Timeframe

| Timeframe | Total Return | Sharpe | Win Rate | Best For |
|-----------|-------------|--------|----------|----------|
| **5 min** | **1,637%** | **2.81** | 48.4% | **Best overall** |
| 15 min | 272% | 1.43 | 44.7% | More conservative |
| 30 min | 21% | 0.21 | 42.4% | Not recommended |
| 60 min | 39% | 0.40 | 42.3% | Not recommended |

**Conclusion: Use 5-minute ORB!**

## Common Questions

### Q: Why is win rate only 48%?
**A:** The strategy has large winners (10R+) and small losers (1R). Average PnL per trade is positive.

### Q: Can I use this on crypto/forex?
**A:** Research was on US stocks only. Would need separate testing for other markets.

### Q: Do I need expensive data?
**A:** You need 1-minute intraday data for US stocks. Many providers offer this.

### Q: What about slippage?
**A:** Research doesn't include slippage. Add 0.01-0.02% buffer per trade.

### Q: Can I trade this manually?
**A:** Possible but challenging. Strategy works best with automation due to:
- Need to scan 7,000+ stocks daily
- Calculate relative volume in real-time
- Execute 20 trades simultaneously
- Monitor all positions intraday

## Integration with Your Existing Bot

This ORB strategy complements your ICT/SMC strategies:

```
ORB Strategy:          ICT/SMC Strategy:
├─ Stock selection     ├─ Entry refinement
├─ Timing (OR)         ├─ Liquidity analysis
├─ Volume filter       ├─ Order flow
└─ Breakout entry      └─ Precise exits
```

**Combined Approach:**
1. Use ORB to find Stocks in Play
2. Use ICT/SMC for entry/exit timing
3. Use ORB risk management
4. Use ICT for trade management

## Next Steps

1. ✅ Run backtests with different parameters
2. ✅ Study top performing stocks
3. ✅ Compare with your existing strategies
4. ✅ Paper trade the strategy
5. ✅ Integrate with your ICT/SMC system

## Files Created

```
backend/
├── orb_academic_strategy.go    # Core strategy
├── orb_backtest_engine.go      # Backtest engine
├── orb_handlers.go             # API handlers
└── routes.go                   # Routes (updated)

public/
└── orb_academic.html           # Web interface

Documentation/
├── ORB_ACADEMIC_STRATEGY.md    # Full documentation
├── ORB_QUICK_START.md          # This file
└── test_orb_strategy.sh        # Test script
```

## Test the Implementation

```bash
./test_orb_strategy.sh
```

This will test all endpoints and show expected results.

## Resources

- **Research Paper:** SSRN #4729284
- **Authors:** Carlo Zarattini, Andrea Barbon, Andrew Aziz
- **Book:** "How to Day Trade for a Living" by Andrew Aziz
- **Website:** Bear Bull Traders (bearbulltraders.com)

## Important Disclaimer

⚠️ **Past performance does not guarantee future results.**

This strategy:
- Requires discipline and proper risk management
- Works best with automation
- Needs quality data and fast execution
- Should be paper traded first
- May not work in all market conditions

Always use proper risk management and never risk more than you can afford to lose.

---

## Support

Questions? Check:
1. Full documentation: `ORB_ACADEMIC_STRATEGY.md`
2. Web interface: `http://localhost:8080/orb_academic.html`
3. Test script: `./test_orb_strategy.sh`

**Happy Trading! 📈**
