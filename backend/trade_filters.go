package main

import (
	"fmt"
	
	"github.com/gofiber/fiber/v2"
)

// TradeFilter represents advanced filtering criteria
type TradeFilter struct {
	MinWinRate        float64 `json:"min_win_rate"`
	MinProfitFactor   float64 `json:"min_profit_factor"`
	MinConfidence     float64 `json:"min_confidence"`
	MinSampleSize     int     `json:"min_sample_size"`
	BestHoursOnly     bool    `json:"best_hours_only"`
	AvoidWorstHours   bool    `json:"avoid_worst_hours"`
	HighPerformOnly   bool    `json:"high_perform_only"`
}

// TradeOpportunity represents a filtered trade opportunity
type TradeOpportunity struct {
	Condition       string  `json:"condition"`
	Type            string  `json:"type"` // kill_zone, pattern, hour, signal_type
	WinRate         float64 `json:"win_rate"`
	AvgProfit       float64 `json:"avg_profit"`
	ProfitFactor    float64 `json:"profit_factor"`
	TotalTrades     int     `json:"total_trades"`
	Confidence      float64 `json:"confidence"`
	Recommendation  string  `json:"recommendation"`
	Score           float64 `json:"score"`
}

// GetBestTradeOpportunities - Find best trading opportunities
func GetBestTradeOpportunities(c *fiber.Ctx) error {
	opportunities := []TradeOpportunity{}
	
	// Best Kill Zones
	killZones := findBestKillZones()
	opportunities = append(opportunities, killZones...)
	
	// Best Patterns
	patterns := findBestPatterns()
	opportunities = append(opportunities, patterns...)
	
	// Best Hours
	hours := findBestHours()
	opportunities = append(opportunities, hours...)
	
	// Best Signal Types
	signalTypes := findBestSignalTypes()
	opportunities = append(opportunities, signalTypes...)
	
	// Sort by score
	sortByScore(opportunities)
	
	return c.JSON(fiber.Map{
		"opportunities": opportunities,
		"total": len(opportunities),
		"recommendation": generateTradeRecommendation(opportunities),
	})
}

// Find best kill zones
func findBestKillZones() []TradeOpportunity {
	opportunities := []TradeOpportunity{}
	
	query := `
		SELECT 
			COALESCE(kill_zone, 'Unknown') as condition,
			COUNT(*) as total,
			SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) as wins,
			ROUND(100.0 * SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) / COUNT(*), 2) as win_rate,
			ROUND(AVG(profit_percent), 2) as avg_profit,
			ROUND(SUM(CASE WHEN status = 'win' THEN profit_percent ELSE 0 END) / 
				NULLIF(ABS(SUM(CASE WHEN status = 'loss' THEN profit_percent ELSE 0 END)), 0), 2) as profit_factor
		FROM trading_signals
		WHERE status IN ('win', 'loss')
		AND created_at >= NOW() - INTERVAL '30 days'
		AND signal_id NOT LIKE 'test_%'
		AND signal_id NOT LIKE 'perm_test_%'
		AND symbol NOT LIKE '%TEST%'
		GROUP BY kill_zone
		HAVING COUNT(*) >= 5
		ORDER BY win_rate DESC, avg_profit DESC
	`
	
	rows, err := DB.Query(query)
	if err != nil {
		return opportunities
	}
	defer rows.Close()
	
	for rows.Next() {
		var opp TradeOpportunity
		var wins int
		rows.Scan(&opp.Condition, &opp.TotalTrades, &wins, &opp.WinRate, &opp.AvgProfit, &opp.ProfitFactor)
		
		opp.Type = "kill_zone"
		opp.Confidence = calculateConfidence(opp.TotalTrades, opp.WinRate)
		opp.Score = calculateOpportunityScore(opp)
		opp.Recommendation = generateRecommendation(opp)
		
		opportunities = append(opportunities, opp)
	}
	
	return opportunities
}

