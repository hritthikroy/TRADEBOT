# ⚡ QUICK COMMANDS - COPY & PASTE

## 🎯 MOST USED COMMANDS

---

## 📊 CHECK PAPER TRADING STATS

```bash
curl -s http://localhost:8080/api/v1/paper-trading/stats | jq '.stats | {
  totalTrades,
  winRate,
  currentBalance,
  netProfit,
  maxDrawdown
}'
```

---

## 🔍 GET CURRENT SIGNAL

```bash
curl -s -X POST http://localhost:8080/api/v1/backtest/live-signal \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTCUSDT","interval":"15m","strategy":"session_trader"}' | jq '.'
```

---

## 📈 VIEW ALL TRADES

```bash
curl -s http://localhost:8080/api/v1/paper-trading/trades | jq '.trades | map({
  id,
  signal,
  entry,
  exitPrice,
  profit,
  status
})'
```

---

## 🚀 START AUTO TRADING

```bash
curl -X POST http://localhost:8080/api/v1/paper-trading/start-auto
```

---

## 🛑 STOP AUTO TRADING

```bash
curl -X POST http://localhost:8080/api/v1/paper-trading/stop-auto
```

---

## 🔄 RESET PAPER TRADING

```bash
curl -X POST http://localhost:8080/api/v1/paper-trading/reset
```

---

## 🧪 TEST EVERYTHING

```bash
./test_paper_trading_api.sh
```

---

## 📱 QUICK STATUS CHECK

```bash
curl -s http://localhost:8080/api/v1/paper-trading/stats | jq -r '
  "Total Trades: \(.stats.totalTrades)",
  "Win Rate: \(.stats.winRate)%",
  "Balance: $\(.stats.currentBalance)",
  "Profit: $\(.stats.netProfit)",
  "Drawdown: \(.stats.maxDrawdown)%"
'
```

---

## 🎯 SIGNAL ONLY (QUICK)

```bash
curl -s -X POST http://localhost:8080/api/v1/backtest/live-signal \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTCUSDT","interval":"15m","strategy":"session_trader"}' | jq -r '.signal'
```

---

## 💰 CURRENT BALANCE

```bash
curl -s http://localhost:8080/api/v1/paper-trading/stats | jq -r '"Balance: $\(.stats.currentBalance)"'
```

---

## 📊 WIN RATE ONLY

```bash
curl -s http://localhost:8080/api/v1/paper-trading/stats | jq -r '"Win Rate: \(.stats.winRate)%"'
```

---

## 🔢 TRADE COUNT

```bash
curl -s http://localhost:8080/api/v1/paper-trading/stats | jq -r '"Total Trades: \(.stats.totalTrades)"'
```

---

## 📈 PROFIT/LOSS

```bash
curl -s http://localhost:8080/api/v1/paper-trading/stats | jq -r '"Net Profit: $\(.stats.netProfit)"'
```

---

## 🎯 COMPLETE SUMMARY

```bash
echo "=== PAPER TRADING SUMMARY ===" && \
curl -s http://localhost:8080/api/v1/paper-trading/stats | jq -r '
  "Total Trades: \(.stats.totalTrades)",
  "Open Trades: \(.stats.openTrades)",
  "Closed Trades: \(.stats.closedTrades)",
  "Winning Trades: \(.stats.winningTrades)",
  "Losing Trades: \(.stats.losingTrades)",
  "Win Rate: \(.stats.winRate)%",
  "Start Balance: $\(.stats.startBalance)",
  "Current Balance: $\(.stats.currentBalance)",
  "Net Profit: $\(.stats.netProfit)",
  "Return: \(.stats.returnPercent)%",
  "Max Drawdown: \(.stats.maxDrawdown)%",
  "Profit Factor: \(.stats.profitFactor)"
'
```

---

## 🔄 DAILY CHECK ROUTINE

