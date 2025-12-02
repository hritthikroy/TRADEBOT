package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	
	"github.com/gofiber/fiber/v2"
)

// BinanceKlineArray represents Binance API response as array
type BinanceKlineArray []interface{}

// HandleBacktestRun handles backtest execution requests
func HandleBacktestRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var config BacktestConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate config
	if config.Symbol == "" {
		config.Symbol = "BTCUSDT"
	}
	if config.Interval == "" {
		config.Interval = "15m"
	}
	if config.Days == 0 {
		config.Days = 30
	}
	if config.StartBalance == 0 {
		config.StartBalance = 500
	}

	// Fetch historical data from Binance
	candles, err := fetchBinanceData(config.Symbol, config.Interval, config.Days)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch data: %v", err), http.StatusInternalServerError)
		return
	}

	if len(candles) < 100 {
		http.Error(w, "Insufficient historical data", http.StatusBadRequest)
		return
	}

	// Check if enhanced backtest is requested
	useEnhanced := config.WindowType != "" || config.UseWalkForward || config.UseMonteCarlo
	
	var result *BacktestResult
	
	if useEnhanced {
		// Run enhanced backtest with simulation windows
		result, err = RunEnhancedBacktest(config, candles)
	} else {
		// Run standard backtest (backward compatible)
		result, err = RunBacktest(config, candles)
	}
	
	if err != nil {
		http.Error(w, fmt.Sprintf("Backtest failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Return results
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// HandleBacktestExport handles CSV export requests
func HandleBacktestExport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var config BacktestConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate config
	if config.Symbol == "" {
		config.Symbol = "BTCUSDT"
	}
	if config.Interval == "" {
		config.Interval = "15m"
	}
	if config.Days == 0 {
		config.Days = 30
	}
	if config.StartBalance == 0 {
		config.StartBalance = 500
	}

	// Fetch historical data
	candles, err := fetchBinanceData(config.Symbol, config.Interval, config.Days)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch data: %v", err), http.StatusInternalServerError)
		return
	}

	// Run backtest
	result, err := RunBacktest(config, candles)
	if err != nil {
		http.Error(w, fmt.Sprintf("Backtest failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Export to CSV
	csv := result.ToCSV()
	
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=backtest_%s_%dd_%d.csv", 
		config.Symbol, config.Days, time.Now().UnixMilli()))
	w.Write([]byte(csv))
}

// fetchBinanceData fetches historical candle data from Binance
func fetchBinanceData(symbol, interval string, days int) ([]Candle, error) {
	// Convert interval to Binance format
	binanceInterval := toBinanceInterval(interval)
	
	// Calculate how many candles we need
	limit := calculateCandleLimit(binanceInterval, days)
	
	// Binance API endpoint
	url := fmt.Sprintf("https://api.binance.com/api/v3/klines?symbol=%s&interval=%s&limit=%d",
		symbol, binanceInterval, limit)
	
	// Make request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from Binance: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("binance API error: %s", string(body))
	}
	
	// Parse response
	var klines []BinanceKlineArray
	if err := json.NewDecoder(resp.Body).Decode(&klines); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	// Convert to Candle format
	candles := make([]Candle, len(klines))
	for i, k := range klines {
		candles[i] = Candle{
			Timestamp: int64(k[0].(float64)),
			Open:      parseFloatBT(k[1]),
			High:      parseFloatBT(k[2]),
			Low:       parseFloatBT(k[3]),
			Close:     parseFloatBT(k[4]),
			Volume:    parseFloatBT(k[5]),
		}
	}
	
	return candles, nil
}

// toBinanceInterval converts interval to Binance format
func toBinanceInterval(interval string) string {
	binanceIntervals := map[string]string{
		"1s": "1m", // Binance doesn't support 1s
		"1m": "1m",
		"3m": "3m",
		"5m": "5m",
		"15m": "15m",
		"30m": "30m",
		"1h": "1h",
		"2h": "2h",
		"4h": "4h",
		"6h": "6h",
		"8h": "8h",
		"12h": "12h",
		"1d": "1d",
		"3d": "3d",
		"1w": "1w",
		"1D": "1d",
		"3D": "3d",
		"1W": "1w",
	}
	
	if bi, ok := binanceIntervals[interval]; ok {
		return bi
	}
	return interval
}

// calculateCandleLimit calculates how many candles to fetch
func calculateCandleLimit(interval string, days int) int {
	candlesPerDay := map[string]int{
		"1m": 1440,
		"3m": 480,
		"5m": 288,
		"15m": 96,
		"30m": 48,
		"1h": 24,
		"2h": 12,
		"4h": 6,
		"6h": 4,
		"8h": 3,
		"12h": 2,
		"1d": 1,
		"3d": 1,
		"1w": 1,
	}
	
	cpd := candlesPerDay[interval]
	if cpd == 0 {
		cpd = 96 // Default to 15m
	}
	
	needed := cpd * days
	
	// Binance limit is 1000, add 50 for indicators
	if needed > 950 {
		return 1000
	}
	return needed + 50
}

// parseFloatBT safely parses interface{} to float64 for backtest
func parseFloatBT(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case string:
		var f float64
		fmt.Sscanf(val, "%f", &f)
		return f
	default:
		return 0
	}
}


// HandleBacktestRunFiber handles backtest execution with Fiber
func HandleBacktestRunFiber(c *fiber.Ctx) error {
	var config BacktestConfig
	if err := c.BodyParser(&config); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate config
	if config.Symbol == "" {
		config.Symbol = "BTCUSDT"
	}
	if config.Interval == "" {
		config.Interval = "15m"
	}
	if config.Days == 0 {
		config.Days = 30
	}
	if config.StartBalance == 0 {
		config.StartBalance = 500
	}

	// Fetch historical data from Binance
	candles, err := fetchBinanceData(config.Symbol, config.Interval, config.Days)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to fetch data: %v", err),
		})
	}

	if len(candles) < 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Insufficient historical data",
		})
	}

	// Run backtest
	result, err := RunBacktest(config, candles)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Backtest failed: %v", err),
		})
	}

	return c.JSON(result)
}

// HandleBacktestExportFiber handles CSV export with Fiber
func HandleBacktestExportFiber(c *fiber.Ctx) error {
	var config BacktestConfig
	if err := c.BodyParser(&config); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate config
	if config.Symbol == "" {
		config.Symbol = "BTCUSDT"
	}
	if config.Interval == "" {
		config.Interval = "15m"
	}
	if config.Days == 0 {
		config.Days = 30
	}
	if config.StartBalance == 0 {
		config.StartBalance = 500
	}

	// Fetch historical data
	candles, err := fetchBinanceData(config.Symbol, config.Interval, config.Days)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to fetch data: %v", err),
		})
	}

	// Run backtest
	result, err := RunBacktest(config, candles)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Backtest failed: %v", err),
		})
	}

	// Export to CSV
	csv := result.ToCSV()
	
	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=backtest_%s_%dd_%d.csv", 
		config.Symbol, config.Days, time.Now().UnixMilli()))
	
	return c.SendString(csv)
}