// Find best patterns
func findBestPatterns() []TradeOpportunity {
	opportunities := []TradeOpportunity{}
	
	query := `
		SELECT 
			pattern_type as condition,
			COUNT(*) as total,
			SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) as wins,
			ROUND(100.0 * SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) / COUNT(*), 2) as win_rate,
			ROUND(AVG(profit_percent), 2) as avg_profit,
			ROUND(SUM(CASE WHEN status = 'win' THEN profit_percent ELSE 0 END) / 
				NULLIF(ABS(SUM(CASE WHEN status = 'loss' THEN profit_percent ELSE 0 END)), 0), 2) as profit_factor
		FROM trading_signals
		WHERE status IN ('win', 'loss') AND pattern_type IS NOT NULL
		AND created_at >= NOW() - INTERVAL '30 days'
		GROUP BY pattern_type
		HAVING COUNT(*) >= 5
		ORDER BY win_rate DESC, avg_profit DESC
	`
	
	rows, err := DB.Query(query)
	if err != nil {
		return opportunities
	}
	defer rows.Close()
	
	for rows.Next() {
		var opp TradeOpportunity
		var wins int
		rows.Scan(&opp.Condition, &opp.TotalTrades, &wins, &opp.WinRate, &opp.AvgProfit, &opp.ProfitFactor)
		
		opp.Type = "pattern"
		opp.Confidence = calculateConfidence(opp.TotalTrades, opp.WinRate)
		opp.Score = calculateOpportunityScore(opp)
		opp.Recommendation = generateRecommendation(opp)
		
		opportunities = append(opportunities, opp)
	}
	
	return opportunities
}

// Find best hours
func findBestHours() []TradeOpportunity {
	opportunities := []TradeOpportunity{}
	
	query := `
		SELECT 
			EXTRACT(HOUR FROM created_at)::text || ':00' as condition,
			COUNT(*) as total,
			SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) as wins,
			ROUND(100.0 * SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) / COUNT(*), 2) as win_rate,
			ROUND(AVG(profit_percent), 2) as avg_profit,
			ROUND(SUM(CASE WHEN status = 'win' THEN profit_percent ELSE 0 END) / 
				NULLIF(ABS(SUM(CASE WHEN status = 'loss' THEN profit_percent ELSE 0 END)), 0), 2) as profit_factor
		FROM trading_signals
		WHERE status IN ('win', 'loss')
		AND created_at >= NOW() - INTERVAL '30 days'
		GROUP BY EXTRACT(HOUR FROM created_at)
		HAVING COUNT(*) >= 3
		ORDER BY win_rate DESC, avg_profit DESC
	`
	
	rows, err := DB.Query(query)
	if err != nil {
		return opportunities
	}
	defer rows.Close()
	
	for rows.Next() {
		var opp TradeOpportunity
		var wins int
		rows.Scan(&opp.Condition, &opp.TotalTrades, &wins, &opp.WinRate, &opp.AvgProfit, &opp.ProfitFactor)
		
		opp.Type = "hour"
		opp.Confidence = calculateConfidence(opp.TotalTrades, opp.WinRate)
		opp.Score = calculateOpportunityScore(opp)
		opp.Recommendation = generateRecommendation(opp)
		
		opportunities = append(opportunities, opp)
	}
	
	return opportunities
}

// Find best signal types
func findBestSignalTypes() []TradeOpportunity {
	opportunities := []TradeOpportunity{}
	
	query := `
		SELECT 
			signal_type as condition,
			COUNT(*) as total,
			SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) as wins,
			ROUND(100.0 * SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) / COUNT(*), 2) as win_rate,
			ROUND(AVG(profit_percent), 2) as avg_profit,
			ROUND(SUM(CASE WHEN status = 'win' THEN profit_percent ELSE 0 END) / 
				NULLIF(ABS(SUM(CASE WHEN status = 'loss' THEN profit_percent ELSE 0 END)), 0), 2) as profit_factor
		FROM trading_signals
		WHERE status IN ('win', 'loss')
		AND created_at >= NOW() - INTERVAL '30 days'
		GROUP BY signal_type
		HAVING COUNT(*) >= 5
		ORDER BY win_rate DESC, avg_profit DESC
	`
	
	rows, err := DB.Query(query)
	if err != nil {
		return opportunities
	}
	defer rows.Close()
	
	for rows.Next() {
		var opp TradeOpportunity
		var wins int
		rows.Scan(&opp.Condition, &opp.TotalTrades, &wins, &opp.WinRate, &opp.AvgProfit, &opp.ProfitFactor)
		
		opp.Type = "signal_type"
		opp.Confidence = calculateConfidence(opp.TotalTrades, opp.WinRate)
		opp.Score = calculateOpportunityScore(opp)
		opp.Recommendation = generateRecommendation(opp)
		
		opportunities = append(opportunities, opp)
	}
	
	return opportunities
}


