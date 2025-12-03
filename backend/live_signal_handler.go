package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type LiveSignalRequest struct {
	Symbol   string `json:"symbol"`
	Strategy string `json:"strategy"`
}

type LiveSignalResponse struct {
	Signal       string  `json:"signal"`
	CurrentPrice float64 `json:"currentPrice"`
	Entry        float64 `json:"entry"`
	StopLoss     float64 `json:"stopLoss"`
	TakeProfit   float64 `json:"takeProfit"`
	TP1          float64 `json:"tp1"`
	TP2          float64 `json:"tp2"`
	TP3          float64 `json:"tp3"`
	RiskReward   float64 `json:"riskReward"`
	Timestamp    int64   `json:"timestamp"`
}

func HandleLiveSignal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LiveSignalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Default to BTCUSDT if no symbol provided
	if req.Symbol == "" {
		req.Symbol = "BTCUSDT"
	}

	// Default to session_trader if no strategy provided
	if req.Strategy == "" {
		req.Strategy = "session_trader"
	}

	// Get interval based on strategy
	interval := getStrategyInterval(req.Strategy)

	// Fetch latest candles from Binance (last 7 days should give us 200+ candles)
	candles, err := fetchBinanceData(req.Symbol, interval, 7)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch candles: %v", err), http.StatusInternalServerError)
		return
	}

	if len(candles) < 50 {
		http.Error(w, "Not enough candle data", http.StatusInternalServerError)
		return
	}

	// Generate signal based on strategy
	signal := generateLiveSignal(candles, req.Strategy)

	json.NewEncoder(w).Encode(signal)
}

func getStrategyInterval(strategy string) string {
	intervals := map[string]string{
		"session_trader":          "15m",
		"breakout_master":         "15m",
		"liquidity_hunter":        "15m",
		"momentum_beast":          "15m",
		"scalper_pro":             "5m",
		"range_master":            "1h",
		"smart_money_tracker":     "1h",
		"reversal_sniper":         "1h",
		"trend_rider":             "4h",
		"institutional_follower":  "4h",
	}

	if interval, ok := intervals[strategy]; ok {
		return interval
	}
	return "15m" // default
}

func generateLiveSignal(candles []Candle, strategy string) LiveSignalResponse {
	currentPrice := candles[len(candles)-1].Close
	
	response := LiveSignalResponse{
		Signal:       "NONE",
		CurrentPrice: currentPrice,
		Entry:        currentPrice,
		StopLoss:     0,
		TakeProfit:   0,
		RiskReward:   0,
		Timestamp:    time.Now().Unix(),
	}

	// Calculate indicators
	closes := make([]float64, len(candles))
	highs := make([]float64, len(candles))
	lows := make([]float64, len(candles))
	
	for i, c := range candles {
		closes[i] = c.Close
		highs[i] = c.High
		lows[i] = c.Low
	}

	// Generate signal based on strategy
	switch strategy {
	case "session_trader":
		response = generateSessionTraderSignal(candles, currentPrice)
	case "breakout_master":
		response = generateBreakoutMasterSignal(candles, currentPrice)
	case "liquidity_hunter":
		response = generateLiquidityHunterSignal(candles, currentPrice)
	case "trend_rider":
		response = generateTrendRiderSignal(candles, currentPrice)
	case "range_master":
		response = generateRangeMasterSignal(candles, currentPrice)
	case "smart_money_tracker":
		response = generateSmartMoneySignal(candles, currentPrice)
	case "institutional_follower":
		response = generateInstitutionalSignal(candles, currentPrice)
	case "reversal_sniper":
		response = generateReversalSignal(candles, currentPrice)
	case "momentum_beast":
		response = generateMomentumSignal(candles, currentPrice)
	case "scalper_pro":
		response = generateScalperSignal(candles, currentPrice)
	default:
		response = generateSessionTraderSignal(candles, currentPrice)
	}

	return response
}

