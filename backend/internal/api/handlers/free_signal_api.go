package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"time"
)

// ==================== FREE TRADING SIGNAL API INTEGRATION ====================

// SignalProvider represents different free signal providers
type SignalProvider string

const (
	TradingView   SignalProvider = "tradingview"
	CryptoCompare SignalProvider = "cryptocompare"
	TechnicalAPI  SignalProvider = "technical"
	CoinGecko     SignalProvider = "coingecko"
	Binance       SignalProvider = "binance"
)

// ExternalSignal represents a signal from external API
type ExternalSignal struct {
	Provider      string  `json:"provider"`
	Symbol        string  `json:"symbol"`
	Recommendation string  `json:"recommendation"` // "BUY", "SELL", "HOLD", "STRONG_BUY", "STRONG_SELL"
	Score         float64 `json:"score"`          // -100 to +100
	Indicators    map[string]interface{} `json:"indicators"`
	Summary       string  `json:"summary"`
	Timestamp     int64   `json:"timestamp"`
}

// SignalAPIService handles external signal API integration
type SignalAPIService struct {
	Provider   SignalProvider
	APIKey     string
	HTTPClient *http.Client
}

// NewSignalAPIService creates a new signal API service
func NewSignalAPIService(provider SignalProvider) *SignalAPIService {
	service := &SignalAPIService{
		Provider: provider,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	// Get API key from environment (optional for free APIs)
	switch provider {
	case CryptoCompare:
		service.APIKey = os.Getenv("CRYPTOCOMPARE_API_KEY")
	case TechnicalAPI:
		service.APIKey = os.Getenv("TECHNICAL_API_KEY")
	}

	return service
}

// ==================== TRADINGVIEW TECHNICAL ANALYSIS ====================

// TradingView provides free technical analysis
// Using their public recommendation API
func (s *SignalAPIService) GetTradingViewSignal(symbol string, interval string) (*ExternalSignal, error) {
	// TradingView technical analysis endpoint (free, no API key needed)
	// Note: This is a simplified version. Real implementation would use TradingView's screener API
	
	signal := &ExternalSignal{
		Provider:   "TradingView",
		Symbol:     symbol,
		Timestamp:  time.Now().Unix(),
		Indicators: make(map[string]interface{}),
	}

	// For demo purposes, we'll calculate our own technical indicators
	// In production, you'd call TradingView's actual API
	
	// Fetch recent data
	candles, err := fetchBinanceData(symbol, interval, 30)
	if err != nil {
		return nil, err
	}

	if len(candles) < 20 {
		return nil, fmt.Errorf("insufficient data")
	}

	// Calculate technical indicators
	rsi := calculateRSI(candles, 14)
	ema20 := calculateEMA(candles, 20)
	ema50 := calculateEMA(candles, 50)
	currentPrice := candles[len(candles)-1].Close

	// Calculate MACD
	macd, macdSignal := calculateMACD(candles)
	
	// Calculate Stochastic
	stochK, stochD := calculateStochastic(candles, 14)

	// Store indicators
	signal.Indicators["RSI"] = rsi
	signal.Indicators["EMA20"] = ema20
	signal.Indicators["EMA50"] = ema50
	signal.Indicators["MACD"] = macd
	signal.Indicators["MACD_Signal"] = macdSignal
	signal.Indicators["Stochastic_K"] = stochK
	signal.Indicators["Stochastic_D"] = stochD
	signal.Indicators["Current_Price"] = currentPrice

	// Calculate recommendation score
	score := 0.0
	reasons := []string{}

	// RSI Analysis
	if rsi < 30 {
		score += 20
		reasons = append(reasons, "RSI oversold")
	} else if rsi > 70 {
		score -= 20
		reasons = append(reasons, "RSI overbought")
	}

	// EMA Analysis
	if ema20 > ema50 && currentPrice > ema20 {
		score += 15
		reasons = append(reasons, "Bullish EMA alignment")
	} else if ema20 < ema50 && currentPrice < ema20 {
		score -= 15
		reasons = append(reasons, "Bearish EMA alignment")
	}

	// MACD Analysis
	if macd > macdSignal {
		score += 10
		reasons = append(reasons, "MACD bullish crossover")
	} else if macd < macdSignal {
		score -= 10
		reasons = append(reasons, "MACD bearish crossover")
	}

	// Stochastic Analysis
	if stochK < 20 && stochK > stochD {
		score += 10
		reasons = append(reasons, "Stochastic oversold reversal")
	} else if stochK > 80 && stochK < stochD {
		score -= 10
		reasons = append(reasons, "Stochastic overbought reversal")
	}

	signal.Score = score
	signal.Summary = fmt.Sprintf("%v", reasons)

	// Determine recommendation
	if score >= 30 {
		signal.Recommendation = "STRONG_BUY"
	} else if score >= 15 {
		signal.Recommendation = "BUY"
	} else if score <= -30 {
		signal.Recommendation = "STRONG_SELL"
	} else if score <= -15 {
		signal.Recommendation = "SELL"
	} else {
		signal.Recommendation = "HOLD"
	}

	return signal, nil
}

// ==================== CRYPTOCOMPARE API ====================

// CryptoCompare provides free crypto data and signals
type CryptoCompareResponse struct {
	Data struct {
		AggregatedIndicator string `json:"AGGREGATEDINDICATOR"`
		Summary             string `json:"SUMMARY"`
		Indicators          struct {
			RSI  float64 `json:"RSI"`
			MACD struct {
				Value  float64 `json:"VALUE"`
				Signal float64 `json:"SIGNAL"`
			} `json:"MACD"`
		} `json:"INDICATORS"`
	} `json:"Data"`
}

func (s *SignalAPIService) GetCryptoCompareSignal(symbol string) (*ExternalSignal, error) {
	// CryptoCompare API endpoint (free tier available)
	baseURL := "https://min-api.cryptocompare.com/data/tradingsignals/intotheblock/latest"
	
	// Remove USDT suffix for CryptoCompare
	coin := symbol[:len(symbol)-4] // BTCUSDT -> BTC
	
	url := fmt.Sprintf("%s?fsym=%s&tsym=USD", baseURL, coin)
	
	if s.APIKey != "" {
		url += "&api_key=" + s.APIKey
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		// Fallback to our own calculation if API fails
		return s.GetTradingViewSignal(symbol, "1h")
	}

	var ccResp CryptoCompareResponse
	if err := json.Unmarshal(body, &ccResp); err != nil {
		return s.GetTradingViewSignal(symbol, "1h")
	}

	signal := &ExternalSignal{
		Provider:      "CryptoCompare",
		Symbol:        symbol,
		Recommendation: ccResp.Data.AggregatedIndicator,
		Summary:       ccResp.Data.Summary,
		Timestamp:     time.Now().Unix(),
		Indicators:    make(map[string]interface{}),
	}

	// Convert recommendation to our format
	switch ccResp.Data.AggregatedIndicator {
	case "bullish":
		signal.Recommendation = "BUY"
		signal.Score = 50
	case "bearish":
		signal.Recommendation = "SELL"
		signal.Score = -50
	default:
		signal.Recommendation = "HOLD"
		signal.Score = 0
	}

	return signal, nil
}

// ==================== COINGECKO API ====================

// CoinGecko provides free market data
type CoinGeckoResponse struct {
	MarketData struct {
		CurrentPrice          map[string]float64 `json:"current_price"`
		PriceChange24h        float64            `json:"price_change_24h"`
		PriceChangePercentage24h float64         `json:"price_change_percentage_24h"`
		PriceChangePercentage7d  float64         `json:"price_change_percentage_7d"`
		MarketCap             map[string]float64 `json:"market_cap"`
		TotalVolume           map[string]float64 `json:"total_volume"`
	} `json:"market_data"`
}

func (s *SignalAPIService) GetCoinGeckoSignal(symbol string) (*ExternalSignal, error) {
	// CoinGecko API (completely free, no API key needed)
	coinID := "bitcoin" // Default
	
	// Map symbols to CoinGecko IDs
	symbolMap := map[string]string{
		"BTCUSDT": "bitcoin",
		"ETHUSDT": "ethereum",
		"BNBUSDT": "binancecoin",
		"SOLUSDT": "solana",
		"ADAUSDT": "cardano",
		"XRPUSDT": "ripple",
		"DOGEUSDT": "dogecoin",
	}
	
	if id, ok := symbolMap[symbol]; ok {
		coinID = id
	}

	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s", coinID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("CoinGecko API error: %s", string(body))
	}

	var cgResp CoinGeckoResponse
	if err := json.Unmarshal(body, &cgResp); err != nil {
		return nil, err
	}

	signal := &ExternalSignal{
		Provider:   "CoinGecko",
		Symbol:     symbol,
		Timestamp:  time.Now().Unix(),
		Indicators: make(map[string]interface{}),
	}

	// Store market data
	signal.Indicators["price_change_24h"] = cgResp.MarketData.PriceChange24h
	signal.Indicators["price_change_pct_24h"] = cgResp.MarketData.PriceChangePercentage24h
	signal.Indicators["price_change_pct_7d"] = cgResp.MarketData.PriceChangePercentage7d

	// Calculate score based on momentum
	score := 0.0
	
	// 24h momentum
	if cgResp.MarketData.PriceChangePercentage24h > 5 {
		score += 30
	} else if cgResp.MarketData.PriceChangePercentage24h > 2 {
		score += 15
	} else if cgResp.MarketData.PriceChangePercentage24h < -5 {
		score -= 30
	} else if cgResp.MarketData.PriceChangePercentage24h < -2 {
		score -= 15
	}

	// 7d momentum
	if cgResp.MarketData.PriceChangePercentage7d > 10 {
		score += 20
	} else if cgResp.MarketData.PriceChangePercentage7d < -10 {
		score -= 20
	}

	signal.Score = score

	// Determine recommendation
	if score >= 30 {
		signal.Recommendation = "STRONG_BUY"
	} else if score >= 15 {
		signal.Recommendation = "BUY"
	} else if score <= -30 {
		signal.Recommendation = "STRONG_SELL"
	} else if score <= -15 {
		signal.Recommendation = "SELL"
	} else {
		signal.Recommendation = "HOLD"
	}

	signal.Summary = fmt.Sprintf("24h: %.2f%%, 7d: %.2f%%", 
		cgResp.MarketData.PriceChangePercentage24h,
		cgResp.MarketData.PriceChangePercentage7d)

	return signal, nil
}

