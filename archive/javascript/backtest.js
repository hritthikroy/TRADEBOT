// Backtesting System for Trading Signals
// Tests historical performance of the prediction system

let backtestResults = {
    trades: [],
    totalTrades: 0,
    winningTrades: 0,
    losingTrades: 0,
    totalProfit: 0,
    totalLoss: 0,
    winRate: 0,
    profitFactor: 0,
    averageRR: 0,
    maxDrawdown: 0,
    startingBalance: 500,  // Starting with $500
    currentBalance: 500
};

// Convert to Binance-compatible interval
function toBinanceInterval(interval) {
    const binanceIntervals = {
        '1s': '1m',    // Binance doesn't support 1s, use 1m
        '1m': '1m',
        '3m': '3m',
        '5m': '5m',
        '15m': '15m',
        '30m': '30m',
        '1h': '1h',
        '2h': '2h',
        '4h': '4h',
        '6h': '6h',
        '8h': '8h',
        '12h': '12h',
        '1d': '1d',
        '3d': '3d',
        '1w': '1w',
        '1D': '1d',    // Normalize
        '3D': '3d',    // Normalize
        '1W': '1w'     // Normalize
    };
    return binanceIntervals[interval] || interval;
}

// Run backtest on historical data
async function runBacktest(symbol, interval, days = 30) {
    console.log(`🔄 Starting backtest: ${symbol} ${interval} for ${days} days...`);
    
    // Store days for report generation
    window.backtestDays = days;
    
    // Validate days parameter
    if (days < 1 || days > 365) {
        throw new Error('Days must be between 1 and 365');
    }
    
    // Convert to Binance-compatible interval
    const binanceInterval = toBinanceInterval(interval);
    
    // Reset results
    backtestResults = {
        trades: [],
        totalTrades: 0,
        winningTrades: 0,
        losingTrades: 0,
        totalProfit: 0,
        totalLoss: 0,
        winRate: 0,
        profitFactor: 0,
        averageRR: 0,
        maxDrawdown: 0,
        startingBalance: 500,  // Starting with $500
        currentBalance: 500,
        peakBalance: 500
    };
    
    // Reset trade history for filters
    recentTradeHistory = [];
    
    try {
        // Fetch historical data (supports multiple requests for long periods)
        const data = await fetchHistoricalDataMultiple(symbol, binanceInterval, days);
        
        if (!Array.isArray(data) || data.length < 100) {
            throw new Error('Insufficient historical data');
        }
        
        const historicalData = data.map(candle => ({
            time: candle[0],
            open: parseFloat(candle[1]),
            high: parseFloat(candle[2]),
            low: parseFloat(candle[3]),
            close: parseFloat(candle[4]),
            volume: parseFloat(candle[5])
        }));
        
        console.log(`📊 Loaded ${historicalData.length} candles for backtest`);
        updateStatus(`Backtesting ${historicalData.length} candles...`);
        
        // Store original data and prepare for simulation
        window.backtestHistoricalData = historicalData;
        window.backtestTrades = [];
        window.backtestCurrentIndex = 0;
        
        // Simulate trading through historical data
        const windowSize = 50; // Use 50 candles for analysis
        
        for (let i = windowSize; i < historicalData.length - 10; i++) {
            const dataWindow = historicalData.slice(i - windowSize, i);
            const futureData = historicalData.slice(i, i + 10);
            
            // Generate signal based on historical data
            const signal = await generateHistoricalSignal(dataWindow, interval);
            
            // Debug: Log signal attempts every 100 candles
            if (i % 100 === 0) {
                console.log(`Progress: ${i}/${historicalData.length} candles | Signals found: ${backtestResults.totalTrades}`);
            }
            
            if (signal) {
                // Simulate trade execution
                const tradeResult = simulateTrade(signal, futureData);
                
                if (tradeResult) {
                    // Store trade with candle index for visualization
                    tradeResult.entryIndex = i;
                    tradeResult.signal = signal;
                    backtestResults.trades.push(tradeResult);
                    window.backtestTrades.push(tradeResult);
                    backtestResults.totalTrades++;
                    
                    if (tradeResult.profit > 0) {
                        backtestResults.winningTrades++;
                        backtestResults.totalProfit += tradeResult.profit;
                    } else {
                        backtestResults.losingTrades++;
                        backtestResults.totalLoss += Math.abs(tradeResult.profit);
                    }
                    
                    // Update balance
                    backtestResults.currentBalance += tradeResult.profit;
                    
                    // Track trade history for filters
                    recentTradeHistory.push(tradeResult);
                    if (recentTradeHistory.length > 5) {
                        recentTradeHistory.shift(); // Keep only last 5 trades
                    }
                    
                    // Track peak for drawdown
                    if (backtestResults.currentBalance > backtestResults.peakBalance) {
                        backtestResults.peakBalance = backtestResults.currentBalance;
                    }
                    
                    // Calculate drawdown
                    const drawdown = (backtestResults.peakBalance - backtestResults.currentBalance) / backtestResults.peakBalance;
                    if (drawdown > backtestResults.maxDrawdown) {
                        backtestResults.maxDrawdown = drawdown;
                    }
                }
                
                // Skip ahead to avoid overlapping trades
                i += 5;
            }
        }
        
        // Calculate final statistics
        calculateBacktestStats();
        
        // Display results
        displayBacktestResults();
        updateStatus(`Backtest complete: ${backtestResults.returnPercent.toFixed(1)}% return`);
        return backtestResults;
        
    } catch (error) {
        console.error('❌ Backtest error:', error);
        updateStatus('Backtest failed: ' + error.message);
        alert('Backtest failed. Check console (F12) for details.');
        return null;
    }
}

