package main

import (
	"math"
)

// TimeframeStrategy defines strategy parameters for each timeframe
type TimeframeStrategy struct {
	Timeframe          string
	MinConfluence      int
	MinRR              float64
	StopLossATR        float64
	TakeProfitATR      []float64
	TrailingStopATR    float64
	RequireSession     bool
	AllowedSessions    []string
	RequireVolume      bool
	VolumeMultiplier   float64
	RequireTrend       bool
	MaxRiskPercent     float64
	RequirePatterns    []string
	RequireICT         bool
	UseSmartMoney      bool
	FilterByVolatility bool
	MinVolatility      float64
	MaxVolatility      float64
}

// GetOptimizedStrategy returns the best strategy for each timeframe
func GetOptimizedStrategy(timeframe string) TimeframeStrategy {
	strategies := map[string]TimeframeStrategy{
		// Scalping - 1m (High frequency, tight stops)
		"1m": {
			Timeframe:          "1m",
			MinConfluence:      5, // Very strict
			MinRR:              2.0,
			StopLossATR:        0.5, // Tight stop
			TakeProfitATR:      []float64{1.0, 1.5, 2.0},
			TrailingStopATR:    0.3,
			RequireSession:     true,
			AllowedSessions:    []string{"London", "NewYork"},
			RequireVolume:      true,
			VolumeMultiplier:   2.0, // High volume required
			RequireTrend:       true,
			MaxRiskPercent:     1.0,
			RequirePatterns:    []string{"Engulfing", "Pin Bar"},
			RequireICT:         true,
			UseSmartMoney:      true,
			FilterByVolatility: true,
			MinVolatility:      0.3,
			MaxVolatility:      2.0,
		},
		
		// Scalping - 3m (Balanced scalping)
		"3m": {
			Timeframe:          "3m",
			MinConfluence:      5,
			MinRR:              2.2,
			StopLossATR:        0.6,
			TakeProfitATR:      []float64{1.3, 2.0, 2.6},
			TrailingStopATR:    0.4,
			RequireSession:     true,
			AllowedSessions:    []string{"London", "NewYork"},
			RequireVolume:      true,
			VolumeMultiplier:   1.8,
			RequireTrend:       true,
			MaxRiskPercent:     1.0,
			RequirePatterns:    []string{"Engulfing", "Order Block"},
			RequireICT:         true,
			UseSmartMoney:      true,
			FilterByVolatility: true,
			MinVolatility:      0.3,
			MaxVolatility:      2.0,
		},
		
		// Scalping - 5m (Most reliable scalping)
		"5m": {
			Timeframe:          "5m",
			MinConfluence:      4,
			MinRR:              2.5,
			StopLossATR:        0.7,
			TakeProfitATR:      []float64{1.75, 2.5, 3.5},
			TrailingStopATR:    0.5,
			RequireSession:     true,
			AllowedSessions:    []string{"London", "NewYork"},
			RequireVolume:      true,
			VolumeMultiplier:   1.5,
			RequireTrend:       true,
			MaxRiskPercent:     1.5,
			RequirePatterns:    []string{"Engulfing", "Order Block", "FVG"},
			RequireICT:         true,
			UseSmartMoney:      true,
			FilterByVolatility: true,
			MinVolatility:      0.2,
			MaxVolatility:      2.5,
		},
		
		// Intraday - 15m (OPTIMIZED - Sweet spot for day trading)
		"15m": {
			Timeframe:          "15m",
			MinConfluence:      5, // Increased from 4 for better filtering
			MinRR:              2.5, // Increased from 2.0 for better RR
			StopLossATR:        1.2, // Increased from 1.0 for less whipsaws
			TakeProfitATR:      []float64{3.0, 4.5, 6.0}, // Better targets
			TrailingStopATR:    1.0, // Increased from 0.8
			RequireSession:     true, // STRICT: Only kill zones
			AllowedSessions:    []string{"London", "NewYork"}, // Removed Asian
			RequireVolume:      true,
			VolumeMultiplier:   1.5, // Increased from 1.3
			RequireTrend:       true, // CHANGED: Require trend
			MaxRiskPercent:     2.0,
			RequirePatterns:    []string{"Order Block", "FVG", "Liquidity Sweep", "Engulfing"},
			RequireICT:         true,
			UseSmartMoney:      true,
			FilterByVolatility: true, // ADDED: Filter volatility
			MinVolatility:      0.2,
			MaxVolatility:      2.5,
		},
		
		// Intraday - 30m (Balanced intraday)
		"30m": {
			Timeframe:          "30m",
			MinConfluence:      4,
			MinRR:              2.2,
			StopLossATR:        1.2,
			TakeProfitATR:      []float64{2.5, 3.5, 5.0},
			TrailingStopATR:    1.0,
			RequireSession:     false,
			AllowedSessions:    []string{},
			RequireVolume:      true,
			VolumeMultiplier:   1.2,
			RequireTrend:       false,
			MaxRiskPercent:     2.0,
			RequirePatterns:    []string{"Order Block", "FVG"},
			RequireICT:         true,
			UseSmartMoney:      true,
			FilterByVolatility: false,
			MinVolatility:      0.0,
			MaxVolatility:      10.0,
		},
		
		// Swing - 1h (OPTIMIZED - High probability swing trades)
		"1h": {
			Timeframe:          "1h",
			MinConfluence:      4, // Increased from 3 for better filtering
			MinRR:              2.8, // Increased from 2.5 for better RR
			StopLossATR:        1.8, // Increased from 1.5 for less whipsaws
			TakeProfitATR:      []float64{5.0, 7.0, 10.0}, // Better targets
			TrailingStopATR:    1.5, // Increased from 1.2
			RequireSession:     true, // ADDED: Session filter
			AllowedSessions:    []string{"London", "NewYork"}, // Kill zones only
			RequireVolume:      true, // ADDED: Volume filter
			VolumeMultiplier:   1.3, // Require above average volume
			RequireTrend:       true, // CHANGED: Require trend
			MaxRiskPercent:     2.0,
			RequirePatterns:    []string{"Order Block", "FVG", "BOS", "Liquidity Sweep"},
			RequireICT:         true,
			UseSmartMoney:      true,
			FilterByVolatility: false,
			MinVolatility:      0.0,
			MaxVolatility:      10.0,
		},
		
		// Swing - 2h (Patient swing trading)
		"2h": {
			Timeframe:          "2h",
			MinConfluence:      3,
			MinRR:              3.0,
			StopLossATR:        1.8,
			TakeProfitATR:      []float64{5.0, 7.0, 10.0},
			TrailingStopATR:    1.5,
			RequireSession:     false,
			AllowedSessions:    []string{},
			RequireVolume:      false,
			VolumeMultiplier:   1.0,
			RequireTrend:       true,
			MaxRiskPercent:     2.0,
			RequirePatterns:    []string{"Order Block", "FVG"},
			RequireICT:         true,
			UseSmartMoney:      true,
			FilterByVolatility: false,
			MinVolatility:      0.0,
			MaxVolatility:      10.0,
		},
		
		// Swing - 4h (Best win rate timeframe)
		"4h": {
			Timeframe:          "4h",
			MinConfluence:      3,
			MinRR:              3.0,
			StopLossATR:        2.0,
			TakeProfitATR:      []float64{6.0, 9.0, 12.0},
			TrailingStopATR:    1.8,
			RequireSession:     false,
			AllowedSessions:    []string{},
			RequireVolume:      false,
			VolumeMultiplier:   1.0,
			RequireTrend:       true,
			MaxRiskPercent:     2.0,
			RequirePatterns:    []string{"Order Block", "FVG", "BOS"},
			RequireICT:         true,
			UseSmartMoney:      true,
			FilterByVolatility: false,
			MinVolatility:      0.0,
			MaxVolatility:      10.0,
		},
		
		// Position - 1d (Long-term high probability)
		"1d": {
			Timeframe:          "1d",
			MinConfluence:      2,
			MinRR:              4.0,
			StopLossATR:        2.5,
			TakeProfitATR:      []float64{10.0, 15.0, 20.0},
			TrailingStopATR:    2.0,
			RequireSession:     false,
			AllowedSessions:    []string{},
			RequireVolume:      false,
			VolumeMultiplier:   1.0,
			RequireTrend:       true,
			MaxRiskPercent:     2.0,
			RequirePatterns:    []string{"Order Block", "BOS"},
			RequireICT:         false,
			UseSmartMoney:      false,
			FilterByVolatility: false,
			MinVolatility:      0.0,
			MaxVolatility:      10.0,
		},
	}
	
	// Return strategy or default to 15m
	if strategy, ok := strategies[timeframe]; ok {
		return strategy
	}
	return strategies["15m"]
}

