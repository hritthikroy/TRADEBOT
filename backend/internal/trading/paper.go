package trading

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// PaperTrade represents a single paper trade
type PaperTrade struct {
	ID           int       `json:"id"`
	Signal       string    `json:"signal"`
	Entry        float64   `json:"entry"`
	StopLoss     float64   `json:"stopLoss"`
	TakeProfit   float64   `json:"takeProfit"`
	TP1          float64   `json:"tp1"`
	TP2          float64   `json:"tp2"`
	TP3          float64   `json:"tp3"`
	EntryTime    time.Time `json:"entryTime"`
	ExitTime     *time.Time `json:"exitTime,omitempty"`
	ExitPrice    float64   `json:"exitPrice"`
	ExitReason   string    `json:"exitReason"` // "TP1", "TP2", "TP3", "StopLoss", "Open"
	Profit       float64   `json:"profit"`
	ProfitPercent float64  `json:"profitPercent"`
	Status       string    `json:"status"` // "open", "won", "lost"
	RiskAmount   float64   `json:"riskAmount"`
}

// PaperTradingStats represents overall statistics
type PaperTradingStats struct {
	TotalTrades    int     `json:"totalTrades"`
	OpenTrades     int     `json:"openTrades"`
	ClosedTrades   int     `json:"closedTrades"`
	WinningTrades  int     `json:"winningTrades"`
	LosingTrades   int     `json:"losingTrades"`
	WinRate        float64 `json:"winRate"`
	TotalProfit    float64 `json:"totalProfit"`
	TotalLoss      float64 `json:"totalLoss"`
	NetProfit      float64 `json:"netProfit"`
	ProfitFactor   float64 `json:"profitFactor"`
	StartBalance   float64 `json:"startBalance"`
	CurrentBalance float64 `json:"currentBalance"`
	ReturnPercent  float64 `json:"returnPercent"`
	MaxDrawdown    float64 `json:"maxDrawdown"`
	AverageWin     float64 `json:"averageWin"`
	AverageLoss    float64 `json:"averageLoss"`
}

// PaperTradingManager manages paper trades
type PaperTradingManager struct {
	trades        []PaperTrade
	startBalance  float64
	currentBalance float64
	riskPercent   float64
	mu            sync.RWMutex
	dataFile      string
}

var paperTradingManager *PaperTradingManager

func init() {
	paperTradingManager = &PaperTradingManager{
		trades:        []PaperTrade{},
		startBalance:  15.0,
		currentBalance: 15.0,
		riskPercent:   0.003, // 0.3%
		dataFile:      "paper_trades.json",
	}
	paperTradingManager.loadTrades()
}

// loadTrades loads trades from file
func (ptm *PaperTradingManager) loadTrades() {
	data, err := os.ReadFile(ptm.dataFile)
	if err != nil {
		return // File doesn't exist yet
	}
	
	var savedData struct {
		Trades        []PaperTrade `json:"trades"`
		StartBalance  float64      `json:"startBalance"`
		CurrentBalance float64     `json:"currentBalance"`
	}
	
	if err := json.Unmarshal(data, &savedData); err == nil {
		ptm.trades = savedData.Trades
		ptm.startBalance = savedData.StartBalance
		ptm.currentBalance = savedData.CurrentBalance
	}
}

// saveTrades saves trades to file
func (ptm *PaperTradingManager) saveTrades() error {
	ptm.mu.RLock()
	defer ptm.mu.RUnlock()
	
	data := struct {
		Trades        []PaperTrade `json:"trades"`
		StartBalance  float64      `json:"startBalance"`
		CurrentBalance float64     `json:"currentBalance"`
	}{
		Trades:        ptm.trades,
		StartBalance:  ptm.startBalance,
		CurrentBalance: ptm.currentBalance,
	}
	
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(ptm.dataFile, jsonData, 0644)
}

// AddTrade adds a new paper trade
func (ptm *PaperTradingManager) AddTrade(signal *AdvancedSignal, currentPrice float64) *PaperTrade {
	ptm.mu.Lock()
	defer ptm.mu.Unlock()
	
	riskAmount := ptm.currentBalance * ptm.riskPercent
	
	trade := PaperTrade{
		ID:         len(ptm.trades) + 1,
		Signal:     signal.Type,
		Entry:      signal.Entry,
		StopLoss:   signal.StopLoss,
		TakeProfit: signal.TP3,
		TP1:        signal.TP1,
		TP2:        signal.TP2,
		TP3:        signal.TP3,
		EntryTime:  time.Now(),
		Status:     "open",
		RiskAmount: riskAmount,
	}
	
	ptm.trades = append(ptm.trades, trade)
	ptm.saveTrades()
	
	return &trade
}

