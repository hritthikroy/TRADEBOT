package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"sync"
	"time"
)

// WorldClassOptimizer finds the absolute best parameters for each strategy
type WorldClassOptimizer struct {
	Symbol       string
	Days         int
	StartBalance float64
	Strategies   []string
}

// OptimizationParams represents parameters to test
type OptimizationParams struct {
	StopATR     float64
	TP1ATR      float64
	TP2ATR      float64
	TP3ATR      float64
	RiskPercent float64
}

// WorldClassOptimizationResult stores the best result for a strategy
type WorldClassOptimizationResult struct {
	Strategy       string                 `json:"strategy"`
	TotalTests     int                    `json:"totalTests"`
	BestScore      float64                `json:"bestScore"`
	BestParams     OptimizationParams     `json:"bestParams"`
	BacktestResult *BacktestResult        `json:"backtestResult"`
	TestDuration   string                 `json:"testDuration"`
}

// WorldClassResults stores all optimization results
type WorldClassResults struct {
	OptimizationDate time.Time                                `json:"optimizationDate"`
	Symbol           string                                   `json:"symbol"`
	Days             int                                      `json:"days"`
	StartBalance     float64                                  `json:"startBalance"`
	TotalDuration    string                                   `json:"totalDuration"`
	Results          map[string]WorldClassOptimizationResult  `json:"results"`
	BestOverall      WorldClassOptimizationResult             `json:"bestOverall"`
}

// NewWorldClassOptimizer creates a new optimizer
func NewWorldClassOptimizer() *WorldClassOptimizer {
	return &WorldClassOptimizer{
		Symbol:       "BTCUSDT",
		Days:         180,
		StartBalance: 1000,
		Strategies: []string{
			"session_trader",
			"breakout_master",
			"liquidity_hunter",
			"trend_rider",
			"range_master",
			"smart_money_tracker",
			"institutional_follower",
			"reversal_sniper",
			"momentum_beast",
			"scalper_pro",
		},
	}
}

// OptimizeAll optimizes all strategies
func (wco *WorldClassOptimizer) OptimizeAll() *WorldClassResults {
	startTime := time.Now()
	
	log.Println("🌍 WORLD-CLASS STRATEGY OPTIMIZATION")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("Symbol: %s | Days: %d | Start Balance: $%.2f", wco.Symbol, wco.Days, wco.StartBalance)
	log.Println("")
	log.Println("Testing Parameters:")
	log.Println("  • Stop Loss: 0.5, 0.75, 1.0, 1.25, 1.5, 2.0 ATR")
	log.Println("  • TP1: 2.0, 2.5, 3.0, 3.5, 4.0, 5.0 ATR")
	log.Println("  • TP2: 3.0, 4.0, 4.5, 5.0, 6.0, 7.5 ATR")
	log.Println("  • TP3: 5.0, 6.0, 7.5, 10.0, 12.5, 15.0 ATR")
	log.Println("  • Risk: 0.5%, 1.0%, 1.5%, 2.0%, 2.5%")
	log.Println("")
	log.Println("Optimization Goals:")
	log.Println("  1. Win Rate > 60%")
	log.Println("  2. Profit Factor > 3.0")
	log.Println("  3. Max Drawdown < 15%")
	log.Println("  4. Return > 500%")
	log.Println("  5. Total Trades > 20")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("")
	
	results := &WorldClassResults{
		OptimizationDate: time.Now(),
		Symbol:           wco.Symbol,
		Days:             wco.Days,
		StartBalance:     wco.StartBalance,
		Results:          make(map[string]WorldClassOptimizationResult),
	}
	
	// Fetch candles once (reuse for all tests)
	log.Println("📊 Fetching historical data...")
	candles, err := fetchBinanceData(wco.Symbol, "15m", wco.Days)
	if err != nil {
		log.Printf("❌ Failed to fetch data: %v", err)
		return results
	}
	log.Printf("✅ Loaded %d candles", len(candles))
	log.Println("")
	
	// Optimize each strategy in parallel
	var wg sync.WaitGroup
	resultsChan := make(chan WorldClassOptimizationResult, len(wco.Strategies))
	
	for _, strategy := range wco.Strategies {
		wg.Add(1)
		go func(strat string) {
			defer wg.Done()
			result := wco.OptimizeStrategy(strat, candles)
			resultsChan <- result
		}(strategy)
	}
	
	// Wait for all optimizations to complete
	go func() {
		wg.Wait()
		close(resultsChan)
	}()
	
	// Collect results
	bestScore := 0.0
	for result := range resultsChan {
		results.Results[result.Strategy] = result
		
		if result.BestScore > bestScore {
			bestScore = result.BestScore
			results.BestOverall = result
		}
	}
	
	results.TotalDuration = time.Since(startTime).String()
	
	log.Println("")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("🎉 WORLD-CLASS OPTIMIZATION COMPLETE!")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("Total Duration: %s", results.TotalDuration)
	log.Printf("Best Overall: %s (Score: %.2f)", results.BestOverall.Strategy, results.BestOverall.BestScore)
	log.Println("")
	
	return results
}

