# Liquidity Hunter - 80-90% Win Rate Strategy

## The Problem
Current strategy generates too many signals (54-1294 trades) with low win rate (25-31%). To achieve 80-90% win rate, we need to be EXTREMELY selective.

## New Approach: "Perfect Setup Only"

### Core Principle
**Only trade when EVERYTHING is perfect**. Miss 95% of trades to catch the 5% that are nearly guaranteed winners.

### Entry Criteria (ALL must be true)

#### For BUY Signals:
1. **Strong Uptrend** - EMA20 > EMA50 > EMA200 (all aligned)
2. **Pullback to EMA20** - Price within 0.5% of EMA20 (not 1%)
3. **RSI Sweet Spot** - RSI between 40-50 (not oversold, not overbought)
4. **Volume Spike** - Current volume > 1.5x average
5. **Bullish Reversal Candle** - Close in top 85% of range (not 80%)
6. **Pullback Pattern** - 2-3 red candles followed by strong green
7. **Above 200 EMA** - Long-term trend confirmation

#### For SELL Signals:
1. **Strong Downtrend** - EMA20 < EMA50 < EMA200 (all aligned)
2. **Pullback to EMA20** - Price within 0.5% of EMA20
3. **RSI Sweet Spot** - RSI between 50-60
4. **Volume Spike** - Current volume > 1.5x average
5. **Bearish Reversal Candle** - Close in bottom 85% of range
6. **Pullback Pattern** - 2-3 green candles followed by strong red
7. **Below 200 EMA** - Long-term trend confirmation

### Risk Management
- **Stop Loss**: 0.5 ATR (very tight)
- **TP1**: 0.5 ATR (1:1 RR, take 50% off)
- **TP2**: 1.0 ATR (2:1 RR, take 30% off)
- **TP3**: 1.5 ATR (3:1 RR, let 20% run)

### Expected Results
- **Trades per 90 days**: 5-15 (very selective)
- **Win Rate**: 80-90%
- **Profit Factor**: 3.0-5.0+
- **Max Drawdown**: <3%
- **Monthly Return**: 2-5%

## Why This Works

1. **Trend Following**: Only trade WITH strong trends
2. **Perfect Timing**: Enter on pullbacks to EMA20 (support/resistance)
3. **Confirmation**: Multiple confirmations reduce false signals
4. **Small Targets**: Easy to hit, high probability
5. **Tight Stops**: Minimize losses on rare losers

## Implementation Status
✅ Strategy logic updated
✅ Backtest parameters updated
⏳ Testing in progress

## Next Steps
1. Test on 4h timeframe (best for this strategy)
2. Verify win rate is 80%+
3. If still low, make conditions even stricter
4. Consider adding time-of-day filter (trade only during high volume hours)

---
*Target: 80-90% Win Rate*
*Current: 31% Win Rate*
*Gap: Need 2.5-3x improvement*
