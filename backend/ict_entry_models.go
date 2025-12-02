package main

import (
	"math"
)

// ==================== ICT ENTRY MODELS ====================

// EntryModel represents an ICT entry model
type EntryModel struct {
	Name       string
	Type       string  // "BUY" or "SELL"
	Entry      float64
	StopLoss   float64
	Target1    float64
	Target2    float64
	Target3    float64
	Confidence float64
	Confluence int
	Reason     string
}

// ==================== MODEL 1: ORDER BLOCK ENTRY ====================

// OrderBlockEntry - Enter at order block with FVG confluence
func OrderBlockEntry(candles []Candle, ict *ICTAnalysis) *EntryModel {
	if len(candles) < 20 || ict == nil {
		return nil
	}
	
	currentPrice := candles[len(candles)-1].Close
	atr := calculateATR(candles[len(candles)-14:], 14)
	
	// Look for unmitigated order blocks
	for _, ob := range ict.OrderBlocks {
		if ob.Mitigated {
			continue
		}
		
		// BULLISH OB ENTRY
		if ob.Type == "bullish" {
			// Price must be at or near the OB
			if currentPrice >= ob.Low && currentPrice <= ob.High*1.01 {
				// Check for FVG confluence
				hasFVG := false
				for _, fvg := range ict.FairValueGaps {
					if fvg.Type == "bullish" && 
					   math.Abs(fvg.MidPoint-ob.MidPoint) < atr {
						hasFVG = true
						break
					}
				}
				
				// Check market structure
				if ict.Structure.Trend != "bearish" {
					confluence := 2 // OB + price at level
					if hasFVG {
						confluence++
					}
					if ict.PremiumDiscount == "discount" {
						confluence++
					}
					if ict.OTE {
						confluence++
					}
					
					if confluence >= 3 {
						return &EntryModel{
							Name:       "Order Block Entry",
							Type:       "BUY",
							Entry:      currentPrice,
							StopLoss:   ob.Low - (atr * 0.5),
							Target1:    currentPrice + (atr * 2),
							Target2:    currentPrice + (atr * 3),
							Target3:    currentPrice + (atr * 4),
							Confidence: 60 + float64(confluence)*8,
							Confluence: confluence,
							Reason:     "Bullish OB with confluence",
						}
					}
				}
			}
		}
		
		// BEARISH OB ENTRY
		if ob.Type == "bearish" {
			// Price must be at or near the OB
			if currentPrice <= ob.High && currentPrice >= ob.Low*0.99 {
				// Check for FVG confluence
				hasFVG := false
				for _, fvg := range ict.FairValueGaps {
					if fvg.Type == "bearish" && 
					   math.Abs(fvg.MidPoint-ob.MidPoint) < atr {
						hasFVG = true
						break
					}
				}
				
				// Check market structure
				if ict.Structure.Trend != "bullish" {
					confluence := 2
					if hasFVG {
						confluence++
					}
					if ict.PremiumDiscount == "premium" {
						confluence++
					}
					if ict.OTE {
						confluence++
					}
					
					if confluence >= 3 {
						return &EntryModel{
							Name:       "Order Block Entry",
							Type:       "SELL",
							Entry:      currentPrice,
							StopLoss:   ob.High + (atr * 0.5),
							Target1:    currentPrice - (atr * 2),
							Target2:    currentPrice - (atr * 3),
							Target3:    currentPrice - (atr * 4),
							Confidence: 60 + float64(confluence)*8,
							Confluence: confluence,
							Reason:     "Bearish OB with confluence",
						}
					}
				}
			}
		}
	}
	
	return nil
}

// ==================== MODEL 2: FVG ENTRY ====================

