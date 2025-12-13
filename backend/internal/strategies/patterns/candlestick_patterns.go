package patterns

import (
	"math"
)

// ==================== CANDLESTICK PATTERNS ====================

// CandlestickPattern represents a recognized pattern
type CandlestickPattern struct {
	Name        string
	Type        string  // "bullish", "bearish", "neutral"
	Strength    float64 // 0-100
	Reliability float64 // Historical success rate
	Timeframe   string
	Index       int
}

// MultiTimeframeAnalysis holds analysis across multiple timeframes
type MultiTimeframeAnalysis struct {
	Higher     TimeframeData // 4H or 1D
	Current    TimeframeData // 1H or 15m
	Lower      TimeframeData // 15m or 5m
	Confluence int           // Number of aligned timeframes
	Direction  string        // "bullish", "bearish", "neutral"
	Strength   float64       // Overall strength 0-100
}

// TimeframeData holds analysis for a specific timeframe
type TimeframeData struct {
	Timeframe string
	Trend     string
	Patterns  []CandlestickPattern
	Structure MarketStructure
	ICT       *ICTAnalysis
	Momentum  float64
	Volume    float64
}

// ==================== SINGLE CANDLE PATTERNS ====================

// RecognizeSingleCandlePatterns identifies single candle patterns
func RecognizeSingleCandlePatterns(candles []Candle, timeframe string) []CandlestickPattern {
	var patterns []CandlestickPattern
	if len(candles) < 5 {
		return patterns
	}

	last := candles[len(candles)-1]
	body := math.Abs(last.Close - last.Open)
	upperShadow := last.High - math.Max(last.Open, last.Close)
	lowerShadow := math.Min(last.Open, last.Close) - last.Low
	totalRange := last.High - last.Low

	if totalRange == 0 {
		return patterns
	}

	// Determine trend context
	trendUp := last.Close > calculateSMA(candles, 10)
	trendDown := last.Close < calculateSMA(candles, 10)

	// DOJI - Indecision
	if body/totalRange < 0.1 {
		patterns = append(patterns, CandlestickPattern{
			Name:        "Doji",
			Type:        "neutral",
			Strength:    70,
			Reliability: 65,
			Timeframe:   timeframe,
			Index:       len(candles) - 1,
		})
	}

	// HAMMER - Bullish reversal (in downtrend)
	if lowerShadow > body*2 && upperShadow < body*0.5 && trendDown {
		patterns = append(patterns, CandlestickPattern{
			Name:        "Hammer",
			Type:        "bullish",
			Strength:    80,
			Reliability: 72,
			Timeframe:   timeframe,
			Index:       len(candles) - 1,
		})
	}

	// INVERTED HAMMER - Bullish reversal (in downtrend)
	if upperShadow > body*2 && lowerShadow < body*0.5 && trendDown && last.Close > last.Open {
		patterns = append(patterns, CandlestickPattern{
			Name:        "Inverted Hammer",
			Type:        "bullish",
			Strength:    75,
			Reliability: 68,
			Timeframe:   timeframe,
			Index:       len(candles) - 1,
		})
	}

	// HANGING MAN - Bearish reversal (in uptrend)
	if lowerShadow > body*2 && upperShadow < body*0.5 && trendUp {
		patterns = append(patterns, CandlestickPattern{
			Name:        "Hanging Man",
			Type:        "bearish",
			Strength:    75,
			Reliability: 70,
			Timeframe:   timeframe,
			Index:       len(candles) - 1,
		})
	}

	// SHOOTING STAR - Bearish reversal (in uptrend)
	if upperShadow > body*2 && lowerShadow < body*0.5 && trendUp {
		patterns = append(patterns, CandlestickPattern{
			Name:        "Shooting Star",
			Type:        "bearish",
			Strength:    82,
			Reliability: 76,
			Timeframe:   timeframe,
			Index:       len(candles) - 1,
		})
	}

	// MARUBOZU - Strong continuation
	if body/totalRange > 0.92 {
		patternType := "bullish"
		if last.Close < last.Open {
			patternType = "bearish"
		}
		patterns = append(patterns, CandlestickPattern{
			Name:        "Marubozu",
			Type:        patternType,
			Strength:    88,
			Reliability: 82,
			Timeframe:   timeframe,
			Index:       len(candles) - 1,
		})
	}

	// SPINNING TOP - Indecision
	if body/totalRange < 0.3 && upperShadow > body && lowerShadow > body {
		patterns = append(patterns, CandlestickPattern{
			Name:        "Spinning Top",
			Type:        "neutral",
			Strength:    60,
			Reliability: 55,
			Timeframe:   timeframe,
			Index:       len(candles) - 1,
		})
	}

	// DRAGONFLY DOJI - Bullish reversal
	if body/totalRange < 0.1 && lowerShadow > totalRange*0.7 && upperShadow < totalRange*0.1 {
		patterns = append(patterns, CandlestickPattern{
			Name:        "Dragonfly Doji",
			Type:        "bullish",
			Strength:    78,
			Reliability: 70,
			Timeframe:   timeframe,
			Index:       len(candles) - 1,
		})
	}

	// GRAVESTONE DOJI - Bearish reversal
	if body/totalRange < 0.1 && upperShadow > totalRange*0.7 && lowerShadow < totalRange*0.1 {
		patterns = append(patterns, CandlestickPattern{
			Name:        "Gravestone Doji",
			Type:        "bearish",
			Strength:    78,
			Reliability: 70,
			Timeframe:   timeframe,
			Index:       len(candles) - 1,
		})
	}

	return patterns
}

