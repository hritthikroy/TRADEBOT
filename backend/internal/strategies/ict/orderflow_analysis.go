package ict

import (
	"math"
)

// ==================== ORDER FLOW & DELTA ANALYSIS ====================

// OrderFlowDelta represents volume delta analysis for order flow
type OrderFlowDelta struct {
	Delta           float64 // Buy volume - Sell volume
	CumulativeDelta float64 // Running total of delta
	DeltaPercent    float64 // Delta as percentage of total volume
	Trend           string  // "bullish", "bearish", "neutral"
	Divergence      bool    // Price/Delta divergence detected
	Strength        float64 // 0-100
}

// FootprintData represents footprint chart analysis
type FootprintData struct {
	BuyVolume       float64
	SellVolume      float64
	Delta           float64
	Imbalance       float64 // Buy/Sell imbalance ratio
	POC             float64 // Point of Control (highest volume price)
	VAH             float64 // Value Area High
	VAL             float64 // Value Area Low
	HighVolumeNodes []float64
	LowVolumeNodes  []float64
}

// OrderFlowAnalysis represents complete order flow analysis
type OrderFlowAnalysis struct {
	Delta           OrderFlowDelta
	Footprint       FootprintData
	Absorption      bool    // Large orders absorbed
	Exhaustion      bool    // Buying/Selling exhaustion
	Imbalance       float64 // Order imbalance
	AggressiveBuys  float64 // Market buy orders
	AggressiveSells float64 // Market sell orders
	Strength        float64 // Overall strength 0-100
}

// VolumeProfile represents volume at price levels
type VolumeProfile struct {
	POC        float64            // Point of Control
	VAH        float64            // Value Area High
	VAL        float64            // Value Area Low
	HVN        []float64          // High Volume Nodes
	LVN        []float64          // Low Volume Nodes
	VolumeMap  map[float64]float64 // Price -> Volume
	TotalVolume float64
}

// ==================== DELTA CALCULATION ====================

// CalculateOrderFlowDelta estimates delta from OHLCV data
func CalculateOrderFlowDelta(candles []Candle) *OrderFlowDelta {
	if len(candles) < 5 {
		return nil
	}

	da := &OrderFlowDelta{}
	cumulativeDelta := 0.0
	totalVolume := 0.0

	for i := 1; i < len(candles); i++ {
		c := candles[i]
		
		// Estimate buy/sell volume from price action
		// If close > open, more buying pressure
		// If close < open, more selling pressure
		candleRange := c.High - c.Low
		if candleRange == 0 {
			continue
		}

		// Calculate buying and selling pressure
		buyPressure := 0.0
		sellPressure := 0.0

		if c.Close > c.Open {
			// Bullish candle - estimate buy volume
			bodyRatio := (c.Close - c.Open) / candleRange
			buyPressure = c.Volume * (0.5 + bodyRatio*0.5)
			sellPressure = c.Volume - buyPressure
		} else if c.Close < c.Open {
			// Bearish candle - estimate sell volume
			bodyRatio := (c.Open - c.Close) / candleRange
			sellPressure = c.Volume * (0.5 + bodyRatio*0.5)
			buyPressure = c.Volume - sellPressure
		} else {
			// Doji - equal pressure
			buyPressure = c.Volume * 0.5
			sellPressure = c.Volume * 0.5
		}

		// Wick analysis for more accuracy
		upperWick := c.High - math.Max(c.Open, c.Close)
		lowerWick := math.Min(c.Open, c.Close) - c.Low

		// Upper wick = selling pressure, Lower wick = buying pressure
		if candleRange > 0 {
			wickSellPressure := (upperWick / candleRange) * c.Volume * 0.3
			wickBuyPressure := (lowerWick / candleRange) * c.Volume * 0.3
			sellPressure += wickSellPressure
			buyPressure += wickBuyPressure
		}

		delta := buyPressure - sellPressure
		cumulativeDelta += delta
		totalVolume += c.Volume
	}

	// Calculate recent delta (last 5 candles)
	recentDelta := 0.0
	recentVolume := 0.0
	startIdx := len(candles) - 5
	if startIdx < 0 {
		startIdx = 0
	}

	for i := startIdx; i < len(candles); i++ {
		c := candles[i]
		candleRange := c.High - c.Low
		if candleRange == 0 {
			continue
		}

		if c.Close > c.Open {
			bodyRatio := (c.Close - c.Open) / candleRange
			buyVol := c.Volume * (0.5 + bodyRatio*0.5)
			recentDelta += buyVol - (c.Volume - buyVol)
		} else {
			bodyRatio := (c.Open - c.Close) / candleRange
			sellVol := c.Volume * (0.5 + bodyRatio*0.5)
			recentDelta += (c.Volume - sellVol) - sellVol
		}
		recentVolume += c.Volume
	}

	da.Delta = recentDelta
	da.CumulativeDelta = cumulativeDelta

	if recentVolume > 0 {
		da.DeltaPercent = (recentDelta / recentVolume) * 100
	}

	// Determine trend
	if da.DeltaPercent > 10 {
		da.Trend = "bullish"
	} else if da.DeltaPercent < -10 {
		da.Trend = "bearish"
	} else {
		da.Trend = "neutral"
	}

	// Check for divergence
	priceChange := candles[len(candles)-1].Close - candles[startIdx].Close
	if (priceChange > 0 && recentDelta < 0) || (priceChange < 0 && recentDelta > 0) {
		da.Divergence = true
	}

	// Calculate strength
	da.Strength = math.Min(math.Abs(da.DeltaPercent)*2, 100)

	return da
}

