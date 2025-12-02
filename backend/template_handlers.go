package main

import (
	"html/template"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// HandleIndexPage renders the main dashboard
func HandleIndexPage(c *fiber.Ctx) error {
	tmpl, err := template.New("index").Funcs(templateFuncs).Parse(GetIndexTemplate())
	if err != nil {
		return c.Status(500).SendString("Template error: " + err.Error())
	}

	data := TemplateData{
		Title:       "Trading Bot - Go Backend",
		CurrentTime: time.Now(),
	}

	c.Set("Content-Type", "text/html; charset=utf-8")
	return tmpl.Execute(c.Response().BodyWriter(), data)
}

// HandleBacktestForm processes backtest form submission
func HandleBacktestForm(c *fiber.Ctx) error {
	// Parse form data
	symbol := c.FormValue("symbol", "BTCUSDT")
	interval := c.FormValue("interval", "15m")
	days, _ := strconv.Atoi(c.FormValue("days", "30"))
	balance, _ := strconv.ParseFloat(c.FormValue("balance", "500"), 64)
	risk, _ := strconv.ParseFloat(c.FormValue("risk", "2"), 64)

	// Create config
	config := BacktestConfig{
		Symbol:       symbol,
		Interval:     interval,
		Days:         days,
		StartBalance: balance,
		RiskPercent:  risk / 100,
	}

	// Fetch data
	candles, err := fetchBinanceData(config.Symbol, config.Interval, config.Days)
	if err != nil {
		return c.Status(500).SendString("Failed to fetch data: " + err.Error())
	}

	// Run backtest
	result, err := RunBacktest(config, candles)
	if err != nil {
		return c.Status(500).SendString("Backtest failed: " + err.Error())
	}

	// Render results
	tmpl, err := template.New("results").Funcs(templateFuncs).Parse(GetBacktestResultsTemplate())
	if err != nil {
		return c.Status(500).SendString("Template error: " + err.Error())
	}

	data := TemplateData{
		Title:       "Backtest Results",
		CurrentTime: time.Now(),
		Data:        result,
	}

	c.Set("Content-Type", "text/html; charset=utf-8")
	return tmpl.Execute(c.Response().BodyWriter(), data)
}

// HandleLiveSignalsPage renders the live signals page
func HandleLiveSignalsPage(c *fiber.Ctx) error {
	tmpl, err := template.New("signals").Parse(GetLiveSignalsTemplate())
	if err != nil {
		return c.Status(500).SendString("Template error: " + err.Error())
	}

	data := TemplateData{
		Title:       "Live Signals",
		CurrentTime: time.Now(),
	}

	c.Set("Content-Type", "text/html; charset=utf-8")
	return tmpl.Execute(c.Response().BodyWriter(), data)
}
