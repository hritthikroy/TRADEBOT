package ict

import (
	"math"
	"time"
)

// ==================== SESSION LIQUIDITY MAPPING ====================
// Track liquidity across trading sessions

// SessionLiquidity represents liquidity for a trading session
type SessionLiquidity struct {
	Session     string  // "Asian", "London", "NewYork"
	High        float64
	Low         float64
	Open        float64
	Close       float64
	Range       float64
	Volume      float64
	Swept       bool
	SweptSide   string // "high", "low", "both", "none"
	Strength    float64
}

// LiquidityMap holds all session liquidity data
type LiquidityMap struct {
	Asian    *SessionLiquidity
	London   *SessionLiquidity
	NewYork  *SessionLiquidity
	Current  *SessionLiquidity
	Previous *SessionLiquidity
}

// ==================== SESSION DETECTION ====================

// GetCurrentSession returns the current trading session
func GetCurrentSession(t time.Time) string {
	hour := t.UTC().Hour()
	
	// Asian Session: 00:00-08:00 UTC
	if hour >= 0 && hour < 8 {
		return "Asian"
	}
	
	// London Session: 08:00-16:00 UTC
	if hour >= 8 && hour < 16 {
		return "London"
	}
	
	// New York Session: 13:00-21:00 UTC
	if hour >= 13 && hour < 21 {
		return "NewYork"
	}
	
	// Overlap: 13:00-16:00 UTC
	if hour >= 13 && hour < 16 {
		return "Overlap"
	}
	
	return "OffHours"
}

// IsHighVolatilitySession checks if current session is high volatility
func IsHighVolatilitySession(t time.Time) bool {
	hour := t.UTC().Hour()
	
	// London Open: 08:00-10:00 UTC (highest volatility)
	if hour >= 8 && hour < 10 {
		return true
	}
	
	// New York Open: 13:00-15:00 UTC (high volatility)
	if hour >= 13 && hour < 15 {
		return true
	}
	
	// London-NY Overlap: 13:00-16:00 UTC (very high volatility)
	if hour >= 13 && hour < 16 {
		return true
	}
	
	return false
}

// ==================== LIQUIDITY MAPPING ====================

// MapSessionLiquidity maps liquidity across sessions
func MapSessionLiquidity(candles []Candle) *LiquidityMap {
	lm := &LiquidityMap{}
	
	if len(candles) < 24 {
		return lm
	}
	
	// Group candles by session
	asianCandles := []Candle{}
	londonCandles := []Candle{}
	nyCandles := []Candle{}
	
	for _, c := range candles {
		t := time.Unix(c.Timestamp/1000, 0).UTC()
		session := GetCurrentSession(t)
		
		switch session {
		case "Asian":
			asianCandles = append(asianCandles, c)
		case "London":
			londonCandles = append(londonCandles, c)
		case "NewYork", "Overlap":
			nyCandles = append(nyCandles, c)
		}
	}
	
	// Calculate Asian session liquidity
	if len(asianCandles) > 0 {
		lm.Asian = calculateSessionLiquidity(asianCandles, "Asian")
	}
	
	// Calculate London session liquidity
	if len(londonCandles) > 0 {
		lm.London = calculateSessionLiquidity(londonCandles, "London")
	}
	
	// Calculate New York session liquidity
	if len(nyCandles) > 0 {
		lm.NewYork = calculateSessionLiquidity(nyCandles, "NewYork")
	}
	
	// Set current session
	currentTime := time.Now().UTC()
	currentSession := GetCurrentSession(currentTime)
	
	switch currentSession {
	case "Asian":
		lm.Current = lm.Asian
		lm.Previous = lm.NewYork
	case "London":
		lm.Current = lm.London
		lm.Previous = lm.Asian
	case "NewYork", "Overlap":
		lm.Current = lm.NewYork
		lm.Previous = lm.London
	}
	
	// Check if previous session liquidity was swept
	if lm.Current != nil && lm.Previous != nil {
		currentPrice := candles[len(candles)-1].Close
		
		// Check if swept high
		if currentPrice > lm.Previous.High {
			lm.Previous.Swept = true
			if lm.Previous.SweptSide == "low" {
				lm.Previous.SweptSide = "both"
			} else {
				lm.Previous.SweptSide = "high"
			}
		}
		
		// Check if swept low
		if currentPrice < lm.Previous.Low {
			lm.Previous.Swept = true
			if lm.Previous.SweptSide == "high" {
				lm.Previous.SweptSide = "both"
			} else {
				lm.Previous.SweptSide = "low"
			}
		}
	}
	
	return lm
}

