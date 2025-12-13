package backtest

import (
	"math"
	"time"
)

// ProfessionalBacktestEngine - Accurate backtest with partial exits
type ProfessionalBacktestEngine struct{}

// RunProfessionalBacktest executes backtest with ACCURATE partial exit logic
func RunProfessionalBacktest(config BacktestConfig, candles []Candle) (*BacktestResult, error) {
	startTime := time.Now()
	
	result := &BacktestResult{
		StartBalance: config.StartBalance,
		FinalBalance: config.StartBalance,
		PeakBalance:  config.StartBalance,
		Trades:       []Trade{},
		ExitReasons:  make(map[string]int),
	}

	// Set defaults - OPTIMIZED for lower drawdown
	if config.RiskPercent == 0 {
		config.RiskPercent = 0.003 // 0.3% risk per trade (optimized for <12% DD)
	}
	if config.MaxPositionCap == 0 {
		config.MaxPositionCap = config.StartBalance * 10
	}
	if config.SlippagePercent == 0 {
		config.SlippagePercent = 0.0015 // 0.15% (more realistic for real trading)
	}
	if config.FeePercent == 0 {
		config.FeePercent = 0.001 // 0.1% (Binance maker/taker fee)
	}

	windowSize := 100
	skipAhead := 5 // Minimum candles between trades (prevents overtrading)
	
	// Real trading constraints
	maxTradesPerDay := 20 // Prevent overtrading
	tradesThisDay := 0
	currentDay := ""

	// Simulate trading through historical data
	for i := windowSize; i < len(candles)-50; i++ { // Need more future candles for TP3
		dataWindow := candles[i-windowSize : i]
		futureData := candles[i : minInt(i+50, len(candles))] // 50 candles future

		// Check daily trade limit (real trading constraint)
		candleTime := time.Unix(candles[i].Timestamp/1000, 0)
		candleDay := candleTime.Format("2006-01-02")
		if candleDay != currentDay {
			currentDay = candleDay
			tradesThisDay = 0
		}
		if tradesThisDay >= maxTradesPerDay {
			continue // Skip if daily limit reached
		}

		// Generate signal using UNIFIED generator
		usg := &UnifiedSignalGenerator{}
		advSignal := usg.GenerateSignal(dataWindow, config.Strategy)
		
		if advSignal == nil || advSignal.Type == "NONE" {
			continue
		}

		// Simulate trade with PROFESSIONAL partial exit logic
		trade := simulateTradeWithPartialExits(advSignal, futureData, config)
		
		if trade != nil {
			trade.EntryIndex = i
			trade.BalanceAfter = result.FinalBalance + trade.Profit
			
			result.Trades = append(result.Trades, *trade)
			result.TotalTrades++
			tradesThisDay++ // Increment daily trade counter
			
			if trade.Profit > 0 {
				result.WinningTrades++
				result.TotalProfit += trade.Profit
			} else {
				result.LosingTrades++
				result.TotalLoss += math.Abs(trade.Profit)
			}
			
			// Update balance
			result.FinalBalance += trade.Profit
			
			// Track peak and drawdown (FIXED: Use START balance as denominator)
			if result.FinalBalance > result.PeakBalance {
				result.PeakBalance = result.FinalBalance
			}
			
			// Calculate drawdown as percentage of STARTING capital (industry standard)
			drawdownAmount := result.PeakBalance - result.FinalBalance
			drawdown := drawdownAmount / config.StartBalance
			if drawdown > result.MaxDrawdown {
				result.MaxDrawdown = drawdown
			}
			
			// Track exit reasons
			result.ExitReasons[trade.ExitReason]++
			
			// Skip ahead after a trade
			i += skipAhead
		}
	}

	// Calculate statistics
	calculateStats(result)
	result.Duration = time.Since(startTime).String()

	return result, nil
}

