package institutional

import (
	"math"
)

// ==================== MIRROR MARKET CONCEPT ====================
// Price tends to mirror previous moves - symmetry in markets

// MirrorPattern represents a mirror/symmetry pattern
type MirrorPattern struct {
	Type           string  // "bullish_mirror", "bearish_mirror"
	OriginalMove   float64 // Size of original move
	MirrorMove     float64 // Expected mirror move size
	SymmetryPoint  float64 // Price where symmetry occurs
	TargetPrice    float64 // Expected target based on mirror
	Accuracy       float64 // How accurate the mirror is (0-100)
	CandleIdx      int
}

// MirrorAnalysis holds mirror market analysis
type MirrorAnalysis struct {
	Patterns        []MirrorPattern
	CurrentMirror   *MirrorPattern
	SymmetryLevel   float64
	ExpectedTarget  float64
	MirrorActive    bool
	MirrorDirection string
}

// ==================== MIRROR PATTERN DETECTION ====================

// FindMirrorPatterns identifies mirror/symmetry patterns
func FindMirrorPatterns(candles []Candle) []MirrorPattern {
	var patterns []MirrorPattern
	
	if len(candles) < 50 {
		return patterns
	}
	
	// Find swing points
	swingHighs := []struct {
		Price float64
		Idx   int
	}{}
	swingLows := []struct {
		Price float64
		Idx   int
	}{}
	
	for i := 2; i < len(candles)-2; i++ {
		// Swing High
		if candles[i].High > candles[i-1].High &&
			candles[i].High > candles[i-2].High &&
			candles[i].High > candles[i+1].High &&
			candles[i].High > candles[i+2].High {
			swingHighs = append(swingHighs, struct {
				Price float64
				Idx   int
			}{candles[i].High, i})
		}
		
		// Swing Low
		if candles[i].Low < candles[i-1].Low &&
			candles[i].Low < candles[i-2].Low &&
			candles[i].Low < candles[i+1].Low &&
			candles[i].Low < candles[i+2].Low {
			swingLows = append(swingLows, struct {
				Price float64
				Idx   int
			}{candles[i].Low, i})
		}
	}
	
	// Find V-Bottom Mirror (Bullish)
	// Pattern: Down move → Low → Up move (mirror of down)
	for i := 0; i < len(swingLows)-1; i++ {
		low := swingLows[i]
		
		// Find high before low
		var highBefore *struct {
			Price float64
			Idx   int
		}
		for j := len(swingHighs) - 1; j >= 0; j-- {
			if swingHighs[j].Idx < low.Idx {
				highBefore = &swingHighs[j]
				break
			}
		}
		
		if highBefore == nil {
			continue
		}
		
		// Find high after low
		var highAfter *struct {
			Price float64
			Idx   int
		}
		for j := 0; j < len(swingHighs); j++ {
			if swingHighs[j].Idx > low.Idx {
				highAfter = &swingHighs[j]
				break
			}
		}
		
		if highAfter == nil {
			continue
		}
		
		// Calculate moves
		downMove := highBefore.Price - low.Price
		upMove := highAfter.Price - low.Price
		
		// Check for symmetry (within 20%)
		symmetry := 1.0 - math.Abs(downMove-upMove)/math.Max(downMove, upMove)
		
		if symmetry > 0.7 { // 70%+ symmetry
			pattern := MirrorPattern{
				Type:          "bullish_mirror",
				OriginalMove:  downMove,
				MirrorMove:    upMove,
				SymmetryPoint: low.Price,
				TargetPrice:   low.Price + downMove, // Mirror target
				Accuracy:      symmetry * 100,
				CandleIdx:     low.Idx,
			}
			patterns = append(patterns, pattern)
		}
	}
	
	// Find Inverted V-Top Mirror (Bearish)
	// Pattern: Up move → High → Down move (mirror of up)
	for i := 0; i < len(swingHighs)-1; i++ {
		high := swingHighs[i]
		
		// Find low before high
		var lowBefore *struct {
			Price float64
			Idx   int
		}
		for j := len(swingLows) - 1; j >= 0; j-- {
			if swingLows[j].Idx < high.Idx {
				lowBefore = &swingLows[j]
				break
			}
		}
		
		if lowBefore == nil {
			continue
		}
		
		// Find low after high
		var lowAfter *struct {
			Price float64
			Idx   int
		}
		for j := 0; j < len(swingLows); j++ {
			if swingLows[j].Idx > high.Idx {
				lowAfter = &swingLows[j]
				break
			}
		}
		
		if lowAfter == nil {
			continue
		}
		
		// Calculate moves
		upMove := high.Price - lowBefore.Price
		downMove := high.Price - lowAfter.Price
		
		// Check for symmetry
		symmetry := 1.0 - math.Abs(upMove-downMove)/math.Max(upMove, downMove)
		
		if symmetry > 0.7 {
			pattern := MirrorPattern{
				Type:          "bearish_mirror",
				OriginalMove:  upMove,
				MirrorMove:    downMove,
				SymmetryPoint: high.Price,
				TargetPrice:   high.Price - upMove, // Mirror target
				Accuracy:      symmetry * 100,
				CandleIdx:     high.Idx,
			}
			patterns = append(patterns, pattern)
		}
	}
	
	return patterns
}

