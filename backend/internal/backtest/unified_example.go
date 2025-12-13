package backtest

// Example usage of the Unified Backtest Engine
// This file demonstrates various use cases

// Example 1: Basic Backtest
func ExampleBasicBacktest() (*UnifiedBacktestResult, error) {
	config := UnifiedBacktestConfig{
		Symbol:       "BTCUSDT",
		Interval:     "15m",
		Days:         30,
		StartBalance: 10000,
		Strategy:     "liquidity_hunter",
	}
	
	candles, err := FetchBinanceData(config.Symbol, config.Interval, config.Days)
	if err != nil {
		return nil, err
	}
	
	return RunUnifiedBacktest(config, candles)
}

// Example 2: Professional Backtest with Partial Exits
func ExampleProfessionalBacktest() (*UnifiedBacktestResult, error) {
	config := UnifiedBacktestConfig{
		Symbol:             "BTCUSDT",
		Interval:           "15m",
		Days:               30,
		StartBalance:       10000,
		Strategy:           "breakout_master",
		EnablePartialExits: true,
		RiskPercent:        0.003, // 0.3% risk per trade
		MaxTradesPerDay:    20,
	}
	
	candles, err := FetchBinanceData(config.Symbol, config.Interval, config.Days)
	if err != nil {
		return nil, err
	}
	
	return RunUnifiedBacktest(config, candles)
}

// Example 3: Advanced Backtest with Monte Carlo and Stress Testing
func ExampleAdvancedBacktest() (*UnifiedBacktestResult, error) {
	config := UnifiedBacktestConfig{
		Symbol:       "BTCUSDT",
		Interval:     "15m",
		Days:         90,
		StartBalance: 10000,
		Strategy:     "session_trader",
		
		// Risk Management
		RiskPercent:        0.003,
		MaxDailyLoss:       5.0,
		MaxConsecutiveLoss: 3,
		MaxTradesPerDay:    20,
		
		// Market Filters
		TradingHoursOnly: true,
		MinVolatility:    0.5,
		MaxVolatility:    3.0,
		
		// Advanced Analysis
		EnableMonteCarlo: true,
		MonteCarloRuns:   1000,
		EnableStressTest: true,
		
		// Simulation
		WindowType:        "expanding",
		RealisticSlippage: true,
	}
	
	candles, err := FetchBinanceData(config.Symbol, config.Interval, config.Days)
	if err != nil {
		return nil, err
	}
	
	return RunUnifiedBacktest(config, candles)
}

// Example 4: Walk-Forward Analysis
func ExampleWalkForwardBacktest() (*UnifiedBacktestResult, error) {
	config := UnifiedBacktestConfig{
		Symbol:         "BTCUSDT",
		Interval:       "15m",
		Days:           180,
		StartBalance:   10000,
		Strategy:       "trend_rider",
		UseWalkForward: true,
		TrainingDays:   60,
		TestingDays:    30,
	}
	
	candles, err := FetchBinanceData(config.Symbol, config.Interval, config.Days)
	if err != nil {
		return nil, err
	}
	
	return RunUnifiedBacktest(config, candles)
}

// Example 5: Parallel Strategy Testing
func ExampleParallelStrategies() (*UnifiedBacktestResult, error) {
	config := UnifiedBacktestConfig{
		Symbol:         "BTCUSDT",
		Interval:       "15m",
		Days:           30,
		StartBalance:   10000,
		EnableParallel: true,
		Strategies: []string{
			"liquidity_hunter",
			"breakout_master",
			"session_trader",
			"trend_rider",
			"smart_money_tracker",
		},
	}
	
	candles, err := FetchBinanceData(config.Symbol, config.Interval, config.Days)
	if err != nil {
		return nil, err
	}
	
	// Returns the best performing strategy
	return RunUnifiedBacktest(config, candles)
}

// Example 6: Ultra-Conservative Backtest
func ExampleConservativeBacktest() (*UnifiedBacktestResult, error) {
	config := UnifiedBacktestConfig{
		Symbol:       "BTCUSDT",
		Interval:     "1h",
		Days:         60,
		StartBalance: 10000,
		Strategy:     "liquidity_hunter",
		
		// Very conservative risk management
		RiskPercent:        0.001, // 0.1% per trade
		MaxDailyLoss:       2.0,   // Stop at 2% daily loss
		MaxConsecutiveLoss: 2,     // Stop after 2 losses
		MaxTradesPerDay:    5,     // Max 5 trades per day
		
		// Strict filters
		TradingHoursOnly: true,
		MinVolatility:    0.8,
		MaxVolatility:    2.0,
		
		// Validation
		EnableMonteCarlo: true,
		UseWalkForward:   true,
	}
	
	candles, err := FetchBinanceData(config.Symbol, config.Interval, config.Days)
	if err != nil {
		return nil, err
	}
	
	return RunUnifiedBacktest(config, candles)
}

// Example 7: Aggressive Scalping Backtest
func ExampleAggressiveScalping() (*UnifiedBacktestResult, error) {
	config := UnifiedBacktestConfig{
		Symbol:       "BTCUSDT",
		Interval:     "5m",
		Days:         14,
		StartBalance: 10000,
		Strategy:     "scalper_pro",
		
		// Aggressive settings
		RiskPercent:     0.005, // 0.5% per trade
		MaxTradesPerDay: 50,    // High frequency
		
		// Fast execution
		WindowType:        "rolling",
		MaxWindow:         50,
		EnablePartialExits: true,
	}
	
	candles, err := FetchBinanceData(config.Symbol, config.Interval, config.Days)
	if err != nil {
		return nil, err
	}
	
	return RunUnifiedBacktest(config, candles)
}

// Example 8: Complete Production-Ready Backtest
func ExampleProductionBacktest() (*UnifiedBacktestResult, error) {
	config := UnifiedBacktestConfig{
		// Basic Configuration
		Symbol:       "BTCUSDT",
		Interval:     "15m",
		Days:         90,
		StartBalance: 10000,
		Strategy:     "liquidity_hunter",
		
		// Risk Management (Production-grade)
		RiskPercent:         0.003, // 0.3% per trade
		MaxPositionCap:      100000,
		MaxDailyLoss:        5.0,
		MaxWeeklyLoss:       10.0,
		MaxConsecutiveLoss:  3,
		DynamicPositionSize: false, // Fixed for consistency
		MaxTradesPerDay:     20,
		
		// Trading Costs (Realistic)
		SlippagePercent:   0.0015, // 0.15%
		FeePercent:        0.001,  // 0.1%
		RealisticSlippage: true,
		IncludeSpread:     true,
		SpreadPercent:     0.0005, // 0.05%
		
		// Market Filters
		MinVolatility:    0.5,
		MaxVolatility:    3.0,
		MinVolume:        1.0,
		TradingHoursOnly: true,
		
		// Simulation (Most realistic)
		WindowType:        "expanding",
		MinWindow:         100,
		MaxWindow:         200,
		UseWalkForward:    true,
		TrainingDays:      60,
		TestingDays:       30,
		
		// Validation (Complete)
		EnableMonteCarlo:   true,
		MonteCarloRuns:     1000,
		EnableStressTest:   true,
		EnablePartialExits: true,
	}
	
	candles, err := FetchBinanceData(config.Symbol, config.Interval, config.Days)
	if err != nil {
		return nil, err
	}
	
	return RunUnifiedBacktest(config, candles)
}
