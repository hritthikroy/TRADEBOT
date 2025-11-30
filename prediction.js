// TradingView widget instance
let widget;
let currentData = [];
let canvas;
let ctx;

// Make currentData globally accessible for backtest visualization
window.currentData = currentData;
let simulationInterval;
let isSimulating = false;
let predictedCandles = [];
let currentInterval = '15m'; // Default interval
let currentSymbol = 'BTCUSDT'; // Default symbol
let candleCloseTime = null;

// Telegram Configuration
const TELEGRAM_ENABLED = true;
const TELEGRAM_BOT_TOKEN = '8235652305:AAEE8i7XfvRJsPppegSu0pwbDSjFryh60Wk';
const TELEGRAM_CHAT_ID = '8145172959';

// Store historical signals for visualization
let historicalSignals = [];
window.historicalSignals = historicalSignals;

// Audio context for sound alerts
let audioContext = null;
let notificationsEnabled = false;

// Track prediction accuracy
let predictionHistory = [];
let lastPredictions = null;

// Drawing tools
let drawings = [];
let currentDrawingTool = null;
let isDrawing = false;
let drawingStart = null;

// Request notification permission on page load
function requestNotificationPermission() {
    const statusEl = document.getElementById('notification-status');
    
    if ('Notification' in window && Notification.permission === 'default') {
        Notification.requestPermission().then(permission => {
            if (permission === 'granted') {
                notificationsEnabled = true;
                console.log('✅ Notifications enabled');
                if (statusEl) statusEl.textContent = '🔔 Alerts: ON';
                if (statusEl) statusEl.style.color = '#26a69a';
            } else {
                if (statusEl) statusEl.textContent = '🔕 Alerts: OFF';
                if (statusEl) statusEl.style.color = '#ef5350';
            }
        });
    } else if (Notification.permission === 'granted') {
        notificationsEnabled = true;
        if (statusEl) statusEl.textContent = '🔔 Alerts: ON';
        if (statusEl) statusEl.style.color = '#26a69a';
    } else {
        if (statusEl) statusEl.textContent = '🔕 Alerts: OFF';
        if (statusEl) statusEl.style.color = '#ef5350';
    }
}

// Play alert sound
function playAlertSound(type) {
    try {
        if (!audioContext) {
            audioContext = new (window.AudioContext || window.webkitAudioContext)();
        }
        
        const oscillator = audioContext.createOscillator();
        const gainNode = audioContext.createGain();
        
        oscillator.connect(gainNode);
        gainNode.connect(audioContext.destination);
        
        // Different tones for BUY vs SELL
        if (type === 'BUY') {
            // Higher pitch for BUY (bullish)
            oscillator.frequency.value = 800;
        } else {
            // Lower pitch for SELL (bearish)
            oscillator.frequency.value = 400;
        }
        
        oscillator.type = 'sine';
        gainNode.gain.setValueAtTime(0.3, audioContext.currentTime);
        gainNode.gain.exponentialRampToValueAtTime(0.01, audioContext.currentTime + 0.5);
        
        oscillator.start(audioContext.currentTime);
        oscillator.stop(audioContext.currentTime + 0.5);
        
        console.log(`🔊 Alert sound played for ${type}`);
    } catch (error) {
        console.warn('Could not play sound:', error);
    }
}

// Show browser notification
function showNotification(signal) {
    if (!notificationsEnabled || !signal) return;
    
    try {
        const title = `${signal.type} Signal - ${currentSymbol}`;
        const body = `Entry: ${signal.entry.toFixed(2)}\nStop Loss: ${signal.stopLoss.toFixed(2)}\nTP1: ${signal.targets[0].price.toFixed(2)} (${signal.targets[0].rr.toFixed(1)}R)\nStrength: ${signal.strength}%`;
        
        const notification = new Notification(title, {
            body: body,
            icon: signal.type === 'BUY' ? '📈' : '📉',
            badge: signal.type === 'BUY' ? '📈' : '📉',
            tag: 'trading-signal',
            requireInteraction: true,
            silent: false
        });
        
        notification.onclick = () => {
            window.focus();
            notification.close();
        };
        
        console.log('📬 Browser notification sent');
    } catch (error) {
        console.warn('Could not show notification:', error);
    }
}

// Send Telegram alert
async function sendTelegramAlert(signal) {
    if (!signal) return;
    
    // Always save signal to tracker first
    saveSignalToTracker(signal);
    
    if (!TELEGRAM_ENABLED) {
        console.log('📊 Signal saved to tracker (Telegram disabled)');
        return;
    }
    
    try {
        const message = `🚨 *${signal.type} SIGNAL*\n\n` +
                       `📊 Symbol: ${currentSymbol}\n` +
                       `💰 Entry: ${signal.entry.toFixed(2)}\n` +
                       `🛑 Stop Loss: ${signal.stopLoss.toFixed(2)}\n\n` +
                       `🎯 Targets:\n` +
                       `TP1: ${signal.targets[0].price.toFixed(2)} (${signal.targets[0].rr.toFixed(1)}R)\n` +
                       `TP2: ${signal.targets[1].price.toFixed(2)} (${signal.targets[1].rr.toFixed(1)}R)\n` +
                       `TP3: ${signal.targets[2].price.toFixed(2)} (${signal.targets[2].rr.toFixed(1)}R)\n\n` +
                       `💪 Strength: ${signal.strength}%\n` +
                       `⏰ Time: ${new Date().toLocaleString()}`;
        
        const url = `https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/sendMessage`;
        
        const response = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                chat_id: TELEGRAM_CHAT_ID,
                text: message,
                parse_mode: 'Markdown'
            })
        });
        
        if (response.ok) {
            console.log('📱 Telegram alert sent successfully!');
        } else {
            console.error('❌ Telegram send failed:', await response.text());
        }
    } catch (error) {
        console.error('Telegram error:', error);
    }
}

// Save signal to tracker
// Debug function to check Supabase status
window.checkSupabaseStatus = function() {
    console.log('=== Supabase Status Check ===');
    console.log('SupabaseDB loaded:', typeof SupabaseDB !== 'undefined');
    console.log('window.supabase loaded:', typeof window.supabase !== 'undefined');
    console.log('SupabaseDB object:', SupabaseDB);
    if (typeof SupabaseDB !== 'undefined') {
        console.log('SupabaseDB.saveSignal:', typeof SupabaseDB.saveSignal);
    }
    console.log('============================');
};

async function saveSignalToTracker(signal) {
    try {
        const signalData = {
            id: Date.now() + Math.random(), // Ensure unique ID
            timestamp: Date.now(),
            type: signal.type,
            symbol: currentSymbol,
            entry: signal.entry,
            stopLoss: signal.stopLoss,
            tp1: signal.targets[0].price,
            tp2: signal.targets[1].price,
            tp3: signal.targets[2].price,
            strength: signal.strength,
            status: 'pending',
            exitPrice: null,
            exitReason: null,
            trailingStopPrice: null,
            trailingStopActive: false,
            livePrice: null
        };
        
        // Save to localStorage
        const signals = JSON.parse(localStorage.getItem('tradingSignals') || '[]');
        signals.unshift(signalData);
        localStorage.setItem('tradingSignals', JSON.stringify(signals));
        console.log('💾 Signal saved to localStorage:', signalData);
        
        // Save to Supabase if available
        if (typeof SupabaseDB !== 'undefined' && typeof window.supabase !== 'undefined') {
            try {
                console.log('🔄 Attempting to save to Supabase...');
                console.log('Signal data to save:', {
                    signal_id: signalData.id.toString(),
                    signal_type: signalData.type,
                    symbol: signalData.symbol,
                    entry_price: signalData.entry,
                    stop_loss: signalData.stopLoss,
                    tp1: signalData.tp1,
                    tp2: signalData.tp2,
                    tp3: signalData.tp3,
                    strength: signalData.strength,
                    status: 'pending'
                });
                
                const result = await SupabaseDB.saveSignal({
                    signal_id: signalData.id.toString(),
                    signal_type: signalData.type,
                    symbol: signalData.symbol,
                    entry_price: signalData.entry,
                    stop_loss: signalData.stopLoss,
                    tp1: signalData.tp1,
                    tp2: signalData.tp2,
                    tp3: signalData.tp3,
                    strength: signalData.strength,
                    status: 'pending'
                });
                console.log('☁️ Signal saved to Supabase!', result);
            } catch (error) {
                console.error('❌ Supabase save failed:', error);
                console.error('Error message:', error.message);
                console.error('Error code:', error.code);
                console.error('Full error:', JSON.stringify(error, null, 2));
            }
        } else {
            console.warn('⚠️ SupabaseDB not loaded. SupabaseDB exists:', typeof SupabaseDB !== 'undefined', 'supabase exists:', typeof window.supabase !== 'undefined');
        }
        
        // Trigger storage event manually for same-window updates
        window.dispatchEvent(new StorageEvent('storage', {
            key: 'tradingSignals',
            newValue: JSON.stringify(signals),
            url: window.location.href
        }));
        
        // Start monitoring this signal
        monitorSignal(signalData);
    } catch (error) {
        console.error('Error saving signal:', error);
    }
}

// Update signal trailing stop in localStorage
function updateSignalTrailingStop(signalId, trailingStopPrice) {
    try {
        const signals = JSON.parse(localStorage.getItem('tradingSignals') || '[]');
        const index = signals.findIndex(s => s.id === signalId);
        
        if (index !== -1) {
            signals[index].trailingStopPrice = trailingStopPrice;
            signals[index].trailingStopActive = true;
            localStorage.setItem('tradingSignals', JSON.stringify(signals));
        }
    } catch (error) {
        console.error('Error updating trailing stop:', error);
    }
}

