# Backend Structure Guide

## 📁 New Structure

```
backend/
├── cmd/server/main.go           # Entry point
├── internal/
│   ├── api/                     # HTTP layer (handlers, routes, middleware)
│   ├── backtest/                # Backtest engines
│   ├── database/                # Data persistence
│   ├── signals/                 # Signal generation
│   ├── strategies/              # Trading strategies
│   │   ├── ict/                 # ICT/SMC strategies
│   │   ├── patterns/            # Pattern recognition
│   │   ├── institutional/       # Institutional trading
│   │   ├── timeframe/           # Multi-timeframe
│   │   └── daily/               # Daily strategies
│   ├── filters/                 # Trade validation
│   ├── optimization/            # Parameter tuning
│   ├── trading/                 # Paper trading
│   ├── activity/                # Activity logging
│   ├── ai/                      # AI integration
│   ├── communication/           # WebSocket & Telegram
│   └── templates/               # HTML rendering
├── tests/unit/                  # Unit tests
├── deployments/                 # Docker, Fly.io, Render
└── scripts/                     # Utility scripts
```

## ⚡ Quick Fix (30 min)

### 1. Export Functions (Capitalize)

**internal/database/connection.go:**
```go
var DB *sql.DB  // Exported
func InitDB() { ... }
func RunMigrations() { ... }
```

**internal/activity/logger.go:**
```go
func GetActivityLogger() *ActivityLogger { ... }
func LogSystemSuccess(msg string, data fiber.Map) { ... }
```

**internal/communication/websocket/server.go:**
```go
var Hub *hub  // Exported
func (h *hub) Run() { ... }
```

**internal/communication/telegram/telegram_bot.go:**
```go
func InitTelegramBot() { ... }
func StartTelegramSignalBot(...) error { ... }
```

**internal/signals/ai_enhanced_signal_generator.go:**
```go
func NewAIEnhancedSignalGenerator() *AIEnhancedSignalGenerator { ... }
func (g *AIEnhancedSignalGenerator) Start() { ... }
```

### 2. Update main.go

Add imports:
```go
import (
    "tradebot-backend/internal/api"
    "tradebot-backend/internal/database"
    "tradebot-backend/internal/activity"
    "tradebot-backend/internal/signals"
    "tradebot-backend/internal/communication/websocket"
    "tradebot-backend/internal/communication/telegram"
)
```

Update calls:
```go
database.InitDB()
database.RunMigrations()
activity.GetActivityLogger()
websocket.Hub.Run()
signals.StartSignalBroadcaster()
telegram.InitTelegramBot()
api.SetupRoutes(app)
```

### 3. Update routes.go

```go
package api

import "tradebot-backend/internal/api/handlers"

func SetupRoutes(app *fiber.App) {
    api.Get("/health", handlers.HealthHandler)
    // ... etc
}
```

### 4. Build

```bash
cd backend
go mod tidy
go build ./cmd/server
go run ./cmd/server
```

## 📊 File Distribution

- api/handlers: 21 files
- backtest: 8 files
- strategies/ict: 7 files
- signals: 6 files
- strategies: 5 files
- database: 4 files
- Other: 35 files

**Total: 86 organized files**

## ✅ Benefits

- Clear domain separation
- Easy to navigate
- Scalable architecture
- Production-ready
- Follows Go standards
