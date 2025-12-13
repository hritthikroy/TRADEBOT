package backtest

import (
	"fmt"
	"sort"
	"time"
)

// ORBBacktestEngine runs backtests for the ORB Academic Strategy
type ORBBacktestEngine struct {
	Strategy       *ORBAcademicStrategy
	InitialCapital float64
	CurrentCapital float64
	ActiveSignals  map[string]*ORBSignal
	ClosedSignals  []*ORBSignal
	EquityCurve    []EquityPoint
	DailyReturns   []float64
}

// NewORBBacktestEngine creates a new backtest engine
func NewORBBacktestEngine(strategy *ORBAcademicStrategy, initialCapital float64) *ORBBacktestEngine {
	return &ORBBacktestEngine{
		Strategy:       strategy,
		InitialCapital: initialCapital,
		CurrentCapital: initialCapital,
		ActiveSignals:  make(map[string]*ORBSignal),
		ClosedSignals:  make([]*ORBSignal, 0),
		EquityCurve:    make([]EquityPoint, 0),
		DailyReturns:   make([]float64, 0),
	}
}

// DailyStockData represents all data needed for a stock on a given day
type DailyStockData struct {
	Symbol            string
	Date              time.Time
	OpenPrice         float64
	HighPrice         float64
	LowPrice          float64
	ClosePrice        float64
	Volume            float64
	IntradayCandles   []Candle // 1-minute candles for the day
	ATR14             float64
	AvgVolume14       float64
	HistoricalORVols  []float64 // Last 14 days of OR volumes
}

// ProcessTradingDay processes a single trading day for the ORB strategy
func (engine *ORBBacktestEngine) ProcessTradingDay(
	date time.Time,
	stocksData []DailyStockData,
) error {
	
	// Step 1: Calculate opening range for all stocks
	candidates := engine.identifyStockCandidates(date, stocksData)
	
	// Step 2: Filter by relative volume and select top N
	topCandidates := engine.selectTopCandidates(candidates)
	
	// Step 3: Generate signals for top candidates
	signals := engine.generateSignals(date, topCandidates)
	
	// Step 4: Simulate intraday trading
	engine.simulateIntradayTrading(date, stocksData, signals)
	
	// Step 5: Close all positions at EOD
	engine.closeAllPositionsEOD(date, stocksData)
	
	// Step 6: Update equity curve
	engine.updateEquityCurve(date)
	
	return nil
}

// identifyStockCandidates identifies stocks that meet basic ORB criteria
func (engine *ORBBacktestEngine) identifyStockCandidates(
	date time.Time,
	stocksData []DailyStockData,
) []StockCandidate {
	
	candidates := make([]StockCandidate, 0)
	
	for _, stock := range stocksData {
		// Check basic filters
		if !engine.Strategy.MeetsBasicFilters(stock.OpenPrice, stock.ATR14, stock.AvgVolume14) {
			continue
		}
		
		// Calculate opening range
		orData := engine.calculateOpeningRange(stock.IntradayCandles, engine.Strategy.TimeFrame)
		if orData == nil {
			continue
		}
		
		// Calculate relative volume
		relVol := CalculateRelativeVolume(orData.Volume, stock.HistoricalORVols)
		if relVol < engine.Strategy.MinRelativeVol {
			continue
		}
		
		// Determine direction
		direction := DetermineDirection(orData.Open, orData.Close)
		if direction == "NONE" {
			continue
		}
		
		candidate := StockCandidate{
			Symbol:         stock.Symbol,
			OpenPrice:      stock.OpenPrice,
			ORHigh:         orData.High,
			ORLow:          orData.Low,
			ORClose:        orData.Close,
			ORVolume:       orData.Volume,
			ATR14:          stock.ATR14,
			AvgVolume14:    stock.AvgVolume14,
			RelativeVolume: relVol,
			Direction:      direction,
		}
		
		candidates = append(candidates, candidate)
	}
	
	return candidates
}

// OpeningRangeData holds the OHLCV for the opening range period
type OpeningRangeData struct {
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

// calculateOpeningRange calculates the opening range from intraday candles
func (engine *ORBBacktestEngine) calculateOpeningRange(
	candles []Candle,
	timeFrameMinutes int,
) *OpeningRangeData {
	
	if len(candles) < timeFrameMinutes {
		return nil
	}
	
	// Get first N minutes of candles (assuming 1-minute candles)
	orCandles := candles[:timeFrameMinutes]
	
	orData := &OpeningRangeData{
		Open:  orCandles[0].Open,
		Close: orCandles[len(orCandles)-1].Close,
		High:  orCandles[0].High,
		Low:   orCandles[0].Low,
	}
	
	// Find high, low, and sum volume
	for _, candle := range orCandles {
		if candle.High > orData.High {
			orData.High = candle.High
		}
		if candle.Low < orData.Low {
			orData.Low = candle.Low
		}
		orData.Volume += candle.Volume
	}
	
	return orData
}

// selectTopCandidates selects the top N candidates by relative volume
func (engine *ORBBacktestEngine) selectTopCandidates(candidates []StockCandidate) []StockCandidate {
	if len(candidates) <= engine.Strategy.TopNStocks {
		return candidates
	}
	
	// Sort by relative volume (descending)
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].RelativeVolume > candidates[j].RelativeVolume
	})
	
	return candidates[:engine.Strategy.TopNStocks]
}

