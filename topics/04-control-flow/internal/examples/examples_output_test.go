package examples

import "os"

func ExampleIfWithInit() {
	IfWithInit(os.Stdout)
	// Output:
	// "42": ok, n*2 = 84
	// "-7": negative -7
	// "abc": not a number (strconv.Atoi: parsing "abc": invalid syntax)
}

func ExampleForForms() {
	ForForms(os.Stdout)
	// Output:
	// 0 1 2 | three-clause
	// 111 | while: Collatz steps from 27
	// 3 | infinite loop left on attempt 3
	// desserts | two indexes
}

func ExampleLoopVarPerIteration() {
	LoopVarPerIteration(os.Stdout)
	// Output:
	// 0 1 2 | closures
	// 0 1 2 | pointers
	// true | different variables
}

func ExampleLoopVarLegacy() {
	LoopVarLegacy(os.Stdout)
	// Output:
	// 3 3 3 | closures (Go 1.21 semantics)
	// 30 30 30 | pointers (Go 1.21 semantics)
}

func ExampleLoopVarCopyBack() {
	LoopVarCopyBack(os.Stdout)
	// Output:
	// 0 3 6 9 | three-clause: the body's change carries over
	// 0 1 2 3 4 5 6 7 8 9 | range: the body's change is lost
}

func ExampleRangeEvaluatedOnce() {
	RangeEvaluatedOnce(os.Stdout)
	// Output:
	// [1 2 3 10 20 30] | append during range
	// 3 | array: the copy still has 3
	// 100 | slice: the shared element is 100
}

func ExampleRangeCopiesValues() {
	RangeCopiesValues(os.Stdout)
	// Output:
	// [{pen 100} {book 900}] | after doubling the copy
	// [{pen 200} {book 1800}] | after doubling through the index
}

func ExampleRangeOverInt() {
	RangeOverInt(os.Stdout)
	// Output:
	// 0 1 2 | range 3
	// 0:uint8 1:uint8 | i has n's type
	// 0 | range -5 runs zero times
	// hi hi | no variable at all
}

func ExampleRangeOverSlice() {
	RangeOverSlice(os.Stdout)
	// Output:
	// 0=go 1=python 2=java | index and value
	// 0 1 2 | index only
	// 2 6 4 | value only
}

func ExampleRangeOverString() {
	RangeOverString(os.Stdout)
	// Output:
	// len: 6
	// 0:h 1:é 3:l 4:l 5:o | byte offset : rune
	// 68 c3 a9 6c 6c 6f | indexing gives bytes
	// 0:U+0061 1:U+FFFD 2:U+0062 | invalid UTF-8
}

func ExampleRangeOverMap() {
	RangeOverMap(os.Stdout)
	// Output:
	// apple=5 fig=0 pear=3 | sorted keys
}

func ExampleMapDeleteDuringRange() {
	MapDeleteDuringRange(os.Stdout)
	// Output:
	// 5 [0 2 4 6 8]
}

func ExampleRangeOverChannel() {
	RangeOverChannel(os.Stdout)
	// Output:
	// a b c | received until closed
}

func ExampleRangeOverFunc() {
	RangeOverFunc(os.Stdout)
	// Output:
	// 5 4 3 | break stops the iterator
}

func ExampleSwitchBasics() {
	SwitchBasics(os.Stdout)
	// Output:
	// Saturday -> weekend
	// Wednesday -> midweek
	// Monday -> start of the week
}

func ExampleSwitchTagless() {
	SwitchTagless(os.Stdout)
	// Output:
	// 204 success
	// 301 redirect
	// 404 client error
	// 503 server error
	// 99 unknown
}

func ExampleSwitchWithInit() {
	SwitchWithInit(os.Stdout)
	// Output:
	// 6 letters: long
}

func ExampleSwitchEvalOrder() {
	SwitchEvalOrder(os.Stdout)
	// Output:
	// eval tag eval a eval b -> first case
}

func ExampleFallthrough() {
	Fallthrough(os.Stdout)
	// Output:
	// n=1: one big?
	// n=2: two
	// n=3: other
}

func ExampleTypeSwitch() {
	TypeSwitch(os.Stdout)
	// Output:
	// int, doubled: 84
	// string of length 2
	// a float: 3.5 (x is float64 here)
	// nil interface
	// something else: []int
}

func ExampleLabeledBreakContinue() {
	LabeledBreakContinue(os.Stdout)
	// Output:
	// row 0 is clean
	// row 1 has a negative, skipped
	// row 2 is clean
	// first value > 5 at [1][2]
}

func ExampleGoto() {
	Goto(os.Stdout)
	// Output:
	// leading letters: "ab", rest: "1c"
}

func ExampleGuardClauses() {
	GuardClauses(os.Stdout)
	// Output:
	// ""     -> missing
	// "x"    -> not a number
	// "150"  -> out of range
	// "30"   -> age 30
}

func ExampleBreakInSwitch() {
	BreakInSwitch(os.Stdout)
	// Output:
	// a b STOP c | unlabeled break: the loop kept going
	// a b | break loop: stopped at STOP
}

func ExampleReadLoop() {
	ReadLoop(os.Stdout)
	// Output:
	// setting 1: port=8080
	// setting 2: host=localhost
}
