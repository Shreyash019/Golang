package examples

import (
	"fmt"
	"io"
)

// LabeledBreakContinue shows labels on an outer loop: continue outer moves
// to the outer loop's next iteration, break outer leaves both loops.
func LabeledBreakContinue(w io.Writer) {
	grid := [][]int{
		{1, 2, 3},
		{4, -1, 6},
		{7, 8, 9},
	}

	// Skip the rest of a row as soon as it holds a negative number.
rows:
	for r, row := range grid {
		for _, v := range row {
			if v < 0 {
				fmt.Fprintf(w, "row %d has a negative, skipped\n", r)
				continue rows
			}
		}
		fmt.Fprintf(w, "row %d is clean\n", r)
	}

	// Find the first value above 5 and stop searching.
	foundR, foundC := -1, -1
search:
	for r, row := range grid {
		for c, v := range row {
			if v > 5 {
				foundR, foundC = r, c
				break search
			}
		}
	}
	fmt.Fprintf(w, "first value > 5 at [%d][%d]\n", foundR, foundC)
}

// Goto shows goto jumping backward to a label: this is all a loop is once
// compiled, which is why for can replace it. Real Go code almost never uses
// goto; it appears in generated code and a few hand-tuned state machines.
func Goto(w io.Writer) {
	input := "ab1c"
	i := 0
letters:
	if i < len(input) && input[i] >= 'a' && input[i] <= 'z' {
		i++
		goto letters
	}
	fmt.Fprintf(w, "leading letters: %q, rest: %q\n", input[:i], input[i:])
}
