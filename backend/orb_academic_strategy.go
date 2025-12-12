package main

import (
	"math"
	"time"
)

// ORBAcademicStrategy implements the Opening Range Breakout strategy
// based on the academic paper "A Profitable Day Trading Strategy For The U.S. Equity Market"
// by Zarattini, Barbon, and Aziz (2024)
type ORBAcademicStrategy struct {
	TimeFrame        int     // Opening range in minutes (5, 15, 30, or 60)
	MinPrice         float64 // Minimum stock price ($5)
	MinAvgVolume     float64 // Minimum 14-day average volume (1,000,000)
	MinATR           float64 // Minimum 14-day ATR ($0.50)
	MinRelativeVol   float64 // Minimum relative volume (100% = 1.0)
	TopNStocks       int     // Trade only top N stocks by relative volume (20)
	StopLossATRPct   float64 // Stop loss as % of ATR (10% = 0.10)
	RiskPerTrade     float64 // Risk per trade as % of capital (1% = 0.01)
	MaxLeverage      float64 // Maximum leverage allowed (4.0)
	CommissionPerShr float64 // Commission per share ($0.0035)
}

// NewORBAcademicStrategy creates a new ORB strategy with default parameters from the paper
func NewORBAcademicStrategy(timeFrame int) *ORBAcademicStrategy {
	return &ORBAcademicStrategy{
		TimeFrame:        timeFrame,
		MinPrice:         5.0,
		MinAvgVolume:     1000000,
		MinATR:           0.50,
		MinRelativeVol:   1.0, // 100%
		TopNStocks:       20,
		StopLossATRPct:   0.10,
		RiskPerTrade:     0.01,
		MaxLeverage:      4.0,
		CommissionPerShr: 0.0035,
	}
}

// StockCandidate represents a stock that meets the ORB criteria
type StockCandidate struct {
	Symbol         string
	OpenPrice      float64
	ORHigh         float64 // Opening Range High
	ORLow          float64 // Opening Range Low
	ORClose        float64 // Opening Range Close
	ORVolume       float64
	ATR14          float64
	AvgVolume14    float64
	RelativeVolume float64
	Direction      string // "LONG" or "SHORT"
	EntryPrice     float64
	StopLoss       float64
	PositionSize   int
}

// ORBSignal represents a trading signal from the ORB strategy
type ORBSignal struct {
	Symbol        string
	Direction     string
	EntryPrice    float64
	StopLoss      float64
	ProfitTarget  float64 // EOD close
	PositionSize  int
	RiskAmount    float64
	RelativeVol   float64
	ATR           float64
	Timestamp     time.Time
	TimeFrame     int
	ORHigh        float64
	ORLow         float64
	Triggered     bool
	ExitPrice     float64
	ExitTime      time.Time
	PnL           float64
	PnLInR        float64 // Profit/Loss in R multiples
}

// CalculateRelativeVolume calculates the relative volume for the opening range
// RelativeVolume = Current OR Volume / Average OR Volume (14 days)
func CalculateRelativeVolume(currentORVolume float64, historicalORVolumes []float64) float64 {
	if len(historicalORVolumes) == 0 {
		return 0
	}

	var sum float64
	for _, vol := range historicalORVolumes {
		sum += vol
	}
	avgVolume := sum / float64(len(historicalORVolumes))

	if avgVolume == 0 {
		return 0
	}

	return currentORVolume / avgVolume
}

// CalculateATR14 calculates the 14-day Average True Range
func CalculateATR14(candles []Candle) float64 {
	if len(candles) < 14 {
		return 0
	}

	var trSum float64
	for i := len(candles) - 14; i < len(candles); i++ {
		tr := CalculateTrueRange(candles[i], candles[i-1])
		trSum += tr
	}

	return trSum / 14.0
}

// CalculateTrueRange calculates the True Range for a candle
func CalculateTrueRange(current, previous Candle) float64 {
	highLow := current.High - current.Low
	highClose := math.Abs(current.High - previous.Close)
	lowClose := math.Abs(current.Low - previous.Close)

	return math.Max(highLow, math.Max(highClose, lowClose))
}

// CalculateAvgVolume14 calculates the 14-day average volume
func CalculateAvgVolume14(volumes []float64) float64 {
	if len(volumes) < 14 {
		return 0
	}

	var sum float64
	for i := len(volumes) - 14; i < len(volumes); i++ {
		sum += volumes[i]
	}

	return sum / 14.0
}

