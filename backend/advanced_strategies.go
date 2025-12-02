package main

import (
	"math"
)

// AdvancedStrategy represents a complete trading strategy
type AdvancedStrategy struct {
	Name               string
	Description        string
	Timeframe          string
	MinConfluence      int
	RequiredConcepts   []string
	TargetWinRate      float64
	TargetProfitFactor float64
}

// GetAdvancedStrategies returns all advanced profitable strategies
func GetAdvancedStrategies() map[string]AdvancedStrategy {
	return map[string]AdvancedStrategy{
		"liquidity_hunter": {
			Name:        "Liquidity Hunter",
			Description: "Hunts liquidity sweeps and traps institutional orders",
			Timeframe:   "15m",
			MinConfluence: 6,
			RequiredConcepts: []string{
				"Liquidity Sweep",
				"Order Block",
				"Fair Value Gap",
				"Break of Structure",
				"Volume Spike",
				"Session Alignment",
			},
			TargetWinRate:      75.0,
			TargetProfitFactor: 2.5,
		},
		"smart_money_tracker": {
			Name:        "Smart Money Tracker",
			Description: "Follows institutional money flow and order blocks",
			Timeframe:   "1h",
			MinConfluence: 7,
			RequiredConcepts: []string{
				"Order Block (Institutional)",
				"Fair Value Gap",
				"Liquidity Void",
				"Market Structure Shift",
				"Volume Profile",
				"Delta Analysis",
				"Premium/Discount Zone",
			},
			TargetWinRate:      80.0,
			TargetProfitFactor: 3.0,
		},
		"breakout_master": {
			Name:        "Breakout Master",
			Description: "Catches explosive breakouts with volume confirmation",
			Timeframe:   "15m",
			MinConfluence: 5,
			RequiredConcepts: []string{
				"Break of Structure",
				"Volume Explosion (2x+)",
				"Consolidation Pattern",
				"Support/Resistance Break",
				"Momentum Confirmation",
			},
			TargetWinRate:      70.0,
			TargetProfitFactor: 2.8,
		},
		"trend_rider": {
			Name:        "Trend Rider",
			Description: "Rides strong trends with pullback entries",
			Timeframe:   "4h",
			MinConfluence: 5,
			RequiredConcepts: []string{
				"Strong Trend (EMA alignment)",
				"Pullback to Key Level",
				"Order Block Support",
				"Higher Timeframe Confirmation",
				"Momentum Divergence",
			},
			TargetWinRate:      75.0,
			TargetProfitFactor: 2.5,
		},
		"scalper_pro": {
			Name:        "Scalper Pro",
			Description: "High-frequency scalping with tight risk management",
			Timeframe:   "5m",
			MinConfluence: 6,
			RequiredConcepts: []string{
				"Micro Order Block",
				"Immediate FVG",
				"Volume Spike",
				"Kill Zone Only",
				"Tight Stop (0.5 ATR)",
				"Quick Target (1.5 ATR)",
			},
			TargetWinRate:      65.0,
			TargetProfitFactor: 2.0,
		},
		"reversal_sniper": {
			Name:        "Reversal Sniper",
			Description: "Catches high-probability reversals at key levels",
			Timeframe:   "1h",
			MinConfluence: 7,
			RequiredConcepts: []string{
				"Divergence (RSI/Price)",
				"Order Block at Extreme",
				"Liquidity Sweep",
				"Fair Value Gap",
				"Volume Climax",
				"Candlestick Pattern",
				"Support/Resistance",
			},
			TargetWinRate:      78.0,
			TargetProfitFactor: 3.2,
		},
		"session_trader": {
			Name:        "Session Trader",
			Description: "Exploits session volatility and liquidity",
			Timeframe:   "15m",
			MinConfluence: 6,
			RequiredConcepts: []string{
				"London/NY Session",
				"Session High/Low Sweep",
				"Order Block",
				"Fair Value Gap",
				"Volume Profile",
				"Time-based Entry",
			},
			TargetWinRate:      72.0,
			TargetProfitFactor: 2.4,
		},
		"momentum_beast": {
			Name:        "Momentum Beast",
			Description: "Rides explosive momentum moves with confirmation",
			Timeframe:   "15m",
			MinConfluence: 5,
			RequiredConcepts: []string{
				"Strong Momentum",
				"Volume Confirmation",
				"Break of Structure",
				"No Resistance Above",
				"Trend Alignment",
			},
			TargetWinRate:      68.0,
			TargetProfitFactor: 2.6,
		},
		"range_master": {
			Name:        "Range Master",
			Description: "Trades ranges with high probability",
			Timeframe:   "1h",
			MinConfluence: 6,
			RequiredConcepts: []string{
				"Clear Range Identified",
				"Support/Resistance Bounce",
				"Order Block at Boundary",
				"Volume Decrease in Middle",
				"Rejection Candle",
				"Mean Reversion",
			},
			TargetWinRate:      73.0,
			TargetProfitFactor: 2.2,
		},
		"institutional_follower": {
			Name:        "Institutional Follower",
			Description: "Follows big money institutional orders",
			Timeframe:   "4h",
			MinConfluence: 8,
			RequiredConcepts: []string{
				"Institutional Order Block",
				"Large Volume Spike",
				"Fair Value Gap",
				"Market Structure Shift",
				"Premium/Discount Zone",
				"Liquidity Grab",
				"Trend Confirmation",
				"Higher TF Alignment",
			},
			TargetWinRate:      82.0,
			TargetProfitFactor: 3.5,
		},
	}
}

