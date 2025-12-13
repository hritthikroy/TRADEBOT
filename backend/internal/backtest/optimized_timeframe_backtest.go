package backtest

import (
	"log"
)

// RunAllTimeframesBacktest tests all timeframes with optimized strategies
func RunAllTimeframesBacktest(symbol string, startBalance float64) (map[string]*BacktestResult, error) {
	timeframes := []string{"1m", "3m", "5m", "15m", "30m", "1h", "2h", "4h", "1d"}
	results := make(map[string]*BacktestResult)
	
	log.Println("🚀 Testing All Timeframes with Optimized Strategies")
	log.Println("=" + string(make([]byte, 70)))
	
	for _, tf := range timeframes {
		// Determine optimal days for each timeframe
		days := getOptimalDays(tf)
		
		log.Printf("\n📊 Testing %s (last %d days)...", tf, days)
		
		// Fetch data
		candles, err := fetchBinanceData(symbol, tf, days)
		if err != nil {
			log.Printf("  ❌ Failed: %v", err)
			continue
		}
		
		if len(candles) < 100 {
			log.Printf("  ⚠️  Insufficient data (%d candles)", len(candles))
			continue
		}
		
		// Get optimized strategy for this timeframe
		strategy := GetOptimizedStrategy(tf)
		
		// Generate signals using optimized strategy
		signals := ApplyTimeframeStrategy(candles, strategy)
		
		if len(signals) == 0 {
			log.Printf("  ⚠️  No signals generated")
			continue
		}
		
		// Run backtest
		config := BacktestConfig{
			Interval:     tf,
			Days:         days,
			StartBalance: startBalance,
		}
		
		result := simulateTradesFromSignals(signals, candles, config, strategy)
		results[tf] = result
		
		// Log results
		if result.TotalTrades > 0 {
			log.Printf("  ✅ Trades: %d | Win Rate: %.1f%% | Return: %.1f%% | PF: %.2f",
				result.TotalTrades, result.WinRate, result.ReturnPercent, result.ProfitFactor)
		}
	}
	
	// Print summary
	printTimeframeSummary(results, startBalance)
	
	return results, nil
}

// simulateTradesFromSignals simulates trade execution from signals
func simulateTradesFromSignals(signals []OptimizedSignal, candles []Candle, config BacktestConfig, strategy TimeframeStrategy) *BacktestResult {
	result := &BacktestResult{
		StartBalance: config.StartBalance,
		FinalBalance: config.StartBalance,
		Trades:       []Trade{},
		ExitReasons:  make(map[string]int),
	}
	
	balance := config.StartBalance
	maxBalance := balance
	
	for _, signal := range signals {
		// Calculate position size based on risk
		riskAmount := balance * (strategy.MaxRiskPercent / 100.0)
		riskPerUnit := signal.StopLoss - signal.Entry
		if signal.Type == "SELL" {
			riskPerUnit = signal.Entry - signal.StopLoss
		}
		
		if riskPerUnit == 0 {
			continue
		}
		
		positionSize := riskAmount / riskPerUnit
		if positionSize <= 0 {
			continue
		}
		
		// Create trade
		trade := Trade{
			Type:     signal.Type,
			Entry:    signal.Entry,
			StopLoss: signal.StopLoss,
		}
		
		// Find exit
		exitFound := false
		var profit float64
		
		for i := 0; i < len(candles); i++ {
			if candles[i].Timestamp <= signal.Timestamp {
				continue
			}
			
			if signal.Type == "BUY" {
				// Check stop loss first
				if candles[i].Low <= signal.StopLoss {
					trade.Exit = signal.StopLoss
					trade.ExitReason = "Stop Loss"
					profit = (trade.Exit - trade.Entry) * positionSize
					exitFound = true
					break
				}
				// Check take profits (TP3 > TP2 > TP1)
				if candles[i].High >= signal.TP3 {
					trade.Exit = signal.TP3
					trade.ExitReason = "TP3"
					profit = (trade.Exit - trade.Entry) * positionSize
					exitFound = true
					break
				} else if candles[i].High >= signal.TP2 {
					trade.Exit = signal.TP2
					trade.ExitReason = "TP2"
					profit = (trade.Exit - trade.Entry) * positionSize
					exitFound = true
					break
				} else if candles[i].High >= signal.TP1 {
					trade.Exit = signal.TP1
					trade.ExitReason = "TP1"
					profit = (trade.Exit - trade.Entry) * positionSize
					exitFound = true
					break
				}
			} else { // SELL
				// Check stop loss first
				if candles[i].High >= signal.StopLoss {
					trade.Exit = signal.StopLoss
					trade.ExitReason = "Stop Loss"
					profit = (trade.Entry - trade.Exit) * positionSize
					exitFound = true
					break
				}
				// Check take profits
				if candles[i].Low <= signal.TP3 {
					trade.Exit = signal.TP3
					trade.ExitReason = "TP3"
					profit = (trade.Entry - trade.Exit) * positionSize
					exitFound = true
					break
				} else if candles[i].Low <= signal.TP2 {
					trade.Exit = signal.TP2
					trade.ExitReason = "TP2"
					profit = (trade.Entry - trade.Exit) * positionSize
					exitFound = true
					break
				} else if candles[i].Low <= signal.TP1 {
					trade.Exit = signal.TP1
					trade.ExitReason = "TP1"
					profit = (trade.Entry - trade.Exit) * positionSize
					exitFound = true
					break
				}
			}
		}
		
		if !exitFound {
			continue
		}
		
		// Update balance
		balance += profit
		trade.Profit = profit
		trade.ProfitPercent = (profit / balance) * 100
		trade.BalanceAfter = balance
		
		if balance > maxBalance {
			maxBalance = balance
		}
		
		// Record trade
		result.Trades = append(result.Trades, trade)
		result.ExitReasons[trade.ExitReason]++
		
		if profit > 0 {
			result.WinningTrades++
			result.TotalProfit += profit
		} else {
			result.LosingTrades++
			result.TotalLoss += profit
		}
	}
	
	// Calculate metrics
	result.TotalTrades = len(result.Trades)
	result.FinalBalance = balance
	
	if result.TotalTrades > 0 {
		result.WinRate = (float64(result.WinningTrades) / float64(result.TotalTrades)) * 100
	}
	
	if result.TotalLoss != 0 {
		result.ProfitFactor = result.TotalProfit / (-result.TotalLoss)
	}
	
	result.ReturnPercent = ((balance - config.StartBalance) / config.StartBalance) * 100
	
	if maxBalance > 0 {
		result.MaxDrawdown = ((maxBalance - balance) / maxBalance) * 100
	}
	
	return result
}

