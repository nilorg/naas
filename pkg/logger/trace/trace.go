package trace

import (
	"fmt"
	"strconv"
)

const (
	defaultSpanID = "0"
)

// NextSpanID ...
func NextSpanID(spanID string) string {
	if spanID == "" {
		return defaultSpanID
	}
	spanIDLen := len(spanID)
	lastID, _ := strconv.Atoi(spanID[spanIDLen-1:])
	if lastID > 0 {
		spanID = spanID[:spanIDLen-2]
	}
	lastID++
	return fmt.Sprintf("%s.%d", spanID, lastID)
}

// StartSpanID ...
func StartSpanID(spanID string) string {
	if spanID == "" {
		return defaultSpanID
	}
	// spanIDLen := len(spanID)
	// lastID, _ := strconv.Atoi(spanID[spanIDLen-1:])
	// lastID++
	return fmt.Sprintf("%s.%d", spanID, 1)
}