// MeetsBasicFilters checks if a stock meets the basic ORB criteria
func (s *ORBAcademicStrategy) MeetsBasicFilters(openPrice, atr14, avgVolume14 float64) bool {
	return openPrice >= s.MinPrice &&
		avgVolume14 >= s.MinAvgVolume &&
		atr14 >= s.MinATR
}

// DetermineDirection determines trade direction based on opening range
// Bullish OR (close > open) = LONG only
// Bearish OR (close < open) = SHORT only
// Doji (close = open) = NO TRADE
func DetermineDirection(orOpen, orClose float64) string {
	if orClose > orOpen {
		return "LONG"
	} else if orClose < orOpen {
		return "SHORT"
	}
	return "NONE"
}

// CalculatePositionSize calculates position size based on risk management
// Position sized so that if stop loss is hit, loss = RiskPerTrade% of capital
// Also respects MaxLeverage constraint
func (s *ORBAcademicStrategy) CalculatePositionSize(
	capital, entryPrice, stopLoss float64,
) int {
	riskPerShare := math.Abs(entryPrice - stopLoss)
	if riskPerShare == 0 {
		return 0
	}

	// Calculate shares based on risk
	riskAmount := capital * s.RiskPerTrade
	sharesFromRisk := int(riskAmount / riskPerShare)

	// Calculate max shares based on leverage
	maxPositionValue := capital * s.MaxLeverage
	maxSharesFromLeverage := int(maxPositionValue / entryPrice)

	// Take the minimum to respect both constraints
	shares := sharesFromRisk
	if maxSharesFromLeverage < shares {
		shares = maxSharesFromLeverage
	}

	return shares
}

// GenerateORBSignal generates a trading signal for a stock candidate
func (s *ORBAcademicStrategy) GenerateORBSignal(
	candidate StockCandidate,
	capital float64,
	timestamp time.Time,
) *ORBSignal {
	if candidate.Direction == "NONE" {
		return nil
	}

	// Set entry price (stop order at OR high/low)
	var entryPrice, stopLoss float64
	if candidate.Direction == "LONG" {
		entryPrice = candidate.ORHigh
		stopLoss = entryPrice - (candidate.ATR14 * s.StopLossATRPct)
	} else { // SHORT
		entryPrice = candidate.ORLow
		stopLoss = entryPrice + (candidate.ATR14 * s.StopLossATRPct)
	}

	// Calculate position size
	positionSize := s.CalculatePositionSize(capital, entryPrice, stopLoss)
	if positionSize == 0 {
		return nil
	}

	riskAmount := math.Abs(entryPrice-stopLoss) * float64(positionSize)

	return &ORBSignal{
		Symbol:       candidate.Symbol,
		Direction:    candidate.Direction,
		EntryPrice:   entryPrice,
		StopLoss:     stopLoss,
		ProfitTarget: 0, // Set at EOD
		PositionSize: positionSize,
		RiskAmount:   riskAmount,
		RelativeVol:  candidate.RelativeVolume,
		ATR:          candidate.ATR14,
		Timestamp:    timestamp,
		TimeFrame:    s.TimeFrame,
		ORHigh:       candidate.ORHigh,
		ORLow:        candidate.ORLow,
		Triggered:    false,
	}
}

// CheckEntryTrigger checks if the entry condition is met
func (signal *ORBSignal) CheckEntryTrigger(currentPrice float64) bool {
	if signal.Triggered {
		return false
	}

	if signal.Direction == "LONG" && currentPrice >= signal.EntryPrice {
		signal.Triggered = true
		return true
	} else if signal.Direction == "SHORT" && currentPrice <= signal.EntryPrice {
		signal.Triggered = true
		return true
	}

	return false
}

// CheckStopLoss checks if stop loss is hit
func (signal *ORBSignal) CheckStopLoss(currentPrice float64) bool {
	if !signal.Triggered {
		return false
	}

	if signal.Direction == "LONG" && currentPrice <= signal.StopLoss {
		return true
	} else if signal.Direction == "SHORT" && currentPrice >= signal.StopLoss {
		return true
	}

	return false
}

