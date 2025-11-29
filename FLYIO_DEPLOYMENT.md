# Deploy to Fly.io - 100% FREE

## ✅ Why Fly.io?

- 🆓 **Completely FREE** - No credit card required
- ⚡ **Fast** - Global edge network
- 🌏 **Singapore region** - Low latency for India
- 🔄 **Auto-deploy** - From GitHub
- 💪 **Always on** - No sleep mode
- 📊 **256MB RAM** - Perfect for Go backend

---

## 🚀 Step-by-Step Deployment

### Step 1: Install Fly CLI (2 minutes)

**macOS/Linux:**
```bash
curl -L https://fly.io/install.sh | sh
```

**Windows (PowerShell):**
```powershell
iwr https://fly.io/install.ps1 -useb | iex
```

**Verify installation:**
```bash
flyctl version
```

### Step 2: Sign Up & Login (1 minute)

```bash
flyctl auth signup
# Or if you have an account:
flyctl auth login
```

- Sign up with GitHub (easiest)
- No credit card required!

### Step 3: Navigate to Backend

```bash
cd backend
```

### Step 4: Launch App (1 minute)

```bash
flyctl launch
```

**Answer the prompts:**
- App name: `tradebot-api` (or your choice)
- Region: Choose `sin` (Singapore) - closest to India
- PostgreSQL: **No** (we're using Supabase)
- Redis: **No**
- Deploy now: **No** (we need to set secrets first)

### Step 5: Set Environment Variables

```bash
# Set your Supabase credentials
flyctl secrets set SUPABASE_HOST=xlxugbqxfrrwutxecwug.supabase.co
flyctl secrets set SUPABASE_PASSWORD=your_database_password_here
```

**Get your Supabase password:**
1. Go to: https://supabase.com/dashboard/project/xlxugbqxfrrwutxecwug
2. Settings → Database → Database Password
3. Copy and use it in the command above

### Step 6: Deploy! 🚀

```bash
flyctl deploy
```

Wait 2-3 minutes for deployment...

### Step 7: Get Your URL

```bash
flyctl info
```

Your API will be at: `https://tradebot-api.fly.dev`

### Step 8: Test Your API

```bash
curl https://tradebot-api.fly.dev/api/v1/health
```

Should return:
```json
{
  "status": "ok",
  "message": "Trading Bot API is running"
}
```

---

## 🎯 Your API Endpoints

Once deployed, your API will be available at:

```
https://tradebot-api.fly.dev/api/v1/
```

**Endpoints:**
- `GET /api/v1/health` - Health check
- `POST /api/v1/signals` - Create signal
- `GET /api/v1/signals` - Get all signals
- `GET /api/v1/signals/pending` - Get pending signals
- `GET /api/v1/analytics/performance` - Performance stats
- `GET /api/v1/analytics/by-killzone` - Kill zone stats
- `GET /api/v1/analytics/by-pattern` - Pattern stats

---

## 🔄 Auto-Deploy from GitHub

Set up automatic deployment:

```bash
# Connect to GitHub
flyctl apps create tradebot-api
flyctl secrets set GITHUB_TOKEN=your_github_token

# Now every push to main branch auto-deploys!
```

---

## 📊 Monitor Your App

```bash
# View logs
flyctl logs

# Check status
flyctl status

# View metrics
flyctl dashboard
```

---

## 💰 Cost Breakdown

**FREE Tier Includes:**
- ✅ 3 shared-cpu VMs (you only need 1)
- ✅ 256MB RAM per VM
- ✅ 3GB storage
- ✅ 160GB outbound data transfer
- ✅ Unlimited inbound data

**Your Usage:**
- Go backend: ~20MB RAM ✅
- API calls: Fast & efficient ✅
- Storage: Minimal ✅
- **Result: 100% FREE forever!** 🎉

---

## 🛠️ Useful Commands

```bash
# View logs in real-time
flyctl logs -a tradebot-api

# Restart app
flyctl apps restart tradebot-api

# Scale (if needed)
flyctl scale count 1 -a tradebot-api

# SSH into your app
flyctl ssh console -a tradebot-api

# Update secrets
flyctl secrets set KEY=VALUE -a tradebot-api

# Deploy new version
flyctl deploy
```

---

## 🐛 Troubleshooting

**Error: "failed to fetch an image"**
```bash
flyctl deploy --local-only
```

**Error: "could not find App"**
```bash
flyctl apps list
# Use the correct app name
```

**Error: "health check failed"**
- Check logs: `flyctl logs`
- Verify database connection
- Check environment variables: `flyctl secrets list`

**App is slow to start:**
- First request after idle takes ~2s (cold start)
- Subsequent requests are instant

---

## 🎉 Next Steps

After deployment:

1. ✅ Test all API endpoints
2. ✅ Update frontend to use new API URL
3. ✅ Monitor logs for any issues
4. ✅ Enjoy your free, fast backend!

---

## 📝 Update Frontend

In your `supabase-config.js` or create new `api-config.js`:

```javascript
const API_URL = 'https://tradebot-api.fly.dev/api/v1';

// Replace Supabase calls with API calls
async function saveSignal(signal) {
    const response = await fetch(`${API_URL}/signals`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(signal)
    });
    return response.json();
}

async function getAllSignals() {
    const response = await fetch(`${API_URL}/signals`);
    return response.json();
}

async function getAnalytics() {
    const response = await fetch(`${API_URL}/analytics/performance`);
    return response.json();
}
```

---

## 🌟 Benefits

- ⚡ **10x faster** than Supabase direct calls
- 🆓 **100% free** - No hidden costs
- 🌏 **Low latency** - Singapore region
- 🔒 **Secure** - HTTPS by default
- 📈 **Scalable** - Can handle thousands of requests
- 🔄 **Always on** - No sleep mode

---

## 🎯 Summary

1. Install Fly CLI
2. Run `flyctl launch` in backend folder
3. Set secrets (Supabase credentials)
4. Deploy with `flyctl deploy`
5. Get your URL and test!

**Total time: 5 minutes**
**Total cost: $0 forever**

Enjoy your ultra-fast, free Go backend! 🚀