// Monitor signal for automatic win/loss detection with trailing stop
function monitorSignal(signal) {
    let highestPrice = signal.entry; // Track highest price for BUY
    let lowestPrice = signal.entry;  // Track lowest price for SELL
    let trailingStopActive = false;
    let trailingStopPrice = signal.stopLoss;
    
    const checkInterval = setInterval(async () => {
        try {
            // Get current price
            const response = await fetch(`https://api.binance.com/api/v3/ticker/price?symbol=${signal.symbol}`);
            const data = await response.json();
            const currentPrice = parseFloat(data.price);
            
            let shouldUpdate = false;
            let newStatus = signal.status;
            let exitPrice = null;
            let exitReason = null;
            
            if (signal.type === 'BUY') {
                // Update highest price
                if (currentPrice > highestPrice) {
                    highestPrice = currentPrice;
                }
                
                // Activate trailing stop after TP1 is hit
                if (currentPrice >= signal.tp1 && !trailingStopActive) {
                    trailingStopActive = true;
                    // Set trailing stop to lock in 50% of profit
                    const profitDistance = signal.tp1 - signal.entry;
                    trailingStopPrice = signal.entry + (profitDistance * 0.5);
                    console.log(`🔒 Trailing stop activated at ${trailingStopPrice.toFixed(2)}`);
                    
                    // Update signal in localStorage with trailing stop
                    updateSignalTrailingStop(signal.id, trailingStopPrice);
                }
                
                // Update trailing stop (lock 50% of gains from highest point)
                if (trailingStopActive) {
                    const newTrailingStop = Math.max(trailingStopPrice, signal.entry + ((highestPrice - signal.entry) * 0.5));
                    if (newTrailingStop > trailingStopPrice) {
                        trailingStopPrice = newTrailingStop;
                        updateSignalTrailingStop(signal.id, trailingStopPrice);
                    }
                }
                
                // Check if trailing stop hit
                if (trailingStopActive && currentPrice <= trailingStopPrice) {
                    newStatus = 'win';
                    exitPrice = trailingStopPrice;
                    exitReason = 'Trailing Stop';
                    shouldUpdate = true;
                }
                // Check if stop loss hit
                else if (currentPrice <= signal.stopLoss) {
                    newStatus = 'loss';
                    exitPrice = signal.stopLoss;
                    exitReason = 'Stop Loss';
                    shouldUpdate = true;
                }
                // Check if TP3 hit
                else if (currentPrice >= signal.tp3) {
                    newStatus = 'win';
                    exitPrice = signal.tp3;
                    exitReason = 'TP3';
                    shouldUpdate = true;
                }
                // Check if TP2 hit (but continue for TP3)
                else if (currentPrice >= signal.tp2 && !trailingStopActive) {
                    // Just activate trailing stop, don't exit yet
                    trailingStopActive = true;
                    const profitDistance = currentPrice - signal.entry;
                    trailingStopPrice = signal.entry + (profitDistance * 0.5);
                    console.log(`🔒 Trailing stop activated at TP2: ${trailingStopPrice.toFixed(2)}`);
                }
            } else { // SELL
                // Update lowest price
                if (currentPrice < lowestPrice) {
                    lowestPrice = currentPrice;
                }
                
                // Activate trailing stop after TP1 is hit
                if (currentPrice <= signal.tp1 && !trailingStopActive) {
                    trailingStopActive = true;
                    // Set trailing stop to lock in 50% of profit
                    const profitDistance = signal.entry - signal.tp1;
                    trailingStopPrice = signal.entry - (profitDistance * 0.5);
                    console.log(`🔒 Trailing stop activated at ${trailingStopPrice.toFixed(2)}`);
                    
                    // Update signal in localStorage with trailing stop
                    updateSignalTrailingStop(signal.id, trailingStopPrice);
                }
                
                // Update trailing stop (lock 50% of gains from lowest point)
                if (trailingStopActive) {
                    const newTrailingStop = Math.min(trailingStopPrice, signal.entry - ((signal.entry - lowestPrice) * 0.5));
                    if (newTrailingStop < trailingStopPrice) {
                        trailingStopPrice = newTrailingStop;
                        updateSignalTrailingStop(signal.id, trailingStopPrice);
                    }
                }
                
                // Check if trailing stop hit
                if (trailingStopActive && currentPrice >= trailingStopPrice) {
                    newStatus = 'win';
                    exitPrice = trailingStopPrice;
                    exitReason = 'Trailing Stop';
                    shouldUpdate = true;
                }
                // Check if stop loss hit
                else if (currentPrice >= signal.stopLoss) {
                    newStatus = 'loss';
                    exitPrice = signal.stopLoss;
                    exitReason = 'Stop Loss';
                    shouldUpdate = true;
                }
                // Check if TP3 hit
                else if (currentPrice <= signal.tp3) {
                    newStatus = 'win';
                    exitPrice = signal.tp3;
                    exitReason = 'TP3';
                    shouldUpdate = true;
                }
                // Check if TP2 hit (but continue for TP3)
                else if (currentPrice <= signal.tp2 && !trailingStopActive) {
                    // Just activate trailing stop, don't exit yet
                    trailingStopActive = true;
                    const profitDistance = signal.entry - currentPrice;
                    trailingStopPrice = signal.entry - (profitDistance * 0.5);
                    console.log(`🔒 Trailing stop activated at TP2: ${trailingStopPrice.toFixed(2)}`);
                }
            }
            
            if (shouldUpdate) {
                // Update signal in localStorage
                const signals = JSON.parse(localStorage.getItem('tradingSignals') || '[]');
                const index = signals.findIndex(s => s.id === signal.id);
                
                if (index !== -1) {
                    signals[index].status = newStatus;
                    signals[index].exitPrice = exitPrice;
                    signals[index].exitReason = exitReason;
                    localStorage.setItem('tradingSignals', JSON.stringify(signals));
                    
                    console.log(`✅ Signal auto-updated: ${newStatus.toUpperCase()} - ${exitReason} at ${exitPrice.toFixed(2)}`);
                    
                    // Send Telegram notification
                    sendTradeResultToTelegram(signal, newStatus, exitPrice, exitReason);
                }
                
                // Stop monitoring
                clearInterval(checkInterval);
            }
            
            // Stop monitoring after 24 hours
            if (Date.now() - signal.timestamp > 24 * 60 * 60 * 1000) {
                clearInterval(checkInterval);
                console.log('⏱️ Signal monitoring stopped (24h timeout)');
            }
            
        } catch (error) {
            console.error('Error monitoring signal:', error);
        }
    }, 10000); // Check every 10 seconds
}

// Send trade result to Telegram
async function sendTradeResultToTelegram(signal, status, exitPrice, exitReason) {
    if (!TELEGRAM_ENABLED) return;
    
    try {
        const emoji = status === 'win' ? '✅' : '❌';
        const statusText = status === 'win' ? 'WIN' : 'LOSS';
        
        const profitPercent = signal.type === 'BUY' 
            ? ((exitPrice - signal.entry) / signal.entry) * 100
            : ((signal.entry - exitPrice) / signal.entry) * 100;
        
        const message = `${emoji} *TRADE ${statusText}*\n\n` +
                       `📊 Symbol: ${signal.symbol}\n` +
                       `📍 Type: ${signal.type}\n` +
                       `💰 Entry: ${signal.entry.toFixed(2)}\n` +
                       `🎯 Exit: ${exitPrice.toFixed(2)}\n` +
                       `📈 Result: ${profitPercent > 0 ? '+' : ''}${profitPercent.toFixed(2)}%\n` +
                       `🏁 Reason: ${exitReason}\n` +
                       `⏰ Time: ${new Date().toLocaleString()}`;
        
        const url = `https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/sendMessage`;
        
        await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                chat_id: TELEGRAM_CHAT_ID,
                text: message,
                parse_mode: 'Markdown'
            })
        });
        
        console.log('📱 Trade result sent to Telegram');
    } catch (error) {
        console.error('Error sending trade result:', error);
    }
}

// Initialize TradingView Advanced Chart
function initChart() {
    widget = new TradingView.widget({
        "autosize": true,
        "symbol": "BINANCE:BTCUSDT",
        "interval": "15",
        "timezone": "Etc/UTC",
        "theme": "dark",
        "style": "1",
        "locale": "en",
        "enable_publishing": false,
        "allow_symbol_change": true,
        "container_id": "tradingview_chart",
        "hide_top_toolbar": false,
        "hide_legend": true,
        "hide_side_toolbar": false,
        "hide_volume": true,
        "save_image": false,
        "studies": [],
        "show_popup_button": false,
        "backgroundColor": "#131722",
        "gridColor": "rgba(0, 0, 0, 0)",
        "overrides": {
            "paneProperties.background": "#131722",
            "paneProperties.backgroundType": "solid",
            "paneProperties.vertGridProperties.color": "#131722",
            "paneProperties.horzGridProperties.color": "#131722",
            "scalesProperties.lineColor": "#2a2e39"
        },
        "disabled_features": [
            "volume_force_overlay", 
            "display_market_status", 
            "header_volume",
            "header_compare",
            "header_screenshot",
            "header_saveload",
            "header_symbol_search",
            "symbol_search_hot_key",
            "header_chart_type",
            "header_settings",
            "header_indicators",
            "header_fullscreen_button"
        ],
        "enabled_features": [
            "side_toolbar_in_fullscreen_mode"
        ]
    });
    
    // Listen for chart changes (symbol and interval)
    if (widget.onChartReady) {
        widget.onChartReady(() => {
        const chart = widget.activeChart();
        console.log('TradingView chart ready');
        
        // Make grid lines same color as background (invisible)
        chart.applyOverrides({
            "paneProperties.vertGridProperties.color": "#131722",
            "paneProperties.horzGridProperties.color": "#131722"
        });
        
        // Remove all volume studies
        chart.getAllStudies().forEach(study => {
            if (study.name.toLowerCase().includes('volume')) {
                chart.removeEntity(study.id);
            }
        });
        
        // Subscribe to symbol changes
        chart.onSymbolChanged().subscribe(null, (symbolData) => {
            const fullSymbol = symbolData.ticker || symbolData.name;
            console.log('TradingView symbol changed:', fullSymbol);
            if (fullSymbol && fullSymbol.includes(':')) {
                const parts = fullSymbol.split(':');
                if (parts[1]) {
                    currentSymbol = parts[1].replace('USDT', 'USDT');
                    console.log('Updated symbol to:', currentSymbol);
                    if (isSimulating) {
                        fetchRealMarketData();
                    }
                }
            }
        });
        
        // Subscribe to interval changes from TradingView
        chart.onIntervalChanged().subscribe(null, (interval) => {
            const newInterval = convertTradingViewInterval(interval);
            console.log('TradingView interval changed from', currentInterval, 'to', newInterval);
            currentInterval = newInterval;
            
            // Update button active state
            document.querySelectorAll('.time-btn').forEach(btn => {
                btn.classList.remove('active');
                if (btn.textContent === newInterval) {
                    btn.classList.add('active');
                }
            });
            
            if (isSimulating) {
                fetchRealMarketData();
            }
        });
        });
    } else {
        console.warn('TradingView onChartReady not available');
    }
}

