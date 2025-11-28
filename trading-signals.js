// Advanced Trading Signals with Risk-Reward Analysis
// Provides precise entry, stop loss, and take profit levels

// Calculate support and resistance levels
function calculateSupportResistance(data) {
    if (data.length < 20) return null;
    
    const recent = data.slice(-50);
    const highs = recent.map(c => c.high);
    const lows = recent.map(c => c.low);
    
    // Find pivot points
    const pivotHigh = Math.max(...highs);
    const pivotLow = Math.min(...lows);
    const pivotClose = recent[recent.length - 1].close;
    
    // Calculate pivot levels
    const pivot = (pivotHigh + pivotLow + pivotClose) / 3;
    const r1 = (2 * pivot) - pivotLow;
    const r2 = pivot + (pivotHigh - pivotLow);
    const r3 = pivotHigh + 2 * (pivot - pivotLow);
    const s1 = (2 * pivot) - pivotHigh;
    const s2 = pivot - (pivotHigh - pivotLow);
    const s3 = pivotLow - 2 * (pivotHigh - pivot);
    
    return {
        resistance: [r1, r2, r3],
        support: [s1, s2, s3],
        pivot: pivot,
        current: pivotClose
    };
}

// Calculate ATR (Average True Range) for volatility
function calculateATR(data, period = 14) {
    if (data.length < period + 1) return null;
    
    const trueRanges = [];
    
    for (let i = 1; i < data.length; i++) {
        const high = data[i].high;
        const low = data[i].low;
        const prevClose = data[i - 1].close;
        
        const tr = Math.max(
            high - low,
            Math.abs(high - prevClose),
            Math.abs(low - prevClose)
        );
        
        trueRanges.push(tr);
    }
    
    const recentTR = trueRanges.slice(-period);
    const atr = recentTR.reduce((sum, tr) => sum + tr, 0) / period;
    
    return atr;
}

// Calculate Fibonacci retracement levels
function calculateFibonacci(data) {
    if (data.length < 20) return null;
    
    const recent = data.slice(-50);
    const high = Math.max(...recent.map(c => c.high));
    const low = Math.min(...recent.map(c => c.low));
    const diff = high - low;
    
    return {
        level_0: high,
        level_236: high - (diff * 0.236),
        level_382: high - (diff * 0.382),
        level_500: high - (diff * 0.500),
        level_618: high - (diff * 0.618),
        level_786: high - (diff * 0.786),
        level_100: low
    };
}

