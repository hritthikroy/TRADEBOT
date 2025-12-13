package daily

import (
	"fmt"
	"log"
	"math"
	"time"
)

// ==================== ULTIMATE DAILY TRADING STRATEGY ====================
// Combines ALL concepts for maximum profitability
// Target: 85-95% win rate, 3:1+ RR, institutional-grade execution

// UltimateDailyStrategy is the most advanced strategy
type UltimateDailyStrategy struct {
	Symbol            string
	ActiveTimeframes  []string
	CheckInterval     time.Duration
	MinConfidence     float64
	MinConfluence     int
	MinRiskReward     float64
	MaxDailyTrades    int
	TradesToday       int
	DailyPnL          float64
	MaxDailyLoss      float64
	IsRunning         bool
	
	// Strategy rotation
	CurrentStrategy   string
	StrategyRotation  []string
}

// DailySetup represents a complete trading setup
type DailySetup struct {
	// Setup identification
	SetupName       string
	SetupType       string
	Timeframe       string
	Timestamp       time.Time
	
	// Direction and confidence
	Direction       string
	Confidence      float64
	Confluence      int
	
	// Entry levels
	Entry           float64
	StopLoss        float64
	TP1             float64
	TP2             float64
	TP3             float64
	RiskReward      float64
	
	// Analysis scores
	Scores          map[string]float64
	
	// Validation
	IsValid         bool
	ValidationMsg   string
	
	// Risk management
	PositionSize    float64
	RiskAmount      float64
	
	// All analysis data
	Delta           *DeltaAnalysis
	Pivots          *PivotPoints
	VWAP            *VWAPAnalysis
	MarketProfile   *MarketProfile
	MTF             *ComprehensiveMTFAnalysis
	PO3             *PO3Analysis
	LiqSweep        *LiquiditySweepAnalysis
	Institutional   *InstitutionalAnalysis
	MMM             *MMMAnalysis
	SD              *SDAnalysis
	ICT             *ICTAnalysis
	Volatility      *VolatilityAnalysis
}

// NewUltimateDailyStrategy creates the ultimate strategy
func NewUltimateDailyStrategy() *UltimateDailyStrategy {
	return &UltimateDailyStrategy{
		Symbol:           "BTCUSDT",
		ActiveTimeframes: []string{"15m", "1h", "4h"},
		CheckInterval:    60 * time.Second,
		MinConfidence:    85.0,
		MinConfluence:    6,
		MinRiskReward:    3.0,
		MaxDailyTrades:   3,
		TradesToday:      0,
		MaxDailyLoss:     2.0, // 2% max daily loss
		IsRunning:        false,
		
		// Strategy rotation for different market conditions
		StrategyRotation: []string{
			"SilverBullet",
			"Unicorn",
			"LiquiditySweep",
			"PowerOf3",
			"SupplyDemand",
			"TurtleSoup",
		},
		CurrentStrategy: "Auto",
	}
}


