package database

import (
	"time"
)

// Candle represents processed candlestick data
type Candle struct {
	Timestamp int64
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}

// LiquiditySweepResult represents a liquidity sweep detection
type LiquiditySweepResult struct {
	Type       string  // "buyside" or "sellside"
	SweepPrice float64
	Confirmed  bool
}

type TradingSignal struct {
	ID                   int64     `json:"id"`
	SignalID             string    `json:"signal_id"`
	CreatedAt            time.Time `json:"created_at"`
	SignalType           string    `json:"signal_type"`
	Symbol               string    `json:"symbol"`
	EntryPrice           float64   `json:"entry_price"`
	StopLoss             float64   `json:"stop_loss"`
	TP1                  float64   `json:"tp1"`
	TP2                  float64   `json:"tp2"`
	TP3                  float64   `json:"tp3"`
	Strength             int       `json:"strength"`
	PatternType          *string   `json:"pattern_type,omitempty"`
	PatternConfidence    *float64  `json:"pattern_confidence,omitempty"`
	KillZone             *string   `json:"kill_zone,omitempty"`
	SessionType          *string   `json:"session_type,omitempty"`
	Status               string    `json:"status"`
	ExitPrice            *float64  `json:"exit_price,omitempty"`
	ExitReason           *string   `json:"exit_reason,omitempty"`
	ExitTime             *time.Time `json:"exit_time,omitempty"`
	ProfitPercent        *float64  `json:"profit_percent,omitempty"`
	ProfitPips           *float64  `json:"profit_pips,omitempty"`
	HoldingTimeMinutes   *int      `json:"holding_time_minutes,omitempty"`
	TrailingStopPrice    *float64  `json:"trailing_stop_price,omitempty"`
	TrailingStopActive   bool      `json:"trailing_stop_active"`
	HighestPrice         *float64  `json:"highest_price,omitempty"`
	LowestPrice          *float64  `json:"lowest_price,omitempty"`
	LivePrice            *float64  `json:"live_price,omitempty"`
	LastUpdated          time.Time `json:"last_updated"`
}

type SignalAnalytics struct {
	KillZone         *string  `json:"kill_zone"`
	SignalType       *string  `json:"signal_type"`
	PatternType      *string  `json:"pattern_type"`
	TotalSignals     int      `json:"total_signals"`
	Wins             int      `json:"wins"`
	Losses           int      `json:"losses"`
	AvgProfitPercent *float64 `json:"avg_profit_percent"`
	AvgWinPercent    *float64 `json:"avg_win_percent"`
	AvgLossPercent   *float64 `json:"avg_loss_percent"`
	AvgHoldingMinutes *float64 `json:"avg_holding_minutes"`
	WinRate          *float64 `json:"win_rate"`
}

type CreateSignalRequest struct {
	SignalID          string   `json:"signal_id"`
	SignalType        string   `json:"signal_type"`
	Symbol            string   `json:"symbol"`
	EntryPrice        float64  `json:"entry_price"`
	StopLoss          float64  `json:"stop_loss"`
	TP1               float64  `json:"tp1"`
	TP2               float64  `json:"tp2"`
	TP3               float64  `json:"tp3"`
	Strength          int      `json:"strength"`
	PatternType       *string  `json:"pattern_type,omitempty"`
	PatternConfidence *float64 `json:"pattern_confidence,omitempty"`
	KillZone          *string  `json:"kill_zone,omitempty"`
	SessionType       *string  `json:"session_type,omitempty"`
}

type UpdateSignalRequest struct {
	Status             string   `json:"status"`
	ExitPrice          *float64 `json:"exit_price,omitempty"`
	ExitReason         *string  `json:"exit_reason,omitempty"`
	ProfitPercent      *float64 `json:"profit_percent,omitempty"`
	ProfitPips         *float64 `json:"profit_pips,omitempty"`
	HoldingTimeMinutes *int     `json:"holding_time_minutes,omitempty"`
	TrailingStopPrice  *float64 `json:"trailing_stop_price,omitempty"`
}

type UpdateLivePriceRequest struct {
	LivePrice         float64  `json:"live_price"`
	TrailingStopPrice *float64 `json:"trailing_stop_price,omitempty"`
}