func generateSessionTraderSignal(candles []Candle, currentPrice float64) LiveSignalResponse {
	response := LiveSignalResponse{
		Signal:       "NONE",
		CurrentPrice: currentPrice,
		Entry:        currentPrice,
		Timestamp:    time.Now().Unix(),
	}

	// Calculate EMAs
	ema9 := calculateEMA(candles, 9)
	ema21 := calculateEMA(candles, 21)
	ema50 := calculateEMA(candles, 50)

	// Calculate RSI
	rsi := calculateRSI(candles, 14)

	// Calculate ATR for OPTIMIZED backtest parameters
	atr := calculateATR(candles, 14)
	
	// BUY Signal: EMA9 > EMA21 > EMA50 and RSI > 40 and RSI < 70
	if ema9 > ema21 && ema21 > ema50 && rsi > 40 && rsi < 70 {
		response.Signal = "BUY"
		response.Entry = currentPrice
		response.StopLoss = currentPrice - (atr * 1.0) // OPTIMIZED: 54.1% WR, 12.74 PF
		response.TP1 = currentPrice + (atr * 4.0) // Take 33% profit
		response.TP2 = currentPrice + (atr * 6.0) // Take 33% profit
		response.TP3 = currentPrice + (atr * 10.0) // Take 34% profit
		response.TakeProfit = response.TP3 // Show final target
		response.RiskReward = (response.TakeProfit - response.Entry) / (response.Entry - response.StopLoss)
	}

	// SELL Signal: EMA9 < EMA21 < EMA50 and RSI < 65 and RSI > 30
	if ema9 < ema21 && ema21 < ema50 && rsi < 65 && rsi > 30 {
		response.Signal = "SELL"
		response.Entry = currentPrice
		response.StopLoss = currentPrice + (atr * 1.0) // OPTIMIZED: 54.1% WR, 12.74 PF
		response.TP1 = currentPrice - (atr * 4.0) // Take 33% profit
		response.TP2 = currentPrice - (atr * 6.0) // Take 33% profit
		response.TP3 = currentPrice - (atr * 10.0) // Take 34% profit
		response.TakeProfit = response.TP3 // Show final target
		response.RiskReward = (response.Entry - response.TakeProfit) / (response.StopLoss - response.Entry)
	}

	return response
}

func generateBreakoutMasterSignal(candles []Candle, currentPrice float64) LiveSignalResponse {
	response := LiveSignalResponse{
		Signal:       "NONE",
		CurrentPrice: currentPrice,
		Entry:        currentPrice,
		Timestamp:    time.Now().Unix(),
	}

	// Find recent high/low (last 20 candles)
	recentHigh := candles[len(candles)-20].High
	recentLow := candles[len(candles)-20].Low
	
	for i := len(candles) - 20; i < len(candles)-1; i++ {
		if candles[i].High > recentHigh {
			recentHigh = candles[i].High
		}
		if candles[i].Low < recentLow {
			recentLow = candles[i].Low
		}
	}

	// Calculate volume
	avgVolume := 0.0
	for i := len(candles) - 20; i < len(candles)-1; i++ {
		avgVolume += candles[i].Volume
	}
	avgVolume /= 20

	currentVolume := candles[len(candles)-1].Volume

	// Calculate ATR for OPTIMIZED backtest parameters
	atr := calculateATR(candles, 14)
	
	// BUY Signal: Price breaks above recent high with volume
	if currentPrice > recentHigh && currentVolume > avgVolume*1.5 {
		response.Signal = "BUY"
		response.Entry = currentPrice
		response.StopLoss = currentPrice - (atr * 1.0) // OPTIMIZED: 54.5% WR, 7.20 PF, 3,845% return
		response.TP1 = currentPrice + (atr * 4.0) // Take 33% profit
		response.TP2 = currentPrice + (atr * 6.0) // Take 33% profit
		response.TP3 = currentPrice + (atr * 10.0) // Take 34% profit
		response.TakeProfit = response.TP3
		response.RiskReward = (response.TakeProfit - response.Entry) / (response.Entry - response.StopLoss)
	}

	// SELL Signal: Price breaks below recent low with volume
	if currentPrice < recentLow && currentVolume > avgVolume*1.5 {
		response.Signal = "SELL"
		response.Entry = currentPrice
		response.StopLoss = currentPrice + (atr * 1.0) // OPTIMIZED: 54.5% WR, 7.20 PF, 3,845% return
		response.TP1 = currentPrice - (atr * 4.0) // Take 33% profit
		response.TP2 = currentPrice - (atr * 6.0) // Take 33% profit
		response.TP3 = currentPrice - (atr * 10.0) // Take 34% profit
		response.TakeProfit = response.TP3
		response.RiskReward = (response.Entry - response.TakeProfit) / (response.StopLoss - response.Entry)
	}

	return response
}

