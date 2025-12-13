package main

import (
	"log"
	"math"
)

// MultiTimeframeSignal represents a signal with multi-timeframe confluence
type MultiTimeframeSignal struct {
	// Entry details
	EntryTimeframe string
	EntryPrice     float64
	EntryTime      int64
	
	// Direction from higher timeframe
	HTFDirection   string  // From 4h
	HTFTrend       string  // bullish/bearish
	HTFConfidence  float64
	
	// Key levels from 1h
	OrderBlocks    []MTFOrderBlock
	FVGs           []MTFFairValueGap
	LiquidityZones []LiquidityZone
	
	// Volume analysis from 15m
	VolumeProfile  MTFVolumeProfile
	DeltaAnalysis  MTFDeltaAnalysis
	InsideFVG      bool
	
	// Precise entry from 3m/1m
	PreciseEntry   float64
	OptimalEntry   bool
	EntryQuality   float64 // 0-100
	
	// Trade setup
	Type           string
	StopLoss       float64
	TP1            float64
	TP2            float64
	TP3            float64
	RiskReward     float64
	Confluence     int
	Strength       float64
}

// MTFOrderBlock represents an order block zone for multi-TF
type MTFOrderBlock struct {
	High       float64
	Low        float64
	Type       string // bullish/bearish
	Strength   float64
	Timeframe  string
	Tested     bool
}

// MTFFairValueGap represents a fair value gap for multi-TF
type MTFFairValueGap struct {
	High       float64
	Low        float64
	Type       string // bullish/bearish
	Filled     bool
	Timeframe  string
}

// LiquidityZone represents a liquidity zone
type LiquidityZone struct {
	Price      float64
	Type       string // high/low
	Swept      bool
	Strength   float64
}

// MTFVolumeProfile represents volume analysis for multi-TF
type MTFVolumeProfile struct {
	POC        float64 // Point of Control
	VAH        float64 // Value Area High
	VAL        float64 // Value Area Low
	HighVolume bool
	BuyVolume  float64
	SellVolume float64
}

// MTFDeltaAnalysis represents order flow delta for multi-TF
type MTFDeltaAnalysis struct {
	CumulativeDelta float64
	DeltaDirection  string // positive/negative
	Strength        float64
	Divergence      bool
}

