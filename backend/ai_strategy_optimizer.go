package main

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
)

// AIStrategyOptimizer uses AI to optimize trading strategies
type AIStrategyOptimizer struct {
	LearningRate    float64
	Generations     int
	PopulationSize  int
	MutationRate    float64
}

// StrategyGene represents a strategy's parameters
type StrategyGene struct {
	StopLossATR     float64
	TakeProfitATR   float64
	ADXThreshold    float64
	RSILow          float64
	RSIHigh         float64
	CooldownCandles int
	MinConfluence   int
	Fitness         float64
}

// AIOptimizationResult holds AI optimization results
type AIOptimizationResult struct {
	BestStrategy      StrategyGene            `json:"bestStrategy"`
	TopStrategies     []StrategyGene          `json:"topStrategies"`
	GenerationHistory []GenerationStats       `json:"generationHistory"`
	ImprovementPct    float64                 `json:"improvementPct"`
	TotalTests        int                     `json:"totalTests"`
	Duration          string                  `json:"duration"`
	Recommendation    string                  `json:"recommendation"`
}

// GenerationStats tracks evolution progress
type GenerationStats struct {
	Generation  int     `json:"generation"`
	BestFitness float64 `json:"bestFitness"`
	AvgFitness  float64 `json:"avgFitness"`
	Diversity   float64 `json:"diversity"`
}

// AIMarketAnalysis provides AI-powered market insights
type AIMarketAnalysis struct {
	MarketRegime      string             `json:"marketRegime"`      // "trending", "ranging", "volatile"
	TrendStrength     float64            `json:"trendStrength"`     // 0-100
	VolatilityLevel   string             `json:"volatilityLevel"`   // "low", "medium", "high"
	SupportLevels     []float64          `json:"supportLevels"`
	ResistanceLevels  []float64          `json:"resistanceLevels"`
	PredictedMove     string             `json:"predictedMove"`     // "up", "down", "sideways"
	Confidence        float64            `json:"confidence"`        // 0-100
	BestStrategy      string             `json:"bestStrategy"`
	RiskLevel         string             `json:"riskLevel"`         // "low", "medium", "high"
	Recommendations   []string           `json:"recommendations"`
}

// RunAIOptimization uses genetic algorithm to find optimal parameters
func RunAIOptimization(config BacktestConfig, candles []Candle) (*AIOptimizationResult, error) {
	optimizer := &AIStrategyOptimizer{
		LearningRate:   0.1,
		Generations:    20,
		PopulationSize: 50,
		MutationRate:   0.2,
	}
	
	fmt.Println("🤖 Starting AI-powered optimization...")
	fmt.Printf("   Generations: %d, Population: %d\n", optimizer.Generations, optimizer.PopulationSize)
	
	// Initialize population with random strategies
	population := optimizer.initializePopulation()
	
	generationHistory := []GenerationStats{}
	initialFitness := 0.0
	
	// Evolution loop
	for gen := 0; gen < optimizer.Generations; gen++ {
		// Evaluate fitness for each strategy
		for i := range population {
			fitness := optimizer.evaluateFitness(&population[i], config, candles)
			population[i].Fitness = fitness
		}
		
		// Sort by fitness
		sort.Slice(population, func(i, j int) bool {
			return population[i].Fitness > population[j].Fitness
		})
		
		// Track progress
		avgFitness := 0.0
		for _, gene := range population {
			avgFitness += gene.Fitness
		}
		avgFitness /= float64(len(population))
		
		if gen == 0 {
			initialFitness = population[0].Fitness
		}
		
		stats := GenerationStats{
			Generation:  gen + 1,
			BestFitness: population[0].Fitness,
			AvgFitness:  avgFitness,
			Diversity:   optimizer.calculateDiversity(population),
		}
		generationHistory = append(generationHistory, stats)
		
		fmt.Printf("   Gen %d: Best=%.2f, Avg=%.2f, Diversity=%.2f\n", 
			gen+1, stats.BestFitness, stats.AvgFitness, stats.Diversity)
		
		// Create next generation
		if gen < optimizer.Generations-1 {
			population = optimizer.evolve(population)
		}
	}
	
	// Get top 5 strategies
	topStrategies := population[:5]
	
	improvementPct := 0.0
	if initialFitness > 0 {
		improvementPct = ((population[0].Fitness - initialFitness) / initialFitness) * 100
	}
	
	result := &AIOptimizationResult{
		BestStrategy:      population[0],
		TopStrategies:     topStrategies,
		GenerationHistory: generationHistory,
		ImprovementPct:    improvementPct,
		TotalTests:        optimizer.Generations * optimizer.PopulationSize,
		Recommendation:    optimizer.getRecommendation(population[0]),
	}
	
	fmt.Printf("✅ AI optimization complete! Improvement: %.2f%%\n", improvementPct)
	
	return result, nil
}

