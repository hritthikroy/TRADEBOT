package ict

import (
	"fmt"
	"log"
	"math"
	"time"
)

// ==================== LIQUIDITY-FIRST UNIFIED STRATEGY ====================
// The Ultimate Truth: Liquidity + Market Structure + Supply/Demand
// Combines: SMC/ICT + Auction Market Theory + Wyckoff
// Core Philosophy: "Wait for liquidity grab, enter on reversal confirmation"

// LiquidityFirstStrategy is the unified institutional strategy
type LiquidityFirstStrategy struct {
	Symbol            string
	Timeframes        []string
	CheckInterval     time.Duration
	MinConfidence     float64
	MinRiskReward     float64
	MaxDailyTrades    int
	TradesToday       int
	IsRunning         bool
}

// UnifiedSetup represents the complete liquidity-first setup
type UnifiedSetup struct {
	// Setup Info
	SetupName       string
	SetupType       string  // "liquidity_sweep", "spring", "value_rejection"
	Timeframe       string
	Timestamp       time.Time
	
	// Core Analysis (The 3 Pillars)
	Liquidity       *LiquidityAnalysisResult
	Structure       *MarketStructureResult
	SupplyDemand    *SupplyDemandResult
	
	// Supporting Analysis
	AuctionTheory   *AuctionTheoryResult
	Wyckoff         *WyckoffPhaseResult
	
	// Trade Parameters
	Direction       string
	Confidence      float64
	Entry           float64
	StopLoss        float64
	TP1             float64
	TP2             float64
	TP3             float64
	RiskReward      float64
	
	// Validation
	IsValid         bool
	Reason          string
}

// LiquidityAnalysisResult holds liquidity analysis
type LiquidityAnalysisResult struct {
	// Liquidity Pools
	BuysideLiquidity  []float64 // Above price (stop losses of shorts)
	SellsideLiquidity []float64 // Below price (stop losses of longs)
	
	// Recent Sweeps
	RecentSweep       bool
	SweepType         string  // "buyside" or "sellside"
	SweepPrice        float64
	SweepConfirmed    bool
	
	// Liquidity Targets
	NearestBuyside    float64
	NearestSellside   float64
	
	// Score
	LiquidityScore    float64
	Direction         string
}

// MarketStructureResult holds market structure analysis
type MarketStructureResult struct {
	// Current Structure
	Trend             string  // "bullish", "bearish", "ranging"
	LastSwingHigh     float64
	LastSwingLow      float64
	
	// Structure Breaks
	BOS               bool    // Break of Structure
	CHOCH             bool    // Change of Character
	MSS               bool    // Market Structure Shift
	
	// Key Levels
	StructureHighs    []float64
	StructureLows     []float64
	
	// Score
	StructureScore    float64
	Direction         string
}

// SupplyDemandResult holds supply/demand analysis
type SupplyDemandResult struct {
	// Zones
	FreshDemandZones  []SDZone
	FreshSupplyZones  []SDZone
	
	// Current Position
	InDemandZone      bool
	InSupplyZone      bool
	NearestDemand     *SDZone
	NearestSupply     *SDZone
	
	// Score
	SDScore           float64
	Direction         string
}

// SDZone represents a supply/demand zone
type SDZone struct {
	High      float64
	Low       float64
	MidPoint  float64
	Strength  float64
	Fresh     bool
	Touches   int
}

// AuctionTheoryResult holds auction market theory analysis
type AuctionTheoryResult struct {
	// Value Area
	POC               float64 // Point of Control
	VAH               float64 // Value Area High
	VAL               float64 // Value Area Low
	
	// Price Position
	AboveValue        bool
	BelowValue        bool
	InValue           bool
	
	// Market Type
	MarketType        string  // "balance", "imbalance", "trend"
	
	// Score
	AuctionScore      float64
	Direction         string
}

// WyckoffPhaseResult holds Wyckoff phase analysis
type WyckoffPhaseResult struct {
	// Current Phase
	Phase             string  // "accumulation", "markup", "distribution", "markdown"
	SubPhase          string  // "spring", "test", "sos", "lpsy", etc.
	
	// Key Events
	SpringDetected    bool
	UpthrustDetected  bool
	SOSDetected       bool
	
	// Score
	WyckoffScore      float64
	Direction         string
}


// NewLiquidityFirstStrategy creates the unified strategy
func NewLiquidityFirstStrategy() *LiquidityFirstStrategy {
	return &LiquidityFirstStrategy{
		Symbol:         "BTCUSDT",
		Timeframes:     []string{"15m", "1h", "4h"},
		CheckInterval:  60 * time.Second,
		MinConfidence:  80.0,
		MinRiskReward:  3.0,
		MaxDailyTrades: 3,
		TradesToday:    0,
		IsRunning:      false,
	}
}

// ==================== CORE ANALYSIS: LIQUIDITY ====================

