package patterns

import (
	"math"
)

// ==================== SUPPLY & DEMAND ZONES ====================
// Fresh zones = high probability trades

// SupplyDemandZone represents a supply or demand zone
type SupplyDemandZone struct {
	Type       string  // "supply" or "demand"
	High       float64
	Low        float64
	MidPoint   float64
	Origin     float64 // Price where zone originated
	Strength   float64 // Zone strength (0-100)
	Fresh      bool    // Has zone been tested?
	TouchCount int     // Number of times tested
	Volume     float64 // Volume at zone creation
	CandleIdx  int
	Quality    string  // "excellent", "good", "moderate", "weak"
}

// SDAnalysis holds supply/demand analysis
type SDAnalysis struct {
	SupplyZones  []SupplyDemandZone
	DemandZones  []SupplyDemandZone
	NearestSupply *SupplyDemandZone
	NearestDemand *SupplyDemandZone
	InZone       bool
	ZoneType     string
}

// ==================== ZONE DETECTION ====================

// FindSupplyDemandZones identifies all S/D zones
func FindSupplyDemandZones(candles []Candle) *SDAnalysis {
	analysis := &SDAnalysis{
		SupplyZones: []SupplyDemandZone{},
		DemandZones: []SupplyDemandZone{},
	}
	
	if len(candles) < 20 {
		return analysis
	}
	
	// Find demand zones (support - strong rally from base)
	for i := 5; i < len(candles)-5; i++ {
		// Look for base (consolidation)
		baseCandles := candles[i-3 : i+1]
		
		// Check if it's a base (low volatility)
		baseHigh := baseCandles[0].High
		baseLow := baseCandles[0].Low
		
		for _, c := range baseCandles {
			if c.High > baseHigh {
				baseHigh = c.High
			}
			if c.Low < baseLow {
				baseLow = c.Low
			}
		}
		
		baseSize := baseHigh - baseLow
		
		// Calculate average candle size
		avgSize := 0.0
		for j := i - 10; j < i; j++ {
			if j >= 0 {
				avgSize += candles[j].High - candles[j].Low
			}
		}
		avgSize /= 10
		
		// Base should be relatively small (consolidation)
		if baseSize < avgSize*1.5 {
			// Look for strong rally after base
			rallySize := 0.0
			rallyCandles := 0
			
			for j := i + 1; j < i + 6 && j < len(candles); j++ {
				if candles[j].Close > candles[j].Open {
					rallySize += candles[j].Close - candles[j].Open
					rallyCandles++
				}
			}
			
			// Strong rally = demand zone
			if rallySize > baseSize*2 && rallyCandles >= 3 {
				// Calculate zone strength
				strength := 50.0
				
				// Rally strength
				rallyStrength := rallySize / baseSize
				strength += math.Min(rallyStrength*10, 30)
				
				// Volume at base
				baseVolume := 0.0
				for _, c := range baseCandles {
					baseVolume += c.Volume
				}
				baseVolume /= float64(len(baseCandles))
				
				avgVolume := 0.0
				for j := i - 10; j < i; j++ {
					if j >= 0 {
						avgVolume += candles[j].Volume
					}
				}
				avgVolume /= 10
				
				if baseVolume > avgVolume*1.2 {
					strength += 15
				}
				
				// Check if zone is fresh (not tested yet)
				fresh := true
				touchCount := 0
				
				for j := i + 6; j < len(candles); j++ {
					if candles[j].Low <= baseHigh && candles[j].Low >= baseLow {
						touchCount++
						if touchCount > 0 {
							fresh = false
						}
					}
				}
				
				// Fresh zones are stronger
				if fresh {
					strength += 20
				} else {
					strength -= float64(touchCount) * 10
				}
				
				// Determine quality
				quality := "weak"
				if strength >= 85 {
					quality = "excellent"
				} else if strength >= 70 {
					quality = "good"
				} else if strength >= 55 {
					quality = "moderate"
				}
				
				zone := SupplyDemandZone{
					Type:       "demand",
					High:       baseHigh,
					Low:        baseLow,
					MidPoint:   (baseHigh + baseLow) / 2,
					Origin:     baseLow,
					Strength:   math.Max(0, math.Min(strength, 100)),
					Fresh:      fresh,
					TouchCount: touchCount,
					Volume:     baseVolume,
					CandleIdx:  i,
					Quality:    quality,
				}
				
				analysis.DemandZones = append(analysis.DemandZones, zone)
			}
		}
		
		// Find supply zones (resistance - strong drop from base)
		if baseSize < avgSize*1.5 {
			dropSize := 0.0
			dropCandles := 0
			
			for j := i + 1; j < i + 6 && j < len(candles); j++ {
				if candles[j].Close < candles[j].Open {
					dropSize += candles[j].Open - candles[j].Close
					dropCandles++
				}
			}
			
			// Strong drop = supply zone
			if dropSize > baseSize*2 && dropCandles >= 3 {
				strength := 50.0
				
				dropStrength := dropSize / baseSize
				strength += math.Min(dropStrength*10, 30)
				
				baseVolume := 0.0
				for _, c := range baseCandles {
					baseVolume += c.Volume
				}
				baseVolume /= float64(len(baseCandles))
				
				avgVolume := 0.0
				for j := i - 10; j < i; j++ {
					if j >= 0 {
						avgVolume += candles[j].Volume
					}
				}
				avgVolume /= 10
				
				if baseVolume > avgVolume*1.2 {
					strength += 15
				}
				
				fresh := true
				touchCount := 0
				
				for j := i + 6; j < len(candles); j++ {
					if candles[j].High >= baseLow && candles[j].High <= baseHigh {
						touchCount++
						if touchCount > 0 {
							fresh = false
						}
					}
				}
				
				if fresh {
					strength += 20
				} else {
					strength -= float64(touchCount) * 10
				}
				
				quality := "weak"
				if strength >= 85 {
					quality = "excellent"
				} else if strength >= 70 {
					quality = "good"
				} else if strength >= 55 {
					quality = "moderate"
				}
				
				zone := SupplyDemandZone{
					Type:       "supply",
					High:       baseHigh,
					Low:        baseLow,
					MidPoint:   (baseHigh + baseLow) / 2,
					Origin:     baseHigh,
					Strength:   math.Max(0, math.Min(strength, 100)),
					Fresh:      fresh,
					TouchCount: touchCount,
					Volume:     baseVolume,
					CandleIdx:  i,
					Quality:    quality,
				}
				
				analysis.SupplyZones = append(analysis.SupplyZones, zone)
			}
		}
	}
	
	// Find nearest zones to current price
	if len(candles) > 0 {
		currentPrice := candles[len(candles)-1].Close
		
		// Find nearest supply above
		minDistSupply := math.MaxFloat64
		for i := range analysis.SupplyZones {
			zone := &analysis.SupplyZones[i]
			if zone.Low > currentPrice {
				dist := zone.Low - currentPrice
				if dist < minDistSupply {
					minDistSupply = dist
					analysis.NearestSupply = zone
				}
			}
		}
		
		// Find nearest demand below
		minDistDemand := math.MaxFloat64
		for i := range analysis.DemandZones {
			zone := &analysis.DemandZones[i]
			if zone.High < currentPrice {
				dist := currentPrice - zone.High
				if dist < minDistDemand {
					minDistDemand = dist
					analysis.NearestDemand = zone
				}
			}
		}
		
		// Check if currently in a zone
		for i := range analysis.SupplyZones {
			zone := &analysis.SupplyZones[i]
			if currentPrice >= zone.Low && currentPrice <= zone.High {
				analysis.InZone = true
				analysis.ZoneType = "supply"
				break
			}
		}
		
		if !analysis.InZone {
			for i := range analysis.DemandZones {
				zone := &analysis.DemandZones[i]
				if currentPrice >= zone.Low && currentPrice <= zone.High {
					analysis.InZone = true
					analysis.ZoneType = "demand"
					break
				}
			}
		}
	}
	
	return analysis
}

