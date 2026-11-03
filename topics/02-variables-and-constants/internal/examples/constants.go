package examples

import (
	"fmt"
	"io"
	"time"
)

// UntypedConstants shows that constant arithmetic is exact and has no type
// until a value is used somewhere that needs one.
func UntypedConstants(w io.Writer) {
	const huge = 1 << 100    // far bigger than any integer type, and legal
	const small = huge >> 98 // 4: computed exactly at compile time
	fmt.Fprintln(w, "small:", small)

	var f float64 = huge // converted when used: representable as a float64
	fmt.Fprintln(w, "huge as float64:", f)

	var d time.Duration = 3 // the same untyped 3 becomes a Duration (3ns)
	var b byte = 'A'        // a rune constant becomes a byte
	fmt.Fprintln(w, "duration:", d, "byte:", b)
}

// ConstantFloatExactness compares 0.1 + 0.2 computed on exact constants
// with the same sum computed on float64 variables.
func ConstantFloatExactness(w io.Writer) {
	const a, b = 0.1, 0.2
	fmt.Fprintln(w, "constants:", a+b == 0.3)

	x, y := 0.1, 0.2
	fmt.Fprintln(w, "variables:", x+y == 0.3)
	fmt.Fprintln(w, "x+y =", x+y)
}

// TypedVsUntyped shows a typed constant (time.Duration) mixing freely with
// an untyped constant, but not with an int variable.
func TypedVsUntyped(w io.Writer) {
	const timeout = 2 * time.Second // typed: its type is time.Duration
	const retries = 3               // untyped: adapts to the other operand
	fmt.Fprintf(w, "timeout: %v (%T)\n", timeout, timeout)
	fmt.Fprintln(w, "retries*timeout:", retries*timeout)

	n := 3 // an int variable has a fixed type
	// n*timeout would not compile: mismatched types int and time.Duration
	fmt.Fprintln(w, "time.Duration(n)*timeout:", time.Duration(n)*timeout)
}

// ConstantIntDivision shows that dividing two untyped integer constants is
// integer division, even when the result is stored in a float64.
func ConstantIntDivision(w io.Writer) {
	const half = 1 / 2        //nolint:staticcheck // SA4025: integer division on purpose; untyped int / untyped int: 0
	const halfF = 1 / 2.0     // one float operand: 0.5
	var ratio float64 = 1 / 2 //nolint:staticcheck // SA4025 on purpose: evaluated as integers first, then converted
	fmt.Fprintln(w, half, halfF, ratio)
}
