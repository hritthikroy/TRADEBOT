package strategies

// MASTER STRATEGY - Combines ALL concepts for maximum profitability
// Uses: ICT, Liquidity, Order Flow, Supply/Demand, Power of 3, Market Maker Model

func generateMasterSignal(data []Candle, interval string) *Signal {
	if len(data) < 100 {
		return nil
	}

	// PHASE 1: BASIC INDICATORS
	ema9 := calculateEMA(data, 9)
	ema20 := calculateEMA(data, 20)
	ema50 := calculateEMA(data, 50)
	rsi := calculateRSI(data, 14)
	atr := calculateATR(data[len(data)-14:], 14)
	
	if atr == 0 {
		return nil
	}

	currentPrice := data[len(data)-1].Close
	
	// PHASE 2: TREND FILTER
	bullTrend := ema9 > ema20 && ema20 > ema50
	bearTrend := ema9 < ema20 && ema20 < ema50
	
	if !bullTrend && !bearTrend {
		return nil
	}

	// PHASE 3: ICT/SMC ANALYSIS
	orderBlocks := FindOrderBlocks(data[len(data)-30:])
	fvgs := FindFairValueGaps(data[len(data)-30:])
	
	hasOrderBlock := len(orderBlocks) > 0
	hasFVG := len(fvgs) > 0

	// PHASE 4: LIQUIDITY ANALYSIS
	liquiditySweep := detectLiquiditySweepMaster(data)
	
	// PHASE 5: ORDER FLOW
	orderFlow := PerformOrderFlowAnalysis(data[len(data)-50:])
	hasOrderFlowConfirmation := false
	if orderFlow != nil {
		if bullTrend && orderFlow.Delta.Trend == "bullish" {
			hasOrderFlowConfirmation = true
		}
		if bearTrend && orderFlow.Delta.Trend == "bearish" {
			hasOrderFlowConfirmation = true
		}
	}

	// PHASE 6: SUPPLY/DEMAND ZONES
	sdZones := FindSupplyDemandZones(data[len(data)-100:])
	nearSDZone := false
	if sdZones != nil {
		for _, zone := range sdZones.DemandZones {
			if bullTrend && currentPrice >= zone.Low && currentPrice <= zone.High {
				nearSDZone = true
				break
			}
		}
		for _, zone := range sdZones.SupplyZones {
			if bearTrend && currentPrice >= zone.Low && currentPrice <= zone.High {
				nearSDZone = true
				break
			}
		}
	}

	// PHASE 7: CANDLESTICK PATTERNS
	patterns := RecognizeAllPatterns(data[len(data)-10:], interval)
	hasPattern := false
	if len(patterns) > 0 {
		for _, p := range patterns {
			if bullTrend && p.Type == "bullish" && p.Strength >= 75 {
				hasPattern = true
				break
			}
			if bearTrend && p.Type == "bearish" && p.Strength >= 75 {
				hasPattern = true
				break
			}
		}
	}

	// PHASE 8: MARKET STRUCTURE
	structure := AnalyzeMarketStructure(data)
	structureConfirmation := false
	if structure.BOS {
		structureConfirmation = true
	}

	// PHASE 9: CONFLUENCE SCORING
	confluenceScore := 0
	if hasOrderBlock {
		confluenceScore++
	}
	if hasFVG {
		confluenceScore++
	}
	if liquiditySweep != nil && liquiditySweep.Confirmed {
		confluenceScore++
	}
	if hasOrderFlowConfirmation {
		confluenceScore++
	}
	if nearSDZone {
		confluenceScore++
	}
	if hasPattern {
		confluenceScore++
	}
	if structureConfirmation {
		confluenceScore++
	}

	// REQUIRE MINIMUM 3 CONFLUENCES
	if confluenceScore < 3 {
		return nil
	}

	// PHASE 10: MOMENTUM CONFIRMATION
	if bullTrend && rsi < 48 {
		return nil
	}
	if bearTrend && rsi > 52 {
		return nil
	}

	// PHASE 11: VOLUME CONFIRMATION
	avgVolume := 0.0
	for i := len(data) - 20; i < len(data); i++ {
		avgVolume += data[i].Volume
	}
	avgVolume /= 20
	
	if data[len(data)-1].Volume < avgVolume*1.1 {
		return nil
	}

	// PHASE 12: GENERATE HIGH REWARD SIGNAL
	var signalType string
	var entry, stopLoss float64
	var targets []Target
	
	// Timeframe-specific R:R
	var rr1, rr2, rr3 float64
	switch interval {
	case "15m", "5m":
		rr1, rr2, rr3 = 3.0, 5.0, 7.0
	case "1h", "30m":
		rr1, rr2, rr3 = 4.0, 6.0, 8.0
	case "2h":
		rr1, rr2, rr3 = 5.0, 7.0, 9.0
	case "4h", "1d":
		rr1, rr2, rr3 = 6.0, 9.0, 12.0
	default:
		rr1, rr2, rr3 = 4.0, 6.0, 8.0
	}
	
	if bullTrend {
		signalType = "BUY"
		entry = currentPrice
		
		// Smart stop placement
		stopLoss = ema50 - atr*0.5
		if liquiditySweep != nil && liquiditySweep.Type == "sellside" {
			stopLoss = liquiditySweep.SweepPrice - atr*0.3
		}
		
		riskAmount := entry - stopLoss
		targets = []Target{
			{Price: entry + riskAmount*rr1, RR: rr1, Percentage: 40},
			{Price: entry + riskAmount*rr2, RR: rr2, Percentage: 35},
			{Price: entry + riskAmount*rr3, RR: rr3, Percentage: 25},
		}
		
	} else {
		signalType = "SELL"
		entry = currentPrice
		
		stopLoss = ema50 + atr*0.5
		if liquiditySweep != nil && liquiditySweep.Type == "buyside" {
			stopLoss = liquiditySweep.SweepPrice + atr*0.3
		}
		
		riskAmount := stopLoss - entry
		targets = []Target{
			{Price: entry - riskAmount*rr1, RR: rr1, Percentage: 40},
			{Price: entry - riskAmount*rr2, RR: rr2, Percentage: 35},
			{Price: entry - riskAmount*rr3, RR: rr3, Percentage: 25},
		}
	}

	// Calculate final confidence
	confidence := 60.0 + float64(confluenceScore)*5.0

	return &Signal{
		Type:      signalType,
		Entry:     entry,
		StopLoss:  stopLoss,
		Targets:   targets,
		Strength:  confidence,
		Timeframe: interval,
	}
}

// detectLiquiditySweepMaster - Enhanced liquidity sweep detection
func detectLiquiditySweepMaster(data []Candle) *LiquiditySweepResult {
	if len(data) < 20 {
		return nil
	}

	// Find swing high/low
	swingHigh := 0.0
	swingLow := data[len(data)-15].Low
	
	for i := len(data) - 20; i < len(data)-5; i++ {
		if data[i].High > swingHigh {
			swingHigh = data[i].High
		}
		if data[i].Low < swingLow {
			swingLow = data[i].Low
		}
	}

	// Check for sweep in recent candles
	for i := len(data) - 5; i < len(data); i++ {
		c := data[i]
		
		// Sellside sweep
		if c.Low < swingLow && c.Close > swingLow {
			return &LiquiditySweepResult{
				Type:        "sellside",
				SweepPrice:  c.Low,
				Confirmed:   c.Close > c.Open,
			}
		}
		
		// Buyside sweep
		if c.High > swingHigh && c.Close < swingHigh {
			return &LiquiditySweepResult{
				Type:        "buyside",
				SweepPrice:  c.High,
				Confirmed:   c.Close < c.Open,
			}
		}
	}

	return nil
}
