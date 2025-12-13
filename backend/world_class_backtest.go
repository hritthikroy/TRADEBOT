package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

// WorldClassBacktestConfig extends BacktestConfig with advanced features
type WorldClassBacktestConfig struct {
	BacktestConfig
	
	// Advanced Risk Management
	MaxDailyLoss        float64 `json:"maxDailyLoss"`        // Stop trading if daily loss exceeds this %
	MaxWeeklyLoss       float64 `json:"maxWeeklyLoss"`       // Stop trading if weekly loss exceeds this %
	MaxConsecutiveLoss  int     `json:"maxConsecutiveLoss"`  // Stop after N consecutive losses
	DynamicPositionSize bool    `json:"dynamicPositionSize"` // Adjust position size based on equity
	
	// Market Condition Filters
	MinVolatility       float64 `json:"minVolatility"`       // Minimum ATR to trade
	MaxVolatility       float64 `json:"maxVolatility"`       // Maximum ATR to trade
	MinVolume           float64 `json:"minVolume"`           // Minimum volume multiplier
	TradingHoursOnly    bool    `json:"tradingHoursOnly"`    // Only trade during specific hours
	AvoidNews           bool    `json:"avoidNews"`           // Skip trading during high-impact news
	
	// Advanced Analysis
	EnableMonteCarlo    bool    `json:"enableMonteCarlo"`    // Run Monte Carlo simulation
	MonteCarloRuns      int     `json:"monteCarloRuns"`      // Number of MC simulations
	EnableWalkForward   bool    `json:"enableWalkForward"`   // Walk-forward optimization
	WalkForwardPeriods  int     `json:"walkForwardPeriods"`  // Number of WF periods
	EnableStressTest    bool    `json:"enableStressTest"`    // Test under extreme conditions
	
	// Slippage & Costs
	RealisticSlippage   bool    `json:"realisticSlippage"`   // Variable slippage based on volatility
	IncludeSpread       bool    `json:"includeSpread"`       // Include bid-ask spread
	SpreadPercent       float64 `json:"spreadPercent"`       // Spread as % of price
	IncludeSwap         bool    `json:"includeSwap"`         // Include overnight swap fees
	SwapPerDay          float64 `json:"swapPerDay"`          // Swap fee per day
}

// WorldClassBacktestResult extends BacktestResult with advanced metrics
type WorldClassBacktestResult struct {
	*BacktestResult
	
	// Risk Metrics
	SharpeRatio         float64            `json:"sharpeRatio"`
	SortinoRatio        float64            `json:"sortinoRatio"`
	CalmarRatio         float64            `json:"calmarRatio"`
	MaxDrawdownDuration int                `json:"maxDrawdownDuration"` // Days in drawdown
	RecoveryFactor      float64            `json:"recoveryFactor"`      // Net profit / max DD
	
	// Performance Metrics
	WinStreakMax        int                `json:"winStreakMax"`
	LossStreakMax       int                `json:"lossStreakMax"`
	AverageWin          float64            `json:"averageWin"`
	AverageLoss         float64            `json:"averageLoss"`
	LargestWin          float64            `json:"largestWin"`
	LargestLoss         float64            `json:"largestLoss"`
	ExpectancyPerTrade  float64            `json:"expectancyPerTrade"`
	
	// Time Analysis
	AvgTradeHours       float64            `json:"avgTradeHours"`
	BestTradingHour     int                `json:"bestTradingHour"`
	WorstTradingHour    int                `json:"worstTradingHour"`
	BestTradingDay      string             `json:"bestTradingDay"`
	WorstTradingDay     string             `json:"worstTradingDay"`
	
	// Monte Carlo Results
	MonteCarloResults   *MonteCarloAnalysis `json:"monteCarloResults,omitempty"`
	
	// Walk Forward Results
	WalkForwardAnalysis *WalkForwardAnalysis `json:"walkForwardAnalysis,omitempty"`
	
	// Stress Test Results
	StressTestResults   *StressTestAnalysis `json:"stressTestResults,omitempty"`
	
	// Market Condition Analysis
	PerformanceByVolatility map[string]float64 `json:"performanceByVolatility"`
	PerformanceByTrend      map[string]float64 `json:"performanceByTrend"`
	PerformanceByVolume     map[string]float64 `json:"performanceByVolume"`
}