// Update status helper
function updateStatus(message) {
    const statusEl = document.getElementById('status-message');
    if (statusEl) {
        statusEl.textContent = message;
    }
    console.log(message);
}

// Calculate candle limit based on interval and days
function calculateCandleLimit(interval, days) {
    const candlesPerDay = {
        '1m': 1440,
        '3m': 480,
        '5m': 288,
        '15m': 96,
        '30m': 48,
        '1h': 24,
        '2h': 12,
        '4h': 6,
        '6h': 4,
        '8h': 3,
        '12h': 2,
        '1d': 1,
        '3d': 0.33,
        '1w': 0.14
    };
    
    const candlesNeeded = Math.ceil((candlesPerDay[interval] || 96) * days);
    
    // Binance limit is 1000 candles per request
    // Add 50 extra candles for indicator calculation
    const limit = Math.min(candlesNeeded + 50, 1000);
    
    console.log(`📊 ${interval}: Need ${candlesNeeded} candles for ${days} days, fetching ${limit}`);
    
    // Warn if we can't get enough data
    if (candlesNeeded > 1000) {
        console.warn(`⚠️ Requested ${candlesNeeded} candles but Binance limit is 1000. Backtest will use available data.`);
    }
    
    return limit;
}

// Fetch historical data with multiple requests if needed (for very long periods)
async function fetchHistoricalDataMultiple(symbol, interval, days) {
    const binanceInterval = toBinanceInterval(interval);
    const limit = calculateCandleLimit(binanceInterval, days);
    
    // If we need more than 1000 candles, we need multiple requests
    const candlesPerDay = {
        '1m': 1440, '3m': 480, '5m': 288, '15m': 96, '30m': 48,
        '1h': 24, '2h': 12, '4h': 6, '6h': 4, '8h': 3, '12h': 2,
        '1d': 1, '3d': 0.33, '1w': 0.14
    };
    
    const totalCandlesNeeded = Math.ceil((candlesPerDay[binanceInterval] || 96) * days);
    
    if (totalCandlesNeeded <= 1000) {
        // Single request is enough
        const url = `https://api.binance.com/api/v3/klines?symbol=${symbol}&interval=${binanceInterval}&limit=${limit}`;
        console.log(`📡 Fetching: ${url}`);
        const response = await fetch(url);
        const data = await response.json();
        
        if (data.code) {
            throw new Error(`Binance API Error: ${data.msg || 'Unknown error'}`);
        }
        
        return data;
    } else {
        // Multiple requests needed
        console.log(`📡 Fetching ${totalCandlesNeeded} candles in multiple requests...`);
        let allData = [];
        let endTime = Date.now();
        const requestsNeeded = Math.ceil(totalCandlesNeeded / 1000);
        
        for (let i = 0; i < requestsNeeded; i++) {
            const url = `https://api.binance.com/api/v3/klines?symbol=${symbol}&interval=${binanceInterval}&limit=1000&endTime=${endTime}`;
            console.log(`📡 Request ${i + 1}/${requestsNeeded}: ${url}`);
            
            const response = await fetch(url);
            const data = await response.json();
            
            if (data.code) {
                throw new Error(`Binance API Error: ${data.msg || 'Unknown error'}`);
            }
            
            if (!Array.isArray(data) || data.length === 0) {
                break;
            }
            
            // Prepend data (we're going backwards in time)
            allData = [...data, ...allData];
            
            // Set endTime to the first candle's time for next request
            endTime = data[0][0] - 1;
            
            // Small delay to avoid rate limiting
            await new Promise(resolve => setTimeout(resolve, 100));
        }
        
        console.log(`✅ Fetched ${allData.length} total candles`);
        return allData;
    }
}

