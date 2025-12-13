package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// GrokAIService handles Grok AI API interactions
type GrokAIService struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

// GrokRequest represents a request to Grok AI
type GrokRequest struct {
	Messages []GrokMessage `json:"messages"`
	Model    string        `json:"model"`
	Stream   bool          `json:"stream"`
	Temperature float64    `json:"temperature"`
}

// GrokMessage represents a message in the conversation
type GrokMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// GrokResponse represents the response from Grok AI
type GrokResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// MarketSentiment represents AI-analyzed market sentiment
type MarketSentiment struct {
	Overall           string  `json:"overall"`            // "bullish", "bearish", "neutral"
	Confidence        float64 `json:"confidence"`         // 0-100
	NewsImpact        string  `json:"news_impact"`        // "positive", "negative", "neutral"
	RiskLevel         string  `json:"risk_level"`         // "low", "medium", "high"
	TradingAdvice     string  `json:"trading_advice"`     // "take_trade", "avoid", "reduce_size"
	KeyFactors        []string `json:"key_factors"`
	ShouldTrade       bool    `json:"should_trade"`
	RecommendedRisk   float64 `json:"recommended_risk"`   // 0.5-2.0 multiplier
}

// NewGrokAIService creates a new Grok AI service
func NewGrokAIService() *GrokAIService {
	apiKey := os.Getenv("GROK_API_KEY")
	if apiKey == "" {
		log.Println("⚠️ GROK_API_KEY not set in environment")
	}

	return &GrokAIService{
		APIKey:  apiKey,
		BaseURL: "https://api.x.ai/v1/chat/completions",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CallGrokAPI makes a request to Grok AI
func (g *GrokAIService) CallGrokAPI(prompt string) (string, error) {
	if g.APIKey == "" {
		return "", fmt.Errorf("Grok API key not configured")
	}

	reqBody := GrokRequest{
		Messages: []GrokMessage{
			{
				Role:    "system",
				Content: "You are a professional cryptocurrency trading analyst. Provide concise, actionable market analysis.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Model:       "grok-beta",
		Stream:      false,
		Temperature: 0.3, // Lower temperature for more consistent analysis
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", g.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.APIKey)

	resp, err := g.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call Grok API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Grok API error (status %d): %s", resp.StatusCode, string(body))
	}

	var grokResp GrokResponse
	if err := json.Unmarshal(body, &grokResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(grokResp.Choices) == 0 {
		return "", fmt.Errorf("no response from Grok AI")
	}

	return grokResp.Choices[0].Message.Content, nil
}

// AnalyzeMarketSentiment analyzes current market sentiment for a trading signal
func (g *GrokAIService) AnalyzeMarketSentiment(symbol string, signalType string, currentPrice float64, strength int) (*MarketSentiment, error) {
	prompt := fmt.Sprintf(`Analyze the current market sentiment for %s trading at $%.2f.

Signal Details:
- Type: %s
- Technical Strength: %d%%
- Current Price: $%.2f

Provide a brief analysis covering:
1. Overall market sentiment (bullish/bearish/neutral)
2. Recent news impact
3. Risk level for this trade
4. Trading recommendation

Format your response as JSON with these fields:
{
  "overall": "bullish|bearish|neutral",
  "confidence": 0-100,
  "news_impact": "positive|negative|neutral",
  "risk_level": "low|medium|high",
  "trading_advice": "take_trade|avoid|reduce_size",
  "key_factors": ["factor1", "factor2"],
  "should_trade": true|false,
  "recommended_risk": 0.5-2.0
}

Keep it concise and actionable.`, symbol, currentPrice, signalType, strength, currentPrice)

	response, err := g.CallGrokAPI(prompt)
	if err != nil {
		return nil, err
	}

	// Parse JSON response
	sentiment := &MarketSentiment{}
	
	// Try to extract JSON from response
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}")
	
	if jsonStart == -1 || jsonEnd == -1 {
		// Fallback: parse text response
		return g.parseTextSentiment(response), nil
	}

	jsonStr := response[jsonStart : jsonEnd+1]
	if err := json.Unmarshal([]byte(jsonStr), sentiment); err != nil {
		// Fallback: parse text response
		return g.parseTextSentiment(response), nil
	}

	return sentiment, nil
}

// parseTextSentiment parses text response when JSON parsing fails
func (g *GrokAIService) parseTextSentiment(text string) *MarketSentiment {
	text = strings.ToLower(text)
	
	sentiment := &MarketSentiment{
		Overall:         "neutral",
		Confidence:      50.0,
		NewsImpact:      "neutral",
		RiskLevel:       "medium",
		TradingAdvice:   "take_trade",
		KeyFactors:      []string{"AI analysis available"},
		ShouldTrade:     true,
		RecommendedRisk: 1.0,
	}

	// Detect sentiment
	if strings.Contains(text, "bullish") || strings.Contains(text, "positive") || strings.Contains(text, "uptrend") {
		sentiment.Overall = "bullish"
		sentiment.Confidence = 70.0
	} else if strings.Contains(text, "bearish") || strings.Contains(text, "negative") || strings.Contains(text, "downtrend") {
		sentiment.Overall = "bearish"
		sentiment.Confidence = 70.0
	}

	// Detect risk
	if strings.Contains(text, "high risk") || strings.Contains(text, "volatile") || strings.Contains(text, "caution") {
		sentiment.RiskLevel = "high"
		sentiment.RecommendedRisk = 0.5
	} else if strings.Contains(text, "low risk") || strings.Contains(text, "stable") {
		sentiment.RiskLevel = "low"
		sentiment.RecommendedRisk = 1.5
	}

	// Detect trading advice
	if strings.Contains(text, "avoid") || strings.Contains(text, "don't trade") || strings.Contains(text, "wait") {
		sentiment.TradingAdvice = "avoid"
		sentiment.ShouldTrade = false
	} else if strings.Contains(text, "reduce") || strings.Contains(text, "smaller position") {
		sentiment.TradingAdvice = "reduce_size"
		sentiment.RecommendedRisk = 0.7
	}

	return sentiment
}

// GetQuickSentiment provides fast sentiment check (cached or simplified)
func (g *GrokAIService) GetQuickSentiment(symbol string) string {
	// Quick sentiment without full AI call (for performance)
	// In production, you might cache this or use a simpler endpoint
	return "neutral"
}

// ValidateSignalWithAI validates a trading signal using AI analysis
func (g *GrokAIService) ValidateSignalWithAI(signal *CreateSignalRequest) (bool, string, float64) {
	sentiment, err := g.AnalyzeMarketSentiment(
		signal.Symbol,
		signal.SignalType,
		signal.EntryPrice,
		signal.Strength,
	)

	if err != nil {
		log.Printf("⚠️ AI validation failed: %v (allowing trade)", err)
		return true, "AI unavailable - proceeding with technical analysis", 1.0
	}

	log.Printf("🤖 AI Sentiment: %s (Confidence: %.0f%%, Risk: %s)",
		sentiment.Overall,
		sentiment.Confidence,
		sentiment.RiskLevel,
	)

	// Decision logic
	if !sentiment.ShouldTrade {
		return false, fmt.Sprintf("AI recommends avoiding trade: %s", sentiment.TradingAdvice), 0.0
	}

	// Check sentiment alignment
	if signal.SignalType == "BUY" && sentiment.Overall == "bearish" && sentiment.Confidence > 70 {
		return false, "AI sentiment conflicts with BUY signal", 0.0
	}

	if signal.SignalType == "SELL" && sentiment.Overall == "bullish" && sentiment.Confidence > 70 {
		return false, "AI sentiment conflicts with SELL signal", 0.0
	}

	// High risk warning
	if sentiment.RiskLevel == "high" {
		log.Printf("⚠️ High risk detected - reducing position size to %.1fx", sentiment.RecommendedRisk)
	}

	reason := fmt.Sprintf("AI validated: %s sentiment, %s risk", sentiment.Overall, sentiment.RiskLevel)
	return true, reason, sentiment.RecommendedRisk
}
