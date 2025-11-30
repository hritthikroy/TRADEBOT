package main

import (
	"database/sql"
	"math"
	"time"

	"github.com/gofiber/fiber/v2"
)

// AIAnalytics provides AI-powered insights
type AIAnalytics struct {
	OverallPerformance PerformanceMetrics      `json:"overall_performance"`
	BestConditions     []TradingCondition      `json:"best_conditions"`
	WorstConditions    []TradingCondition      `json:"worst_conditions"`
	Recommendations    []Recommendation        `json:"recommendations"`
	PredictedWinRate   float64                 `json:"predicted_win_rate"`
	OptimalSettings    OptimalSettings         `json:"optimal_settings"`
	RiskAnalysis       RiskAnalysis            `json:"risk_analysis"`
	TimeAnalysis       TimeAnalysis            `json:"time_analysis"`
}

type PerformanceMetrics struct {
	TotalSignals    int     `json:"total_signals"`
	WinRate         float64 `json:"win_rate"`
	AvgProfit       float64 `json:"avg_profit"`
	ProfitFactor    float64 `json:"profit_factor"`
	Sharpe          float64 `json:"sharpe_ratio"`
	MaxDrawdown     float64 `json:"max_drawdown"`
	ConsecutiveWins int     `json:"consecutive_wins"`
	ConsecutiveLoss int     `json:"consecutive_loss"`
}

type TradingCondition struct {
	Condition   string  `json:"condition"`
	WinRate     float64 `json:"win_rate"`
	AvgProfit   float64 `json:"avg_profit"`
	TotalTrades int     `json:"total_trades"`
	Confidence  float64 `json:"confidence"`
}

type Recommendation struct {
	Priority    string  `json:"priority"` // high, medium, low
	Category    string  `json:"category"` // risk, timing, pattern, etc
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Impact      float64 `json:"impact"` // estimated improvement %
}

type OptimalSettings struct {
	BestKillZones    []string `json:"best_kill_zones"`
	BestPatterns     []string `json:"best_patterns"`
	BestSignalType   string   `json:"best_signal_type"`
	MinStrength      float64  `json:"min_strength"`
	OptimalStopLoss  float64  `json:"optimal_stop_loss_percent"`
	OptimalTakeProfit float64 `json:"optimal_take_profit_percent"`
}

type RiskAnalysis struct {
	RiskScore          float64 `json:"risk_score"` // 0-100
	VolatilityLevel    string  `json:"volatility_level"`
	AvgRiskReward      float64 `json:"avg_risk_reward"`
	LargestLoss        float64 `json:"largest_loss"`
	LargestWin         float64 `json:"largest_win"`
	RecommendedRisk    float64 `json:"recommended_risk_percent"`
}

type TimeAnalysis struct {
	BestHours       []int   `json:"best_hours"`
	WorstHours      []int   `json:"worst_hours"`
	AvgHoldingTime  float64 `json:"avg_holding_time_minutes"`
	FastestWin      float64 `json:"fastest_win_minutes"`
	SlowestWin      float64 `json:"slowest_win_minutes"`
}

// GetAIAnalytics - Main AI analytics endpoint
func GetAIAnalytics(c *fiber.Ctx) error {
	analytics := AIAnalytics{}
	
	// Get overall performance
	perf, err := calculatePerformanceMetrics()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to calculate performance"})
	}
	analytics.OverallPerformance = perf
	
	// Find best and worst conditions
	analytics.BestConditions = findBestConditions()
	analytics.WorstConditions = findWorstConditions()
	
	// Generate recommendations
	analytics.Recommendations = generateRecommendations(perf)
	
	// Predict future win rate
	analytics.PredictedWinRate = predictWinRate()
	
	// Find optimal settings
	analytics.OptimalSettings = findOptimalSettings()
	
	// Risk analysis
	analytics.RiskAnalysis = analyzeRisk()
	
	// Time analysis
	analytics.TimeAnalysis = analyzeTimePatterns()
	
	return c.JSON(analytics)
}

