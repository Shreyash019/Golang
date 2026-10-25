// Package gomod parses the subset of go.mod files that every Go module has:
// the module, go, toolchain and require directives.
//
// Exercise E3 — implement. It trains reading go.mod (notes §5.5) and writing
// a small, table-tested parser with good error messages.
package gomod

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strings"
)

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

// Parse parses the go.mod content in data. Directives it doesn't model
// (replace, exclude, retract, tool, godebug, ignore, ...) are skipped,
// including their block form. Errors name the 1-based line they occur on.
func Parse(data []byte) (Module, error) {
	var (
		m     Module
		block string // directive of the block we're inside, or ""
	)
	sc := bufio.NewScanner(bytes.NewReader(data))
	for n := 1; sc.Scan(); n++ {
		line, comment := splitComment(sc.Text())
		if line == "" {
			continue
		}
		if block != "" {
			if line == ")" {
				block = ""
				continue
			}
			if block == "require" {
				r, err := parseRequire(strings.Fields(line), comment)
				if err != nil {
					return Module{}, fmt.Errorf("line %d: %w", n, err)
				}
				m.Require = append(m.Require, r)
			}
			continue
		}

		verb, rest, _ := strings.Cut(line, " ")
		rest = strings.TrimSpace(rest)
		if rest == "(" {
			block = verb
			continue
		}
		args := strings.Fields(rest)
		switch verb {
		case "module", "go", "toolchain":
			if len(args) != 1 {
				return Module{}, fmt.Errorf("line %d: %s directive needs exactly one argument", n, verb)
			}
			switch verb {
			case "module":
				m.Path = args[0]
			case "go":
				m.GoVersion = args[0]
			case "toolchain":
				m.Toolchain = args[0]
			}
		case "require":
			r, err := parseRequire(args, comment)
			if err != nil {
				return Module{}, fmt.Errorf("line %d: %w", n, err)
			}
			m.Require = append(m.Require, r)
		}
	}
	if err := sc.Err(); err != nil {
		return Module{}, err
	}
	if m.Path == "" {
		return Module{}, ErrNoModule
	}
	return m, nil
}

// splitComment returns the line without its // comment, and the comment text.
func splitComment(s string) (line, comment string) {
	line, comment, _ = strings.Cut(s, "//")
	return strings.TrimSpace(line), strings.TrimSpace(comment)
}

func parseRequire(args []string, comment string) (Require, error) {
	if len(args) != 2 {
		return Require{}, fmt.Errorf("require needs a module path and a version, got %q", strings.Join(args, " "))
	}
	return Require{Path: args[0], Version: args[1], Indirect: comment == "indirect"}, nil
}
