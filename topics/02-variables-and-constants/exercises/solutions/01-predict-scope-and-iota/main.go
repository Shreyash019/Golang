// Exercise E1 — predict the output.
//
// Goal: predict, line by line, what this program prints. It trains iota and
// implicit repetition (notes §5.5), scope and shadowing, including where a
// variable's scope starts (§5.6), := redeclaration (§5.7), and untyped
// constants with their default types (§5.2, §5.3).
//
// Write your prediction in prediction.go BEFORE running anything, then:
//
//	make ex E=01-predict-scope-and-iota
//	go run ./exercises/01-predict-scope-and-iota
package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

const (
	low = iota * 5
	mid
	_
	high
)

// level is a package-level variable. It is never assigned after
// initialization; the program only shadows it.
var level = 1

func main() { run(os.Stdout) }

func run(w io.Writer) {
	fmt.Fprintln(w, low, mid, high)

	level := level + 1
	if level := level * 10; level > 10 {
		fmt.Fprintln(w, "if:", level)
	}
	fmt.Fprintln(w, "function:", level)
	fmt.Fprintln(w, "package:", packageLevel())

	n, err := strconv.Atoi("12") //nolint:ineffassign,staticcheck // the next := reuses err on purpose
	m, err := strconv.Atoi("x")
	fmt.Fprintln(w, n, m, err != nil)

	fmt.Fprintln(w, 7/2, 7/2.0, 7.0/2)
	fmt.Fprintf(w, "%T %T\n", 'a', 1<<3)

	const a, b = 0.1, 0.2
	p, q := a, b
	fmt.Fprintln(w, a+b == 0.3, p+q == 0.3)
}

func packageLevel() int { return level }
