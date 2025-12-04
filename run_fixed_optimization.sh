#!/bin/bash

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🔥 RUNNING FIXED OPTIMIZER - NOW TESTS REAL PARAMETERS!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "⏱️  Testing Period: 180 days"
echo "💰 Start Balance: $1000"
echo "🎯 Goal: Find PROFITABLE parameters for all 10 strategies"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "🚀 Starting optimization..."
echo ""

# Run the optimization
curl -X POST http://localhost:8080/api/v1/backtest/world-class-optimize \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "days": 180,
    "startBalance": 1000
  }' 2>/dev/null | python3 -m json.tool > WORLD_CLASS_OPTIMIZATION_RESULTS_FIXED.json

if [ $? -eq 0 ]; then
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "✅ OPTIMIZATION COMPLETE!"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "📄 Results saved to: WORLD_CLASS_OPTIMIZATION_RESULTS_FIXED.json"
    echo ""
    echo "📊 PROFITABLE STRATEGIES:"
    echo ""
    
    # Extract and display profitable strategies
    python3 << 'EOF'
import json
import sys

try:
    with open('WORLD_CLASS_OPTIMIZATION_RESULTS_FIXED.json', 'r') as f:
        data = json.load(f)
    
    results = data.get('results', {})
    profitable = []
    
    for strategy, result in results.items():
        backtest = result.get('backtestResult')
        if backtest and backtest.get('returnPercent', 0) > 0:
            profitable.append({
                'name': strategy,
                'return': backtest.get('returnPercent', 0),
                'winRate': backtest.get('winRate', 0),
                'profitFactor': backtest.get('profitFactor', 0),
                'trades': backtest.get('totalTrades', 0),
                'params': result.get('bestParams', {})
            })
    
    if profitable:
        # Sort by return
        profitable.sort(key=lambda x: x['return'], reverse=True)
        
        for i, s in enumerate(profitable, 1):
            print(f"{i}. {s['name'].upper()}")
            print(f"   Return: {s['return']:.1f}% | WR: {s['winRate']:.1f}% | PF: {s['profitFactor']:.2f} | Trades: {s['trades']}")
            params = s['params']
            print(f"   Stop: {params.get('StopATR', 0):.2f} | TP1: {params.get('TP1ATR', 0):.1f} | TP2: {params.get('TP2ATR', 0):.1f} | TP3: {params.get('TP3ATR', 0):.1f} | Risk: {params.get('RiskPercent', 0):.1f}%")
            print()
        
        print(f"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
        print(f"🎉 Found {len(profitable)} profitable strategies!")
        print(f"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
    else:
        print("⚠️  No profitable strategies found.")
        print("   Try:")
        print("   - Different time period (60-90 days)")
        print("   - Different symbol (ETHUSDT, BNBUSDT)")
        print("   - Check server logs for errors")
        
except Exception as e:
    print(f"Error reading results: {e}")
    sys.exit(1)
EOF

else
    echo ""
    echo "❌ Optimization failed!"
    echo ""
    echo "Please check:"
    echo "1. Server is running: cd backend && go run ."
    echo "2. Server logs for errors"
    echo ""
fi