// GenerateMultiTimeframeSignal generates signal using top-down analysis
// OPTIMIZED VERSION with balanced parameters
func GenerateMultiTimeframeSignal(
	candles4h []Candle,
	candles1h []Candle,
	candles15m []Candle,
	candles3m []Candle,
	candles1m []Candle,
) *MultiTimeframeSignal {
	
	// Minimum data requirements
	if len(candles4h) < 50 || len(candles1h) < 50 || len(candles15m) < 50 {
		return nil
	}
	if len(candles3m) < 20 || len(candles1m) < 20 {
		return nil
	}
	
	// Step 1: Get direction from 4h
	htfDirection, htfTrend, htfConfidence := analyze4hDirection(candles4h)
	if htfDirection == "" {
		return nil // No clear direction
	}
	
	log.Printf("📊 4h Direction: %s (Trend: %s, Confidence: %.1f%%)", 
		htfDirection, htfTrend, htfConfidence)
	
	// Step 2: Find key levels on 1h (OB, FVG, Liquidity)
	orderBlocks := findOrderBlocks1h(candles1h, htfDirection)
	fvgs := findFVGs1h(candles1h, htfDirection)
	liquidityZones := findLiquidityZones1h(candles1h)
	
	if len(orderBlocks) == 0 && len(fvgs) == 0 {
		return nil // No key levels found
	}
	
	log.Printf("🎯 1h Levels: %d OBs, %d FVGs, %d Liquidity Zones", 
		len(orderBlocks), len(fvgs), len(liquidityZones))
	
	// Step 3: Analyze volume and delta on 15m
	volumeProfile := analyzeVolumeProfile15m(candles15m)
	deltaAnalysis := analyzeDelta15m(candles15m)
	insideFVG := checkInsideFVG15m(candles15m, fvgs)
	
	// Check if volume confirms direction
	if !volumeConfirmsDirection(volumeProfile, deltaAnalysis, htfDirection) {
		return nil // Volume doesn't confirm
	}
	
	log.Printf("📈 15m Analysis: Delta %s, Inside FVG: %v", 
		deltaAnalysis.DeltaDirection, insideFVG)
	
	// Step 4: Find precise entry on 3m
	preciseEntry3m, quality3m := findPreciseEntry3m(candles3m, htfDirection, orderBlocks, fvgs)
	
	// Step 5: Refine entry on 1m for best execution
	optimalEntry, entryPrice, entryQuality := refineEntry1m(
		candles1m, htfDirection, preciseEntry3m, orderBlocks, fvgs)
	
	if !optimalEntry {
		return nil // No optimal entry found
	}
	
	log.Printf("🎯 Entry Found: %.2f (Quality: %.1f%%, 3m: %.1f%%)", 
		entryPrice, entryQuality, quality3m)
	
	// Calculate stops and targets
	stopLoss, tp1, tp2, tp3 := calculateMultiTimeframeTargets(
		entryPrice, htfDirection, candles4h, candles1h)
	
	// Calculate risk/reward
	risk := math.Abs(entryPrice - stopLoss)
	reward := math.Abs(entryPrice - tp1)
	rr := reward / risk
	
	// Minimum RR requirement (OPTIMIZED: Balanced)
	if rr < 2.5 { // Reduced from 3.0 for more trades
		return nil // RR too low for multi-timeframe setup
	}
	
	// Calculate confluence score
	confluence := calculateMultiTimeframeConfluence(
		htfConfidence, len(orderBlocks), len(fvgs), 
		volumeProfile, deltaAnalysis, insideFVG, entryQuality)
	
	// Minimum confluence for multi-timeframe (OPTIMIZED: Balanced)
	if confluence < 5 { // Reduced from 6 for more trades
		return nil // Not enough confluence
	}
	
	signal := &MultiTimeframeSignal{
		EntryTimeframe: "1m",
		EntryPrice:     entryPrice,
		EntryTime:      candles1m[len(candles1m)-1].Timestamp,
		HTFDirection:   htfDirection,
		HTFTrend:       htfTrend,
		HTFConfidence:  htfConfidence,
		OrderBlocks:    orderBlocks,
		FVGs:           fvgs,
		LiquidityZones: liquidityZones,
		VolumeProfile:  volumeProfile,
		DeltaAnalysis:  deltaAnalysis,
		InsideFVG:      insideFVG,
		PreciseEntry:   preciseEntry3m,
		OptimalEntry:   optimalEntry,
		EntryQuality:   entryQuality,
		Type:           htfDirection,
		StopLoss:       stopLoss,
		TP1:            tp1,
		TP2:            tp2,
		TP3:            tp3,
		RiskReward:     rr,
		Confluence:     confluence,
		Strength:       (htfConfidence + entryQuality) / 2,
	}
	
	log.Printf("✅ Multi-TF Signal: %s | Entry: %.2f | RR: %.2f | Confluence: %d", 
		signal.Type, signal.EntryPrice, signal.RiskReward, signal.Confluence)
	
	return signal
}

// analyze4hDirection determines overall market direction from 4h
// OPTIMIZED: Lower confidence threshold for more signals
func analyze4hDirection(candles []Candle) (string, string, float64) {
	if len(candles) < 50 {
		return "", "", 0
	}
	
	idx := len(candles) - 1
	
	// Calculate EMAs
	ema20 := calculateEMA(candles[:idx+1], 20)
	ema50 := calculateEMA(candles[:idx+1], 50)
	ema200 := calculateEMA(candles[:idx+1], 200)
	
	currentPrice := candles[idx].Close
	
	// Determine trend (OPTIMIZED: More lenient)
	var trend string
	var confidence float64
	
	if ema20 > ema50 && ema50 > ema200 && currentPrice > ema20 {
		trend = "strong_bullish"
		confidence = 85.0 // Reduced from 90
	} else if ema20 > ema50 && currentPrice > ema20 {
		trend = "bullish"
		confidence = 70.0 // Reduced from 75
	} else if ema20 < ema50 && ema50 < ema200 && currentPrice < ema20 {
		trend = "strong_bearish"
		confidence = 85.0 // Reduced from 90
	} else if ema20 < ema50 && currentPrice < ema20 {
		trend = "bearish"
		confidence = 70.0 // Reduced from 75
	} else {
		trend = "ranging"
		confidence = 45.0 // Reduced from 50
	}
	
	// Check for break of structure
	bos := checkBreakOfStructure(candles, idx)
	if bos != "" {
		confidence += 10.0
	}
	
	// Determine direction
	var direction string
	if trend == "strong_bullish" || trend == "bullish" {
		direction = "BUY"
	} else if trend == "strong_bearish" || trend == "bearish" {
		direction = "SELL"
	}
	
	return direction, trend, confidence
}

