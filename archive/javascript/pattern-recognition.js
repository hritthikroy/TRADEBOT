// Advanced Pattern Recognition for Trading Signals
// Detects candlestick patterns, retests, and reversals

// Detect candlestick patterns
function detectCandlePatterns(candles) {
    if (candles.length < 3) return [];
    
    const patterns = [];
    const current = candles[candles.length - 1];
    const prev = candles[candles.length - 2];
    const prev2 = candles[candles.length - 3];
    
    // Calculate candle properties
    const currentBody = Math.abs(current.close - current.open);
    const currentRange = current.high - current.low;
    const currentUpperWick = current.high - Math.max(current.open, current.close);
    const currentLowerWick = Math.min(current.open, current.close) - current.low;
    
    const prevBody = Math.abs(prev.close - prev.open);
    const prevRange = prev.high - prev.low;
    
    // 1. HAMMER (Bullish Reversal)
    if (currentLowerWick > currentBody * 2 && 
        currentUpperWick < currentBody * 0.3 &&
        current.close < prev.close) {
        patterns.push({
            name: 'HAMMER',
            type: 'BULLISH_REVERSAL',
            strength: 75,
            description: 'Strong buying pressure at lows'
        });
    }
    
    // 2. SHOOTING STAR (Bearish Reversal)
    if (currentUpperWick > currentBody * 2 && 
        currentLowerWick < currentBody * 0.3 &&
        current.close > prev.close) {
        patterns.push({
            name: 'SHOOTING_STAR',
            type: 'BEARISH_REVERSAL',
            strength: 75,
            description: 'Strong selling pressure at highs'
        });
    }
    
    // 3. ENGULFING BULLISH
    if (current.close > current.open && // Current is bullish
        prev.close < prev.open && // Previous is bearish
        current.open < prev.close && // Opens below prev close
        current.close > prev.open) { // Closes above prev open
        patterns.push({
            name: 'BULLISH_ENGULFING',
            type: 'BULLISH_REVERSAL',
            strength: 85,
            description: 'Bulls overwhelm bears'
        });
    }
    
    // 4. ENGULFING BEARISH
    if (current.close < current.open && // Current is bearish
        prev.close > prev.open && // Previous is bullish
        current.open > prev.close && // Opens above prev close
        current.close < prev.open) { // Closes below prev open
        patterns.push({
            name: 'BEARISH_ENGULFING',
            type: 'BEARISH_REVERSAL',
            strength: 85,
            description: 'Bears overwhelm bulls'
        });
    }
    
    // 5. MORNING STAR (Bullish Reversal - 3 candles)
    if (prev2.close < prev2.open && // First is bearish
        Math.abs(prev.close - prev.open) < prevBody * 0.3 && // Middle is doji/small
        current.close > current.open && // Last is bullish
        current.close > (prev2.open + prev2.close) / 2) { // Closes above midpoint
        patterns.push({
            name: 'MORNING_STAR',
            type: 'BULLISH_REVERSAL',
            strength: 90,
            description: 'Strong 3-candle reversal'
        });
    }
    
    // 6. EVENING STAR (Bearish Reversal - 3 candles)
    if (prev2.close > prev2.open && // First is bullish
        Math.abs(prev.close - prev.open) < prevBody * 0.3 && // Middle is doji/small
        current.close < current.open && // Last is bearish
        current.close < (prev2.open + prev2.close) / 2) { // Closes below midpoint
        patterns.push({
            name: 'EVENING_STAR',
            type: 'BEARISH_REVERSAL',
            strength: 90,
            description: 'Strong 3-candle reversal'
        });
    }
    
    // 7. DOJI (Indecision)
    if (currentBody < currentRange * 0.1) {
        patterns.push({
            name: 'DOJI',
            type: 'INDECISION',
            strength: 60,
            description: 'Market indecision, potential reversal'
        });
    }
    
    // 8. MARUBOZU (Strong Trend)
    if (currentBody > currentRange * 0.9) {
        const type = current.close > current.open ? 'BULLISH_CONTINUATION' : 'BEARISH_CONTINUATION';
        patterns.push({
            name: 'MARUBOZU',
            type: type,
            strength: 80,
            description: 'Strong directional momentum'
        });
    }
    
    return patterns;
}

