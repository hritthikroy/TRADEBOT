package main

import (
	"fmt"
	"log"
	"time"
)

// AIEnhancedSignalGenerator extends SignalGenerator with AI capabilities
type AIEnhancedSignalGenerator struct {
	*SignalGenerator
	GrokService       *GrokAIService
	UseAI             bool
	AIFilterEnabled   bool
	MinAIConfidence   float64
	RejectedByAI      int
	EnhancedSignals   int
}

// NewAIEnhancedSignalGenerator creates a new AI-enhanced signal generator
func NewAIEnhancedSignalGenerator() *AIEnhancedSignalGenerator {
	return &AIEnhancedSignalGenerator{
		SignalGenerator:  NewSignalGenerator(),
		GrokService:      NewGrokAIService(),
		UseAI:            true,
		AIFilterEnabled:  true,
		MinAIConfidence:  40.0, // Minimum AI confidence to proceed
		RejectedByAI:     0,
		EnhancedSignals:  0,
	}
}

// GenerateAIEnhancedSignals generates signals with AI validation
func (aisg *AIEnhancedSignalGenerator) GenerateAIEnhancedSignals() {
	// Reset daily counter at midnight
	now := time.Now()
	if now.Hour() == 0 && now.Minute() == 0 {
		aisg.SignalsToday = 0
		aisg.RejectedByAI = 0
		aisg.EnhancedSignals = 0
		log.Println("🔄 Daily counters reset")
	}

	// Check if we've hit daily limit
	if aisg.SignalsToday >= aisg.MaxSignalsPerDay {
		log.Printf("⏸️ Daily signal limit reached (%d/%d)", aisg.SignalsToday, aisg.MaxSignalsPerDay)
		return
	}

	log.Printf("\n🔍 AI-Enhanced Scan: %d timeframes...", len(aisg.Timeframes))
	log.Printf("🤖 AI Filter: %v | Min Confidence: %.0f%%\n", aisg.AIFilterEnabled, aisg.MinAIConfidence)

	// Check each timeframe
	for _, timeframe := range aisg.Timeframes {
		signal, err := aisg.AnalyzeMarket(timeframe)

		if err != nil {
			log.Printf("❌ [%s] Error analyzing market: %v", timeframe, err)
			continue
		}

		if signal == nil {
			continue
		}

		// AI Validation
		if aisg.UseAI && aisg.AIFilterEnabled {
			log.Printf("\n🤖 [%s] Running AI validation for %s signal...", timeframe, signal.SignalType)
			
			approved, reason, riskMultiplier := aisg.GrokService.ValidateSignalWithAI(signal)

			if !approved {
				aisg.RejectedByAI++
				log.Printf("❌ [%s] Signal REJECTED by AI: %s", timeframe, reason)
				log.Printf("📊 AI Stats: Rejected: %d | Enhanced: %d\n", aisg.RejectedByAI, aisg.EnhancedSignals)
				continue
			}

			log.Printf("✅ [%s] Signal APPROVED by AI: %s", timeframe, reason)
			
			// Adjust risk based on AI recommendation
			if riskMultiplier != 1.0 {
				aisg.adjustSignalRisk(signal, riskMultiplier)
				log.Printf("⚖️ Risk adjusted by %.1fx based on AI analysis", riskMultiplier)
			}

			aisg.EnhancedSignals++
		}

		// Save to database
		if err := aisg.SaveSignal(signal, timeframe); err != nil {
			log.Printf("❌ [%s] Failed to save signal: %v", timeframe, err)
			continue
		}

		aisg.SignalsToday++
		aisg.LastSignalTime = time.Now()

		log.Printf("\n✅ [%s] %s Signal Generated!", timeframe, signal.SignalType)
		log.Printf("   Entry: %.2f", signal.EntryPrice)
		log.Printf("   Stop Loss: %.2f", signal.StopLoss)
		log.Printf("   TP1: %.2f | TP2: %.2f | TP3: %.2f", signal.TP1, signal.TP2, signal.TP3)
		log.Printf("   Strength: %d%%", signal.Strength)
		log.Printf("   📊 Today: %d signals | %d AI-enhanced | %d AI-rejected\n", 
			aisg.SignalsToday, aisg.EnhancedSignals, aisg.RejectedByAI)

		// Check if we hit daily limit
		if aisg.SignalsToday >= aisg.MaxSignalsPerDay {
			log.Println("⏸️ Daily limit reached, stopping scan")
			break
		}
	}
}

