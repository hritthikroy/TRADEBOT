#!/bin/bash

# Comprehensive Test Script for Trading Bot
# Usage: ./test.sh [option]
# Options: signal, supabase, telegram, all

echo "🧪 Trading Bot Test Suite"
echo "=========================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

API_URL="http://localhost:8080/api/v1"

# Test 1: Signal Generation
test_signal() {
    echo "📊 Testing Signal Generation..."
    response=$(curl -s -X POST ${API_URL}/backtest/live-signal \
        -H "Content-Type: application/json" \
        -d '{"symbol":"BTCUSDT","strategy":"session_trader"}')
    
    if echo "$response" | grep -q "signal"; then
        echo -e "${GREEN}✅ Signal generation working${NC}"
        echo "$response" | jq '.'
    else
        echo -e "${RED}❌ Signal generation failed${NC}"
        echo "$response"
    fi
    echo ""
}

# Test 2: Supabase Connection
test_supabase() {
    echo "💾 Testing Supabase Connection..."
    
    if [ -z "$SUPABASE_URL" ] || [ -z "$SUPABASE_KEY" ]; then
        echo -e "${YELLOW}⚠️  Loading from .env file${NC}"
        if [ -f "backend/.env" ]; then
            export $(cat backend/.env | grep -v '^#' | xargs)
        fi
    fi
    
    if [ -z "$SUPABASE_URL" ]; then
        echo -e "${RED}❌ SUPABASE_URL not set${NC}"
        return 1
    fi
    
    response=$(curl -s "${SUPABASE_URL}/rest/v1/trading_signals?limit=1" \
        -H "apikey: ${SUPABASE_KEY}" \
        -H "Authorization: Bearer ${SUPABASE_KEY}")
    
    if echo "$response" | grep -q "\["; then
        echo -e "${GREEN}✅ Supabase connection working${NC}"
        echo "Recent signals:"
        echo "$response" | jq '.'
    else
        echo -e "${RED}❌ Supabase connection failed${NC}"
        echo "$response"
    fi
    echo ""
}

# Test 3: Telegram Bot Status
test_telegram() {
    echo "📱 Testing Telegram Bot..."
    response=$(curl -s ${API_URL}/telegram/status)
    
    if echo "$response" | grep -q "running"; then
        echo -e "${GREEN}✅ Telegram bot status retrieved${NC}"
        echo "$response" | jq '.'
    else
        echo -e "${YELLOW}⚠️  Telegram bot may not be running${NC}"
        echo "$response"
    fi
    echo ""
}

# Test 4: Backend Health
test_health() {
    echo "🏥 Testing Backend Health..."
    response=$(curl -s ${API_URL}/health)
    
    if echo "$response" | grep -q "status"; then
        echo -e "${GREEN}✅ Backend is healthy${NC}"
        echo "$response" | jq '.'
    else
        echo -e "${RED}❌ Backend health check failed${NC}"
        echo "$response"
    fi
    echo ""
}

# Test 5: Filter Settings
test_filters() {
    echo "🎯 Testing Filter Settings..."
    response=$(curl -s ${API_URL}/settings)
    
    if echo "$response" | grep -q "filterBuy"; then
        echo -e "${GREEN}✅ Filter settings retrieved${NC}"
        echo "$response" | jq '.'
    else
        echo -e "${RED}❌ Filter settings failed${NC}"
        echo "$response"
    fi
    echo ""
}

# Main execution
case "$1" in
    signal)
        test_signal
        ;;
    supabase)
        test_supabase
        ;;
    telegram)
        test_telegram
        ;;
    health)
        test_health
        ;;
    filters)
        test_filters
        ;;
    all|"")
        test_health
        test_signal
        test_filters
        test_supabase
        test_telegram
        ;;
    *)
        echo "Usage: $0 {signal|supabase|telegram|health|filters|all}"
        echo ""
        echo "Examples:"
        echo "  ./test.sh           # Run all tests"
        echo "  ./test.sh signal    # Test signal generation only"
        echo "  ./test.sh supabase  # Test Supabase connection only"
        exit 1
        ;;
esac

echo "=========================="
echo "✅ Test suite complete"