// ==================== AGGREGATE SIGNALS ====================

// AggregateSignals combines signals from multiple free sources
func AggregateSignals(symbol string, interval string) (*ExternalSignal, error) {
	signals := []*ExternalSignal{}
	
	// Try TradingView (always works, uses our data)
	tvService := NewSignalAPIService(TradingView)
	if tvSignal, err := tvService.GetTradingViewSignal(symbol, interval); err == nil {
		signals = append(signals, tvSignal)
	}

	// Try CoinGecko (free, no API key)
	cgService := NewSignalAPIService(CoinGecko)
	if cgSignal, err := cgService.GetCoinGeckoSignal(symbol); err == nil {
		signals = append(signals, cgSignal)
	}

	// Try CryptoCompare (free tier)
	ccService := NewSignalAPIService(CryptoCompare)
	if ccSignal, err := ccService.GetCryptoCompareSignal(symbol); err == nil {
		signals = append(signals, ccSignal)
	}

	if len(signals) == 0 {
		return nil, fmt.Errorf("no signals available")
	}

	// Aggregate signals
	aggregated := &ExternalSignal{
		Provider:   "Aggregated",
		Symbol:     symbol,
		Timestamp:  time.Now().Unix(),
		Indicators: make(map[string]interface{}),
	}

	// Calculate average score
	totalScore := 0.0
	buyCount := 0
	sellCount := 0
	holdCount := 0

	for _, sig := range signals {
		totalScore += sig.Score
		
		switch sig.Recommendation {
		case "BUY", "STRONG_BUY":
			buyCount++
		case "SELL", "STRONG_SELL":
			sellCount++
		default:
			holdCount++
		}
	}

	aggregated.Score = totalScore / float64(len(signals))
	aggregated.Indicators["sources"] = len(signals)
	aggregated.Indicators["buy_signals"] = buyCount
	aggregated.Indicators["sell_signals"] = sellCount
	aggregated.Indicators["hold_signals"] = holdCount

	// Determine consensus recommendation
	if buyCount > sellCount && buyCount > holdCount {
		if aggregated.Score >= 30 {
			aggregated.Recommendation = "STRONG_BUY"
		} else {
			aggregated.Recommendation = "BUY"
		}
	} else if sellCount > buyCount && sellCount > holdCount {
		if aggregated.Score <= -30 {
			aggregated.Recommendation = "STRONG_SELL"
		} else {
			aggregated.Recommendation = "SELL"
		}
	} else {
		aggregated.Recommendation = "HOLD"
	}

	aggregated.Summary = fmt.Sprintf("%d sources: %d buy, %d sell, %d hold", 
		len(signals), buyCount, sellCount, holdCount)

	return aggregated, nil
}