// AnalyzeLiquidity performs comprehensive liquidity analysis
func AnalyzeLiquidity(candles []Candle) *LiquidityAnalysisResult {
	if len(candles) < 50 {
		return nil
	}
	
	result := &LiquidityAnalysisResult{
		BuysideLiquidity:  []float64{},
		SellsideLiquidity: []float64{},
		Direction:         "neutral",
	}
	
	currentPrice := candles[len(candles)-1].Close
	
	// Find all swing highs (buyside liquidity - shorts' stops)
	// Find all swing lows (sellside liquidity - longs' stops)
	for i := 2; i < len(candles)-2; i++ {
		// Swing High
		if candles[i].High > candles[i-1].High &&
			candles[i].High > candles[i-2].High &&
			candles[i].High > candles[i+1].High &&
			candles[i].High > candles[i+2].High {
			result.BuysideLiquidity = append(result.BuysideLiquidity, candles[i].High)
		}
		
		// Swing Low
		if candles[i].Low < candles[i-1].Low &&
			candles[i].Low < candles[i-2].Low &&
			candles[i].Low < candles[i+1].Low &&
			candles[i].Low < candles[i+2].Low {
			result.SellsideLiquidity = append(result.SellsideLiquidity, candles[i].Low)
		}
	}
	
	// Find equal highs/lows (STRONG liquidity)
	tolerance := currentPrice * 0.001 // 0.1%
	for i := 0; i < len(candles)-10; i++ {
		for j := i + 5; j < len(candles); j++ {
			// Equal highs
			if math.Abs(candles[i].High-candles[j].High) < tolerance {
				result.BuysideLiquidity = append(result.BuysideLiquidity, (candles[i].High+candles[j].High)/2)
			}
			// Equal lows
			if math.Abs(candles[i].Low-candles[j].Low) < tolerance {
				result.SellsideLiquidity = append(result.SellsideLiquidity, (candles[i].Low+candles[j].Low)/2)
			}
		}
	}
	
	// Find nearest liquidity levels
	for _, level := range result.BuysideLiquidity {
		if level > currentPrice {
			if result.NearestBuyside == 0 || level < result.NearestBuyside {
				result.NearestBuyside = level
			}
		}
	}
	
	for _, level := range result.SellsideLiquidity {
		if level < currentPrice {
			if result.NearestSellside == 0 || level > result.NearestSellside {
				result.NearestSellside = level
			}
		}
	}
	
	// Detect recent liquidity sweep
	recentCandles := candles[len(candles)-10:]
	for i := 1; i < len(recentCandles); i++ {
		c := recentCandles[i]
		_ = recentCandles[i-1] // prev unused
		
		// Buyside sweep (wick above, close below)
		for _, level := range result.BuysideLiquidity {
			if c.High > level && c.Close < level && c.Close < c.Open {
				result.RecentSweep = true
				result.SweepType = "buyside"
				result.SweepPrice = c.High
				
				// Confirm if next candle continues down
				if i < len(recentCandles)-1 && recentCandles[i+1].Close < c.Close {
					result.SweepConfirmed = true
				}
			}
		}
		
		// Sellside sweep (wick below, close above)
		for _, level := range result.SellsideLiquidity {
			if c.Low < level && c.Close > level && c.Close > c.Open {
				result.RecentSweep = true
				result.SweepType = "sellside"
				result.SweepPrice = c.Low
				
				if i < len(recentCandles)-1 && recentCandles[i+1].Close > c.Close {
					result.SweepConfirmed = true
				}
			}
		}
	}
	
	// Calculate score and direction
	if result.RecentSweep && result.SweepConfirmed {
		result.LiquidityScore = 90
		if result.SweepType == "sellside" {
			result.Direction = "bullish" // Swept lows = bullish
		} else {
			result.Direction = "bearish" // Swept highs = bearish
		}
	} else if result.RecentSweep {
		result.LiquidityScore = 70
		if result.SweepType == "sellside" {
			result.Direction = "bullish"
		} else {
			result.Direction = "bearish"
		}
	} else {
		result.LiquidityScore = 50
	}
	
	return result
}

// ==================== CORE ANALYSIS: MARKET STRUCTURE ====================