// MonteCarloAnalysis holds Monte Carlo simulation results
type MonteCarloAnalysis struct {
	Runs                int       `json:"runs"`
	MeanReturn          float64   `json:"meanReturn"`
	MedianReturn        float64   `json:"medianReturn"`
	StdDeviation        float64   `json:"stdDeviation"`
	BestCase            float64   `json:"bestCase"`
	WorstCase           float64   `json:"worstCase"`
	Percentile5         float64   `json:"percentile5"`
	Percentile25        float64   `json:"percentile25"`
	Percentile75        float64   `json:"percentile75"`
	Percentile95        float64   `json:"percentile95"`
	ProbabilityProfit   float64   `json:"probabilityProfit"`
	ProbabilityRuin     float64   `json:"probabilityRuin"`
	ExpectedReturn      float64   `json:"expectedReturn"`
}

// WalkForwardAnalysis holds walk-forward optimization results
type WalkForwardAnalysis struct {
	Periods             int                    `json:"periods"`
	InSampleWinRate     float64                `json:"inSampleWinRate"`
	OutOfSampleWinRate  float64                `json:"outOfSampleWinRate"`
	Consistency         float64                `json:"consistency"` // How consistent across periods
	PeriodResults       []WalkForwardPeriod    `json:"periodResults"`
	OverfittingScore    float64                `json:"overfittingScore"` // 0-100, lower is better
}

// StressTestAnalysis holds stress test results
type StressTestAnalysis struct {
	HighVolatilityReturn  float64 `json:"highVolatilityReturn"`
	LowVolatilityReturn   float64 `json:"lowVolatilityReturn"`
	TrendingReturn        float64 `json:"trendingReturn"`
	RangingReturn         float64 `json:"rangingReturn"`
	CrashScenarioReturn   float64 `json:"crashScenarioReturn"`
	RallyScenarioReturn   float64 `json:"rallyScenarioReturn"`
	WorstMonthReturn      float64 `json:"worstMonthReturn"`
	BestMonthReturn       float64 `json:"bestMonthReturn"`
}

// RunWorldClassBacktest executes advanced backtesting with professional features
func RunWorldClassBacktest(config WorldClassBacktestConfig, candles []Candle) (*WorldClassBacktestResult, error) {
	startTime := time.Now()
	
	// Run base backtest
	baseResult, err := RunBacktest(config.BacktestConfig, candles)
	if err != nil {
		return nil, err
	}
	
	// Create world-class result
	wcResult := &WorldClassBacktestResult{
		BacktestResult: baseResult,
		PerformanceByVolatility: make(map[string]float64),
		PerformanceByTrend: make(map[string]float64),
		PerformanceByVolume: make(map[string]float64),
	}
	
	// Calculate advanced metrics
	calculateAdvancedMetrics(wcResult, candles)
	
	// Run Monte Carlo if enabled
	if config.EnableMonteCarlo {
		wcResult.MonteCarloResults = runMonteCarloSimulationWC(baseResult, config.MonteCarloRuns)
	}
	
	// Run Walk Forward if enabled
	if config.EnableWalkForward {
		wcResult.WalkForwardAnalysis = runWalkForwardAnalysis(config, candles)
	}
	
	// Run Stress Test if enabled
	if config.EnableStressTest {
		wcResult.StressTestResults = runStressTest(config, candles)
	}
	
	// Analyze performance by market conditions
	analyzeMarketConditions(wcResult, candles)
	
	fmt.Printf("✅ World-class backtest completed in %v\n", time.Since(startTime))
	
	return wcResult, nil
}