// initializePopulation creates random strategies
func (ai *AIStrategyOptimizer) initializePopulation() []StrategyGene {
	population := make([]StrategyGene, ai.PopulationSize)
	
	for i := range population {
		population[i] = StrategyGene{
			StopLossATR:     randomFloat(0.5, 2.0),
			TakeProfitATR:   randomFloat(2.0, 8.0),
			ADXThreshold:    randomFloat(20.0, 35.0),
			RSILow:          randomFloat(30.0, 45.0),
			RSIHigh:         randomFloat(55.0, 70.0),
			CooldownCandles: randomInt(20, 50),
			MinConfluence:   randomInt(5, 9),
		}
	}
	
	return population
}

// evaluateFitness tests a strategy and returns fitness score
func (ai *AIStrategyOptimizer) evaluateFitness(gene *StrategyGene, config BacktestConfig, candles []Candle) float64 {
	// Run backtest with these parameters
	result, err := RunBacktest(config, candles)
	if err != nil || result.TotalTrades < 10 {
		return 0
	}
	
	// Multi-objective fitness function
	// Balances: profit factor, win rate, drawdown, trade count
	
	profitScore := math.Min(result.ProfitFactor*20, 100)
	winRateScore := result.WinRate
	drawdownScore := math.Max(0, 100-(result.MaxDrawdown*5))
	tradeScore := math.Min(float64(result.TotalTrades)/2, 50)
	
	// Weighted fitness
	fitness := (profitScore * 0.4) + (winRateScore * 0.3) + (drawdownScore * 0.2) + (tradeScore * 0.1)
	
	return fitness
}

// evolve creates next generation using selection, crossover, mutation
func (ai *AIStrategyOptimizer) evolve(population []StrategyGene) []StrategyGene {
	nextGen := make([]StrategyGene, ai.PopulationSize)
	
	// Elitism: keep top 10%
	eliteCount := ai.PopulationSize / 10
	copy(nextGen[:eliteCount], population[:eliteCount])
	
	// Create rest through crossover and mutation
	for i := eliteCount; i < ai.PopulationSize; i++ {
		parent1 := ai.selectParent(population)
		parent2 := ai.selectParent(population)
		
		child := ai.crossover(parent1, parent2)
		child = ai.mutate(child)
		
		nextGen[i] = child
	}
	
	return nextGen
}

// selectParent uses tournament selection
func (ai *AIStrategyOptimizer) selectParent(population []StrategyGene) StrategyGene {
	tournamentSize := 5
	best := population[randomInt(0, len(population))]
	
	for i := 1; i < tournamentSize; i++ {
		candidate := population[randomInt(0, len(population))]
		if candidate.Fitness > best.Fitness {
			best = candidate
		}
	}
	
	return best
}

// crossover combines two parents
func (ai *AIStrategyOptimizer) crossover(parent1, parent2 StrategyGene) StrategyGene {
	return StrategyGene{
		StopLossATR:     (parent1.StopLossATR + parent2.StopLossATR) / 2,
		TakeProfitATR:   (parent1.TakeProfitATR + parent2.TakeProfitATR) / 2,
		ADXThreshold:    (parent1.ADXThreshold + parent2.ADXThreshold) / 2,
		RSILow:          (parent1.RSILow + parent2.RSILow) / 2,
		RSIHigh:         (parent1.RSIHigh + parent2.RSIHigh) / 2,
		CooldownCandles: (parent1.CooldownCandles + parent2.CooldownCandles) / 2,
		MinConfluence:   (parent1.MinConfluence + parent2.MinConfluence) / 2,
	}
}

// mutate randomly changes parameters
func (ai *AIStrategyOptimizer) mutate(gene StrategyGene) StrategyGene {
	if randomFloat(0, 1) < ai.MutationRate {
		gene.StopLossATR += randomFloat(-0.2, 0.2)
		gene.StopLossATR = math.Max(0.5, math.Min(2.0, gene.StopLossATR))
	}
	if randomFloat(0, 1) < ai.MutationRate {
		gene.TakeProfitATR += randomFloat(-1.0, 1.0)
		gene.TakeProfitATR = math.Max(2.0, math.Min(8.0, gene.TakeProfitATR))
	}
	if randomFloat(0, 1) < ai.MutationRate {
		gene.ADXThreshold += randomFloat(-3.0, 3.0)
		gene.ADXThreshold = math.Max(20.0, math.Min(35.0, gene.ADXThreshold))
	}
	
	return gene
}

// calculateDiversity measures population diversity
func (ai *AIStrategyOptimizer) calculateDiversity(population []StrategyGene) float64 {
	if len(population) < 2 {
		return 0
	}
	
	totalDistance := 0.0
	comparisons := 0
	
	for i := 0; i < len(population)-1; i++ {
		for j := i + 1; j < len(population); j++ {
			distance := ai.geneticDistance(population[i], population[j])
			totalDistance += distance
			comparisons++
		}
	}
	
	return totalDistance / float64(comparisons)
}

