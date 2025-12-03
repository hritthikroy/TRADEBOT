-- Fix user_settings to enable both filters by default
-- Run this in your Supabase SQL Editor

UPDATE user_settings 
SET 
    filter_buy = true,
    filter_sell = true,
    updated_at = NOW()
WHERE id = 1;

-- Verify the update
SELECT * FROM user_settings WHERE id = 1;
