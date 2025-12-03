# 🔧 Complete Supabase Fix Guide

## 🎯 Quick Fix (Most Common Issue)

**Problem:** Signals not saving to Supabase
**Cause:** Table doesn't exist yet (90% of cases)
**Solution:** Run the SQL schema

### 3-Step Fix:

1. **Open Supabase SQL Editor:**
   - Go to: https://supabase.com/dashboard/project/elqhqhjevajzjoghiiss/sql
   - Click "New Query"

2. **Run the Schema:**
   - Open `supabase-setup.sql` in your project
   - Copy ALL contents (Ctrl+A, Ctrl+C)
   - Paste into SQL Editor
   - Click "Run" (or Ctrl+Enter)

3. **Verify:**
   ```bash
   ./diagnose_supabase.sh
   ```

Done! ✅

---

## 🔍 Detailed Diagnosis

### Step 1: Run Diagnostic Tool

```bash
./diagnose_supabase.sh
```

This will check:
- ✅ Environment variables
- ✅ Supabase connection
- ✅ Table exists
- ✅ INSERT permission
- ✅ SELECT permission

**If all checks pass:** Your Supabase is working!
**If any check fails:** Follow the fix instructions shown

### Step 2: Check Backend Logs

Start your backend and watch for errors:

```bash
cd backend
go run . 2>&1 | grep -i supabase
```

**Look for:**
- ✅ `✅ Signal saved to Supabase: BUY BTCUSDT @ $50000.00`
- ❌ `⚠️ Supabase not configured, skipping signal storage`
- ❌ `❌ Supabase error (status 404): relation does not exist`
- ❌ `❌ Supabase error (status 401): Invalid API key`

### Step 3: Manual Test

Test the API directly:

```bash
./test_supabase_connection.sh
```

This will:
1. Check if table exists
2. Try to insert a test signal
3. Retrieve recent signals

---

## 🛠️ Common Issues & Fixes

### Issue 1: Table Doesn't Exist

**Symptoms:**
- Error: `relation "public.trading_signals" does not exist`
- HTTP 404 when querying table

**Fix:**
```sql
-- Run this in Supabase SQL Editor
-- (or just run the entire supabase-setup.sql file)

CREATE TABLE trading_signals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    symbol TEXT NOT NULL,
    strategy TEXT NOT NULL,
    signal_type TEXT NOT NULL,
    entry_price DECIMAL(20, 8) NOT NULL,
    stop_loss DECIMAL(20, 8) NOT NULL,
    take_profit DECIMAL(20, 8) NOT NULL,
    current_price DECIMAL(20, 8) NOT NULL,
    risk_reward DECIMAL(10, 2) NOT NULL,
    profit_loss DECIMAL(20, 8),
    profit_loss_percent DECIMAL(10, 4),
    status TEXT DEFAULT 'ACTIVE',
    result TEXT,
    progress DECIMAL(5, 2) DEFAULT 0,
    filter_buy BOOLEAN DEFAULT true,
    filter_sell BOOLEAN DEFAULT true,
    signal_time TIMESTAMP WITH TIME ZONE NOT NULL,
    closed_at TIMESTAMP WITH TIME ZONE
);

-- Enable RLS
ALTER TABLE trading_signals ENABLE ROW LEVEL SECURITY;

-- Create policy
CREATE POLICY "Allow all operations for anon"
ON trading_signals FOR ALL TO anon
USING (true) WITH CHECK (true);
```

### Issue 2: RLS Blocking Inserts

**Symptoms:**
- Error: `new row violates row-level security policy`
- HTTP 403 Forbidden

**Fix:**
```sql
-- Run in Supabase SQL Editor
ALTER TABLE trading_signals ENABLE ROW LEVEL SECURITY;

-- Drop existing policies
DROP POLICY IF EXISTS "Allow all operations for anon" ON trading_signals;
DROP POLICY IF EXISTS "Allow all operations for authenticated" ON trading_signals;
DROP POLICY IF EXISTS "Allow all operations for service_role" ON trading_signals;

-- Create new policies
CREATE POLICY "Allow all operations for anon"
ON trading_signals FOR ALL TO anon
USING (true) WITH CHECK (true);

CREATE POLICY "Allow all operations for authenticated"
ON trading_signals FOR ALL TO authenticated
USING (true) WITH CHECK (true);

CREATE POLICY "Allow all operations for service_role"
ON trading_signals FOR ALL TO service_role
USING (true) WITH CHECK (true);
```

### Issue 3: Wrong API Key