// GenerateUltimateSetup generates the best possible setup
func (uds *UltimateDailyStrategy) GenerateUltimateSetup(timeframe string) *DailySetup {
	// Fetch market data
	candles, err := uds.FetchMarketData(uds.Symbol, timeframe, 200)
	if err != nil || len(candles) < 100 {
		return nil
	}
	
	setup := &DailySetup{
		Timeframe: timeframe,
		Timestamp: time.Now(),
		Scores:    make(map[string]float64),
		IsValid:   false,
	}
	
	currentPrice := candles[len(candles)-1].Close
	currentTime := time.Now()
	
	// ==================== PHASE 1: PRE-FILTERS ====================
	
	// Time filter - only trade optimal sessions
	if !uds.isOptimalTradingTime(currentTime) {
		setup.ValidationMsg = "Outside optimal trading time"
		return setup
	}
	
	// Volatility filter
	setup.Volatility = AnalyzeVolatility(candles)
	if setup.Volatility == nil || !ShouldTradeVolatility(setup.Volatility) {
		setup.ValidationMsg = "Volatility not optimal"
		return setup
	}
	
	// ==================== PHASE 2: COMPREHENSIVE ANALYSIS ====================
	
	// 1. Delta Analysis
	setup.Delta = CalculateDelta(candles)
	
	// 2. Pivot Points
	setup.Pivots = CalculatePivotPoints(candles)
	
	// 3. VWAP
	setup.VWAP = CalculateVWAP(candles)
	
	// 4. Market Profile
	setup.MarketProfile = CalculateMarketProfile(candles, 20)
	
	// 5. Multi-Timeframe (All 13 TFs)
	allTFData, _ := FetchAllTimeframeData(uds.Symbol, 100)
	if len(allTFData) > 0 {
		setup.MTF = PerformComprehensiveMTFAnalysis(allTFData)
	}
	
	// 6. Power of 3
	setup.PO3 = AnalyzeDailyPO3(candles)
	
	// 7. Liquidity Sweep
	setup.LiqSweep = PerformLiquiditySweepAnalysis(candles)
	
	// 8. Institutional Setups
	setup.Institutional = PerformInstitutionalAnalysis(candles)
	
	// 9. Market Maker Model
	setup.MMM = PerformMMMAnalysis(candles)
	
	// 10. Supply & Demand
	setup.SD = FindSupplyDemandZones(candles)
	
	// 11. ICT/SMC
	setup.ICT = PerformICTAnalysis(candles)
	
	// ==================== PHASE 3: SCORING ====================
	
	bullishScore := 0.0
	bearishScore := 0.0
	confluenceCount := 0
	
	// Delta Score (15 points)
	if setup.Delta != nil {
		if setup.Delta.DeltaTrend == "bullish" && !setup.Delta.DeltaDivergence {
			setup.Scores["delta"] = 15
			bullishScore += 15
			confluenceCount++
		} else if setup.Delta.DeltaTrend == "bearish" && !setup.Delta.DeltaDivergence {
			setup.Scores["delta"] = -15
			bearishScore += 15
			confluenceCount++
		}
		
		// Divergence penalty
		if setup.Delta.DeltaDivergence {
			setup.Scores["delta_divergence"] = -20
		}
	}
	
	// Pivot Score (12 points)
	if setup.Pivots != nil {
		pivotDir, pivotStr := GetPivotSignal(setup.Pivots, currentPrice)
		if pivotDir == "bullish" {
			setup.Scores["pivot"] = pivotStr * 0.12
			bullishScore += pivotStr * 0.12
			if pivotStr >= 70 {
				confluenceCount++
			}
		} else if pivotDir == "bearish" {
			setup.Scores["pivot"] = -pivotStr * 0.12
			bearishScore += pivotStr * 0.12
			if pivotStr >= 70 {
				confluenceCount++
			}
		}
	}
	
	// VWAP Score (10 points)
	if setup.VWAP != nil {
		if setup.VWAP.Signal == "bullish" {
			setup.Scores["vwap"] = setup.VWAP.Strength * 0.1
			bullishScore += setup.VWAP.Strength * 0.1
			if setup.VWAP.Strength >= 70 {
				confluenceCount++
			}
		} else if setup.VWAP.Signal == "bearish" {
			setup.Scores["vwap"] = -setup.VWAP.Strength * 0.1
			bearishScore += setup.VWAP.Strength * 0.1
			if setup.VWAP.Strength >= 70 {
				confluenceCount++
			}
		}
	}
	
	// Market Profile Score (10 points)
	if setup.MarketProfile != nil {
		if setup.MarketProfile.Signal == "bullish" {
			setup.Scores["profile"] = setup.MarketProfile.Strength * 0.1
			bullishScore += setup.MarketProfile.Strength * 0.1
			confluenceCount++
		} else if setup.MarketProfile.Signal == "bearish" {
			setup.Scores["profile"] = -setup.MarketProfile.Strength * 0.1
			bearishScore += setup.MarketProfile.Strength * 0.1
			confluenceCount++
		}
	}
	
	// MTF Score (20 points - highest weight)
	if setup.MTF != nil {
		mtfScore := GetMTFConfluenceScore(setup.MTF, "bullish")
		mtfScoreBear := GetMTFConfluenceScore(setup.MTF, "bearish")
		
		if mtfScore > mtfScoreBear {
			setup.Scores["mtf"] = mtfScore * 0.2
			bullishScore += mtfScore * 0.2
			if mtfScore >= 70 {
				confluenceCount++
			}
		} else if mtfScoreBear > mtfScore {
			setup.Scores["mtf"] = -mtfScoreBear * 0.2
			bearishScore += mtfScoreBear * 0.2
			if mtfScoreBear >= 70 {
				confluenceCount++
			}
		}
		
		// All TF aligned bonus
		if setup.MTF.AllAligned {
			if setup.MTF.Direction == "bullish" {
				setup.Scores["mtf_all"] = 15
				bullishScore += 15
			} else {
				setup.Scores["mtf_all"] = -15
				bearishScore += 15
			}
			confluenceCount++
		}
	}
	
	// Power of 3 Score (15 points)
	if setup.PO3 != nil {
		po3Dir, po3Str := GetPO3Signal(setup.PO3, currentPrice)
		if po3Dir == "bullish" {
			setup.Scores["po3"] = po3Str * 0.15
			bullishScore += po3Str * 0.15
			if po3Str >= 70 {
				confluenceCount++
			}
		} else if po3Dir == "bearish" {
			setup.Scores["po3"] = -po3Str * 0.15
			bearishScore += po3Str * 0.15
			if po3Str >= 70 {
				confluenceCount++
			}
		}
		
		// Optimal entry bonus
		if setup.PO3.OptimalEntry {
			setup.Scores["po3_optimal"] = 10
			confluenceCount++
		}
	}
	
	// Liquidity Sweep Score (15 points)
	if setup.LiqSweep != nil {
		lsDir, lsStr := GetLiquiditySweepSignal(setup.LiqSweep, currentPrice)
		if lsDir == "bullish" {
			setup.Scores["liq_sweep"] = lsStr * 0.15
			bullishScore += lsStr * 0.15
			if lsStr >= 70 {
				confluenceCount++
			}
		} else if lsDir == "bearish" {
			setup.Scores["liq_sweep"] = -lsStr * 0.15
			bearishScore += lsStr * 0.15
			if lsStr >= 70 {
				confluenceCount++
			}
		}
		
		// Confirmed sweep bonus
		if setup.LiqSweep.RecentSweep != nil && setup.LiqSweep.RecentSweep.Confirmed {
			if setup.LiqSweep.RecentSweep.Type == "sellside" {
				setup.Scores["sweep_confirmed"] = 12
				bullishScore += 12
			} else {
				setup.Scores["sweep_confirmed"] = -12
				bearishScore += 12
			}
			confluenceCount++
		}
	}
	
	// Institutional Score (18 points)
	if setup.Institutional != nil {
		instDir, instStr := GetInstitutionalSignal(setup.Institutional)
		if instDir == "bullish" {
			setup.Scores["institutional"] = instStr * 0.18
			bullishScore += instStr * 0.18
			if instStr >= 70 {
				confluenceCount++
			}
		} else if instDir == "bearish" {
			setup.Scores["institutional"] = -instStr * 0.18
			bearishScore += instStr * 0.18
			if instStr >= 70 {
				confluenceCount++
			}
		}
		
		// Specific setup bonuses
		if setup.Institutional.SilverBullet != nil && setup.Institutional.SilverBullet.Valid {
			setup.SetupName = "Silver Bullet"
			setup.SetupType = "institutional"
			confluenceCount++
		}
		if setup.Institutional.Unicorn != nil && setup.Institutional.Unicorn.Valid {
			setup.SetupName = "Unicorn"
			setup.SetupType = "institutional"
			confluenceCount++
		}
	}
	
	// MMM Score (12 points)
	if setup.MMM != nil {
		mmmStr := GetMMMSignalStrength(setup.MMM, "bullish")
		mmmStrBear := GetMMMSignalStrength(setup.MMM, "bearish")
		
		if mmmStr > mmmStrBear {
			setup.Scores["mmm"] = mmmStr * 0.12
			bullishScore += mmmStr * 0.12
			if mmmStr >= 60 {
				confluenceCount++
			}
		} else if mmmStrBear > mmmStr {
			setup.Scores["mmm"] = -mmmStrBear * 0.12
			bearishScore += mmmStrBear * 0.12
			if mmmStrBear >= 60 {
				confluenceCount++
			}
		}
	}
	
	// Supply/Demand Score (10 points)
	if setup.SD != nil {
		sdStrBull := GetSDSignalStrength(setup.SD, "bullish", currentPrice)
		sdStrBear := GetSDSignalStrength(setup.SD, "bearish", currentPrice)
		
		if sdStrBull > sdStrBear {
			setup.Scores["sd"] = sdStrBull * 0.1
			bullishScore += sdStrBull * 0.1
			if sdStrBull >= 70 {
				confluenceCount++
			}
		} else if sdStrBear > sdStrBull {
			setup.Scores["sd"] = -sdStrBear * 0.1
			bearishScore += sdStrBear * 0.1
			if sdStrBear >= 70 {
				confluenceCount++
			}
		}
	}
	
	// ICT Score (10 points)
	if setup.ICT != nil {
		if setup.ICT.Structure.Trend == "bullish" {
			setup.Scores["ict"] = float64(setup.ICT.Confluence) * 2
			bullishScore += float64(setup.ICT.Confluence) * 2
			if setup.ICT.Confluence >= 3 {
				confluenceCount++
			}
		} else if setup.ICT.Structure.Trend == "bearish" {
			setup.Scores["ict"] = -float64(setup.ICT.Confluence) * 2
			bearishScore += float64(setup.ICT.Confluence) * 2
			if setup.ICT.Confluence >= 3 {
				confluenceCount++
			}
		}
		
		// OTE bonus
		if setup.ICT.OTE {
			setup.Scores["ote"] = 8
			confluenceCount++
		}
	}
	
	// ==================== PHASE 4: DECISION ====================
	
	setup.Confluence = confluenceCount
	
	// Determine direction
	if bullishScore > bearishScore+15 {
		setup.Direction = "bullish"
		setup.Confidence = math.Min(bullishScore, 100)
	} else if bearishScore > bullishScore+15 {
		setup.Direction = "bearish"
		setup.Confidence = math.Min(bearishScore, 100)
	} else {
		setup.ValidationMsg = fmt.Sprintf("No clear direction (Bull: %.1f, Bear: %.1f)", bullishScore, bearishScore)
		return setup
	}
	
	// Validate confluence
	if setup.Confluence < uds.MinConfluence {
		setup.ValidationMsg = fmt.Sprintf("Insufficient confluence: %d/%d", setup.Confluence, uds.MinConfluence)
		return setup
	}
	
	// Validate confidence
	if setup.Confidence < uds.MinConfidence {
		setup.ValidationMsg = fmt.Sprintf("Confidence too low: %.1f%% (min: %.1f%%)", setup.Confidence, uds.MinConfidence)
		return setup
	}
	
	// ==================== PHASE 5: ENTRY CALCULATION ====================
	
	atr := calculateATR(candles[len(candles)-14:], 14)
	setup.Entry = currentPrice
	
	if setup.Direction == "bullish" {
		// Find optimal stop loss
		setup.StopLoss = currentPrice - atr*1.5
		
		// Use better levels if available
		if setup.SD != nil && setup.SD.NearestDemand != nil {
			demandStop := setup.SD.NearestDemand.Low - atr*0.2
			if demandStop > setup.StopLoss && demandStop < currentPrice-atr*0.5 {
				setup.StopLoss = demandStop
			}
		}
		
		if setup.Pivots != nil && setup.Pivots.S1 < currentPrice {
			pivotStop := setup.Pivots.S1 - atr*0.2
			if pivotStop > setup.StopLoss {
				setup.StopLoss = pivotStop
			}
		}
		
		// Calculate targets
		risk := setup.Entry - setup.StopLoss
		setup.TP1 = setup.Entry + risk*3.0
		setup.TP2 = setup.Entry + risk*5.0
		setup.TP3 = setup.Entry + risk*8.0
		
		// Adjust targets to key levels
		if setup.Pivots != nil {
			if setup.Pivots.R1 > setup.Entry && setup.Pivots.R1 < setup.TP1 {
				setup.TP1 = setup.Pivots.R1
			}
			if setup.Pivots.R2 > setup.TP1 && setup.Pivots.R2 < setup.TP2 {
				setup.TP2 = setup.Pivots.R2
			}
		}
		
	} else { // bearish
		setup.StopLoss = currentPrice + atr*1.5
		
		if setup.SD != nil && setup.SD.NearestSupply != nil {
			supplyStop := setup.SD.NearestSupply.High + atr*0.2
			if supplyStop < setup.StopLoss && supplyStop > currentPrice+atr*0.5 {
				setup.StopLoss = supplyStop
			}
		}
		
		if setup.Pivots != nil && setup.Pivots.R1 > currentPrice {
			pivotStop := setup.Pivots.R1 + atr*0.2
			if pivotStop < setup.StopLoss {
				setup.StopLoss = pivotStop
			}
		}
		
		risk := setup.StopLoss - setup.Entry
		setup.TP1 = setup.Entry - risk*3.0
		setup.TP2 = setup.Entry - risk*5.0
		setup.TP3 = setup.Entry - risk*8.0
		
		if setup.Pivots != nil {
			if setup.Pivots.S1 < setup.Entry && setup.Pivots.S1 > setup.TP1 {
				setup.TP1 = setup.Pivots.S1
			}
			if setup.Pivots.S2 < setup.TP1 && setup.Pivots.S2 > setup.TP2 {
				setup.TP2 = setup.Pivots.S2
			}
		}
	}
	
	// Calculate risk:reward
	risk := math.Abs(setup.Entry - setup.StopLoss)
	reward := math.Abs(setup.TP1 - setup.Entry)
	setup.RiskReward = reward / risk
	
	// Validate RR
	if setup.RiskReward < uds.MinRiskReward {
		setup.ValidationMsg = fmt.Sprintf("RR too low: %.2f:1 (min: %.1f:1)", setup.RiskReward, uds.MinRiskReward)
		return setup
	}
	
	// Set setup name if not already set
	if setup.SetupName == "" {
		setup.SetupName = uds.determineSetupName(setup)
	}
	
	setup.IsValid = true
	setup.ValidationMsg = "VALID - All criteria met"
	
	return setup
}