// GetSDSignalStrength returns signal strength based on S/D analysis
func GetSDSignalStrength(analysis *SDAnalysis, direction string, currentPrice float64) float64 {
	if analysis == nil {
		return 0
	}
	
	strength := 0.0
	
	// In zone bonus
	if analysis.InZone {
		if direction == "bullish" && analysis.ZoneType == "demand" {
			// Find the zone we're in
			for _, zone := range analysis.DemandZones {
				if currentPrice >= zone.Low && currentPrice <= zone.High {
					strength += zone.Strength * 0.5
					
					// Fresh zone bonus
					if zone.Fresh {
						strength += 25
					}
					
					// Quality bonus
					if zone.Quality == "excellent" {
						strength += 15
					} else if zone.Quality == "good" {
						strength += 10
					}
					break
				}
			}
		} else if direction == "bearish" && analysis.ZoneType == "supply" {
			for _, zone := range analysis.SupplyZones {
				if currentPrice >= zone.Low && currentPrice <= zone.High {
					strength += zone.Strength * 0.5
					
					if zone.Fresh {
						strength += 25
					}
					
					if zone.Quality == "excellent" {
						strength += 15
					} else if zone.Quality == "good" {
						strength += 10
					}
					break
				}
			}
		}
	}
	
	// Near zone bonus
	if !analysis.InZone {
		if direction == "bullish" && analysis.NearestDemand != nil {
			dist := currentPrice - analysis.NearestDemand.High
			avgRange := (analysis.NearestDemand.High - analysis.NearestDemand.Low) * 2
			
			// Within 2x zone size
			if dist < avgRange {
				proximity := 1.0 - (dist / avgRange)
				strength += analysis.NearestDemand.Strength * proximity * 0.3
				
				if analysis.NearestDemand.Fresh {
					strength += 15 * proximity
				}
			}
		} else if direction == "bearish" && analysis.NearestSupply != nil {
			dist := analysis.NearestSupply.Low - currentPrice
			avgRange := (analysis.NearestSupply.High - analysis.NearestSupply.Low) * 2
			
			if dist < avgRange {
				proximity := 1.0 - (dist / avgRange)
				strength += analysis.NearestSupply.Strength * proximity * 0.3
				
				if analysis.NearestSupply.Fresh {
					strength += 15 * proximity
				}
			}
		}
	}
	
	return math.Max(0, math.Min(strength, 100))
}

// ShouldTradeSDZone determines if we should trade a S/D zone
func ShouldTradeSDZone(zone *SupplyDemandZone) bool {
	if zone == nil {
		return false
	}
	
	// Only trade fresh zones
	if !zone.Fresh {
		return false
	}
	
	// Only trade good quality or better
	if zone.Quality == "weak" || zone.Quality == "moderate" {
		return false
	}
	
	// Strength must be high
	if zone.Strength < 70 {
		return false
	}
	
	return true
}
