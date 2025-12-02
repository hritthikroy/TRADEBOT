package main

import (
	"math"
)

// ==================== LIQUIDITY SWEEP ====================
// Professional institutional liquidity hunting patterns

// LiquiditySweep represents a liquidity sweep event
type LiquiditySweep struct {
	Type           string  // "buyside" or "sellside"
	SweepPrice     float64 // Price where liquidity was swept
	ReturnPrice    float64 // Price returned to after sweep
	SweepSize      float64 // How far past the level
	Volume         float64 // Volume on sweep
	Strength       float64 // Sweep strength (0-100)
	CandleIdx      int
	Confirmed      bool    // Reversal confirmed
	TargetLevel    float64 // Original liquidity level
	TimeToReverse  int     // Candles to reverse
}

// LiquidityPool represents a pool of liquidity
type LiquidityPool struct {
	Type       string  // "buyside" (above) or "sellside" (below)
	Price      float64
	Strength   float64
	TouchCount int
	Swept      bool
	SweptAt    int // Candle index when swept
}

// LiquiditySweepAnalysis holds all liquidity sweep analysis
type LiquiditySweepAnalysis struct {
	Sweeps          []LiquiditySweep
	Pools           []LiquidityPool
	RecentSweep     *LiquiditySweep
	NearestBuyside  *LiquidityPool
	NearestSellside *LiquidityPool
	SweepInProgress bool
	ExpectedMove    string
}

// ==================== LIQUIDITY POOL DETECTION ====================

// FindLiquidityPools identifies all liquidity pools
func FindLiquidityPools(candles []Candle) []LiquidityPool {
	var pools []LiquidityPool
	
	if len(candles) < 20 {
		return pools
	}
	
	// Find swing highs (buyside liquidity above)
	for i := 2; i < len(candles)-2; i++ {
		// Swing High
		if candles[i].High > candles[i-1].High &&
			candles[i].High > candles[i-2].High &&
			candles[i].High > candles[i+1].High &&
			candles[i].High > candles[i+2].High {
			
			pool := LiquidityPool{
				Type:       "buyside",
				Price:      candles[i].High,
				Strength:   70,
				TouchCount: 1,
				Swept:      false,
			}
			
			// Check for equal highs (stronger liquidity)
			for j := i + 3; j < len(candles); j++ {
				if math.Abs(candles[j].High-pool.Price) < pool.Price*0.001 {
					pool.TouchCount++
					pool.Strength += 10
				}
			}
			
			pools = append(pools, pool)
		}
		
		// Swing Low
		if candles[i].Low < candles[i-1].Low &&
			candles[i].Low < candles[i-2].Low &&
			candles[i].Low < candles[i+1].Low &&
			candles[i].Low < candles[i+2].Low {
			
			pool := LiquidityPool{
				Type:       "sellside",
				Price:      candles[i].Low,
				Strength:   70,
				TouchCount: 1,
				Swept:      false,
			}
			
			// Check for equal lows
			for j := i + 3; j < len(candles); j++ {
				if math.Abs(candles[j].Low-pool.Price) < pool.Price*0.001 {
					pool.TouchCount++
					pool.Strength += 10
				}
			}
			
			pools = append(pools, pool)
		}
	}
	
	// Find trendline liquidity
	pools = append(pools, findTrendlineLiquidity(candles)...)
	
	// Find round number liquidity
	pools = append(pools, findRoundNumberLiquidity(candles)...)
	
	return pools
}

// findTrendlineLiquidity finds liquidity along trendlines
func findTrendlineLiquidity(candles []Candle) []LiquidityPool {
	var pools []LiquidityPool
	
	if len(candles) < 20 {
		return pools
	}
	
	// Find ascending trendline (support)
	lows := []float64{}
	for _, c := range candles {
		lows = append(lows, c.Low)
	}
	
	// Simple trendline: connect first and last significant lows
	firstLow := lows[0]
	lastLow := lows[len(lows)-1]
	
	if lastLow > firstLow {
		// Ascending trendline = sellside liquidity below
		pool := LiquidityPool{
			Type:     "sellside",
			Price:    lastLow,
			Strength: 65,
			Swept:    false,
		}
		pools = append(pools, pool)
	}
	
	// Find descending trendline (resistance)
	highs := []float64{}
	for _, c := range candles {
		highs = append(highs, c.High)
	}
	
	firstHigh := highs[0]
	lastHigh := highs[len(highs)-1]
	
	if lastHigh < firstHigh {
		// Descending trendline = buyside liquidity above
		pool := LiquidityPool{
			Type:     "buyside",
			Price:    lastHigh,
			Strength: 65,
			Swept:    false,
		}
		pools = append(pools, pool)
	}
	
	return pools
}