// ==================== TWO CANDLE PATTERNS ====================

// RecognizeTwoCandlePatterns identifies two candle patterns
func RecognizeTwoCandlePatterns(candles []Candle, timeframe string) []CandlestickPattern {
	var patterns []CandlestickPattern
	if len(candles) < 3 {
		return patterns
	}

	prev := candles[len(candles)-2]
	last := candles[len(candles)-1]

	prevBody := math.Abs(prev.Close - prev.Open)
	lastBody := math.Abs(last.Close - last.Open)
	prevRange := prev.High - prev.Low

	// Trend context
	sma := calculateSMA(candles, 10)
	trendUp := last.Close > sma
	trendDown := last.Close < sma

	// BULLISH ENGULFING
	if prev.Close < prev.Open && last.Close > last.Open &&
		last.Open <= prev.Close && last.Close >= prev.Open &&
		lastBody > prevBody*1.1 {
		strength := 85.0
		if trendDown {
			strength += 5 // Stronger in downtrend
		}
		patterns = append(patterns, CandlestickPattern{
			Name:        "Bullish Engulfing",
			Type:        "bullish",
			Strength:    strength,
			Reliability: 78,
			Timeframe:   timeframe,
			Index:       len(candles) - 1,
		})
	}

	// BEARISH ENGULFING
	if prev.Close > prev.Open && last.Close < last.Open &&
		last.Open >= prev.Close && last.Close <= prev.Open &&
		lastBody > prevBody*1.1 {
		strength := 85.0
		if trendUp {
			strength += 5 // Stronger in uptrend
		}
		patterns = append(patterns, CandlestickPattern{
			Name:        "Bearish Engulfing",
			Type:        "bearish",
			Strength:    strength,
			Reliability: 78,
			Timeframe:   timeframe,
			Index:       len(candles) - 1,
		})
	}

	// PIERCING PATTERN - Bullish reversal
	if prev.Close < prev.Open && last.Close > last.Open &&
		last.Open < prev.Low && last.Close > (prev.Open+prev.Close)/2 &&
		last.Close < prev.Open {
		patterns = append(patterns, CandlestickPattern{
			Name:        "Piercing Pattern",
			Type:        "bullish",
			Strength:    80,
			Reliability: 72,
			Timeframe:   timeframe,
			Index:       len(candles) - 1,
		})
	}

	// DARK CLOUD COVER - Bearish reversal
	if prev.Close > prev.Open && last.Close < last.Open &&
		last.Open > prev.High && last.Close < (prev.Open+prev.Close)/2 &&
		last.Close > prev.Open {
		patterns = append(patterns, CandlestickPattern{
			Name:        "Dark Cloud Cover",
			Type:        "bearish",
			Strength:    80,
			Reliability: 72,
			Timeframe:   timeframe,
			Index:       len(candles) - 1,
		})
	}

	// TWEEZER BOTTOM - Bullish reversal
	tolerance := prevRange * 0.05
	if math.Abs(prev.Low-last.Low) < tolerance && prev.Close < prev.Open && last.Close > last.Open {
		patterns = append(patterns, CandlestickPattern{
			Name:        "Tweezer Bottom",
			Type:        "bullish",
			Strength:    75,
			Reliability: 68,
			Timeframe:   timeframe,
			Index:       len(candles) - 1,
		})
	}

	// TWEEZER TOP - Bearish reversal
	if math.Abs(prev.High-last.High) < tolerance && prev.Close > prev.Open && last.Close < last.Open {
		patterns = append(patterns, CandlestickPattern{
			Name:        "Tweezer Top",
			Type:        "bearish",
			Strength:    75,
			Reliability: 68,
			Timeframe:   timeframe,
			Index:       len(candles) - 1,
		})
	}

	// HARAMI BULLISH
	if prev.Close < prev.Open && last.Close > last.Open &&
		last.Open > prev.Close && last.Close < prev.Open &&
		lastBody < prevBody*0.5 {
		patterns = append(patterns, CandlestickPattern{
			Name:        "Bullish Harami",
			Type:        "bullish",
			Strength:    72,
			Reliability: 65,
			Timeframe:   timeframe,
			Index:       len(candles) - 1,
		})
	}

	// HARAMI BEARISH
	if prev.Close > prev.Open && last.Close < last.Open &&
		last.Open < prev.Close && last.Close > prev.Open &&
		lastBody < prevBody*0.5 {
		patterns = append(patterns, CandlestickPattern{
			Name:        "Bearish Harami",
			Type:        "bearish",
			Strength:    72,
			Reliability: 65,
			Timeframe:   timeframe,
			Index:       len(candles) - 1,
		})
	}

	return patterns
}


