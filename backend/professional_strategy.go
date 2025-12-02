package main

import (
	"fmt"
	"log"
	"math"
	"time"
)

// ==================== PROFESSIONAL INSTITUTIONAL STRATEGY ====================
// Combines ALL advanced concepts for maximum performance
// Target: 80-90% win rate with institutional-grade setups

// ProfessionalSignalGenerator is the ultimate signal generator
type ProfessionalSignalGenerator struct {
	Symbol            string
	Timeframes        []string
	CheckInterval     time.Duration
	MinConfidence     float64
	MinConfluence     int
	MaxSignalsPerDay  int
	SignalsToday      int
	IsRunning         bool
	
	// Strategy weights
	Weights           StrategyWeights
}

// StrategyWeights holds weights for each analysis component
type StrategyWeights struct {
	PO3             float64
	LiquiditySweep  float64
	Mirror          float64
	Institutional   float64
	MMM             float64
	SupplyDemand    float64
	SessionLiq      float64
	Volatility      float64
	ICT             float64
	OrderFlow       float64
	MTF             float64
	Patterns        float64
}

// NewProfessionalSignalGenerator creates the professional generator
func NewProfessionalSignalGenerator() *ProfessionalSignalGenerator {
	return &ProfessionalSignalGenerator{
		Symbol:           "BTCUSDT",
		Timeframes:       []string{"15m", "1h", "4h"},
		CheckInterval:    60 * time.Second,
		MinConfidence:    80.0, // Very high threshold
		MinConfluence:    5,    // Require 5+ factors
		MaxSignalsPerDay: 3,    // Quality over quantity
		SignalsToday:     0,
		IsRunning:        false,
		
		// Optimized weights
		Weights: StrategyWeights{
			PO3:            15,
			LiquiditySweep: 15,
			Mirror:         10,
			Institutional:  20,
			MMM:            15,
			SupplyDemand:   10,
			SessionLiq:     5,
			Volatility:     5,
			ICT:            10,
			OrderFlow:      10,
			MTF:            10,
			Patterns:       5,
		},
	}
}

// ComprehensiveAnalysis holds all analysis results
type ComprehensiveAnalysis struct {
	// Phase 1 Analysis
	MMM           *MMMAnalysis
	SD            *SDAnalysis
	SessionLiq    *LiquidityMap
	Volatility    *VolatilityAnalysis
	
	// Phase 2 Analysis (New)
	PO3           *PO3Analysis
	LiqSweep      *LiquiditySweepAnalysis
	Mirror        *MirrorAnalysis
	Institutional *InstitutionalAnalysis
	
	// Existing Analysis
	ICT           *ICTAnalysis
	OrderFlow     *OrderFlowAnalysis
	MTF           *MultiTimeframeAnalysis
	ComprehensiveMTF *ComprehensiveMTFAnalysis // NEW: All timeframes
	Patterns      []CandlestickPattern
	
	// Scores
	Scores        map[string]float64
	BullishScore  float64
	BearishScore  float64
	
	// Final Decision
	Direction     string
	Confidence    float64
	Confluence    int
	SetupType     string
	Entry         float64
	StopLoss      float64
	TP1           float64
	TP2           float64
	TP3           float64
	RiskReward    float64
}