// CalculatePnL calculates the profit/loss for a closed position
func (s *ORBAcademicStrategy) CalculatePnL(signal *ORBSignal) {
	if !signal.Triggered {
		signal.PnL = 0
		signal.PnLInR = 0
		return
	}

	var pnlPerShare float64
	if signal.Direction == "LONG" {
		pnlPerShare = signal.ExitPrice - signal.EntryPrice
	} else { // SHORT
		pnlPerShare = signal.EntryPrice - signal.ExitPrice
	}

	// Subtract commission
	totalCommission := s.CommissionPerShr * float64(signal.PositionSize) * 2 // Entry + Exit
	grossPnL := pnlPerShare * float64(signal.PositionSize)
	signal.PnL = grossPnL - totalCommission

	// Calculate PnL in R (risk units)
	riskPerShare := math.Abs(signal.EntryPrice - signal.StopLoss)
	if riskPerShare > 0 {
		signal.PnLInR = pnlPerShare / riskPerShare
	}
}

// ORBBacktestResult stores the results of an ORB backtest
type ORBBacktestResult struct {
	Strategy         string
	TimeFrame        int
	StartDate        time.Time
	EndDate          time.Time
	InitialCapital   float64
	FinalCapital     float64
	TotalReturn      float64
	AnnualizedReturn float64
	Volatility       float64
	SharpeRatio      float64
	MaxDrawdown      float64
	WinRate          float64
	TotalTrades      int
	WinningTrades    int
	LosingTrades     int
	AvgWin           float64
	AvgLoss          float64
	AvgPnLInR        float64
	LargestWin       float64
	LargestLoss      float64
	Alpha            float64
	Beta             float64
	Signals          []*ORBSignal
	EquityCurve      []EquityPoint
}

// EquityPoint represents a point in the equity curve
type EquityPoint struct {
	Date   time.Time
	Equity float64
	Return float64
}

// CalculateBacktestMetrics calculates comprehensive backtest statistics
func CalculateBacktestMetrics(result *ORBBacktestResult) {
	if len(result.Signals) == 0 {
		return
	}

	// Calculate basic metrics
	result.TotalReturn = (result.FinalCapital - result.InitialCapital) / result.InitialCapital

	// Calculate win/loss statistics
	var totalWin, totalLoss, totalPnLInR float64
	result.LargestWin = math.Inf(-1)
	result.LargestLoss = math.Inf(1)

	for _, signal := range result.Signals {
		if !signal.Triggered {
			continue
		}

		result.TotalTrades++
		totalPnLInR += signal.PnLInR

		if signal.PnL > 0 {
			result.WinningTrades++
			totalWin += signal.PnL
			if signal.PnL > result.LargestWin {
				result.LargestWin = signal.PnL
			}
		} else {
			result.LosingTrades++
			totalLoss += signal.PnL
			if signal.PnL < result.LargestLoss {
				result.LargestLoss = signal.PnL
			}
		}
	}

	if result.TotalTrades > 0 {
		result.WinRate = float64(result.WinningTrades) / float64(result.TotalTrades)
		result.AvgPnLInR = totalPnLInR / float64(result.TotalTrades)
	}

	if result.WinningTrades > 0 {
		result.AvgWin = totalWin / float64(result.WinningTrades)
	}

	if result.LosingTrades > 0 {
		result.AvgLoss = totalLoss / float64(result.LosingTrades)
	}

	// Calculate annualized return
	years := result.EndDate.Sub(result.StartDate).Hours() / 24 / 365.25
	if years > 0 {
		result.AnnualizedReturn = math.Pow(1+result.TotalReturn, 1/years) - 1
	}

	// Calculate volatility and Sharpe ratio from equity curve
	if len(result.EquityCurve) > 1 {
		var returns []float64
		for i := 1; i < len(result.EquityCurve); i++ {
			returns = append(returns, result.EquityCurve[i].Return)
		}

		// Calculate mean first
		var sum float64
		for _, r := range returns {
			sum += r
		}
		avgReturn := sum / float64(len(returns))
		
		// Calculate standard deviation
		var sumSquares float64
		for _, r := range returns {
			diff := r - avgReturn
			sumSquares += diff * diff
		}
		stdDev := math.Sqrt(sumSquares / float64(len(returns)))
		
		result.Volatility = stdDev * math.Sqrt(252) // Annualized
		avgReturn = avgReturn * 252                  // Annualized

		if result.Volatility > 0 {
			result.SharpeRatio = avgReturn / result.Volatility
		}
	}

	// Calculate maximum drawdown
	if len(result.EquityCurve) > 0 {
		maxDrawdown := 0.0
		peak := result.EquityCurve[0].Equity

		for _, point := range result.EquityCurve {
			if point.Equity > peak {
				peak = point.Equity
			}
			drawdown := (peak - point.Equity) / peak
			if drawdown > maxDrawdown {
				maxDrawdown = drawdown
			}
		}
		result.MaxDrawdown = maxDrawdown
	}
}
