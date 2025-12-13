package institutional

import (
	"math"
	"time"
)

// ==================== INSTITUTIONAL TRADING SETUPS ====================
// Professional-grade trading setups used by institutions

// ==================== SILVER BULLET SETUP ====================
// ICT Silver Bullet: High probability setup during specific times

// SilverBullet represents a Silver Bullet setup
type SilverBullet struct {
	Type       string  // "bullish" or "bearish"
	TimeWindow string  // "london", "ny_am", "ny_pm"
	FVG        *FairValueGap
	OrderBlock *OrderBlock
	Entry      float64
	StopLoss   float64
	Target     float64
	Strength   float64
	Valid      bool
}

// DetectSilverBullet detects Silver Bullet setups
func DetectSilverBullet(candles []Candle, currentTime time.Time) *SilverBullet {
	hour := currentTime.UTC().Hour()
	minute := currentTime.UTC().Minute()
	
	// Silver Bullet Windows (UTC):
	// 1. London: 10:00-11:00 UTC
	// 2. NY AM: 14:00-15:00 UTC (10:00-11:00 EST)
	// 3. NY PM: 19:00-20:00 UTC (15:00-16:00 EST)
	
	var timeWindow string
	isValidTime := false
	
	if hour == 10 && minute >= 0 && minute <= 60 {
		timeWindow = "london"
		isValidTime = true
	} else if hour == 14 && minute >= 0 && minute <= 60 {
		timeWindow = "ny_am"
		isValidTime = true
	} else if hour == 19 && minute >= 0 && minute <= 60 {
		timeWindow = "ny_pm"
		isValidTime = true
	}
	
	if !isValidTime || len(candles) < 20 {
		return nil
	}
	
	// Find FVG in recent candles
	fvgs := FindFairValueGaps(candles[len(candles)-20:])
	if len(fvgs) == 0 {
		return nil
	}
	
	// Find Order Blocks
	obs := FindOrderBlocks(candles[len(candles)-20:])
	
	// Get most recent FVG
	recentFVG := fvgs[len(fvgs)-1]
	
	// Find matching Order Block
	var matchingOB *OrderBlock
	for i := range obs {
		ob := &obs[i]
		if ob.Type == recentFVG.Type {
			// OB should be near FVG
			if math.Abs(ob.MidPoint-recentFVG.MidPoint) < (recentFVG.High-recentFVG.Low)*2 {
				matchingOB = ob
				break
			}
		}
	}
	
	currentPrice := candles[len(candles)-1].Close
	atr := calculateATR(candles[len(candles)-14:], 14)
	
	sb := &SilverBullet{
		TimeWindow: timeWindow,
		FVG:        &recentFVG,
		OrderBlock: matchingOB,
		Strength:   70,
		Valid:      false,
	}
	
	// Bullish Silver Bullet
	if recentFVG.Type == "bullish" {
		// Price should be at or below FVG
		if currentPrice <= recentFVG.High {
			sb.Type = "bullish"
			sb.Entry = recentFVG.MidPoint
			sb.StopLoss = recentFVG.Low - atr*0.3
			sb.Target = recentFVG.High + atr*2
			sb.Valid = true
			
			if matchingOB != nil {
				sb.Strength += 15
			}
		}
	}
	
	// Bearish Silver Bullet
	if recentFVG.Type == "bearish" {
		if currentPrice >= recentFVG.Low {
			sb.Type = "bearish"
			sb.Entry = recentFVG.MidPoint
			sb.StopLoss = recentFVG.High + atr*0.3
			sb.Target = recentFVG.Low - atr*2
			sb.Valid = true
			
			if matchingOB != nil {
				sb.Strength += 15
			}
		}
	}
	
	return sb
}

// ==================== UNICORN MODEL ====================
// ICT Unicorn: Breaker + FVG + OTE confluence

// UnicornSetup represents a Unicorn setup
type UnicornSetup struct {
	Type         string
	Breaker      *BreakerBlock
	FVG          *FairValueGap
	OTELevel     float64
	Entry        float64
	StopLoss     float64
	Target1      float64
	Target2      float64
	Strength     float64
	Confluence   int
	Valid        bool
}