// OptimizedSignal represents a signal with additional metadata
type OptimizedSignal struct {
	Type       string
	Entry      float64
	StopLoss   float64
	TP1        float64
	TP2        float64
	TP3        float64
	Strength   float64
	Confluence int
	Reasons    []string
	Timeframe  string
	Timestamp  int64
	RR         float64
}

// ApplyTimeframeStrategy applies optimized strategy to signal generation
func ApplyTimeframeStrategy(candles []Candle, strategy TimeframeStrategy) []OptimizedSignal {
	signals := []OptimizedSignal{}
	
	if len(candles) < 100 {
		return signals
	}
	
	for i := 50; i < len(candles)-1; i++ {
		confluence := 0
		reasons := []string{}
		
		// Calculate ATR for dynamic stops
		atr := calculateATR(candles, i)
		
		// Check volatility filter
		if strategy.FilterByVolatility {
			volatility := (candles[i].High - candles[i].Low) / candles[i].Close * 100
			if volatility < strategy.MinVolatility || volatility > strategy.MaxVolatility {
				continue
			}
		}
		
		// Check volume
		if strategy.RequireVolume {
			avgVolume := calculateAverageVolume(candles, i, 20)
			if candles[i].Volume < avgVolume*strategy.VolumeMultiplier {
				continue
			}
			confluence++
			reasons = append(reasons, "High Volume")
		}
		
		// Check trend
		if strategy.RequireTrend {
			trend := detectTrend(candles, i, 20)
			if trend == "none" {
				continue
			}
			confluence++
			reasons = append(reasons, "Strong Trend")
		}
		
		// Check patterns
		for _, pattern := range strategy.RequirePatterns {
			if hasPattern(candles, i, pattern) {
				confluence++
				reasons = append(reasons, pattern)
			}
		}
		
		// Check ICT concepts
		if strategy.RequireICT {
			if detectOrderBlock(candles, i) {
				confluence++
				reasons = append(reasons, "Order Block")
			}
			if detectFVG(candles, i) {
				confluence++
				reasons = append(reasons, "Fair Value Gap")
			}
		}
		
		// Check smart money concepts
		if strategy.UseSmartMoney {
			if detectLiquiditySweep(candles, i) {
				confluence++
				reasons = append(reasons, "Liquidity Sweep")
			}
		}
		
		// Check confluence threshold
		if confluence < strategy.MinConfluence {
			continue
		}
		
		// Determine signal type
		signalType := determineSignalType(candles, i)
		if signalType == "" {
			continue
		}
		
		// Calculate entry, stop loss, and take profits
		entry := candles[i].Close
		var stopLoss, tp1, tp2, tp3 float64
		
		if signalType == "BUY" {
			stopLoss = entry - (atr * strategy.StopLossATR)
			tp1 = entry + (atr * strategy.TakeProfitATR[0])
			tp2 = entry + (atr * strategy.TakeProfitATR[1])
			tp3 = entry + (atr * strategy.TakeProfitATR[2])
		} else {
			stopLoss = entry + (atr * strategy.StopLossATR)
			tp1 = entry - (atr * strategy.TakeProfitATR[0])
			tp2 = entry - (atr * strategy.TakeProfitATR[1])
			tp3 = entry - (atr * strategy.TakeProfitATR[2])
		}
		
		// Calculate RR
		risk := math.Abs(entry - stopLoss)
		reward := math.Abs(entry - tp1)
		rr := reward / risk
		
		// Check minimum RR
		if rr < strategy.MinRR {
			continue
		}
		
		// Create signal
		signal := OptimizedSignal{
			Type:       signalType,
			Entry:      entry,
			StopLoss:   stopLoss,
			TP1:        tp1,
			TP2:        tp2,
			TP3:        tp3,
			Strength:   float64(confluence) * 15.0,
			Confluence: confluence,
			Reasons:    reasons,
			Timeframe:  strategy.Timeframe,
			Timestamp:  candles[i].Timestamp,
			RR:         rr,
		}
		
		signals = append(signals, signal)
	}
	
	return signals
}