// Generate precise trading signal with entry, SL, TP
function generateTradingSignal(data, aiPrediction, volumeAnalysis, mtfTrendFilter) {
    if (data.length < 20) return null;
    
    const currentPrice = data[data.length - 1].close;
    const atr = calculateATR(data);
    const srLevels = calculateSupportResistance(data);
    const fibonacci = calculateFibonacci(data);
    
    if (!atr || !srLevels) return null;
    
    // Pattern analysis disabled for now (causing issues in backtest)
    let patternAnalysis = null;
    
    // Determine signal strength with pattern boost
    let signalStrength = aiPrediction ? aiPrediction.confidence : 50;
    let isBullish = aiPrediction && aiPrediction.signal === 'BUY';
    
    // Pattern boost disabled (was causing negative returns)
    
    // Skip NEUTRAL signals
    if (!aiPrediction || aiPrediction.signal === 'NEUTRAL') {
        console.log('Skipping NEUTRAL signal');
        return null;
    }
    
    // Multi-timeframe trend filter - prefer alignment
    if (mtfTrendFilter && mtfTrendFilter.higherTFTrend !== null) {
        const higherTFBullish = mtfTrendFilter.higherTFTrend;
        
        // Reject if strongly against higher timeframe trend
        if (isBullish !== higherTFBullish && mtfTrendFilter.confidence > 0.75) {
            console.log(`❌ Signal rejected: Against strong higher TF trend (${higherTFBullish ? 'UP' : 'DOWN'}, ${(mtfTrendFilter.confidence * 100).toFixed(0)}%)`);
            return null;
        }
        
        // Boost confidence if aligned with higher TF
        if (isBullish === higherTFBullish && mtfTrendFilter.confidence > 0.65) {
            signalStrength = Math.min(signalStrength + 15, 100);
            console.log(`✅ Signal boosted: Aligned with higher TF trend (+15%)`);
        }
    }
    
    let signal = null;
    
    // ICT/SMC ANALYSIS
    // 1. Identify Order Blocks (last strong move before reversal)
    const orderBlocks = findOrderBlocks(data);
    
    // 2. Find Fair Value Gaps (imbalance zones)
    const fvg = findFairValueGaps(data);
    
    // 3. Check Break of Structure
    const bos = detectBreakOfStructure(data);
    
    // 4. Calculate cumulative delta
    const deltaInfo = calculateCumulativeDelta(data);
    
    // 5. Check if price is at key level (S/R retest)
    const atKeyLevel = checkKeyLevelRetest(currentPrice, srLevels, atr);
    
    // 6. ADVANCED: Detect Liquidity Sweeps
    const liquiditySweep = detectLiquiditySweep(data);
    
    // 7. ADVANCED: Find Breaker Blocks
    const breakerBlocks = findBreakerBlocks(data, orderBlocks);
    
    // 8. ADVANCED: Detect Displacement
    const displacement = detectDisplacement(data, atr);
    
    // 9. ADVANCED: Find Liquidity Voids
    const liquidityVoids = findLiquidityVoids(data);
    
    // 10. ADVANCED: Power of 3 (PO3) - Market cycle detection
    const po3 = detectPowerOf3(data);
    
    // 11. ADVANCED: AMD (Accumulation, Manipulation, Distribution)
    const amd = detectAMD(data, atr);
    
    console.log(`ICT Check: ${isBullish ? 'BULLISH' : 'BEARISH'} | Strength: ${signalStrength}% | Delta: ${deltaInfo.delta > 0 ? '+' : ''}${deltaInfo.delta.toFixed(0)} | PO3: ${po3.phase} | AMD: ${amd.phase} | Sweep: ${liquiditySweep.bullish ? 'BULL' : liquiditySweep.bearish ? 'BEAR' : 'NONE'} | Breaker: ${breakerBlocks.bullish.length}/${breakerBlocks.bearish.length}`);
    
    // ICT/SMC BUY SETUP
    if (isBullish && signalStrength >= 55) {
        const entry = currentPrice;
        
        // Check for bullish confluence
        let confluence = 0;
        
        // Delta confirmation (most important)
        if (deltaInfo.delta > 0) {
            if (deltaInfo.strength > 0.7) confluence += 3;
            else if (deltaInfo.strength > 0.5) confluence += 2;
            else confluence += 1;
        }
        
        // Order block support (very important)
        if (orderBlocks.bullish.length > 0) {
            const nearOB = orderBlocks.bullish.some(ob => 
                currentPrice >= ob.low * 0.98 && currentPrice <= ob.high * 1.03
            );
            if (nearOB) confluence += 4;
        }
        
        // Breaker block support (VERY STRONG - failed resistance)
        if (breakerBlocks.bearish.length > 0) {
            const nearBreaker = breakerBlocks.bearish.some(bb => 
                currentPrice >= bb.low * 0.98 && currentPrice <= bb.high * 1.03
            );
            if (nearBreaker) confluence += 5; // Strongest signal
        }
        
        // Fair Value Gap (important)
        if (fvg.bullish.length > 0) {
            const inFVG = fvg.bullish.some(gap => 
                currentPrice >= gap.low * 0.98 && currentPrice <= gap.high * 1.02
            );
            if (inFVG) confluence += 3;
        }
        
        // Liquidity Sweep (strong reversal signal)
        if (liquiditySweep.bullish) confluence += 4;
        
        // Displacement (institutional move)
        if (displacement.bullish && displacement.strength > 2.5) confluence += 3;
        
        // Power of 3 - Distribution phase (best for buys after manipulation)
        if (po3.phase === 'DISTRIBUTION' && po3.direction === 'BULLISH') confluence += 4;
        
        // AMD - Distribution phase (smart money buying)
        if (amd.phase === 'DISTRIBUTION' && amd.direction === 'BULLISH') confluence += 3;
        
        // Break of Structure (confirmation)
        if (bos.type === 'BULLISH') confluence += 2;
        
        // Key level retest (strong)
        if (atKeyLevel.type === 'SUPPORT') confluence += 3;
        
        console.log(`BUY Confluence: ${confluence}/38 points | PO3: ${po3.phase} | AMD: ${amd.phase} | Sweep: ${liquiditySweep.bullish}`);
        
        // Require minimum confluence (raised for quality)
        if (confluence >= 8) {
            // Stop below order block or recent low
            let stopLoss;
            if (orderBlocks.bullish.length > 0) {
                const nearestOB = orderBlocks.bullish[orderBlocks.bullish.length - 1];
                stopLoss = nearestOB.low - (atr * 0.3);
            } else {
                const recentLows = data.slice(-15).map(c => c.low);
                stopLoss = Math.min(...recentLows) - (atr * 0.5);
            }
            
            // Target next liquidity zone (optimized for better RR)
            const takeProfit1 = currentPrice + (atr * 2.5);
            const takeProfit2 = currentPrice + (atr * 4.5);
            const takeProfit3 = currentPrice + (atr * 7.0);
            
            const risk = entry - stopLoss;
            const reward1 = takeProfit1 - entry;
            const reward2 = takeProfit2 - entry;
            const reward3 = takeProfit3 - entry;
            
            const rr1 = reward1 / risk;
            const rr2 = reward2 / risk;
            const rr3 = reward3 / risk;
            
            if (rr1 >= 1.5) {
                signal = {
                    type: 'BUY',
                    strength: signalStrength,
                    entry: entry,
                    stopLoss: stopLoss,
                    targets: [
                        { price: takeProfit1, rr: rr1, percentage: 40 },
                        { price: takeProfit2, rr: rr2, percentage: 30 },
                        { price: takeProfit3, rr: rr3, percentage: 30 }
                    ],
                    risk: risk,
                    riskPercent: (risk / entry) * 100,
                    bestRR: Math.max(rr1, rr2, rr3),
                    atr: atr,
                    support: srLevels.support[0],
                    resistance: srLevels.resistance[0]
                };
            }
        }
    // ICT/SMC SELL SETUP
    } else if (!isBullish && signalStrength >= 55) {
        const entry = currentPrice;
        
        // Check for bearish confluence
        let confluence = 0;
        
        // Delta confirmation (most important)
        if (deltaInfo.delta < 0) {
            if (deltaInfo.strength > 0.7) confluence += 3;
            else if (deltaInfo.strength > 0.5) confluence += 2;
            else confluence += 1;
        }
        
        // Order block resistance (very important)
        if (orderBlocks.bearish.length > 0) {
            const nearOB = orderBlocks.bearish.some(ob => 
                currentPrice <= ob.high * 1.02 && currentPrice >= ob.low * 0.97
            );
            if (nearOB) confluence += 4;
        }
        
        // Breaker block resistance (VERY STRONG - failed support)
        if (breakerBlocks.bullish.length > 0) {
            const nearBreaker = breakerBlocks.bullish.some(bb => 
                currentPrice <= bb.high * 1.02 && currentPrice >= bb.low * 0.97
            );
            if (nearBreaker) confluence += 5; // Strongest signal
        }
        
        // Fair Value Gap (important)
        if (fvg.bearish.length > 0) {
            const inFVG = fvg.bearish.some(gap => 
                currentPrice >= gap.low * 0.98 && currentPrice <= gap.high * 1.02
            );
            if (inFVG) confluence += 3;
        }
        
        // Liquidity Sweep (strong reversal signal)
        if (liquiditySweep.bearish) confluence += 4;
        
        // Displacement (institutional move)
        if (displacement.bearish && displacement.strength > 2.5) confluence += 3;
        
        // Power of 3 - Distribution phase (best for sells after manipulation)
        if (po3.phase === 'DISTRIBUTION' && po3.direction === 'BEARISH') confluence += 4;
        
        // AMD - Distribution phase (smart money selling)
        if (amd.phase === 'DISTRIBUTION' && amd.direction === 'BEARISH') confluence += 3;
        
        // Break of Structure (confirmation)
        if (bos.type === 'BEARISH') confluence += 2;
        
        // Key level retest (strong)
        if (atKeyLevel.type === 'RESISTANCE') confluence += 3;
        
        console.log(`SELL Confluence: ${confluence}/38 points | PO3: ${po3.phase} | AMD: ${amd.phase} | Sweep: ${liquiditySweep.bearish}`);
        
        // Require minimum confluence (raised for quality)
        if (confluence >= 8) {
            // Stop above order block or recent high
            let stopLoss;
            if (orderBlocks.bearish.length > 0) {
                const nearestOB = orderBlocks.bearish[orderBlocks.bearish.length - 1];
                stopLoss = nearestOB.high + (atr * 0.3);
            } else {
                const recentHighs = data.slice(-15).map(c => c.high);
                stopLoss = Math.max(...recentHighs) + (atr * 0.5);
            }
            
            // Target next liquidity zone (optimized for better RR)
            const takeProfit1 = currentPrice - (atr * 2.5);
            const takeProfit2 = currentPrice - (atr * 4.5);
            const takeProfit3 = currentPrice - (atr * 7.0);
            
            const risk = stopLoss - entry;
            const reward1 = entry - takeProfit1;
            const reward2 = entry - takeProfit2;
            const reward3 = entry - takeProfit3;
            
            const rr1 = reward1 / risk;
            const rr2 = reward2 / risk;
            const rr3 = reward3 / risk;
            
            if (rr1 >= 1.5) {
                signal = {
                    type: 'SELL',
                    strength: signalStrength,
                    entry: entry,
                    stopLoss: stopLoss,
                    targets: [
                        { price: takeProfit1, rr: rr1, percentage: 40 },
                        { price: takeProfit2, rr: rr2, percentage: 30 },
                        { price: takeProfit3, rr: rr3, percentage: 30 }
                    ],
                    risk: risk,
                    riskPercent: (risk / entry) * 100,
                    bestRR: Math.max(rr1, rr2, rr3),
                    atr: atr,
                    support: srLevels.support[0],
                    resistance: srLevels.resistance[0]
                };
            }
        }
    }
    
    return signal;
}

