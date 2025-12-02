package main

import (
	"math"
)

// ==================== ADVANCED CANDLESTICK PATTERNS ====================
// Professional-grade pattern recognition

// AdvancedPattern represents an advanced pattern
type AdvancedPattern struct {
	Name        string
	Type        string  // "bullish", "bearish", "neutral"
	Strength    float64
	Reliability float64
	Description string
	Entry       float64
	StopLoss    float64
	Target      float64
	CandleIdx   int
}

// RecognizeAdvancedPatterns identifies advanced patterns
func RecognizeAdvancedPatterns(candles []Candle) []AdvancedPattern {
	var patterns []AdvancedPattern
	
	if len(candles) < 10 {
		return patterns
	}
	
	// Three White Soldiers
	if p := detectThreeWhiteSoldiers(candles); p != nil {
		patterns = append(patterns, *p)
	}
	
	// Three Black Crows
	if p := detectThreeBlackCrows(candles); p != nil {
		patterns = append(patterns, *p)
	}
	
	// Rising Three Methods
	if p := detectRisingThreeMethods(candles); p != nil {
		patterns = append(patterns, *p)
	}
	
	// Falling Three Methods
	if p := detectFallingThreeMethods(candles); p != nil {
		patterns = append(patterns, *p)
	}
	
	// Abandoned Baby
	if p := detectAbandonedBaby(candles); p != nil {
		patterns = append(patterns, *p)
	}
	
	// Three Inside Up/Down
	if p := detectThreeInside(candles); p != nil {
		patterns = append(patterns, *p)
	}
	
	// Three Outside Up/Down
	if p := detectThreeOutside(candles); p != nil {
		patterns = append(patterns, *p)
	}
	
	// Kicker Pattern
	if p := detectKicker(candles); p != nil {
		patterns = append(patterns, *p)
	}
	
	// Island Reversal
	if p := detectIslandReversal(candles); p != nil {
		patterns = append(patterns, *p)
	}
	
	// Hook Reversal
	if p := detectHookReversal(candles); p != nil {
		patterns = append(patterns, *p)
	}
	
	// San-Ku (Three Gaps)
	if p := detectSanKu(candles); p != nil {
		patterns = append(patterns, *p)
	}
	
	// Squeeze Alert
	if p := detectSqueezeAlert(candles); p != nil {
		patterns = append(patterns, *p)
	}
	
	return patterns
}

// detectThreeWhiteSoldiers detects Three White Soldiers pattern
func detectThreeWhiteSoldiers(candles []Candle) *AdvancedPattern {
	if len(candles) < 5 {
		return nil
	}
	
	c1 := candles[len(candles)-3]
	c2 := candles[len(candles)-2]
	c3 := candles[len(candles)-1]
	
	// All three must be bullish
	if c1.Close <= c1.Open || c2.Close <= c2.Open || c3.Close <= c3.Open {
		return nil
	}
	
	// Each opens within previous body and closes higher
	if c2.Open < c1.Open || c2.Open > c1.Close {
		return nil
	}
	if c3.Open < c2.Open || c3.Open > c2.Close {
		return nil
	}
	
	// Each closes higher than previous
	if c2.Close <= c1.Close || c3.Close <= c2.Close {
		return nil
	}
	
	// Small upper wicks
	body1 := c1.Close - c1.Open
	body2 := c2.Close - c2.Open
	body3 := c3.Close - c3.Open
	
	wick1 := c1.High - c1.Close
	wick2 := c2.High - c2.Close
	wick3 := c3.High - c3.Close
	
	if wick1 > body1*0.3 || wick2 > body2*0.3 || wick3 > body3*0.3 {
		return nil
	}
	
	return &AdvancedPattern{
		Name:        "Three White Soldiers",
		Type:        "bullish",
		Strength:    90,
		Reliability: 85,
		Description: "Strong bullish reversal - three consecutive bullish candles",
		Entry:       c3.Close,
		StopLoss:    c1.Low,
		Target:      c3.Close + (c3.Close - c1.Low) * 2,
		CandleIdx:   len(candles) - 1,
	}
}

