package utils

import "time"

var (
	defaultString string
	defaultInt    int
)

// DefaultString return defaultVal if val is empty
func DefaultString(val string, defaultVal string) string {
	if val == defaultString {
		return defaultVal
	}

	return val
}

// DefaultInt return defaultVal if val is empty
func DefaultInt(val int, defaultVal int) int {
	if val == defaultInt {
		return defaultVal
	}

	return val
}

// DefaultDuration return defaultVal if val is empty fpr time.Duration
func DefaultDuration(val time.Duration, defaultVal time.Duration) time.Duration {
	if val == 0 {
		return defaultVal
	}

	return val
}

// DefaultInterface return defaultVal if val is empty
func DefaultInterface(val interface{}, defaultVal interface{}) interface{} {
	if val == nil {
		return defaultVal
	}

	return val
}