```bash
# Run this once per day
echo "📊 Daily Paper Trading Report" && \
echo "=============================" && \
curl -s http://localhost:8080/api/v1/paper-trading/stats | jq -r '
  "Date: \(now | strftime("%Y-%m-%d %H:%M"))",
  "",
  "📈 Performance:",
  "  Total Trades: \(.stats.totalTrades)",
  "  Win Rate: \(.stats.winRate)%",
  "  Current Balance: $\(.stats.currentBalance)",
  "  Net Profit: $\(.stats.netProfit) (\(.stats.returnPercent)%)",
  "",
  "⚠️  Risk:",
  "  Max Drawdown: \(.stats.maxDrawdown)%",
  "",
  "✅ Status: \(if .stats.winRate >= 75 then "EXCELLENT" elif .stats.winRate >= 60 then "GOOD" else "NEEDS IMPROVEMENT" end)"
'
```

---

## 🎯 SAVE TO FILE

```bash
# Save stats to file
curl -s http://localhost:8080/api/v1/paper-trading/stats | jq '.stats' > paper_trading_stats_$(date +%Y%m%d).json
echo "Stats saved to paper_trading_stats_$(date +%Y%m%d).json"
```

---

## 📱 MOBILE FRIENDLY (SHORT)

```bash
# Quick check on mobile
curl -s http://localhost:8080/api/v1/paper-trading/stats | jq -r '"\(.stats.totalTrades) trades | \(.stats.winRate)% WR | $\(.stats.currentBalance)"'
```

---

## 🔍 CHECK IF RUNNING

```bash
# Check if auto trading is active
if curl -s http://localhost:8080/api/v1/health | grep -q "healthy"; then
  echo "✅ Backend is running"
  curl -s http://localhost:8080/api/v1/paper-trading/stats | jq -r '"Total Trades: \(.stats.totalTrades)"'
else
  echo "❌ Backend is not running"
fi
```

---

## 🎯 ALIASES (ADD TO ~/.zshrc)

```bash
# Add these to your ~/.zshrc for quick access
alias pts='curl -s http://localhost:8080/api/v1/paper-trading/stats | jq ".stats"'
alias ptq='curl -s http://localhost:8080/api/v1/paper-trading/stats | jq -r "\(.stats.totalTrades) trades | \(.stats.winRate)% WR | $\(.stats.currentBalance)"'
alias sig='curl -s -X POST http://localhost:8080/api/v1/backtest/live-signal -H "Content-Type: application/json" -d "{\"symbol\":\"BTCUSDT\",\"interval\":\"15m\",\"strategy\":\"session_trader\"}" | jq "."'
alias sigq='curl -s -X POST http://localhost:8080/api/v1/backtest/live-signal -H "Content-Type: application/json" -d "{\"symbol\":\"BTCUSDT\",\"interval\":\"15m\",\"strategy\":\"session_trader\"}" | jq -r ".signal"'
```

**Then use:**
- `pts` - Full stats
- `ptq` - Quick stats
- `sig` - Get signal
- `sigq` - Signal only

---

## 📞 TROUBLESHOOTING

**Backend not responding?**
```bash
curl -s http://localhost:8080/api/v1/health
```

**Check backend logs:**
```bash
cd backend && go run . 2>&1 | tail -20
```

**Restart backend:**
```bash
# Stop current process (Ctrl+C)
# Then start again
cd backend && go run .
```

---

## ✅ COPY THESE TO YOUR NOTES

**Daily check:**
```bash
curl -s http://localhost:8080/api/v1/paper-trading/stats | jq '.stats | {totalTrades, winRate, currentBalance, netProfit}'
```

**Get signal:**
```bash
curl -s -X POST http://localhost:8080/api/v1/backtest/live-signal -H "Content-Type: application/json" -d '{"symbol":"BTCUSDT","interval":"15m","strategy":"session_trader"}' | jq '.'
```

**TradingView:**
https://www.tradingview.com/chart/

---

## 🎯 DONE!

Save this file for quick reference! 📋
