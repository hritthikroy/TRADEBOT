package backtest

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"
)

// UnifiedBacktestConfig - One config to rule them all
type UnifiedBacktestConfig struct {
	// Basic Configuration
	Symbol          string  `json:"symbol"`
	Interval        string  `json:"interval"`
	Days            int     `json:"days"`
	StartBalance    float64 `json:"startBalance"`
	Strategy        string  `json:"strategy"`
	
	// Risk Management
	RiskPercent         float64 `json:"riskPercent"`         // Risk per trade (default: 0.3%)
	MaxPositionCap      float64 `json:"maxPositionCap"`      // Max position size
	MaxDailyLoss        float64 `json:"maxDailyLoss"`        // Stop trading if daily loss exceeds %
	MaxWeeklyLoss       float64 `json:"maxWeeklyLoss"`       // Stop trading if weekly loss exceeds %
	MaxConsecutiveLoss  int     `json:"maxConsecutiveLoss"`  // Stop after N consecutive losses
	DynamicPositionSize bool    `json:"dynamicPositionSize"` // Adjust position based on equity
	
	// Trading Costs
	SlippagePercent     float64 `json:"slippagePercent"`     // Slippage (default: 0.15%)
	FeePercent          float64 `json:"feePercent"`          // Trading fees (default: 0.1%)
	RealisticSlippage   bool    `json:"realisticSlippage"`   // Variable slippage based on volatility
	IncludeSpread       bool    `json:"includeSpread"`       // Include bid-ask spread
	SpreadPercent       float64 `json:"spreadPercent"`       // Spread as % of price
	
	// Market Filters
	MinVolatility       float64 `json:"minVolatility"`       // Minimum ATR to trade
	MaxVolatility       float64 `json:"maxVolatility"`       // Maximum ATR to trade
	MinVolume           float64 `json:"minVolume"`           // Minimum volume multiplier
	TradingHoursOnly    bool    `json:"tradingHoursOnly"`    // Only trade during specific hours
	MaxTradesPerDay     int     `json:"maxTradesPerDay"`     // Prevent overtrading
	
	// Simulation Methods
	WindowType          string  `json:"windowType"`          // "expanding", "rolling", "fixed"
	MinWindow           int     `json:"minWindow"`           // Minimum candles needed
	MaxWindow           int     `json:"maxWindow"`           // Maximum window size
	UseWalkForward      bool    `json:"useWalkForward"`      // Walk-forward analysis
	TrainingDays        int     `json:"trainingDays"`        // Days for training
	TestingDays         int     `json:"testingDays"`         // Days for testing
	
	// Advanced Analysis
	EnableMonteCarlo    bool    `json:"enableMonteCarlo"`    // Monte Carlo simulation
	MonteCarloRuns      int     `json:"monteCarloRuns"`      // Number of MC simulations
	EnableStressTest    bool    `json:"enableStressTest"`    // Test under extreme conditions
	EnableMultiTF       bool    `json:"enableMultiTF"`       // Multi-timeframe analysis
	EnablePartialExits  bool    `json:"enablePartialExits"`  // Use partial exit logic
	
	// Parallel Processing
	EnableParallel      bool    `json:"enableParallel"`      // Run multiple strategies in parallel
	Strategies          []string `json:"strategies"`          // List of strategies to test
}

// UnifiedBacktestResult - Comprehensive results
type UnifiedBacktestResult struct {
	// Basic Metrics
	TotalTrades         int                 `json:"totalTrades"`
	WinningTrades       int                 `json:"winningTrades"`
	LosingTrades        int                 `json:"losingTrades"`
	WinRate             float64             `json:"winRate"`
	TotalProfit         float64             `json:"totalProfit"`
	TotalLoss           float64             `json:"totalLoss"`
	NetProfit           float64             `json:"netProfit"`
	ReturnPercent       float64             `json:"returnPercent"`
	ProfitFactor        float64             `json:"profitFactor"`
	AverageRR           float64             `json:"averageRR"`
	MaxDrawdown         float64             `json:"maxDrawdown"`
	StartBalance        float64             `json:"startBalance"`
	FinalBalance        float64             `json:"finalBalance"`
	PeakBalance         float64             `json:"peakBalance"`
	
	// Advanced Risk Metrics
	SharpeRatio         float64             `json:"sharpeRatio"`
	SortinoRatio        float64             `json:"sortinoRatio"`
	CalmarRatio         float64             `json:"calmarRatio"`
	MaxDrawdownDuration int                 `json:"maxDrawdownDuration"`
	RecoveryFactor      float64             `json:"recoveryFactor"`
	MaxConsecutiveLosses int                `json:"maxConsecutiveLosses"`
	
	// Performance Breakdown
	WinStreakMax        int                 `json:"winStreakMax"`
	LossStreakMax       int                 `json:"lossStreakMax"`
	AverageWin          float64             `json:"averageWin"`
	AverageLoss         float64             `json:"averageLoss"`
	LargestWin          float64             `json:"largestWin"`
	LargestLoss         float64             `json:"largestLoss"`
	ExpectancyPerTrade  float64             `json:"expectancyPerTrade"`
	
	// Time Analysis
	AvgTradeHours       float64             `json:"avgTradeHours"`
	BestTradingHour     int                 `json:"bestTradingHour"`
	WorstTradingHour    int                 `json:"worstTradingHour"`
	
	// Trade Details
	Trades              []Trade             `json:"trades"`
	ExitReasons         map[string]int      `json:"exitReasons"`
	
	// Advanced Analysis Results
	MonteCarloResults   *MonteCarloAnalysis  `json:"monteCarloResults,omitempty"`
	WalkForwardAnalysis *WalkForwardAnalysis `json:"walkForwardAnalysis,omitempty"`
	StressTestResults   *StressTestAnalysis  `json:"stressTestResults,omitempty"`
	
	// Market Condition Analysis
	PerformanceByVolatility map[string]float64 `json:"performanceByVolatility"`
	PerformanceByTrend      map[string]float64 `json:"performanceByTrend"`
	
	// Metadata
	StrategyName        string              `json:"strategyName"`
	Duration            string              `json:"duration"`
	WindowType          string              `json:"windowType,omitempty"`
}

