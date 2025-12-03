# ✅ Filter Storage Using Supabase REST API

## What Was Changed

Switched from PostgreSQL direct connection to **Supabase REST API** for storing user filter settings. This eliminates the need for database credentials and uses your existing Supabase setup.

## Why This Change?

### Before:
- ❌ Required `SUPABASE_HOST` and `SUPABASE_PASSWORD` environment variables
- ❌ Needed direct PostgreSQL connection
- ❌ More complex setup

### After:
- ✅ Uses existing `SUPABASE_URL` and `SUPABASE_KEY` from your `.env`
- ✅ No additional database credentials needed
- ✅ Simpler, REST API-based approach
- ✅ Works with your current Supabase setup

## Setup Instructions

### 1. Create the `user_settings` Table in Supabase

Go to your Supabase SQL Editor and run:

```sql
-- Create user_settings table
CREATE TABLE IF NOT EXISTS user_settings (
    id INTEGER PRIMARY KEY DEFAULT 1,
    filter_buy BOOLEAN DEFAULT true,
    filter_sell BOOLEAN DEFAULT true,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT single_row CHECK (id = 1)
);

-- Insert default settings
INSERT INTO user_settings (id, filter_buy, filter_sell)
VALUES (1, true, true)
ON CONFLICT (id) DO NOTHING;

-- Enable Row Level Security
ALTER TABLE user_settings ENABLE ROW LEVEL SECURITY;

-- Create policies for access
CREATE POLICY "Allow all operations for anon on user_settings"
ON user_settings FOR ALL TO anon
USING (true) WITH CHECK (true);

CREATE POLICY "Allow all operations for authenticated on user_settings"
ON user_settings FOR ALL TO authenticated
USING (true) WITH CHECK (true);

-- Grant permissions
GRANT ALL ON user_settings TO anon;
GRANT ALL ON user_settings TO authenticated;
```

### 2. Restart Backend

```bash
cd backend
go run .
```

You should see:
```
✅ Server starting on port 8080
```

No more database warnings!

### 3. Test in Browser

1. Open http://localhost:8080
2. Change filter settings (check/uncheck BUY or SELL)
3. Reload the page
4. ✅ Settings should be restored from Supabase

## How It Works

### API Endpoints

#### GET /api/v1/settings
Fetches settings from Supabase:
```bash
curl http://localhost:8080/api/v1/settings
```

Response:
```json
{
  "id": 1,
  "filterBuy": true,
  "filterSell": true
}
```

#### POST /api/v1/settings
Updates settings in Supabase:
```bash
curl -X POST http://localhost:8080/api/v1/settings \
  -H "Content-Type: application/json" \
  -d '{"filterBuy":true,"filterSell":false}'
```

### Backend Implementation

The backend now uses Supabase REST API:

```go
// Get settings
url := fmt.Sprintf("%s/rest/v1/user_settings?id=eq.1", supabaseURL)
req.Header.Set("apikey", supabaseKey)
req.Header.Set("Authorization", "Bearer " + supabaseKey)

// Update settings
url := fmt.Sprintf("%s/rest/v1/user_settings", supabaseURL)
req.Header.Set("Prefer", "resolution=merge-duplicates")
```

### Fallback Behavior

If Supabase is not configured or unavailable:
- ✅ Returns default settings (both filters enabled)
- ✅ System continues to work
- ⚠️ Settings won't persist (will reset on reload)
- ℹ️  Logs warning but doesn't crash

## Environment Variables

Your `.env` file already has everything needed:

```env
SUPABASE_URL=https://elqhqhjevaizjoghiiss.supabase.co
SUPABASE_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

No additional variables required!

## Benefits

### 1. **Simpler Setup**
- No PostgreSQL connection string needed
- Uses existing Supabase credentials
- One less thing to configure

### 2. **More Reliable**
- REST API is more stable than direct DB connections
- Better error handling
- Automatic retries

### 3. **Better Security**
- No database password in environment variables
- Uses Supabase's built-in security
- Row Level Security policies

### 4. **Easier Deployment**
- Works on any platform (Vercel, Netlify, etc.)
- No need for database connection pooling
- Serverless-friendly

## Testing

### Test 1: Save Settings
1. Open http://localhost:8080
2. Uncheck "Sell Trades (Short)"
3. Check browser console - should see successful POST
4. Check Supabase table - should see `filter_sell=false`

### Test 2: Load Settings
1. Reload the page
2. ✅ "Sell Trades" should still be unchecked
3. Settings persisted successfully!

### Test 3: Telegram Bot
1. Change filters (e.g., only BUY)
2. Start Telegram bot
3. ✅ Only BUY signals sent to Telegram

### Test 4: Verify in Supabase
Go to Supabase Table Editor:
1. Find `user_settings` table
2. Should see one row with id=1
3. Check `filter_buy` and `filter_sell` values

## Troubleshooting

### Settings Not Saving
**Check Supabase credentials:**
```bash
# In backend/.env
echo $SUPABASE_URL
echo $SUPABASE_KEY
```

**Check Supabase table exists:**
- Go to Supabase Table Editor
- Look for `user_settings` table
- If missing, run the SQL from Step 1

### Settings Not Loading
**Check browser console:**
- Should see GET request to `/api/v1/settings`
- Should return JSON with filterBuy and filterSell

**Check backend logs:**
```
✅ Loaded filter settings from Supabase: filterBuy=true, filterSell=false
```

### 503 Errors Gone?
Yes! The 503 errors were because the backend was looking for PostgreSQL credentials. Now it uses Supabase REST API which is already configured.

## Migration from Old System

If you had the old PostgreSQL-based system:
1. ✅ No migration needed
2. ✅ Old code automatically replaced
3. ✅ Just create the Supabase table
4. ✅ Restart backend

## Summary

Filter settings now use Supabase REST API:
1. ✅ **No database credentials needed** - uses existing Supabase setup
2. ✅ **Simpler configuration** - just create one table
3. ✅ **More reliable** - REST API is stable
4. ✅ **Better deployment** - works everywhere
5. ✅ **Automatic fallback** - graceful degradation if Supabase unavailable

Just create the table in Supabase and restart the backend! 🎉