// PerformComprehensiveAnalysis performs ALL analysis
func (psg *ProfessionalSignalGenerator) PerformComprehensiveAnalysis(candles []Candle, timeframe string) *ComprehensiveAnalysis {
	if len(candles) < 100 {
		return nil
	}
	
	analysis := &ComprehensiveAnalysis{
		Scores: make(map[string]float64),
	}
	
	// ==================== PHASE 1: FILTERS ====================
	
	currentTime := time.Now()
	
	// Session Filter
	if !ShouldTradeSession(currentTime) {
		log.Printf("⏸️  [%s] Outside optimal trading session", timeframe)
		return nil
	}
	
	// Volatility Filter
	analysis.Volatility = AnalyzeVolatility(candles)
	if !ShouldTradeVolatility(analysis.Volatility) {
		log.Printf("⏸️  [%s] Volatility not optimal: %s", timeframe, analysis.Volatility.Regime)
		return nil
	}
	
	// ==================== PHASE 2: ALL ANALYSIS ====================
	
	// 1. Power of 3
	analysis.PO3 = AnalyzeDailyPO3(candles)
	
	// 2. Liquidity Sweep
	analysis.LiqSweep = PerformLiquiditySweepAnalysis(candles)
	
	// 3. Mirror Market
	analysis.Mirror = PerformMirrorAnalysis(candles)
	
	// 4. Institutional Setups
	analysis.Institutional = PerformInstitutionalAnalysis(candles)
	
	// 5. Market Maker Model
	analysis.MMM = PerformMMMAnalysis(candles)
	
	// 6. Supply & Demand
	analysis.SD = FindSupplyDemandZones(candles)
	
	// 7. Session Liquidity
	analysis.SessionLiq = MapSessionLiquidity(candles)
	
	// 8. ICT/SMC
	analysis.ICT = PerformICTAnalysis(candles)
	
	// 9. Order Flow
	analysis.OrderFlow = PerformOrderFlowAnalysis(candles)
	
	// 10. Multi-Timeframe (Basic)
	analysis.MTF = PerformMultiTimeframeAnalysis(candles, timeframe)
	
	// 10b. Comprehensive MTF (All 13 timeframes)
	// Fetch all timeframe data for comprehensive analysis
	allTFData, err := FetchAllTimeframeData(psg.Symbol, 100)
	if err == nil && len(allTFData) > 0 {
		analysis.ComprehensiveMTF = PerformComprehensiveMTFAnalysis(allTFData)
	}
	
	// 11. Candlestick Patterns
	analysis.Patterns = RecognizeAllPatterns(candles, timeframe)
	
	// ==================== PHASE 3: SCORING ====================
	
	currentPrice := candles[len(candles)-1].Close
	
	// Score each component
	psg.scoreAllComponents(analysis, currentPrice)
	
	// Calculate total scores
	analysis.BullishScore = 0
	analysis.BearishScore = 0
	
	for _, score := range analysis.Scores {
		if score > 0 {
			analysis.BullishScore += score
		} else {
			analysis.BearishScore += math.Abs(score)
		}
	}
	
	// ==================== PHASE 4: DECISION ====================
	
	// Count confluence factors
	analysis.Confluence = psg.countConfluence(analysis)
	
	// Require minimum confluence
	if analysis.Confluence < psg.MinConfluence {
		log.Printf("⏸️  [%s] Insufficient confluence: %d/%d", 
			timeframe, analysis.Confluence, psg.MinConfluence)
		return nil
	}
	
	// Determine direction
	if analysis.BullishScore > analysis.BearishScore && analysis.BullishScore > 50 {
		analysis.Direction = "bullish"
		analysis.Confidence = math.Min(analysis.BullishScore, 95)
	} else if analysis.BearishScore > analysis.BullishScore && analysis.BearishScore > 50 {
		analysis.Direction = "bearish"
		analysis.Confidence = math.Min(analysis.BearishScore, 95)
	} else {
		log.Printf("⏸️  [%s] No clear direction (Bull: %.1f, Bear: %.1f)", 
			timeframe, analysis.BullishScore, analysis.BearishScore)
		return nil
	}
	
	// Require minimum confidence
	if analysis.Confidence < psg.MinConfidence {
		log.Printf("⏸️  [%s] Confidence too low: %.1f%% (min: %.1f%%)", 
			timeframe, analysis.Confidence, psg.MinConfidence)
		return nil
	}
	
	// ==================== PHASE 5: ENTRY CALCULATION ====================
	
	psg.calculateEntryLevels(analysis, candles)
	
	// Verify minimum RR
	if analysis.RiskReward < 2.0 {
		log.Printf("⏸️  [%s] Risk:Reward too low: %.2f:1", timeframe, analysis.RiskReward)
		return nil
	}
	
	// Determine setup type
	analysis.SetupType = psg.determineSetupType(analysis)
	
	return analysis
}

