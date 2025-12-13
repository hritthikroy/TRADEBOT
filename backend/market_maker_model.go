package main

import (
	"math"
)

// ==================== MARKET MAKER MODEL ====================
// Detect institutional manipulation patterns

// StopHunt represents a stop hunt pattern
type StopHunt struct {
	Type          string  // "bullish" or "bearish"
	SweepPrice    float64 // Price where stops were swept
	ReversalPrice float64 // Price where reversal occurred
	Volume        float64 // Volume on reversal
	Strength      float64 // Pattern strength (0-100)
	CandleIdx     int
	Confirmed     bool
}

// LiquidityGrab represents a liquidity grab pattern
type LiquidityGrab struct {
	Type       string  // "buyside" or "sellside"
	GrabPrice  float64 // Price where liquidity was grabbed
	ReturnPrice float64 // Price returned to after grab
	WickSize   float64 // Size of the wick
	Strength   float64 // Pattern strength (0-100)
	CandleIdx  int
}

// ManipulationPhase represents accumulation/manipulation/distribution
type ManipulationPhase struct {
	Phase      string  // "accumulation", "manipulation", "distribution", "markup", "markdown"
	StartIdx   int
	EndIdx     int
	Confidence float64
}

// MMM Analysis holds all market maker analysis
type MMMAnalysis struct {
	StopHunts       []StopHunt
	LiquidityGrabs  []LiquidityGrab
	CurrentPhase    *ManipulationPhase
	TrapDetected    bool
	InstitutionalBias string // "bullish", "bearish", "neutral"
}

// ==================== STOP HUNT DETECTION ====================

// DetectStopHunts identifies stop hunt patterns
func DetectStopHunts(candles []Candle) []StopHunt {
	var stopHunts []StopHunt
	
	if len(candles) < 20 {
		return stopHunts
	}
	
	// Find recent swing highs and lows (potential stop locations)
	for i := 10; i < len(candles)-3; i++ {
		curr := candles[i]
		
		// Find recent swing high (last 10 candles)
		swingHigh := 0.0
		for j := i - 10; j < i; j++ {
			if candles[j].High > swingHigh {
				swingHigh = candles[j].High
			}
		}
		
		// Find recent swing low
		swingLow := candles[i-10].Low
		for j := i - 10; j < i; j++ {
			if candles[j].Low < swingLow {
				swingLow = candles[j].Low
			}
		}
		
		// BULLISH STOP HUNT
		// Price sweeps below swing low, then reverses up
		if curr.Low < swingLow && curr.Close > curr.Open {
			// Check for strong reversal in next 1-3 candles
			reversalStrength := 0.0
			for j := i + 1; j < i + 4 && j < len(candles); j++ {
				if candles[j].Close > swingLow {
					reversalStrength += (candles[j].Close - swingLow) / (swingHigh - swingLow) * 100
				}
			}
			
			if reversalStrength > 30 {
				wickSize := curr.Open - curr.Low
				bodySize := curr.Close - curr.Open
				
				strength := 50.0
				
				// Long wick below = strong rejection
				if wickSize > bodySize*2 {
					strength += 20
				}
				
				// High volume = institutional activity
				avgVolume := 0.0
				for j := i - 10; j < i; j++ {
					avgVolume += candles[j].Volume
				}
				avgVolume /= 10
				
				if curr.Volume > avgVolume*1.5 {
					strength += 15
				}
				
				// Strong close above open
				if curr.Close > curr.Open {
					strength += 10
				}
				
				// Reversal strength
				strength += reversalStrength * 0.3
				
				stopHunt := StopHunt{
					Type:          "bullish",
					SweepPrice:    curr.Low,
					ReversalPrice: curr.Close,
					Volume:        curr.Volume,
					Strength:      math.Min(strength, 100),
					CandleIdx:     i,
					Confirmed:     reversalStrength > 50,
				}
				
				stopHunts = append(stopHunts, stopHunt)
			}
		}
		
		// BEARISH STOP HUNT
		// Price sweeps above swing high, then reverses down
		if curr.High > swingHigh && curr.Close < curr.Open {
			reversalStrength := 0.0
			for j := i + 1; j < i + 4 && j < len(candles); j++ {
				if candles[j].Close < swingHigh {
					reversalStrength += (swingHigh - candles[j].Close) / (swingHigh - swingLow) * 100
				}
			}
			
			if reversalStrength > 30 {
				wickSize := curr.High - curr.Open
				bodySize := curr.Open - curr.Close
				
				strength := 50.0
				
				if wickSize > bodySize*2 {
					strength += 20
				}
				
				avgVolume := 0.0
				for j := i - 10; j < i; j++ {
					avgVolume += candles[j].Volume
				}
				avgVolume /= 10
				
				if curr.Volume > avgVolume*1.5 {
					strength += 15
				}
				
				if curr.Close < curr.Open {
					strength += 10
				}
				
				strength += reversalStrength * 0.3
				
				stopHunt := StopHunt{
					Type:          "bearish",
					SweepPrice:    curr.High,
					ReversalPrice: curr.Close,
					Volume:        curr.Volume,
					Strength:      math.Min(strength, 100),
					CandleIdx:     i,
					Confirmed:     reversalStrength > 50,
				}
				
				stopHunts = append(stopHunts, stopHunt)
			}
		}
	}
	
	return stopHunts
}

