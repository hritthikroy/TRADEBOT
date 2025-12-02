package main

import (
	"math"
	"time"
)

// ==================== POWER OF 3 (AMD) ====================
// ICT's Power of 3: Accumulation → Manipulation → Distribution
// The core institutional trading model

// PO3Phase represents a Power of 3 phase
type PO3Phase struct {
	Phase       string  // "accumulation", "manipulation", "distribution"
	StartIdx    int
	EndIdx      int
	StartPrice  float64
	EndPrice    float64
	Range       float64
	Volume      float64
	Confidence  float64
	Direction   string  // "bullish" or "bearish"
}

// PO3Analysis holds Power of 3 analysis
type PO3Analysis struct {
	CurrentPhase    *PO3Phase
	Phases          []PO3Phase
	DailyBias       string  // "bullish", "bearish", "neutral"
	WeeklyBias      string
	ExpectedMove    string  // "up", "down", "sideways"
	OptimalEntry    bool
	ManipulationComplete bool
	DistributionStarted  bool
}

// ==================== DAILY POWER OF 3 ====================

// AnalyzeDailyPO3 analyzes daily Power of 3 structure
func AnalyzeDailyPO3(candles []Candle) *PO3Analysis {
	analysis := &PO3Analysis{
		Phases:     []PO3Phase{},
		DailyBias:  "neutral",
		WeeklyBias: "neutral",
	}
	
	if len(candles) < 24 {
		return analysis
	}
	
	// Get current time in UTC
	now := time.Now().UTC()
	hour := now.Hour()
	
	// Group candles by session
	asianCandles := []Candle{}
	londonCandles := []Candle{}
	nyCandles := []Candle{}
	
	for _, c := range candles {
		t := time.Unix(c.Timestamp/1000, 0).UTC()
		h := t.Hour()
		
		// Asian: 00:00-08:00 UTC (Accumulation)
		if h >= 0 && h < 8 {
			asianCandles = append(asianCandles, c)
		}
		// London: 08:00-13:00 UTC (Manipulation)
		if h >= 8 && h < 13 {
			londonCandles = append(londonCandles, c)
		}
		// NY: 13:00-21:00 UTC (Distribution)
		if h >= 13 && h < 21 {
			nyCandles = append(nyCandles, c)
		}
	}
	
	// Analyze Asian Session (Accumulation)
	if len(asianCandles) > 0 {
		asianHigh := asianCandles[0].High
		asianLow := asianCandles[0].Low
		asianVolume := 0.0
		
		for _, c := range asianCandles {
			if c.High > asianHigh {
				asianHigh = c.High
			}
			if c.Low < asianLow {
				asianLow = c.Low
			}
			asianVolume += c.Volume
		}
		
		accumulation := PO3Phase{
			Phase:      "accumulation",
			StartPrice: asianCandles[0].Open,
			EndPrice:   asianCandles[len(asianCandles)-1].Close,
			Range:      asianHigh - asianLow,
			Volume:     asianVolume,
			Direction:  "neutral",
			Confidence: 70,
		}
		
		// Determine accumulation bias
		if accumulation.EndPrice > accumulation.StartPrice {
			accumulation.Direction = "bullish"
		} else if accumulation.EndPrice < accumulation.StartPrice {
			accumulation.Direction = "bearish"
		}
		
		analysis.Phases = append(analysis.Phases, accumulation)
	}
	
	// Analyze London Session (Manipulation)
	if len(londonCandles) > 0 {
		londonHigh := londonCandles[0].High
		londonLow := londonCandles[0].Low
		londonVolume := 0.0
		
		for _, c := range londonCandles {
			if c.High > londonHigh {
				londonHigh = c.High
			}
			if c.Low < londonLow {
				londonLow = c.Low
			}
			londonVolume += c.Volume
		}
		
		manipulation := PO3Phase{
			Phase:      "manipulation",
			StartPrice: londonCandles[0].Open,
			EndPrice:   londonCandles[len(londonCandles)-1].Close,
			Range:      londonHigh - londonLow,
			Volume:     londonVolume,
			Direction:  "neutral",
			Confidence: 75,
		}
		
		// Manipulation often goes AGAINST the true direction first
		// Then reverses for distribution
		if manipulation.EndPrice > manipulation.StartPrice {
			manipulation.Direction = "bullish"
		} else {
			manipulation.Direction = "bearish"
		}
		
		analysis.Phases = append(analysis.Phases, manipulation)
		
		// Check if manipulation is complete
		if len(asianCandles) > 0 {
			asianHigh := asianCandles[0].High
			asianLow := asianCandles[0].Low
			for _, c := range asianCandles {
				if c.High > asianHigh {
					asianHigh = c.High
				}
				if c.Low < asianLow {
					asianLow = c.Low
				}
			}
			
			// Manipulation sweeps Asian high or low
			if londonHigh > asianHigh || londonLow < asianLow {
				analysis.ManipulationComplete = true
			}
		}
	}
	
	// Analyze NY Session (Distribution)
	if len(nyCandles) > 0 {
		nyHigh := nyCandles[0].High
		nyLow := nyCandles[0].Low
		nyVolume := 0.0
		
		for _, c := range nyCandles {
			if c.High > nyHigh {
				nyHigh = c.High
			}
			if c.Low < nyLow {
				nyLow = c.Low
			}
			nyVolume += c.Volume
		}
		
		distribution := PO3Phase{
			Phase:      "distribution",
			StartPrice: nyCandles[0].Open,
			EndPrice:   nyCandles[len(nyCandles)-1].Close,
			Range:      nyHigh - nyLow,
			Volume:     nyVolume,
			Direction:  "neutral",
			Confidence: 80,
		}
		
		if distribution.EndPrice > distribution.StartPrice {
			distribution.Direction = "bullish"
		} else {
			distribution.Direction = "bearish"
		}
		
		analysis.Phases = append(analysis.Phases, distribution)
		analysis.DistributionStarted = true
	}
	
	// Determine current phase based on time
	if hour >= 0 && hour < 8 {
		if len(analysis.Phases) > 0 {
			analysis.CurrentPhase = &analysis.Phases[0]
		}
	} else if hour >= 8 && hour < 13 {
		if len(analysis.Phases) > 1 {
			analysis.CurrentPhase = &analysis.Phases[1]
		}
	} else if hour >= 13 && hour < 21 {
		if len(analysis.Phases) > 2 {
			analysis.CurrentPhase = &analysis.Phases[2]
		}
	}
	
	// Determine daily bias
	analysis.DailyBias = determinePO3Bias(analysis)
	
	// Determine expected move
	if analysis.ManipulationComplete && !analysis.DistributionStarted {
		// After manipulation, expect opposite move
		if len(analysis.Phases) > 1 && analysis.Phases[1].Direction == "bullish" {
			analysis.ExpectedMove = "down"
		} else if len(analysis.Phases) > 1 && analysis.Phases[1].Direction == "bearish" {
			analysis.ExpectedMove = "up"
		}
	}
	
	// Check for optimal entry
	analysis.OptimalEntry = isOptimalPO3Entry(analysis, hour)
	
	return analysis
}

