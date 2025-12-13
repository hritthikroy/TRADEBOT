package timeframe

import (
	"fmt"
	"math"
	"sort"
)

// ==================== COMPREHENSIVE MULTI-TIMEFRAME CONFLUENCE ====================
// Analyzes ALL timeframes: 1m, 3m, 5m, 15m, 30m, 45m, 1h, 2h, 4h, 6h, 8h, 12h, 1D

// AllTimeframes contains all supported timeframes
var AllTimeframes = []string{
	"1m", "3m", "5m", "15m", "30m", "45m",
	"1h", "2h", "4h", "6h", "8h", "12h", "1d",
}

// TimeframeWeight holds weight for each timeframe
var TimeframeWeights = map[string]float64{
	"1m":  0.5,  // Scalping - lowest weight
	"3m":  0.6,
	"5m":  0.7,
	"15m": 1.0,  // Day trading - base weight
	"30m": 1.2,
	"45m": 1.3,
	"1h":  1.5,  // Swing trading
	"2h":  1.7,
	"4h":  2.0,  // Position trading - high weight
	"6h":  2.2,
	"8h":  2.3,
	"12h": 2.5,
	"1d":  3.0,  // Daily - highest weight
}

// TimeframeMinutes converts timeframe to minutes
var TimeframeMinutes = map[string]int{
	"1m":  1,
	"3m":  3,
	"5m":  5,
	"15m": 15,
	"30m": 30,
	"45m": 45,
	"1h":  60,
	"2h":  120,
	"4h":  240,
	"6h":  360,
	"8h":  480,
	"12h": 720,
	"1d":  1440,
}


// ComprehensiveMTFAnalysis holds analysis for all timeframes
type ComprehensiveMTFAnalysis struct {
	Timeframes      map[string]*TimeframeAnalysis
	TotalTimeframes int
	BullishCount    int
	BearishCount    int
	NeutralCount    int
	
	// Weighted scores
	WeightedBullish float64
	WeightedBearish float64
	TotalWeight     float64
	
	// Confluence
	ConfluenceScore   float64 // 0-100
	ConfluencePercent float64 // Percentage of aligned TFs
	
	// Direction
	Direction       string  // "bullish", "bearish", "neutral"
	Strength        float64 // 0-100
	Confidence      float64 // 0-100
	
	// Grouped analysis
	ScalpingBias    string  // 1m, 3m, 5m
	DayTradingBias  string  // 15m, 30m, 45m
	SwingBias       string  // 1h, 2h, 4h
	PositionBias    string  // 6h, 8h, 12h, 1d
	
	// Key levels from all timeframes
	KeyResistance   []float64
	KeySupport      []float64
	
	// Alignment
	AllAligned      bool
	HigherTFAligned bool // 4h+ aligned
	LowerTFAligned  bool // Below entry TF aligned
}

// TimeframeAnalysis holds analysis for a single timeframe
type TimeframeAnalysis struct {
	Timeframe   string
	Direction   string  // "bullish", "bearish", "neutral"
	Strength    float64 // 0-100
	Weight      float64
	
	// Technical indicators
	Trend       string  // EMA trend
	RSI         float64
	RSISignal   string  // "oversold", "overbought", "neutral"
	MACD        string  // "bullish", "bearish", "neutral"
	
	// Price action
	AboveEMA20  bool
	AboveEMA50  bool
	AboveEMA200 bool
	
	// Structure
	Structure   string  // "uptrend", "downtrend", "ranging"
	SwingHigh   float64
	SwingLow    float64
	
	// Patterns
	BullishPatterns int
	BearishPatterns int
	
	// Volume
	VolumeStrength float64
	
	// Order blocks
	NearBullishOB bool
	NearBearishOB bool
}


// ==================== TIMEFRAME ANALYSIS ====================

