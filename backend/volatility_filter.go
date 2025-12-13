package main

import (
	"math"
	"sort"
)

// ==================== VOLATILITY FILTERING ====================
// Trade only in optimal volatility conditions

// VolatilityAnalysis holds volatility metrics
type VolatilityAnalysis struct {
	ATR            float64
	ATRPercent     float64
	ATRPercentile  float64
	BBWidth        float64
	BBPercentile   float64
	Regime         string  // "low", "normal", "high", "extreme"
	Optimal        bool
	TradingAdvice  string
	SizeMultiplier float64
}

// ==================== ATR ANALYSIS ====================

// CalculateATRPercentile calculates ATR percentile (0-100)
func CalculateATRPercentile(candles []Candle, period int) float64 {
	if len(candles) < period*3 {
		return 50
	}
	
	// Calculate ATR for multiple periods
	atrs := []float64{}
	
	for i := period; i < len(candles); i++ {
		subset := candles[i-period : i]
		atr := calculateATR(subset, period)
		if atr > 0 {
			atrs = append(atrs, atr)
		}
	}
	
	if len(atrs) == 0 {
		return 50
	}
	
	// Get current ATR
	currentATR := calculateATR(candles[len(candles)-period:], period)
	
	// Calculate percentile
	sort.Float64s(atrs)
	
	rank := 0
	for _, atr := range atrs {
		if currentATR > atr {
			rank++
		}
	}
	
	percentile := float64(rank) / float64(len(atrs)) * 100
	return percentile
}

// ==================== BOLLINGER BAND WIDTH ====================

// CalculateBollingerBandWidth calculates BB width percentile
func CalculateBollingerBandWidth(candles []Candle, period int, stdDev float64) (float64, float64) {
	if len(candles) < period*2 {
		return 0, 50
	}
	
	// Calculate BB width for multiple periods
	widths := []float64{}
	
	for i := period; i < len(candles); i++ {
		subset := candles[i-period : i]
		
		// Calculate SMA
		sum := 0.0
		for _, c := range subset {
			sum += c.Close
		}
		sma := sum / float64(period)
		
		// Calculate standard deviation
		variance := 0.0
		for _, c := range subset {
			variance += math.Pow(c.Close-sma, 2)
		}
		std := math.Sqrt(variance / float64(period))
		
		// BB width = (upper - lower) / middle
		upper := sma + (stdDev * std)
		lower := sma - (stdDev * std)
		width := (upper - lower) / sma * 100
		
		widths = append(widths, width)
	}
	
	if len(widths) == 0 {
		return 0, 50
	}
	
	// Get current width
	currentWidth := widths[len(widths)-1]
	
	// Calculate percentile
	sort.Float64s(widths)
	
	rank := 0
	for _, w := range widths {
		if currentWidth > w {
			rank++
		}
	}
	
	percentile := float64(rank) / float64(len(widths)) * 100
	return currentWidth, percentile
}

// ==================== VOLATILITY REGIME ====================

// AnalyzeVolatility performs complete volatility analysis
func AnalyzeVolatility(candles []Candle) *VolatilityAnalysis {
	va := &VolatilityAnalysis{
		Optimal:        false,
		SizeMultiplier: 1.0,
	}
	
	if len(candles) < 50 {
		va.Regime = "unknown"
		va.TradingAdvice = "Not enough data"
		return va
	}
	
	// Calculate ATR
	va.ATR = calculateATR(candles[len(candles)-14:], 14)
	currentPrice := candles[len(candles)-1].Close
	va.ATRPercent = (va.ATR / currentPrice) * 100
	
	// Calculate ATR percentile
	va.ATRPercentile = CalculateATRPercentile(candles, 14)
	
	// Calculate Bollinger Band width
	va.BBWidth, va.BBPercentile = CalculateBollingerBandWidth(candles, 20, 2.0)
	
	// Determine volatility regime
	if va.ATRPercentile < 20 {
		va.Regime = "low"
		va.TradingAdvice = "Low volatility - Expect breakout soon"
		va.SizeMultiplier = 0.8 // Reduce size (choppy)
		
	} else if va.ATRPercentile >= 20 && va.ATRPercentile < 40 {
		va.Regime = "normal-low"
		va.TradingAdvice = "Below average volatility - Good for scalping"
		va.Optimal = true
		va.SizeMultiplier = 1.0
		
	} else if va.ATRPercentile >= 40 && va.ATRPercentile < 70 {
		va.Regime = "normal"
		va.TradingAdvice = "Optimal volatility - Best trading conditions"
		va.Optimal = true
		va.SizeMultiplier = 1.2 // Increase size (optimal)
		
	} else if va.ATRPercentile >= 70 && va.ATRPercentile < 90 {
		va.Regime = "high"
		va.TradingAdvice = "High volatility - Reduce position size"
		va.Optimal = true
		va.SizeMultiplier = 0.9
		
	} else {
		va.Regime = "extreme"
		va.TradingAdvice = "Extreme volatility - Avoid trading or hedge"
		va.Optimal = false
		va.SizeMultiplier = 0.5 // Significantly reduce size
	}
	
	// BB squeeze detection (potential breakout)
	if va.BBPercentile < 20 {
		va.TradingAdvice += " | BB Squeeze - Breakout imminent"
		va.SizeMultiplier *= 1.2 // Increase for breakout
	}
	
	// BB expansion (trend in progress)
	if va.BBPercentile > 80 {
		va.TradingAdvice += " | BB Expansion - Strong trend"
	}
	
	return va
}