// determinePO3Bias determines the daily bias from PO3 analysis
func determinePO3Bias(analysis *PO3Analysis) string {
	if len(analysis.Phases) < 2 {
		return "neutral"
	}
	
	// If accumulation and manipulation agree, strong bias
	if analysis.Phases[0].Direction == analysis.Phases[1].Direction {
		return analysis.Phases[0].Direction
	}
	
	// If manipulation reversed accumulation, bias is manipulation direction
	if analysis.ManipulationComplete {
		// After manipulation sweep, expect reversal
		if analysis.Phases[1].Direction == "bullish" {
			return "bearish" // Manipulation up = expect down
		}
		return "bullish" // Manipulation down = expect up
	}
	
	return "neutral"
}

// isOptimalPO3Entry checks if current time is optimal for PO3 entry
func isOptimalPO3Entry(analysis *PO3Analysis, hour int) bool {
	// Best entries:
	// 1. After manipulation (08:00-10:00 UTC) - London open
	// 2. Start of distribution (13:00-15:00 UTC) - NY open
	
	if hour >= 8 && hour < 10 && analysis.ManipulationComplete {
		return true
	}
	
	if hour >= 13 && hour < 15 && analysis.ManipulationComplete {
		return true
	}
	
	return false
}

// GetPO3Signal returns trading signal based on Power of 3
func GetPO3Signal(analysis *PO3Analysis, currentPrice float64) (string, float64) {
	if analysis == nil {
		return "neutral", 0
	}
	
	direction := "neutral"
	strength := 0.0
	
	// Strong signal when manipulation is complete
	if analysis.ManipulationComplete && analysis.OptimalEntry {
		if analysis.ExpectedMove == "up" {
			direction = "bullish"
			strength = 80
		} else if analysis.ExpectedMove == "down" {
			direction = "bearish"
			strength = 80
		}
	}
	
	// Moderate signal based on daily bias
	if direction == "neutral" && analysis.DailyBias != "neutral" {
		direction = analysis.DailyBias
		strength = 60
	}
	
	// Bonus for distribution phase alignment
	if analysis.DistributionStarted && len(analysis.Phases) > 2 {
		if analysis.Phases[2].Direction == direction {
			strength += 15
		}
	}
	
	return direction, math.Min(strength, 100)
}

// ==================== WEEKLY POWER OF 3 ====================

