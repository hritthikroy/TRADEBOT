# 🔍 Troubleshooting: No Signals Being Saved/Sent

## 🎯 Problem

- Supabase not storing signals
- Telegram not sending messages

## 🔧 Enhanced Logging Added

I've added detailed logging to help diagnose the issue. Now you'll see:

### Live Signal Handler Logs:
```
🔍 Generated signal: BUY for BTCUSDT using session_trader strategy
✅ Signal saved to Supabase: BUY BTCUSDT @ $50000.00
📤 Sent BUY signal to Telegram for BTCUSDT
```

OR if there's an issue:
```
⚠️  Failed to save signal to Supabase: [error details]
⚠️  Telegram bot is nil, cannot send signal
⚠️  Telegram bot token is empty, cannot send signal
ℹ️  Signal is NONE, not sending to Telegram
```

### Telegram Bot Logs:
```
🔄 Telegram bot checking market for BTCUSDT...
🔍 Telegram bot generated signal: BUY
✅ Signal saved to Supabase: BUY BTCUSDT @ $50000.00
📤 Sent BUY signal to Telegram for BTCUSDT
```

OR:
```
ℹ️  No signal (NONE), waiting for next check...
⏭️  Signal BUY filtered out (filterBuy=false, filterSell=true)
❌ Error fetching candles for Telegram: [error]
⚠️  Not enough candles (45), skipping
```

## 📋 Step-by-Step Diagnosis

### Step 1: Restart Backend with Logging

```bash
cd backend
go run .
```

Watch for these startup messages:
```
✅ Telegram bot initialized
✅ Telegram bot auto-started for BTCUSDT with session_trader strategy
🤖 Telegram signal bot started for BTCUSDT with session_trader strategy (checking every 15 seconds)
```

### Step 2: Wait 15 Seconds

The bot checks every 15 seconds. You should see:
```
🔄 Telegram bot checking market for BTCUSDT...
🔍 Telegram bot generated signal: NONE
ℹ️  No signal (NONE), waiting for next check...
```

This is NORMAL! Most of the time, there's no signal.

### Step 3: Manually Generate a Signal

```bash
./test_signal_generation.sh
```

This will:
1. Generate a signal via API
2. Check if it was saved to Supabase
3. Check Telegram bot status

### Step 4: Check Backend Logs

Look for:
- ✅ `Signal saved to Supabase` - Good!
- ✅ `Sent signal to Telegram` - Good!
- ⚠️ `Failed to save` - Problem with Supabase
- ⚠️ `Telegram bot is nil` - Bot not initialized
- ℹ️ `Signal is NONE` - No trading opportunity

## 🐛 Common Issues & Solutions

### Issue 1: Signal is Always NONE

**Symptom:**
```
ℹ️  No signal (NONE), waiting for next check...
```

**Cause:** Market conditions don't match strategy criteria

**Solution:** This is NORMAL! Strategies are selective. Try:
1. Wait longer (signals can take minutes/hours)
2. Try different strategy: `breakout_master`, `liquidity_hunter`
3. Try different symbol: `ETHUSDT`, `SOLUSDT`
4. Lower timeframe strategies generate more signals

**Test with multiple strategies:**
```bash
# Test session_trader
curl -X POST http://localhost:8080/api/v1/backtest/live-signal \
  -H "Content-Type: application/json" \
  -d '{"symbol": "BTCUSDT", "strategy": "session_trader"}'

# Test breakout_master
curl -X POST http://localhost:8080/api/v1/backtest/live-signal \
  -H "Content-Type: application/json" \
  -d '{"symbol": "BTCUSDT", "strategy": "breakout_master"}'

# Test scalper_pro (more signals)
curl -X POST http://localhost:8080/api/v1/backtest/live-signal \
  -H "Content-Type: application/json" \
  -d '{"symbol": "BTCUSDT", "strategy": "scalper_pro"}'
```

### Issue 2: Supabase Not Saving

**Symptom:**
```
⚠️  Failed to save signal to Supabase: [error]
```

**Check:**
1. Is Supabase URL correct in `.env`?
   ```bash
   grep SUPABASE backend/.env
   ```
   Should show: `https://elqhqhjevaizjoghiiss.supabase.co`

2. Test connection:
   ```bash
   ./test_supabase_working.sh
   ```

3. Check Supabase dashboard:
   - Go to: https://supabase.com/dashboard/project/elqhqhjevaizjoghiiss/editor
   - Check if `trading_signals` table exists

**Fix:**
```bash
# If table doesn't exist, run the SQL
# Copy contents of supabase-setup.sql
# Paste in Supabase SQL Editor and run
```

### Issue 3: Telegram Not Sending

