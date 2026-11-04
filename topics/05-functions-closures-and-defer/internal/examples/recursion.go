package examples

import (
	"fmt"
	"io"
	"strings"
)

// Recursion shows a function calling itself: a base case that stops it, and
// a recursive case that moves toward the base case.
func Recursion(w io.Writer) {
	fmt.Fprintln(w, "factorial(10) =", factorial(10))
	fmt.Fprintln(w, "gcd(84, 36) =", gcd(84, 36))

	tree := &node{"/", []*node{
		{"etc", []*node{{"hosts", nil}}},
		{"usr", []*node{{"bin", []*node{{"go", nil}}}, {"lib", nil}}},
	}}
	printTree(w, tree, 0)
	fmt.Fprintln(w, "files:", count(tree))
}

func factorial(n int) int {
	if n <= 1 { // base case
		return 1
	}
	return n * factorial(n-1) // recursive case: n shrinks toward 1
}

// gcd is tail-recursive, but Go doesn't eliminate tail calls: each call
// still gets its own stack frame.
func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

type node struct {
	name     string
	children []*node
}

func printTree(w io.Writer, n *node, depth int) {
	fmt.Fprintf(w, "%s%s\n", strings.Repeat("  ", depth), n.name)
	for _, c := range n.children {
		printTree(w, c, depth+1)
	}
}

func count(n *node) int {
	total := 1
	for _, c := range n.children {
		total += count(c)
	}
	return total
}

// StackOverflow recurses with no base case. Goroutine stacks grow, but only
// up to a limit (1 GB on 64-bit platforms), and hitting it is a fatal error
// that recover can't catch.
func StackOverflow(w io.Writer) {
	fmt.Fprintln(w, "recursing forever...")
	fmt.Fprintln(w, depth(0))
}

func depth(n int) int {
	return depth(n+1) + 1 //nolint:staticcheck // SA5007: the infinite recursion is the point
}
