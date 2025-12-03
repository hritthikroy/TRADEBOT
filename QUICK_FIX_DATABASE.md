# 🔧 Quick Fix: Database Setup

## Error You're Seeing
```
GET http://localhost:8080/api/v1/settings 503 (Service Unavailable)
POST http://localhost:8080/api/v1/settings 503 (Service Unavailable)
```

## What This Means
The `user_settings` table doesn't exist in your database yet.

## Quick Fix (Automatic)

### Option 1: Restart Backend (Easiest)
The backend now has automatic migrations. Just restart it:

```bash
cd backend
go run .
```

You should see:
```
🔄 Running database migrations...
✅ user_settings table ready
✅ Default filter settings initialized
✅ Database migrations complete
```

### Option 2: Manual SQL (If Option 1 Fails)

Run this SQL in your database:

```sql
-- Create user_settings table
CREATE TABLE IF NOT EXISTS user_settings (
    id INTEGER PRIMARY KEY DEFAULT 1,
    filter_buy BOOLEAN DEFAULT true,
    filter_sell BOOLEAN DEFAULT true,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT single_row CHECK (id = 1)
);

-- Insert default settings
INSERT INTO user_settings (id, filter_buy, filter_sell)
VALUES (1, true, true)
ON CONFLICT (id) DO NOTHING;
```

## What Was Fixed

### 1. **Automatic Migrations** (`backend/migrations.go`)
✅ Backend now creates the table automatically on startup
✅ No manual SQL needed

### 2. **Better Error Handling** (`backend/user_settings.go`)
✅ Gracefully handles missing table
✅ Falls back to defaults if database unavailable
✅ Better logging

### 3. **Startup Integration** (`backend/main.go`)
✅ Runs migrations after database initialization
✅ Ensures table exists before use

## Verification

After restarting the backend, check the logs:

### Success Logs:
```
✅ Database connected successfully
🔄 Running database migrations...
✅ user_settings table ready
✅ Default filter settings initialized
✅ Database migrations complete
```

### If Database Not Available:
```
⚠️  Database not available, skipping migrations
⚠️  Database not available, using default filters (both enabled)
```

This is OK - the system will work with default filters (both BUY and SELL enabled).

## Testing

1. **Restart backend:**
   ```bash
   cd backend
   go run .
   ```

2. **Open browser console** (F12)

3. **Reload the page**

4. **Check for errors:**
   - ✅ No more 503 errors
   - ✅ Filters load successfully
   - ✅ Settings save successfully

## Fallback Behavior

If the database is not available, the system will:
- ✅ Use default filters (BUY=true, SELL=true)
- ✅ Continue working normally
- ✅ Log warnings but not crash
- ⚠️ Settings won't persist (will reset on reload)

## Database Options

### PostgreSQL (Local)
```bash
# Create table
psql -U your_user -d your_database -f database_setup_complete.sql
```

### Supabase (Cloud)
1. Go to Supabase SQL Editor
2. Paste contents of `database_setup_complete.sql`
3. Click "Run"

### SQLite (File-based)
The migrations will create the table automatically.

## Summary

The issue is fixed with automatic migrations. Just:
1. ✅ Restart the backend
2. ✅ Table will be created automatically
3. ✅ Default settings will be initialized
4. ✅ Everything will work!

No manual SQL needed anymore! 🎉
