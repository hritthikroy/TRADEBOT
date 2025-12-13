package optimization

import (
	"fmt"
	"log"
	"math"
	"sort"
)

// ParameterSet represents a set of strategy parameters
type ParameterSet struct {
	MinConfluence int
	StopATR       float64
	TP1ATR        float64
	TP2ATR        float64
	TP3ATR        float64
	RiskPercent   float64
}

// OptimizationResult holds optimization results
type OptimizationResult struct {
	StrategyName  string       `json:"strategyName"`
	Parameters    ParameterSet `json:"parameters"`
	WinRate       float64      `json:"winRate"`
	ProfitFactor  float64      `json:"profitFactor"`
	ReturnPercent float64      `json:"returnPercent"`
	TotalTrades   int          `json:"totalTrades"`
	Score         float64      `json:"score"`
}

// OptimizeStrategyParameters finds best parameters for a strategy
func OptimizeStrategyParameters(strategyName string, symbol string, startBalance float64, days int) ([]OptimizationResult, error) {
	log.Printf("🔬 Optimizing parameters for: %s", strategyName)
	
	strategies := GetAdvancedStrategies()
	strategy, exists := strategies[strategyName]
	if !exists {
		return nil, fmt.Errorf("strategy not found: %s", strategyName)
	}
	
	// Fetch data
	candles, err := fetchBinanceData(symbol, strategy.Timeframe, days)
	if err != nil {
		return nil, err
	}
	
	if len(candles) < 100 {
		return nil, fmt.Errorf("insufficient data")
	}
	
	results := []OptimizationResult{}
	
	// Parameter ranges to test
	confluenceLevels := []int{4, 5, 6, 7, 8}
	stopATRs := []float64{0.5, 1.0, 1.5, 2.0}
	tp1ATRs := []float64{2.0, 3.0, 4.0, 5.0}
	riskPercents := []float64{1.0, 1.5, 2.0, 2.5}
	
	totalTests := len(confluenceLevels) * len(stopATRs) * len(tp1ATRs) * len(riskPercents)
	testCount := 0
	
	log.Printf("   Testing %d parameter combinations...", totalTests)
	
	for _, minConf := range confluenceLevels {
		for _, stopATR := range stopATRs {
			for _, tp1ATR := range tp1ATRs {
				for _, riskPct := range riskPercents {
					testCount++
					
					// Calculate TP2 and TP3 based on TP1
					tp2ATR := tp1ATR * 1.5
					tp3ATR := tp1ATR * 2.5
					
					params := ParameterSet{
						MinConfluence: minConf,
						StopATR:       stopATR,
						TP1ATR:        tp1ATR,
						TP2ATR:        tp2ATR,
						TP3ATR:        tp3ATR,
						RiskPercent:   riskPct,
					}
					
					// Test with these parameters
					result := testStrategyWithParameters(strategyName, strategy, candles, startBalance, params)
					
					if result.TotalTrades >= 5 { // Minimum trades for valid test
						results = append(results, result)
					}
					
					if testCount%10 == 0 {
						log.Printf("   Progress: %d/%d tests completed", testCount, totalTests)
					}
				}
			}
		}
	}
	
	// Sort by score
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	
	log.Printf("   ✅ Found %d valid parameter sets", len(results))
	
	if len(results) > 0 {
		best := results[0]
		log.Printf("   🏆 Best Parameters:")
		log.Printf("      Min Confluence: %d", best.Parameters.MinConfluence)
		log.Printf("      Stop ATR: %.1f", best.Parameters.StopATR)
		log.Printf("      TP1 ATR: %.1f", best.Parameters.TP1ATR)
		log.Printf("      Risk: %.1f%%", best.Parameters.RiskPercent)
		log.Printf("      Win Rate: %.1f%%", best.WinRate)
		log.Printf("      Return: %.1f%%", best.ReturnPercent)
		log.Printf("      Profit Factor: %.2f", best.ProfitFactor)
		log.Printf("      Score: %.1f", best.Score)
	}
	
	return results, nil
}

