package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"time"
)

// BinanceKline represents a candlestick from Binance API
type BinanceKline struct {
	OpenTime  int64
	Open      string
	High      string
	Low       string
	Close     string
	Volume    string
	CloseTime int64
}

// Candle represents processed candlestick data
type Candle struct {
	Timestamp int64
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}

// SignalGenerator handles automatic signal generation
type SignalGenerator struct {
	Symbol            string
	Timeframes        []string
	CheckInterval     time.Duration
	MinSignalStrength int
	MaxSignalsPerDay  int
	SignalsToday      int
	LastSignalTime    time.Time
	IsRunning         bool
}

// NewSignalGenerator creates a new signal generator
func NewSignalGenerator() *SignalGenerator {
	return &SignalGenerator{
		Symbol:            "BTCUSDT",
		Timeframes:        []string{"1m", "3m", "15m", "1h", "2h", "4h"}, // Added 1m for scalping
		CheckInterval:     60 * time.Second,  // ✅ Every 60 seconds (1 minute)
		MinSignalStrength: 60,
		MaxSignalsPerDay:  999999,  // ♾️ No daily limit - unlimited signals
		SignalsToday:      0,
		LastSignalTime:    time.Time{},
		IsRunning:         false,
	}
}

// FetchMarketData fetches candlestick data from Binance
func (sg *SignalGenerator) FetchMarketData(symbol, interval string, limit int) ([]Candle, error) {
	url := fmt.Sprintf("https://api.binance.com/api/v3/klines?symbol=%s&interval=%s&limit=%d", symbol, interval, limit)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var klines [][]interface{}
	if err := json.Unmarshal(body, &klines); err != nil {
		return nil, err
	}
	
	candles := make([]Candle, 0, len(klines))
	for _, k := range klines {
		candle := Candle{
			Timestamp: int64(k[0].(float64)),
			Open:      parseFloat(k[1].(string)),
			High:      parseFloat(k[2].(string)),
			Low:       parseFloat(k[3].(string)),
			Close:     parseFloat(k[4].(string)),
			Volume:    parseFloat(k[5].(string)),
		}
		candles = append(candles, candle)
	}
	
	return candles, nil
}

// CalculateRSI calculates Relative Strength Index
func (sg *SignalGenerator) CalculateRSI(data []Candle, period int) float64 {
	if len(data) < period+1 {
		return 50.0
	}
	
	var gains, losses float64
	
	for i := len(data) - period; i < len(data); i++ {
		change := data[i].Close - data[i-1].Close
		if change > 0 {
			gains += change
		} else {
			losses += math.Abs(change)
		}
	}
	
	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)
	
	if avgLoss == 0 {
		return 100.0
	}
	
	rs := avgGain / avgLoss
	return 100.0 - (100.0 / (1.0 + rs))
}

// CalculateEMA calculates Exponential Moving Average
func (sg *SignalGenerator) CalculateEMA(data []Candle, period int) float64 {
	if len(data) < period {
		return 0
	}
	
	k := 2.0 / float64(period+1)
	ema := data[0].Close
	
	for i := 1; i < len(data); i++ {
		ema = (data[i].Close * k) + (ema * (1 - k))
	}
	
	return ema
}

// CalculateATR calculates Average True Range
func (sg *SignalGenerator) CalculateATR(data []Candle, period int) float64 {
	if len(data) < period+1 {
		return 0
	}
	
	trueRanges := make([]float64, 0)
	
	for i := 1; i < len(data); i++ {
		high := data[i].High
		low := data[i].Low
		prevClose := data[i-1].Close
		
		tr := math.Max(
			high-low,
			math.Max(
				math.Abs(high-prevClose),
				math.Abs(low-prevClose),
			),
		)
		
		trueRanges = append(trueRanges, tr)
	}
	
	// Get recent true ranges
	start := len(trueRanges) - period
	if start < 0 {
		start = 0
	}
	recentTR := trueRanges[start:]
	
	var sum float64
	for _, tr := range recentTR {
		sum += tr
	}
	
	return sum / float64(len(recentTR))
}