// ==================== FOOTPRINT ANALYSIS ====================

// CalculateFootprint creates footprint analysis from candles
func CalculateFootprint(candles []Candle) *FootprintData {
	if len(candles) < 10 {
		return nil
	}

	fp := &FootprintData{}
	
	totalBuyVol := 0.0
	totalSellVol := 0.0
	priceVolume := make(map[float64]float64)

	for _, c := range candles {
		candleRange := c.High - c.Low
		if candleRange == 0 {
			continue
		}

		// Estimate buy/sell volume
		if c.Close > c.Open {
			bodyRatio := (c.Close - c.Open) / candleRange
			buyVol := c.Volume * (0.5 + bodyRatio*0.5)
			totalBuyVol += buyVol
			totalSellVol += c.Volume - buyVol
		} else {
			bodyRatio := (c.Open - c.Close) / candleRange
			sellVol := c.Volume * (0.5 + bodyRatio*0.5)
			totalSellVol += sellVol
			totalBuyVol += c.Volume - sellVol
		}

		// Build volume profile
		midPrice := (c.High + c.Low) / 2
		priceVolume[midPrice] += c.Volume
	}

	fp.BuyVolume = totalBuyVol
	fp.SellVolume = totalSellVol
	fp.Delta = totalBuyVol - totalSellVol

	// Calculate imbalance
	if totalSellVol > 0 {
		fp.Imbalance = totalBuyVol / totalSellVol
	}

	// Find POC (Point of Control)
	maxVol := 0.0
	for price, vol := range priceVolume {
		if vol > maxVol {
			maxVol = vol
			fp.POC = price
		}
	}

	// Calculate Value Area (70% of volume)
	totalVol := totalBuyVol + totalSellVol
	_ = totalVol * 0.7 // targetVol for future use
	
	// Simple VAH/VAL calculation
	prices := make([]float64, 0, len(priceVolume))
	for p := range priceVolume {
		prices = append(prices, p)
	}
	
	if len(prices) > 0 {
		// Sort prices
		for i := 0; i < len(prices)-1; i++ {
			for j := i + 1; j < len(prices); j++ {
				if prices[i] > prices[j] {
					prices[i], prices[j] = prices[j], prices[i]
				}
			}
		}
		
		fp.VAL = prices[0]
		fp.VAH = prices[len(prices)-1]
		
		// Find HVN and LVN
		avgVol := totalVol / float64(len(prices))
		for _, p := range prices {
			if priceVolume[p] > avgVol*1.5 {
				fp.HighVolumeNodes = append(fp.HighVolumeNodes, p)
			} else if priceVolume[p] < avgVol*0.5 {
				fp.LowVolumeNodes = append(fp.LowVolumeNodes, p)
			}
		}
	}

	return fp
}

// ==================== VOLUME PROFILE ====================