// isOptimalTradingTime checks if current time is optimal for trading
func (uds *UltimateDailyStrategy) isOptimalTradingTime(t time.Time) bool {
	hour := t.UTC().Hour()
	weekday := t.Weekday()
	
	// No weekend trading
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}
	
	// No Friday afternoon
	if weekday == time.Friday && hour >= 18 {
		return false
	}
	
	// Optimal times (UTC):
	// London Open: 08:00-12:00
	// NY Open: 13:00-17:00
	// London-NY Overlap: 13:00-16:00 (BEST)
	
	if hour >= 8 && hour < 12 {
		return true // London
	}
	if hour >= 13 && hour < 17 {
		return true // NY
	}
	
	return false
}

// determineSetupName determines the primary setup name
func (uds *UltimateDailyStrategy) determineSetupName(setup *DailySetup) string {
	// Priority order
	if setup.Institutional != nil {
		if setup.Institutional.SilverBullet != nil && setup.Institutional.SilverBullet.Valid {
			return "Silver Bullet"
		}
		if setup.Institutional.Unicorn != nil && setup.Institutional.Unicorn.Valid {
			return "Unicorn"
		}
		if setup.Institutional.TurtleSoup != nil && setup.Institutional.TurtleSoup.Valid {
			return "Turtle Soup"
		}
	}
	
	if setup.LiqSweep != nil && setup.LiqSweep.RecentSweep != nil && setup.LiqSweep.RecentSweep.Confirmed {
		return "Liquidity Sweep"
	}
	
	if setup.PO3 != nil && setup.PO3.OptimalEntry {
		return "Power of 3"
	}
	
	if setup.SD != nil && setup.SD.InZone {
		return "Supply/Demand"
	}
	
	if setup.MTF != nil && setup.MTF.AllAligned {
		return "MTF Confluence"
	}
	
	return "Multi-Factor"
}