// findOrderBlocks1h finds order blocks on 1h timeframe
func findOrderBlocks1h(candles []Candle, direction string) []MTFOrderBlock {
	orderBlocks := []MTFOrderBlock{}
	
	if len(candles) < 20 {
		return orderBlocks
	}
	
	// Look for last 10 candles
	for i := len(candles) - 10; i < len(candles)-1; i++ {
		if i < 2 {
			continue
		}
		
		prev := candles[i-1]
		curr := candles[i]
		next := candles[i+1]
		
		// Bullish order block (for BUY direction)
		if direction == "BUY" {
			if prev.Close < prev.Open && // Bearish candle
			   curr.Close > curr.Open && // Bullish engulfing
			   curr.Close > prev.Open &&
			   next.Close > curr.Close { // Continuation
				
				ob := MTFOrderBlock{
					High:      curr.High,
					Low:       curr.Low,
					Type:      "bullish",
					Strength:  calculateOBStrength(prev, curr, next),
					Timeframe: "1h",
					Tested:    false,
				}
				orderBlocks = append(orderBlocks, ob)
			}
		}
		
		// Bearish order block (for SELL direction)
		if direction == "SELL" {
			if prev.Close > prev.Open && // Bullish candle
			   curr.Close < curr.Open && // Bearish engulfing
			   curr.Close < prev.Open &&
			   next.Close < curr.Close { // Continuation
				
				ob := MTFOrderBlock{
					High:      curr.High,
					Low:       curr.Low,
					Type:      "bearish",
					Strength:  calculateOBStrength(prev, curr, next),
					Timeframe: "1h",
					Tested:    false,
				}
				orderBlocks = append(orderBlocks, ob)
			}
		}
	}
	
	return orderBlocks
}

// findFVGs1h finds fair value gaps on 1h
func findFVGs1h(candles []Candle, direction string) []MTFFairValueGap {
	fvgs := []MTFFairValueGap{}
	
	if len(candles) < 10 {
		return fvgs
	}
	
	for i := len(candles) - 10; i < len(candles)-1; i++ {
		if i < 2 {
			continue
		}
		
		// Bullish FVG
		if direction == "BUY" && candles[i].Low > candles[i-2].High {
			fvg := MTFFairValueGap{
				High:      candles[i].Low,
				Low:       candles[i-2].High,
				Type:      "bullish",
				Filled:    false,
				Timeframe: "1h",
			}
			fvgs = append(fvgs, fvg)
		}
		
		// Bearish FVG
		if direction == "SELL" && candles[i].High < candles[i-2].Low {
			fvg := MTFFairValueGap{
				High:      candles[i-2].Low,
				Low:       candles[i].High,
				Type:      "bearish",
				Filled:    false,
				Timeframe: "1h",
			}
			fvgs = append(fvgs, fvg)
		}
	}
	
	return fvgs
}

// findLiquidityZones1h finds liquidity zones
func findLiquidityZones1h(candles []Candle) []LiquidityZone {
	zones := []LiquidityZone{}
	
	if len(candles) < 20 {
		return zones
	}
	
	// Find swing highs and lows
	for i := len(candles) - 15; i < len(candles)-5; i++ {
		if i < 5 {
			continue
		}
		
		// Swing high
		if candles[i].High > candles[i-1].High &&
		   candles[i].High > candles[i+1].High {
			zone := LiquidityZone{
				Price:    candles[i].High,
				Type:     "high",
				Swept:    false,
				Strength: calculateLiquidityStrength(candles, i),
			}
			zones = append(zones, zone)
		}
		
		// Swing low
		if candles[i].Low < candles[i-1].Low &&
		   candles[i].Low < candles[i+1].Low {
			zone := LiquidityZone{
				Price:    candles[i].Low,
				Type:     "low",
				Swept:    false,
				Strength: calculateLiquidityStrength(candles, i),
			}
			zones = append(zones, zone)
		}
	}
	
	return zones
}

