// Command backoff prints an exponential backoff schedule: how long a client
// waits before each retry attempt, and the total.
//
//	go run ./exercises/solutions/04-refactor-backoff      # the default 5 attempts
//	go run ./exercises/solutions/04-refactor-backoff 3    # 3 attempts
//
// Solution E4. What changed, and why:
//   - Settings that never change are constants, named in MixedCaps. A
//     constant can't be reassigned by accident, is safe to read from any
//     goroutine, and costs nothing at runtime.
//   - baseDelay is a typed time.Duration constant, so the unit lives in the
//     type, not in a "Ms" suffix, and no conversions are needed.
//   - The getters and the BackoffConstants type are gone: Go reads a
//     constant directly.
//   - The running total is a local variable, so two runs can't interfere.
//   - Variables are declared where they are first needed, with := when
//     there is a value and var when the zero value is the starting point.
package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

const (
	maxAttempts = 5
	baseDelay   = 100 * time.Millisecond
	multiplier  = 2
)

func main() {
	if err := run(os.Stdout, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "backoff:", err)
		os.Exit(1)
	}
}

func run(w io.Writer, args []string) error {
	attempts := maxAttempts
	if len(args) > 0 {
		n, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid attempts %q: %w", args[0], err)
		}
		attempts = n
	}
	if attempts < 1 {
		return fmt.Errorf("attempts must be at least 1, got %d", attempts)
	}
	if attempts > maxAttempts {
		fmt.Fprintf(w, "requested %d attempts, capped at %d\n", attempts, maxAttempts)
		attempts = maxAttempts
	}

	var total time.Duration
	for i := 1; i <= attempts; i++ {
		d := delay(i)
		total += d
		fmt.Fprintf(w, "attempt %d: wait %v\n", i, d)
	}
	fmt.Fprintf(w, "total wait: %v\n", total)
	return nil
}

// delay returns the wait before the given attempt (1-based): baseDelay,
// multiplied by multiplier for every attempt after the first.
func delay(attempt int) time.Duration {
	d := baseDelay
	for i := 1; i < attempt; i++ {
		d *= multiplier // untyped constant: no conversion needed
	}
	return d
}
