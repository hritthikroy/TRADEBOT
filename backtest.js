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
    startingBalance: 500,
    currentBalance: 500
};

// Run backtest on historical data
async function runBacktest(symbol, interval, days = 30) {
    console.log(`🔄 Starting backtest: ${symbol} ${interval} for ${days} days...`);
    
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
        startingBalance: 500,
        currentBalance: 500,
        peakBalance: 500
    };
    
    try {
        // Fetch historical data
        const limit = calculateCandleLimit(interval, days);
        const url = `https://api.binance.com/api/v3/klines?symbol=${symbol}&interval=${interval}&limit=${limit}`;
        
        const response = await fetch(url);
        const data = await response.json();
        
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
        '4h': 6
    };
    
    const limit = (candlesPerDay[interval] || 96) * days;
    return Math.min(limit, 1000); // Binance limit
}

// Generate signal from historical data
async function generateHistoricalSignal(data, interval) {
    try {
        // Use same logic as live trading
        const volumeAnalysis = analyzeHistoricalVolume(data);
        
        if (!volumeAnalysis) return null;
        
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
        
        // Simple multi-timeframe trend filter
        const mtfTrendFilter = {
            higherTFTrend: upCount > downCount,
            confidence: Math.abs(upCount - downCount) / 10
        };
        
        if (window.generateTradingSignal) {
            return generateTradingSignal(data, aiPrediction, volumeAnalysis, mtfTrendFilter);
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
    const riskAmount = 100; // Risk $100 per trade (1% of $10k account)
    
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
                
                // Activate trailing stop after 1.0R profit (optimized)
                const profitPercent = (highestPrice - entry) / (entry - signal.stopLoss);
                if (profitPercent > 1.0) {
                    trailingActive = true;
                    // Trail stop to breakeven + 50% of profit (more aggressive)
                    stopLoss = Math.max(stopLoss, entry + (highestPrice - entry) * 0.50);
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
                
                // Activate trailing stop after 1.0R profit (optimized)
                const profitPercent = (entry - lowestPrice) / (signal.stopLoss - entry);
                if (profitPercent > 1.0) {
                    trailingActive = true;
                    // Trail stop to breakeven + 50% of profit (more aggressive)
                    stopLoss = Math.min(stopLoss, entry - (entry - lowestPrice) * 0.50);
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
    
    // Calculate profit/loss
    const priceDiff = signal.type === 'BUY' ? 
        (exitPrice - entry) : (entry - exitPrice);
    
    const riskDiff = Math.abs(entry - signal.stopLoss);
    const positionSize = riskAmount / riskDiff;
    const profit = priceDiff * positionSize;
    const rrAchieved = Math.abs(priceDiff / riskDiff);
    
    return {
        type: signal.type,
        entry: entry,
        exit: exitPrice,
        stopLoss: signal.stopLoss,
        exitReason: exitReason,
        profit: profit,
        profitPercent: (profit / riskAmount) * 100,
        candlesHeld: candlesHeld,
        rr: rrAchieved
    };
}

// Calculate backtest statistics
function calculateBacktestStats() {
    backtestResults.netProfit = backtestResults.currentBalance - backtestResults.startingBalance;
    backtestResults.returnPercent = (backtestResults.netProfit / backtestResults.startingBalance) * 100;
    
    if (backtestResults.totalTrades === 0) {
        backtestResults.winRate = 0;
        backtestResults.profitFactor = 0;
        backtestResults.averageRR = 0;
        return;
    }
    
    backtestResults.winRate = (backtestResults.winningTrades / backtestResults.totalTrades) * 100;
    backtestResults.profitFactor = backtestResults.totalLoss > 0 ? 
        backtestResults.totalProfit / backtestResults.totalLoss : 0;
    
    const totalRR = backtestResults.trades.reduce((sum, t) => sum + t.rr, 0);
    backtestResults.averageRR = totalRR / backtestResults.totalTrades;
}

// Display backtest results
function displayBacktestResults() {
    console.log('\n📊 ===== BACKTEST RESULTS =====');
    console.log(`Period: Last 30 days on ${currentInterval || '15m'}`);
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

// Visualize backtest on prediction chart
function visualizeBacktest() {
    console.log('🎬 Replay button clicked!');
    
    if (!window.backtestHistoricalData || window.backtestHistoricalData.length === 0) {
        console.error('❌ No backtest data available. Run backtest first.');
        alert('Please run backtest first!');
        return;
    }
    
    if (!window.backtestTrades || window.backtestTrades.length === 0) {
        console.warn('⚠️ No trades to visualize');
        alert('No trades found in backtest results!');
        return;
    }
    
    console.log('📊 Starting backtest visualization...');
    console.log(`Total trades to visualize: ${window.backtestTrades.length}`);
    
    // Stop live updates
    if (window.isSimulating) {
        window.isSimulating = false;
        console.log('⏸️ Stopped live updates');
    }
    
    // Check if required functions exist
    if (!window.currentData) {
        console.error('❌ window.currentData not found');
        alert('Chart data not available. Please refresh the page.');
        return;
    }
    
    if (!window.drawChart) {
        console.error('❌ window.drawChart function not found');
        alert('Chart drawing function not available. Please refresh the page.');
        return;
    }
    
    console.log('✅ Chart functions available, starting replay...');
    
    // Replace current data with backtest data
    if (window.currentData) {
        // Start from where trades begin (after the initial window used for analysis)
        const startIndex = 50; // Start from beginning after initial analysis window
        
        // Update the global currentData array (don't reassign, modify in place)
        currentData.length = 0; // Clear array
        const initialData = window.backtestHistoricalData.slice(0, startIndex);
        initialData.forEach(candle => currentData.push(candle));
        
        // Also update window reference
        window.currentData = currentData;
        window.predictedCandles = [];
        window.backtestVisualizationMode = true;
        window.backtestCurrentTrades = [];
        
        console.log(`📊 Starting from beginning, showing first ${startIndex} candles`);
        console.log(`Total candles to replay: ${window.backtestHistoricalData.length - startIndex}`);
        console.log(`📊 Total trades in backtest: ${window.backtestTrades.length}`);
        
        // Draw initial chart
        try {
            window.drawChart();
            console.log('✅ Initial chart drawn');
        } catch (error) {
            console.error('❌ Error drawing chart:', error);
            alert('Error drawing chart: ' + error.message);
            return;
        }
        
        updateStatus('Starting backtest replay...');
        
        // Start animation
        let currentIndex = startIndex;
        let tradesFound = 0;
        console.log(`🎬 Starting animation from index ${currentIndex}`);
        console.log(`📊 Looking for ${window.backtestTrades.length} trades from index ${startIndex} to ${window.backtestHistoricalData.length}`);
        
        // Log all trade entry indices for debugging
        if (window.backtestTrades.length > 0) {
            const tradeIndices = window.backtestTrades.map(t => t.entryIndex).slice(0, 10);
            console.log(`📍 First 10 trade indices: ${tradeIndices.join(', ')}`);
            console.log(`📍 Min trade index: ${Math.min(...window.backtestTrades.map(t => t.entryIndex))}`);
            console.log(`📍 Max trade index: ${Math.max(...window.backtestTrades.map(t => t.entryIndex))}`);
        } else {
            console.error('❌ No trades in backtestTrades array!');
        }
        
        const animationInterval = setInterval(() => {
            try {
                if (currentIndex >= window.backtestHistoricalData.length) {
                    clearInterval(animationInterval);
                    console.log('✅ Backtest visualization complete');
                    console.log(`📊 Final Results: ${backtestResults.returnPercent.toFixed(1)}% return | ${backtestResults.totalTrades} trades | ${backtestResults.winRate.toFixed(1)}% win rate`);
                    updateStatus(`Replay complete: ${backtestResults.returnPercent.toFixed(1)}% return`);
                    window.backtestVisualizationMode = false;
                    return;
                }
            
                // Add next candle
                currentData.push(window.backtestHistoricalData[currentIndex]);
                
                // Keep only last 50 candles visible
                if (currentData.length > 50) {
                    currentData.shift();
                }
            
                // Debug: Log every 100 candles
                if (currentIndex % 100 === 0) {
                    console.log(`📊 Progress: Index ${currentIndex}, Trades found so far: ${tradesFound}`);
                }
            
                // Check if there's a trade at this index
                const trade = window.backtestTrades.find(t => t.entryIndex === currentIndex);
                if (trade) {
                    tradesFound++;
                    const tradeNum = window.backtestTrades.indexOf(trade) + 1;
                const profitSign = trade.profit > 0 ? '+' : '';
                const emoji = trade.profit > 0 ? '✅' : '❌';
                console.log(`${emoji} Trade ${tradeNum}/${window.backtestTrades.length}: ${trade.type} @ ${trade.entry.toFixed(2)} → ${trade.exitReason} @ ${trade.exit.toFixed(2)} | P/L: ${profitSign}$${trade.profit.toFixed(2)} (${profitSign}${trade.profitPercent.toFixed(1)}%)`);
                
                // Store active trade for visualization
                window.backtestCurrentTrades.push({
                    ...trade,
                    candleIndex: window.currentData.length - 1
                });
            }
            
                // Redraw chart with trade markers
                window.drawChart();
                
                currentIndex++;
                const progress = ((currentIndex - startIndex) / (window.backtestHistoricalData.length - startIndex) * 100).toFixed(0);
                updateStatus(`Replaying: ${progress}% | Balance: $${backtestResults.currentBalance.toFixed(0)}`);
                
            } catch (error) {
                console.error('❌ Error in animation loop:', error);
                clearInterval(animationInterval);
                alert('Animation error: ' + error.message);
            }
        }, 50); // 50ms per candle (faster replay)
        
        console.log('✅ Animation interval started');
    } else {
        console.error('❌ window.currentData is not available');
        alert('Chart data not available!');
    }
}

// Export functions
window.runBacktest = runBacktest;
window.backtestResults = backtestResults;

// Simple backtest function that works immediately
window.startBacktest = function() {
    console.log('🔄 Starting backtest...');
    
    const btn = document.getElementById('backtest-btn');
    
    if (btn) {
        btn.disabled = true;
        btn.textContent = '⏳ Running...';
    }
    
    // Run backtest
    runBacktest(window.currentSymbol || 'BTCUSDT', window.currentInterval || '15m', 30)
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

// Initialize backtest button
setTimeout(() => {
    const backtestBtn = document.getElementById('backtest-btn');
    
    if (backtestBtn) {
        backtestBtn.onclick = window.startBacktest;
        console.log('✅ Backtest button ready');
    } else {
        console.error('❌ Backtest button not found');
    }
}, 2000);
