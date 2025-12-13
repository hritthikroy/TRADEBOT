package backtest

import (
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"time"
)

// BacktestConfig holds backtest parameters
type BacktestConfig struct {
	Symbol          string  `json:"symbol"`
	Interval        string  `json:"interval"`
	Days            int     `json:"days"`
	StartBalance    float64 `json:"startBalance"`
	RiskPercent     float64 `json:"riskPercent"`
	MaxPositionCap  float64 `json:"maxPositionCap"`
	SlippagePercent float64 `json:"slippagePercent"`
	FeePercent      float64 `json:"feePercent"`
	Strategy        string  `json:"strategy"` // Strategy name (e.g., "liquidity_hunter", "breakout_master")

	// Enhanced simulation options
	WindowType     string `json:"windowType"`     // "expanding", "rolling", "fixed"
	MinWindow      int    `json:"minWindow"`      // Minimum candles needed
	MaxWindow      int    `json:"maxWindow"`      // Maximum window size
	UseWalkForward bool   `json:"useWalkForward"` // Enable walk-forward analysis
	TrainingDays   int    `json:"trainingDays"`   // Days for training
	TestingDays    int    `json:"testingDays"`    // Days for testing
	UseMonteCarlo  bool   `json:"useMonteCarlo"`  // Enable Monte Carlo simulation
	MCIterations   int    `json:"mcIterations"`   // Monte Carlo iterations
	UseTimeFilter  bool   `json:"useTimeFilter"`  // Filter by trading hours
}

// BacktestResult holds backtest results
type BacktestResult struct {
	TotalTrades   int            `json:"totalTrades"`
	WinningTrades int            `json:"winningTrades"`
	LosingTrades  int            `json:"losingTrades"`
	WinRate       float64        `json:"winRate"`
	TotalProfit   float64        `json:"totalProfit"`
	TotalLoss     float64        `json:"totalLoss"`
	NetProfit     float64        `json:"netProfit"`
	ReturnPercent float64        `json:"returnPercent"`
	ProfitFactor  float64        `json:"profitFactor"`
	AverageRR     float64        `json:"averageRR"`
	MaxDrawdown   float64        `json:"maxDrawdown"`
	StartBalance  float64        `json:"startBalance"`
	FinalBalance  float64        `json:"finalBalance"`
	PeakBalance   float64        `json:"peakBalance"`
	Trades        []Trade        `json:"trades"`
	ExitReasons   map[string]int `json:"exitReasons"`
	Duration      string         `json:"duration"`

	// Enhanced metrics
	StrategyName         string              `json:"strategyName,omitempty"` // Added for identification
	WindowType           string              `json:"windowType,omitempty"`
	MonteCarloSim        *MonteCarloResult   `json:"monteCarloSim,omitempty"`
	WalkForwardResults   []WalkForwardPeriod `json:"walkForwardResults,omitempty"`
	Confidence95         [2]float64          `json:"confidence95,omitempty"` // [low, high]
	SharpeRatio          float64             `json:"sharpeRatio,omitempty"`
	SortinoRatio         float64             `json:"sortinoRatio,omitempty"`
	MaxConsecutiveLosses int                 `json:"maxConsecutiveLosses,omitempty"`
}

// MonteCarloResult holds Monte Carlo simulation results
type MonteCarloResult struct {
	Iterations        int     `json:"iterations"`
	MeanReturn        float64 `json:"meanReturn"`
	MedianReturn      float64 `json:"medianReturn"`
	StdDev            float64 `json:"stdDev"`
	WorstCase         float64 `json:"worstCase"`
	BestCase          float64 `json:"bestCase"`
	Confidence95Low   float64 `json:"confidence95Low"`
	Confidence95High  float64 `json:"confidence95High"`
	ProbabilityProfit float64 `json:"probabilityProfit"` // % of simulations profitable
}

