package signals

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type StoredSignal struct {
	ID              string    `json:"id,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
	Symbol          string    `json:"symbol"`
	Strategy        string    `json:"strategy"`
	SignalType      string    `json:"signal_type"`
	EntryPrice      float64   `json:"entry_price"`
	StopLoss        float64   `json:"stop_loss"`
	TakeProfit      float64   `json:"take_profit"`
	TP1             float64   `json:"tp1,omitempty"`
	TP2             float64   `json:"tp2,omitempty"`
	TP3             float64   `json:"tp3,omitempty"`
	CurrentPrice    float64   `json:"current_price"`
	RiskReward      float64   `json:"risk_reward"`
	ProfitLoss      float64   `json:"profit_loss,omitempty"`
	ProfitLossPercent float64 `json:"profit_loss_percent,omitempty"`
	Status          string    `json:"status"`
	Result          string    `json:"result,omitempty"`
	Progress        float64   `json:"progress"`
	FilterBuy       bool      `json:"filter_buy"`
	FilterSell      bool      `json:"filter_sell"`
	SignalTime      time.Time `json:"signal_time"`
	ClosedAt        *time.Time `json:"closed_at,omitempty"`
}

// SaveSignalToSupabase saves a trading signal to Supabase
func SaveSignalToSupabase(signal LiveSignalResponse, symbol, strategy string, filterBuy, filterSell bool) error {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")
	
	log.Printf("🔍 Supabase URL from env: %s", supabaseURL)
	log.Printf("🔍 Supabase Key from env: %s...", supabaseKey[:30])
	
	if supabaseURL == "" || supabaseKey == "" {
		log.Println("❌ Supabase not configured! URL or KEY is empty")
		log.Printf("   SUPABASE_URL empty: %v", supabaseURL == "")
		log.Printf("   SUPABASE_KEY empty: %v", supabaseKey == "")
		return fmt.Errorf("supabase not configured")
	}
	
	storedSignal := StoredSignal{
		Symbol:       symbol,
		Strategy:     strategy,
		SignalType:   signal.Signal,
		EntryPrice:   signal.Entry,
		StopLoss:     signal.StopLoss,
		TakeProfit:   signal.TakeProfit,
		TP1:          signal.TP1,
		TP2:          signal.TP2,
		TP3:          signal.TP3,
		CurrentPrice: signal.CurrentPrice,
		RiskReward:   signal.RiskReward,
		Status:       "ACTIVE",
		Progress:     0,
		FilterBuy:    filterBuy,
		FilterSell:   filterSell,
		SignalTime:   time.Now(),
	}
	
	jsonData, err := json.Marshal(storedSignal)
	if err != nil {
		return fmt.Errorf("failed to marshal signal: %v", err)
	}
	
	log.Printf("🔍 Saving to Supabase: %s", string(jsonData))
	
	url := fmt.Sprintf("%s/rest/v1/trading_signals", supabaseURL)
	log.Printf("🔍 Supabase URL: %s", url)
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", supabaseKey)
	req.Header.Set("Authorization", "Bearer "+supabaseKey)
	req.Header.Set("Prefer", "return=representation")
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Failed to send request to Supabase: %v", err)
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()
	
	// Read response body for debugging
	var responseBody bytes.Buffer
	responseBody.ReadFrom(resp.Body)
	
	log.Printf("🔍 Supabase response status: %d", resp.StatusCode)
	log.Printf("🔍 Supabase response body: %s", responseBody.String())
	
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		log.Printf("❌ Supabase error (status %d): %s", resp.StatusCode, responseBody.String())
		return fmt.Errorf("supabase returned status %d: %s", resp.StatusCode, responseBody.String())
	}
	
	log.Printf("✅ Signal saved to Supabase: %s %s @ $%.2f", signal.Signal, symbol, signal.Entry)
	return nil
}

// GetRecentSignals retrieves recent signals from Supabase
func GetRecentSignals(limit int) ([]StoredSignal, error) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")
	
	if supabaseURL == "" || supabaseKey == "" {
		return nil, fmt.Errorf("supabase not configured")
	}
	
	url := fmt.Sprintf("%s/rest/v1/trading_signals?order=created_at.desc&limit=%d", supabaseURL, limit)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("apikey", supabaseKey)
	req.Header.Set("Authorization", "Bearer "+supabaseKey)
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var signals []StoredSignal
	if err := json.NewDecoder(resp.Body).Decode(&signals); err != nil {
		return nil, err
	}
	
	return signals, nil
}

// UpdateSignalStatus updates the status of a signal in Supabase
func UpdateSignalStatus(signalID string, status string, currentPrice float64, profitLoss float64, profitLossPercent float64) error {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")
	
	if supabaseURL == "" || supabaseKey == "" {
		return fmt.Errorf("supabase not configured")
	}
	
	updateData := map[string]interface{}{
		"status":              status,
		"current_price":       currentPrice,
		"profit_loss":         profitLoss,
		"profit_loss_percent": profitLossPercent,
	}
	
	if status == "HIT_TP" || status == "HIT_SL" || status == "CLOSED" {
		now := time.Now()
		updateData["closed_at"] = now
	}
	
	jsonData, err := json.Marshal(updateData)
	if err != nil {
		return err
	}
	
	url := fmt.Sprintf("%s/rest/v1/trading_signals?id=eq.%s", supabaseURL, signalID)
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", supabaseKey)
	req.Header.Set("Authorization", "Bearer "+supabaseKey)
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("supabase returned status %d", resp.StatusCode)
	}
	
	return nil
}

// GetSignalPerformance retrieves performance metrics from Supabase
func GetSignalPerformance() ([]map[string]interface{}, error) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")
	
	if supabaseURL == "" || supabaseKey == "" {
		return nil, fmt.Errorf("supabase not configured")
	}
	
	url := fmt.Sprintf("%s/rest/v1/signal_performance", supabaseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("apikey", supabaseKey)
	req.Header.Set("Authorization", "Bearer "+supabaseKey)
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var performance []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&performance); err != nil {
		return nil, err
	}
	
	return performance, nil
}
