// Example: How to integrate Go backtest into your existing index.html

// Add this to your backtest.js or create a new file

// Detect if Go backend is available
async function isGoBackendAvailable() {
    try {
        const response = await fetch('http://localhost:8080/api/v1/health', {
            method: 'GET',
            signal: AbortSignal.timeout(1000) // 1 second timeout
        });
        return response.ok;
    } catch {
        return false;
    }
}

// Smart backtest router - uses Go if available, falls back to JS
async function smartBacktest(symbol, interval, days) {
    const useGo = await isGoBackendAvailable();
    
    if (useGo && days > 7) {
        console.log('🚀 Using Go backend for fast backtest');
        return runGoBacktest(symbol, interval, days);
    } else {
        console.log('📊 Using JavaScript backtest');
        return runJSBacktest(symbol, interval, days);
    }
}

// Go-powered backtest
async function runGoBacktest(symbol, interval, days) {
    try {
        const response = await fetch('http://localhost:8080/api/v1/backtest/run', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ symbol, interval, days })
        });

        if (!response.ok) {
            throw new Error('Go backtest failed');
        }

        const result = await response.json();
        
        // Convert Go format to match your existing JS format
        return {
            trades: result.trades,
            totalTrades: result.totalTrades,
            winningTrades: result.winningTrades,
            losingTrades: result.losingTrades,
            totalProfit: result.totalProfit,
            totalLoss: result.totalLoss,
            winRate: result.winRate,
            profitFactor: result.profitFactor,
            averageRR: result.averageRR,
            maxDrawdown: result.maxDrawdown,
            startingBalance: result.startBalance,
            currentBalance: result.finalBalance,
            peakBalance: result.peakBalance,
            netProfit: result.netProfit,
            returnPercent: result.returnPercent
        };
    } catch (error) {
        console.error('Go backtest error:', error);
        // Fallback to JS
        return runJSBacktest(symbol, interval, days);
    }
}

// Your existing JavaScript backtest (keep as fallback)
async function runJSBacktest(symbol, interval, days) {
    // Your existing backtest code here
    return window.runBacktest(symbol, interval, days);
}

// Parallel backtest - Go only feature
async function runParallelBacktest(symbol, interval, daysList = [7, 30, 60, 90]) {
    const useGo = await isGoBackendAvailable();
    
    if (!useGo) {
        alert('⚠️ Go backend not available. Start it with: cd backend && go run *.go');
        return null;
    }

    try {
        const response = await fetch('http://localhost:8080/api/v1/backtest/parallel', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ symbol, interval, daysList })
        });

        if (!response.ok) {
            throw new Error('Parallel backtest failed');
        }

        const data = await response.json();
        
        console.log(`✅ Completed ${daysList.length} backtests in ${data.duration}`);
        
        return data.results;
    } catch (error) {
        console.error('Parallel backtest error:', error);
        return null;
    }
}

// Update your existing backtest button handler
window.startBacktest = async function(days = 30) {
    console.log(`🔄 Starting backtest for ${days} days...`);
    
    const btn = document.getElementById('backtest-btn');
    if (btn) {
        btn.disabled = true;
        btn.textContent = '⏳ Running...';
    }
    
    try {
        // Use smart router
        const result = await smartBacktest(
            window.currentSymbol || 'BTCUSDT',
            window.currentInterval || '15m',
            days
        );
        
        // Display results using your existing function
        window.backtestResults = result;
        displayBacktestResults();
        
    } catch (error) {
        console.error('❌ Backtest error:', error);
    } finally {
        if (btn) {
            btn.disabled = false;
            btn.textContent = '📊 Backtest';
        }
    }
};

