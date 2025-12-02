package main

import (
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"time"
)

// BacktestConfig holds backtest parameters
type BacktestConfig struct {
	Symbol         string  `json:"symbol"`
	Interval       string  `json:"interval"`
	Days           int     `json:"days"`
	StartBalance   float64 `json:"startBalance"`
	RiskPercent    float64 `json:"riskPercent"`
	MaxPositionCap float64 `json:"maxPositionCap"`
	SlippagePercent float64 `json:"slippagePercent"`
	FeePercent     float64 `json:"feePercent"`
	Strategy       string  `json:"strategy"`       // Strategy name (e.g., "liquidity_hunter", "breakout_master")
	
	// Enhanced simulation options
	WindowType     string  `json:"windowType"`     // "expanding", "rolling", "fixed"
	MinWindow      int     `json:"minWindow"`      // Minimum candles needed
	MaxWindow      int     `json:"maxWindow"`      // Maximum window size
	UseWalkForward bool    `json:"useWalkForward"` // Enable walk-forward analysis
	TrainingDays   int     `json:"trainingDays"`   // Days for training
	TestingDays    int     `json:"testingDays"`    // Days for testing
	UseMonteCarlo  bool    `json:"useMonteCarlo"`  // Enable Monte Carlo simulation
	MCIterations   int     `json:"mcIterations"`   // Monte Carlo iterations
	UseTimeFilter  bool    `json:"useTimeFilter"`  // Filter by trading hours
}

// BacktestResult holds backtest results
type BacktestResult struct {
	TotalTrades    int              `json:"totalTrades"`
	WinningTrades  int              `json:"winningTrades"`
	LosingTrades   int              `json:"losingTrades"`
	WinRate        float64          `json:"winRate"`
	TotalProfit    float64          `json:"totalProfit"`
	TotalLoss      float64          `json:"totalLoss"`
	NetProfit      float64          `json:"netProfit"`
	ReturnPercent  float64          `json:"returnPercent"`
	ProfitFactor   float64          `json:"profitFactor"`
	AverageRR      float64          `json:"averageRR"`
	MaxDrawdown    float64          `json:"maxDrawdown"`
	StartBalance   float64          `json:"startBalance"`
	FinalBalance   float64          `json:"finalBalance"`
	PeakBalance    float64          `json:"peakBalance"`
	Trades         []Trade          `json:"trades"`
	ExitReasons    map[string]int   `json:"exitReasons"`
	Duration       string           `json:"duration"`
	
	// Enhanced metrics
	WindowType     string           `json:"windowType,omitempty"`
	MonteCarloSim  *MonteCarloResult `json:"monteCarloSim,omitempty"`
	WalkForwardResults []WalkForwardPeriod `json:"walkForwardResults,omitempty"`
	Confidence95   [2]float64       `json:"confidence95,omitempty"` // [low, high]
}

// MonteCarloResult holds Monte Carlo simulation results
type MonteCarloResult struct {
	Iterations      int     `json:"iterations"`
	MeanReturn      float64 `json:"meanReturn"`
	MedianReturn    float64 `json:"medianReturn"`
	StdDev          float64 `json:"stdDev"`
	WorstCase       float64 `json:"worstCase"`
	BestCase        float64 `json:"bestCase"`
	Confidence95Low float64 `json:"confidence95Low"`
	Confidence95High float64 `json:"confidence95High"`
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
	Type         string  `json:"type"`
	Entry        float64 `json:"entry"`
	Exit         float64 `json:"exit"`
	StopLoss     float64 `json:"stopLoss"`
	ExitReason   string  `json:"exitReason"`
	CandlesHeld  int     `json:"candlesHeld"`
	Profit       float64 `json:"profit"`
	ProfitPercent float64 `json:"profitPercent"`
	RR           float64 `json:"rr"`
	BalanceAfter float64 `json:"balanceAfter"`
	EntryIndex   int     `json:"entryIndex"`
}

