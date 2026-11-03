// Exercise E1 — predict the output.
//
// Goal: predict, line by line, what this program prints. It trains integer
// wraparound (notes §5.2), division, remainder and shifts on negative
// numbers (§5.3), which decimals binary floats represent exactly (§5.4),
// conversions (§5.5), operator precedence and ^ / &^ (§5.7, §6.1), and
// strings as bytes (§6.5).
//
// Write your prediction in prediction.go BEFORE running anything, then:
//
//	make ex E=01-predict-int-math
//	go run ./exercises/01-predict-int-math
package main

import (
	"fmt"
	"io"
	"os"
)

func main() { run(os.Stdout) }

func run(w io.Writer) {
	var small int8 = 120
	small += 10
	fmt.Fprintln(w, small)

	a, b := -9, 4
	fmt.Fprintln(w, a/b, a%b, a>>2)

	fmt.Fprintln(w, 1<<3+2, 2^4, 7&^5)

	f := 3.99
	fmt.Fprintln(w, int(f), int(-f))

	x := 1000
	fmt.Fprintln(w, uint8(x), int8(x))

	fmt.Fprintln(w, 9/2*2.0, 9/2.0*2)

	s := "naïve"
	fmt.Fprintln(w, len(s), s[2], string(s[0]))

	var u uint = 3
	fmt.Fprintln(w, u-5 > 0)

	p, q := 0.5, 0.25
	fmt.Fprintln(w, p+q == 0.75)
}