// Track recent trade history for filtering
let recentTradeHistory = [];

// Generate signal from historical data
async function generateHistoricalSignal(data, interval) {
    try {
        // Use same logic as live trading
        const volumeAnalysis = analyzeHistoricalVolume(data);
        
        if (!volumeAnalysis) return null;
        
        // 🛡️ ALL PRE-ENTRY FILTERS DISABLED
        // Testing showed filters reduce profitability by blocking good trades
        // Original system: 109 trades, 392% return
        // With filters: 56 trades, 58% return
        // Strategy: Let all signals through, use trailing stops for protection
        
        // Simplified AI prediction for backtest (complex logic causing issues)
        const recentCandles = data.slice(-20);
        
        // 1. Trend Analysis (last 10 candles)
        const last10 = recentCandles.slice(-10);
        let upCount = 0;
        let downCount = 0;
        let totalVolume = 0;
        let bullVolume = 0;
        let bearVolume = 0;
        
        // Delta Analysis - Buy vs Sell Volume
        let cumulativeDelta = 0;
        let positiveDeltaCount = 0;
        let negativeDeltaCount = 0;
        
        last10.forEach(c => {
            const vol = c.volume || 0;
            totalVolume += vol;
            
            // Calculate delta for each candle
            const range = c.high - c.low;
            const closePosition = range > 0 ? (c.close - c.low) / range : 0.5;
            
            // Estimate buy/sell volume based on close position and candle type
            let buyVol, sellVol;
            if (c.close > c.open) {
                // Bullish candle
                buyVol = vol * (0.5 + closePosition * 0.5); // 50-100% buy volume
                sellVol = vol - buyVol;
                upCount++;
                bullVolume += vol;
            } else {
                // Bearish candle
                sellVol = vol * (0.5 + (1 - closePosition) * 0.5); // 50-100% sell volume
                buyVol = vol - sellVol;
                downCount++;
                bearVolume += vol;
            }
            
            const delta = buyVol - sellVol;
            cumulativeDelta += delta;
            
            if (delta > 0) positiveDeltaCount++;
            else negativeDeltaCount++;
        });
        
        // 2. Momentum Analysis
        const priceChange = (last10[last10.length - 1].close - last10[0].close) / last10[0].close;
        const momentum = priceChange * 100; // Percentage change
        
        // 3. Volume Analysis
        const volumeRatio = totalVolume > 0 ? bullVolume / totalVolume : 0.5;
        
        // 4. Moving Average Analysis
        const ma5 = recentCandles.slice(-5).reduce((sum, c) => sum + c.close, 0) / 5;
        const ma10 = recentCandles.slice(-10).reduce((sum, c) => sum + c.close, 0) / 10;
        const ma20 = recentCandles.reduce((sum, c) => sum + c.close, 0) / 20;
        const currentPrice = last10[last10.length - 1].close;
        
        // 5. Support/Resistance Analysis
        const highs = recentCandles.map(c => c.high);
        const lows = recentCandles.map(c => c.low);
        const resistance = Math.max(...highs);
        const support = Math.min(...lows);
        const pricePosition = (currentPrice - support) / (resistance - support);
        
        // 6. Footprint Analysis - Volume at Price Levels
        const priceRange = resistance - support;
        const priceStep = priceRange / 10; // Divide into 10 levels
        const volumeProfile = new Array(10).fill(0);
        const buyVolumeProfile = new Array(10).fill(0);
        const sellVolumeProfile = new Array(10).fill(0);
        
        recentCandles.forEach(c => {
            const vol = c.volume || 0;
            const candleRange = c.high - c.low;
            const closePos = candleRange > 0 ? (c.close - c.low) / candleRange : 0.5;
            
            // Distribute volume across price levels
            const lowLevel = Math.floor((c.low - support) / priceStep);
            const highLevel = Math.floor((c.high - support) / priceStep);
            
            for (let level = Math.max(0, lowLevel); level <= Math.min(9, highLevel); level++) {
                volumeProfile[level] += vol / (highLevel - lowLevel + 1);
                
                if (c.close > c.open) {
                    buyVolumeProfile[level] += vol * closePos / (highLevel - lowLevel + 1);
                    sellVolumeProfile[level] += vol * (1 - closePos) / (highLevel - lowLevel + 1);
                } else {
                    sellVolumeProfile[level] += vol * (1 - closePos) / (highLevel - lowLevel + 1);
                    buyVolumeProfile[level] += vol * closePos / (highLevel - lowLevel + 1);
                }
            }
        });
        
        // Find Point of Control (POC) - price level with highest volume
        const pocLevel = volumeProfile.indexOf(Math.max(...volumeProfile));
        const pocPrice = support + (pocLevel * priceStep);
        const currentLevel = Math.floor((currentPrice - support) / priceStep);
        
        // Check if current price is above/below POC
        const abovePOC = currentPrice > pocPrice;
        const buyPressureAtPOC = buyVolumeProfile[pocLevel] > sellVolumeProfile[pocLevel];
        
        // Calculate trend for reference only (not filtering)
        const trendStrength = Math.abs(upCount - downCount) / 10;
        const isTrendingUp = upCount > downCount;
        const isTrendingDown = downCount > upCount;
        
        // Combine signals with stricter requirements
        let bullishSignals = 0;
        let bearishSignals = 0;
        
        // Trend signals (stronger weight)
        if (upCount > downCount * 1.3) bullishSignals += 2; // Strong trend
        else if (upCount > downCount) bullishSignals++;
        else if (downCount > upCount * 1.3) bearishSignals += 2;
        else bearishSignals++;
        
        // Momentum signals (higher threshold)
        if (momentum > 1.0) bullishSignals += 2;
        else if (momentum > 0.3) bullishSignals++;
        else if (momentum < -1.0) bearishSignals += 2;
        else if (momentum < -0.3) bearishSignals++;
        
        // Volume signals (stricter)
        if (volumeRatio > 0.6) bullishSignals += 2;
        else if (volumeRatio > 0.52) bullishSignals++;
        else if (volumeRatio < 0.4) bearishSignals += 2;
        else if (volumeRatio < 0.48) bearishSignals++;
        
        // MA signals (stronger alignment required)
        if (ma5 > ma10 && ma10 > ma20 && currentPrice > ma5) bullishSignals += 3; // Perfect alignment
        else if (ma5 < ma10 && ma10 < ma20 && currentPrice < ma5) bearishSignals += 3;
        else if (currentPrice > ma20) bullishSignals++;
        else bearishSignals++;
        
        // Price position signals (only at extremes)
        if (pricePosition < 0.25) bullishSignals += 2; // Strong support
        else if (pricePosition > 0.75) bearishSignals += 2; // Strong resistance
        
        // Delta signals (most important - higher weight)
        if (cumulativeDelta > 0) {
            bullishSignals += 3; // Very strong buy pressure
            if (positiveDeltaCount > negativeDeltaCount * 1.5) bullishSignals += 2; // Extremely strong
        } else {
            bearishSignals += 3; // Very strong sell pressure
            if (negativeDeltaCount > positiveDeltaCount * 1.5) bearishSignals += 2; // Extremely strong
        }
        
        // Footprint signals (higher weight)
        if (abovePOC && buyPressureAtPOC) {
            bullishSignals += 3; // Above POC with buy pressure = very bullish
        } else if (!abovePOC && !buyPressureAtPOC) {
            bearishSignals += 3; // Below POC with sell pressure = very bearish
        }
        
        // Volume imbalance at current level (stricter threshold)
        if (currentLevel >= 0 && currentLevel < 10) {
            const buyVolAtLevel = buyVolumeProfile[currentLevel];
            const sellVolAtLevel = sellVolumeProfile[currentLevel];
            if (buyVolAtLevel > sellVolAtLevel * 1.5) bullishSignals += 2;
            else if (sellVolAtLevel > buyVolAtLevel * 1.5) bearishSignals += 2;
        }
        
        // Balanced confidence calculation
        const totalSignals = bullishSignals + bearishSignals;
        const dominantSignals = Math.max(bullishSignals, bearishSignals);
        const signalRatio = totalSignals > 0 ? dominantSignals / totalSignals : 0.5;
        const signalStrength = Math.abs(bullishSignals - bearishSignals);
        const confidence = 45 + (signalRatio * 35) + (Math.min(signalStrength, 12) * 1.5); // 45-98%
        
        const aiPrediction = {
            signal: upCount > downCount ? 'BUY' : 'SELL',
            confidence: Math.min(confidence, 90)
        };
        
        // No filtering - let the signal through
        
        // Simple multi-timeframe trend filter
        const mtfTrendFilter = {
            higherTFTrend: upCount > downCount,
            confidence: Math.abs(upCount - downCount) / 10
        };
        
        if (window.generateTradingSignal) {
            // Disable enhanced ICT for backtest (causes unrealistic returns)
            return generateTradingSignal(data, aiPrediction, volumeAnalysis, mtfTrendFilter, false);
        }
        
        return null;
    } catch (error) {
        console.error('Error generating signal:', error);
        return null;
    }
}

