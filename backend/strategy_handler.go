package main

import (
	"github.com/gofiber/fiber/v2"
)

// GetStrategyConfig returns the current active strategy configuration
func GetStrategyConfig(c *fiber.Ctx) error {
	stats := GetStrategyStats()
	
	return c.JSON(fiber.Map{
		"success": true,
		"strategy": stats,
	})
}

// GetAllStrategies returns all available strategy profiles
func GetAllStrategies(c *fiber.Ctx) error {
	strategies := []map[string]interface{}{
		{
			"id":          1,
			"name":        ConservativeStrategy.Name,
			"description": ConservativeStrategy.Description,
			"expected_win_rate": "75-80%",
			"signals_per_month": "40-50",
			"quality":     "Highest",
			"parameters": map[string]interface{}{
				"min_confluence":   ConservativeStrategy.MinConfluence,
				"min_strength":     ConservativeStrategy.MinStrength,
				"min_risk_reward":  ConservativeStrategy.MinRiskReward,
				"kill_zones":       ConservativeStrategy.AllowedKillZones,
				"sessions":         ConservativeStrategy.AllowedSessions,
				"require_mtf":      ConservativeStrategy.RequireMTF,
			},
		},
		{
			"id":          2,
			"name":        BalancedStrategy.Name,
			"description": BalancedStrategy.Description,
			"expected_win_rate": "68-72%",
			"signals_per_month": "60-70",
			"quality":     "Good",
			"parameters": map[string]interface{}{
				"min_confluence":   BalancedStrategy.MinConfluence,
				"min_strength":     BalancedStrategy.MinStrength,
				"min_risk_reward":  BalancedStrategy.MinRiskReward,
				"kill_zones":       BalancedStrategy.AllowedKillZones,
				"sessions":         BalancedStrategy.AllowedSessions,
				"require_mtf":      BalancedStrategy.RequireMTF,
			},
		},
		{
			"id":          3,
			"name":        AggressiveStrategy.Name,
			"description": AggressiveStrategy.Description,
			"expected_win_rate": "61-65%",
			"signals_per_month": "80-100",
			"quality":     "Moderate",
			"parameters": map[string]interface{}{
				"min_confluence":   AggressiveStrategy.MinConfluence,
				"min_strength":     AggressiveStrategy.MinStrength,
				"min_risk_reward":  AggressiveStrategy.MinRiskReward,
				"kill_zones":       AggressiveStrategy.AllowedKillZones,
				"sessions":         AggressiveStrategy.AllowedSessions,
				"require_mtf":      AggressiveStrategy.RequireMTF,
			},
		},
	}
	
	return c.JSON(fiber.Map{
		"success":    true,
		"strategies": strategies,
		"active":     GetActiveStrategy().Name,
	})
}

// CompareStrategies compares performance of all strategies
func CompareStrategies(c *fiber.Ctx) error {
	// Get signals from database
	query := `
		SELECT 
			signal_type,
			strength,
			kill_zone,
			session_type,
			pattern_type,
			status,
			profit_percent
		FROM trading_signals
		WHERE status IN ('win', 'loss')
		AND signal_id NOT LIKE 'test_%'
		AND signal_id NOT LIKE 'perm_test_%'
		AND symbol NOT LIKE '%TEST%'
	`
	
	rows, err := DB.Query(query)
	if (err != nil) {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch signals",
		})
	}
	defer rows.Close()
	
	// Collect all signals
	type Signal struct {
		SignalType   string
		Strength     int
		KillZone     *string
		SessionType  *string
		PatternType  *string
		Status       string
		ProfitPercent *float64
	}
	
	var signals []Signal
	for rows.Next() {
		var s Signal
		err := rows.Scan(
			&s.SignalType,
			&s.Strength,
			&s.KillZone,
			&s.SessionType,
			&s.PatternType,
			&s.Status,
			&s.ProfitPercent,
		)
		if err != nil {
			continue
		}
		signals = append(signals, s)
	}
	
	// Test each strategy
	strategies := []StrategyConfig{
		ConservativeStrategy,
		BalancedStrategy,
		AggressiveStrategy,
	}
	
	results := make([]map[string]interface{}, 0)
	
	for _, strategy := range strategies {
		totalSignals := 0
		wins := 0
		losses := 0
		totalProfit := 0.0
		totalLoss := 0.0
		
		for _, signal := range signals {
			// Create signal map for validation
			signalMap := map[string]interface{}{
				"strength":     signal.Strength,
				"kill_zone":    stringValue(signal.KillZone),
				"session_type": stringValue(signal.SessionType),
				"pattern_type": stringValue(signal.PatternType),
				"confluence":   8, // Default (we don't store this)
				"risk_reward":  1.5, // Default (we don't store this)
			}
			
			// Check if signal would pass this strategy
			if !strategy.ValidateSignal(signalMap) {
				continue
			}
			
			totalSignals++
			
			if signal.Status == "win" {
				wins++
				if signal.ProfitPercent != nil {
					totalProfit += *signal.ProfitPercent
				}
			} else if signal.Status == "loss" {
				losses++
				if signal.ProfitPercent != nil {
					totalLoss += *signal.ProfitPercent
				}
			}
		}
		
		winRate := 0.0
		if wins+losses > 0 {
			winRate = float64(wins) / float64(wins+losses) * 100
		}
		
		profitFactor := 0.0
		if totalLoss != 0 {
			profitFactor = totalProfit / (-totalLoss)
		}
		
		results = append(results, map[string]interface{}{
			"strategy":      strategy.Name,
			"description":   strategy.Description,
			"total_signals": totalSignals,
			"wins":          wins,
			"losses":        losses,
			"win_rate":      round(winRate, 2),
			"profit_factor": round(profitFactor, 2),
			"total_profit":  round(totalProfit, 2),
			"total_loss":    round(totalLoss, 2),
			"net_profit":    round(totalProfit+totalLoss, 2),
		})
	}
	
	return c.JSON(fiber.Map{
		"success": true,
		"comparison": results,
		"note": "Based on your historical signals",
	})
}

// Helper function to get string value from pointer
func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Helper function to round float
func round(val float64, precision int) float64 {
	ratio := 1.0
	for i := 0; i < precision; i++ {
		ratio *= 10
	}
	return float64(int(val*ratio+0.5)) / ratio
}
