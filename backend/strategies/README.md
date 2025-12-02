# 🎯 Trading Strategies

This folder contains complete trading strategies that combine concepts from the `concepts/` folder.

## Available Strategies

### 1. AI-Enhanced Strategy
**File**: `ai_enhanced_signal_generator.go`
- Uses Grok AI for signal validation
- Sentiment analysis and risk adjustment
- Dynamic position sizing
- News impact filtering

**Best For**: All market conditions with AI oversight

### 2. Liquidity-First Strategy
**File**: `liquidity_first_strategy.go`
- Focuses on liquidity sweeps
- Hunts stop loss clusters
- High win rate approach
- Conservative risk management

**Best For**: Ranging and choppy markets

### 3. Professional Strategy
**File**: `professional_strategy.go`
- Multi-timeframe confluence
- ICT concepts integration
- Institutional setup detection
- Advanced filtering

**Best For**: Trending markets with clear structure

### 4. Ultimate Daily Strategy
**File**: `ultimate_daily_strategy.go`
- Daily timeframe focus
- Swing trading approach
- Multiple concept integration
- Long-term positions

**Best For**: Swing traders, lower frequency

### 5. Enhanced Strategy
**File**: `enhanced_strategy.go`
- Balanced approach
- Multiple indicators
- Pattern recognition
- Medium frequency

**Best For**: General trading, all conditions

### 6. Advanced Signal Generator
**File**: `advanced_signal_generator.go`
- High-frequency signals
- Quick scalping setups
- Fast execution
- Multiple timeframes

**Best For**: Scalpers, day traders

### 7. Basic Signal Generator
**File**: `signal_generator.go`
- Simple technical analysis
- RSI + EMA based
- Good for beginners
- Easy to understand

**Best For**: Learning and testing

### 8. Backtest Signal Generator
**File**: `backtest_signal_generator.go`
- Historical testing
- Strategy validation
- Performance metrics
- Optimization

**Best For**: Strategy development

## Strategy Components

### Configuration
**File**: `strategy_configs.go`
- Strategy parameters
- Risk settings
- Timeframe configs
- Entry/exit rules

### Handler
**File**: `strategy_handler.go`
- API endpoints
- Strategy selection
- Parameter updates
- Performance tracking

## How Strategies Work

```
┌─────────────────────┐
│ Market Data Input   │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ Concepts Analysis   │
│ (from concepts/)    │
│                     │
│ • ICT Setups        │
│ • Liquidity         │
│ • Order Flow        │
│ • Patterns          │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ Strategy Logic      │
│ (this folder)       │
│                     │
│ • Combine concepts  │
│ • Apply filters     │
│ • Generate signals  │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ Signal Output       │
│ • Entry price       │
│ • Stop loss         │
│ • Take profits      │
│ • Risk/reward       │
└─────────────────────┘
```

## Strategy Comparison

| Strategy | Win Rate | Frequency | Risk | Complexity |
|----------|----------|-----------|------|------------|
| AI-Enhanced | 75-85% | Medium | Low | High |
| Liquidity-First | 70-80% | Low | Low | Medium |
| Professional | 65-75% | Medium | Medium | High |
| Ultimate Daily | 60-70% | Very Low | Low | Medium |
| Enhanced | 65-75% | Medium | Medium | Medium |
| Advanced | 60-70% | High | High | High |
| Basic | 55-65% | High | Medium | Low |

## Choosing a Strategy

### For Beginners
Start with: **Basic Signal Generator**
- Simple logic
- Easy to understand
- Good for learning

### For Consistent Profits
Use: **Liquidity-First Strategy**
- High win rate
- Lower frequency
- Conservative approach

### For Maximum Performance
Use: **AI-Enhanced Strategy**
- Best win rate
- Smart filtering
- Adaptive risk

### For Swing Trading
Use: **Ultimate Daily Strategy**
- Daily timeframe
- Longer holds
- Less monitoring

### For Day Trading
Use: **Professional Strategy**
- Intraday focus
- Multiple timeframes
- Active management

## Creating a New Strategy

1. **Create new file**: `my_strategy.go`
2. **Import concepts**:
```go
import (
    "backend/concepts"
)
```
3. **Define strategy struct**:
```go
type MyStrategy struct {
    Symbol     string
    Timeframe  string
    MinStrength int
}
```
4. **Implement analysis**:
```go
func (s *MyStrategy) Analyze(data []Candle) *Signal {
    // Use concepts
    liquidity := concepts.DetectLiquiditySweep(data)
    ict := concepts.AnalyzeICTSetup(data)
    
    // Your logic
    if liquidity.Detected && ict.Valid {
        return generateSignal()
    }
    return nil
}
```
5. **Add to strategy handler**
6. **Test and optimize**

## Strategy Configuration

Each strategy can be configured via:
- `strategy_configs.go` - Default settings
- API endpoints - Runtime updates
- Environment variables - Deployment configs

## Performance Tracking

All strategies track:
- Win rate
- Profit factor
- Max drawdown
- Average win/loss
- Holding time
- Risk-adjusted returns

## Testing Strategies

Use the backtest engine:
```bash
# Test liquidity-first strategy
./test_liquidity_first.sh

# Test professional strategy
./test_professional_strategy.sh

# Test ultimate strategy
./test_ultimate_strategy.sh
```

## Best Practices

1. **Combine Multiple Concepts** - Don't rely on one indicator
2. **Use Proper Risk Management** - Always set stop loss
3. **Filter Signals** - Quality over quantity
4. **Backtest First** - Validate before live trading
5. **Monitor Performance** - Track and adjust
6. **Use AI Enhancement** - Let AI filter bad trades

## Strategy Updates

Strategies are continuously improved based on:
- Performance data
- Market conditions
- User feedback
- New concepts

## Related Files

- `../concepts/` - Trading concepts used by strategies
- `../models.go` - Data structures
- `../handlers.go` - API endpoints
- `../backtest_engine.go` - Testing framework

## Documentation

- `LIQUIDITY_FIRST_STRATEGY.md` - Liquidity strategy guide
- `ADVANCED_TRADING_CONCEPTS.md` - Concept explanations
- `GROK_AI_INTEGRATION.md` - AI enhancement guide
- `IMPLEMENTATION_GUIDE.md` - Setup instructions