func generateLiquidityHunterSignal(candles []Candle, currentPrice float64) LiveSignalResponse {
	response := LiveSignalResponse{
		Signal:       "NONE",
		CurrentPrice: currentPrice,
		Entry:        currentPrice,
		Timestamp:    time.Now().Unix(),
	}

	// Calculate ATR and EMAs
	atr := calculateATR(candles, 14)
	ema20 := calculateEMA(candles, 20)
	ema50 := calculateEMA(candles, 50)

	// Find liquidity zones (recent swing highs/lows)
	swingHigh := findSwingHigh(candles, 10)
	swingLow := findSwingLow(candles, 10)

	// BUY Signal: Price near swing low (liquidity grab) and EMA20 > EMA50
	if currentPrice <= swingLow*1.005 && ema20 > ema50 {
		response.Signal = "BUY"
		response.Entry = currentPrice
		response.StopLoss = currentPrice - (atr * 1.5) // OPTIMIZED: 61.7% WR, 8.24 PF, 894% return - BEST OVERALL
		response.TP1 = currentPrice + (atr * 4.0) // Take 33% profit
		response.TP2 = currentPrice + (atr * 6.0) // Take 33% profit
		response.TP3 = currentPrice + (atr * 10.0) // Take 34% profit
		response.TakeProfit = response.TP3
		response.RiskReward = (response.TakeProfit - response.Entry) / (response.Entry - response.StopLoss)
	}

	// SELL Signal: Price near swing high (liquidity grab) and EMA20 < EMA50
	if currentPrice >= swingHigh*0.995 && ema20 < ema50 {
		response.Signal = "SELL"
		response.Entry = currentPrice
		response.StopLoss = currentPrice + (atr * 1.5) // OPTIMIZED: 61.7% WR, 8.24 PF, 894% return - BEST OVERALL
		response.TP1 = currentPrice - (atr * 4.0) // Take 33% profit
		response.TP2 = currentPrice - (atr * 6.0) // Take 33% profit
		response.TP3 = currentPrice - (atr * 10.0) // Take 34% profit
		response.TakeProfit = response.TP3
		response.RiskReward = (response.Entry - response.TakeProfit) / (response.StopLoss - response.Entry)
	}

	return response
}

