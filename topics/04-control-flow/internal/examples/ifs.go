package examples

import (
	"fmt"
	"io"
	"strconv"
)

// IfWithInit shows that a variable declared in an if statement's init
// statement is visible in the condition, the if block and every else / else
// if block, and nowhere after the statement.
func IfWithInit(w io.Writer) {
	for _, in := range []string{"42", "-7", "abc"} {
		if n, err := strconv.Atoi(in); err != nil {
			fmt.Fprintf(w, "%q: not a number (%v)\n", in, err)
		} else if n < 0 {
			fmt.Fprintf(w, "%q: negative %d\n", in, n)
		} else {
			fmt.Fprintf(w, "%q: ok, n*2 = %d\n", in, n*2)
		}
		// n and err don't exist here: the if statement is their scope.
	}
}

// GuardClauses shows the idiomatic shape of a function with several
// failure cases: handle each one and return early, so the happy path stays
// at the left margin with no else.
func GuardClauses(w io.Writer) {
	for _, in := range []string{"", "x", "150", "30"} {
		fmt.Fprintf(w, "%-6q -> %s\n", in, describeAge(in))
	}
}

func describeAge(s string) string {
	if s == "" {
		return "missing"
	}
	age, err := strconv.Atoi(s)
	if err != nil {
		return "not a number"
	}
	if age < 0 || age > 130 {
		return "out of range"
	}
	return "age " + strconv.Itoa(age)
}