// DetectUnicornSetup detects Unicorn setups
func DetectUnicornSetup(candles []Candle) *UnicornSetup {
	if len(candles) < 50 {
		return nil
	}
	
	// Find Breaker Blocks
	obs := FindEnhancedOrderBlocks(candles)
	breakers := FindBreakerBlocks(candles, obs)
	
	if len(breakers) == 0 {
		return nil
	}
	
	// Find FVGs
	fvgs := FindFairValueGaps(candles)
	
	// Check OTE zone
	inOTE := IsInOTE(candles)
	
	_ = candles[len(candles)-1].Close // currentPrice
	atr := calculateATR(candles[len(candles)-14:], 14)
	
	// Find recent breaker (use last one if available)
	var recentBreaker *BreakerBlock
	if len(breakers) > 0 {
		recentBreaker = &breakers[len(breakers)-1]
	}
	
	if recentBreaker == nil {
		return nil
	}
	
	unicorn := &UnicornSetup{
		Breaker:    recentBreaker,
		Confluence: 1, // Breaker
		Strength:   60,
		Valid:      false,
	}
	
	// Find FVG near breaker
	for i := range fvgs {
		fvg := &fvgs[i]
		if fvg.Type == recentBreaker.Type {
			dist := math.Abs(fvg.MidPoint - recentBreaker.MidPoint)
			if dist < atr*2 {
				unicorn.FVG = fvg
				unicorn.Confluence++
				unicorn.Strength += 15
				break
			}
		}
	}
	
	// OTE confluence
	if inOTE {
		unicorn.Confluence++
		unicorn.Strength += 15
		
		// Calculate OTE level (61.8% - 78.6%)
		high := candles[0].High
		low := candles[0].Low
		for _, c := range candles[:len(candles)/2] {
			if c.High > high {
				high = c.High
			}
			if c.Low < low {
				low = c.Low
			}
		}
		unicorn.OTELevel = low + (high-low)*0.618
	}
	
	// Validate setup
	if unicorn.Confluence >= 2 {
		unicorn.Valid = true
		unicorn.Type = recentBreaker.Type
		
		if unicorn.Type == "bullish" {
			unicorn.Entry = recentBreaker.MidPoint
			unicorn.StopLoss = recentBreaker.Low - atr*0.3
			unicorn.Target1 = recentBreaker.High + atr*2
			unicorn.Target2 = recentBreaker.High + atr*4
		} else {
			unicorn.Entry = recentBreaker.MidPoint
			unicorn.StopLoss = recentBreaker.High + atr*0.3
			unicorn.Target1 = recentBreaker.Low - atr*2
			unicorn.Target2 = recentBreaker.Low - atr*4
		}
	}
	
	return unicorn
}

// ==================== TURTLE SOUP ====================
// Turtle Soup: Fade failed breakouts

// TurtleSoup represents a Turtle Soup setup
type TurtleSoup struct {
	Type           string
	BreakoutLevel  float64
	FailurePrice   float64
	Entry          float64
	StopLoss       float64
	Target         float64
	DaysToBreakout int
	Strength       float64
	Valid          bool
}

// DetectTurtleSoup detects Turtle Soup setups
func DetectTurtleSoup(candles []Candle, lookback int) *TurtleSoup {
	if len(candles) < lookback+5 {
		return nil
	}
	
	// Find 20-day high and low
	periodCandles := candles[len(candles)-lookback-5 : len(candles)-5]
	
	periodHigh := periodCandles[0].High
	periodLow := periodCandles[0].Low
	
	for _, c := range periodCandles {
		if c.High > periodHigh {
			periodHigh = c.High
		}
		if c.Low < periodLow {
			periodLow = c.Low
		}
	}
	
	// Check recent candles for breakout and failure
	recentCandles := candles[len(candles)-5:]
	_ = candles[len(candles)-1].Close // currentPrice
	atr := calculateATR(candles[len(candles)-14:], 14)
	
	ts := &TurtleSoup{
		Valid: false,
	}
	
	// Bullish Turtle Soup: Break below low, then reverse up
	for i, c := range recentCandles {
		if c.Low < periodLow && c.Close > periodLow {
			// Failed breakdown
			ts.Type = "bullish"
			ts.BreakoutLevel = periodLow
			ts.FailurePrice = c.Low
			ts.Entry = periodLow + atr*0.2
			ts.StopLoss = c.Low - atr*0.5
			ts.Target = periodLow + atr*3
			ts.DaysToBreakout = i + 1
			ts.Strength = 75
			ts.Valid = true
			break
		}
	}
	
	// Bearish Turtle Soup: Break above high, then reverse down
	if !ts.Valid {
		for i, c := range recentCandles {
			if c.High > periodHigh && c.Close < periodHigh {
				// Failed breakout
				ts.Type = "bearish"
				ts.BreakoutLevel = periodHigh
				ts.FailurePrice = c.High
				ts.Entry = periodHigh - atr*0.2
				ts.StopLoss = c.High + atr*0.5
				ts.Target = periodHigh - atr*3
				ts.DaysToBreakout = i + 1
				ts.Strength = 75
				ts.Valid = true
				break
			}
		}
	}
	
	return ts
}

