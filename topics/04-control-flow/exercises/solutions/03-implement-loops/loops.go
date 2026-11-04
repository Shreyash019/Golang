// Package loops is a set of small functions built from Go's control flow.
package loops

import (
	"strconv"
	"strings"
)

// StatusClass returns the class of an HTTP status code: "informational"
// (100–199), "success" (200–299), "redirect" (300–399), "client error"
// (400–499), "server error" (500–599), and "invalid" for anything else.
//
// A tagless switch reads like a table of ranges. Rejecting invalid codes
// first lets every later case test only its lower bound, because cases are
// tried top to bottom and the first true one wins.
func StatusClass(code int) string {
	switch {
	case code < 100 || code > 599:
		return "invalid"
	case code >= 500:
		return "server error"
	case code >= 400:
		return "client error"
	case code >= 300:
		return "redirect"
	case code >= 200:
		return "success"
	default:
		return "informational"
	}
}

// CollatzSteps returns how many steps the Collatz sequence takes to reach 1
// from n: halve an even number, replace an odd number with 3n+1. It returns
// ok == false, and 0 steps, when n < 1.
//
// The number of iterations isn't known in advance, so this is the
// condition-only form of for: Go's while loop. The guard clause comes first
// because the loop would never end for n <= 0.
func CollatzSteps(n int) (steps int, ok bool) {
	if n < 1 {
		return 0, false
	}
	for n != 1 {
		if n%2 == 0 {
			n /= 2
		} else {
			n = 3*n + 1
		}
		steps++
	}
	return steps, true
}

// Primes returns the prime numbers less than or equal to limit, in
// increasing order, by trial division. It returns nil when there are none.
//
// continue candidates abandons the inner loop and moves the outer loop on as
// soon as a divisor is found, with no "isPrime" flag to set and test.
// d*d <= c stops at the square root: a composite number always has a divisor
// no larger than its square root.
func Primes(limit int) []int {
	var primes []int
candidates:
	for c := 2; c <= limit; c++ {
		for d := 2; d*d <= c; d++ {
			if c%d == 0 {
				continue candidates
			}
		}
		primes = append(primes, c)
	}
	return primes
}

// Find returns the position of the first occurrence of target in grid,
// scanning rows top to bottom and each row left to right. col counts runes,
// not bytes: in "héllo", 'l' is at column 2. When target isn't in grid,
// Find returns 0, 0, false.
//
// Ranging over a string yields byte offsets, so the column is counted
// separately. Returning from inside the nested loops is the simplest way
// out of both at once; a labeled break is the alternative when code after
// the loops has to run.
func Find(grid []string, target rune) (row, col int, ok bool) {
	for r, line := range grid {
		c := 0
		for _, ch := range line {
			if ch == target {
				return r, c, true
			}
			c++
		}
	}
	return 0, 0, false
}

// Compress run-length encodes s: each run of the same rune becomes the rune
// followed by the run's length, so "aaabcc" is "a3b1c2" and "ééé!" is "é3!1".
// Compress("") is "".
//
// The loop carries two pieces of state, the current rune and its count, and
// writes a run out when the rune changes. The last run has no "next rune"
// to trigger it, so it's written after the loop. strings.Builder avoids
// building a new string on every +=.
func Compress(s string) string {
	if s == "" {
		return ""
	}
	var b strings.Builder
	var cur rune
	count := 0
	for _, r := range s {
		if count > 0 && r != cur {
			b.WriteRune(cur)
			b.WriteString(strconv.Itoa(count))
			count = 0
		}
		cur = r
		count++
	}
	b.WriteRune(cur)
	b.WriteString(strconv.Itoa(count))
	return b.String()
}
