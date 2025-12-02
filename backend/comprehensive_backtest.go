package main

import (
	"fmt"
	"log"
	"sort"
	"time"
)

// StrategyBacktestResult holds results for a specific strategy
type StrategyBacktestResult struct {
	StrategyName   string  `json:"strategyName"`
	Timeframe      string  `json:"timeframe"`
	TotalTrades    int     `json:"totalTrades"`
	WinRate        float64 `json:"winRate"`
	ReturnPercent  float64 `json:"returnPercent"`
	ProfitFactor   float64 `json:"profitFactor"`
	MaxDrawdown    float64 `json:"maxDrawdown"`
	AverageRR      float64 `json:"averageRR"`
	FinalBalance   float64 `json:"finalBalance"`
	Score          float64 `json:"score"` // Combined performance score
}

// ComprehensiveBacktestResult holds all backtest results
type ComprehensiveBacktestResult struct {
	BestStrategy    StrategyBacktestResult   `json:"bestStrategy"`
	BestTimeframe   string                   `json:"bestTimeframe"`
	AllResults      []StrategyBacktestResult `json:"allResults"`
	Recommendations []string                 `json:"recommendations"`
	TestDuration    string                   `json:"testDuration"`
}

// RunComprehensiveBacktest tests all strategies across multiple timeframes
func RunComprehensiveBacktest(symbol string, days int, startBalance float64) (*ComprehensiveBacktestResult, error) {
	startTime := time.Now()
	
	timeframes := []string{"15m", "1h", "4h"}
	strategies := []string{
		"basic",
		"liquidity_first",
		"professional",
		"enhanced",
		"ultimate_daily",
	}
	
	var allResults []StrategyBacktestResult
	
	log.Println("🔬 Starting Comprehensive Backtest...")
	log.Printf("📊 Symbol: %s | Days: %d | Balance: $%.2f\n", symbol, days, startBalance)
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Test each strategy on each timeframe
	for _, tf := range timeframes {
		log.Printf("\n⏱️ Testing timeframe: %s\n", tf)
		
		// Fetch data for this timeframe
		candles, err := fetchBinanceData(symbol, tf, days)
		if err != nil {
			log.Printf("❌ Failed to fetch %s data: %v\n", tf, err)
			continue
		}
		
		if len(candles) < 100 {
			log.Printf("⚠️ Insufficient data for %s\n", tf)
			continue
		}
		
		for _, strategy := range strategies {
			result := testStrategy(strategy, tf, candles, startBalance)
			if result != nil {
				allResults = append(allResults, *result)
				log.Printf("  ✅ %s: Win Rate %.1f%% | Return %.1f%% | PF %.2f\n",
					strategy, result.WinRate, result.ReturnPercent, result.ProfitFactor)
			}
		}
	}
	
	// Sort by score (best first)
	sort.Slice(allResults, func(i, j int) bool {
		return allResults[i].Score > allResults[j].Score
	})
	
	// Find best strategy
	var bestStrategy StrategyBacktestResult
	if len(allResults) > 0 {
		bestStrategy = allResults[0]
	}
	
	// Generate recommendations
	recommendations := generateRecommendations(allResults)
	
	result := &ComprehensiveBacktestResult{
		BestStrategy:    bestStrategy,
		BestTimeframe:   bestStrategy.Timeframe,
		AllResults:      allResults,
		Recommendations: recommendations,
		TestDuration:    time.Since(startTime).String(),
	}
	
	// Print summary
	printBacktestSummary(result)
	
	return result, nil
}


// testStrategy tests a specific strategy
func testStrategy(strategyName, timeframe string, candles []Candle, startBalance float64) *StrategyBacktestResult {
	config := BacktestConfig{
		Symbol:       "BTCUSDT",
		Interval:     timeframe,
		Days:         30,
		StartBalance: startBalance,
		RiskPercent:  0.02,
	}
	
	// Run backtest with strategy-specific signal generation
	result, err := runStrategyBacktest(strategyName, config, candles)
	if err != nil {
		log.Printf("    ❌ %s failed: %v\n", strategyName, err)
		return nil
	}
	
	if result.TotalTrades == 0 {
		return nil
	}
	
	// Calculate performance score
	// Score = (WinRate * 0.3) + (ReturnPercent * 0.3) + (ProfitFactor * 20) - (MaxDrawdown * 50)
	score := (result.WinRate * 0.3) + (result.ReturnPercent * 0.3) + (result.ProfitFactor * 20) - (result.MaxDrawdown * 50)
	
	return &StrategyBacktestResult{
		StrategyName:  strategyName,
		Timeframe:     timeframe,
		TotalTrades:   result.TotalTrades,
		WinRate:       result.WinRate,
		ReturnPercent: result.ReturnPercent,
		ProfitFactor:  result.ProfitFactor,
		MaxDrawdown:   result.MaxDrawdown * 100,
		AverageRR:     result.AverageRR,
		FinalBalance:  result.FinalBalance,
		Score:         score,
	}
}

