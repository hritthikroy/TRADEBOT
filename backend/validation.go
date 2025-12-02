package main

import (
	"fmt"
	"strings"
)

// ValidateBacktestRequest validates backtest request parameters
func ValidateBacktestRequest(req BacktestConfig) error {
	// Validate symbol
	if req.Symbol == "" {
		return fmt.Errorf("symbol is required")
	}
	req.Symbol = strings.ToUpper(req.Symbol)
	if !strings.HasSuffix(req.Symbol, "USDT") {
		return fmt.Errorf("only USDT pairs are supported")
	}

	// Validate interval
	validIntervals := map[string]bool{
		"1m": true, "5m": true, "15m": true,
		"1h": true, "4h": true, "1d": true,
	}
	if !validIntervals[req.Interval] {
		return fmt.Errorf("invalid interval: %s (valid: 1m, 5m, 15m, 1h, 4h, 1d)", req.Interval)
	}

	// Validate days
	if req.Days < 1 || req.Days > 365 {
		return fmt.Errorf("days must be between 1 and 365")
	}

	// Validate start balance
	if req.StartBalance < 10 || req.StartBalance > 1000000 {
		return fmt.Errorf("start balance must be between 10 and 1,000,000")
	}

	return nil
}

// ValidateSignalRequest validates signal request parameters
func ValidateSignalRequest(symbol, timeframe string) error {
	if symbol == "" {
		return fmt.Errorf("symbol is required")
	}

	validTimeframes := map[string]bool{
		"1m": true, "5m": true, "15m": true,
		"1h": true, "4h": true, "1d": true,
	}
	if timeframe != "" && !validTimeframes[timeframe] {
		return fmt.Errorf("invalid timeframe: %s", timeframe)
	}

	return nil
}
