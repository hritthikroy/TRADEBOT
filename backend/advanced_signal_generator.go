package main

import (
	"math"
)

// ==================== ADVANCED SIGNAL GENERATOR ====================
// Combines: Multi-Timeframe, Candlestick Patterns, ICT/SMC, Order Flow,
// Delta Analysis, Footprint, Volume Profile, Enhanced Order Blocks

// AdvancedAnalysis holds all analysis results
type AdvancedAnalysis struct {
	MTF           *MultiTimeframeAnalysis
	ICT           *ICTAnalysis
	OrderFlow     *OrderFlowAnalysis
	VolumeProfile *VolumeProfile
	EnhancedOBs   []EnhancedOrderBlock
	BreakerBlocks []BreakerBlock
	MitigationBlocks []MitigationBlock
	Patterns      []CandlestickPattern
	
	// Scores
	TrendScore      float64
	MomentumScore   float64
	VolumeScore     float64
	StructureScore  float64
	ConfluenceScore float64
	TotalScore      float64
	
	// Direction
	Direction string
	Confidence float64
}

// PerformAdvancedAnalysis performs complete market analysis
func PerformAdvancedAnalysis(candles []Candle, interval string) *AdvancedAnalysis {
	if len(candles) < 50 {
		return nil
	}

	aa := &AdvancedAnalysis{}

	// 1. Multi-Timeframe Analysis
	aa.MTF = PerformMultiTimeframeAnalysis(candles, interval)

	// 2. ICT/SMC Analysis
	aa.ICT = PerformICTAnalysis(candles)

	// 3. Order Flow Analysis
	aa.OrderFlow = PerformOrderFlowAnalysis(candles)

	// 4. Volume Profile
	aa.VolumeProfile = CalculateVolumeProfile(candles, 20)

	// 5. Enhanced Order Blocks
	aa.EnhancedOBs = FindEnhancedOrderBlocks(candles)

	// 6. Breaker Blocks
	aa.BreakerBlocks = FindBreakerBlocks(candles, aa.EnhancedOBs)

	// 7. Mitigation Blocks
	aa.MitigationBlocks = FindMitigationBlocks(candles)

	// 8. Candlestick Patterns
	aa.Patterns = RecognizeAllPatterns(candles, interval)

	// Calculate scores
	aa.calculateScores(candles)

	return aa
}