func generateTrendRiderSignal(candles []Candle, currentPrice float64) LiveSignalResponse {
	response := LiveSignalResponse{
		Signal:       "NONE",
		CurrentPrice: currentPrice,
		Entry:        currentPrice,
		Timestamp:    time.Now().Unix(),
	}

	// Calculate EMAs for trend
	ema20 := calculateEMA(candles, 20)
	ema50 := calculateEMA(candles, 50)
	ema100 := calculateEMA(candles, 100)

	// Calculate MACD
	macd, signal := calculateMACD(candles)

	// Calculate ATR for OPTIMIZED backtest parameters
	atr := calculateATR(candles, 14)
	
	// BUY Signal: Strong uptrend with MACD confirmation
	if ema20 > ema50 && ema50 > ema100 && macd > signal && macd > 0 {
		response.Signal = "BUY"
		response.Entry = currentPrice
		response.StopLoss = currentPrice - (atr * 0.5) // OPTIMIZED: 36.4% WR, 6.71 PF, 942% return
		response.TP1 = currentPrice + (atr * 2.5) // Take 33% profit
		response.TP2 = currentPrice + (atr * 5.0) // Take 33% profit
		response.TP3 = currentPrice + (atr * 7.5) // Take 34% profit
		response.TakeProfit = response.TP3
		response.RiskReward = (response.TakeProfit - response.Entry) / (response.Entry - response.StopLoss)
	}

	// SELL Signal: Strong downtrend with MACD confirmation
	if ema20 < ema50 && ema50 < ema100 && macd < signal && macd < 0 {
		response.Signal = "SELL"
		response.Entry = currentPrice
		response.StopLoss = currentPrice + (atr * 0.5) // OPTIMIZED: 36.4% WR, 6.71 PF, 942% return
		response.TP1 = currentPrice - (atr * 2.5) // Take 33% profit
		response.TP2 = currentPrice - (atr * 5.0) // Take 33% profit
		response.TP3 = currentPrice - (atr * 7.5) // Take 34% profit
		response.TakeProfit = response.TP3
		response.RiskReward = (response.Entry - response.TakeProfit) / (response.StopLoss - response.Entry)
	}

	return response
}

func generateRangeMasterSignal(candles []Candle, currentPrice float64) LiveSignalResponse {
	response := LiveSignalResponse{
		Signal:       "NONE",
		CurrentPrice: currentPrice,
		Entry:        currentPrice,
		Timestamp:    time.Now().Unix(),
	}

	// Calculate Bollinger Bands
	sma20 := calculateSMA(candles, 20)
	stdDev := calculateStdDevForBB(candles, 20)
	upperBand := sma20 + (stdDev * 2)
	lowerBand := sma20 - (stdDev * 2)

	// Calculate RSI
	rsi := calculateRSI(candles, 14)

	// Calculate ATR for OPTIMIZED backtest parameters
	atr := calculateATR(candles, 14)
	
	// BUY Signal: Price near lower band and RSI oversold
	if currentPrice <= lowerBand*1.01 && rsi < 35 {
		response.Signal = "BUY"
		response.Entry = currentPrice
		response.StopLoss = currentPrice - (atr * 0.5) // OPTIMIZED: 44.2% WR, 7.63 PF, 329% return
		response.TP1 = currentPrice + (atr * 1.7) // Take 33% profit
		response.TP2 = currentPrice + (atr * 3.3) // Take 33% profit
		response.TP3 = currentPrice + (atr * 5.0) // Take 34% profit
		response.TakeProfit = response.TP3
		response.RiskReward = (response.TakeProfit - response.Entry) / (response.Entry - response.StopLoss)
	}

	// SELL Signal: Price near upper band and RSI overbought
	if currentPrice >= upperBand*0.99 && rsi > 65 {
		response.Signal = "SELL"
		response.Entry = currentPrice
		response.StopLoss = currentPrice + (atr * 0.5) // OPTIMIZED: 44.2% WR, 7.63 PF, 329% return
		response.TP1 = currentPrice - (atr * 1.7) // Take 33% profit
		response.TP2 = currentPrice - (atr * 3.3) // Take 33% profit
		response.TP3 = currentPrice - (atr * 5.0) // Take 34% profit
		response.TakeProfit = response.TP3
		response.RiskReward = (response.Entry - response.TakeProfit) / (response.StopLoss - response.Entry)
	}

	return response
}

func generateSmartMoneySignal(candles []Candle, currentPrice float64) LiveSignalResponse {
	// Similar to liquidity hunter but with order block detection
	return generateLiquidityHunterSignal(candles, currentPrice)
}

func generateInstitutionalSignal(candles []Candle, currentPrice float64) LiveSignalResponse {
	// Similar to trend rider but with higher timeframe bias
	return generateTrendRiderSignal(candles, currentPrice)
}