// testStrategyWithParameters tests strategy with specific parameters
func testStrategyWithParameters(strategyName string, strategy AdvancedStrategy, candles []Candle, startBalance float64, params ParameterSet) OptimizationResult {
	result := OptimizationResult{
		StrategyName: strategyName,
		Parameters:   params,
	}
	
	// Generate signals with custom confluence
	signals := []AdvancedSignal{}
	for i := 100; i < len(candles)-1; i++ {
		confluence := 0
		reasons := []string{}
		
		for _, concept := range strategy.RequiredConcepts {
			if checkConcept(candles, i, concept) {
				confluence++
				reasons = append(reasons, concept)
			}
		}
		
		if confluence < params.MinConfluence {
			continue
		}
		
		signalType := determineSignalTypeAdvanced(candles, i, strategyName)
		if signalType == "" {
			continue
		}
		
		atr := calculateATR(candles, i)
		entry := candles[i].Close
		
		var stopLoss, tp1, tp2, tp3 float64
		if signalType == "BUY" {
			stopLoss = entry - (atr * params.StopATR)
			tp1 = entry + (atr * params.TP1ATR)
			tp2 = entry + (atr * params.TP2ATR)
			tp3 = entry + (atr * params.TP3ATR)
		} else {
			stopLoss = entry + (atr * params.StopATR)
			tp1 = entry - (atr * params.TP1ATR)
			tp2 = entry - (atr * params.TP2ATR)
			tp3 = entry - (atr * params.TP3ATR)
		}
		
		signal := AdvancedSignal{
			Type:       signalType,
			Entry:      entry,
			StopLoss:   stopLoss,
			TP1:        tp1,
			TP2:        tp2,
			TP3:        tp3,
			Confluence: confluence,
			Reasons:    reasons,
		}
		
		signals = append(signals, signal)
	}
	
	if len(signals) == 0 {
		return result
	}
	
	// Simulate trades with custom risk
	balance := startBalance
	maxBalance := balance
	totalProfit := 0.0
	totalLoss := 0.0
	winningTrades := 0
	losingTrades := 0
	
	for _, signal := range signals {
		riskAmount := balance * (params.RiskPercent / 100.0)
		riskPerUnit := math.Abs(signal.Entry - signal.StopLoss)
		
		if riskPerUnit == 0 {
			continue
		}
		
		positionSize := riskAmount / riskPerUnit
		if positionSize <= 0 {
			continue
		}
		
		// Find exit
		exitFound := false
		var profit float64
		
		for i := 0; i < len(candles); i++ {
			if signal.Type == "BUY" {
				if candles[i].Low <= signal.StopLoss {
					profit = (signal.StopLoss - signal.Entry) * positionSize
					exitFound = true
					break
				}
				if candles[i].High >= signal.TP3 {
					profit = (signal.TP3 - signal.Entry) * positionSize
					exitFound = true
					break
				} else if candles[i].High >= signal.TP2 {
					profit = (signal.TP2 - signal.Entry) * positionSize
					exitFound = true
					break
				} else if candles[i].High >= signal.TP1 {
					profit = (signal.TP1 - signal.Entry) * positionSize
					exitFound = true
					break
				}
			} else {
				if candles[i].High >= signal.StopLoss {
					profit = (signal.Entry - signal.StopLoss) * positionSize
					exitFound = true
					break
				}
				if candles[i].Low <= signal.TP3 {
					profit = (signal.Entry - signal.TP3) * positionSize
					exitFound = true
					break
				} else if candles[i].Low <= signal.TP2 {
					profit = (signal.Entry - signal.TP2) * positionSize
					exitFound = true
					break
				} else if candles[i].Low <= signal.TP1 {
					profit = (signal.Entry - signal.TP1) * positionSize
					exitFound = true
					break
				}
			}
		}
		
		if !exitFound {
			continue
		}
		
		balance += profit
		if balance > maxBalance {
			maxBalance = balance
		}
		
		result.TotalTrades++
		if profit > 0 {
			winningTrades++
			totalProfit += profit
		} else {
			losingTrades++
			totalLoss += profit
		}
	}
	
	// Calculate metrics
	if result.TotalTrades > 0 {
		result.WinRate = (float64(winningTrades) / float64(result.TotalTrades)) * 100
	}
	if totalLoss != 0 {
		result.ProfitFactor = totalProfit / (-totalLoss)
	}
	result.ReturnPercent = ((balance - startBalance) / startBalance) * 100
	
	// Calculate score
	result.Score = calculateOptimizationScore(result)
	
	return result
}