// Display trading signal on UI
// Draw trading signal directly on chart
function drawTradingSignalOnChart(signal, leftPadding, topPadding, chartHeight, maxPrice, minPrice, padding) {
    if (!signal || !ctx) return;
    
    const priceToY = (price) => {
        return topPadding + chartHeight - ((price - (minPrice - padding)) / (maxPrice + padding - (minPrice - padding))) * chartHeight;
    };
    
    const isBuy = signal.type === 'BUY';
    const color = isBuy ? '#26a69a' : '#ef5350';
    
    const chartWidth = canvas.width - leftPadding - 150; // Leave space for labels
    
    // Draw entry line with price label
    const entryY = priceToY(signal.entry);
    ctx.strokeStyle = color;
    ctx.lineWidth = 3;
    ctx.setLineDash([]);
    ctx.beginPath();
    ctx.moveTo(leftPadding, entryY);
    ctx.lineTo(leftPadding + chartWidth, entryY);
    ctx.stroke();
    ctx.setLineDash([]);
    
    // Entry price label (highlighted)
    ctx.fillStyle = color;
    ctx.fillRect(leftPadding + chartWidth + 5, entryY - 12, 90, 24);
    ctx.fillStyle = '#ffffff';
    ctx.font = 'bold 12px Arial';
    ctx.textAlign = 'left';
    ctx.fillText(`ENTRY`, leftPadding + chartWidth + 10, entryY - 1);
    ctx.fillText(`${signal.entry.toFixed(2)}`, leftPadding + chartWidth + 10, entryY + 11);
    
    // Draw stop loss line with price label
    const slY = priceToY(signal.stopLoss);
    ctx.strokeStyle = '#ef5350';
    ctx.lineWidth = 2;
    ctx.setLineDash([5, 5]);
    ctx.beginPath();
    ctx.moveTo(leftPadding, slY);
    ctx.lineTo(leftPadding + chartWidth, slY);
    ctx.stroke();
    ctx.setLineDash([]);
    
    // Stop loss price label (highlighted)
    ctx.fillStyle = '#ef5350';
    ctx.fillRect(leftPadding + chartWidth + 5, slY - 12, 90, 24);
    ctx.fillStyle = '#ffffff';
    ctx.font = 'bold 11px Arial';
    ctx.textAlign = 'left';
    ctx.fillText(`STOP LOSS`, leftPadding + chartWidth + 10, slY - 1);
    ctx.fillText(`${signal.stopLoss.toFixed(2)}`, leftPadding + chartWidth + 10, slY + 11);
    
    // Draw take profit lines with price labels
    signal.targets.forEach((target, i) => {
        const tpY = priceToY(target.price);
        const tpColor = i === 0 ? '#26a69a' : i === 1 ? '#4caf50' : '#66bb6a';
        
        ctx.strokeStyle = tpColor;
        ctx.lineWidth = 2;
        ctx.setLineDash([5, 5]);
        ctx.beginPath();
        ctx.moveTo(leftPadding, tpY);
        ctx.lineTo(leftPadding + chartWidth, tpY);
        ctx.stroke();
        ctx.setLineDash([]);
        
        // TP price label (highlighted)
        ctx.fillStyle = tpColor;
        ctx.fillRect(leftPadding + chartWidth + 5, tpY - 12, 90, 24);
        ctx.fillStyle = '#ffffff';
        ctx.font = 'bold 11px Arial';
        ctx.textAlign = 'left';
        ctx.fillText(`TP${i+1} (${target.rr.toFixed(1)}R)`, leftPadding + chartWidth + 10, tpY - 1);
        ctx.fillText(`${target.price.toFixed(2)}`, leftPadding + chartWidth + 10, tpY + 11);
    });
    
    // Draw signal box
    const boxX = leftPadding + 10;
    const boxY = topPadding + 10;
    const boxWidth = 180;
    const boxHeight = 80;
    
    // Box background
    ctx.fillStyle = 'rgba(30, 34, 45, 0.9)';
    ctx.fillRect(boxX, boxY, boxWidth, boxHeight);
    ctx.strokeStyle = color;
    ctx.lineWidth = 2;
    ctx.strokeRect(boxX, boxY, boxWidth, boxHeight);
    
    // Signal text
    ctx.fillStyle = color;
    ctx.font = 'bold 16px Arial';
    ctx.textAlign = 'left';
    ctx.fillText(`${isBuy ? '📈' : '📉'} ${signal.type}`, boxX + 10, boxY + 25);
    
    ctx.fillStyle = '#d1d4dc';
    ctx.font = '12px Arial';
    ctx.fillText(`Entry: $${signal.entry.toFixed(2)}`, boxX + 10, boxY + 45);
    ctx.fillText(`SL: $${signal.stopLoss.toFixed(2)}`, boxX + 10, boxY + 60);
    ctx.fillText(`RR: ${signal.bestRR.toFixed(1)}:1 | ${signal.strength}%`, boxX + 10, boxY + 75);
}