// analyzeVolumeProfile15m analyzes volume on 15m
func analyzeVolumeProfile15m(candles []Candle) MTFVolumeProfile {
	if len(candles) < 20 {
		return MTFVolumeProfile{}
	}
	
	// Calculate average volume
	var totalVolume, buyVolume, sellVolume float64
	for i := len(candles) - 20; i < len(candles); i++ {
		totalVolume += candles[i].Volume
		if candles[i].Close > candles[i].Open {
			buyVolume += candles[i].Volume
		} else {
			sellVolume += candles[i].Volume
		}
	}
	
	avgVolume := totalVolume / 20
	currentVolume := candles[len(candles)-1].Volume
	
	return MTFVolumeProfile{
		POC:        candles[len(candles)-1].Close,
		VAH:        candles[len(candles)-1].High,
		VAL:        candles[len(candles)-1].Low,
		HighVolume: currentVolume > avgVolume*1.5,
		BuyVolume:  buyVolume,
		SellVolume: sellVolume,
	}
}

// analyzeDelta15m analyzes order flow delta
func analyzeDelta15m(candles []Candle) MTFDeltaAnalysis {
	if len(candles) < 10 {
		return MTFDeltaAnalysis{}
	}
	
	var cumulativeDelta float64
	for i := len(candles) - 10; i < len(candles); i++ {
		if candles[i].Close > candles[i].Open {
			cumulativeDelta += candles[i].Volume
		} else {
			cumulativeDelta -= candles[i].Volume
		}
	}
	
	direction := "positive"
	if cumulativeDelta < 0 {
		direction = "negative"
	}
	
	strength := math.Abs(cumulativeDelta) / 10
	
	return MTFDeltaAnalysis{
		CumulativeDelta: cumulativeDelta,
		DeltaDirection:  direction,
		Strength:        strength,
		Divergence:      false,
	}
}

// checkInsideFVG15m checks if price is inside an FVG
func checkInsideFVG15m(candles []Candle, fvgs []MTFFairValueGap) bool {
	if len(candles) == 0 || len(fvgs) == 0 {
		return false
	}
	
	currentPrice := candles[len(candles)-1].Close
	
	for _, fvg := range fvgs {
		if currentPrice >= fvg.Low && currentPrice <= fvg.High {
			return true
		}
	}
	
	return false
}

// volumeConfirmsDirection checks if volume confirms the direction
func volumeConfirmsDirection(vp MTFVolumeProfile, delta MTFDeltaAnalysis, direction string) bool {
	if direction == "BUY" {
		return delta.DeltaDirection == "positive" && vp.BuyVolume > vp.SellVolume
	}
	return delta.DeltaDirection == "negative" && vp.SellVolume > vp.BuyVolume
}

// findPreciseEntry3m finds precise entry on 3m
func findPreciseEntry3m(candles []Candle, direction string, obs []MTFOrderBlock, fvgs []MTFFairValueGap) (float64, float64) {
	if len(candles) < 5 {
		return 0, 0
	}
	
	currentPrice := candles[len(candles)-1].Close
	quality := 50.0
	
	// Check if near order block
	for _, ob := range obs {
		if currentPrice >= ob.Low && currentPrice <= ob.High {
			quality += 20.0
		}
	}
	
	// Check if near FVG
	for _, fvg := range fvgs {
		if currentPrice >= fvg.Low && currentPrice <= fvg.High {
			quality += 15.0
		}
	}
	
	// Check for bullish/bearish structure
	if direction == "BUY" && candles[len(candles)-1].Close > candles[len(candles)-1].Open {
		quality += 10.0
	} else if direction == "SELL" && candles[len(candles)-1].Close < candles[len(candles)-1].Open {
		quality += 10.0
	}
	
	if quality > 100 {
		quality = 100
	}
	
	return currentPrice, quality
}