// Convert TradingView interval format to Binance format
function convertTradingViewInterval(tvInterval) {
    const intervalMap = {
        '1': '1m',
        '3': '3m',
        '5': '5m',
        '15': '15m',
        '30': '30m',
        '60': '1h',
        '120': '2h',
        '240': '4h',
        '1D': '1d',
        'D': '1d',
        '1W': '1w',
        'W': '1w',
        '1M': '1M',
        'M': '1M'
    };
    return intervalMap[tvInterval] || '15m';
}

// Convert Binance interval to TradingView format
function convertToTradingViewInterval(interval) {
    const map = {
        '1m': '1',
        '3m': '3',
        '5m': '5',
        '15m': '15',
        '30m': '30',
        '1h': '60',
        '2h': '120',
        '4h': '240',
        '1d': 'D',
        '1w': 'W',
        '1M': 'M'
    };
    return map[interval] || '15';
}

// Get interval in milliseconds
function getIntervalMs(interval) {
    const map = {
        '1m': 60000,
        '3m': 180000,
        '5m': 300000,
        '15m': 900000,
        '30m': 1800000,
        '1h': 3600000,
        '2h': 7200000,
        '4h': 14400000,
        '1d': 86400000,
        '1w': 604800000,
        '1M': 2592000000
    };
    return map[interval] || 900000;
}

// Initialize canvas
function initCanvas() {
    canvas = document.getElementById('prediction-chart');
    if (!canvas) {
        console.error('Canvas not found');
        return;
    }
    ctx = canvas.getContext('2d');
    
    // Set canvas size properly
    const rect = canvas.getBoundingClientRect();
    canvas.width = rect.width || canvas.offsetWidth || 1200;
    canvas.height = rect.height || canvas.offsetHeight || 500;
    
    console.log('Canvas initialized:', canvas.width, 'x', canvas.height);
    
    // Add mouse wheel event listener for zoom
    canvas.addEventListener('wheel', handleWheelZoom, { passive: false });
    
    // Initial draw
    drawChart();
}

// Multi-timeframe data storage
let multiTimeframeData = {};

// Fetch multi-timeframe data for better predictions
async function fetchMultiTimeframeData() {
    try {
        // Get all higher timeframes to analyze
        const higherTFs = getHigherTimeframes(currentInterval);
        const timeframes = [currentInterval, ...higherTFs];
        
        // Fetch data for all timeframes in parallel
        const promises = timeframes.map(tf => 
            fetch(`https://api.binance.com/api/v3/klines?symbol=${currentSymbol}&interval=${tf}&limit=100`)
                .then(res => res.json())
                .then(data => ({ timeframe: tf, data }))
                .catch(err => {
                    console.warn(`Failed to fetch ${tf} data:`, err);
                    return { timeframe: tf, data: [] };
                })
        );
        
        const results = await Promise.all(promises);
        
        // Store multi-timeframe data (include volume)
        results.forEach(result => {
            if (result.data.length > 0) {
                multiTimeframeData[result.timeframe] = result.data.map(candle => ({
                    time: candle[0],
                    open: parseFloat(candle[1]),
                    high: parseFloat(candle[2]),
                    low: parseFloat(candle[3]),
                    close: parseFloat(candle[4]),
                    volume: parseFloat(candle[5])
                }));
            }
        });
        
        console.log('Multi-timeframe data loaded:', Object.keys(multiTimeframeData).join(', '));
        return true;
    } catch (error) {
        console.error('Error fetching multi-timeframe data:', error);
        return false;
    }
}

// Get multiple higher timeframes for better analysis
function getHigherTimeframes(interval) {
    const map = {
        '1m': ['3m', '5m', '15m'],  // 1m uses 3 higher timeframes
        '3m': ['5m', '15m', '1h'],
        '5m': ['15m', '30m', '1h'],
        '15m': ['30m', '1h', '4h'],
        '30m': ['1h', '4h', '1d'],
        '1h': ['4h', '1d'],
        '2h': ['4h', '1d'],
        '4h': ['1d', '1w'],
        '1d': ['1w']
    };
    return map[interval] || [];
}

// Get single higher timeframe (for backward compatibility)
function getHigherTimeframe(interval) {
    const timeframes = getHigherTimeframes(interval);
    return timeframes.length > 0 ? timeframes[0] : null;
}

// Fetch real market data from Binance API
async function fetchRealMarketData() {
    try {
        updateStatus('Fetching multi-timeframe data...');
        
        // Fetch multi-timeframe data
        await fetchMultiTimeframeData();
        
        // Get current timeframe data
        const url = `https://api.binance.com/api/v3/klines?symbol=${currentSymbol}&interval=${currentInterval}&limit=200`;
        console.log('Fetching:', url);
        
        const response = await fetch(url);
        console.log('Response status:', response.status);
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const data = await response.json();
        console.log('Received data:', data.length, 'candles');
        
        if (!Array.isArray(data)) {
            throw new Error('Invalid data received from API');
        }
        
        if (data.length === 0) {
            throw new Error('No data received from API');
        }
        
        // Convert Binance data to our format (include volume)
        currentData = data.map(candle => ({
            time: candle[0],
            open: parseFloat(candle[1]),
            high: parseFloat(candle[2]),
            low: parseFloat(candle[3]),
            close: parseFloat(candle[4]),
            volume: parseFloat(candle[5])
        }));
        
        console.log('Converted data:', currentData.length, 'candles');
        
        // Check prediction accuracy before generating new ones
        checkPredictionAccuracy();
        
        // Calculate next candle close time
        const lastCandle = currentData[currentData.length - 1];
        const intervalMs = getIntervalMs(currentInterval);
        candleCloseTime = lastCandle.time + intervalMs;
        
        // Auto-generate predictions with multi-timeframe and AI analysis
        await autoPredictNextCandles();
        
        drawChart();
        updateStatus(`✓ ${currentSymbol} (${currentInterval}) - Multi-TF predictions active`);
        return true;
    } catch (error) {
        console.error('Error fetching market data:', error);
        updateStatus('❌ Failed: ' + error.message);
        
        // Show error on canvas
        if (ctx) {
            ctx.fillStyle = '#131722';
            ctx.fillRect(0, 0, canvas.width, canvas.height);
            ctx.fillStyle = '#ef5350';
            ctx.font = '16px Arial';
            ctx.textAlign = 'center';
            ctx.fillText('Error loading data: ' + error.message, canvas.width / 2, canvas.height / 2);
        }
        
        return false;
    }
}

// Update status message
function updateStatus(message) {
    const statusEl = document.getElementById('status-message');
    if (statusEl) {
        statusEl.textContent = message;
        statusEl.style.color = '#2962ff';
    }
}

// Update countdown timer
function updateCountdown() {
    if (!candleCloseTime) return;
    
    const now = Date.now();
    const timeLeft = candleCloseTime - now;
    
    if (timeLeft <= 0) {
        updateStatus('Candle closing... fetching new data');
        return;
    }
    
    const minutes = Math.floor(timeLeft / 60000);
    const seconds = Math.floor((timeLeft % 60000) / 1000);
    
    const countdownEl = document.getElementById('countdown-timer');
    if (countdownEl) {
        countdownEl.textContent = `${minutes}:${seconds.toString().padStart(2, '0')}`;
    }
}

// Start countdown timer
function startCountdownTimer() {
    setInterval(updateCountdown, 1000);
}

// Start real-time simulation with real data
async function startSimulation() {
    console.log('Starting simulation...');
    
    if (isSimulating) {
        console.log('Already simulating');
        return;
    }
    
    isSimulating = true;
    predictedCandles = [];
    
    // Fetch initial real data
    console.log('Fetching initial data...');
    const success = await fetchRealMarketData();
    
    if (!success) {
        console.error('Failed to fetch initial data');
        isSimulating = false;
        return;
    }
    
    console.log('Initial data loaded successfully');
    
    // Update data every 30 seconds for stable long-term predictions
    simulationInterval = setInterval(async () => {
        console.log('Auto-updating data...');
        await fetchRealMarketData();
    }, 30000);
    
    console.log('Simulation started, will update every 3 seconds');
}

// Stop simulation
function stopSimulation() {
    isSimulating = false;
    if (simulationInterval) {
        clearInterval(simulationInterval);
    }
}

// Draw the chart with candles
function drawChart() {
    if (!ctx) return;
    
    // Performance optimization: Use requestAnimationFrame for smoother rendering
    if (window.drawChartPending) return;
    window.drawChartPending = true;
    
    requestAnimationFrame(() => {
        window.drawChartPending = false;
        drawChartInternal();
    });
}

