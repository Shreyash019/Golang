package examples

import (
	"fmt"
	"io"
)

// ClosureCounter shows a closure keeping state between calls: each call to
// counter creates a new count variable, and the returned function keeps it
// alive.
func ClosureCounter(w io.Writer) {
	a := counter()
	b := counter()
	fmt.Fprintln(w, a(), a(), a(), "| a has its own count")
	fmt.Fprintln(w, b(), b(), "| b has another")
	fmt.Fprintln(w, a(), "| a continues")
}

func counter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

// ClosureSharesVariable shows that a closure captures the variable itself,
// not a copy of its value: changes made after the closure was created are
// visible to it, and two closures over the same variable share it.
func ClosureSharesVariable(w io.Writer) {
	greeting := "hello"
	say := func() string { return greeting + ", world" }
	greeting = "goodbye" // changed after the closure was created
	fmt.Fprintln(w, say())

	total := 0
	addTo := func(n int) { total += n }
	read := func() int { return total }
	addTo(5)
	addTo(10)
	fmt.Fprintln(w, "read():", read(), "| total:", total)

	// To capture the current value, copy it into a new variable first.
	v := 1
	snapshot := v
	get := func() int { return snapshot }
	v = 2
	fmt.Fprintln(w, "snapshot:", get(), "| v:", v)
}

// ClosureState shows closures sharing private state, the closure-based
// alternative to a small struct: next and reset both see the same id.
func ClosureState(w io.Writer) {
	next, reset := idGenerator("order-")
	fmt.Fprintln(w, next(), next(), next())
	reset()
	fmt.Fprintln(w, next())
}

func idGenerator(prefix string) (next func() string, reset func()) {
	id := 0
	next = func() string {
		id++
		return fmt.Sprintf("%s%d", prefix, id)
	}
	reset = func() { id = 0 }
	return next, reset
}

// RecursiveClosure shows that a function literal can't refer to itself by
// the name being declared with :=, because that name isn't in scope until
// the statement ends. Declare the variable first, then assign.
func RecursiveClosure(w io.Writer) {
	var fib func(n int) int
	fib = func(n int) int {
		if n < 2 {
			return n
		}
		return fib(n-1) + fib(n-2)
	}
	fmt.Fprintln(w, "fib(10) =", fib(10))
}
