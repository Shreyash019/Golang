// Package examples holds the runnable worked examples from this topic's notes.
//
// Each example writes its output to an io.Writer, so the same code is run by
// cmd/examples and checked by the Example tests in examples_test.go: the
// output shown in the notes is the output the tests verify.
package examples

import "io"

// Example is one worked example from the notes.
type Example struct {
	// Name is the kebab-case name used on the command line, e.g. "slice-aliasing".
	Name string
	// Section is the notes section the example belongs to, e.g. "5.2".
	Section string
	// Summary is a one-line description shown when listing examples.
	Summary string
	// Crashes marks examples that panic, deadlock or race on purpose.
	// They run only when named explicitly, never as part of "all".
	Crashes bool
	// Run executes the example, writing its output to w.
	Run func(w io.Writer)
}

// All returns the examples in the order they appear in the notes.
func All() []Example {
	return []Example{
		{Name: "pass-by-value", Section: "5.1", Summary: "every argument is a copy; what the copy shares depends on the type", Run: PassByValue},
		{Name: "multiple-returns", Section: "5.2", Summary: "several results, comma-ok, _, f(g()) chaining", Run: MultipleReturns},
		{Name: "named-results", Section: "5.2", Summary: "named results start at zero; bare return", Run: NamedResults},
		{Name: "variadic", Section: "5.3", Summary: "...T is a []T; f(), f(1, 2), f(s...)", Run: Variadic},
		{Name: "variadic-aliasing", Section: "5.3", Summary: "f(s...) passes the caller's slice, not a copy", Run: VariadicAliasing},
		{Name: "function-values", Section: "5.4", Summary: "func types, nil func values, anonymous functions", Run: FunctionValues},
		{Name: "higher-order", Section: "5.4", Summary: "functions taking and returning functions: filter, compose", Run: HigherOrder},
		{Name: "nil-func-call", Section: "5.4", Summary: "calling a nil func value panics", Crashes: true, Run: NilFuncCall},
		{Name: "closure-counter", Section: "5.5", Summary: "each closure keeps its own captured state alive", Run: ClosureCounter},
		{Name: "closure-shares-variable", Section: "5.5", Summary: "closures capture variables, not values", Run: ClosureSharesVariable},
		{Name: "closure-state", Section: "5.5", Summary: "two closures sharing private state", Run: ClosureState},
		{Name: "recursive-closure", Section: "5.5", Summary: "declare the variable first, then assign the literal", Run: RecursiveClosure},
		{Name: "recursion", Section: "5.6", Summary: "base case, recursive case, walking a tree", Run: Recursion},
		{Name: "stack-overflow", Section: "5.6", Summary: "unbounded recursion: fatal error at the 1 GB stack limit", Crashes: true, Run: StackOverflow},
		{Name: "defer-order", Section: "5.7", Summary: "deferred calls run LIFO when the function returns", Run: DeferOrder},
		{Name: "defer-args-evaluated", Section: "5.7", Summary: "arguments are evaluated at the defer statement", Run: DeferArgsEvaluated},
		{Name: "defer-modifies-result", Section: "5.7", Summary: "a deferred closure can change named results", Run: DeferModifiesResult},
		{Name: "defer-runs-on-panic", Section: "5.7", Summary: "deferred calls run while a panic unwinds", Crashes: true, Run: DeferRunsOnPanic},
		{Name: "defer-in-loop", Section: "5.7", Summary: "defer in a loop waits for the function; a func per iteration fixes it", Run: DeferInLoop},
		{Name: "defer-trace", Section: "6.1", Summary: "defer trace(name)(): enter now, leave on return", Run: DeferTrace},
		{Name: "error-first-look", Section: "5.8", Summary: "error as the last result, %w wrapping, errors.Is", Run: ErrorFirstLook},
	}
}
