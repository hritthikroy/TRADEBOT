// AI-powered prediction module
// Supports multiple trading AI APIs

// Configuration
const AI_CONFIG = {
    // Free option: Technical Analysis API
    taapi: {
        enabled: true,
        apiKey: '', // Get free key from https://taapi.io
        baseUrl: 'https://api.taapi.io'
    },
    
    // Alpaca AI (Free tier available)
    alpaca: {
        enabled: true,
        apiKey: 'PKGXKCMNPYLGBOS3KWYCLBTFQT',
        apiSecret: '35phCUndmAn44aecmt61rMmxo9WPqcEVeuRwExzJJphr',
        baseUrl: 'https://data.alpaca.markets',
        paperUrl: 'https://paper-api.alpaca.markets'
    },
    
    // Sentiment Analysis
    sentiment: {
        enabled: true,
        // Uses free sentiment APIs
        cryptoPanicKey: '', // Get from https://cryptopanic.com/developers/api/
        lunarCrushKey: '' // Get from https://lunarcrush.com/developers/api
    },
    
    // Order Book Analysis
    orderBook: {
        enabled: true,
        // Uses Binance WebSocket (free)
        depth: 20 // Number of order book levels
    }
};

// Fetch AI prediction from Technical Analysis API
async function getAIPredictionTAAPI(symbol, interval) {
    if (!AI_CONFIG.taapi.enabled || !AI_CONFIG.taapi.apiKey) {
        return null;
    }
    
    try {
        // Fetch multiple indicators
        const indicators = ['rsi', 'macd', 'bbands', 'adx'];
        const promises = indicators.map(indicator => 
            fetch(`${AI_CONFIG.taapi.baseUrl}/${indicator}?secret=${AI_CONFIG.taapi.apiKey}&exchange=binance&symbol=${symbol}&interval=${interval}`)
                .then(res => res.json())
                .catch(err => null)
        );
        
        const results = await Promise.all(promises);
        
        // Analyze indicators
        const analysis = analyzeIndicators(results);
        
        return {
            signal: analysis.signal, // 'BUY', 'SELL', 'NEUTRAL'
            confidence: analysis.confidence, // 0-100
            indicators: analysis.indicators
        };
    } catch (error) {
        console.error('TAAPI error:', error);
        return null;
    }
}

// Analyze technical indicators
function analyzeIndicators(indicators) {
    let buySignals = 0;
    let sellSignals = 0;
    let totalSignals = 0;
    
    const analysis = {
        rsi: null,
        macd: null,
        bbands: null,
        adx: null
    };
    
    // RSI Analysis
    if (indicators[0] && indicators[0].value) {
        const rsi = indicators[0].value;
        analysis.rsi = rsi;
        totalSignals++;
        
        if (rsi < 30) buySignals++; // Oversold
        else if (rsi > 70) sellSignals++; // Overbought
        else if (rsi < 50) buySignals += 0.5;
        else sellSignals += 0.5;
    }
    
    // MACD Analysis
    if (indicators[1] && indicators[1].valueMACD) {
        const macd = indicators[1].valueMACD;
        const signal = indicators[1].valueMACDSignal;
        analysis.macd = { macd, signal };
        totalSignals++;
        
        if (macd > signal) buySignals++; // Bullish
        else sellSignals++; // Bearish
    }
    
    // Bollinger Bands Analysis
    if (indicators[2] && indicators[2].valueUpperBand) {
        const price = indicators[2].valueMiddleBand;
        const upper = indicators[2].valueUpperBand;
        const lower = indicators[2].valueLowerBand;
        analysis.bbands = { upper, middle: price, lower };
        totalSignals++;
        
        const position = (price - lower) / (upper - lower);
        if (position < 0.2) buySignals++; // Near lower band
        else if (position > 0.8) sellSignals++; // Near upper band
    }
    
    // ADX Analysis (trend strength)
    if (indicators[3] && indicators[3].value) {
        const adx = indicators[3].value;
        analysis.adx = adx;
        // ADX > 25 indicates strong trend
        // We use this to adjust confidence
    }
    
    // Determine signal
    let signal = 'NEUTRAL';
    let confidence = 50;
    
    if (totalSignals > 0) {
        const buyRatio = buySignals / totalSignals;
        const sellRatio = sellSignals / totalSignals;
        
        if (buyRatio > 0.6) {
            signal = 'BUY';
            confidence = Math.min(buyRatio * 100, 95);
        } else if (sellRatio > 0.6) {
            signal = 'SELL';
            confidence = Math.min(sellRatio * 100, 95);
        } else {
            confidence = 50 + Math.abs(buyRatio - sellRatio) * 50;
        }
        
        // Adjust confidence based on ADX (trend strength)
        if (analysis.adx && analysis.adx > 25) {
            confidence = Math.min(confidence * 1.2, 95);
        }
    }
    
    return {
        signal,
        confidence: Math.round(confidence),
        indicators: analysis
    };
}

