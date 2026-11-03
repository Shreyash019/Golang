package examples

import (
	"fmt"
	"io"
	"math"
)

// DeclarationForms shows the ways to declare a variable and the type each
// one ends up with.
func DeclarationForms(w io.Writer) {
	var a int           // explicit type, zero value
	var b = 42          // type inferred from the initializer
	var c, d = "x", 3.5 // several at once, each inferred separately
	e := true           // short variable declaration: inside functions only
	var (
		host = "localhost"
		port = 8080
	)
	fmt.Fprintf(w, "a=%v b=%v c=%v d=%v e=%v\n", a, b, c, d, e)
	fmt.Fprintf(w, "%T %T %T %T %T\n", a, b, c, d, e)
	fmt.Fprintf(w, "%s:%d\n", host, port)
}

// DefaultTypes shows the type an untyped constant gets when nothing else
// decides it: in := declarations and when stored in an interface (any).
func DefaultTypes(w io.Writer) {
	i := 7
	f := 7.0
	r := 'x'
	c := 2i
	s := "go"
	b := 1 < 2
	fmt.Fprintf(w, "%T %T %T %T %T %T\n", i, f, r, c, s, b)

	const mixed = 1 + 2.0 // untyped int + untyped float = untyped float
	fmt.Fprintf(w, "mixed=%v %T\n", mixed, mixed)
	fmt.Fprintf(w, "any(1<<10) is %T\n", any(1<<10))
}

// MultipleAssignment shows tuple assignment: the swap idiom, the two-phase
// evaluation order from the spec, and the blank identifier.
func MultipleAssignment(w io.Writer) {
	x, y := 1, 2
	x, y = y, x // the right side is evaluated before anything is assigned
	fmt.Fprintln(w, "swapped:", x, y)

	i := 0
	nums := []int{10, 20, 30}
	i, nums[i] = 2, 99 // the index i is evaluated (0) before i becomes 2
	fmt.Fprintln(w, "i:", i, "nums:", nums)

	_, frac := math.Modf(3.75) // _ discards a value you must receive
	fmt.Fprintln(w, "fraction:", frac)
}