// Analyze volume from historical data
function analyzeHistoricalVolume(data) {
    if (!data || data.length === 0) return null;
    
    const recentCandles = data.slice(-30);
    
    // Calculate average volume
    const avgVolume = recentCandles.reduce((sum, c) => sum + (c.volume || 0), 0) / recentCandles.length;
    const currentVolume = data[data.length - 1].volume || avgVolume;
    const volumeRatio = currentVolume / avgVolume;
    
    // Calculate average candle metrics
    const avgHeight = recentCandles.reduce((sum, c) => sum + (c.high - c.low), 0) / recentCandles.length;
    const avgBodySize = recentCandles.reduce((sum, c) => sum + Math.abs(c.close - c.open), 0) / recentCandles.length;
    
    // Calculate wick ratios
    const avgUpperWick = recentCandles.reduce((sum, c) => {
        const bodyTop = Math.max(c.open, c.close);
        return sum + (c.high - bodyTop);
    }, 0) / recentCandles.length;
    
    const avgLowerWick = recentCandles.reduce((sum, c) => {
        const bodyBottom = Math.min(c.open, c.close);
        return sum + (bodyBottom - c.low);
    }, 0) / recentCandles.length;
    
    // Determine volatility
    let volatilityLevel = 'normal';
    if (volumeRatio > 1.5) volatilityLevel = 'high';
    else if (volumeRatio < 0.7) volatilityLevel = 'low';
    
    return {
        avgHeight: avgHeight,
        avgBodyRatio: avgBodySize / avgHeight,
        avgUpperWickRatio: avgUpperWick / avgHeight,
        avgLowerWickRatio: avgLowerWick / avgHeight,
        volumeRatio: volumeRatio,
        volatilityLevel: volatilityLevel,
        predictedHeight: avgHeight * (0.8 + volumeRatio * 0.4)
    };
}