function displayTradingSignal(signal) {
    // Store signal globally for chart drawing
    window.currentTradingSignal = signal;
    
    if (!signal) {
        console.log('⏳ Waiting for Quality Setup (Confidence ≥60% + RR ≥1.5:1 + Higher TF Check)');
        return;
    }
    
    const isBuy = signal.type === 'BUY';
    const color = isBuy ? '#26a69a' : '#ef5350';
    const arrow = isBuy ? '📈' : '📉';
    
    let html = `
        <div style="background: ${color}22; border-left: 4px solid ${color}; padding: 15px; border-radius: 4px; margin-top: 10px;">
            <div style="font-size: 20px; font-weight: bold; color: ${color}; margin-bottom: 10px;">
                ${arrow} ${signal.type} SIGNAL
                <span style="font-size: 14px; color: #787b86; margin-left: 10px;">
                    Strength: ${signal.strength}%
                </span>
            </div>
            
            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 10px; margin-top: 10px;">
                <div>
                    <div style="color: #787b86; font-size: 11px;">ENTRY</div>
                    <div style="color: ${color}; font-size: 16px; font-weight: bold;">
                        $${signal.entry.toFixed(2)}
                    </div>
                </div>
                
                <div>
                    <div style="color: #787b86; font-size: 11px;">STOP LOSS</div>
                    <div style="color: #ef5350; font-size: 16px; font-weight: bold;">
                        $${signal.stopLoss.toFixed(2)}
                    </div>
                </div>
            </div>
            
            <div style="margin-top: 10px; padding-top: 10px; border-top: 1px solid #2a2e39;">
                <div style="color: #787b86; font-size: 11px; margin-bottom: 5px;">TAKE PROFIT TARGETS</div>
                ${signal.targets.map((target, i) => `
                    <div style="display: flex; justify-content: space-between; margin: 5px 0; font-size: 13px;">
                        <span style="color: #d1d4dc;">TP${i + 1} (${target.percentage}%)</span>
                        <span style="color: #26a69a; font-weight: bold;">
                            $${target.price.toFixed(2)}
                            <span style="color: #787b86; font-size: 11px; margin-left: 5px;">
                                RR: ${target.rr.toFixed(2)}:1
                            </span>
                        </span>
                    </div>
                `).join('')}
            </div>
            
            <div style="margin-top: 10px; padding-top: 10px; border-top: 1px solid #2a2e39; display: flex; justify-content: space-between; font-size: 12px;">
                <div>
                    <span style="color: #787b86;">Risk:</span>
                    <span style="color: #ef5350; font-weight: bold;">${signal.riskPercent.toFixed(2)}%</span>
                </div>
                <div>
                    <span style="color: #787b86;">Best RR:</span>
                    <span style="color: #26a69a; font-weight: bold;">${signal.bestRR.toFixed(2)}:1</span>
                </div>
                <div>
                    <span style="color: #787b86;">ATR:</span>
                    <span style="color: #2962ff; font-weight: bold;">$${signal.atr.toFixed(2)}</span>
                </div>
            </div>
        </div>
    `;
    
    updateSignalDisplay(signal.type, html, signal);
}

