package examples

import (
	"fmt"
	"io"
	"strings"
)

// FunctionValues shows that functions are values: they have types, can be
// stored in variables, struct fields and slices, passed and returned. The
// zero value of a function type is nil, and func values compare only to nil.
func FunctionValues(w io.Writer) {
	var op func(int, int) int // nil
	fmt.Fprintln(w, "zero value is nil:", op == nil)

	op = add
	fmt.Fprintln(w, "op = add:", op(2, 3))

	op = func(a, b int) int { return a * b } // an anonymous function
	fmt.Fprintln(w, "op = multiply:", op(2, 3))

	// A named function type documents intent and can be reused.
	type transform func(string) string
	steps := []transform{strings.TrimSpace, strings.ToUpper, func(s string) string { return s + "!" }}
	s := "  hello "
	for _, step := range steps {
		s = step(s)
	}
	fmt.Fprintf(w, "pipeline: %q\n", s)
	fmt.Fprintf(w, "type of add: %T\n", add)
}

// HigherOrder shows functions that take or return functions: filter takes a
// predicate, and compose builds a new function from two others.
func HigherOrder(w io.Writer) {
	nums := []int{1, 2, 3, 4, 5, 6}
	even := func(n int) bool { return n%2 == 0 }
	fmt.Fprintln(w, "filter even:", filter(nums, even))
	fmt.Fprintln(w, "filter > 3: ", filter(nums, func(n int) bool { return n > 3 }))

	inc := func(n int) int { return n + 1 }
	sq := func(n int) int { return n * n }
	fmt.Fprintln(w, "compose(sq, inc)(3) =", compose(sq, inc)(3), "(inc first, then sq)")
	fmt.Fprintln(w, "compose(inc, sq)(3) =", compose(inc, sq)(3))
}

func filter(xs []int, keep func(int) bool) []int {
	var out []int
	for _, x := range xs {
		if keep(x) {
			out = append(out, x)
		}
	}
	return out
}

// compose returns f∘g: a function that applies g, then f.
func compose(f, g func(int) int) func(int) int {
	return func(x int) int { return f(g(x)) }
}

// NilFuncCall calls a nil function value, which panics with a nil pointer
// dereference.
func NilFuncCall(w io.Writer) {
	var onDone func(string)
	fmt.Fprintln(w, "calling a nil func value...")
	onDone("finished")
}
