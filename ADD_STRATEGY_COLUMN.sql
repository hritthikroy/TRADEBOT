-- Add selected_strategies column to user_settings table
ALTER TABLE user_settings 
ADD COLUMN IF NOT EXISTS selected_strategies TEXT[] DEFAULT ARRAY['session_trader'];

-- Update existing row to have default strategy
UPDATE user_settings 
SET selected_strategies = ARRAY['session_trader']
WHERE id = 1;

-- Verify the update
SELECT * FROM user_settings WHERE id = 1;
