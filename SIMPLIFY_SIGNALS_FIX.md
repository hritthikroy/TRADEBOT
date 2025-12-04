# Signal Generation - Too Strict Problem

## Issue
Strategies require ALL 5-6 conditions to be true simultaneously, resulting in almost zero trades.

## Solution
Change from "ALL conditions" to "AT LEAST 3 out of 5 conditions" using a scoring system.

## Example Fix Applied to session_trader

### Before (Too Strict):
```go
if ema9 > ema21 && ema21 > ema50 && currentPrice > ema200 && 
   rsi > 40 && rsi < 70 && macd > signal && volumeConfirm {
    // Signal - requires ALL 5 conditions
}
```

### After (More Flexible):
```go
buyScore := 0
if ema9 > ema21 && ema21 > ema50 { buyScore++ }
if currentPrice > ema200 { buyScore++ }
if rsi > 40 && rsi < 70 { buyScore++ }
if macd > signal { buyScore++ }
if volumeConfirm { buyScore++ }

if buyScore >= 3 {  // At least 3 out of 5
    // Signal
}
```

## Need to Apply to All Strategies

The same fix needs to be applied to:
1. ✅ session_trader (DONE)
2. ❌ breakout_master
3. ❌ liquidity_hunter  
4. ❌ trend_rider
5. ❌ range_master
6. ❌ smart_money_tracker
7. ❌ institutional_follower
8. ❌ reversal_sniper
9. ❌ momentum_beast
10. ❌ scalper_pro

This will generate more trades and allow the optimizer to find profitable parameters.