// AnalyzeMarketStructureFull performs comprehensive structure analysis
func AnalyzeMarketStructureFull(candles []Candle) *MarketStructureResult {
	if len(candles) < 50 {
		return nil
	}
	
	result := &MarketStructureResult{
		Trend:          "ranging",
		StructureHighs: []float64{},
		StructureLows:  []float64{},
		Direction:      "neutral",
	}
	
	// Find swing points
	for i := 2; i < len(candles)-2; i++ {
		if candles[i].High > candles[i-1].High &&
			candles[i].High > candles[i-2].High &&
			candles[i].High > candles[i+1].High &&
			candles[i].High > candles[i+2].High {
			result.StructureHighs = append(result.StructureHighs, candles[i].High)
		}
		
		if candles[i].Low < candles[i-1].Low &&
			candles[i].Low < candles[i-2].Low &&
			candles[i].Low < candles[i+1].Low &&
			candles[i].Low < candles[i+2].Low {
			result.StructureLows = append(result.StructureLows, candles[i].Low)
		}
	}
	
	if len(result.StructureHighs) < 2 || len(result.StructureLows) < 2 {
		return result
	}
	
	// Set last swing points
	result.LastSwingHigh = result.StructureHighs[len(result.StructureHighs)-1]
	result.LastSwingLow = result.StructureLows[len(result.StructureLows)-1]
	
	// Determine trend (HH/HL = bullish, LH/LL = bearish)
	lastHH := result.StructureHighs[len(result.StructureHighs)-1]
	prevHH := result.StructureHighs[len(result.StructureHighs)-2]
	lastLL := result.StructureLows[len(result.StructureLows)-1]
	prevLL := result.StructureLows[len(result.StructureLows)-2]
	
	if lastHH > prevHH && lastLL > prevLL {
		result.Trend = "bullish"
	} else if lastHH < prevHH && lastLL < prevLL {
		result.Trend = "bearish"
	}
	
	// Detect BOS (Break of Structure)
	currentPrice := candles[len(candles)-1].Close
	
	if result.Trend == "bullish" && currentPrice > result.LastSwingHigh {
		result.BOS = true
	}
	if result.Trend == "bearish" && currentPrice < result.LastSwingLow {
		result.BOS = true
	}
	
	// Detect CHOCH (Change of Character) - first break against trend
	if result.Trend == "bullish" && currentPrice < result.LastSwingLow {
		result.CHOCH = true
		result.MSS = true
	}
	if result.Trend == "bearish" && currentPrice > result.LastSwingHigh {
		result.CHOCH = true
		result.MSS = true
	}
	
	// Calculate score
	if result.CHOCH {
		result.StructureScore = 85
		// CHOCH = reversal signal
		if result.Trend == "bullish" {
			result.Direction = "bearish"
		} else {
			result.Direction = "bullish"
		}
	} else if result.BOS {
		result.StructureScore = 75
		result.Direction = result.Trend
	} else {
		result.StructureScore = 60
		result.Direction = result.Trend
	}
	
	return result
}


// ==================== CORE ANALYSIS: SUPPLY & DEMAND ====================

// AnalyzeSupplyDemandFull performs comprehensive S/D analysis
func AnalyzeSupplyDemandFull(candles []Candle) *SupplyDemandResult {
	if len(candles) < 50 {
		return nil
	}
	
	result := &SupplyDemandResult{
		FreshDemandZones: []SDZone{},
		FreshSupplyZones: []SDZone{},
		Direction:        "neutral",
	}
	
	currentPrice := candles[len(candles)-1].Close
	
	// Find demand zones (base before rally)
	for i := 5; i < len(candles)-5; i++ {
		// Look for consolidation followed by strong up move
		baseHigh := candles[i].High
		baseLow := candles[i].Low
		
		// Check for base (small range)
		for j := i - 2; j <= i; j++ {
			if candles[j].High > baseHigh {
				baseHigh = candles[j].High
			}
			if candles[j].Low < baseLow {
				baseLow = candles[j].Low
			}
		}
		
		baseSize := baseHigh - baseLow
		
		// Check for rally after base
		rallySize := 0.0
		for j := i + 1; j < i+5 && j < len(candles); j++ {
			if candles[j].Close > candles[j].Open {
				rallySize += candles[j].Close - candles[j].Open
			}
		}
		
		if rallySize > baseSize*2 {
			// Check if zone is fresh
			fresh := true
			touches := 0
			for j := i + 5; j < len(candles); j++ {
				if candles[j].Low <= baseHigh && candles[j].Low >= baseLow {
					touches++
					if touches > 0 {
						fresh = false
					}
				}
			}
			
			zone := SDZone{
				High:     baseHigh,
				Low:      baseLow,
				MidPoint: (baseHigh + baseLow) / 2,
				Strength: 70 + rallySize/baseSize*10,
				Fresh:    fresh,
				Touches:  touches,
			}
			
			if fresh {
				zone.Strength += 20
			}
			
			result.FreshDemandZones = append(result.FreshDemandZones, zone)
		}
	}
	
	// Find supply zones (base before drop)
	for i := 5; i < len(candles)-5; i++ {
		baseHigh := candles[i].High
		baseLow := candles[i].Low
		
		for j := i - 2; j <= i; j++ {
			if candles[j].High > baseHigh {
				baseHigh = candles[j].High
			}
			if candles[j].Low < baseLow {
				baseLow = candles[j].Low
			}
		}
		
		baseSize := baseHigh - baseLow
		
		// Check for drop after base
		dropSize := 0.0
		for j := i + 1; j < i+5 && j < len(candles); j++ {
			if candles[j].Close < candles[j].Open {
				dropSize += candles[j].Open - candles[j].Close
			}
		}
		
		if dropSize > baseSize*2 {
			fresh := true
			touches := 0
			for j := i + 5; j < len(candles); j++ {
				if candles[j].High >= baseLow && candles[j].High <= baseHigh {
					touches++
					if touches > 0 {
						fresh = false
					}
				}
			}
			
			zone := SDZone{
				High:     baseHigh,
				Low:      baseLow,
				MidPoint: (baseHigh + baseLow) / 2,
				Strength: 70 + dropSize/baseSize*10,
				Fresh:    fresh,
				Touches:  touches,
			}
			
			if fresh {
				zone.Strength += 20
			}
			
			result.FreshSupplyZones = append(result.FreshSupplyZones, zone)
		}
	}
	
	// Check current position
	for i := range result.FreshDemandZones {
		zone := &result.FreshDemandZones[i]
		if currentPrice >= zone.Low && currentPrice <= zone.High {
			result.InDemandZone = true
		}
		if zone.Fresh && (result.NearestDemand == nil || zone.High > result.NearestDemand.High) {
			if zone.High < currentPrice {
				result.NearestDemand = zone
			}
		}
	}
	
	for i := range result.FreshSupplyZones {
		zone := &result.FreshSupplyZones[i]
		if currentPrice >= zone.Low && currentPrice <= zone.High {
			result.InSupplyZone = true
		}
		if zone.Fresh && (result.NearestSupply == nil || zone.Low < result.NearestSupply.Low) {
			if zone.Low > currentPrice {
				result.NearestSupply = zone
			}
		}
	}
	
	// Calculate score and direction
	if result.InDemandZone {
		result.SDScore = 85
		result.Direction = "bullish"
	} else if result.InSupplyZone {
		result.SDScore = 85
		result.Direction = "bearish"
	} else if result.NearestDemand != nil && result.NearestDemand.Fresh {
		result.SDScore = 70
		result.Direction = "bullish"
	} else if result.NearestSupply != nil && result.NearestSupply.Fresh {
		result.SDScore = 70
		result.Direction = "bearish"
	} else {
		result.SDScore = 50
	}
	
	return result
}

