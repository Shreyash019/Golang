package examples

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// ForForms shows the three shapes of Go's only loop keyword: the
// three-clause loop, the condition-only ("while") loop and the infinite loop
// left with break.
func ForForms(w io.Writer) {
	// 1. three-clause: init; condition; post
	for i := 0; i < 3; i++ {
		fmt.Fprint(w, i, " ")
	}
	fmt.Fprintln(w, "| three-clause")

	// 2. condition only: Go's "while"
	n := 27
	steps := 0
	for n != 1 {
		if n%2 == 0 {
			n /= 2
		} else {
			n = 3*n + 1
		}
		steps++
	}
	fmt.Fprintln(w, steps, "| while: Collatz steps from 27")

	// 3. infinite: exit with break (or return)
	attempt := 0
	for {
		attempt++
		if attempt == 3 {
			break
		}
	}
	fmt.Fprintln(w, attempt, "| infinite loop left on attempt 3")

	// Parallel assignment in the post statement: two indexes moving inward.
	s := []byte("stressed")
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	fmt.Fprintln(w, string(s), "| two indexes")
}

// ReadLoop shows the most common while loop in Go programs: reading input
// with bufio.Scanner, whose Scan method is the loop condition. continue
// skips lines, break stops early, and the error is checked after the loop.
func ReadLoop(w io.Writer) {
	input := "# config\nport=8080\n\nhost=localhost\nEND\nignored=true\n"
	sc := bufio.NewScanner(strings.NewReader(input))
	n := 0
	for sc.Scan() {
		line := sc.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if line == "END" {
			break
		}
		n++
		fmt.Fprintf(w, "setting %d: %s\n", n, line)
	}
	if err := sc.Err(); err != nil { // Scan returns false on EOF *and* on error
		fmt.Fprintln(w, "read error:", err)
	}
}

// LoopVarPerIteration shows that since Go 1.22 every iteration of a for
// loop gets a fresh variable: closures and pointers taken in different
// iterations see different variables.
func LoopVarPerIteration(w io.Writer) {
	var prints []func()
	var ptrs []*int
	for i := 0; i < 3; i++ {
		prints = append(prints, func() { fmt.Fprint(w, i, " ") })
		ptrs = append(ptrs, &i)
	}
	for _, p := range prints {
		p()
	}
	fmt.Fprintln(w, "| closures")
	fmt.Fprintln(w, *ptrs[0], *ptrs[1], *ptrs[2], "| pointers")
	fmt.Fprintln(w, ptrs[0] != ptrs[1], "| different variables")
}

// LoopVarCopyBack shows the difference between the two loop forms when the
// body changes the loop variable. A three-clause loop copies the value into
// the next iteration's variable before the post statement, so the change
// carries over. A range loop assigns the next value from the range, so the
// change is lost.
func LoopVarCopyBack(w io.Writer) {
	for i := 0; i < 10; i++ {
		fmt.Fprint(w, i, " ")
		i += 2 // carried into the next iteration, then i++ runs
	}
	fmt.Fprintln(w, "| three-clause: the body's change carries over")

	for i := range 10 {
		fmt.Fprint(w, i, " ")
		i += 2 // overwritten: the next iteration gets the next value
		_ = i
	}
	fmt.Fprintln(w, "| range: the body's change is lost")
}