// scoreAllComponents scores all analysis components
func (psg *ProfessionalSignalGenerator) scoreAllComponents(analysis *ComprehensiveAnalysis, currentPrice float64) {
	w := psg.Weights
	
	// 1. Power of 3 Score
	if analysis.PO3 != nil {
		po3Dir, po3Str := GetPO3Signal(analysis.PO3, currentPrice)
		if po3Dir == "bullish" {
			analysis.Scores["po3"] = po3Str * w.PO3 / 100
		} else if po3Dir == "bearish" {
			analysis.Scores["po3"] = -po3Str * w.PO3 / 100
		}
		
		// Bonus for optimal entry
		if analysis.PO3.OptimalEntry {
			analysis.Scores["po3_optimal"] = w.PO3 * 0.5
		}
	}
	
	// 2. Liquidity Sweep Score
	if analysis.LiqSweep != nil {
		lsDir, lsStr := GetLiquiditySweepSignal(analysis.LiqSweep, currentPrice)
		if lsDir == "bullish" {
			analysis.Scores["liq_sweep"] = lsStr * w.LiquiditySweep / 100
		} else if lsDir == "bearish" {
			analysis.Scores["liq_sweep"] = -lsStr * w.LiquiditySweep / 100
		}
		
		// Bonus for confirmed sweep
		if analysis.LiqSweep.RecentSweep != nil && analysis.LiqSweep.RecentSweep.Confirmed {
			analysis.Scores["liq_sweep_confirmed"] = w.LiquiditySweep * 0.5
		}
	}
	
	// 3. Mirror Market Score
	if analysis.Mirror != nil {
		mirDir, mirStr := GetMirrorSignal(analysis.Mirror, currentPrice)
		if mirDir == "bullish" {
			analysis.Scores["mirror"] = mirStr * w.Mirror / 100
		} else if mirDir == "bearish" {
			analysis.Scores["mirror"] = -mirStr * w.Mirror / 100
		}
	}
	
	// 4. Institutional Score
	if analysis.Institutional != nil {
		instDir, instStr := GetInstitutionalSignal(analysis.Institutional)
		if instDir == "bullish" {
			analysis.Scores["institutional"] = instStr * w.Institutional / 100
		} else if instDir == "bearish" {
			analysis.Scores["institutional"] = -instStr * w.Institutional / 100
		}
		
		// Bonus for specific setups
		if analysis.Institutional.SilverBullet != nil && analysis.Institutional.SilverBullet.Valid {
			if analysis.Institutional.SilverBullet.Type == "bullish" {
				analysis.Scores["silver_bullet"] = w.Institutional * 0.3
			} else {
				analysis.Scores["silver_bullet"] = -w.Institutional * 0.3
			}
		}
		
		if analysis.Institutional.Unicorn != nil && analysis.Institutional.Unicorn.Valid {
			if analysis.Institutional.Unicorn.Type == "bullish" {
				analysis.Scores["unicorn"] = w.Institutional * 0.4
			} else {
				analysis.Scores["unicorn"] = -w.Institutional * 0.4
			}
		}
	}
	
	// 5. Market Maker Model Score
	if analysis.MMM != nil {
		mmmStr := GetMMMSignalStrength(analysis.MMM, "bullish")
		mmmStrBear := GetMMMSignalStrength(analysis.MMM, "bearish")
		
		if mmmStr > mmmStrBear {
			analysis.Scores["mmm"] = mmmStr * w.MMM / 100
		} else if mmmStrBear > mmmStr {
			analysis.Scores["mmm"] = -mmmStrBear * w.MMM / 100
		}
	}
	
	// 6. Supply & Demand Score
	if analysis.SD != nil {
		sdStrBull := GetSDSignalStrength(analysis.SD, "bullish", currentPrice)
		sdStrBear := GetSDSignalStrength(analysis.SD, "bearish", currentPrice)
		
		if sdStrBull > sdStrBear {
			analysis.Scores["sd"] = sdStrBull * w.SupplyDemand / 100
		} else if sdStrBear > sdStrBull {
			analysis.Scores["sd"] = -sdStrBear * w.SupplyDemand / 100
		}
	}
	
	// 7. Session Liquidity Score
	if analysis.SessionLiq != nil {
		sessDir, sessStr := GetSessionLiquiditySignal(analysis.SessionLiq, currentPrice)
		if sessDir == "bullish" {
			analysis.Scores["session"] = sessStr * w.SessionLiq / 100
		} else if sessDir == "bearish" {
			analysis.Scores["session"] = -sessStr * w.SessionLiq / 100
		}
	}
	
	// 8. Volatility Score
	if analysis.Volatility != nil {
		volScore := GetVolatilityScore(analysis.Volatility)
		analysis.Scores["volatility"] = volScore * w.Volatility / 100
	}
	
	// 9. ICT Score
	if analysis.ICT != nil {
		if analysis.ICT.Structure.Trend == "bullish" {
			analysis.Scores["ict"] = float64(analysis.ICT.Confluence) * w.ICT / 5
		} else if analysis.ICT.Structure.Trend == "bearish" {
			analysis.Scores["ict"] = -float64(analysis.ICT.Confluence) * w.ICT / 5
		}
		
		// OTE bonus
		if analysis.ICT.OTE {
			analysis.Scores["ote"] = w.ICT * 0.3
		}
	}
	
	// 10. Order Flow Score
	if analysis.OrderFlow != nil {
		if analysis.OrderFlow.Delta.Trend == "bullish" && !analysis.OrderFlow.Delta.Divergence {
			analysis.Scores["orderflow"] = analysis.OrderFlow.Delta.Strength * w.OrderFlow / 100
		} else if analysis.OrderFlow.Delta.Trend == "bearish" && !analysis.OrderFlow.Delta.Divergence {
			analysis.Scores["orderflow"] = -analysis.OrderFlow.Delta.Strength * w.OrderFlow / 100
		}
		
		// Divergence penalty
		if analysis.OrderFlow.Delta.Divergence {
			analysis.Scores["divergence"] = -w.OrderFlow * 0.5
		}
	}
	
	// 11. MTF Score (Basic)
	if analysis.MTF != nil {
		if analysis.MTF.Direction == "bullish" {
			analysis.Scores["mtf"] = float64(analysis.MTF.Confluence) * w.MTF / 3
		} else if analysis.MTF.Direction == "bearish" {
			analysis.Scores["mtf"] = -float64(analysis.MTF.Confluence) * w.MTF / 3
		}
	}
	
	// 11b. Comprehensive MTF Score (All 13 timeframes)
	if analysis.ComprehensiveMTF != nil {
		mtfScore := GetMTFConfluenceScore(analysis.ComprehensiveMTF, "bullish")
		mtfScoreBear := GetMTFConfluenceScore(analysis.ComprehensiveMTF, "bearish")
		
		if mtfScore > mtfScoreBear {
			analysis.Scores["comprehensive_mtf"] = mtfScore * w.MTF / 100 * 1.5 // 1.5x weight for comprehensive
		} else if mtfScoreBear > mtfScore {
			analysis.Scores["comprehensive_mtf"] = -mtfScoreBear * w.MTF / 100 * 1.5
		}
		
		// All timeframes aligned bonus
		if analysis.ComprehensiveMTF.AllAligned {
			if analysis.ComprehensiveMTF.Direction == "bullish" {
				analysis.Scores["mtf_all_aligned"] = 25
			} else if analysis.ComprehensiveMTF.Direction == "bearish" {
				analysis.Scores["mtf_all_aligned"] = -25
			}
		}
		
		// Higher TF aligned bonus
		if analysis.ComprehensiveMTF.HigherTFAligned {
			if analysis.ComprehensiveMTF.PositionBias == "bullish" {
				analysis.Scores["mtf_higher_aligned"] = 15
			} else if analysis.ComprehensiveMTF.PositionBias == "bearish" {
				analysis.Scores["mtf_higher_aligned"] = -15
			}
		}
	}
	
	// 12. Pattern Score
	bullishPatterns := 0
	bearishPatterns := 0
	for _, p := range analysis.Patterns {
		if p.Strength >= 75 {
			if p.Type == "bullish" {
				bullishPatterns++
			} else {
				bearishPatterns++
			}
		}
	}
	if bullishPatterns > bearishPatterns {
		analysis.Scores["patterns"] = float64(bullishPatterns) * w.Patterns / 3
	} else if bearishPatterns > bullishPatterns {
		analysis.Scores["patterns"] = -float64(bearishPatterns) * w.Patterns / 3
	}
}