// RunUnifiedBacktest - The ONE backtest engine to rule them all
func RunUnifiedBacktest(config UnifiedBacktestConfig, candles []Candle) (*UnifiedBacktestResult, error) {
	startTime := time.Now()
	
	log.Println("🚀 Starting Unified Backtest Engine")
	log.Printf("📊 Symbol: %s | Interval: %s | Days: %d | Strategy: %s", 
		config.Symbol, config.Interval, config.Days, config.Strategy)
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	// Set intelligent defaults
	applyDefaults(&config)
	
	// Choose execution path based on configuration
	var result *UnifiedBacktestResult
	var err error
	
	if config.EnableParallel && len(config.Strategies) > 1 {
		// Parallel multi-strategy testing
		result, err = runParallelStrategies(config, candles)
	} else if config.UseWalkForward {
		// Walk-forward analysis
		result, err = runWalkForwardUnified(config, candles)
	} else if config.EnablePartialExits {
		// Professional partial exits
		result, err = runWithPartialExits(config, candles)
	} else {
		// Standard backtest with all features
		result, err = runStandardUnified(config, candles)
	}
	
	if err != nil {
		return nil, err
	}
	
	// Calculate advanced metrics
	calculateAdvancedMetricsUnified(result, candles)
	
	// Run Monte Carlo if enabled
	if config.EnableMonteCarlo && len(result.Trades) > 10 {
		result.MonteCarloResults = runMonteCarloUnified(result.Trades, config.MonteCarloRuns, config.StartBalance)
	}
	
	// Run stress test if enabled
	if config.EnableStressTest {
		result.StressTestResults = runStressTestUnified(config, candles)
	}
	
	result.Duration = time.Since(startTime).String()
	
	// Print comprehensive summary
	printUnifiedSummary(result)
	
	return result, nil
}

// applyDefaults sets intelligent defaults
func applyDefaults(config *UnifiedBacktestConfig) {
	if config.RiskPercent == 0 {
		config.RiskPercent = 0.003 // 0.3% - optimized for low drawdown
	}
	if config.MaxPositionCap == 0 {
		config.MaxPositionCap = config.StartBalance * 10
	}
	if config.SlippagePercent == 0 {
		config.SlippagePercent = 0.0015 // 0.15%
	}
	if config.FeePercent == 0 {
		config.FeePercent = 0.001 // 0.1%
	}
	if config.MinWindow == 0 {
		config.MinWindow = 100
	}
	if config.MaxWindow == 0 {
		config.MaxWindow = 200
	}
	if config.MonteCarloRuns == 0 {
		config.MonteCarloRuns = 1000
	}
	if config.MaxTradesPerDay == 0 {
		config.MaxTradesPerDay = 20
	}
	if config.WindowType == "" {
		config.WindowType = "expanding"
	}
}

// runStandardUnified - Standard backtest with all features
func runStandardUnified(config UnifiedBacktestConfig, candles []Candle) (*UnifiedBacktestResult, error) {
	result := &UnifiedBacktestResult{
		StartBalance:            config.StartBalance,
		FinalBalance:            config.StartBalance,
		PeakBalance:             config.StartBalance,
		Trades:                  []Trade{},
		ExitReasons:             make(map[string]int),
		PerformanceByVolatility: make(map[string]float64),
		PerformanceByTrend:      make(map[string]float64),
		StrategyName:            config.Strategy,
		WindowType:              config.WindowType,
	}
	
	skipAhead := 5
	tradesThisDay := 0
	currentDay := ""
	consecutiveLosses := 0
	
	// Simulate trading
	for i := config.MinWindow; i < len(candles)-50; i++ {
		// Calculate window based on type
		var dataWindow []Candle
		switch config.WindowType {
		case "expanding":
			windowSize := minIntUnified(i, config.MaxWindow)
			dataWindow = candles[i-windowSize : i]
		case "rolling":
			windowSize := config.MaxWindow
			if i < windowSize {
				windowSize = i
			}
			dataWindow = candles[i-windowSize : i]
		default:
			dataWindow = candles[i-config.MinWindow : i]
		}
		
		futureData := candles[i : minIntUnified(i+50, len(candles))]
		
		// Daily trade limit check
		candleTime := time.Unix(candles[i].Timestamp/1000, 0)
		candleDay := candleTime.Format("2006-01-02")
		if candleDay != currentDay {
			currentDay = candleDay
			tradesThisDay = 0
		}
		if tradesThisDay >= config.MaxTradesPerDay {
			continue
		}
		
		// Consecutive loss check
		if config.MaxConsecutiveLoss > 0 && consecutiveLosses >= config.MaxConsecutiveLoss {
			continue
		}
		
		// Trading hours filter
		if config.TradingHoursOnly && !shouldTradeAtTimeUnified(candles[i].Timestamp) {
			continue
		}
		
		// Volatility filter
		if config.MinVolatility > 0 || config.MaxVolatility > 0 {
			atr := calculateATR(dataWindow, 14)
			avgPrice := (candles[i].High + candles[i].Low) / 2
			volatilityPct := (atr / avgPrice) * 100
			
			if config.MinVolatility > 0 && volatilityPct < config.MinVolatility {
				continue
			}
			if config.MaxVolatility > 0 && volatilityPct > config.MaxVolatility {
				continue
			}
		}
		
		// Generate signal
		usg := &UnifiedSignalGenerator{}
		advSignal := usg.GenerateSignal(dataWindow, config.Strategy)
		
		if advSignal == nil || advSignal.Type == "NONE" {
			continue
		}
		
		// Simulate trade
		trade := simulateTradeUnified(advSignal, futureData, config)
		
		if trade != nil {
			trade.EntryIndex = i
			trade.BalanceAfter = result.FinalBalance + trade.Profit
			
			result.Trades = append(result.Trades, *trade)
			result.TotalTrades++
			tradesThisDay++
			
			if trade.Profit > 0 {
				result.WinningTrades++
				result.TotalProfit += trade.Profit
				consecutiveLosses = 0
			} else {
				result.LosingTrades++
				result.TotalLoss += math.Abs(trade.Profit)
				consecutiveLosses++
			}
			
			// Update balance
			result.FinalBalance += trade.Profit
			
			// Track peak and drawdown
			if result.FinalBalance > result.PeakBalance {
				result.PeakBalance = result.FinalBalance
			}
			
			drawdownAmount := result.PeakBalance - result.FinalBalance
			drawdown := drawdownAmount / config.StartBalance
			if drawdown > result.MaxDrawdown {
				result.MaxDrawdown = drawdown
			}
			
			result.ExitReasons[trade.ExitReason]++
			
			i += skipAhead
		}
	}
	
	// Calculate statistics
	calculateStatsUnified(result)
	
	return result, nil
}