// Calculate comprehensive performance metrics
func calculatePerformanceMetrics() (PerformanceMetrics, error) {
	var metrics PerformanceMetrics
	
	query := `
		SELECT 
			COUNT(*) as total,
			COALESCE(ROUND(100.0 * SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) / 
				NULLIF(SUM(CASE WHEN status IN ('win', 'loss') THEN 1 ELSE 0 END), 0), 2), 0) as win_rate,
			COALESCE(ROUND(AVG(CASE WHEN status IN ('win', 'loss') THEN profit_percent END), 2), 0) as avg_profit,
			COALESCE(SUM(CASE WHEN status = 'win' THEN profit_percent ELSE 0 END), 0) as total_wins,
			COALESCE(ABS(SUM(CASE WHEN status = 'loss' THEN profit_percent ELSE 0 END)), 1) as total_losses
		FROM trading_signals
		WHERE status IN ('win', 'loss')
		AND signal_id NOT LIKE 'test_%'
		AND signal_id NOT LIKE 'perm_test_%'
		AND symbol NOT LIKE '%TEST%'
	`
	
	var totalWins, totalLosses float64
	err := DB.QueryRow(query).Scan(
		&metrics.TotalSignals,
		&metrics.WinRate,
		&metrics.AvgProfit,
		&totalWins,
		&totalLosses,
	)
	
	if err != nil {
		return metrics, err
	}
	
	// Calculate profit factor
	if totalLosses > 0 {
		metrics.ProfitFactor = totalWins / totalLosses
	}
	
	// Calculate Sharpe ratio (simplified)
	metrics.Sharpe = calculateSharpeRatio()
	
	// Calculate max drawdown
	metrics.MaxDrawdown = calculateMaxDrawdown()
	
	// Calculate consecutive wins/losses
	metrics.ConsecutiveWins, metrics.ConsecutiveLoss = calculateStreaks()
	
	return metrics, nil
}

// Calculate Sharpe ratio
func calculateSharpeRatio() float64 {
	query := `
		SELECT COALESCE(STDDEV(profit_percent), 0), COALESCE(AVG(profit_percent), 0)
		FROM trading_signals
		WHERE status IN ('win', 'loss')
	`
	
	var stdDev, avgReturn float64
	DB.QueryRow(query).Scan(&stdDev, &avgReturn)
	
	if stdDev == 0 {
		return 0
	}
	
	// Simplified Sharpe: (avg return) / std dev
	return math.Round((avgReturn/stdDev)*100) / 100
}

// Calculate maximum drawdown
func calculateMaxDrawdown() float64 {
	query := `
		SELECT profit_percent, created_at
		FROM trading_signals
		WHERE status IN ('win', 'loss')
		AND signal_id NOT LIKE 'test_%'
		AND signal_id NOT LIKE 'perm_test_%'
		AND symbol NOT LIKE '%TEST%'
		ORDER BY created_at ASC
	`
	
	rows, err := DB.Query(query)
	if err != nil {
		return 0
	}
	defer rows.Close()
	
	var peak, maxDD float64
	var cumulative float64
	
	for rows.Next() {
		var profit float64
		var createdAt time.Time
		rows.Scan(&profit, &createdAt)
		
		cumulative += profit
		if cumulative > peak {
			peak = cumulative
		}
		
		drawdown := peak - cumulative
		if drawdown > maxDD {
			maxDD = drawdown
		}
	}
	
	return math.Round(maxDD*100) / 100
}

// Calculate consecutive win/loss streaks
func calculateStreaks() (int, int) {
	query := `
		SELECT status
		FROM trading_signals
		WHERE status IN ('win', 'loss')
		AND signal_id NOT LIKE 'test_%'
		AND signal_id NOT LIKE 'perm_test_%'
		AND symbol NOT LIKE '%TEST%'
		ORDER BY created_at DESC
	`
	
	rows, err := DB.Query(query)
	if err != nil {
		return 0, 0
	}
	defer rows.Close()
	
	var maxWins, maxLosses, currentWins, currentLosses int
	var lastStatus string
	
	for rows.Next() {
		var status string
		rows.Scan(&status)
		
		if status == "win" {
			if lastStatus == "win" {
				currentWins++
			} else {
				currentWins = 1
				if currentLosses > maxLosses {
					maxLosses = currentLosses
				}
				currentLosses = 0
			}
		} else {
			if lastStatus == "loss" {
				currentLosses++
			} else {
				currentLosses = 1
				if currentWins > maxWins {
					maxWins = currentWins
				}
				currentWins = 0
			}
		}
		
		lastStatus = status
	}
	
	if currentWins > maxWins {
		maxWins = currentWins
	}
	if currentLosses > maxLosses {
		maxLosses = currentLosses
	}
	
	return maxWins, maxLosses
}

