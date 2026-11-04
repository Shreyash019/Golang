// Package funcs is a set of small functions built from Go's function
// features.
package funcs

import (
	"errors"
	"fmt"
)

// ErrDivideByZero is returned by Divide when the divisor is zero.
var ErrDivideByZero = errors.New("divide by zero")

// Divide returns the quotient and remainder of a / b, with Go's truncating
// division (Divide(-7, 2) is -3, -1). When b is 0 it returns 0, 0 and
// ErrDivideByZero.
//
// Returning an error instead of letting a / b panic turns a programmer bug
// into a value the caller can handle. On error the other results are zero
// values, the convention callers rely on.
func Divide(a, b int) (q, r int, err error) {
	if b == 0 {
		return 0, 0, ErrDivideByZero
	}
	return a / b, a % b, nil
}

// Sum returns the sum of nums; Sum() is 0.
//
// Inside the function nums is an ordinary []int (nil when called with no
// arguments), and ranging over a nil slice runs zero times.
func Sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// Pipeline returns a function that passes its argument through fs from left
// to right: Pipeline(f, g)(x) is g(f(x)). With no functions it returns its
// argument unchanged.
//
// The returned closure captures fs, so the loop happens on every call, not
// once here. A caller who passes a slice with fs... and later changes it
// would see the change: copy fs first (slices.Clone) if that matters.
func Pipeline(fs ...func(int) int) func(int) int {
	return func(x int) int {
		for _, f := range fs {
			x = f(x)
		}
		return x
	}
}

// Counter returns two functions sharing one counter. next returns start on
// its first call, then start+step, start+2*step, and so on. reset makes the
// next call to next return start again. Counters from different calls to
// Counter are independent.
//
// Both closures capture the same variable, current, which escapes to the
// heap; each call to Counter creates a new one, so counters are
// independent. This is the closure alternative to a small struct with two
// methods (topic 10).
func Counter(start, step int) (next func() int, reset func()) {
	current := start
	next = func() int {
		v := current
		current += step
		return v
	}
	reset = func() { current = start }
	return next, reset
}

// Retry calls fn with attempt numbers 1, 2, … until fn returns nil or
// attempts calls have been made (an attempts value below 1 counts as 1). It
// returns nil on success; otherwise an error reading
// "after N attempts: <last error>" that wraps the last error, so errors.Is
// sees it.
//
// fn is a func value, so Retry knows nothing about what is being retried.
// Real retries also wait between attempts, with backoff and jitter, and
// stop when a context is cancelled (topics 26 and 42).
func Retry(attempts int, fn func(attempt int) error) error {
	attempts = max(attempts, 1)
	var err error
	for attempt := 1; attempt <= attempts; attempt++ {
		if err = fn(attempt); err == nil {
			return nil
		}
	}
	return fmt.Errorf("after %d attempts: %w", attempts, err)
}

// Node is a node in a tree, such as a directory and its entries.
type Node struct {
	Name     string
	Children []*Node
}

// Walk visits n and its descendants depth first, parents before children,
// children in order. visit gets the slash-separated path from the root
// ("root/a/b") and the depth (the root is 0). When visit returns false, the
// children of that node are skipped. Walk(nil, visit) does nothing.
//
// A tree is a recursive structure, so a recursive function is the natural
// fit: the exported function handles nil and starts the recursion, and the
// helper carries the path and depth down. The recursion depth equals the
// tree's depth, which is fine for file trees; for untrusted, possibly very
// deep input, an explicit stack (a slice) bounds memory instead.
func Walk(n *Node, visit func(path string, depth int) bool) {
	if n == nil {
		return
	}
	walk(n, n.Name, 0, visit)
}

func walk(n *Node, path string, depth int, visit func(string, int) bool) {
	if !visit(path, depth) {
		return
	}
	for _, c := range n.Children {
		walk(c, path+"/"+c.Name, depth+1, visit)
	}
}