function updateSignalDisplay(type, content, signal) {
    const signalDiv = document.getElementById('trading-signal-display');
    
    if (signalDiv) {
        signalDiv.innerHTML = content;
        console.log('✅ Signal display updated');
    } else {
        console.error('❌ Signal display element not found');
    }
    
    // Log to console
    if (signal) {
        console.log('🎯 TRADING SIGNAL:', signal);
    } else {
        console.log('⏳ Waiting for signal...');
    }
}

// ICT/SMC Helper Functions

// Find Order Blocks (last candle before strong move)
function findOrderBlocks(data) {
    const bullishOB = [];
    const bearishOB = [];
    
    for (let i = 3; i < data.length - 1; i++) {
        const prev = data[i - 1];
        const curr = data[i];
        const next = data[i + 1];
        
        // Bullish OB: Down candle followed by strong up move
        if (curr.close < curr.open && next.close > next.open) {
            const moveSize = next.close - next.open;
            if (moveSize > (curr.open - curr.close) * 2) {
                bullishOB.push({
                    index: i,
                    high: curr.high,
                    low: curr.low,
                    strength: moveSize
                });
            }
        }
        
        // Bearish OB: Up candle followed by strong down move
        if (curr.close > curr.open && next.close < next.open) {
            const moveSize = next.open - next.close;
            if (moveSize > (curr.close - curr.open) * 2) {
                bearishOB.push({
                    index: i,
                    high: curr.high,
                    low: curr.low,
                    strength: moveSize
                });
            }
        }
    }
    
    return {
        bullish: bullishOB.slice(-3), // Keep last 3
        bearish: bearishOB.slice(-3)
    };
}

