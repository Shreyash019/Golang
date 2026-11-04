package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"testing"
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
	inspect := func(fn func(n ast.Node) bool) {
		for _, f := range files {
			ast.Inspect(f, fn)
		}
	}

	t.Run("errors are returned, not panicked and recovered", func(t *testing.T) {
		inspect(func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			if id, ok := call.Fun.(*ast.Ident); ok && (id.Name == "panic" || id.Name == "recover") {
				t.Errorf("%s: %s(): bad input is an expected outcome, so return an error (fmt.Errorf) and let the caller decide; panic is for programmer bugs", fset.Position(call.Pos()), id.Name)
			}
			return true
		})
	})

	t.Run("no interface for a family of functions", func(t *testing.T) {
		inspect(func(n ast.Node) bool {
			ts, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}
			if _, isIface := ts.Type.(*ast.InterfaceType); isIface {
				t.Errorf("%s: interface %s: a one-method \"strategy\" with no state is just a function; use a func type and closures (interfaces are declared by their consumers, topic 17)", fset.Position(ts.Pos()), ts.Name.Name)
			}
			return true
		})
	})

	t.Run("no Impl, Factory or Strategy types", func(t *testing.T) {
		inspect(func(n ast.Node) bool {
			ts, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}
			for _, suffix := range []string{"Impl", "Factory", "Strategy"} {
				if strings.HasSuffix(ts.Name.Name, suffix) {
					t.Errorf("%s: type %s: Java naming for an abstraction Go doesn't need; a function that returns (func, error) replaces the factory", fset.Position(ts.Pos()), ts.Name.Name)
				}
			}
			return true
		})
	})

	t.Run("results are returned, not written through out-parameters", func(t *testing.T) {
		basic := map[string]bool{"int": true, "int64": true, "int32": true, "string": true, "bool": true, "float64": true}
		inspect(func(n ast.Node) bool {
			fd, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}
			for _, field := range fd.Type.Params.List {
				star, ok := field.Type.(*ast.StarExpr)
				if !ok {
					continue
				}
				if id, ok := star.X.(*ast.Ident); ok && basic[id.Name] {
					t.Errorf("%s: %s takes a *%s to fill in: return the value and an error instead, e.g. func parsePrice(s string) (int64, error)", fset.Position(field.Pos()), fd.Name.Name, id.Name)
				}
			}
			return true
		})
	})
}