// runWithPartialExits - Professional partial exit logic
func runWithPartialExits(config UnifiedBacktestConfig, candles []Candle) (*UnifiedBacktestResult, error) {
	result := &UnifiedBacktestResult{
		StartBalance:            config.StartBalance,
		FinalBalance:            config.StartBalance,
		PeakBalance:             config.StartBalance,
		Trades:                  []Trade{},
		ExitReasons:             make(map[string]int),
		PerformanceByVolatility: make(map[string]float64),
		PerformanceByTrend:      make(map[string]float64),
		StrategyName:            config.Strategy,
	}
	
	skipAhead := 5
	
	for i := config.MinWindow; i < len(candles)-50; i++ {
		dataWindow := candles[i-config.MinWindow : i]
		futureData := candles[i : minIntUnified(i+50, len(candles))]
		
		usg := &UnifiedSignalGenerator{}
		advSignal := usg.GenerateSignal(dataWindow, config.Strategy)
		
		if advSignal == nil || advSignal.Type == "NONE" {
			continue
		}
		
		// Simulate with partial exits (30%, 30%, 40%)
		trade := simulateTradeWithPartialExitsUnified(advSignal, futureData, config)
		
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
			
			result.FinalBalance += trade.Profit
			
			if result.FinalBalance > result.PeakBalance {
				result.PeakBalance = result.FinalBalance
			}
			
			drawdown := (result.PeakBalance - result.FinalBalance) / config.StartBalance
			if drawdown > result.MaxDrawdown {
				result.MaxDrawdown = drawdown
			}
			
			result.ExitReasons[trade.ExitReason]++
			i += skipAhead
		}
	}
	
	calculateStatsUnified(result)
	return result, nil
}

// runWalkForwardUnified - Walk-forward analysis
func runWalkForwardUnified(config UnifiedBacktestConfig, candles []Candle) (*UnifiedBacktestResult, error) {
	if config.TrainingDays == 0 {
		config.TrainingDays = 60
	}
	if config.TestingDays == 0 {
		config.TestingDays = 30
	}
	
	candlesPerDay := getCandlesPerDayUnified(config.Interval)
	trainingCandles := config.TrainingDays * candlesPerDay
	testingCandles := config.TestingDays * candlesPerDay
	stepSize := testingCandles / 2
	
	aggregatedResult := &UnifiedBacktestResult{
		StartBalance:            config.StartBalance,
		FinalBalance:            config.StartBalance,
		PeakBalance:             config.StartBalance,
		Trades:                  []Trade{},
		ExitReasons:             make(map[string]int),
		PerformanceByVolatility: make(map[string]float64),
		PerformanceByTrend:      make(map[string]float64),
		StrategyName:            config.Strategy,
		WalkForwardAnalysis:     &WalkForwardAnalysis{PeriodResults: []WalkForwardPeriod{}},
	}
	
	periodNum := 1
	
	for i := 0; i+trainingCandles+testingCandles < len(candles); i += stepSize {
		testStart := i + trainingCandles
		testEnd := testStart + testingCandles
		
		if testEnd > len(candles) {
			break
		}
		
		testData := candles[testStart:testEnd]
		
		testConfig := config
		testConfig.UseWalkForward = false
		testConfig.StartBalance = aggregatedResult.FinalBalance
		
		periodResult, err := runStandardUnified(testConfig, testData)
		if err != nil {
			continue
		}
		
		period := WalkForwardPeriod{
			PeriodNum:     periodNum,
			TrainStart:    i,
			TrainEnd:      i + trainingCandles,
			TestStart:     testStart,
			TestEnd:       testEnd,
			WinRate:       periodResult.WinRate,
			ReturnPercent: periodResult.ReturnPercent,
			ProfitFactor:  periodResult.ProfitFactor,
			TotalTrades:   periodResult.TotalTrades,
		}
		aggregatedResult.WalkForwardAnalysis.PeriodResults = append(
			aggregatedResult.WalkForwardAnalysis.PeriodResults, period)
		
		for _, trade := range periodResult.Trades {
			trade.EntryIndex += testStart
			aggregatedResult.Trades = append(aggregatedResult.Trades, trade)
		}
		
		aggregatedResult.TotalTrades += periodResult.TotalTrades
		aggregatedResult.WinningTrades += periodResult.WinningTrades
		aggregatedResult.LosingTrades += periodResult.LosingTrades
		aggregatedResult.TotalProfit += periodResult.TotalProfit
		aggregatedResult.TotalLoss += periodResult.TotalLoss
		aggregatedResult.FinalBalance = periodResult.FinalBalance
		
		if aggregatedResult.FinalBalance > aggregatedResult.PeakBalance {
			aggregatedResult.PeakBalance = aggregatedResult.FinalBalance
		}
		
		drawdown := (aggregatedResult.PeakBalance - aggregatedResult.FinalBalance) / aggregatedResult.PeakBalance
		if drawdown > aggregatedResult.MaxDrawdown {
			aggregatedResult.MaxDrawdown = drawdown
		}
		
		for reason, count := range periodResult.ExitReasons {
			aggregatedResult.ExitReasons[reason] += count
		}
		
		periodNum++
	}
	
	calculateStatsUnified(aggregatedResult)
	return aggregatedResult, nil
}

