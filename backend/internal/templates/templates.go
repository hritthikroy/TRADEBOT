package templates

import (
	"html/template"
	"time"
)

// TemplateData holds data for template rendering
type TemplateData struct {
	Title       string
	CurrentTime time.Time
	Data        interface{}
}

// GetIndexTemplate returns the main dashboard template
func GetIndexTemplate() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }
        .container { max-width: 1400px; margin: 0 auto; }
        .header { text-align: center; color: white; margin-bottom: 30px; }
        .header h1 { font-size: 2.5em; margin-bottom: 10px; }
        .header p { font-size: 1.2em; opacity: 0.9; }
        .card {
            background: white;
            border-radius: 12px;
            padding: 24px;
            margin-bottom: 20px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        }
        .form-group { margin-bottom: 16px; }
        .form-group label {
            display: block;
            margin-bottom: 8px;
            font-weight: 600;
            color: #555;
        }
        .form-group input, .form-group select {
            width: 100%;
            padding: 10px;
            border: 2px solid #ddd;
            border-radius: 6px;
            font-size: 14px;
        }
        .btn {
            padding: 12px 24px;
            background: #667eea;
            color: white;
            border: none;
            border-radius: 6px;
            cursor: pointer;
            font-size: 16px;
            font-weight: 600;
            transition: all 0.3s;
        }
        .btn:hover { background: #5568d3; transform: translateY(-2px); }
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 16px;
            margin-top: 20px;
        }
        .stat-card {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 20px;
            border-radius: 8px;
            text-align: center;
        }
        .stat-card h3 { font-size: 14px; opacity: 0.9; margin-bottom: 8px; }
        .stat-card p { font-size: 28px; font-weight: bold; }
        .hidden { display: none; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚀 Trading Bot</h1>
            <p>Powered by Go Backend - Server-Side Rendered</p>
            <p style="font-size: 0.9em; opacity: 0.8;">Generated at: {{.CurrentTime.Format "2006-01-02 15:04:05"}}</p>
        </div>
        
        <div class="card">
            <h2>📊 Backtest Configuration</h2>
            <form action="/backtest/run" method="POST">
                <div class="form-group">
                    <label>Symbol</label>
                    <input type="text" name="symbol" value="BTCUSDT" required>
                </div>
                <div class="form-group">
                    <label>Interval</label>
                    <select name="interval">
                        <option value="1m">1 Minute</option>
                        <option value="5m">5 Minutes</option>
                        <option value="15m" selected>15 Minutes</option>
                        <option value="1h">1 Hour</option>
                        <option value="4h">4 Hours</option>
                        <option value="1d">1 Day</option>
                    </select>
                </div>
                <div class="form-group">
                    <label>Days to Test</label>
                    <input type="number" name="days" value="30" min="1" max="365" required>
                </div>
                <div class="form-group">
                    <label>Starting Balance ($)</label>
                    <input type="number" name="balance" value="500" min="100" required>
                </div>
                <div class="form-group">
                    <label>Risk Per Trade (%)</label>
                    <input type="number" name="risk" value="2" min="0.5" max="10" step="0.5" required>
                </div>
                <button type="submit" class="btn">Run Backtest</button>
            </form>
        </div>
    </div>
</body>
</html>`
}

// GetBacktestResultsTemplate returns the backtest results template
func GetBacktestResultsTemplate() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Backtest Results</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }
        .container { max-width: 1600px; margin: 0 auto; }
        .header { text-align: center; color: white; margin-bottom: 30px; }
        .header h1 { font-size: 2.5em; margin-bottom: 10px; }
        .card {
            background: white;
            border-radius: 12px;
            padding: 24px;
            margin-bottom: 20px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        }
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 16px;
        }
        .stat-card {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 20px;
            border-radius: 8px;
            text-align: center;
        }
        .stat-card h3 { font-size: 14px; opacity: 0.9; margin-bottom: 8px; }
        .stat-card p { font-size: 28px; font-weight: bold; }
        .profit { color: #4CAF50; }
        .loss { color: #f44336; }
        .table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 20px;
        }
        .table th, .table td {
            padding: 12px;
            text-align: left;
            border-bottom: 1px solid #ddd;
        }
        .table th {
            background: #f5f5f5;
            font-weight: 600;
        }
        .table tr:hover { background: #f9f9f9; }
        .btn {
            padding: 12px 24px;
            background: #667eea;
            color: white;
            border: none;
            border-radius: 6px;
            cursor: pointer;
            font-size: 16px;
            font-weight: 600;
            text-decoration: none;
            display: inline-block;
            margin-right: 10px;
        }
        .btn:hover { background: #5568d3; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>📈 Backtest Results</h1>
            <p>Executed in {{.Data.Duration}}</p>
        </div>
        
        <div class="card">
            <h2>Statistics</h2>
            <div class="stats-grid">
                <div class="stat-card">
                    <h3>Total Trades</h3>
                    <p>{{.Data.TotalTrades}}</p>
                </div>
                <div class="stat-card">
                    <h3>Win Rate</h3>
                    <p>{{printf "%.1f" .Data.WinRate}}%</p>
                </div>
                <div class="stat-card">
                    <h3>Return</h3>
                    <p class="{{if ge .Data.ReturnPercent 0.0}}profit{{else}}loss{{end}}">
                        {{printf "%.1f" .Data.ReturnPercent}}%
                    </p>
                </div>
                <div class="stat-card">
                    <h3>Final Balance</h3>
                    <p>${{printf "%.2f" .Data.FinalBalance}}</p>
                </div>
                <div class="stat-card">
                    <h3>Profit Factor</h3>
                    <p>{{printf "%.2f" .Data.ProfitFactor}}</p>
                </div>
                <div class="stat-card">
                    <h3>Max Drawdown</h3>
                    <p>{{printf "%.1f" (mul .Data.MaxDrawdown 100)}}%</p>
                </div>
            </div>
        </div>
        
        <div class="card">
            <h2>Trades ({{len .Data.Trades}} total)</h2>
            <div style="overflow-x: auto;">
                <table class="table">
                    <thead>
                        <tr>
                            <th>#</th>
                            <th>Type</th>
                            <th>Entry</th>
                            <th>Exit</th>
                            <th>Exit Reason</th>
                            <th>Profit</th>
                            <th>Profit %</th>
                            <th>RR</th>
                            <th>Balance</th>
                        </tr>
                    </thead>
                    <tbody>
                        {{range $index, $trade := .Data.Trades}}
                        <tr>
                            <td>{{add $index 1}}</td>
                            <td><strong>{{$trade.Type}}</strong></td>
                            <td>${{printf "%.2f" $trade.Entry}}</td>
                            <td>${{printf "%.2f" $trade.Exit}}</td>
                            <td>{{$trade.ExitReason}}</td>
                            <td class="{{if ge $trade.Profit 0.0}}profit{{else}}loss{{end}}">
                                ${{printf "%.2f" $trade.Profit}}
                            </td>
                            <td class="{{if ge $trade.ProfitPercent 0.0}}profit{{else}}loss{{end}}">
                                {{printf "%.1f" $trade.ProfitPercent}}%
                            </td>
                            <td>{{printf "%.2f" $trade.RR}}</td>
                            <td>${{printf "%.2f" $trade.BalanceAfter}}</td>
                        </tr>
                        {{end}}
                    </tbody>
                </table>
            </div>
        </div>
        
        <div style="text-align: center; margin-top: 20px;">
            <a href="/" class="btn">Run Another Backtest</a>
            <a href="/signals/live" class="btn">View Live Signals</a>
        </div>
    </div>
</body>
</html>`
}