function drawChartInternal() {
    if (!ctx) return;
    
    // Clear canvas with TradingView-like background
    ctx.fillStyle = '#131722';
    ctx.fillRect(0, 0, canvas.width, canvas.height);
    
    if (currentData.length === 0) {
        ctx.fillStyle = '#787b86';
        ctx.font = '16px Arial';
        ctx.textAlign = 'center';
        ctx.fillText('Loading market data...', canvas.width / 2, canvas.height / 2);
        return;
    }
    
    // Show candles based on zoom level
    const realCandlesToShow = Math.min(zoomLevel, currentData.length);
    const displayCandles = [...currentData.slice(-realCandlesToShow), ...predictedCandles];
    
    // Calculate price range
    const allPrices = displayCandles.flatMap(c => [c.high, c.low]);
    const maxPrice = Math.max(...allPrices);
    const minPrice = Math.min(...allPrices);
    const priceRange = maxPrice - minPrice;
    const padding = priceRange * 0.15;
    
    // Chart dimensions
    const rightPadding = 80;
    const leftPadding = 10;
    const topPadding = 20;
    const bottomPadding = 40;
    const chartWidth = canvas.width - leftPadding - rightPadding;
    const chartHeight = canvas.height - topPadding - bottomPadding;
    
    // Fixed candle width for consistent sizing
    const totalCandles = displayCandles.length;
    const candleWidth = chartWidth / (totalCandles + 1);
    // Use fixed body width ratio for all candles
    const candleBodyWidth = candleWidth * 0.65;
    
    // Draw price labels on the right (no grid lines)
    const priceLabels = 6;
    for (let i = 0; i <= priceLabels; i++) {
        const y = topPadding + (chartHeight / priceLabels) * i;
        const price = maxPrice + padding - ((maxPrice + padding - (minPrice - padding)) / priceLabels) * i;
        
        // Price text only
        ctx.fillStyle = '#787b86';
        ctx.font = '12px monospace';
        ctx.textAlign = 'left';
        ctx.fillText(price.toFixed(2), canvas.width - rightPadding + 10, y + 4);
    }
    

    
    // Draw all candles with consistent sizing
    const realCandleCount = realCandlesToShow;
    
    displayCandles.forEach((candle, index) => {
        const isPredicted = index >= realCandleCount;
        drawCandle(candle, index, candleWidth, candleBodyWidth, leftPadding, topPadding, chartHeight, maxPrice, minPrice, padding, isPredicted, false);
    });
    
    // Draw current price line
    if (currentData.length > 0) {
        const currentPrice = currentData[currentData.length - 1].close;
        const priceToY = (price) => {
            return topPadding + chartHeight - ((price - (minPrice - padding)) / (maxPrice + padding - (minPrice - padding))) * chartHeight;
        };
        
        const priceY = priceToY(currentPrice);
        
        // Draw horizontal line
        ctx.strokeStyle = '#2196f3';
        ctx.lineWidth = 2;
        ctx.setLineDash([5, 5]);
        ctx.beginPath();
        ctx.moveTo(leftPadding, priceY);
        ctx.lineTo(canvas.width - rightPadding, priceY);
        ctx.stroke();
        ctx.setLineDash([]);
        
        // Draw price label
        const labelWidth = 80;
        const labelHeight = 24;
        const labelX = canvas.width - rightPadding + 5;
        const labelY = priceY - labelHeight / 2;
        
        // Label background
        ctx.fillStyle = '#2196f3';
        ctx.fillRect(labelX, labelY, labelWidth, labelHeight);
        
        // Label text
        ctx.fillStyle = '#ffffff';
        ctx.font = 'bold 13px monospace';
        ctx.textAlign = 'left';
        ctx.fillText(`$${currentPrice.toFixed(2)}`, labelX + 5, labelY + 16);
        
        // Small arrow pointing to line
        ctx.beginPath();
        ctx.moveTo(labelX, priceY);
        ctx.lineTo(labelX - 5, priceY - 5);
        ctx.lineTo(labelX - 5, priceY + 5);
        ctx.closePath();
        ctx.fillStyle = '#2196f3';
        ctx.fill();
    }
    
    // Draw support and resistance levels
    if (currentData.length > 20) {
        const recent50 = currentData.slice(-50);
        const highs = recent50.map(c => c.high);
        const lows = recent50.map(c => c.low);
        
        // Calculate pivot levels
        const pivotHigh = Math.max(...highs);
        const pivotLow = Math.min(...lows);
        const pivotClose = currentData[currentData.length - 1].close;
        const pivot = (pivotHigh + pivotLow + pivotClose) / 3;
        
        // Calculate S/R levels
        const r1 = (2 * pivot) - pivotLow;
        const r2 = pivot + (pivotHigh - pivotLow);
        const s1 = (2 * pivot) - pivotHigh;
        const s2 = pivot - (pivotHigh - pivotLow);
        
        const priceToY = (price) => {
            return topPadding + chartHeight - ((price - (minPrice - padding)) / (maxPrice + padding - (minPrice - padding))) * chartHeight;
        };
        
        // Draw resistance levels (RED = Price likely to FALL here)
        if (r1 >= minPrice && r1 <= maxPrice) {
            const y = priceToY(r1);
            ctx.strokeStyle = '#ef5350';
            ctx.lineWidth = 2;
            ctx.setLineDash([]);
            ctx.beginPath();
            ctx.moveTo(leftPadding, y);
            ctx.lineTo(leftPadding + chartWidth, y);
            ctx.stroke();
            
            // Label with background
            ctx.fillStyle = '#ef5350';
            ctx.fillRect(leftPadding + chartWidth + 5, y - 10, 80, 16);
            ctx.fillStyle = '#ffffff';
            ctx.font = 'bold 11px Arial';
            ctx.textAlign = 'left';
            ctx.fillText(`RESIST ${r1.toFixed(2)}`, leftPadding + chartWidth + 10, y + 2);
        }
        
        // Draw support levels (GREEN = Price likely to RISE here)
        if (s1 >= minPrice && s1 <= maxPrice) {
            const y = priceToY(s1);
            ctx.strokeStyle = '#26a69a';
            ctx.lineWidth = 2;
            ctx.setLineDash([]);
            ctx.beginPath();
            ctx.moveTo(leftPadding, y);
            ctx.lineTo(leftPadding + chartWidth, y);
            ctx.stroke();
            
            // Label with background
            ctx.fillStyle = '#26a69a';
            ctx.fillRect(leftPadding + chartWidth + 5, y - 10, 85, 16);
            ctx.fillStyle = '#ffffff';
            ctx.font = 'bold 11px Arial';
            ctx.textAlign = 'left';
            ctx.fillText(`SUPPORT ${s1.toFixed(2)}`, leftPadding + chartWidth + 10, y + 2);
        }
    }
    
    // Draw historical signal markers (small arrows)
    if (historicalSignals.length > 0) {
        const startIndex = Math.max(0, currentData.length - realCandlesToShow);
        
        const priceToYFunc = (price) => {
            return topPadding + chartHeight - ((price - (minPrice - padding)) / (maxPrice + padding - (minPrice - padding))) * chartHeight;
        };
        
        historicalSignals.forEach(signal => {
            const signalIndex = signal.candleIndex - startIndex;
            
            // Only draw if signal is in visible range
            if (signalIndex >= 0 && signalIndex < realCandlesToShow) {
                const x = leftPadding + (signalIndex + 0.5) * candleWidth;
                const signalPrice = signal.entry;
                const y = priceToYFunc(signalPrice);
                
                // Draw small arrow
                ctx.fillStyle = signal.type === 'BUY' ? '#26a69a' : '#ef5350';
                ctx.beginPath();
                if (signal.type === 'BUY') {
                    // Up arrow
                    ctx.moveTo(x, y + 15);
                    ctx.lineTo(x - 5, y + 25);
                    ctx.lineTo(x + 5, y + 25);
                } else {
                    // Down arrow
                    ctx.moveTo(x, y - 15);
                    ctx.lineTo(x - 5, y - 25);
                    ctx.lineTo(x + 5, y - 25);
                }
                ctx.closePath();
                ctx.fill();
                
                // Draw small label
                ctx.fillStyle = signal.type === 'BUY' ? '#26a69a' : '#ef5350';
                ctx.font = 'bold 9px Arial';
                ctx.textAlign = 'center';
                ctx.fillText(signal.type, x, signal.type === 'BUY' ? y + 38 : y - 28);
            }
        });
    }
    
    // Draw trading signals on chart (current signal - larger)
    if (window.currentTradingSignal) {
        drawTradingSignalOnChart(window.currentTradingSignal, leftPadding, topPadding, chartHeight, maxPrice, minPrice, padding);
    }
    
    // Draw time labels at bottom (show fewer labels for clarity)
    const timeStep = Math.max(1, Math.floor(displayCandles.length / 8));
    for (let i = 0; i < displayCandles.length; i += timeStep) {
        const x = leftPadding + (i + 0.5) * candleWidth;
        const time = new Date(displayCandles[i].time);
        const timeStr = time.getHours().toString().padStart(2, '0') + ':' + time.getMinutes().toString().padStart(2, '0');
        
        ctx.fillStyle = '#787b86';
        ctx.font = '10px Arial';
        ctx.textAlign = 'center';
        ctx.fillText(timeStr, x, canvas.height - bottomPadding + 20);
    }
}

// Export drawChart for backtest visualization
window.drawChart = drawChart;

// Draw individual candle (TradingView style)
function drawCandle(candle, index, candleWidth, candleBodyWidth, leftPadding, topPadding, chartHeight, maxPrice, minPrice, padding, isPredicted, isRunning) {
    // Calculate exact center position for this candle
    const x = leftPadding + (index + 1) * candleWidth;
    
    const priceToY = (price) => {
        return topPadding + chartHeight - ((price - (minPrice - padding)) / (maxPrice + padding - (minPrice - padding))) * chartHeight;
    };
    
    const highY = priceToY(candle.high);
    const lowY = priceToY(candle.low);
    const openY = priceToY(candle.open);
    const closeY = priceToY(candle.close);
    
    const isGreen = candle.close >= candle.open;
    
    // Colors - TradingView style
    // Predicted candles use white color
    let wickColor, bodyColor;
    
    if (isPredicted) {
        // White color for predicted candles
        wickColor = '#ffffff';
        bodyColor = '#ffffff';
    } else {
        // Green/Red for real candles
        wickColor = isGreen ? '#26a69a' : '#ef5350';
        bodyColor = isGreen ? '#26a69a' : '#ef5350';
    }
    
    // Draw wick (thin line) - always centered
    ctx.strokeStyle = wickColor;
    ctx.lineWidth = 1;
    ctx.beginPath();
    ctx.moveTo(x, highY);
    ctx.lineTo(x, lowY);
    ctx.stroke();
    
    // Draw body with EXACT same width for all candles
    const bodyTop = Math.min(openY, closeY);
    const bodyHeight = Math.max(Math.abs(closeY - openY), 1);
    
    // Use exact same body width for all candles
    const exactBodyWidth = candleBodyWidth;
    const bodyLeft = x - exactBodyWidth / 2;
    
    ctx.fillStyle = bodyColor;
    ctx.fillRect(bodyLeft, bodyTop, exactBodyWidth, bodyHeight);
    
    // Add white border for predicted candles (no border, just white color)
    if (isPredicted) {
        // No border needed - predicted candles are already white
    }
}