// ==================== QUASIMODO (QM) ====================
// Quasimodo: Over and Under pattern

// QuasimodoSetup represents a Quasimodo setup
type QuasimodoSetup struct {
	Type       string
	LeftHigh   float64
	Head       float64
	RightHigh  float64
	Neckline   float64
	Entry      float64
	StopLoss   float64
	Target     float64
	Strength   float64
	Valid      bool
}

// DetectQuasimodo detects Quasimodo patterns
func DetectQuasimodo(candles []Candle) *QuasimodoSetup {
	if len(candles) < 30 {
		return nil
	}
	
	// Find swing points
	type SwingPoint struct {
		Price float64
		Idx   int
		Type  string
	}
	
	var swings []SwingPoint
	
	for i := 2; i < len(candles)-2; i++ {
		if candles[i].High > candles[i-1].High &&
			candles[i].High > candles[i-2].High &&
			candles[i].High > candles[i+1].High &&
			candles[i].High > candles[i+2].High {
			swings = append(swings, SwingPoint{candles[i].High, i, "high"})
		}
		
		if candles[i].Low < candles[i-1].Low &&
			candles[i].Low < candles[i-2].Low &&
			candles[i].Low < candles[i+1].Low &&
			candles[i].Low < candles[i+2].Low {
			swings = append(swings, SwingPoint{candles[i].Low, i, "low"})
		}
	}
	
	if len(swings) < 5 {
		return nil
	}
	
	atr := calculateATR(candles[len(candles)-14:], 14)
	
	// Bullish QM: Higher Low after Lower Low
	// Pattern: High → Low → Higher High → Lower Low → Entry
	for i := 0; i < len(swings)-4; i++ {
		if swings[i].Type == "high" &&
			swings[i+1].Type == "low" &&
			swings[i+2].Type == "high" &&
			swings[i+3].Type == "low" {
			
			leftHigh := swings[i].Price
			leftLow := swings[i+1].Price
			head := swings[i+2].Price
			rightLow := swings[i+3].Price
			
			// QM conditions:
			// 1. Head > Left High (higher high)
			// 2. Right Low < Left Low (lower low)
			if head > leftHigh && rightLow < leftLow {
				qm := &QuasimodoSetup{
					Type:      "bullish",
					LeftHigh:  leftHigh,
					Head:      head,
					RightHigh: 0, // Will be formed
					Neckline:  leftLow,
					Entry:     leftLow,
					StopLoss:  rightLow - atr*0.3,
					Target:    leftLow + (head - rightLow),
					Strength:  80,
					Valid:     true,
				}
				return qm
			}
		}
	}
	
	// Bearish QM: Lower High after Higher High
	for i := 0; i < len(swings)-4; i++ {
		if swings[i].Type == "low" &&
			swings[i+1].Type == "high" &&
			swings[i+2].Type == "low" &&
			swings[i+3].Type == "high" {
			
			leftLow := swings[i].Price
			leftHigh := swings[i+1].Price
			head := swings[i+2].Price
			rightHigh := swings[i+3].Price
			
			// QM conditions:
			// 1. Head < Left Low (lower low)
			// 2. Right High > Left High (higher high)
			if head < leftLow && rightHigh > leftHigh {
				qm := &QuasimodoSetup{
					Type:      "bearish",
					LeftHigh:  leftHigh,
					Head:      head,
					RightHigh: rightHigh,
					Neckline:  leftHigh,
					Entry:     leftHigh,
					StopLoss:  rightHigh + atr*0.3,
					Target:    leftHigh - (rightHigh - head),
					Strength:  80,
					Valid:     true,
				}
				return qm
			}
		}
	}
	
	return nil
}

