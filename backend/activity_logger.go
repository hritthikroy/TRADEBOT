package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// ActivityType represents different types of activities
type ActivityType string

const (
	ActivitySignal     ActivityType = "SIGNAL"
	ActivityTrade      ActivityType = "TRADE"
	ActivityBacktest   ActivityType = "BACKTEST"
	ActivityOptimize   ActivityType = "OPTIMIZE"
	ActivityAI         ActivityType = "AI_ANALYSIS"
	ActivityTelegram   ActivityType = "TELEGRAM"
	ActivityWebSocket  ActivityType = "WEBSOCKET"
	ActivityAPI        ActivityType = "API_CALL"
	ActivityError      ActivityType = "ERROR"
	ActivitySystem     ActivityType = "SYSTEM"
	ActivityPaperTrade ActivityType = "PAPER_TRADE"
)

// ActivityLevel represents the severity/importance of an activity
type ActivityLevel string

const (
	LevelInfo    ActivityLevel = "INFO"
	LevelSuccess ActivityLevel = "SUCCESS"
	LevelWarning ActivityLevel = "WARNING"
	LevelError   ActivityLevel = "ERROR"
	LevelDebug   ActivityLevel = "DEBUG"
)

// Activity represents a single logged activity
type Activity struct {
	ID        string        `json:"id"`
	Timestamp time.Time     `json:"timestamp"`
	Type      ActivityType  `json:"type"`
	Level     ActivityLevel `json:"level"`
	Message   string        `json:"message"`
	Details   interface{}   `json:"details,omitempty"`
	Duration  string        `json:"duration,omitempty"`
	User      string        `json:"user,omitempty"`
	IP        string        `json:"ip,omitempty"`
}

// ActivityLogger manages activity logging
type ActivityLogger struct {
	activities []Activity
	mu         sync.RWMutex
	maxSize    int
	clients    map[string]chan Activity
	clientsMu  sync.RWMutex
}

var (
	activityLogger *ActivityLogger
	loggerOnce     sync.Once
)

// GetActivityLogger returns the singleton activity logger
func GetActivityLogger() *ActivityLogger {
	loggerOnce.Do(func() {
		activityLogger = &ActivityLogger{
			activities: make([]Activity, 0),
			maxSize:    1000, // Keep last 1000 activities
			clients:    make(map[string]chan Activity),
		}
	})
	return activityLogger
}

// Log adds a new activity to the log
func (al *ActivityLogger) Log(activityType ActivityType, level ActivityLevel, message string, details interface{}) {
	al.mu.Lock()
	defer al.mu.Unlock()

	activity := Activity{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Timestamp: time.Now(),
		Type:      activityType,
		Level:     level,
		Message:   message,
		Details:   details,
	}

	// Add to activities list
	al.activities = append(al.activities, activity)

	// Keep only last maxSize activities
	if len(al.activities) > al.maxSize {
		al.activities = al.activities[len(al.activities)-al.maxSize:]
	}

	// Broadcast to all connected clients
	al.broadcast(activity)

	// Also log to console with emoji
	al.logToConsole(activity)
}

// LogWithDuration logs an activity with duration
func (al *ActivityLogger) LogWithDuration(activityType ActivityType, level ActivityLevel, message string, duration time.Duration, details interface{}) {
	al.mu.Lock()
	defer al.mu.Unlock()

	activity := Activity{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Timestamp: time.Now(),
		Type:      activityType,
		Level:     level,
		Message:   message,
		Details:   details,
		Duration:  duration.String(),
	}

	al.activities = append(al.activities, activity)

	if len(al.activities) > al.maxSize {
		al.activities = al.activities[len(al.activities)-al.maxSize:]
	}

	al.broadcast(activity)
	al.logToConsole(activity)
}

// GetActivities returns recent activities
func (al *ActivityLogger) GetActivities(limit int) []Activity {
	al.mu.RLock()
	defer al.mu.RUnlock()

	if limit <= 0 || limit > len(al.activities) {
		limit = len(al.activities)
	}

	start := len(al.activities) - limit
	if start < 0 {
		start = 0
	}

	result := make([]Activity, limit)
	copy(result, al.activities[start:])

	// Reverse to show newest first
	for i := 0; i < len(result)/2; i++ {
		j := len(result) - i - 1
		result[i], result[j] = result[j], result[i]
	}

	return result
}

// GetActivitiesByType returns activities filtered by type
func (al *ActivityLogger) GetActivitiesByType(activityType ActivityType, limit int) []Activity {
	al.mu.RLock()
	defer al.mu.RUnlock()

	filtered := make([]Activity, 0)
	for i := len(al.activities) - 1; i >= 0 && len(filtered) < limit; i-- {
		if al.activities[i].Type == activityType {
			filtered = append(filtered, al.activities[i])
		}
	}

	return filtered
}

// GetActivitiesByLevel returns activities filtered by level
func (al *ActivityLogger) GetActivitiesByLevel(level ActivityLevel, limit int) []Activity {
	al.mu.RLock()
	defer al.mu.RUnlock()

	filtered := make([]Activity, 0)
	for i := len(al.activities) - 1; i >= 0 && len(filtered) < limit; i-- {
		if al.activities[i].Level == level {
			filtered = append(filtered, al.activities[i])
		}
	}

	return filtered
}

