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
		{Name: "type-sizes", Section: "5.1", Summary: "the size of every basic type and the range of int", Run: TypeSizes},
		{Name: "integer-overflow", Section: "5.2", Summary: "integer arithmetic wraps around silently (two's complement)", Run: IntegerOverflow},
		{Name: "division-and-remainder", Section: "5.3", Summary: "/ truncates toward zero, % takes the dividend's sign, >> floors", Run: DivisionAndRemainder},
		{Name: "float-precision", Section: "5.4", Summary: "what float32 and float64 can't represent exactly", Run: FloatPrecision},
		{Name: "float-special-values", Section: "5.4", Summary: "+Inf, -Inf, NaN and -0: float division by zero doesn't panic", Run: FloatSpecialValues},
		{Name: "conversions", Section: "5.5", Summary: "T(x) truncates, wraps and loses precision, silently", Run: Conversions},
		{Name: "alias-vs-defined", Section: "5.6", Summary: "an alias is the same type; a defined type needs a conversion", Run: AliasVsDefined},
		{Name: "operator-precedence", Section: "5.7", Summary: "five precedence levels: << and & bind tighter than in C/Java", Run: OperatorPrecedence},
		{Name: "bitwise-ops", Section: "6.1", Summary: "&, |, ^, &^, shifts and unary ^ on signed and unsigned values", Run: BitwiseOps},
		{Name: "short-circuit", Section: "6.2", Summary: "&& and || skip the right operand when the result is known", Run: ShortCircuit},
		{Name: "assignment-ops", Section: "6.3", Summary: "op= assignments, ++ and -- statements", Run: AssignmentOps},
		{Name: "comparison-rules", Section: "6.4", Summary: "string ordering, comparing mixed types, NaN, min and max", Run: ComparisonRules},
		{Name: "string-basics", Section: "6.5", Summary: "strings are immutable bytes: len, indexing vs range", Run: StringBasics},
		{Name: "string-conversions", Section: "6.6", Summary: "string(rune) vs strconv for numbers and text", Run: StringConversions},
		{Name: "complex-numbers", Section: "6.7", Summary: "complex128, real, imag, math/cmplx", Run: ComplexNumbers},
		{Name: "overflow-detection", Section: "6.8", Summary: "math/bits: carries, 128-bit products, bit counts", Run: OverflowDetection},
		{Name: "float-comparison", Section: "6.9", Summary: "== vs a relative tolerance for floats", Run: FloatComparison},
		{Name: "number-literals", Section: "6.10", Summary: "0b, 0o, 0x, _ separators, hex floats, rune escapes, legacy octal", Run: NumberLiterals},
		{Name: "unsigned-underflow", Section: "7.2", Summary: "an unsigned len-1 wraps to a huge number", Run: UnsignedUnderflow},
		{Name: "int-division-by-zero", Section: "7.3", Summary: "integer division by a zero variable panics", Crashes: true, Run: IntDivisionByZero},
		{Name: "negative-shift", Section: "7.4", Summary: "a negative shift count panics", Crashes: true, Run: NegativeShift},
	}
}