// GenerateSignalWithStrategy generates signal using specific strategy
func GenerateSignalWithStrategy(candles []Candle, strategyName string) *AdvancedSignal {
	strategies := GetAdvancedStrategies()
	strategy, exists := strategies[strategyName]
	if !exists {
		return nil
	}
	
	if len(candles) < 100 {
		return nil
	}
	
	idx := len(candles) - 1
	confluence := 0
	reasons := []string{}
	
	// Check each required concept
	for _, concept := range strategy.RequiredConcepts {
		if checkConcept(candles, idx, concept) {
			confluence++
			reasons = append(reasons, concept)
		}
	}
	
	// Must meet minimum confluence
	if confluence < strategy.MinConfluence {
		return nil
	}
	
	// Determine signal type
	signalType := determineSignalTypeAdvanced(candles, idx, strategyName)
	if signalType == "" {
		return nil
	}
	
	// Calculate entry, stops, targets
	atr := calculateATR(candles, idx)
	entry := candles[idx].Close
	
	var stopLoss, tp1, tp2, tp3 float64
	var stopATR, tp1ATR, tp2ATR, tp3ATR float64
	
	// Strategy-specific risk/reward
	switch strategyName {
	case "scalper_pro":
		stopATR, tp1ATR, tp2ATR, tp3ATR = 0.5, 1.5, 2.5, 3.5
	case "liquidity_hunter":
		stopATR, tp1ATR, tp2ATR, tp3ATR = 1.0, 3.0, 5.0, 7.0
	case "smart_money_tracker", "institutional_follower":
		stopATR, tp1ATR, tp2ATR, tp3ATR = 1.5, 5.0, 8.0, 12.0
	case "reversal_sniper":
		stopATR, tp1ATR, tp2ATR, tp3ATR = 1.2, 4.0, 6.0, 9.0
	case "breakout_master", "momentum_beast":
		stopATR, tp1ATR, tp2ATR, tp3ATR = 1.0, 3.5, 6.0, 9.0
	default:
		stopATR, tp1ATR, tp2ATR, tp3ATR = 1.5, 4.0, 6.0, 8.0
	}
	
	if signalType == "BUY" {
		stopLoss = entry - (atr * stopATR)
		tp1 = entry + (atr * tp1ATR)
		tp2 = entry + (atr * tp2ATR)
		tp3 = entry + (atr * tp3ATR)
	} else {
		stopLoss = entry + (atr * stopATR)
		tp1 = entry - (atr * tp1ATR)
		tp2 = entry - (atr * tp2ATR)
		tp3 = entry - (atr * tp3ATR)
	}
	
	// Calculate RR
	risk := math.Abs(entry - stopLoss)
	reward := math.Abs(entry - tp1)
	rr := reward / risk
	
	// Minimum RR based on strategy
	minRR := 2.0
	if strategyName == "smart_money_tracker" || strategyName == "institutional_follower" {
		minRR = 3.0
	}
	
	if rr < minRR {
		return nil
	}
	
	signal := &AdvancedSignal{
		Strategy:   strategyName,
		Type:       signalType,
		Entry:      entry,
		StopLoss:   stopLoss,
		TP1:        tp1,
		TP2:        tp2,
		TP3:        tp3,
		Confluence: confluence,
		Reasons:    reasons,
		Strength:   float64(confluence) * 12.0,
		RR:         rr,
		Timeframe:  strategy.Timeframe,
	}
	
	return signal
}