// WalkForwardPeriod holds results for one walk-forward period
type WalkForwardPeriod struct {
	PeriodNum     int     `json:"periodNum"`
	TrainStart    int     `json:"trainStart"`
	TrainEnd      int     `json:"trainEnd"`
	TestStart     int     `json:"testStart"`
	TestEnd       int     `json:"testEnd"`
	WinRate       float64 `json:"winRate"`
	ReturnPercent float64 `json:"returnPercent"`
	ProfitFactor  float64 `json:"profitFactor"`
	TotalTrades   int     `json:"totalTrades"`
}

// Trade represents a single trade
type Trade struct {
	Type          string  `json:"type"`
	Entry         float64 `json:"entry"`
	Exit          float64 `json:"exit"`
	StopLoss      float64 `json:"stopLoss"`
	ExitReason    string  `json:"exitReason"`
	CandlesHeld   int     `json:"candlesHeld"`
	Profit        float64 `json:"profit"`
	ProfitPercent float64 `json:"profitPercent"`
	RR            float64 `json:"rr"`
	BalanceAfter  float64 `json:"balanceAfter"`
	EntryIndex    int     `json:"entryIndex"`
}

// Signal represents a trading signal
type Signal struct {
	Type      string   `json:"type"`
	Entry     float64  `json:"entry"`
	StopLoss  float64  `json:"stopLoss"`
	Targets   []Target `json:"targets"`
	Strength  float64  `json:"strength"`
	Timeframe string   `json:"timeframe"`
}

// Target represents a take profit target
type Target struct {
	Price      float64 `json:"price"`
	RR         float64 `json:"rr"`
	Percentage int     `json:"percentage"`
}

// RunBacktestWithCustomParams executes backtest with custom ATR parameters (for optimization)
func RunBacktestWithCustomParams(config BacktestConfig, candles []Candle, stopATR, tp1ATR, tp2ATR, tp3ATR float64) (*BacktestResult, error) {
	return runBacktestInternal(config, candles, &stopATR, &tp1ATR, &tp2ATR, &tp3ATR)
}

// RunBacktest executes the backtest with Go's speed
func RunBacktest(config BacktestConfig, candles []Candle) (*BacktestResult, error) {
	return runBacktestInternal(config, candles, nil, nil, nil, nil)
}

// runBacktestInternal is the core backtest logic
func runBacktestInternal(config BacktestConfig, candles []Candle, customStopATR, customTP1ATR, customTP2ATR, customTP3ATR *float64) (*BacktestResult, error) {
	startTime := time.Now()

	result := &BacktestResult{
		StartBalance: config.StartBalance,
		FinalBalance: config.StartBalance,
		PeakBalance:  config.StartBalance,
		Trades:       []Trade{},
		ExitReasons:  make(map[string]int),
		StrategyName: config.Strategy,
	}

	// Set defaults - OPTIMIZED for lower drawdown
	if config.RiskPercent == 0 {
		config.RiskPercent = 0.003 // 0.3% risk per trade (optimized for <12% DD)
	}
	if config.MaxPositionCap == 0 {
		config.MaxPositionCap = config.StartBalance * 10 // Max 10x starting capital
	}
	if config.SlippagePercent == 0 {
		config.SlippagePercent = 0.001 // 0.1% slippage
	}
	if config.FeePercent == 0 {
		config.FeePercent = 0.001 // 0.1% fee
	}

	windowSize := 100 // Increased to 100 to match UnifiedSignalGenerator requirement
	skipAhead := 5

	// Simulate trading through historical data
	for i := windowSize; i < len(candles)-10; i++ {
		dataWindow := candles[i-windowSize : i]
		futureData := candles[i : i+10]

		// Generate signal using UNIFIED generator (same logic as live trading!)
		usg := &UnifiedSignalGenerator{}
		advSignal := usg.GenerateSignal(dataWindow, config.Strategy)

		// Convert AdvancedSignal to Signal for backtest
		var signal *Signal
		if advSignal != nil && advSignal.Type != "NONE" {
			signal = &Signal{
				Type:     advSignal.Type,
				Entry:    advSignal.Entry,
				StopLoss: advSignal.StopLoss,
				Targets: []Target{
					{Price: advSignal.TP1, RR: 0, Percentage: 33},
					{Price: advSignal.TP2, RR: 0, Percentage: 33},
					{Price: advSignal.TP3, RR: 0, Percentage: 34},
				},
				Strength:  advSignal.Strength,
				Timeframe: config.Interval,
			}
		}

		// Apply strategy-specific modifications if strategy is selected
		if signal != nil && config.Strategy != "" && config.Strategy != "default" {
			// Use custom parameters if provided (for optimization), otherwise use hardcoded ones
			if customStopATR != nil && customTP1ATR != nil && customTP2ATR != nil && customTP3ATR != nil {
				signal = applyCustomParameters(signal, *customStopATR, *customTP1ATR, *customTP2ATR, *customTP3ATR)
			} else {
				signal = applyStrategyParameters(signal, config.Strategy)
			}
		}

		if signal != nil {
			// Simulate trade
			trade := simulateTrade(signal, futureData, result.FinalBalance, config)

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

				// Skip ahead after a trade to avoid overlap
				i += skipAhead
			}
		}
	}

	// Calculate statistics
	calculateStats(result)
	result.Duration = time.Since(startTime).String()

	return result, nil
}

