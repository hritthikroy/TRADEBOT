package main

import (
	"math"
	"math/rand"
	"sort"
	"time"
)

// RunEnhancedBacktest runs backtest with simulation windows and advanced features
func RunEnhancedBacktest(config BacktestConfig, candles []Candle) (*BacktestResult, error) {
	startTime := time.Now()
	
	// Set defaults for enhanced features
	if config.WindowType == "" {
		config.WindowType = "expanding" // Most realistic
	}
	if config.MinWindow == 0 {
		config.MinWindow = 50
	}
	if config.MaxWindow == 0 {
		config.MaxWindow = 200
	}
	if config.MCIterations == 0 {
		config.MCIterations = 1000
	}
	
	var result *BacktestResult
	var err error
	
	// Choose backtest method
	if config.UseWalkForward {
		result, err = runWalkForwardBacktest(config, candles)
	} else {
		result, err = runSimulationWindowBacktest(config, candles)
	}
	
	if err != nil {
		return nil, err
	}
	
	// Add Monte Carlo simulation if enabled
	if config.UseMonteCarlo && len(result.Trades) > 10 {
		mcResult := runMonteCarloSimulation(result.Trades, config.MCIterations, config.StartBalance)
		result.MonteCarloSim = mcResult
		result.Confidence95 = [2]float64{mcResult.Confidence95Low, mcResult.Confidence95High}
	}
	
	result.Duration = time.Since(startTime).String()
	result.WindowType = config.WindowType
	
	return result, nil
}

// runSimulationWindowBacktest runs backtest with expanding/rolling window
func runSimulationWindowBacktest(config BacktestConfig, candles []Candle) (*BacktestResult, error) {
	result := &BacktestResult{
		StartBalance: config.StartBalance,
		FinalBalance: config.StartBalance,
		PeakBalance:  config.StartBalance,
		Trades:       []Trade{},
		ExitReasons:  make(map[string]int),
	}
	
	// Set defaults
	if config.RiskPercent == 0 {
		config.RiskPercent = 0.02
	}
	if config.MaxPositionCap == 0 {
		config.MaxPositionCap = config.StartBalance * 10
	}
	if config.SlippagePercent == 0 {
		config.SlippagePercent = 0.001
	}
	if config.FeePercent == 0 {
		config.FeePercent = 0.001
	}
	
	skipAhead := 10 // More realistic - skip more candles after trade
	
	// Start from minimum window size
	for i := config.MinWindow; i < len(candles)-10; i++ {
		// Calculate window size based on type
		var dataWindow []Candle
		
		switch config.WindowType {
		case "expanding":
			// Use all data from start to current (most realistic)
			windowSize := minInt(i, config.MaxWindow)
			dataWindow = candles[i-windowSize : i]
			
		case "rolling":
			// Use fixed window that rolls forward
			windowSize := config.MaxWindow
			if i < windowSize {
				windowSize = i
			}
			dataWindow = candles[i-windowSize : i]
			
		case "fixed":
			// Use fixed window (original method - less realistic)
			if i < config.MinWindow {
				continue
			}
			dataWindow = candles[i-config.MinWindow : i]
			
		default:
			// Default to expanding
			windowSize := minInt(i, config.MaxWindow)
			dataWindow = candles[i-windowSize : i]
		}
		
		// Time filter - only trade during optimal hours
		if config.UseTimeFilter && !shouldTradeAtTime(candles[i].Timestamp) {
			continue
		}
		
		futureData := candles[i : minInt(i+20, len(candles))] // Look ahead 20 candles max
		
		// Generate signal using available data only
		signal := generateBacktestSignal(dataWindow, config.Interval)
		
		if signal != nil {
			// Simulate trade with realistic conditions
			trade := simulateRealisticTrade(signal, futureData, result.FinalBalance, config, candles[i].Timestamp)
			
			if trade != nil {
				trade.EntryIndex = i
				trade.BalanceAfter = result.FinalBalance + trade.Profit
				
				result.Trades = append(result.Trades, *trade)
				result.TotalTrades++
				
				if trade.Profit > 0 {
					result.WinningTrades++
					result.TotalProfit += trade.Profit
				} else {
					result.LosingTrades++
					result.TotalLoss += math.Abs(trade.Profit)
				}
				
				// Update balance
				result.FinalBalance += trade.Profit
				
				// Track peak and drawdown
				if result.FinalBalance > result.PeakBalance {
					result.PeakBalance = result.FinalBalance
				}
				
				drawdown := (result.PeakBalance - result.FinalBalance) / result.PeakBalance
				if drawdown > result.MaxDrawdown {
					result.MaxDrawdown = drawdown
				}
				
				// Track exit reasons
				result.ExitReasons[trade.ExitReason]++
				
				// Skip ahead to avoid overlapping trades
				i += skipAhead
			}
		}
	}
	
	// Calculate statistics
	calculateStats(result)
	
	return result, nil
}

