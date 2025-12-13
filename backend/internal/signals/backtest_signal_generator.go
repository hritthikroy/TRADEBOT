package signals

import (
	"math"
)

// generateBacktestSignal - ADVANCED MULTI-FACTOR STRATEGY
// Combines: MTF Analysis, Candlestick Patterns, ICT/SMC, Order Flow,
// Delta Analysis, Footprint, Volume Profile, Enhanced Order Blocks
func generateBacktestSignal(data []Candle, interval string) *Signal {
	if len(data) < 50 {
		return nil
	}

	// Use high-probability trend following strategy
	return generateTrendFollowingSignal(data, interval)
}

// generateTrendFollowingSignal - High probability trend following with advanced filters
func generateTrendFollowingSignal(data []Candle, interval string) *Signal {
	if len(data) < 50 {
		return nil
	}

	currentPrice := data[len(data)-1].Close
	currentCandle := data[len(data)-1]
	prevCandle := data[len(data)-2]
	
	// Timeframe-specific parameters - keep strict for all
	minConfirmations := 4
	switch interval {
	case "1m", "3m", "5m":
		minConfirmations = 5 // More confirmations for noisy timeframes
	case "15m", "30m":
		minConfirmations = 4
	case "1h", "2h":
		minConfirmations = 4 // Keep strict
	case "4h", "1d":
		minConfirmations = 4
	}
	
	// Calculate indicators
	ema9 := calculateEMA(data, 9)
	ema20 := calculateEMA(data, 20)
	ema50 := calculateEMA(data, 50)
	atr := calculateATR(data[len(data)-14:], 14)
	rsi := calculateRSI(data, 14)
	
	if atr == 0 {
		return nil
	}

	// Perform advanced analysis for confluence
	analysis := PerformAdvancedAnalysis(data, interval)
	if analysis == nil {
		return nil
	}

	// Calculate trend strength with more nuance
	trendStrength := 0.0
	
	// EMA alignment
	if ema9 > ema20 && ema20 > ema50 {
		trendStrength = 1.0 // Strong uptrend
	} else if ema9 < ema20 && ema20 < ema50 {
		trendStrength = -1.0 // Strong downtrend
	} else if ema9 > ema20 && currentPrice > ema50 {
		trendStrength = 0.7 // Moderate uptrend
	} else if ema9 < ema20 && currentPrice < ema50 {
		trendStrength = -0.7 // Moderate downtrend
	} else if currentPrice > ema20 && currentPrice > ema50 {
		trendStrength = 0.5 // Weak uptrend
	} else if currentPrice < ema20 && currentPrice < ema50 {
		trendStrength = -0.5 // Weak downtrend
	}

	// Volume analysis
	avgVolume := 0.0
	for i := len(data) - 20; i < len(data); i++ {
		avgVolume += data[i].Volume
	}
	avgVolume /= 20
	volumeRatio := currentCandle.Volume / avgVolume

	// Volatility filter - avoid low volatility periods
	recentATR := calculateATR(data[len(data)-5:], 5)
	volatilityOK := recentATR >= atr*0.7

	// Candle analysis
	candleBody := math.Abs(currentCandle.Close - currentCandle.Open)
	candleRange := currentCandle.High - currentCandle.Low
	bodyRatio := 0.0
	if candleRange > 0 {
		bodyRatio = candleBody / candleRange
	}

	var signalType string
	var confidence float64

	// ==================== BULLISH SETUPS ====================
	if trendStrength >= 0.5 && rsi > 40 && rsi < 75 && volatilityOK {
		confirmations := 0
		bonusConfidence := 0.0

		// 1. EMA Pullback Entry
		pullbackToEMA := currentCandle.Low <= ema20*1.005 && currentCandle.Close > ema20
		if pullbackToEMA {
			confirmations++
			bonusConfidence += 5
		}

		// 2. Bullish candle confirmation
		bullishCandle := currentCandle.Close > currentCandle.Open && bodyRatio > 0.5
		if bullishCandle {
			confirmations++
		}

		// 3. Previous candle setup (rejection or small body)
		prevSetup := prevCandle.Close < prevCandle.Open || 
			math.Abs(prevCandle.Close-prevCandle.Open)/(prevCandle.High-prevCandle.Low+0.0001) < 0.3
		if prevSetup {
			confirmations++
		}

		// 4. Order flow confirmation
		if analysis.OrderFlow != nil {
			if analysis.OrderFlow.Delta.Trend == "bullish" {
				confirmations++
				bonusConfidence += 5
			}
			if !analysis.OrderFlow.Delta.Divergence {
				bonusConfidence += 5
			}
			if analysis.OrderFlow.Imbalance > 0.2 {
				confirmations++
			}
		}

		// 5. Volume confirmation
		if volumeRatio > 1.3 {
			confirmations++
			bonusConfidence += 5
		}

		// 6. Pattern confirmation
		for _, p := range analysis.Patterns {
			if p.Type == "bullish" && p.Strength >= 75 {
				confirmations++
				bonusConfidence += p.Strength * 0.1
				break
			}
		}

		// 7. ICT/SMC confluence
		if analysis.ICT != nil {
			if analysis.ICT.PremiumDiscount == "discount" {
				confirmations++
				bonusConfidence += 10
			}
			if analysis.ICT.OTE {
				confirmations++
				bonusConfidence += 8
			}
			// Near bullish order block
			for _, ob := range analysis.EnhancedOBs {
				if ob.Type == "bullish" && !ob.Mitigated {
					if currentPrice >= ob.Low && currentPrice <= ob.High+atr*0.5 {
						confirmations++
						bonusConfidence += 10
						break
					}
				}
			}
		}

		// 8. Structure confirmation
		if analysis.ICT != nil && analysis.ICT.Structure.Trend == "bullish" {
			confirmations++
		}

		// Require minimum confirmations (timeframe-adaptive)
		if confirmations >= minConfirmations {
			signalType = "BUY"
			confidence = 55 + float64(confirmations)*5 + bonusConfidence
		}
	}

	// ==================== BEARISH SETUPS ====================
	if trendStrength <= -0.5 && rsi < 60 && rsi > 25 && volatilityOK && signalType == "" {
		confirmations := 0
		bonusConfidence := 0.0

		// 1. EMA Pullback Entry
		pullbackToEMA := currentCandle.High >= ema20*0.995 && currentCandle.Close < ema20
		if pullbackToEMA {
			confirmations++
			bonusConfidence += 5
		}

		// 2. Bearish candle confirmation
		bearishCandle := currentCandle.Close < currentCandle.Open && bodyRatio > 0.5
		if bearishCandle {
			confirmations++
		}

		// 3. Previous candle setup
		prevSetup := prevCandle.Close > prevCandle.Open || 
			math.Abs(prevCandle.Close-prevCandle.Open)/(prevCandle.High-prevCandle.Low+0.0001) < 0.3
		if prevSetup {
			confirmations++
		}

		// 4. Order flow confirmation
		if analysis.OrderFlow != nil {
			if analysis.OrderFlow.Delta.Trend == "bearish" {
				confirmations++
				bonusConfidence += 5
			}
			if !analysis.OrderFlow.Delta.Divergence {
				bonusConfidence += 5
			}
			if analysis.OrderFlow.Imbalance < -0.2 {
				confirmations++
			}
		}

		// 5. Volume confirmation
		if volumeRatio > 1.3 {
			confirmations++
			bonusConfidence += 5
		}

		// 6. Pattern confirmation
		for _, p := range analysis.Patterns {
			if p.Type == "bearish" && p.Strength >= 75 {
				confirmations++
				bonusConfidence += p.Strength * 0.1
				break
			}
		}

		// 7. ICT/SMC confluence
		if analysis.ICT != nil {
			if analysis.ICT.PremiumDiscount == "premium" {
				confirmations++
				bonusConfidence += 10
			}
			if analysis.ICT.OTE {
				confirmations++
				bonusConfidence += 8
			}
			// Near bearish order block
			for _, ob := range analysis.EnhancedOBs {
				if ob.Type == "bearish" && !ob.Mitigated {
					if currentPrice <= ob.High && currentPrice >= ob.Low-atr*0.5 {
						confirmations++
						bonusConfidence += 10
						break
					}
				}
			}
		}

		// 8. Structure confirmation
		if analysis.ICT != nil && analysis.ICT.Structure.Trend == "bearish" {
			confirmations++
		}

		// Require minimum confirmations (timeframe-adaptive)
		if confirmations >= minConfirmations {
			signalType = "SELL"
			confidence = 55 + float64(confirmations)*5 + bonusConfidence
		}
	}

	if signalType == "" {
		return nil
	}

	// Cap confidence
	confidence = math.Min(confidence, 95)

	// Calculate entry, stop loss, targets
	entry := currentPrice
	var stopLoss, target1, target2, target3 float64

	if signalType == "BUY" {
		// Smart stop loss placement
		swingLow := findRecentSwingLow(data, 10)
		
		// Use the tighter of: swing low, EMA50, or ATR-based
		stopLoss = swingLow - atr*0.2
		if ema50-atr*0.3 > stopLoss && ema50 < entry {
			stopLoss = ema50 - atr*0.3
		}
		
		risk := entry - stopLoss
		if risk <= 0 || risk < atr*0.5 {
			stopLoss = entry - atr*1.2
			risk = atr * 1.2
		}
		if risk > atr*2.5 {
			stopLoss = entry - atr*2.0
			risk = atr * 2.0
		}

		// Dynamic targets based on confidence
		targetMult := 2.0 + (confidence-60)*0.02
		target1 = entry + risk*targetMult
		target2 = entry + risk*targetMult*1.5
		target3 = entry + risk*targetMult*2.0
	} else {
		// Smart stop loss placement
		swingHigh := findRecentSwingHigh(data, 10)
		
		stopLoss = swingHigh + atr*0.2
		if ema50+atr*0.3 < stopLoss && ema50 > entry {
			stopLoss = ema50 + atr*0.3
		}
		
		risk := stopLoss - entry
		if risk <= 0 || risk < atr*0.5 {
			stopLoss = entry + atr*1.2
			risk = atr * 1.2
		}
		if risk > atr*2.5 {
			stopLoss = entry + atr*2.0
			risk = atr * 2.0
		}

		targetMult := 2.0 + (confidence-60)*0.02
		target1 = entry - risk*targetMult
		target2 = entry - risk*targetMult*1.5
		target3 = entry - risk*targetMult*2.0
	}

	// Verify minimum RR
	rr := math.Abs(target1-entry) / math.Abs(entry-stopLoss)
	if rr < 1.8 {
		return nil
	}

	return &Signal{
		Type:     signalType,
		Entry:    entry,
		StopLoss: stopLoss,
		Targets: []Target{
			{Price: target1, RR: rr, Percentage: 100},
			{Price: target2, RR: rr * 1.5, Percentage: 0},
			{Price: target3, RR: rr * 2.0, Percentage: 0},
		},
		Strength:  confidence,
		Timeframe: interval,
	}
}

