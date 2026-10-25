// Package examples holds the runnable worked examples from this topic's notes.
//
// Each example writes its output to an io.Writer, so the same code is run by
// cmd/examples and checked by the Example tests in examples_test.go: the
// output shown in the notes is the output the tests verify.
package examples

import "io"

// Example is one worked example from the notes.
type Example struct {
	// Name is the kebab-case name used on the command line, e.g. "slice-aliasing".
	Name string
	// Section is the notes section the example belongs to, e.g. "5.2".
	Section string
	// Summary is a one-line description shown when listing examples.
	Summary string
	// Crashes marks examples that panic, deadlock or race on purpose.
	// They run only when named explicitly, never as part of "all".
	Crashes bool
	// Run executes the example, writing its output to w.
	Run func(w io.Writer)
}

// All returns the examples in the order they appear in the notes.
func All() []Example {
	return []Example{
		{Name: "import-path-vs-package-name", Section: "6.3", Summary: "math/rand/v2 is package rand; an alias renames path", Run: ImportPathVsPackageName},
		{Name: "blank-import-registers", Section: "6.4", Summary: "a blank import runs image/png's init to register a decoder", Run: BlankImportRegisters},
		{Name: "embed-directive", Section: "6.6", Summary: "//go:embed bakes a file into the binary", Run: EmbedDirective},
		{Name: "build-info", Section: "6.7", Summary: "what the toolchain records inside every binary", Run: BuildInfo},
	}
}