// Calculate opportunity score (0-100)
func calculateOpportunityScore(opp TradeOpportunity) float64 {
	score := 0.0
	
	// Win rate weight (40%)
	score += (opp.WinRate / 100.0) * 40.0
	
	// Profit factor weight (30%)
	if opp.ProfitFactor > 0 {
		profitScore := (opp.ProfitFactor / 3.0) * 30.0
		if profitScore > 30 {
			profitScore = 30
		}
		score += profitScore
	}
	
	// Average profit weight (20%)
	if opp.AvgProfit > 0 {
		profitPercent := (opp.AvgProfit / 5.0) * 20.0
		if profitPercent > 20 {
			profitPercent = 20
		}
		score += profitPercent
	}
	
	// Confidence weight (10%)
	score += (opp.Confidence / 100.0) * 10.0
	
	return score
}

// Generate recommendation text
func generateRecommendation(opp TradeOpportunity) string {
	if opp.Score >= 80 {
		return "🟢 EXCELLENT - Highly recommended"
	} else if opp.Score >= 65 {
		return "🟡 GOOD - Recommended"
	} else if opp.Score >= 50 {
		return "🟠 MODERATE - Use with caution"
	} else {
		return "🔴 POOR - Avoid"
	}
}

// Sort opportunities by score
func sortByScore(opportunities []TradeOpportunity) {
	for i := 0; i < len(opportunities)-1; i++ {
		for j := i + 1; j < len(opportunities); j++ {
			if opportunities[j].Score > opportunities[i].Score {
				opportunities[i], opportunities[j] = opportunities[j], opportunities[i]
			}
		}
	}
}

// Generate overall trade recommendation
func generateTradeRecommendation(opportunities []TradeOpportunity) string {
	if len(opportunities) == 0 {
		return "Not enough data to generate recommendations. Need at least 5 completed signals."
	}
	
	excellent := 0
	good := 0
	
	for _, opp := range opportunities {
		if opp.Score >= 80 {
			excellent++
		} else if opp.Score >= 65 {
			good++
		}
	}
	
	if excellent > 0 {
		return "Found " + string(rune(excellent+'0')) + " excellent trading opportunities. Focus on these for best results."
	} else if good > 0 {
		return "Found " + string(rune(good+'0')) + " good trading opportunities. These have solid win rates."
	}
	
	return "Current conditions show moderate performance. Wait for better setups or adjust strategy."
}

// GetTradeRules - Get specific trading rules based on analysis
func GetTradeRules(c *fiber.Ctx) error {
	rules := generateTradeRules()
	
	return c.JSON(fiber.Map{
		"rules": rules,
		"total": len(rules),
	})
}

type TradeRule struct {
	Rule       string  `json:"rule"`
	Reason     string  `json:"reason"`
	Impact     string  `json:"impact"`
	Priority   string  `json:"priority"`
	Confidence float64 `json:"confidence"`
}