// findRoundNumberLiquidity finds liquidity at round numbers
func findRoundNumberLiquidity(candles []Candle) []LiquidityPool {
	var pools []LiquidityPool
	
	if len(candles) == 0 {
		return pools
	}
	
	currentPrice := candles[len(candles)-1].Close
	
	// Find nearest round numbers
	// For BTC: $1000 levels
	// For smaller assets: $10 or $100 levels
	
	roundLevel := 1000.0
	if currentPrice < 1000 {
		roundLevel = 100.0
	}
	if currentPrice < 100 {
		roundLevel = 10.0
	}
	
	// Round number above
	roundAbove := math.Ceil(currentPrice/roundLevel) * roundLevel
	pools = append(pools, LiquidityPool{
		Type:     "buyside",
		Price:    roundAbove,
		Strength: 60,
		Swept:    false,
	})
	
	// Round number below
	roundBelow := math.Floor(currentPrice/roundLevel) * roundLevel
	pools = append(pools, LiquidityPool{
		Type:     "sellside",
		Price:    roundBelow,
		Strength: 60,
		Swept:    false,
	})
	
	return pools
}

// ==================== LIQUIDITY SWEEP DETECTION ====================

// DetectLiquiditySweeps identifies liquidity sweep events
func DetectLiquiditySweeps(candles []Candle, pools []LiquidityPool) []LiquiditySweep {
	var sweeps []LiquiditySweep
	
	if len(candles) < 10 || len(pools) == 0 {
		return sweeps
	}
	
	for i := 5; i < len(candles)-1; i++ {
		curr := candles[i]
		next := candles[i+1]
		
		for _, pool := range pools {
			// BUYSIDE SWEEP (sweep above, reverse down)
			if pool.Type == "buyside" && !pool.Swept {
				// Price sweeps above pool level
				if curr.High > pool.Price && curr.Close < pool.Price {
					// Confirm reversal
					if next.Close < curr.Close {
						sweep := LiquiditySweep{
							Type:          "buyside",
							SweepPrice:    curr.High,
							ReturnPrice:   curr.Close,
							SweepSize:     curr.High - pool.Price,
							Volume:        curr.Volume,
							CandleIdx:     i,
							TargetLevel:   pool.Price,
							TimeToReverse: 1,
							Confirmed:     true,
						}
						
						// Calculate strength
						sweep.Strength = calculateSweepStrength(sweep, candles, i)
						sweeps = append(sweeps, sweep)
					}
				}
			}
			
			// SELLSIDE SWEEP (sweep below, reverse up)
			if pool.Type == "sellside" && !pool.Swept {
				// Price sweeps below pool level
				if curr.Low < pool.Price && curr.Close > pool.Price {
					// Confirm reversal
					if next.Close > curr.Close {
						sweep := LiquiditySweep{
							Type:          "sellside",
							SweepPrice:    curr.Low,
							ReturnPrice:   curr.Close,
							SweepSize:     pool.Price - curr.Low,
							Volume:        curr.Volume,
							CandleIdx:     i,
							TargetLevel:   pool.Price,
							TimeToReverse: 1,
							Confirmed:     true,
						}
						
						sweep.Strength = calculateSweepStrength(sweep, candles, i)
						sweeps = append(sweeps, sweep)
					}
				}
			}
		}
	}
	
	return sweeps
}

// calculateSweepStrength calculates the strength of a liquidity sweep
func calculateSweepStrength(sweep LiquiditySweep, candles []Candle, idx int) float64 {
	strength := 50.0
	
	// Volume confirmation
	avgVolume := 0.0
	for i := idx - 5; i < idx; i++ {
		if i >= 0 {
			avgVolume += candles[i].Volume
		}
	}
	avgVolume /= 5
	
	if sweep.Volume > avgVolume*1.5 {
		strength += 20
	}
	
	// Quick reversal = stronger
	if sweep.TimeToReverse <= 2 {
		strength += 15
	}
	
	// Large sweep size = stronger
	avgRange := 0.0
	for i := idx - 5; i < idx; i++ {
		if i >= 0 {
			avgRange += candles[i].High - candles[i].Low
		}
	}
	avgRange /= 5
	
	if sweep.SweepSize > avgRange*0.5 {
		strength += 10
	}
	
	// Confirmed reversal
	if sweep.Confirmed {
		strength += 10
	}
	
	return math.Min(strength, 100)
}

// ==================== FULL LIQUIDITY SWEEP ANALYSIS ====================