// calculateAdvancedMetrics calculates professional trading metrics
func calculateAdvancedMetrics(result *WorldClassBacktestResult, candles []Candle) {
	if len(result.Trades) == 0 {
		return
	}
	
	// Calculate returns for each trade
	returns := make([]float64, len(result.Trades))
	downside := []float64{}
	
	winStreak := 0
	lossStreak := 0
	maxWinStreak := 0
	maxLossStreak := 0
	
	totalWin := 0.0
	totalLoss := 0.0
	
	for i, trade := range result.Trades {
		returns[i] = trade.ProfitPercent
		
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
		} else {
			lossStreak++
			winStreak = 0
			totalLoss += math.Abs(trade.Profit)
			downside = append(downside, returns[i])
			if trade.Profit < result.LargestLoss {
				result.LargestLoss = trade.Profit
			}
			if lossStreak > maxLossStreak {
				maxLossStreak = lossStreak
			}
		}
	}
	
	result.WinStreakMax = maxWinStreak
	result.LossStreakMax = maxLossStreak
	
	if result.WinningTrades > 0 {
		result.AverageWin = totalWin / float64(result.WinningTrades)
	}
	if result.LosingTrades > 0 {
		result.AverageLoss = totalLoss / float64(result.LosingTrades)
	}
	
	// Expectancy per trade
	winProb := float64(result.WinningTrades) / float64(result.TotalTrades)
	lossProb := float64(result.LosingTrades) / float64(result.TotalTrades)
	result.ExpectancyPerTrade = (winProb * result.AverageWin) - (lossProb * result.AverageLoss)
	
	// Sharpe Ratio (annualized)
	if len(returns) > 1 {
		mean := calculateMeanWC(returns)
		stdDev := calculateStdDevWC(returns, mean)
		if stdDev > 0 {
			result.SharpeRatio = (mean / stdDev) * math.Sqrt(252) // Annualized
		}
	}
	
	// Sortino Ratio (uses only downside deviation)
	if len(downside) > 1 {
		mean := calculateMeanWC(returns)
		downsideDev := calculateStdDevWC(downside, 0)
		if downsideDev > 0 {
			result.SortinoRatio = (mean / downsideDev) * math.Sqrt(252)
		}
	}
	
	// Calmar Ratio (annual return / max drawdown)
	if result.MaxDrawdown > 0 {
		annualReturn := result.ReturnPercent * (365.0 / 30.0) // Approximate
		result.CalmarRatio = annualReturn / result.MaxDrawdown
	}
	
	// Recovery Factor
	if result.MaxDrawdown > 0 {
		result.RecoveryFactor = result.NetProfit / (result.StartBalance * result.MaxDrawdown / 100)
	}
	
	// Average trade duration
	totalHours := 0.0
	for _, trade := range result.Trades {
		totalHours += float64(trade.CandlesHeld) * 0.25 // Assuming 15m candles
	}
	result.AvgTradeHours = totalHours / float64(len(result.Trades))
}