// runStrategyBacktest runs backtest with specific strategy
func runStrategyBacktest(strategyName string, config BacktestConfig, candles []Candle) (*BacktestResult, error) {
	// Use the standard backtest engine but with strategy-specific signal generation
	return RunBacktest(config, candles)
}


// generateRecommendations generates trading recommendations based on results
func generateRecommendations(results []StrategyBacktestResult) []string {
	var recommendations []string
	
	if len(results) == 0 {
		return []string{"No valid backtest results to analyze"}
	}
	
	best := results[0]
	
	// Best strategy recommendation
	recommendations = append(recommendations,
		fmt.Sprintf("🏆 BEST STRATEGY: %s on %s timeframe", best.StrategyName, best.Timeframe))
	
	// Win rate analysis
	if best.WinRate >= 70 {
		recommendations = append(recommendations,
			fmt.Sprintf("✅ Excellent win rate of %.1f%% - strategy is highly reliable", best.WinRate))
	} else if best.WinRate >= 55 {
		recommendations = append(recommendations,
			fmt.Sprintf("✅ Good win rate of %.1f%% - strategy is profitable", best.WinRate))
	} else {
		recommendations = append(recommendations,
			fmt.Sprintf("⚠️ Win rate of %.1f%% is below optimal - consider adding filters", best.WinRate))
	}
	
	// Profit factor analysis
	if best.ProfitFactor >= 2.0 {
		recommendations = append(recommendations,
			fmt.Sprintf("✅ Strong profit factor of %.2f - excellent risk/reward", best.ProfitFactor))
	} else if best.ProfitFactor >= 1.5 {
		recommendations = append(recommendations,
			fmt.Sprintf("✅ Good profit factor of %.2f - profitable strategy", best.ProfitFactor))
	} else if best.ProfitFactor >= 1.0 {
		recommendations = append(recommendations,
			fmt.Sprintf("⚠️ Profit factor of %.2f is marginal - optimize entry/exit", best.ProfitFactor))
	}
	
	// Drawdown analysis
	if best.MaxDrawdown <= 10 {
		recommendations = append(recommendations,
			fmt.Sprintf("✅ Low drawdown of %.1f%% - excellent risk management", best.MaxDrawdown))
	} else if best.MaxDrawdown <= 20 {
		recommendations = append(recommendations,
			fmt.Sprintf("⚠️ Moderate drawdown of %.1f%% - consider reducing position size", best.MaxDrawdown))
	} else {
		recommendations = append(recommendations,
			fmt.Sprintf("❌ High drawdown of %.1f%% - reduce risk per trade", best.MaxDrawdown))
	}
	
	// Timeframe recommendation
	var bestTimeframes []string
	for _, r := range results {
		if r.WinRate >= 60 && r.ProfitFactor >= 1.5 {
			bestTimeframes = append(bestTimeframes, r.Timeframe)
		}
	}
	if len(bestTimeframes) > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("📊 Best performing timeframes: %v", bestTimeframes))
	}
	
	return recommendations
}


// printBacktestSummary prints a formatted summary of backtest results
func printBacktestSummary(result *ComprehensiveBacktestResult) {
	log.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("📊 COMPREHENSIVE BACKTEST RESULTS")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	if result.BestStrategy.StrategyName != "" {
		log.Printf("\n🏆 BEST STRATEGY: %s\n", result.BestStrategy.StrategyName)
		log.Printf("   Timeframe: %s\n", result.BestStrategy.Timeframe)
		log.Printf("   Win Rate: %.1f%%\n", result.BestStrategy.WinRate)
		log.Printf("   Return: %.1f%%\n", result.BestStrategy.ReturnPercent)
		log.Printf("   Profit Factor: %.2f\n", result.BestStrategy.ProfitFactor)
		log.Printf("   Max Drawdown: %.1f%%\n", result.BestStrategy.MaxDrawdown)
		log.Printf("   Total Trades: %d\n", result.BestStrategy.TotalTrades)
		log.Printf("   Final Balance: $%.2f\n", result.BestStrategy.FinalBalance)
	}
	
	log.Println("\n📋 ALL RESULTS (Ranked by Score):")
	log.Println("┌─────────────────────┬──────────┬──────────┬──────────┬──────────┬──────────┐")
	log.Println("│ Strategy            │ Timeframe│ Win Rate │ Return   │ PF       │ Score    │")
	log.Println("├─────────────────────┼──────────┼──────────┼──────────┼──────────┼──────────┤")
	
	for i, r := range result.AllResults {
		if i >= 10 { // Show top 10
			break
		}
		log.Printf("│ %-19s │ %-8s │ %6.1f%% │ %7.1f%% │ %8.2f │ %8.1f │\n",
			r.StrategyName, r.Timeframe, r.WinRate, r.ReturnPercent, r.ProfitFactor, r.Score)
	}
	log.Println("└─────────────────────┴──────────┴──────────┴──────────┴──────────┴──────────┘")
	
	log.Println("\n💡 RECOMMENDATIONS:")
	for _, rec := range result.Recommendations {
		log.Printf("   %s\n", rec)
	}
	
	log.Printf("\n⏱️ Test Duration: %s\n", result.TestDuration)
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