// Helper functions
func calculateAverageVolume(candles []Candle, index, period int) float64 {
	if index < period {
		return 0
	}
	sum := 0.0
	for i := index - period; i < index; i++ {
		sum += candles[i].Volume
	}
	return sum / float64(period)
}

func detectTrend(candles []Candle, index, period int) string {
	if index < period {
		return "none"
	}
	
	ema20 := calculateEMA(candles, 20)
	ema50 := calculateEMA(candles, 50)
	
	if ema20 > ema50 && candles[index].Close > ema20 {
		return "bullish"
	} else if ema20 < ema50 && candles[index].Close < ema20 {
		return "bearish"
	}
	return "none"
}

func hasPattern(candles []Candle, index int, pattern string) bool {
	if index < 2 {
		return false
	}
	
	switch pattern {
	case "Engulfing":
		return isBullishEngulfing(candles, index) || isBearishEngulfing(candles, index)
	case "Pin Bar":
		return isPinBar(candles, index)
	case "Order Block":
		return detectOrderBlock(candles, index)
	case "FVG":
		return detectFVG(candles, index)
	case "Liquidity Sweep":
		return detectLiquiditySweep(candles, index)
	case "BOS":
		return hasBreakOfStructure(candles, index)
	}
	return false
}

func determineSignalType(candles []Candle, index int) string {
	if index < 2 {
		return ""
	}
	
	// Check for bullish signals
	if isBullishEngulfing(candles, index) || 
	   (candles[index].Close > candles[index].Open && candles[index-1].Close < candles[index-1].Open) {
		return "BUY"
	}
	
	// Check for bearish signals
	if isBearishEngulfing(candles, index) || 
	   (candles[index].Close < candles[index].Open && candles[index-1].Close > candles[index-1].Open) {
		return "SELL"
	}
	
	return ""
}