// ==================== TIME SYMMETRY ====================

// TimeSymmetry represents time-based symmetry
type TimeSymmetry struct {
	OriginalDuration int     // Candles for original move
	MirrorDuration   int     // Expected candles for mirror
	Accuracy         float64 // Time symmetry accuracy
	ExpectedEnd      int     // Expected candle index for mirror completion
}

// CalculateTimeSymmetry calculates time symmetry for a pattern
func CalculateTimeSymmetry(pattern MirrorPattern, currentIdx int) *TimeSymmetry {
	// Time symmetry: Mirror move takes similar time as original
	
	// This is a simplified calculation
	// In reality, you'd track the actual candle counts
	
	ts := &TimeSymmetry{
		OriginalDuration: 10, // Placeholder
		MirrorDuration:   10,
		Accuracy:         80,
		ExpectedEnd:      currentIdx + 10,
	}
	
	return ts
}

// ==================== MEASURED MOVE ====================

// MeasuredMove represents a measured move pattern
type MeasuredMove struct {
	Type        string  // "bullish" or "bearish"
	FirstLeg    float64 // Size of first leg
	Correction  float64 // Size of correction
	SecondLeg   float64 // Expected second leg (= first leg)
	EntryPrice  float64 // Entry at end of correction
	TargetPrice float64 // Target based on measured move
	StopLoss    float64 // Stop below correction low
	RiskReward  float64
	Confidence  float64
}

// FindMeasuredMoves identifies measured move patterns
func FindMeasuredMoves(candles []Candle) []MeasuredMove {
	var moves []MeasuredMove
	
	if len(candles) < 30 {
		return moves
	}
	
	// Find swing points
	type SwingPoint struct {
		Price float64
		Idx   int
		Type  string // "high" or "low"
	}
	
	var swings []SwingPoint
	
	for i := 2; i < len(candles)-2; i++ {
		if candles[i].High > candles[i-1].High &&
			candles[i].High > candles[i-2].High &&
			candles[i].High > candles[i+1].High &&
			candles[i].High > candles[i+2].High {
			swings = append(swings, SwingPoint{candles[i].High, i, "high"})
		}
		
		if candles[i].Low < candles[i-1].Low &&
			candles[i].Low < candles[i-2].Low &&
			candles[i].Low < candles[i+1].Low &&
			candles[i].Low < candles[i+2].Low {
			swings = append(swings, SwingPoint{candles[i].Low, i, "low"})
		}
	}
	
	// Look for ABCD pattern (measured move)
	// Bullish: Low(A) → High(B) → Low(C) → High(D)
	// AB = CD (measured move)
	
	for i := 0; i < len(swings)-3; i++ {
		a := swings[i]
		b := swings[i+1]
		c := swings[i+2]
		d := swings[i+3]
		
		// Bullish ABCD
		if a.Type == "low" && b.Type == "high" && c.Type == "low" && d.Type == "high" {
			ab := b.Price - a.Price
			bc := b.Price - c.Price
			cd := d.Price - c.Price
			
			// AB should roughly equal CD
			if math.Abs(ab-cd)/ab < 0.2 { // Within 20%
				move := MeasuredMove{
					Type:        "bullish",
					FirstLeg:    ab,
					Correction:  bc,
					SecondLeg:   cd,
					EntryPrice:  c.Price,
					TargetPrice: c.Price + ab, // Measured move target
					StopLoss:    c.Price - bc*0.5,
					Confidence:  80,
				}
				
				risk := move.EntryPrice - move.StopLoss
				reward := move.TargetPrice - move.EntryPrice
				move.RiskReward = reward / risk
				
				moves = append(moves, move)
			}
		}
		
		// Bearish ABCD
		if a.Type == "high" && b.Type == "low" && c.Type == "high" && d.Type == "low" {
			ab := a.Price - b.Price
			bc := c.Price - b.Price
			cd := c.Price - d.Price
			
			if math.Abs(ab-cd)/ab < 0.2 {
				move := MeasuredMove{
					Type:        "bearish",
					FirstLeg:    ab,
					Correction:  bc,
					SecondLeg:   cd,
					EntryPrice:  c.Price,
					TargetPrice: c.Price - ab,
					StopLoss:    c.Price + bc*0.5,
					Confidence:  80,
				}
				
				risk := move.StopLoss - move.EntryPrice
				reward := move.EntryPrice - move.TargetPrice
				move.RiskReward = reward / risk
				
				moves = append(moves, move)
			}
		}
	}
	
	return moves
}