// AnalyzeWeeklyPO3 analyzes weekly Power of 3 structure
func AnalyzeWeeklyPO3(candles []Candle) *PO3Analysis {
	analysis := &PO3Analysis{
		Phases:     []PO3Phase{},
		WeeklyBias: "neutral",
	}
	
	if len(candles) < 100 {
		return analysis
	}
	
	// Weekly PO3:
	// Monday-Tuesday: Accumulation
	// Wednesday: Manipulation
	// Thursday-Friday: Distribution
	
	mondayTuesday := []Candle{}
	wednesday := []Candle{}
	thursdayFriday := []Candle{}
	
	for _, c := range candles {
		t := time.Unix(c.Timestamp/1000, 0).UTC()
		weekday := t.Weekday()
		
		switch weekday {
		case time.Monday, time.Tuesday:
			mondayTuesday = append(mondayTuesday, c)
		case time.Wednesday:
			wednesday = append(wednesday, c)
		case time.Thursday, time.Friday:
			thursdayFriday = append(thursdayFriday, c)
		}
	}
	
	// Analyze each phase
	if len(mondayTuesday) > 0 {
		phase := analyzePO3PhaseCandles(mondayTuesday, "accumulation")
		analysis.Phases = append(analysis.Phases, phase)
	}
	
	if len(wednesday) > 0 {
		phase := analyzePO3PhaseCandles(wednesday, "manipulation")
		analysis.Phases = append(analysis.Phases, phase)
	}
	
	if len(thursdayFriday) > 0 {
		phase := analyzePO3PhaseCandles(thursdayFriday, "distribution")
		analysis.Phases = append(analysis.Phases, phase)
	}
	
	// Determine weekly bias
	analysis.WeeklyBias = determinePO3Bias(analysis)
	
	return analysis
}

// analyzePO3PhaseCandles analyzes candles for a PO3 phase
func analyzePO3PhaseCandles(candles []Candle, phaseName string) PO3Phase {
	phase := PO3Phase{
		Phase:      phaseName,
		Confidence: 70,
	}
	
	if len(candles) == 0 {
		return phase
	}
	
	high := candles[0].High
	low := candles[0].Low
	volume := 0.0
	
	for _, c := range candles {
		if c.High > high {
			high = c.High
		}
		if c.Low < low {
			low = c.Low
		}
		volume += c.Volume
	}
	
	phase.StartPrice = candles[0].Open
	phase.EndPrice = candles[len(candles)-1].Close
	phase.Range = high - low
	phase.Volume = volume
	
	if phase.EndPrice > phase.StartPrice {
		phase.Direction = "bullish"
	} else if phase.EndPrice < phase.StartPrice {
		phase.Direction = "bearish"
	} else {
		phase.Direction = "neutral"
	}
	
	return phase
}

// ==================== CANDLE POWER OF 3 ====================

// CandlePO3 represents Power of 3 within a single candle
type CandlePO3 struct {
	Accumulation float64 // Opening range
	Manipulation float64 // Wick (fake move)
	Distribution float64 // Body (real move)
	Direction    string
	Strength     float64
	IsValid      bool
}

// AnalyzeCandlePO3 analyzes Power of 3 within a single candle
func AnalyzeCandlePO3(candle Candle) *CandlePO3 {
	po3 := &CandlePO3{
		IsValid: false,
	}
	
	body := math.Abs(candle.Close - candle.Open)
	upperWick := candle.High - math.Max(candle.Open, candle.Close)
	lowerWick := math.Min(candle.Open, candle.Close) - candle.Low
	totalRange := candle.High - candle.Low
	
	if totalRange == 0 {
		return po3
	}
	
	// Bullish PO3 Candle:
	// - Lower wick (manipulation down)
	// - Body up (distribution up)
	if candle.Close > candle.Open && lowerWick > body*0.3 {
		po3.Accumulation = (candle.Open - candle.Low) / totalRange * 100
		po3.Manipulation = lowerWick / totalRange * 100
		po3.Distribution = body / totalRange * 100
		po3.Direction = "bullish"
		po3.Strength = (body / totalRange) * 100
		po3.IsValid = true
	}
	
	// Bearish PO3 Candle:
	// - Upper wick (manipulation up)
	// - Body down (distribution down)
	if candle.Close < candle.Open && upperWick > body*0.3 {
		po3.Accumulation = (candle.High - candle.Open) / totalRange * 100
		po3.Manipulation = upperWick / totalRange * 100
		po3.Distribution = body / totalRange * 100
		po3.Direction = "bearish"
		po3.Strength = (body / totalRange) * 100
		po3.IsValid = true
	}
	
	return po3
}

// FindPO3Candles finds candles with Power of 3 structure
func FindPO3Candles(candles []Candle) []CandlePO3 {
	var po3Candles []CandlePO3
	
	for _, c := range candles {
		po3 := AnalyzeCandlePO3(c)
		if po3.IsValid && po3.Strength >= 50 {
			po3Candles = append(po3Candles, *po3)
		}
	}
	
	return po3Candles
}