// Find Fair Value Gaps (price imbalances)
function findFairValueGaps(data) {
    const bullishFVG = [];
    const bearishFVG = [];
    
    for (let i = 2; i < data.length; i++) {
        const candle1 = data[i - 2];
        const candle2 = data[i - 1];
        const candle3 = data[i];
        
        // Bullish FVG: Gap between candle1 high and candle3 low
        if (candle3.low > candle1.high) {
            bullishFVG.push({
                low: candle1.high,
                high: candle3.low,
                index: i
            });
        }
        
        // Bearish FVG: Gap between candle1 low and candle3 high
        if (candle3.high < candle1.low) {
            bearishFVG.push({
                low: candle3.high,
                high: candle1.low,
                index: i
            });
        }
    }
    
    return {
        bullish: bullishFVG.slice(-2), // Keep last 2
        bearish: bearishFVG.slice(-2)
    };
}

// Detect Break of Structure
function detectBreakOfStructure(data) {
    if (data.length < 20) return { type: 'NONE' };
    
    const recent = data.slice(-20);
    const highs = recent.map(c => c.high);
    const lows = recent.map(c => c.low);
    
    const recentHigh = Math.max(...highs.slice(-10));
    const prevHigh = Math.max(...highs.slice(0, 10));
    const recentLow = Math.min(...lows.slice(-10));
    const prevLow = Math.min(...lows.slice(0, 10));
    
    // Bullish BOS: Break above previous high
    if (recentHigh > prevHigh * 1.005) {
        return { type: 'BULLISH', strength: (recentHigh - prevHigh) / prevHigh };
    }
    
    // Bearish BOS: Break below previous low
    if (recentLow < prevLow * 0.995) {
        return { type: 'BEARISH', strength: (prevLow - recentLow) / prevLow };
    }
    
    return { type: 'NONE' };
}

// Calculate Cumulative Delta
function calculateCumulativeDelta(data) {
    const recent = data.slice(-10);
    let cumulativeDelta = 0;
    let totalVolume = 0;
    
    recent.forEach(c => {
        const range = c.high - c.low;
        const closePos = range > 0 ? (c.close - c.low) / range : 0.5;
        const vol = c.volume || 0;
        totalVolume += vol;
        
        let buyVol, sellVol;
        if (c.close > c.open) {
            buyVol = vol * (0.5 + closePos * 0.5);
            sellVol = vol - buyVol;
        } else {
            sellVol = vol * (0.5 + (1 - closePos) * 0.5);
            buyVol = vol - sellVol;
        }
        
        cumulativeDelta += (buyVol - sellVol);
    });
    
    const strength = totalVolume > 0 ? Math.abs(cumulativeDelta) / totalVolume : 0;
    
    return { delta: cumulativeDelta, strength: strength };
}

// Check if price is retesting key level
function checkKeyLevelRetest(price, srLevels, atr) {
    const tolerance = atr * 0.5;
    
    // Check support retest
    for (let support of srLevels.support) {
        if (Math.abs(price - support) < tolerance) {
            return { type: 'SUPPORT', level: support };
        }
    }
    
    // Check resistance retest
    for (let resistance of srLevels.resistance) {
        if (Math.abs(price - resistance) < tolerance) {
            return { type: 'RESISTANCE', level: resistance };
        }
    }
    
    return { type: 'NONE' };
}