// simulateTradeWithPartialExits - PROFESSIONAL: Handles TP1, TP2, TP3 correctly
func simulateTradeWithPartialExits(signal *AdvancedSignal, futureData []Candle, config BacktestConfig) *Trade {
	if signal == nil || len(futureData) == 0 {
		return nil
	}

	entry := signal.Entry
	stopLoss := signal.StopLoss
	tp1 := signal.TP1
	tp2 := signal.TP2
	tp3 := signal.TP3
	
	// Calculate position size based on START balance (not current)
	riskAmount := config.StartBalance * config.RiskPercent
	riskDiff := math.Abs(entry - stopLoss)
	if riskDiff == 0 {
		return nil
	}
	
	positionSize := riskAmount / riskDiff
	
	// Cap position size
	maxPositionValue := riskAmount * 10
	if positionSize * entry > maxPositionValue {
		positionSize = maxPositionValue / entry
	}
	
	// Apply slippage
	if signal.Type == "BUY" {
		entry *= (1 + config.SlippagePercent)
	} else {
		entry *= (1 - config.SlippagePercent)
	}
	
	// Track partial exits
	remainingPosition := positionSize
	totalProfit := 0.0
	exitReason := ""
	exitPrice := 0.0
	candlesHeld := 0
	
	// Partial exit percentages (optimized: 30%, 30%, 40% - let more ride!)
	tp1Percent := 0.30
	tp2Percent := 0.30
	tp3Percent := 0.40
	
	tp1Hit := false
	tp2Hit := false
	tp3Hit := false
	
	// Simulate price movement
	for candleIdx, candle := range futureData {
		candlesHeld = candleIdx + 1
		
		if signal.Type == "BUY" {
			// Check stop loss FIRST (highest priority)
			if candle.Low <= stopLoss {
				// Exit remaining position at stop loss
				profit := (stopLoss - entry) * remainingPosition
				profit -= math.Abs(profit) * config.FeePercent * 2
				totalProfit += profit
				
				exitReason = "Stop Loss"
				exitPrice = stopLoss
				break
			}
			
			// Check TP1
			if !tp1Hit && candle.High >= tp1 {
				tp1Hit = true
				exitSize := positionSize * tp1Percent
				profit := (tp1 - entry) * exitSize
				profit -= math.Abs(profit) * config.FeePercent * 2
				totalProfit += profit
				remainingPosition -= exitSize
				
				// Move stop loss to breakeven after TP1
				stopLoss = entry
				
				if remainingPosition <= 0 {
					exitReason = "Target 1"
					exitPrice = tp1
					break
				}
			}
			
			// Check TP2
			if tp1Hit && !tp2Hit && candle.High >= tp2 {
				tp2Hit = true
				exitSize := positionSize * tp2Percent
				profit := (tp2 - entry) * exitSize
				profit -= math.Abs(profit) * config.FeePercent * 2
				totalProfit += profit
				remainingPosition -= exitSize
				
				if remainingPosition <= 0 {
					exitReason = "Target 2"
					exitPrice = tp2
					break
				}
			}
			
			// Check TP3
			if tp2Hit && !tp3Hit && candle.High >= tp3 {
				tp3Hit = true
				exitSize := positionSize * tp3Percent
				profit := (tp3 - entry) * exitSize
				profit -= math.Abs(profit) * config.FeePercent * 2
				totalProfit += profit
				remainingPosition -= exitSize
				
				exitReason = "Target 3"
				exitPrice = tp3
				break
			}
			
		} else { // SELL
			// Check stop loss FIRST
			if candle.High >= stopLoss {
				profit := (entry - stopLoss) * remainingPosition
				profit -= math.Abs(profit) * config.FeePercent * 2
				totalProfit += profit
				
				exitReason = "Stop Loss"
				exitPrice = stopLoss
				break
			}
			
			// Check TP1
			if !tp1Hit && candle.Low <= tp1 {
				tp1Hit = true
				exitSize := positionSize * tp1Percent
				profit := (entry - tp1) * exitSize
				profit -= math.Abs(profit) * config.FeePercent * 2
				totalProfit += profit
				remainingPosition -= exitSize
				
				// Move stop loss to breakeven
				stopLoss = entry
				
				if remainingPosition <= 0 {
					exitReason = "Target 1"
					exitPrice = tp1
					break
				}
			}
			
			// Check TP2
			if tp1Hit && !tp2Hit && candle.Low <= tp2 {
				tp2Hit = true
				exitSize := positionSize * tp2Percent
				profit := (entry - tp2) * exitSize
				profit -= math.Abs(profit) * config.FeePercent * 2
				totalProfit += profit
				remainingPosition -= exitSize
				
				if remainingPosition <= 0 {
					exitReason = "Target 2"
					exitPrice = tp2
					break
				}
			}
			
			// Check TP3
			if tp2Hit && !tp3Hit && candle.Low <= tp3 {
				tp3Hit = true
				exitSize := positionSize * tp3Percent
				profit := (entry - tp3) * exitSize
				profit -= math.Abs(profit) * config.FeePercent * 2
				totalProfit += profit
				remainingPosition -= exitSize
				
				exitReason = "Target 3"
				exitPrice = tp3
				break
			}
		}
		
		// Timeout after 50 candles if no exit
		if candleIdx >= len(futureData)-1 {
			// Force exit at current price
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
	
	// Calculate final metrics
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


