# 🚀 Quick Reference Guide

## Essential Commands

### Start Server
```bash
cd backend
./trading-bot
```

### Build
```bash
make build
```

### Run Tests
```bash
make test
```

### Check Health
```bash
curl http://localhost:8080/api/v1/health
```

### Run Backtest
```bash
curl -X POST http://localhost:8080/api/v1/backtest/run \
  -H "Content-Type: application/json" \
  -d '{"symbol":"BTCUSDT","interval":"15m","days":30,"startBalance":500}'
```

## Key Endpoints

- **Health**: `GET /api/v1/health`
- **Backtest**: `POST /api/v1/backtest/run`
- **Signals**: `GET /api/v1/signals`
- **WebSocket**: `ws://localhost:8080/ws/signals`

## Configuration

Edit `backend/.env`:
```env
SUPABASE_HOST=your-host.supabase.co
SUPABASE_PASSWORD=your-password
PORT=8080
ALLOWED_ORIGINS=http://localhost:8080
```

## Documentation

- **README.md** - Start here
- **API_DOCUMENTATION.md** - API reference
- **TROUBLESHOOTING.md** - Fix issues
- **ARCHITECTURE.md** - System design

## Common Issues

### Server won't start
```bash
# Check if port is in use
lsof -i :8080
# Kill process
kill -9 <PID>
```

### Database connection failed
```bash
# Verify credentials in .env
cat backend/.env
```

### No signals generated
```bash
# Try different timeframe
# 4h usually has better results
```

## Development

```bash
make dev      # Run with auto-reload
make lint     # Check code quality
make fmt      # Format code
make clean    # Clean build artifacts
```

## Docker

```bash
make docker-build    # Build image
make docker-run      # Run container
make docker-stop     # Stop container
```

## Testing

```bash
make test              # Run all tests
make test-coverage     # With coverage report
./test_all_features.sh # Integration tests
```

---

**Need Help?** See TROUBLESHOOTING.md
