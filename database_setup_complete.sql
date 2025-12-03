-- ============================================
-- COMPLETE DATABASE SETUP FOR TRADING SIGNALS
-- This file includes both fresh setup and migration
-- ============================================

-- ============================================
-- OPTION 1: FRESH INSTALLATION
-- Use this if you're setting up for the first time
-- ============================================

-- Drop existing table if needed (CAUTION: This deletes all data!)
-- Uncomment the line below only if you want to start fresh
-- DROP TABLE IF EXISTS trading_signals CASCADE;

-- Create the trading_signals table
CREATE TABLE IF NOT EXISTS trading_signals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Signal Details
    symbol TEXT NOT NULL,
    strategy TEXT NOT NULL,
    signal_type TEXT NOT NULL,              -- 'BUY' or 'SELL'
    entry_price DECIMAL(20, 8) NOT NULL,
    stop_loss DECIMAL(20, 8) NOT NULL,
    take_profit DECIMAL(20, 8) NOT NULL,
    tp1 DECIMAL(20, 8),                     -- Take Profit 1 (33%)
    tp2 DECIMAL(20, 8),                     -- Take Profit 2 (33%)
    tp3 DECIMAL(20, 8),                     -- Take Profit 3 (34%)
    current_price DECIMAL(20, 8) NOT NULL,
    risk_reward DECIMAL(10, 2) NOT NULL,
    
    -- Performance Tracking
    profit_loss DECIMAL(20, 8),
    profit_loss_percent DECIMAL(10, 4),
    status TEXT DEFAULT 'ACTIVE',           -- 'ACTIVE', 'HIT_TP', 'HIT_SL', 'CLOSED'
    result TEXT,                            -- 'WIN', 'LOSS', 'BREAKEVEN'
    progress DECIMAL(5, 2) DEFAULT 0,
    
    -- Filters
    filter_buy BOOLEAN DEFAULT true,
    filter_sell BOOLEAN DEFAULT true,
    
    -- Timestamps
    signal_time TIMESTAMP WITH TIME ZONE NOT NULL,
    closed_at TIMESTAMP WITH TIME ZONE
);

-- ============================================
-- OPTION 2: MIGRATION FOR EXISTING TABLE
-- Use this if you already have trading_signals table
-- ============================================

-- Add TP columns if they don't exist
ALTER TABLE trading_signals 
ADD COLUMN IF NOT EXISTS tp1 DECIMAL(20, 8),
ADD COLUMN IF NOT EXISTS tp2 DECIMAL(20, 8),
ADD COLUMN IF NOT EXISTS tp3 DECIMAL(20, 8);

-- ============================================
-- INDEXES FOR PERFORMANCE
-- ============================================