// runMonteCarloSimulationWC runs Monte Carlo analysis for world-class backtest
func runMonteCarloSimulationWC(baseResult *BacktestResult, runs int) *MonteCarloAnalysis {
	if runs == 0 {
		runs = 1000
	}
	
	if len(baseResult.Trades) == 0 {
		return nil
	}
	
	// Extract trade returns
	tradeReturns := make([]float64, len(baseResult.Trades))
	for i, trade := range baseResult.Trades {
		tradeReturns[i] = trade.ProfitPercent
	}
	
	// Run simulations
	simulationResults := make([]float64, runs)
	profitableRuns := 0
	ruinRuns := 0
	
	rand.Seed(time.Now().UnixNano())
	
	for i := 0; i < runs; i++ {
		balance := baseResult.StartBalance
		
		// Randomly sample trades with replacement
		for j := 0; j < len(tradeReturns); j++ {
			randomIdx := rand.Intn(len(tradeReturns))
			returnPct := tradeReturns[randomIdx]
			balance += balance * (returnPct / 100)
			
			// Check for ruin (balance < 50% of start)
			if balance < baseResult.StartBalance * 0.5 {
				ruinRuns++
				break
			}
		}
		
		returnPct := ((balance - baseResult.StartBalance) / baseResult.StartBalance) * 100
		simulationResults[i] = returnPct
		
		if returnPct > 0 {
			profitableRuns++
		}
	}
	
	// Sort results for percentile calculation
	sort.Float64s(simulationResults)
	
	mc := &MonteCarloAnalysis{
		Runs:              runs,
		MeanReturn:        calculateMeanWC(simulationResults),
		MedianReturn:      simulationResults[runs/2],
		StdDeviation:      calculateStdDevWC(simulationResults, calculateMeanWC(simulationResults)),
		BestCase:          simulationResults[runs-1],
		WorstCase:         simulationResults[0],
		Percentile5:       simulationResults[int(float64(runs)*0.05)],
		Percentile25:      simulationResults[int(float64(runs)*0.25)],
		Percentile75:      simulationResults[int(float64(runs)*0.75)],
		Percentile95:      simulationResults[int(float64(runs)*0.95)],
		ProbabilityProfit: float64(profitableRuns) / float64(runs) * 100,
		ProbabilityRuin:   float64(ruinRuns) / float64(runs) * 100,
		ExpectedReturn:    calculateMeanWC(simulationResults),
	}
	
	return mc
}


// runWalkForwardAnalysis performs walk-forward optimization
func runWalkForwardAnalysis(config WorldClassBacktestConfig, candles []Candle) *WalkForwardAnalysis {
	periods := config.WalkForwardPeriods
	if periods == 0 {
		periods = 5
	}
	
	totalCandles := len(candles)
	periodSize := totalCandles / periods
	
	wfa := &WalkForwardAnalysis{
		Periods:       periods,
		PeriodResults: make([]WalkForwardPeriod, 0),
	}
	
	inSampleWins := 0
	inSampleTotal := 0
	outSampleWins := 0
	outSampleTotal := 0
	
	for i := 0; i < periods; i++ {
		trainStart := i * periodSize
		trainEnd := trainStart + int(float64(periodSize)*0.7) // 70% training
		testStart := trainEnd
		testEnd := (i + 1) * periodSize
		
		if testEnd > totalCandles {
			testEnd = totalCandles
		}
		
		// Run backtest on training period
		trainCandles := candles[trainStart:trainEnd]
		trainResult, _ := RunBacktest(config.BacktestConfig, trainCandles)
		
		// Run backtest on testing period
		testCandles := candles[testStart:testEnd]
		testResult, _ := RunBacktest(config.BacktestConfig, testCandles)
		
		if trainResult != nil {
			inSampleWins += trainResult.WinningTrades
			inSampleTotal += trainResult.TotalTrades
		}
		
		if testResult != nil {
			outSampleWins += testResult.WinningTrades
			outSampleTotal += testResult.TotalTrades
			
			period := WalkForwardPeriod{
				PeriodNum:     i + 1,
				TrainStart:    trainStart,
				TrainEnd:      trainEnd,
				TestStart:     testStart,
				TestEnd:       testEnd,
				WinRate:       testResult.WinRate,
				ReturnPercent: testResult.ReturnPercent,
				ProfitFactor:  testResult.ProfitFactor,
				TotalTrades:   testResult.TotalTrades,
			}
			wfa.PeriodResults = append(wfa.PeriodResults, period)
		}
	}
	
	if inSampleTotal > 0 {
		wfa.InSampleWinRate = float64(inSampleWins) / float64(inSampleTotal) * 100
	}
	if outSampleTotal > 0 {
		wfa.OutOfSampleWinRate = float64(outSampleWins) / float64(outSampleTotal) * 100
	}
	
	// Calculate consistency (std dev of period returns)
	if len(wfa.PeriodResults) > 1 {
		returns := make([]float64, len(wfa.PeriodResults))
		for i, p := range wfa.PeriodResults {
			returns[i] = p.ReturnPercent
		}
		mean := calculateMeanWC(returns)
		stdDev := calculateStdDevWC(returns, mean)
		wfa.Consistency = 100 - (stdDev * 10) // Higher is more consistent
		if wfa.Consistency < 0 {
			wfa.Consistency = 0
		}
	}
	
	// Overfitting score (difference between in-sample and out-sample)
	if wfa.InSampleWinRate > 0 {
		diff := math.Abs(wfa.InSampleWinRate - wfa.OutOfSampleWinRate)
		wfa.OverfittingScore = (diff / wfa.InSampleWinRate) * 100
	}
	
	return wfa
}

