// Package shadowing parses a timeout setting.
//
// Solution E2. The bug: inside the if block, `timeout, err :=` declared a
// NEW timeout scoped to that block (shadowing the outer one), so the parsed
// value was dropped when the block ended and the function always returned
// def. The fix declares err on its own and assigns timeout with =.
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
		var err error
		timeout, err = time.ParseDuration(s)
		if err != nil {
			return 0, fmt.Errorf("parse timeout %q: %w", s, err)
		}
		if timeout <= 0 {
			return 0, fmt.Errorf("timeout %q must be greater than zero", s)
		}
	}
	return timeout, nil
}