// AnalyzeSingleTimeframe analyzes a single timeframe
func AnalyzeSingleTimeframe(candles []Candle, timeframe string) *TimeframeAnalysis {
	if len(candles) < 50 {
		return nil
	}
	
	ta := &TimeframeAnalysis{
		Timeframe: timeframe,
		Direction: "neutral",
		Weight:    TimeframeWeights[timeframe],
	}
	
	currentPrice := candles[len(candles)-1].Close
	
	// Calculate EMAs
	ema20 := calculateEMA(candles, 20)
	ema50 := calculateEMA(candles, 50)
	ema200 := 0.0
	if len(candles) >= 200 {
		ema200 = calculateEMA(candles, 200)
	}
	
	ta.AboveEMA20 = currentPrice > ema20
	ta.AboveEMA50 = currentPrice > ema50
	ta.AboveEMA200 = currentPrice > ema200
	
	// Calculate RSI
	ta.RSI = calculateRSI(candles, 14)
	if ta.RSI < 30 {
		ta.RSISignal = "oversold"
	} else if ta.RSI > 70 {
		ta.RSISignal = "overbought"
	} else {
		ta.RSISignal = "neutral"
	}
	
	// Determine trend
	if ema20 > ema50 && currentPrice > ema20 {
		ta.Trend = "bullish"
	} else if ema20 < ema50 && currentPrice < ema20 {
		ta.Trend = "bearish"
	} else {
		ta.Trend = "neutral"
	}
	
	// Calculate MACD
	ta.MACD = calculateMACDSignal(candles)
	
	// Find swing points
	ta.SwingHigh, ta.SwingLow = findRecentSwings(candles, 20)
	
	// Determine structure
	ta.Structure = determineStructure(candles)
	
	// Count patterns
	patterns := RecognizeAllPatterns(candles, timeframe)
	for _, p := range patterns {
		if p.Type == "bullish" {
			ta.BullishPatterns++
		} else if p.Type == "bearish" {
			ta.BearishPatterns++
		}
	}
	
	// Volume analysis
	ta.VolumeStrength = calculateVolumeStrength(candles)
	
	// Check for nearby order blocks
	obs := FindOrderBlocks(candles)
	atr := calculateATR(candles[len(candles)-14:], 14)
	for _, ob := range obs {
		dist := math.Abs(currentPrice - ob.MidPoint)
		if dist < atr*2 {
			if ob.Type == "bullish" {
				ta.NearBullishOB = true
			} else {
				ta.NearBearishOB = true
			}
		}
	}
	
	// Calculate overall direction and strength
	bullishScore := 0.0
	bearishScore := 0.0
	
	// Trend contribution
	if ta.Trend == "bullish" {
		bullishScore += 25
	} else if ta.Trend == "bearish" {
		bearishScore += 25
	}
	
	// EMA alignment
	if ta.AboveEMA20 && ta.AboveEMA50 {
		bullishScore += 15
	} else if !ta.AboveEMA20 && !ta.AboveEMA50 {
		bearishScore += 15
	}
	
	if ta.AboveEMA200 {
		bullishScore += 10
	} else if len(candles) >= 200 {
		bearishScore += 10
	}
	
	// RSI contribution
	if ta.RSISignal == "oversold" {
		bullishScore += 10 // Potential reversal up
	} else if ta.RSISignal == "overbought" {
		bearishScore += 10 // Potential reversal down
	}
	
	// MACD contribution
	if ta.MACD == "bullish" {
		bullishScore += 15
	} else if ta.MACD == "bearish" {
		bearishScore += 15
	}
	
	// Pattern contribution
	bullishScore += float64(ta.BullishPatterns) * 5
	bearishScore += float64(ta.BearishPatterns) * 5
	
	// Order block contribution
	if ta.NearBullishOB {
		bullishScore += 10
	}
	if ta.NearBearishOB {
		bearishScore += 10
	}
	
	// Volume contribution
	if ta.VolumeStrength > 1.2 {
		// High volume confirms direction
		if bullishScore > bearishScore {
			bullishScore += 10
		} else {
			bearishScore += 10
		}
	}
	
	// Determine direction
	if bullishScore > bearishScore+10 {
		ta.Direction = "bullish"
		ta.Strength = math.Min(bullishScore, 100)
	} else if bearishScore > bullishScore+10 {
		ta.Direction = "bearish"
		ta.Strength = math.Min(bearishScore, 100)
	} else {
		ta.Direction = "neutral"
		ta.Strength = 50
	}
	
	return ta
}


// ==================== HELPER FUNCTIONS ====================
// Note: calculateEMA and calculateRSI are defined in candlestick_patterns.go

