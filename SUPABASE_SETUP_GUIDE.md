# Supabase Setup Guide for Trading Bot

## Step 1: Create Supabase Account (5 minutes)

1. Go to https://supabase.com
2. Click "Start your project"
3. Sign up with GitHub
4. Create new project:
   - **Name:** `tradebot`
   - **Database Password:** Create a strong password (SAVE THIS!)
   - **Region:** Choose closest to you (e.g., Mumbai/Singapore for India)
5. Wait 2 minutes for database initialization

## Step 2: Get API Credentials

1. In your Supabase dashboard, go to **Settings** → **API**
2. Copy these 2 values:
   - **Project URL** (e.g., `https://xxxxx.supabase.co`)
   - **anon public key** (long string starting with `eyJ...`)

## Step 3: Create Database Table

1. In Supabase dashboard, go to **SQL Editor**
2. Click **New Query**
3. Copy the entire content from `supabase-setup.sql` file
4. Paste it into the SQL editor
5. Click **Run** (bottom right)
6. You should see "Success. No rows returned"

## Step 4: Configure Your App

1. Open `supabase-config.js`
2. Replace these lines:
   ```javascript
   const SUPABASE_URL = 'YOUR_SUPABASE_URL';
   const SUPABASE_ANON_KEY = 'YOUR_SUPABASE_ANON_KEY';
   ```
   With your actual values:
   ```javascript
   const SUPABASE_URL = 'https://xxxxx.supabase.co';
   const SUPABASE_ANON_KEY = 'eyJhbGc...your-key-here';
   ```

## Step 5: Add Supabase to HTML Files

Add this script tag BEFORE your other scripts in both `index.html` and `signal-tracker.html`:

```html
<!-- Add this in the <head> section -->
<script src="https://cdn.jsdelivr.net/npm/@supabase/supabase-js@2"></script>
<script src="supabase-config.js"></script>
```

## Step 6: Test the Connection

1. Open browser console (F12)
2. Type: `SupabaseDB.getAllSignals()`
3. You should see: `[]` (empty array - no signals yet)
4. If you see an error, check your URL and API key

## What You Get

### ✅ Features
- **Cross-device sync** - Access signals from any device
- **Historical data** - Never lose old signals
- **Real-time updates** - See changes instantly
- **Advanced analytics** - Query your performance
- **Backup** - Data stored safely in cloud

### 📊 Analytics Queries

Run these in Supabase SQL Editor:

**Win Rate by Kill Zone:**
```sql
SELECT * FROM signal_analytics 
ORDER BY win_rate DESC;
```

**Best Performing Patterns:**
```sql
SELECT pattern_type, 
       COUNT(*) as total,
       AVG(profit_percent) as avg_profit,
       ROUND(100.0 * SUM(CASE WHEN status='win' THEN 1 ELSE 0 END) / COUNT(*), 2) as win_rate
FROM trading_signals
WHERE status IN ('win', 'loss')
GROUP BY pattern_type
ORDER BY win_rate DESC;
```

**Performance by Time of Day:**
```sql
SELECT 
    EXTRACT(HOUR FROM created_at) as hour,
    COUNT(*) as signals,
    AVG(profit_percent) as avg_profit
FROM trading_signals
WHERE status IN ('win', 'loss')
GROUP BY hour
ORDER BY hour;
```

**Best Currency Pairs:**
```sql
SELECT symbol,
       COUNT(*) as total_trades,
       SUM(CASE WHEN status='win' THEN 1 ELSE 0 END) as wins,
       ROUND(AVG(profit_percent), 2) as avg_profit,
       ROUND(100.0 * SUM(CASE WHEN status='win' THEN 1 ELSE 0 END) / COUNT(*), 2) as win_rate
FROM trading_signals
WHERE status IN ('win', 'loss')
GROUP BY symbol
ORDER BY win_rate DESC;
```

## Troubleshooting

**Error: "Invalid API key"**
- Check you copied the `anon public` key, not the `service_role` key
- Make sure there are no extra spaces

**Error: "relation does not exist"**
- Run the SQL setup script again
- Make sure you're in the correct project

**No data showing:**
- Check browser console for errors
- Verify Supabase URL is correct
- Try: `SupabaseDB.getAllSignals()` in console

## Next Steps

After setup, your signals will automatically:
1. Save to Supabase when generated
2. Sync across all devices
3. Update in real-time
4. Be available for analysis

## Free Tier Limits

- **Database:** 500 MB
- **Bandwidth:** 2 GB/month
- **API Requests:** Unlimited
- **Real-time connections:** 200 concurrent

This is more than enough for trading signals!