// countConfluence counts the number of confluence factors
func (psg *ProfessionalSignalGenerator) countConfluence(analysis *ComprehensiveAnalysis) int {
	count := 0
	threshold := 5.0
	
	for _, score := range analysis.Scores {
		if math.Abs(score) >= threshold {
			count++
		}
	}
	
	return count
}

// calculateEntryLevels calculates entry, stop loss, and targets
func (psg *ProfessionalSignalGenerator) calculateEntryLevels(analysis *ComprehensiveAnalysis, candles []Candle) {
	currentPrice := candles[len(candles)-1].Close
	atr := calculateATR(candles[len(candles)-14:], 14)
	
	analysis.Entry = currentPrice
	
	if analysis.Direction == "bullish" {
		// Find optimal stop loss
		analysis.StopLoss = currentPrice - atr*1.5
		
		// Use S/D zone if available
		if analysis.SD != nil && analysis.SD.NearestDemand != nil {
			demandStop := analysis.SD.NearestDemand.Low - atr*0.2
			if demandStop > analysis.StopLoss && demandStop < currentPrice {
				analysis.StopLoss = demandStop
			}
		}
		
		// Use institutional setup if available
		if analysis.Institutional != nil {
			if analysis.Institutional.SilverBullet != nil && analysis.Institutional.SilverBullet.Valid {
				analysis.Entry = analysis.Institutional.SilverBullet.Entry
				analysis.StopLoss = analysis.Institutional.SilverBullet.StopLoss
			}
			if analysis.Institutional.Unicorn != nil && analysis.Institutional.Unicorn.Valid {
				analysis.Entry = analysis.Institutional.Unicorn.Entry
				analysis.StopLoss = analysis.Institutional.Unicorn.StopLoss
			}
		}
		
		// Calculate targets
		risk := analysis.Entry - analysis.StopLoss
		analysis.TP1 = analysis.Entry + risk*2.5
		analysis.TP2 = analysis.Entry + risk*4.0
		analysis.TP3 = analysis.Entry + risk*6.0
		
	} else { // bearish
		analysis.StopLoss = currentPrice + atr*1.5
		
		if analysis.SD != nil && analysis.SD.NearestSupply != nil {
			supplyStop := analysis.SD.NearestSupply.High + atr*0.2
			if supplyStop < analysis.StopLoss && supplyStop > currentPrice {
				analysis.StopLoss = supplyStop
			}
		}
		
		if analysis.Institutional != nil {
			if analysis.Institutional.SilverBullet != nil && analysis.Institutional.SilverBullet.Valid {
				analysis.Entry = analysis.Institutional.SilverBullet.Entry
				analysis.StopLoss = analysis.Institutional.SilverBullet.StopLoss
			}
			if analysis.Institutional.Unicorn != nil && analysis.Institutional.Unicorn.Valid {
				analysis.Entry = analysis.Institutional.Unicorn.Entry
				analysis.StopLoss = analysis.Institutional.Unicorn.StopLoss
			}
		}
		
		risk := analysis.StopLoss - analysis.Entry
		analysis.TP1 = analysis.Entry - risk*2.5
		analysis.TP2 = analysis.Entry - risk*4.0
		analysis.TP3 = analysis.Entry - risk*6.0
	}
	
	// Calculate risk:reward
	risk := math.Abs(analysis.Entry - analysis.StopLoss)
	reward := math.Abs(analysis.TP1 - analysis.Entry)
	analysis.RiskReward = reward / risk
}