// ==================== AUCTION MARKET THEORY ====================

// AnalyzeAuctionTheory performs AMT analysis
func AnalyzeAuctionTheory(candles []Candle) *AuctionTheoryResult {
	if len(candles) < 50 {
		return nil
	}
	
	result := &AuctionTheoryResult{
		Direction: "neutral",
	}
	
	// Calculate Volume Profile
	mp := CalculateMarketProfile(candles, 20)
	if mp == nil {
		return result
	}
	
	result.POC = mp.POC
	result.VAH = mp.VAH
	result.VAL = mp.VAL
	
	currentPrice := candles[len(candles)-1].Close
	
	// Determine price position
	if currentPrice > result.VAH {
		result.AboveValue = true
	} else if currentPrice < result.VAL {
		result.BelowValue = true
	} else {
		result.InValue = true
	}
	
	// Determine market type
	// Balance = price staying in value area
	// Imbalance = price moving away from value
	
	recentInValue := 0
	for i := len(candles) - 10; i < len(candles); i++ {
		if candles[i].Close >= result.VAL && candles[i].Close <= result.VAH {
			recentInValue++
		}
	}
	
	if recentInValue >= 7 {
		result.MarketType = "balance"
	} else if recentInValue <= 3 {
		result.MarketType = "imbalance"
	} else {
		result.MarketType = "transition"
	}
	
	// Value Area Rejection strategy
	// If price opens outside VA and fails to get back in = trend day
	// If price falls back in = traverse to other side
	
	if result.AboveValue {
		// Above value = look for rejection or continuation
		if candles[len(candles)-1].Close < candles[len(candles)-1].Open {
			// Rejection candle above value = bearish
			result.Direction = "bearish"
			result.AuctionScore = 75
		} else {
			result.Direction = "bullish" // Continuation
			result.AuctionScore = 65
		}
	} else if result.BelowValue {
		if candles[len(candles)-1].Close > candles[len(candles)-1].Open {
			// Rejection candle below value = bullish
			result.Direction = "bullish"
			result.AuctionScore = 75
		} else {
			result.Direction = "bearish"
			result.AuctionScore = 65
		}
	} else {
		// In value = wait for breakout
		result.AuctionScore = 50
	}
	
	return result
}

// ==================== WYCKOFF PHASE ANALYSIS ====================