// GetStats returns activity statistics
func (al *ActivityLogger) GetStats() map[string]interface{} {
	al.mu.RLock()
	defer al.mu.RUnlock()

	stats := map[string]interface{}{
		"total_activities": len(al.activities),
		"by_type":          make(map[ActivityType]int),
		"by_level":         make(map[ActivityLevel]int),
		"connected_clients": len(al.clients),
	}

	byType := make(map[ActivityType]int)
	byLevel := make(map[ActivityLevel]int)

	for _, activity := range al.activities {
		byType[activity.Type]++
		byLevel[activity.Level]++
	}

	stats["by_type"] = byType
	stats["by_level"] = byLevel

	return stats
}

// Subscribe adds a client to receive real-time activity updates
func (al *ActivityLogger) Subscribe(clientID string) chan Activity {
	al.clientsMu.Lock()
	defer al.clientsMu.Unlock()

	ch := make(chan Activity, 100)
	al.clients[clientID] = ch
	return ch
}

// Unsubscribe removes a client from receiving updates
func (al *ActivityLogger) Unsubscribe(clientID string) {
	al.clientsMu.Lock()
	defer al.clientsMu.Unlock()

	if ch, exists := al.clients[clientID]; exists {
		close(ch)
		delete(al.clients, clientID)
	}
}

// broadcast sends activity to all subscribed clients
func (al *ActivityLogger) broadcast(activity Activity) {
	al.clientsMu.RLock()
	defer al.clientsMu.RUnlock()

	for clientID, ch := range al.clients {
		select {
		case ch <- activity:
		default:
			// Channel full, skip this client
			fmt.Printf("⚠️  Activity channel full for client %s\n", clientID)
		}
	}
}

// logToConsole prints activity to console with formatting
func (al *ActivityLogger) logToConsole(activity Activity) {
	emoji := al.getEmoji(activity.Type, activity.Level)
	timestamp := activity.Timestamp.Format("15:04:05")
	
	detailsStr := ""
	if activity.Details != nil {
		if bytes, err := json.Marshal(activity.Details); err == nil {
			detailsStr = string(bytes)
			if len(detailsStr) > 100 {
				detailsStr = detailsStr[:100] + "..."
			}
		}
	}

	durationStr := ""
	if activity.Duration != "" {
		durationStr = fmt.Sprintf(" [%s]", activity.Duration)
	}

	fmt.Printf("%s [%s] %s: %s%s %s\n", 
		emoji, 
		timestamp, 
		activity.Type, 
		activity.Message,
		durationStr,
		detailsStr,
	)
}

// getEmoji returns appropriate emoji for activity type and level
func (al *ActivityLogger) getEmoji(activityType ActivityType, level ActivityLevel) string {
	if level == LevelError {
		return "❌"
	}
	if level == LevelWarning {
		return "⚠️"
	}
	if level == LevelSuccess {
		return "✅"
	}

	switch activityType {
	case ActivitySignal:
		return "📊"
	case ActivityTrade:
		return "💰"
	case ActivityBacktest:
		return "🧪"
	case ActivityOptimize:
		return "⚙️"
	case ActivityAI:
		return "🤖"
	case ActivityTelegram:
		return "📱"
	case ActivityWebSocket:
		return "🔌"
	case ActivityAPI:
		return "🌐"
	case ActivityPaperTrade:
		return "📝"
	case ActivitySystem:
		return "🖥️"
	default:
		return "ℹ️"
	}
}

// Clear removes all activities
func (al *ActivityLogger) Clear() {
	al.mu.Lock()
	defer al.mu.Unlock()

	al.activities = make([]Activity, 0)
	fmt.Println("🗑️  Activity log cleared")
}

// Helper functions for common logging patterns

func LogSignal(message string, details interface{}) {
	GetActivityLogger().Log(ActivitySignal, LevelInfo, message, details)
}

func LogSignalSuccess(message string, details interface{}) {
	GetActivityLogger().Log(ActivitySignal, LevelSuccess, message, details)
}

func LogTrade(message string, details interface{}) {
	GetActivityLogger().Log(ActivityTrade, LevelInfo, message, details)
}

func LogTradeSuccess(message string, details interface{}) {
	GetActivityLogger().Log(ActivityTrade, LevelSuccess, message, details)
}

func LogBacktest(message string, details interface{}) {
	GetActivityLogger().Log(ActivityBacktest, LevelInfo, message, details)
}

func LogBacktestWithDuration(message string, duration time.Duration, details interface{}) {
	GetActivityLogger().LogWithDuration(ActivityBacktest, LevelSuccess, message, duration, details)
}

func LogAI(message string, details interface{}) {
	GetActivityLogger().Log(ActivityAI, LevelInfo, message, details)
}

func LogAISuccess(message string, details interface{}) {
	GetActivityLogger().Log(ActivityAI, LevelSuccess, message, details)
}

func LogError(activityType ActivityType, message string, details interface{}) {
	GetActivityLogger().Log(activityType, LevelError, message, details)
}

func LogWarning(activityType ActivityType, message string, details interface{}) {
	GetActivityLogger().Log(activityType, LevelWarning, message, details)
}

func LogSystem(message string, details interface{}) {
	GetActivityLogger().Log(ActivitySystem, LevelInfo, message, details)
}

func LogSystemSuccess(message string, details interface{}) {
	GetActivityLogger().Log(ActivitySystem, LevelSuccess, message, details)
}

func LogAPI(message string, details interface{}) {
	GetActivityLogger().Log(ActivityAPI, LevelInfo, message, details)
}