// determineSetupType determines the primary setup type
func (psg *ProfessionalSignalGenerator) determineSetupType(analysis *ComprehensiveAnalysis) string {
	// Priority order for setup types
	if analysis.Institutional != nil {
		if analysis.Institutional.SilverBullet != nil && analysis.Institutional.SilverBullet.Valid {
			return "Silver Bullet"
		}
		if analysis.Institutional.Unicorn != nil && analysis.Institutional.Unicorn.Valid {
			return "Unicorn"
		}
		if analysis.Institutional.TurtleSoup != nil && analysis.Institutional.TurtleSoup.Valid {
			return "Turtle Soup"
		}
		if analysis.Institutional.Quasimodo != nil && analysis.Institutional.Quasimodo.Valid {
			return "Quasimodo"
		}
	}
	
	if analysis.LiqSweep != nil && analysis.LiqSweep.RecentSweep != nil && analysis.LiqSweep.RecentSweep.Confirmed {
		return "Liquidity Sweep"
	}
	
	if analysis.PO3 != nil && analysis.PO3.OptimalEntry {
		return "Power of 3"
	}
	
	if analysis.MMM != nil && analysis.MMM.InstitutionalBias != "neutral" {
		return "Market Maker Model"
	}
	
	if analysis.SD != nil && analysis.SD.InZone {
		return "Supply/Demand Zone"
	}
	
	return "Multi-Factor Confluence"
}