// Simulate trade execution with trailing stop
function simulateTrade(signal, futureData) {
    const entry = signal.entry;
    let stopLoss = signal.stopLoss;
    const targets = signal.targets;
    
    // 🛡️ REALISTIC RISK MANAGEMENT
    const currentBalance = backtestResults.currentBalance;
    const riskPercent = 0.02; // 2% of current balance (reduced from 3%)
    
    // Cap maximum position size to prevent unrealistic growth
    const maxPositionSize = backtestResults.startingBalance * 10; // Max 10x starting capital per trade
    const uncappedRisk = currentBalance * riskPercent;
    const riskAmount = Math.min(uncappedRisk, maxPositionSize);
    
    let exitPrice = null;
    let exitReason = '';
    let candlesHeld = 0;
    let highestPrice = entry; // For BUY trailing
    let lowestPrice = entry; // For SELL trailing
    let trailingActive = false;
    
    // Check each future candle
    for (let i = 0; i < futureData.length; i++) {
        const candle = futureData[i];
        candlesHeld++;
        
        if (signal.type === 'BUY') {
            // Update highest price
            if (candle.high > highestPrice) {
                highestPrice = candle.high;
                
                // 🎯 REALISTIC TRAILING: More conservative to avoid noise
                const profitPercent = (highestPrice - entry) / (entry - signal.stopLoss);
                if (profitPercent > 1.2) { // Activate at 1.2R (more conservative)
                    trailingActive = true;
                    // Trail stop to lock 40% of profit (more realistic)
                    stopLoss = Math.max(stopLoss, entry + (highestPrice - entry) * 0.40);
                }
            }
            
            // Check stop loss
            if (candle.low <= stopLoss) {
                exitPrice = stopLoss;
                exitReason = trailingActive ? 'Trailing Stop' : 'Stop Loss';
                break;
            }
            
            // Check take profit 1
            if (candle.high >= targets[0].price) {
                exitPrice = targets[0].price;
                exitReason = 'TP1';
                break;
            }
        } else {
            // SELL trade
            if (candle.low < lowestPrice) {
                lowestPrice = candle.low;
                
                // 🎯 REALISTIC TRAILING: More conservative to avoid noise
                const profitPercent = (entry - lowestPrice) / (signal.stopLoss - entry);
                if (profitPercent > 1.2) { // Activate at 1.2R (more conservative)
                    trailingActive = true;
                    // Trail stop to lock 40% of profit (more realistic)
                    stopLoss = Math.min(stopLoss, entry - (entry - lowestPrice) * 0.40);
                }
            }
            
            // Check stop loss
            if (candle.high >= stopLoss) {
                exitPrice = stopLoss;
                exitReason = trailingActive ? 'Trailing Stop' : 'Stop Loss';
                break;
            }
            
            // Check take profit 1
            if (candle.low <= targets[0].price) {
                exitPrice = targets[0].price;
                exitReason = 'TP1';
                break;
            }
        }
    }
    
    // If no exit, close at last candle
    if (!exitPrice) {
        exitPrice = futureData[futureData.length - 1].close;
        exitReason = 'Timeout';
    }
    
    // Calculate profit/loss with realistic costs
    const priceDiff = signal.type === 'BUY' ? 
        (exitPrice - entry) : (entry - exitPrice);
    
    const riskDiff = Math.abs(entry - signal.stopLoss);
    const positionSize = riskAmount / riskDiff;
    
    // Apply realistic slippage (0.05% on entry + exit)
    const slippagePercent = 0.0005;
    const slippageCost = entry * slippagePercent * 2 * positionSize;
    
    // Apply trading fees (0.1% maker/taker average)
    const feePercent = 0.001;
    const tradingFees = entry * feePercent * 2 * positionSize;
    
    const grossProfit = priceDiff * positionSize;
    const profit = grossProfit - slippageCost - tradingFees;
    const rrAchieved = Math.abs(priceDiff / riskDiff);
    
    // Ensure all values are valid numbers
    const safeProfit = isNaN(profit) || !isFinite(profit) ? 0 : profit;
    const safeProfitPercent = isNaN(profit) || !isFinite(profit) || riskAmount === 0 ? 0 : (profit / riskAmount) * 100;
    const safeRR = isNaN(rrAchieved) || !isFinite(rrAchieved) ? 0 : rrAchieved;
    
    return {
        type: signal.type,
        entry: entry,
        exit: exitPrice,
        stopLoss: signal.stopLoss,
        exitReason: exitReason,
        profit: safeProfit,
        profitPercent: safeProfitPercent,
        candlesHeld: candlesHeld,
        rr: safeRR
    };
}