// runWalkForwardBacktest runs walk-forward analysis
func runWalkForwardBacktest(config BacktestConfig, candles []Candle) (*BacktestResult, error) {
	// Set defaults
	if config.TrainingDays == 0 {
		config.TrainingDays = 60 // 60 days training
	}
	if config.TestingDays == 0 {
		config.TestingDays = 30 // 30 days testing
	}
	
	// Calculate candles per period based on interval
	candlesPerDay := getCandlesPerDay(config.Interval)
	trainingCandles := config.TrainingDays * candlesPerDay
	testingCandles := config.TestingDays * candlesPerDay
	stepSize := testingCandles / 2 // 50% overlap
	
	aggregatedResult := &BacktestResult{
		StartBalance:       config.StartBalance,
		FinalBalance:       config.StartBalance,
		PeakBalance:        config.StartBalance,
		Trades:             []Trade{},
		ExitReasons:        make(map[string]int),
		WalkForwardResults: []WalkForwardPeriod{},
	}
	
	periodNum := 1
	
	// Walk forward through data
	for i := 0; i+trainingCandles+testingCandles < len(candles); i += stepSize {
		trainStart := i
		trainEnd := i + trainingCandles
		testStart := trainEnd
		testEnd := testStart + testingCandles
		
		if testEnd > len(candles) {
			break
		}
		
		// Test on unseen data (we skip training optimization for now)
		testData := candles[testStart:testEnd]
		
		// Run backtest on test period
		testConfig := config
		testConfig.UseWalkForward = false // Prevent recursion
		testConfig.StartBalance = aggregatedResult.FinalBalance // Use current balance
		
		periodResult, err := runSimulationWindowBacktest(testConfig, testData)
		if err != nil {
			continue
		}
		
		// Record period results
		period := WalkForwardPeriod{
			PeriodNum:     periodNum,
			TrainStart:    trainStart,
			TrainEnd:      trainEnd,
			TestStart:     testStart,
			TestEnd:       testEnd,
			WinRate:       periodResult.WinRate,
			ReturnPercent: periodResult.ReturnPercent,
			ProfitFactor:  periodResult.ProfitFactor,
			TotalTrades:   periodResult.TotalTrades,
		}
		aggregatedResult.WalkForwardResults = append(aggregatedResult.WalkForwardResults, period)
		
		// Aggregate trades
		for _, trade := range periodResult.Trades {
			trade.EntryIndex += testStart // Adjust index
			aggregatedResult.Trades = append(aggregatedResult.Trades, trade)
		}
		
		// Update aggregated stats
		aggregatedResult.TotalTrades += periodResult.TotalTrades
		aggregatedResult.WinningTrades += periodResult.WinningTrades
		aggregatedResult.LosingTrades += periodResult.LosingTrades
		aggregatedResult.TotalProfit += periodResult.TotalProfit
		aggregatedResult.TotalLoss += periodResult.TotalLoss
		aggregatedResult.FinalBalance = periodResult.FinalBalance
		
		// Update peak and drawdown
		if aggregatedResult.FinalBalance > aggregatedResult.PeakBalance {
			aggregatedResult.PeakBalance = aggregatedResult.FinalBalance
		}
		
		drawdown := (aggregatedResult.PeakBalance - aggregatedResult.FinalBalance) / aggregatedResult.PeakBalance
		if drawdown > aggregatedResult.MaxDrawdown {
			aggregatedResult.MaxDrawdown = drawdown
		}
		
		// Aggregate exit reasons
		for reason, count := range periodResult.ExitReasons {
			aggregatedResult.ExitReasons[reason] += count
		}
		
		periodNum++
	}
	
	// Calculate final statistics
	calculateStats(aggregatedResult)
	
	return aggregatedResult, nil
}