// ==================== LIQUIDITY GRAB DETECTION ====================

// DetectLiquidityGrabs identifies liquidity grab patterns
func DetectLiquidityGrabs(candles []Candle) []LiquidityGrab {
	var grabs []LiquidityGrab
	
	if len(candles) < 10 {
		return grabs
	}
	
	for i := 5; i < len(candles)-1; i++ {
		curr := candles[i]
		next := candles[i+1]
		
		// Calculate average body size
		avgBody := 0.0
		for j := i - 5; j < i; j++ {
			avgBody += math.Abs(candles[j].Close - candles[j].Open)
		}
		avgBody /= 5
		
		// BUYSIDE LIQUIDITY GRAB
		// Long upper wick, close back inside range
		upperWick := curr.High - math.Max(curr.Open, curr.Close)
		lowerWick := math.Min(curr.Open, curr.Close) - curr.Low
		body := math.Abs(curr.Close - curr.Open)
		
		// Upper wick > 2x body, close below high
		if upperWick > body*2 && upperWick > avgBody*1.5 {
			// Check if next candle confirms rejection
			if next.Close < curr.High-upperWick*0.5 {
				strength := 50.0
				
				// Wick size relative to body
				wickRatio := upperWick / body
				strength += math.Min(wickRatio*10, 30)
				
				// Volume confirmation
				avgVolume := 0.0
				for j := i - 5; j < i; j++ {
					avgVolume += candles[j].Volume
				}
				avgVolume /= 5
				
				if curr.Volume > avgVolume*1.3 {
					strength += 15
				}
				
				// Next candle confirms
				if next.Close < curr.Close {
					strength += 10
				}
				
				grab := LiquidityGrab{
					Type:        "buyside",
					GrabPrice:   curr.High,
					ReturnPrice: curr.Close,
					WickSize:    upperWick,
					Strength:    math.Min(strength, 100),
					CandleIdx:   i,
				}
				
				grabs = append(grabs, grab)
			}
		}
		
		// SELLSIDE LIQUIDITY GRAB
		// Long lower wick, close back inside range
		if lowerWick > body*2 && lowerWick > avgBody*1.5 {
			if next.Close > curr.Low+lowerWick*0.5 {
				strength := 50.0
				
				wickRatio := lowerWick / body
				strength += math.Min(wickRatio*10, 30)
				
				avgVolume := 0.0
				for j := i - 5; j < i; j++ {
					avgVolume += candles[j].Volume
				}
				avgVolume /= 5
				
				if curr.Volume > avgVolume*1.3 {
					strength += 15
				}
				
				if next.Close > curr.Close {
					strength += 10
				}
				
				grab := LiquidityGrab{
					Type:        "sellside",
					GrabPrice:   curr.Low,
					ReturnPrice: curr.Close,
					WickSize:    lowerWick,
					Strength:    math.Min(strength, 100),
					CandleIdx:   i,
				}
				
				grabs = append(grabs, grab)
			}
		}
	}
	
	return grabs
}

// ==================== MANIPULATION PHASE DETECTION ====================

// DetectManipulationPhase identifies current market phase
func DetectManipulationPhase(candles []Candle) *ManipulationPhase {
	if len(candles) < 50 {
		return nil
	}
	
	// Analyze last 50 candles for phase detection
	recentCandles := candles[len(candles)-50:]
	
	// Calculate range and volatility
	high := recentCandles[0].High
	low := recentCandles[0].Low
	
	for _, c := range recentCandles {
		if c.High > high {
			high = c.High
		}
		if c.Low < low {
			low = c.Low
		}
	}
	
	rangeSize := high - low
	currentPrice := candles[len(candles)-1].Close
	
	// Calculate position in range
	position := (currentPrice - low) / rangeSize
	
	// Calculate volatility (ATR)
	atr := calculateATR(recentCandles, 14)
	avgATR := atr / currentPrice * 100
	
	// Calculate volume trend
	firstHalfVolume := 0.0
	secondHalfVolume := 0.0
	
	for i := 0; i < 25; i++ {
		firstHalfVolume += recentCandles[i].Volume
	}
	for i := 25; i < 50; i++ {
		secondHalfVolume += recentCandles[i].Volume
	}
	
	volumeIncrease := (secondHalfVolume - firstHalfVolume) / firstHalfVolume
	
	// Determine phase
	phase := &ManipulationPhase{}
	
	// ACCUMULATION: Low volatility, ranging, increasing volume
	if avgATR < 2.0 && position > 0.3 && position < 0.7 && volumeIncrease > 0.1 {
		phase.Phase = "accumulation"
		phase.Confidence = 70 + volumeIncrease*100
		
	// MANIPULATION: High volatility, stop hunts, fake breakouts
	} else if avgATR > 3.0 {
		phase.Phase = "manipulation"
		phase.Confidence = 60 + (avgATR-3.0)*10
		
	// MARKUP: Strong uptrend, higher highs, increasing volume
	} else if position > 0.7 && volumeIncrease > 0.2 {
		phase.Phase = "markup"
		phase.Confidence = 75 + volumeIncrease*50
		
	// MARKDOWN: Strong downtrend, lower lows, increasing volume
	} else if position < 0.3 && volumeIncrease > 0.2 {
		phase.Phase = "markdown"
		phase.Confidence = 75 + volumeIncrease*50
		
	// DISTRIBUTION: High prices, ranging, decreasing volume
	} else if position > 0.7 && volumeIncrease < -0.1 {
		phase.Phase = "distribution"
		phase.Confidence = 70 + math.Abs(volumeIncrease)*100
		
	} else {
		phase.Phase = "ranging"
		phase.Confidence = 50
	}
	
	phase.StartIdx = len(candles) - 50
	phase.EndIdx = len(candles) - 1
	phase.Confidence = math.Min(phase.Confidence, 95)
	
	return phase
}