// calculateScores calculates all component scores
func (aa *AdvancedAnalysis) calculateScores(candles []Candle) {
	currentPrice := candles[len(candles)-1].Close

	// 1. TREND SCORE (0-100)
	aa.TrendScore = 50
	if aa.MTF != nil {
		if aa.MTF.Direction == "bullish" {
			aa.TrendScore = 50 + float64(aa.MTF.Confluence)*15
		} else if aa.MTF.Direction == "bearish" {
			aa.TrendScore = 50 - float64(aa.MTF.Confluence)*15
		}
		aa.TrendScore += (aa.MTF.Strength - 50) * 0.3
	}

	// 2. MOMENTUM SCORE (0-100)
	aa.MomentumScore = 50
	if aa.OrderFlow != nil {
		// Delta contribution
		if aa.OrderFlow.Delta.Trend == "bullish" {
			aa.MomentumScore += aa.OrderFlow.Delta.Strength * 0.3
		} else if aa.OrderFlow.Delta.Trend == "bearish" {
			aa.MomentumScore -= aa.OrderFlow.Delta.Strength * 0.3
		}

		// Imbalance contribution
		aa.MomentumScore += aa.OrderFlow.Imbalance * 20

		// Exhaustion/Absorption signals
		if aa.OrderFlow.Exhaustion {
			// Exhaustion suggests reversal
			if aa.MomentumScore > 50 {
				aa.MomentumScore -= 15
			} else {
				aa.MomentumScore += 15
			}
		}
	}

	// 3. VOLUME SCORE (0-100)
	aa.VolumeScore = 50
	if aa.VolumeProfile != nil {
		// Price relative to POC
		if currentPrice > aa.VolumeProfile.POC {
			aa.VolumeScore += 10
		} else {
			aa.VolumeScore -= 10
		}

		// Price in Value Area
		if currentPrice >= aa.VolumeProfile.VAL && currentPrice <= aa.VolumeProfile.VAH {
			aa.VolumeScore += 5 // In value area = consolidation
		}

		// Near LVN (potential breakout)
		for _, lvn := range aa.VolumeProfile.LVN {
			if math.Abs(currentPrice-lvn) < (aa.VolumeProfile.VAH-aa.VolumeProfile.VAL)*0.1 {
				aa.VolumeScore += 15 // Near low volume = potential fast move
				break
			}
		}
	}

	if aa.OrderFlow != nil && aa.OrderFlow.Footprint.Imbalance > 1.2 {
		aa.VolumeScore += 15 // Buy imbalance
	} else if aa.OrderFlow != nil && aa.OrderFlow.Footprint.Imbalance < 0.8 {
		aa.VolumeScore -= 15 // Sell imbalance
	}

	// 4. STRUCTURE SCORE (0-100)
	aa.StructureScore = 50
	if aa.ICT != nil {
		// Market structure
		if aa.ICT.Structure.Trend == "bullish" {
			aa.StructureScore += 15
		} else if aa.ICT.Structure.Trend == "bearish" {
			aa.StructureScore -= 15
		}

		// BOS/CHOCH
		if aa.ICT.Structure.BOS {
			aa.StructureScore += 10
		}
		if aa.ICT.Structure.CHOCH {
			aa.StructureScore += 20 // Strong reversal signal
		}

		// Premium/Discount
		if aa.ICT.PremiumDiscount == "discount" {
			aa.StructureScore += 10 // Good for longs
		} else if aa.ICT.PremiumDiscount == "premium" {
			aa.StructureScore -= 10 // Good for shorts
		}

		// OTE
		if aa.ICT.OTE {
			aa.StructureScore += 15
		}
	}

	// Enhanced Order Blocks near price
	atr := calculateATR(candles[len(candles)-14:], 14)
	for _, ob := range aa.EnhancedOBs {
		if !ob.Mitigated {
			distance := math.Abs(currentPrice - ob.MidPoint)
			if distance < atr*2 {
				if ob.Type == "bullish" && currentPrice >= ob.Low && currentPrice <= ob.High+atr*0.5 {
					aa.StructureScore += ob.Strength * 0.2
				} else if ob.Type == "bearish" && currentPrice <= ob.High && currentPrice >= ob.Low-atr*0.5 {
					aa.StructureScore -= ob.Strength * 0.2
				}
			}
		}
	}

	// 5. CONFLUENCE SCORE (0-100)
	aa.ConfluenceScore = 0
	confluenceFactors := 0

	// Pattern confluence
	bullishPatterns := 0
	bearishPatterns := 0
	for _, p := range aa.Patterns {
		if p.Type == "bullish" {
			bullishPatterns++
		} else if p.Type == "bearish" {
			bearishPatterns++
		}
	}
	if bullishPatterns > bearishPatterns {
		aa.ConfluenceScore += float64(bullishPatterns) * 5
		confluenceFactors++
	} else if bearishPatterns > bullishPatterns {
		aa.ConfluenceScore -= float64(bearishPatterns) * 5
		confluenceFactors++
	}

	// MTF confluence
	if aa.MTF != nil && aa.MTF.Confluence >= 2 {
		if aa.MTF.Direction == "bullish" {
			aa.ConfluenceScore += 20
		} else {
			aa.ConfluenceScore -= 20
		}
		confluenceFactors++
	}

	// Order flow confluence
	if aa.OrderFlow != nil {
		if aa.OrderFlow.Delta.Trend == "bullish" && !aa.OrderFlow.Delta.Divergence {
			aa.ConfluenceScore += 15
			confluenceFactors++
		} else if aa.OrderFlow.Delta.Trend == "bearish" && !aa.OrderFlow.Delta.Divergence {
			aa.ConfluenceScore -= 15
			confluenceFactors++
		}
	}

	// ICT confluence
	if aa.ICT != nil && aa.ICT.Confluence >= 3 {
		if aa.ICT.Structure.Trend == "bullish" {
			aa.ConfluenceScore += 15
		} else {
			aa.ConfluenceScore -= 15
		}
		confluenceFactors++
	}

	// Normalize confluence score
	aa.ConfluenceScore = 50 + aa.ConfluenceScore

	// TOTAL SCORE (weighted average)
	aa.TotalScore = (aa.TrendScore*0.25 + aa.MomentumScore*0.25 + 
		aa.VolumeScore*0.15 + aa.StructureScore*0.20 + aa.ConfluenceScore*0.15)

	// Determine direction and confidence - require stronger bias
	if aa.TotalScore > 60 {
		aa.Direction = "bullish"
		aa.Confidence = math.Min((aa.TotalScore-50)*2.5, 95)
	} else if aa.TotalScore < 40 {
		aa.Direction = "bearish"
		aa.Confidence = math.Min((50-aa.TotalScore)*2.5, 95)
	} else {
		aa.Direction = "neutral"
		aa.Confidence = 0
	}

	// Clamp scores
	aa.TrendScore = math.Max(0, math.Min(100, aa.TrendScore))
	aa.MomentumScore = math.Max(0, math.Min(100, aa.MomentumScore))
	aa.VolumeScore = math.Max(0, math.Min(100, aa.VolumeScore))
	aa.StructureScore = math.Max(0, math.Min(100, aa.StructureScore))
	aa.ConfluenceScore = math.Max(0, math.Min(100, aa.ConfluenceScore))
	aa.TotalScore = math.Max(0, math.Min(100, aa.TotalScore))
}

