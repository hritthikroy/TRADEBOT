# Deploy AI Analytics Backend

## Prerequisites
- Go 1.21 or higher
- PostgreSQL database (Supabase)
- Environment variables configured

## Installation

### 1. Install Go (if not installed)
```bash
# macOS
brew install go

# Or download from: https://golang.org/dl/
```

### 2. Install Dependencies
```bash
cd backend
go mod tidy
```

### 3. Configure Environment
Create `.env` file in backend directory:
```env
DATABASE_URL=postgresql://user:password@host:port/database
PORT=8080
```

For Supabase, use your connection string from Supabase dashboard.

### 4. Run Backend
```bash
go run .
```

Or build and run:
```bash
go build -o tradebot-api
./tradebot-api
```

## Testing AI Analytics

### 1. Test Health Endpoint
```bash
curl http://localhost:8080/api/v1/health
```

### 2. Test AI Analytics
```bash
curl http://localhost:8080/api/v1/analytics/ai
```

### 3. Open Dashboard
Open `ai-analytics.html` in your browser

## API Endpoints

### AI Analytics
- `GET /api/v1/analytics/ai` - Get comprehensive AI analysis

### Performance Stats
- `GET /api/v1/analytics/performance` - Overall performance
- `GET /api/v1/analytics/by-killzone` - Stats by kill zone
- `GET /api/v1/analytics/by-pattern` - Stats by pattern

### Signals
- `GET /api/v1/signals` - Get all signals
- `GET /api/v1/signals/pending` - Get pending signals
- `POST /api/v1/signals` - Create signal
- `PUT /api/v1/signals/:id` - Update signal
- `DELETE /api/v1/signals/:id` - Delete signal

## Deployment Options

### Option 1: Fly.io
```bash
fly launch
fly deploy
```

### Option 2: Render
1. Connect GitHub repo
2. Select "Go" environment
3. Build command: `go build -o app`
4. Start command: `./app`

### Option 3: Railway
1. Connect GitHub repo
2. Auto-detects Go
3. Deploys automatically

### Option 4: Local Development
```bash
go run .
```

## Environment Variables

Required:
- `DATABASE_URL` - PostgreSQL connection string
- `PORT` - Server port (default: 8080)

Optional:
- `GO_ENV` - Environment (development/production)

## Troubleshooting

### Port Already in Use
```bash
# Kill process on port 8080
lsof -ti:8080 | xargs kill -9
```

### Database Connection Error
- Check DATABASE_URL format
- Verify Supabase credentials
- Ensure database is accessible

### CORS Issues
Backend already configured with CORS:
```go
AllowOrigins: "*"
```

## Performance Tips

1. **Database Indexing**: Ensure indexes on:
   - `created_at`
   - `status`
   - `kill_zone`
   - `pattern_type`

2. **Caching**: Consider adding Redis for analytics caching

3. **Connection Pooling**: Already configured in database.go

## Security

1. **API Keys**: Add authentication middleware
2. **Rate Limiting**: Add rate limiter
3. **HTTPS**: Use TLS in production
4. **Environment**: Never commit .env file

## Monitoring

### Logs
Backend logs all requests via Fiber logger middleware

### Health Check
```bash
curl http://localhost:8080/api/v1/health
```

### Database Stats
Check Supabase dashboard for:
- Query performance
- Connection count
- Storage usage

## Next Steps

1. Deploy backend to cloud
2. Update `ai-analytics.html` with production URL
3. Set up monitoring
4. Configure alerts
5. Add authentication

## Support

For issues:
1. Check backend logs
2. Verify database connection
3. Test endpoints with curl
4. Check browser console