// Detect Liquidity Sweeps (stop hunts before reversal)
function detectLiquiditySweep(data) {
    if (data.length < 15) return { bullish: false, bearish: false };
    
    const recent = data.slice(-15);
    const last3 = recent.slice(-3);
    
    // Find recent swing high/low
    const swingHigh = Math.max(...recent.slice(0, 12).map(c => c.high));
    const swingLow = Math.min(...recent.slice(0, 12).map(c => c.low));
    
    // Bullish sweep: Price breaks below swing low then reverses up
    const brokeBelow = last3.some(c => c.low < swingLow);
    const reversedUp = last3[last3.length - 1].close > last3[0].close;
    const bullishSweep = brokeBelow && reversedUp;
    
    // Bearish sweep: Price breaks above swing high then reverses down
    const brokeAbove = last3.some(c => c.high > swingHigh);
    const reversedDown = last3[last3.length - 1].close < last3[0].close;
    const bearishSweep = brokeAbove && reversedDown;
    
    return { 
        bullish: bullishSweep, 
        bearish: bearishSweep,
        swingHigh: swingHigh,
        swingLow: swingLow
    };
}

// Find Breaker Blocks (failed OB that becomes opposite OB)
function findBreakerBlocks(data, orderBlocks) {
    const breakerBlocks = { bullish: [], bearish: [] };
    
    if (data.length < 10) return breakerBlocks;
    
    const currentPrice = data[data.length - 1].close;
    
    // Check if bullish OB was broken (becomes bearish breaker)
    orderBlocks.bullish.forEach(ob => {
        const obIndex = ob.index;
        if (obIndex < data.length - 5) {
            // Check if price broke below OB and is now above it
            const brokeBelowOB = data.slice(obIndex + 1).some(c => c.close < ob.low);
            const nowAboveOB = currentPrice > ob.high;
            
            if (brokeBelowOB && nowAboveOB) {
                breakerBlocks.bullish.push({
                    high: ob.high,
                    low: ob.low,
                    strength: 'STRONG' // Failed support becomes resistance
                });
            }
        }
    });
    
    // Check if bearish OB was broken (becomes bullish breaker)
    orderBlocks.bearish.forEach(ob => {
        const obIndex = ob.index;
        if (obIndex < data.length - 5) {
            // Check if price broke above OB and is now below it
            const brokeAboveOB = data.slice(obIndex + 1).some(c => c.close > ob.high);
            const nowBelowOB = currentPrice < ob.low;
            
            if (brokeAboveOB && nowBelowOB) {
                breakerBlocks.bearish.push({
                    high: ob.high,
                    low: ob.low,
                    strength: 'STRONG' // Failed resistance becomes support
                });
            }
        }
    });
    
    return breakerBlocks;
}

// Detect Displacement (strong institutional move)
function detectDisplacement(data, atr) {
    if (data.length < 5) return { bullish: false, bearish: false, strength: 0 };
    
    const last5 = data.slice(-5);
    const avgCandleSize = last5.slice(0, 4).reduce((sum, c) => sum + Math.abs(c.close - c.open), 0) / 4;
    
    const lastCandle = last5[last5.length - 1];
    const lastCandleSize = Math.abs(lastCandle.close - lastCandle.open);
    
    // Displacement = candle 2-3x larger than average
    const isDisplacement = lastCandleSize > avgCandleSize * 2;
    const strength = lastCandleSize / avgCandleSize;
    
    return {
        bullish: isDisplacement && lastCandle.close > lastCandle.open,
        bearish: isDisplacement && lastCandle.close < lastCandle.open,
        strength: strength
    };
}

// Find Liquidity Voids (gaps with no trading)
function findLiquidityVoids(data) {
    const voids = { bullish: [], bearish: [] };
    
    for (let i = 1; i < data.length; i++) {
        const prev = data[i - 1];
        const curr = data[i];
        
        // Bullish void: Gap up (curr.low > prev.high)
        if (curr.low > prev.high) {
            const gapSize = curr.low - prev.high;
            voids.bullish.push({
                top: curr.low,
                bottom: prev.high,
                size: gapSize,
                index: i
            });
        }
        
        // Bearish void: Gap down (curr.high < prev.low)
        if (curr.high < prev.low) {
            const gapSize = prev.low - curr.high;
            voids.bearish.push({
                top: prev.low,
                bottom: curr.high,
                size: gapSize,
                index: i
            });
        }
    }
    
    return {
        bullish: voids.bullish.slice(-2),
        bearish: voids.bearish.slice(-2)
    };
}