// Calculate candle height
function getCandleHeight(candle) {
    return candle.high - candle.low;
}

// Auto-predict next 3 candles with multi-timeframe, volume, and AI analysis
async function autoPredictNextCandles() {
    if (currentData.length === 0) {
        return;
    }
    
    const currentCandle = currentData[currentData.length - 1];
    
    // Get Ensemble AI prediction (combines all sources)
    let aiPrediction = null;
    if (window.getAIEnhancedPrediction) {
        aiPrediction = await getAIEnhancedPrediction(currentSymbol, currentInterval, currentData);
    }
    
    // Multi-timeframe trend analysis
    const mtfAnalysis = analyzeMultiTimeframeTrend();
    
    // Volume analysis for volatility prediction
    const volumeAnalysis = analyzeVolume();
    
    // Current timeframe analysis with Delta (last 10 candles)
    const recentCandles = currentData.slice(-10);
    let upCount = 0;
    let downCount = 0;
    let totalPriceChange = 0;
    let cumulativeDelta = 0;
    
    for (let i = 1; i < recentCandles.length; i++) {
        const c = recentCandles[i];
        const priceChange = c.close - recentCandles[i-1].close;
        totalPriceChange += priceChange;
        
        // Delta Analysis - estimate buy vs sell volume
        const range = c.high - c.low;
        const closePosition = range > 0 ? (c.close - c.low) / range : 0.5;
        const vol = c.volume || 0;
        
        let buyVol, sellVol;
        if (c.close > c.open) {
            buyVol = vol * (0.5 + closePosition * 0.5);
            sellVol = vol - buyVol;
            upCount++;
        } else {
            sellVol = vol * (0.5 + (1 - closePosition) * 0.5);
            buyVol = vol - sellVol;
            downCount++;
        }
        
        cumulativeDelta += (buyVol - sellVol);
    }
    
    // Delta confirmation
    const deltaConfirmsTrend = (upCount > downCount && cumulativeDelta > 0) || 
                                (downCount > upCount && cumulativeDelta < 0);
    
    console.log(`📊 Delta Analysis: ${cumulativeDelta > 0 ? 'BUY' : 'SELL'} pressure | Confirms trend: ${deltaConfirmsTrend}`);
    
    // Combine current TF, higher TF, delta, and AI trends with stability filter
    const currentTFTrend = upCount > downCount;
    const deltaTrend = cumulativeDelta > 0;
    
    // Use delta + higher timeframe for stronger signals
    let isUptrend = mtfAnalysis.higherTFTrend !== null 
        ? mtfAnalysis.higherTFTrend
        : (deltaConfirmsTrend ? deltaTrend : currentTFTrend);
    
    // Only override with AI if confidence is high (65%+)
    if (aiPrediction && aiPrediction.confidence > 65) {
        // Check if AI agrees with higher timeframe
        const aiTrend = aiPrediction.signal === 'BUY';
        if (mtfAnalysis.higherTFTrend === null || aiTrend === mtfAnalysis.higherTFTrend) {
            isUptrend = aiTrend;
            console.log(`Using AI signal: ${aiPrediction.signal} (${aiPrediction.confidence}% confidence)`);
        } else {
            console.log(`AI conflicts with higher TF, using higher TF trend`);
        }
    }
    
    // Store previous trend for stability check
    if (!window.previousTrend) {
        window.previousTrend = isUptrend;
    }
    
    // Require strong confirmation to change trend direction
    const trendStrengthRatio = upCount / downCount;
    if (window.previousTrend !== isUptrend) {
        // Only change trend if there's strong confirmation
        if (Math.abs(trendStrengthRatio - 1) < 0.3) {
            // Trend is weak, keep previous trend
            isUptrend = window.previousTrend;
            console.log('Trend change rejected - insufficient confirmation');
        } else {
            // Strong trend change confirmed
            window.previousTrend = isUptrend;
            console.log('Trend change confirmed:', isUptrend ? 'UP' : 'DOWN');
        }
    }
    
    const avgPriceChange = totalPriceChange / (recentCandles.length - 1);
    const trendStrength = Math.abs(avgPriceChange) * mtfAnalysis.confidence;
    
    // Generate predictions: 3 candles for all timeframes
    const numPredictions = 3;
    predictedCandles = [];
    
    // Analyze market structure for ALL timeframes (not just 1m)
    const marketStructure = analyzeMarketStructure(currentData);
    if (marketStructure) {
        console.log(`📊 ${currentInterval} Market Structure:`, {
            phase: marketStructure.phase,
            momentum: marketStructure.momentum.toFixed(2),
            pricePosition: (marketStructure.pricePosition * 100).toFixed(0) + '%',
            swingPoints: `${marketStructure.swingHighs}H / ${marketStructure.swingLows}L`,
            resistance: marketStructure.resistance.toFixed(2),
            support: marketStructure.support.toFixed(2)
        });
    }
    
    // Predict candles with varying heights based on volume and structure
    let previousCandle = currentCandle;
    for (let i = 1; i <= numPredictions; i++) {
        const predictedCandle = generateRealisticPredictedCandle(
            previousCandle, 
            isUptrend, 
            i,
            trendStrength,
            volumeAnalysis,
            marketStructure
        );
        predictedCandles.push(predictedCandle);
        previousCandle = predictedCandle;
    }
    
    console.log(`✅ Generated ${predictedCandles.length} predicted candles:`, predictedCandles.map(c => `${c.close.toFixed(2)}`).join(' → '));
    
    // Store predictions for accuracy tracking
    lastPredictions = {
        time: currentCandle.time,
        predictions: predictedCandles.map(c => ({
            close: c.close,
            high: c.high,
            low: c.low
        })),
        trend: isUptrend ? 'UP' : 'DOWN'
    };
    
    const aiInfo = aiPrediction ? ` | AI: ${aiPrediction.signal} (${aiPrediction.confidence}%)` : '';
    const structureInfo = marketStructure ? ` | Structure: ${marketStructure.phase} @ ${marketStructure.keyLevel}` : '';
    const deltaInfo = ` | Delta: ${cumulativeDelta > 0 ? '+' : ''}${cumulativeDelta.toFixed(0)} ${deltaConfirmsTrend ? '✓' : '✗'}`;
    console.log(`Multi-TF + Delta + Volume + AI | Trend: ${isUptrend ? 'UP' : 'DOWN'} | Confidence: ${(mtfAnalysis.confidence * 100).toFixed(0)}% | Volatility: ${volumeAnalysis.volatilityLevel}${deltaInfo}${aiInfo}${structureInfo}`);
    
    // Update prediction count in legend
    const countEl = document.getElementById('prediction-count');
    if (countEl) {
        countEl.textContent = `${numPredictions} Predictions`;
    }
    
    // Generate trading signal with entry, SL, TP, and RR
    if (window.generateTradingSignal && window.displayTradingSignal) {
        console.log('Generating trading signal...');
        
        // Add multi-timeframe trend filter
        const mtfTrendFilter = {
            higherTFTrend: mtfAnalysis.higherTFTrend,
            confidence: mtfAnalysis.confidence
        };
        
        const tradingSignal = generateTradingSignal(currentData, aiPrediction, volumeAnalysis, mtfTrendFilter);
        console.log('Trading signal generated:', tradingSignal);
        
        // Only update signal if:
        // 1. There's a new signal and no current signal exists, OR
        // 2. The new signal is different from the current one
        const shouldUpdateSignal = !window.currentTradingSignal || 
                                   (tradingSignal && 
                                    (!window.currentTradingSignal || 
                                     tradingSignal.type !== window.currentTradingSignal.type ||
                                     Math.abs(tradingSignal.entry - window.currentTradingSignal.entry) > 1));
        
        // Store signal in history with timestamp and candle index (only if new signal)
        if (tradingSignal && shouldUpdateSignal) {
            historicalSignals.push({
                ...tradingSignal,
                timestamp: Date.now(),
                candleIndex: currentData.length - 1,
                price: currentData[currentData.length - 1].close
            });
            
            // Keep only last 50 signals
            if (historicalSignals.length > 50) {
                historicalSignals.shift();
            }
            
            // Update the display with new signal
            displayTradingSignal(tradingSignal);
            console.log('✅ New signal displayed and locked');
            
            // Send Telegram alert only (browser alerts disabled)
            sendTelegramAlert(tradingSignal);
            
            // Send WhatsApp/Telegram alert (if configured)
            if (window.sendWhatsAppAlert) {
                sendWhatsAppAlert(tradingSignal);
            }
            if (window.sendTelegramAlert) {
                sendTelegramAlert(tradingSignal);
            }
        } else if (!tradingSignal && window.currentTradingSignal) {
            // Signal conditions no longer met - clear it
            displayTradingSignal(null);
            console.log('⏳ Signal cleared - waiting for new setup');
        } else {
            console.log('🔒 Keeping existing signal (SL/TP locked)');
        }
    } else {
        console.error('Trading signal functions not loaded');
    }
}