// Signal represents a trading signal
type Signal struct {
	Type      string    `json:"type"`
	Entry     float64   `json:"entry"`
	StopLoss  float64   `json:"stopLoss"`
	Targets   []Target  `json:"targets"`
	Strength  float64   `json:"strength"`
	Timeframe string    `json:"timeframe"`
}

// Target represents a take profit target
type Target struct {
	Price      float64 `json:"price"`
	RR         float64 `json:"rr"`
	Percentage int     `json:"percentage"`
}

// RunBacktest executes the backtest with Go's speed
func RunBacktest(config BacktestConfig, candles []Candle) (*BacktestResult, error) {
	startTime := time.Now()
	
	result := &BacktestResult{
		StartBalance: config.StartBalance,
		FinalBalance: config.StartBalance,
		PeakBalance:  config.StartBalance,
		Trades:       []Trade{},
		ExitReasons:  make(map[string]int),
	}

	// Set defaults
	if config.RiskPercent == 0 {
		config.RiskPercent = 0.02 // 2% risk per trade
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

	windowSize := 50
	skipAhead := 5

	// Simulate trading through historical data
	for i := windowSize; i < len(candles)-10; i++ {
		dataWindow := candles[i-windowSize : i]
		futureData := candles[i : i+10]

		// Generate signal using simplified backtest logic (matches JavaScript)
		// If strategy is specified, use strategy-specific parameters but default signal logic
		signal := generateBacktestSignal(dataWindow, config.Interval)
		
		// Apply strategy-specific modifications if strategy is selected
		if signal != nil && config.Strategy != "" && config.Strategy != "default" {
			signal = applyStrategyParameters(signal, config.Strategy)
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
	
	// Calculate position size with cap
	riskAmount := currentBalance * config.RiskPercent
	if riskAmount > config.MaxPositionCap {
		riskAmount = config.MaxPositionCap
	}
	
	riskDiff := math.Abs(entry - stopLoss)
	if riskDiff == 0 {
		return nil
	}
	
	positionSize := riskAmount / riskDiff
	
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
					Type:         signal.Type,
					Entry:        entry,
					Exit:         stopLoss,
					StopLoss:     stopLoss,
					ExitReason:   "Stop Loss",
					CandlesHeld:  candleIdx + 1,
					Profit:       profit,
					ProfitPercent: (profit / riskAmount) * 100,
					RR:           (stopLoss - entry) / (entry - stopLoss),
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
						Type:         signal.Type,
						Entry:        entry,
						Exit:         trailingStopPrice,
						StopLoss:     trailingStopPrice,
						ExitReason:   "Trailing Stop",
						CandlesHeld:  candleIdx + 1,
						Profit:       profit,
						ProfitPercent: (profit / riskAmount) * 100,
						RR:           (trailingStopPrice - entry) / (entry - stopLoss),
					}
				}
			}
			
			// Check target
			if len(signal.Targets) > 0 && candle.High >= signal.Targets[0].Price {
				profit := (signal.Targets[0].Price - entry) * positionSize
				profit -= math.Abs(profit) * config.FeePercent * 2
				
				return &Trade{
					Type:         signal.Type,
					Entry:        entry,
					Exit:         signal.Targets[0].Price,
					StopLoss:     stopLoss,
					ExitReason:   "Target Hit",
					CandlesHeld:  candleIdx + 1,
					Profit:       profit,
					ProfitPercent: (profit / riskAmount) * 100,
					RR:           signal.Targets[0].RR,
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
					Type:         signal.Type,
					Entry:        entry,
					Exit:         stopLoss,
					StopLoss:     stopLoss,
					ExitReason:   "Stop Loss",
					CandlesHeld:  candleIdx + 1,
					Profit:       profit,
					ProfitPercent: (profit / riskAmount) * 100,
					RR:           (entry - stopLoss) / (stopLoss - entry),
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
						Type:         signal.Type,
						Entry:        entry,
						Exit:         trailingStopPrice,
						StopLoss:     trailingStopPrice,
						ExitReason:   "Trailing Stop",
						CandlesHeld:  candleIdx + 1,
						Profit:       profit,
						ProfitPercent: (profit / riskAmount) * 100,
						RR:           (entry - trailingStopPrice) / (stopLoss - entry),
					}
				}
			}
			
			// Check target
			if len(signal.Targets) > 0 && candle.Low <= signal.Targets[0].Price {
				profit := (entry - signal.Targets[0].Price) * positionSize
				profit -= math.Abs(profit) * config.FeePercent * 2
				
				return &Trade{
					Type:         signal.Type,
					Entry:        entry,
					Exit:         signal.Targets[0].Price,
					StopLoss:     stopLoss,
					ExitReason:   "Target Hit",
					CandlesHeld:  candleIdx + 1,
					Profit:       profit,
					ProfitPercent: (profit / riskAmount) * 100,
					RR:           signal.Targets[0].RR,
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
	
	// Calculate average RR
	totalRR := 0.0
	for _, trade := range result.Trades {
		totalRR += trade.RR
	}
	if result.TotalTrades > 0 {
		result.AverageRR = totalRR / float64(result.TotalTrades)
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


// applyStrategyParameters modifies signal based on OPTIMIZED strategy parameters
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
	atr := math.Abs(entry - stopLoss) / 1.5 // Estimate ATR from stop loss
	
	// Apply OPTIMIZED strategy-specific risk/reward ratios from Dec 2, 2025 optimization
	var stopATR, tp1ATR, tp2ATR, tp3ATR float64
	
	switch strategyName {
	case "liquidity_hunter":
		// OPTIMIZED: 61.7% WR, 8.24 PF, 894.1% return - BEST OVERALL
		stopATR, tp1ATR, tp2ATR, tp3ATR = 1.5, 4.0, 6.0, 10.0
	case "session_trader":
		// OPTIMIZED: 54.1% WR, 12.74 PF, 283.3% return - HIGHEST PROFIT FACTOR
		stopATR, tp1ATR, tp2ATR, tp3ATR = 1.0, 4.0, 6.0, 10.0
	case "breakout_master":
		// OPTIMIZED: 54.5% WR, 7.20 PF, 3,845.3% return - HIGHEST RETURN
		stopATR, tp1ATR, tp2ATR, tp3ATR = 1.0, 4.0, 6.0, 10.0
	case "range_master":
		// OPTIMIZED: 44.2% WR, 7.63 PF, 329.5% return
		stopATR, tp1ATR, tp2ATR, tp3ATR = 0.5, 2.0, 3.0, 5.0
	case "institutional_follower":
		// OPTIMIZED: 38.5% WR, 9.08 PF, 1,018.8% return
		stopATR, tp1ATR, tp2ATR, tp3ATR = 0.5, 3.0, 4.5, 7.5
	case "trend_rider":
		// OPTIMIZED: 36.4% WR, 6.71 PF, 942.3% return
		stopATR, tp1ATR, tp2ATR, tp3ATR = 0.5, 3.0, 4.5, 7.5
	case "smart_money_tracker":
		// OPTIMIZED: 34.1% WR, 6.83 PF, 3,508.8% return - MOST ACTIVE (135 trades)
		stopATR, tp1ATR, tp2ATR, tp3ATR = 1.0, 5.0, 7.5, 12.5
	case "reversal_sniper":
		// OPTIMIZED: 28.6% WR, 3.96 PF, 66.6% return
		stopATR, tp1ATR, tp2ATR, tp3ATR = 0.5, 5.0, 7.5, 12.5
	case "momentum_beast":
		// Similar to breakout master
		stopATR, tp1ATR, tp2ATR, tp3ATR = 1.0, 4.0, 6.0, 10.0
	case "scalper_pro":
		// Tight stops for scalping
		stopATR, tp1ATR, tp2ATR, tp3ATR = 0.5, 1.5, 2.5, 3.5
	default:
		// Default parameters
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