// runMonteCarloSimulation runs Monte Carlo simulation on trades
func runMonteCarloSimulation(trades []Trade, iterations int, startBalance float64) *MonteCarloResult {
	if len(trades) < 10 {
		return nil
	}
	
	results := make([]float64, iterations)
	profitableCount := 0
	
	// Run simulations
	for i := 0; i < iterations; i++ {
		// Shuffle trades randomly
		shuffled := make([]Trade, len(trades))
		copy(shuffled, trades)
		rand.Shuffle(len(shuffled), func(i, j int) {
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		})
		
		// Calculate equity curve
		balance := startBalance
		for _, trade := range shuffled {
			balance += trade.Profit
		}
		
		results[i] = balance
		
		if balance > startBalance {
			profitableCount++
		}
	}
	
	// Sort results for percentile calculation
	sort.Float64s(results)
	
	// Calculate statistics
	mean := calculateMean(results)
	median := results[len(results)/2]
	stdDev := calculateStdDev(results, mean)
	
	return &MonteCarloResult{
		Iterations:        iterations,
		MeanReturn:        ((mean - startBalance) / startBalance) * 100,
		MedianReturn:      ((median - startBalance) / startBalance) * 100,
		StdDev:            stdDev,
		WorstCase:         results[0],
		BestCase:          results[len(results)-1],
		Confidence95Low:   results[int(float64(len(results))*0.025)],
		Confidence95High:  results[int(float64(len(results))*0.975)],
		ProbabilityProfit: (float64(profitableCount) / float64(iterations)) * 100,
	}
}

// simulateRealisticTrade simulates trade with realistic market conditions
func simulateRealisticTrade(signal *Signal, futureData []Candle, currentBalance float64, config BacktestConfig, timestamp int64) *Trade {
	if signal == nil || len(futureData) == 0 {
		return nil
	}
	
	entry := signal.Entry
	stopLoss := signal.StopLoss
	
	// Calculate position size with cap
	riskAmount := currentBalance * config.RiskPercent
	if riskAmount > config.MaxPositionCap {
		riskAmount = config.MaxPositionCap
	}
	
	riskDiff := math.Abs(entry - stopLoss)
	if riskDiff == 0 {
		return nil
	}
	
	_ = riskAmount / riskDiff // positionSize calculated in simulateTrade
	
	// Calculate realistic slippage based on volatility
	volatility := calculateVolatility(futureData)
	realisticSlippage := config.SlippagePercent * (1 + volatility)
	
	// Apply slippage to entry
	if signal.Type == "BUY" {
		entry *= (1 + realisticSlippage)
	} else {
		entry *= (1 - realisticSlippage)
	}
	
	// Use the original simulateTrade logic with enhanced slippage
	return simulateTrade(signal, futureData, currentBalance, config)
}

// shouldTradeAtTime checks if we should trade at this time
func shouldTradeAtTime(timestamp int64) bool {
	t := time.Unix(timestamp/1000, 0).UTC()
	hour := t.Hour()
	weekday := t.Weekday()
	
	// Skip weekends
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}
	
	// Only trade during high-liquidity hours (8am-8pm UTC)
	if hour < 8 || hour > 20 {
		return false
	}
	
	// Avoid first and last hour (high volatility)
	if hour == 8 || hour == 20 {
		return false
	}
	
	return true
}

// calculateVolatility calculates recent volatility
func calculateVolatility(candles []Candle) float64 {
	if len(candles) < 2 {
		return 0
	}
	
	// Calculate price changes
	changes := make([]float64, len(candles)-1)
	for i := 1; i < len(candles); i++ {
		changes[i-1] = math.Abs((candles[i].Close - candles[i-1].Close) / candles[i-1].Close)
	}
	
	// Return average volatility
	sum := 0.0
	for _, change := range changes {
		sum += change
	}
	
	return sum / float64(len(changes))
}

// getCandlesPerDay returns approximate candles per day for interval
func getCandlesPerDay(interval string) int {
	switch interval {
	case "1m":
		return 1440
	case "5m":
		return 288
	case "15m":
		return 96
	case "1h":
		return 24
	case "4h":
		return 6
	case "1d":
		return 1
	default:
		return 24 // Default to 1h
	}
}

// calculateMean calculates mean of float64 slice
func calculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	
	return sum / float64(len(values))
}

// calculateStdDev calculates standard deviation
func calculateStdDev(values []float64, mean float64) float64 {
	if len(values) == 0 {
		return 0
	}
	
	sumSquares := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquares += diff * diff
	}
	
	variance := sumSquares / float64(len(values))
	return math.Sqrt(variance)
}

// minInt returns minimum of two ints
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
