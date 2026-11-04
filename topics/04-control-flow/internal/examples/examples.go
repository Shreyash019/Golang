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
		{Name: "if-with-init", Section: "5.1", Summary: "an if's init variable is scoped to the if and all its else branches", Run: IfWithInit},
		{Name: "for-forms", Section: "5.2", Summary: "three-clause, condition-only (while), infinite, two indexes", Run: ForForms},
		{Name: "loopvar-per-iteration", Section: "5.3", Summary: "Go 1.22+: a fresh loop variable per iteration", Run: LoopVarPerIteration},
		{Name: "loopvar-legacy", Section: "5.3", Summary: "Go 1.21 semantics: one variable shared by every iteration", Run: LoopVarLegacy},
		{Name: "loopvar-copy-back", Section: "5.3", Summary: "changing i in the body: carried over (3-clause) vs lost (range)", Run: LoopVarCopyBack},
		{Name: "range-evaluated-once", Section: "5.4", Summary: "the range expression is evaluated once; arrays are copied", Run: RangeEvaluatedOnce},
		{Name: "range-copies-values", Section: "5.4", Summary: "the iteration value is a copy; write through the index", Run: RangeCopiesValues},
		{Name: "range-over-int", Section: "5.5", Summary: "Go 1.22: for i := range n, typed like n, zero times for n <= 0", Run: RangeOverInt},
		{Name: "range-over-slice", Section: "5.5", Summary: "index and value, index only, value only", Run: RangeOverSlice},
		{Name: "range-over-string", Section: "5.5", Summary: "byte offsets and runes; invalid UTF-8 becomes U+FFFD", Run: RangeOverString},
		{Name: "range-over-map", Section: "5.5", Summary: "unspecified order; sort the keys for a stable one", Run: RangeOverMap},
		{Name: "map-delete-during-range", Section: "5.5", Summary: "deleting while ranging over a map is safe", Run: MapDeleteDuringRange},
		{Name: "range-over-channel", Section: "5.5", Summary: "receive until the channel is closed", Run: RangeOverChannel},
		{Name: "range-over-func", Section: "5.5", Summary: "Go 1.23: range over an iterator; break stops it", Run: RangeOverFunc},
		{Name: "switch-basics", Section: "5.6", Summary: "several values per case, no fallthrough, default anywhere", Run: SwitchBasics},
		{Name: "switch-tagless", Section: "5.6", Summary: "switch with no tag: the first true case wins", Run: SwitchTagless},
		{Name: "switch-with-init", Section: "5.6", Summary: "a switch's init statement, scoped to the switch", Run: SwitchWithInit},
		{Name: "switch-eval-order", Section: "5.6", Summary: "cases are evaluated lazily, top to bottom, left to right", Run: SwitchEvalOrder},
		{Name: "fallthrough", Section: "5.7", Summary: "fallthrough runs the next body without testing its case", Run: Fallthrough},
		{Name: "type-switch", Section: "5.7", Summary: "a first look: branching on an interface's dynamic type", Run: TypeSwitch},
		{Name: "labeled-break-continue", Section: "5.8", Summary: "continue rows / break search on an outer loop", Run: LabeledBreakContinue},
		{Name: "goto", Section: "5.8", Summary: "goto jumping backward: what a loop compiles to", Run: Goto},
		{Name: "guard-clauses", Section: "6.1", Summary: "return early; keep the happy path unindented", Run: GuardClauses},
		{Name: "read-loop", Section: "6.2", Summary: "for sc.Scan(): continue, break, then check sc.Err()", Run: ReadLoop},
		{Name: "break-in-switch", Section: "7.1", Summary: "break in a switch leaves the switch, not the loop", Run: BreakInSwitch},
	}
}