// AdvancedSignal represents a signal from advanced strategy
type AdvancedSignal struct {
	Strategy   string
	Type       string
	Entry      float64
	StopLoss   float64
	TP1        float64
	TP2        float64
	TP3        float64
	Confluence int
	Reasons    []string
	Strength   float64
	RR         float64
	Timeframe  string
}

// checkConcept checks if a specific concept is present
func checkConcept(candles []Candle, idx int, concept string) bool {
	if idx < 20 {
		return false
	}
	
	switch concept {
	case "Liquidity Sweep":
		return detectLiquiditySweep(candles, idx)
	case "Order Block", "Order Block (Institutional)", "Micro Order Block", "Order Block at Extreme", "Order Block Support", "Order Block at Boundary":
		return detectOrderBlock(candles, idx)
	case "Fair Value Gap", "Immediate FVG":
		return detectFVG(candles, idx)
	case "Break of Structure", "Market Structure Shift":
		return hasBreakOfStructure(candles, idx)
	case "Volume Spike", "Volume Explosion (2x+)", "Large Volume Spike":
		return hasVolumeSpike(candles, idx, 2.0)
	case "Session Alignment", "Kill Zone Only", "London/NY Session":
		return isKillZone(candles[idx].Timestamp)
	case "Strong Trend (EMA alignment)", "Trend Alignment", "Trend Confirmation":
		return hasStrongTrend(candles, idx)
	case "Pullback to Key Level":
		return isPullbackToKeyLevel(candles, idx)
	case "Higher Timeframe Confirmation", "Higher TF Alignment":
		return true // Simplified - would need HTF data
	case "Momentum Divergence", "Divergence (RSI/Price)":
		return hasDivergence(candles, idx)
	case "Volume Profile", "Delta Analysis":
		return hasVolumeConfirmation(candles, idx)
	case "Premium/Discount Zone":
		return isInPremiumDiscountZone(candles, idx)
	case "Consolidation Pattern", "Clear Range Identified":
		return hasConsolidation(candles, idx)
	case "Support/Resistance Break", "Support/Resistance Bounce", "Support/Resistance":
		return isAtSupportResistance(candles, idx)
	case "Momentum Confirmation", "Strong Momentum":
		return hasStrongMomentum(candles, idx)
	case "Volume Climax":
		return hasVolumeClimax(candles, idx)
	case "Candlestick Pattern", "Rejection Candle":
		return hasSignificantPattern(candles, idx)
	case "Session High/Low Sweep":
		return hasSessionSweep(candles, idx)
	case "Time-based Entry":
		return isOptimalEntryTime(candles[idx].Timestamp)
	case "No Resistance Above":
		return hasNoResistanceAbove(candles, idx)
	case "Volume Decrease in Middle":
		return hasLowVolume(candles, idx)
	case "Mean Reversion":
		return isMeanReversion(candles, idx)
	case "Liquidity Grab", "Liquidity Void":
		return detectLiquiditySweep(candles, idx)
	default:
		return false
	}
}

// Helper functions for concept detection
func hasVolumeSpike(candles []Candle, idx int, multiplier float64) bool {
	if idx < 20 {
		return false
	}
	avgVol := calculateAverageVolume(candles, idx, 20)
	return candles[idx].Volume > avgVol*multiplier
}

func isKillZone(timestamp int64) bool {
	// Simplified - would need proper timezone handling
	return true // Assume always in kill zone for now
}

func hasStrongTrend(candles []Candle, idx int) bool {
	if idx < 50 {
		return false
	}
	ema20 := calculateEMA(candles[:idx+1], 20)
	ema50 := calculateEMA(candles[:idx+1], 50)
	return math.Abs(ema20-ema50) > ema50*0.01
}

func isPullbackToKeyLevel(candles []Candle, idx int) bool {
	if idx < 20 {
		return false
	}
	// Check if price pulled back to EMA20
	ema20 := calculateEMA(candles[:idx+1], 20)
	return math.Abs(candles[idx].Close-ema20) < ema20*0.005
}

func hasDivergence(candles []Candle, idx int) bool {
	// Simplified divergence detection
	if idx < 30 {
		return false
	}
	// Price making lower low but momentum increasing
	return candles[idx].Low < candles[idx-10].Low && candles[idx].Volume > candles[idx-10].Volume
}

func hasVolumeConfirmation(candles []Candle, idx int) bool {
	if idx < 10 {
		return false
	}
	avgVol := calculateAverageVolume(candles, idx, 10)
	return candles[idx].Volume > avgVol*1.3
}

