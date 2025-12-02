# 🏗️ Architecture Documentation

## System Overview

The Trading Bot is a high-performance, real-time cryptocurrency trading system built with Go and PostgreSQL.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                         Frontend                             │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │Dashboard │  │ Signals  │  │Backtest  │  │Analytics │   │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘   │
└───────┼─────────────┼─────────────┼─────────────┼──────────┘
        │             │             │             │
        └─────────────┴─────────────┴─────────────┘
                      │
        ┌─────────────▼─────────────────────────────┐
        │         API Gateway (Fiber)                │
        │  ┌──────────────────────────────────┐     │
        │  │  Middleware Stack                 │     │
        │  │  • Rate Limiting                  │     │
        │  │  • CORS                           │     │
        │  │  • Logging                        │     │
        │  │  • Recovery                       │     │
        │  │  • Validation                     │     │
        │  └──────────────────────────────────┘     │
        └───────────┬───────────────────────────────┘
                    │
        ┌───────────┴───────────────────────────────┐
        │                                            │
┌───────▼────────┐  ┌──────────────┐  ┌────────────▼─────┐
│  REST API      │  │  WebSocket   │  │  Health Checks   │
│  • Signals     │  │  Hub         │  │  • Liveness      │
│  • Backtest    │  │  • Broadcast │  │  • Readiness     │
│  • Analytics   │  │  • Clients   │  │  • Metrics       │
└───────┬────────┘  └──────┬───────┘  └──────────────────┘
        │                  │
        └──────────┬───────┘
                   │
        ┌──────────▼──────────────────────────────┐
        │         Business Logic Layer            │
        │  ┌────────────────────────────────┐    │
        │  │  Signal Generation             │    │
        │  │  • AI Enhanced                 │    │
        │  │  • Pattern Recognition         │    │
        │  │  • Multi-Timeframe Analysis    │    │
        │  └────────────────────────────────┘    │
        │  ┌────────────────────────────────┐    │
        │  │  Backtest Engine               │    │
        │  │  • Historical Data Processing  │    │
        │  │  • Performance Metrics         │    │
        │  │  • Risk Management             │    │
        │  └────────────────────────────────┘    │
        │  ┌────────────────────────────────┐    │
        │  │  Trading Strategies            │    │
        │  │  • ICT/SMC Concepts            │    │
        │  │  • Order Flow Analysis         │    │
        │  │  • Liquidity Detection         │    │
        │  └────────────────────────────────┘    │
        └─────────────────┬───────────────────────┘
                          │
        ┌─────────────────▼───────────────────────┐
        │         Data Access Layer               │
        │  ┌────────────────────────────────┐    │
        │  │  Database (PostgreSQL)         │    │
        │  │  • Connection Pooling          │    │
        │  │  • Retry Logic                 │    │
        │  │  • Health Monitoring           │    │
        │  └────────────────────────────────┘    │
        │  ┌────────────────────────────────┐    │
        │  │  External APIs                 │    │
        │  │  • Binance                     │    │
        │  │  • Grok AI                     │    │
        │  └────────────────────────────────┘    │
        └─────────────────────────────────────────┘
```

## Component Details

### 1. API Gateway (Fiber Framework)

**Responsibilities:**
- HTTP request routing
- Middleware execution
- WebSocket management
- Static file serving

**Key Features:**
- Rate limiting (100 req/min)
- CORS protection
- Request validation
- Error handling
- Panic recovery

### 2. Signal Generation System

**Components:**
- `signal_generator.go` - Base signal generation
- `ai_enhanced_signal_generator.go` - AI-powered signals
- `advanced_signal_generator.go` - Multi-factor analysis

**Process Flow:**
```
Market Data → Pattern Recognition → ICT Analysis → 
Confluence Check → Risk Calculation → Signal Output
```

**Confirmations Required:**
- Minimum 4+ factors
- Risk/Reward > 1.8:1
- Session alignment
- Volume confirmation

### 3. Backtest Engine

**Features:**
- Historical data simulation
- Multiple timeframe support
- Performance metrics calculation
- Trade execution simulation

**Metrics Tracked:**
- Win rate
- Profit factor
- Max drawdown
- Sharpe ratio
- Average R-multiple

### 4. Database Layer

**Schema:**
- `trading_signals` - Signal storage
- `signal_analytics` - Performance view

**Optimizations:**
- Connection pooling (25 max, 5 idle)
- Indexed queries
- Prepared statements
- Health monitoring

### 5. WebSocket System

**Architecture:**
- Hub-based broadcasting
- Client management
- Connection limits (1000 max)
- Automatic cleanup

**Message Types:**
- New signals
- Signal updates
- Market data
- System alerts

## Data Flow

### Signal Generation Flow

```
1. Market Data Fetch
   ↓