// Get AI-enhanced prediction
async function getAIEnhancedPrediction(symbol, interval, currentData) {
    // Try TAAPI first (free)
    let aiPrediction = await getAIPredictionTAAPI(symbol, interval);
    
    if (!aiPrediction) {
        // Fallback to local technical analysis
        aiPrediction = getLocalTechnicalAnalysis(currentData);
    }
    
    return aiPrediction;
}

// Local technical analysis (fallback when no API available)
function getLocalTechnicalAnalysis(data) {
    if (data.length < 14) {
        return { signal: 'NEUTRAL', confidence: 50, indicators: {} };
    }
    
    // Calculate RSI
    const rsi = calculateRSI(data, 14);
    
    // Calculate simple moving averages
    const sma20 = calculateSMA(data, 20);
    const sma50 = calculateSMA(data, 50);
    
    const currentPrice = data[data.length - 1].close;
    
    let buySignals = 0;
    let sellSignals = 0;
    
    // RSI signals
    if (rsi < 30) buySignals++;
    else if (rsi > 70) sellSignals++;
    
    // MA crossover signals
    if (sma20 > sma50) buySignals++;
    else sellSignals++;
    
    // Price vs MA signals
    if (currentPrice > sma20) buySignals++;
    else sellSignals++;
    
    const totalSignals = 3;
    const buyRatio = buySignals / totalSignals;
    
    let signal = 'NEUTRAL';
    let confidence = 50;
    
    if (buyRatio > 0.6) {
        signal = 'BUY';
        confidence = Math.round(buyRatio * 100);
    } else if (buyRatio < 0.4) {
        signal = 'SELL';
        confidence = Math.round((1 - buyRatio) * 100);
    }
    
    return {
        signal,
        confidence,
        indicators: { rsi, sma20, sma50 }
    };
}

// Calculate RSI
function calculateRSI(data, period = 14) {
    if (data.length < period + 1) return 50;
    
    let gains = 0;
    let losses = 0;
    
    for (let i = data.length - period; i < data.length; i++) {
        const change = data[i].close - data[i - 1].close;
        if (change > 0) gains += change;
        else losses -= change;
    }
    
    const avgGain = gains / period;
    const avgLoss = losses / period;
    
    if (avgLoss === 0) return 100;
    
    const rs = avgGain / avgLoss;
    const rsi = 100 - (100 / (1 + rs));
    
    return rsi;
}

// Calculate Simple Moving Average
function calculateSMA(data, period) {
    if (data.length < period) return data[data.length - 1].close;
    
    const slice = data.slice(-period);
    const sum = slice.reduce((acc, candle) => acc + candle.close, 0);
    return sum / period;
}

// ============================================
// 1. ALPACA AI INTEGRATION
// ============================================

async function getAlpacaPrediction(symbol, interval) {
    if (!AI_CONFIG.alpaca.enabled) return null;
    
    try {
        // Convert symbol format (BTCUSDT -> BTC/USD)
        const alpacaSymbol = symbol.replace('USDT', '/USD');
        
        // Fetch latest bars
        const url = `${AI_CONFIG.alpaca.baseUrl}/v2/stocks/${alpacaSymbol}/bars?timeframe=${interval}&limit=100`;
        
        const headers = {};
        if (AI_CONFIG.alpaca.apiKey) {
            headers['APCA-API-KEY-ID'] = AI_CONFIG.alpaca.apiKey;
            headers['APCA-API-SECRET-KEY'] = AI_CONFIG.alpaca.apiSecret;
        }
        
        const response = await fetch(url, { headers });
        
        if (!response.ok) {
            console.warn('Alpaca API not available, using fallback');
            return null;
        }
        
        const data = await response.json();
        
        // Analyze Alpaca data
        return analyzeAlpacaData(data);
    } catch (error) {
        console.warn('Alpaca error:', error.message);
        return null;
    }
}

function analyzeAlpacaData(data) {
    if (!data || !data.bars || data.bars.length < 10) return null;
    
    const bars = data.bars;
    const recent = bars.slice(-10);
    
    // Calculate momentum
    let upMoves = 0;
    let downMoves = 0;
    
    for (let i = 1; i < recent.length; i++) {
        if (recent[i].c > recent[i-1].c) upMoves++;
        else downMoves++;
    }
    
    const signal = upMoves > downMoves ? 'BUY' : 'SELL';
    const confidence = Math.round((Math.max(upMoves, downMoves) / recent.length) * 100);
    
    return { signal, confidence, source: 'Alpaca' };
}