// detectThreeBlackCrows detects Three Black Crows pattern
func detectThreeBlackCrows(candles []Candle) *AdvancedPattern {
	if len(candles) < 5 {
		return nil
	}
	
	c1 := candles[len(candles)-3]
	c2 := candles[len(candles)-2]
	c3 := candles[len(candles)-1]
	
	// All three must be bearish
	if c1.Close >= c1.Open || c2.Close >= c2.Open || c3.Close >= c3.Open {
		return nil
	}
	
	// Each opens within previous body and closes lower
	if c2.Open > c1.Open || c2.Open < c1.Close {
		return nil
	}
	if c3.Open > c2.Open || c3.Open < c2.Close {
		return nil
	}
	
	// Each closes lower than previous
	if c2.Close >= c1.Close || c3.Close >= c2.Close {
		return nil
	}
	
	return &AdvancedPattern{
		Name:        "Three Black Crows",
		Type:        "bearish",
		Strength:    90,
		Reliability: 85,
		Description: "Strong bearish reversal - three consecutive bearish candles",
		Entry:       c3.Close,
		StopLoss:    c1.High,
		Target:      c3.Close - (c1.High - c3.Close) * 2,
		CandleIdx:   len(candles) - 1,
	}
}

// detectRisingThreeMethods detects Rising Three Methods
func detectRisingThreeMethods(candles []Candle) *AdvancedPattern {
	if len(candles) < 7 {
		return nil
	}
	
	c1 := candles[len(candles)-5] // First long bullish
	c5 := candles[len(candles)-1] // Last long bullish
	
	// First and last must be bullish
	if c1.Close <= c1.Open || c5.Close <= c5.Open {
		return nil
	}
	
	// Middle candles should be small and within first candle's range
	for i := len(candles) - 4; i < len(candles)-1; i++ {
		c := candles[i]
		if c.High > c1.High || c.Low < c1.Low {
			return nil
		}
	}
	
	// Last candle closes above first
	if c5.Close <= c1.Close {
		return nil
	}
	
	return &AdvancedPattern{
		Name:        "Rising Three Methods",
		Type:        "bullish",
		Strength:    85,
		Reliability: 80,
		Description: "Bullish continuation - consolidation within uptrend",
		Entry:       c5.Close,
		StopLoss:    c1.Low,
		Target:      c5.Close + (c5.Close - c1.Low),
		CandleIdx:   len(candles) - 1,
	}
}

// detectFallingThreeMethods detects Falling Three Methods
func detectFallingThreeMethods(candles []Candle) *AdvancedPattern {
	if len(candles) < 7 {
		return nil
	}
	
	c1 := candles[len(candles)-5]
	c5 := candles[len(candles)-1]
	
	// First and last must be bearish
	if c1.Close >= c1.Open || c5.Close >= c5.Open {
		return nil
	}
	
	// Middle candles within first candle's range
	for i := len(candles) - 4; i < len(candles)-1; i++ {
		c := candles[i]
		if c.High > c1.High || c.Low < c1.Low {
			return nil
		}
	}
	
	// Last candle closes below first
	if c5.Close >= c1.Close {
		return nil
	}
	
	return &AdvancedPattern{
		Name:        "Falling Three Methods",
		Type:        "bearish",
		Strength:    85,
		Reliability: 80,
		Description: "Bearish continuation - consolidation within downtrend",
		Entry:       c5.Close,
		StopLoss:    c1.High,
		Target:      c5.Close - (c1.High - c5.Close),
		CandleIdx:   len(candles) - 1,
	}
}


