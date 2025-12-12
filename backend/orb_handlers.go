package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ORBBacktestRequest represents the request for an ORB backtest
type ORBBacktestRequest struct {
	TimeFrame      int       `json:"timeFrame"`      // 5, 15, 30, or 60 minutes
	StartDate      string    `json:"startDate"`      // YYYY-MM-DD
	EndDate        string    `json:"endDate"`        // YYYY-MM-DD
	InitialCapital float64   `json:"initialCapital"` // Starting capital
	Symbols        []string  `json:"symbols"`        // List of symbols to backtest (optional, empty = all)
	TopNStocks     int       `json:"topNStocks"`     // Number of top stocks to trade (default 20)
	MinRelativeVol float64   `json:"minRelativeVol"` // Minimum relative volume (default 1.0)
}

// ORBBacktestResponse represents the response from an ORB backtest
type ORBBacktestResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message,omitempty"`
	Result  *ORBBacktestResult  `json:"result,omitempty"`
	Summary *ORBBacktestSummary `json:"summary,omitempty"`
}

// ORBBacktestSummary provides a high-level summary of backtest results
type ORBBacktestSummary struct {
	Strategy         string  `json:"strategy"`
	TimeFrame        int     `json:"timeFrame"`
	Period           string  `json:"period"`
	TotalReturn      string  `json:"totalReturn"`
	AnnualizedReturn string  `json:"annualizedReturn"`
	SharpeRatio      float64 `json:"sharpeRatio"`
	MaxDrawdown      string  `json:"maxDrawdown"`
	WinRate          string  `json:"winRate"`
	TotalTrades      int     `json:"totalTrades"`
	AvgPnLInR        float64 `json:"avgPnLInR"`
	Alpha            float64 `json:"alpha"`
	Beta             float64 `json:"beta"`
}

