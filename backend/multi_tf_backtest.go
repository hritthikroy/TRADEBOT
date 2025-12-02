package main

import (
	"fmt"
	"log"
)

// RunMultiTimeframeBacktest runs backtest with multi-TF analysis
func RunMultiTimeframeBacktest(symbol string, days int, startBalance float64) (*BacktestResult, error) {
	log.Println("🚀 Starting Multi-Timeframe Backtest")
	log.Printf("Symbol: %s, Days: %d, Balance: $%.2f", symbol, days, startBalance)
	log.Println("Strategy: 4h→1h→15m→3m→1m Top-Down Analysis")
	log.Println("=" + string(make([]byte, 70)))
	
	// Fetch all timeframes
	log.Println("\n📊 Fetching multi-timeframe data...")
	
	candles4h, err := fetchBinanceData(symbol, "4h", days)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch 4h data: %w", err)
	}
	log.Printf("✅ 4h: %d candles", len(candles4h))
	
	candles1h, err := fetchBinanceData(symbol, "1h", days)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch 1h data: %w", err)
	}
	log.Printf("✅ 1h: %d candles", len(candles1h))
	
	candles15m, err := fetchBinanceData(symbol, "15m", days)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch 15m data: %w", err)
	}
	log.Printf("✅ 15m: %d candles", len(candles15m))
	
	candles3m, err := fetchBinanceData(symbol, "3m", days)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch 3m data: %w", err)
	}
	log.Printf("✅ 3m: %d candles", len(candles3m))
	
	candles1m, err := fetchBinanceData(symbol, "1m", days)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch 1m data: %w", err)
	}
	log.Printf("✅ 1m: %d candles", len(candles1m))
	
	// Generate signals using multi-TF analysis
	log.Println("\n🔍 Analyzing multi-timeframe confluence...")
	signals := []MultiTimeframeSignal{}
	
	// Scan through 1m candles and generate signals
	for i := 200; i < len(candles1m)-1; i++ {
		// Get corresponding indices for other timeframes
		idx4h := minIntMTF(i/240, len(candles4h)-1)
		idx1h := minIntMTF(i/60, len(candles1h)-1)
		idx15m := minIntMTF(i/15, len(candles15m)-1)
		idx3m := minIntMTF(i/3, len(candles3m)-1)
		
		// Generate signal
		signal := GenerateMultiTimeframeSignal(
			candles4h[:idx4h+1],
			candles1h[:idx1h+1],
			candles15m[:idx15m+1],
			candles3m[:idx3m+1],
			candles1m[:i+1],
		)
		
		if signal != nil {
			signals = append(signals, *signal)
		}
	}
	
	log.Printf("\n✅ Generated %d multi-timeframe signals", len(signals))
	
	if len(signals) == 0 {
		log.Println("⚠️  No signals met the strict multi-TF criteria")
		return &BacktestResult{
			StartBalance: startBalance,
			FinalBalance: startBalance,
			TotalTrades:  0,
		}, nil
	}
	
	// Simulate trades
	log.Println("\n💼 Simulating trades...")
	result := simulateMultiTFTrades(signals, candles1m, startBalance)
	
	// Print results
	printMultiTFResults(result)
	
	return result, nil
}

