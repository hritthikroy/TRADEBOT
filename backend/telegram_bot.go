package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type TelegramBot struct {
	Token   string
	ChatID  string
	Running bool
}

type TelegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

var telegramBot *TelegramBot

// InitTelegramBot initializes the Telegram bot
func InitTelegramBot() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")
	
	if token == "" || chatID == "" {
		log.Println("⚠️  Telegram bot not configured (missing TELEGRAM_BOT_TOKEN or TELEGRAM_CHAT_ID)")
		return
	}
	
	telegramBot = &TelegramBot{
		Token:   token,
		ChatID:  chatID,
		Running: false,
	}
	
	log.Println("✅ Telegram bot initialized")
}

// StartTelegramSignalBot starts the 24/7 signal monitoring
func StartTelegramSignalBot(symbol, strategy string, filterBuy, filterSell bool) error {
	if telegramBot == nil {
		return fmt.Errorf("telegram bot not initialized")
	}
	
	if telegramBot.Running {
		return fmt.Errorf("telegram bot already running")
	}
	
	telegramBot.Running = true
	
	// Send startup message
	startMsg := fmt.Sprintf("🤖 *Trading Signal Bot Started*\n\n"+
		"📊 Symbol: `%s`\n"+
		"🎯 Strategy: `%s`\n"+
		"🟢 Buy Signals: %v\n"+
		"🔴 Sell Signals: %v\n"+
		"⏰ Checking every 15 seconds\n\n"+
		"_Bot will run 24/7 until stopped_\n"+
		"_All signals saved to Supabase automatically_",
		symbol, strategy, filterBuy, filterSell)
	
	telegramBot.SendMessage(startMsg)
	
	// Start background goroutine
	go func() {
		ticker := time.NewTicker(15 * time.Second) // Check every 15 seconds - balanced speed and safety
		defer ticker.Stop()
		
		log.Printf("🤖 Telegram signal bot started for %s with %s strategy (checking every 15 seconds)", symbol, strategy)
		
		lastSignalTime := time.Now().Add(-1 * time.Minute) // Allow first signal immediately
		lastSignalType := ""
		
		for telegramBot.Running {
			select {
			case <-ticker.C:
				log.Printf("🔄 Telegram bot checking market for %s...", symbol)
				
				// Generate signal
				interval := getStrategyInterval(strategy)
				candles, err := fetchBinanceData(symbol, interval, 7)
				if err != nil {
					log.Printf("❌ Error fetching candles for Telegram: %v", err)
					continue
				}
				
				if len(candles) < 50 {
					log.Printf("⚠️  Not enough candles (%d), skipping", len(candles))
					continue
				}
				
				signal := generateLiveSignal(candles, strategy)
				log.Printf("🔍 Telegram bot generated signal: %s", signal.Signal)
				
				// Check if signal matches filter
				if signal.Signal == "NONE" {
					log.Printf("ℹ️  No signal (NONE), waiting for next check...")
					continue
				}
				
				// Get current filter settings from database (not startup parameters)
				currentFilterBuy, currentFilterSell := GetCurrentFilterSettings()
				log.Printf("🔍 Current filter settings: filterBuy=%v, filterSell=%v", currentFilterBuy, currentFilterSell)
				
				if (signal.Signal == "BUY" && !currentFilterBuy) || (signal.Signal == "SELL" && !currentFilterSell) {
					log.Printf("⏭️  Signal %s filtered out (filterBuy=%v, filterSell=%v)", signal.Signal, currentFilterBuy, currentFilterSell)
					continue
				}
				
				// Rate limiting: Only send if signal changed (not just time passed)
				// This prevents saving the same signal multiple times
				if signal.Signal == lastSignalType {
					timeSinceLastSignal := time.Since(lastSignalTime)
					log.Printf("⏭️  Skipping duplicate %s signal (same as last, %v ago)", signal.Signal, timeSinceLastSignal)
					continue // Skip duplicate signals completely
				}
				
				log.Printf("✅ New signal detected: %s for %s", signal.Signal, symbol)
				
				// Save signal to Supabase FIRST (use current filter settings)
				log.Printf("💾 Attempting to save signal to Supabase...")
				err = SaveSignalToSupabase(signal, symbol, strategy, currentFilterBuy, currentFilterSell)
				if err != nil {
					log.Printf("❌ FAILED to save signal to Supabase: %v", err)
					log.Printf("   Signal: %s %s @ $%.2f", signal.Signal, symbol, signal.Entry)
				} else {
					log.Printf("✅ Signal successfully saved to Supabase")
				}
				
				// Send signal to Telegram (even if Supabase fails)
				log.Printf("📱 Sending signal to Telegram...")
				telegramBot.SendSignal(signal, symbol, strategy)
				
				// Update last signal tracking
				lastSignalTime = time.Now()
				lastSignalType = signal.Signal
				log.Printf("✅ Signal processing complete for %s", signal.Signal)
			}
		}
		
		log.Println("🤖 Telegram signal bot stopped")
	}()
	
	return nil
}