// Find best performing conditions
func findBestConditions() []TradingCondition {
	conditions := []TradingCondition{}
	
	// Best kill zones
	query := `
		SELECT 
			COALESCE(kill_zone, 'Unknown') as condition,
			ROUND(100.0 * SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) / 
				NULLIF(COUNT(*), 0), 2) as win_rate,
			ROUND(AVG(profit_percent), 2) as avg_profit,
			COUNT(*) as total
		FROM trading_signals
		WHERE status IN ('win', 'loss')
		AND signal_id NOT LIKE 'test_%'
		AND signal_id NOT LIKE 'perm_test_%'
		AND symbol NOT LIKE '%TEST%'
		GROUP BY kill_zone
		HAVING COUNT(*) >= 3
		ORDER BY win_rate DESC, avg_profit DESC
		LIMIT 3
	`
	
	rows, _ := DB.Query(query)
	defer rows.Close()
	
	for rows.Next() {
		var tc TradingCondition
		rows.Scan(&tc.Condition, &tc.WinRate, &tc.AvgProfit, &tc.TotalTrades)
		tc.Confidence = calculateConfidence(tc.TotalTrades, tc.WinRate)
		conditions = append(conditions, tc)
	}
	
	return conditions
}

// Find worst performing conditions
func findWorstConditions() []TradingCondition {
	conditions := []TradingCondition{}
	
	query := `
		SELECT 
			COALESCE(kill_zone, 'Unknown') as condition,
			ROUND(100.0 * SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) / 
				NULLIF(COUNT(*), 0), 2) as win_rate,
			ROUND(AVG(profit_percent), 2) as avg_profit,
			COUNT(*) as total
		FROM trading_signals
		WHERE status IN ('win', 'loss')
		GROUP BY kill_zone
		HAVING COUNT(*) >= 3
		ORDER BY win_rate ASC, avg_profit ASC
		LIMIT 3
	`
	
	rows, _ := DB.Query(query)
	defer rows.Close()
	
	for rows.Next() {
		var tc TradingCondition
		rows.Scan(&tc.Condition, &tc.WinRate, &tc.AvgProfit, &tc.TotalTrades)
		tc.Confidence = calculateConfidence(tc.TotalTrades, tc.WinRate)
		conditions = append(conditions, tc)
	}
	
	return conditions
}

// Calculate confidence score based on sample size
func calculateConfidence(sampleSize int, winRate float64) float64 {
	// Wilson score interval for confidence
	if sampleSize == 0 {
		return 0
	}
	
	// Simple confidence: more samples = higher confidence
	confidence := math.Min(float64(sampleSize)/20.0, 1.0) * 100
	
	// Adjust for extreme win rates (less confident)
	if winRate > 90 || winRate < 10 {
		confidence *= 0.8
	}
	
	return math.Round(confidence*100) / 100
}

// Generate AI-powered recommendations
func generateRecommendations(perf PerformanceMetrics) []Recommendation {
	recommendations := []Recommendation{}
	
	// Win rate recommendations
	if perf.WinRate <= 50 {
		recommendations = append(recommendations, Recommendation{
			Priority:    "high",
			Category:    "strategy",
			Title:       "Improve Signal Quality",
			Description: "Win rate is at or below 50%. Increase minimum signal strength to 75%+ and only trade during high-confidence setups. Focus on quality over quantity.",
			Impact:      20.0,
		})
	}
	
	// Profit factor recommendations
	if perf.ProfitFactor < 2.0 {
		recommendations = append(recommendations, Recommendation{
			Priority:    "high",
			Category:    "risk",
			Title:       "Optimize Risk/Reward Ratio",
			Description: "Target profit factor of 2.0+. Widen take-profit targets to 2-3x your stop-loss distance. Let winners run longer.",
			Impact:      25.0,
		})
	}
	
	// Drawdown recommendations
	if perf.MaxDrawdown > 5 {
		recommendations = append(recommendations, Recommendation{
			Priority:    "high",
			Category:    "risk",
			Title:       "Reduce Maximum Drawdown",
			Description: "Drawdown is " + formatFloat(perf.MaxDrawdown) + "%. Reduce position size to 0.5-1% per trade and implement daily loss limits.",
			Impact:      30.0,
		})
	}
	
	// Consecutive loss recommendations
	if perf.ConsecutiveLoss >= 3 {
		recommendations = append(recommendations, Recommendation{
			Priority:    "medium",
			Category:    "psychology",
			Title:       "Implement Loss Limits",
			Description: "Stop trading after 3 consecutive losses. Take a break, review your strategy, and only resume when market conditions improve.",
			Impact:      15.0,
		})
	}
	
	// Small sample size warning
	if perf.TotalSignals < 30 {
		recommendations = append(recommendations, Recommendation{
			Priority:    "medium",
			Category:    "data",
			Title:       "Collect More Data",
			Description: "You have " + formatInt(perf.TotalSignals) + " signals. Collect at least 30-50 signals for statistically significant insights. Be patient and consistent.",
			Impact:      10.0,
		})
	}
	
	// Sharpe ratio recommendations
	if perf.Sharpe < 1.0 {
		recommendations = append(recommendations, Recommendation{
			Priority:    "medium",
			Category:    "risk",
			Title:       "Improve Risk-Adjusted Returns",
			Description: "Sharpe ratio is low. Focus on consistent profits with lower volatility. Avoid high-risk trades and stick to proven setups.",
			Impact:      18.0,
		})
	}
	
	// Check for missing data
	missingDataRecs := checkMissingData()
	recommendations = append(recommendations, missingDataRecs...)
	
	// Check for pattern-specific issues
	patternRecs := analyzePatternPerformance()
	recommendations = append(recommendations, patternRecs...)
	
	// Check for time-specific issues
	timeRecs := analyzeTimePerformance()
	recommendations = append(recommendations, timeRecs...)
	
	return recommendations
}

