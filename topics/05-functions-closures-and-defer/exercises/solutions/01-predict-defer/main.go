// Exercise E1 — predict the output.
//
// Goal: predict, line by line, what this program prints. It trains defer
// order and argument evaluation (notes §5.7), deferred closures and named
// results (§5.7), closures capturing variables (§5.5), variadic aliasing
// (§5.3), what a function can and can't change in its caller's slice
// (§5.1) and composing function values (§5.4).
//
// Write your prediction in prediction.go BEFORE running anything, then:
//
//	make ex E=01-predict-defer
//	go run ./exercises/01-predict-defer
package main

import (
	"fmt"
	"io"
	"os"
)

func main() { run(os.Stdout) }

func run(w io.Writer) {
	// 1
	func() {
		for i := range 3 {
			defer fmt.Fprint(w, i, ";")
		}
		fmt.Fprint(w, "loop;")
	}()
	fmt.Fprintln(w)

	// 2
	func() {
		s := "a"
		defer fmt.Fprintln(w, "arg:", s)
		defer func() { fmt.Fprint(w, "closure:", s, " ") }()
		s = "b"
	}()

	// 3
	fmt.Fprintln(w, named(), unnamed())

	// 4
	a, b := makeCounter(), makeCounter()
	a()
	a()
	b()
	fmt.Fprintln(w, a(), b())

	// 5
	x := 1
	get := func() int { return x }
	x = 5
	fmt.Fprintln(w, get())

	// 6
	nums := []int{1, 2, 3}
	zeroFirst(nums...)
	zeroFirst(nums[1], nums[2])
	fmt.Fprintln(w, nums)

	// 7
	s := []int{1}
	grow(s)
	fmt.Fprintln(w, s)

	// 8
	inc := func(n int) int { return n + 1 }
	double := func(n int) int { return n * 2 }
	fmt.Fprintln(w, apply(3, inc, double), apply(3, double, inc), apply(3))
}

func named() (n int) {
	defer func() { n *= 10 }()
	n = 1
	return n + 1
}

func unnamed() int {
	n := 1
	defer func() { n *= 10 }()
	return n + 1
}

func makeCounter() func() int {
	c := 0
	return func() int {
		c++
		return c
	}
}

func zeroFirst(xs ...int) {
	if len(xs) > 0 {
		xs[0] = 0
	}
}

func grow(s []int) {
	s = append(s, 2)
	s[0] = 9
}

func apply(x int, fs ...func(int) int) int {
	for _, f := range fs {
		x = f(x)
	}
	return x
}