// GenerateAdvancedSignal generates a trading signal using all analysis
func GenerateAdvancedSignal(candles []Candle, interval string) *Signal {
	if len(candles) < 50 {
		return nil
	}

	// Perform advanced analysis
	analysis := PerformAdvancedAnalysis(candles, interval)
	if analysis == nil {
		return nil
	}

	// Require high minimum confidence for quality signals
	if analysis.Confidence < 55 {
		return nil
	}

	// Don't trade against divergence
	if analysis.OrderFlow != nil && analysis.OrderFlow.Delta.Divergence {
		return nil
	}

	// Require multiple confluence factors
	confluenceCount := 0
	
	// MTF alignment
	if analysis.MTF != nil && analysis.MTF.Confluence >= 2 {
		confluenceCount++
	}
	
	// ICT confluence
	if analysis.ICT != nil && analysis.ICT.Confluence >= 2 {
		confluenceCount++
	}
	
	// Order flow alignment
	if analysis.OrderFlow != nil {
		if (analysis.Direction == "bullish" && analysis.OrderFlow.Delta.Trend == "bullish") ||
			(analysis.Direction == "bearish" && analysis.OrderFlow.Delta.Trend == "bearish") {
			confluenceCount++
		}
	}
	
	// Strong patterns present
	for _, p := range analysis.Patterns {
		if p.Type == analysis.Direction && p.Strength >= 80 {
			confluenceCount++
			break
		}
	}
	
	// Require at least 2 confluence factors
	if confluenceCount < 2 {
		return nil
	}

	currentPrice := candles[len(candles)-1].Close
	atr := calculateATR(candles[len(candles)-14:], 14)
	if atr == 0 {
		return nil
	}

	// Find optimal entry, stop loss, and targets
	var entry, stopLoss, target1, target2, target3 float64
	entry = currentPrice

	if analysis.Direction == "bullish" {
		// Find stop loss level
		stopLoss = findOptimalStopLoss(candles, analysis, "bullish", atr)
		
		// Calculate risk
		risk := entry - stopLoss
		if risk <= 0 {
			stopLoss = entry - atr*1.5
			risk = atr * 1.5
		}

		// Dynamic targets based on analysis strength
		targetMultiplier := 1.5 + (analysis.Confidence/100)*1.0
		target1 = entry + risk*targetMultiplier
		target2 = entry + risk*targetMultiplier*1.8
		target3 = entry + risk*targetMultiplier*2.5

		// Adjust targets based on resistance levels
		if analysis.VolumeProfile != nil && analysis.VolumeProfile.VAH > entry {
			// First target near VAH
			if analysis.VolumeProfile.VAH < target1 {
				target1 = analysis.VolumeProfile.VAH
			}
		}

	} else if analysis.Direction == "bearish" {
		// Find stop loss level
		stopLoss = findOptimalStopLoss(candles, analysis, "bearish", atr)
		
		// Calculate risk
		risk := stopLoss - entry
		if risk <= 0 {
			stopLoss = entry + atr*1.5
			risk = atr * 1.5
		}

		// Dynamic targets
		targetMultiplier := 1.5 + (analysis.Confidence/100)*1.0
		target1 = entry - risk*targetMultiplier
		target2 = entry - risk*targetMultiplier*1.8
		target3 = entry - risk*targetMultiplier*2.5

		// Adjust targets based on support levels
		if analysis.VolumeProfile != nil && analysis.VolumeProfile.VAL < entry {
			if analysis.VolumeProfile.VAL > target1 {
				target1 = analysis.VolumeProfile.VAL
			}
		}

	} else {
		return nil
	}

	// Verify minimum RR - require higher RR for profitability
	rr := math.Abs(target1-entry) / math.Abs(entry-stopLoss)
	if rr < 2.0 {
		return nil
	}

	signalType := "BUY"
	if analysis.Direction == "bearish" {
		signalType = "SELL"
	}

	return &Signal{
		Type:     signalType,
		Entry:    entry,
		StopLoss: stopLoss,
		Targets: []Target{
			{Price: target1, RR: rr, Percentage: 100},
			{Price: target2, RR: rr * 1.8, Percentage: 0},
			{Price: target3, RR: rr * 2.5, Percentage: 0},
		},
		Strength:  analysis.Confidence,
		Timeframe: interval,
	}
}

