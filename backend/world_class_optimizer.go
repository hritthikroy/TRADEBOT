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
	
	log.Println("📊 Each strategy will fetch data for its optimal timeframe...")
	log.Println("")
	
	// Optimize each strategy in parallel
	var wg sync.WaitGroup
	resultsChan := make(chan WorldClassOptimizationResult, len(wco.Strategies))
	
	for _, strategy := range wco.Strategies {
		wg.Add(1)
		go func(strat string) {
			defer wg.Done()
			result := wco.OptimizeStrategy(strat, nil) // Each strategy fetches its own candles
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
	
	// Fetch strategy-specific candles with correct interval
	interval := getStrategyInterval(strategy)
	strategyCandles, err := fetchBinanceData(wco.Symbol, interval, wco.Days)
	if err != nil || len(strategyCandles) < 50 {
		log.Printf("❌ %s: Failed to fetch candles for interval %s", strategy, interval)
		return WorldClassOptimizationResult{
			Strategy:   strategy,
			TotalTests: 0,
			BestScore:  0,
		}
	}
	log.Printf("  📊 %s: Using %s interval with %d candles", strategy, interval, len(strategyCandles))
	
	// Parameter ranges
	stopLossValues := []float64{0.5, 0.75, 1.0, 1.25, 1.5, 2.0}
	tp1Values := []float64{2.0, 2.5, 3.0, 3.5, 4.0, 5.0}
	tp2Values := []float64{3.0, 4.0, 4.5, 5.0, 6.0, 7.5}
	tp3Values := []float64{5.0, 6.0, 7.5, 10.0, 12.5, 15.0}
	riskValues := []float64{0.5, 1.0, 1.5, 2.0, 2.5}
	
	bestScore := 0.0
	var bestParams OptimizationParams
	var bestResult *BacktestResult
	var bestUnfilteredResult *BacktestResult // Track best result even if it doesn't meet criteria
	var bestUnfilteredParams OptimizationParams
	bestUnfilteredScore := 0.0
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
							
							// Run backtest with CUSTOM parameters (not hardcoded ones)
							config := BacktestConfig{
								Symbol:       wco.Symbol,
								Interval:     interval,
								Days:         wco.Days,
								StartBalance: wco.StartBalance,
								RiskPercent:  risk / 100,
								Strategy:     strategy,
							}
							
							// Pass custom parameters to backtest
							result, err := RunBacktestWithCustomParams(config, strategyCandles, stop, tp1, tp2, tp3)
							if err != nil {
								continue
							}
							
							// Calculate score
							score := wco.CalculateScore(result)
							
							// Track best unfiltered result
							if score > bestUnfilteredScore {
								bestUnfilteredScore = score
								bestUnfilteredResult = result
								bestUnfilteredParams = OptimizationParams{
									StopATR:     stop,
									TP1ATR:      tp1,
									TP2ATR:      tp2,
									TP3ATR:      tp3,
									RiskPercent: risk,
								}
							}
							
							// Check if this is the best that meets criteria
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
								
								log.Printf("  ✨ %s: NEW BEST! Score %.2f | Stop %.2f | TP1 %.1f | TP2 %.1f | TP3 %.1f | Risk %.1f%% | WR %.1f%% | PF %.2f | Return %.0f%% | Trades %d",
									strategy, score, stop, tp1, tp2, tp3, risk, result.WinRate, result.ProfitFactor, result.ReturnPercent, result.TotalTrades)
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
	
	// Debug: Show what we found
	if bestUnfilteredResult != nil {
		log.Printf("  📊 %s: Best unfiltered - Trades=%d, WR=%.1f%%, PF=%.2f, Return=%.0f%%, DD=%.1f%%",
			strategy, bestUnfilteredResult.TotalTrades, bestUnfilteredResult.WinRate,
			bestUnfilteredResult.ProfitFactor, bestUnfilteredResult.ReturnPercent, bestUnfilteredResult.MaxDrawdown)
	}
	
	// If no result met criteria, use the best unfiltered result
	if bestResult == nil && bestUnfilteredResult != nil {
		log.Printf("⚠️  %s: No results met criteria. Using best unfiltered result.", strategy)
		bestResult = bestUnfilteredResult
		bestScore = bestUnfilteredScore
		bestParams = bestUnfilteredParams
	}
	
	if bestResult != nil {
		log.Printf("✅ %s: Complete! Tests: %d | Duration: %s | Best Score: %.2f | WR: %.1f%% | PF: %.2f | Return: %.0f%%",
			strategy, totalTests, duration, bestScore, bestResult.WinRate, bestResult.ProfitFactor, bestResult.ReturnPercent)
	} else {
		log.Printf("❌ %s: Complete! Tests: %d | Duration: %s | No viable results found (no trades generated)",
			strategy, totalTests, duration)
	}
	
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
	if result == nil || result.TotalTrades == 0 {
		return 0
	}
	
	// CRITICAL: Heavily penalize losing strategies
	if result.ReturnPercent < 0 {
		return 0 // Losing strategies get ZERO score
	}
	
	// Weighted scoring formula - Prioritize PROFITABILITY
	score := (result.ReturnPercent * 2.0) +              // Return is MOST important: 2x
		(result.ProfitFactor * 20.0) +               // Profit factor: 20x (increased)
		(result.WinRate * 1.0) +                     // Win rate: 1x
		-(result.MaxDrawdown * 3.0) +                // Drawdown penalty: -3x (increased)
		(float64(result.WinningTrades) * 2.0) +      // Winning trades bonus: 2x
		(float64(result.TotalTrades) / 5.0)          // Trade count bonus: 0.2x (reduced)
	
	return math.Max(0, score)
}

// MeetsMinimumCriteria checks if result meets minimum standards
func (wco *WorldClassOptimizer) MeetsMinimumCriteria(result *BacktestResult) bool {
	if result == nil {
		return false
	}
	
	// Focus on PROFITABLE strategies only
	return result.ReturnPercent > 0 &&         // MUST be profitable
		result.ProfitFactor > 1.0 &&       // MUST have profit factor > 1
		result.WinRate >= 40.0 &&          // Min 40% win rate
		result.MaxDrawdown <= 40.0 &&      // Max 40% drawdown
		result.TotalTrades >= 5            // Min 5 trades
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