// AnalyzeWyckoffPhase performs Wyckoff phase analysis
func AnalyzeWyckoffPhase(candles []Candle) *WyckoffPhaseResult {
	if len(candles) < 100 {
		return nil
	}
	
	result := &WyckoffPhaseResult{
		Phase:     "unknown",
		Direction: "neutral",
	}
	
	// Analyze price and volume patterns
	// Accumulation: Price ranging at lows with increasing volume
	// Distribution: Price ranging at highs with increasing volume
	// Markup: Uptrend
	// Markdown: Downtrend
	
	// Find range
	high := candles[0].High
	low := candles[0].Low
	for _, c := range candles {
		if c.High > high {
			high = c.High
		}
		if c.Low < low {
			low = c.Low
		}
	}
	
	rangeSize := high - low
	currentPrice := candles[len(candles)-1].Close
	pricePosition := (currentPrice - low) / rangeSize
	
	// Volume analysis
	firstHalfVol := 0.0
	secondHalfVol := 0.0
	mid := len(candles) / 2
	
	for i := 0; i < mid; i++ {
		firstHalfVol += candles[i].Volume
	}
	for i := mid; i < len(candles); i++ {
		secondHalfVol += candles[i].Volume
	}
	
	volumeIncreasing := secondHalfVol > firstHalfVol*1.1
	
	// Determine phase
	if pricePosition < 0.3 && volumeIncreasing {
		result.Phase = "accumulation"
		result.Direction = "bullish"
		result.WyckoffScore = 75
	} else if pricePosition > 0.7 && volumeIncreasing {
		result.Phase = "distribution"
		result.Direction = "bearish"
		result.WyckoffScore = 75
	} else if pricePosition > 0.5 {
		// Check for markup
		upCandles := 0
		for i := len(candles) - 20; i < len(candles); i++ {
			if candles[i].Close > candles[i].Open {
				upCandles++
			}
		}
		if upCandles > 12 {
			result.Phase = "markup"
			result.Direction = "bullish"
			result.WyckoffScore = 70
		}
	} else {
		// Check for markdown
		downCandles := 0
		for i := len(candles) - 20; i < len(candles); i++ {
			if candles[i].Close < candles[i].Open {
				downCandles++
			}
		}
		if downCandles > 12 {
			result.Phase = "markdown"
			result.Direction = "bearish"
			result.WyckoffScore = 70
		}
	}
	
	// Detect Spring (false breakdown in accumulation)
	if result.Phase == "accumulation" {
		// Look for wick below support followed by close above
		for i := len(candles) - 10; i < len(candles); i++ {
			c := candles[i]
			// Spring: Low below range, close inside
			if c.Low < low && c.Close > low {
				result.SpringDetected = true
				result.SubPhase = "spring"
				result.WyckoffScore = 90
				result.Direction = "bullish"
				break
			}
		}
	}
	
	// Detect Upthrust (false breakout in distribution)
	if result.Phase == "distribution" {
		for i := len(candles) - 10; i < len(candles); i++ {
			c := candles[i]
			// Upthrust: High above range, close inside
			if c.High > high && c.Close < high {
				result.UpthrustDetected = true
				result.SubPhase = "upthrust"
				result.WyckoffScore = 90
				result.Direction = "bearish"
				break
			}
		}
	}
	
	return result
}


// ==================== UNIFIED SETUP GENERATION ====================

