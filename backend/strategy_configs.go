package main

// Strategy Configuration Profiles
// Test different parameters to find the best performance

// StrategyConfig holds all strategy parameters
type StrategyConfig struct {
	Name                string
	Description         string
	MinConfluence       int
	MinStrength         int
	MinRiskReward       float64
	AllowedKillZones    []string
	AllowedSessions     []string
	RequireMTF          bool
	MinMTFAlignment     int
	AllowedPatterns     []string
	MinATRPercentile    int
	MaxATRPercentile    int
	MaxTradesPerDay     int
	MaxConsecutiveLoss  int
	DailyLossLimit      float64
}

// Strategy Profile 1: CONSERVATIVE (Maximum Accuracy)
// Goal: 75-80% win rate, fewer signals, highest quality
var ConservativeStrategy = StrategyConfig{
	Name:        "Conservative",
	Description: "Maximum accuracy - 75-80% win rate expected",
	
	// Signal Quality (STRICT)
	MinConfluence:   12, // Increase from 8
	MinStrength:     75, // Increase from 55
	MinRiskReward:   2.0, // Increase from 1.5
	
	// Kill Zone Filter (STRICT)
	AllowedKillZones: []string{
		"London Open",
		"London Close",
		"New York Open",
	},
	
	// Session Filter (STRICT)
	AllowedSessions: []string{
		"London",
		"New York",
	},
	
	// Multi-Timeframe (REQUIRED)
	RequireMTF:      true,
	MinMTFAlignment: 75, // Require 75%+ alignment
	
	// Pattern Filter (BEST ONLY)
	AllowedPatterns: []string{
		"Engulfing",
		"Pin Bar",
		"Breaker Block",
		"Order Block",
	},
	
	// Volatility Filter (MODERATE)
	MinATRPercentile: 30,
	MaxATRPercentile: 90,
	
	// Risk Management (STRICT)
	MaxTradesPerDay:    3,  // Max 3 trades per day
	MaxConsecutiveLoss: 2,  // Stop after 2 losses
	DailyLossLimit:     1.5, // 1.5% max loss per day
}

// Strategy Profile 2: BALANCED (Good Quality + Quantity)
// Goal: 68-72% win rate, moderate signals, balanced approach
var BalancedStrategy = StrategyConfig{
	Name:        "Balanced",
	Description: "Balance of quality and quantity - 68-72% win rate expected",
	
	// Signal Quality (MODERATE)
	MinConfluence:   10, // Moderate increase
	MinStrength:     70, // Good quality
	MinRiskReward:   1.8, // Better RR
	
	// Kill Zone Filter (MODERATE)
	AllowedKillZones: []string{
		"London Open",
		"London Close",
		"New York Open",
		"New York Close",
	},
	
	// Session Filter (MODERATE)
	AllowedSessions: []string{
		"London",
		"New York",
		"Overlap", // London-NY overlap
	},
	
	// Multi-Timeframe (PREFERRED)
	RequireMTF:      false, // Optional
	MinMTFAlignment: 60,    // Lower requirement
	
	// Pattern Filter (GOOD PATTERNS)
	AllowedPatterns: []string{
		"Engulfing",
		"Pin Bar",
		"Breaker Block",
		"Order Block",
		"Fair Value Gap",
		"Breakout",
	},
	
	// Volatility Filter (RELAXED)
	MinATRPercentile: 25,
	MaxATRPercentile: 95,
	
	// Risk Management (MODERATE)
	MaxTradesPerDay:    5,  // Max 5 trades per day
	MaxConsecutiveLoss: 3,  // Stop after 3 losses
	DailyLossLimit:     2.0, // 2% max loss per day
}

