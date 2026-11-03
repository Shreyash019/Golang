package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"unicode"
)

// TestIdioms fails while the habits are still in the code. Each subtest
// names one habit and the idiomatic alternative (notes §8).
func TestIdioms(t *testing.T) {
	fset := token.NewFileSet()
	paths, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	var files []*ast.File
	for _, p := range paths {
		if strings.HasSuffix(p, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, p, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		files = append(files, f)
	}

	t.Run("no floating point for money or counts", func(t *testing.T) {
		for _, f := range files {
			ast.Inspect(f, func(n ast.Node) bool {
				if id, ok := n.(*ast.Ident); ok && (id.Name == "float64" || id.Name == "float32") {
					t.Errorf("%s: %s: keep money as an integer number of cents (int64) and counts as int; binary floats can't represent 0.10 exactly", fset.Position(id.Pos()), id.Name)
				}
				return true
			})
		}
	})

	t.Run("no math package for integer arithmetic", func(t *testing.T) {
		for _, f := range files {
			for _, imp := range f.Imports {
				if path, _ := strconv.Unquote(imp.Path.Value); path == "math" {
					t.Errorf("%s: import \"math\": ceiling division is (a + b - 1) / b, and rounding cents is integer arithmetic too", fset.Position(imp.Pos()))
				}
			}
		}
	})

	t.Run("no comparisons with true or false", func(t *testing.T) {
		for _, f := range files {
			ast.Inspect(f, func(n ast.Node) bool {
				be, ok := n.(*ast.BinaryExpr)
				if !ok || (be.Op != token.EQL && be.Op != token.NEQ) {
					return true
				}
				for _, side := range []ast.Expr{be.X, be.Y} {
					if id, ok := side.(*ast.Ident); ok && (id.Name == "true" || id.Name == "false") {
						t.Errorf("%s: comparing with %s: a bool is already a condition, write `if ok` or `if !ok`", fset.Position(be.Pos()), id.Name)
					}
				}
				return true
			})
		}
	})

	t.Run("no Get prefix", func(t *testing.T) {
		for _, f := range files {
			for _, d := range f.Decls {
				fd, ok := d.(*ast.FuncDecl)
				if !ok {
					continue
				}
				if rest, ok := strings.CutPrefix(fd.Name.Name, "Get"); ok && rest != "" && unicode.IsUpper([]rune(rest)[0]) {
					t.Errorf("%s: %s: Go doesn't prefix functions with Get (subtotal, not GetSubtotal)", fset.Position(fd.Pos()), fd.Name.Name)
				}
			}
		}
	})
}