// ==================== THREE CANDLE PATTERNS ====================

// RecognizeThreeCandlePatterns identifies three candle patterns
func RecognizeThreeCandlePatterns(candles []Candle, timeframe string) []CandlestickPattern {
	var patterns []CandlestickPattern
	if len(candles) < 4 {
		return patterns
	}

	c1 := candles[len(candles)-3]
	c2 := candles[len(candles)-2]
	c3 := candles[len(candles)-1]

	c2Body := math.Abs(c2.Close - c2.Open)
	c1Range := c1.High - c1.Low

	// MORNING STAR - Strong bullish reversal
	if c1.Close < c1.Open && // First bearish
		c2Body < c1Range*0.3 && // Second small body (star)
		c3.Close > c3.Open && // Third bullish
		c3.Close > (c1.Open+c1.Close)/2 { // Third closes above midpoint
		patterns = append(patterns, CandlestickPattern{
			Name:        "Morning Star",
			Type:        "bullish",
			Strength:    92,
			Reliability: 85,
			Timeframe:   timeframe,
			Index:       len(candles) - 1,
		})
	}

	// EVENING STAR - Strong bearish reversal
	if c1.Close > c1.Open && // First bullish
		c2Body < c1Range*0.3 && // Second small body (star)
		c3.Close < c3.Open && // Third bearish
		c3.Close < (c1.Open+c1.Close)/2 { // Third closes below midpoint
		patterns = append(patterns, CandlestickPattern{
			Name:        "Evening Star",
			Type:        "bearish",
			Strength:    92,
			Reliability: 85,
			Timeframe:   timeframe,
			Index:       len(candles) - 1,
		})
	}

	// THREE WHITE SOLDIERS - Strong bullish continuation
	if c1.Close > c1.Open && c2.Close > c2.Open && c3.Close > c3.Open &&
		c2.Close > c1.Close && c3.Close > c2.Close &&
		c2.Open > c1.Open && c2.Open < c1.Close &&
		c3.Open > c2.Open && c3.Open < c2.Close {
		patterns = append(patterns, CandlestickPattern{
			Name:        "Three White Soldiers",
			Type:        "bullish",
			Strength:    90,
			Reliability: 82,
			Timeframe:   timeframe,
			Index:       len(candles) - 1,
		})
	}

	// THREE BLACK CROWS - Strong bearish continuation
	if c1.Close < c1.Open && c2.Close < c2.Open && c3.Close < c3.Open &&
		c2.Close < c1.Close && c3.Close < c2.Close &&
		c2.Open < c1.Open && c2.Open > c1.Close &&
		c3.Open < c2.Open && c3.Open > c2.Close {
		patterns = append(patterns, CandlestickPattern{
			Name:        "Three Black Crows",
			Type:        "bearish",
			Strength:    90,
			Reliability: 82,
			Timeframe:   timeframe,
			Index:       len(candles) - 1,
		})
	}

	// THREE INSIDE UP - Bullish reversal
	if c1.Close < c1.Open && // First bearish
		c2.Close > c2.Open && // Second bullish
		c2.Open > c1.Close && c2.Close < c1.Open && // Second inside first
		c3.Close > c3.Open && c3.Close > c1.Open { // Third bullish, closes above first
		patterns = append(patterns, CandlestickPattern{
			Name:        "Three Inside Up",
			Type:        "bullish",
			Strength:    85,
			Reliability: 78,
			Timeframe:   timeframe,
			Index:       len(candles) - 1,
		})
	}

	// THREE INSIDE DOWN - Bearish reversal
	if c1.Close > c1.Open && // First bullish
		c2.Close < c2.Open && // Second bearish
		c2.Open < c1.Close && c2.Close > c1.Open && // Second inside first
		c3.Close < c3.Open && c3.Close < c1.Open { // Third bearish, closes below first
		patterns = append(patterns, CandlestickPattern{
			Name:        "Three Inside Down",
			Type:        "bearish",
			Strength:    85,
			Reliability: 78,
			Timeframe:   timeframe,
			Index:       len(candles) - 1,
		})
	}

	// ABANDONED BABY BULLISH
	if c1.Close < c1.Open && // First bearish
		c2Body < c1Range*0.1 && // Second doji
		c2.High < c1.Low && // Gap down
		c3.Close > c3.Open && // Third bullish
		c3.Low > c2.High { // Gap up
		patterns = append(patterns, CandlestickPattern{
			Name:        "Abandoned Baby Bullish",
			Type:        "bullish",
			Strength:    95,
			Reliability: 88,
			Timeframe:   timeframe,
			Index:       len(candles) - 1,
		})
	}

	// ABANDONED BABY BEARISH
	if c1.Close > c1.Open && // First bullish
		c2Body < c1Range*0.1 && // Second doji
		c2.Low > c1.High && // Gap up
		c3.Close < c3.Open && // Third bearish
		c3.High < c2.Low { // Gap down
		patterns = append(patterns, CandlestickPattern{
			Name:        "Abandoned Baby Bearish",
			Type:        "bearish",
			Strength:    95,
			Reliability: 88,
			Timeframe:   timeframe,
			Index:       len(candles) - 1,
		})
	}

	return patterns
}