// PerformLiquiditySweepAnalysis performs complete liquidity sweep analysis
func PerformLiquiditySweepAnalysis(candles []Candle) *LiquiditySweepAnalysis {
	analysis := &LiquiditySweepAnalysis{
		ExpectedMove: "neutral",
	}
	
	// Find liquidity pools
	analysis.Pools = FindLiquidityPools(candles)
	
	// Detect sweeps
	analysis.Sweeps = DetectLiquiditySweeps(candles, analysis.Pools)
	
	// Find most recent sweep
	if len(analysis.Sweeps) > 0 {
		for i := len(analysis.Sweeps) - 1; i >= 0; i-- {
			if analysis.Sweeps[i].CandleIdx >= len(candles)-10 {
				analysis.RecentSweep = &analysis.Sweeps[i]
				break
			}
		}
	}
	
	// Find nearest pools
	if len(candles) > 0 {
		currentPrice := candles[len(candles)-1].Close
		
		minDistBuy := math.MaxFloat64
		minDistSell := math.MaxFloat64
		
		for i := range analysis.Pools {
			pool := &analysis.Pools[i]
			
			if pool.Type == "buyside" && pool.Price > currentPrice {
				dist := pool.Price - currentPrice
				if dist < minDistBuy {
					minDistBuy = dist
					analysis.NearestBuyside = pool
				}
			}
			
			if pool.Type == "sellside" && pool.Price < currentPrice {
				dist := currentPrice - pool.Price
				if dist < minDistSell {
					minDistSell = dist
					analysis.NearestSellside = pool
				}
			}
		}
	}
	
	// Determine expected move based on recent sweep
	if analysis.RecentSweep != nil && analysis.RecentSweep.Confirmed {
		if analysis.RecentSweep.Type == "buyside" {
			analysis.ExpectedMove = "down" // Swept buyside = expect down
		} else {
			analysis.ExpectedMove = "up" // Swept sellside = expect up
		}
	}
	
	return analysis
}

// GetLiquiditySweepSignal returns trading signal based on liquidity sweep
func GetLiquiditySweepSignal(analysis *LiquiditySweepAnalysis, currentPrice float64) (string, float64) {
	if analysis == nil {
		return "neutral", 0
	}
	
	direction := "neutral"
	strength := 0.0
	
	// Recent confirmed sweep = strong signal
	if analysis.RecentSweep != nil && analysis.RecentSweep.Confirmed {
		if analysis.RecentSweep.Type == "sellside" {
			direction = "bullish"
			strength = analysis.RecentSweep.Strength
		} else if analysis.RecentSweep.Type == "buyside" {
			direction = "bearish"
			strength = analysis.RecentSweep.Strength
		}
	}
	
	// Near liquidity pool = potential sweep coming
	if analysis.NearestBuyside != nil {
		dist := analysis.NearestBuyside.Price - currentPrice
		threshold := currentPrice * 0.005 // 0.5%
		
		if dist < threshold {
			// Near buyside liquidity = potential short after sweep
			if direction == "neutral" {
				direction = "bearish"
				strength = 50
			} else if direction == "bearish" {
				strength += 15
			}
		}
	}
	
	if analysis.NearestSellside != nil {
		dist := currentPrice - analysis.NearestSellside.Price
		threshold := currentPrice * 0.005
		
		if dist < threshold {
			// Near sellside liquidity = potential long after sweep
			if direction == "neutral" {
				direction = "bullish"
				strength = 50
			} else if direction == "bullish" {
				strength += 15
			}
		}
	}
	
	return direction, math.Min(strength, 100)
}

// ==================== JUDAS SWING ====================
// ICT Judas Swing: Fake move to trap traders before real move

// JudasSwing represents a Judas swing pattern
type JudasSwing struct {
	Type       string  // "bullish" or "bearish"
	FakeMove   float64 // Size of fake move
	RealMove   float64 // Size of real move (expected)
	TrapPrice  float64 // Price where traders got trapped
	EntryPrice float64 // Optimal entry after Judas
	Strength   float64
	CandleIdx  int
}

// DetectJudasSwing detects Judas swing patterns
func DetectJudasSwing(candles []Candle) *JudasSwing {
	if len(candles) < 20 {
		return nil
	}
	
	// Judas swing typically happens at session open
	// Look for fake breakout followed by reversal
	
	recentCandles := candles[len(candles)-10:]
	
	// Find the swing
	high := recentCandles[0].High
	low := recentCandles[0].Low
	highIdx := 0
	lowIdx := 0
	
	for i, c := range recentCandles {
		if c.High > high {
			high = c.High
			highIdx = i
		}
		if c.Low < low {
			low = c.Low
			lowIdx = i
		}
	}
	
	currentPrice := candles[len(candles)-1].Close
	
	// BULLISH JUDAS: Fake move down, then up
	// Low happens first, then price reverses up
	if lowIdx < highIdx && currentPrice > (high+low)/2 {
		judas := &JudasSwing{
			Type:       "bullish",
			FakeMove:   recentCandles[0].Open - low,
			RealMove:   high - low,
			TrapPrice:  low,
			EntryPrice: (high + low) / 2,
			Strength:   70,
			CandleIdx:  len(candles) - 10 + lowIdx,
		}
		
		// Strength based on reversal size
		if judas.RealMove > judas.FakeMove*1.5 {
			judas.Strength += 20
		}
		
		return judas
	}
	
	// BEARISH JUDAS: Fake move up, then down
	// High happens first, then price reverses down
	if highIdx < lowIdx && currentPrice < (high+low)/2 {
		judas := &JudasSwing{
			Type:       "bearish",
			FakeMove:   high - recentCandles[0].Open,
			RealMove:   high - low,
			TrapPrice:  high,
			EntryPrice: (high + low) / 2,
			Strength:   70,
			CandleIdx:  len(candles) - 10 + highIdx,
		}
		
		if judas.RealMove > judas.FakeMove*1.5 {
			judas.Strength += 20
		}
		
		return judas
	}
	
	return nil
}