// Check for missing data that affects analysis
func checkMissingData() []Recommendation {
	recs := []Recommendation{}
	
	// Check for kill zone data
	query := `
		SELECT COUNT(*) 
		FROM trading_signals 
		WHERE kill_zone IS NULL 
		AND status IN ('win', 'loss')
		AND signal_id NOT LIKE 'test_%'
	`
	var nullKillZones int
	DB.QueryRow(query).Scan(&nullKillZones)
	
	if nullKillZones > 0 {
		recs = append(recs, Recommendation{
			Priority:    "low",
			Category:    "data",
			Title:       "Add Kill Zone Data",
			Description: "Some signals are missing kill zone information. The sync service will automatically fill this data. Refresh in a few minutes.",
			Impact:      5.0,
		})
	}
	
	// Check for pattern data
	query = `
		SELECT COUNT(*) 
		FROM trading_signals 
		WHERE pattern_type IS NULL 
		AND status IN ('win', 'loss')
		AND signal_id NOT LIKE 'test_%'
	`
	var nullPatterns int
	DB.QueryRow(query).Scan(&nullPatterns)
	
	if nullPatterns > 0 {
		recs = append(recs, Recommendation{
			Priority:    "low",
			Category:    "data",
			Title:       "Add Pattern Data",
			Description: "Some signals are missing pattern information. The sync service will automatically detect patterns. Refresh in a few minutes.",
			Impact:      5.0,
		})
	}
	
	return recs
}

// Analyze pattern performance
func analyzePatternPerformance() []Recommendation {
	recs := []Recommendation{}
	
	query := `
		SELECT 
			pattern_type,
			ROUND(100.0 * SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) / 
				NULLIF(COUNT(*), 0), 2) as win_rate,
			COUNT(*) as total
		FROM trading_signals
		WHERE status IN ('win', 'loss') AND pattern_type IS NOT NULL
		GROUP BY pattern_type
		HAVING COUNT(*) >= 5
		ORDER BY win_rate ASC
		LIMIT 2
	`
	
	rows, err := DB.Query(query)
	if err != nil {
		return recs
	}
	defer rows.Close()
	
	for rows.Next() {
		var pattern string
		var winRate float64
		var total int
		rows.Scan(&pattern, &winRate, &total)
		
		if winRate < 40 {
			recs = append(recs, Recommendation{
				Priority:    "medium",
				Category:    "pattern",
				Title:       "Avoid Low-Performing Pattern: " + pattern,
				Description: "This pattern has a " + formatFloat(winRate) + "% win rate. Consider filtering it out or requiring higher confirmation.",
				Impact:      12.0,
			})
		}
	}
	
	return recs
}

// Analyze time-based performance
func analyzeTimePerformance() []Recommendation {
	recs := []Recommendation{}
	
	query := `
		SELECT 
			EXTRACT(HOUR FROM created_at) as hour,
			ROUND(100.0 * SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) / 
				NULLIF(COUNT(*), 0), 2) as win_rate,
			COUNT(*) as total
		FROM trading_signals
		WHERE status IN ('win', 'loss')
		GROUP BY EXTRACT(HOUR FROM created_at)
		HAVING COUNT(*) >= 3
		ORDER BY win_rate ASC
		LIMIT 2
	`
	
	rows, err := DB.Query(query)
	if err != nil {
		return recs
	}
	defer rows.Close()
	
	for rows.Next() {
		var hour int
		var winRate float64
		var total int
		rows.Scan(&hour, &winRate, &total)
		
		if winRate < 35 {
			recs = append(recs, Recommendation{
				Priority:    "low",
				Category:    "timing",
				Title:       "Avoid Trading at Hour: " + formatInt(hour) + ":00",
				Description: "Low win rate during this hour. Market conditions may be unfavorable.",
				Impact:      8.0,
			})
		}
	}
	
	return recs
}

