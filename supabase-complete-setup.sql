-- ============================================
-- COMPLETE SUPABASE SETUP FOR TRADING SIGNALS
-- ============================================
-- Run this ENTIRE file in Supabase SQL Editor
-- This will DROP and RECREATE the table fresh
-- ============================================

-- STEP 1: Drop existing table and policies
DROP TABLE IF EXISTS trading_signals CASCADE;
DROP VIEW IF EXISTS signal_analytics CASCADE;

-- STEP 2: Create the trading_signals table
CREATE TABLE trading_signals (
    id BIGSERIAL PRIMARY KEY,
    signal_id TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Signal Details
    signal_type TEXT NOT NULL,              -- 'BUY' or 'SELL'
    symbol TEXT NOT NULL DEFAULT 'BTCUSDT',
    timeframe TEXT DEFAULT '15m',           -- '1m', '5m', '15m', '1h', '4h'
    entry_price DECIMAL(20, 8) NOT NULL,
    stop_loss DECIMAL(20, 8) NOT NULL,
    tp1 DECIMAL(20, 8) NOT NULL,
    tp2 DECIMAL(20, 8) NOT NULL,
    tp3 DECIMAL(20, 8) NOT NULL,
    strength INTEGER NOT NULL DEFAULT 70,
    
    -- Pattern Recognition
    pattern_type TEXT,
    pattern_confidence DECIMAL(5, 2),
    
    -- Session Info
    kill_zone TEXT,                         -- 'Asian', 'London', 'NewYork', 'Off-Hours'
    session_type TEXT,
    
    -- Status Tracking
    status TEXT DEFAULT 'pending',          -- 'pending', 'win', 'loss'
    exit_price DECIMAL(20, 8),
    exit_reason TEXT,                       -- 'TP1', 'TP2', 'TP3', 'Trailing Stop', 'Stop Loss'
    exit_time TIMESTAMP WITH TIME ZONE,
    
    -- Performance Metrics
    profit_percent DECIMAL(10, 4),
    profit_pips DECIMAL(10, 2),
    holding_time_minutes INTEGER,
    
    -- Trailing Stop
    trailing_stop_price DECIMAL(20, 8),
    trailing_stop_active BOOLEAN DEFAULT FALSE,
    highest_price DECIMAL(20, 8),
    lowest_price DECIMAL(20, 8),
    
    -- Live Tracking
    live_price DECIMAL(20, 8),
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Confluence
    confluence_score INTEGER DEFAULT 0
);

-- STEP 3: Create indexes for faster queries
CREATE INDEX idx_signals_status ON trading_signals(status);
CREATE INDEX idx_signals_symbol ON trading_signals(symbol);
CREATE INDEX idx_signals_timeframe ON trading_signals(timeframe);
CREATE INDEX idx_signals_created_at ON trading_signals(created_at DESC);
CREATE INDEX idx_signals_signal_type ON trading_signals(signal_type);

-- STEP 4: Enable Row Level Security
ALTER TABLE trading_signals ENABLE ROW LEVEL SECURITY;

-- STEP 5: Create policies for full public access
CREATE POLICY "Allow public select" ON trading_signals
    FOR SELECT USING (true);

CREATE POLICY "Allow public insert" ON trading_signals
    FOR INSERT WITH CHECK (true);

CREATE POLICY "Allow public update" ON trading_signals
    FOR UPDATE USING (true) WITH CHECK (true);

CREATE POLICY "Allow public delete" ON trading_signals
    FOR DELETE USING (true);

-- STEP 6: Grant permissions to anon and authenticated roles
GRANT ALL ON trading_signals TO anon;
GRANT ALL ON trading_signals TO authenticated;
GRANT USAGE, SELECT ON SEQUENCE trading_signals_id_seq TO anon;
GRANT USAGE, SELECT ON SEQUENCE trading_signals_id_seq TO authenticated;

-- STEP 7: Create analytics view
CREATE VIEW signal_analytics AS
SELECT 
    timeframe,
    signal_type,
    kill_zone,
    COUNT(*) as total_signals,
    SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) as wins,
    SUM(CASE WHEN status = 'loss' THEN 1 ELSE 0 END) as losses,
    SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END) as pending,
    ROUND(
        100.0 * SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) / 
        NULLIF(SUM(CASE WHEN status IN ('win', 'loss') THEN 1 ELSE 0 END), 0),
        2
    ) as win_rate,
    ROUND(AVG(CASE WHEN status = 'win' THEN profit_percent END), 2) as avg_win_percent,
    ROUND(AVG(CASE WHEN status = 'loss' THEN profit_percent END), 2) as avg_loss_percent
FROM trading_signals
GROUP BY timeframe, signal_type, kill_zone;

-- STEP 8: Create function to auto-update timestamp
CREATE OR REPLACE FUNCTION update_last_updated()
RETURNS TRIGGER AS $$
BEGIN
    NEW.last_updated = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- STEP 9: Create trigger for auto-update
DROP TRIGGER IF EXISTS update_signals_timestamp ON trading_signals;
CREATE TRIGGER update_signals_timestamp
    BEFORE UPDATE ON trading_signals
    FOR EACH ROW
    EXECUTE FUNCTION update_last_updated();

-- ============================================
-- VERIFICATION - Run these to confirm setup
-- ============================================

-- Check table exists
SELECT 'Table created successfully!' as status 
WHERE EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'trading_signals');

-- Check policies
SELECT policyname, cmd FROM pg_policies WHERE tablename = 'trading_signals';

-- Show table structure
SELECT column_name, data_type, is_nullable 
FROM information_schema.columns 
WHERE table_name = 'trading_signals'
ORDER BY ordinal_position;

-- ============================================
-- DONE! Your table is ready to use.
-- ============================================