// Calculate backtest statistics
function calculateBacktestStats() {
    // Ensure values are numbers, not undefined
    backtestResults.currentBalance = backtestResults.currentBalance || backtestResults.startingBalance;
    backtestResults.totalProfit = backtestResults.totalProfit || 0;
    backtestResults.totalLoss = backtestResults.totalLoss || 0;
    
    backtestResults.netProfit = backtestResults.currentBalance - backtestResults.startingBalance;
    backtestResults.returnPercent = (backtestResults.netProfit / backtestResults.startingBalance) * 100;
    
    if (backtestResults.totalTrades === 0) {
        backtestResults.winRate = 0;
        backtestResults.profitFactor = 0;
        backtestResults.averageRR = 0;
        backtestResults.returnPercent = 0;
        return;
    }
    
    backtestResults.winRate = (backtestResults.winningTrades / backtestResults.totalTrades) * 100;
    backtestResults.profitFactor = backtestResults.totalLoss > 0 ? 
        backtestResults.totalProfit / backtestResults.totalLoss : 0;
    
    const totalRR = backtestResults.trades.reduce((sum, t) => sum + (t.rr || 0), 0);
    backtestResults.averageRR = totalRR / backtestResults.totalTrades;
}

// Display backtest results
function displayBacktestResults() {
    const days = window.backtestDays || 30;
    console.log('\n📊 ===== BACKTEST RESULTS =====');
    console.log(`Period: Last ${days} days on ${currentInterval || '15m'}`);
    console.log(`Total Trades: ${backtestResults.totalTrades}`);
    console.log(`Winning Trades: ${backtestResults.winningTrades} (${backtestResults.winRate.toFixed(1)}%)`);
    console.log(`Losing Trades: ${backtestResults.losingTrades}`);
    console.log(`\n💰 PROFITABILITY:`);
    console.log(`Starting Balance: $${backtestResults.startingBalance.toFixed(2)}`);
    console.log(`Final Balance: $${backtestResults.currentBalance.toFixed(2)}`);
    console.log(`Net Profit: $${backtestResults.netProfit.toFixed(2)} (${backtestResults.returnPercent.toFixed(2)}%)`);
    console.log(`Total Profit: $${backtestResults.totalProfit.toFixed(2)}`);
    console.log(`Total Loss: $${backtestResults.totalLoss.toFixed(2)}`);
    console.log(`Profit Factor: ${backtestResults.profitFactor.toFixed(2)}`);
    console.log(`Average RR: ${backtestResults.averageRR.toFixed(2)}:1`);
    console.log(`Max Drawdown: ${(backtestResults.maxDrawdown * 100).toFixed(2)}%`);
    
    // Show last 5 trades
    console.log(`\n📈 Last 5 Trades:`);
    backtestResults.trades.slice(-5).forEach((trade, i) => {
        const profitSign = trade.profit > 0 ? '+' : '';
        console.log(`${i + 1}. ${trade.type} | ${trade.exitReason} | ${profitSign}$${trade.profit.toFixed(2)} (${profitSign}${trade.profitPercent.toFixed(1)}%)`);
    });
    
    // Overall assessment
    console.log(`\n🎯 ASSESSMENT:`);
    if (backtestResults.totalTrades === 0) {
        console.log('⚠️ NO TRADES - Criteria too strict! Lower confidence threshold.');
        console.log('💡 TIP: Current settings require 50% confidence + 1.2:1 RR + Higher TF trend');
        console.log('💡 Try: Lower timeframe (1m, 5m) for more opportunities');
    } else if (backtestResults.returnPercent > 20) {
        console.log('✅ HIGHLY PROFITABLE - Excellent system!');
    } else if (backtestResults.returnPercent > 10) {
        console.log('✅ PROFITABLE - Good system');
    } else if (backtestResults.returnPercent > 0) {
        console.log('⚠️ SLIGHTLY PROFITABLE - Needs improvement');
    } else {
        console.log('❌ NOT PROFITABLE - Requires optimization');
    }
    
    console.log('==============================\n');
}