// geneticDistance calculates difference between two strategies
func (ai *AIStrategyOptimizer) geneticDistance(g1, g2 StrategyGene) float64 {
	d1 := math.Abs(g1.StopLossATR - g2.StopLossATR)
	d2 := math.Abs(g1.TakeProfitATR - g2.TakeProfitATR)
	d3 := math.Abs(g1.ADXThreshold - g2.ADXThreshold)
	d4 := math.Abs(g1.RSILow - g2.RSILow)
	d5 := math.Abs(g1.RSIHigh - g2.RSIHigh)
	
	return (d1 + d2 + d3 + d4 + d5) / 5
}

// getRecommendation provides AI recommendation
func (ai *AIStrategyOptimizer) getRecommendation(gene StrategyGene) string {
	if gene.Fitness > 80 {
		return "⭐⭐⭐⭐⭐ Excellent - AI found optimal parameters"
	} else if gene.Fitness > 70 {
		return "⭐⭐⭐⭐ Good - Strong parameters found"
	} else if gene.Fitness > 60 {
		return "⭐⭐⭐ Fair - Decent parameters, room for improvement"
	}
	return "⭐⭐ Poor - Consider different strategy or more data"
}

// AnalyzeMarketWithAI provides AI-powered market analysis
func AnalyzeMarketWithAI(candles []Candle) (*AIMarketAnalysis, error) {
	if len(candles) < 200 {
		return nil, fmt.Errorf("need at least 200 candles")
	}
	
	analysis := &AIMarketAnalysis{
		SupportLevels:    []float64{},
		ResistanceLevels: []float64{},
		Recommendations:  []string{},
	}
	
	// Detect market regime
	atr := calculateATR(candles, 14)
	avgPrice := (candles[len(candles)-1].High + candles[len(candles)-1].Low) / 2
	volatility := (atr / avgPrice) * 100
	
	if volatility > 2.0 {
		analysis.VolatilityLevel = "high"
		analysis.RiskLevel = "high"
	} else if volatility > 1.0 {
		analysis.VolatilityLevel = "medium"
		analysis.RiskLevel = "medium"
	} else {
		analysis.VolatilityLevel = "low"
		analysis.RiskLevel = "low"
	}
	
	// Detect trend
	ema20 := calculateEMA(candles, 20)
	ema50 := calculateEMA(candles, 50)
	ema200 := calculateEMA(candles, 200)
	
	currentPrice := candles[len(candles)-1].Close
	
	if ema20 > ema50 && ema50 > ema200 && currentPrice > ema20 {
		analysis.MarketRegime = "trending"
		analysis.PredictedMove = "up"
		analysis.TrendStrength = 80
		analysis.Confidence = 75
		analysis.BestStrategy = "trend_rider"
		analysis.Recommendations = append(analysis.Recommendations, "Strong uptrend - Use trend-following strategies")
	} else if ema20 < ema50 && ema50 < ema200 && currentPrice < ema20 {
		analysis.MarketRegime = "trending"
		analysis.PredictedMove = "down"
		analysis.TrendStrength = 80
		analysis.Confidence = 75
		analysis.BestStrategy = "trend_rider"
		analysis.Recommendations = append(analysis.Recommendations, "Strong downtrend - Use trend-following strategies")
	} else {
		analysis.MarketRegime = "ranging"
		analysis.PredictedMove = "sideways"
		analysis.TrendStrength = 30
		analysis.Confidence = 60
		analysis.BestStrategy = "range_master"
		analysis.Recommendations = append(analysis.Recommendations, "Ranging market - Use mean-reversion strategies")
	}
	
	// Find support/resistance
	analysis.SupportLevels = findSupportLevels(candles, 3)
	analysis.ResistanceLevels = findResistanceLevels(candles, 3)
	
	// Add recommendations based on analysis
	if analysis.VolatilityLevel == "high" {
		analysis.Recommendations = append(analysis.Recommendations, "High volatility - Widen stop losses")
	}
	if analysis.RiskLevel == "high" {
		analysis.Recommendations = append(analysis.Recommendations, "High risk - Reduce position size")
	}
	
	return analysis, nil
}

// Helper functions
func randomFloat(min, max float64) float64 {
	return min + (max-min)*float64(randomInt(0, 1000))/1000.0
}

func randomInt(min, max int) int {
	if max <= min {
		return min
	}
	return min + int(math.Abs(float64(max-min))) % (max - min)
}

func findSupportLevels(candles []Candle, count int) []float64 {
	levels := []float64{}
	for i := 50; i < len(candles)-2; i++ {
		if candles[i].Low < candles[i-1].Low && candles[i].Low < candles[i+1].Low {
			levels = append(levels, candles[i].Low)
		}
	}
	sort.Float64s(levels)
	if len(levels) > count {
		return levels[:count]
	}
	return levels
}

func findResistanceLevels(candles []Candle, count int) []float64 {
	levels := []float64{}
	for i := 50; i < len(candles)-2; i++ {
		if candles[i].High > candles[i-1].High && candles[i].High > candles[i+1].High {
			levels = append(levels, candles[i].High)
		}
	}
	sort.Float64s(levels)
	if len(levels) > count {
		return levels[len(levels)-count:]
	}
	return levels
}

// ToJSON exports AI results
func (r *AIOptimizationResult) ToJSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}
