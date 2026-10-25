package examples

import (
	"fmt"
	"io"
	"math/rand/v2" // import path ends in "v2", but the package clause says "package rand"
	pathpkg "path" // an alias: the file now refers to package path as pathpkg
)

// ImportPathVsPackageName shows that an import path and a package name are
// different things: the path locates the package, and the name in its
// package clause is what the importing file uses (unless it sets an alias).
func ImportPathVsPackageName(w io.Writer) {
	// A seeded PCG source makes the "random" numbers reproducible.
	r := rand.New(rand.NewPCG(1, 2))
	fmt.Fprintln(w, "rand.IntN:", r.IntN(100), r.IntN(100), r.IntN(100))

	fmt.Fprintln(w, "pathpkg.Join:", pathpkg.Join("cmd", "api", "main.go"))
}