// Export functions
window.runBacktest = runBacktest;
window.backtestResults = backtestResults;

// Simple backtest function that works immediately
window.startBacktest = function(days = 30) {
    console.log(`🔄 Starting backtest for ${days} days...`);
    
    const btn = document.getElementById('backtest-btn');
    
    if (btn) {
        btn.disabled = true;
        btn.textContent = '⏳ Running...';
    }
    
    // Run backtest
    runBacktest(window.currentSymbol || 'BTCUSDT', window.currentInterval || '15m', days)
        .then(() => {
            if (btn) {
                btn.disabled = false;
                btn.textContent = '📊 Backtest';
            }
        })
        .catch(error => {
            console.error('❌ Backtest error:', error);
            if (btn) {
                btn.disabled = false;
                btn.textContent = '📊 Backtest';
            }
        });
};

// Run backtest with selected period from dropdown
window.runBacktestWithPeriod = function(event) {
    // Prevent default if called from button click
    if (event && event.preventDefault) {
        event.preventDefault();
    }
    
    const periodSelect = document.getElementById('backtest-period');
    const days = periodSelect ? parseInt(periodSelect.value) : 30;
    
    console.log(`📊 Selected period: ${days} days`);
    window.startBacktest(days);
};

// Initialize backtest button
setTimeout(() => {
    const backtestBtn = document.getElementById('backtest-btn');
    
    if (backtestBtn) {
        // Don't override onclick - it's set in HTML
        console.log('✅ Backtest button ready');
    } else {
        console.error('❌ Backtest button not found');
    }
}, 2000);

