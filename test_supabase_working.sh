#!/bin/bash

echo "🧪 Testing Supabase Connection"
echo "=============================="
echo ""

# Test connection
echo "📡 Testing connection to Supabase..."
RESPONSE=$(curl -s "https://elqhqhjevaizjoghiiss.supabase.co/rest/v1/trading_signals?limit=1" \
  -H "apikey: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImVscWhxaGpldmFpempvZ2hpaXNzIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NjQ3MDEwNTksImV4cCI6MjA4MDI3NzA1OX0.02-CgybOf7PiQSaQ-uJhojKob5Rw_2vkFdyurPNqLvA" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImVscWhxaGpldmFpempvZ2hpaXNzIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NjQ3MDEwNTksImV4cCI6MjA4MDI3NzA1OX0.02-CgybOf7PiQSaQ-uJhojKob5Rw_2vkFdyurPNqLvA")

if [ $? -eq 0 ]; then
    echo "✅ Connection successful!"
    echo "Response: $RESPONSE"
else
    echo "❌ Connection failed!"
    exit 1
fi

echo ""
echo "=============================="
echo "✅ Supabase is working!"
echo ""
echo "Next steps:"
echo "1. Start backend: cd backend && go run ."
echo "2. Generate signals from UI"
echo "3. Check Supabase dashboard to see signals"
echo ""