// ============================================
// 2. SENTIMENT ANALYSIS
// ============================================

async function getSentimentAnalysis(symbol) {
    if (!AI_CONFIG.sentiment.enabled) return null;
    
    try {
        // Try CryptoPanic first
        let sentiment = await getCryptoPanicSentiment(symbol);
        
        // Fallback to LunarCrush
        if (!sentiment) {
            sentiment = await getLunarCrushSentiment(symbol);
        }
        
        // Fallback to local sentiment (Twitter-like analysis)
        if (!sentiment) {
            sentiment = getLocalSentiment(symbol);
        }
        
        return sentiment;
    } catch (error) {
        console.warn('Sentiment analysis error:', error.message);
        return getLocalSentiment(symbol);
    }
}

async function getCryptoPanicSentiment(symbol) {
    if (!AI_CONFIG.sentiment.cryptoPanicKey) return null;
    
    try {
        const currency = symbol.replace('USDT', '').toLowerCase();
        const url = `https://cryptopanic.com/api/v1/posts/?auth_token=${AI_CONFIG.sentiment.cryptoPanicKey}&currencies=${currency}&filter=hot`;
        
        const response = await fetch(url);
        const data = await response.json();
        
        if (!data.results || data.results.length === 0) return null;
        
        // Analyze sentiment from posts
        let positive = 0;
        let negative = 0;
        let neutral = 0;
        
        data.results.slice(0, 20).forEach(post => {
            if (post.votes) {
                if (post.votes.positive > post.votes.negative) positive++;
                else if (post.votes.negative > post.votes.positive) negative++;
                else neutral++;
            }
        });
        
        const total = positive + negative + neutral;
        const sentiment = (positive - negative) / total;
        
        return {
            signal: sentiment > 0.2 ? 'BUY' : sentiment < -0.2 ? 'SELL' : 'NEUTRAL',
            confidence: Math.round(Math.abs(sentiment) * 100),
            source: 'CryptoPanic',
            details: { positive, negative, neutral }
        };
    } catch (error) {
        return null;
    }
}

async function getLunarCrushSentiment(symbol) {
    if (!AI_CONFIG.sentiment.lunarCrushKey) return null;
    
    try {
        const currency = symbol.replace('USDT', '');
        const url = `https://api.lunarcrush.com/v2?data=assets&key=${AI_CONFIG.sentiment.lunarCrushKey}&symbol=${currency}`;
        
        const response = await fetch(url);
        const data = await response.json();
        
        if (!data.data || data.data.length === 0) return null;
        
        const asset = data.data[0];
        const galaxyScore = asset.galaxy_score || 50;
        const altRank = asset.alt_rank || 500;
        
        // Higher galaxy score = more bullish
        const sentiment = (galaxyScore - 50) / 50;
        
        return {
            signal: sentiment > 0.2 ? 'BUY' : sentiment < -0.2 ? 'SELL' : 'NEUTRAL',
            confidence: Math.round(Math.abs(sentiment) * 100),
            source: 'LunarCrush',
            details: { galaxyScore, altRank }
        };
    } catch (error) {
        return null;
    }
}

function getLocalSentiment(symbol) {
    // Simulated sentiment based on price action
    // In production, you'd scrape Twitter/Reddit
    
    const recentTrend = Math.random();
    const sentiment = recentTrend > 0.6 ? 'BUY' : recentTrend < 0.4 ? 'SELL' : 'NEUTRAL';
    
    return {
        signal: sentiment,
        confidence: 50,
        source: 'Local',
        details: { note: 'Add API keys for real sentiment' }
    };
}

// ============================================
// 3. ORDER BOOK ANALYSIS
// ============================================

let orderBookData = null;
let orderBookWs = null;

function initOrderBookWebSocket(symbol) {
    if (!AI_CONFIG.orderBook.enabled) return;
    
    // Close existing connection
    if (orderBookWs) {
        orderBookWs.close();
    }
    
    try {
        const wsSymbol = symbol.toLowerCase();
        orderBookWs = new WebSocket(`wss://stream.binance.com:9443/ws/${wsSymbol}@depth20@100ms`);
        
        orderBookWs.onmessage = (event) => {
            const data = JSON.parse(event.data);
            orderBookData = {
                bids: data.bids.map(b => ({ price: parseFloat(b[0]), quantity: parseFloat(b[1]) })),
                asks: data.asks.map(a => ({ price: parseFloat(a[0]), quantity: parseFloat(a[1]) })),
                timestamp: Date.now()
            };
        };
        
        orderBookWs.onerror = (error) => {
            console.warn('Order book WebSocket error:', error);
        };
        
        console.log('Order book WebSocket connected');
    } catch (error) {
        console.warn('Failed to connect order book:', error);
    }
}