// Predict future win rate using linear regression
func predictWinRate() float64 {
	query := `
		SELECT 
			ROW_NUMBER() OVER (ORDER BY created_at) as x,
			CASE WHEN status = 'win' THEN 1.0 ELSE 0.0 END as y
		FROM trading_signals
		WHERE status IN ('win', 'loss')
		ORDER BY created_at DESC
		LIMIT 50
	`
	
	rows, err := DB.Query(query)
	if err != nil {
		return 50.0
	}
	defer rows.Close()
	
	var sumX, sumY, sumXY, sumX2 float64
	var n float64
	
	for rows.Next() {
		var x, y float64
		rows.Scan(&x, &y)
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
		n++
	}
	
	if n == 0 {
		return 50.0
	}
	
	// Linear regression: y = mx + b
	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	intercept := (sumY - slope*sumX) / n
	
	// Predict next value
	nextX := n + 1
	predicted := slope*nextX + intercept
	
	// Convert to percentage and clamp
	predictedWinRate := math.Max(0, math.Min(100, predicted*100))
	
	return math.Round(predictedWinRate*100) / 100
}

// Find optimal settings
func findOptimalSettings() OptimalSettings {
	settings := OptimalSettings{}
	
	// Best kill zones
	query := `
		SELECT kill_zone
		FROM trading_signals
		WHERE status IN ('win', 'loss')
		GROUP BY kill_zone
		HAVING COUNT(*) >= 3
		ORDER BY 
			ROUND(100.0 * SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) / COUNT(*), 2) DESC,
			AVG(profit_percent) DESC
		LIMIT 3
	`
	
	rows, _ := DB.Query(query)
	for rows.Next() {
		var kz sql.NullString
		rows.Scan(&kz)
		if kz.Valid {
			settings.BestKillZones = append(settings.BestKillZones, kz.String)
		}
	}
	rows.Close()
	
	// Best patterns
	query = `
		SELECT pattern_type
		FROM trading_signals
		WHERE status IN ('win', 'loss') AND pattern_type IS NOT NULL
		GROUP BY pattern_type
		HAVING COUNT(*) >= 3
		ORDER BY 
			ROUND(100.0 * SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) / COUNT(*), 2) DESC
		LIMIT 3
	`
	
	rows, _ = DB.Query(query)
	for rows.Next() {
		var pt sql.NullString
		rows.Scan(&pt)
		if pt.Valid {
			settings.BestPatterns = append(settings.BestPatterns, pt.String)
		}
	}
	rows.Close()
	
	// Best signal type
	query = `
		SELECT signal_type
		FROM trading_signals
		WHERE status IN ('win', 'loss')
		GROUP BY signal_type
		ORDER BY 
			ROUND(100.0 * SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) / COUNT(*), 2) DESC
		LIMIT 1
	`
	DB.QueryRow(query).Scan(&settings.BestSignalType)
	
	// Optimal strength threshold
	query = `
		SELECT ROUND(AVG(strength), 0)
		FROM trading_signals
		WHERE status = 'win'
	`
	DB.QueryRow(query).Scan(&settings.MinStrength)
	
	// Optimal stop loss
	query = `
		SELECT ROUND(AVG(ABS((stop_loss - entry_price) / entry_price * 100)), 2)
		FROM trading_signals
		WHERE status = 'win'
	`
	DB.QueryRow(query).Scan(&settings.OptimalStopLoss)
	
	// Optimal take profit
	query = `
		SELECT ROUND(AVG(ABS((tp1 - entry_price) / entry_price * 100)), 2)
		FROM trading_signals
		WHERE status = 'win'
	`
	DB.QueryRow(query).Scan(&settings.OptimalTakeProfit)
	
	return settings
}

