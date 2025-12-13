#!/bin/bash

echo "🚀 Starting Trading Bot AI Analytics Backend..."
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed!"
    echo "📥 Install Go from: https://golang.org/dl/"
    echo "   Or use: brew install go (on macOS)"
    exit 1
fi

echo "✅ Go is installed: $(go version)"
echo ""

# Check if .env exists
if [ ! -f .env ]; then
    echo "⚠️  .env file not found!"
    echo "📝 Creating .env from .env.example..."
    if [ -f .env.example ]; then
        cp .env.example .env
        echo "✅ Created .env file"
        echo "⚠️  Please edit .env with your database credentials"
        exit 1
    else
        echo "❌ .env.example not found!"
        echo "Please create .env file with DATABASE_URL"
        exit 1
    fi
fi

echo "✅ .env file found"
echo ""

# Install dependencies
echo "📦 Installing dependencies..."
go mod tidy
if [ $? -ne 0 ]; then
    echo "❌ Failed to install dependencies"
    exit 1
fi
echo "✅ Dependencies installed"
echo ""

# Run the server
echo "🚀 Starting server on port 8080..."
echo "📊 AI Analytics: http://localhost:8080/api/v1/analytics/ai"
echo "🏥 Health Check: http://localhost:8080/api/v1/health"
echo ""
echo "Press Ctrl+C to stop"
echo ""

go run .