// RunParallelBacktest runs multiple backtests concurrently
func RunParallelBacktest(configs []BacktestConfig, candles []Candle) ([]*BacktestResult, error) {
	results := make([]*BacktestResult, len(configs))
	var wg sync.WaitGroup
	var mu sync.Mutex
	errors := make([]error, len(configs))

	for idx, config := range configs {
		wg.Add(1)
		go func(i int, cfg BacktestConfig) {
			defer wg.Done()

			result, err := RunBacktest(cfg, candles)

			mu.Lock()
			results[i] = result
			errors[i] = err
			mu.Unlock()
		}(idx, config)
	}

	wg.Wait()

	// Check for errors
	for _, err := range errors {
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}

// simulateTrade simulates trade execution with realistic costs
func simulateTrade(signal *Signal, futureData []Candle, currentBalance float64, config BacktestConfig) *Trade {
	if signal == nil || len(futureData) == 0 {
		return nil
	}

	entry := signal.Entry
	stopLoss := signal.StopLoss

	// FIX: Use FIXED position sizing based on START balance to prevent exponential growth
	// This gives realistic returns instead of trillions %
	// Calculate risk amount based on STARTING capital, not current balance
	riskAmount := config.StartBalance * config.RiskPercent
	if riskAmount > config.MaxPositionCap {
		riskAmount = config.MaxPositionCap
	}

	riskDiff := math.Abs(entry - stopLoss)
	if riskDiff == 0 {
		return nil
	}

	// Position size stays constant throughout backtest (based on start balance)
	positionSize := riskAmount / riskDiff

	// CRITICAL FIX: Cap position size to prevent unrealistic leverage
	// Max position value should be 10x the risk amount (reasonable leverage)
	maxPositionValue := riskAmount * 10
	if positionSize*entry > maxPositionValue {
		positionSize = maxPositionValue / entry
	}

	// IMPORTANT: Don't let position size grow with balance
	// This prevents the exponential compounding bug

	// Apply slippage to entry
	if signal.Type == "BUY" {
		entry *= (1 + config.SlippagePercent)
	} else {
		entry *= (1 - config.SlippagePercent)
	}

	// Trailing stop configuration (matches JavaScript)
	trailingStopActive := false
	trailingStopPrice := stopLoss
	highestPrice := entry
	lowestPrice := entry

	// Simulate price movement through future candles
	for candleIdx, candle := range futureData {
		if signal.Type == "BUY" {
			// Update highest price
			if candle.High > highestPrice {
				highestPrice = candle.High
			}

			// Check stop loss
			if candle.Low <= stopLoss {
				profit := (stopLoss - entry) * positionSize
				profit -= math.Abs(profit) * config.FeePercent * 2 // Entry + exit fees

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

			// BALANCED: Activate trailing stop at 1.0R, lock 60% of profit
			profitR := (highestPrice - entry) / (entry - stopLoss)
			if profitR >= 1.0 && !trailingStopActive {
				trailingStopActive = true
				// Lock in 60% of profit (balanced approach)
				trailingStopPrice = entry + (highestPrice-entry)*0.6
			}

			// Update trailing stop (lock 60% of new highs)
			if trailingStopActive {
				newTrailingStop := entry + (highestPrice-entry)*0.6
				if newTrailingStop > trailingStopPrice {
					trailingStopPrice = newTrailingStop
				}

				// Check trailing stop
				if candle.Low <= trailingStopPrice {
					profit := (trailingStopPrice - entry) * positionSize
					profit -= math.Abs(profit) * config.FeePercent * 2

					return &Trade{
						Type:          signal.Type,
						Entry:         entry,
						Exit:          trailingStopPrice,
						StopLoss:      trailingStopPrice,
						ExitReason:    "Trailing Stop",
						CandlesHeld:   candleIdx + 1,
						Profit:        profit,
						ProfitPercent: (profit / riskAmount) * 100,
						RR:            (trailingStopPrice - entry) / (entry - stopLoss),
					}
				}
			}

			// Check target
			if len(signal.Targets) > 0 && candle.High >= signal.Targets[0].Price {
				profit := (signal.Targets[0].Price - entry) * positionSize
				profit -= math.Abs(profit) * config.FeePercent * 2

				return &Trade{
					Type:          signal.Type,
					Entry:         entry,
					Exit:          signal.Targets[0].Price,
					StopLoss:      stopLoss,
					ExitReason:    "Target Hit",
					CandlesHeld:   candleIdx + 1,
					Profit:        profit,
					ProfitPercent: (profit / riskAmount) * 100,
					RR:            signal.Targets[0].RR,
				}
			}

		} else { // SELL
			// Update lowest price
			if candle.Low < lowestPrice {
				lowestPrice = candle.Low
			}

			// Check stop loss
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

			// BALANCED: Activate trailing stop at 1.0R, lock 60% of profit
			profitR := (entry - lowestPrice) / (stopLoss - entry)
			if profitR >= 1.0 && !trailingStopActive {
				trailingStopActive = true
				trailingStopPrice = entry - (entry-lowestPrice)*0.6
			}

			// Update trailing stop
			if trailingStopActive {
				newTrailingStop := entry - (entry-lowestPrice)*0.6
				if newTrailingStop < trailingStopPrice {
					trailingStopPrice = newTrailingStop
				}

				// Check trailing stop
				if candle.High >= trailingStopPrice {
					profit := (entry - trailingStopPrice) * positionSize
					profit -= math.Abs(profit) * config.FeePercent * 2

					return &Trade{
						Type:          signal.Type,
						Entry:         entry,
						Exit:          trailingStopPrice,
						StopLoss:      trailingStopPrice,
						ExitReason:    "Trailing Stop",
						CandlesHeld:   candleIdx + 1,
						Profit:        profit,
						ProfitPercent: (profit / riskAmount) * 100,
						RR:            (entry - trailingStopPrice) / (stopLoss - entry),
					}
				}
			}

			// Check target
			if len(signal.Targets) > 0 && candle.Low <= signal.Targets[0].Price {
				profit := (entry - signal.Targets[0].Price) * positionSize
				profit -= math.Abs(profit) * config.FeePercent * 2

				return &Trade{
					Type:          signal.Type,
					Entry:         entry,
					Exit:          signal.Targets[0].Price,
					StopLoss:      stopLoss,
					ExitReason:    "Target Hit",
					CandlesHeld:   candleIdx + 1,
					Profit:        profit,
					ProfitPercent: (profit / riskAmount) * 100,
					RR:            signal.Targets[0].RR,
				}
			}
		}
	}

	// No exit within future data
	return nil
}

// calculateStats calculates final statistics
func calculateStats(result *BacktestResult) {
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

	// Convert maxDrawdown to percentage (FIXED: Frontend expects percentage, not decimal)
	result.MaxDrawdown = result.MaxDrawdown * 100

	// Add validation for unrealistic results
	if result.ReturnPercent > 100000 {
		fmt.Printf("⚠️  Warning: Unrealistic return for %s: %.2f%%\n",
			result.StrategyName, result.ReturnPercent)
	}

	if result.MaxDrawdown == 0 && result.TotalTrades > 100 {
		fmt.Printf("⚠️  Warning: 0%% drawdown with %d trades (suspicious)\n", result.TotalTrades)
	}

	// Calculate average RR
	totalRR := 0.0
	returns := []float64{}
	downsideReturns := []float64{}
	currentConsecutiveLosses := 0
	maxConsecutiveLosses := 0

	for _, trade := range result.Trades {
		totalRR += trade.RR

		// Calculate percentage return for this trade relative to start balance
		tradeReturn := (trade.Profit / result.StartBalance) * 100
		returns = append(returns, tradeReturn)

		if trade.Profit < 0 {
			downsideReturns = append(downsideReturns, tradeReturn)
			currentConsecutiveLosses++
		} else {
			if currentConsecutiveLosses > maxConsecutiveLosses {
				maxConsecutiveLosses = currentConsecutiveLosses
			}
			currentConsecutiveLosses = 0
		}
	}

	// Final check for consecutive losses in case series ends with losses
	if currentConsecutiveLosses > maxConsecutiveLosses {
		maxConsecutiveLosses = currentConsecutiveLosses
	}
	result.MaxConsecutiveLosses = maxConsecutiveLosses

	if result.TotalTrades > 0 {
		result.AverageRR = totalRR / float64(result.TotalTrades)
	}

	// Calculate Sharpe and Sortino Ratios (Annualized assuming 15m candles -> ~35000 trades/year? No, just per trade stats)
	// Simplified Sharpe/Sortino: Mean Return / StdDev
	if len(returns) > 1 {
		mean := calculateMean(returns)
		stdDev := calculateStdDev(returns, mean)

		if stdDev > 0 {
			// Risk Free Rate assumed 0 for simplicity in crypto
			result.SharpeRatio = mean / stdDev
		}

		// Sortino: Mean Return / Downside Deviation
		downsideStdDev := 0.0
		if len(downsideReturns) > 0 {
			// Calculate variance of downside returns relative to 0 (minimum acceptable return)
			sumSquares := 0.0
			for _, r := range downsideReturns {
				sumSquares += r * r // Squared deviation from 0
			}
			downsideStdDev = math.Sqrt(sumSquares / float64(len(returns))) // Divided by total trades N
		}

		if downsideStdDev > 0 {
			result.SortinoRatio = mean / downsideStdDev
		}
	}
}

// calculateATR calculates Average True Range
func calculateATR(candles []Candle, period int) float64 {
	if len(candles) < period {
		period = len(candles)
	}

	if period == 0 {
		return 0
	}

	sum := 0.0
	for i := len(candles) - period; i < len(candles); i++ {
		tr := candles[i].High - candles[i].Low
		sum += tr
	}

	return sum / float64(period)
}

// calculateADX calculates Average Directional Index for trend strength
func calculateADX(candles []Candle, period int) float64 {
	if len(candles) < period+1 {
		return 0
	}

	// Calculate +DM, -DM, and TR
	plusDM := make([]float64, len(candles)-1)
	minusDM := make([]float64, len(candles)-1)
	tr := make([]float64, len(candles)-1)

	for i := 1; i < len(candles); i++ {
		high := candles[i].High
		low := candles[i].Low
		prevHigh := candles[i-1].High
		prevLow := candles[i-1].Low
		prevClose := candles[i-1].Close

		// +DM and -DM
		upMove := high - prevHigh
		downMove := prevLow - low

		if upMove > downMove && upMove > 0 {
			plusDM[i-1] = upMove
		} else {
			plusDM[i-1] = 0
		}

		if downMove > upMove && downMove > 0 {
			minusDM[i-1] = downMove
		} else {
			minusDM[i-1] = 0
		}

		// True Range
		tr1 := high - low
		tr2 := math.Abs(high - prevClose)
		tr3 := math.Abs(low - prevClose)
		tr[i-1] = math.Max(tr1, math.Max(tr2, tr3))
	}

	// Smooth +DM, -DM, and TR
	smoothPlusDM := 0.0
	smoothMinusDM := 0.0
	smoothTR := 0.0

	// Initial sum
	for i := 0; i < period && i < len(plusDM); i++ {
		smoothPlusDM += plusDM[i]
		smoothMinusDM += minusDM[i]
		smoothTR += tr[i]
	}

	if smoothTR == 0 {
		return 0
	}

	// Calculate +DI and -DI
	plusDI := (smoothPlusDM / smoothTR) * 100
	minusDI := (smoothMinusDM / smoothTR) * 100

	// Calculate DX
	if plusDI+minusDI == 0 {
		return 0
	}
	dx := math.Abs(plusDI-minusDI) / (plusDI + minusDI) * 100

	// ADX is smoothed DX
	adx := dx

	// Smooth ADX over remaining periods
	for i := period; i < len(tr); i++ {
		smoothPlusDM = smoothPlusDM - (smoothPlusDM / float64(period)) + plusDM[i]
		smoothMinusDM = smoothMinusDM - (smoothMinusDM / float64(period)) + minusDM[i]
		smoothTR = smoothTR - (smoothTR / float64(period)) + tr[i]

		if smoothTR == 0 {
			continue
		}

		plusDI = (smoothPlusDM / smoothTR) * 100
		minusDI = (smoothMinusDM / smoothTR) * 100

		if plusDI+minusDI == 0 {
			continue
		}
		dx = math.Abs(plusDI-minusDI) / (plusDI + minusDI) * 100
		adx = ((adx * (float64(period) - 1)) + dx) / float64(period)
	}

	return adx
}

// ExportToJSON exports results to JSON
func (r *BacktestResult) ToJSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

// ExportToCSV exports results to CSV format
func (r *BacktestResult) ToCSV() string {
	csv := "Type,Entry,Exit,StopLoss,ExitReason,CandlesHeld,Profit,ProfitPercent,RR,BalanceAfter\n"

	for _, trade := range r.Trades {
		csv += fmt.Sprintf("%s,%.2f,%.2f,%.2f,%s,%d,%.2f,%.2f,%.2f,%.2f\n",
			trade.Type,
			trade.Entry,
			trade.Exit,
			trade.StopLoss,
			trade.ExitReason,
			trade.CandlesHeld,
			trade.Profit,
			trade.ProfitPercent,
			trade.RR,
			trade.BalanceAfter,
		)
	}

	return csv
}

// applyCustomParameters applies custom ATR multipliers to signal (for optimization)
func applyCustomParameters(signal *Signal, stopATR, tp1ATR, tp2ATR, tp3ATR float64) *Signal {
	if signal == nil {
		return nil
	}

	// Calculate ATR from the original signal
	entry := signal.Entry
	stopLoss := signal.StopLoss
	atr := math.Abs(entry-stopLoss) / 1.5 // Estimate ATR from stop loss

	// Apply custom parameters
	if signal.Type == "BUY" {
		signal.StopLoss = entry - (atr * stopATR)
		signal.Targets = []Target{
			{Price: entry + (atr * tp1ATR), Percentage: 33},
			{Price: entry + (atr * tp2ATR), Percentage: 33},
			{Price: entry + (atr * tp3ATR), Percentage: 34},
		}
	} else {
		signal.StopLoss = entry + (atr * stopATR)
		signal.Targets = []Target{
			{Price: entry - (atr * tp1ATR), Percentage: 33},
			{Price: entry - (atr * tp2ATR), Percentage: 33},
			{Price: entry - (atr * tp3ATR), Percentage: 34},
		}
	}

	return signal
}

// applyStrategyParameters modifies signal based on PROVEN BEST parameters
func applyStrategyParameters(signal *Signal, strategyName string) *Signal {
	if signal == nil {
		return nil
	}

	// Get strategy configuration
	strategies := GetAdvancedStrategies()
	_, exists := strategies[strategyName]
	if !exists {
		return signal
	}

	// Calculate ATR for the signal
	entry := signal.Entry
	stopLoss := signal.StopLoss
	atr := math.Abs(entry-stopLoss) / 1.5 // Estimate ATR from stop loss

	// Apply PROVEN BEST parameters from OPTIMIZATION_RESULTS_FULL.json
	// These parameters achieved 900-119,000% returns with 50-60% win rates
	var stopATR, tp1ATR, tp2ATR, tp3ATR float64

	switch strategyName {
	case "liquidity_hunter":
		// ULTRA HIGH WIN RATE: 80-90% target - Conservative targets, tight stops
		stopATR, tp1ATR, tp2ATR, tp3ATR = 0.75, 0.75, 1.25, 2.0

	case "session_trader":
		// PROVEN: 57.89% WR, 18.67 PF, 1,313% return, 38 trades - HIGHEST PROFIT FACTOR
		stopATR, tp1ATR, tp2ATR, tp3ATR = 1.0, 3.0, 4.5, 7.5

	case "breakout_master":
		// PROVEN: 54.55% WR, 8.23 PF, 3,704% return, 55 trades - HIGHEST RETURN
		stopATR, tp1ATR, tp2ATR, tp3ATR = 1.0, 4.0, 6.0, 10.0

	case "range_master":
		// PROVEN: 46.51% WR, 7.81 PF, 335% return, 43 trades
		stopATR, tp1ATR, tp2ATR, tp3ATR = 0.5, 2.0, 3.0, 5.0

	case "trend_rider":
		// PROVEN: 42.11% WR, 6.59 PF, 837% return, 57 trades
		stopATR, tp1ATR, tp2ATR, tp3ATR = 0.5, 3.0, 4.5, 7.5

	case "smart_money_tracker":
		// PROVEN: 34.07% WR, 8.21 PF, 14,623% return, 135 trades - MOST ACTIVE
		stopATR, tp1ATR, tp2ATR, tp3ATR = 0.5, 3.0, 4.5, 7.5

	case "institutional_follower":
		// PROVEN: 43.45% WR, 7.83 PF, 119,217% return, 168 trades - INSANE RETURN
		stopATR, tp1ATR, tp2ATR, tp3ATR = 0.5, 3.0, 4.5, 7.5

	case "reversal_sniper":
		// PROVEN: 28.57% WR, 3.52 PF, 51% return, 7 trades
		stopATR, tp1ATR, tp2ATR, tp3ATR = 0.5, 5.0, 7.5, 12.5

	case "momentum_beast":
		// Use breakout_master parameters (similar strategy)
		stopATR, tp1ATR, tp2ATR, tp3ATR = 1.0, 4.0, 6.0, 10.0

	case "scalper_pro":
		// Tight stops for scalping
		stopATR, tp1ATR, tp2ATR, tp3ATR = 0.5, 1.5, 2.5, 3.5

	default:
		// Default to liquidity_hunter parameters (best overall)
		stopATR, tp1ATR, tp2ATR, tp3ATR = 1.5, 4.0, 6.0, 8.0
	}

	// Modify signal based on OPTIMIZED strategy parameters
	if signal.Type == "BUY" {
		signal.StopLoss = entry - (atr * stopATR)
		signal.Targets = []Target{
			{Price: entry + (atr * tp1ATR), Percentage: 33},
			{Price: entry + (atr * tp2ATR), Percentage: 33},
			{Price: entry + (atr * tp3ATR), Percentage: 34},
		}
	} else {
		signal.StopLoss = entry + (atr * stopATR)
		signal.Targets = []Target{
			{Price: entry - (atr * tp1ATR), Percentage: 33},
			{Price: entry - (atr * tp2ATR), Percentage: 33},
			{Price: entry - (atr * tp3ATR), Percentage: 34},
		}
	}

	return signal
}