// GenerateProfessionalSignal generates a professional-grade signal
func (psg *ProfessionalSignalGenerator) GenerateProfessionalSignal(timeframe string) (*CreateSignalRequest, error) {
	// Fetch market data
	candles, err := psg.FetchMarketData(psg.Symbol, timeframe, 200)
	if err != nil {
		return nil, err
	}
	
	// Perform comprehensive analysis
	analysis := psg.PerformComprehensiveAnalysis(candles, timeframe)
	if analysis == nil {
		return nil, nil
	}
	
	// Create signal
	signalType := "BUY"
	if analysis.Direction == "bearish" {
		signalType = "SELL"
	}
	
	signalID := fmt.Sprintf("%d.%d", time.Now().UnixMilli(), time.Now().Nanosecond())
	killZone := psg.DetectKillZone(time.Now())
	sessionType := GetCurrentSession(time.Now())
	
	signal := &CreateSignalRequest{
		SignalID:    signalID,
		SignalType:  signalType,
		Symbol:      psg.Symbol,
		EntryPrice:  analysis.Entry,
		StopLoss:    analysis.StopLoss,
		TP1:         analysis.TP1,
		TP2:         analysis.TP2,
		TP3:         analysis.TP3,
		Strength:    int(analysis.Confidence),
		KillZone:    &killZone,
		SessionType: &sessionType,
	}
	
	// Log signal details
	log.Printf("\n🎯 [%s] PROFESSIONAL %s SIGNAL GENERATED!", timeframe, signalType)
	log.Printf("   Setup Type: %s", analysis.SetupType)
	log.Printf("   Confidence: %.1f%%", analysis.Confidence)
	log.Printf("   Confluence: %d factors", analysis.Confluence)
	log.Printf("   Entry: %.2f", analysis.Entry)
	log.Printf("   Stop Loss: %.2f", analysis.StopLoss)
	log.Printf("   TP1: %.2f (RR: %.2f:1)", analysis.TP1, analysis.RiskReward)
	log.Printf("   TP2: %.2f", analysis.TP2)
	log.Printf("   TP3: %.2f", analysis.TP3)
	log.Printf("   Session: %s", sessionType)
	log.Printf("   Scores: Bull=%.1f, Bear=%.1f\n", analysis.BullishScore, analysis.BearishScore)
	
	return signal, nil
}