// FairValueGapEntry - Enter at fair value gap fill
func FairValueGapEntry(candles []Candle, ict *ICTAnalysis) *EntryModel {
	if len(candles) < 20 || ict == nil {
		return nil
	}
	
	currentPrice := candles[len(candles)-1].Close
	atr := calculateATR(candles[len(candles)-14:], 14)
	
	for _, fvg := range ict.FairValueGaps {
		if fvg.Filled {
			continue
		}
		
		// BULLISH FVG ENTRY
		if fvg.Type == "bullish" {
			// Price entering the FVG from above
			if currentPrice <= fvg.High && currentPrice >= fvg.Low {
				confluence := 2
				
				// Check for OB confluence
				for _, ob := range ict.OrderBlocks {
					if ob.Type == "bullish" && 
					   math.Abs(ob.MidPoint-fvg.MidPoint) < atr {
						confluence++
						break
					}
				}
				
				if ict.PremiumDiscount == "discount" {
					confluence++
				}
				if ict.Structure.Trend == "bullish" {
					confluence++
				}
				
				if confluence >= 3 {
					return &EntryModel{
						Name:       "FVG Entry",
						Type:       "BUY",
						Entry:      currentPrice,
						StopLoss:   fvg.Low - (atr * 0.5),
						Target1:    currentPrice + (atr * 2),
						Target2:    currentPrice + (atr * 3),
						Target3:    currentPrice + (atr * 4),
						Confidence: 55 + float64(confluence)*8,
						Confluence: confluence,
						Reason:     "Bullish FVG fill",
					}
				}
			}
		}
		
		// BEARISH FVG ENTRY
		if fvg.Type == "bearish" {
			// Price entering the FVG from below
			if currentPrice >= fvg.Low && currentPrice <= fvg.High {
				confluence := 2
				
				for _, ob := range ict.OrderBlocks {
					if ob.Type == "bearish" && 
					   math.Abs(ob.MidPoint-fvg.MidPoint) < atr {
						confluence++
						break
					}
				}
				
				if ict.PremiumDiscount == "premium" {
					confluence++
				}
				if ict.Structure.Trend == "bearish" {
					confluence++
				}
				
				if confluence >= 3 {
					return &EntryModel{
						Name:       "FVG Entry",
						Type:       "SELL",
						Entry:      currentPrice,
						StopLoss:   fvg.High + (atr * 0.5),
						Target1:    currentPrice - (atr * 2),
						Target2:    currentPrice - (atr * 3),
						Target3:    currentPrice - (atr * 4),
						Confidence: 55 + float64(confluence)*8,
						Confluence: confluence,
						Reason:     "Bearish FVG fill",
					}
				}
			}
		}
	}
	
	return nil
}

// ==================== MODEL 3: LIQUIDITY SWEEP ENTRY ====================

// LiquiditySweepEntry - Enter after liquidity sweep
func LiquiditySweepEntry(candles []Candle, ict *ICTAnalysis) *EntryModel {
	if len(candles) < 20 || ict == nil {
		return nil
	}
	
	currentPrice := candles[len(candles)-1].Close
	prevCandle := candles[len(candles)-2]
	atr := calculateATR(candles[len(candles)-14:], 14)
	
	for _, liq := range ict.Liquidity {
		// SELL-SIDE LIQUIDITY SWEEP (bullish entry)
		if liq.Type == "sellside" && !liq.Swept {
			// Check if previous candle swept the liquidity
			if prevCandle.Low < liq.Price && currentPrice > liq.Price {
				// Liquidity was swept, now reversing
				confluence := 3 // Liquidity sweep + reversal
				
				if ict.PremiumDiscount == "discount" {
					confluence++
				}
				if ict.Structure.Trend == "bullish" || ict.Structure.CHOCH {
					confluence++
				}
				
				// Check for OB below
				for _, ob := range ict.OrderBlocks {
					if ob.Type == "bullish" && ob.High < currentPrice {
						confluence++
						break
					}
				}
				
				if confluence >= 4 {
					return &EntryModel{
						Name:       "Liquidity Sweep Entry",
						Type:       "BUY",
						Entry:      currentPrice,
						StopLoss:   prevCandle.Low - (atr * 0.3),
						Target1:    currentPrice + (atr * 2.5),
						Target2:    currentPrice + (atr * 4),
						Target3:    currentPrice + (atr * 6),
						Confidence: 70 + float64(confluence)*5,
						Confluence: confluence,
						Reason:     "Sell-side liquidity swept",
					}
				}
			}
		}
		
		// BUY-SIDE LIQUIDITY SWEEP (bearish entry)
		if liq.Type == "buyside" && !liq.Swept {
			// Check if previous candle swept the liquidity
			if prevCandle.High > liq.Price && currentPrice < liq.Price {
				// Liquidity was swept, now reversing
				confluence := 3
				
				if ict.PremiumDiscount == "premium" {
					confluence++
				}
				if ict.Structure.Trend == "bearish" || ict.Structure.CHOCH {
					confluence++
				}
				
				for _, ob := range ict.OrderBlocks {
					if ob.Type == "bearish" && ob.Low > currentPrice {
						confluence++
						break
					}
				}
				
				if confluence >= 4 {
					return &EntryModel{
						Name:       "Liquidity Sweep Entry",
						Type:       "SELL",
						Entry:      currentPrice,
						StopLoss:   prevCandle.High + (atr * 0.3),
						Target1:    currentPrice - (atr * 2.5),
						Target2:    currentPrice - (atr * 4),
						Target3:    currentPrice - (atr * 6),
						Confidence: 70 + float64(confluence)*5,
						Confluence: confluence,
						Reason:     "Buy-side liquidity swept",
					}
				}
			}
		}
	}
	
	return nil
}