// UpdateOpenTrades checks and updates all open trades
func (ptm *PaperTradingManager) UpdateOpenTrades(currentPrice float64) []PaperTrade {
	ptm.mu.Lock()
	defer ptm.mu.Unlock()
	
	var closedTrades []PaperTrade
	
	for i := range ptm.trades {
		if ptm.trades[i].Status != "open" {
			continue
		}
		
		trade := &ptm.trades[i]
		
		if trade.Signal == "BUY" {
			// Check if TP or SL hit
			if currentPrice >= trade.TP3 {
				trade.ExitPrice = trade.TP3
				trade.ExitReason = "TP3"
				trade.Status = "won"
				trade.Profit = trade.RiskAmount * 4.5 // Approximate
			} else if currentPrice >= trade.TP2 {
				trade.ExitPrice = trade.TP2
				trade.ExitReason = "TP2"
				trade.Status = "won"
				trade.Profit = trade.RiskAmount * 3.0
			} else if currentPrice >= trade.TP1 {
				trade.ExitPrice = trade.TP1
				trade.ExitReason = "TP1"
				trade.Status = "won"
				trade.Profit = trade.RiskAmount * 2.0
			} else if currentPrice <= trade.StopLoss {
				trade.ExitPrice = trade.StopLoss
				trade.ExitReason = "StopLoss"
				trade.Status = "lost"
				trade.Profit = -trade.RiskAmount
			}
		} else if trade.Signal == "SELL" {
			// Check if TP or SL hit
			if currentPrice <= trade.TP3 {
				trade.ExitPrice = trade.TP3
				trade.ExitReason = "TP3"
				trade.Status = "won"
				trade.Profit = trade.RiskAmount * 4.5
			} else if currentPrice <= trade.TP2 {
				trade.ExitPrice = trade.TP2
				trade.ExitReason = "TP2"
				trade.Status = "won"
				trade.Profit = trade.RiskAmount * 3.0
			} else if currentPrice <= trade.TP1 {
				trade.ExitPrice = trade.TP1
				trade.ExitReason = "TP1"
				trade.Status = "won"
				trade.Profit = trade.RiskAmount * 2.0
			} else if currentPrice >= trade.StopLoss {
				trade.ExitPrice = trade.StopLoss
				trade.ExitReason = "StopLoss"
				trade.Status = "lost"
				trade.Profit = -trade.RiskAmount
			}
		}
		
		if trade.Status != "open" {
			now := time.Now()
			trade.ExitTime = &now
			trade.ProfitPercent = (trade.Profit / ptm.currentBalance) * 100
			ptm.currentBalance += trade.Profit
			closedTrades = append(closedTrades, *trade)
		}
	}
	
	if len(closedTrades) > 0 {
		ptm.saveTrades()
	}
	
	return closedTrades
}

// GetStats returns paper trading statistics
func (ptm *PaperTradingManager) GetStats() PaperTradingStats {
	ptm.mu.RLock()
	defer ptm.mu.RUnlock()
	
	stats := PaperTradingStats{
		TotalTrades:    len(ptm.trades),
		StartBalance:   ptm.startBalance,
		CurrentBalance: ptm.currentBalance,
	}
	
	var totalProfit, totalLoss float64
	var winningTrades, losingTrades, openTrades int
	
	for _, trade := range ptm.trades {
		if trade.Status == "open" {
			openTrades++
		} else if trade.Status == "won" {
			winningTrades++
			totalProfit += trade.Profit
		} else if trade.Status == "lost" {
			losingTrades++
			totalLoss += -trade.Profit
		}
	}
	
	stats.OpenTrades = openTrades
	stats.ClosedTrades = winningTrades + losingTrades
	stats.WinningTrades = winningTrades
	stats.LosingTrades = losingTrades
	stats.TotalProfit = totalProfit
	stats.TotalLoss = totalLoss
	stats.NetProfit = totalProfit - totalLoss
	
	if stats.ClosedTrades > 0 {
		stats.WinRate = (float64(winningTrades) / float64(stats.ClosedTrades)) * 100
	}
	
	if totalLoss > 0 {
		stats.ProfitFactor = totalProfit / totalLoss
	}
	
	if stats.StartBalance > 0 {
		stats.ReturnPercent = ((stats.CurrentBalance - stats.StartBalance) / stats.StartBalance) * 100
	}
	
	if winningTrades > 0 {
		stats.AverageWin = totalProfit / float64(winningTrades)
	}
	
	if losingTrades > 0 {
		stats.AverageLoss = totalLoss / float64(losingTrades)
	}
	
	// Calculate max drawdown
	peak := stats.StartBalance
	maxDD := 0.0
	balance := stats.StartBalance
	
	for _, trade := range ptm.trades {
		if trade.Status != "open" {
			balance += trade.Profit
			if balance > peak {
				peak = balance
			}
			dd := ((peak - balance) / peak) * 100
			if dd > maxDD {
				maxDD = dd
			}
		}
	}
	stats.MaxDrawdown = maxDD
	
	return stats
}

// GetAllTrades returns all trades
func (ptm *PaperTradingManager) GetAllTrades() []PaperTrade {
	ptm.mu.RLock()
	defer ptm.mu.RUnlock()
	return ptm.trades
}

// ResetPaperTrading resets all paper trading data
func (ptm *PaperTradingManager) ResetPaperTrading() {
	ptm.mu.Lock()
	defer ptm.mu.Unlock()
	
	ptm.trades = []PaperTrade{}
	ptm.currentBalance = ptm.startBalance
	ptm.saveTrades()
}
