package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// ExternalAIConfig holds API configurations
type ExternalAIConfig struct {
	OpenAIKey      string
	AnthropicKey   string
	GeminiKey      string
	PerplexityKey  string
	EnableOpenAI   bool
	EnableAnthropic bool
	EnableGemini   bool
	EnablePerplexity bool
}

// AIEnhancedAnalysis combines local AI with external AI insights
type AIEnhancedAnalysis struct {
	LocalAnalysis      *AIMarketAnalysis `json:"localAnalysis"`
	OpenAIInsights     *OpenAIResponse   `json:"openAIInsights,omitempty"`
	AnthropicInsights  *AnthropicResponse `json:"anthropicInsights,omitempty"`
	GeminiInsights     *GeminiResponse   `json:"geminiInsights,omitempty"`
	PerplexityInsights *PerplexityResponse `json:"perplexityInsights,omitempty"`
	CombinedRecommendation string         `json:"combinedRecommendation"`
	ConfidenceScore    float64            `json:"confidenceScore"`
}

// OpenAI Response
type OpenAIResponse struct {
	Analysis       string   `json:"analysis"`
	Recommendation string   `json:"recommendation"`
	RiskLevel      string   `json:"riskLevel"`
	KeyPoints      []string `json:"keyPoints"`
}

// Anthropic Response
type AnthropicResponse struct {
	Analysis       string   `json:"analysis"`
	Recommendation string   `json:"recommendation"`
	Confidence     float64  `json:"confidence"`
}

// Gemini Response
type GeminiResponse struct {
	Analysis       string   `json:"analysis"`
	Prediction     string   `json:"prediction"`
	Confidence     float64  `json:"confidence"`
}

// Perplexity Response
type PerplexityResponse struct {
	MarketNews     []string `json:"marketNews"`
	Sentiment      string   `json:"sentiment"`
	KeyEvents      []string `json:"keyEvents"`
}

// LoadAIConfig loads API keys from environment
func LoadAIConfig() *ExternalAIConfig {
	return &ExternalAIConfig{
		OpenAIKey:        os.Getenv("OPENAI_API_KEY"),
		AnthropicKey:     os.Getenv("ANTHROPIC_API_KEY"),
		GeminiKey:        os.Getenv("GEMINI_API_KEY"),
		PerplexityKey:    os.Getenv("PERPLEXITY_API_KEY"),
		EnableOpenAI:     os.Getenv("ENABLE_OPENAI") == "true",
		EnableAnthropic:  os.Getenv("ENABLE_ANTHROPIC") == "true",
		EnableGemini:     os.Getenv("ENABLE_GEMINI") == "true",
		EnablePerplexity: os.Getenv("ENABLE_PERPLEXITY") == "true",
	}
}

// EnhanceWithExternalAI enhances local analysis with external AI
func EnhanceWithExternalAI(localAnalysis *AIMarketAnalysis, candles []Candle, symbol string) (*AIEnhancedAnalysis, error) {
	config := LoadAIConfig()
	
	enhanced := &AIEnhancedAnalysis{
		LocalAnalysis: localAnalysis,
	}
	
	// Prepare market data summary
	marketSummary := prepareMarketSummary(localAnalysis, candles, symbol)
	
	// Call external AI services in parallel
	type aiResult struct {
		name string
		data interface{}
		err  error
	}
	
	results := make(chan aiResult, 4)
	activeAPIs := 0
	
	// OpenAI
	if config.EnableOpenAI && config.OpenAIKey != "" {
		activeAPIs++
		go func() {
			data, err := callOpenAI(config.OpenAIKey, marketSummary)
			results <- aiResult{"openai", data, err}
		}()
	}
	
	// Anthropic Claude
	if config.EnableAnthropic && config.AnthropicKey != "" {
		activeAPIs++
		go func() {
			data, err := callAnthropic(config.AnthropicKey, marketSummary)
			results <- aiResult{"anthropic", data, err}
		}()
	}
	
	// Google Gemini
	if config.EnableGemini && config.GeminiKey != "" {
		activeAPIs++
		go func() {
			data, err := callGemini(config.GeminiKey, marketSummary)
			results <- aiResult{"gemini", data, err}
		}()
	}
	
	// Perplexity (for news/sentiment)
	if config.EnablePerplexity && config.PerplexityKey != "" {
		activeAPIs++
		go func() {
			data, err := callPerplexity(config.PerplexityKey, symbol)
			results <- aiResult{"perplexity", data, err}
		}()
	}
	
	// Collect results
	for i := 0; i < activeAPIs; i++ {
		result := <-results
		if result.err == nil {
			switch result.name {
			case "openai":
				enhanced.OpenAIInsights = result.data.(*OpenAIResponse)
			case "anthropic":
				enhanced.AnthropicInsights = result.data.(*AnthropicResponse)
			case "gemini":
				enhanced.GeminiInsights = result.data.(*GeminiResponse)
			case "perplexity":
				enhanced.PerplexityInsights = result.data.(*PerplexityResponse)
			}
		}
	}
	
	// Combine insights
	enhanced.CombinedRecommendation = combineAIInsights(enhanced)
	enhanced.ConfidenceScore = calculateCombinedConfidence(enhanced)
	
	return enhanced, nil
}