// GenerateUnifiedSetup generates the ultimate liquidity-first setup
func (lfs *LiquidityFirstStrategy) GenerateUnifiedSetup(timeframe string) *UnifiedSetup {
	candles, err := lfs.FetchMarketData(lfs.Symbol, timeframe, 200)
	if err != nil || len(candles) < 100 {
		return nil
	}
	
	setup := &UnifiedSetup{
		Timeframe: timeframe,
		Timestamp: time.Now(),
		IsValid:   false,
	}
	
	currentPrice := candles[len(candles)-1].Close
	currentTime := time.Now()
	
	// Check trading time
	if !lfs.isOptimalTime(currentTime) {
		setup.Reason = "Outside optimal trading time"
		return setup
	}
	
	// ==================== THE 3 PILLARS ====================
	
	// 1. LIQUIDITY ANALYSIS
	setup.Liquidity = AnalyzeLiquidity(candles)
	
	// 2. MARKET STRUCTURE
	setup.Structure = AnalyzeMarketStructureFull(candles)
	
	// 3. SUPPLY & DEMAND
	setup.SupplyDemand = AnalyzeSupplyDemandFull(candles)
	
	// ==================== SUPPORTING ANALYSIS ====================
	
	// 4. AUCTION THEORY
	setup.AuctionTheory = AnalyzeAuctionTheory(candles)
	
	// 5. WYCKOFF
	setup.Wyckoff = AnalyzeWyckoffPhase(candles)
	
	// ==================== CONFLUENCE CHECK ====================
	
	bullishScore := 0.0
	bearishScore := 0.0
	confluenceCount := 0
	
	// Liquidity (highest weight - 30%)
	if setup.Liquidity != nil {
		if setup.Liquidity.RecentSweep && setup.Liquidity.SweepConfirmed {
			confluenceCount++
			if setup.Liquidity.Direction == "bullish" {
				bullishScore += setup.Liquidity.LiquidityScore * 0.30
			} else {
				bearishScore += setup.Liquidity.LiquidityScore * 0.30
			}
		}
	}
	
	// Structure (25%)
	if setup.Structure != nil {
		confluenceCount++
		if setup.Structure.Direction == "bullish" {
			bullishScore += setup.Structure.StructureScore * 0.25
		} else if setup.Structure.Direction == "bearish" {
			bearishScore += setup.Structure.StructureScore * 0.25
		}
		
		// CHOCH bonus
		if setup.Structure.CHOCH {
			confluenceCount++
			if setup.Structure.Direction == "bullish" {
				bullishScore += 10
			} else {
				bearishScore += 10
			}
		}
	}
	
	// Supply/Demand (25%)
	if setup.SupplyDemand != nil {
		if setup.SupplyDemand.InDemandZone || setup.SupplyDemand.InSupplyZone {
			confluenceCount++
		}
		if setup.SupplyDemand.Direction == "bullish" {
			bullishScore += setup.SupplyDemand.SDScore * 0.25
		} else if setup.SupplyDemand.Direction == "bearish" {
			bearishScore += setup.SupplyDemand.SDScore * 0.25
		}
	}
	
	// Auction Theory (10%)
	if setup.AuctionTheory != nil {
		if setup.AuctionTheory.Direction == "bullish" {
			bullishScore += setup.AuctionTheory.AuctionScore * 0.10
		} else if setup.AuctionTheory.Direction == "bearish" {
			bearishScore += setup.AuctionTheory.AuctionScore * 0.10
		}
		
		// Value rejection bonus
		if (setup.AuctionTheory.BelowValue && setup.AuctionTheory.Direction == "bullish") ||
			(setup.AuctionTheory.AboveValue && setup.AuctionTheory.Direction == "bearish") {
			confluenceCount++
		}
	}
	
	// Wyckoff (10%)
	if setup.Wyckoff != nil {
		if setup.Wyckoff.Direction == "bullish" {
			bullishScore += setup.Wyckoff.WyckoffScore * 0.10
		} else if setup.Wyckoff.Direction == "bearish" {
			bearishScore += setup.Wyckoff.WyckoffScore * 0.10
		}
		
		// Spring/Upthrust bonus
		if setup.Wyckoff.SpringDetected || setup.Wyckoff.UpthrustDetected {
			confluenceCount++
			if setup.Wyckoff.Direction == "bullish" {
				bullishScore += 15
			} else {
				bearishScore += 15
			}
		}
	}
	
	// ==================== DECISION ====================
	
	// Determine direction
	if bullishScore > bearishScore+10 {
		setup.Direction = "bullish"
		setup.Confidence = math.Min(bullishScore, 100)
	} else if bearishScore > bullishScore+10 {
		setup.Direction = "bearish"
		setup.Confidence = math.Min(bearishScore, 100)
	} else {
		setup.Reason = fmt.Sprintf("No clear direction (Bull: %.1f, Bear: %.1f)", bullishScore, bearishScore)
		return setup
	}
	
	// Check minimum confidence
	if setup.Confidence < lfs.MinConfidence {
		setup.Reason = fmt.Sprintf("Confidence too low: %.1f%% (min: %.1f%%)", setup.Confidence, lfs.MinConfidence)
		return setup
	}
	
	// Check confluence
	if confluenceCount < 3 {
		setup.Reason = fmt.Sprintf("Insufficient confluence: %d/3", confluenceCount)
		return setup
	}
	
	// ==================== ENTRY CALCULATION ====================
	
	atr := calculateATR(candles[len(candles)-14:], 14)
	setup.Entry = currentPrice
	
	if setup.Direction == "bullish" {
		// Stop below liquidity sweep or demand zone
		setup.StopLoss = currentPrice - atr*1.5
		
		if setup.Liquidity != nil && setup.Liquidity.SweepPrice > 0 {
			sweepStop := setup.Liquidity.SweepPrice - atr*0.2
			if sweepStop > setup.StopLoss && sweepStop < currentPrice {
				setup.StopLoss = sweepStop
			}
		}
		
		if setup.SupplyDemand != nil && setup.SupplyDemand.NearestDemand != nil {
			demandStop := setup.SupplyDemand.NearestDemand.Low - atr*0.2
			if demandStop > setup.StopLoss {
				setup.StopLoss = demandStop
			}
		}
		
		// Targets
		risk := setup.Entry - setup.StopLoss
		setup.TP1 = setup.Entry + risk*3.0
		setup.TP2 = setup.Entry + risk*5.0
		setup.TP3 = setup.Entry + risk*8.0
		
		// Adjust to liquidity targets
		if setup.Liquidity != nil && setup.Liquidity.NearestBuyside > setup.Entry {
			if setup.Liquidity.NearestBuyside < setup.TP1 {
				setup.TP1 = setup.Liquidity.NearestBuyside
			}
		}
		
	} else { // bearish
		setup.StopLoss = currentPrice + atr*1.5
		
		if setup.Liquidity != nil && setup.Liquidity.SweepPrice > 0 {
			sweepStop := setup.Liquidity.SweepPrice + atr*0.2
			if sweepStop < setup.StopLoss && sweepStop > currentPrice {
				setup.StopLoss = sweepStop
			}
		}
		
		if setup.SupplyDemand != nil && setup.SupplyDemand.NearestSupply != nil {
			supplyStop := setup.SupplyDemand.NearestSupply.High + atr*0.2
			if supplyStop < setup.StopLoss {
				setup.StopLoss = supplyStop
			}
		}
		
		risk := setup.StopLoss - setup.Entry
		setup.TP1 = setup.Entry - risk*3.0
		setup.TP2 = setup.Entry - risk*5.0
		setup.TP3 = setup.Entry - risk*8.0
		
		if setup.Liquidity != nil && setup.Liquidity.NearestSellside < setup.Entry {
			if setup.Liquidity.NearestSellside > setup.TP1 {
				setup.TP1 = setup.Liquidity.NearestSellside
			}
		}
	}
	
	// Calculate RR
	risk := math.Abs(setup.Entry - setup.StopLoss)
	reward := math.Abs(setup.TP1 - setup.Entry)
	setup.RiskReward = reward / risk
	
	if setup.RiskReward < lfs.MinRiskReward {
		setup.Reason = fmt.Sprintf("RR too low: %.2f:1 (min: %.1f:1)", setup.RiskReward, lfs.MinRiskReward)
		return setup
	}
	
	// Determine setup name
	setup.SetupName = lfs.determineSetupName(setup)
	setup.IsValid = true
	setup.Reason = "VALID - All criteria met"
	
	return setup
}

