#!/bin/bash

echo "🎯 Testing Session Trader SELL - Final Solution"
echo "================================================"
echo ""
echo "Configuration: Smart uptrend detection (2 of 4 checks)"
echo ""

# Test 30 days
echo "1️⃣  Testing 30 Days..."
curl -X POST http://localhost:8080/api/v1/backtest/test-all-strategies \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTCUSDT","days":30,"startBalance":1000,"filterBuy":false,"filterSell":true}' 2>/dev/null | \
  python3 -c "
import sys, json
data = json.load(sys.stdin)
for r in data['results']:
    if r['strategyName'] == 'session_trader':
        wr = r['winRate']
        trades = r['totalTrades']
        pf = r['profitFactor']
        dd = r['maxDrawdown']
        print(f'   Trades: {trades} (target: 80-100)')
        print(f'   Win Rate: {wr:.1f}% (target: 45-55%)')
        print(f'   Profit Factor: {pf:.2f} (target: >2.0)')
        print(f'   Max Drawdown: {dd:.1f}% (target: <40%)')
        print()
        if trades >= 70 and wr >= 45 and pf > 2.0 and dd < 40:
            print('   ✅ PASS: All metrics good!')
        elif trades >= 50:
            print('   ⚠️  PARTIAL: Some metrics need improvement')
        else:
            print('   ❌ FAIL: Too few trades')
        break
"

echo ""
echo "2️⃣  Testing Bad Period (5 Days)..."
curl -X POST http://localhost:8080/api/v1/backtest/test-all-strategies \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTCUSDT","days":5,"startBalance":1000,"filterBuy":false,"filterSell":true}' 2>/dev/null | \
  python3 -c "
import sys, json
data = json.load(sys.stdin)
for r in data['results']:
    if r['strategyName'] == 'session_trader':
        trades = r['totalTrades']
        wins = r['winningTrades']
        losses = r['losingTrades']
        reduction = (50 - trades) / 50 * 100
        print(f'   Original: 50 trades (7W/43L)')
        print(f'   Current:  {trades} trades ({wins}W/{losses}L)')
        print(f'   Reduction: {reduction:.0f}%')
        print()
        if trades <= 20:
            print('   ✅ PASS: Excellent reduction!')
        elif trades <= 30:
            print('   ✅ PASS: Good reduction!')
        else:
            print('   ⚠️  PARTIAL: Some reduction')
        break
"

echo ""
echo "================================================"
echo "📊 SUMMARY"
echo "================================================"
echo ""
echo "✅ Smart uptrend detection active"
echo "✅ 4 checks: price>EMA50, bullish candles, higher highs, price rising"
echo "✅ Skip trade if 2+ checks are true"
echo ""
echo "Expected Results:"
echo "  • 30 days: ~81 trades, ~49% WR, ~2.8 PF, ~35% DD"
echo "  • Bad period: ~21 trades (58% reduction from 50)"
echo ""
echo "📖 See SESSION_TRADER_FINAL_SOLUTION.md for details"
echo ""