// ==================== FULL MMM ANALYSIS ====================

// PerformMMMAnalysis performs complete market maker analysis
func PerformMMMAnalysis(candles []Candle) *MMMAnalysis {
	analysis := &MMMAnalysis{
		InstitutionalBias: "neutral",
	}
	
	// Detect stop hunts
	analysis.StopHunts = DetectStopHunts(candles)
	
	// Detect liquidity grabs
	analysis.LiquidityGrabs = DetectLiquidityGrabs(candles)
	
	// Detect manipulation phase
	analysis.CurrentPhase = DetectManipulationPhase(candles)
	
	// Determine institutional bias
	bullishSignals := 0
	bearishSignals := 0
	
	// Count recent stop hunts (last 20 candles)
	for _, sh := range analysis.StopHunts {
		if sh.CandleIdx >= len(candles)-20 {
			if sh.Type == "bullish" && sh.Confirmed {
				bullishSignals += 2
			} else if sh.Type == "bearish" && sh.Confirmed {
				bearishSignals += 2
			}
		}
	}
	
	// Count recent liquidity grabs
	for _, lg := range analysis.LiquidityGrabs {
		if lg.CandleIdx >= len(candles)-20 {
			if lg.Type == "sellside" && lg.Strength > 70 {
				bullishSignals++
			} else if lg.Type == "buyside" && lg.Strength > 70 {
				bearishSignals++
			}
		}
	}
	
	// Phase bias
	if analysis.CurrentPhase != nil {
		if analysis.CurrentPhase.Phase == "accumulation" || analysis.CurrentPhase.Phase == "markup" {
			bullishSignals += 2
		} else if analysis.CurrentPhase.Phase == "distribution" || analysis.CurrentPhase.Phase == "markdown" {
			bearishSignals += 2
		}
	}
	
	// Determine bias
	if bullishSignals > bearishSignals+1 {
		analysis.InstitutionalBias = "bullish"
	} else if bearishSignals > bullishSignals+1 {
		analysis.InstitutionalBias = "bearish"
	}
	
	// Detect trap (manipulation phase with conflicting signals)
	if analysis.CurrentPhase != nil && analysis.CurrentPhase.Phase == "manipulation" {
		analysis.TrapDetected = true
	}
	
	return analysis
}

// GetMMMSignalStrength returns signal strength based on MMM analysis
func GetMMMSignalStrength(analysis *MMMAnalysis, direction string) float64 {
	if analysis == nil {
		return 0
	}
	
	strength := 0.0
	
	// Institutional bias alignment
	if analysis.InstitutionalBias == direction {
		strength += 25
	} else if analysis.InstitutionalBias != "neutral" {
		strength -= 15 // Against institutional bias
	}
	
	// Recent stop hunts
	for _, sh := range analysis.StopHunts {
		if sh.Type == direction && sh.Confirmed {
			strength += sh.Strength * 0.2
		}
	}
	
	// Recent liquidity grabs
	for _, lg := range analysis.LiquidityGrabs {
		if (direction == "bullish" && lg.Type == "sellside") ||
			(direction == "bearish" && lg.Type == "buyside") {
			strength += lg.Strength * 0.15
		}
	}
	
	// Phase alignment
	if analysis.CurrentPhase != nil {
		if (direction == "bullish" && (analysis.CurrentPhase.Phase == "accumulation" || analysis.CurrentPhase.Phase == "markup")) ||
			(direction == "bearish" && (analysis.CurrentPhase.Phase == "distribution" || analysis.CurrentPhase.Phase == "markdown")) {
			strength += analysis.CurrentPhase.Confidence * 0.2
		}
	}
	
	// Penalty for trap detection
	if analysis.TrapDetected {
		strength -= 20
	}
	
	return math.Max(0, math.Min(strength, 100))
}