// ==================== FULL MIRROR ANALYSIS ====================

// PerformMirrorAnalysis performs complete mirror market analysis
func PerformMirrorAnalysis(candles []Candle) *MirrorAnalysis {
	analysis := &MirrorAnalysis{
		MirrorActive:    false,
		MirrorDirection: "neutral",
	}
	
	// Find mirror patterns
	analysis.Patterns = FindMirrorPatterns(candles)
	
	// Find current active mirror
	if len(analysis.Patterns) > 0 {
		// Get most recent pattern
		for i := len(analysis.Patterns) - 1; i >= 0; i-- {
			pattern := analysis.Patterns[i]
			if pattern.CandleIdx >= len(candles)-20 {
				analysis.CurrentMirror = &pattern
				analysis.MirrorActive = true
				analysis.ExpectedTarget = pattern.TargetPrice
				
				if pattern.Type == "bullish_mirror" {
					analysis.MirrorDirection = "bullish"
				} else {
					analysis.MirrorDirection = "bearish"
				}
				break
			}
		}
	}
	
	// Calculate overall symmetry level
	if len(analysis.Patterns) > 0 {
		totalAccuracy := 0.0
		for _, p := range analysis.Patterns {
			totalAccuracy += p.Accuracy
		}
		analysis.SymmetryLevel = totalAccuracy / float64(len(analysis.Patterns))
	}
	
	return analysis
}

// GetMirrorSignal returns trading signal based on mirror analysis
func GetMirrorSignal(analysis *MirrorAnalysis, currentPrice float64) (string, float64) {
	if analysis == nil || !analysis.MirrorActive {
		return "neutral", 0
	}
	
	direction := analysis.MirrorDirection
	strength := 0.0
	
	if analysis.CurrentMirror != nil {
		strength = analysis.CurrentMirror.Accuracy * 0.8
		
		// Bonus if price is near symmetry point
		dist := math.Abs(currentPrice - analysis.CurrentMirror.SymmetryPoint)
		threshold := analysis.CurrentMirror.OriginalMove * 0.1
		
		if dist < threshold {
			strength += 15
		}
	}
	
	return direction, math.Min(strength, 100)
}

// ==================== FRACTAL MARKET ====================

// FractalPattern represents a fractal pattern
type FractalPattern struct {
	Type      string  // "bullish" or "bearish"
	High      float64
	Low       float64
	MidPoint  float64
	Strength  float64
	CandleIdx int
}

// FindFractals identifies Williams fractals
func FindFractals(candles []Candle) []FractalPattern {
	var fractals []FractalPattern
	
	if len(candles) < 5 {
		return fractals
	}
	
	for i := 2; i < len(candles)-2; i++ {
		// Bullish Fractal (swing low)
		if candles[i].Low < candles[i-1].Low &&
			candles[i].Low < candles[i-2].Low &&
			candles[i].Low < candles[i+1].Low &&
			candles[i].Low < candles[i+2].Low {
			
			fractal := FractalPattern{
				Type:      "bullish",
				High:      candles[i].High,
				Low:       candles[i].Low,
				MidPoint:  (candles[i].High + candles[i].Low) / 2,
				Strength:  70,
				CandleIdx: i,
			}
			fractals = append(fractals, fractal)
		}
		
		// Bearish Fractal (swing high)
		if candles[i].High > candles[i-1].High &&
			candles[i].High > candles[i-2].High &&
			candles[i].High > candles[i+1].High &&
			candles[i].High > candles[i+2].High {
			
			fractal := FractalPattern{
				Type:      "bearish",
				High:      candles[i].High,
				Low:       candles[i].Low,
				MidPoint:  (candles[i].High + candles[i].Low) / 2,
				Strength:  70,
				CandleIdx: i,
			}
			fractals = append(fractals, fractal)
		}
	}
	
	return fractals
}
