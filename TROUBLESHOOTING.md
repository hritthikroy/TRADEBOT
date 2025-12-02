# 🔧 Troubleshooting Guide

## Common Issues and Solutions

### 1. Server Won't Start

**Symptom**: Server fails to start or crashes immediately

**Possible Causes:**
- Missing environment variables
- Port already in use
- Database connection failure

**Solutions:**
```bash
# Check if port is in use
lsof -i :8080

# Kill process using port
kill -9 <PID>

# Verify environment variables
cat backend/.env

# Check logs
tail -f backend/server.log
```

### 2. Database Connection Failed

**Symptom**: "Database connection failed" errors

**Solutions:**
```bash
# Verify credentials
echo $SUPABASE_HOST
echo $SUPABASE_PASSWORD

# Test connection manually
psql "postgresql://postgres:PASSWORD@HOST:5432/postgres?sslmode=require"

# Check Supabase dashboard
# Ensure database is running and accessible
```

### 3. No Signals Generated

**Symptom**: Backtest returns 0 trades

**Possible Causes:**
- Filters too strict
- Market conditions don't meet criteria
- Wrong timeframe/symbol

**Solutions:**
```bash
# Try different timeframes
curl -X POST http://localhost:8080/api/v1/backtest/run \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTCUSDT","interval":"4h","days":90,"startBalance":500}'

# Check signal generation logs
grep "Signal" backend/server.log
```

### 4. WebSocket Not Connecting

**Symptom**: Frontend shows "Disconnected"

**Solutions:**
```javascript
# Check browser console for errors
# Verify WebSocket URL
ws://localhost:8080/ws

# Test WebSocket manually
wscat -c ws://localhost:8080/ws
```

### 5. High Memory Usage

**Symptom**: Server using too much RAM

**Solutions:**
```bash
# Check memory stats
curl http://localhost:8080/api/v1/health

# Restart server
pkill trading-bot
./trading-bot

# Reduce connection pool
# Edit database.go: SetMaxOpenConns(10)
```

### 6. Rate Limit Errors

**Symptom**: 429 Too Many Requests

**Solutions:**
```bash
# Wait 1 minute
# Or increase rate limit in main.go

# Check current rate limit
grep "Max:" backend/main.go
```

### 7. Docker Build Fails

**Symptom**: Docker build errors

**Solutions:**
```bash
# Clean Docker cache
docker system prune -a

# Rebuild from scratch
cd backend
docker build --no-cache -t trading-bot .

# Check Dockerfile syntax
docker build --dry-run -t trading-bot .
```

### 8. Tests Failing

**Symptom**: `go test` fails

**Solutions:**
```bash
cd backend

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -v -run TestHealthEndpoint

# Check test coverage
go test -cover ./...
```

### 9. CORS Errors

**Symptom**: Browser blocks requests

**Solutions:**
```bash
# Add your domain to ALLOWED_ORIGINS
echo "ALLOWED_ORIGINS=http://localhost:8080,https://yourdomain.com" >> .env

# Restart server
```

### 10. Slow Performance

**Symptom**: API responses are slow

**Solutions:**
```bash
# Check health metrics
curl http://localhost:8080/api/v1/health

# Monitor goroutines
# If > 1000, there may be a leak

# Check database connections
# Should be < 25 open connections

# Profile the application
go tool pprof http://localhost:8080/debug/pprof/profile
```

## Debugging Commands

### Check Server Status
```bash
# Is server running?
ps aux | grep trading-bot

# Check port
netstat -an | grep 8080

# Test health endpoint
curl http://localhost:8080/api/v1/health
```

### View Logs
```bash
# Real-time logs
tail -f backend/server.log

# Search for errors
grep "ERROR\|FATAL" backend/server.log

# Last 100 lines
tail -n 100 backend/server.log
```

### Database Debugging
```bash
# Connect to database
psql "postgresql://postgres:PASSWORD@HOST:5432/postgres?sslmode=require"

# Check signals table
SELECT COUNT(*) FROM trading_signals;

# View recent signals
SELECT * FROM trading_signals ORDER BY created_at DESC LIMIT 10;

# Check analytics
SELECT * FROM signal_analytics;
```

### Performance Profiling
```bash
# CPU profile
go tool pprof http://localhost:8080/debug/pprof/profile

# Memory profile
go tool pprof http://localhost:8080/debug/pprof/heap

# Goroutine profile
go tool pprof http://localhost:8080/debug/pprof/goroutine
```

## Error Messages

### "Database credentials not configured"
- Set SUPABASE_HOST and SUPABASE_PASSWORD in .env

### "Rate limit exceeded"
- Wait 1 minute or increase RATE_LIMIT_MAX

### "Invalid interval"
- Use: 1m, 5m, 15m, 1h, 4h, 1d

### "Symbol is required"
- Provide symbol in request body

### "Only USDT pairs are supported"
- Use symbols ending in USDT (e.g., BTCUSDT)

## Getting Help

1. Check logs: `tail -f backend/server.log`
2. Run health check: `curl http://localhost:8080/api/v1/health`
3. Run tests: `./test_all_features.sh`
4. Check documentation: See INDEX.md
5. Review architecture: See ARCHITECTURE.md

## Prevention

### Before Deployment
- [ ] Run all tests
- [ ] Check environment variables
- [ ] Test database connection
- [ ] Verify CORS settings
- [ ] Review logs for warnings

### Monitoring
- [ ] Set up health check monitoring
- [ ] Monitor memory usage
- [ ] Track error rates
- [ ] Watch database connections
- [ ] Monitor API latency

---

**Last Updated**: December 2024