// determineSetupName determines the primary setup type
func (lfs *LiquidityFirstStrategy) determineSetupName(setup *UnifiedSetup) string {
	// Priority: Liquidity Sweep > Wyckoff Spring > Value Rejection > S/D Zone
	
	if setup.Liquidity != nil && setup.Liquidity.RecentSweep && setup.Liquidity.SweepConfirmed {
		if setup.Liquidity.SweepType == "sellside" {
			return "Liquidity Sweep (Sellside)"
		}
		return "Liquidity Sweep (Buyside)"
	}
	
	if setup.Wyckoff != nil {
		if setup.Wyckoff.SpringDetected {
			return "Wyckoff Spring"
		}
		if setup.Wyckoff.UpthrustDetected {
			return "Wyckoff Upthrust"
		}
	}
	
	if setup.AuctionTheory != nil {
		if setup.AuctionTheory.BelowValue && setup.AuctionTheory.Direction == "bullish" {
			return "Value Area Rejection (Long)"
		}
		if setup.AuctionTheory.AboveValue && setup.AuctionTheory.Direction == "bearish" {
			return "Value Area Rejection (Short)"
		}
	}
	
	if setup.SupplyDemand != nil {
		if setup.SupplyDemand.InDemandZone {
			return "Fresh Demand Zone"
		}
		if setup.SupplyDemand.InSupplyZone {
			return "Fresh Supply Zone"
		}
	}
	
	if setup.Structure != nil && setup.Structure.CHOCH {
		return "Change of Character (CHOCH)"
	}
	
	return "Multi-Factor Confluence"
}

// isOptimalTime checks if current time is optimal
func (lfs *LiquidityFirstStrategy) isOptimalTime(t time.Time) bool {
	hour := t.UTC().Hour()
	weekday := t.Weekday()
	
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}
	
	// London: 08:00-12:00 UTC
	// NY: 13:00-17:00 UTC
	// Silver Bullet times: 10:00, 14:00, 19:00 UTC
	
	if hour >= 8 && hour < 12 {
		return true
	}
	if hour >= 13 && hour < 17 {
		return true
	}
	
	return false
}

// FetchMarketData fetches candlestick data
func (lfs *LiquidityFirstStrategy) FetchMarketData(symbol, interval string, limit int) ([]Candle, error) {
	sg := &SignalGenerator{}
	return sg.FetchMarketData(symbol, interval, limit)
}

