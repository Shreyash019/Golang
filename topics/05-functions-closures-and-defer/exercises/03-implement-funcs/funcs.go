// Exercise E3 — implement.
//
// Goal: implement the six functions below until every test passes. They
// train multiple results and a first error (notes §5.2, §5.8), variadic
// parameters (§5.3), functions as values (§5.4), closures sharing state
// (§5.5), a higher-order helper with error wrapping (§5.4, §5.8) and
// recursion (§5.6).
//
//	make ex E=03-implement-funcs

// Package funcs is a set of small functions built from Go's function
// features.
package funcs

import "errors"

// ErrDivideByZero is returned by Divide when the divisor is zero.
var ErrDivideByZero = errors.New("divide by zero")

// Divide returns the quotient and remainder of a / b, with Go's truncating
// division (Divide(-7, 2) is -3, -1). When b is 0 it returns 0, 0 and
// ErrDivideByZero.
func Divide(a, b int) (q, r int, err error) {
	return 0, 0, nil
}

// Sum returns the sum of nums; Sum() is 0.
func Sum(nums ...int) int {
	return 0
}

// Pipeline returns a function that passes its argument through fs from left
// to right: Pipeline(f, g)(x) is g(f(x)). With no functions it returns its
// argument unchanged.
func Pipeline(fs ...func(int) int) func(int) int {
	return nil
}

// Counter returns two functions sharing one counter. next returns start on
// its first call, then start+step, start+2*step, and so on. reset makes the
// next call to next return start again. Counters from different calls to
// Counter are independent.
func Counter(start, step int) (next func() int, reset func()) {
	return nil, nil
}

// Retry calls fn with attempt numbers 1, 2, … until fn returns nil or
// attempts calls have been made (an attempts value below 1 counts as 1). It
// returns nil on success; otherwise an error reading
// "after N attempts: <last error>" that wraps the last error, so errors.Is
// sees it.
func Retry(attempts int, fn func(attempt int) error) error {
	return nil
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
func Walk(n *Node, visit func(path string, depth int) bool) {
}
