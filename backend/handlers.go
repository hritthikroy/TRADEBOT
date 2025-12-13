package main

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
)

// CreateSignal creates a new trading signal
func CreateSignal(c *fiber.Ctx) error {
	req := new(CreateSignalRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	query := `
		INSERT INTO trading_signals (
			signal_id, signal_type, symbol, entry_price, stop_loss, 
			tp1, tp2, tp3, strength, pattern_type, pattern_confidence,
			kill_zone, session_type, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, 'pending')
		RETURNING id, created_at
	`

	var id int64
	var createdAt time.Time
	err := DB.QueryRow(
		query, req.SignalID, req.SignalType, req.Symbol, req.EntryPrice,
		req.StopLoss, req.TP1, req.TP2, req.TP3, req.Strength,
		req.PatternType, req.PatternConfidence, req.KillZone, req.SessionType,
	).Scan(&id, &createdAt)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create signal"})
	}

	return c.Status(201).JSON(fiber.Map{
		"id":         id,
		"signal_id":  req.SignalID,
		"created_at": createdAt,
		"message":    "Signal created successfully",
	})
}

// GetAllSignals retrieves all signals
func GetAllSignals(c *fiber.Ctx) error {
	if DB == nil {
		return c.JSON([]TradingSignal{})
	}

	query := `
		SELECT * FROM trading_signals 
		ORDER BY created_at DESC
	`

	rows, err := DB.Query(query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch signals"})
	}
	defer rows.Close()

	signals := []TradingSignal{}
	for rows.Next() {
		var signal TradingSignal
		err := rows.Scan(
			&signal.ID, &signal.SignalID, &signal.CreatedAt, &signal.SignalType,
			&signal.Symbol, &signal.EntryPrice, &signal.StopLoss, &signal.TP1,
			&signal.TP2, &signal.TP3, &signal.Strength, &signal.PatternType,
			&signal.PatternConfidence, &signal.KillZone, &signal.SessionType,
			&signal.Status, &signal.ExitPrice, &signal.ExitReason, &signal.ExitTime,
			&signal.ProfitPercent, &signal.ProfitPips, &signal.HoldingTimeMinutes,
			&signal.TrailingStopPrice, &signal.TrailingStopActive, &signal.HighestPrice,
			&signal.LowestPrice, &signal.LivePrice, &signal.LastUpdated,
		)
		if err != nil {
			continue
		}
		signals = append(signals, signal)
	}

	return c.JSON(signals)
}

// GetPendingSignals retrieves only pending signals
func GetPendingSignals(c *fiber.Ctx) error {
	if DB == nil {
		return c.JSON([]TradingSignal{})
	}

	query := `
		SELECT * FROM trading_signals 
		WHERE status = 'pending'
		ORDER BY created_at DESC
	`

	rows, err := DB.Query(query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch pending signals"})
	}
	defer rows.Close()

	signals := []TradingSignal{}
	for rows.Next() {
		var signal TradingSignal
		err := rows.Scan(
			&signal.ID, &signal.SignalID, &signal.CreatedAt, &signal.SignalType,
			&signal.Symbol, &signal.EntryPrice, &signal.StopLoss, &signal.TP1,
			&signal.TP2, &signal.TP3, &signal.Strength, &signal.PatternType,
			&signal.PatternConfidence, &signal.KillZone, &signal.SessionType,
			&signal.Status, &signal.ExitPrice, &signal.ExitReason, &signal.ExitTime,
			&signal.ProfitPercent, &signal.ProfitPips, &signal.HoldingTimeMinutes,
			&signal.TrailingStopPrice, &signal.TrailingStopActive, &signal.HighestPrice,
			&signal.LowestPrice, &signal.LivePrice, &signal.LastUpdated,
		)
		if err != nil {
			continue
		}
		signals = append(signals, signal)
	}

	return c.JSON(signals)
}

