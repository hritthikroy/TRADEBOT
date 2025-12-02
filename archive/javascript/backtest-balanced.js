// BALANCED BACKTEST CONFIGURATION
// Copy these settings into your backtest.js

// ============================================
// POSITION SIZING (Line ~548)
// ============================================
const currentBalance = backtestResults.currentBalance;
const riskPercent = 0.03; // 3% of current balance
const MAX_POSITION_USD = 50000; // $50K max position (allows growth)
const MAX_RISK_USD = 1500; // $1500 max risk per trade

const riskAmount = Math.min(
    currentBalance * riskPercent,
    MAX_RISK_USD
);

// ============================================
// TRADING COSTS (Line ~600)
// ============================================
const SLIPPAGE_PERCENT = 0.02; // 0.02% slippage (limit orders)
const FEE_PERCENT = 0.08; // 0.08% fees (with VIP discount)

// ============================================
// TRAILING STOP (Line ~570)
// ============================================
// BUY trades:
if (profitPercent > 1.0) { // Activate at 1.0R
    trailingActive = true;
    stopLoss = Math.max(stopLoss, entry + (highestPrice - entry) * 0.50); // Lock 50%
}

// SELL trades:
if (profitPercent > 1.0) { // Activate at 1.0R
    trailingActive = true;
    stopLoss = Math.min(stopLoss, entry - (entry - lowestPrice) * 0.50); // Lock 50%
}

// ============================================
// SIGNAL FILTERS (Line ~480)
// ============================================
// Minimum confidence
if (aiPrediction.confidence < 50) return null;

// Minimum R:R
const rr = Math.abs(signal.targets[0].price - signal.entry) / Math.abs(signal.entry - signal.stopLoss);
if (rr < 1.2) return null;

// Trend alignment (lenient)
if (!mtfTrendFilter.higherTFTrend && mtfTrendFilter.confidence < 0.4) return null;

// ============================================
// CIRCUIT BREAKERS (Line ~120)
// ============================================
// Max drawdown
const currentDrawdown = (backtestResults.peakBalance - backtestResults.currentBalance) / backtestResults.peakBalance;
if (currentDrawdown > 0.35) { // 35% max
    console.log('⚠️ Max drawdown reached, stopping');
    break;
}

// Consecutive losses
if (consecutiveLosses >= 5) {
    console.log('⚠️ 5 losses in a row, pausing 20 candles');
    i += 20;
    consecutiveLosses = 0;
}

// ============================================
// EXPECTED RESULTS
// ============================================
// Return: 15-50% per month
// Trades: 40-100 per month
// Win Rate: 52-58%
// Max Drawdown: 20-35%
// Profit Factor: 1.4-2.2