// CalculateVolumeProfile creates a volume profile
func CalculateVolumeProfile(candles []Candle, levels int) *VolumeProfile {
	if len(candles) < 10 || levels < 5 {
		return nil
	}

	vp := &VolumeProfile{
		VolumeMap: make(map[float64]float64),
	}

	// Find price range
	highestHigh := candles[0].High
	lowestLow := candles[0].Low

	for _, c := range candles {
		if c.High > highestHigh {
			highestHigh = c.High
		}
		if c.Low < lowestLow {
			lowestLow = c.Low
		}
		vp.TotalVolume += c.Volume
	}

	priceRange := highestHigh - lowestLow
	if priceRange == 0 {
		return nil
	}

	levelSize := priceRange / float64(levels)

	// Distribute volume across price levels
	for _, c := range candles {
		// Distribute candle volume across its range
		candleLevels := int((c.High - c.Low) / levelSize)
		if candleLevels < 1 {
			candleLevels = 1
		}
		volPerLevel := c.Volume / float64(candleLevels)

		for price := c.Low; price <= c.High; price += levelSize {
			roundedPrice := math.Round(price/levelSize) * levelSize
			vp.VolumeMap[roundedPrice] += volPerLevel
		}
	}

	// Find POC
	maxVol := 0.0
	for price, vol := range vp.VolumeMap {
		if vol > maxVol {
			maxVol = vol
			vp.POC = price
		}
	}

	// Calculate Value Area (70% of volume)
	targetVol := vp.TotalVolume * 0.7
	accumulatedVol := vp.VolumeMap[vp.POC]
	vp.VAH = vp.POC
	vp.VAL = vp.POC

	// Expand from POC until 70% volume captured
	for accumulatedVol < targetVol {
		upperVol := vp.VolumeMap[vp.VAH+levelSize]
		lowerVol := vp.VolumeMap[vp.VAL-levelSize]

		if upperVol >= lowerVol && upperVol > 0 {
			vp.VAH += levelSize
			accumulatedVol += upperVol
		} else if lowerVol > 0 {
			vp.VAL -= levelSize
			accumulatedVol += lowerVol
		} else {
			break
		}
	}

	// Find HVN and LVN
	avgVol := vp.TotalVolume / float64(len(vp.VolumeMap))
	for price, vol := range vp.VolumeMap {
		if vol > avgVol*1.5 {
			vp.HVN = append(vp.HVN, price)
		} else if vol < avgVol*0.5 {
			vp.LVN = append(vp.LVN, price)
		}
	}

	return vp
}


// ==================== ORDER FLOW ANALYSIS ====================

// PerformOrderFlowAnalysis performs complete order flow analysis
func PerformOrderFlowAnalysis(candles []Candle) *OrderFlowAnalysis {
	if len(candles) < 20 {
		return nil
	}

	ofa := &OrderFlowAnalysis{}

	// Delta Analysis
	delta := CalculateOrderFlowDelta(candles)
	if delta != nil {
		ofa.Delta = *delta
	}

	// Footprint Analysis
	footprint := CalculateFootprint(candles)
	if footprint != nil {
		ofa.Footprint = *footprint
	}

	// Calculate aggressive orders (market orders)
	recentCandles := candles[len(candles)-10:]
	for _, c := range recentCandles {
		candleRange := c.High - c.Low
		if candleRange == 0 {
			continue
		}

		// Large body candles indicate aggressive orders
		body := math.Abs(c.Close - c.Open)
		bodyRatio := body / candleRange

		if bodyRatio > 0.7 { // Strong candle
			if c.Close > c.Open {
				ofa.AggressiveBuys += c.Volume * bodyRatio
			} else {
				ofa.AggressiveSells += c.Volume * bodyRatio
			}
		}
	}

	// Calculate imbalance
	totalAggressive := ofa.AggressiveBuys + ofa.AggressiveSells
	if totalAggressive > 0 {
		ofa.Imbalance = (ofa.AggressiveBuys - ofa.AggressiveSells) / totalAggressive
	}

	// Detect absorption (large volume with small price movement)
	lastCandle := candles[len(candles)-1]
	avgVolume := 0.0
	avgRange := 0.0
	for _, c := range candles[len(candles)-20:] {
		avgVolume += c.Volume
		avgRange += c.High - c.Low
	}
	avgVolume /= 20
	avgRange /= 20

	lastRange := lastCandle.High - lastCandle.Low
	if lastCandle.Volume > avgVolume*1.5 && lastRange < avgRange*0.5 {
		ofa.Absorption = true
	}

	// Detect exhaustion (decreasing volume with price continuation)
	if len(candles) >= 5 {
		last5 := candles[len(candles)-5:]
		volumeDecreasing := true
		priceDirection := last5[4].Close - last5[0].Close

		for i := 1; i < len(last5); i++ {
			if last5[i].Volume > last5[i-1].Volume {
				volumeDecreasing = false
				break
			}
		}

		if volumeDecreasing && math.Abs(priceDirection) > avgRange {
			ofa.Exhaustion = true
		}
	}

	// Calculate overall strength
	strength := 50.0

	// Delta contribution
	strength += ofa.Delta.Strength * 0.3

	// Imbalance contribution
	strength += math.Abs(ofa.Imbalance) * 30

	// Absorption/Exhaustion signals
	if ofa.Absorption {
		strength += 10
	}
	if ofa.Exhaustion {
		strength += 15
	}

	ofa.Strength = math.Min(strength, 100)

	return ofa
}