// FetchMarketData fetches candlestick data
func (psg *ProfessionalSignalGenerator) FetchMarketData(symbol, interval string, limit int) ([]Candle, error) {
	sg := &SignalGenerator{}
	return sg.FetchMarketData(symbol, interval, limit)
}

// DetectKillZone detects current kill zone
func (psg *ProfessionalSignalGenerator) DetectKillZone(t time.Time) string {
	sg := &SignalGenerator{}
	return sg.DetectKillZone(t)
}

// SaveSignal saves signal to database
func (psg *ProfessionalSignalGenerator) SaveSignal(signal *CreateSignalRequest, timeframe string) error {
	sg := &SignalGenerator{}
	return sg.SaveSignal(signal, timeframe)
}

// Start begins professional signal generation
func (psg *ProfessionalSignalGenerator) Start() {
	if psg.IsRunning {
		log.Println("⚠️  Professional signal generator already running")
		return
	}
	
	psg.IsRunning = true
	log.Println("\n🏆 PROFESSIONAL INSTITUTIONAL SIGNAL GENERATOR STARTED")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("📊 Symbol: %s", psg.Symbol)
	log.Printf("⏱️  Timeframes: %v", psg.Timeframes)
	log.Printf("🎯 Min Confidence: %.1f%%", psg.MinConfidence)
	log.Printf("🔗 Min Confluence: %d factors", psg.MinConfluence)
	log.Printf("📈 Max Signals/Day: %d", psg.MaxSignalsPerDay)
	log.Println("")
	log.Println("📚 Active Concepts:")
	log.Println("   ✅ Power of 3 (AMD)")
	log.Println("   ✅ Liquidity Sweep")
	log.Println("   ✅ Mirror Market")
	log.Println("   ✅ Silver Bullet")
	log.Println("   ✅ Unicorn Setup")
	log.Println("   ✅ Turtle Soup")
	log.Println("   ✅ Quasimodo")
	log.Println("   ✅ Market Maker Model")
	log.Println("   ✅ Supply & Demand")
	log.Println("   ✅ Session Liquidity")
	log.Println("   ✅ ICT/SMC Concepts")
	log.Println("   ✅ Order Flow Analysis")
	log.Println("   ✅ Multi-Timeframe")
	log.Println("   ✅ Volatility Filter")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	// Initial scan
	psg.ScanMarkets()
	
	// Periodic scans
	ticker := time.NewTicker(psg.CheckInterval)
	go func() {
		for range ticker.C {
			if psg.IsRunning {
				psg.ScanMarkets()
			}
		}
	}()
}

// ScanMarkets scans all timeframes for signals
func (psg *ProfessionalSignalGenerator) ScanMarkets() {
	// Reset daily counter at midnight
	now := time.Now()
	if now.Hour() == 0 && now.Minute() == 0 {
		psg.SignalsToday = 0
		log.Println("🔄 Daily signal counter reset")
	}
	
	// Check daily limit
	if psg.SignalsToday >= psg.MaxSignalsPerDay {
		log.Printf("⏸️  Daily limit reached (%d/%d)", psg.SignalsToday, psg.MaxSignalsPerDay)
		return
	}
	
	log.Printf("\n🔍 Professional scan: %d timeframes...", len(psg.Timeframes))
	
	for _, tf := range psg.Timeframes {
		signal, err := psg.GenerateProfessionalSignal(tf)
		
		if err != nil {
			log.Printf("❌ [%s] Error: %v", tf, err)
			continue
		}
		
		if signal != nil {
			// Save to database
			if err := psg.SaveSignal(signal, tf); err != nil {
				log.Printf("❌ [%s] Failed to save: %v", tf, err)
				continue
			}
			
			psg.SignalsToday++
			
			if psg.SignalsToday >= psg.MaxSignalsPerDay {
				log.Println("⏸️  Daily limit reached")
				break
			}
		}
	}
}