// adjustSignalRisk adjusts stop loss and take profit levels based on AI risk assessment
func (aisg *AIEnhancedSignalGenerator) adjustSignalRisk(signal *CreateSignalRequest, multiplier float64) {
	entry := signal.EntryPrice
	
	if signal.SignalType == "BUY" {
		// Adjust stop loss
		slDistance := entry - signal.StopLoss
		signal.StopLoss = entry - (slDistance * multiplier)
		
		// Adjust take profits
		tp1Distance := signal.TP1 - entry
		tp2Distance := signal.TP2 - entry
		tp3Distance := signal.TP3 - entry
		
		signal.TP1 = entry + (tp1Distance * multiplier)
		signal.TP2 = entry + (tp2Distance * multiplier)
		signal.TP3 = entry + (tp3Distance * multiplier)
	} else {
		// Adjust stop loss
		slDistance := signal.StopLoss - entry
		signal.StopLoss = entry + (slDistance * multiplier)
		
		// Adjust take profits
		tp1Distance := entry - signal.TP1
		tp2Distance := entry - signal.TP2
		tp3Distance := entry - signal.TP3
		
		signal.TP1 = entry - (tp1Distance * multiplier)
		signal.TP2 = entry - (tp2Distance * multiplier)
		signal.TP3 = entry - (tp3Distance * multiplier)
	}
}

// Start begins the AI-enhanced signal generation loop
func (aisg *AIEnhancedSignalGenerator) Start() {
	if aisg.IsRunning {
		log.Println("⚠️ AI-Enhanced Signal Generator already running")
		return
	}

	aisg.IsRunning = true
	log.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("🚀 AI-ENHANCED SIGNAL GENERATOR STARTED")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("📊 Symbol: %s", aisg.Symbol)
	log.Printf("⏱️ Timeframes: %v", aisg.Timeframes)
	log.Printf("🔄 Check interval: %v", aisg.CheckInterval)
	log.Printf("📈 Min strength: %d%%", aisg.MinSignalStrength)
	log.Printf("🎯 Max signals/day: %d", aisg.MaxSignalsPerDay)
	log.Printf("🤖 AI Filter: %v", aisg.AIFilterEnabled)
	log.Printf("🎓 Min AI Confidence: %.0f%%", aisg.MinAIConfidence)
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	// Initial scan
	aisg.GenerateAIEnhancedSignals()

	// Periodic scans
	ticker := time.NewTicker(aisg.CheckInterval)
	go func() {
		for range ticker.C {
			if aisg.IsRunning {
				aisg.GenerateAIEnhancedSignals()
			}
		}
	}()
}

// GetAIStats returns AI enhancement statistics
func (aisg *AIEnhancedSignalGenerator) GetAIStats() map[string]interface{} {
	totalProcessed := aisg.EnhancedSignals + aisg.RejectedByAI
	approvalRate := 0.0
	if totalProcessed > 0 {
		approvalRate = float64(aisg.EnhancedSignals) / float64(totalProcessed) * 100
	}

	return map[string]interface{}{
		"signals_generated": aisg.SignalsToday,
		"ai_enhanced":       aisg.EnhancedSignals,
		"ai_rejected":       aisg.RejectedByAI,
		"total_processed":   totalProcessed,
		"approval_rate":     fmt.Sprintf("%.1f%%", approvalRate),
		"ai_enabled":        aisg.AIFilterEnabled,
	}
}

// ToggleAI enables or disables AI filtering
func (aisg *AIEnhancedSignalGenerator) ToggleAI(enabled bool) {
	aisg.AIFilterEnabled = enabled
	status := "disabled"
	if enabled {
		status = "enabled"
	}
	log.Printf("🤖 AI filtering %s", status)
}
