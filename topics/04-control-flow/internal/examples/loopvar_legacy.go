// The build constraint below sets this file's language version to Go 1.21
// (in a module whose go line is 1.21 or later, a //go:build go1.N line sets
// the file's language version; see go help buildconstraint), so its loops have the old, pre-1.22 semantics: ONE variable for the whole
// loop. Modules whose go.mod says "go 1.21" or lower get this behaviour in
// every file.

//go:build go1.21

package examples

import (
	"fmt"
	"io"
)

// LoopVarLegacy shows the pre-Go 1.22 loop variable: every closure and every
// pointer refers to the same variable, so they all see its final value.
func LoopVarLegacy(w io.Writer) {
	var prints []func()
	for i := 0; i < 3; i++ {
		prints = append(prints, func() { fmt.Fprint(w, i, " ") })
	}
	for _, p := range prints {
		p()
	}
	fmt.Fprintln(w, "| closures (Go 1.21 semantics)")

	var ptrs []*int
	for _, v := range []int{10, 20, 30} {
		ptrs = append(ptrs, &v)
	}
	fmt.Fprintln(w, *ptrs[0], *ptrs[1], *ptrs[2], "| pointers (Go 1.21 semantics)")
}