// calculateMACDSignal calculates MACD signal
func calculateMACDSignal(candles []Candle) string {
	if len(candles) < 26 {
		return "neutral"
	}
	
	ema12 := calculateEMA(candles, 12)
	ema26 := calculateEMA(candles, 26)
	macd := ema12 - ema26
	
	// Simple signal line approximation
	prevMACD := 0.0
	if len(candles) > 27 {
		prevEMA12 := calculateEMA(candles[:len(candles)-1], 12)
		prevEMA26 := calculateEMA(candles[:len(candles)-1], 26)
		prevMACD = prevEMA12 - prevEMA26
	}
	
	if macd > 0 && macd > prevMACD {
		return "bullish"
	} else if macd < 0 && macd < prevMACD {
		return "bearish"
	}
	
	return "neutral"
}

// findRecentSwings finds recent swing high and low
func findRecentSwings(candles []Candle, lookback int) (float64, float64) {
	if len(candles) < lookback {
		lookback = len(candles)
	}
	
	recent := candles[len(candles)-lookback:]
	
	high := recent[0].High
	low := recent[0].Low
	
	for _, c := range recent {
		if c.High > high {
			high = c.High
		}
		if c.Low < low {
			low = c.Low
		}
	}
	
	return high, low
}

// determineStructure determines market structure
func determineStructure(candles []Candle) string {
	if len(candles) < 20 {
		return "ranging"
	}
	
	// Find swing highs and lows
	var highs, lows []float64
	
	for i := 2; i < len(candles)-2; i++ {
		if candles[i].High > candles[i-1].High &&
			candles[i].High > candles[i-2].High &&
			candles[i].High > candles[i+1].High &&
			candles[i].High > candles[i+2].High {
			highs = append(highs, candles[i].High)
		}
		
		if candles[i].Low < candles[i-1].Low &&
			candles[i].Low < candles[i-2].Low &&
			candles[i].Low < candles[i+1].Low &&
			candles[i].Low < candles[i+2].Low {
			lows = append(lows, candles[i].Low)
		}
	}
	
	if len(highs) < 2 || len(lows) < 2 {
		return "ranging"
	}
	
	// Check for higher highs and higher lows (uptrend)
	lastHigh := highs[len(highs)-1]
	prevHigh := highs[len(highs)-2]
	lastLow := lows[len(lows)-1]
	prevLow := lows[len(lows)-2]
	
	if lastHigh > prevHigh && lastLow > prevLow {
		return "uptrend"
	}
	
	if lastHigh < prevHigh && lastLow < prevLow {
		return "downtrend"
	}
	
	return "ranging"
}

// calculateVolumeStrength calculates volume relative to average
func calculateVolumeStrength(candles []Candle) float64 {
	if len(candles) < 20 {
		return 1.0
	}
	
	avgVolume := 0.0
	for i := len(candles) - 20; i < len(candles)-1; i++ {
		avgVolume += candles[i].Volume
	}
	avgVolume /= 19
	
	if avgVolume == 0 {
		return 1.0
	}
	
	currentVolume := candles[len(candles)-1].Volume
	return currentVolume / avgVolume
}


// ==================== COMPREHENSIVE MTF ANALYSIS ====================