// refineEntry1m refines entry on 1m for optimal execution
func refineEntry1m(candles []Candle, direction string, entry3m float64, obs []MTFOrderBlock, fvgs []MTFFairValueGap) (bool, float64, float64) {
	if len(candles) < 3 {
		return false, 0, 0
	}
	
	currentPrice := candles[len(candles)-1].Close
	quality := 60.0
	
	// Look for optimal entry pattern
	optimal := false
	
	if direction == "BUY" {
		// Look for bullish engulfing or pin bar at support
		if isBullishEngulfing(candles, len(candles)-1) {
			optimal = true
			quality += 20.0
		}
		if isPinBar(candles, len(candles)-1) && candles[len(candles)-1].Close > candles[len(candles)-1].Open {
			optimal = true
			quality += 15.0
		}
	} else {
		// Look for bearish engulfing or pin bar at resistance
		if isBearishEngulfing(candles, len(candles)-1) {
			optimal = true
			quality += 20.0
		}
		if isPinBar(candles, len(candles)-1) && candles[len(candles)-1].Close < candles[len(candles)-1].Open {
			optimal = true
			quality += 15.0
		}
	}
	
	// Check proximity to key levels
	for _, ob := range obs {
		distance := math.Abs(currentPrice - (ob.High+ob.Low)/2)
		if distance < (ob.High-ob.Low)*0.3 {
			quality += 10.0
		}
	}
	
	if quality > 100 {
		quality = 100
	}
	
	return optimal, currentPrice, quality
}

// calculateMultiTimeframeTargets calculates stops and targets
func calculateMultiTimeframeTargets(entry float64, direction string, candles4h, candles1h []Candle) (float64, float64, float64, float64) {
	atr4h := calculateATR(candles4h, len(candles4h)-1)
	atr1h := calculateATR(candles1h, len(candles1h)-1)
	
	var stopLoss, tp1, tp2, tp3 float64
	
	if direction == "BUY" {
		stopLoss = entry - (atr1h * 1.5) // Tighter stop on 1h ATR
		tp1 = entry + (atr4h * 4.0)      // Targets based on 4h ATR
		tp2 = entry + (atr4h * 6.0)
		tp3 = entry + (atr4h * 8.0)
	} else {
		stopLoss = entry + (atr1h * 1.5)
		tp1 = entry - (atr4h * 4.0)
		tp2 = entry - (atr4h * 6.0)
		tp3 = entry - (atr4h * 8.0)
	}
	
	return stopLoss, tp1, tp2, tp3
}

// calculateMultiTimeframeConfluence calculates total confluence
func calculateMultiTimeframeConfluence(htfConf float64, numOBs, numFVGs int, vp MTFVolumeProfile, delta MTFDeltaAnalysis, insideFVG bool, entryQuality float64) int {
	confluence := 0
	
	// HTF direction
	if htfConf >= 80 {
		confluence += 2
	} else if htfConf >= 70 {
		confluence++
	}
	
	// Order blocks
	confluence += numOBs
	
	// FVGs
	confluence += numFVGs
	
	// Volume
	if vp.HighVolume {
		confluence++
	}
	
	// Delta
	if delta.Strength > 1000 {
		confluence++
	}
	
	// Inside FVG
	if insideFVG {
		confluence++
	}
	
	// Entry quality
	if entryQuality >= 80 {
		confluence += 2
	} else if entryQuality >= 70 {
		confluence++
	}
	
	return confluence
}

// Helper functions
func calculateOBStrength(prev, curr, next Candle) float64 {
	bodySize := math.Abs(curr.Close - curr.Open)
	prevBodySize := math.Abs(prev.Close - prev.Open)
	
	if prevBodySize == 0 {
		return 50.0
	}
	
	ratio := bodySize / prevBodySize
	return math.Min(ratio*50, 100.0)
}

func calculateLiquidityStrength(candles []Candle, idx int) float64 {
	// Count how many times price tested this level
	tests := 0
	level := candles[idx].High
	
	for i := idx + 1; i < len(candles) && i < idx+10; i++ {
		if math.Abs(candles[i].High-level) < level*0.001 {
			tests++
		}
	}
	
	return float64(tests) * 20.0
}

func checkBreakOfStructure(candles []Candle, idx int) string {
	if idx < 20 {
		return ""
	}
	
	// Find recent high/low
	recentHigh := candles[idx-20].High
	recentLow := candles[idx-20].Low
	
	for i := idx - 19; i < idx; i++ {
		if candles[i].High > recentHigh {
			recentHigh = candles[i].High
		}
		if candles[i].Low < recentLow {
			recentLow = candles[i].Low
		}
	}
	
	// Check if current breaks structure
	if candles[idx].Close > recentHigh {
		return "bullish_bos"
	}
	if candles[idx].Close < recentLow {
		return "bearish_bos"
	}
	
	return ""
}
