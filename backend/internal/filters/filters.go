package filters

// SQL filter to exclude test data
const TEST_DATA_FILTER = `
	AND signal_id NOT LIKE 'test_%'
	AND signal_id NOT LIKE 'perm_test_%'
	AND symbol NOT LIKE '%TEST%'
	AND entry_price NOT IN (1, 50000)
`

// Helper function to add test filter to WHERE clause
func addTestFilter(baseQuery string) string {
	return baseQuery + TEST_DATA_FILTER
}