// ==================== MODEL 4: BOS/CHOCH ENTRY ====================

// StructureBreakEntry - Enter on break of structure or change of character
func StructureBreakEntry(candles []Candle, ict *ICTAnalysis) *EntryModel {
	if len(candles) < 20 || ict == nil {
		return nil
	}
	
	currentPrice := candles[len(candles)-1].Close
	atr := calculateATR(candles[len(candles)-14:], 14)
	
	// BOS ENTRY (continuation)
	if ict.Structure.BOS {
		if ict.Structure.Trend == "bullish" {
			confluence := 3 // BOS + trend
			
			if ict.PremiumDiscount == "discount" {
				confluence++
			}
			
			// Look for pullback to OB or FVG
			for _, ob := range ict.OrderBlocks {
				if ob.Type == "bullish" && 
				   currentPrice >= ob.Low && currentPrice <= ob.High {
					confluence += 2
					break
				}
			}
			
			if confluence >= 4 {
				return &EntryModel{
					Name:       "BOS Entry",
					Type:       "BUY",
					Entry:      currentPrice,
					StopLoss:   ict.Structure.LastSwingLow - (atr * 0.3),
					Target1:    currentPrice + (atr * 2),
					Target2:    currentPrice + (atr * 3),
					Target3:    ict.Structure.LastSwingHigh + (atr * 2),
					Confidence: 65 + float64(confluence)*6,
					Confluence: confluence,
					Reason:     "Bullish BOS continuation",
				}
			}
		}
		
		if ict.Structure.Trend == "bearish" {
			confluence := 3
			
			if ict.PremiumDiscount == "premium" {
				confluence++
			}
			
			for _, ob := range ict.OrderBlocks {
				if ob.Type == "bearish" && 
				   currentPrice <= ob.High && currentPrice >= ob.Low {
					confluence += 2
					break
				}
			}
			
			if confluence >= 4 {
				return &EntryModel{
					Name:       "BOS Entry",
					Type:       "SELL",
					Entry:      currentPrice,
					StopLoss:   ict.Structure.LastSwingHigh + (atr * 0.3),
					Target1:    currentPrice - (atr * 2),
					Target2:    currentPrice - (atr * 3),
					Target3:    ict.Structure.LastSwingLow - (atr * 2),
					Confidence: 65 + float64(confluence)*6,
					Confluence: confluence,
					Reason:     "Bearish BOS continuation",
				}
			}
		}
	}
	
	// CHOCH ENTRY (reversal)
	if ict.Structure.CHOCH {
		// Bullish CHOCH (was bearish, now bullish)
		if currentPrice > ict.Structure.LastSwingHigh {
			confluence := 4 // CHOCH is strong signal
			
			if ict.PremiumDiscount == "discount" {
				confluence++
			}
			
			return &EntryModel{
				Name:       "CHOCH Entry",
				Type:       "BUY",
				Entry:      currentPrice,
				StopLoss:   ict.Structure.LastSwingLow - (atr * 0.5),
				Target1:    currentPrice + (atr * 3),
				Target2:    currentPrice + (atr * 5),
				Target3:    currentPrice + (atr * 7),
				Confidence: 75 + float64(confluence)*4,
				Confluence: confluence,
				Reason:     "Bullish CHOCH reversal",
			}
		}
		
		// Bearish CHOCH
		if currentPrice < ict.Structure.LastSwingLow {
			confluence := 4
			
			if ict.PremiumDiscount == "premium" {
				confluence++
			}
			
			return &EntryModel{
				Name:       "CHOCH Entry",
				Type:       "SELL",
				Entry:      currentPrice,
				StopLoss:   ict.Structure.LastSwingHigh + (atr * 0.5),
				Target1:    currentPrice - (atr * 3),
				Target2:    currentPrice - (atr * 5),
				Target3:    currentPrice - (atr * 7),
				Confidence: 75 + float64(confluence)*4,
				Confluence: confluence,
				Reason:     "Bearish CHOCH reversal",
			}
		}
	}
	
	return nil
}