// FetchMarketData fetches candlestick data
func (uds *UltimateDailyStrategy) FetchMarketData(symbol, interval string, limit int) ([]Candle, error) {
	sg := &SignalGenerator{}
	return sg.FetchMarketData(symbol, interval, limit)
}

// PrintSetupReport prints a detailed setup report
func (uds *UltimateDailyStrategy) PrintSetupReport(setup *DailySetup) {
	if setup == nil {
		fmt.Println("No setup available")
		return
	}
	
	fmt.Println("\n╔══════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║              🎯 ULTIMATE DAILY TRADING SETUP                         ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
	
	if !setup.IsValid {
		fmt.Printf("\n❌ INVALID SETUP: %s\n", setup.ValidationMsg)
		return
	}
	
	icon := "🟢"
	if setup.Direction == "bearish" {
		icon = "🔴"
	}
	
	fmt.Printf("\n%s SETUP: %s (%s)\n", icon, setup.SetupName, setup.Direction)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	fmt.Printf("📊 Timeframe: %s\n", setup.Timeframe)
	fmt.Printf("🎯 Confidence: %.1f%%\n", setup.Confidence)
	fmt.Printf("🔗 Confluence: %d factors\n", setup.Confluence)
	fmt.Printf("📈 Risk:Reward: %.2f:1\n", setup.RiskReward)
	
	fmt.Println("\n💰 ENTRY LEVELS:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("   Entry:     $%.2f\n", setup.Entry)
	fmt.Printf("   Stop Loss: $%.2f\n", setup.StopLoss)
	fmt.Printf("   TP1:       $%.2f (3:1)\n", setup.TP1)
	fmt.Printf("   TP2:       $%.2f (5:1)\n", setup.TP2)
	fmt.Printf("   TP3:       $%.2f (8:1)\n", setup.TP3)
	
	fmt.Println("\n📊 ANALYSIS SCORES:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	for name, score := range setup.Scores {
		icon := "⚪"
		if score > 5 {
			icon = "🟢"
		} else if score < -5 {
			icon = "🔴"
		}
		fmt.Printf("   %s %-20s: %+.1f\n", icon, name, score)
	}
	
	// MTF Summary
	if setup.MTF != nil {
		fmt.Println("\n📈 MTF CONFLUENCE:")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Printf("   Direction: %s (%.1f%% confidence)\n", setup.MTF.Direction, setup.MTF.Confidence)
		fmt.Printf("   Aligned: %d/%d timeframes\n", 
			int(math.Max(float64(setup.MTF.BullishCount), float64(setup.MTF.BearishCount))),
			setup.MTF.TotalTimeframes)
		fmt.Printf("   Scalping: %s | Day: %s | Swing: %s | Position: %s\n",
			setup.MTF.ScalpingBias, setup.MTF.DayTradingBias, 
			setup.MTF.SwingBias, setup.MTF.PositionBias)
	}
	
	// Pivot Levels
	if setup.Pivots != nil {
		fmt.Println("\n📍 PIVOT LEVELS:")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Printf("   R2: $%.2f | R1: $%.2f | Pivot: $%.2f | S1: $%.2f | S2: $%.2f\n",
			setup.Pivots.R2, setup.Pivots.R1, setup.Pivots.Pivot, 
			setup.Pivots.S1, setup.Pivots.S2)
		fmt.Printf("   Position: %s | Nearest: %s ($%.2f)\n",
			setup.Pivots.PricePosition, setup.Pivots.NearestType, setup.Pivots.NearestLevel)
	}
	
	// Delta
	if setup.Delta != nil {
		fmt.Println("\n📊 DELTA ANALYSIS:")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Printf("   Trend: %s | Momentum: %s | Divergence: %v\n",
			setup.Delta.DeltaTrend, setup.Delta.DeltaMomentum, setup.Delta.DeltaDivergence)
	}
	
	fmt.Println("\n╔══════════════════════════════════════════════════════════════════════╗")
	fmt.Printf("║  ✅ VALID SETUP - %s                                    \n", setup.ValidationMsg)
	fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
}

// Start begins the ultimate daily strategy
func (uds *UltimateDailyStrategy) Start() {
	if uds.IsRunning {
		log.Println("⚠️  Ultimate strategy already running")
		return
	}
	
	uds.IsRunning = true
	log.Println("\n🏆 ULTIMATE DAILY TRADING STRATEGY STARTED")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("📊 Symbol: %s", uds.Symbol)
	log.Printf("⏱️  Timeframes: %v", uds.ActiveTimeframes)
	log.Printf("🎯 Min Confidence: %.1f%%", uds.MinConfidence)
	log.Printf("🔗 Min Confluence: %d factors", uds.MinConfluence)
	log.Printf("📈 Min RR: %.1f:1", uds.MinRiskReward)
	log.Printf("📉 Max Daily Loss: %.1f%%", uds.MaxDailyLoss)
	log.Printf("🔄 Max Daily Trades: %d", uds.MaxDailyTrades)
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("\n📚 Active Analysis:")
	log.Println("   ✅ Delta Analysis")
	log.Println("   ✅ Pivot Points (Standard, Fib, Camarilla, Woodie, DeMark)")
	log.Println("   ✅ VWAP & Bands")
	log.Println("   ✅ Market Profile (POC, VAH, VAL)")
	log.Println("   ✅ 13-Timeframe MTF Confluence")
	log.Println("   ✅ Power of 3 (AMD)")
	log.Println("   ✅ Liquidity Sweep")
	log.Println("   ✅ Institutional Setups (Silver Bullet, Unicorn, etc.)")
	log.Println("   ✅ Market Maker Model")
	log.Println("   ✅ Supply & Demand")
	log.Println("   ✅ ICT/SMC Concepts")
	log.Println("   ✅ Volatility Filter")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	// Initial scan
	uds.ScanForSetups()
	
	// Periodic scans
	ticker := time.NewTicker(uds.CheckInterval)
	go func() {
		for range ticker.C {
			if uds.IsRunning {
				uds.ScanForSetups()
			}
		}
	}()
}

// ScanForSetups scans all timeframes for setups
func (uds *UltimateDailyStrategy) ScanForSetups() {
	// Reset daily counter at midnight
	now := time.Now()
	if now.Hour() == 0 && now.Minute() == 0 {
		uds.TradesToday = 0
		uds.DailyPnL = 0
		log.Println("🔄 Daily counters reset")
	}
	
	// Check daily limits
	if uds.TradesToday >= uds.MaxDailyTrades {
		log.Printf("⏸️  Daily trade limit reached (%d/%d)", uds.TradesToday, uds.MaxDailyTrades)
		return
	}
	
	if uds.DailyPnL <= -uds.MaxDailyLoss {
		log.Printf("⏸️  Daily loss limit reached (%.2f%%)", uds.DailyPnL)
		return
	}
	
	log.Printf("\n🔍 Scanning for ultimate setups...")
	
	for _, tf := range uds.ActiveTimeframes {
		setup := uds.GenerateUltimateSetup(tf)
		
		if setup != nil && setup.IsValid {
			uds.PrintSetupReport(setup)
			uds.TradesToday++
			
			// Save signal
			signal := &CreateSignalRequest{
				SignalID:   fmt.Sprintf("%d.%d", time.Now().UnixMilli(), time.Now().Nanosecond()),
				SignalType: "BUY",
				Symbol:     uds.Symbol,
				EntryPrice: setup.Entry,
				StopLoss:   setup.StopLoss,
				TP1:        setup.TP1,
				TP2:        setup.TP2,
				TP3:        setup.TP3,
				Strength:   int(setup.Confidence),
			}
			
			if setup.Direction == "bearish" {
				signal.SignalType = "SELL"
			}
			
			sg := &SignalGenerator{}
			sg.SaveSignal(signal, tf)
			
			if uds.TradesToday >= uds.MaxDailyTrades {
				break
			}
		}
	}
}