// ==================== INSTITUTIONAL CANDLE ====================
// Large range candles indicating institutional activity

// InstitutionalCandle represents an institutional candle
type InstitutionalCandle struct {
	Type       string
	Open       float64
	High       float64
	Low        float64
	Close      float64
	Range      float64
	Body       float64
	Volume     float64
	Strength   float64
	CandleIdx  int
}

// FindInstitutionalCandles finds large institutional candles
func FindInstitutionalCandles(candles []Candle) []InstitutionalCandle {
	var instCandles []InstitutionalCandle
	
	if len(candles) < 20 {
		return instCandles
	}
	
	// Calculate average range
	avgRange := 0.0
	avgVolume := 0.0
	for _, c := range candles {
		avgRange += c.High - c.Low
		avgVolume += c.Volume
	}
	avgRange /= float64(len(candles))
	avgVolume /= float64(len(candles))
	
	// Find candles with range > 2x average
	for i, c := range candles {
		candleRange := c.High - c.Low
		body := math.Abs(c.Close - c.Open)
		
		if candleRange > avgRange*2 && c.Volume > avgVolume*1.5 {
			ic := InstitutionalCandle{
				Open:      c.Open,
				High:      c.High,
				Low:       c.Low,
				Close:     c.Close,
				Range:     candleRange,
				Body:      body,
				Volume:    c.Volume,
				CandleIdx: i,
			}
			
			// Determine type
			if c.Close > c.Open {
				ic.Type = "bullish"
			} else {
				ic.Type = "bearish"
			}
			
			// Calculate strength
			ic.Strength = 50 + (candleRange/avgRange-1)*20 + (c.Volume/avgVolume-1)*15
			ic.Strength = math.Min(ic.Strength, 100)
			
			instCandles = append(instCandles, ic)
		}
	}
	
	return instCandles
}

// ==================== DISPLACEMENT ====================
// Strong directional move indicating institutional intent

// Displacement represents a displacement move
type Displacement struct {
	Type       string
	StartPrice float64
	EndPrice   float64
	Size       float64
	Candles    int
	Volume     float64
	Strength   float64
	StartIdx   int
	EndIdx     int
}

// FindDisplacements finds displacement moves
func FindDisplacements(candles []Candle) []Displacement {
	var displacements []Displacement
	
	if len(candles) < 10 {
		return displacements
	}
	
	// Calculate average move
	avgMove := 0.0
	for i := 1; i < len(candles); i++ {
		avgMove += math.Abs(candles[i].Close - candles[i-1].Close)
	}
	avgMove /= float64(len(candles) - 1)
	
	// Find consecutive strong moves in same direction
	for i := 0; i < len(candles)-3; i++ {
		// Check for 3+ consecutive bullish candles
		bullishCount := 0
		totalMove := 0.0
		totalVolume := 0.0
		
		for j := i; j < len(candles) && j < i+5; j++ {
			if candles[j].Close > candles[j].Open {
				bullishCount++
				totalMove += candles[j].Close - candles[j].Open
				totalVolume += candles[j].Volume
			} else {
				break
			}
		}
		
		if bullishCount >= 3 && totalMove > avgMove*3 {
			disp := Displacement{
				Type:       "bullish",
				StartPrice: candles[i].Open,
				EndPrice:   candles[i+bullishCount-1].Close,
				Size:       totalMove,
				Candles:    bullishCount,
				Volume:     totalVolume,
				StartIdx:   i,
				EndIdx:     i + bullishCount - 1,
			}
			disp.Strength = 50 + (totalMove/avgMove-1)*15 + float64(bullishCount)*5
			disp.Strength = math.Min(disp.Strength, 100)
			
			displacements = append(displacements, disp)
		}
		
		// Check for bearish displacement
		bearishCount := 0
		totalMove = 0
		totalVolume = 0
		
		for j := i; j < len(candles) && j < i+5; j++ {
			if candles[j].Close < candles[j].Open {
				bearishCount++
				totalMove += candles[j].Open - candles[j].Close
				totalVolume += candles[j].Volume
			} else {
				break
			}
		}
		
		if bearishCount >= 3 && totalMove > avgMove*3 {
			disp := Displacement{
				Type:       "bearish",
				StartPrice: candles[i].Open,
				EndPrice:   candles[i+bearishCount-1].Close,
				Size:       totalMove,
				Candles:    bearishCount,
				Volume:     totalVolume,
				StartIdx:   i,
				EndIdx:     i + bearishCount - 1,
			}
			disp.Strength = 50 + (totalMove/avgMove-1)*15 + float64(bearishCount)*5
			disp.Strength = math.Min(disp.Strength, 100)
			
			displacements = append(displacements, disp)
		}
	}
	
	return displacements
}