// ==================== ENHANCED ORDER BLOCKS ====================

// EnhancedOrderBlock represents an order block with volume analysis
type EnhancedOrderBlock struct {
	Type           string  // "bullish" or "bearish"
	High           float64
	Low            float64
	MidPoint       float64
	Volume         float64
	Delta          float64 // Volume delta at this level
	Imbalance      float64 // Buy/Sell imbalance
	Strength       float64
	Mitigated      bool
	TouchCount     int     // Number of times price returned
	LastTouchIdx   int
	VolumeProfile  *VolumeProfile
	InstitutionalActivity bool // High volume + specific patterns
}

// FindEnhancedOrderBlocks finds order blocks with volume/delta analysis
func FindEnhancedOrderBlocks(candles []Candle) []EnhancedOrderBlock {
	var orderBlocks []EnhancedOrderBlock

	if len(candles) < 15 {
		return orderBlocks
	}

	// Calculate average volume
	avgVolume := 0.0
	for _, c := range candles {
		avgVolume += c.Volume
	}
	avgVolume /= float64(len(candles))

	for i := 3; i < len(candles)-2; i++ {
		prev := candles[i-1]
		curr := candles[i]
		next := candles[i+1]

		// Calculate delta for this candle
		candleRange := curr.High - curr.Low
		delta := 0.0
		if candleRange > 0 {
			if curr.Close > curr.Open {
				bodyRatio := (curr.Close - curr.Open) / candleRange
				buyVol := curr.Volume * (0.5 + bodyRatio*0.5)
				delta = buyVol - (curr.Volume - buyVol)
			} else {
				bodyRatio := (curr.Open - curr.Close) / candleRange
				sellVol := curr.Volume * (0.5 + bodyRatio*0.5)
				delta = (curr.Volume - sellVol) - sellVol
			}
		}

		// BULLISH ORDER BLOCK with enhanced criteria
		if prev.Close < prev.Open && // Previous bearish
			curr.Close < curr.Open && // Current bearish (the OB)
			next.Close > next.Open && // Next bullish
			next.Close > curr.High { // Strong move up

			// Check for institutional activity
			institutional := curr.Volume > avgVolume*1.5 || next.Volume > avgVolume*2

			// Calculate imbalance
			imbalance := 0.0
			if curr.Volume > 0 {
				imbalance = delta / curr.Volume
			}

			strength := 60.0
			if institutional {
				strength += 20
			}
			if curr.Volume > avgVolume*1.3 {
				strength += 10
			}
			if math.Abs(imbalance) > 0.3 {
				strength += 10
			}

			ob := EnhancedOrderBlock{
				Type:                  "bullish",
				High:                  curr.High,
				Low:                   curr.Low,
				MidPoint:              (curr.High + curr.Low) / 2,
				Volume:                curr.Volume,
				Delta:                 delta,
				Imbalance:             imbalance,
				Strength:              math.Min(strength, 100),
				Mitigated:             false,
				InstitutionalActivity: institutional,
			}

			orderBlocks = append(orderBlocks, ob)
		}

		// BEARISH ORDER BLOCK with enhanced criteria
		if prev.Close > prev.Open && // Previous bullish
			curr.Close > curr.Open && // Current bullish (the OB)
			next.Close < next.Open && // Next bearish
			next.Close < curr.Low { // Strong move down

			institutional := curr.Volume > avgVolume*1.5 || next.Volume > avgVolume*2

			imbalance := 0.0
			if curr.Volume > 0 {
				imbalance = delta / curr.Volume
			}

			strength := 60.0
			if institutional {
				strength += 20
			}
			if curr.Volume > avgVolume*1.3 {
				strength += 10
			}
			if math.Abs(imbalance) > 0.3 {
				strength += 10
			}

			ob := EnhancedOrderBlock{
				Type:                  "bearish",
				High:                  curr.High,
				Low:                   curr.Low,
				MidPoint:              (curr.High + curr.Low) / 2,
				Volume:                curr.Volume,
				Delta:                 delta,
				Imbalance:             imbalance,
				Strength:              math.Min(strength, 100),
				Mitigated:             false,
				InstitutionalActivity: institutional,
			}

			orderBlocks = append(orderBlocks, ob)
		}
	}

	// Check for mitigation and touch count
	_ = candles[len(candles)-1].Close // currentPrice for reference
	for i := range orderBlocks {
		ob := &orderBlocks[i]
		
		// Count touches
		for j := ob.LastTouchIdx + 1; j < len(candles); j++ {
			c := candles[j]
			if c.Low <= ob.High && c.High >= ob.Low {
				ob.TouchCount++
				ob.LastTouchIdx = j
				
				// Check if mitigated (price closed through)
				if ob.Type == "bullish" && c.Close < ob.Low {
					ob.Mitigated = true
				} else if ob.Type == "bearish" && c.Close > ob.High {
					ob.Mitigated = true
				}
			}
		}

		// Reduce strength for multiple touches
		if ob.TouchCount > 2 {
			ob.Strength *= 0.8
		}
	}

	return orderBlocks
}