func isBullishEngulfing(candles []Candle, index int) bool {
	if index < 1 {
		return false
	}
	prev := candles[index-1]
	curr := candles[index]
	return prev.Close < prev.Open && 
	       curr.Close > curr.Open && 
	       curr.Close > prev.Open && 
	       curr.Open < prev.Close
}

func isBearishEngulfing(candles []Candle, index int) bool {
	if index < 1 {
		return false
	}
	prev := candles[index-1]
	curr := candles[index]
	return prev.Close > prev.Open && 
	       curr.Close < curr.Open && 
	       curr.Close < prev.Open && 
	       curr.Open > prev.Close
}

func isPinBar(candles []Candle, index int) bool {
	c := candles[index]
	body := math.Abs(c.Close - c.Open)
	upperWick := c.High - math.Max(c.Open, c.Close)
	lowerWick := math.Min(c.Open, c.Close) - c.Low
	
	// Bullish pin bar
	if lowerWick > body*2 && upperWick < body {
		return true
	}
	// Bearish pin bar
	if upperWick > body*2 && lowerWick < body {
		return true
	}
	return false
}

func hasBreakOfStructure(candles []Candle, index int) bool {
	if index < 20 {
		return false
	}
	
	// Find recent high/low
	recentHigh := candles[index-20].High
	recentLow := candles[index-20].Low
	
	for i := index - 19; i < index; i++ {
		if candles[i].High > recentHigh {
			recentHigh = candles[i].High
		}
		if candles[i].Low < recentLow {
			recentLow = candles[i].Low
		}
	}
	
	// Check if current candle breaks structure
	return candles[index].Close > recentHigh || candles[index].Close < recentLow
}


// detectOrderBlock detects order block patterns
func detectOrderBlock(candles []Candle, index int) bool {
	if index < 5 {
		return false
	}
	
	// Look for strong move followed by consolidation
	prev := candles[index-1]
	curr := candles[index]
	
	// Bullish order block
	if prev.Close < prev.Open && curr.Close > curr.Open {
		bodySize := curr.Close - curr.Open
		prevBodySize := prev.Open - prev.Close
		if bodySize > prevBodySize*1.5 {
			return true
		}
	}
	
	// Bearish order block
	if prev.Close > prev.Open && curr.Close < curr.Open {
		bodySize := curr.Open - curr.Close
		prevBodySize := prev.Close - prev.Open
		if bodySize > prevBodySize*1.5 {
			return true
		}
	}
	
	return false
}

// detectFVG detects fair value gaps
func detectFVG(candles []Candle, index int) bool {
	if index < 2 {
		return false
	}
	
	// Bullish FVG: gap between candle[i-2].High and candle[i].Low
	if candles[index].Low > candles[index-2].High {
		return true
	}
	
	// Bearish FVG: gap between candles[i-2].Low and candle[i].High
	if candles[index].High < candles[index-2].Low {
		return true
	}
	
	return false
}

// detectLiquiditySweep detects liquidity sweeps
func detectLiquiditySweep(candles []Candle, index int) bool {
	if index < 10 {
		return false
	}
	
	// Find recent high/low
	recentHigh := candles[index-10].High
	recentLow := candles[index-10].Low
	
	for i := index - 9; i < index; i++ {
		if candles[i].High > recentHigh {
			recentHigh = candles[i].High
		}
		if candles[i].Low < recentLow {
			recentLow = candles[i].Low
		}
	}
	
	// Check if current candle swept liquidity and reversed
	curr := candles[index]
	
	// Bullish sweep: wick below recent low, close above
	if curr.Low < recentLow && curr.Close > recentLow {
		return true
	}
	
	// Bearish sweep: wick above recent high, close below
	if curr.High > recentHigh && curr.Close < recentHigh {
		return true
	}
	
	return false
}
