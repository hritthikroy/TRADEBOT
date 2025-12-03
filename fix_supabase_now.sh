#!/bin/bash

echo "🔧 SUPABASE QUICK FIX"
echo "===================="
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}This script will help you fix Supabase storage issues.${NC}"
echo ""

# Step 1: Check if diagnostic script exists
if [ ! -f "diagnose_supabase.sh" ]; then
    echo -e "${RED}❌ diagnose_supabase.sh not found${NC}"
    exit 1
fi

# Step 2: Run diagnostic
echo -e "${YELLOW}Running diagnostic...${NC}"
echo ""

if ./diagnose_supabase.sh; then
    echo ""
    echo -e "${GREEN}✅ Everything is working!${NC}"
    echo ""
    echo "Your Supabase is configured correctly."
    echo "Signals should be saving automatically."
    echo ""
    exit 0
fi

# If diagnostic failed, show fix instructions
echo ""
echo -e "${RED}❌ Diagnostic found issues${NC}"
echo ""
echo -e "${YELLOW}📋 FIX INSTRUCTIONS:${NC}"
echo ""
echo "1. Open Supabase SQL Editor:"
echo -e "   ${BLUE}https://supabase.com/dashboard/project/elqhqhjevajzjoghiiss/sql${NC}"
echo ""
echo "2. Click 'New Query'"
echo ""
echo "3. Copy the contents of this file:"
echo -e "   ${BLUE}supabase-setup.sql${NC}"
echo ""
echo "4. Paste into SQL Editor and click 'Run'"
echo ""
echo "5. Run this script again to verify:"
echo -e "   ${BLUE}./fix_supabase_now.sh${NC}"
echo ""
echo -e "${YELLOW}Need help? Check these guides:${NC}"
echo "  - FIX_SUPABASE_STORAGE.md"
echo "  - SUPABASE_FIX_COMPLETE.md"
echo ""