// calculateSessionLiquidity calculates liquidity for a session
func calculateSessionLiquidity(candles []Candle, session string) *SessionLiquidity {
	if len(candles) == 0 {
		return nil
	}
	
	sl := &SessionLiquidity{
		Session: session,
		Open:    candles[0].Open,
		Close:   candles[len(candles)-1].Close,
		High:    candles[0].High,
		Low:     candles[0].Low,
		Volume:  0,
		Swept:   false,
		SweptSide: "none",
	}
	
	// Find high, low, and total volume
	for _, c := range candles {
		if c.High > sl.High {
			sl.High = c.High
		}
		if c.Low < sl.Low {
			sl.Low = c.Low
		}
		sl.Volume += c.Volume
	}
	
	sl.Range = sl.High - sl.Low
	
	// Calculate strength based on range and volume
	avgVolume := sl.Volume / float64(len(candles))
	sl.Strength = 50.0
	
	// Large range = strong liquidity
	avgCandleRange := 0.0
	for _, c := range candles {
		avgCandleRange += c.High - c.Low
	}
	avgCandleRange /= float64(len(candles))
	
	if sl.Range > avgCandleRange*3 {
		sl.Strength += 25
	}
	
	// High volume = strong liquidity
	if avgVolume > 0 {
		sl.Strength += 15
	}
	
	// Session-specific strength
	switch session {
	case "London":
		sl.Strength += 10 // London is high liquidity
	case "NewYork":
		sl.Strength += 15 // NY is highest liquidity
	case "Asian":
		sl.Strength -= 10 // Asian is lower liquidity
	}
	
	sl.Strength = math.Max(0, math.Min(sl.Strength, 100))
	
	return sl
}

// GetSessionLiquiditySignal returns trading signal based on session liquidity
func GetSessionLiquiditySignal(lm *LiquidityMap, currentPrice float64) (string, float64) {
	if lm == nil || lm.Previous == nil {
		return "neutral", 0
	}
	
	strength := 0.0
	direction := "neutral"
	
	// Liquidity sweep = reversal opportunity
	if lm.Previous.Swept {
		if lm.Previous.SweptSide == "high" {
			// Swept high = bearish reversal
			direction = "bearish"
			strength = lm.Previous.Strength * 0.7
			
		} else if lm.Previous.SweptSide == "low" {
			// Swept low = bullish reversal
			direction = "bullish"
			strength = lm.Previous.Strength * 0.7
			
		} else if lm.Previous.SweptSide == "both" {
			// Both sides swept = strong reversal
			if currentPrice > (lm.Previous.High+lm.Previous.Low)/2 {
				direction = "bearish"
			} else {
				direction = "bullish"
			}
			strength = lm.Previous.Strength * 0.9
		}
	}
	
	// Near session high/low = potential reversal
	if lm.Current != nil {
		distToHigh := math.Abs(currentPrice - lm.Current.High)
		distToLow := math.Abs(currentPrice - lm.Current.Low)
		threshold := lm.Current.Range * 0.1
		
		if distToHigh < threshold {
			// Near session high = potential short
			if direction == "bearish" {
				strength += 15 // Confluence
			} else {
				direction = "bearish"
				strength = 40
			}
		}
		
		if distToLow < threshold {
			// Near session low = potential long
			if direction == "bullish" {
				strength += 15 // Confluence
			} else {
				direction = "bullish"
				strength = 40
			}
		}
	}
	
	return direction, math.Min(strength, 100)
}

// ShouldTradeSession determines if we should trade in current session
func ShouldTradeSession(t time.Time) bool {
	hour := t.UTC().Hour()
	
	// Best trading times
	// London Open: 08:00-12:00 UTC
	if hour >= 8 && hour < 12 {
		return true
	}
	
	// New York Open: 13:00-17:00 UTC
	if hour >= 13 && hour < 17 {
		return true
	}
	
	// Avoid:
	// - Asian session (low liquidity)
	// - Late NY session (low liquidity)
	// - Weekends
	
	weekday := t.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}
	
	// Avoid Friday after 16:00 UTC (profit taking)
	if weekday == time.Friday && hour >= 16 {
		return false
	}
	
	return false
}

// GetSessionMultiplier returns position size multiplier based on session
func GetSessionMultiplier(t time.Time) float64 {
	hour := t.UTC().Hour()
	
	// London-NY Overlap: 13:00-16:00 UTC (best time)
	if hour >= 13 && hour < 16 {
		return 1.5 // Increase position size
	}
	
	// London Open: 08:00-13:00 UTC (good time)
	if hour >= 8 && hour < 13 {
		return 1.2
	}
	
	// NY Session: 16:00-20:00 UTC (good time)
	if hour >= 16 && hour < 20 {
		return 1.2
	}
	
	// Asian Session: 00:00-08:00 UTC (reduce size)
	if hour >= 0 && hour < 8 {
		return 0.7
	}
	
	// Off hours
	return 0.5
}