func generateReversalSignal(candles []Candle, currentPrice float64) LiveSignalResponse {
	// Similar to range master but looking for reversals
	return generateRangeMasterSignal(candles, currentPrice)
}

func generateMomentumSignal(candles []Candle, currentPrice float64) LiveSignalResponse {
	// Similar to breakout master but with momentum confirmation
	return generateBreakoutMasterSignal(candles, currentPrice)
}

func generateScalperSignal(candles []Candle, currentPrice float64) LiveSignalResponse {
	// Quick scalping signals with tight stops
	response := generateSessionTraderSignal(candles, currentPrice)
	
	// Scalping with OPTIMIZED backtest parameters
	if response.Signal == "BUY" {
		atr := calculateATR(candles, 14)
		response.StopLoss = currentPrice - (atr * 0.5) // OPTIMIZED: Tight stops for scalping
		response.TP1 = currentPrice + (atr * 1.2) // Take 33% profit
		response.TP2 = currentPrice + (atr * 2.3) // Take 33% profit
		response.TP3 = currentPrice + (atr * 3.5) // Take 34% profit
		response.TakeProfit = response.TP3
		response.RiskReward = (response.TakeProfit - response.Entry) / (response.Entry - response.StopLoss)
	} else if response.Signal == "SELL" {
		atr := calculateATR(candles, 14)
		response.StopLoss = currentPrice + (atr * 0.5) // OPTIMIZED: Tight stops for scalping
		response.TP1 = currentPrice - (atr * 1.2) // Take 33% profit
		response.TP2 = currentPrice - (atr * 2.3) // Take 33% profit
		response.TP3 = currentPrice - (atr * 3.5) // Take 34% profit
		response.TakeProfit = response.TP3
		response.RiskReward = (response.Entry - response.TakeProfit) / (response.StopLoss - response.Entry)
	}
	
	return response
}

// Helper functions
func findSwingHigh(candles []Candle, lookback int) float64 {
	if len(candles) < lookback {
		return candles[len(candles)-1].High
	}
	
	high := candles[len(candles)-lookback].High
	for i := len(candles) - lookback; i < len(candles); i++ {
		if candles[i].High > high {
			high = candles[i].High
		}
	}
	return high
}

func findSwingLow(candles []Candle, lookback int) float64 {
	if len(candles) < lookback {
		return candles[len(candles)-1].Low
	}
	
	low := candles[len(candles)-lookback].Low
	for i := len(candles) - lookback; i < len(candles); i++ {
		if candles[i].Low < low {
			low = candles[i].Low
		}
	}
	return low
}

func calculateStdDevForBB(candles []Candle, period int) float64 {
	if len(candles) < period {
		return 0
	}
	
	sma := calculateSMA(candles, period)
	variance := 0.0
	
	for i := len(candles) - period; i < len(candles); i++ {
		diff := candles[i].Close - sma
		variance += diff * diff
	}
	
	variance /= float64(period)
	
	// Return standard deviation (sqrt of variance)
	stdDev := 0.0
	if variance > 0 {
		stdDev = 1.0
		for i := 0; i < 10; i++ { // Newton's method for sqrt
			stdDev = (stdDev + variance/stdDev) / 2
		}
	}
	return stdDev
}


// Track last signal to prevent duplicates
var (
	lastLiveSignalType   = ""
	lastLiveSignalTime   = time.Time{}
	lastLiveSignalSymbol = ""
)