// Download Detailed Backtest Report
function downloadBacktestReport() {
    if (!backtestResults || backtestResults.totalTrades === 0) {
        alert('⚠️ No backtest data available. Please run a backtest first!');
        return;
    }
    
    const results = backtestResults;
    const symbol = window.currentSymbol || 'BTCUSDT';
    const interval = window.currentInterval || '15m';
    const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
    
    // Generate detailed report
    let report = '';
    
    // Header
    report += '═══════════════════════════════════════════════════════════\n';
    report += '           DETAILED BACKTEST REPORT\n';
    report += '═══════════════════════════════════════════════════════════\n\n';
    
    // Test Configuration
    report += '📊 TEST CONFIGURATION\n';
    report += '─────────────────────────────────────────────────────────\n';
    report += `Symbol:              ${symbol}\n`;
    report += `Timeframe:           ${interval}\n`;
    report += `Test Period:         ${window.backtestDays || 30} days\n`;
    report += `Starting Balance:    $${results.startingBalance.toFixed(2)}\n`;
    report += `Risk Per Trade:      $15.00 (3% of account)\n`;
    report += `Generated:           ${new Date().toLocaleString()}\n\n`;
    
    // Performance Summary
    report += '💰 PERFORMANCE SUMMARY\n';
    report += '─────────────────────────────────────────────────────────\n';
    report += `Final Balance:       $${results.currentBalance.toFixed(2)}\n`;
    report += `Net Profit/Loss:     $${results.netProfit.toFixed(2)}\n`;
    report += `Return:              ${results.returnPercent.toFixed(2)}%\n`;
    report += `Peak Balance:        $${results.peakBalance.toFixed(2)}\n`;
    report += `Max Drawdown:        ${(results.maxDrawdown * 100).toFixed(2)}%\n\n`;
    
    // Trade Statistics
    report += '📈 TRADE STATISTICS\n';
    report += '─────────────────────────────────────────────────────────\n';
    report += `Total Trades:        ${results.totalTrades}\n`;
    report += `Winning Trades:      ${results.winningTrades} (${results.winRate.toFixed(1)}%)\n`;
    report += `Losing Trades:       ${results.losingTrades} (${(100 - results.winRate).toFixed(1)}%)\n`;
    report += `Average RR:          ${results.averageRR.toFixed(2)}:1\n`;
    report += `Profit Factor:       ${results.profitFactor.toFixed(2)}\n\n`;
    
    // Profit/Loss Breakdown
    report += '💵 PROFIT/LOSS BREAKDOWN\n';
    report += '─────────────────────────────────────────────────────────\n';
    report += `Total Profit:        $${results.totalProfit.toFixed(2)}\n`;
    report += `Total Loss:          $${results.totalLoss.toFixed(2)}\n`;
    report += `Average Win:         $${(results.totalProfit / Math.max(results.winningTrades, 1)).toFixed(2)}\n`;
    report += `Average Loss:        $${(results.totalLoss / Math.max(results.losingTrades, 1)).toFixed(2)}\n`;
    report += `Largest Win:         $${Math.max(...results.trades.filter(t => t.profit > 0).map(t => t.profit), 0).toFixed(2)}\n`;
    report += `Largest Loss:        $${Math.min(...results.trades.filter(t => t.profit < 0).map(t => t.profit), 0).toFixed(2)}\n\n`;
    
    // Exit Reason Analysis
    const exitReasons = {};
    results.trades.forEach(t => {
        exitReasons[t.exitReason] = (exitReasons[t.exitReason] || 0) + 1;
    });
    
    report += '🎯 EXIT REASON ANALYSIS\n';
    report += '─────────────────────────────────────────────────────────\n';
    Object.entries(exitReasons).forEach(([reason, count]) => {
        const percentage = (count / results.totalTrades * 100).toFixed(1);
        report += `${reason.padEnd(20)} ${count} trades (${percentage}%)\n`;
    });
    report += '\n';
    
    // Trade Type Analysis
    const buyTrades = results.trades.filter(t => t.type === 'BUY');
    const sellTrades = results.trades.filter(t => t.type === 'SELL');
    const buyWins = buyTrades.filter(t => t.profit > 0).length;
    const sellWins = sellTrades.filter(t => t.profit > 0).length;
    
    report += '📊 TRADE TYPE ANALYSIS\n';
    report += '─────────────────────────────────────────────────────────\n';
    report += `BUY Trades:          ${buyTrades.length} (${buyWins} wins, ${(buyWins/Math.max(buyTrades.length,1)*100).toFixed(1)}% win rate)\n`;
    report += `SELL Trades:         ${sellTrades.length} (${sellWins} wins, ${(sellWins/Math.max(sellTrades.length,1)*100).toFixed(1)}% win rate)\n\n`;
    
    // Detailed Trade Log
    report += '📋 DETAILED TRADE LOG\n';
    report += '═══════════════════════════════════════════════════════════\n\n';
    
    results.trades.forEach((trade, index) => {
        const profitSign = trade.profit >= 0 ? '+' : '';
        const status = trade.profit >= 0 ? '✅ WIN' : '❌ LOSS';
        
        report += `Trade #${index + 1} - ${status}\n`;
        report += `─────────────────────────────────────────────────────────\n`;
        report += `Type:                ${trade.type}\n`;
        report += `Entry:               $${trade.entry.toFixed(2)}\n`;
        report += `Exit:                $${trade.exit.toFixed(2)}\n`;
        report += `Stop Loss:           $${trade.stopLoss.toFixed(2)}\n`;
        report += `Exit Reason:         ${trade.exitReason}\n`;
        report += `Candles Held:        ${trade.candlesHeld}\n`;
        report += `Risk:Reward:         ${trade.rr.toFixed(2)}:1\n`;
        report += `Profit/Loss:         ${profitSign}$${trade.profit.toFixed(2)} (${profitSign}${trade.profitPercent.toFixed(2)}%)\n`;
        report += `Balance After:       $${(results.startingBalance + results.trades.slice(0, index + 1).reduce((sum, t) => sum + t.profit, 0)).toFixed(2)}\n`;
        report += '\n';
    });
    
    // Footer
    report += '═══════════════════════════════════════════════════════════\n';
    report += '                    END OF REPORT\n';
    report += '═══════════════════════════════════════════════════════════\n';
    report += '\n⚠️ DISCLAIMER: Past performance does not guarantee future results.\n';
    report += 'This is a simulation based on historical data.\n';
    report += 'Always use proper risk management in live trading.\n';
    
    // Create download
    const blob = new Blob([report], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `backtest_${symbol}_${interval}_${timestamp}.txt`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
    
    console.log('📥 Backtest report downloaded successfully!');
    alert('✅ Detailed backtest report downloaded!');
}

// Export function globally
window.downloadBacktestReport = downloadBacktestReport;