// PrintUnifiedSetup prints the setup report
func (lfs *LiquidityFirstStrategy) PrintUnifiedSetup(setup *UnifiedSetup) {
	if setup == nil {
		return
	}
	
	fmt.Println("\n╔══════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║           🎯 LIQUIDITY-FIRST UNIFIED SETUP                           ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
	
	if !setup.IsValid {
		fmt.Printf("\n❌ INVALID: %s\n", setup.Reason)
		return
	}
	
	icon := "🟢"
	if setup.Direction == "bearish" {
		icon = "🔴"
	}
	
	fmt.Printf("\n%s %s - %s\n", icon, setup.SetupName, setup.Direction)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Confidence: %.1f%% | RR: %.2f:1\n", setup.Confidence, setup.RiskReward)
	
	fmt.Println("\n💰 ENTRY LEVELS:")
	fmt.Printf("   Entry: $%.2f | SL: $%.2f | TP1: $%.2f | TP2: $%.2f | TP3: $%.2f\n",
		setup.Entry, setup.StopLoss, setup.TP1, setup.TP2, setup.TP3)
	
	fmt.Println("\n📊 THE 3 PILLARS:")
	
	if setup.Liquidity != nil {
		sweepStatus := "No sweep"
		if setup.Liquidity.RecentSweep {
			sweepStatus = fmt.Sprintf("%s sweep", setup.Liquidity.SweepType)
			if setup.Liquidity.SweepConfirmed {
				sweepStatus += " ✓"
			}
		}
		fmt.Printf("   1. LIQUIDITY: %s | Score: %.0f | Dir: %s\n",
			sweepStatus, setup.Liquidity.LiquidityScore, setup.Liquidity.Direction)
	}
	
	if setup.Structure != nil {
		structStatus := setup.Structure.Trend
		if setup.Structure.CHOCH {
			structStatus += " + CHOCH"
		} else if setup.Structure.BOS {
			structStatus += " + BOS"
		}
		fmt.Printf("   2. STRUCTURE: %s | Score: %.0f | Dir: %s\n",
			structStatus, setup.Structure.StructureScore, setup.Structure.Direction)
	}
	
	if setup.SupplyDemand != nil {
		sdStatus := "Neutral"
		if setup.SupplyDemand.InDemandZone {
			sdStatus = "In Demand Zone"
		} else if setup.SupplyDemand.InSupplyZone {
			sdStatus = "In Supply Zone"
		}
		fmt.Printf("   3. S/D: %s | Score: %.0f | Dir: %s\n",
			sdStatus, setup.SupplyDemand.SDScore, setup.SupplyDemand.Direction)
	}
	
	fmt.Println("\n📈 SUPPORTING:")
	
	if setup.AuctionTheory != nil {
		fmt.Printf("   AMT: %s | POC: $%.2f | Dir: %s\n",
			setup.AuctionTheory.MarketType, setup.AuctionTheory.POC, setup.AuctionTheory.Direction)
	}
	
	if setup.Wyckoff != nil {
		wyckoffStatus := setup.Wyckoff.Phase
		if setup.Wyckoff.SpringDetected {
			wyckoffStatus += " + SPRING"
		}
		if setup.Wyckoff.UpthrustDetected {
			wyckoffStatus += " + UPTHRUST"
		}
		fmt.Printf("   Wyckoff: %s | Dir: %s\n", wyckoffStatus, setup.Wyckoff.Direction)
	}
	
	fmt.Println("\n╔══════════════════════════════════════════════════════════════════════╗")
	fmt.Printf("║  ✅ %s                                                    \n", setup.Reason)
	fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
}

// Start begins the strategy
func (lfs *LiquidityFirstStrategy) Start() {
	if lfs.IsRunning {
		return
	}
	
	lfs.IsRunning = true
	log.Println("\n🎯 LIQUIDITY-FIRST UNIFIED STRATEGY STARTED")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("Core Philosophy: Wait for liquidity grab, enter on reversal confirmation")
	log.Println("")
	log.Println("THE 3 PILLARS:")
	log.Println("   1. LIQUIDITY - Where are the stops?")
	log.Println("   2. STRUCTURE - What is the trend?")
	log.Println("   3. SUPPLY/DEMAND - Where is value?")
	log.Println("")
	log.Println("SUPPORTING:")
	log.Println("   • Auction Market Theory (POC, Value Area)")
	log.Println("   • Wyckoff (Accumulation, Distribution, Spring)")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	// Initial scan
	for _, tf := range lfs.Timeframes {
		setup := lfs.GenerateUnifiedSetup(tf)
		if setup != nil && setup.IsValid {
			lfs.PrintUnifiedSetup(setup)
		}
	}
	
	// Periodic scans
	ticker := time.NewTicker(lfs.CheckInterval)
	go func() {
		for range ticker.C {
			if lfs.IsRunning && lfs.TradesToday < lfs.MaxDailyTrades {
				for _, tf := range lfs.Timeframes {
					setup := lfs.GenerateUnifiedSetup(tf)
					if setup != nil && setup.IsValid {
						lfs.PrintUnifiedSetup(setup)
						lfs.TradesToday++
					}
				}
			}
		}
	}()
}