// Analyze volume to predict volatility and candle sizes
function analyzeVolume() {
    const recentCandles = currentData.slice(-30); // Use more candles for stability
    
    // Calculate average volume
    const avgVolume = recentCandles.reduce((sum, c) => sum + c.volume, 0) / recentCandles.length;
    const currentVolume = currentData[currentData.length - 1].volume;
    
    // Calculate volume ratio
    const volumeRatio = currentVolume / avgVolume;
    
    // Calculate ATR (Average True Range) for more accurate height prediction
    let atr = 0;
    for (let i = 1; i < recentCandles.length; i++) {
        const high = recentCandles[i].high;
        const low = recentCandles[i].low;
        const prevClose = recentCandles[i - 1].close;
        
        const tr = Math.max(
            high - low,
            Math.abs(high - prevClose),
            Math.abs(low - prevClose)
        );
        atr += tr;
    }
    atr = atr / (recentCandles.length - 1);
    
    // Calculate average candle height (last 10 candles for recent volatility)
    const recent10 = recentCandles.slice(-10);
    const avgHeight = recent10.reduce((sum, c) => sum + getCandleHeight(c), 0) / recent10.length;
    
    // Calculate average body size
    const avgBodySize = recent10.reduce((sum, c) => sum + Math.abs(c.close - c.open), 0) / recent10.length;
    const avgBodyRatio = avgBodySize / avgHeight;
    
    // Calculate average wick sizes
    const avgUpperWick = recentCandles.reduce((sum, c) => {
        const bodyTop = Math.max(c.open, c.close);
        return sum + (c.high - bodyTop);
    }, 0) / recentCandles.length;
    
    const avgLowerWick = recentCandles.reduce((sum, c) => {
        const bodyBottom = Math.min(c.open, c.close);
        return sum + (bodyBottom - c.low);
    }, 0) / recentCandles.length;
    
    // Determine volatility level based on ATR
    const atrRatio = atr / avgHeight;
    let volatilityLevel = 'normal';
    if (atrRatio > 1.2 || volumeRatio > 1.5) volatilityLevel = 'high';
    else if (atrRatio < 0.8 || volumeRatio < 0.7) volatilityLevel = 'low';
    
    // Use ATR for more accurate height prediction
    // ATR is more reliable than simple average height
    const predictedHeight = atr * 0.9; // Use 90% of ATR for conservative estimate
    
    return {
        avgHeight: avgHeight,
        atr: atr,
        avgBodyRatio: Math.max(avgBodyRatio, 0.4), // Minimum 40% body
        avgUpperWickRatio: avgUpperWick / avgHeight,
        avgLowerWickRatio: avgLowerWick / avgHeight,
        volumeRatio: volumeRatio,
        volatilityLevel: volatilityLevel,
        predictedHeight: predictedHeight
    };
}

// Analyze multi-timeframe trend for better accuracy
function analyzeMultiTimeframeTrend() {
    const higherTFs = getHigherTimeframes(currentInterval);
    
    if (higherTFs.length === 0) {
        return { higherTFTrend: null, confidence: 1.0 };
    }
    
    let totalUpVotes = 0;
    let totalDownVotes = 0;
    let analyzedTFs = 0;
    let totalConfidence = 0;
    
    // Analyze all available higher timeframes
    higherTFs.forEach(tf => {
        if (!multiTimeframeData[tf] || multiTimeframeData[tf].length < 5) {
            return;
        }
        
        const tfData = multiTimeframeData[tf];
        const recentCandles = tfData.slice(-5);
        
        let upCount = 0;
        let downCount = 0;
        
        recentCandles.forEach(candle => {
            if (candle.close > candle.open) upCount++;
            else downCount++;
        });
        
        // Weight: closer timeframes have more weight
        const weight = 1.0 / (analyzedTFs + 1);
        totalUpVotes += upCount * weight;
        totalDownVotes += downCount * weight;
        
        // Calculate confidence for this timeframe
        const trendStrength = Math.abs(upCount - downCount) / recentCandles.length;
        totalConfidence += trendStrength * weight;
        
        analyzedTFs++;
    });
    
    if (analyzedTFs === 0) {
        return { higherTFTrend: null, confidence: 1.0 };
    }
    
    // Determine overall trend from all timeframes
    const isUptrend = totalUpVotes > totalDownVotes;
    
    // Calculate final confidence (70-100%)
    const avgConfidence = totalConfidence / analyzedTFs;
    const confidence = 0.7 + (avgConfidence * 0.3);
    
    console.log(`Multi-TF Analysis: ${analyzedTFs} timeframes | Trend: ${isUptrend ? 'UP' : 'DOWN'} | Confidence: ${(confidence * 100).toFixed(0)}%`);
    
    return {
        higherTFTrend: isUptrend,
        confidence: confidence
    };
}

// Predict where the current running candle will close (with volume analysis)
function predictRunningCandleWithVolume(currentCandle, isUptrend, trendStrength, volumeAnalysis) {
    const targetHeight = volumeAnalysis.predictedHeight;
    // The running candle already has open, high, low, close
    // We predict where it will actually close based on trend
    
    const currentOpen = currentCandle.open;
    const currentClose = currentCandle.close;
    const currentHigh = currentCandle.high;
    const currentLow = currentCandle.low;
    
    // Calculate how much the candle has moved so far
    const currentMove = currentClose - currentOpen;
    const currentProgress = Math.abs(currentMove) / targetHeight;
    
    // Predict final close based on trend
    let predictedClose;
    if (isUptrend) {
        // In uptrend, expect close to move higher
        const remainingMove = targetHeight * (0.65 - currentProgress);
        predictedClose = currentClose + Math.max(remainingMove, 0);
    } else {
        // In downtrend, expect close to move lower
        const remainingMove = targetHeight * (0.65 - currentProgress);
        predictedClose = currentClose - Math.max(remainingMove, 0);
    }
    
    // Predict final high and low
    let predictedHigh = currentHigh;
    let predictedLow = currentLow;
    
    if (isUptrend) {
        predictedHigh = Math.max(currentHigh, predictedClose + targetHeight * 0.15);
        predictedLow = Math.min(currentLow, currentOpen - targetHeight * 0.20);
    } else {
        predictedLow = Math.min(currentLow, predictedClose - targetHeight * 0.15);
        predictedHigh = Math.max(currentHigh, currentOpen + targetHeight * 0.20);
    }
    
    // Ensure exact height
    const actualHeight = predictedHigh - predictedLow;
    const heightDiff = targetHeight - actualHeight;
    if (Math.abs(heightDiff) > 0.01) {
        predictedHigh += heightDiff / 2;
        predictedLow -= heightDiff / 2;
    }
    
    return {
        open: currentOpen,
        high: predictedHigh,
        low: predictedLow,
        close: predictedClose,
        time: currentCandle.time,
        candleHeight: predictedHigh - predictedLow,
        isRunning: true  // Mark as running candle
    };
}

// Analyze market structure for all timeframes
function analyzeMarketStructure(data) {
    if (data.length < 30) return null;
    
    // Use more candles for higher timeframes
    const lookback = currentInterval === '1m' ? 30 : 
                     currentInterval === '3m' ? 40 :
                     currentInterval === '5m' ? 50 :
                     currentInterval === '15m' ? 60 :
                     currentInterval === '30m' ? 70 :
                     currentInterval === '1h' ? 80 :
                     currentInterval === '4h' ? 100 : 50;
    
    const recent = data.slice(-Math.min(lookback, data.length));
    
    // Identify swing highs and lows (more sensitive for lower timeframes)
    const swingHighs = [];
    const swingLows = [];
    const swingWindow = currentInterval === '1m' ? 2 : 
                       currentInterval === '3m' ? 3 :
                       currentInterval === '5m' ? 3 : 4;
    
    for (let i = swingWindow; i < recent.length - swingWindow; i++) {
        // Swing high: higher than surrounding candles
        let isSwingHigh = true;
        for (let j = 1; j <= swingWindow; j++) {
            if (recent[i].high <= recent[i-j].high || recent[i].high <= recent[i+j].high) {
                isSwingHigh = false;
                break;
            }
        }
        if (isSwingHigh) {
            swingHighs.push({ index: i, price: recent[i].high });
        }
        
        // Swing low: lower than surrounding candles
        let isSwingLow = true;
        for (let j = 1; j <= swingWindow; j++) {
            if (recent[i].low >= recent[i-j].low || recent[i].low >= recent[i+j].low) {
                isSwingLow = false;
                break;
            }
        }
        if (isSwingLow) {
            swingLows.push({ index: i, price: recent[i].low });
        }
    }
    
    // Determine market phase
    const currentPrice = recent[recent.length - 1].close;
    const recentHigh = Math.max(...recent.slice(-10).map(c => c.high));
    const recentLow = Math.min(...recent.slice(-10).map(c => c.low));
    
    let phase = 'RANGING';
    if (swingHighs.length >= 2 && swingLows.length >= 2) {
        const lastTwoHighs = swingHighs.slice(-2);
        const lastTwoLows = swingLows.slice(-2);
        
        // Higher highs and higher lows = uptrend
        if (lastTwoHighs[1].price > lastTwoHighs[0].price && 
            lastTwoLows[1].price > lastTwoLows[0].price) {
            phase = 'UPTREND';
        }
        // Lower highs and lower lows = downtrend
        else if (lastTwoHighs[1].price < lastTwoHighs[0].price && 
                 lastTwoLows[1].price < lastTwoLows[0].price) {
            phase = 'DOWNTREND';
        }
    }
    
    // Calculate momentum (use more candles for higher timeframes)
    const momentumPeriod = currentInterval === '1m' ? 5 :
                          currentInterval === '3m' ? 7 :
                          currentInterval === '5m' ? 10 :
                          currentInterval === '15m' ? 12 : 15;
    
    const momentum = recent.slice(-momentumPeriod).reduce((sum, c, i, arr) => {
        if (i === 0) return 0;
        return sum + (c.close - arr[i-1].close);
    }, 0);
    
    // Identify key levels
    const resistance = swingHighs.length > 0 ? swingHighs[swingHighs.length - 1].price : recentHigh;
    const support = swingLows.length > 0 ? swingLows[swingLows.length - 1].price : recentLow;
    
    // Calculate trend strength
    const trendStrength = Math.abs(momentum) / momentumPeriod;
    
    // Identify if price is at key level
    const pricePosition = (currentPrice - support) / (resistance - support);
    let keyLevel = 'MIDDLE';
    if (pricePosition > 0.9) keyLevel = 'RESISTANCE';
    else if (pricePosition < 0.1) keyLevel = 'SUPPORT';
    else if (pricePosition > 0.7) keyLevel = 'NEAR_RESISTANCE';
    else if (pricePosition < 0.3) keyLevel = 'NEAR_SUPPORT';
    
    return {
        phase: phase,
        momentum: momentum,
        trendStrength: trendStrength,
        resistance: resistance,
        support: support,
        swingHighs: swingHighs.length,
        swingLows: swingLows.length,
        pricePosition: pricePosition, // 0-1 range
        keyLevel: keyLevel,
        timeframe: currentInterval
    };
}