// ShouldTradeVolatility determines if volatility is suitable for trading
func ShouldTradeVolatility(va *VolatilityAnalysis) bool {
	if va == nil {
		return false
	}
	
	// Don't trade in extreme volatility
	if va.Regime == "extreme" {
		return false
	}
	
	// Don't trade in very low volatility (choppy)
	if va.ATRPercentile < 15 {
		return false
	}
	
	// Optimal trading conditions
	if va.Optimal {
		return true
	}
	
	return false
}

// GetVolatilityAdjustedStopLoss calculates stop loss based on volatility
func GetVolatilityAdjustedStopLoss(entry float64, direction string, va *VolatilityAnalysis) float64 {
	if va == nil {
		// Default 2% stop
		if direction == "bullish" {
			return entry * 0.98
		}
		return entry * 1.02
	}
	
	// Base stop loss multiplier
	multiplier := 1.5
	
	// Adjust based on volatility regime
	switch va.Regime {
	case "low":
		multiplier = 1.0 // Tight stop in low volatility
	case "normal-low":
		multiplier = 1.3
	case "normal":
		multiplier = 1.5
	case "high":
		multiplier = 2.0 // Wider stop in high volatility
	case "extreme":
		multiplier = 2.5 // Very wide stop
	}
	
	// Calculate stop loss
	if direction == "bullish" {
		return entry - (va.ATR * multiplier)
	}
	return entry + (va.ATR * multiplier)
}

// GetVolatilityAdjustedTargets calculates targets based on volatility
func GetVolatilityAdjustedTargets(entry float64, stopLoss float64, direction string, va *VolatilityAnalysis) (float64, float64, float64) {
	risk := math.Abs(entry - stopLoss)
	
	// Base risk:reward ratios
	rr1 := 2.0
	rr2 := 3.5
	rr3 := 5.0
	
	// Adjust based on volatility
	if va != nil {
		switch va.Regime {
		case "low":
			// Smaller targets in low volatility
			rr1 = 1.5
			rr2 = 2.5
			rr3 = 3.5
			
		case "normal-low":
			rr1 = 1.8
			rr2 = 3.0
			rr3 = 4.5
			
		case "normal":
			// Standard targets
			rr1 = 2.0
			rr2 = 3.5
			rr3 = 5.0
			
		case "high":
			// Larger targets in high volatility
			rr1 = 2.5
			rr2 = 4.0
			rr3 = 6.0
			
		case "extreme":
			// Very large targets (or avoid trading)
			rr1 = 3.0
			rr2 = 5.0
			rr3 = 7.5
		}
	}
	
	// Calculate targets
	var tp1, tp2, tp3 float64
	
	if direction == "bullish" {
		tp1 = entry + (risk * rr1)
		tp2 = entry + (risk * rr2)
		tp3 = entry + (risk * rr3)
	} else {
		tp1 = entry - (risk * rr1)
		tp2 = entry - (risk * rr2)
		tp3 = entry - (risk * rr3)
	}
	
	return tp1, tp2, tp3
}

// GetVolatilityScore returns a score (0-100) for signal strength
func GetVolatilityScore(va *VolatilityAnalysis) float64 {
	if va == nil {
		return 50
	}
	
	score := 50.0
	
	// Optimal regime bonus
	if va.Optimal {
		score += 20
	}
	
	// Regime-specific scoring
	switch va.Regime {
	case "normal":
		score += 15 // Best regime
	case "normal-low", "high":
		score += 10 // Good regimes
	case "low":
		score -= 10 // Choppy
	case "extreme":
		score -= 30 // Dangerous
	}
	
	// BB squeeze bonus (breakout opportunity)
	if va.BBPercentile < 20 {
		score += 15
	}
	
	// BB expansion bonus (trend)
	if va.BBPercentile > 70 && va.BBPercentile < 90 {
		score += 10
	}
	
	// Extreme expansion penalty
	if va.BBPercentile > 90 {
		score -= 10
	}
	
	return math.Max(0, math.Min(score, 100))
}

// ==================== VOLATILITY-BASED FILTERS ====================

// FilterByVolatility filters signals based on volatility
func FilterByVolatility(signal map[string]interface{}, va *VolatilityAnalysis) bool {
	if va == nil {
		return true // No filter if no analysis
	}
	
	// Reject extreme volatility
	if va.Regime == "extreme" {
		return false
	}
	
	// Reject very low volatility
	if va.ATRPercentile < 15 {
		return false
	}
	
	// Accept optimal volatility
	if va.Optimal {
		return true
	}
	
	// Accept if volatility score is good
	score := GetVolatilityScore(va)
	return score >= 55
}