// findOptimalStopLoss finds the best stop loss level
func findOptimalStopLoss(candles []Candle, analysis *AdvancedAnalysis, direction string, atr float64) float64 {
	currentPrice := candles[len(candles)-1].Close
	
	if direction == "bullish" {
		// Start with swing low
		swingLow := findRecentSwingLow(candles, 15)
		stopLoss := swingLow - atr*0.3

		// Check for bullish order blocks below
		for _, ob := range analysis.EnhancedOBs {
			if ob.Type == "bullish" && !ob.Mitigated && ob.Low < currentPrice {
				// Place stop below order block
				obStop := ob.Low - atr*0.2
				if obStop > stopLoss && obStop < currentPrice-atr*0.5 {
					stopLoss = obStop
				}
			}
		}

		// Check volume profile VAL
		if analysis.VolumeProfile != nil && analysis.VolumeProfile.VAL < currentPrice {
			valStop := analysis.VolumeProfile.VAL - atr*0.2
			if valStop > stopLoss && valStop < currentPrice-atr*0.5 {
				stopLoss = valStop
			}
		}

		return stopLoss

	} else { // bearish
		// Start with swing high
		swingHigh := findRecentSwingHigh(candles, 15)
		stopLoss := swingHigh + atr*0.3

		// Check for bearish order blocks above
		for _, ob := range analysis.EnhancedOBs {
			if ob.Type == "bearish" && !ob.Mitigated && ob.High > currentPrice {
				obStop := ob.High + atr*0.2
				if obStop < stopLoss && obStop > currentPrice+atr*0.5 {
					stopLoss = obStop
				}
			}
		}

		// Check volume profile VAH
		if analysis.VolumeProfile != nil && analysis.VolumeProfile.VAH > currentPrice {
			vahStop := analysis.VolumeProfile.VAH + atr*0.2
			if vahStop < stopLoss && vahStop > currentPrice+atr*0.5 {
				stopLoss = vahStop
			}
		}

		return stopLoss
	}
}