// Generate realistic predicted candle with strategy-based logic
function generateRealisticPredictedCandle(previousCandle, isUptrend, candleNumber, trendStrength, volumeAnalysis, marketStructure) {
    const timeInterval = getIntervalMs(currentInterval);
    
    // Start from previous close
    const open = previousCandle.close;
    
    // Calculate delta-based momentum from recent candles
    const recentCandles = currentData.slice(-5);
    let cumulativeDelta = 0;
    recentCandles.forEach(c => {
        const range = c.high - c.low;
        const closePosition = range > 0 ? (c.close - c.low) / range : 0.5;
        const vol = c.volume || 0;
        
        let buyVol, sellVol;
        if (c.close > c.open) {
            buyVol = vol * (0.5 + closePosition * 0.5);
            sellVol = vol - buyVol;
        } else {
            sellVol = vol * (0.5 + (1 - closePosition) * 0.5);
            buyVol = vol - sellVol;
        }
        cumulativeDelta += (buyVol - sellVol);
    });
    
    // Delta confirms or weakens the trend
    const deltaConfirmsTrend = (isUptrend && cumulativeDelta > 0) || (!isUptrend && cumulativeDelta < 0);
    const deltaStrength = Math.abs(cumulativeDelta) / (recentCandles.reduce((sum, c) => sum + (c.volume || 0), 0) || 1);
    
    // Use volume analysis for realistic body and wick ratios
    let bodyRatio = Math.max(volumeAnalysis.avgBodyRatio, 0.5); // Ensure minimum 50% body
    let upperWickRatio = volumeAnalysis.avgUpperWickRatio;
    let lowerWickRatio = volumeAnalysis.avgLowerWickRatio;
    
    // Adjust based on market structure and delta
    if (marketStructure) {
        // Strong trend with delta confirmation = larger bodies
        if (marketStructure.phase === 'UPTREND' && isUptrend && deltaConfirmsTrend) {
            bodyRatio = Math.min(bodyRatio * 1.5, 0.85); // Stronger moves
            upperWickRatio *= 0.6; // Very small upper wicks
            lowerWickRatio *= 1.3; // Strong buying pressure
        } else if (marketStructure.phase === 'DOWNTREND' && !isUptrend && deltaConfirmsTrend) {
            bodyRatio = Math.min(bodyRatio * 1.5, 0.85);
            lowerWickRatio *= 0.6; // Very small lower wicks
            upperWickRatio *= 1.3; // Strong selling pressure
        }
        // Trend without delta confirmation = weaker moves
        else if (marketStructure.phase === 'UPTREND' && isUptrend && !deltaConfirmsTrend) {
            bodyRatio *= 0.7; // Smaller bodies (weakening)
            upperWickRatio *= 1.5; // Larger rejection wicks
        } else if (marketStructure.phase === 'DOWNTREND' && !isUptrend && !deltaConfirmsTrend) {
            bodyRatio *= 0.7;
            lowerWickRatio *= 1.5;
        }
        // In ranging markets, wicks dominate
        else if (marketStructure.phase === 'RANGING') {
            upperWickRatio *= 1.6;
            lowerWickRatio *= 1.6;
            bodyRatio *= 0.6;
        }
        
        // Near resistance/support with opposing delta = strong rejection
        if (marketStructure.pricePosition > 0.85) {
            if (isUptrend && !deltaConfirmsTrend) {
                upperWickRatio *= 2.5; // Very strong rejection
                bodyRatio *= 0.4; // Small body
            } else if (isUptrend && deltaConfirmsTrend) {
                // Breakout attempt
                bodyRatio *= 1.2;
                upperWickRatio *= 0.8;
            }
        } else if (marketStructure.pricePosition < 0.15) {
            if (!isUptrend && !deltaConfirmsTrend) {
                lowerWickRatio *= 2.5; // Very strong rejection
                bodyRatio *= 0.4;
            } else if (!isUptrend && deltaConfirmsTrend) {
                // Breakdown attempt
                bodyRatio *= 1.2;
                lowerWickRatio *= 0.8;
            }
        }
    }
    
    // Progressive movement based on candle number and delta strength
    // Use ATR-based height for accuracy
    let heightVariation = 1.0;
    
    // First candle: most accurate (very close to ATR)
    if (candleNumber === 1) {
        heightVariation = 0.98 + (deltaStrength * 0.1); // 98-108% - very accurate
    }
    // Second candle: slight variation
    else if (candleNumber === 2) {
        heightVariation = 0.95 + (deltaStrength * 0.15); // 95-110%
    }
    // Third candle: moderate variation
    else {
        heightVariation = 0.92 + (deltaStrength * 0.2); // 92-112%
    }
    
    // Use ATR-based predicted height (more accurate than simple average)
    const targetHeight = volumeAnalysis.predictedHeight * heightVariation;
    
    // Calculate realistic body and wick sizes
    // Body should be a portion of total height, wicks fill the rest
    const bodySize = targetHeight * bodyRatio;
    const upperWickSize = targetHeight * upperWickRatio;
    const lowerWickSize = targetHeight * lowerWickRatio;
    
    // Adjust body size with delta influence (but keep proportions)
    let adjustedBodySize = bodySize;
    if (deltaConfirmsTrend) {
        adjustedBodySize *= (1 + deltaStrength * 0.2); // Up to 20% larger
    } else {
        adjustedBodySize *= (1 - deltaStrength * 0.15); // Up to 15% smaller
    }
    
    let high, low, close;
    
    // Check for key level proximity
    let atKeyLevel = false;
    if (marketStructure) {
        const distanceToResistance = (marketStructure.resistance - open) / open;
        const distanceToSupport = (open - marketStructure.support) / open;
        atKeyLevel = distanceToResistance < 0.005 || distanceToSupport < 0.005; // Within 0.5%
    }
    
    if (isUptrend) {
        // Uptrend: close higher than open
        close = open + adjustedBodySize;
        
        // Check if close would exceed resistance
        if (marketStructure && close > marketStructure.resistance && !deltaConfirmsTrend) {
            // Likely rejection - reduce close
            close = open + (adjustedBodySize * 0.5);
        }
        
        // Add wicks proportionally
        high = close + upperWickSize;
        low = open - lowerWickSize;
        
        // Ensure high is actually the highest point
        high = Math.max(high, close, open);
        low = Math.min(low, close, open);
        
        // Verify total height matches target (adjust if needed)
        const actualHeight = high - low;
        if (Math.abs(actualHeight - targetHeight) > targetHeight * 0.1) {
            // Adjust wicks to match target height
            const heightDiff = targetHeight - actualHeight;
            high += heightDiff * 0.5;
            low -= heightDiff * 0.5;
        }
        
        // Cap at resistance if no breakout momentum
        if (marketStructure && high > marketStructure.resistance && !deltaConfirmsTrend) {
            high = marketStructure.resistance + (upperWickSize * 0.2); // Small overshoot
        }
    } else {
        // Downtrend: close lower than open
        close = open - adjustedBodySize;
        
        // Check if close would break support
        if (marketStructure && close < marketStructure.support && !deltaConfirmsTrend) {
            // Likely rejection - reduce close
            close = open - (adjustedBodySize * 0.5);
        }
        
        // Add wicks proportionally
        low = close - lowerWickSize;
        high = open + upperWickSize;
        
        // Ensure high/low are correct
        high = Math.max(high, close, open);
        low = Math.min(low, close, open);
        
        // Verify total height matches target (adjust if needed)
        const actualHeight = high - low;
        if (Math.abs(actualHeight - targetHeight) > targetHeight * 0.1) {
            // Adjust wicks to match target height
            const heightDiff = targetHeight - actualHeight;
            high += heightDiff * 0.5;
            low -= heightDiff * 0.5;
        }
        
        // Cap at support if no breakdown momentum
        if (marketStructure && low < marketStructure.support && !deltaConfirmsTrend) {
            low = marketStructure.support - (lowerWickSize * 0.2); // Small undershoot
        }
    }
    
    return {
        open: open,
        high: high,
        low: low,
        close: close,
        time: previousCandle.time + (timeInterval * candleNumber),
        candleHeight: high - low,
        volume: previousCandle.volume || 0
    };
}

// Zoom level for candle display (like TradingView)
let zoomLevel = 100; // Number of real candles to show (default - increased to see more history)
const MIN_ZOOM = 5;   // Maximum zoom in (5 candles)
const MAX_ZOOM = 200; // Maximum zoom out (200 candles)

// Handle mouse wheel zoom (TradingView-style)
function handleWheelZoom(event) {
    event.preventDefault();
    
    const delta = Math.sign(event.deltaY);
    
    // Dynamic zoom speed based on current zoom level (faster when zoomed out)
    let zoomSpeed;
    if (zoomLevel < 20) {
        zoomSpeed = 1; // Slow zoom when very zoomed in
    } else if (zoomLevel < 50) {
        zoomSpeed = 2; // Medium zoom
    } else if (zoomLevel < 100) {
        zoomSpeed = 5; // Fast zoom
    } else {
        zoomSpeed = 10; // Very fast zoom when zoomed out
    }
    
    if (delta < 0) {
        // Scroll up = Zoom in (show fewer candles)
        if (zoomLevel > MIN_ZOOM) {
            zoomLevel = Math.max(MIN_ZOOM, zoomLevel - zoomSpeed);
            drawChart();
            showZoomIndicator();
        }
    } else {
        // Scroll down = Zoom out (show more candles)
        if (zoomLevel < MAX_ZOOM) {
            zoomLevel = Math.min(MAX_ZOOM, zoomLevel + zoomSpeed);
            drawChart();
            showZoomIndicator();
        }
    }
}

// Toggle minimize mode
function toggleMinimize() {
    const overlay = document.getElementById('prediction-overlay');
    overlay.classList.toggle('minimized');
}

// Toggle fullscreen mode
function toggleFullscreen() {
    const overlay = document.getElementById('prediction-overlay');
    if (overlay.classList.contains('fullscreen')) {
        overlay.classList.remove('fullscreen');
        console.log('Exited fullscreen');
    } else {
        overlay.classList.add('fullscreen');
        console.log('Entered fullscreen');
    }
    
    // Redraw chart after resize
    setTimeout(() => {
        const canvas = document.getElementById('prediction-chart');
        const rect = canvas.getBoundingClientRect();
        canvas.width = rect.width;
        canvas.height = rect.height;
        drawChart();
    }, 100);
}