**Symptom:**
```
⚠️  Telegram bot is nil, cannot send signal
⚠️  Telegram bot token is empty, cannot send signal
```

**Check:**
1. Is Telegram configured in `.env`?
   ```bash
   grep TELEGRAM backend/.env
   ```

2. Check bot status:
   ```bash
   curl http://localhost:8080/api/v1/telegram/status
   ```

3. Test Telegram API:
   ```bash
   curl "https://api.telegram.org/bot8582809296:AAFkw9Qv_PunAuto-x03HY57441M-AJQ3W8/getMe"
   ```

**Fix:**
- Make sure `TELEGRAM_BOT_TOKEN` and `TELEGRAM_CHAT_ID` are set in `.env`
- Restart backend after updating `.env`

### Issue 4: Bot Not Running

**Symptom:**
No logs about checking market

**Check:**
```bash
curl http://localhost:8080/api/v1/telegram/status
```

Should show:
```json
{
  "configured": true,
  "running": true,
  "message": "Telegram bot ready"
}
```

**Fix:**
1. Check `TELEGRAM_AUTO_START=true` in `.env`
2. Manually start bot:
   ```bash
   curl -X POST http://localhost:8080/api/v1/telegram/start \
     -H "Content-Type: application/json" \
     -d '{
       "symbol": "BTCUSDT",
       "strategy": "session_trader",
       "filterBuy": true,
       "filterSell": true
     }'
   ```

### Issue 5: Signals Filtered Out

**Symptom:**
```
⏭️  Signal BUY filtered out (filterBuy=false, filterSell=true)
```

**Cause:** Filter settings blocking signals

**Fix:**
Update `.env`:
```env
TELEGRAM_FILTER_BUY=true
TELEGRAM_FILTER_SELL=true
```

Restart backend.

## 🧪 Force Generate a Signal

To test if everything works, you can force generate signals:

```bash
# Try multiple strategies until you get a BUY/SELL
for strategy in session_trader breakout_master liquidity_hunter scalper_pro; do
  echo "Testing $strategy..."
  curl -s -X POST http://localhost:8080/api/v1/backtest/live-signal \
    -H "Content-Type: application/json" \
    -d "{\"symbol\": \"BTCUSDT\", \"strategy\": \"$strategy\"}" | jq '.signal'
done
```

## ✅ Expected Behavior

### When Signal is NONE (90% of the time):
```
🔄 Telegram bot checking market for BTCUSDT...
🔍 Telegram bot generated signal: NONE
ℹ️  No signal (NONE), waiting for next check...
```
**This is NORMAL!** Strategies wait for good opportunities.

### When Signal is BUY/SELL (10% of the time):
```
🔄 Telegram bot checking market for BTCUSDT...
🔍 Telegram bot generated signal: BUY
✅ Signal saved to Supabase: BUY BTCUSDT @ $50000.00
📤 Sent BUY signal to Telegram for BTCUSDT
```
**This is what you want!** Check:
1. Supabase dashboard - see new row
2. Telegram app - see message

## 📊 Monitoring

### Watch Backend Logs:
```bash
cd backend
go run . 2>&1 | grep -E "(🔄|🔍|✅|⚠️|📤)"
```

### Check Supabase:
```bash
curl -s "https://elqhqhjevaizjoghiiss.supabase.co/rest/v1/trading_signals?order=created_at.desc&limit=5" \
  -H "apikey: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImVscWhxaGpldmFpempvZ2hpaXNzIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NjQ3MDEwNTksImV4cCI6MjA4MDI3NzA1OX0.02-CgybOf7PiQSaQ-uJhojKob5Rw_2vkFdyurPNqLvA" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImVscWhxaGpldmFpempvZ2hpaXNzIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NjQ3MDEwNTksImV4cCI6MjA4MDI3NzA1OX0.02-CgybOf7PiQSaQ-uJhojKob5Rw_2vkFdyurPNqLvA" | jq '.[] | {signal_type, symbol, created_at}'
```

### Check Telegram Bot:
```bash
curl http://localhost:8080/api/v1/telegram/status | jq '.'
```

## 🎯 Summary

**Most likely reason:** Signal is NONE (no trading opportunity)

**How to verify:**
1. Restart backend
2. Watch logs for 1-2 minutes
3. Run `./test_signal_generation.sh`
4. Check logs for detailed error messages

**What to look for:**
- ✅ `Signal saved to Supabase` - Working!
- ✅ `Sent signal to Telegram` - Working!
- ℹ️ `Signal is NONE` - Normal, wait for opportunity
- ⚠️ Any error message - Follow fix above

The enhanced logging will tell you exactly what's happening! 🔍