// StopTelegramSignalBot stops the signal monitoring
func StopTelegramSignalBot() {
	if telegramBot != nil && telegramBot.Running {
		telegramBot.Running = false
		telegramBot.SendMessage("🛑 *Trading Signal Bot Stopped*\n\n_Signal monitoring has been disabled_")
		log.Println("🤖 Telegram signal bot stopped")
	}
}

// SendSignal sends a trading signal to Telegram
func (bot *TelegramBot) SendSignal(signal LiveSignalResponse, symbol, strategy string) {
	var emoji string
	var signalType string
	
	if signal.Signal == "BUY" {
		emoji = "🟢"
		signalType = "BUY SIGNAL"
	} else {
		emoji = "🔴"
		signalType = "SELL SIGNAL"
	}
	
	message := fmt.Sprintf(
		"%s *%s*\n\n"+
			"📊 *Symbol:* `%s`\n"+
			"🎯 *Strategy:* `%s`\n"+
			"💰 *Current Price:* `$%.2f`\n\n"+
			"📍 *Entry:* `$%.2f`\n"+
			"🛑 *Stop Loss:* `$%.2f`\n\n"+
			"🎯 *Take Profit Levels:*\n"+
			"   TP1 (33%%): `$%.2f`\n"+
			"   TP2 (33%%): `$%.2f`\n"+
			"   TP3 (34%%): `$%.2f`\n\n"+
			"📊 *Risk/Reward:* `%.2f:1`\n"+
			"⏰ *Time:* `%s`\n\n"+
			"_Automated signal from Trading Bot_",
		emoji, signalType,
		symbol,
		strategy,
		signal.CurrentPrice,
		signal.Entry,
		signal.StopLoss,
		signal.TP1,
		signal.TP2,
		signal.TP3,
		signal.RiskReward,
		time.Now().Format("2006-01-02 15:04:05 MST"),
	)
	
	bot.SendMessage(message)
	log.Printf("📤 Sent %s signal to Telegram for %s", signal.Signal, symbol)
}

// SendMessage sends a message to Telegram
func (bot *TelegramBot) SendMessage(text string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", bot.Token)
	
	log.Printf("🔍 Sending to Telegram - ChatID: %s", bot.ChatID)
	
	msg := TelegramMessage{
		ChatID:    bot.ChatID,
		Text:      text,
		ParseMode: "Markdown",
	}
	
	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("❌ Failed to marshal Telegram message: %v", err)
		return err
	}
	
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("❌ Failed to send to Telegram API: %v", err)
		return err
	}
	defer resp.Body.Close()
	
	log.Printf("🔍 Telegram API response status: %d", resp.StatusCode)
	
	if resp.StatusCode != http.StatusOK {
		// Read error response
		var responseBody bytes.Buffer
		responseBody.ReadFrom(resp.Body)
		log.Printf("❌ Telegram API error (status %d): %s", resp.StatusCode, responseBody.String())
		return fmt.Errorf("telegram API returned status %d: %s", resp.StatusCode, responseBody.String())
	}
	
	log.Printf("✅ Message sent to Telegram successfully")
	return nil
}

// GetTelegramBotStatus returns the current bot status
func GetTelegramBotStatus() map[string]interface{} {
	if telegramBot == nil {
		return map[string]interface{}{
			"configured": false,
			"running":    false,
			"message":    "Telegram bot not configured",
		}
	}
	
	return map[string]interface{}{
		"configured": true,
		"running":    telegramBot.Running,
		"message":    "Telegram bot ready",
	}
}
