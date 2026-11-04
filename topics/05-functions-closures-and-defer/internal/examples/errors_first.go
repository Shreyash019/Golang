package examples

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

// ErrNegative is returned by parseQuantity for quantities below zero. A
// package-level error value like this is called a sentinel error.
var ErrNegative = errors.New("quantity is negative")

// ErrorFirstLook shows the basic error convention: return error as the last
// result, nil means success, check it right after the call, add context
// with fmt.Errorf and %w, and test for a specific error with errors.Is.
// Errors in depth are topic 18.
func ErrorFirstLook(w io.Writer) {
	for _, in := range []string{"3", "x", "-2"} {
		q, err := parseQuantity(in)
		if err != nil {
			fmt.Fprintf(w, "%-4q error: %v (negative? %t)\n", in, err, errors.Is(err, ErrNegative))
			continue
		}
		fmt.Fprintf(w, "%-4q ok: %d\n", in, q)
	}
}

func parseQuantity(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("parse quantity %q: %w", s, err)
	}
	if n < 0 {
		return 0, fmt.Errorf("parse quantity %q: %w", s, ErrNegative)
	}
	return n, nil
}