// HandleORBBacktest handles the ORB backtest endpoint
func HandleORBBacktest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ORBBacktestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := validateORBBacktestRequest(&req); err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse dates
	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	// Set defaults
	if req.InitialCapital == 0 {
		req.InitialCapital = 25000
	}
	if req.TopNStocks == 0 {
		req.TopNStocks = 20
	}
	if req.MinRelativeVol == 0 {
		req.MinRelativeVol = 1.0
	}

	// Run backtest
	result, err := runORBBacktestFromRequest(&req, startDate, endDate)
	if err != nil {
		respondWithError(w, fmt.Sprintf("Backtest failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Create summary
	summary := createORBBacktestSummary(result)

	response := ORBBacktestResponse{
		Success: true,
		Result:  result,
		Summary: summary,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleORBLiveSignals handles the endpoint for generating live ORB signals
func HandleORBLiveSignals(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get query parameters
	timeFrameStr := r.URL.Query().Get("timeframe")
	if timeFrameStr == "" {
		timeFrameStr = "5"
	}

	var timeFrame int
	fmt.Sscanf(timeFrameStr, "%d", &timeFrame)

	if timeFrame != 5 && timeFrame != 15 && timeFrame != 30 && timeFrame != 60 {
		respondWithError(w, "Invalid timeframe. Must be 5, 15, 30, or 60", http.StatusBadRequest)
		return
	}

	// Generate live signals for today
	signals, err := generateLiveORBSignals(timeFrame)
	if err != nil {
		respondWithError(w, fmt.Sprintf("Failed to generate signals: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success":   true,
		"timeFrame": timeFrame,
		"date":      time.Now().Format("2006-01-02"),
		"signals":   signals,
		"count":     len(signals),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleORBCompareTimeframes compares ORB performance across different timeframes
func HandleORBCompareTimeframes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ORBBacktestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	if req.InitialCapital == 0 {
		req.InitialCapital = 25000
	}

	// Run backtests for all timeframes
	timeframes := []int{5, 15, 30, 60}
	results := make(map[string]*ORBBacktestSummary)

	for _, tf := range timeframes {
		req.TimeFrame = tf
		result, err := runORBBacktestFromRequest(&req, startDate, endDate)
		if err != nil {
			continue
		}
		summary := createORBBacktestSummary(result)
		results[fmt.Sprintf("%dm", tf)] = summary
	}

	response := map[string]interface{}{
		"success": true,
		"period":  fmt.Sprintf("%s to %s", req.StartDate, req.EndDate),
		"results": results,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleORBTopPerformers returns the top performing stocks for ORB strategy
func HandleORBTopPerformers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	timeFrameStr := r.URL.Query().Get("timeframe")
	if timeFrameStr == "" {
		timeFrameStr = "5"
	}

	var timeFrame int
	fmt.Sscanf(timeFrameStr, "%d", &timeFrame)

	// This would query historical data to find top performers
	// For now, return the top performers from the paper
	topPerformers := getTopPerformersFromPaper(timeFrame)

	response := map[string]interface{}{
		"success":       true,
		"timeFrame":     timeFrame,
		"topPerformers": topPerformers,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper functions

func validateORBBacktestRequest(req *ORBBacktestRequest) error {
	if req.TimeFrame != 5 && req.TimeFrame != 15 && req.TimeFrame != 30 && req.TimeFrame != 60 {
		return fmt.Errorf("invalid timeframe: must be 5, 15, 30, or 60 minutes")
	}

	if req.StartDate == "" || req.EndDate == "" {
		return fmt.Errorf("start date and end date are required")
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return fmt.Errorf("invalid start date format: use YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return fmt.Errorf("invalid end date format: use YYYY-MM-DD")
	}

	if endDate.Before(startDate) {
		return fmt.Errorf("end date must be after start date")
	}

	return nil
}

func runORBBacktestFromRequest(req *ORBBacktestRequest, startDate, endDate time.Time) (*ORBBacktestResult, error) {
	// This is a placeholder - in production, you would:
	// 1. Fetch historical intraday data for all stocks
	// 2. Calculate ATR, average volume, etc.
	// 3. Run the backtest engine
	
	// For now, return a mock result based on the paper's findings
	return createMockORBResult(req.TimeFrame, startDate, endDate, req.InitialCapital), nil
}

func createORBBacktestSummary(result *ORBBacktestResult) *ORBBacktestSummary {
	return &ORBBacktestSummary{
		Strategy:         result.Strategy,
		TimeFrame:        result.TimeFrame,
		Period:           fmt.Sprintf("%s to %s", result.StartDate.Format("2006-01-02"), result.EndDate.Format("2006-01-02")),
		TotalReturn:      fmt.Sprintf("%.2f%%", result.TotalReturn*100),
		AnnualizedReturn: fmt.Sprintf("%.2f%%", result.AnnualizedReturn*100),
		SharpeRatio:      result.SharpeRatio,
		MaxDrawdown:      fmt.Sprintf("%.2f%%", result.MaxDrawdown*100),
		WinRate:          fmt.Sprintf("%.2f%%", result.WinRate*100),
		TotalTrades:      result.TotalTrades,
		AvgPnLInR:        result.AvgPnLInR,
		Alpha:            result.Alpha,
		Beta:             result.Beta,
	}
}

func generateLiveORBSignals(timeFrame int) ([]*ORBSignal, error) {
	// This would fetch live market data and generate signals
	// Placeholder implementation
	return []*ORBSignal{}, nil
}

func getTopPerformersFromPaper(timeFrame int) []map[string]interface{} {
	// Data from Table 4 in the paper
	performers := make([]map[string]interface{}, 0)

	switch timeFrame {
	case 5:
		performers = []map[string]interface{}{
			{"symbol": "DDD", "pnlR": 385, "winRate": 21},
			{"symbol": "FSLR", "pnlR": 370, "winRate": 20},
			{"symbol": "NVDA", "pnlR": 309, "winRate": 19},
			{"symbol": "SWBI", "pnlR": 272, "winRate": 24},
			{"symbol": "RCL", "pnlR": 271, "winRate": 20},
		}
	case 15:
		performers = []map[string]interface{}{
			{"symbol": "CAR", "pnlR": 233, "winRate": 21},
			{"symbol": "NVDA", "pnlR": 200, "winRate": 20},
			{"symbol": "AMD", "pnlR": 189, "winRate": 19},
			{"symbol": "LITE", "pnlR": 187, "winRate": 22},
			{"symbol": "FOSL", "pnlR": 187, "winRate": 20},
		}
	case 30:
		performers = []map[string]interface{}{
			{"symbol": "MXIM", "pnlR": 214, "winRate": 26},
			{"symbol": "SAVE", "pnlR": 213, "winRate": 23},
			{"symbol": "ACAD", "pnlR": 201, "winRate": 21},
			{"symbol": "CDNS", "pnlR": 196, "winRate": 24},
			{"symbol": "WOLF", "pnlR": 185, "winRate": 25},
		}
	case 60:
		performers = []map[string]interface{}{
			{"symbol": "DDD", "pnlR": 185, "winRate": 24},
			{"symbol": "THC", "pnlR": 183, "winRate": 26},
			{"symbol": "TKAT", "pnlR": 177, "winRate": 40},
			{"symbol": "DISH", "pnlR": 169, "winRate": 25},
			{"symbol": "EXEL", "pnlR": 162, "winRate": 25},
		}
	}

	return performers
}

func createMockORBResult(timeFrame int, startDate, endDate time.Time, initialCapital float64) *ORBBacktestResult {
	// Mock results based on the paper's findings (Table 3)
	var totalReturn, annualizedReturn, volatility, sharpeRatio, winRate, maxDrawdown, alpha, beta float64

	switch timeFrame {
	case 5:
		totalReturn = 16.37      // 1,637%
		annualizedReturn = 0.416 // 41.6%
		volatility = 0.148       // 14.8%
		sharpeRatio = 2.81
		winRate = 0.484 // 48.4%
		maxDrawdown = 0.12
		alpha = 0.358 // 35.8%
		beta = 0.00
	case 15:
		totalReturn = 2.72       // 272%
		annualizedReturn = 0.174 // 17.4%
		volatility = 0.122       // 12.2%
		sharpeRatio = 1.43
		winRate = 0.447 // 44.7%
		maxDrawdown = 0.11
		alpha = 0.169 // 16.9%
		beta = -0.01
	case 30:
		totalReturn = 0.21      // 21%
		annualizedReturn = 0.023 // 2.3%
		volatility = 0.111      // 11.1%
		sharpeRatio = 0.21
		winRate = 0.424 // 42.4%
		maxDrawdown = 0.35
		alpha = 0.028 // 2.8%
		beta = 0.01
	case 60:
		totalReturn = 0.39      // 39%
		annualizedReturn = 0.041 // 4.1%
		volatility = 0.102      // 10.2%
		sharpeRatio = 0.40
		winRate = 0.423 // 42.3%
		maxDrawdown = 0.21
		alpha = 0.044 // 4.4%
		beta = 0.01
	}

	finalCapital := initialCapital * (1 + totalReturn)

	return &ORBBacktestResult{
		Strategy:         fmt.Sprintf("%d-minute ORB + Relative Volume", timeFrame),
		TimeFrame:        timeFrame,
		StartDate:        startDate,
		EndDate:          endDate,
		InitialCapital:   initialCapital,
		FinalCapital:     finalCapital,
		TotalReturn:      totalReturn,
		AnnualizedReturn: annualizedReturn,
		Volatility:       volatility,
		SharpeRatio:      sharpeRatio,
		MaxDrawdown:      maxDrawdown,
		WinRate:          winRate,
		Alpha:            alpha,
		Beta:             beta,
	}
}

func respondWithError(w http.ResponseWriter, message string, statusCode int) {
	response := ORBBacktestResponse{
		Success: false,
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