// generateSignals generates trading signals for the selected candidates
func (engine *ORBBacktestEngine) generateSignals(
	date time.Time,
	candidates []StockCandidate,
) []*ORBSignal {
	
	signals := make([]*ORBSignal, 0)
	
	for _, candidate := range candidates {
		signal := engine.Strategy.GenerateORBSignal(candidate, engine.CurrentCapital, date)
		if signal != nil {
			signals = append(signals, signal)
		}
	}
	
	return signals
}

// simulateIntradayTrading simulates the intraday price action and trade execution
func (engine *ORBBacktestEngine) simulateIntradayTrading(
	date time.Time,
	stocksData []DailyStockData,
	signals []*ORBSignal,
) {
	
	// Create a map for quick stock data lookup
	stockMap := make(map[string]DailyStockData)
	for _, stock := range stocksData {
		stockMap[stock.Symbol] = stock
	}
	
	// Process each signal
	for _, signal := range signals {
		stock, exists := stockMap[signal.Symbol]
		if !exists {
			continue
		}
		
		// Skip candles in the opening range
		startIdx := engine.Strategy.TimeFrame
		if startIdx >= len(stock.IntradayCandles) {
			continue
		}
		
		// Simulate tick-by-tick through remaining candles
		for i := startIdx; i < len(stock.IntradayCandles); i++ {
			candle := stock.IntradayCandles[i]
			
			// Check if entry is triggered
			if !signal.Triggered {
				if signal.CheckEntryTrigger(candle.High) || signal.CheckEntryTrigger(candle.Low) {
					signal.Triggered = true
					engine.ActiveSignals[signal.Symbol] = signal
				}
				continue
			}
			
			// Check if stop loss is hit
			if signal.CheckStopLoss(candle.Low) || signal.CheckStopLoss(candle.High) {
				signal.ExitPrice = signal.StopLoss
				signal.ExitTime = time.Unix(candle.Timestamp, 0)
				engine.Strategy.CalculatePnL(signal)
				engine.CurrentCapital += signal.PnL
				engine.ClosedSignals = append(engine.ClosedSignals, signal)
				delete(engine.ActiveSignals, signal.Symbol)
				break
			}
		}
	}
}

// closeAllPositionsEOD closes all active positions at end of day
func (engine *ORBBacktestEngine) closeAllPositionsEOD(
	date time.Time,
	stocksData []DailyStockData,
) {
	
	stockMap := make(map[string]DailyStockData)
	for _, stock := range stocksData {
		stockMap[stock.Symbol] = stock
	}
	
	for symbol, signal := range engine.ActiveSignals {
		stock, exists := stockMap[symbol]
		if !exists {
			continue
		}
		
		// Close at EOD price
		signal.ExitPrice = stock.ClosePrice
		signal.ExitTime = date.Add(16 * time.Hour) // 4:00 PM ET
		engine.Strategy.CalculatePnL(signal)
		engine.CurrentCapital += signal.PnL
		engine.ClosedSignals = append(engine.ClosedSignals, signal)
	}
	
	// Clear active signals
	engine.ActiveSignals = make(map[string]*ORBSignal)
}

// updateEquityCurve updates the equity curve with current day's performance
func (engine *ORBBacktestEngine) updateEquityCurve(date time.Time) {
	var dailyReturn float64
	if len(engine.EquityCurve) > 0 {
		prevEquity := engine.EquityCurve[len(engine.EquityCurve)-1].Equity
		dailyReturn = (engine.CurrentCapital - prevEquity) / prevEquity
	}
	
	point := EquityPoint{
		Date:   date,
		Equity: engine.CurrentCapital,
		Return: dailyReturn,
	}
	
	engine.EquityCurve = append(engine.EquityCurve, point)
	engine.DailyReturns = append(engine.DailyReturns, dailyReturn)
}

// GetResults returns the final backtest results
func (engine *ORBBacktestEngine) GetResults(startDate, endDate time.Time) *ORBBacktestResult {
	result := &ORBBacktestResult{
		Strategy:       fmt.Sprintf("%d-minute ORB + Relative Volume", engine.Strategy.TimeFrame),
		TimeFrame:      engine.Strategy.TimeFrame,
		StartDate:      startDate,
		EndDate:        endDate,
		InitialCapital: engine.InitialCapital,
		FinalCapital:   engine.CurrentCapital,
		Signals:        engine.ClosedSignals,
		EquityCurve:    engine.EquityCurve,
	}
	
	CalculateBacktestMetrics(result)
	
	return result
}

// RunORBBacktest is a convenience function to run a complete backtest
func RunORBBacktest(
	timeFrame int,
	startDate, endDate time.Time,
	stocksDataByDay map[time.Time][]DailyStockData,
	initialCapital float64,
) (*ORBBacktestResult, error) {
	
	strategy := NewORBAcademicStrategy(timeFrame)
	engine := NewORBBacktestEngine(strategy, initialCapital)
	
	// Get sorted list of trading days
	tradingDays := make([]time.Time, 0, len(stocksDataByDay))
	for day := range stocksDataByDay {
		if day.After(startDate) && day.Before(endDate) || day.Equal(startDate) || day.Equal(endDate) {
			tradingDays = append(tradingDays, day)
		}
	}
	sort.Slice(tradingDays, func(i, j int) bool {
		return tradingDays[i].Before(tradingDays[j])
	})
	
	// Process each trading day
	for _, day := range tradingDays {
		stocksData := stocksDataByDay[day]
		err := engine.ProcessTradingDay(day, stocksData)
		if err != nil {
			return nil, fmt.Errorf("error processing day %s: %v", day.Format("2006-01-02"), err)
		}
	}
	
	return engine.GetResults(startDate, endDate), nil
}
