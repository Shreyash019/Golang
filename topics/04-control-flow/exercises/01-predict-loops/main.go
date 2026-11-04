// Exercise E1 — predict the output.
//
// Goal: predict, line by line, what this program prints. It trains the two
// loop forms when the body changes the loop variable (notes §5.3), break and
// continue inside a switch (§7.1), fallthrough (§5.7), labeled continue
// (§5.8), what range evaluates and copies (§5.4), ranging over a string
// (§5.5) and per-iteration loop variables (§5.3).
//
// Write your prediction in prediction.go BEFORE running anything, then:
//
//	make ex E=01-predict-loops
//	go run ./exercises/01-predict-loops
package main

import (
	"fmt"
	"io"
	"os"
)

func main() { run(os.Stdout) }

func run(w io.Writer) {
	// 1
	for i := range 4 {
		i *= 10
		fmt.Fprint(w, i, ";")
	}
	fmt.Fprintln(w)

	// 2
	for i := 0; i < 10; i++ {
		fmt.Fprint(w, i, ";")
		i *= 2
	}
	fmt.Fprintln(w)

	// 3
	for i := range 6 {
		switch {
		case i == 1:
			continue
		case i == 3:
			break
		case i == 5:
			fallthrough
		case i > 10:
			fmt.Fprint(w, "big;")
		}
		fmt.Fprint(w, i, ";")
	}
	fmt.Fprintln(w)

	// 4
outer:
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			if j > i {
				continue outer
			}
			fmt.Fprintf(w, "%d%d;", i, j)
		}
	}
	fmt.Fprintln(w)

	// 5
	arr := [3]int{1, 2, 3}
	sl := []int{1, 2, 3}
	sum1, sum2 := 0, 0
	for i, v := range arr {
		arr[2] = 10
		sum1 += v + i
	}
	for i, v := range sl {
		sl[2] = 10
		sum2 += v + i
	}
	fmt.Fprintln(w, sum1, sum2)

	// 6
	nums := []int{1, 2, 3}
	for _, n := range nums {
		n *= 2
		nums = append(nums, n)
	}
	fmt.Fprintln(w, nums)

	// 7
	for i, r := range "añb" {
		fmt.Fprintf(w, "%d%c;", i, r)
	}
	fmt.Fprintln(w, len("añb"))

	// 8
	var fs []func() int
	for i := 0; i < 3; i++ {
		fs = append(fs, func() int { return i * i })
	}
	for _, f := range fs {
		fmt.Fprint(w, f(), ";")
	}
	fmt.Fprintln(w)
}