// detectAbandonedBaby detects Abandoned Baby pattern
func detectAbandonedBaby(candles []Candle) *AdvancedPattern {
	if len(candles) < 5 {
		return nil
	}
	
	c1 := candles[len(candles)-3]
	c2 := candles[len(candles)-2] // Doji with gap
	c3 := candles[len(candles)-1]
	
	// Middle candle should be doji-like
	body2 := math.Abs(c2.Close - c2.Open)
	range2 := c2.High - c2.Low
	if range2 == 0 || body2/range2 > 0.1 {
		return nil
	}
	
	// Bullish Abandoned Baby
	if c1.Close < c1.Open && c3.Close > c3.Open {
		// Gap down then gap up
		if c2.High < c1.Low && c2.High < c3.Low {
			return &AdvancedPattern{
				Name:        "Bullish Abandoned Baby",
				Type:        "bullish",
				Strength:    95,
				Reliability: 90,
				Description: "Rare and powerful bullish reversal",
				Entry:       c3.Close,
				StopLoss:    c2.Low,
				Target:      c3.Close + (c3.Close - c2.Low) * 2,
				CandleIdx:   len(candles) - 1,
			}
		}
	}
	
	// Bearish Abandoned Baby
	if c1.Close > c1.Open && c3.Close < c3.Open {
		// Gap up then gap down
		if c2.Low > c1.High && c2.Low > c3.High {
			return &AdvancedPattern{
				Name:        "Bearish Abandoned Baby",
				Type:        "bearish",
				Strength:    95,
				Reliability: 90,
				Description: "Rare and powerful bearish reversal",
				Entry:       c3.Close,
				StopLoss:    c2.High,
				Target:      c3.Close - (c2.High - c3.Close) * 2,
				CandleIdx:   len(candles) - 1,
			}
		}
	}
	
	return nil
}

// detectThreeInside detects Three Inside Up/Down
func detectThreeInside(candles []Candle) *AdvancedPattern {
	if len(candles) < 5 {
		return nil
	}
	
	c1 := candles[len(candles)-3]
	c2 := candles[len(candles)-2]
	c3 := candles[len(candles)-1]
	
	// Three Inside Up (bullish)
	if c1.Close < c1.Open { // First bearish
		// Second is bullish harami
		if c2.Close > c2.Open && c2.Open > c1.Close && c2.Close < c1.Open {
			// Third confirms
			if c3.Close > c3.Open && c3.Close > c1.Open {
				return &AdvancedPattern{
					Name:        "Three Inside Up",
					Type:        "bullish",
					Strength:    82,
					Reliability: 78,
					Description: "Bullish reversal confirmation",
					Entry:       c3.Close,
					StopLoss:    c1.Low,
					Target:      c3.Close + (c3.Close - c1.Low) * 1.5,
					CandleIdx:   len(candles) - 1,
				}
			}
		}
	}
	
	// Three Inside Down (bearish)
	if c1.Close > c1.Open { // First bullish
		// Second is bearish harami
		if c2.Close < c2.Open && c2.Open < c1.Close && c2.Close > c1.Open {
			// Third confirms
			if c3.Close < c3.Open && c3.Close < c1.Open {
				return &AdvancedPattern{
					Name:        "Three Inside Down",
					Type:        "bearish",
					Strength:    82,
					Reliability: 78,
					Description: "Bearish reversal confirmation",
					Entry:       c3.Close,
					StopLoss:    c1.High,
					Target:      c3.Close - (c1.High - c3.Close) * 1.5,
					CandleIdx:   len(candles) - 1,
				}
			}
		}
	}
	
	return nil
}

// detectThreeOutside detects Three Outside Up/Down
func detectThreeOutside(candles []Candle) *AdvancedPattern {
	if len(candles) < 5 {
		return nil
	}
	
	c1 := candles[len(candles)-3]
	c2 := candles[len(candles)-2]
	c3 := candles[len(candles)-1]
	
	// Three Outside Up (bullish)
	if c1.Close < c1.Open { // First bearish
		// Second is bullish engulfing
		if c2.Close > c2.Open && c2.Open < c1.Close && c2.Close > c1.Open {
			// Third confirms
			if c3.Close > c3.Open && c3.Close > c2.Close {
				return &AdvancedPattern{
					Name:        "Three Outside Up",
					Type:        "bullish",
					Strength:    85,
					Reliability: 80,
					Description: "Strong bullish reversal with confirmation",
					Entry:       c3.Close,
					StopLoss:    c1.Low,
					Target:      c3.Close + (c3.Close - c1.Low) * 2,
					CandleIdx:   len(candles) - 1,
				}
			}
		}
	}
	
	// Three Outside Down (bearish)
	if c1.Close > c1.Open { // First bullish
		// Second is bearish engulfing
		if c2.Close < c2.Open && c2.Open > c1.Close && c2.Close < c1.Open {
			// Third confirms
			if c3.Close < c3.Open && c3.Close < c2.Close {
				return &AdvancedPattern{
					Name:        "Three Outside Down",
					Type:        "bearish",
					Strength:    85,
					Reliability: 80,
					Description: "Strong bearish reversal with confirmation",
					Entry:       c3.Close,
					StopLoss:    c1.High,
					Target:      c3.Close - (c1.High - c3.Close) * 2,
					CandleIdx:   len(candles) - 1,
				}
			}
		}
	}
	
	return nil
}

