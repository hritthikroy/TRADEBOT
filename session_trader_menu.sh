#!/bin/bash

# ============================================
# SESSION TRADER - MAIN MENU
# ============================================

while true; do
    clear
    echo "╔════════════════════════════════════════════╗"
    echo "║   SESSION TRADER - TRADING ASSISTANT      ║"
    echo "╚════════════════════════════════════════════╝"
    echo ""
    echo "📊 Daily Workflow:"
    echo "  1) Morning Check (Should I trade today?)"
    echo "  2) Start Paper Trading"
    echo "  3) Check Current Stats"
    echo "  4) Stop Paper Trading"
    echo ""
    echo "🧪 Testing & Analysis:"
    echo "  5) Run Full Backtest (7, 14, 30, 60, 90 days)"
    echo "  6) Run Diagnostic"
    echo "  7) View Dashboard (Browser)"
    echo ""
    echo "📚 Documentation:"
    echo "  8) View Quick Start Guide"
    echo "  9) View Full Strategy Guide"
    echo ""
    echo "  0) Exit"
    echo ""
    echo -n "Select option: "
    read -r choice
    
    case $choice in
        1)
            echo ""
            ./daily_session_trader_check.sh
            echo ""
            echo "Press Enter to continue..."
            read -r
            ;;
        2)
            echo ""
            ./start_paper_trading_session_trader.sh
            echo ""
            echo "Press Enter to continue..."
            read -r
            ;;
        3)
            echo ""
            ./check_paper_trading_stats.sh
            echo ""
            echo "Press Enter to continue..."
            read -r
            ;;
        4)
            echo ""
            ./stop_paper_trading.sh
            echo ""
            echo "Press Enter to continue..."
            read -r
            ;;
        5)
            echo ""
            ./test_session_trader_simple.sh
            echo ""
            echo "Press Enter to continue..."
            read -r
            ;;
        6)
            echo ""
            ./diagnose_session_trader.sh
            echo ""
            echo "Press Enter to continue..."
            read -r
            ;;
        7)
            echo ""
            echo "Opening dashboard in browser..."
            open http://localhost:8080/paper-trading 2>/dev/null || \
            xdg-open http://localhost:8080/paper-trading 2>/dev/null || \
            echo "Please open: http://localhost:8080/paper-trading"
            echo ""
            echo "Press Enter to continue..."
            read -r
            ;;
        8)
            echo ""
            cat SESSION_TRADER_QUICK_START.md | less
            ;;
        9)
            echo ""
            cat SESSION_TRADER_PROFITABLE_VERSION.md | less
            ;;
        0)
            echo ""
            echo "Goodbye! 👋"
            echo ""
            exit 0
            ;;
        *)
            echo ""
            echo "Invalid option. Please try again."
            sleep 2
            ;;
    esac
done