// simulateMultiTFTrades simulates trade execution
func simulateMultiTFTrades(signals []MultiTimeframeSignal, candles []Candle, startBalance float64) *BacktestResult {
	result := &BacktestResult{
		StartBalance: startBalance,
		FinalBalance: startBalance,
		Trades:       []Trade{},
		ExitReasons:  make(map[string]int),
	}
	
	balance := startBalance
	maxBalance := balance
	
	for _, signal := range signals {
		// Calculate position size (1.5% risk for multi-TF setups)
		riskPercent := 1.5
		riskAmount := balance * (riskPercent / 100.0)
		riskPerUnit := signal.StopLoss - signal.EntryPrice
		if signal.Type == "SELL" {
			riskPerUnit = signal.EntryPrice - signal.StopLoss
		}
		
		if riskPerUnit == 0 {
			continue
		}
		
		positionSize := riskAmount / riskPerUnit
		if positionSize <= 0 {
			continue
		}
		
		// Find exit point
		trade := Trade{
			Type:     signal.Type,
			Entry:    signal.EntryPrice,
			StopLoss: signal.StopLoss,
		}
		
		exitFound := false
		var profit float64
		
		// Look for exit in next candles
		startIdx := findCandleIndex(candles, signal.EntryTime)
		if startIdx == -1 {
			continue
		}
		
		for i := startIdx + 1; i < len(candles); i++ {
			if signal.Type == "BUY" {
				// Check stop loss
				if candles[i].Low <= signal.StopLoss {
					trade.Exit = signal.StopLoss
					trade.ExitReason = "Stop Loss"
					profit = (trade.Exit - trade.Entry) * positionSize
					exitFound = true
					break
				}
				// Check TPs
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
				if candles[i].High >= signal.StopLoss {
					trade.Exit = signal.StopLoss
					trade.ExitReason = "Stop Loss"
					profit = (trade.Entry - trade.Exit) * positionSize
					exitFound = true
					break
				}
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
	
	result.ReturnPercent = ((balance - startBalance) / startBalance) * 100
	
	if maxBalance > 0 {
		result.MaxDrawdown = ((maxBalance - balance) / maxBalance) * 100
	}
	
	return result
}

// findCandleIndex finds the index of a candle by timestamp
func findCandleIndex(candles []Candle, timestamp int64) int {
	for i, c := range candles {
		if c.Timestamp >= timestamp {
			return i
		}
	}
	return -1
}

// printMultiTFResults prints formatted results
func printMultiTFResults(result *BacktestResult) {
	log.Println("\n" + string(make([]byte, 70)))
	log.Println("📊 MULTI-TIMEFRAME BACKTEST RESULTS")
	log.Println(string(make([]byte, 70)))
	
	log.Printf("\n💰 PERFORMANCE:")
	log.Printf("  Start Balance:    $%.2f", result.StartBalance)
	log.Printf("  Final Balance:    $%.2f", result.FinalBalance)
	log.Printf("  Return:           %.1f%%", result.ReturnPercent)
	log.Printf("  Profit Factor:    %.2f", result.ProfitFactor)
	
	log.Printf("\n📈 TRADE STATISTICS:")
	log.Printf("  Total Trades:     %d", result.TotalTrades)
	log.Printf("  Winning Trades:   %d", result.WinningTrades)
	log.Printf("  Losing Trades:    %d", result.LosingTrades)
	log.Printf("  Win Rate:         %.1f%%", result.WinRate)
	
	log.Printf("\n💵 PROFIT/LOSS:")
	log.Printf("  Total Profit:     $%.2f", result.TotalProfit)
	log.Printf("  Total Loss:       $%.2f", result.TotalLoss)
	log.Printf("  Net Profit:       $%.2f", result.TotalProfit+result.TotalLoss)
	log.Printf("  Max Drawdown:     %.1f%%", result.MaxDrawdown)
	
	if len(result.ExitReasons) > 0 {
		log.Printf("\n🎯 EXIT REASONS:")
		for reason, count := range result.ExitReasons {
			percentage := (float64(count) / float64(result.TotalTrades)) * 100
			log.Printf("  %-15s: %d (%.1f%%)", reason, count, percentage)
		}
	}
	
	// Performance rating
	log.Println("\n" + string(make([]byte, 70)))
	if result.WinRate >= 75 && result.ProfitFactor >= 2.0 {
		log.Println("🏆 RATING: EXCELLENT - Multi-TF strategy working perfectly!")
	} else if result.WinRate >= 65 && result.ProfitFactor >= 1.5 {
		log.Println("✅ RATING: GOOD - Multi-TF strategy is profitable")
	} else if result.WinRate >= 55 && result.ProfitFactor >= 1.2 {
		log.Println("⚠️  RATING: MODERATE - Needs optimization")
	} else {
		log.Println("❌ RATING: POOR - Strategy needs significant improvement")
	}
	log.Println(string(make([]byte, 70)))
}

func minIntMTF(a, b int) int {
	if a < b {
		return a
	}
	return b
}