// runParallelStrategies - Test multiple strategies in parallel
func runParallelStrategies(config UnifiedBacktestConfig, candles []Candle) (*UnifiedBacktestResult, error) {
	results := make([]*UnifiedBacktestResult, len(config.Strategies))
	var wg sync.WaitGroup
	var mu sync.Mutex
	
	for idx, strategy := range config.Strategies {
		wg.Add(1)
		go func(i int, strat string) {
			defer wg.Done()
			
			cfg := config
			cfg.Strategy = strat
			cfg.EnableParallel = false
			
			result, err := runStandardUnified(cfg, candles)
			if err != nil {
				return
			}
			
			mu.Lock()
			results[i] = result
			mu.Unlock()
		}(idx, strategy)
	}
	
	wg.Wait()
	
	// Find best strategy
	var bestResult *UnifiedBacktestResult
	bestScore := -999999.0
	
	for _, result := range results {
		if result == nil {
			continue
		}
		score := calculateStrategyScore(result)
		if score > bestScore {
			bestScore = score
			bestResult = result
		}
	}
	
	return bestResult, nil
}

// simulateTradeUnified - Unified trade simulation
func simulateTradeUnified(signal *AdvancedSignal, futureData []Candle, config UnifiedBacktestConfig) *Trade {
	if signal == nil || len(futureData) == 0 {
		return nil
	}
	
	entry := signal.Entry
	stopLoss := signal.StopLoss
	
	// Calculate position size
	riskAmount := config.StartBalance * config.RiskPercent
	riskDiff := math.Abs(entry - stopLoss)
	if riskDiff == 0 {
		return nil
	}
	
	positionSize := riskAmount / riskDiff
	maxPositionValue := riskAmount * 10
	if positionSize*entry > maxPositionValue {
		positionSize = maxPositionValue / entry
	}
	
	// Apply slippage
	slippage := config.SlippagePercent
	if config.RealisticSlippage {
		volatility := calculateVolatilityUnified(futureData)
		slippage = slippage * (1 + volatility)
	}
	
	if signal.Type == "BUY" {
		entry *= (1 + slippage)
	} else {
		entry *= (1 - slippage)
	}
	
	// Simulate price movement
	for candleIdx, candle := range futureData {
		if signal.Type == "BUY" {
			// Check stop loss
			if candle.Low <= stopLoss {
				profit := (stopLoss - entry) * positionSize
				profit -= math.Abs(profit) * config.FeePercent * 2
				
				return &Trade{
					Type:          signal.Type,
					Entry:         entry,
					Exit:          stopLoss,
					StopLoss:      stopLoss,
					ExitReason:    "Stop Loss",
					CandlesHeld:   candleIdx + 1,
					Profit:        profit,
					ProfitPercent: (profit / riskAmount) * 100,
					RR:            (stopLoss - entry) / (entry - stopLoss),
				}
			}
			
			// Check TP3
			if candle.High >= signal.TP3 {
				profit := (signal.TP3 - entry) * positionSize
				profit -= math.Abs(profit) * config.FeePercent * 2
				
				return &Trade{
					Type:          signal.Type,
					Entry:         entry,
					Exit:          signal.TP3,
					StopLoss:      stopLoss,
					ExitReason:    "Target 3",
					CandlesHeld:   candleIdx + 1,
					Profit:        profit,
					ProfitPercent: (profit / riskAmount) * 100,
					RR:            (signal.TP3 - entry) / (entry - stopLoss),
				}
			}
		} else { // SELL
			if candle.High >= stopLoss {
				profit := (entry - stopLoss) * positionSize
				profit -= math.Abs(profit) * config.FeePercent * 2
				
				return &Trade{
					Type:          signal.Type,
					Entry:         entry,
					Exit:          stopLoss,
					StopLoss:      stopLoss,
					ExitReason:    "Stop Loss",
					CandlesHeld:   candleIdx + 1,
					Profit:        profit,
					ProfitPercent: (profit / riskAmount) * 100,
					RR:            (entry - stopLoss) / (stopLoss - entry),
				}
			}
			
			if candle.Low <= signal.TP3 {
				profit := (entry - signal.TP3) * positionSize
				profit -= math.Abs(profit) * config.FeePercent * 2
				
				return &Trade{
					Type:          signal.Type,
					Entry:         entry,
					Exit:          signal.TP3,
					StopLoss:      stopLoss,
					ExitReason:    "Target 3",
					CandlesHeld:   candleIdx + 1,
					Profit:        profit,
					ProfitPercent: (profit / riskAmount) * 100,
					RR:            (entry - signal.TP3) / (stopLoss - entry),
				}
			}
		}
	}
	
	return nil
}