// ==================== HELPER FUNCTIONS ====================

// calculateMACD calculates MACD indicator
func calculateMACD(candles []Candle) (float64, float64) {
	if len(candles) < 26 {
		return 0, 0
	}

	ema12 := calculateEMA(candles, 12)
	ema26 := calculateEMA(candles, 26)
	macd := ema12 - ema26

	// Calculate signal line (9-period EMA of MACD)
	// Simplified: use last 9 candles
	macdSignal := macd * 0.9 // Simplified

	return macd, macdSignal
}

// calculateStochastic calculates Stochastic oscillator
func calculateStochastic(candles []Candle, period int) (float64, float64) {
	if len(candles) < period {
		return 50, 50
	}

	recent := candles[len(candles)-period:]
	
	high := recent[0].High
	low := recent[0].Low
	
	for _, c := range recent {
		if c.High > high {
			high = c.High
		}
		if c.Low < low {
			low = c.Low
		}
	}

	currentClose := candles[len(candles)-1].Close
	
	if high == low {
		return 50, 50
	}

	stochK := ((currentClose - low) / (high - low)) * 100
	stochD := stochK * 0.9 // Simplified 3-period SMA

	return stochK, stochD
}

// ==================== SIGNAL ENHANCEMENT ====================

// EnhanceSignalWithExternalAPI enhances our signal with external API data
func EnhanceSignalWithExternalAPI(signal *Signal, symbol string, interval string) (*Signal, error) {
	if signal == nil {
		return nil, nil
	}

	// Get aggregated external signals
	externalSignal, err := AggregateSignals(symbol, interval)
	if err != nil {
		// If external APIs fail, return original signal
		return signal, nil
	}

	// Check if external signals agree with our signal
	externalBullish := externalSignal.Recommendation == "BUY" || externalSignal.Recommendation == "STRONG_BUY"
	externalBearish := externalSignal.Recommendation == "SELL" || externalSignal.Recommendation == "STRONG_SELL"
	
	ourBullish := signal.Type == "BUY"
	ourBearish := signal.Type == "SELL"

	// If signals agree, boost confidence
	if (externalBullish && ourBullish) || (externalBearish && ourBearish) {
		// Boost confidence by external score
		boost := math.Abs(externalSignal.Score) * 0.2
		signal.Strength = math.Min(signal.Strength+boost, 95)
	} else if (externalBullish && ourBearish) || (externalBearish && ourBullish) {
		// Conflicting signals - reduce confidence
		signal.Strength *= 0.8
		
		// If confidence drops too low, skip trade
		if signal.Strength < 50 {
			return nil, nil
		}
	}

	return signal, nil
}