// ==================== BREAKER BLOCKS ====================

// BreakerBlock represents a failed order block that becomes support/resistance
type BreakerBlock struct {
	Type       string  // "bullish" or "bearish"
	High       float64
	Low        float64
	MidPoint   float64
	Strength   float64
	OriginalOB *EnhancedOrderBlock
}

// FindBreakerBlocks finds breaker blocks (mitigated OBs that flip)
func FindBreakerBlocks(candles []Candle, orderBlocks []EnhancedOrderBlock) []BreakerBlock {
	var breakerBlocks []BreakerBlock

	for _, ob := range orderBlocks {
		if ob.Mitigated {
			// A mitigated bullish OB becomes bearish breaker
			// A mitigated bearish OB becomes bullish breaker
			breakerType := "bearish"
			if ob.Type == "bearish" {
				breakerType = "bullish"
			}

			bb := BreakerBlock{
				Type:       breakerType,
				High:       ob.High,
				Low:        ob.Low,
				MidPoint:   ob.MidPoint,
				Strength:   ob.Strength * 0.8, // Slightly weaker than original
				OriginalOB: &ob,
			}

			breakerBlocks = append(breakerBlocks, bb)
		}
	}

	return breakerBlocks
}

// ==================== MITIGATION BLOCKS ====================

// MitigationBlock represents a zone where orders were filled
type MitigationBlock struct {
	Type     string
	High     float64
	Low      float64
	Volume   float64
	Strength float64
}

// FindMitigationBlocks finds zones where large orders were filled
func FindMitigationBlocks(candles []Candle) []MitigationBlock {
	var mitigationBlocks []MitigationBlock

	if len(candles) < 20 {
		return mitigationBlocks
	}

	avgVolume := 0.0
	for _, c := range candles {
		avgVolume += c.Volume
	}
	avgVolume /= float64(len(candles))

	for i := 1; i < len(candles)-1; i++ {
		c := candles[i]

		// High volume candle with rejection (long wicks)
		candleRange := c.High - c.Low
		if candleRange == 0 {
			continue
		}

		body := math.Abs(c.Close - c.Open)
		upperWick := c.High - math.Max(c.Open, c.Close)
		lowerWick := math.Min(c.Open, c.Close) - c.Low

		// Mitigation = high volume + rejection
		if c.Volume > avgVolume*1.5 {
			// Upper rejection = bearish mitigation
			if upperWick > body*1.5 && upperWick > candleRange*0.4 {
				mb := MitigationBlock{
					Type:     "bearish",
					High:     c.High,
					Low:      c.High - upperWick,
					Volume:   c.Volume,
					Strength: 70 + (c.Volume/avgVolume-1)*10,
				}
				mitigationBlocks = append(mitigationBlocks, mb)
			}

			// Lower rejection = bullish mitigation
			if lowerWick > body*1.5 && lowerWick > candleRange*0.4 {
				mb := MitigationBlock{
					Type:     "bullish",
					High:     c.Low + lowerWick,
					Low:      c.Low,
					Volume:   c.Volume,
					Strength: 70 + (c.Volume/avgVolume-1)*10,
				}
				mitigationBlocks = append(mitigationBlocks, mb)
			}
		}
	}

	return mitigationBlocks
}
