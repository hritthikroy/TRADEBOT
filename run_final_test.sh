#!/bin/bash

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🎯 FINAL COMPREHENSIVE TEST - FINDING BEST PROFITABLE PARAMETERS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "⏱️  Testing Period: 90 days (good balance of speed & reliability)"
echo "💰 Start Balance: $1000"
echo "🎯 Testing: All 10 strategies × 3,990 parameters each"
echo "🔧 Fixes Applied:"
echo "   ✅ Optimizer tests real parameters (not hardcoded)"
echo "   ✅ Signal generation simplified (3/5 conditions)"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check if server is running
if ! curl -s http://localhost:8080/api/v1/health > /dev/null 2>&1; then
    echo "❌ Server is not running!"
    echo ""
    echo "Please start the server first:"
    echo "   cd backend && go run ."
    echo ""
    exit 1
fi

echo "✅ Server is running"
echo ""
echo "🚀 Starting optimization... (this will take 3-5 minutes)"
echo ""

# Run the optimization
curl -X POST http://localhost:8080/api/v1/backtest/world-class-optimize \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "days": 90,
    "startBalance": 1000
  }' 2>/dev/null | python3 -m json.tool > FINAL_OPTIMIZATION_RESULTS.json

if [ $? -eq 0 ]; then
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "✅ OPTIMIZATION COMPLETE!"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "📄 Results saved to: FINAL_OPTIMIZATION_RESULTS.json"
    echo ""
    
    # Analyze and display results
    python3 << 'EOF'
import json
import sys

try:
    with open('FINAL_OPTIMIZATION_RESULTS.json', 'r') as f:
        data = json.load(f)
    
    results = data.get('results', {})
    profitable = []
    unprofitable = []
    
    print("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
    print("📊 DETAILED RESULTS BY STRATEGY")
    print("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
    print()
    
    for strategy, result in sorted(results.items()):
        backtest = result.get('backtestResult')
        params = result.get('bestParams', {})
        
        print(f"🎯 {strategy.upper()}")
        
        if backtest and backtest.get('totalTrades', 0) > 0:
            ret = backtest.get('returnPercent', 0)
            wr = backtest.get('winRate', 0)
            pf = backtest.get('profitFactor', 0)
            trades = backtest.get('totalTrades', 0)
            dd = backtest.get('maxDrawdown', 0) * 100
            
            if ret > 0:
                print(f"   ✅ PROFITABLE")
                profitable.append({
                    'name': strategy,
                    'return': ret,
                    'winRate': wr,
                    'profitFactor': pf,
                    'trades': trades,
                    'drawdown': dd,
                    'params': params
                })
            else:
                print(f"   ⚠️  LOSING")
                unprofitable.append(strategy)
            
            print(f"   Return: {ret:.1f}% | Win Rate: {wr:.1f}% | Profit Factor: {pf:.2f}")
            print(f"   Trades: {trades} | Max Drawdown: {dd:.1f}%")
            print(f"   Parameters:")
            print(f"      Stop Loss: {params.get('StopATR', 0):.2f} ATR")
            print(f"      TP1: {params.get('TP1ATR', 0):.1f} ATR | TP2: {params.get('TP2ATR', 0):.1f} ATR | TP3: {params.get('TP3ATR', 0):.1f} ATR")
            print(f"      Risk per Trade: {params.get('RiskPercent', 0):.1f}%")
        else:
            print(f"   ❌ NO TRADES GENERATED")
            unprofitable.append(strategy)
        
        print()
    
    print("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
    print("🏆 PROFITABLE STRATEGIES RANKED")
    print("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
    print()
    
    if profitable:
        # Sort by return
        profitable.sort(key=lambda x: x['return'], reverse=True)
        
        for i, s in enumerate(profitable, 1):
            print(f"{i}. {s['name'].upper()}")
            print(f"   💰 Return: {s['return']:.1f}% | Win Rate: {s['winRate']:.1f}% | Profit Factor: {s['profitFactor']:.2f}")
            print(f"   📊 Trades: {s['trades']} | Max Drawdown: {s['drawdown']:.1f}%")
            params = s['params']
            print(f"   🎯 Best Parameters:")
            print(f"      Stop: {params.get('StopATR', 0):.2f} | TP1: {params.get('TP1ATR', 0):.1f} | TP2: {params.get('TP2ATR', 0):.1f} | TP3: {params.get('TP3ATR', 0):.1f} | Risk: {params.get('RiskPercent', 0):.1f}%")
            print()
        
        print("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
        print(f"🎉 SUCCESS! Found {len(profitable)} profitable strategies!")
        print("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
        print()
        
        # Calculate portfolio stats
        total_return = sum(s['return'] for s in profitable)
        avg_return = total_return / len(profitable)
        avg_wr = sum(s['winRate'] for s in profitable) / len(profitable)
        avg_pf = sum(s['profitFactor'] for s in profitable) / len(profitable)
        
        print("📈 PORTFOLIO STATISTICS (if using all profitable strategies):")
        print(f"   Average Return: {avg_return:.1f}%")
        print(f"   Average Win Rate: {avg_wr:.1f}%")
        print(f"   Average Profit Factor: {avg_pf:.2f}")
        print()
        
    else:
        print("⚠️  No profitable strategies found in this test period.")
        print()
        print("Possible reasons:")
        print("   • Market conditions in this 90-day period weren't favorable")
        print("   • Signal generation still too strict for some strategies")
        print("   • Try different time period or symbol")
        print()
    
    if unprofitable:
        print("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
        print(f"⚠️  STRATEGIES NEEDING IMPROVEMENT ({len(unprofitable)}):")
        print("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
        for s in unprofitable:
            print(f"   • {s}")
        print()
        
except Exception as e:
    print(f"Error analyzing results: {e}")
    import traceback
    traceback.print_exc()
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
