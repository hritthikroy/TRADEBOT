# 📅 Calendar Date Picker Feature

## What's New?

The "Days to Test" parameter now has **TWO MODES**:

### 1. **Days Mode (Default)** 📊
- Simple number input (e.g., 30, 60, 90 days)
- Tests the last X days from today
- Shows date range automatically
- **Example**: Enter "30" → Tests last 30 days

### 2. **Calendar Mode** 📆
- Toggle the "Use Calendar" checkbox
- Choose from preset periods OR custom dates
- Automatically calculates days between dates

## How to Use

### Quick Start (Days Mode)
1. Enter number of days (e.g., 30)
2. See the date range displayed below
3. Click "Run Backtest"

### Advanced (Calendar Mode)
1. ✅ Check "Use Calendar"
2. Choose an option:
   - **Recent Data**: Last 30 days
   - **2024 Bull Run**: Jan-Mar 2024 (90 days, +74%)
   - **2023 Bull Run**: Oct-Dec 2023 (92 days, +63%)
   - **2021 Bull Run**: Jan-Apr 2021 (120 days, +120%)
   - **2020 Bull Run**: Oct-Dec 2020 (92 days, +190%)
   - **Custom Date Range**: Pick any dates you want!

3. For custom dates:
   - Select **Start Date** (when to begin testing)
   - Select **End Date** (when to stop testing)
   - Days are calculated automatically
   - See the summary below the date picker

## Features

✅ **Visual Date Display**: See exactly what period you're testing
✅ **Automatic Calculation**: Days calculated from calendar dates
✅ **Preset Bull Runs**: Test historical profitable periods
✅ **Custom Ranges**: Pick any date range you want
✅ **Validation**: Prevents invalid date ranges
✅ **Real-time Updates**: See changes as you select dates

## Examples

### Example 1: Test Last 7 Days
- Mode: Days Mode
- Input: `7`
- Result: Tests from Dec 27, 2024 to Dec 4, 2024

### Example 2: Test 2024 Bull Run
- Mode: Calendar Mode
- Select: "🐂 2024 Bull Run"
- Result: Tests Jan 1 - Mar 31, 2024 (90 days)

### Example 3: Custom Period
- Mode: Calendar Mode
- Select: "Custom Date Range"
- Start: Nov 1, 2024
- End: Dec 1, 2024
- Result: Tests 30 days in November 2024

## Benefits

🎯 **More Intuitive**: Pick dates visually instead of counting days
📊 **Historical Testing**: Easy access to known bull/bear markets
🔍 **Precise Control**: Test exact date ranges
✨ **Better UX**: See what you're testing before you run it

## Technical Details

- Both modes update the same `days` parameter in the backend
- Calendar mode calculates days automatically
- Backend receives the number of days (no changes needed)
- Date validation prevents errors
- Works with all strategies and filters

## Tips

💡 **Tip 1**: Use preset bull runs to see how strategies perform in trending markets
💡 **Tip 2**: Use custom dates to test specific events (e.g., halving, major news)
💡 **Tip 3**: Days mode is faster for quick tests
💡 **Tip 4**: Calendar mode is better for precise historical analysis

---

**Status**: ✅ Fully Implemented and Working
**Compatibility**: Works with all 10 strategies
**Backend Changes**: None required (uses existing days parameter)