// AnalyzeMarket analyzes market data and generates signal if conditions are met
func (sg *SignalGenerator) AnalyzeMarket(timeframe string) (*CreateSignalRequest, error) {
	data, err := sg.FetchMarketData(sg.Symbol, timeframe, 100)
	if err != nil {
		return nil, err
	}
	
	if len(data) < 50 {
		return nil, fmt.Errorf("not enough data")
	}
	
	currentPrice := data[len(data)-1].Close
	rsi := sg.CalculateRSI(data, 14)
	
	// Calculate EMAs
	ema20Data := data
	if len(data) > 20 {
		ema20Data = data[len(data)-20:]
	}
	ema20 := sg.CalculateEMA(ema20Data, 20)
	
	ema50Data := data
	if len(data) > 50 {
		ema50Data = data[len(data)-50:]
	}
	ema50 := sg.CalculateEMA(ema50Data, 50)
	
	atr := sg.CalculateATR(data, 14)
	
	// Determine trend
	isBullish := ema20 > ema50 && currentPrice > ema20
	isBearish := ema20 < ema50 && currentPrice < ema20
	
	// Calculate signal strength
	strength := 50
	
	// RSI contribution
	if isBullish && rsi < 40 {
		strength += 15 // Oversold in uptrend
	}
	if isBearish && rsi > 60 {
		strength += 15 // Overbought in downtrend
	}
	
	// EMA alignment
	if isBullish && currentPrice > ema20 && ema20 > ema50 {
		strength += 15
	}
	if isBearish && currentPrice < ema20 && ema20 < ema50 {
		strength += 15
	}
	
	// Volume confirmation
	var avgVolume float64
	volumeData := data
	if len(data) > 20 {
		volumeData = data[len(data)-20:]
	}
	for _, c := range volumeData {
		avgVolume += c.Volume
	}
	avgVolume /= float64(len(volumeData))
	
	currentVolume := data[len(data)-1].Volume
	if currentVolume > avgVolume*1.2 {
		strength += 10
	}
	
	// Price action
	lastCandle := data[len(data)-1]
	candleSize := math.Abs(lastCandle.Close - lastCandle.Open)
	
	var avgCandleSize float64
	recentCandles := data
	if len(data) > 10 {
		recentCandles = data[len(data)-10:]
	}
	for _, c := range recentCandles {
		avgCandleSize += math.Abs(c.Close - c.Open)
	}
	avgCandleSize /= float64(len(recentCandles))
	
	if candleSize > avgCandleSize*1.5 {
		strength += 10
	}
	
	// Check if signal meets minimum strength
	if strength < sg.MinSignalStrength {
		log.Printf("⏳ [%s] Signal strength too low: %d%%", timeframe, strength)
		return nil, nil
	}
	
	// Generate signal
	var signal *CreateSignalRequest
	
	if isBullish {
		entry := currentPrice
		stopLoss := entry - (atr * 2)
		tp1 := entry + (atr * 2.5)
		tp2 := entry + (atr * 4.5)
		tp3 := entry + (atr * 7.0)
		
		signalID := fmt.Sprintf("%d.%d", time.Now().UnixMilli(), time.Now().Nanosecond())
		signalType := "BUY"
		killZone := sg.DetectKillZone(time.Now())
		sessionType := sg.DetectSession(time.Now())
		
		signal = &CreateSignalRequest{
			SignalID:    signalID,
			SignalType:  signalType,
			Symbol:      sg.Symbol,
			EntryPrice:  entry,
			StopLoss:    stopLoss,
			TP1:         tp1,
			TP2:         tp2,
			TP3:         tp3,
			Strength:    strength,
			KillZone:    &killZone,
			SessionType: &sessionType,
		}
	} else if isBearish {
		entry := currentPrice
		stopLoss := entry + (atr * 2)
		tp1 := entry - (atr * 2.5)
		tp2 := entry - (atr * 4.5)
		tp3 := entry - (atr * 7.0)
		
		signalID := fmt.Sprintf("%d.%d", time.Now().UnixMilli(), time.Now().Nanosecond())
		signalType := "SELL"
		killZone := sg.DetectKillZone(time.Now())
		sessionType := sg.DetectSession(time.Now())
		
		signal = &CreateSignalRequest{
			SignalID:    signalID,
			SignalType:  signalType,
			Symbol:      sg.Symbol,
			EntryPrice:  entry,
			StopLoss:    stopLoss,
			TP1:         tp1,
			TP2:         tp2,
			TP3:         tp3,
			Strength:    strength,
			KillZone:    &killZone,
			SessionType: &sessionType,
		}
	}
	
	return signal, nil
}