// ==================== MULTI-CANDLE PATTERNS ====================

// RecognizeMultiCandlePatterns identifies patterns with 4+ candles
func RecognizeMultiCandlePatterns(candles []Candle, timeframe string) []CandlestickPattern {
	var patterns []CandlestickPattern
	if len(candles) < 6 {
		return patterns
	}

	// RISING THREE METHODS - Bullish continuation
	if len(candles) >= 5 {
		c1 := candles[len(candles)-5]
		c2 := candles[len(candles)-4]
		c3 := candles[len(candles)-3]
		c4 := candles[len(candles)-2]
		c5 := candles[len(candles)-1]

		if c1.Close > c1.Open && c5.Close > c5.Open && // First and last bullish
			c2.Close < c2.Open && c3.Close < c3.Open && c4.Close < c4.Open && // Middle bearish
			c2.Low > c1.Low && c3.Low > c1.Low && c4.Low > c1.Low && // Stay above first low
			c5.Close > c1.High { // Last closes above first high
			patterns = append(patterns, CandlestickPattern{
				Name:        "Rising Three Methods",
				Type:        "bullish",
				Strength:    88,
				Reliability: 78,
				Timeframe:   timeframe,
				Index:       len(candles) - 1,
			})
		}

		// FALLING THREE METHODS - Bearish continuation
		if c1.Close < c1.Open && c5.Close < c5.Open && // First and last bearish
			c2.Close > c2.Open && c3.Close > c3.Open && c4.Close > c4.Open && // Middle bullish
			c2.High < c1.High && c3.High < c1.High && c4.High < c1.High && // Stay below first high
			c5.Close < c1.Low { // Last closes below first low
			patterns = append(patterns, CandlestickPattern{
				Name:        "Falling Three Methods",
				Type:        "bearish",
				Strength:    88,
				Reliability: 78,
				Timeframe:   timeframe,
				Index:       len(candles) - 1,
			})
		}
	}

	// BULLISH BREAKAWAY
	if len(candles) >= 5 {
		c1 := candles[len(candles)-5]
		c5 := candles[len(candles)-1]

		// Long bearish, then small bodies, then bullish breakout
		if c1.Close < c1.Open && c5.Close > c5.Open && c5.Close > c1.Open {
			smallBodies := true
			for i := len(candles) - 4; i < len(candles)-1; i++ {
				body := math.Abs(candles[i].Close - candles[i].Open)
				range_ := candles[i].High - candles[i].Low
				if range_ > 0 && body/range_ > 0.5 {
					smallBodies = false
					break
				}
			}
			if smallBodies {
				patterns = append(patterns, CandlestickPattern{
					Name:        "Bullish Breakaway",
					Type:        "bullish",
					Strength:    82,
					Reliability: 72,
					Timeframe:   timeframe,
					Index:       len(candles) - 1,
				})
			}
		}
	}

	return patterns
}

