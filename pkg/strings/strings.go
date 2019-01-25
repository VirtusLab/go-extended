package strings

// Ellipsis returns the string in  an abbreviated form with maximum number of characters
func Ellipsis(s string, max int) string {
	if max <= 4 {
		return s[:max]
	} else if len(s) <= max {
		return s
	}
	offset := max - 3
	return s[:offset] + "..."
}