// Detect Power of 3 (PO3) - Market Cycle
// Phase 1: Accumulation (consolidation)
// Phase 2: Manipulation (liquidity grab/fake move)
// Phase 3: Distribution (real move)
function detectPowerOf3(data) {
    if (data.length < 20) return { phase: 'UNKNOWN', direction: null };
    
    const recent = data.slice(-20);
    const last10 = recent.slice(-10);
    const prev10 = recent.slice(0, 10);
    
    // Calculate range and volatility
    const last10High = Math.max(...last10.map(c => c.high));
    const last10Low = Math.min(...last10.map(c => c.low));
    const last10Range = last10High - last10Low;
    
    const prev10High = Math.max(...prev10.map(c => c.high));
    const prev10Low = Math.min(...prev10.map(c => c.low));
    const prev10Range = prev10High - prev10Low;
    
    // Phase 1: ACCUMULATION - Low volatility, tight range
    const isAccumulation = last10Range < prev10Range * 0.7;
    
    // Phase 2: MANIPULATION - Quick spike then reversal
    const last3 = last10.slice(-3);
    const spikedUp = last3[0].high > prev10High && last3[2].close < last3[0].close;
    const spikedDown = last3[0].low < prev10Low && last3[2].close > last3[0].close;
    const isManipulation = spikedUp || spikedDown;
    
    // Phase 3: DISTRIBUTION - Strong directional move
    const strongUpMove = last10[last10.length - 1].close > prev10High * 1.01;
    const strongDownMove = last10[last10.length - 1].close < prev10Low * 0.99;
    const isDistribution = strongUpMove || strongDownMove;
    
    let phase = 'UNKNOWN';
    let direction = null;
    
    if (isDistribution) {
        phase = 'DISTRIBUTION';
        direction = strongUpMove ? 'BULLISH' : 'BEARISH';
    } else if (isManipulation) {
        phase = 'MANIPULATION';
        direction = spikedUp ? 'BEARISH' : 'BULLISH'; // Opposite of spike
    } else if (isAccumulation) {
        phase = 'ACCUMULATION';
    }
    
    return { phase, direction };
}

// Detect AMD (Accumulation, Manipulation, Distribution)
// Similar to PO3 but focuses on volume and delta
function detectAMD(data, atr) {
    if (data.length < 15) return { phase: 'UNKNOWN', direction: null };
    
    const recent = data.slice(-15);
    const last5 = recent.slice(-5);
    const prev10 = recent.slice(0, 10);
    
    // Calculate volume profile
    const avgVolume = prev10.reduce((sum, c) => sum + (c.volume || 0), 0) / prev10.length;
    const last5Volume = last5.reduce((sum, c) => sum + (c.volume || 0), 0) / last5.length;
    
    // Calculate delta for last 5 candles
    let cumulativeDelta = 0;
    last5.forEach(c => {
        const range = c.high - c.low;
        const closePos = range > 0 ? (c.close - c.low) / range : 0.5;
        const vol = c.volume || 0;
        
        let buyVol, sellVol;
        if (c.close > c.open) {
            buyVol = vol * (0.5 + closePos * 0.5);
            sellVol = vol - buyVol;
        } else {
            sellVol = vol * (0.5 + (1 - closePos) * 0.5);
            buyVol = vol - sellVol;
        }
        cumulativeDelta += (buyVol - sellVol);
    });
    
    // Calculate price movement
    const priceChange = last5[last5.length - 1].close - last5[0].close;
    const priceChangePercent = Math.abs(priceChange / last5[0].close) * 100;
    
    // ACCUMULATION: Low volume, tight range
    const isAccumulation = last5Volume < avgVolume * 0.8 && priceChangePercent < 0.5;
    
    // MANIPULATION: High volume spike with reversal
    const hasVolumeSpike = last5Volume > avgVolume * 1.5;
    const hasReversal = (cumulativeDelta > 0 && priceChange < 0) || (cumulativeDelta < 0 && priceChange > 0);
    const isManipulation = hasVolumeSpike && hasReversal;
    
    // DISTRIBUTION: High volume with strong directional move
    const strongMove = priceChangePercent > 1.0;
    const deltaConfirms = (cumulativeDelta > 0 && priceChange > 0) || (cumulativeDelta < 0 && priceChange < 0);
    const isDistribution = hasVolumeSpike && strongMove && deltaConfirms;
    
    let phase = 'UNKNOWN';
    let direction = null;
    
    if (isDistribution) {
        phase = 'DISTRIBUTION';
        direction = priceChange > 0 ? 'BULLISH' : 'BEARISH';
    } else if (isManipulation) {
        phase = 'MANIPULATION';
        direction = cumulativeDelta > 0 ? 'BULLISH' : 'BEARISH';
    } else if (isAccumulation) {
        phase = 'ACCUMULATION';
    }
    
    return { phase, direction, delta: cumulativeDelta };
}

// Export functions
window.generateTradingSignal = generateTradingSignal;
window.displayTradingSignal = displayTradingSignal;
window.drawTradingSignalOnChart = drawTradingSignalOnChart;
window.calculateATR = calculateATR;
window.calculateSupportResistance = calculateSupportResistance;
