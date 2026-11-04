package examples

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// DeferOrder shows that deferred calls run when the surrounding function
// returns, in last-in-first-out order.
func DeferOrder(w io.Writer) {
	fmt.Fprintln(w, "start")
	defer fmt.Fprintln(w, "deferred 1")
	defer fmt.Fprintln(w, "deferred 2")
	defer fmt.Fprintln(w, "deferred 3")
	fmt.Fprintln(w, "end of function body")
}

// DeferArgsEvaluated shows that a deferred call's function value and
// arguments are evaluated when the defer statement runs, not when the call
// runs. A deferred closure reads variables when it runs.
func DeferArgsEvaluated(w io.Writer) {
	x := 1
	defer fmt.Fprintln(w, "deferred call, argument evaluated at defer time: x =", x)
	defer func() { fmt.Fprintln(w, "deferred closure, reads x when it runs:     x =", x) }()
	x = 2
	fmt.Fprintln(w, "x set to", x)
}

// DeferModifiesResult shows a deferred closure reading and changing a named
// result after the return statement has set it: the classic way to annotate
// every error a function returns.
func DeferModifiesResult(w io.Writer) {
	fmt.Fprintln(w, "triple(5) =", triple(5))
	_, err := loadUser("")
	fmt.Fprintln(w, "loadUser(\"\"):", err)
	_, err = loadUser("gopher")
	fmt.Fprintln(w, "loadUser(\"gopher\"):", err)
}

func triple(n int) (result int) {
	defer func() { result *= 3 }() // runs after "return n" sets result = n
	return n
}

func loadUser(name string) (user string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("load user %q: %w", name, err)
		}
	}()
	if name == "" {
		return "", errors.New("empty name")
	}
	return strings.ToUpper(name), nil
}

// DeferRunsOnPanic shows that deferred calls still run while a panic unwinds
// the stack, which is why defer is the place for cleanup. Recovering from
// the panic is topic 18.
func DeferRunsOnPanic(w io.Writer) {
	defer fmt.Fprintln(w, "deferred: cleanup ran during the panic")
	fmt.Fprintln(w, "about to index out of range")
	var xs []int
	_ = xs[3]
}

// DeferInLoop shows the fix for defer in a loop: move the loop body into a
// function, so each deferred call runs at the end of each iteration's call.
func DeferInLoop(w io.Writer) {
	open := 0
	acquire := func(name string) func() {
		open++
		fmt.Fprintf(w, "open %s (open now: %d)\n", name, open)
		return func() {
			open--
			fmt.Fprintf(w, "close %s\n", name)
		}
	}

	fmt.Fprintln(w, "-- defer directly in the loop:")
	func() {
		for _, name := range []string{"a", "b", "c"} {
			release := acquire(name)
			defer release() // runs when this func returns, after the loop
		}
		fmt.Fprintln(w, "loop finished with", open, "open")
	}()

	fmt.Fprintln(w, "-- body moved into a function:")
	for _, name := range []string{"a", "b", "c"} {
		func() {
			release := acquire(name)
			defer release() // runs at the end of this iteration
		}()
	}
	fmt.Fprintln(w, "loop finished with", open, "open")
}

// DeferTrace shows the Effective Go trace idiom: the expression in a defer
// statement is evaluated immediately, so trace runs on entry and the
// function it returns runs on exit.
func DeferTrace(w io.Writer) {
	depth := 0
	trace := func(name string) func() {
		fmt.Fprintf(w, "%senter %s\n", strings.Repeat("  ", depth), name)
		depth++
		return func() {
			depth--
			fmt.Fprintf(w, "%sleave %s\n", strings.Repeat("  ", depth), name)
		}
	}

	inner := func() {
		defer trace("inner")()
		fmt.Fprintf(w, "%s(inner's work)\n", strings.Repeat("  ", depth))
	}
	outer := func() {
		defer trace("outer")()
		inner()
		inner()
	}
	outer()
}
