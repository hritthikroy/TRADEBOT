#!/bin/bash

echo "🎯 Testing Session Trader 99.6% SELL Win Rate Parameters"
echo "=========================================================="
echo ""
echo "📋 Parameters from Git Commit 79da2b7:"
echo "  - SELL Signal: EMA9 < EMA21 < EMA50 AND RSI(30-65)"
echo "  - Stop Loss: 1.0 ATR"
echo "  - TP1: 4.0 ATR"
echo "  - TP2: 6.0 ATR"
echo "  - TP3: 10.0 ATR"
echo ""
echo "✅ Status: ALREADY IMPLEMENTED in backend/unified_signal_generator.go"
echo ""
echo "🧪 Testing SELL trades only..."
echo ""

# Test via API
curl -X POST http://localhost:8080/api/v1/backtest/test-all-strategies \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "startBalance": 1000,
    "filterBuy": false,
    "filterSell": true
  }' 2>/dev/null | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    if 'results' in data:
        for result in data['results']:
            if result.get('strategyName') == 'session_trader':
                print('📊 Session Trader SELL-Only Results:')
                print(f\"   Win Rate: {result.get('winRate', 0):.1f}%\")
                print(f\"   Return: {result.get('returnPercent', 0):.1f}%\")
                print(f\"   Total Trades: {result.get('totalTrades', 0)}\")
                print(f\"   Profit Factor: {result.get('profitFactor', 0):.2f}\")
                print('')
                if result.get('winRate', 0) > 95:
                    print('✅ SUCCESS: Win rate > 95% (Expected: 99.6%)')
                else:
                    print('⚠️  Win rate lower than expected')
                break
    else:
        print('❌ No results found')
        print(json.dumps(data, indent=2))
except Exception as e:
    print(f'❌ Error: {e}')
"

echo ""
echo "📖 How to test in browser:"
echo "  1. Open http://localhost:8080"
echo "  2. UNCHECK '🟢 Buy Trades (Long)'"
echo "  3. KEEP CHECKED '🔴 Sell Trades (Short)'"
echo "  4. Click '🏆 Test All Strategies'"
echo "  5. Look for Session Trader results"
echo ""
echo "Expected Results:"
echo "  - Win Rate: ~99.6%"
echo "  - Return: ~3,200%"
echo "  - Trades: ~118"
echo ""
