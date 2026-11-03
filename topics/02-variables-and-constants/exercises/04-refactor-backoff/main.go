// Exercise E4 — refactor to idiomatic Go.
//
// Command backoff prints an exponential backoff schedule: how long a client
// waits before each retry attempt, and the total.
//
//	go run ./exercises/04-refactor-backoff      # the default 5 attempts
//	go run ./exercises/04-refactor-backoff 3    # 3 attempts
//
// It works, and behaviour_test.go passes. But it's written the way a Java
// "constants class" or a Python settings module would be: mutable globals in
// SCREAMING_SNAKE_CASE, getters, explicit zero initializers, durations as
// ints with a unit suffix, and global state reset on every call.
//
// Rewrite it as idiomatic Go (notes §8). behaviour_test.go must keep
// passing; idiom_test.go fails until the habits are gone:
//
//	make ex E=04-refactor-backoff
package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

// MAX_RETRIES is the maximum number of attempts.
var MAX_RETRIES int = 5

// BASE_DELAY_MS is the delay before the first attempt, in milliseconds.
var BASE_DELAY_MS int = 100

// BACKOFF_MULTIPLIER multiplies the delay after each attempt.
var BACKOFF_MULTIPLIER int = 2

// totalWaitMs accumulates the schedule's total, in milliseconds.
var totalWaitMs int = 0

// BackoffConstants gives access to the settings.
type BackoffConstants struct{}

// GetMaxRetries returns MAX_RETRIES.
func (c BackoffConstants) GetMaxRetries() int { return MAX_RETRIES }

// GetBaseDelayMs returns BASE_DELAY_MS.
func (c BackoffConstants) GetBaseDelayMs() int { return BASE_DELAY_MS }

func main() {
	if err := run(os.Stdout, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "backoff:", err)
		os.Exit(1)
	}
}

func run(w io.Writer, args []string) error {
	var constants BackoffConstants = BackoffConstants{}
	var attempts int = 0
	var err error = nil
	attempts = constants.GetMaxRetries()
	if len(args) > 0 {
		attempts, err = strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid attempts %q: %w", args[0], err)
		}
	}
	if attempts < 1 {
		return fmt.Errorf("attempts must be at least 1, got %d", attempts)
	}
	if attempts > constants.GetMaxRetries() {
		fmt.Fprintf(w, "requested %d attempts, capped at %d\n", attempts, constants.GetMaxRetries())
		attempts = constants.GetMaxRetries()
	}

	totalWaitMs = 0
	for i := 1; i <= attempts; i++ {
		var delayMs int = GetDelayMs(i, constants)
		totalWaitMs = totalWaitMs + delayMs
		fmt.Fprintf(w, "attempt %d: wait %v\n", i, time.Duration(delayMs)*time.Millisecond)
	}
	fmt.Fprintf(w, "total wait: %v\n", time.Duration(totalWaitMs)*time.Millisecond)
	return nil
}

// GetDelayMs returns the delay before the given attempt, in milliseconds.
func GetDelayMs(attempt int, constants BackoffConstants) int {
	var delayMs int = 0
	delayMs = constants.GetBaseDelayMs()
	for i := 1; i < attempt; i++ {
		delayMs = delayMs * BACKOFF_MULTIPLIER
	}
	return delayMs
}
