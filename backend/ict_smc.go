package main

import (
	"math"
)

// ==================== ICT/SMC CONCEPTS ====================

// OrderBlock represents a bullish or bearish order block
type OrderBlock struct {
	Type       string  // "bullish" or "bearish"
	High       float64
	Low        float64
	MidPoint   float64
	Volume     float64
	Strength   float64
	Mitigated  bool
	CandleIdx  int
}

// FairValueGap represents an imbalance/FVG
type FairValueGap struct {
	Type      string  // "bullish" or "bearish"
	High      float64
	Low       float64
	MidPoint  float64
	Size      float64
	Filled    bool
	CandleIdx int
}

// LiquidityLevel represents buy-side or sell-side liquidity
type LiquidityLevel struct {
	Type     string  // "buyside" or "sellside"
	Price    float64
	Strength float64
	Swept    bool
}

// MarketStructure represents market structure shift/break
type MarketStructure struct {
	Trend          string  // "bullish", "bearish", "ranging"
	LastSwingHigh  float64
	LastSwingLow   float64
	BOS            bool    // Break of Structure
	CHOCH          bool    // Change of Character
	MSS            bool    // Market Structure Shift
}

// KillZone represents ICT kill zones
type KillZone struct {
	Name   string
	Active bool
	Hour   int
}

// ICTAnalysis holds all ICT/SMC analysis results
type ICTAnalysis struct {
	OrderBlocks     []OrderBlock
	FairValueGaps   []FairValueGap
	Liquidity       []LiquidityLevel
	Structure       MarketStructure
	KillZone        *KillZone
	PremiumDiscount string // "premium", "discount", "equilibrium"
	OTE             bool   // Optimal Trade Entry zone
	Confluence      int    // Number of confluent factors
}

// ==================== ORDER BLOCKS ====================

// FindOrderBlocks identifies bullish and bearish order blocks
func FindOrderBlocks(candles []Candle) []OrderBlock {
	var orderBlocks []OrderBlock
	
	if len(candles) < 10 {
		return orderBlocks
	}
	
	for i := 2; i < len(candles)-1; i++ {
		prev := candles[i-1]
		curr := candles[i]
		next := candles[i+1]
		
		// BULLISH ORDER BLOCK
		// Last down candle before a strong up move
		if prev.Close < prev.Open && // Previous bearish
		   curr.Close < curr.Open && // Current bearish (the OB)
		   next.Close > next.Open && // Next bullish
		   next.Close > curr.High {  // Strong move up
			
			avgVolume := (prev.Volume + curr.Volume + next.Volume) / 3
			strength := (next.Close - curr.Low) / (curr.High - curr.Low)
			
			ob := OrderBlock{
				Type:      "bullish",
				High:      curr.High,
				Low:       curr.Low,
				MidPoint:  (curr.High + curr.Low) / 2,
				Volume:    curr.Volume,
				Strength:  math.Min(strength*100, 100),
				Mitigated: false,
				CandleIdx: i,
			}
			
			// Higher volume = stronger OB
			if curr.Volume > avgVolume*1.2 {
				ob.Strength += 20
			}
			
			orderBlocks = append(orderBlocks, ob)
		}
		
		// BEARISH ORDER BLOCK
		// Last up candle before a strong down move
		if prev.Close > prev.Open && // Previous bullish
		   curr.Close > curr.Open && // Current bullish (the OB)
		   next.Close < next.Open && // Next bearish
		   next.Close < curr.Low {   // Strong move down
			
			avgVolume := (prev.Volume + curr.Volume + next.Volume) / 3
			strength := (curr.High - next.Close) / (curr.High - curr.Low)
			
			ob := OrderBlock{
				Type:      "bearish",
				High:      curr.High,
				Low:       curr.Low,
				MidPoint:  (curr.High + curr.Low) / 2,
				Volume:    curr.Volume,
				Strength:  math.Min(strength*100, 100),
				Mitigated: false,
				CandleIdx: i,
			}
			
			if curr.Volume > avgVolume*1.2 {
				ob.Strength += 20
			}
			
			orderBlocks = append(orderBlocks, ob)
		}
	}
	
	return orderBlocks
}

