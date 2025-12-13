package main

import (
	"testing"
)

func TestValidateBacktestRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     BacktestConfig
		wantErr bool
	}{
		{
			name: "valid request",
			req: BacktestConfig{
				Symbol:       "BTCUSDT",
				Interval:     "15m",
				Days:         30,
				StartBalance: 500,
			},
			wantErr: false,
		},
		{
			name: "empty symbol",
			req: BacktestConfig{
				Symbol:       "",
				Interval:     "15m",
				Days:         30,
				StartBalance: 500,
			},
			wantErr: true,
		},
		{
			name: "invalid interval",
			req: BacktestConfig{
				Symbol:       "BTCUSDT",
				Interval:     "30m",
				Days:         30,
				StartBalance: 500,
			},
			wantErr: true,
		},
		{
			name: "days too high",
			req: BacktestConfig{
				Symbol:       "BTCUSDT",
				Interval:     "15m",
				Days:         400,
				StartBalance: 500,
			},
			wantErr: true,
		},
		{
			name: "balance too low",
			req: BacktestConfig{
				Symbol:       "BTCUSDT",
				Interval:     "15m",
				Days:         30,
				StartBalance: 5,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBacktestRequest(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateBacktestRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateSignalRequest(t *testing.T) {
	tests := []struct {
		name      string
		symbol    string
		timeframe string
		wantErr   bool
	}{
		{
			name:      "valid request",
			symbol:    "BTCUSDT",
			timeframe: "15m",
			wantErr:   false,
		},
		{
			name:      "empty symbol",
			symbol:    "",
			timeframe: "15m",
			wantErr:   true,
		},
		{
			name:      "invalid timeframe",
			symbol:    "BTCUSDT",
			timeframe: "30m",
			wantErr:   true,
		},
		{
			name:      "empty timeframe is ok",
			symbol:    "BTCUSDT",
			timeframe: "",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSignalRequest(tt.symbol, tt.timeframe)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSignalRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