2. Technical Analysis
   • Candlestick patterns
   • Support/Resistance
   • Trend analysis
   ↓
3. ICT/SMC Analysis
   • Order blocks
   • Fair value gaps
   • Liquidity zones
   ↓
4. Confluence Calculation
   • Score each factor
   • Weight by importance
   • Filter by threshold
   ↓
5. Risk Management
   • Calculate stop loss
   • Set take profits
   • Position sizing
   ↓
6. Signal Output
   • Store in database
   • Broadcast via WebSocket
   • Return via API
```

### Backtest Flow

```
1. Request Validation
   ↓
2. Historical Data Fetch
   ↓
3. Candle-by-Candle Simulation
   • Generate signals
   • Execute trades
   • Track positions
   ↓
4. Performance Calculation
   • Win/loss tracking
   • Profit calculation
   • Metrics aggregation
   ↓
5. Results Return
```

## Security Architecture

### Layers of Protection

1. **Network Layer**
   - CORS restrictions
   - Rate limiting
   - Request size limits

2. **Application Layer**
   - Input validation
   - SQL injection prevention
   - Error sanitization

3. **Data Layer**
   - Encrypted connections
   - Environment variables
   - No hardcoded secrets

## Performance Optimizations

### Backend
- Goroutine pooling
- Connection pooling
- Efficient data structures
- Minimal allocations

### Database
- Indexed queries
- Prepared statements
- Connection reuse
- Query optimization

### WebSocket
- Buffered channels
- Non-blocking sends
- Automatic cleanup
- Connection limits

## Scalability Considerations

### Horizontal Scaling
- Stateless API design
- Database connection pooling
- WebSocket hub per instance

### Vertical Scaling
- Efficient memory usage
- CPU-bound optimizations
- Goroutine management

## Monitoring & Observability

### Health Checks
- `/api/v1/health` - Detailed health
- `/api/v1/ready` - Readiness probe
- `/api/v1/live` - Liveness probe

### Metrics
- Request latency
- Error rates
- Database stats
- WebSocket connections
- Memory usage
- Goroutine count

### Logging
- Structured logging
- Request tracing
- Error tracking
- Performance logs

## Deployment Architecture

### Docker
```
Container:
  - Go binary
  - Static files
  - Non-root user
  - Health checks
```

### Environment Variables
```
Required:
  - SUPABASE_HOST
  - SUPABASE_PASSWORD

Optional:
  - PORT
  - ALLOWED_ORIGINS
  - GROK_API_KEY
```

## Technology Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Fiber v2
- **Database**: PostgreSQL (Supabase)
- **WebSocket**: Fiber WebSocket

### Frontend
- **HTML5**: Semantic markup
- **CSS3**: Modern styling
- **JavaScript**: Vanilla JS
- **WebSocket**: Native API

### Infrastructure
- **Container**: Docker
- **CI/CD**: GitHub Actions
- **Hosting**: VPS/Cloud

## Design Patterns

### Used Patterns
1. **Hub Pattern** - WebSocket management
2. **Repository Pattern** - Data access
3. **Strategy Pattern** - Trading strategies
4. **Factory Pattern** - Signal generation
5. **Middleware Pattern** - Request processing

## Error Handling Strategy

### Levels
1. **Panic Recovery** - Goroutine level
2. **Error Returns** - Function level
3. **HTTP Errors** - API level
4. **User Messages** - Frontend level

### Logging
- Errors logged with context
- Request IDs for tracing
- Structured log format
- UTC timestamps

## Future Enhancements

### Planned
- [ ] Prometheus metrics
- [ ] Distributed tracing
- [ ] Redis caching
- [ ] Message queue (NATS/RabbitMQ)
- [ ] Multi-region deployment
- [ ] Advanced analytics dashboard

### Under Consideration
- [ ] GraphQL API
- [ ] gRPC for internal services
- [ ] Kubernetes deployment
- [ ] Service mesh integration

---

**Last Updated**: December 2024  
**Version**: 1.0.0
