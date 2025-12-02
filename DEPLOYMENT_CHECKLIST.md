# 🚀 Deployment Checklist

Use this checklist before deploying to production.

## Pre-Deployment

### Code Quality
- [ ] All tests passing (`make test`)
- [ ] Linter passing (`make lint`)
- [ ] Code formatted (`make fmt`)
- [ ] No TODO/FIXME in critical paths
- [ ] Error handling reviewed
- [ ] Logging configured properly

### Security
- [ ] Environment variables set (not hardcoded)
- [ ] CORS configured for production domains
- [ ] Rate limiting enabled
- [ ] Database credentials secured
- [ ] No sensitive data in logs
- [ ] API keys rotated
- [ ] HTTPS enabled (if applicable)

### Configuration
- [ ] `.env` file created from `.env.example`
- [ ] `SUPABASE_HOST` set
- [ ] `SUPABASE_PASSWORD` set
- [ ] `PORT` configured
- [ ] `ALLOWED_ORIGINS` set to production domains
- [ ] `GROK_API_KEY` set (if using AI)

### Database
- [ ] Database schema deployed (`supabase-complete-setup.sql`)
- [ ] Database connection tested
- [ ] Indexes created
- [ ] RLS policies enabled
- [ ] Backup strategy in place

### Performance
- [ ] Connection pooling configured
- [ ] Rate limits appropriate for load
- [ ] Memory limits set
- [ ] CPU limits set (if containerized)
- [ ] Caching strategy reviewed

## Deployment Steps

### 1. Build
```bash
make build
# or
make docker-build
```

### 2. Test Build
```bash
# Local test
./backend/trading-bot

# Docker test
make docker-run
```

### 3. Verify Health
```bash
curl http://localhost:8080/api/v1/health
curl http://localhost:8080/api/v1/ready
```

### 4. Run Integration Tests
```bash
./test_all_features.sh
```

### 5. Deploy
```bash
# Docker
docker push your-registry/trading-bot:latest

# VPS
scp backend/trading-bot user@server:/app/
ssh user@server 'systemctl restart trading-bot'

# Cloud Platform
# Follow platform-specific instructions
```

## Post-Deployment

### Immediate Checks (0-5 minutes)
- [ ] Server started successfully
- [ ] Health endpoint responding
- [ ] Database connected
- [ ] WebSocket connections working
- [ ] No error logs

### Short-term Monitoring (5-30 minutes)
- [ ] API response times normal
- [ ] Memory usage stable
- [ ] CPU usage normal
- [ ] No connection leaks
- [ ] Signals generating correctly

### Long-term Monitoring (1-24 hours)
- [ ] No memory leaks
- [ ] Database performance stable
- [ ] Error rate acceptable
- [ ] User feedback positive
- [ ] Metrics within expected ranges

## Rollback Plan

If issues occur:

### Quick Rollback
```bash
# Docker
docker pull your-registry/trading-bot:previous
docker stop trading-bot
docker run -d --name trading-bot your-registry/trading-bot:previous

# VPS
ssh user@server 'systemctl stop trading-bot'
ssh user@server 'cp /app/trading-bot.backup /app/trading-bot'
ssh user@server 'systemctl start trading-bot'
```

### Database Rollback
```bash
# Restore from backup
psql "postgresql://..." < backup.sql
```

## Monitoring Setup

### Health Checks
```bash
# Add to monitoring system
curl http://your-domain.com/api/v1/health
curl http://your-domain.com/api/v1/ready
```

### Alerts
Set up alerts for:
- [ ] Server down (health check fails)
- [ ] High error rate (> 5%)
- [ ] High memory usage (> 80%)
- [ ] High CPU usage (> 80%)
- [ ] Database connection failures
- [ ] Slow response times (> 1s)

### Logs
```bash
# Tail logs
tail -f /var/log/trading-bot/server.log

# Search for errors
grep "ERROR\|FATAL" /var/log/trading-bot/server.log
```

## Environment-Specific

### Development
- [ ] Debug logging enabled
- [ ] CORS allows localhost
- [ ] Test database used
- [ ] Rate limits relaxed

### Staging
- [ ] Production-like configuration
- [ ] Separate database
- [ ] Monitoring enabled
- [ ] Load testing performed

### Production
- [ ] All security measures enabled
- [ ] Monitoring and alerting active
- [ ] Backup strategy implemented
- [ ] Disaster recovery plan documented
- [ ] On-call rotation established

## Documentation

- [ ] API documentation updated
- [ ] Architecture diagrams current
- [ ] Runbook created
- [ ] Troubleshooting guide available
- [ ] Team trained on new features

## Compliance

- [ ] Data privacy requirements met
- [ ] Logging compliant with regulations
- [ ] Security audit completed
- [ ] Terms of service updated
- [ ] Privacy policy updated

## Performance Benchmarks

Expected metrics:
- API response time: < 100ms (p95)
- Backtest completion: < 5s for 30 days
- Memory usage: < 500MB
- CPU usage: < 50% average
- Database connections: < 25
- WebSocket connections: < 1000

## Success Criteria

Deployment is successful when:
- [ ] All health checks passing
- [ ] Zero critical errors in logs
- [ ] Performance within benchmarks
- [ ] All features working as expected
- [ ] User feedback positive
- [ ] Monitoring shows green status

## Emergency Contacts

- **On-call Engineer**: [Contact]
- **Database Admin**: [Contact]
- **DevOps Lead**: [Contact]
- **Product Owner**: [Contact]

## Notes

Add deployment-specific notes here:
- Date:
- Version:
- Changes:
- Known issues:
- Special considerations:

---

**Last Updated**: December 2024  
**Version**: 1.0.0