// GetSignalByID retrieves a single signal by ID
func GetSignalByID(c *fiber.Ctx) error {
	id := c.Params("id")

	query := `SELECT * FROM trading_signals WHERE signal_id = $1`

	var signal TradingSignal
	err := DB.QueryRow(query, id).Scan(
		&signal.ID, &signal.SignalID, &signal.CreatedAt, &signal.SignalType,
		&signal.Symbol, &signal.EntryPrice, &signal.StopLoss, &signal.TP1,
		&signal.TP2, &signal.TP3, &signal.Strength, &signal.PatternType,
		&signal.PatternConfidence, &signal.KillZone, &signal.SessionType,
		&signal.Status, &signal.ExitPrice, &signal.ExitReason, &signal.ExitTime,
		&signal.ProfitPercent, &signal.ProfitPips, &signal.HoldingTimeMinutes,
		&signal.TrailingStopPrice, &signal.TrailingStopActive, &signal.HighestPrice,
		&signal.LowestPrice, &signal.LivePrice, &signal.LastUpdated,
	)

	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{"error": "Signal not found"})
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch signal"})
	}

	return c.JSON(signal)
}

// UpdateSignal updates a signal's status
func UpdateSignal(c *fiber.Ctx) error {
	id := c.Params("id")
	req := new(UpdateSignalRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	query := `
		UPDATE trading_signals 
		SET status = $1, exit_price = $2, exit_reason = $3, exit_time = NOW(),
		    profit_percent = $4, profit_pips = $5, holding_time_minutes = $6,
		    trailing_stop_price = $7
		WHERE signal_id = $8
		RETURNING id
	`

	var signalID int64
	err := DB.QueryRow(
		query, req.Status, req.ExitPrice, req.ExitReason, req.ProfitPercent,
		req.ProfitPips, req.HoldingTimeMinutes, req.TrailingStopPrice, id,
	).Scan(&signalID)

	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{"error": "Signal not found"})
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update signal"})
	}

	return c.JSON(fiber.Map{"message": "Signal updated successfully"})
}

// UpdateLivePrice updates the live price of a pending signal
func UpdateLivePrice(c *fiber.Ctx) error {
	id := c.Params("id")
	req := new(UpdateLivePriceRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	query := `
		UPDATE trading_signals 
		SET live_price = $1, trailing_stop_price = $2, 
		    trailing_stop_active = CASE WHEN $2 IS NOT NULL THEN true ELSE trailing_stop_active END
		WHERE signal_id = $3 AND status = 'pending'
		RETURNING id
	`

	var signalID int64
	err := DB.QueryRow(query, req.LivePrice, req.TrailingStopPrice, id).Scan(&signalID)

	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{"error": "Signal not found or not pending"})
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update live price"})
	}

	return c.JSON(fiber.Map{"message": "Live price updated successfully"})
}

// DeleteSignal deletes a signal
func DeleteSignal(c *fiber.Ctx) error {
	id := c.Params("id")

	query := `DELETE FROM trading_signals WHERE signal_id = $1`

	result, err := DB.Exec(query, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete signal"})
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Signal not found"})
	}

	return c.JSON(fiber.Map{"message": "Signal deleted successfully"})
}

// GetAnalytics retrieves analytics data
func GetAnalytics(c *fiber.Ctx) error {
	if DB == nil {
		return c.JSON([]SignalAnalytics{})
	}

	query := `SELECT * FROM signal_analytics`

	rows, err := DB.Query(query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch analytics"})
	}
	defer rows.Close()

	analytics := []SignalAnalytics{}
	for rows.Next() {
		var a SignalAnalytics
		err := rows.Scan(
			&a.KillZone, &a.SignalType, &a.PatternType, &a.TotalSignals,
			&a.Wins, &a.Losses, &a.AvgProfitPercent, &a.AvgWinPercent,
			&a.AvgLossPercent, &a.AvgHoldingMinutes, &a.WinRate,
		)
		if err != nil {
			continue
		}
		analytics = append(analytics, a)
	}

	return c.JSON(analytics)
}