func generateTradeRules() []TradeRule {
	rules := []TradeRule{}
	
	// Rule 1: Best kill zones (last 30 days)
	query := `
		SELECT kill_zone, 
			ROUND(100.0 * SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) / COUNT(*), 2) as win_rate,
			COUNT(*) as total
		FROM trading_signals
		WHERE status IN ('win', 'loss') AND kill_zone IS NOT NULL
		AND created_at >= NOW() - INTERVAL '30 days'
		GROUP BY kill_zone
		HAVING COUNT(*) >= 5
		ORDER BY win_rate DESC
		LIMIT 1
	`
	
	var bestKZ string
	var winRate float64
	var total int
	err := DB.QueryRow(query).Scan(&bestKZ, &winRate, &total)
	if err == nil && winRate > 60 {
		rules = append(rules, TradeRule{
			Rule:       "Trade only during " + bestKZ,
			Reason:     "This kill zone has " + formatFloat(winRate) + "% win rate",
			Impact:     "High",
			Priority:   "High",
			Confidence: calculateConfidence(total, winRate),
		})
	}
	
	// Rule 2: Avoid worst kill zones (last 30 days)
	query = `
		SELECT kill_zone, 
			ROUND(100.0 * SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) / COUNT(*), 2) as win_rate,
			COUNT(*) as total
		FROM trading_signals
		WHERE status IN ('win', 'loss') AND kill_zone IS NOT NULL
		AND created_at >= NOW() - INTERVAL '30 days'
		GROUP BY kill_zone
		HAVING COUNT(*) >= 5
		ORDER BY win_rate ASC
		LIMIT 1
	`
	
	var worstKZ string
	err = DB.QueryRow(query).Scan(&worstKZ, &winRate, &total)
	if err == nil && winRate < 40 {
		rules = append(rules, TradeRule{
			Rule:       "Avoid trading during " + worstKZ,
			Reason:     "This kill zone has only " + formatFloat(winRate) + "% win rate",
			Impact:     "High",
			Priority:   "High",
			Confidence: calculateConfidence(total, winRate),
		})
	}
	
	// Rule 3: Minimum signal strength (last 30 days)
	query = `
		SELECT ROUND(AVG(strength), 0)
		FROM trading_signals
		WHERE status = 'win'
		AND created_at >= NOW() - INTERVAL '30 days'
	`
	
	var avgStrength float64
	err = DB.QueryRow(query).Scan(&avgStrength)
	if err == nil && avgStrength > 0 {
		rules = append(rules, TradeRule{
			Rule:       "Only take signals with strength >= " + formatFloat(avgStrength) + "%",
			Reason:     "Winning signals average " + formatFloat(avgStrength) + "% strength",
			Impact:     "Medium",
			Priority:   "Medium",
			Confidence: 85.0,
		})
	}
	
	// Rule 4: Best signal type (last 30 days)
	query = `
		SELECT signal_type,
			ROUND(100.0 * SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) / COUNT(*), 2) as win_rate,
			COUNT(*) as total
		FROM trading_signals
		WHERE status IN ('win', 'loss')
		AND created_at >= NOW() - INTERVAL '30 days'
		GROUP BY signal_type
		ORDER BY win_rate DESC
		LIMIT 1
	`
	
	var bestType string
	err = DB.QueryRow(query).Scan(&bestType, &winRate, &total)
	if err == nil && total >= 5 {
		if winRate > 55 {
			rules = append(rules, TradeRule{
				Rule:       "Prefer " + bestType + " signals",
				Reason:     bestType + " signals have " + formatFloat(winRate) + "% win rate",
				Impact:     "Medium",
				Priority:   "Medium",
				Confidence: calculateConfidence(total, winRate),
			})
		}
	}
	
	// Rule 5: Best hours (last 30 days)
	query = `
		SELECT EXTRACT(HOUR FROM created_at)::int as hour,
			ROUND(100.0 * SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) / COUNT(*), 2) as win_rate,
			COUNT(*) as total
		FROM trading_signals
		WHERE status IN ('win', 'loss')
		AND created_at >= NOW() - INTERVAL '30 days'
		GROUP BY EXTRACT(HOUR FROM created_at)
		HAVING COUNT(*) >= 3
		ORDER BY win_rate DESC
		LIMIT 1
	`
	
	var bestHour int
	err = DB.QueryRow(query).Scan(&bestHour, &winRate, &total)
	if err == nil && winRate > 60 {
		rules = append(rules, TradeRule{
			Rule:       "Best trading hour: " + formatInt(bestHour) + ":00",
			Reason:     "This hour has " + formatFloat(winRate) + "% win rate",
			Impact:     "Low",
			Priority:   "Low",
			Confidence: calculateConfidence(total, winRate),
		})
	}
	
	return rules
}