// ==================== HELPER FUNCTIONS ====================

// calculateSMA calculates Simple Moving Average
func calculateSMA(candles []Candle, period int) float64 {
	if len(candles) < period {
		period = len(candles)
	}
	sum := 0.0
	for i := len(candles) - period; i < len(candles); i++ {
		sum += candles[i].Close
	}
	return sum / float64(period)
}

// calculateEMA calculates Exponential Moving Average
func calculateEMA(candles []Candle, period int) float64 {
	if len(candles) < period {
		return calculateSMA(candles, period)
	}

	multiplier := 2.0 / float64(period+1)
	ema := calculateSMA(candles[:period], period)

	for i := period; i < len(candles); i++ {
		ema = (candles[i].Close-ema)*multiplier + ema
	}
	return ema
}

// calculateRSI calculates Relative Strength Index
func calculateRSI(candles []Candle, period int) float64 {
	if len(candles) < period+1 {
		return 50
	}

	gains := 0.0
	losses := 0.0

	for i := len(candles) - period; i < len(candles); i++ {
		change := candles[i].Close - candles[i-1].Close
		if change > 0 {
			gains += change
		} else {
			losses -= change
		}
	}

	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	if avgLoss == 0 {
		return 100
	}

	rs := avgGain / avgLoss
	return 100 - (100 / (1 + rs))
}