// ==================== FAIR VALUE GAPS ====================

// FindFairValueGaps identifies imbalances/FVGs
func FindFairValueGaps(candles []Candle) []FairValueGap {
	var fvgs []FairValueGap
	
	if len(candles) < 3 {
		return fvgs
	}
	
	for i := 1; i < len(candles)-1; i++ {
		prev := candles[i-1]
		curr := candles[i]
		next := candles[i+1]
		
		// BULLISH FVG (gap up)
		// Previous candle high < Next candle low
		if prev.High < next.Low {
			fvg := FairValueGap{
				Type:      "bullish",
				High:      next.Low,
				Low:       prev.High,
				MidPoint:  (next.Low + prev.High) / 2,
				Size:      next.Low - prev.High,
				Filled:    false,
				CandleIdx: i,
			}
			fvgs = append(fvgs, fvg)
		}
		
		// BEARISH FVG (gap down)
		// Previous candle low > Next candle high
		if prev.Low > next.High {
			fvg := FairValueGap{
				Type:      "bearish",
				High:      prev.Low,
				Low:       next.High,
				MidPoint:  (prev.Low + next.High) / 2,
				Size:      prev.Low - next.High,
				Filled:    false,
				CandleIdx: i,
			}
			fvgs = append(fvgs, fvg)
		}
		
		// Check for FVG with wicks (less strict)
		bodyPrev := math.Max(prev.Open, prev.Close)
		bodyNext := math.Min(next.Open, next.Close)
		
		if bodyPrev < bodyNext && curr.Close > curr.Open {
			// Bullish FVG with body
			gap := bodyNext - bodyPrev
			if gap > (curr.High-curr.Low)*0.3 { // Significant gap
				fvg := FairValueGap{
					Type:      "bullish",
					High:      bodyNext,
					Low:       bodyPrev,
					MidPoint:  (bodyNext + bodyPrev) / 2,
					Size:      gap,
					Filled:    false,
					CandleIdx: i,
				}
				fvgs = append(fvgs, fvg)
			}
		}
	}
	
	return fvgs
}

// ==================== LIQUIDITY ====================

// FindLiquidity identifies liquidity pools (equal highs/lows, swing points)
func FindLiquidity(candles []Candle) []LiquidityLevel {
	var liquidity []LiquidityLevel
	
	if len(candles) < 20 {
		return liquidity
	}
	
	// Find swing highs and lows
	for i := 2; i < len(candles)-2; i++ {
		// Swing High (potential sell-side liquidity above)
		if candles[i].High > candles[i-1].High &&
		   candles[i].High > candles[i-2].High &&
		   candles[i].High > candles[i+1].High &&
		   candles[i].High > candles[i+2].High {
			
			liq := LiquidityLevel{
				Type:     "buyside",
				Price:    candles[i].High,
				Strength: 70,
				Swept:    false,
			}
			liquidity = append(liquidity, liq)
		}
		
		// Swing Low (potential buy-side liquidity below)
		if candles[i].Low < candles[i-1].Low &&
		   candles[i].Low < candles[i-2].Low &&
		   candles[i].Low < candles[i+1].Low &&
		   candles[i].Low < candles[i+2].Low {
			
			liq := LiquidityLevel{
				Type:     "sellside",
				Price:    candles[i].Low,
				Strength: 70,
				Swept:    false,
			}
			liquidity = append(liquidity, liq)
		}
	}
	
	// Find equal highs (strong liquidity)
	tolerance := (candles[len(candles)-1].High - candles[len(candles)-1].Low) * 0.1
	for i := 0; i < len(candles)-5; i++ {
		for j := i + 3; j < len(candles); j++ {
			if math.Abs(candles[i].High-candles[j].High) < tolerance {
				liq := LiquidityLevel{
					Type:     "buyside",
					Price:    (candles[i].High + candles[j].High) / 2,
					Strength: 90, // Equal highs are strong liquidity
					Swept:    false,
				}
				liquidity = append(liquidity, liq)
				break
			}
		}
	}
	
	// Find equal lows (strong liquidity)
	for i := 0; i < len(candles)-5; i++ {
		for j := i + 3; j < len(candles); j++ {
			if math.Abs(candles[i].Low-candles[j].Low) < tolerance {
				liq := LiquidityLevel{
					Type:     "sellside",
					Price:    (candles[i].Low + candles[j].Low) / 2,
					Strength: 90, // Equal lows are strong liquidity
					Swept:    false,
				}
				liquidity = append(liquidity, liq)
				break
			}
		}
	}
	
	return liquidity
}

