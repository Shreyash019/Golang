package main

import "strings"

// StringUtils groups string helper methods.
type StringUtils struct{}

// IsNullOrEmpty reports whether s is empty.
func (StringUtils) IsNullOrEmpty(s string) bool { return len(strings.TrimSpace(s)) == 0 }

// SplitKeyValue splits "key=value" into its parts.
func (StringUtils) SplitKeyValue(s string) (string, string) {
	parts := strings.SplitN(s, "=", 2)
	if len(parts) != 2 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}
