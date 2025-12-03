-- ============================================
-- COMPLETE DATABASE SETUP FOR TRADING BOT
-- Run this in your Supabase SQL Editor
-- ============================================

-- ============================================
-- 1. TRADING SIGNALS TABLE
-- ============================================

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

-- Add TP columns if they don't exist (for existing tables)
ALTER TABLE trading_signals 
ADD COLUMN IF NOT EXISTS tp1 DECIMAL(20, 8),
ADD COLUMN IF NOT EXISTS tp2 DECIMAL(20, 8),
ADD COLUMN IF NOT EXISTS tp3 DECIMAL(20, 8);

-- ============================================
-- 2. USER SETTINGS TABLE
-- ============================================

-- Create user_settings table for filter preferences and strategy selections
CREATE TABLE IF NOT EXISTS user_settings (
    id INTEGER PRIMARY KEY DEFAULT 1,
    filter_buy BOOLEAN DEFAULT true,
    filter_sell BOOLEAN DEFAULT true,
    selected_strategies TEXT[] DEFAULT ARRAY['session_trader'],
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT single_row CHECK (id = 1)
);

-- Add selected_strategies column if it doesn't exist (for existing tables)
ALTER TABLE user_settings 
ADD COLUMN IF NOT EXISTS selected_strategies TEXT[] DEFAULT ARRAY['session_trader'];

-- Insert default settings
INSERT INTO user_settings (id, filter_buy, filter_sell, selected_strategies)
VALUES (1, true, true, ARRAY['session_trader'])
ON CONFLICT (id) DO NOTHING;

-- ============================================
-- 3. INDEXES FOR PERFORMANCE
-- ============================================

CREATE INDEX IF NOT EXISTS idx_signals_symbol ON trading_signals(symbol);
CREATE INDEX IF NOT EXISTS idx_signals_strategy ON trading_signals(strategy);
CREATE INDEX IF NOT EXISTS idx_signals_status ON trading_signals(status);
CREATE INDEX IF NOT EXISTS idx_signals_signal_type ON trading_signals(signal_type);
CREATE INDEX IF NOT EXISTS idx_signals_created_at ON trading_signals(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_signals_signal_time ON trading_signals(signal_time DESC);

-- ============================================
-- 4. ROW LEVEL SECURITY (RLS)
-- ============================================

-- Enable RLS for trading_signals
ALTER TABLE trading_signals ENABLE ROW LEVEL SECURITY;

-- Drop existing policies if they exist
DROP POLICY IF EXISTS "Allow all operations for anon" ON trading_signals;
DROP POLICY IF EXISTS "Allow all operations for authenticated" ON trading_signals;
DROP POLICY IF EXISTS "Allow all operations for service_role" ON trading_signals;

-- Create policies for trading_signals
CREATE POLICY "Allow all operations for anon"
ON trading_signals FOR ALL TO anon
USING (true) WITH CHECK (true);

CREATE POLICY "Allow all operations for authenticated"
ON trading_signals FOR ALL TO authenticated
USING (true) WITH CHECK (true);

CREATE POLICY "Allow all operations for service_role"
ON trading_signals FOR ALL TO service_role
USING (true) WITH CHECK (true);

-- Grant permissions for trading_signals
GRANT ALL ON trading_signals TO anon;
GRANT ALL ON trading_signals TO authenticated;
GRANT ALL ON trading_signals TO service_role;

-- Enable RLS for user_settings
ALTER TABLE user_settings ENABLE ROW LEVEL SECURITY;

-- Drop existing policies if they exist
DROP POLICY IF EXISTS "Allow all operations for anon on user_settings" ON user_settings;
DROP POLICY IF EXISTS "Allow all operations for authenticated on user_settings" ON user_settings;
DROP POLICY IF EXISTS "Allow all operations for service_role on user_settings" ON user_settings;

-- Create policies for user_settings
CREATE POLICY "Allow all operations for anon on user_settings"
ON user_settings FOR ALL TO anon
USING (true) WITH CHECK (true);

CREATE POLICY "Allow all operations for authenticated on user_settings"
ON user_settings FOR ALL TO authenticated
USING (true) WITH CHECK (true);

CREATE POLICY "Allow all operations for service_role on user_settings"
ON user_settings FOR ALL TO service_role
USING (true) WITH CHECK (true);

-- Grant permissions for user_settings
GRANT ALL ON user_settings TO anon;
GRANT ALL ON user_settings TO authenticated;
GRANT ALL ON user_settings TO service_role;

-- ============================================
-- 5. AUTO-UPDATE TIMESTAMP TRIGGER
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

-- Create trigger for trading_signals
CREATE TRIGGER update_signals_updated_at
    BEFORE UPDATE ON trading_signals
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

-- ============================================
-- 6. PERFORMANCE VIEW
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
-- 7. VERIFICATION
-- ============================================

-- Check tables exist
SELECT 'trading_signals' as table_name, COUNT(*) as row_count FROM trading_signals
UNION ALL
SELECT 'user_settings' as table_name, COUNT(*) as row_count FROM user_settings;

-- Check user_settings data
SELECT * FROM user_settings WHERE id = 1;

-- ============================================
-- SETUP COMPLETE!
-- ============================================

SELECT '✅ Database setup complete!' as status;
