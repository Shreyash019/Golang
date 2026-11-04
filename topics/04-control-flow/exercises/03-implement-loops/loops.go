// Exercise E3 — implement.
//
// Goal: implement the five functions below until every test passes. Each one
// trains a different control-flow tool: a tagless switch (notes §5.6), the
// condition-only for (§5.2), a labeled continue (§5.8), ranging over a string
// by rune (§5.5) and a loop that carries state between iterations.
//
//	make ex E=03-implement-loops

// Package loops is a set of small functions built from Go's control flow.
package loops

// StatusClass returns the class of an HTTP status code: "informational"
// (100–199), "success" (200–299), "redirect" (300–399), "client error"
// (400–499), "server error" (500–599), and "invalid" for anything else.
func StatusClass(code int) string {
	return ""
}

// CollatzSteps returns how many steps the Collatz sequence takes to reach 1
// from n: halve an even number, replace an odd number with 3n+1. It returns
// ok == false, and 0 steps, when n < 1.
func CollatzSteps(n int) (steps int, ok bool) {
	return 0, false
}

// Primes returns the prime numbers less than or equal to limit, in
// increasing order, by trial division. It returns nil when there are none.
func Primes(limit int) []int {
	return nil
}

// Find returns the position of the first occurrence of target in grid,
// scanning rows top to bottom and each row left to right. col counts runes,
// not bytes: in "héllo", 'l' is at column 2. When target isn't in grid,
// Find returns 0, 0, false.
func Find(grid []string, target rune) (row, col int, ok bool) {
	return 0, 0, false
}

// Compress run-length encodes s: each run of the same rune becomes the rune
// followed by the run's length, so "aaabcc" is "a3b1c2" and "ééé!" is "é3!1".
// Compress("") is "".
func Compress(s string) string {
	return ""
}