// getOptimalDays returns optimal backtest period for each timeframe
func getOptimalDays(timeframe string) int {
	days := map[string]int{
		"1m":  7,
		"3m":  7,
		"5m":  14,
		"15m": 30,
		"30m": 30,
		"1h":  60,
		"2h":  60,
		"4h":  90,
		"1d":  180,
	}
	
	if d, ok := days[timeframe]; ok {
		return d
	}
	return 30
}

// printTimeframeSummary prints comparison of all timeframes
func printTimeframeSummary(results map[string]*BacktestResult, startBalance float64) {
	log.Println("\n" + string(make([]byte, 90)))
	log.Println("📊 TIMEFRAME COMPARISON SUMMARY")
	log.Println(string(make([]byte, 90)))
	
	// Find best performers
	var bestWinRate, bestReturn, bestPF string
	var maxWinRate, maxReturn, maxPF float64
	
	log.Println("\n┌──────────┬─────────┬──────────┬──────────┬──────────┬──────────┬──────────┐")
	log.Println("│Timeframe │ Trades  │ Win Rate │ Return % │ Profit F │ End Bal  │ Status   │")
	log.Println("├──────────┼─────────┼──────────┼──────────┼──────────┼──────────┼──────────┤")
	
	timeframes := []string{"1m", "3m", "5m", "15m", "30m", "1h", "2h", "4h", "1d"}
	
	for _, tf := range timeframes {
		result, ok := results[tf]
		if !ok || result.TotalTrades == 0 {
			log.Printf("│ %-8s │ %-7s │ %-8s │ %-8s │ %-8s │ %-8s │ %-8s │\n",
				tf, "N/A", "N/A", "N/A", "N/A", "N/A", "NO DATA")
			continue
		}
		
		status := "❌ LOSS"
		if result.ReturnPercent > 20 {
			status = "🔥 GREAT"
		} else if result.ReturnPercent > 10 {
			status = "✅ GOOD"
		} else if result.ReturnPercent > 0 {
			status = "✅ PROFIT"
		}
		
		log.Printf("│ %-8s │ %7d │ %7.1f%% │ %7.1f%% │ %8.2f │ $%7.0f │ %-8s │\n",
			tf, result.TotalTrades, result.WinRate, result.ReturnPercent, 
			result.ProfitFactor, result.FinalBalance, status)
		
		// Track best performers
		if result.WinRate > maxWinRate {
			maxWinRate = result.WinRate
			bestWinRate = tf
		}
		if result.ReturnPercent > maxReturn {
			maxReturn = result.ReturnPercent
			bestReturn = tf
		}
		if result.ProfitFactor > maxPF {
			maxPF = result.ProfitFactor
			bestPF = tf
		}
	}
	
	log.Println("└──────────┴─────────┴──────────┴──────────┴──────────┴──────────┴──────────┘")
	
	log.Println("\n🏆 BEST PERFORMERS:")
	log.Printf("  🎯 Best Win Rate:      %s (%.1f%%)\n", bestWinRate, maxWinRate)
	log.Printf("  💰 Best Return:        %s (%.1f%%)\n", bestReturn, maxReturn)
	log.Printf("  📊 Best Profit Factor: %s (%.2f)\n", bestPF, maxPF)
	
	log.Println("\n💡 RECOMMENDATIONS:")
	log.Println("  • For SCALPING: Use 5m or 15m timeframes")
	log.Println("  • For DAY TRADING: Use 15m or 1h timeframes")
	log.Println("  • For SWING TRADING: Use 4h or 1d timeframes")
	log.Printf("  • BEST OVERALL: %s timeframe\n", bestReturn)
	
	log.Println("\n" + string(make([]byte, 90)))
}