// GetPerformanceStats retrieves overall performance statistics
func GetPerformanceStats(c *fiber.Ctx) error {
	// Check if DB is available
	if DB == nil {
		return c.JSON(fiber.Map{
			"total_signals": 0,
			"wins":          0,
			"losses":        0,
			"pending":       0,
			"avg_profit":    nil,
			"win_rate":      nil,
		})
	}

	query := `
		SELECT 
			COUNT(*) as total_signals,
			SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) as wins,
			SUM(CASE WHEN status = 'loss' THEN 1 ELSE 0 END) as losses,
			SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END) as pending,
			ROUND(AVG(CASE WHEN status IN ('win', 'loss') THEN profit_percent END), 2) as avg_profit,
			ROUND(100.0 * SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) / 
				NULLIF(SUM(CASE WHEN status IN ('win', 'loss') THEN 1 ELSE 0 END), 0), 2) as win_rate
		FROM trading_signals
	`

	var stats struct {
		TotalSignals int      `json:"total_signals"`
		Wins         int      `json:"wins"`
		Losses       int      `json:"losses"`
		Pending      int      `json:"pending"`
		AvgProfit    *float64 `json:"avg_profit"`
		WinRate      *float64 `json:"win_rate"`
	}

	err := DB.QueryRow(query).Scan(
		&stats.TotalSignals, &stats.Wins, &stats.Losses,
		&stats.Pending, &stats.AvgProfit, &stats.WinRate,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch performance stats"})
	}

	return c.JSON(stats)
}

// GetStatsByKillZone retrieves statistics grouped by kill zone
func GetStatsByKillZone(c *fiber.Ctx) error {
	query := `
		SELECT 
			kill_zone,
			COUNT(*) as total,
			SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) as wins,
			ROUND(AVG(profit_percent), 2) as avg_profit,
			ROUND(100.0 * SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) / 
				NULLIF(SUM(CASE WHEN status IN ('win', 'loss') THEN 1 ELSE 0 END), 0), 2) as win_rate
		FROM trading_signals
		WHERE status IN ('win', 'loss')
		GROUP BY kill_zone
		ORDER BY win_rate DESC
	`

	rows, err := DB.Query(query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch kill zone stats"})
	}
	defer rows.Close()

	type KillZoneStats struct {
		KillZone  *string  `json:"kill_zone"`
		Total     int      `json:"total"`
		Wins      int      `json:"wins"`
		AvgProfit *float64 `json:"avg_profit"`
		WinRate   *float64 `json:"win_rate"`
	}

	stats := []KillZoneStats{}
	for rows.Next() {
		var s KillZoneStats
		err := rows.Scan(&s.KillZone, &s.Total, &s.Wins, &s.AvgProfit, &s.WinRate)
		if err != nil {
			continue
		}
		stats = append(stats, s)
	}

	return c.JSON(stats)
}

// GetStatsByPattern retrieves statistics grouped by pattern type
func GetStatsByPattern(c *fiber.Ctx) error {
	query := `
		SELECT 
			pattern_type,
			COUNT(*) as total,
			SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) as wins,
			ROUND(AVG(profit_percent), 2) as avg_profit,
			ROUND(100.0 * SUM(CASE WHEN status = 'win' THEN 1 ELSE 0 END) / 
				NULLIF(SUM(CASE WHEN status IN ('win', 'loss') THEN 1 ELSE 0 END), 0), 2) as win_rate
		FROM trading_signals
		WHERE status IN ('win', 'loss') AND pattern_type IS NOT NULL
		GROUP BY pattern_type
		ORDER BY win_rate DESC
	`

	rows, err := DB.Query(query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch pattern stats"})
	}
	defer rows.Close()

	type PatternStats struct {
		PatternType *string  `json:"pattern_type"`
		Total       int      `json:"total"`
		Wins        int      `json:"wins"`
		AvgProfit   *float64 `json:"avg_profit"`
		WinRate     *float64 `json:"win_rate"`
	}

	stats := []PatternStats{}
	for rows.Next() {
		var s PatternStats
		err := rows.Scan(&s.PatternType, &s.Total, &s.Wins, &s.AvgProfit, &s.WinRate)
		if err != nil {
			continue
		}
		stats = append(stats, s)
	}

	return c.JSON(stats)
}