// Strategy Profile 3: AGGRESSIVE (Current Settings)
// Goal: 61-65% win rate, more signals, more opportunities
var AggressiveStrategy = StrategyConfig{
	Name:        "Aggressive",
	Description: "More signals, current performance - 61-65% win rate expected",
	
	// Signal Quality (CURRENT)
	MinConfluence:   8,   // Current setting
	MinStrength:     55,  // Current setting
	MinRiskReward:   1.5, // Current setting
	
	// Kill Zone Filter (RELAXED)
	AllowedKillZones: []string{
		"Asian Session",
		"London Open",
		"London Close",
		"New York Open",
		"New York Close",
		"Off Hours",
	},
	
	// Session Filter (ALL)
	AllowedSessions: []string{
		"Asian",
		"London",
		"New York",
		"Overlap",
		"Off Hours",
	},
	
	// Multi-Timeframe (OPTIONAL)
	RequireMTF:      false,
	MinMTFAlignment: 50, // Low requirement
	
	// Pattern Filter (ALL PATTERNS)
	AllowedPatterns: []string{
		"Engulfing",
		"Pin Bar",
		"Breaker Block",
		"Order Block",
		"Fair Value Gap",
		"Breakout",
		"Trend Following",
		"Continuation",
		"Inside Bar",
	},
	
	// Volatility Filter (WIDE)
	MinATRPercentile: 10,
	MaxATRPercentile: 100,
	
	// Risk Management (RELAXED)
	MaxTradesPerDay:    10, // Max 10 trades per day
	MaxConsecutiveLoss: 5,  // Stop after 5 losses
	DailyLossLimit:     3.0, // 3% max loss per day
}

// GetActiveStrategy returns the currently active strategy
// Change this to switch between strategies
func GetActiveStrategy() StrategyConfig {
	// CHANGE THIS LINE TO SWITCH STRATEGIES:
	// return ConservativeStrategy  // For maximum accuracy
	// return BalancedStrategy      // For balance
	return AggressiveStrategy       // For more signals (current)
}

// ValidateSignal checks if a signal meets the strategy requirements
func (config *StrategyConfig) ValidateSignal(signal map[string]interface{}) bool {
	// Check confluence
	if confluence, ok := signal["confluence"].(int); ok {
		if confluence < config.MinConfluence {
			return false
		}
	}
	
	// Check strength
	if strength, ok := signal["strength"].(int); ok {
		if strength < config.MinStrength {
			return false
		}
	}
	
	// Check risk-reward
	if rr, ok := signal["risk_reward"].(float64); ok {
		if rr < config.MinRiskReward {
			return false
		}
	}
	
	// Check kill zone
	if killZone, ok := signal["kill_zone"].(string); ok {
		if !contains(config.AllowedKillZones, killZone) {
			return false
		}
	}
	
	// Check session
	if session, ok := signal["session_type"].(string); ok {
		if !contains(config.AllowedSessions, session) {
			return false
		}
	}
	
	// Check MTF alignment
	if config.RequireMTF {
		if mtf, ok := signal["mtf_alignment"].(int); ok {
			if mtf < config.MinMTFAlignment {
				return false
			}
		} else {
			return false // MTF required but not present
		}
	}
	
	// Check pattern
	if pattern, ok := signal["pattern_type"].(string); ok {
		if !contains(config.AllowedPatterns, pattern) {
			return false
		}
	}
	
	// All checks passed
	return true
}

// Helper function to check if slice contains string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// GetStrategyStats returns statistics for the active strategy
func GetStrategyStats() map[string]interface{} {
	config := GetActiveStrategy()
	
	return map[string]interface{}{
		"name":                config.Name,
		"description":         config.Description,
		"min_confluence":      config.MinConfluence,
		"min_strength":        config.MinStrength,
		"min_risk_reward":     config.MinRiskReward,
		"allowed_kill_zones":  config.AllowedKillZones,
		"allowed_sessions":    config.AllowedSessions,
		"require_mtf":         config.RequireMTF,
		"min_mtf_alignment":   config.MinMTFAlignment,
		"allowed_patterns":    config.AllowedPatterns,
		"max_trades_per_day":  config.MaxTradesPerDay,
		"max_consecutive_loss": config.MaxConsecutiveLoss,
		"daily_loss_limit":    config.DailyLossLimit,
	}
}