// generateCombinedStrategy combines pattern + ICT + breakout
func generateCombinedStrategy(data []Candle, interval string) *Signal {
	if len(data) < 30 {
		return nil
	}

	// Perform multi-timeframe analysis
	mta := PerformMultiTimeframeAnalysis(data, interval)

	// If neutral, try to determine direction from recent price action
	if mta.Direction == "neutral" {
		// Check recent trend
		if len(data) >= 10 {
			recentChange := data[len(data)-1].Close - data[len(data)-10].Close
			if recentChange > 0 {
				mta.Direction = "bullish"
			} else if recentChange < 0 {
				mta.Direction = "bearish"
			} else {
				return nil
			}
		} else {
			return nil
		}
	}

	// Get ICT analysis
	ict := PerformICTAnalysis(data)

	// Calculate scores
	patternScore := 0.0
	ictScore := 0.0
	breakoutScore := 0.0

	// Pattern score from MTF analysis
	for _, tf := range []TimeframeData{mta.Higher, mta.Current, mta.Lower} {
		for _, p := range tf.Patterns {
			if p.Type == mta.Direction {
				patternScore += p.Strength * p.Reliability / 1000
			}
		}
	}

	// ICT score
	if ict != nil {
		ictScore = float64(ict.Confluence) * 15

		// Order block near price
		currentPrice := data[len(data)-1].Close
		atr := calculateATR(data[len(data)-14:], 14)

		for _, ob := range ict.OrderBlocks {
			if !ob.Mitigated {
				if mta.Direction == "bullish" && ob.Type == "bullish" {
					if currentPrice >= ob.Low && currentPrice <= ob.High+atr*0.5 {
						ictScore += 20
					}
				} else if mta.Direction == "bearish" && ob.Type == "bearish" {
					if currentPrice <= ob.High && currentPrice >= ob.Low-atr*0.5 {
						ictScore += 20
					}
				}
			}
		}

		// FVG near price
		for _, fvg := range ict.FairValueGaps {
			if !fvg.Filled {
				if mta.Direction == "bullish" && fvg.Type == "bullish" {
					if currentPrice >= fvg.Low && currentPrice <= fvg.High {
						ictScore += 15
					}
				} else if mta.Direction == "bearish" && fvg.Type == "bearish" {
					if currentPrice <= fvg.High && currentPrice >= fvg.Low {
						ictScore += 15
					}
				}
			}
		}
	}

	// Breakout score
	breakoutSignal := checkBreakoutConditions(data, mta.Direction)
	if breakoutSignal {
		breakoutScore = 30
	}

	// Total score
	totalScore := patternScore + ictScore + breakoutScore

	// Lower minimum score threshold
	if totalScore < 30 {
		return nil
	}

	// Generate signal
	currentPrice := data[len(data)-1].Close
	atr := calculateATR(data[len(data)-14:], 14)

	if atr == 0 {
		return nil
	}

	var stopLoss, target1, target2, target3 float64
	signalType := "BUY"

	if mta.Direction == "bullish" {
		stopLoss = findRecentSwingLow(data, 15) - atr*0.3
		riskAmount := currentPrice - stopLoss
		target1 = currentPrice + riskAmount*2.0
		target2 = currentPrice + riskAmount*3.5
		target3 = currentPrice + riskAmount*5.0
	} else {
		signalType = "SELL"
		stopLoss = findRecentSwingHigh(data, 15) + atr*0.3
		riskAmount := stopLoss - currentPrice
		target1 = currentPrice - riskAmount*2.0
		target2 = currentPrice - riskAmount*3.5
		target3 = currentPrice - riskAmount*5.0
	}

	rr := math.Abs(target1-currentPrice) / math.Abs(currentPrice-stopLoss)
	if rr < 1.8 {
		return nil
	}

	confidence := math.Min(totalScore, 95)

	return &Signal{
		Type:     signalType,
		Entry:    currentPrice,
		StopLoss: stopLoss,
		Targets: []Target{
			{Price: target1, RR: rr, Percentage: 100},
			{Price: target2, RR: rr * 1.75, Percentage: 0},
			{Price: target3, RR: rr * 2.5, Percentage: 0},
		},
		Strength:  confidence,
		Timeframe: interval,
	}
}

