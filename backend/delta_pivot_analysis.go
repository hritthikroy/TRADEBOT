package main

import (
	"math"
)

// ==================== DELTA ANALYSIS ====================
// Professional delta and cumulative delta analysis

// DeltaAnalysis holds comprehensive delta analysis
type DeltaAnalysis struct {
	CurrentDelta     float64
	CumulativeDelta  float64
	DeltaTrend       string  // "bullish", "bearish", "neutral"
	DeltaDivergence  bool
	DeltaStrength    float64
	BuyVolume        float64
	SellVolume       float64
	DeltaPercentage  float64
	SessionDelta     float64
	DeltaMomentum    string  // "increasing", "decreasing", "stable"
}

// CalculateDelta calculates buy/sell delta from candles
func CalculateDelta(candles []Candle) *DeltaAnalysis {
	if len(candles) < 20 {
		return nil
	}
	
	da := &DeltaAnalysis{}
	
	// Calculate delta for each candle
	// Approximation: bullish candle = more buy volume, bearish = more sell
	totalBuyVol := 0.0
	totalSellVol := 0.0
	cumulativeDelta := 0.0
	
	for _, c := range candles {
		candleRange := c.High - c.Low
		if candleRange == 0 {
			continue
		}
		
		// Calculate buy/sell ratio based on close position
		closePosition := (c.Close - c.Low) / candleRange
		
		buyVol := c.Volume * closePosition
		sellVol := c.Volume * (1 - closePosition)
		
		totalBuyVol += buyVol
		totalSellVol += sellVol
		cumulativeDelta += (buyVol - sellVol)
	}
	
	// Current candle delta
	last := candles[len(candles)-1]
	lastRange := last.High - last.Low
	if lastRange > 0 {
		closePos := (last.Close - last.Low) / lastRange
		da.CurrentDelta = last.Volume * (2*closePos - 1)
	}
	
	da.BuyVolume = totalBuyVol
	da.SellVolume = totalSellVol
	da.CumulativeDelta = cumulativeDelta
	
	// Calculate delta percentage
	totalVol := totalBuyVol + totalSellVol
	if totalVol > 0 {
		da.DeltaPercentage = ((totalBuyVol - totalSellVol) / totalVol) * 100
	}
	
	// Determine delta trend
	recentDelta := 0.0
	for i := len(candles) - 10; i < len(candles); i++ {
		c := candles[i]
		r := c.High - c.Low
		if r > 0 {
			pos := (c.Close - c.Low) / r
			recentDelta += c.Volume * (2*pos - 1)
		}
	}
	
	if recentDelta > 0 && cumulativeDelta > 0 {
		da.DeltaTrend = "bullish"
	} else if recentDelta < 0 && cumulativeDelta < 0 {
		da.DeltaTrend = "bearish"
	} else {
		da.DeltaTrend = "neutral"
	}
	
	// Check for divergence
	priceUp := candles[len(candles)-1].Close > candles[len(candles)-10].Close
	deltaUp := recentDelta > 0
	
	da.DeltaDivergence = (priceUp && !deltaUp) || (!priceUp && deltaUp)
	
	// Calculate strength
	da.DeltaStrength = math.Min(math.Abs(da.DeltaPercentage)*2, 100)
	
	// Delta momentum
	firstHalfDelta := 0.0
	secondHalfDelta := 0.0
	mid := len(candles) / 2
	
	for i := 0; i < mid; i++ {
		c := candles[i]
		r := c.High - c.Low
		if r > 0 {
			pos := (c.Close - c.Low) / r
			firstHalfDelta += c.Volume * (2*pos - 1)
		}
	}
	
	for i := mid; i < len(candles); i++ {
		c := candles[i]
		r := c.High - c.Low
		if r > 0 {
			pos := (c.Close - c.Low) / r
			secondHalfDelta += c.Volume * (2*pos - 1)
		}
	}
	
	if secondHalfDelta > firstHalfDelta*1.2 {
		da.DeltaMomentum = "increasing"
	} else if secondHalfDelta < firstHalfDelta*0.8 {
		da.DeltaMomentum = "decreasing"
	} else {
		da.DeltaMomentum = "stable"
	}
	
	return da
}


// ==================== PIVOT POINTS ====================
// Multiple pivot point calculation methods

