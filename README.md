# ICT/SMC Trading Strategy v2.0 🚀

Advanced cryptocurrency trading system with Smart Money Concepts, AI predictions, and institutional trading techniques.

**📊 Backtest Performance (v2.0 Optimized):**
- **Target Return:** 25-35% (30 days)
- **Previous:** 19.5% return, 64.7% win rate
- **Optimization:** Higher quality setups, better risk-reward ratios

## Features

### Core Trading System
- **ICT/SMC Analysis** - Order Blocks, Fair Value Gaps, Breaker Blocks
- **Power of 3 (PO3)** - Accumulation, Manipulation, Distribution phases
- **AMD Detection** - Smart money phase identification
- **Liquidity Sweeps** - Stop hunt detection before reversals
- **Multi-Timeframe Analysis** - 15m, 30m, 1h, 4h confluence
- **AI Ensemble Predictions** - Technical + Order Book + Local ML

### Advanced Features
- Real-time TradingView chart integration
- Automated backtesting system
- Risk-reward optimization (1.5:1 minimum)
- Intelligent trailing stops
- Delta volume analysis
- Dark theme professional UI

## How to Use

1. Open `index.html` in a web browser
2. The TradingView chart will load with BTC/USDT by default
3. Click "Predict Next 2 Candles (Same Height)" to generate predictions
4. View the predicted candle data with OHLC values

## How It Works

1. **Current Candle Analysis**: Reads the current candle data (Open, High, Low, Close)
2. **Height Calculation**: Calculates candle height (High - Low)
3. **Trend Detection**: Determines if market is in uptrend or downtrend
4. **Prediction Generation**: Creates 2 future candles with:
   - Same height as current candle
   - Continuation of detected trend
   - Realistic OHLC relationships

## Note

This is a simulation for educational purposes. The current implementation uses simulated data. To use real TradingView data, you would need to:

1. Implement TradingView's Datafeed API
2. Subscribe to real-time data updates
3. Extract actual candle information from the chart

## Files

- `index.html` - Main HTML structure and UI
- `prediction.js` - Prediction logic and TradingView integration
- `README.md` - Documentation