**Symptoms:**
- Error: `Invalid API key`
- HTTP 401 Unauthorized

**Fix:**
1. Go to Supabase Dashboard → Settings → API
2. Copy the `anon` `public` key (NOT the service_role key)
3. Update `backend/.env`:
   ```
   SUPABASE_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
   ```
4. Restart backend

### Issue 4: Wrong URL

**Symptoms:**
- Error: `connection refused` or `timeout`
- Can't connect to Supabase

**Fix:**
1. Go to Supabase Dashboard → Settings → API
2. Copy the URL (should be: `https://elqhqhjevajzjoghiiss.supabase.co`)
3. Update `backend/.env`:
   ```
   SUPABASE_URL=https://elqhqhjevajzjoghiiss.supabase.co
   ```
4. Restart backend

### Issue 5: Column Mismatch

**Symptoms:**
- Error: `column "xyz" does not exist`
- HTTP 400 Bad Request

**Fix:**
Drop and recreate table:
```sql
-- Run in Supabase SQL Editor
DROP TABLE IF EXISTS trading_signals CASCADE;

-- Then run the entire supabase-setup.sql file
```

---

## ✅ Verification Steps

### 1. Check Table in Dashboard

1. Go to Supabase Dashboard
2. Click **Table Editor** in left sidebar
3. Look for `trading_signals` table
4. Click on it to see structure

**Expected columns:**
- id (uuid)
- created_at (timestamptz)
- updated_at (timestamptz)
- symbol (text)
- strategy (text)
- signal_type (text)
- entry_price (numeric)
- stop_loss (numeric)
- take_profit (numeric)
- current_price (numeric)
- risk_reward (numeric)
- profit_loss (numeric)
- profit_loss_percent (numeric)
- status (text)
- result (text)
- progress (numeric)
- filter_buy (bool)
- filter_sell (bool)
- signal_time (timestamptz)
- closed_at (timestamptz)

### 2. Test Manual Insert

In Table Editor:
1. Click **Insert** → **Insert row**
2. Fill in required fields:
   - symbol: `BTCUSDT`
   - strategy: `test`
   - signal_type: `BUY`
   - entry_price: `50000`
   - stop_loss: `49500`
   - take_profit: `51000`
   - current_price: `50000`
   - risk_reward: `2`
   - signal_time: `now()`
3. Click **Save**

**If this works:** Table is set up correctly!

### 3. Test via Backend

```bash
# Start backend
cd backend
go run .

# In another terminal, generate signal
curl -X POST http://localhost:8080/api/v1/live-signal \
  -H "Content-Type: application/json" \
  -d '{"symbol": "BTCUSDT", "strategy": "session_trader"}'

# Check if saved
curl http://localhost:8080/api/v1/signals/recent?limit=5
```

### 4. Check Backend Logs

Look for:
```
✅ Signal saved to Supabase: BUY BTCUSDT @ $50000.00
```

### 5. Verify in Supabase

1. Go to Table Editor → trading_signals
2. Refresh the page
3. You should see your signals

---

## 🚀 Complete Setup Checklist

- [ ] Run `./diagnose_supabase.sh` - all checks pass
- [ ] Table `trading_signals` exists in Supabase
- [ ] RLS policies are enabled
- [ ] Environment variables are correct in `backend/.env`
- [ ] Backend starts without errors
- [ ] Manual insert works in Table Editor
- [ ] API insert works via curl
- [ ] Backend logs show "✅ Signal saved"
- [ ] Signals visible in Supabase Table Editor
- [ ] Telegram bot sends signals (if enabled)

---

## 🆘 Still Not Working?

### Enable Debug Logging

Add to `backend/signal_storage.go` after line 62:

```go
log.Printf("🔍 Supabase URL: %s", url)
log.Printf("🔍 Signal JSON: %s", string(jsonData))
log.Printf("🔍 API Key: %s...", supabaseKey[:20])
```

Restart backend and check logs for exact error.

### Check Supabase Logs

1. Go to Supabase Dashboard
2. Click **Logs** → **API Logs**
3. Look for POST requests to `/rest/v1/trading_signals`
4. Check for errors

### Contact Support

If nothing works, provide:
1. Output of `./diagnose_supabase.sh`
2. Backend logs (last 50 lines)
3. Supabase API logs
4. Screenshot of Table Editor

---

## 📝 Summary

**Most common issue:** Table doesn't exist
**Quick fix:** Run `supabase-setup.sql` in SQL Editor
**Verify:** Run `./diagnose_supabase.sh`

That's it! 🎉
