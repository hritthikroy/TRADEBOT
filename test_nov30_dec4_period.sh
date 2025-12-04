#!/bin/bash

echo "🔍 Testing Session Trader SELL during Nov 30 - Dec 4 period"
echo "============================================================"
echo ""
echo "This period had 47 consecutive losses (100% loss rate)"
echo "Testing with current ultra-strict filters..."
echo ""

# Test with days parameter to cover Nov 30 - Dec 4 (about 5 days)
curl -X POST http://localhost:8080/api/v1/backtest/test-all-strategies \
  -H "Content-Type: application/json" \
  -d '{
    "symbol":"BTCUSDT",
    "days":5,
    "startBalance":1000,
    "filterBuy":false,
    "filterSell":true
  }' 2>/dev/null | python3 -c "
import sys, json
data = json.load(sys.stdin)
for r in data['results']:
    if r['strategyName'] == 'session_trader':
        wr = r['winRate']
        trades = r['totalTrades']
        wins = r['winningTrades']
        losses = r['losingTrades']
        ret = r['returnPercent']
        pf = r['profitFactor']
        
        print(f'🎯 Session Trader SELL - Last 5 Days')
        print(f'=' * 60)
        print(f'Win Rate: {wr:.1f}%')
        print(f'Trades: {trades} (Wins: {wins}, Losses: {losses})')
        print(f'Return: {ret:,.0f}%')
        print(f'Profit Factor: {pf:.2f}')
        print('')
        
        if trades == 0:
            print('✅ PERFECT! No trades = avoided bad period completely!')
        elif wr >= 70:
            print('✅ GOOD! Win rate >= 70% even in tough period')
        elif wr >= 50:
            print('✅ Profitable in tough period')
        elif wr > 0:
            print(f'⚠️  Some losses ({losses} losses vs {wins} wins)')
        else:
            print('❌ Still all losses - need stronger filters')
        break
"

echo ""
echo "Note: If trades = 0, filters successfully avoided the bad period!"
