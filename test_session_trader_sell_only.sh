#!/bin/bash

echo "🔍 Testing Session Trader - SELL Trades Only (30 days, 15m)"
echo "============================================================"
echo ""

# Test with SELL trades only
curl -s -X POST http://localhost:8080/api/v1/backtest/test-all-strategies \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "startBalance": 1000,
    "filterBuy": false,
    "filterSell": true
  }' | python3 << 'PYTHON_SCRIPT'
import json
import sys

try:
    data = json.load(sys.stdin)
    
    # Find session_trader results
    results = data.get('results', [])
    session_trader = None
    
    for strategy in results:
        if strategy.get('strategyName') == 'session_trader':
            session_trader = strategy
            break
    
    if session_trader:
        print("\n✅ SESSION TRADER - SELL TRADES ONLY RESULTS:")
        print("=" * 60)
        print(f"Win Rate:      {session_trader.get('winRate', 0):.2f}%")
        print(f"Sell Win Rate: {session_trader.get('sellWinRate', 0):.2f}%")
        print(f"Total Trades:  {session_trader.get('totalTrades', 0)}")
        print(f"Sell Trades:   {session_trader.get('sellTrades', 0)}")
        print(f"Sell Wins:     {session_trader.get('sellWins', 0)}")
        print(f"Return:        {session_trader.get('returnPercent', 0):.2f}%")
        print(f"Profit Factor: {session_trader.get('profitFactor', 0):.2f}")
        print(f"Max Drawdown:  {session_trader.get('maxDrawdown', 0):.2f}%")
        print("=" * 60)
        
        # Check if it's close to 99%
        sell_wr = session_trader.get('sellWinRate', 0)
        if sell_wr > 95:
            print(f"\n🎉 EXCELLENT! {sell_wr:.1f}% win rate achieved!")
        elif sell_wr > 80:
            print(f"\n✅ GOOD! {sell_wr:.1f}% win rate")
        elif sell_wr > 60:
            print(f"\n⚠️  MODERATE: {sell_wr:.1f}% win rate")
        else:
            print(f"\n❌ LOW: {sell_wr:.1f}% win rate - parameters may need adjustment")
    else:
        print("\n❌ Session Trader not found in results")
        print("\nAll strategies found:")
        for s in results:
            print(f"  - {s.get('strategyName')}")
    
except Exception as e:
    print(f"\n❌ Error: {e}")
    print("\nRaw response:")
    print(sys.stdin.read())
PYTHON_SCRIPT

echo ""
echo "============================================================"
echo "Test complete!"
echo ""
echo "To test in browser:"
echo "1. Open http://localhost:8080"
echo "2. UNCHECK '🟢 Buy Trades'"
echo "3. KEEP CHECKED '🔴 Sell Trades'"
echo "4. Click '🏆 Test All Strategies'"
echo ""