// SaveSignal saves a signal to the database
func (sg *SignalGenerator) SaveSignal(signal *CreateSignalRequest, timeframe string) error {
	query := `
		INSERT INTO trading_signals (
			signal_id, signal_type, symbol, entry_price, stop_loss, 
			tp1, tp2, tp3, strength, kill_zone, session_type, 
			status, trailing_stop_active, timeframe
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	
	_, err := DB.Exec(
		query,
		signal.SignalID,
		signal.SignalType,
		signal.Symbol,
		signal.EntryPrice,
		signal.StopLoss,
		signal.TP1,
		signal.TP2,
		signal.TP3,
		signal.Strength,
		signal.KillZone,
		signal.SessionType,
		"pending",
		false,
		timeframe,
	)
	
	return err
}

// DetectKillZone detects the current kill zone
func (sg *SignalGenerator) DetectKillZone(t time.Time) string {
	hour := t.UTC().Hour()
	
	// London Open: 8:00-10:00 UTC
	if hour >= 8 && hour < 10 {
		return "London Open"
	}
	
	// London Close: 16:00-18:00 UTC
	if hour >= 16 && hour < 18 {
		return "London Close"
	}
	
	// New York Open: 13:00-15:00 UTC
	if hour >= 13 && hour < 15 {
		return "New York Open"
	}
	
	// New York Close: 21:00-23:00 UTC
	if hour >= 21 && hour < 23 {
		return "New York Close"
	}
	
	// Asian Session: 0:00-8:00 UTC
	if hour >= 0 && hour < 8 {
		return "Asian Session"
	}
	
	return "Off Hours"
}

// DetectSession detects the current trading session
func (sg *SignalGenerator) DetectSession(t time.Time) string {
	hour := t.UTC().Hour()
	
	if hour >= 0 && hour < 8 {
		return "Asian"
	}
	if hour >= 8 && hour < 16 {
		return "London"
	}
	if hour >= 13 && hour < 21 {
		return "New York"
	}
	
	return "Off Hours"
}

// GenerateSignals is the main loop for signal generation
func (sg *SignalGenerator) GenerateSignals() {
	// Reset daily counter at midnight
	now := time.Now()
	if now.Hour() == 0 && now.Minute() == 0 {
		sg.SignalsToday = 0
		log.Println("🔄 Daily signal counter reset")
	}
	
	// Check if we've hit daily limit
	if sg.SignalsToday >= sg.MaxSignalsPerDay {
		log.Printf("⏸️ Daily signal limit reached (%d/%d)", sg.SignalsToday, sg.MaxSignalsPerDay)
		return
	}
	
	// No cooldown - generate signals every scan if conditions are met
	
	log.Printf("\n🔍 Scanning %d timeframes...", len(sg.Timeframes))
	
	// Check each timeframe
	for _, timeframe := range sg.Timeframes {
		signal, err := sg.AnalyzeMarket(timeframe)
		
		if err != nil {
			log.Printf("❌ [%s] Error analyzing market: %v", timeframe, err)
			continue
		}
		
		if signal != nil {
			// Save to database
			if err := sg.SaveSignal(signal, timeframe); err != nil {
				log.Printf("❌ [%s] Failed to save signal: %v", timeframe, err)
				continue
			}
			
			sg.SignalsToday++
			sg.LastSignalTime = time.Now()
			
			log.Printf("\n✅ [%s] %s Signal Generated!", timeframe, signal.SignalType)
			log.Printf("   Entry: %.2f", signal.EntryPrice)
			log.Printf("   Stop Loss: %.2f", signal.StopLoss)
			log.Printf("   TP1: %.2f", signal.TP1)
			log.Printf("   TP2: %.2f", signal.TP2)
			log.Printf("   TP3: %.2f", signal.TP3)
			log.Printf("   Strength: %d%%", signal.Strength)
			log.Printf("   Signals today: %d/%d\n", sg.SignalsToday, sg.MaxSignalsPerDay)
			
			// Check if we hit daily limit
			if sg.SignalsToday >= sg.MaxSignalsPerDay {
				log.Println("⏸️ Daily limit reached, stopping scan")
				break
			}
			
			// Continue scanning other timeframes (removed break to allow multiple signals per scan)
		}
	}
}

// Start begins the signal generation loop
func (sg *SignalGenerator) Start() {
	if sg.IsRunning {
		log.Println("⚠️ Signal generator already running")
		return
	}
	
	sg.IsRunning = true
	log.Println("🚀 Automatic Signal Generator Started")
	log.Printf("📊 Symbol: %s", sg.Symbol)
	log.Printf("⏱️ Timeframes: %v", sg.Timeframes)
	log.Printf("🔄 Check interval: %v", sg.CheckInterval)
	log.Printf("📈 Min strength: %d%%", sg.MinSignalStrength)
	log.Printf("🎯 Max signals/day: %d", sg.MaxSignalsPerDay)
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	// Initial scan
	sg.GenerateSignals()
	
	// Periodic scans
	ticker := time.NewTicker(sg.CheckInterval)
	go func() {
		for range ticker.C {
			if sg.IsRunning {
				sg.GenerateSignals()
			}
		}
	}()
}

// Stop stops the signal generation loop
func (sg *SignalGenerator) Stop() {
	sg.IsRunning = false
	log.Println("⏹️ Signal generator stopped")
}

// Helper function to parse float from string
func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}