// detectKicker detects Kicker pattern
func detectKicker(candles []Candle) *AdvancedPattern {
	if len(candles) < 4 {
		return nil
	}
	
	c1 := candles[len(candles)-2]
	c2 := candles[len(candles)-1]
	
	body1 := math.Abs(c1.Close - c1.Open)
	body2 := math.Abs(c2.Close - c2.Open)
	avgBody := (body1 + body2) / 2
	
	// Bullish Kicker
	if c1.Close < c1.Open && c2.Close > c2.Open {
		// Gap up open
		if c2.Open > c1.Open && body2 > avgBody*1.5 {
			return &AdvancedPattern{
				Name:        "Bullish Kicker",
				Type:        "bullish",
				Strength:    92,
				Reliability: 88,
				Description: "Very strong bullish signal - gap reversal",
				Entry:       c2.Close,
				StopLoss:    c1.Low,
				Target:      c2.Close + (c2.Close - c1.Low) * 2,
				CandleIdx:   len(candles) - 1,
			}
		}
	}
	
	// Bearish Kicker
	if c1.Close > c1.Open && c2.Close < c2.Open {
		// Gap down open
		if c2.Open < c1.Open && body2 > avgBody*1.5 {
			return &AdvancedPattern{
				Name:        "Bearish Kicker",
				Type:        "bearish",
				Strength:    92,
				Reliability: 88,
				Description: "Very strong bearish signal - gap reversal",
				Entry:       c2.Close,
				StopLoss:    c1.High,
				Target:      c2.Close - (c1.High - c2.Close) * 2,
				CandleIdx:   len(candles) - 1,
			}
		}
	}
	
	return nil
}

// detectIslandReversal detects Island Reversal
func detectIslandReversal(candles []Candle) *AdvancedPattern {
	if len(candles) < 6 {
		return nil
	}
	
	// Look for gap up, consolidation, gap down (or vice versa)
	for i := len(candles) - 5; i < len(candles)-2; i++ {
		prev := candles[i-1]
		island := candles[i]
		next := candles[i+1]
		
		// Bearish Island (gap up then gap down)
		if island.Low > prev.High && next.High < island.Low {
			return &AdvancedPattern{
				Name:        "Bearish Island Reversal",
				Type:        "bearish",
				Strength:    88,
				Reliability: 85,
				Description: "Isolated price action - strong reversal",
				Entry:       candles[len(candles)-1].Close,
				StopLoss:    island.High,
				Target:      candles[len(candles)-1].Close - (island.High - candles[len(candles)-1].Close),
				CandleIdx:   len(candles) - 1,
			}
		}
		
		// Bullish Island (gap down then gap up)
		if island.High < prev.Low && next.Low > island.High {
			return &AdvancedPattern{
				Name:        "Bullish Island Reversal",
				Type:        "bullish",
				Strength:    88,
				Reliability: 85,
				Description: "Isolated price action - strong reversal",
				Entry:       candles[len(candles)-1].Close,
				StopLoss:    island.Low,
				Target:      candles[len(candles)-1].Close + (candles[len(candles)-1].Close - island.Low),
				CandleIdx:   len(candles) - 1,
			}
		}
	}
	
	return nil
}

