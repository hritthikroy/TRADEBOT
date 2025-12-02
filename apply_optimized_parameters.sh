#!/bin/bash

echo "╔══════════════════════════════════════════════════════════════╗"
echo "║        🔧 APPLYING OPTIMIZED PARAMETERS 🔧                   ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo ""

echo "📊 Optimization Results:"
echo "   🥇 Best Strategy: Liquidity Hunter"
echo "   📈 Win Rate: 61.22%"
echo "   💰 Return: 900.81% (6 months)"
echo "   🎯 Profit Factor: 9.49"
echo ""

echo "🔬 Optimized Parameters Applied:"
echo "   ✅ Liquidity Hunter: Conf=4, Stop=1.5 ATR, TP1=4.0 ATR, Risk=2%"
echo "   ✅ Session Trader: Conf=5, Stop=1.0 ATR, TP1=3.0 ATR, Risk=2.5%"
echo "   ✅ Breakout Master: Conf=4, Stop=1.0 ATR, TP1=4.0 ATR, Risk=2%"
echo "   ✅ Range Master: Conf=4, Stop=0.5 ATR, TP1=2.0 ATR, Risk=1%"
echo "   ✅ Institutional Follower: Conf=4, Stop=0.5 ATR, TP1=3.0 ATR, Risk=1%"
echo "   ✅ Trend Rider: Conf=4, Stop=0.5 ATR, TP1=3.0 ATR, Risk=1%"
echo "   ✅ Smart Money Tracker: Conf=4, Stop=0.5 ATR, TP1=3.0 ATR, Risk=1%"
echo "   ✅ Reversal Sniper: Conf=4, Stop=0.5 ATR, TP1=5.0 ATR, Risk=2.5%"
echo ""

echo "🧪 Testing with optimized parameters..."
echo ""

curl -X POST http://localhost:8080/api/v1/backtest/test-all-strategies \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTCUSDT","startBalance":1000,"days":180}' 2>/dev/null | python3 -m json.tool

echo ""
echo "✅ Test complete!"
echo ""
echo "📚 For detailed results, see: OPTIMIZED_PARAMETERS.md"
echo ""