// simulateTradeWithPartialExitsUnified - Partial exit logic
func simulateTradeWithPartialExitsUnified(signal *AdvancedSignal, futureData []Candle, config UnifiedBacktestConfig) *Trade {
	if signal == nil || len(futureData) == 0 {
		return nil
	}
	
	entry := signal.Entry
	stopLoss := signal.StopLoss
	
	riskAmount := config.StartBalance * config.RiskPercent
	riskDiff := math.Abs(entry - stopLoss)
	if riskDiff == 0 {
		return nil
	}
	
	positionSize := riskAmount / riskDiff
	maxPositionValue := riskAmount * 10
	if positionSize*entry > maxPositionValue {
		positionSize = maxPositionValue / entry
	}
	
	// Apply slippage
	if signal.Type == "BUY" {
		entry *= (1 + config.SlippagePercent)
	} else {
		entry *= (1 - config.SlippagePercent)
	}
	
	remainingPosition := positionSize
	totalProfit := 0.0
	exitReason := ""
	exitPrice := 0.0
	candlesHeld := 0
	
	tp1Percent := 0.30
	tp2Percent := 0.30
	tp3Percent := 0.40
	
	tp1Hit := false
	tp2Hit := false
	
	for candleIdx, candle := range futureData {
		candlesHeld = candleIdx + 1
		
		if signal.Type == "BUY" {
			if candle.Low <= stopLoss {
				profit := (stopLoss - entry) * remainingPosition
				profit -= math.Abs(profit) * config.FeePercent * 2
				totalProfit += profit
				exitReason = "Stop Loss"
				exitPrice = stopLoss
				break
			}
			
			if !tp1Hit && candle.High >= signal.TP1 {
				tp1Hit = true
				exitSize := positionSize * tp1Percent
				profit := (signal.TP1 - entry) * exitSize
				profit -= math.Abs(profit) * config.FeePercent * 2
				totalProfit += profit
				remainingPosition -= exitSize
				stopLoss = entry // Move to breakeven
			}
			
			if tp1Hit && !tp2Hit && candle.High >= signal.TP2 {
				tp2Hit = true
				exitSize := positionSize * tp2Percent
				profit := (signal.TP2 - entry) * exitSize
				profit -= math.Abs(profit) * config.FeePercent * 2
				totalProfit += profit
				remainingPosition -= exitSize
			}
			
			if tp2Hit && candle.High >= signal.TP3 {
				exitSize := positionSize * tp3Percent
				profit := (signal.TP3 - entry) * exitSize
				profit -= math.Abs(profit) * config.FeePercent * 2
				totalProfit += profit
				exitReason = "Target 3"
				exitPrice = signal.TP3
				break
			}
		} else { // SELL
			if candle.High >= stopLoss {
				profit := (entry - stopLoss) * remainingPosition
				profit -= math.Abs(profit) * config.FeePercent * 2
				totalProfit += profit
				exitReason = "Stop Loss"
				exitPrice = stopLoss
				break
			}
			
			if !tp1Hit && candle.Low <= signal.TP1 {
				tp1Hit = true
				exitSize := positionSize * tp1Percent
				profit := (entry - signal.TP1) * exitSize
				profit -= math.Abs(profit) * config.FeePercent * 2
				totalProfit += profit
				remainingPosition -= exitSize
				stopLoss = entry
			}
			
			if tp1Hit && !tp2Hit && candle.Low <= signal.TP2 {
				tp2Hit = true
				exitSize := positionSize * tp2Percent
				profit := (entry - signal.TP2) * exitSize
				profit -= math.Abs(profit) * config.FeePercent * 2
				totalProfit += profit
				remainingPosition -= exitSize
			}
			
			if tp2Hit && candle.Low <= signal.TP3 {
				exitSize := positionSize * tp3Percent
				profit := (entry - signal.TP3) * exitSize
				profit -= math.Abs(profit) * config.FeePercent * 2
				totalProfit += profit
				exitReason = "Target 3"
				exitPrice = signal.TP3
				break
			}
		}
		
		if candleIdx >= len(futureData)-1 {
			currentPrice := candle.Close
			profit := 0.0
			if signal.Type == "BUY" {
				profit = (currentPrice - entry) * remainingPosition
			} else {
				profit = (entry - currentPrice) * remainingPosition
			}
			profit -= math.Abs(profit) * config.FeePercent * 2
			totalProfit += profit
			exitReason = "Timeout"
			exitPrice = currentPrice
			break
		}
	}
	
	rr := 0.0
	if signal.Type == "BUY" {
		rr = (exitPrice - entry) / (entry - signal.StopLoss)
	} else {
		rr = (entry - exitPrice) / (signal.StopLoss - entry)
	}
	
	return &Trade{
		Type:          signal.Type,
		Entry:         entry,
		Exit:          exitPrice,
		StopLoss:      signal.StopLoss,
		ExitReason:    exitReason,
		CandlesHeld:   candlesHeld,
		Profit:        totalProfit,
		ProfitPercent: (totalProfit / riskAmount) * 100,
		RR:            rr,
	}
}