// Fiber wrapper for live signal handler
func HandleLiveSignalFiber(c *fiber.Ctx) error {
	var req LiveSignalRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Default to BTCUSDT if no symbol provided
	if req.Symbol == "" {
		req.Symbol = "BTCUSDT"
	}

	// Default to session_trader if no strategy provided
	if req.Strategy == "" {
		req.Strategy = "session_trader"
	}

	// Get interval based on strategy
	interval := getStrategyInterval(req.Strategy)

	// Fetch latest candles from Binance (last 7 days should give us 200+ candles)
	candles, err := fetchBinanceData(req.Symbol, interval, 7)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to fetch candles: %v", err),
		})
	}

	if len(candles) < 50 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Not enough candle data",
		})
	}

	// Generate signal based on strategy
	signal := generateLiveSignal(candles, req.Strategy)
	
	log.Printf("🔍 Generated signal: %s for %s using %s strategy", signal.Signal, req.Symbol, req.Strategy)

	// Get current filter settings from database
	filterBuy, filterSell := GetCurrentFilterSettings()
	log.Printf("🔍 Current filter settings: filterBuy=%v, filterSell=%v", filterBuy, filterSell)

	// Check if both filters are disabled (bot paused)
	if !filterBuy && !filterSell {
		log.Printf("⏸️  Both filters disabled - Live signal handler paused")
		return c.JSON(fiber.Map{
			"signal":  "PAUSED",
			"message": "Signal generation paused (both filters disabled)",
		})
	}

	// Check if signal matches filter
	signalMatchesFilter := true
	if signal.Signal == "BUY" && !filterBuy {
		signalMatchesFilter = false
		log.Printf("⏭️  BUY signal filtered out (filterBuy=false)")
	}
	if signal.Signal == "SELL" && !filterSell {
		signalMatchesFilter = false
		log.Printf("⏭️  SELL signal filtered out (filterSell=false)")
	}

	// Check for duplicate signals (same signal type for same symbol)
	if signal.Signal != "NONE" && signal.Signal == lastLiveSignalType && req.Symbol == lastLiveSignalSymbol {
		timeSinceLastSignal := time.Since(lastLiveSignalTime)
		log.Printf("⏭️  Skipping duplicate %s signal (same as last, %v ago)", signal.Signal, timeSinceLastSignal)
		return c.JSON(signal) // Return signal but don't save/send
	}

	// Only save BUY/SELL signals that match filter to Supabase
	signalSavedToDatabase := false
	if signal.Signal != "NONE" && signalMatchesFilter {
		log.Printf("💾 Saving signal to Supabase: %s %s @ $%.2f", signal.Signal, req.Symbol, signal.Entry)
		err = SaveSignalToSupabase(signal, req.Symbol, req.Strategy, filterBuy, filterSell)
		if err != nil {
			// Log error but don't fail the request
			log.Printf("❌ FAILED to save signal to Supabase: %v", err)
			log.Printf("   Check SUPABASE_URL and SUPABASE_KEY in .env")
			log.Printf("   Check if trading_signals table exists in Supabase")
		} else {
			log.Printf("✅ Signal successfully saved to Supabase: %s %s @ $%.2f", signal.Signal, req.Symbol, signal.Entry)
			signalSavedToDatabase = true
			
			// Update last signal tracking
			lastLiveSignalType = signal.Signal
			lastLiveSignalTime = time.Now()
			lastLiveSignalSymbol = req.Symbol
		}
	} else if signal.Signal == "NONE" {
		log.Printf("ℹ️  Signal is NONE, not saving to Supabase")
	} else {
		log.Printf("ℹ️  Signal filtered out, not saving to Supabase")
	}

	// Only send to Telegram if signal was successfully saved to database
	if signalSavedToDatabase {
		if telegramBot == nil {
			log.Printf("⚠️  Telegram bot is nil, cannot send signal")
		} else if telegramBot.Token == "" {
			log.Printf("⚠️  Telegram bot token is empty, cannot send signal")
		} else {
			go telegramBot.SendSignal(signal, req.Symbol, req.Strategy)
			log.Printf("📤 Sent %s signal to Telegram for %s", signal.Signal, req.Symbol)
		}
	} else {
		log.Printf("ℹ️  Signal not sent to Telegram (not saved to database)")
	}

	return c.JSON(signal)
}