// PerformComprehensiveMTFAnalysis analyzes all timeframes
func PerformComprehensiveMTFAnalysis(candlesByTF map[string][]Candle) *ComprehensiveMTFAnalysis {
	analysis := &ComprehensiveMTFAnalysis{
		Timeframes:    make(map[string]*TimeframeAnalysis),
		KeyResistance: []float64{},
		KeySupport:    []float64{},
	}
	
	// Analyze each timeframe
	for tf, candles := range candlesByTF {
		if len(candles) >= 50 {
			ta := AnalyzeSingleTimeframe(candles, tf)
			if ta != nil {
				analysis.Timeframes[tf] = ta
				analysis.TotalTimeframes++
				
				// Count directions
				switch ta.Direction {
				case "bullish":
					analysis.BullishCount++
					analysis.WeightedBullish += ta.Strength * ta.Weight
				case "bearish":
					analysis.BearishCount++
					analysis.WeightedBearish += ta.Strength * ta.Weight
				default:
					analysis.NeutralCount++
				}
				
				analysis.TotalWeight += ta.Weight
				
				// Collect key levels
				analysis.KeyResistance = append(analysis.KeyResistance, ta.SwingHigh)
				analysis.KeySupport = append(analysis.KeySupport, ta.SwingLow)
			}
		}
	}
	
	if analysis.TotalTimeframes == 0 {
		return analysis
	}
	
	// Calculate weighted scores
	if analysis.TotalWeight > 0 {
		analysis.WeightedBullish /= analysis.TotalWeight
		analysis.WeightedBearish /= analysis.TotalWeight
	}
	
	// Calculate confluence
	maxAligned := math.Max(float64(analysis.BullishCount), float64(analysis.BearishCount))
	analysis.ConfluencePercent = (maxAligned / float64(analysis.TotalTimeframes)) * 100
	analysis.ConfluenceScore = analysis.ConfluencePercent
	
	// Determine overall direction
	if analysis.WeightedBullish > analysis.WeightedBearish+10 {
		analysis.Direction = "bullish"
		analysis.Strength = analysis.WeightedBullish
		analysis.Confidence = analysis.ConfluencePercent
	} else if analysis.WeightedBearish > analysis.WeightedBullish+10 {
		analysis.Direction = "bearish"
		analysis.Strength = analysis.WeightedBearish
		analysis.Confidence = analysis.ConfluencePercent
	} else {
		analysis.Direction = "neutral"
		analysis.Strength = 50
		analysis.Confidence = 50
	}
	
	// Calculate group biases
	analysis.ScalpingBias = calculateGroupBias(analysis.Timeframes, []string{"1m", "3m", "5m"})
	analysis.DayTradingBias = calculateGroupBias(analysis.Timeframes, []string{"15m", "30m", "45m"})
	analysis.SwingBias = calculateGroupBias(analysis.Timeframes, []string{"1h", "2h", "4h"})
	analysis.PositionBias = calculateGroupBias(analysis.Timeframes, []string{"6h", "8h", "12h", "1d"})
	
	// Check alignment
	analysis.AllAligned = analysis.BullishCount == analysis.TotalTimeframes || 
		analysis.BearishCount == analysis.TotalTimeframes
	
	// Check higher TF alignment (4h+)
	higherTFs := []string{"4h", "6h", "8h", "12h", "1d"}
	higherBullish := 0
	higherBearish := 0
	higherTotal := 0
	
	for _, tf := range higherTFs {
		if ta, ok := analysis.Timeframes[tf]; ok {
			higherTotal++
			if ta.Direction == "bullish" {
				higherBullish++
			} else if ta.Direction == "bearish" {
				higherBearish++
			}
		}
	}
	
	if higherTotal > 0 {
		analysis.HigherTFAligned = higherBullish == higherTotal || higherBearish == higherTotal
	}
	
	// Sort and deduplicate key levels
	analysis.KeyResistance = deduplicateLevels(analysis.KeyResistance)
	analysis.KeySupport = deduplicateLevels(analysis.KeySupport)
	
	return analysis
}

// calculateGroupBias calculates bias for a group of timeframes
func calculateGroupBias(timeframes map[string]*TimeframeAnalysis, group []string) string {
	bullish := 0
	bearish := 0
	
	for _, tf := range group {
		if ta, ok := timeframes[tf]; ok {
			if ta.Direction == "bullish" {
				bullish++
			} else if ta.Direction == "bearish" {
				bearish++
			}
		}
	}
	
	if bullish > bearish {
		return "bullish"
	} else if bearish > bullish {
		return "bearish"
	}
	return "neutral"
}

// deduplicateLevels removes duplicate price levels
func deduplicateLevels(levels []float64) []float64 {
	if len(levels) == 0 {
		return levels
	}
	
	sort.Float64s(levels)
	
	result := []float64{levels[0]}
	threshold := levels[0] * 0.001 // 0.1% threshold
	
	for i := 1; i < len(levels); i++ {
		if levels[i]-result[len(result)-1] > threshold {
			result = append(result, levels[i])
		}
	}
	
	return result
}


// ==================== MTF SIGNAL GENERATION ====================

