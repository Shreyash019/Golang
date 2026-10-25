// Exercise E1 — predict the output.
//
// Goal: predict, line by line, what this program prints. It trains Go's
// program initialization order (notes §5.3): package-level variables are
// initialized in dependency order, across files in the order the go command
// passes them to the compiler (sorted by file name), then every init function
// runs in that same file order, and only then does main start.
//
// Write your prediction in prediction.go BEFORE running anything, then:
//
//	make ex E=01-predict-init-order
//	go run ./exercises/01-predict-init-order
package main

import (
	"fmt"
	"io"
	"os"
)

// trace records initialization steps. A package-level mutable slice is not
// something to copy into real code; it's here only to observe init order.
var trace []string

//nolint:unused // initialized only to observe the order
var (
	total = add("total", price+tax)
	price = add("price", 100)
	tax   = add("tax", price/5)
	label = add("label", 7)
)

func add(name string, v int) int {
	trace = append(trace, fmt.Sprintf("%s=%d", name, v))
	return v
}

func init() { trace = append(trace, "init A (main.go)") }

func init() { trace = append(trace, fmt.Sprintf("init B (main.go) sees total=%d", total)) }

func main() { run(os.Stdout) }

func run(w io.Writer) {
	trace = append(trace, "run")
	for _, step := range trace {
		fmt.Fprintln(w, step)
	}
}
