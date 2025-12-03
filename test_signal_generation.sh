#!/bin/bash

echo "🧪 Testing Signal Generation & Storage"
echo "======================================"
echo ""

# Test 1: Generate a live signal
echo "📡 Test 1: Generating live signal..."
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/backtest/live-signal \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "strategy": "session_trader"
  }')

echo "Response: $RESPONSE" | jq '.'
echo ""

# Extract signal type
SIGNAL_TYPE=$(echo "$RESPONSE" | jq -r '.signal // .Signal')
echo "Signal Type: $SIGNAL_TYPE"
echo ""

# Test 2: Check if it was saved to Supabase
echo "💾 Test 2: Checking Supabase for recent signals..."
curl -s "https://elqhqhjevaizjoghiiss.supabase.co/rest/v1/trading_signals?order=created_at.desc&limit=3" \
  -H "apikey: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImVscWhxaGpldmFpempvZ2hpaXNzIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NjQ3MDEwNTksImV4cCI6MjA4MDI3NzA1OX0.02-CgybOf7PiQSaQ-uJhojKob5Rw_2vkFdyurPNqLvA" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImVscWhxaGpldmFpempvZ2hpaXNzIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NjQ3MDEwNTksImV4cCI6MjA4MDI3NzA1OX0.02-CgybOf7PiQSaQ-uJhojKob5Rw_2vkFdyurPNqLvA" | jq '.'

echo ""

# Test 3: Check Telegram bot status
echo "📱 Test 3: Checking Telegram bot status..."
curl -s http://localhost:8080/api/v1/telegram/status | jq '.'

echo ""
echo "======================================"
echo ""
echo "📋 What to check:"
echo "1. If signal is BUY or SELL (not NONE), it should be saved"
echo "2. Check Supabase response - should show recent signals"
echo "3. Check Telegram bot status - should show running=true"
echo "4. Check your Telegram app for messages"
echo ""