// GetSmartSignalFilter - Returns filter criteria for better trades
func GetSmartSignalFilter(c *fiber.Ctx) error {
	filter := SmartFilter{}
	
	// Calculate optimal strength threshold (last 30 days)
	query := `
		SELECT ROUND(AVG(strength), 0)
		FROM trading_signals
		WHERE status = 'win'
		AND created_at >= NOW() - INTERVAL '30 days'
	`
	DB.QueryRow(query).Scan(&filter.MinStrength)
	
	// Get best kill zones (last 30 days)
	query = `
		SELECT kill_zone
		FROM trading_signals
		WHERE status IN ('win', 'loss') AND kill_zone IS NOT NULL
		AND created_at >= NOW() - INTERVAL '30 days'
		GROUP BY kill_zone
		HAVING COUNT(*) >= 3
		ORDER BY ROUND(100.0 * SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) / COUNT(*), 2) DESC
		LIMIT 3
	`
	rows, _ := DB.Query(query)
	for rows.Next() {
		var kz string
		rows.Scan(&kz)
		filter.AllowedKillZones = append(filter.AllowedKillZones, kz)
	}
	rows.Close()
	
	// Get best patterns (last 30 days)
	query = `
		SELECT pattern_type
		FROM trading_signals
		WHERE status IN ('win', 'loss') AND pattern_type IS NOT NULL
		AND created_at >= NOW() - INTERVAL '30 days'
		GROUP BY pattern_type
		HAVING COUNT(*) >= 3
		ORDER BY ROUND(100.0 * SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) / COUNT(*), 2) DESC
		LIMIT 3
	`
	rows, _ = DB.Query(query)
	for rows.Next() {
		var pt string
		rows.Scan(&pt)
		filter.AllowedPatterns = append(filter.AllowedPatterns, pt)
	}
	rows.Close()
	
	// Get best hours (last 30 days)
	query = `
		SELECT EXTRACT(HOUR FROM created_at)::int
		FROM trading_signals
		WHERE status IN ('win', 'loss')
		AND created_at >= NOW() - INTERVAL '30 days'
		GROUP BY EXTRACT(HOUR FROM created_at)
		HAVING COUNT(*) >= 2
		ORDER BY ROUND(100.0 * SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) / COUNT(*), 2) DESC
		LIMIT 5
	`
	rows, _ = DB.Query(query)
	for rows.Next() {
		var hour int
		rows.Scan(&hour)
		filter.BestHours = append(filter.BestHours, hour)
	}
	rows.Close()
	
	// Calculate min win rate (last 30 days)
	query = `
		SELECT ROUND(AVG(win_rate), 2)
		FROM (
			SELECT ROUND(100.0 * SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) / COUNT(*), 2) as win_rate
			FROM trading_signals
			WHERE status IN ('win', 'loss')
			AND created_at >= NOW() - INTERVAL '30 days'
			GROUP BY kill_zone
		) as subquery
	`
	DB.QueryRow(query).Scan(&filter.MinWinRate)
	
	return c.JSON(filter)
}

type SmartFilter struct {
	MinStrength       float64  `json:"min_strength"`
	MinWinRate        float64  `json:"min_win_rate"`
	AllowedKillZones  []string `json:"allowed_kill_zones"`
	AllowedPatterns   []string `json:"allowed_patterns"`
	BestHours         []int    `json:"best_hours"`
}


// calculateConfidence calculates confidence based on sample size and win rate
func calculateConfidence(sampleSize int, winRate float64) float64 {
	// Base confidence on sample size
	sizeConfidence := 0.0
	if sampleSize >= 50 {
		sizeConfidence = 100.0
	} else if sampleSize >= 20 {
		sizeConfidence = 80.0
	} else if sampleSize >= 10 {
		sizeConfidence = 60.0
	} else if sampleSize >= 5 {
		sizeConfidence = 40.0
	} else {
		sizeConfidence = 20.0
	}
	
	// Adjust by win rate deviation from 50%
	winRateDeviation := (winRate - 50.0) / 50.0
	adjustment := winRateDeviation * 20.0
	
	confidence := sizeConfidence + adjustment
	
	// Clamp between 0 and 100
	if confidence > 100 {
		confidence = 100
	} else if confidence < 0 {
		confidence = 0
	}
	
	return confidence
}

// formatFloat formats float to string with 2 decimals
func formatFloat(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

// formatInt formats int to string
func formatInt(i int) string {
	return fmt.Sprintf("%d", i)
}