// Add parallel backtest button to your HTML
function addParallelBacktestButton() {
    const backtestBtn = document.getElementById('backtest-btn');
    if (!backtestBtn) return;
    
    const parallelBtn = document.createElement('button');
    parallelBtn.id = 'parallel-backtest-btn';
    parallelBtn.className = 'btn btn-secondary';
    parallelBtn.textContent = '⚡ Parallel Backtest';
    parallelBtn.onclick = async () => {
        parallelBtn.disabled = true;
        parallelBtn.textContent = '⏳ Running...';
        
        const results = await runParallelBacktest(
            window.currentSymbol || 'BTCUSDT',
            window.currentInterval || '15m',
            [7, 30, 60, 90]
        );
        
        if (results) {
            displayParallelResults(results);
        }
        
        parallelBtn.disabled = false;
        parallelBtn.textContent = '⚡ Parallel Backtest';
    };
    
    backtestBtn.parentNode.insertBefore(parallelBtn, backtestBtn.nextSibling);
}

// Display parallel results
function displayParallelResults(results) {
    console.log('\n📊 ===== PARALLEL BACKTEST RESULTS =====');
    
    results.forEach((result, index) => {
        const days = [7, 30, 60, 90][index];
        console.log(`\n${days} Days:`);
        console.log(`  Trades: ${result.totalTrades}`);
        console.log(`  Win Rate: ${result.winRate.toFixed(1)}%`);
        console.log(`  Return: ${result.returnPercent.toFixed(2)}%`);
        console.log(`  Profit Factor: ${result.profitFactor.toFixed(2)}`);
        console.log(`  Max Drawdown: ${(result.maxDrawdown * 100).toFixed(2)}%`);
    });
    
    console.log('\n=======================================\n');
    
    // You can also create a visual comparison table
    createComparisonTable(results);
}

// Create comparison table
function createComparisonTable(results) {
    const table = document.createElement('table');
    table.style.cssText = 'width: 100%; border-collapse: collapse; margin-top: 20px;';
    
    table.innerHTML = `
        <thead>
            <tr style="background: #f5f5f5;">
                <th style="padding: 10px; border: 1px solid #ddd;">Period</th>
                <th style="padding: 10px; border: 1px solid #ddd;">Trades</th>
                <th style="padding: 10px; border: 1px solid #ddd;">Win Rate</th>
                <th style="padding: 10px; border: 1px solid #ddd;">Return</th>
                <th style="padding: 10px; border: 1px solid #ddd;">Profit Factor</th>
                <th style="padding: 10px; border: 1px solid #ddd;">Max DD</th>
            </tr>
        </thead>
        <tbody>
            ${results.map((r, i) => {
                const days = [7, 30, 60, 90][i];
                const returnColor = r.returnPercent > 0 ? '#26a69a' : '#ef5350';
                return `
                    <tr>
                        <td style="padding: 10px; border: 1px solid #ddd;">${days} days</td>
                        <td style="padding: 10px; border: 1px solid #ddd;">${r.totalTrades}</td>
                        <td style="padding: 10px; border: 1px solid #ddd;">${r.winRate.toFixed(1)}%</td>
                        <td style="padding: 10px; border: 1px solid #ddd; color: ${returnColor}; font-weight: bold;">
                            ${r.returnPercent.toFixed(2)}%
                        </td>
                        <td style="padding: 10px; border: 1px solid #ddd;">${r.profitFactor.toFixed(2)}</td>
                        <td style="padding: 10px; border: 1px solid #ddd;">${(r.maxDrawdown * 100).toFixed(2)}%</td>
                    </tr>
                `;
            }).join('')}
        </tbody>
    `;
    
    // Insert after backtest results
    const resultsDiv = document.querySelector('.backtest-results') || document.body;
    resultsDiv.appendChild(table);
}

// Initialize on page load
document.addEventListener('DOMContentLoaded', async () => {
    // Check if Go backend is available
    const goAvailable = await isGoBackendAvailable();
    
    if (goAvailable) {
        console.log('✅ Go backend detected - Fast backtests enabled!');
        addParallelBacktestButton();
    } else {
        console.log('ℹ️ Go backend not available - Using JavaScript backtests');
        console.log('💡 Start Go backend with: cd backend && go run *.go');
    }
});

// Export functions
window.smartBacktest = smartBacktest;
window.runGoBacktest = runGoBacktest;
window.runParallelBacktest = runParallelBacktest;
window.isGoBackendAvailable = isGoBackendAvailable;
