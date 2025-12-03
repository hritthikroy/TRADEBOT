#!/bin/bash

echo "🐛 Debug Signal Flow"
echo "===================="
echo ""

echo "This will generate a signal and show you EXACTLY what happens."
echo ""
echo "1. Start your backend in one terminal: cd backend && go run ."
echo "2. Run this script in another terminal"
echo ""
read -p "Press Enter when backend is running..."

echo ""
echo "📡 Generating signal..."
echo ""

curl -X POST http://localhost:8080/api/v1/backtest/live-signal \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "strategy": "session_trader"
  }' | jq '.'

echo ""
echo "===================="
echo ""
echo "Now check your backend logs for:"
echo ""
echo "✅ GOOD SIGNS:"
echo "  🔍 Generated signal: BUY/SELL"
echo "  🔍 Saving to Supabase: {...}"
echo "  🔍 Supabase response status: 201"
echo "  ✅ Signal saved to Supabase"
echo "  🔍 Sending to Telegram - ChatID: ..."
echo "  🔍 Telegram API response status: 200"
echo "  ✅ Message sent to Telegram successfully"
echo ""
echo "❌ BAD SIGNS:"
echo "  ℹ️  Signal is NONE (means no trading opportunity)"
echo "  ❌ Supabase error (status XXX)"
echo "  ❌ Failed to send to Telegram API"
echo "  ⚠️  Telegram bot is nil"
echo ""
echo "Copy the error messages and we can fix them!"
echo ""