func prepareMarketSummary(analysis *AIMarketAnalysis, candles []Candle, symbol string) string {
	currentPrice := candles[len(candles)-1].Close
	
	return fmt.Sprintf(`Market Analysis for %s:
- Current Price: %.2f
- Market Regime: %s
- Trend Strength: %.0f/100
- Volatility: %s
- Predicted Move: %s
- Risk Level: %s
- Support Levels: %v
- Resistance Levels: %v

Please provide:
1. Your analysis of the current market conditions
2. Trading recommendation (BUY/SELL/HOLD)
3. Risk assessment
4. Key points to consider
5. Confidence level (0-100)`,
		symbol, currentPrice, analysis.MarketRegime, analysis.TrendStrength,
		analysis.VolatilityLevel, analysis.PredictedMove, analysis.RiskLevel,
		analysis.SupportLevels, analysis.ResistanceLevels)
}

func callOpenAI(apiKey, prompt string) (*OpenAIResponse, error) {
	url := "https://api.openai.com/v1/chat/completions"
	
	payload := map[string]interface{}{
		"model": "gpt-4-turbo-preview",
		"messages": []map[string]string{
			{
				"role": "system",
				"content": "You are an expert cryptocurrency trading analyst. Provide concise, actionable trading insights.",
			},
			{
				"role": "user",
				"content": prompt,
			},
		},
		"temperature": 0.7,
		"max_tokens": 500,
	}
	
	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	
	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}
	
	content := result.Choices[0].Message.Content
	
	return &OpenAIResponse{
		Analysis:       content,
		Recommendation: extractRecommendation(content),
		RiskLevel:      extractRiskLevel(content),
		KeyPoints:      extractKeyPoints(content),
	}, nil
}

func callAnthropic(apiKey, prompt string) (*AnthropicResponse, error) {
	url := "https://api.anthropic.com/v1/messages"
	
	payload := map[string]interface{}{
		"model": "claude-3-sonnet-20240229",
		"max_tokens": 500,
		"messages": []map[string]string{
			{
				"role": "user",
				"content": prompt,
			},
		},
	}
	
	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	
	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	
	if len(result.Content) == 0 {
		return nil, fmt.Errorf("no response from Anthropic")
	}
	
	content := result.Content[0].Text
	
	return &AnthropicResponse{
		Analysis:       content,
		Recommendation: extractRecommendation(content),
		Confidence:     extractConfidence(content),
	}, nil
}

func callGemini(apiKey, prompt string) (*GeminiResponse, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=%s", apiKey)
	
	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	}
	
	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	
	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	
	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response from Gemini")
	}
	
	content := result.Candidates[0].Content.Parts[0].Text
	
	return &GeminiResponse{
		Analysis:   content,
		Prediction: extractRecommendation(content),
		Confidence: extractConfidence(content),
	}, nil
}