// ==================== FULL INSTITUTIONAL ANALYSIS ====================

// InstitutionalAnalysis holds all institutional analysis
type InstitutionalAnalysis struct {
	SilverBullet    *SilverBullet
	Unicorn         *UnicornSetup
	TurtleSoup      *TurtleSoup
	Quasimodo       *QuasimodoSetup
	InstCandles     []InstitutionalCandle
	Displacements   []Displacement
	OverallBias     string
	SetupCount      int
	TotalStrength   float64
}

// PerformInstitutionalAnalysis performs complete institutional analysis
func PerformInstitutionalAnalysis(candles []Candle) *InstitutionalAnalysis {
	analysis := &InstitutionalAnalysis{
		OverallBias: "neutral",
	}
	
	currentTime := time.Now()
	
	// Detect setups
	analysis.SilverBullet = DetectSilverBullet(candles, currentTime)
	analysis.Unicorn = DetectUnicornSetup(candles)
	analysis.TurtleSoup = DetectTurtleSoup(candles, 20)
	analysis.Quasimodo = DetectQuasimodo(candles)
	analysis.InstCandles = FindInstitutionalCandles(candles)
	analysis.Displacements = FindDisplacements(candles)
	
	// Count valid setups and calculate strength
	bullishStrength := 0.0
	bearishStrength := 0.0
	
	if analysis.SilverBullet != nil && analysis.SilverBullet.Valid {
		analysis.SetupCount++
		if analysis.SilverBullet.Type == "bullish" {
			bullishStrength += analysis.SilverBullet.Strength
		} else {
			bearishStrength += analysis.SilverBullet.Strength
		}
	}
	
	if analysis.Unicorn != nil && analysis.Unicorn.Valid {
		analysis.SetupCount++
		if analysis.Unicorn.Type == "bullish" {
			bullishStrength += analysis.Unicorn.Strength
		} else {
			bearishStrength += analysis.Unicorn.Strength
		}
	}
	
	if analysis.TurtleSoup != nil && analysis.TurtleSoup.Valid {
		analysis.SetupCount++
		if analysis.TurtleSoup.Type == "bullish" {
			bullishStrength += analysis.TurtleSoup.Strength
		} else {
			bearishStrength += analysis.TurtleSoup.Strength
		}
	}
	
	if analysis.Quasimodo != nil && analysis.Quasimodo.Valid {
		analysis.SetupCount++
		if analysis.Quasimodo.Type == "bullish" {
			bullishStrength += analysis.Quasimodo.Strength
		} else {
			bearishStrength += analysis.Quasimodo.Strength
		}
	}
	
	// Recent institutional candles
	for _, ic := range analysis.InstCandles {
		if ic.CandleIdx >= len(candles)-10 {
			if ic.Type == "bullish" {
				bullishStrength += ic.Strength * 0.3
			} else {
				bearishStrength += ic.Strength * 0.3
			}
		}
	}
	
	// Recent displacements
	for _, d := range analysis.Displacements {
		if d.EndIdx >= len(candles)-10 {
			if d.Type == "bullish" {
				bullishStrength += d.Strength * 0.4
			} else {
				bearishStrength += d.Strength * 0.4
			}
		}
	}
	
	// Determine overall bias
	if bullishStrength > bearishStrength+20 {
		analysis.OverallBias = "bullish"
		analysis.TotalStrength = bullishStrength
	} else if bearishStrength > bullishStrength+20 {
		analysis.OverallBias = "bearish"
		analysis.TotalStrength = bearishStrength
	} else {
		analysis.TotalStrength = (bullishStrength + bearishStrength) / 2
	}
	
	return analysis
}

// GetInstitutionalSignal returns trading signal based on institutional analysis
func GetInstitutionalSignal(analysis *InstitutionalAnalysis) (string, float64) {
	if analysis == nil {
		return "neutral", 0
	}
	
	return analysis.OverallBias, math.Min(analysis.TotalStrength, 100)
}