// calculateMomentumScore calculates momentum strength 0-100
func calculateMomentumScore(candles []Candle) float64 {
	if len(candles) < 14 {
		return 50
	}

	// RSI component
	rsi := calculateRSI(candles, 14)

	// Price change component
	priceChange := (candles[len(candles)-1].Close - candles[len(candles)-10].Close) / candles[len(candles)-10].Close * 100

	// EMA alignment
	ema9 := calculateEMA(candles, 9)
	ema21 := calculateEMA(candles, 21)
	emaScore := 50.0
	if ema9 > ema21 {
		emaScore = 70
	} else if ema9 < ema21 {
		emaScore = 30
	}

	// Combine scores
	momentum := (rsi*0.4 + (50+priceChange*5)*0.3 + emaScore*0.3)
	return math.Max(0, math.Min(100, momentum))
}

// calculateVolumeScore calculates volume strength 0-100
func calculateVolumeScore(candles []Candle) float64 {
	if len(candles) < 20 {
		return 50
	}

	recent := candles[len(candles)-5:]
	older := candles[len(candles)-20 : len(candles)-5]

	recentAvg := 0.0
	olderAvg := 0.0

	for _, c := range recent {
		recentAvg += c.Volume
	}
	for _, c := range older {
		olderAvg += c.Volume
	}

	recentAvg /= float64(len(recent))
	olderAvg /= float64(len(older))

	if olderAvg == 0 {
		return 50
	}

	ratio := recentAvg / olderAvg
	strength := 50 + ((ratio - 1) * 30)
	return math.Max(0, math.Min(100, strength))
}


// ==================== MULTI-TIMEFRAME ANALYSIS ====================

// RecognizeAllPatterns identifies all candlestick patterns
func RecognizeAllPatterns(candles []Candle, timeframe string) []CandlestickPattern {
	var patterns []CandlestickPattern

	patterns = append(patterns, RecognizeSingleCandlePatterns(candles, timeframe)...)
	patterns = append(patterns, RecognizeTwoCandlePatterns(candles, timeframe)...)
	patterns = append(patterns, RecognizeThreeCandlePatterns(candles, timeframe)...)
	patterns = append(patterns, RecognizeMultiCandlePatterns(candles, timeframe)...)

	return patterns
}

// AnalyzeTimeframe performs complete analysis for a single timeframe
func AnalyzeTimeframe(candles []Candle, timeframe string) TimeframeData {
	td := TimeframeData{
		Timeframe: timeframe,
		Trend:     "neutral",
	}

	if len(candles) < 30 {
		return td
	}

	// Pattern recognition
	td.Patterns = RecognizeAllPatterns(candles, timeframe)

	// Market structure
	td.Structure = AnalyzeMarketStructure(candles)
	td.Trend = td.Structure.Trend

	// ICT analysis
	td.ICT = PerformICTAnalysis(candles)

	// Momentum calculation
	td.Momentum = calculateMomentumScore(candles)

	// Volume analysis
	td.Volume = calculateVolumeScore(candles)

	return td
}

// ConvertToHigherTimeframe converts candles to higher timeframe
func ConvertToHigherTimeframe(candles []Candle, multiplier int) []Candle {
	if len(candles) < multiplier {
		return candles
	}

	var htfCandles []Candle
	for i := 0; i <= len(candles)-multiplier; i += multiplier {
		chunk := candles[i : i+multiplier]

		htfCandle := Candle{
			Open:      chunk[0].Open,
			High:      chunk[0].High,
			Low:       chunk[0].Low,
			Close:     chunk[len(chunk)-1].Close,
			Volume:    0,
			Timestamp: chunk[0].Timestamp,
		}

		for _, c := range chunk {
			if c.High > htfCandle.High {
				htfCandle.High = c.High
			}
			if c.Low < htfCandle.Low {
				htfCandle.Low = c.Low
			}
			htfCandle.Volume += c.Volume
		}

		htfCandles = append(htfCandles, htfCandle)
	}

	return htfCandles
}