// PivotPoints holds all pivot levels
type PivotPoints struct {
	// Standard Pivot
	Pivot float64
	R1    float64
	R2    float64
	R3    float64
	S1    float64
	S2    float64
	S3    float64
	
	// Fibonacci Pivot
	FibR1 float64
	FibR2 float64
	FibR3 float64
	FibS1 float64
	FibS2 float64
	FibS3 float64
	
	// Camarilla Pivot
	CamR1 float64
	CamR2 float64
	CamR3 float64
	CamR4 float64
	CamS1 float64
	CamS2 float64
	CamS3 float64
	CamS4 float64
	
	// Woodie Pivot
	WoodPivot float64
	WoodR1    float64
	WoodR2    float64
	WoodS1    float64
	WoodS2    float64
	
	// DeMark Pivot
	DeMarkHigh float64
	DeMarkLow  float64
	
	// Current price position
	AbovePivot    bool
	NearestLevel  float64
	NearestType   string
	PricePosition string // "above_r1", "between_pivot_r1", etc.
}

// CalculatePivotPoints calculates all pivot point types
func CalculatePivotPoints(candles []Candle) *PivotPoints {
	if len(candles) < 2 {
		return nil
	}
	
	// Use previous day/period for pivot calculation
	prev := candles[len(candles)-2]
	current := candles[len(candles)-1]
	
	high := prev.High
	low := prev.Low
	close := prev.Close
	open := prev.Open
	
	pp := &PivotPoints{}
	
	// Standard Pivot Points
	pp.Pivot = (high + low + close) / 3
	pp.R1 = (2 * pp.Pivot) - low
	pp.S1 = (2 * pp.Pivot) - high
	pp.R2 = pp.Pivot + (high - low)
	pp.S2 = pp.Pivot - (high - low)
	pp.R3 = high + 2*(pp.Pivot-low)
	pp.S3 = low - 2*(high-pp.Pivot)
	
	// Fibonacci Pivot Points
	diff := high - low
	pp.FibR1 = pp.Pivot + (0.382 * diff)
	pp.FibR2 = pp.Pivot + (0.618 * diff)
	pp.FibR3 = pp.Pivot + diff
	pp.FibS1 = pp.Pivot - (0.382 * diff)
	pp.FibS2 = pp.Pivot - (0.618 * diff)
	pp.FibS3 = pp.Pivot - diff
	
	// Camarilla Pivot Points
	pp.CamR1 = close + (diff * 1.1 / 12)
	pp.CamR2 = close + (diff * 1.1 / 6)
	pp.CamR3 = close + (diff * 1.1 / 4)
	pp.CamR4 = close + (diff * 1.1 / 2)
	pp.CamS1 = close - (diff * 1.1 / 12)
	pp.CamS2 = close - (diff * 1.1 / 6)
	pp.CamS3 = close - (diff * 1.1 / 4)
	pp.CamS4 = close - (diff * 1.1 / 2)
	
	// Woodie Pivot Points
	pp.WoodPivot = (high + low + 2*close) / 4
	pp.WoodR1 = (2 * pp.WoodPivot) - low
	pp.WoodR2 = pp.WoodPivot + diff
	pp.WoodS1 = (2 * pp.WoodPivot) - high
	pp.WoodS2 = pp.WoodPivot - diff
	
	// DeMark Pivot Points
	var x float64
	if close < open {
		x = high + (2 * low) + close
	} else if close > open {
		x = (2 * high) + low + close
	} else {
		x = high + low + (2 * close)
	}
	pp.DeMarkHigh = x/2 - low
	pp.DeMarkLow = x/2 - high
	
	// Determine current price position
	currentPrice := current.Close
	pp.AbovePivot = currentPrice > pp.Pivot
	
	// Find nearest level
	levels := []struct {
		price float64
		name  string
	}{
		{pp.R3, "R3"}, {pp.R2, "R2"}, {pp.R1, "R1"},
		{pp.Pivot, "Pivot"},
		{pp.S1, "S1"}, {pp.S2, "S2"}, {pp.S3, "S3"},
	}
	
	minDist := math.MaxFloat64
	for _, l := range levels {
		dist := math.Abs(currentPrice - l.price)
		if dist < minDist {
			minDist = dist
			pp.NearestLevel = l.price
			pp.NearestType = l.name
		}
	}
	
	// Determine price position zone
	if currentPrice > pp.R2 {
		pp.PricePosition = "above_r2"
	} else if currentPrice > pp.R1 {
		pp.PricePosition = "between_r1_r2"
	} else if currentPrice > pp.Pivot {
		pp.PricePosition = "between_pivot_r1"
	} else if currentPrice > pp.S1 {
		pp.PricePosition = "between_s1_pivot"
	} else if currentPrice > pp.S2 {
		pp.PricePosition = "between_s2_s1"
	} else {
		pp.PricePosition = "below_s2"
	}
	
	return pp
}