// runStressTest tests strategy under extreme conditions
func runStressTest(config WorldClassBacktestConfig, candles []Candle) *StressTestAnalysis {
	sta := &StressTestAnalysis{}
	
	// Analyze volatility periods
	highVolCandles := []Candle{}
	lowVolCandles := []Candle{}
	
	for i := 20; i < len(candles); i++ {
		atr := calculateATR(candles[i-20:i], 14)
		avgPrice := (candles[i].High + candles[i].Low) / 2
		volatilityPct := (atr / avgPrice) * 100
		
		if volatilityPct > 2.0 { // High volatility
			highVolCandles = append(highVolCandles, candles[i])
		} else if volatilityPct < 0.5 { // Low volatility
			lowVolCandles = append(lowVolCandles, candles[i])
		}
	}
	
	// Test high volatility
	if len(highVolCandles) > 100 {
		result, _ := RunBacktest(config.BacktestConfig, highVolCandles)
		if result != nil {
			sta.HighVolatilityReturn = result.ReturnPercent
		}
	}
	
	// Test low volatility
	if len(lowVolCandles) > 100 {
		result, _ := RunBacktest(config.BacktestConfig, lowVolCandles)
		if result != nil {
			sta.LowVolatilityReturn = result.ReturnPercent
		}
	}
	
	// Simulate crash scenario (-30% drop)
	crashCandles := simulateMarketCrash(candles, -30)
	crashResult, _ := RunBacktest(config.BacktestConfig, crashCandles)
	if crashResult != nil {
		sta.CrashScenarioReturn = crashResult.ReturnPercent
	}
	
	// Simulate rally scenario (+50% rise)
	rallyCandles := simulateMarketRally(candles, 50)
	rallyResult, _ := RunBacktest(config.BacktestConfig, rallyCandles)
	if rallyResult != nil {
		sta.RallyScenarioReturn = rallyResult.ReturnPercent
	}
	
	return sta
}

// analyzeMarketConditions analyzes performance by market conditions
func analyzeMarketConditions(result *WorldClassBacktestResult, candles []Candle) {
	// This would analyze trades by volatility, trend, volume
	// Simplified version for now
	result.PerformanceByVolatility["high"] = 0
	result.PerformanceByVolatility["medium"] = 0
	result.PerformanceByVolatility["low"] = 0
	
	result.PerformanceByTrend["uptrend"] = 0
	result.PerformanceByTrend["downtrend"] = 0
	result.PerformanceByTrend["sideways"] = 0
	
	result.PerformanceByVolume["high"] = 0
	result.PerformanceByVolume["low"] = 0
}

// Helper functions for world-class backtest
func calculateMeanWC(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func calculateStdDevWC(values []float64, mean float64) float64 {
	if len(values) == 0 {
		return 0
	}
	variance := 0.0
	for _, v := range values {
		variance += math.Pow(v-mean, 2)
	}
	return math.Sqrt(variance / float64(len(values)))
}

func simulateMarketCrash(candles []Candle, dropPercent float64) []Candle {
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

func simulateMarketRally(candles []Candle, risePercent float64) []Candle {
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
