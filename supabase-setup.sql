-- Trading Signals Table
CREATE TABLE trading_signals (
    id BIGSERIAL PRIMARY KEY,
    signal_id TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Signal Details
    signal_type TEXT NOT NULL, -- 'BUY' or 'SELL'
    symbol TEXT NOT NULL,
    entry_price DECIMAL(20, 8) NOT NULL,
    stop_loss DECIMAL(20, 8) NOT NULL,
    tp1 DECIMAL(20, 8) NOT NULL,
    tp2 DECIMAL(20, 8) NOT NULL,
    tp3 DECIMAL(20, 8) NOT NULL,
    strength INTEGER NOT NULL,
    
    -- Pattern Recognition
    pattern_type TEXT,
    pattern_confidence DECIMAL(5, 2),
    
    -- Kill Zone Info
    kill_zone TEXT, -- 'Asian', 'London', 'NewYork', 'Off-Hours'
    session_type TEXT, -- 'Asian', 'London', 'NewYork', 'Overlap', 'Off-Hours'
    
    -- Status Tracking
    status TEXT DEFAULT 'pending', -- 'pending', 'win', 'loss'
    exit_price DECIMAL(20, 8),
    exit_reason TEXT, -- 'TP1', 'TP2', 'TP3', 'Trailing Stop', 'Stop Loss'
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
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for faster queries
CREATE INDEX idx_signals_status ON trading_signals(status);
CREATE INDEX idx_signals_symbol ON trading_signals(symbol);
CREATE INDEX idx_signals_created_at ON trading_signals(created_at DESC);
CREATE INDEX idx_signals_kill_zone ON trading_signals(kill_zone);
CREATE INDEX idx_signals_pattern ON trading_signals(pattern_type);
CREATE INDEX idx_signals_signal_type ON trading_signals(signal_type);

-- Enable Row Level Security
ALTER TABLE trading_signals ENABLE ROW LEVEL SECURITY;

-- Policy: Allow all operations (you can restrict this later)
CREATE POLICY "Enable all access for authenticated users" ON trading_signals
    FOR ALL USING (true);

-- Create a view for analytics
CREATE VIEW signal_analytics AS
SELECT 
    kill_zone,
    signal_type,
    pattern_type,
    COUNT(*) as total_signals,
    SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) as wins,
    SUM(CASE WHEN status = 'loss' THEN 1 ELSE 0 END) as losses,
    ROUND(AVG(CASE WHEN status IN ('win', 'loss') THEN profit_percent END), 2) as avg_profit_percent,
    ROUND(AVG(CASE WHEN status = 'win' THEN profit_percent END), 2) as avg_win_percent,
    ROUND(AVG(CASE WHEN status = 'loss' THEN profit_percent END), 2) as avg_loss_percent,
    ROUND(AVG(holding_time_minutes), 0) as avg_holding_minutes,
    ROUND(
        100.0 * SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) / 
        NULLIF(SUM(CASE WHEN status IN ('win', 'loss') THEN 1 ELSE 0 END), 0),
        2
    ) as win_rate
FROM trading_signals
WHERE status IN ('win', 'loss', 'pending')
GROUP BY kill_zone, signal_type, pattern_type;

-- Function to update last_updated timestamp
CREATE OR REPLACE FUNCTION update_last_updated()
RETURNS TRIGGER AS $$
BEGIN
    NEW.last_updated = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to auto-update last_updated
CREATE TRIGGER update_signals_timestamp
    BEFORE UPDATE ON trading_signals
    FOR EACH ROW
    EXECUTE FUNCTION update_last_updated();