-- Create indexes for faster queries (IF NOT EXISTS requires PostgreSQL 9.5+)
CREATE INDEX IF NOT EXISTS idx_signals_symbol ON trading_signals(symbol);
CREATE INDEX IF NOT EXISTS idx_signals_strategy ON trading_signals(strategy);
CREATE INDEX IF NOT EXISTS idx_signals_status ON trading_signals(status);
CREATE INDEX IF NOT EXISTS idx_signals_signal_type ON trading_signals(signal_type);
CREATE INDEX IF NOT EXISTS idx_signals_created_at ON trading_signals(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_signals_signal_time ON trading_signals(signal_time DESC);

-- ============================================
-- ROW LEVEL SECURITY (RLS)
-- ============================================

-- Enable Row Level Security
ALTER TABLE trading_signals ENABLE ROW LEVEL SECURITY;

-- Drop existing policies if they exist
DROP POLICY IF EXISTS "Allow all operations for anon" ON trading_signals;
DROP POLICY IF EXISTS "Allow all operations for authenticated" ON trading_signals;
DROP POLICY IF EXISTS "Allow all operations for service_role" ON trading_signals;

-- Create policies for full access
CREATE POLICY "Allow all operations for anon"
ON trading_signals
FOR ALL
TO anon
USING (true)
WITH CHECK (true);

CREATE POLICY "Allow all operations for authenticated"
ON trading_signals
FOR ALL
TO authenticated
USING (true)
WITH CHECK (true);

CREATE POLICY "Allow all operations for service_role"
ON trading_signals
FOR ALL
TO service_role
USING (true)
WITH CHECK (true);

-- Grant permissions
GRANT ALL ON trading_signals TO anon;
GRANT ALL ON trading_signals TO authenticated;
GRANT ALL ON trading_signals TO service_role;

-- ============================================
-- AUTO-UPDATE TIMESTAMP TRIGGER
-- ============================================

-- Create function to auto-update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Drop existing trigger if it exists
DROP TRIGGER IF EXISTS update_signals_updated_at ON trading_signals;

-- Create trigger for auto-update
CREATE TRIGGER update_signals_updated_at
    BEFORE UPDATE ON trading_signals
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

-- ============================================
-- PERFORMANCE VIEW
-- ============================================

-- Drop existing view if it exists
DROP VIEW IF EXISTS signal_performance;

-- Create view for signal performance
CREATE OR REPLACE VIEW signal_performance AS
SELECT 
    strategy,
    signal_type,
    COUNT(*) as total_signals,
    SUM(CASE WHEN status = 'HIT_TP' THEN 1 ELSE 0 END) as wins,
    SUM(CASE WHEN status = 'HIT_SL' THEN 1 ELSE 0 END) as losses,
    SUM(CASE WHEN status = 'ACTIVE' THEN 1 ELSE 0 END) as active,
    ROUND(
        100.0 * SUM(CASE WHEN status = 'HIT_TP' THEN 1 ELSE 0 END) / 
        NULLIF(SUM(CASE WHEN status IN ('HIT_TP', 'HIT_SL') THEN 1 ELSE 0 END), 0),
        2
    ) as win_rate,
    ROUND(AVG(CASE WHEN status = 'HIT_TP' THEN profit_loss_percent END), 2) as avg_win_percent,
    ROUND(AVG(CASE WHEN status = 'HIT_SL' THEN profit_loss_percent END), 2) as avg_loss_percent,
    ROUND(AVG(risk_reward), 2) as avg_risk_reward
FROM trading_signals
GROUP BY strategy, signal_type;

-- Grant view permissions
GRANT SELECT ON signal_performance TO anon;
GRANT SELECT ON signal_performance TO authenticated;
GRANT SELECT ON signal_performance TO service_role;

-- ============================================
-- VERIFICATION QUERIES
-- ============================================

-- Check table structure
SELECT 
    column_name, 
    data_type, 
    is_nullable,
    column_default
FROM information_schema.columns 
WHERE table_name = 'trading_signals'
ORDER BY ordinal_position;

-- Check if TP columns exist
SELECT 
    column_name,
    data_type
FROM information_schema.columns 
WHERE table_name = 'trading_signals'
AND column_name IN ('tp1', 'tp2', 'tp3');

-- Check indexes
SELECT 
    indexname,
    indexdef
FROM pg_indexes 
WHERE tablename = 'trading_signals';

-- Check policies
SELECT 
    policyname, 
    cmd, 
    roles 
FROM pg_policies 
WHERE tablename = 'trading_signals';

-- Check row count
SELECT COUNT(*) as total_signals FROM trading_signals;

-- Check recent signals with TP levels
SELECT 
    id,
    symbol,
    strategy,
    signal_type,
    entry_price,
    stop_loss,
    tp1,
    tp2,
    tp3,
    take_profit,
    risk_reward,
    status,
    created_at
FROM trading_signals
ORDER BY created_at DESC
LIMIT 5;

-- ============================================
-- TEST DATA (OPTIONAL - UNCOMMENT TO TEST)
-- ============================================

-- Insert test signal with TP levels
/*
INSERT INTO trading_signals (
    symbol, 
    strategy, 
    signal_type, 
    entry_price, 
    stop_loss, 
    take_profit,
    tp1,
    tp2,
    tp3,
    current_price, 
    risk_reward, 
    signal_time
)
VALUES (
    'BTCUSDT', 
    'session_trader', 
    'BUY', 
    50000.00, 
    49500.00, 
    55000.00,
    52000.00,
    53000.00,
    55000.00,
    50000.00, 
    10.0, 
    NOW()
);
*/

-- ============================================
-- CLEANUP OLD DATA (OPTIONAL)
-- ============================================

-- Delete signals older than 30 days (UNCOMMENT TO USE)
/*
DELETE FROM trading_signals 
WHERE created_at < NOW() - INTERVAL '30 days';
*/

-- Delete NONE signals (if any were accidentally saved)
/*
DELETE FROM trading_signals 
WHERE signal_type = 'NONE';
*/

-- ============================================
-- USER SETTINGS TABLE
-- ============================================

-- Create user_settings table for filter preferences
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

-- Enable RLS for user_settings
ALTER TABLE user_settings ENABLE ROW LEVEL SECURITY;

-- Drop existing policies if they exist
DROP POLICY IF EXISTS "Allow all operations for anon on user_settings" ON user_settings;
DROP POLICY IF EXISTS "Allow all operations for authenticated on user_settings" ON user_settings;
DROP POLICY IF EXISTS "Allow all operations for service_role on user_settings" ON user_settings;

-- Create policies
CREATE POLICY "Allow all operations for anon on user_settings"
ON user_settings
FOR ALL
TO anon
USING (true)
WITH CHECK (true);

CREATE POLICY "Allow all operations for authenticated on user_settings"
ON user_settings
FOR ALL
TO authenticated
USING (true)
WITH CHECK (true);

CREATE POLICY "Allow all operations for service_role on user_settings"
ON user_settings
FOR ALL
TO service_role
USING (true)
WITH CHECK (true);

-- Grant permissions
GRANT ALL ON user_settings TO anon;
GRANT ALL ON user_settings TO authenticated;
GRANT ALL ON user_settings TO service_role;

-- ============================================
-- SETUP COMPLETE!
-- ============================================

-- Summary of what was done:
-- ✅ Created/Updated trading_signals table with TP1, TP2, TP3 columns
-- ✅ Added indexes for performance
-- ✅ Configured Row Level Security (RLS)
-- ✅ Set up auto-update timestamp trigger
-- ✅ Created signal_performance view
-- ✅ Created user_settings table for filter persistence
-- ✅ Verified table structure and permissions

-- Next steps:
-- 1. Update your .env file with Supabase credentials
-- 2. Test signal generation: POST to /api/live-signal
-- 3. Check signals.html to see TP levels displayed
-- 4. Verify Telegram bot shows all three TP levels
-- 5. Filter settings now persist in database

SELECT 'Database setup complete! ✅' as status;