// Detect support/resistance retests
function detectRetests(candles, srLevels) {
    if (!srLevels || candles.length < 10) return [];
    
    const retests = [];
    const current = candles[candles.length - 1];
    const recent = candles.slice(-10);
    
    // Check for support retest
    srLevels.support.forEach((supportLevel, index) => {
        // Check if price recently touched support
        const touchedSupport = recent.some(c => 
            Math.abs(c.low - supportLevel) / supportLevel < 0.003 // Within 0.3%
        );
        
        // Check if price is bouncing from support
        const bouncing = current.close > current.low && 
                        Math.abs(current.low - supportLevel) / supportLevel < 0.005;
        
        if (touchedSupport && bouncing) {
            retests.push({
                type: 'SUPPORT_RETEST',
                level: supportLevel,
                strength: 80 - (index * 10), // S1 stronger than S2
                description: `Support retest at ${supportLevel.toFixed(2)}`,
                signal: 'BULLISH'
            });
        }
    });
    
    // Check for resistance retest
    srLevels.resistance.forEach((resistanceLevel, index) => {
        // Check if price recently touched resistance
        const touchedResistance = recent.some(c => 
            Math.abs(c.high - resistanceLevel) / resistanceLevel < 0.003
        );
        
        // Check if price is rejecting from resistance
        const rejecting = current.close < current.high && 
                         Math.abs(current.high - resistanceLevel) / resistanceLevel < 0.005;
        
        if (touchedResistance && rejecting) {
            retests.push({
                type: 'RESISTANCE_RETEST',
                level: resistanceLevel,
                strength: 80 - (index * 10),
                description: `Resistance retest at ${resistanceLevel.toFixed(2)}`,
                signal: 'BEARISH'
            });
        }
    });
    
    return retests;
}

// Detect reversal zones
function detectReversalZones(candles) {
    if (candles.length < 20) return [];
    
    const reversals = [];
    const recent = candles.slice(-20);
    
    // Find swing highs and lows
    for (let i = 2; i < recent.length - 2; i++) {
        const candle = recent[i];
        
        // Swing High (potential resistance)
        if (candle.high > recent[i-1].high && 
            candle.high > recent[i-2].high &&
            candle.high > recent[i+1].high && 
            candle.high > recent[i+2].high) {
            
            reversals.push({
                type: 'SWING_HIGH',
                price: candle.high,
                index: i,
                strength: 70,
                description: 'Previous swing high - resistance zone'
            });
        }
        
        // Swing Low (potential support)
        if (candle.low < recent[i-1].low && 
            candle.low < recent[i-2].low &&
            candle.low < recent[i+1].low && 
            candle.low < recent[i+2].low) {
            
            reversals.push({
                type: 'SWING_LOW',
                price: candle.low,
                index: i,
                strength: 70,
                description: 'Previous swing low - support zone'
            });
        }
    }
    
    return reversals;
}

// Combine all pattern analysis
function analyzePatterns(candles, srLevels) {
    const candlePatterns = detectCandlePatterns(candles);
    const retests = detectRetests(candles, srLevels);
    const reversalZones = detectReversalZones(candles);
    
    // Calculate overall signal strength
    let bullishScore = 0;
    let bearishScore = 0;
    
    candlePatterns.forEach(pattern => {
        if (pattern.type.includes('BULLISH')) {
            bullishScore += pattern.strength;
        } else if (pattern.type.includes('BEARISH')) {
            bearishScore += pattern.strength;
        }
    });
    
    retests.forEach(retest => {
        if (retest.signal === 'BULLISH') {
            bullishScore += retest.strength;
        } else if (retest.signal === 'BEARISH') {
            bearishScore += retest.strength;
        }
    });
    
    return {
        candlePatterns: candlePatterns,
        retests: retests,
        reversalZones: reversalZones,
        bullishScore: bullishScore,
        bearishScore: bearishScore,
        signal: bullishScore > bearishScore ? 'BULLISH' : bearishScore > bullishScore ? 'BEARISH' : 'NEUTRAL',
        confidence: Math.abs(bullishScore - bearishScore) / Math.max(bullishScore, bearishScore, 1) * 100
    };
}

// Export functions
window.detectCandlePatterns = detectCandlePatterns;
window.detectRetests = detectRetests;
window.detectReversalZones = detectReversalZones;
window.analyzePatterns = analyzePatterns;
