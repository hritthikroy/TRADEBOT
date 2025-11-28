# Market Prediction Simulator

A simulation system that predicts the next two candles with the same height using TradingView's Advanced Chart widget.

## Features

- TradingView Advanced Chart integration
- Predicts next 2 candles maintaining the same height as current candle
- Trend detection (uptrend/downtrend)
- Real-time visualization
- Dark theme UI

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