// ==================== MARKET STRUCTURE ====================

// AnalyzeMarketStructure determines trend and structure breaks
func AnalyzeMarketStructure(candles []Candle) MarketStructure {
	ms := MarketStructure{
		Trend: "ranging",
	}
	
	if len(candles) < 20 {
		return ms
	}
	
	// Find swing highs and lows
	var swingHighs, swingLows []float64
	var swingHighIdx, swingLowIdx []int
	
	for i := 2; i < len(candles)-2; i++ {
		// Swing High
		if candles[i].High > candles[i-1].High &&
		   candles[i].High > candles[i-2].High &&
		   candles[i].High > candles[i+1].High &&
		   candles[i].High > candles[i+2].High {
			swingHighs = append(swingHighs, candles[i].High)
			swingHighIdx = append(swingHighIdx, i)
		}
		
		// Swing Low
		if candles[i].Low < candles[i-1].Low &&
		   candles[i].Low < candles[i-2].Low &&
		   candles[i].Low < candles[i+1].Low &&
		   candles[i].Low < candles[i+2].Low {
			swingLows = append(swingLows, candles[i].Low)
			swingLowIdx = append(swingLowIdx, i)
		}
	}
	
	if len(swingHighs) < 2 || len(swingLows) < 2 {
		return ms
	}
	
	// Set last swing points
	ms.LastSwingHigh = swingHighs[len(swingHighs)-1]
	ms.LastSwingLow = swingLows[len(swingLows)-1]
	
	// Determine trend
	// Higher highs and higher lows = bullish
	// Lower highs and lower lows = bearish
	lastHH := swingHighs[len(swingHighs)-1]
	prevHH := swingHighs[len(swingHighs)-2]
	lastLL := swingLows[len(swingLows)-1]
	prevLL := swingLows[len(swingLows)-2]
	
	if lastHH > prevHH && lastLL > prevLL {
		ms.Trend = "bullish"
	} else if lastHH < prevHH && lastLL < prevLL {
		ms.Trend = "bearish"
	}
	
	// Check for BOS (Break of Structure)
	currentPrice := candles[len(candles)-1].Close
	
	if ms.Trend == "bullish" && currentPrice > ms.LastSwingHigh {
		ms.BOS = true
	}
	if ms.Trend == "bearish" && currentPrice < ms.LastSwingLow {
		ms.BOS = true
	}
	
	// Check for CHOCH (Change of Character)
	if ms.Trend == "bullish" && currentPrice < ms.LastSwingLow {
		ms.CHOCH = true
		ms.MSS = true
	}
	if ms.Trend == "bearish" && currentPrice > ms.LastSwingHigh {
		ms.CHOCH = true
		ms.MSS = true
	}
	
	return ms
}

// ==================== PREMIUM/DISCOUNT ====================

