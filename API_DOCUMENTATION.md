# 📡 API Documentation

Complete API reference for the Trading Bot.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

Currently no authentication required. Add API keys in production.

## Rate Limiting

- **Limit**: 100 requests per minute per IP
- **Response**: 429 Too Many Requests when exceeded

## Health Endpoints

### GET /health

Detailed health check with system metrics.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-12-02T10:30:00Z",
  "uptime": "2h15m30s",
  "version": "1.0.0",
  "database": {
    "status": "connected",
    "open_connections": 5,
    "in_use": 2,
    "idle": 3
  },
  "system": {
    "go_version": "go1.21",
    "goroutines": 45,
    "memory_mb": 128,
    "num_cpu": 8
  }
}
```

### GET /ready

Readiness probe for load balancers.

**Response:**
```json
{
  "status": "ready"
}
```

### GET /live

Liveness probe for health monitoring.

**Response:**
```json
{
  "status": "alive"
}
```

## Backtest Endpoints

### POST /backtest/run

Run a historical backtest.

**Request:**
```json
{
  "symbol": "BTCUSDT",
  "interval": "15m",
  "days": 30,
  "startBalance": 500
}
```

**Parameters:**
- `symbol` (string, required): Trading pair (must end in USDT)
- `interval` (string, required): Timeframe (1m, 5m, 15m, 1h, 4h, 1d)
- `days` (integer, required): Historical days (1-365)
- `startBalance` (float, required): Starting capital (10-1000000)

**Response:**
```json
{
  "symbol": "BTCUSDT",
  "interval": "15m",
  "days": 30,
  "startBalance": 500,
  "endBalance": 687.50,
  "totalTrades": 45,
  "winningTrades": 32,
  "losingTrades": 13,
  "winRate": 71.11,
  "profitFactor": 2.15,
  "returnPercent": 37.50,
  "maxDrawdown": 8.25,
  "averageRR": 2.1,
  "trades": [...],
  "exitReasons": {
    "TP1": 15,
    "TP2": 10,
    "TP3": 7,
    "Trailing Stop": 8,
    "Stop Loss": 5
  }
}
```

### POST /backtest/enhanced

Run backtest with external signal validation.

**Request:**
```json
{
  "symbol": "BTCUSDT",
  "interval": "15m",
  "days": 30,
  "startBalance": 500,
  "useExternalSignals": true
}
```

## Signal Endpoints

### POST /signals

Create a new signal.

**Request:**
```json
{
  "signal_type": "BUY",
  "symbol": "BTCUSDT",
  "timeframe": "15m",
  "entry_price": 91234.56,
  "stop_loss": 90800.00,
  "tp1": 92000.00,
  "tp2": 92500.00,
  "tp3": 93000.00,
  "strength": 78.5,
  "pattern_type": "Order Block",
  "kill_zone": "London"
}
```

**Response:**
```json
{
  "id": 123,
  "signal_id": "uuid-here",
  "created_at": "2024-12-02T10:30:00Z",
  ...
}
```

### GET /signals

Get all signals with optional filters.

**Query Parameters:**
- `status` (string): Filter by status (pending, win, loss)
- `symbol` (string): Filter by symbol
- `timeframe` (string): Filter by timeframe
- `limit` (integer): Max results (default: 100)
- `offset` (integer): Pagination offset

**Response:**
```json
{
  "signals": [...],
  "total": 150,
  "limit": 100,
  "offset": 0
}
```

### GET /signals/pending

Get only pending (active) signals.

**Response:**
```json
{
  "signals": [...],
  "count": 5
}
```

### GET /signals/:id

Get a specific signal by ID.

**Response:**
```json
{
  "id": 123,
  "signal_id": "uuid-here",
  ...
}
```

### PUT /signals/:id

Update a signal.

**Request:**
```json
{
  "status": "win",
  "exit_price": 92150.00,
  "exit_reason": "TP2",
  "profit_percent": 2.15
}
```

### PUT /signals/:id/live-price

Update live price for tracking.

**Request:**
```json
{
  "live_price": 91500.00
}
```

### DELETE /signals/:id

Delete a signal.

**Response:**
```json
{
  "message": "Signal deleted successfully"
}
```

## Analytics Endpoints

### GET /analytics

Get overall analytics.

**Response:**
```json
{
  "total_signals": 150,
  "pending": 5,
  "wins": 105,
  "losses": 40,
  "win_rate": 72.41,
  "avg_profit": 2.35,
  "avg_loss": -1.12,
  "profit_factor": 2.10
}
```

### GET /analytics/performance

Get detailed performance statistics.

**Response:**
```json
{
  "daily": {...},
  "weekly": {...},
  "monthly": {...},
  "all_time": {...}
}
```

### GET /analytics/by-killzone

Get statistics grouped by kill zone.

**Response:**
```json
{
  "London": {
    "total": 50,
    "win_rate": 75.0,
    "avg_rr": 2.3
  },
  "NewYork": {...},
  "Asian": {...}
}
```

### GET /analytics/by-pattern

Get statistics grouped by pattern type.

**Response:**
```json
{
  "Order Block": {
    "total": 30,
    "win_rate": 80.0
  },
  "Fair Value Gap": {...}
}
```

## Filter Endpoints

### GET /filters/opportunities

Get best trade opportunities based on current market.

**Response:**
```json
{
  "opportunities": [
    {
      "symbol": "BTCUSDT",
      "score": 85.5,
      "reasons": ["Strong OB", "FVG present", "London session"]
    }
  ]
}
```

### GET /filters/rules

Get current trading rules and filters.

**Response:**
```json
{
  "min_confirmations": 4,
  "min_rr": 1.8,
  "max_risk_percent": 2.0,
  "allowed_sessions": ["London", "NewYork"]
}
```

### GET /filters/smart

Get smart filter criteria.

**Response:**
```json
{
  "filters": {
    "pattern_confidence": ">= 70",
    "confluence_score": ">= 4",
    "session": "Kill zones only"
  }
}
```

## AI Endpoints

### GET /ai/stats

Get AI enhancement statistics.

**Response:**
```json
{
  "enabled": true,
  "total_analyzed": 150,
  "filtered_out": 25,
  "filter_rate": 16.67,
  "avg_confidence": 78.5
}
```

### POST /ai/toggle

Enable or disable AI filtering.

**Request:**
```json
{
  "enabled": true
}
```

### GET /ai/test

Test Grok AI connection.

**Response:**
```json
{
  "status": "connected",
  "latency_ms": 250,
  "api_version": "v1"
}
```

### POST /ai/sentiment

Analyze market sentiment for a symbol.

**Request:**
```json
{
  "symbol": "BTCUSDT"
}
```

**Response:**
```json
{
  "symbol": "BTCUSDT",
  "sentiment": "bullish",
  "confidence": 75.5,
  "factors": [
    "Strong buying pressure",
    "Breaking resistance"
  ]
}
```

## External Signal Endpoints

### POST /external-signals/get

Get signals from external providers.

**Request:**
```json
{
  "symbol": "BTCUSDT",
  "providers": ["provider1", "provider2"]
}
```

### POST /external-signals/enhanced

Get enhanced signal with multiple sources.

**Request:**
```json
{
  "symbol": "BTCUSDT",
  "timeframe": "15m"
}
```

### POST /external-signals/compare

Compare signals from multiple sources.

**Request:**
```json
{
  "symbol": "BTCUSDT",
  "signals": [...]
}
```

### GET /external-signals/providers

List available free signal providers.

**Response:**
```json
{
  "providers": [
    {
      "name": "Provider 1",
      "status": "active",
      "rate_limit": "100/day"
    }
  ]
}
```

## WebSocket

### Connection

```javascript
const ws = new WebSocket('ws://localhost:8080/ws/signals');