// calculateOptimizationScore calculates optimization score
func calculateOptimizationScore(result OptimizationResult) float64 {
	if result.TotalTrades < 5 {
		return 0
	}
	
	// Weighted scoring
	winRateScore := result.WinRate * 0.35
	pfScore := math.Min(result.ProfitFactor*15, 100) * 0.30
	returnScore := math.Min(result.ReturnPercent, 200) * 0.25
	tradesScore := math.Min(float64(result.TotalTrades)*2, 50) * 0.10
	
	return winRateScore + pfScore + returnScore + tradesScore
}

// OptimizeAllStrategies optimizes all strategies
func OptimizeAllStrategies(symbol string, startBalance float64, days int) (map[string][]OptimizationResult, error) {
	log.Println("🔬 COMPREHENSIVE PARAMETER OPTIMIZATION")
	log.Println("=" + string(make([]byte, 70)))
	
	strategies := GetAdvancedStrategies()
	allResults := make(map[string][]OptimizationResult)
	
	for name := range strategies {
		results, err := OptimizeStrategyParameters(name, symbol, startBalance, days)
		if err != nil {
			log.Printf("  ❌ Failed to optimize %s: %v", name, err)
			continue
		}
		
		if len(results) > 0 {
			allResults[name] = results
		}
	}
	
	// Print summary
	printOptimizationSummary(allResults)
	
	return allResults, nil
}

// printOptimizationSummary prints optimization summary
func printOptimizationSummary(allResults map[string][]OptimizationResult) {
	log.Println("\n" + string(make([]byte, 90)))
	log.Println("🏆 OPTIMIZATION RESULTS - BEST PARAMETERS FOR EACH STRATEGY")
	log.Println(string(make([]byte, 90)))
	
	// Collect best result for each strategy
	bestResults := []OptimizationResult{}
	for _, results := range allResults {
		if len(results) > 0 {
			bestResults = append(bestResults, results[0])
		}
	}
	
	// Sort by score
	sort.Slice(bestResults, func(i, j int) bool {
		return bestResults[i].Score > bestResults[j].Score
	})
	
	log.Println("\n┌────┬─────────────────────────┬──────┬──────┬──────┬──────┬──────────┬──────────┬──────────┐")
	log.Println("│Rank│ Strategy                │ Conf │ Stop │ TP1  │ Risk │ Win Rate │ Return % │ Score    │")
	log.Println("├────┼─────────────────────────┼──────┼──────┼──────┼──────┼──────────┼──────────┼──────────┤")
	
	for i, r := range bestResults {
		if i >= 10 {
			break
		}
		
		log.Printf("│ %-2d │ %-23s │  %2d  │ %.1f  │ %.1f  │ %.1f%% │  %6.1f%% │  %6.1f%% │  %6.1f  │",
			i+1, truncate(r.StrategyName, 23),
			r.Parameters.MinConfluence,
			r.Parameters.StopATR,
			r.Parameters.TP1ATR,
			r.Parameters.RiskPercent,
			r.WinRate,
			r.ReturnPercent,
			r.Score)
	}
	
	log.Println("└────┴─────────────────────────┴──────┴──────┴──────┴──────┴──────────┴──────────┴──────────┘")
	
	if len(bestResults) > 0 {
		best := bestResults[0]
		log.Println("\n🥇 OVERALL BEST OPTIMIZED STRATEGY:")
		log.Printf("   Strategy: %s", best.StrategyName)
		log.Printf("   Min Confluence: %d", best.Parameters.MinConfluence)
		log.Printf("   Stop Loss: %.1f ATR", best.Parameters.StopATR)
		log.Printf("   Take Profit 1: %.1f ATR", best.Parameters.TP1ATR)
		log.Printf("   Take Profit 2: %.1f ATR", best.Parameters.TP2ATR)
		log.Printf("   Take Profit 3: %.1f ATR", best.Parameters.TP3ATR)
		log.Printf("   Risk Per Trade: %.1f%%", best.Parameters.RiskPercent)
		log.Printf("   Win Rate: %.1f%%", best.WinRate)
		log.Printf("   Return: %.1f%%", best.ReturnPercent)
		log.Printf("   Profit Factor: %.2f", best.ProfitFactor)
		log.Printf("   Total Trades: %d", best.TotalTrades)
		log.Printf("   Score: %.1f", best.Score)
	}
	
	log.Println("\n" + string(make([]byte, 90)))
}