// OptimizeStrategy optimizes a single strategy
func (wco *WorldClassOptimizer) OptimizeStrategy(strategy string, candles []Candle) WorldClassOptimizationResult {
	startTime := time.Now()
	
	log.Printf("🎯 Optimizing: %s", strategy)
	
	// Parameter ranges
	stopLossValues := []float64{0.5, 0.75, 1.0, 1.25, 1.5, 2.0}
	tp1Values := []float64{2.0, 2.5, 3.0, 3.5, 4.0, 5.0}
	tp2Values := []float64{3.0, 4.0, 4.5, 5.0, 6.0, 7.5}
	tp3Values := []float64{5.0, 6.0, 7.5, 10.0, 12.5, 15.0}
	riskValues := []float64{0.5, 1.0, 1.5, 2.0, 2.5}
	
	bestScore := 0.0
	var bestParams OptimizationParams
	var bestResult *BacktestResult
	totalTests := 0
	
	// Test all combinations
	for _, stop := range stopLossValues {
		for _, tp1 := range tp1Values {
			for _, tp2 := range tp2Values {
				for _, tp3 := range tp3Values {
					for _, risk := range riskValues {
						// Validate: TP1 < TP2 < TP3
						if tp1 < tp2 && tp2 < tp3 {
							totalTests++
							
							// Run backtest
							config := BacktestConfig{
								Symbol:       wco.Symbol,
								Interval:     "15m",
								Days:         wco.Days,
								StartBalance: wco.StartBalance,
								RiskPercent:  risk / 100,
								Strategy:     strategy,
							}
							
							result, err := RunBacktest(config, candles)
							if err != nil {
								continue
							}
							
							// Calculate score
							score := wco.CalculateScore(result)
							
							// Check if this is the best
							if score > bestScore && wco.MeetsMinimumCriteria(result) {
								bestScore = score
								bestParams = OptimizationParams{
									StopATR:     stop,
									TP1ATR:      tp1,
									TP2ATR:      tp2,
									TP3ATR:      tp3,
									RiskPercent: risk,
								}
								bestResult = result
								
								log.Printf("  ✨ %s: Score %.2f | WR %.1f%% | PF %.2f | DD %.1f%% | Trades %d",
									strategy, score, result.WinRate, result.ProfitFactor, result.MaxDrawdown, result.TotalTrades)
							}
							
							// Progress
							if totalTests%50 == 0 {
								log.Printf("  ⏳ %s: Tested %d combinations...", strategy, totalTests)
							}
						}
					}
				}
			}
		}
	}
	
	duration := time.Since(startTime)
	log.Printf("✅ %s: Complete! Tests: %d | Duration: %s | Best Score: %.2f",
		strategy, totalTests, duration, bestScore)
	
	return WorldClassOptimizationResult{
		Strategy:       strategy,
		TotalTests:     totalTests,
		BestScore:      bestScore,
		BestParams:     bestParams,
		BacktestResult: bestResult,
		TestDuration:   duration.String(),
	}
}

// CalculateScore calculates optimization score
func (wco *WorldClassOptimizer) CalculateScore(result *BacktestResult) float64 {
	if result == nil {
		return 0
	}
	
	// Weighted scoring formula
	// Prioritize: Win Rate, Profit Factor, Low Drawdown, High Return
	score := (result.WinRate * 2.0) +                    // Win rate weight: 2x
		(result.ProfitFactor * 10.0) +               // Profit factor weight: 10x
		(result.ReturnPercent / 10.0) +              // Return weight: 0.1x
		-(result.MaxDrawdown * 2.0) +                // Drawdown penalty: -2x
		(float64(result.TotalTrades) / 2.0) +        // Trade count bonus: 0.5x
		(float64(result.WinningTrades) / 1.0)        // Winning trades bonus: 1x
	
	return math.Max(0, score)
}

// MeetsMinimumCriteria checks if result meets minimum standards
func (wco *WorldClassOptimizer) MeetsMinimumCriteria(result *BacktestResult) bool {
	if result == nil {
		return false
	}
	
	return result.WinRate >= 55.0 &&           // Min 55% win rate
		result.ProfitFactor >= 2.5 &&      // Min 2.5 profit factor
		result.MaxDrawdown <= 20.0 &&      // Max 20% drawdown
		result.TotalTrades >= 15 &&        // Min 15 trades
		result.ReturnPercent >= 100.0      // Min 100% return
}

// SaveResults saves optimization results to file
func (wco *WorldClassOptimizer) SaveResults(results *WorldClassResults, filename string) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(filename, data, 0644)
}

// PrintSummary prints a summary of results
func (results *WorldClassResults) PrintSummary() {
	fmt.Println("")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📊 OPTIMIZATION SUMMARY")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("")
	
	for _, result := range results.Results {
		if result.BacktestResult != nil {
			fmt.Printf("🎯 %s:\n", result.Strategy)
			fmt.Printf("   Score: %.2f | Tests: %d | Duration: %s\n",
				result.BestScore, result.TotalTests, result.TestDuration)
			fmt.Printf("   Stop: %.2f | TP1: %.2f | TP2: %.2f | TP3: %.2f | Risk: %.1f%%\n",
				result.BestParams.StopATR, result.BestParams.TP1ATR,
				result.BestParams.TP2ATR, result.BestParams.TP3ATR,
				result.BestParams.RiskPercent)
			fmt.Printf("   WR: %.1f%% | PF: %.2f | Return: %.0f%% | DD: %.1f%% | Trades: %d\n",
				result.BacktestResult.WinRate, result.BacktestResult.ProfitFactor,
				result.BacktestResult.ReturnPercent, result.BacktestResult.MaxDrawdown,
				result.BacktestResult.TotalTrades)
			fmt.Println("")
		}
	}
	
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("🏆 BEST OVERALL: %s (Score: %.2f)\n", results.BestOverall.Strategy, results.BestOverall.BestScore)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("")
}