// GetPivotSignal returns trading signal based on pivot analysis
func GetPivotSignal(pp *PivotPoints, currentPrice float64) (string, float64) {
	if pp == nil {
		return "neutral", 0
	}
	
	direction := "neutral"
	strength := 0.0
	
	// Price at support = bullish
	// Price at resistance = bearish
	
	distToS1 := math.Abs(currentPrice - pp.S1)
	distToS2 := math.Abs(currentPrice - pp.S2)
	distToR1 := math.Abs(currentPrice - pp.R1)
	distToR2 := math.Abs(currentPrice - pp.R2)
	distToPivot := math.Abs(currentPrice - pp.Pivot)
	
	threshold := (pp.R1 - pp.S1) * 0.05 // 5% of range
	
	// Near support levels = bullish
	if distToS1 < threshold || distToS2 < threshold {
		direction = "bullish"
		strength = 70
		if distToS2 < threshold {
			strength = 80 // Stronger support
		}
	}
	
	// Near resistance levels = bearish
	if distToR1 < threshold || distToR2 < threshold {
		direction = "bearish"
		strength = 70
		if distToR2 < threshold {
			strength = 80
		}
	}
	
	// Near pivot = watch for breakout
	if distToPivot < threshold {
		if pp.AbovePivot {
			direction = "bullish"
			strength = 60
		} else {
			direction = "bearish"
			strength = 60
		}
	}
	
	// Camarilla breakout signals
	if currentPrice > pp.CamR4 {
		direction = "bullish"
		strength = 85 // Strong breakout
	} else if currentPrice < pp.CamS4 {
		direction = "bearish"
		strength = 85
	}
	
	return direction, strength
}


// ==================== VWAP ANALYSIS ====================
// Volume Weighted Average Price

// VWAPAnalysis holds VWAP data
type VWAPAnalysis struct {
	VWAP          float64
	UpperBand1    float64 // +1 std dev
	UpperBand2    float64 // +2 std dev
	LowerBand1    float64 // -1 std dev
	LowerBand2    float64 // -2 std dev
	AboveVWAP     bool
	Deviation     float64 // How far from VWAP
	Signal        string
	Strength      float64
}

// CalculateVWAP calculates VWAP and bands
func CalculateVWAP(candles []Candle) *VWAPAnalysis {
	if len(candles) < 10 {
		return nil
	}
	
	va := &VWAPAnalysis{}
	
	// Calculate VWAP
	sumPV := 0.0 // Price * Volume
	sumV := 0.0  // Volume
	
	for _, c := range candles {
		typicalPrice := (c.High + c.Low + c.Close) / 3
		sumPV += typicalPrice * c.Volume
		sumV += c.Volume
	}
	
	if sumV == 0 {
		return nil
	}
	
	va.VWAP = sumPV / sumV
	
	// Calculate standard deviation
	sumSqDiff := 0.0
	for _, c := range candles {
		typicalPrice := (c.High + c.Low + c.Close) / 3
		diff := typicalPrice - va.VWAP
		sumSqDiff += diff * diff * c.Volume
	}
	
	variance := sumSqDiff / sumV
	stdDev := math.Sqrt(variance)
	
	// Calculate bands
	va.UpperBand1 = va.VWAP + stdDev
	va.UpperBand2 = va.VWAP + (2 * stdDev)
	va.LowerBand1 = va.VWAP - stdDev
	va.LowerBand2 = va.VWAP - (2 * stdDev)
	
	// Current price analysis
	currentPrice := candles[len(candles)-1].Close
	va.AboveVWAP = currentPrice > va.VWAP
	va.Deviation = ((currentPrice - va.VWAP) / va.VWAP) * 100
	
	// Generate signal
	if currentPrice < va.LowerBand2 {
		va.Signal = "bullish" // Oversold
		va.Strength = 85
	} else if currentPrice < va.LowerBand1 {
		va.Signal = "bullish"
		va.Strength = 70
	} else if currentPrice > va.UpperBand2 {
		va.Signal = "bearish" // Overbought
		va.Strength = 85
	} else if currentPrice > va.UpperBand1 {
		va.Signal = "bearish"
		va.Strength = 70
	} else if va.AboveVWAP {
		va.Signal = "bullish"
		va.Strength = 55
	} else {
		va.Signal = "bearish"
		va.Strength = 55
	}
	
	return va
}

// ==================== MARKET PROFILE ====================
// Value Area and POC analysis

