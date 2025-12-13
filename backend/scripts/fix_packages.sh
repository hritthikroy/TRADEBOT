#!/bin/bash
# fix_packages.sh - Automatically update package declarations

echo "🔧 Fixing package declarations..."

# API handlers
echo "  → Updating API handlers..."
find internal/api/handlers -name "*.go" -exec sed -i '' 's/^package main$/package handlers/' {} \;

# API middleware
echo "  → Updating API middleware..."
find internal/api/middleware -name "*.go" -exec sed -i '' 's/^package main$/package middleware/' {} \;

# API routes
echo "  → Updating API routes..."
sed -i '' 's/^package main$/package api/' internal/api/routes.go

# Backtest
echo "  → Updating backtest package..."
find internal/backtest -name "*.go" -exec sed -i '' 's/^package main$/package backtest/' {} \;

# Database
echo "  → Updating database package..."
find internal/database -name "*.go" -exec sed -i '' 's/^package main$/package database/' {} \;

# Signals
echo "  → Updating signals package..."
find internal/signals -name "*.go" -exec sed -i '' 's/^package main$/package signals/' {} \;

# Strategies
echo "  → Updating strategies packages..."
find internal/strategies -maxdepth 1 -name "*.go" -exec sed -i '' 's/^package main$/package strategies/' {} \;
find internal/strategies/ict -name "*.go" -exec sed -i '' 's/^package main$/package ict/' {} \;
find internal/strategies/patterns -name "*.go" -exec sed -i '' 's/^package main$/package patterns/' {} \;
find internal/strategies/institutional -name "*.go" -exec sed -i '' 's/^package main$/package institutional/' {} \;
find internal/strategies/timeframe -name "*.go" -exec sed -i '' 's/^package main$/package timeframe/' {} \;
find internal/strategies/daily -name "*.go" -exec sed -i '' 's/^package main$/package daily/' {} \;

# Filters
echo "  → Updating filters package..."
find internal/filters -name "*.go" -exec sed -i '' 's/^package main$/package filters/' {} \;

# Optimization
echo "  → Updating optimization package..."
find internal/optimization -name "*.go" -exec sed -i '' 's/^package main$/package optimization/' {} \;

# Trading
echo "  → Updating trading package..."
find internal/trading -name "*.go" -exec sed -i '' 's/^package main$/package trading/' {} \;

# Activity
echo "  → Updating activity package..."
find internal/activity -name "*.go" -exec sed -i '' 's/^package main$/package activity/' {} \;

# AI
echo "  → Updating AI package..."
find internal/ai -name "*.go" -exec sed -i '' 's/^package main$/package ai/' {} \;

# Communication
echo "  → Updating communication packages..."
find internal/communication/telegram -name "*.go" -exec sed -i '' 's/^package main$/package telegram/' {} \;
find internal/communication/websocket -name "*.go" -exec sed -i '' 's/^package main$/package websocket/' {} \;

# Templates
echo "  → Updating templates package..."
find internal/templates -name "*.go" -exec sed -i '' 's/^package main$/package templates/' {} \;

echo ""
echo "✅ Package declarations updated!"
echo ""
echo "📋 Next steps:"
echo "  1. Review changes: git diff"
echo "  2. Update imports in main.go and other files"
echo "  3. Run: go mod tidy"
echo "  4. Build: go build ./cmd/server"
echo "  5. Test: go test ./..."
echo ""