// PerformMultiTimeframeAnalysis analyzes multiple timeframes
func PerformMultiTimeframeAnalysis(candles []Candle, baseInterval string) *MultiTimeframeAnalysis {
	mta := &MultiTimeframeAnalysis{}

	// Determine multipliers based on base interval
	var higherMult, currentMult int

	switch baseInterval {
	case "1m":
		higherMult = 60  // 1H
		currentMult = 15 // 15m
	case "5m":
		higherMult = 48 // 4H
		currentMult = 12 // 1H
	case "15m":
		higherMult = 16 // 4H
		currentMult = 4  // 1H
	case "1h":
		higherMult = 24 // 1D
		currentMult = 4  // 4H
	case "4h":
		higherMult = 6  // 1D
		currentMult = 1 // 4H (same)
	default:
		higherMult = 4
		currentMult = 1
	}

	// Convert to higher timeframes
	higherCandles := ConvertToHigherTimeframe(candles, higherMult)
	currentCandles := ConvertToHigherTimeframe(candles, currentMult)

	// Analyze each timeframe
	mta.Higher = AnalyzeTimeframe(higherCandles, getHigherTimeframeName(baseInterval))
	mta.Current = AnalyzeTimeframe(currentCandles, getCurrentTimeframeName(baseInterval))
	mta.Lower = AnalyzeTimeframe(candles, baseInterval)

	// Calculate confluence
	directions := []string{mta.Higher.Trend, mta.Current.Trend, mta.Lower.Trend}
	bullishCount := 0
	bearishCount := 0

	for _, dir := range directions {
		if dir == "bullish" {
			bullishCount++
		} else if dir == "bearish" {
			bearishCount++
		}
	}

	// Determine overall direction and confluence
	if bullishCount >= 2 {
		mta.Direction = "bullish"
		mta.Confluence = bullishCount
	} else if bearishCount >= 2 {
		mta.Direction = "bearish"
		mta.Confluence = bearishCount
	} else {
		mta.Direction = "neutral"
		mta.Confluence = 0
	}

	// Calculate overall strength (weighted by timeframe importance)
	totalStrength := mta.Higher.Momentum*0.45 + mta.Current.Momentum*0.35 + mta.Lower.Momentum*0.20
	mta.Strength = math.Min(totalStrength, 100)

	return mta
}

// Helper functions for timeframe names
func getHigherTimeframeName(base string) string {
	switch base {
	case "1m":
		return "1H"
	case "5m":
		return "4H"
	case "15m":
		return "4H"
	case "1h":
		return "1D"
	case "4h":
		return "1D"
	default:
		return "4H"
	}
}

func getCurrentTimeframeName(base string) string {
	switch base {
	case "1m":
		return "15m"
	case "5m":
		return "1H"
	case "15m":
		return "1H"
	case "1h":
		return "4H"
	case "4h":
		return "4H"
	default:
		return "1H"
	}
}

// ==================== PATTERN PREDICTION SIGNAL GENERATOR ====================