// MarketProfile holds market profile data
type MarketProfile struct {
	POC           float64   // Point of Control
	VAH           float64   // Value Area High
	VAL           float64   // Value Area Low
	ValueArea     float64   // Value Area percentage
	PriceLevels   []float64
	VolumeLevels  []float64
	HVN           []float64 // High Volume Nodes
	LVN           []float64 // Low Volume Nodes
	InValueArea   bool
	Signal        string
	Strength      float64
}

// CalculateMarketProfile calculates market profile
func CalculateMarketProfile(candles []Candle, levels int) *MarketProfile {
	if len(candles) < 20 || levels < 10 {
		return nil
	}
	
	mp := &MarketProfile{
		PriceLevels:  make([]float64, levels),
		VolumeLevels: make([]float64, levels),
	}
	
	// Find price range
	high := candles[0].High
	low := candles[0].Low
	for _, c := range candles {
		if c.High > high {
			high = c.High
		}
		if c.Low < low {
			low = c.Low
		}
	}
	
	priceRange := high - low
	levelSize := priceRange / float64(levels)
	
	// Initialize price levels
	for i := 0; i < levels; i++ {
		mp.PriceLevels[i] = low + (float64(i)+0.5)*levelSize
	}
	
	// Distribute volume to levels
	for _, c := range candles {
		candleRange := c.High - c.Low
		if candleRange == 0 {
			continue
		}
		
		// Distribute candle volume across its range
		for i := 0; i < levels; i++ {
			levelLow := low + float64(i)*levelSize
			levelHigh := levelLow + levelSize
			
			// Check overlap
			overlapLow := math.Max(c.Low, levelLow)
			overlapHigh := math.Min(c.High, levelHigh)
			
			if overlapHigh > overlapLow {
				overlap := (overlapHigh - overlapLow) / candleRange
				mp.VolumeLevels[i] += c.Volume * overlap
			}
		}
	}
	
	// Find POC (highest volume level)
	maxVol := 0.0
	pocIdx := 0
	for i, vol := range mp.VolumeLevels {
		if vol > maxVol {
			maxVol = vol
			pocIdx = i
		}
	}
	mp.POC = mp.PriceLevels[pocIdx]
	
	// Calculate Value Area (70% of volume)
	totalVol := 0.0
	for _, vol := range mp.VolumeLevels {
		totalVol += vol
	}
	
	targetVol := totalVol * 0.7
	vaVol := mp.VolumeLevels[pocIdx]
	vaLowIdx := pocIdx
	vaHighIdx := pocIdx
	
	for vaVol < targetVol {
		// Expand value area
		expandUp := false
		expandDown := false
		
		if vaHighIdx < levels-1 {
			expandUp = true
		}
		if vaLowIdx > 0 {
			expandDown = true
		}
		
		if expandUp && expandDown {
			if mp.VolumeLevels[vaHighIdx+1] > mp.VolumeLevels[vaLowIdx-1] {
				vaHighIdx++
				vaVol += mp.VolumeLevels[vaHighIdx]
			} else {
				vaLowIdx--
				vaVol += mp.VolumeLevels[vaLowIdx]
			}
		} else if expandUp {
			vaHighIdx++
			vaVol += mp.VolumeLevels[vaHighIdx]
		} else if expandDown {
			vaLowIdx--
			vaVol += mp.VolumeLevels[vaLowIdx]
		} else {
			break
		}
	}
	
	mp.VAH = low + float64(vaHighIdx+1)*levelSize
	mp.VAL = low + float64(vaLowIdx)*levelSize
	mp.ValueArea = (vaVol / totalVol) * 100
	
	// Find HVN and LVN
	avgVol := totalVol / float64(levels)
	for i, vol := range mp.VolumeLevels {
		if vol > avgVol*1.5 {
			mp.HVN = append(mp.HVN, mp.PriceLevels[i])
		} else if vol < avgVol*0.5 {
			mp.LVN = append(mp.LVN, mp.PriceLevels[i])
		}
	}
	
	// Current price analysis
	currentPrice := candles[len(candles)-1].Close
	mp.InValueArea = currentPrice >= mp.VAL && currentPrice <= mp.VAH
	
	// Generate signal
	if currentPrice < mp.VAL {
		mp.Signal = "bullish" // Below value = potential long
		mp.Strength = 70
	} else if currentPrice > mp.VAH {
		mp.Signal = "bearish" // Above value = potential short
		mp.Strength = 70
	} else if currentPrice < mp.POC {
		mp.Signal = "bullish"
		mp.Strength = 55
	} else {
		mp.Signal = "bearish"
		mp.Strength = 55
	}
	
	return mp
}
