package main

import (
	"log"
	"math"
	"sort"
)

// StrategyTestResult holds test results for a strategy
type StrategyTestResult struct {
	StrategyName   string  `json:"strategyName"`
	Description    string  `json:"description"`
	Timeframe      string  `json:"timeframe"`
	TotalTrades    int     `json:"totalTrades"`
	WinningTrades  int     `json:"winningTrades"`
	LosingTrades   int     `json:"losingTrades"`
	WinRate        float64 `json:"winRate"`
	ProfitFactor   float64 `json:"profitFactor"`
	ReturnPercent  float64 `json:"returnPercent"`
	FinalBalance   float64 `json:"finalBalance"`
	MaxDrawdown    float64 `json:"maxDrawdown"`
	AverageRR      float64 `json:"averageRR"`
	Score          float64 `json:"score"`
	TargetWinRate  float64 `json:"targetWinRate"`
	TargetPF       float64 `json:"targetProfitFactor"`
}

// TestAllStrategies tests all advanced strategies
func TestAllStrategies(symbol string, startBalance float64) ([]StrategyTestResult, error) {
	strategies := GetAdvancedStrategies()
	results := []StrategyTestResult{}
	
	log.Println("🚀 Testing All Advanced Strategies")
	log.Println("=" + string(make([]byte, 70)))
	
	for name, strategy := range strategies {
		log.Printf("\n📊 Testing: %s (%s)", strategy.Name, strategy.Timeframe)
		
		// Determine days based on timeframe
		days := getOptimalDays(strategy.Timeframe)
		
		// Fetch data
		candles, err := fetchBinanceData(symbol, strategy.Timeframe, days)
		if err != nil {
			log.Printf("  ❌ Failed to fetch data: %v", err)
			continue
		}
		
		if len(candles) < 100 {
			log.Printf("  ⚠️  Insufficient data")
			continue
		}
		
		// Generate signals
		signals := []AdvancedSignal{}
		for i := 100; i < len(candles)-1; i++ {
			signal := GenerateSignalWithStrategy(candles[:i+1], name)
			if signal != nil {
				signals = append(signals, *signal)
			}
		}
		
		if len(signals) == 0 {
			log.Printf("  ⚠️  No signals generated")
			continue
		}
		
		// Simulate trades
		result := simulateAdvancedTrades(signals, candles, startBalance, strategy)
		result.StrategyName = name
		result.Description = strategy.Description
		result.Timeframe = strategy.Timeframe
		result.TargetWinRate = strategy.TargetWinRate
		result.TargetPF = strategy.TargetProfitFactor
		
		// Calculate score
		result.Score = calculateStrategyScore(result)
		
		results = append(results, result)
		
		log.Printf("  ✅ Trades: %d | WR: %.1f%% | Return: %.1f%% | PF: %.2f | Score: %.1f",
			result.TotalTrades, result.WinRate, result.ReturnPercent, result.ProfitFactor, result.Score)
	}
	
	// Sort by score
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	
	// Print summary
	printStrategySummary(results)
	
	return results, nil
}

// simulateAdvancedTrades simulates trades from advanced signals
func simulateAdvancedTrades(signals []AdvancedSignal, candles []Candle, startBalance float64, strategy AdvancedStrategy) StrategyTestResult {
	result := StrategyTestResult{
		FinalBalance: startBalance,
	}
	
	balance := startBalance
	maxBalance := balance
	totalProfit := 0.0
	totalLoss := 0.0
	
	for _, signal := range signals {
		// Calculate position size (2% risk)
		riskPercent := 2.0
		riskAmount := balance * (riskPercent / 100.0)
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
			result.WinningTrades++
			totalProfit += profit
		} else {
			result.LosingTrades++
			totalLoss += profit
		}
	}
	
	// Calculate metrics
	result.FinalBalance = balance
	if result.TotalTrades > 0 {
		result.WinRate = (float64(result.WinningTrades) / float64(result.TotalTrades)) * 100
	}
	if totalLoss != 0 {
		result.ProfitFactor = totalProfit / (-totalLoss)
	}
	result.ReturnPercent = ((balance - startBalance) / startBalance) * 100
	if maxBalance > 0 {
		result.MaxDrawdown = ((maxBalance - balance) / maxBalance) * 100
	}
	
	return result
}

// calculateStrategyScore calculates overall performance score
func calculateStrategyScore(result StrategyTestResult) float64 {
	if result.TotalTrades == 0 {
		return 0
	}
	
	// Weighted scoring
	winRateScore := result.WinRate * 0.3
	pfScore := result.ProfitFactor * 20 * 0.3
	returnScore := result.ReturnPercent * 0.2
	tradesScore := float64(result.TotalTrades) * 0.1
	drawdownScore := (100 - result.MaxDrawdown) * 0.1
	
	return winRateScore + pfScore + returnScore + tradesScore + drawdownScore
}

// printStrategySummary prints formatted summary
func printStrategySummary(results []StrategyTestResult) {
	log.Println("\n" + string(make([]byte, 90)))
	log.Println("🏆 STRATEGY PERFORMANCE RANKING")
	log.Println(string(make([]byte, 90)))
	
	log.Println("\n┌────┬─────────────────────────┬──────────┬─────────┬──────────┬──────────┬──────────┬────────┐")
	log.Println("│Rank│ Strategy                │Timeframe │ Trades  │ Win Rate │ Return % │ Profit F │ Score  │")
	log.Println("├────┼─────────────────────────┼──────────┼─────────┼──────────┼──────────┼──────────┼────────┤")
	
	for i, r := range results {
		if i >= 10 {
			break
		}
		
		status := "❌"
		if r.WinRate >= r.TargetWinRate*0.9 && r.ProfitFactor >= r.TargetPF*0.9 {
			status = "✅"
		} else if r.WinRate >= r.TargetWinRate*0.8 {
			status = "⚠️"
		}
		
		log.Printf("│ %-2d │ %-23s │ %-8s │ %7d │ %7.1f%% │ %7.1f%% │ %8.2f │ %6.1f │ %s",
			i+1, truncate(r.StrategyName, 23), r.Timeframe, r.TotalTrades,
			r.WinRate, r.ReturnPercent, r.ProfitFactor, r.Score, status)
	}
	
	log.Println("└────┴─────────────────────────┴──────────┴─────────┴──────────┴──────────┴──────────┴────────┘")
	
	if len(results) > 0 {
		best := results[0]
		log.Println("\n🥇 BEST STRATEGY:")
		log.Printf("   Name: %s", best.StrategyName)
		log.Printf("   Description: %s", best.Description)
		log.Printf("   Timeframe: %s", best.Timeframe)
		log.Printf("   Win Rate: %.1f%% (Target: %.1f%%)", best.WinRate, best.TargetWinRate)
		log.Printf("   Return: %.1f%%", best.ReturnPercent)
		log.Printf("   Profit Factor: %.2f (Target: %.2f)", best.ProfitFactor, best.TargetPF)
		log.Printf("   Total Trades: %d", best.TotalTrades)
		log.Printf("   Score: %.1f", best.Score)
	}
	
	log.Println("\n" + string(make([]byte, 90)))
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