// calculateStatsUnified - Calculate all statistics
func calculateStatsUnified(result *UnifiedBacktestResult) {
	if result.TotalTrades > 0 {
		result.WinRate = (float64(result.WinningTrades) / float64(result.TotalTrades)) * 100
	}
	
	result.NetProfit = result.TotalProfit - result.TotalLoss
	
	if result.StartBalance > 0 {
		result.ReturnPercent = (result.NetProfit / result.StartBalance) * 100
	}
	
	if result.TotalLoss > 0 {
		result.ProfitFactor = result.TotalProfit / result.TotalLoss
	}
	
	result.MaxDrawdown = result.MaxDrawdown * 100
	
	// Calculate detailed metrics
	totalRR := 0.0
	returns := []float64{}
	downsideReturns := []float64{}
	currentConsecutiveLosses := 0
	maxConsecutiveLosses := 0
	winStreak := 0
	lossStreak := 0
	maxWinStreak := 0
	maxLossStreak := 0
	totalWin := 0.0
	totalLoss := 0.0
	
	for _, trade := range result.Trades {
		totalRR += trade.RR
		tradeReturn := (trade.Profit / result.StartBalance) * 100
		returns = append(returns, tradeReturn)
		
		if trade.Profit > 0 {
			winStreak++
			lossStreak = 0
			totalWin += trade.Profit
			if trade.Profit > result.LargestWin {
				result.LargestWin = trade.Profit
			}
			if winStreak > maxWinStreak {
				maxWinStreak = winStreak
			}
			if currentConsecutiveLosses > maxConsecutiveLosses {
				maxConsecutiveLosses = currentConsecutiveLosses
			}
			currentConsecutiveLosses = 0
		} else {
			lossStreak++
			winStreak = 0
			totalLoss += math.Abs(trade.Profit)
			downsideReturns = append(downsideReturns, tradeReturn)
			if trade.Profit < result.LargestLoss {
				result.LargestLoss = trade.Profit
			}
			if lossStreak > maxLossStreak {
				maxLossStreak = lossStreak
			}
			currentConsecutiveLosses++
		}
	}
	
	if currentConsecutiveLosses > maxConsecutiveLosses {
		maxConsecutiveLosses = currentConsecutiveLosses
	}
	
	result.MaxConsecutiveLosses = maxConsecutiveLosses
	result.WinStreakMax = maxWinStreak
	result.LossStreakMax = maxLossStreak
	
	if result.WinningTrades > 0 {
		result.AverageWin = totalWin / float64(result.WinningTrades)
	}
	if result.LosingTrades > 0 {
		result.AverageLoss = totalLoss / float64(result.LosingTrades)
	}
	
	winProb := float64(result.WinningTrades) / float64(result.TotalTrades)
	lossProb := float64(result.LosingTrades) / float64(result.TotalTrades)
	result.ExpectancyPerTrade = (winProb * result.AverageWin) - (lossProb * result.AverageLoss)
	
	if result.TotalTrades > 0 {
		result.AverageRR = totalRR / float64(result.TotalTrades)
	}
	
	// Sharpe and Sortino ratios
	if len(returns) > 1 {
		mean := calculateMeanUnified(returns)
		stdDev := calculateStdDevUnified(returns, mean)
		
		if stdDev > 0 {
			result.SharpeRatio = (mean / stdDev) * math.Sqrt(252)
		}
		
		if len(downsideReturns) > 0 {
			downsideStdDev := calculateStdDevUnified(downsideReturns, 0)
			if downsideStdDev > 0 {
				result.SortinoRatio = (mean / downsideStdDev) * math.Sqrt(252)
			}
		}
	}
	
	// Calmar ratio
	if result.MaxDrawdown > 0 {
		annualReturn := result.ReturnPercent * (365.0 / 30.0)
		result.CalmarRatio = annualReturn / result.MaxDrawdown
	}
	
	// Recovery factor
	if result.MaxDrawdown > 0 {
		result.RecoveryFactor = result.NetProfit / (result.StartBalance * result.MaxDrawdown / 100)
	}
}

// calculateAdvancedMetricsUnified - Calculate advanced metrics
func calculateAdvancedMetricsUnified(result *UnifiedBacktestResult, candles []Candle) {
	if len(result.Trades) == 0 {
		return
	}
	
	// Time analysis
	totalHours := 0.0
	hourlyPerformance := make(map[int]float64)
	hourlyCount := make(map[int]int)
	
	for _, trade := range result.Trades {
		totalHours += float64(trade.CandlesHeld) * 0.25 // Assuming 15m candles
		
		// Track performance by hour
		if trade.EntryIndex < len(candles) {
			t := time.Unix(candles[trade.EntryIndex].Timestamp/1000, 0)
			hour := t.Hour()
			hourlyPerformance[hour] += trade.Profit
			hourlyCount[hour]++
		}
	}
	
	result.AvgTradeHours = totalHours / float64(len(result.Trades))
	
	// Find best/worst trading hours
	bestHour := 0
	worstHour := 0
	bestProfit := -999999.0
	worstProfit := 999999.0
	
	for hour, profit := range hourlyPerformance {
		if profit > bestProfit {
			bestProfit = profit
			bestHour = hour
		}
		if profit < worstProfit {
			worstProfit = profit
			worstHour = hour
		}
	}
	
	result.BestTradingHour = bestHour
	result.WorstTradingHour = worstHour
}