// GeneratePatternPredictionSignal creates signals based on multi-timeframe pattern analysis
func GeneratePatternPredictionSignal(candles []Candle, interval string) *Signal {
	if len(candles) < 50 {
		return nil
	}

	// Perform multi-timeframe analysis
	mta := PerformMultiTimeframeAnalysis(candles, interval)

	// Allow signals even with 1 timeframe alignment if patterns are strong
	if mta.Direction == "neutral" {
		return nil
	}

	// Get current price and ATR
	currentPrice := candles[len(candles)-1].Close
	atr := calculateATR(candles[len(candles)-14:], 14)
	if atr == 0 {
		return nil
	}

	// Find best pattern across timeframes
	var bestPattern *CandlestickPattern
	maxScore := 0.0

	// Weight patterns by timeframe (higher = more important)
	timeframes := []TimeframeData{mta.Higher, mta.Current, mta.Lower}
	weights := []float64{0.40, 0.35, 0.25}

	for i, tf := range timeframes {
		for _, pattern := range tf.Patterns {
			// Consider patterns aligned with direction OR strong reversal patterns
			if pattern.Type == mta.Direction || pattern.Strength >= 85 {
				score := (pattern.Strength * pattern.Reliability * weights[i]) / 100
				if score > maxScore {
					maxScore = score
					p := pattern
					bestPattern = &p
				}
			}
		}
	}

	// Lower threshold for pattern score
	if bestPattern == nil || maxScore < 25 {
		return nil
	}

	// Calculate confidence based on multiple factors
	confidence := maxScore

	// Confluence bonus
	confidence += float64(mta.Confluence) * 8

	// Momentum alignment bonus
	if mta.Direction == "bullish" && mta.Strength > 55 {
		confidence += (mta.Strength - 55) * 0.4
	} else if mta.Direction == "bearish" && mta.Strength < 45 {
		confidence += (45 - mta.Strength) * 0.4
	}

	// ICT confluence bonus
	if mta.Current.ICT != nil {
		confidence += float64(mta.Current.ICT.Confluence) * 4

		// Premium/Discount alignment
		if mta.Direction == "bullish" && mta.Current.ICT.PremiumDiscount == "discount" {
			confidence += 8
		} else if mta.Direction == "bearish" && mta.Current.ICT.PremiumDiscount == "premium" {
			confidence += 8
		}

		// OTE bonus
		if mta.Current.ICT.OTE {
			confidence += 6
		}
	}

	// Volume confirmation
	if mta.Lower.Volume > 60 {
		confidence += 5
	}

	// Pattern reliability bonus
	confidence += (bestPattern.Reliability - 70) * 0.3

	confidence = math.Min(confidence, 95)

	// Lower minimum confidence threshold
	if confidence < 50 {
		return nil
	}

	// Calculate entry, stop loss, and targets
	entry := currentPrice
	var stopLoss, target1, target2, target3 float64

	// Dynamic ATR multiplier based on pattern strength
	atrMult := 1.2 + (bestPattern.Strength-70)*0.02

	if mta.Direction == "bullish" {
		// Find recent swing low for stop
		swingLow := findRecentSwingLow(candles, 10)
		stopLoss = math.Min(entry-(atr*atrMult), swingLow-atr*0.2)

		// Targets based on pattern projection
		riskAmount := entry - stopLoss
		target1 = entry + riskAmount*2.0  // 2R
		target2 = entry + riskAmount*3.5  // 3.5R
		target3 = entry + riskAmount*5.0  // 5R

	} else { // bearish
		// Find recent swing high for stop
		swingHigh := findRecentSwingHigh(candles, 10)
		stopLoss = math.Max(entry+(atr*atrMult), swingHigh+atr*0.2)

		// Targets based on pattern projection
		riskAmount := stopLoss - entry
		target1 = entry - riskAmount*2.0
		target2 = entry - riskAmount*3.5
		target3 = entry - riskAmount*5.0
	}

	// Verify minimum RR
	rr := math.Abs(target1-entry) / math.Abs(entry-stopLoss)
	if rr < 1.8 {
		return nil
	}

	signalType := "BUY"
	if mta.Direction == "bearish" {
		signalType = "SELL"
	}

	return &Signal{
		Type:     signalType,
		Entry:    entry,
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

// findRecentSwingLow finds the lowest low in recent candles
func findRecentSwingLow(candles []Candle, lookback int) float64 {
	if len(candles) < lookback {
		lookback = len(candles)
	}

	swingLow := candles[len(candles)-1].Low
	for i := len(candles) - lookback; i < len(candles); i++ {
		if candles[i].Low < swingLow {
			swingLow = candles[i].Low
		}
	}
	return swingLow
}

// findRecentSwingHigh finds the highest high in recent candles
func findRecentSwingHigh(candles []Candle, lookback int) float64 {
	if len(candles) < lookback {
		lookback = len(candles)
	}

	swingHigh := candles[len(candles)-1].High
	for i := len(candles) - lookback; i < len(candles); i++ {
		if candles[i].High > swingHigh {
			swingHigh = candles[i].High
		}
	}
	return swingHigh
}