// GetMTFSignal returns trading signal based on MTF analysis
func GetMTFSignal(analysis *ComprehensiveMTFAnalysis, entryTF string) (string, float64, string) {
	if analysis == nil || analysis.TotalTimeframes == 0 {
		return "neutral", 0, "No MTF data"
	}
	
	direction := analysis.Direction
	strength := analysis.Strength
	reason := ""
	
	// Build reason string
	reason = fmt.Sprintf("%d/%d TFs aligned (%d bull, %d bear, %d neutral)",
		int(math.Max(float64(analysis.BullishCount), float64(analysis.BearishCount))),
		analysis.TotalTimeframes,
		analysis.BullishCount,
		analysis.BearishCount,
		analysis.NeutralCount)
	
	// Bonus for all aligned
	if analysis.AllAligned {
		strength += 20
		reason += " | ALL ALIGNED"
	}
	
	// Bonus for higher TF alignment
	if analysis.HigherTFAligned {
		strength += 15
		reason += " | Higher TFs aligned"
	}
	
	// Check entry TF alignment
	if ta, ok := analysis.Timeframes[entryTF]; ok {
		if ta.Direction == direction {
			strength += 10
			reason += " | Entry TF confirms"
		} else if ta.Direction != "neutral" {
			strength -= 15
			reason += " | Entry TF conflicts"
		}
	}
	
	// Group alignment bonuses
	groupsAligned := 0
	if analysis.ScalpingBias == direction {
		groupsAligned++
	}
	if analysis.DayTradingBias == direction {
		groupsAligned++
	}
	if analysis.SwingBias == direction {
		groupsAligned++
	}
	if analysis.PositionBias == direction {
		groupsAligned++
	}
	
	if groupsAligned >= 3 {
		strength += 15
		reason += fmt.Sprintf(" | %d/4 groups aligned", groupsAligned)
	}
	
	return direction, math.Min(strength, 100), reason
}

// GetMTFConfluenceScore returns confluence score for signal validation
func GetMTFConfluenceScore(analysis *ComprehensiveMTFAnalysis, direction string) float64 {
	if analysis == nil {
		return 0
	}
	
	score := 0.0
	
	// Base score from confluence percent
	score += analysis.ConfluencePercent * 0.5
	
	// Weighted direction alignment
	if direction == "bullish" {
		score += analysis.WeightedBullish * 0.3
	} else if direction == "bearish" {
		score += analysis.WeightedBearish * 0.3
	}
	
	// All aligned bonus
	if analysis.AllAligned && analysis.Direction == direction {
		score += 20
	}
	
	// Higher TF alignment bonus
	if analysis.HigherTFAligned {
		higherDir := "neutral"
		for _, tf := range []string{"4h", "6h", "8h", "12h", "1d"} {
			if ta, ok := analysis.Timeframes[tf]; ok {
				if ta.Direction != "neutral" {
					higherDir = ta.Direction
					break
				}
			}
		}
		if higherDir == direction {
			score += 15
		}
	}
	
	return math.Min(score, 100)
}

// ShouldTradeMTF determines if MTF confluence is sufficient for trading
func ShouldTradeMTF(analysis *ComprehensiveMTFAnalysis, direction string, minConfluence float64) bool {
	if analysis == nil {
		return false
	}
	
	// Check minimum confluence
	if analysis.ConfluencePercent < minConfluence {
		return false
	}
	
	// Direction must match
	if analysis.Direction != direction {
		return false
	}
	
	// Higher TFs should not conflict
	for _, tf := range []string{"4h", "6h", "8h", "12h", "1d"} {
		if ta, ok := analysis.Timeframes[tf]; ok {
			if ta.Direction != "neutral" && ta.Direction != direction {
				return false // Higher TF conflicts
			}
		}
	}
	
	return true
}


// ==================== MTF DATA FETCHING ====================

// FetchAllTimeframeData fetches data for all timeframes
func FetchAllTimeframeData(symbol string, limit int) (map[string][]Candle, error) {
	sg := &SignalGenerator{}
	result := make(map[string][]Candle)
	
	for _, tf := range AllTimeframes {
		candles, err := sg.FetchMarketData(symbol, tf, limit)
		if err != nil {
			continue // Skip failed timeframes
		}
		if len(candles) >= 50 {
			result[tf] = candles
		}
	}
	
	return result, nil
}

// FetchSelectedTimeframeData fetches data for selected timeframes
func FetchSelectedTimeframeData(symbol string, timeframes []string, limit int) (map[string][]Candle, error) {
	sg := &SignalGenerator{}
	result := make(map[string][]Candle)
	
	for _, tf := range timeframes {
		candles, err := sg.FetchMarketData(symbol, tf, limit)
		if err != nil {
			continue
		}
		if len(candles) >= 50 {
			result[tf] = candles
		}
	}
	
	return result, nil
}