ws.onopen = () => {
  console.log('Connected');
};

ws.onmessage = (event) => {
  const signals = JSON.parse(event.data);
  console.log('New signals:', signals);
};

ws.onerror = (error) => {
  console.error('WebSocket error:', error);
};

ws.onclose = () => {
  console.log('Disconnected');
};
```

### Message Format

```json
[
  {
    "type": "BUY",
    "entry": 91234.56,
    "stopLoss": 90800.00,
    "targets": [
      {"price": 92000.00, "rr": 1.76}
    ],
    "strength": 78.5,
    "timeframe": "15m"
  }
]
```

## Error Responses

### 400 Bad Request
```json
{
  "error": "Validation failed",
  "code": 400,
  "path": "/api/v1/backtest/run",
  "message": "days must be between 1 and 365"
}
```

### 429 Too Many Requests
```json
{
  "error": "Rate limit exceeded. Please try again later."
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal Server Error",
  "code": 500,
  "path": "/api/v1/signals",
  "message": "database connection failed"
}
```

### 503 Service Unavailable
```json
{
  "status": "not_ready",
  "reason": "database_unavailable"
}
```

## Examples

### cURL Examples

```bash
# Health check
curl http://localhost:8080/api/v1/health

# Run backtest
curl -X POST http://localhost:8080/api/v1/backtest/run \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "BTCUSDT",
    "interval": "15m",
    "days": 30,
    "startBalance": 500
  }'

# Get pending signals
curl http://localhost:8080/api/v1/signals/pending

# Get analytics
curl http://localhost:8080/api/v1/analytics
```

### JavaScript Examples

```javascript
// Run backtest
const response = await fetch('http://localhost:8080/api/v1/backtest/run', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    symbol: 'BTCUSDT',
    interval: '15m',
    days: 30,
    startBalance: 500
  })
});

const result = await response.json();
console.log(result);

// Get signals
const signals = await fetch('http://localhost:8080/api/v1/signals/pending')
  .then(res => res.json());
```

---

**Last Updated**: December 2024  
**Version**: 1.0.0