// runMonteCarloUnified - Monte Carlo simulation
func runMonteCarloUnified(trades []Trade, runs int, startBalance float64) *MonteCarloAnalysis {
	if runs == 0 {
		runs = 1000
	}
	
	if len(trades) == 0 {
		return nil
	}
	
	tradeReturns := make([]float64, len(trades))
	for i, trade := range trades {
		tradeReturns[i] = trade.ProfitPercent
	}
	
	simulationResults := make([]float64, runs)
	profitableRuns := 0
	ruinRuns := 0
	
	rand.Seed(time.Now().UnixNano())
	
	for i := 0; i < runs; i++ {
		balance := startBalance
		
		for j := 0; j < len(tradeReturns); j++ {
			randomIdx := rand.Intn(len(tradeReturns))
			returnPct := tradeReturns[randomIdx]
			balance += balance * (returnPct / 100)
			
			if balance < startBalance*0.5 {
				ruinRuns++
				break
			}
		}
		
		returnPct := ((balance - startBalance) / startBalance) * 100
		simulationResults[i] = returnPct
		
		if returnPct > 0 {
			profitableRuns++
		}
	}
	
	sort.Float64s(simulationResults)
	
	mc := &MonteCarloAnalysis{
		Runs:              runs,
		MeanReturn:        calculateMeanUnified(simulationResults),
		MedianReturn:      simulationResults[runs/2],
		StdDeviation:      calculateStdDevUnified(simulationResults, calculateMeanUnified(simulationResults)),
		BestCase:          simulationResults[runs-1],
		WorstCase:         simulationResults[0],
		Percentile5:       simulationResults[int(float64(runs)*0.05)],
		Percentile25:      simulationResults[int(float64(runs)*0.25)],
		Percentile75:      simulationResults[int(float64(runs)*0.75)],
		Percentile95:      simulationResults[int(float64(runs)*0.95)],
		ProbabilityProfit: float64(profitableRuns) / float64(runs) * 100,
		ProbabilityRuin:   float64(ruinRuns) / float64(runs) * 100,
		ExpectedReturn:    calculateMeanUnified(simulationResults),
	}
	
	return mc
}

// runStressTestUnified - Stress test under extreme conditions
func runStressTestUnified(config UnifiedBacktestConfig, candles []Candle) *StressTestAnalysis {
	sta := &StressTestAnalysis{}
	
	// Test high volatility periods
	highVolCandles := []Candle{}
	lowVolCandles := []Candle{}
	
	for i := 20; i < len(candles); i++ {
		atr := calculateATR(candles[i-20:i], 14)
		avgPrice := (candles[i].High + candles[i].Low) / 2
		volatilityPct := (atr / avgPrice) * 100
		
		if volatilityPct > 2.0 {
			highVolCandles = append(highVolCandles, candles[i])
		} else if volatilityPct < 0.5 {
			lowVolCandles = append(lowVolCandles, candles[i])
		}
	}
	
	if len(highVolCandles) > 100 {
		result, _ := runStandardUnified(config, highVolCandles)
		if result != nil {
			sta.HighVolatilityReturn = result.ReturnPercent
		}
	}
	
	if len(lowVolCandles) > 100 {
		result, _ := runStandardUnified(config, lowVolCandles)
		if result != nil {
			sta.LowVolatilityReturn = result.ReturnPercent
		}
	}
	
	// Simulate crash scenario
	crashCandles := simulateMarketCrashUnified(candles, -30)
	crashResult, _ := runStandardUnified(config, crashCandles)
	if crashResult != nil {
		sta.CrashScenarioReturn = crashResult.ReturnPercent
	}
	
	// Simulate rally scenario
	rallyCandles := simulateMarketRallyUnified(candles, 50)
	rallyResult, _ := runStandardUnified(config, rallyCandles)
	if rallyResult != nil {
		sta.RallyScenarioReturn = rallyResult.ReturnPercent
	}
	
	return sta
}

