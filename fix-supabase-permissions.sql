-- Fix Supabase Permissions for Public Access
-- Run this in Supabase SQL Editor

-- Drop existing policies
DROP POLICY IF EXISTS "Enable all access for authenticated users" ON trading_signals;

-- Create new policy that allows all operations without authentication
CREATE POLICY "Enable all access for everyone" ON trading_signals
    FOR ALL 
    USING (true)
    WITH CHECK (true);

-- Verify RLS is enabled
ALTER TABLE trading_signals ENABLE ROW LEVEL SECURITY;

-- Test: This should work now
SELECT * FROM trading_signals LIMIT 1;
