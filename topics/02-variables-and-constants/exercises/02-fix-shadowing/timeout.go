// Package shadowing parses a timeout setting.
//
// Exercise E2 — fix the bug. ParseTimeout compiles, passes go vet, and
// handles bad input correctly, yet every valid timeout is ignored. Find the
// root cause, fix it with the smallest change, and explain the gotcha in one
// sentence (notes §7.1).
//
//	make ex E=02-fix-shadowing
package shadowing

import (
	"fmt"
	"time"
)

// ParseTimeout parses s as a time.Duration, such as "1.5s" or "300ms".
// An empty s means "not configured" and returns def. A timeout must be
// greater than zero.
func ParseTimeout(s string, def time.Duration) (time.Duration, error) {
	timeout := def
	if s != "" {
		timeout, err := time.ParseDuration(s)
		if err != nil {
			return 0, fmt.Errorf("parse timeout %q: %w", s, err)
		}
		if timeout <= 0 {
			return 0, fmt.Errorf("timeout %q must be greater than zero", s)
		}
	}
	return timeout, nil
}
