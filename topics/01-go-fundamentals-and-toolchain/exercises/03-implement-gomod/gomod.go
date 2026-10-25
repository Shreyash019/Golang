// Package gomod parses the subset of go.mod files that every Go module has:
// the module, go, toolchain and require directives.
//
// Exercise E3 — implement Parse so that gomod_test.go passes. It trains
// reading go.mod (notes §5.5) and writing a small, table-tested parser with
// good error messages.
//
//	make ex E=03-implement-gomod
//
// The rules your parser must follow:
//   - Blank lines and // comments are ignored. A trailing "// indirect"
//     comment on a requirement sets Require.Indirect.
//   - "module", "go" and "toolchain" each take exactly one argument.
//   - "require" appears as a single line (require path version) or as a
//     block: "require (" ... ")" with one "path version" per line.
//   - Other directives (replace, exclude, retract, tool, godebug, ...) are
//     skipped, including their block form.
//   - A malformed line returns an error that names its 1-based line number,
//     e.g. "line 6: ...". A file without a module directive returns
//     ErrNoModule.
package gomod

import "errors"

// Module is the parsed content of a go.mod file.
type Module struct {
	Path      string    // from the module directive
	GoVersion string    // from the go directive, e.g. "1.22.0"
	Toolchain string    // from the toolchain directive, e.g. "go1.27.1"; "" if absent
	Require   []Require // in file order, from single-line and block requires
}

// Require is one module requirement.
type Require struct {
	Path     string
	Version  string
	Indirect bool // the line carries a "// indirect" comment
}

// ErrNoModule is returned when the data has no module directive.
var ErrNoModule = errors.New("missing module directive")

// Parse parses the go.mod content in data.
func Parse(data []byte) (Module, error) {
	return Module{}, nil
}
