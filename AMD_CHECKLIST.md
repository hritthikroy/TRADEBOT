# ✅ AMD Phase Implementation - Checklist

## Pre-Test Checklist

- [x] Code implemented in `backend/unified_signal_generator.go`
- [x] Backup created at `backend/unified_signal_generator.go.backup`
- [x] Code compiles without errors
- [x] Binary built successfully (`tradebot`)
- [x] Test scripts created and executable
- [x] Documentation complete

## Testing Checklist

### Step 1: Run Comparison Test
```bash
./compare_before_after_amd.sh
```

- [ ] Script runs without errors
- [ ] Original results displayed
- [ ] AMD results displayed
- [ ] Comparison table shown
- [ ] Recommendation provided

### Step 2: Review Results

Check these metrics:

- [ ] **Win Rate:** AMD > Original + 5%?
- [ ] **Profit Factor:** AMD > Original × 1.2?
- [ ] **Max Drawdown:** AMD < Original × 0.9?
- [ ] **Trade Count:** AMD = 40-60 trades/month?
- [ ] **Signal Quality:** Phase indicators visible?

### Step 3: Make Decision

**If 3+ metrics improved:**
- [ ] Keep AMD phases (already active)
- [ ] Update documentation with actual results
- [ ] Test in paper trading
- [ ] Monitor for 1 week

**If < 3 metrics improved:**
- [ ] Rollback to original
- [ ] Document why it didn't work
- [ ] Consider alternative approaches

## Rollback Checklist (If Needed)

```bash
cp backend/unified_signal_generator.go.backup backend/unified_signal_generator.go
cd backend && go build -o ../tradebot && cd ..
pkill tradebot
./tradebot &
```

- [ ] Original code restored
- [ ] Recompiled successfully
- [ ] Backend restarted
- [ ] Health check passed
- [ ] Test original version works

## Post-Implementation Checklist

### If Keeping AMD Phases

- [ ] Document actual results in `AMD_RESULTS.md`
- [ ] Update `SESSION_TRADER_FINAL_SOLUTION.md`
- [ ] Test in paper trading for 3-7 days
- [ ] Monitor win rate and drawdown
- [ ] Compare with original over time
- [ ] Consider live trading after validation

### If Rolling Back

- [ ] Document why AMD didn't work
- [ ] Save test results for reference
- [ ] Consider alternative improvements
- [ ] Keep documentation for future reference

## Quick Commands

### Test AMD
```bash
./compare_before_after_amd.sh
```

### Rollback
```bash
cp backend/unified_signal_generator.go.backup backend/unified_signal_generator.go
cd backend && go build -o ../tradebot && cd ..
pkill tradebot && ./tradebot &
```

### Check Status
```bash
pgrep tradebot && echo "✅ Running" || echo "❌ Not running"
curl -s http://localhost:8080/health | jq
```

### View Results
```bash
cat /tmp/original_results.json | jq '.strategies[] | select(.name == "Session Trader")'
cat /tmp/amd_results.json | jq '.strategies[] | select(.name == "Session Trader")'
```

## Success Criteria

### Minimum Requirements (Keep AMD)
- ✅ Win rate ≥ 55%
- ✅ Profit factor ≥ 3.5
- ✅ Max drawdown ≤ 30%
- ✅ Trades ≥ 20/month

### Rollback Triggers
- ❌ Win rate < 45%
- ❌ Profit factor < 2.5
- ❌ Max drawdown > 40%
- ❌ Trades < 20/month

## Documentation Files

- [x] `AMD_PHASES_IMPLEMENTATION.md` - Technical details
- [x] `AMD_PHASES_VISUAL_GUIDE.md` - Visual examples
- [x] `START_HERE_AMD_PHASES.md` - Quick start
- [x] `AMD_QUICK_REFERENCE.md` - One-page reference
- [x] `AMD_IMPLEMENTATION_SUMMARY.md` - Complete summary
- [x] `AMD_CHECKLIST.md` - This file

## Test Scripts

- [x] `test_amd_phases.sh` - Quick test
- [x] `compare_before_after_amd.sh` - Full comparison

## Backup Files

- [x] `backend/unified_signal_generator.go.backup` - Original

## Status

- [x] Implementation complete
- [x] Code compiles
- [x] Documentation complete
- [x] Test scripts ready
- [x] Backup available
- [ ] **Testing pending** ← YOU ARE HERE
- [ ] Results reviewed
- [ ] Decision made

## Next Action

**Run this command now:**
```bash
./compare_before_after_amd.sh
```

Then check the boxes above based on results!

---

**Last Updated:** Dec 7, 2025  
**Status:** ✅ Ready for testing  
**Next Step:** Run comparison test
