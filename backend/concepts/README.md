# 📚 Trading Concepts

This folder contains the core trading concepts and analysis tools used across all strategies.

## ICT/Smart Money Concepts

### Core ICT Files
- **ict_smc.go** - Smart Money Concepts (SMC) implementation
- **ict_entry_models.go** - ICT entry models and setups
- **institutional_setups.go** - Institutional trading setups

### Liquidity Analysis
- **liquidity_sweep.go** - Liquidity sweep detection
- **session_liquidity.go** - Session-based liquidity analysis

### Market Structure
- **orderflow_analysis.go** - Order flow and delta analysis
- **delta_pivot_analysis.go** - Delta pivot points
- **supply_demand.go** - Supply and demand zones
- **market_maker_model.go** - Market maker behavior patterns

### Advanced Concepts
- **power_of_3.go** - Power of 3 (Accumulation, Manipulation, Distribution)
- **mirror_market.go** - Mirror market analysis
- **multi_timeframe_confluence.go** - Multi-timeframe analysis

### Pattern Recognition
- **candlestick_patterns.go** - Candlestick pattern detection
- **advanced_patterns.go** - Advanced chart patterns

### Filters & Tools
- **volatility_filter.go** - Volatility-based filtering

## Purpose

These files contain **reusable trading concepts** that can be:
- Used by multiple strategies
- Combined in different ways
- Extended with new features
- Tested independently

## Usage

Strategies import and use these concepts:
```go
import (
    "backend/concepts"
)

// Use ICT concepts
setup := concepts.DetectInstitutionalSetup(candles)
liquidity := concepts.AnalyzeLiquiditySweep(data)
```

## File Organization

### Analysis Tools
- Order flow analysis
- Delta analysis
- Liquidity detection
- Market structure

### Pattern Detection
- Candlestick patterns
- Chart patterns
- ICT setups

### Market Models
- Smart Money Concepts
- Market maker behavior
- Institutional trading

### Filters
- Volatility filters
- Time-based filters
- Confluence checks

## Key Concepts Explained

### 1. Smart Money Concepts (SMC)
Understanding how institutions trade and following their footprints.

### 2. Liquidity Sweeps
Detecting when price hunts stop losses before reversing.

### 3. Order Flow
Analyzing buying/selling pressure through delta and volume.

### 4. Supply & Demand
Identifying institutional supply and demand zones.

### 5. Power of 3
The three phases: Accumulation → Manipulation → Distribution

### 6. Multi-Timeframe Confluence
Aligning signals across multiple timeframes for higher probability.

## Adding New Concepts

To add a new concept:
1. Create a new `.go` file in this folder
2. Implement the concept with clear functions
3. Add tests
4. Document in this README
5. Use in strategies

## Dependencies

Concepts should be:
- ✅ Independent (minimal dependencies)
- ✅ Reusable (used by multiple strategies)
- ✅ Testable (unit tests)
- ✅ Well-documented (clear comments)

## Related Documentation

- `ADVANCED_TRADING_CONCEPTS.md` - Detailed concept explanations
- `LIQUIDITY_FIRST_STRATEGY.md` - Liquidity-focused approach
- Strategy folder - How concepts are used
