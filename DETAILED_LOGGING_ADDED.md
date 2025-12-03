# 🔍 Detailed Logging Added - Debug Signal Flow

## ✅ What I Added

I've added **extremely detailed logging** to every step of the signal flow. Now you'll see EXACTLY what's happening and where it fails.

## 📊 What You'll See in Logs

### 1. Signal Generation
```
🔍 Generated signal: BUY for BTCUSDT using session_trader strategy
```
OR
```
ℹ️  Signal is NONE, not sending to Telegram
```

### 2. Supabase Save Attempt
```
🔍 Saving to Supabase: {"symbol":"BTCUSDT","strategy":"session_trader",...}
🔍 Supabase URL: https://elqhqhjevaizjoghiiss.supabase.co/rest/v1/trading_signals
🔍 Supabase response status: 201
🔍 Supabase response body: [{"id":"...","symbol":"BTCUSDT",...}]
✅ Signal saved to Supabase: BUY BTCUSDT @ $50000.00
```

OR if error:
```
❌ Failed to send request to Supabase: connection refused
❌ Supabase error (status 404): {"message":"relation does not exist"}
❌ Supabase error (status 401): {"message":"Invalid API key"}
```

### 3. Telegram Send Attempt
```
🔍 Sending to Telegram - ChatID: 8145172959
🔍 Telegram API response status: 200
✅ Message sent to Telegram successfully
📤 Sent BUY signal to Telegram for BTCUSDT
```

OR if error:
```
⚠️  Telegram bot is nil, cannot send signal
⚠️  Telegram bot token is empty, cannot send signal
❌ Failed to send to Telegram API: connection refused
❌ Telegram API error (status 400): {"description":"Bad Request: chat not found"}
❌ Telegram API error (status 401): {"description":"Unauthorized"}
```

## 🧪 How to Debug

### Step 1: Restart Backend
```bash
cd backend
go run .
```

Watch for startup messages:
```
✅ Telegram bot initialized
✅ Telegram bot auto-started for BTCUSDT with session_trader strategy
```

### Step 2: Generate a Signal
In another terminal:
```bash
./debug_signal_flow.sh
```

### Step 3: Read the Logs

The backend logs will show you EXACTLY what happened at each step.

## 🔍 Common Error Messages & Fixes

### Error 1: "Signal is NONE"
```
ℹ️  Signal is NONE, not sending to Telegram
```

**Meaning:** No trading opportunity right now
**Fix:** This is NORMAL! Wait or try different strategy/symbol

### Error 2: "Supabase error (status 404)"
```
❌ Supabase error (status 404): relation "public.trading_signals" does not exist
```

**Meaning:** Table doesn't exist
**Fix:** Run `supabase-setup.sql` in Supabase SQL Editor

### Error 3: "Supabase error (status 401)"
```
❌ Supabase error (status 401): Invalid API key
```

**Meaning:** Wrong API key
**Fix:** Check `SUPABASE_KEY` in `backend/.env`

### Error 4: "connection refused"
```
❌ Failed to send request to Supabase: connection refused
```

**Meaning:** Wrong URL or network issue
**Fix:** Check `SUPABASE_URL` in `backend/.env`

### Error 5: "Telegram bot is nil"
```
⚠️  Telegram bot is nil, cannot send signal
```

**Meaning:** Telegram not initialized
**Fix:** Check `TELEGRAM_BOT_TOKEN` and `TELEGRAM_CHAT_ID` in `backend/.env`

### Error 6: "Telegram API error (status 400)"
```
❌ Telegram API error (status 400): Bad Request: chat not found
```

**Meaning:** Wrong chat ID
**Fix:** Verify `TELEGRAM_CHAT_ID` in `backend/.env`

### Error 7: "Telegram API error (status 401)"
```
❌ Telegram API error (status 401): Unauthorized
```

**Meaning:** Wrong bot token
**Fix:** Verify `TELEGRAM_BOT_TOKEN` in `backend/.env`

## 📋 Complete Debug Checklist

### Before Running:
- [ ] Backend `.env` file has correct `SUPABASE_URL`
- [ ] Backend `.env` file has correct `SUPABASE_KEY`
- [ ] Backend `.env` file has correct `TELEGRAM_BOT_TOKEN`
- [ ] Backend `.env` file has correct `TELEGRAM_CHAT_ID`
- [ ] Supabase table `trading_signals` exists
- [ ] Backend is running

### After Generating Signal:
- [ ] Check logs for "Generated signal: BUY/SELL" (not NONE)
- [ ] Check logs for "Supabase response status: 201"
- [ ] Check logs for "Signal saved to Supabase"
- [ ] Check logs for "Telegram API response status: 200"
- [ ] Check logs for "Message sent to Telegram"
- [ ] Check Supabase dashboard for new row
- [ ] Check Telegram app for message

## 🎯 Example: Perfect Flow

When everything works, you'll see:
```
🔍 Generated signal: BUY for BTCUSDT using session_trader strategy
🔍 Saving to Supabase: {"symbol":"BTCUSDT",...}
🔍 Supabase URL: https://elqhqhjevaizjoghiiss.supabase.co/rest/v1/trading_signals
🔍 Supabase response status: 201
🔍 Supabase response body: [{"id":"abc-123",...}]
✅ Signal saved to Supabase: BUY BTCUSDT @ $50000.00
🔍 Sending to Telegram - ChatID: 8145172959
🔍 Telegram API response status: 200
✅ Message sent to Telegram successfully
📤 Sent BUY signal to Telegram for BTCUSDT
```

## 🎯 Example: Supabase Error

If Supabase fails:
```
🔍 Generated signal: BUY for BTCUSDT using session_trader strategy
🔍 Saving to Supabase: {"symbol":"BTCUSDT",...}
🔍 Supabase URL: https://elqhqhjevaizjoghiiss.supabase.co/rest/v1/trading_signals
🔍 Supabase response status: 404
🔍 Supabase response body: {"message":"relation \"public.trading_signals\" does not exist"}
❌ Supabase error (status 404): relation "public.trading_signals" does not exist
⚠️  Failed to save signal to Supabase: supabase returned status 404
```

**Fix:** Run `supabase-setup.sql`

## 🎯 Example: Telegram Error

If Telegram fails:
```
🔍 Generated signal: BUY for BTCUSDT using session_trader strategy
✅ Signal saved to Supabase: BUY BTCUSDT @ $50000.00
🔍 Sending to Telegram - ChatID: 8145172959
🔍 Telegram API response status: 400
❌ Telegram API error (status 400): {"description":"Bad Request: chat not found"}
```

**Fix:** Verify `TELEGRAM_CHAT_ID`

## 🚀 Quick Test

```bash
# 1. Start backend
cd backend
go run .

# 2. In another terminal, run debug script
./debug_signal_flow.sh

# 3. Copy any error messages from backend logs
# 4. We can fix them based on the error!
```

## 📝 Summary

**Added logging to:**
1. ✅ Signal generation
2. ✅ Supabase save (URL, request, response, status)
3. ✅ Telegram send (ChatID, response, status)

**Now you can see:**
- Exactly what data is being sent
- Exact error messages from APIs
- Which step is failing
- Why it's failing

**The logs will tell you everything!** 🔍

Just restart your backend and run `./debug_signal_flow.sh` to see the detailed logs!