func isInPremiumDiscountZone(candles []Candle, idx int) bool {
	if idx < 50 {
		return false
	}
	high := candles[idx-50].High
	low := candles[idx-50].Low
	for i := idx - 49; i < idx; i++ {
		if candles[i].High > high {
			high = candles[i].High
		}
		if candles[i].Low < low {
			low = candles[i].Low
		}
	}
	mid := (high + low) / 2
	// Premium if above mid, discount if below
	return math.Abs(candles[idx].Close-mid) > (high-low)*0.2
}

func hasConsolidation(candles []Candle, idx int) bool {
	if idx < 20 {
		return false
	}
	// Check if range is tight
	high := candles[idx-20].High
	low := candles[idx-20].Low
	for i := idx - 19; i < idx; i++ {
		if candles[i].High > high {
			high = candles[i].High
		}
		if candles[i].Low < low {
			low = candles[i].Low
		}
	}
	rangeSize := (high - low) / candles[idx].Close
	return rangeSize < 0.02 // 2% range
}

func isAtSupportResistance(candles []Candle, idx int) bool {
	if idx < 50 {
		return false
	}
	// Find recent highs/lows
	for i := idx - 50; i < idx-5; i++ {
		if math.Abs(candles[idx].Close-candles[i].High) < candles[idx].Close*0.005 {
			return true
		}
		if math.Abs(candles[idx].Close-candles[i].Low) < candles[idx].Close*0.005 {
			return true
		}
	}
	return false
}

func hasStrongMomentum(candles []Candle, idx int) bool {
	if idx < 5 {
		return false
	}
	// Check if last 5 candles are in same direction
	bullish := 0
	for i := idx - 4; i <= idx; i++ {
		if candles[i].Close > candles[i].Open {
			bullish++
		}
	}
	return bullish >= 4 || bullish <= 1
}

func hasVolumeClimax(candles []Candle, idx int) bool {
	if idx < 20 {
		return false
	}
	avgVol := calculateAverageVolume(candles, idx, 20)
	return candles[idx].Volume > avgVol*3.0
}

func hasSignificantPattern(candles []Candle, idx int) bool {
	return isBullishEngulfing(candles, idx) || isBearishEngulfing(candles, idx) || isPinBar(candles, idx)
}

func hasSessionSweep(candles []Candle, idx int) bool {
	// Simplified - check if swept recent high/low
	return detectLiquiditySweep(candles, idx)
}

func isOptimalEntryTime(timestamp int64) bool {
	// Simplified - would need proper time analysis
	return true
}

func hasNoResistanceAbove(candles []Candle, idx int) bool {
	if idx < 50 {
		return false
	}
	currentPrice := candles[idx].Close
	for i := idx - 50; i < idx; i++ {
		if candles[i].High > currentPrice*1.02 {
			return false
		}
	}
	return true
}

func hasLowVolume(candles []Candle, idx int) bool {
	if idx < 20 {
		return false
	}
	avgVol := calculateAverageVolume(candles, idx, 20)
	return candles[idx].Volume < avgVol*0.7
}

func isMeanReversion(candles []Candle, idx int) bool {
	if idx < 20 {
		return false
	}
	ema20 := calculateEMA(candles[:idx+1], 20)
	deviation := math.Abs(candles[idx].Close-ema20) / ema20
	return deviation > 0.02 // 2% deviation
}

func determineSignalTypeAdvanced(candles []Candle, idx int, strategy string) string {
	if idx < 2 {
		return ""
	}
	
	// Strategy-specific signal determination
	switch strategy {
	case "reversal_sniper":
		// Look for reversal patterns
		if isBullishEngulfing(candles, idx) && candles[idx].Low < candles[idx-10].Low {
			return "BUY"
		}
		if isBearishEngulfing(candles, idx) && candles[idx].High > candles[idx-10].High {
			return "SELL"
		}
	case "breakout_master", "momentum_beast":
		// Look for breakouts
		if hasBreakOfStructure(candles, idx) && candles[idx].Close > candles[idx].Open {
			return "BUY"
		}
		if hasBreakOfStructure(candles, idx) && candles[idx].Close < candles[idx].Open {
			return "SELL"
		}
	default:
		// Standard determination
		if candles[idx].Close > candles[idx].Open && candles[idx].Close > candles[idx-1].Close {
			return "BUY"
		}
		if candles[idx].Close < candles[idx].Open && candles[idx].Close < candles[idx-1].Close {
			return "SELL"
		}
	}
	
	return ""
}