// GetLiveSignalsTemplate returns the live signals template with WebSocket
func GetLiveSignalsTemplate() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Live Signals - Real-time</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }
        .container { max-width: 1600px; margin: 0 auto; }
        .header { text-align: center; color: white; margin-bottom: 30px; }
        .header h1 { font-size: 2.5em; margin-bottom: 10px; }
        .status {
            display: inline-block;
            padding: 4px 12px;
            border-radius: 12px;
            font-size: 14px;
            font-weight: 600;
        }
        .status.connected { background: #4CAF50; color: white; }
        .status.disconnected { background: #f44336; color: white; }
        .card {
            background: white;
            border-radius: 12px;
            padding: 24px;
            margin-bottom: 20px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        }
        .signals-grid {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
            gap: 20px;
        }
        .signal-card {
            border: 2px solid #ddd;
            border-radius: 8px;
            padding: 20px;
            background: white;
            animation: slideIn 0.3s ease-out;
        }
        @keyframes slideIn {
            from { opacity: 0; transform: translateY(-10px); }
            to { opacity: 1; transform: translateY(0); }
        }
        .signal-card.buy { border-left: 4px solid #4CAF50; }
        .signal-card.sell { border-left: 4px solid #f44336; }
        .signal-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 16px;
        }
        .signal-type {
            font-size: 20px;
            font-weight: bold;
        }
        .signal-type.buy { color: #4CAF50; }
        .signal-type.sell { color: #f44336; }
        .signal-info { margin-bottom: 12px; }
        .signal-info label {
            font-size: 12px;
            color: #666;
            display: block;
            margin-bottom: 4px;
        }
        .signal-info value {
            font-size: 16px;
            font-weight: 600;
            color: #333;
        }
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 16px;
            margin-bottom: 20px;
        }
        .stat-card {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 20px;
            border-radius: 8px;
            text-align: center;
        }
        .stat-card h3 { font-size: 14px; opacity: 0.9; margin-bottom: 8px; }
        .stat-card p { font-size: 28px; font-weight: bold; }
        .empty {
            text-align: center;
            padding: 60px;
            color: #999;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>📡 Live Signal Tracker</h1>
            <p>Real-time Updates via WebSocket</p>
            <p style="margin-top: 10px;">
                <span id="status" class="status disconnected">Connecting...</span>
            </p>
        </div>
        
        <div class="stats-grid">
            <div class="stat-card">
                <h3>Active Signals</h3>
                <p id="activeCount">0</p>
            </div>
            <div class="stat-card">
                <h3>Total Today</h3>
                <p id="totalCount">0</p>
            </div>
            <div class="stat-card">
                <h3>Win Rate</h3>
                <p id="winRate">0%</p>
            </div>
            <div class="stat-card">
                <h3>Updates/sec</h3>
                <p id="updateRate">0</p>
            </div>
        </div>
        
        <div class="card">
            <h2>Active Signals</h2>
            <div id="signalsContainer" class="signals-grid">
                <div class="empty">
                    <h3>⏳ Waiting for signals...</h3>
                    <p>Signals will appear here in real-time</p>
                </div>
            </div>
        </div>
    </div>
    
    <script>
        let ws;
        let updateCount = 0;
        let lastUpdateTime = Date.now();
        
        function connect() {
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            ws = new WebSocket(protocol + '//' + window.location.host + '/ws/signals');
            
            ws.onopen = () => {
                document.getElementById('status').textContent = 'Connected';
                document.getElementById('status').className = 'status connected';
                console.log('WebSocket connected');
            };
            
            ws.onmessage = (event) => {
                const data = JSON.parse(event.data);
                updateSignals(data);
                updateStats();
                
                // Calculate update rate
                updateCount++;
                const now = Date.now();
                if (now - lastUpdateTime >= 1000) {
                    document.getElementById('updateRate').textContent = updateCount;
                    updateCount = 0;
                    lastUpdateTime = now;
                }
            };
            
            ws.onclose = () => {
                document.getElementById('status').textContent = 'Disconnected';
                document.getElementById('status').className = 'status disconnected';
                console.log('WebSocket disconnected, reconnecting...');
                setTimeout(connect, 3000);
            };
            
            ws.onerror = (error) => {
                console.error('WebSocket error:', error);
            };
        }
        
        function updateSignals(signals) {
            const container = document.getElementById('signalsContainer');
            
            if (!signals || signals.length === 0) {
                container.innerHTML = '<div class="empty"><h3>📭 No active signals</h3></div>';
                return;
            }
            
            container.innerHTML = signals.map(signal => ` + "`" + `
                <div class="signal-card ${signal.type.toLowerCase()}">
                    <div class="signal-header">
                        <span class="signal-type ${signal.type.toLowerCase()}">${signal.type}</span>
                        <span style="font-size: 12px; color: #666;">${signal.timeframe}</span>
                    </div>
                    <div class="signal-info">
                        <label>Entry</label>
                        <value>$${signal.entry.toFixed(2)}</value>
                    </div>
                    <div class="signal-info">
                        <label>Stop Loss</label>
                        <value>$${signal.stopLoss.toFixed(2)}</value>
                    </div>
                    <div class="signal-info">
                        <label>Target</label>
                        <value>$${signal.targets[0].price.toFixed(2)}</value>
                    </div>
                    <div class="signal-info">
                        <label>Strength</label>
                        <value>${signal.strength.toFixed(1)}%</value>
                    </div>
                    <div class="signal-info">
                        <label>RR Ratio</label>
                        <value>${signal.targets[0].rr.toFixed(2)}:1</value>
                    </div>
                </div>
            ` + "`" + `).join('');
        }
        
        function updateStats() {
            // Update stats from current signals
            const signals = document.querySelectorAll('.signal-card:not(.empty)');
            document.getElementById('activeCount').textContent = signals.length;
        }
        
        // Connect on page load
        connect();
    </script>
</body>
</html>`
}

// Template functions
var templateFuncs = template.FuncMap{
	"add": func(a, b int) int {
		return a + b
	},
	"mul": func(a, b float64) float64 {
		return a * b
	},
}