// Make toggleMinimize globally accessible
window.toggleMinimize = toggleMinimize;

// Change interval for both charts
function changeInterval(interval) {
    console.log('Changing interval to:', interval);
    currentInterval = interval;
    
    // Clear backtest results when changing timeframe
    if (window.backtestResults) {
        window.backtestResults.totalTrades = 0;
        window.backtestResults.returnPercent = 0;
    }
    
    // Update status to show new timeframe
    updateStatus(`Switched to ${interval} timeframe`);
    
    // Update active button
    document.querySelectorAll('.time-btn').forEach(btn => {
        btn.classList.remove('active');
    });
    event.target.classList.add('active');
    
    // Update TradingView chart
    if (widget && widget.activeChart) {
        try {
            const tvInterval = convertToTradingViewInterval(interval);
            widget.activeChart().setResolution(tvInterval);
            console.log('TradingView interval changed to:', tvInterval);
        } catch (error) {
            console.error('Error changing TradingView interval:', error);
        }
    }
    
    // Fetch new data for prediction chart
    if (isSimulating) {
        fetchRealMarketData();
    }
}

// Check prediction accuracy
function checkPredictionAccuracy() {
    if (!lastPredictions || currentData.length < 3) return;
    
    // Get the last 3 actual candles
    const actualCandles = currentData.slice(-3);
    const predictions = lastPredictions.predictions;
    
    let correctDirection = 0;
    let correctPrice = 0;
    let totalError = 0;
    
    for (let i = 0; i < Math.min(3, predictions.length, actualCandles.length); i++) {
        const predicted = predictions[i];
        const actual = actualCandles[i];
        
        // Check direction accuracy (up/down)
        const predictedTrend = lastPredictions.trend;
        const actualTrend = actual.close > currentData[currentData.length - 4 - i].close ? 'UP' : 'DOWN';
        if (predictedTrend === actualTrend) correctDirection++;
        
        // Check price accuracy (within 0.5%)
        const priceError = Math.abs(predicted.close - actual.close) / actual.close;
        totalError += priceError;
        if (priceError < 0.005) correctPrice++; // Within 0.5%
    }
    
    const directionAccuracy = (correctDirection / 3) * 100;
    const priceAccuracy = (correctPrice / 3) * 100;
    const avgError = (totalError / 3) * 100;
    
    // Store in history
    predictionHistory.push({
        time: Date.now(),
        directionAccuracy,
        priceAccuracy,
        avgError
    });
    
    // Keep only last 20 predictions
    if (predictionHistory.length > 20) {
        predictionHistory.shift();
    }
    
    // Calculate overall accuracy
    if (predictionHistory.length > 0) {
        const avgDirectionAccuracy = predictionHistory.reduce((sum, p) => sum + p.directionAccuracy, 0) / predictionHistory.length;
        const avgPriceAccuracy = predictionHistory.reduce((sum, p) => sum + p.priceAccuracy, 0) / predictionHistory.length;
        const overallError = predictionHistory.reduce((sum, p) => sum + p.avgError, 0) / predictionHistory.length;
        
        console.log(`📊 3-Candle Prediction Accuracy (last ${predictionHistory.length} checks):`);
        console.log(`   Direction: ${avgDirectionAccuracy.toFixed(1)}% | Price: ${avgPriceAccuracy.toFixed(1)}% | Avg Error: ${overallError.toFixed(2)}%`);
        
        // Update UI
        const statusEl = document.getElementById('status');
        if (statusEl) {
            const accuracyText = ` | 📊 Prediction Accuracy: ${avgDirectionAccuracy.toFixed(0)}%`;
            if (!statusEl.textContent.includes('Prediction Accuracy')) {
                statusEl.textContent += accuracyText;
            }
        }
    }
}

// Initialize chart on page load
window.addEventListener('DOMContentLoaded', () => {
    console.log('=== DOM loaded, initializing... ===');
    
    // Browser notifications disabled - using Telegram only
    // requestNotificationPermission();
    
    try {
        initChart();
        console.log('TradingView chart initialized');
    } catch (error) {
        console.error('Error initializing TradingView:', error);
    }
    
    // Initialize order book WebSocket
    if (window.initOrderBookWebSocket) {
        initOrderBookWebSocket(currentSymbol);
    }
    
    // Wait a moment for canvas to be ready
    setTimeout(() => {
        try {
            initCanvas();
            console.log('Canvas initialized');
            
            startSimulation();
            console.log('Simulation start triggered');
            
            // Start countdown timer
            startCountdownTimer();
        } catch (error) {
            console.error('Error in initialization:', error);
        }
    }, 100);
});

// Handle window resize
window.addEventListener('resize', () => {
    if (canvas) {
        const rect = canvas.getBoundingClientRect();
        canvas.width = rect.width || canvas.offsetWidth;
        canvas.height = rect.height || canvas.offsetHeight;
        drawChart();
    }
});

// Show zoom level indicator (TradingView-style)
function showZoomIndicator() {
    const indicator = document.getElementById('zoom-indicator');
    if (!indicator) {
        const div = document.createElement('div');
        div.id = 'zoom-indicator';
        div.style.cssText = 'position: absolute; bottom: 50px; right: 20px; background: rgba(33, 150, 243, 0.95); color: white; padding: 8px 16px; border-radius: 4px; font-size: 13px; font-weight: bold; pointer-events: none; z-index: 9999; box-shadow: 0 2px 8px rgba(0,0,0,0.3);';
        document.querySelector('.prediction-overlay').appendChild(div);
    }
    
    const ind = document.getElementById('zoom-indicator');
    
    // Calculate zoom percentage (5 candles = 100% zoomed in, 200 candles = 0% zoomed in)
    const zoomPercent = Math.round(((MAX_ZOOM - zoomLevel) / (MAX_ZOOM - MIN_ZOOM)) * 100);
    
    // Show zoom level with bar
    ind.innerHTML = `
        <div style="display: flex; align-items: center; gap: 10px;">
            <span>🔍 ${zoomLevel} candles</span>
            <div style="width: 100px; height: 4px; background: rgba(255,255,255,0.3); border-radius: 2px; overflow: hidden;">
                <div style="width: ${zoomPercent}%; height: 100%; background: white; transition: width 0.1s;"></div>
            </div>
        </div>
    `;
    ind.style.display = 'block';
    
    clearTimeout(window.zoomTimeout);
    window.zoomTimeout = setTimeout(() => {
        ind.style.display = 'none';
    }, 1500);
}


// Drawing Tools Functions
function toggleDrawingTool(tool) {
    if (currentDrawingTool === tool) {
        currentDrawingTool = null;
        document.getElementById(`draw-${tool}-btn`)?.classList.remove('active');
    } else {
        // Remove active from all drawing buttons
        document.querySelectorAll('[id^="draw-"]').forEach(btn => btn.classList.remove('active'));
        
        currentDrawingTool = tool;
        document.getElementById(`draw-${tool}-btn`)?.classList.add('active');
    }
    canvas.style.cursor = currentDrawingTool ? 'crosshair' : 'default';
}

function clearDrawings() {
    drawings = [];
    drawChart();
}

// Add mouse event listeners for drawing
function initDrawingEvents() {
    if (!canvas) return;
    
    canvas.addEventListener('mousedown', (e) => {
        if (!currentDrawingTool) return;
        
        const rect = canvas.getBoundingClientRect();
        const x = e.clientX - rect.left;
        const y = e.clientY - rect.top;
        
        isDrawing = true;
        drawingStart = { x, y };
    });
    
    canvas.addEventListener('mousemove', (e) => {
        if (!isDrawing || !drawingStart) return;
        
        const rect = canvas.getBoundingClientRect();
        const x = e.clientX - rect.left;
        const y = e.clientY - rect.top;
        
        // Redraw chart with temporary drawing
        drawChart();
        drawTemporaryShape(drawingStart.x, drawingStart.y, x, y);
    });
    
    canvas.addEventListener('mouseup', (e) => {
        if (!isDrawing || !drawingStart) return;
        
        const rect = canvas.getBoundingClientRect();
        const x = e.clientX - rect.left;
        const y = e.clientY - rect.top;
        
        // Save the drawing
        drawings.push({
            type: currentDrawingTool,
            x1: drawingStart.x,
            y1: drawingStart.y,
            x2: x,
            y2: y
        });
        
        isDrawing = false;
        drawingStart = null;
        drawChart();
    });
    
    canvas.addEventListener('mouseleave', () => {
        if (isDrawing) {
            isDrawing = false;
            drawingStart = null;
            drawChart();
        }
    });
}

function drawTemporaryShape(x1, y1, x2, y2) {
    if (!ctx || !currentDrawingTool) return;
    
    ctx.strokeStyle = '#2962ff';
    ctx.lineWidth = 2;
    ctx.setLineDash([5, 5]);
    
    ctx.beginPath();
    
    if (currentDrawingTool === 'line') {
        ctx.moveTo(x1, y1);
        ctx.lineTo(x2, y2);
    } else if (currentDrawingTool === 'hline') {
        ctx.moveTo(0, y1);
        ctx.lineTo(canvas.width, y1);
    } else if (currentDrawingTool === 'rect') {
        ctx.rect(x1, y1, x2 - x1, y2 - y1);
    }
    
    ctx.stroke();
    ctx.setLineDash([]);
}

function drawSavedDrawings() {
    if (!ctx || drawings.length === 0) return;
    
    drawings.forEach(drawing => {
        ctx.strokeStyle = '#2962ff';
        ctx.lineWidth = 2;
        
        ctx.beginPath();
        
        if (drawing.type === 'line') {
            ctx.moveTo(drawing.x1, drawing.y1);
            ctx.lineTo(drawing.x2, drawing.y2);
        } else if (drawing.type === 'hline') {
            ctx.moveTo(0, drawing.y1);
            ctx.lineTo(canvas.width, drawing.y1);
        } else if (drawing.type === 'rect') {
            ctx.strokeRect(drawing.x1, drawing.y1, drawing.x2 - drawing.x1, drawing.y2 - drawing.y1);
        }
        
        ctx.stroke();
    });
}

// Make functions globally accessible
window.toggleDrawingTool = toggleDrawingTool;
window.clearDrawings = clearDrawings;