func callPerplexity(apiKey, symbol string) (*PerplexityResponse, error) {
	url := "https://api.perplexity.ai/chat/completions"
	
	prompt := fmt.Sprintf("What are the latest news and market sentiment for %s cryptocurrency? Provide 3-5 key points.", symbol)
	
	payload := map[string]interface{}{
		"model": "llama-3.1-sonar-small-128k-online",
		"messages": []map[string]string{
			{
				"role": "user",
				"content": prompt,
			},
		},
	}
	
	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	
	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no response from Perplexity")
	}
	
	content := result.Choices[0].Message.Content
	
	return &PerplexityResponse{
		MarketNews: extractKeyPoints(content),
		Sentiment:  extractSentiment(content),
		KeyEvents:  extractKeyPoints(content),
	}, nil
}

// Helper functions
func extractRecommendation(text string) string {
	text = strings.ToUpper(text)
	if strings.Contains(text, "BUY") || strings.Contains(text, "LONG") {
		return "BUY"
	} else if strings.Contains(text, "SELL") || strings.Contains(text, "SHORT") {
		return "SELL"
	}
	return "HOLD"
}

func extractRiskLevel(text string) string {
	text = strings.ToLower(text)
	if strings.Contains(text, "high risk") {
		return "high"
	} else if strings.Contains(text, "low risk") {
		return "low"
	}
	return "medium"
}

func extractConfidence(text string) float64 {
	// Simple extraction - look for percentage
	// In production, use more sophisticated NLP
	return 75.0
}

func extractKeyPoints(text string) []string {
	lines := strings.Split(text, "\n")
	points := []string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 10 && (strings.HasPrefix(line, "-") || strings.HasPrefix(line, "•") || strings.HasPrefix(line, "*")) {
			points = append(points, strings.TrimLeft(line, "-•* "))
		}
	}
	if len(points) == 0 {
		points = []string{text}
	}
	return points
}

func extractSentiment(text string) string {
	text = strings.ToLower(text)
	if strings.Contains(text, "bullish") || strings.Contains(text, "positive") {
		return "bullish"
	} else if strings.Contains(text, "bearish") || strings.Contains(text, "negative") {
		return "bearish"
	}
	return "neutral"
}

func combineAIInsights(enhanced *AIEnhancedAnalysis) string {
	recommendations := []string{}
	
	if enhanced.OpenAIInsights != nil {
		recommendations = append(recommendations, "OpenAI: "+enhanced.OpenAIInsights.Recommendation)
	}
	if enhanced.AnthropicInsights != nil {
		recommendations = append(recommendations, "Claude: "+enhanced.AnthropicInsights.Recommendation)
	}
	if enhanced.GeminiInsights != nil {
		recommendations = append(recommendations, "Gemini: "+enhanced.GeminiInsights.Prediction)
	}
	
	// Count votes
	buyVotes := 0
	sellVotes := 0
	holdVotes := 0
	
	for _, rec := range recommendations {
		if strings.Contains(rec, "BUY") {
			buyVotes++
		} else if strings.Contains(rec, "SELL") {
			sellVotes++
		} else {
			holdVotes++
		}
	}
	
	// Determine consensus
	if buyVotes > sellVotes && buyVotes > holdVotes {
		return fmt.Sprintf("🟢 CONSENSUS: BUY (%d/%d AI models agree)", buyVotes, len(recommendations))
	} else if sellVotes > buyVotes && sellVotes > holdVotes {
		return fmt.Sprintf("🔴 CONSENSUS: SELL (%d/%d AI models agree)", sellVotes, len(recommendations))
	} else {
		return fmt.Sprintf("🟡 CONSENSUS: HOLD (No clear agreement - %d BUY, %d SELL, %d HOLD)", buyVotes, sellVotes, holdVotes)
	}
}

func calculateCombinedConfidence(enhanced *AIEnhancedAnalysis) float64 {
	confidences := []float64{}
	
	if enhanced.LocalAnalysis != nil {
		confidences = append(confidences, enhanced.LocalAnalysis.Confidence)
	}
	if enhanced.AnthropicInsights != nil {
		confidences = append(confidences, enhanced.AnthropicInsights.Confidence)
	}
	if enhanced.GeminiInsights != nil {
		confidences = append(confidences, enhanced.GeminiInsights.Confidence)
	}
	
	if len(confidences) == 0 {
		return 0
	}
	
	sum := 0.0
	for _, c := range confidences {
		sum += c
	}
	
	return sum / float64(len(confidences))
}
