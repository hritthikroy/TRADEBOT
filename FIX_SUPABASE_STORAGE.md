# 🔧 Fix Supabase Storage Issue

## Problem
Signals are not being saved to Supabase database.

## Most Common Causes

### 1. **Table Doesn't Exist** (90% of cases)
The `trading_signals` table hasn't been created in Supabase yet.

### 2. **Wrong Credentials**
SUPABASE_URL or SUPABASE_KEY is incorrect.

### 3. **RLS Policies**
Row Level Security is blocking inserts.

### 4. **Column Mismatch**
Table columns don't match the Go struct.

## 🔍 Step-by-Step Diagnosis

### Step 1: Test Supabase Connection

Run this command:
```bash
./test_supabase_connection.sh
```

**Expected Output:**
- ✅ If table exists: You'll see `[]` or existing signals
- ❌ If table doesn't exist: You'll see error like `relation "public.trading_signals" does not exist`

### Step 2: Check Environment Variables

```bash
cd backend
cat .env | grep SUPABASE
```

**Should show:**
```
SUPABASE_URL=https://elqhqhjevajzjoghiiss.supabase.co
SUPABASE_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Step 3: Verify Backend Logs

Start the backend and look for these messages:
```bash
cd backend
go run .
```

**Look for:**
- ✅ `✅ Signal saved to Supabase: BUY BTCUSDT @ $50000.00`
- ❌ `⚠️ Supabase not configured, skipping signal storage`
- ❌ `❌ Supabase error (status 404): relation does not exist`

## ✅ Solution: Create the Table

### Option 1: Using Supabase Dashboard (Recommended)

1. **Go to Supabase Dashboard:**
   - Open: https://supabase.com/dashboard/project/elqhqhjevajzjoghiiss
   - Login if needed

2. **Open SQL Editor:**
   - Click **SQL Editor** in the left sidebar
   - Click **New Query**

3. **Copy the SQL:**
   - Open the file `supabase-setup.sql` in your project
   - Copy ALL the contents (Ctrl+A, Ctrl+C)

4. **Paste and Run:**
   - Paste into the SQL Editor
   - Click **Run** (or press Ctrl+Enter)

5. **Verify Success:**
   - You should see: "Success. No rows returned"
   - Click **Table Editor** in left sidebar
   - You should see `trading_signals` table

### Option 2: Using Supabase CLI

```bash
# Install Supabase CLI (if not installed)
brew install supabase/tap/supabase

# Login
supabase login

# Link to your project
supabase link --project-ref elqhqhjevajzjoghiiss

# Run the migration
supabase db push
```

## 🧪 Test After Setup

### Test 1: Manual Insert via Dashboard

1. Go to **Table Editor** → `trading_signals`
2. Click **Insert** → **Insert row**
3. Fill in:
   - symbol: `BTCUSDT`
   - strategy: `test`
   - signal_type: `BUY`
   - entry_price: `50000`
   - stop_loss: `49500`
   - take_profit: `51000`
   - current_price: `50000`
   - risk_reward: `2`
   - status: `ACTIVE`
   - progress: `0`
   - filter_buy: `true`
   - filter_sell: `true`
   - signal_time: `now()`
4. Click **Save**

**If this works:** Table is set up correctly!

### Test 2: Test via Backend

```bash
# Generate a signal
curl -X POST http://localhost:8080/api/v1/live-signal \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "strategy": "session_trader"
  }'

# Check if it was saved
curl http://localhost:8080/api/v1/signals/recent?limit=5
```

**Expected:** You should see the signal in the response.

### Test 3: Check Supabase Dashboard

1. Go to **Table Editor** → `trading_signals`
2. You should see your signals appearing in real-time
3. Refresh the page to see new signals

## 🔧 Common Fixes

### Fix 1: Table Doesn't Exist
**Solution:** Run `supabase-setup.sql` in SQL Editor (see above)

### Fix 2: Wrong API Key
**Solution:** 
1. Go to Supabase Dashboard → Settings → API
2. Copy the `anon` `public` key
3. Update `backend/.env`:
   ```
   SUPABASE_KEY=your_new_key_here
   ```
4. Restart backend

### Fix 3: RLS Blocking Inserts
**Solution:** The `supabase-setup.sql` already includes RLS policies. But if needed:

```sql
-- Run this in SQL Editor
ALTER TABLE trading_signals ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Allow all operations for anon"
ON trading_signals
FOR ALL
TO anon
USING (true)
WITH CHECK (true);
```

### Fix 4: Column Mismatch
**Solution:** Drop and recreate table:

```sql
-- Run in SQL Editor
DROP TABLE IF EXISTS trading_signals CASCADE;

-- Then run the full supabase-setup.sql
```

## 📊 Verify It's Working

### Check 1: Backend Logs
You should see:
```
✅ Signal saved to Supabase: BUY BTCUSDT @ $50000.00
```

### Check 2: Supabase Dashboard
- Go to Table Editor
- See signals appearing
- Check timestamps are recent

### Check 3: API Response
```bash
curl http://localhost:8080/api/v1/signals/recent?limit=5
```
Should return JSON with signals.

### Check 4: Telegram
If Telegram bot is running, you should get messages AND see them in Supabase.

## 🎯 Quick Checklist

- [ ] Supabase table created (run `supabase-setup.sql`)
- [ ] Environment variables set correctly in `backend/.env`
- [ ] Backend restarted after changes
- [ ] RLS policies enabled (included in setup SQL)
- [ ] Test insert works via dashboard
- [ ] Test insert works via API
- [ ] Backend logs show "✅ Signal saved"
- [ ] Signals visible in Supabase Table Editor

## 🆘 Still Not Working?

### Enable Debug Mode

Add this to `backend/signal_storage.go` after line 62:

```go
log.Printf("🔍 Attempting to save signal to: %s", url)
log.Printf("🔍 Signal data: %s", string(jsonData))
```

Then check the logs for the exact error message.

### Check Supabase Logs

1. Go to Supabase Dashboard
2. Click **Logs** in left sidebar
3. Look for errors related to `trading_signals`

### Manual Test

```bash
# Test direct API call
curl -X POST "https://elqhqhjevajzjoghiiss.supabase.co/rest/v1/trading_signals" \
  -H "apikey: YOUR_SUPABASE_KEY" \
  -H "Authorization: Bearer YOUR_SUPABASE_KEY" \
  -H "Content-Type: application/json" \
  -H "Prefer: return=representation" \
  -d '{
    "symbol": "BTCUSDT",
    "strategy": "test",
    "signal_type": "BUY",
    "entry_price": 50000,
    "stop_loss": 49500,
    "take_profit": 51000,
    "current_price": 50000,
    "risk_reward": 2,
    "status": "ACTIVE",
    "progress": 0,
    "filter_buy": true,
    "filter_sell": true,
    "signal_time": "2024-01-15T10:00:00Z"
  }'
```

If this works but the backend doesn't, there's an issue with the Go code.
If this doesn't work, there's an issue with Supabase setup.

## 📝 Summary

**Most likely issue:** Table doesn't exist yet.

**Quick fix:**
1. Open Supabase SQL Editor
2. Run `supabase-setup.sql`
3. Restart backend
4. Test with `./test_supabase_connection.sh`

Done! 🎉