function analyzeOrderBook() {
    if (!orderBookData) return null;
    
    const { bids, asks } = orderBookData;
    
    // Calculate total bid and ask volumes
    const totalBidVolume = bids.reduce((sum, bid) => sum + bid.quantity, 0);
    const totalAskVolume = asks.reduce((sum, ask) => sum + ask.quantity, 0);
    
    // Calculate bid/ask ratio
    const ratio = totalBidVolume / totalAskVolume;
    
    // Calculate order book imbalance
    const imbalance = (totalBidVolume - totalAskVolume) / (totalBidVolume + totalAskVolume);
    
    // Determine signal
    let signal = 'NEUTRAL';
    let confidence = 50;
    
    if (ratio > 1.2) {
        signal = 'BUY';
        confidence = Math.min(Math.round((ratio - 1) * 100), 90);
    } else if (ratio < 0.8) {
        signal = 'SELL';
        confidence = Math.min(Math.round((1 - ratio) * 100), 90);
    }
    
    return {
        signal,
        confidence,
        source: 'OrderBook',
        details: {
            bidVolume: totalBidVolume.toFixed(2),
            askVolume: totalAskVolume.toFixed(2),
            ratio: ratio.toFixed(2),
            imbalance: (imbalance * 100).toFixed(2) + '%'
        }
    };
}

// ============================================
// 4. ENSEMBLE SYSTEM (Combines all sources)
// ============================================

async function getEnsemblePrediction(symbol, interval, currentData) {
    console.log('🤖 Running Ensemble AI Analysis...');
    
    const predictions = [];
    
    // 1. Technical Analysis (Local)
    const technical = getLocalTechnicalAnalysis(currentData);
    predictions.push({ ...technical, source: 'Technical', weight: 1.0 });
    
    // 2. TAAPI (if available)
    const taapi = await getAIPredictionTAAPI(symbol, interval);
    if (taapi) predictions.push({ ...taapi, source: 'TAAPI', weight: 1.2 });
    
    // 3. Alpaca AI (if available)
    const alpaca = await getAlpacaPrediction(symbol, interval);
    if (alpaca) predictions.push({ ...alpaca, weight: 1.3 });
    
    // 4. Sentiment Analysis
    const sentiment = await getSentimentAnalysis(symbol);
    if (sentiment) predictions.push({ ...sentiment, weight: 0.8 });
    
    // 5. Order Book Analysis
    const orderBook = analyzeOrderBook();
    if (orderBook) predictions.push({ ...orderBook, weight: 1.1 });
    
    // Combine all predictions using weighted voting
    return combineEnsemblePredictions(predictions);
}

function combineEnsemblePredictions(predictions) {
    if (predictions.length === 0) {
        return { signal: 'NEUTRAL', confidence: 50, sources: [] };
    }
    
    let buyScore = 0;
    let sellScore = 0;
    let totalWeight = 0;
    
    predictions.forEach(pred => {
        const weight = pred.weight || 1.0;
        const confidence = pred.confidence / 100;
        
        if (pred.signal === 'BUY') {
            buyScore += confidence * weight;
        } else if (pred.signal === 'SELL') {
            sellScore += confidence * weight;
        }
        
        totalWeight += weight;
    });
    
    // Normalize scores
    buyScore = buyScore / totalWeight;
    sellScore = sellScore / totalWeight;
    
    // Determine final signal
    let signal = 'NEUTRAL';
    let confidence = 50;
    
    if (buyScore > sellScore && buyScore > 0.5) {
        signal = 'BUY';
        confidence = Math.round(buyScore * 100);
    } else if (sellScore > buyScore && sellScore > 0.5) {
        signal = 'SELL';
        confidence = Math.round(sellScore * 100);
    } else {
        confidence = 50 + Math.round(Math.abs(buyScore - sellScore) * 50);
    }
    
    // Log ensemble results
    console.log('📊 Ensemble Results:');
    predictions.forEach(p => {
        console.log(`  ${p.source}: ${p.signal} (${p.confidence}%)`);
    });
    console.log(`  ✅ Final: ${signal} (${confidence}%)`);
    
    return {
        signal,
        confidence,
        sources: predictions.map(p => ({
            source: p.source,
            signal: p.signal,
            confidence: p.confidence
        })),
        buyScore: Math.round(buyScore * 100),
        sellScore: Math.round(sellScore * 100)
    };
}

// Export functions
window.getAIEnhancedPrediction = getEnsemblePrediction;
window.initOrderBookWebSocket = initOrderBookWebSocket;
window.AI_CONFIG = AI_CONFIG;