// checkBreakoutConditions checks for breakout setup
func checkBreakoutConditions(data []Candle, direction string) bool {
	if len(data) < 30 {
		return false
	}

	rangeCandles := data[len(data)-25 : len(data)-5]
	rangeHigh := rangeCandles[0].High
	rangeLow := rangeCandles[0].Low

	for _, c := range rangeCandles {
		if c.High > rangeHigh {
			rangeHigh = c.High
		}
		if c.Low < rangeLow {
			rangeLow = c.Low
		}
	}

	currentPrice := data[len(data)-1].Close
	currentCandle := data[len(data)-1]

	// Volume check
	avgVolume := 0.0
	for _, c := range rangeCandles {
		avgVolume += c.Volume
	}
	avgVolume /= float64(len(rangeCandles))

	if direction == "bullish" {
		return currentPrice > rangeHigh &&
			currentCandle.Close > currentCandle.Open &&
			currentCandle.Volume > avgVolume*1.2
	} else {
		return currentPrice < rangeLow &&
			currentCandle.Close < currentCandle.Open &&
			currentCandle.Volume > avgVolume*1.2
	}
}

// generateBreakoutSignal - Breakout strategy as fallback
func generateBreakoutSignal(data []Candle, interval string) *Signal {
	if len(data) < 50 {
		return nil
	}

	recentCandles := data[len(data)-30:]
	last10 := recentCandles[len(recentCandles)-10:]
	currentCandle := last10[len(last10)-1]
	currentPrice := currentCandle.Close

	// Calculate ATR
	atr := calculateATR(data[len(data)-14:], 14)
	if atr == 0 {
		return nil
	}

	// Find recent high and low (last 20 candles for range)
	rangeCandles := recentCandles[:20]
	rangeHigh := findMax(rangeCandles)
	rangeLow := findMin(rangeCandles)
	rangeSize := rangeHigh - rangeLow

	if rangeSize == 0 {
		return nil
	}

	// Check for breakout
	var signalType string
	var confidence float64

	// Calculate average volume
	avgVolume := 0.0
	for _, c := range rangeCandles {
		avgVolume += c.Volume
	}
	avgVolume /= float64(len(rangeCandles))

	// Check trend alignment (last 10 candles)
	bullishCandles := 0
	bearishCandles := 0
	for _, c := range last10 {
		if c.Close > c.Open {
			bullishCandles++
		} else {
			bearishCandles++
		}
	}

	// BULLISH BREAKOUT - Balanced selectivity
	if currentPrice > rangeHigh && currentCandle.Close > currentCandle.Open {
		// Require good volume (1.3x average)
		if currentCandle.Volume > avgVolume*1.3 {
			// Require decent candle body (50%+)
			bodySize := currentCandle.Close - currentCandle.Open
			candleRange := currentCandle.High - currentCandle.Low
			// Require trend alignment (5+ bullish candles)
			if candleRange > 0 && bodySize/candleRange > 0.5 && bullishCandles >= 5 {
				signalType = "BUY"
				breakoutStrength := (currentPrice - rangeHigh) / atr
				confidence = 65 + math.Min(breakoutStrength*15, 25)
			}
		}
	}

	// BEARISH BREAKOUT - Balanced selectivity
	if currentPrice < rangeLow && currentCandle.Close < currentCandle.Open {
		// Require good volume (1.3x average)
		if currentCandle.Volume > avgVolume*1.3 {
			// Require decent candle body (50%+)
			bodySize := currentCandle.Open - currentCandle.Close
			candleRange := currentCandle.High - currentCandle.Low
			// Require trend alignment (5+ bearish candles)
			if candleRange > 0 && bodySize/candleRange > 0.5 && bearishCandles >= 5 {
				signalType = "SELL"
				breakoutStrength := (rangeLow - currentPrice) / atr
				confidence = 65 + math.Min(breakoutStrength*15, 25)
			}
		}
	}

	if signalType == "" {
		return nil
	}

	confidence = math.Min(confidence, 90)

	// Calculate entry, stop loss, targets
	entry := currentPrice
	var stopLoss, target1, target2, target3 float64

	if signalType == "BUY" {
		stopLoss = rangeHigh - (atr * 0.5)
		target1 = entry + rangeSize*0.5
		target2 = entry + rangeSize
		target3 = entry + rangeSize*1.5
	} else {
		stopLoss = rangeLow + (atr * 0.5)
		target1 = entry - rangeSize*0.5
		target2 = entry - rangeSize
		target3 = entry - rangeSize*1.5
	}

	// Verify minimum RR (higher for better profitability)
	rr := math.Abs(target1-entry) / math.Abs(entry-stopLoss)
	if rr < 2.0 {
		return nil
	}

	return &Signal{
		Type:     signalType,
		Entry:    entry,
		StopLoss: stopLoss,
		Targets: []Target{
			{Price: target1, RR: rr, Percentage: 100},
			{Price: target2, RR: rr * 2, Percentage: 0},
			{Price: target3, RR: rr * 3, Percentage: 0},
		},
		Strength:  confidence,
		Timeframe: interval,
	}
}

// Helper functions
func calculateMA(candles []Candle, period int) float64 {
	if len(candles) < period {
		period = len(candles)
	}
	sum := 0.0
	for i := len(candles) - period; i < len(candles); i++ {
		sum += candles[i].Close
	}
	return sum / float64(period)
}

func findMax(candles []Candle) float64 {
	maxVal := candles[0].High
	for _, c := range candles {
		if c.High > maxVal {
			maxVal = c.High
		}
	}
	return maxVal
}

func findMin(candles []Candle) float64 {
	minVal := candles[0].Low
	for _, c := range candles {
		if c.Low < minVal {
			minVal = c.Low
		}
	}
	return minVal
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