// CalculatePremiumDiscount determines if price is in premium or discount
func CalculatePremiumDiscount(candles []Candle) (string, float64) {
	if len(candles) < 20 {
		return "equilibrium", 0.5
	}
	
	// Find range high and low
	rangeHigh := candles[0].High
	rangeLow := candles[0].Low
	
	for _, c := range candles {
		if c.High > rangeHigh {
			rangeHigh = c.High
		}
		if c.Low < rangeLow {
			rangeLow = c.Low
		}
	}
	
	currentPrice := candles[len(candles)-1].Close
	rangeSize := rangeHigh - rangeLow
	
	if rangeSize == 0 {
		return "equilibrium", 0.5
	}
	
	// Calculate position (0 = low, 1 = high)
	position := (currentPrice - rangeLow) / rangeSize
	
	// Premium = above 0.5 (above equilibrium)
	// Discount = below 0.5 (below equilibrium)
	if position > 0.7 {
		return "premium", position
	} else if position < 0.3 {
		return "discount", position
	}
	
	return "equilibrium", position
}

// ==================== OPTIMAL TRADE ENTRY ====================

// IsInOTE checks if price is in Optimal Trade Entry zone (61.8% - 78.6% fib)
func IsInOTE(candles []Candle) bool {
	if len(candles) < 20 {
		return false
	}
	
	// Find recent swing high and low
	swingHigh := candles[0].High
	swingLow := candles[0].Low
	
	for _, c := range candles[:len(candles)/2] {
		if c.High > swingHigh {
			swingHigh = c.High
		}
		if c.Low < swingLow {
			swingLow = c.Low
		}
	}
	
	currentPrice := candles[len(candles)-1].Close
	rangeSize := swingHigh - swingLow
	
	if rangeSize == 0 {
		return false
	}
	
	// Calculate retracement level
	retracement := (swingHigh - currentPrice) / rangeSize
	
	// OTE zone is 61.8% - 78.6% retracement
	return retracement >= 0.618 && retracement <= 0.786
}

// ==================== KILL ZONES ====================

// GetKillZone returns the current kill zone based on hour (UTC)
func GetKillZone(hour int) *KillZone {
	// Asian Kill Zone: 00:00 - 04:00 UTC
	if hour >= 0 && hour < 4 {
		return &KillZone{Name: "Asian", Active: true, Hour: hour}
	}
	
	// London Kill Zone: 07:00 - 10:00 UTC
	if hour >= 7 && hour < 10 {
		return &KillZone{Name: "London", Active: true, Hour: hour}
	}
	
	// New York Kill Zone: 13:00 - 16:00 UTC
	if hour >= 13 && hour < 16 {
		return &KillZone{Name: "NewYork", Active: true, Hour: hour}
	}
	
	// London Close: 15:00 - 17:00 UTC
	if hour >= 15 && hour < 17 {
		return &KillZone{Name: "LondonClose", Active: true, Hour: hour}
	}
	
	return nil
}

// ==================== FULL ICT ANALYSIS ====================

// PerformICTAnalysis performs complete ICT/SMC analysis
func PerformICTAnalysis(candles []Candle) *ICTAnalysis {
	analysis := &ICTAnalysis{
		Confluence: 0,
	}
	
	// Find Order Blocks
	analysis.OrderBlocks = FindOrderBlocks(candles)
	
	// Find Fair Value Gaps
	analysis.FairValueGaps = FindFairValueGaps(candles)
	
	// Find Liquidity
	analysis.Liquidity = FindLiquidity(candles)
	
	// Analyze Market Structure
	analysis.Structure = AnalyzeMarketStructure(candles)
	
	// Calculate Premium/Discount
	analysis.PremiumDiscount, _ = CalculatePremiumDiscount(candles)
	
	// Check OTE
	analysis.OTE = IsInOTE(candles)
	
	// Calculate confluence
	if len(analysis.OrderBlocks) > 0 {
		analysis.Confluence++
	}
	if len(analysis.FairValueGaps) > 0 {
		analysis.Confluence++
	}
	if analysis.Structure.BOS || analysis.Structure.CHOCH {
		analysis.Confluence++
	}
	if analysis.OTE {
		analysis.Confluence++
	}
	if analysis.PremiumDiscount == "discount" || analysis.PremiumDiscount == "premium" {
		analysis.Confluence++
	}
	
	return analysis
}