// ==================== MTF REPORT ====================

// MTFReport holds a formatted report of MTF analysis
type MTFReport struct {
	Summary       string
	Direction     string
	Confidence    float64
	Confluence    float64
	Details       []MTFTimeframeDetail
	Recommendation string
}

// MTFTimeframeDetail holds detail for each timeframe
type MTFTimeframeDetail struct {
	Timeframe   string
	Direction   string
	Strength    float64
	Trend       string
	RSI         float64
	MACD        string
	Structure   string
}

// GenerateMTFReport generates a detailed MTF report
func GenerateMTFReport(analysis *ComprehensiveMTFAnalysis) *MTFReport {
	if analysis == nil {
		return nil
	}
	
	report := &MTFReport{
		Direction:  analysis.Direction,
		Confidence: analysis.Confidence,
		Confluence: analysis.ConfluencePercent,
		Details:    []MTFTimeframeDetail{},
	}
	
	// Generate summary
	report.Summary = fmt.Sprintf(
		"MTF Analysis: %s (%.1f%% confidence, %.1f%% confluence)\n"+
		"Bullish: %d | Bearish: %d | Neutral: %d\n"+
		"Scalping: %s | Day Trading: %s | Swing: %s | Position: %s",
		analysis.Direction,
		analysis.Confidence,
		analysis.ConfluencePercent,
		analysis.BullishCount,
		analysis.BearishCount,
		analysis.NeutralCount,
		analysis.ScalpingBias,
		analysis.DayTradingBias,
		analysis.SwingBias,
		analysis.PositionBias,
	)
	
	// Add details for each timeframe
	for _, tf := range AllTimeframes {
		if ta, ok := analysis.Timeframes[tf]; ok {
			detail := MTFTimeframeDetail{
				Timeframe: tf,
				Direction: ta.Direction,
				Strength:  ta.Strength,
				Trend:     ta.Trend,
				RSI:       ta.RSI,
				MACD:      ta.MACD,
				Structure: ta.Structure,
			}
			report.Details = append(report.Details, detail)
		}
	}
	
	// Generate recommendation
	if analysis.AllAligned {
		report.Recommendation = fmt.Sprintf("STRONG %s - All timeframes aligned", analysis.Direction)
	} else if analysis.HigherTFAligned && analysis.ConfluencePercent >= 70 {
		report.Recommendation = fmt.Sprintf("GOOD %s - Higher TFs aligned, %.0f%% confluence", 
			analysis.Direction, analysis.ConfluencePercent)
	} else if analysis.ConfluencePercent >= 60 {
		report.Recommendation = fmt.Sprintf("MODERATE %s - %.0f%% confluence", 
			analysis.Direction, analysis.ConfluencePercent)
	} else {
		report.Recommendation = "WAIT - Insufficient confluence"
	}
	
	return report
}

// PrintMTFReport prints a formatted MTF report
func PrintMTFReport(report *MTFReport) {
	if report == nil {
		fmt.Println("No MTF report available")
		return
	}
	
	fmt.Println("\n" + "═══════════════════════════════════════════════════════════")
	fmt.Println("                 MULTI-TIMEFRAME CONFLUENCE REPORT")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println(report.Summary)
	fmt.Println("───────────────────────────────────────────────────────────")
	fmt.Println("TIMEFRAME DETAILS:")
	fmt.Println("───────────────────────────────────────────────────────────")
	
	for _, d := range report.Details {
		icon := "⚪"
		if d.Direction == "bullish" {
			icon = "🟢"
		} else if d.Direction == "bearish" {
			icon = "🔴"
		}
		
		fmt.Printf("%s %-4s | %-8s | Str: %5.1f | RSI: %5.1f | MACD: %-8s | %s\n",
			icon, d.Timeframe, d.Direction, d.Strength, d.RSI, d.MACD, d.Structure)
	}
	
	fmt.Println("───────────────────────────────────────────────────────────")
	fmt.Printf("RECOMMENDATION: %s\n", report.Recommendation)
	fmt.Println("═══════════════════════════════════════════════════════════")
}