// printUnifiedSummary - Print comprehensive summary
func printUnifiedSummary(result *UnifiedBacktestResult) {
	log.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("📊 UNIFIED BACKTEST RESULTS")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	log.Printf("\n🏆 STRATEGY: %s", result.StrategyName)
	log.Printf("⏱️  Duration: %s", result.Duration)
	
	log.Println("\n💰 PERFORMANCE:")
	log.Printf("  Start Balance:    $%.2f", result.StartBalance)
	log.Printf("  Final Balance:    $%.2f", result.FinalBalance)
	log.Printf("  Net Profit:       $%.2f", result.NetProfit)
	log.Printf("  Return:           %.2f%%", result.ReturnPercent)
	log.Printf("  Profit Factor:    %.2f", result.ProfitFactor)
	
	log.Println("\n📈 TRADE STATISTICS:")
	log.Printf("  Total Trades:     %d", result.TotalTrades)
	log.Printf("  Winning Trades:   %d", result.WinningTrades)
	log.Printf("  Losing Trades:    %d", result.LosingTrades)
	log.Printf("  Win Rate:         %.2f%%", result.WinRate)
	log.Printf("  Average RR:       %.2f", result.AverageRR)
	
	log.Println("\n⚠️  RISK METRICS:")
	log.Printf("  Max Drawdown:     %.2f%%", result.MaxDrawdown)
	log.Printf("  Sharpe Ratio:     %.2f", result.SharpeRatio)
	log.Printf("  Sortino Ratio:    %.2f", result.SortinoRatio)
	log.Printf("  Calmar Ratio:     %.2f", result.CalmarRatio)
	log.Printf("  Recovery Factor:  %.2f", result.RecoveryFactor)
	log.Printf("  Max Consecutive Losses: %d", result.MaxConsecutiveLosses)
	
	log.Println("\n💵 WIN/LOSS ANALYSIS:")
	log.Printf("  Average Win:      $%.2f", result.AverageWin)
	log.Printf("  Average Loss:     $%.2f", result.AverageLoss)
	log.Printf("  Largest Win:      $%.2f", result.LargestWin)
	log.Printf("  Largest Loss:     $%.2f", result.LargestLoss)
	log.Printf("  Expectancy/Trade: $%.2f", result.ExpectancyPerTrade)
	log.Printf("  Win Streak Max:   %d", result.WinStreakMax)
	log.Printf("  Loss Streak Max:  %d", result.LossStreakMax)
	
	if len(result.ExitReasons) > 0 {
		log.Println("\n🎯 EXIT REASONS:")
		for reason, count := range result.ExitReasons {
			percentage := (float64(count) / float64(result.TotalTrades)) * 100
			log.Printf("  %-20s: %d (%.1f%%)", reason, count, percentage)
		}
	}
	
	if result.MonteCarloResults != nil {
		mc := result.MonteCarloResults
		log.Println("\n🎲 MONTE CARLO ANALYSIS:")
		log.Printf("  Runs:             %d", mc.Runs)
		log.Printf("  Mean Return:      %.2f%%", mc.MeanReturn)
		log.Printf("  Median Return:    %.2f%%", mc.MedianReturn)
		log.Printf("  Best Case:        %.2f%%", mc.BestCase)
		log.Printf("  Worst Case:       %.2f%%", mc.WorstCase)
		log.Printf("  95%% Confidence:   %.2f%% to %.2f%%", mc.Percentile5, mc.Percentile95)
		log.Printf("  Probability Profit: %.1f%%", mc.ProbabilityProfit)
		log.Printf("  Probability Ruin:   %.1f%%", mc.ProbabilityRuin)
	}
	
	if result.StressTestResults != nil {
		st := result.StressTestResults
		log.Println("\n🔥 STRESS TEST RESULTS:")
		log.Printf("  High Volatility:  %.2f%%", st.HighVolatilityReturn)
		log.Printf("  Low Volatility:   %.2f%%", st.LowVolatilityReturn)
		log.Printf("  Crash Scenario:   %.2f%%", st.CrashScenarioReturn)
		log.Printf("  Rally Scenario:   %.2f%%", st.RallyScenarioReturn)
	}
	
	if result.WalkForwardAnalysis != nil && len(result.WalkForwardAnalysis.PeriodResults) > 0 {
		wf := result.WalkForwardAnalysis
		log.Println("\n🚶 WALK-FORWARD ANALYSIS:")
		log.Printf("  Periods:          %d", len(wf.PeriodResults))
		log.Printf("  In-Sample WR:     %.2f%%", wf.InSampleWinRate)
		log.Printf("  Out-Sample WR:    %.2f%%", wf.OutOfSampleWinRate)
		log.Printf("  Consistency:      %.2f", wf.Consistency)
		log.Printf("  Overfitting Score: %.2f", wf.OverfittingScore)
	}
	
	// Performance rating
	log.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	rating := calculatePerformanceRating(result)
	log.Printf("🏆 OVERALL RATING: %s", rating)
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

// Helper functions
func calculateStrategyScore(result *UnifiedBacktestResult) float64 {
	return (result.WinRate * 0.3) + (result.ReturnPercent * 0.3) + 
		(result.ProfitFactor * 20) - (result.MaxDrawdown * 50)
}

func calculatePerformanceRating(result *UnifiedBacktestResult) string {
	score := calculateStrategyScore(result)
	
	if score > 100 && result.WinRate >= 70 && result.ProfitFactor >= 2.5 {
		return "🔥 EXCEPTIONAL - World-class performance!"
	} else if score > 50 && result.WinRate >= 60 && result.ProfitFactor >= 2.0 {
		return "⭐ EXCELLENT - Professional-grade strategy"
	} else if score > 20 && result.WinRate >= 55 && result.ProfitFactor >= 1.5 {
		return "✅ GOOD - Profitable and reliable"
	} else if score > 0 && result.WinRate >= 50 {
		return "⚠️  MODERATE - Needs optimization"
	} else {
		return "❌ POOR - Requires significant improvement"
	}
}

func minIntUnified(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func calculateMeanUnified(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func calculateStdDevUnified(values []float64, mean float64) float64 {
	if len(values) == 0 {
		return 0
	}
	variance := 0.0
	for _, v := range values {
		variance += math.Pow(v-mean, 2)
	}
	return math.Sqrt(variance / float64(len(values)))
}

func calculateVolatilityUnified(candles []Candle) float64 {
	if len(candles) < 2 {
		return 0
	}
	changes := make([]float64, len(candles)-1)
	for i := 1; i < len(candles); i++ {
		changes[i-1] = math.Abs((candles[i].Close - candles[i-1].Close) / candles[i-1].Close)
	}
	sum := 0.0
	for _, change := range changes {
		sum += change
	}
	return sum / float64(len(changes))
}

func shouldTradeAtTimeUnified(timestamp int64) bool {
	t := time.Unix(timestamp/1000, 0).UTC()
	hour := t.Hour()
	weekday := t.Weekday()
	
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}
	if hour < 8 || hour > 20 {
		return false
	}
	if hour == 8 || hour == 20 {
		return false
	}
	return true
}

func getCandlesPerDayUnified(interval string) int {
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
		return 24
	}
}

func simulateMarketCrashUnified(candles []Candle, dropPercent float64) []Candle {
	crashed := make([]Candle, len(candles))
	copy(crashed, candles)
	multiplier := 1.0 + (dropPercent / 100.0)
	for i := range crashed {
		crashed[i].Open *= multiplier
		crashed[i].High *= multiplier
		crashed[i].Low *= multiplier
		crashed[i].Close *= multiplier
	}
	return crashed
}

func simulateMarketRallyUnified(candles []Candle, risePercent float64) []Candle {
	rallied := make([]Candle, len(candles))
	copy(rallied, candles)
	multiplier := 1.0 + (risePercent / 100.0)
	for i := range rallied {
		rallied[i].Open *= multiplier
		rallied[i].High *= multiplier
		rallied[i].Low *= multiplier
		rallied[i].Close *= multiplier
	}
	return rallied
}
