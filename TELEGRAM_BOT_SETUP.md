# 📱 Telegram Bot Setup Guide

## ✅ What Was Added

A 24/7 Telegram bot that automatically sends trading signals to your Telegram channel every 5 minutes.

## 🚀 Features

- **24/7 Operation**: Runs continuously in the background
- **Auto Signal Detection**: Checks for signals every 5 minutes
- **Filtered Signals**: Respects Buy/Sell filters
- **Rich Formatting**: Beautiful formatted messages with emojis
- **Real-time Updates**: Instant signal delivery to Telegram
- **Easy Control**: Start/Stop from the web interface

## 📋 Setup Instructions

### Step 1: Create a Telegram Bot

1. Open Telegram and search for `@BotFather`
2. Send `/newbot` command
3. Follow the instructions to create your bot
4. Copy the **Bot Token** (looks like: `123456789:ABCdefGHIjklMNOpqrsTUVwxyz`)

### Step 2: Get Your Chat ID

**Option A: Using Your Bot**
1. Start a chat with your new bot
2. Send any message to it
3. Visit: `https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates`
4. Look for `"chat":{"id":` - that's your Chat ID

**Option B: Using a Channel**
1. Create a Telegram channel
2. Add your bot as an administrator
3. Send a message to the channel
4. Visit: `https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates`
5. Look for `"chat":{"id":` - that's your Channel ID (starts with `-100`)

### Step 3: Configure Environment Variables

Edit `backend/.env` file:

```env
# Telegram Bot Configuration
TELEGRAM_BOT_TOKEN=123456789:ABCdefGHIjklMNOpqrsTUVwxyz
TELEGRAM_CHAT_ID=your_chat_id_here
```

**Important Notes:**
- For personal chat: Use positive number (e.g., `123456789`)
- For channel: Use negative number (e.g., `-1001234567890`)
- For group: Use negative number (e.g., `-123456789`)

### Step 4: Restart Backend

```bash
cd backend
go run .
```

You should see: `✅ Telegram bot initialized`

## 🎯 How to Use

### Start the Bot

1. Go to **Live Signals** tab
2. Configure your settings:
   - Symbol (e.g., BTCUSDT)
   - Strategy (e.g., Session Trader)
   - Trade Type Filter (Buy/Sell)
3. Click **"📱 Start Telegram Bot (24/7)"**
4. Bot will start monitoring and sending signals

### Stop the Bot

1. Click **"🛑 Stop Telegram Bot"**
2. Bot will stop sending signals

## 📨 Signal Format

Signals are sent in this format:

```
🟢 BUY SIGNAL

📊 Symbol: BTCUSDT
🎯 Strategy: session_trader
💰 Current Price: $43,250.50

📍 Entry: $43,250.50
🛑 Stop Loss: $42,800.00
🎯 Take Profit: $44,375.75
📊 Risk/Reward: 2.5:1

⏰ Time: 2024-12-02 15:30:45 UTC

Automated signal from Trading Bot
```

## ⚙️ Configuration

### Signal Frequency
- Default: Every 5 minutes
- Can be changed in `telegram_bot.go` (line 48)

```go
ticker := time.NewTicker(5 * time.Minute)
```

### Filters
- **Buy Only**: Only sends BUY signals
- **Sell Only**: Only sends SELL signals
- **Both**: Sends all signals (default)

## 🔧 API Endpoints

### Start Bot
```
POST /api/v1/telegram/start
Content-Type: application/json

{
  "symbol": "BTCUSDT",
  "strategy": "session_trader",
  "filterBuy": true,
  "filterSell": true
}
```

### Stop Bot
```
POST /api/v1/telegram/stop
```

### Get Status
```
GET /api/v1/telegram/status
```

## 🛡️ Security Best Practices

1. **Keep Token Secret**: Never share your bot token
2. **Use Private Channel**: Don't use public channels for signals
3. **Restrict Bot Access**: Only add trusted admins to your channel
4. **Monitor Activity**: Check bot logs regularly
5. **Rotate Tokens**: Change bot token if compromised

## 🐛 Troubleshooting

### Bot Not Sending Messages

**Check 1: Token and Chat ID**
```bash
# Test your configuration
curl "https://api.telegram.org/bot<YOUR_TOKEN>/sendMessage?chat_id=<YOUR_CHAT_ID>&text=Test"
```

**Check 2: Bot Permissions**
- For channels: Bot must be admin
- For groups: Bot must be member
- For personal: Just start the bot

**Check 3: Backend Logs**
```bash
# Check for errors
tail -f backend/server.log
```

### Bot Shows "Not Configured"

- Verify `.env` file has correct values
- Restart backend server
- Check environment variables are loaded

### Signals Not Appearing

- Check if bot is actually running (status should show "Running")
- Verify strategy is generating signals
- Check filter settings (Buy/Sell)
- Look at backend logs for errors

## 📊 Monitoring

### Check Bot Status
- Green indicator: Bot is running
- Status card shows: "🟢 Running"
- Backend logs show: "🤖 Telegram signal bot started"

### Signal Delivery
- Each signal sent shows in logs: "📤 Sent BUY signal to Telegram"
- Check your Telegram channel for messages
- Signals appear every 5 minutes (if conditions met)

## 🔄 Advanced Usage

### Multiple Strategies
Run multiple bots for different strategies:
1. Create multiple bot tokens
2. Use different chat IDs
3. Start each bot with different strategy

### Custom Intervals
Modify check frequency in `telegram_bot.go`:
```go
// Check every 1 minute
ticker := time.NewTicker(1 * time.Minute)

// Check every 15 minutes
ticker := time.NewTicker(15 * time.Minute)
```

### Custom Message Format
Edit `SendSignal()` function in `telegram_bot.go` to customize message format.

## 📱 Telegram Features

### Commands (Future Enhancement)
- `/start` - Start receiving signals
- `/stop` - Stop receiving signals
- `/status` - Check bot status
- `/help` - Show help message

### Buttons (Future Enhancement)
- "✅ Take Trade" button
- "❌ Skip Trade" button
- "📊 View Chart" button

## 🎉 Benefits

### For Traders
- **Never Miss a Signal**: 24/7 monitoring
- **Mobile Alerts**: Get signals on your phone
- **Instant Notifications**: Real-time delivery
- **Hands-Free**: Fully automated

### For Strategy Testing
- **Live Validation**: See signals in real-time
- **Performance Tracking**: Monitor signal quality
- **Easy Sharing**: Share signals with team
- **Historical Record**: All signals saved in Telegram

## 🔮 Future Enhancements

Planned features:
- Signal performance tracking
- Trade execution integration
- Multiple symbol monitoring
- Custom alert conditions
- Signal analytics dashboard
- Backtesting integration
- Portfolio management

---

**Status**: ✅ FULLY FUNCTIONAL
**Last Updated**: December 2, 2024
**Version**: 1.0.0

## 🆘 Support

If you encounter issues:
1. Check this guide first
2. Review backend logs
3. Test Telegram API manually
4. Verify environment variables
5. Restart backend server

**Note**: Make sure your Telegram bot token and chat ID are correctly configured in the `.env` file before starting the bot!