// detectHookReversal detects Hook Reversal
func detectHookReversal(candles []Candle) *AdvancedPattern {
	if len(candles) < 4 {
		return nil
	}
	
	c1 := candles[len(candles)-2]
	c2 := candles[len(candles)-1]
	
	// Bullish Hook
	if c1.Close < c1.Open && c2.Close > c2.Open {
		if c2.Low < c1.Low && c2.Close > c1.Open {
			return &AdvancedPattern{
				Name:        "Bullish Hook Reversal",
				Type:        "bullish",
				Strength:    78,
				Reliability: 72,
				Description: "Lower low followed by strong close",
				Entry:       c2.Close,
				StopLoss:    c2.Low,
				Target:      c2.Close + (c2.Close - c2.Low) * 2,
				CandleIdx:   len(candles) - 1,
			}
		}
	}
	
	// Bearish Hook
	if c1.Close > c1.Open && c2.Close < c2.Open {
		if c2.High > c1.High && c2.Close < c1.Open {
			return &AdvancedPattern{
				Name:        "Bearish Hook Reversal",
				Type:        "bearish",
				Strength:    78,
				Reliability: 72,
				Description: "Higher high followed by strong close down",
				Entry:       c2.Close,
				StopLoss:    c2.High,
				Target:      c2.Close - (c2.High - c2.Close) * 2,
				CandleIdx:   len(candles) - 1,
			}
		}
	}
	
	return nil
}

// detectSanKu detects San-Ku (Three Gaps)
func detectSanKu(candles []Candle) *AdvancedPattern {
	if len(candles) < 6 {
		return nil
	}
	
	gapCount := 0
	gapUp := true
	
	for i := len(candles) - 5; i < len(candles)-1; i++ {
		c1 := candles[i]
		c2 := candles[i+1]
		
		// Gap up
		if c2.Low > c1.High {
			if gapUp {
				gapCount++
			} else {
				gapCount = 1
				gapUp = true
			}
		}
		
		// Gap down
		if c2.High < c1.Low {
			if !gapUp {
				gapCount++
			} else {
				gapCount = 1
				gapUp = false
			}
		}
	}
	
	if gapCount >= 3 {
		if gapUp {
			return &AdvancedPattern{
				Name:        "San-Ku (Three Gaps Up)",
				Type:        "bearish", // Exhaustion signal
				Strength:    80,
				Reliability: 75,
				Description: "Three consecutive gaps up - exhaustion likely",
				Entry:       candles[len(candles)-1].Close,
				StopLoss:    candles[len(candles)-1].High,
				Target:      candles[len(candles)-1].Close - (candles[len(candles)-1].High - candles[len(candles)-1].Close) * 2,
				CandleIdx:   len(candles) - 1,
			}
		} else {
			return &AdvancedPattern{
				Name:        "San-Ku (Three Gaps Down)",
				Type:        "bullish", // Exhaustion signal
				Strength:    80,
				Reliability: 75,
				Description: "Three consecutive gaps down - exhaustion likely",
				Entry:       candles[len(candles)-1].Close,
				StopLoss:    candles[len(candles)-1].Low,
				Target:      candles[len(candles)-1].Close + (candles[len(candles)-1].Close - candles[len(candles)-1].Low) * 2,
				CandleIdx:   len(candles) - 1,
			}
		}
	}
	
	return nil
}

// detectSqueezeAlert detects Bollinger Band squeeze
func detectSqueezeAlert(candles []Candle) *AdvancedPattern {
	if len(candles) < 20 {
		return nil
	}
	
	// Calculate Bollinger Band width
	_, bbPercentile := CalculateBollingerBandWidth(candles, 20, 2.0)
	
	// Squeeze = very narrow bands
	if bbPercentile < 10 {
		// Determine likely direction from recent price action
		recentUp := candles[len(candles)-1].Close > candles[len(candles)-5].Close
		
		if recentUp {
			return &AdvancedPattern{
				Name:        "Squeeze Alert (Bullish)",
				Type:        "bullish",
				Strength:    75,
				Reliability: 70,
				Description: "Bollinger squeeze - breakout imminent",
				Entry:       candles[len(candles)-1].Close,
				StopLoss:    candles[len(candles)-1].Low,
				Target:      candles[len(candles)-1].Close * 1.03, // 3% target
				CandleIdx:   len(candles) - 1,
			}
		} else {
			return &AdvancedPattern{
				Name:        "Squeeze Alert (Bearish)",
				Type:        "bearish",
				Strength:    75,
				Reliability: 70,
				Description: "Bollinger squeeze - breakout imminent",
				Entry:       candles[len(candles)-1].Close,
				StopLoss:    candles[len(candles)-1].High,
				Target:      candles[len(candles)-1].Close * 0.97, // 3% target
				CandleIdx:   len(candles) - 1,
			}
		}
	}
	
	return nil
}