// Analyze risk metrics
func analyzeRisk() RiskAnalysis {
	risk := RiskAnalysis{}
	
	// Calculate volatility
	query := `
		SELECT 
			COALESCE(STDDEV(profit_percent), 0) as volatility,
			COALESCE(AVG(ABS((tp1 - entry_price) / (stop_loss - entry_price))), 0) as rr_ratio,
			COALESCE(MIN(profit_percent), 0) as largest_loss,
			COALESCE(MAX(profit_percent), 0) as largest_win
		FROM trading_signals
		WHERE status IN ('win', 'loss')
	`
	
	var volatility float64
	DB.QueryRow(query).Scan(&volatility, &risk.AvgRiskReward, &risk.LargestLoss, &risk.LargestWin)
	
	// Determine volatility level
	if volatility < 2 {
		risk.VolatilityLevel = "low"
	} else if volatility < 5 {
		risk.VolatilityLevel = "medium"
	} else {
		risk.VolatilityLevel = "high"
	}
	
	// Calculate risk score (0-100, lower is better)
	riskScore := 50.0
	
	// Adjust for volatility
	riskScore += volatility * 5
	
	// Adjust for drawdown
	maxDD := calculateMaxDrawdown()
	riskScore += maxDD
	
	// Adjust for consecutive losses
	_, consecutiveLoss := calculateStreaks()
	riskScore += float64(consecutiveLoss) * 2
	
	risk.RiskScore = math.Min(100, math.Max(0, riskScore))
	
	// Recommended risk per trade
	if risk.RiskScore < 30 {
		risk.RecommendedRisk = 2.0 // Low risk = can risk more
	} else if risk.RiskScore < 60 {
		risk.RecommendedRisk = 1.0 // Medium risk
	} else {
		risk.RecommendedRisk = 0.5 // High risk = risk less
	}
	
	return risk
}

// Analyze time patterns
func analyzeTimePatterns() TimeAnalysis {
	analysis := TimeAnalysis{}
	
	// Best hours
	query := `
		SELECT EXTRACT(HOUR FROM created_at)::int
		FROM trading_signals
		WHERE status IN ('win', 'loss')
		GROUP BY EXTRACT(HOUR FROM created_at)
		HAVING COUNT(*) >= 2
		ORDER BY 
			ROUND(100.0 * SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) / COUNT(*), 2) DESC
		LIMIT 3
	`
	
	rows, _ := DB.Query(query)
	for rows.Next() {
		var hour int
		rows.Scan(&hour)
		analysis.BestHours = append(analysis.BestHours, hour)
	}
	rows.Close()
	
	// Worst hours
	query = `
		SELECT EXTRACT(HOUR FROM created_at)::int
		FROM trading_signals
		WHERE status IN ('win', 'loss')
		GROUP BY EXTRACT(HOUR FROM created_at)
		HAVING COUNT(*) >= 2
		ORDER BY 
			ROUND(100.0 * SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) / COUNT(*), 2) ASC
		LIMIT 3
	`
	
	rows, _ = DB.Query(query)
	for rows.Next() {
		var hour int
		rows.Scan(&hour)
		analysis.WorstHours = append(analysis.WorstHours, hour)
	}
	rows.Close()
	
	// Average holding time
	query = `
		SELECT COALESCE(AVG(holding_time_minutes), 0)
		FROM trading_signals
		WHERE status IN ('win', 'loss') AND holding_time_minutes IS NOT NULL
	`
	DB.QueryRow(query).Scan(&analysis.AvgHoldingTime)
	
	// Fastest win
	query = `
		SELECT COALESCE(MIN(holding_time_minutes), 0)
		FROM trading_signals
		WHERE status = 'win' AND holding_time_minutes IS NOT NULL
	`
	DB.QueryRow(query).Scan(&analysis.FastestWin)
	
	// Slowest win
	query = `
		SELECT COALESCE(MAX(holding_time_minutes), 0)
		FROM trading_signals
		WHERE status = 'win' AND holding_time_minutes IS NOT NULL
	`
	DB.QueryRow(query).Scan(&analysis.SlowestWin)
	
	return analysis
}

// Helper functions
func formatFloat(f float64) string {
	if f == float64(int(f)) {
		return formatInt(int(f))
	}
	// Format with 2 decimal places
	s := ""
	if f < 0 {
		s = "-"
		f = -f
	}
	whole := int(f)
	decimal := int((f - float64(whole)) * 100)
	return s + formatInt(whole) + "." + formatInt(decimal/10) + formatInt(decimal%10)
}

func formatInt(i int) string {
	if i < 0 {
		return "-" + formatInt(-i)
	}
	if i == 0 {
		return "0"
	}
	
	result := ""
	for i > 0 {
		result = string(rune('0'+i%10)) + result
		i /= 10
	}
	return result
}
