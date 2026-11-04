package examples

import (
	"fmt"
	"io"
	"time"
)

// SwitchBasics shows an expression switch: cases are compared with ==, a case
// can list several values, the matching case's body runs and the switch
// ends (no fallthrough by default), and default can go anywhere.
func SwitchBasics(w io.Writer) {
	for _, d := range []time.Weekday{time.Saturday, time.Wednesday, time.Monday} {
		switch d {
		default:
			fmt.Fprintln(w, d, "-> midweek")
		case time.Saturday, time.Sunday:
			fmt.Fprintln(w, d, "-> weekend")
		case time.Monday:
			fmt.Fprintln(w, d, "-> start of the week")
		}
	}
}

// SwitchTagless shows a switch with no expression ("switch true"): the first
// case whose condition is true wins, so order matters.
func SwitchTagless(w io.Writer) {
	for _, code := range []int{204, 301, 404, 503, 99} {
		var class string
		switch {
		case code >= 500 && code <= 599:
			class = "server error"
		case code >= 400:
			class = "client error"
		case code >= 300:
			class = "redirect"
		case code >= 200:
			class = "success"
		default:
			class = "unknown"
		}
		fmt.Fprintln(w, code, class)
	}
}

// SwitchWithInit shows a switch with an init statement, scoped to the switch
// like an if's.
func SwitchWithInit(w io.Writer) {
	switch n := len("gopher"); {
	case n > 5:
		fmt.Fprintln(w, n, "letters: long")
	default:
		fmt.Fprintln(w, n, "letters: short")
	}
}

// SwitchEvalOrder shows that case expressions are evaluated lazily, top to
// bottom and left to right, and evaluation stops at the first match.
func SwitchEvalOrder(w io.Writer) {
	check := func(name string, v int) int {
		fmt.Fprintf(w, "eval %s ", name)
		return v
	}
	switch check("tag", 2) {
	case check("a", 1), check("b", 2), check("c", 3):
		fmt.Fprintln(w, "-> first case")
	case check("d", 2):
		fmt.Fprintln(w, "-> second case")
	}
}

// Fallthrough shows that fallthrough transfers control to the next case's
// body WITHOUT testing its condition.
func Fallthrough(w io.Writer) {
	for _, n := range []int{1, 2, 3} {
		fmt.Fprintf(w, "n=%d:", n)
		switch {
		case n == 1:
			fmt.Fprint(w, " one")
			fallthrough
		case n > 100: // not tested when we fall into it
			fmt.Fprint(w, " big?")
		case n == 2:
			fmt.Fprint(w, " two")
		default:
			fmt.Fprint(w, " other")
		}
		fmt.Fprintln(w)
	}
}

// TypeSwitch gives a first look at a type switch: it branches on the dynamic
// type stored in an interface value. Interfaces and type switches are topic 12.
func TypeSwitch(w io.Writer) {
	for _, v := range []any{42, "go", 3.5, nil, []int{1}} {
		switch x := v.(type) {
		case int:
			fmt.Fprintln(w, "int, doubled:", x*2)
		case string:
			fmt.Fprintln(w, "string of length", len(x))
		case float64, float32:
			fmt.Fprintf(w, "a float: %v (x is %T here)\n", x, x)
		case nil:
			fmt.Fprintln(w, "nil interface")
		default:
			fmt.Fprintf(w, "something else: %T\n", x)
		}
	}
}

// BreakInSwitch shows the classic trap: break inside a switch (or select)
// leaves the switch, not the enclosing loop.
func BreakInSwitch(w io.Writer) {
	tokens := []string{"a", "b", "STOP", "c"}
	for _, t := range tokens {
		switch t {
		case "STOP":
			break //nolint:staticcheck // SA4011 on purpose: this is the trap being shown
		}
		fmt.Fprint(w, t, " ")
	}
	fmt.Fprintln(w, "| unlabeled break: the loop kept going")

loop:
	for _, t := range tokens {
		switch t {
		case "STOP":
			break loop
		}
		fmt.Fprint(w, t, " ")
	}
	fmt.Fprintln(w, "| break loop: stopped at STOP")
}