// ==================== MODEL 5: OTE ENTRY ====================

// OTEEntry - Enter in Optimal Trade Entry zone
func OTEEntry(candles []Candle, ict *ICTAnalysis) *EntryModel {
	if len(candles) < 20 || ict == nil || !ict.OTE {
		return nil
	}
	
	currentPrice := candles[len(candles)-1].Close
	atr := calculateATR(candles[len(candles)-14:], 14)
	
	// Must have trend direction
	if ict.Structure.Trend == "bullish" {
		confluence := 3 // OTE + bullish trend
		
		// Check for additional confluence
		for _, ob := range ict.OrderBlocks {
			if ob.Type == "bullish" {
				confluence++
				break
			}
		}
		for _, fvg := range ict.FairValueGaps {
			if fvg.Type == "bullish" {
				confluence++
				break
			}
		}
		
		if confluence >= 4 {
			return &EntryModel{
				Name:       "OTE Entry",
				Type:       "BUY",
				Entry:      currentPrice,
				StopLoss:   ict.Structure.LastSwingLow - (atr * 0.3),
				Target1:    ict.Structure.LastSwingHigh,
				Target2:    ict.Structure.LastSwingHigh + (atr * 2),
				Target3:    ict.Structure.LastSwingHigh + (atr * 4),
				Confidence: 70 + float64(confluence)*5,
				Confluence: confluence,
				Reason:     "Bullish OTE zone entry",
			}
		}
	}
	
	if ict.Structure.Trend == "bearish" {
		confluence := 3
		
		for _, ob := range ict.OrderBlocks {
			if ob.Type == "bearish" {
				confluence++
				break
			}
		}
		for _, fvg := range ict.FairValueGaps {
			if fvg.Type == "bearish" {
				confluence++
				break
			}
		}
		
		if confluence >= 4 {
			return &EntryModel{
				Name:       "OTE Entry",
				Type:       "SELL",
				Entry:      currentPrice,
				StopLoss:   ict.Structure.LastSwingHigh + (atr * 0.3),
				Target1:    ict.Structure.LastSwingLow,
				Target2:    ict.Structure.LastSwingLow - (atr * 2),
				Target3:    ict.Structure.LastSwingLow - (atr * 4),
				Confidence: 70 + float64(confluence)*5,
				Confluence: confluence,
				Reason:     "Bearish OTE zone entry",
			}
		}
	}
	
	return nil
}

// ==================== BEST ENTRY SELECTOR ====================

// GetBestICTEntry evaluates all entry models and returns the best one
func GetBestICTEntry(candles []Candle) *EntryModel {
	if len(candles) < 50 {
		return nil
	}
	
	// Perform ICT analysis
	ict := PerformICTAnalysis(candles)
	
	// Try all entry models
	entries := []*EntryModel{
		LiquiditySweepEntry(candles, ict),  // Highest priority
		StructureBreakEntry(candles, ict),
		OrderBlockEntry(candles, ict),
		FairValueGapEntry(candles, ict),
		OTEEntry(candles, ict),
	}
	
	// Find the best entry (highest confidence with most confluence)
	var bestEntry *EntryModel
	bestScore := 0.0
	
	for _, entry := range entries {
		if entry != nil {
			score := entry.Confidence + float64(entry.Confluence)*10
			if score > bestScore {
				bestScore = score
				bestEntry = entry
			}
		}
	}
	
	return bestEntry
}
